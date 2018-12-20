package server

import (
	"encoding/json"
	"net/http"
	"net/url"
)

const errorIsRequired = "is_required"
const errorMinLength = "min_length"
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
