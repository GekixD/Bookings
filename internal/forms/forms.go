package forms

import (
	"net/http"
	"net/url"
)

// Form creates a custom form struct and it embeds a url.Values object
type Form struct {
	url.Values
	Errors errors
}

// New initializes an empty form struct
func New(data url.Values) *Form {
	return &Form{
		data,
		errors(map[string][]string{}),
	}
}

// Has checks if a form field is in Post and non-empty
func (f *Form) Has(field string, req *http.Request) bool {
	contents := req.Form.Get(field)
	if contents == "" {
		return false
	}
	return true
}
