package server

import "net/url"

type loginUserPayload struct {
	Email    string
	Password string
	payloadValidator
}

func (p *loginUserPayload) isValid() bool {
	p.errs = url.Values{}
	if p.Email == "" {
		p.errs.Add("email", errorIsRequired)
	}

	if p.Password == "" {
		p.errs.Add("password", errorIsRequired)
	}

	return len(p.errs) == 0
}

type createUserPayload struct {
	Email    string
	Name     string
	Password string
	payloadValidator
}

func (p *createUserPayload) isValid() bool {
	p.errs = url.Values{}

	if p.Email == "" {
		p.errs.Add("email", errorIsRequired)
	}

	if p.Password == "" {
		p.errs.Add("password", errorIsRequired)
	} else if len(p.Password) < 6 {
		p.errs.Add("password", errorMinLength)
	}

	if p.Name == "" {
		p.errs.Add("name", errorIsRequired)
	}

	return len(p.errs) == 0
}
