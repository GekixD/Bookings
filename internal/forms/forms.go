package forms

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/asaskevich/govalidator"
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

// Required checks the form's required fields to make sure they're populated
func (f *Form) Required(fields ...string) {
	for _, field := range fields {
		value := f.Get(field)
		if strings.TrimSpace(value) == "" {
			f.Errors.Add(field, "This field can not be empty!")
		}
	}
}

// Has checks if a value from a field is in Post and non-empty
func (f *Form) Has(field string, req *http.Request) bool {
	contents := req.Form.Get(field)
	return contents != ""
	// if contents == "" {
	// 	return false
	// }
	// return true
}

// Valid returns true if there are no errors, otherwise false
func (f *Form) Valid() bool {
	return len(f.Errors) == 0
}

// MinLength checks a string if it has the minimum required length in order to validate
func (f *Form) MinLength(field string, length int, req *http.Request) bool {
	value := req.Form.Get(field)
	if len(value) < length {
		f.Errors.Add(field, fmt.Sprintf("This field must be at least %d characters long", length))
		return false
	}
	return true
}

// IsEmail uses govalidate to check if the field contains a valid email address
func (f *Form) IsEmail(field string) {
	if !govalidator.IsEmail(f.Get(field)) {
		f.Errors.Add(field, "This is not a valid email address!")
	}
}
