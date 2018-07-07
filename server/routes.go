package server

import (
	"github.com/gorilla/mux"
)

// Routes is the link to all routes.
type Routes struct {
	Root      *mux.Router
	ApiRoot   *mux.Router
	Users     *mux.Router
	AuthToken *mux.Router
	Expenses  *mux.Router
}

// NewRoutes returns new Routes object by passing in a mux.Router
func NewRoutes(root *mux.Router) *Routes {
	routes := &Routes{Root: root}
	routes.ApiRoot = root.PathPrefix("/api").Subrouter()
	routes.Users = routes.ApiRoot.PathPrefix("/users").Subrouter()
	routes.Expenses = routes.ApiRoot.PathPrefix("/expenses").Subrouter()
	return routes
}
