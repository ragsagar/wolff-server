package server

import (
	"net/url"
	"time"
)

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
