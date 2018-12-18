package server

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"time"

	"github.com/go-pg/pg"
	"github.com/gorilla/mux"
	"github.com/ragsagar/wolff/model"
)

func (srv *Server) InitExpenseAPIs() {
	srv.Routes.Expenses.Handle("/", srv.ApiWithTokenValidation(createExpense)).Methods("POST")
	srv.Routes.Expenses.Handle("/", srv.ApiWithTokenValidation(getExpenses)).Methods("GET")
	srv.Routes.Expenses.Handle("/accounts/", srv.ApiWithTokenValidation(getExpenseAccounts)).Methods("GET")
	srv.Routes.Expenses.Handle("/accounts/", srv.ApiWithTokenValidation(createExpenseAccount)).Methods("POST")
	srv.Routes.Expenses.Handle("/accounts/{id}/", srv.ApiWithTokenValidation(deleteExpenseAccount)).Methods("DELETE")
}

type expenseRequestBody struct {
	AccountID  string    `json:"account_id"`
	Date       time.Time `json:"date"`
	CategoryID string    `json:"category_id"`
	Amount     float64   `json:"amount"`
	Title      string    `json:"title"`
}

func (e *expenseRequestBody) validate() url.Values {
	errs := url.Values{}

	if e.AccountID == "" {
		errs.Add("title", "title field is required.")
	}

	nilTime := time.Time{}
	if e.Date == nilTime {
		errs.Add("date", "date field is required.")
	}

	if e.Amount == 0 {
		errs.Add("amount", "amount field is required.")
	}

	if e.Title == "" {
		errs.Add("title", "title field is required.")
	}

	return errs
}

func createExpense(c *Context, w http.ResponseWriter, r *http.Request) {
	expenseData := &expenseRequestBody{}
	defer r.Body.Close()
	if err := json.NewDecoder(r.Body).Decode(expenseData); err != nil {
		log.Println(err.Error())
		panic(err)
	}
	if validationErrors := expenseData.validate(); len(validationErrors) > 0 {
		errorMsg := map[string]interface{}{"errors": validationErrors}
		writeJSONResponse(errorMsg, http.StatusBadRequest, w)
		return
	}
	log.Println(expenseData)
	expense := model.Expense{
		AccountID:  expenseData.AccountID,
		Date:       expenseData.Date,
		CategoryID: expenseData.CategoryID,
		Amount:     expenseData.Amount,
		UserID:     c.User.ID,
		Title:      expenseData.Title,
	}
	if err := c.Srv.Store.Expense().Store(&expense); err != nil {
		response := map[string]string{"errors": err.Error()}
		WriteJsonResponse(response, http.StatusInternalServerError, w)
	}

	jsonData, err := expense.ToJSON()
	if err != nil {
		log.Println("Error in json generation ", err.Error())
		response := map[string]string{"errors": "Error in expense json generation."}
		WriteJsonResponse(response, http.StatusInternalServerError, w)
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonData)

}

func getExpenses(c *Context, w http.ResponseWriter, r *http.Request) {
	expenses, err := c.Srv.Store.Expense().GetExpenses(c.User.ID)
	if err != nil {
		log.Println("Erorr in getting expenses: ", err.Error())
		response := map[string]string{"error_message": err.Error()}
		WriteJsonResponse(response, http.StatusInternalServerError, w)
		return
	}

	jsonData, err := json.Marshal(expenses)
	if err != nil {
		response := map[string]string{"error_message": err.Error()}
		WriteJsonResponse(response, http.StatusInternalServerError, w)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonData)
}

func deleteExpense(c *Context, w http.ResponseWriter, r *http.Request) {
	// expense, err := c.Srv.Store.Expense().DeleteExpenseWithUserID(id, userId)
}

func getExpenseAccounts(c *Context, w http.ResponseWriter, r *http.Request) {
	if c.User == nil {
		error_message := map[string]string{"error_message": "No user found in context."}
		WriteJsonResponse(error_message, http.StatusUnauthorized, w)
		return
	}

	expenseAccounts, err := c.Srv.Store.Expense().GetExpenseAccounts(c.User.ID)
	if err != nil {
		response := map[string]string{"error_message": err.Error()}
		WriteJsonResponse(response, http.StatusBadRequest, w)
		return
	}

	jsonData, err := json.Marshal(expenseAccounts)
	if err != nil {
		response := map[string]string{"error_message": err.Error()}
		WriteJsonResponse(response, http.StatusBadRequest, w)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonData)
}

func createExpenseAccount(c *Context, w http.ResponseWriter, r *http.Request) {
	// Validate and parse the json request body
	if c.User == nil {
		error_message := map[string]string{"error_message": "No user found in context."}
		WriteJsonResponse(error_message, http.StatusUnauthorized, w)
		return
	}
	var requestData struct {
		Name string
	}
	b, _ := ioutil.ReadAll(r.Body)
	json.Unmarshal(b, &requestData)
	if requestData.Name == "" {
		response := map[string]string{"error_message": "Name can't be blank."}
		WriteJsonResponse(response, http.StatusBadRequest, w)
		return
	}

	// Create the expense account in database.
	expenseAccount := model.ExpenseAccount{Name: requestData.Name, UserID: c.User.ID}
	expenseAccount.PreSave()
	err := c.Srv.Store.Expense().StoreAccount(expenseAccount)
	if err != nil {
		log.Println(err.Error())
		response := map[string]string{"error_message": err.Error()}
		WriteJsonResponse(response, http.StatusInternalServerError, w)
		return
	}
	log.Println("Successfully created account with id", expenseAccount.ID)

	// Construct the json response
	jsonData, err := expenseAccount.ToJSON()
	if err != nil {
		log.Println("Error in expense account json creation: ", err.Error())
		error_message := map[string]string{"error_message": err.Error()}
		WriteJsonResponse(error_message, http.StatusInternalServerError, w)
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonData)
}

func deleteExpenseAccount(c *Context, w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	expenseAccount, err := c.Srv.Store.Expense().GetAccountByID(id)
	if err != nil {
		if err == pg.ErrNoRows {
			error_message := map[string]string{"error_message": "No account exist with this id."}
			WriteJsonResponse(error_message, http.StatusNotFound, w)
		} else {
			error_message := map[string]string{"error_message": err.Error()}
			WriteJsonResponse(error_message, http.StatusInternalServerError, w)
		}
		return
	}
	if expenseAccount.UserID != c.User.ID {
		error_message := map[string]string{"error_message": "Not authorized to delete this account."}
		WriteJsonResponse(error_message, http.StatusUnauthorized, w)
		return
	}

	err = c.Srv.Store.Expense().DeleteAccount(expenseAccount)
	if err != nil {
		error_message := map[string]string{"error_message": err.Error()}
		WriteJsonResponse(error_message, http.StatusInternalServerError, w)
		return
	}
	w.WriteHeader(http.StatusNoContent)

}

func writeJSONResponse(res map[string]interface{}, s int, w http.ResponseWriter) {
	w.Header().Set("Content-type", "applciation/json")
	w.WriteHeader(s)
	json.NewEncoder(w).Encode(res)
}
