package forms

import (
	"fmt"
	"net/url"
	"strings"  
	"unicode/utf8"
)

type Form struct {
	url.Values
	Errors errors
}

// define func to initialize a custom Form struct
func New(data url.Values) *Form {
	return &Form{
		data,
		errors(map[string][]string{}),
	}
}

// method which check that specific field in the form data are not blank
func (f *Form) Required(fields ...string) {
	for _, field := range fields {
		value := f.Get(field)
		if strings.TrimSpace(value) == "" {
			f.Errors.Add(field, "This field cannot be blank")
		}
	}
}

// method which check that specific field in the form data 
// contains no more than d charecters
func (f *Form) MaxLength(field string, d int) {
	value := f.Get(field)
	if value == "" {
		return
	}
	if utf8.RuneCountInString(value) > d {
		f.Errors.Add(field, fmt.Sprintf("This field is too long (maximum is %d characters)", d))
	}
}

// method which check that specific field int the form data
// matches one of a set of specific permitted values.
// if the check fails then add the appropirate message to the form erros
func (f *Form) PermittedValues(field string, opts ...string) {
	value := f.Get(field)
	if value == "" {
		return
	}
	for _, opt := range opts {
		if value == opt {
			return
		}
	}
	f.Errors.Add(field, "This field is invalid")
}

// method which return true if there are no errors
func (f *Form) Valid() bool {
	return len(f.Errors) == 0
}