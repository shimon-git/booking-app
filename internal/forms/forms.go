package forms

import (
	"net/http"
	"net/url"
)

// Form creates a custom from struct, embed a url.Values object
type Form struct {
	url.Values
	Errors errors
}

// New initialize a Form struct
func New(data url.Values) *Form {
	return &Form{
		data,
		errors(map[string][]string{}),
	}
}

// Has check if form field is in post and not empty
func (f *Form) Has(field string, r *http.Request) bool {
	x := r.Form.Get(field)
	if x == "" {
		f.Errors.Add(field, "This field cennot be blank")
		return false
	} else {
		return true
	}
}

// Valid  return true if there are no errors, otherwise false
func (f *Form) Valid() bool {
	return len(f.Errors) == 0
}