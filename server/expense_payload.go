package server

import (
	"encoding/json"
	"net/http"
	"net/url"
	"time"
)

const errorIsRequired = "is_required"
const errorMinimumLength = "min_length"
const errorMaxLength = "max_length"
const errorInvalidJSON = "invalid_json_input"
const errorDbWrite = "db_write_error"
const errorDbFetch = "db_fetch_error"
const errorDbDelete = "db_delete_error"
const errorJSONGeneration = "error_in_generating_json"
const errorNotAuthorized = "not_authorized"
const errorNotFound = "not_found"

type payloadValidator struct {
	errs url.Values
}

func (p payloadValidator) getErrorMessage() map[string]interface{} {
	return map[string]interface{}{"errors": p.errs}
}

func (p payloadValidator) writeErrorMessage(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusBadRequest)
	json.NewEncoder(w).Encode(p.getErrorMessage())
}

type createExpensePayload struct {
	AccountID  string    `json:"account_id"`
	Date       time.Time `json:"date"`
	CategoryID string    `json:"category_id"`
	Amount     float64   `json:"amount"`
	Title      string    `json:"title"`
	payloadValidator
}

func (e *createExpensePayload) isValid() bool {
	e.errs = url.Values{}

	if e.AccountID == "" {
		e.errs.Add("account_id", errorIsRequired)
	}

	nilTime := time.Time{}
	if e.Date == nilTime {
		e.errs.Add("date", errorIsRequired)
	}

	if e.Amount == 0 {
		e.errs.Add("amount", errorIsRequired)
	}

	if e.Title == "" {
		e.errs.Add("title", errorIsRequired)
	}

	return len(e.errs) == 0
}

type createExpenseAccountPayload struct {
	Name string
	payloadValidator
}

func (p *createExpenseAccountPayload) isValid() bool {
	p.errs = url.Values{}
	if p.Name == "" {
		p.errs.Add("name", errorIsRequired)
	}
	return len(p.errs) == 0
}
