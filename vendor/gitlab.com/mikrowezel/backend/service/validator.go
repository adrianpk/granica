package service

import (
	"regexp"
	"strings"
	"unicode/utf8"
)

// NOTE: Experiment, preferring to avoid reflection
// at the expense of greater verbosity.
// The cli generator willi have to do the heavy lifting.
// Generic code will be moved to the most appropriate package.

type (
	Changes map[string][]string

	ErrorSet map[string][]string

	Validator struct {
		Errors ErrorSet
	}
)

const (
	RequiredErrMsg   = "required"
	MinLengthErrMsg  = "too short"
	MaxLengthErrMsg  = "too long"
	NotAllowedErrMsg = "not in allowed list"
	NotEmailErrMsg   = "not an email address"
	NoMatchErrMsg    = "confirmation does not match"
)

var (
	matchFirstCap = regexp.MustCompile("(.)([A-Z][a-z]+)")
	matchAllCap   = regexp.MustCompile("([a-z0-9])([A-Z])")
	emailRegex    = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+\\/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")
)

func NewValidator() Validator {
	return Validator{
		ErrorSet(map[string][]string{}),
	}
}

// ValidateRequired value.
func (v *Validator) ValidateRequired(val string, errMsg ...string) (ok bool) {
	val = strings.Trim(val, " ")
	return utf8.RuneCountInString(val) > 0
}

// ValidateMinLength value.
func (v *Validator) ValidateMinLength(val string, min int, errMsg ...string) (ok bool) {
	return utf8.RuneCountInString(val) > min
}

// ValidateMaxLength value.
func (v *Validator) ValidateMaxLength(val string, max int) (ok bool) {
	return utf8.RuneCountInString(val) < max
}

// ValidateEmail value.
func (v *Validator) ValidateEmail(val string) (ok bool) {
	return len(val) < 254 && emailRegex.MatchString(val)
}

// ValidateConfirmation value.
func (v *Validator) ValidateConfirmation(val, confirmation string) (ok bool) {
	return val == confirmation
}

func (v *Validator) HasErrors() bool {
	return !v.IsValid()
}

func (v *Validator) IsValid() bool {
	return len(v.Errors) == 0
}

func (es ErrorSet) Add(field, msg string) {
	es[field] = append(es[field], msg)
}

func (es ErrorSet) FieldErrors(field string) []string {
	return es[field]
}

func (es ErrorSet) IsEmpty() bool {
	return len(es) < 1
}
