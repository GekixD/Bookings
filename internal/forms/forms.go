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
		f.Errors.Add(field, "This field can not be empty!")
		return false
	}
	return true
}

// Valid returns true if there are no errors, otherwise false
func (f *Form) Valid() bool {
	return len(f.Errors) == 0
}
