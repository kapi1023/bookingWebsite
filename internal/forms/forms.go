package forms

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/asaskevich/govalidator"
)

// Form create a custom form struct, have a url.Values object
type Form struct {
	url.Values
	Error errors
}

// Valid returns true if there are no errors
func (f *Form) Valid() bool {
	return len(f.Error) == 0
}

// New initialize a form struct
func New(data url.Values) *Form {
	return &Form{
		data,
		errors(map[string][]string{}),
	}
}

// checks for Required fields
func (f *Form) Required(fields ...string) {
	for _, field := range fields {
		value := f.Get(field)
		if strings.TrimSpace(value) == "" {
			f.Error.Add(field, "This field cannot be blank")
		}
	}
}

// Has checks if form field are empty
func (f *Form) Has(field string, r *http.Request) bool {

	requestedField := r.Form.Get(field)

	if requestedField == "" {
		f.Error.Add(field, "This field cannot be blank")
		return false
	}
	return true

}

// MinLength check for string min length
func (f *Form) MinLength(field string, lenghth int, r *http.Request) bool {
	requestedField := r.Form.Get(field)
	if len(requestedField) < lenghth {
		f.Error.Add(field, fmt.Sprintf("This field must be at least %d characters long", lenghth))
		return false
	}
	return true
}

func (f *Form) IsEmail(field string) {
	if !govalidator.IsEmail(f.Get(field)) {
		f.Error.Add(field, "Invalid email address")
	}
}
