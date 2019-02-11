package server

import (
	"encoding/json"
	"io"
	"log"
	"net/http"

	"github.com/go-pg/pg"
	"github.com/gorilla/mux"
	"github.com/ragsagar/wolff/model"
	"github.com/ragsagar/wolff/store"
)

func (srv *Server) InitExpenseAPIs() {
	srv.Routes.Expenses.Handle("/", srv.ApiWithTokenValidation(createExpense)).Methods("POST")
	srv.Routes.Expenses.Handle("/", srv.ApiWithTokenValidation(getExpenses)).Methods("GET")
	srv.Routes.Expenses.Handle("/accounts/", srv.ApiWithTokenValidation(getExpenseAccounts)).Methods("GET")
	srv.Routes.Expenses.Handle("/accounts/", srv.ApiWithTokenValidation(createExpenseAccount)).Methods("POST")
	srv.Routes.Expenses.Handle("/accounts/{id}/", srv.ApiWithTokenValidation(deleteExpenseAccount)).Methods("DELETE")
}

func errorResponse(errorType string) map[string]interface{} {
	return map[string]interface{}{"error": errorType}
}

func loadJSON(payload interface{}, body io.Reader) error {
	return json.NewDecoder(body).Decode(payload)
}

func createExpense(c *Context, w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	payload := &createExpensePayload{}

	if err := loadJSON(payload, r.Body); err != nil {
		writeJSONResponse(errorResponse(errorInvalidJSON), http.StatusInternalServerError, w)
		return
	}

	if !payload.isValid() {
		payload.writeErrorMessage(w)
		return
	}
	expense := model.Expense{
		AccountID:  payload.AccountID,
		Date:       payload.Date,
		CategoryID: payload.CategoryID,
		Amount:     payload.Amount,
		UserID:     c.User.ID,
		Title:      payload.Title,
	}
	if err := c.Srv.Store.Expense().Store(&expense); err != nil {
		// TODO: Log this properly
		writeJSONResponse(errorResponse(errorDbWrite), http.StatusInternalServerError, w)
	}

	jsonData, err := expense.ToJSON()
	if err != nil {
		// TODO: Log this properly
		writeJSONResponse(errorResponse(errorJSONGeneration), http.StatusInternalServerError, w)
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonData)

}

func getExpenses(c *Context, w http.ResponseWriter, r *http.Request) {
	filter := store.ExpenseFilter{}
	filter.ParseURLValues(r.URL.Query())
	expenses, err := c.Srv.Store.Expense().GetExpenses(c.User.ID, filter)
	if err != nil {
		log.Println("Erorr in getting expenses: ", err.Error())
		writeJSONResponse(errorResponse(errorDbFetch), http.StatusInternalServerError, w)
		return
	}

	jsonData, err := json.Marshal(expenses)
	if err != nil {
		writeJSONResponse(errorResponse(errorJSONGeneration), http.StatusInternalServerError, w)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if string(jsonData) == "null" {
		w.Write([]byte("{}"))
	} else {
		w.Write(jsonData)
	}
}

func deleteExpense(c *Context, w http.ResponseWriter, r *http.Request) {
	// expense, err := c.Srv.Store.Expense().DeleteExpenseWithUserID(id, userId)
}

func getExpenseAccounts(c *Context, w http.ResponseWriter, r *http.Request) {

	expenseAccounts, err := c.Srv.Store.Expense().GetExpenseAccounts(c.User.ID)
	if err != nil {
		writeJSONResponse(errorResponse(errorDbFetch), http.StatusBadRequest, w)
		return
	}

	jsonData, err := json.Marshal(expenseAccounts)
	if err != nil {
		writeJSONResponse(errorResponse(errorJSONGeneration), http.StatusBadRequest, w)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonData)
}

func createExpenseAccount(c *Context, w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	payload := &createExpenseAccountPayload{}
	if err := loadJSON(payload, r.Body); err != nil {
		writeJSONResponse(errorResponse(errorInvalidJSON), http.StatusInternalServerError, w)
		return
	}

	if !payload.isValid() {
		payload.writeErrorMessage(w)
		return
	}

	// Create the expense account in database.
	expenseAccount := model.ExpenseAccount{Name: payload.Name, UserID: c.User.ID}
	expenseAccount.PreSave()
	err := c.Srv.Store.Expense().StoreAccount(expenseAccount)
	if err != nil {
		writeJSONResponse(errorResponse(errorDbWrite), http.StatusInternalServerError, w)
		return
	}

	log.Println("Successfully created account with id", expenseAccount.ID)

	// Construct the json response
	jsonData, err := expenseAccount.ToJSON()
	if err != nil {
		log.Println("Error in expense account json creation: ", err.Error())
		writeJSONResponse(errorResponse(errorJSONGeneration), http.StatusInternalServerError, w)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonData)
}

func deleteExpenseAccount(c *Context, w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	expenseAccount, err := c.Srv.Store.Expense().GetAccountByID(id)
	if err != nil {
		if err == pg.ErrNoRows {
			writeJSONResponse(errorResponse(errorNotFound), http.StatusNotFound, w)
		} else {
			writeJSONResponse(errorResponse(errorDbFetch), http.StatusInternalServerError, w)
		}
		return
	}
	if expenseAccount.UserID != c.User.ID {
		writeJSONResponse(errorResponse(errorNotAuthorized), http.StatusUnauthorized, w)
		return
	}

	err = c.Srv.Store.Expense().DeleteAccount(expenseAccount)
	if err != nil {
		writeJSONResponse(errorResponse(errorDbDelete), http.StatusInternalServerError, w)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func writeJSONResponse(res map[string]interface{}, s int, w http.ResponseWriter) {
	w.Header().Set("Content-type", "applciation/json")
	w.WriteHeader(s)
	json.NewEncoder(w).Encode(res)
}
