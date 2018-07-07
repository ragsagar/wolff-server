package server

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/ragsagar/wolff/model"
	"github.com/ragsagar/wolff/store"
)

type Server struct {
	Routes *Routes
	Store  store.Store
}

func NewServer(store store.Store) *Server {
	router := mux.NewRouter()
	srv := &Server{
		Routes: NewRoutes(router),
		Store:  store,
	}
	srv.InitUsers()
	return srv
}

func (srv *Server) ApiWithTokenValidation(hf handlerFunc) *handler {
	return &handler{
		hf:                hf,
		doTokenValidation: true,
		srv:               srv,
	}
}

func (srv *Server) OpenAPI(hf handlerFunc) *handler {
	return &handler{
		hf:                hf,
		doTokenValidation: false,
		srv:               srv,
	}
}

func (srv *Server) Run(addr string) {
	log.Println("Starting server at", addr)
	s := &http.Server{
		Addr:           addr,
		Handler:        srv.Routes.Root,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20, // 1 MB
	}
	log.Fatal(s.ListenAndServe())
}

type handlerFunc func(*Context, http.ResponseWriter, *http.Request)

type handler struct {
	hf                handlerFunc
	doTokenValidation bool
	srv               *Server
}

func (h handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	c := &Context{Srv: h.srv}
	log.Println(r.RemoteAddr, r.Method, r.URL)
	if h.doTokenValidation {
		auth_header := r.Header.Get("Authorization")
		user, err := validateToken(h.srv.Store, auth_header)
		if err != nil {
			log.Println(err)
			errorData := map[string]string{
				"error_message": err.Error(),
			}
			responseJson, _ := json.Marshal(errorData)
			w.WriteHeader(http.StatusUnauthorized)
			w.Write(responseJson)
			//http.Error(w, string(responseJson), http.StatusUnauthorized)
			return
		}
		c.User = user
	}
	// Call the http handler func
	h.hf(c, w, r)
}

func validateToken(store store.Store, token string) (*model.User, error) {
	if token == "" {
		return nil, errors.New("Token missing")
	}
	// TODO: Check token in db
	authToken, err := store.AuthToken().Find(token)
	if err != nil {
		return nil, errors.New("Token not found")
	}

	return authToken.User, nil
}
