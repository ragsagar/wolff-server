package server

import "net/http"

func (srv *Server) InitExpenseAPIs() {
	// srv.Routes.Expenses.Handle("/", srv.(createExpense)).Methods("POST")
}

func createExpense(c *Context, w http.ResponseWriter, r *http.Request) {
	// var requestData struct { }
}
