package server

import (
	"log"
	"net/http"

	"github.com/ragsagar/wolff/model"
)

func (srv *Server) InitUsers() {
	srv.Routes.Users.Handle("/", srv.OpenAPI(createUser)).Methods("POST")
	srv.Routes.Users.Handle("/login/", srv.OpenAPI(loginUser)).Methods("POST")
	srv.Routes.Users.Handle("/profile/", srv.ApiWithTokenValidation(getUserProfile)).Methods("GET")
}

func loginUser(c *Context, w http.ResponseWriter, r *http.Request) {
	payload := &loginUserPayload{}
	if err := loadJSON(payload, r.Body); err != nil {
		writeJSONResponse(errorResponse(errorInvalidJSON), http.StatusInternalServerError, w)
		return
	}

	if !payload.isValid() {
		payload.writeErrorMessage(w)
		return
	}

	user, err := c.Srv.Store.User().GetUserByEmail(payload.Email)
	if err != nil {
		log.Println("Error in fetching user by email. ", err.Error())
		writeJSONResponse(errorResponse(errorDbFetch), http.StatusInternalServerError, w)
	}
	if !model.CheckPasswordHash(payload.Password, user.Password) {
		writeJSONResponse(errorResponse("Invalid email or password"), http.StatusBadRequest, w)
		return
	}

	authToken, err := c.Srv.Store.AuthToken().Create(user)
	if err != nil {
		log.Println("Error in generating token.")
		writeJSONResponse(errorResponse(errorDbWrite), http.StatusInternalServerError, w)
		return
	}
	response := map[string]interface{}{
		"user_id":    user.ID,
		"auth_token": authToken.Key}
	writeJSONResponse(response, http.StatusCreated, w)
}

func getUserProfile(c *Context, w http.ResponseWriter, r *http.Request) {
	jsonData, err := c.User.ToJSON()
	if err != nil {
		log.Println("Error in user json creation: ", err.Error())
		writeJSONResponse(errorResponse(errorJSONGeneration), http.StatusInternalServerError, w)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonData)
}

func createUser(c *Context, w http.ResponseWriter, r *http.Request) {
	payload := &createUserPayload{}
	if err := loadJSON(payload, r.Body); err != nil {
		writeJSONResponse(errorResponse(errorInvalidJSON), http.StatusInternalServerError, w)
		return
	}
	if !payload.isValid() {
		payload.writeErrorMessage(w)
		return
	}

	user := model.User{Email: payload.Email, Name: payload.Name, Active: true}
	user.SetPassword(payload.Password)
	user.PreSave()
	err := c.Srv.Store.User().StoreUser(user)
	if err != nil {
		writeJSONResponse(errorResponse(errorDbWrite), http.StatusInternalServerError, w)
		return
	}
	log.Println("Successfully created user with id", user.ID)
	authToken, err := c.Srv.Store.AuthToken().Create(&user)
	if err != nil {
		log.Println("Error in generating token.")
		writeJSONResponse(errorResponse(errorDbWrite), http.StatusInternalServerError, w)
		return
	}
	response := map[string]interface{}{
		"user_id":    user.ID,
		"auth_token": authToken.Key}
	writeJSONResponse(response, http.StatusCreated, w)
}
