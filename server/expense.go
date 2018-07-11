package server

import (
	"encoding/json"
	"log"
	"net/http"
	"net/url"
	"time"

	"github.com/ragsagar/wolff/model"
)

func (srv *Server) InitExpenseAPIs() {
	srv.Routes.Expenses.Handle("/", srv.ApiWithTokenValidation(createExpense)).Methods("POST")
}

type expenseRequestBody struct {
	AccountID  string    `json:"account_id"`
	Date       time.Time `json:"date"`
	CategoryID string    `json:"category_id"`
	Amount     float64   `json:"amount"`
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

func writeJSONResponse(res map[string]interface{}, s int, w http.ResponseWriter) {
	w.Header().Set("Content-type", "applciation/json")
	w.WriteHeader(s)
	json.NewEncoder(w).Encode(res)
}
