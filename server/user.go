package server

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/ragsagar/wolff/model"
)

func (srv *Server) InitUsers() {
	srv.Routes.Users.Handle("/", srv.OpenAPI(createUser)).Methods("POST")
	srv.Routes.Users.Handle("/login/", srv.OpenAPI(loginUser)).Methods("POST")
	srv.Routes.Users.Handle("/profile/", srv.ApiWithTokenValidation(getUserProfile)).Methods("GET")
}

func loginUser(c *Context, w http.ResponseWriter, r *http.Request) {
	var requestData struct {
		Email    string
		Password string
	}
	b, _ := ioutil.ReadAll(r.Body)
	json.Unmarshal(b, &requestData)
	if requestData.Email == "" || requestData.Password == "" {
		response := map[string]string{"error_message": "Email or password is missing."}
		WriteJsonResponse(response, http.StatusBadRequest, w)
		return
	}
	user, err := c.Srv.Store.User().GetUserByEmail(requestData.Email)
	if err != nil {
		log.Println("Error in fetching user by email. ", err.Error())
		response := map[string]string{"error_message": err.Error()}
		WriteJsonResponse(response, http.StatusInternalServerError, w)
	}
	if !model.CheckPasswordHash(requestData.Password, user.Password) {
		response := map[string]string{"error_message": "Invalid email or password."}
		WriteJsonResponse(response, http.StatusBadRequest, w)
		return
	}

	authToken, err := c.Srv.Store.AuthToken().Create(user)
	if err != nil {
		log.Println("Error in generating token.")
		response := map[string]string{"error_message": err.Error()}
		WriteJsonResponse(response, http.StatusInternalServerError, w)
		return
	}
	response := map[string]string{
		"user_id":    user.Id,
		"auth_token": authToken.Key}
	WriteJsonResponse(response, http.StatusCreated, w)
}

func getUserById(c *Context, w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	log.Println("Id: ", vars["id"])
	user, err := c.Srv.Store.User().GetUserByID(vars["id"])
	if err != nil {
		log.Println(err)
		w.Write([]byte("Error"))
		return
	}
	jsonData, err := user.ToJSON()
	if err != nil {
		log.Println(err)
	}
	w.Write(jsonData)
}

func getUserProfile(c *Context, w http.ResponseWriter, r *http.Request) {
	if c.User == nil {
		log.Println("No user found in the context.")
		error_message := map[string]string{"error_message": "No user found in context."}
		WriteJsonResponse(error_message, http.StatusUnauthorized, w)
		return
	}
	jsonData, err := c.User.ToJSON()
	if err != nil {
		log.Println("Error in user json creation: ", err.Error())
		error_message := map[string]string{"error_message": "Error in user json generation."}
		WriteJsonResponse(error_message, http.StatusInternalServerError, w)
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonData)
}

func createUser(c *Context, w http.ResponseWriter, r *http.Request) {
	var requestData struct {
		Email    string
		Name     string
		Password string
	}
	b, _ := ioutil.ReadAll(r.Body)
	json.Unmarshal(b, &requestData)
	if requestData.Email == "" || requestData.Password == "" {
		response := map[string]string{"error_message": "Email or password is missing."}
		WriteJsonResponse(response, http.StatusBadRequest, w)
		return
	}
	user := model.User{Email: requestData.Email, Name: requestData.Name, Active: true}
	user.SetPassword(requestData.Password)
	user.PreSave()
	err := c.Srv.Store.User().StoreUser(user)
	if err != nil {
		response := map[string]string{"error_message": err.Error()}
		WriteJsonResponse(response, http.StatusInternalServerError, w)
		return
	}
	log.Println("Successfully created user with id", user.Id)
	authToken, err := c.Srv.Store.AuthToken().Create(&user)
	if err != nil {
		log.Println("Error in generating token.")
		response := map[string]string{"error_message": err.Error()}
		WriteJsonResponse(response, http.StatusInternalServerError, w)
		return
	}
	response := map[string]string{
		"user_id":    user.Id,
		"auth_token": authToken.Key}
	WriteJsonResponse(response, http.StatusCreated, w)
}

func WriteJsonResponse(response map[string]string, statusCode int, w http.ResponseWriter) {
	jsonData, err := json.Marshal(response)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	w.Write(jsonData)
}
