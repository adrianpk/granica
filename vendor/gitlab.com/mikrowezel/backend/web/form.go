package web

import (
	"net/url"
	"regexp"
	"strings"
	"unicode"
	"unicode/utf8"
)

type (
	Changes map[string][]string

	ErrorSet map[string][]string

	Form struct {
		Changes
		Errors ErrorSet
	}
)

const (
	requiredErrMsg   = "Required"
	minLengthErrMsg  = "Too short"
	maxLengthErrMsg  = "Too long"
	notAllowedErrMsg = "Not in allowed list"
	notEmailErrMsg   = "Not an email address"
	notMatchErrMsg   = "Confirmation does not match"
)

var (
	matchFirstCap = regexp.MustCompile("(.)([A-Z][a-z]+)")
	matchAllCap   = regexp.MustCompile("([a-z0-9])([A-Z])")
	emailRegex    = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+\\/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")
)

func NewForm(f url.Values) *Form {
	return &Form{
		Changes(f),
		ErrorSet(map[string][]string{}),
	}
}

// NOTE: Pascal to dashed case conversion hinders and removes flexibility.
// It is desirable to operate stright with form field names.
// Remove this case conversion after implementing model validator in the service layer (*)
// and use this methods only for trivial web form validations.
// This is a temporary solution that dirties functionality of this package a bit
// just to keep a more neat code where it is used.
// (*) This way we can use same business rules for JSON REST and GRPC endpoints.

// ValidateRequired value.
// Field name is a Pascal case string.
// Associated form field is this name in dashed case format.
// i.e.: field 'FamilyName' > mapKey 'family-name'
func (f *Form) ValidateRequired(fields []string, errMsg ...string) {
	for _, field := range fields {
		mapKey := pascalToDashed(field)
		val := f.Get(mapKey)

		msg := requiredErrMsg
		if len(errMsg) > 0 {
			msg = errMsg[0]
		}

		if strings.TrimSpace(val) == "" {
			f.Errors.Add(field, msg)
		}
	}
}

// ValidateMinLength value.
// Field name is a Pascal case string.
// Associated form field is this name in dashed case format.
// i.e.: field 'FamilyName' > mapKey 'family-name'
func (f *Form) ValidateMinLength(field string, max int, errMsg ...string) {
	mapKey := pascalToDashed(field)
	val := f.Get(mapKey)

	if val == "" {
		return
	}

	msg := minLengthErrMsg
	if len(errMsg) > 0 {
		msg = errMsg[0]
	}

	if utf8.RuneCountInString(val) > max {
		f.Errors.Add(field, msg)
	}
}

// ValidateMaxLength value.
// Field name is a Pascal case string.
// Associated form field is this name in dashed case format.
// i.e.: field 'FamilyName' > mapKey 'family-name'
func (f *Form) ValidateMaxLength(field string, max int, errMsg ...string) {
	mapKey := pascalToDashed(field)
	val := f.Get(mapKey)

	if val == "" {
		return
	}

	msg := maxLengthErrMsg
	if len(errMsg) > 0 {
		msg = errMsg[0]
	}

	if utf8.RuneCountInString(val) > max {
		f.Errors.Add(field, msg)
	}
}

// ValidateAllowed value.
// Field name is a Pascal case string.
// Associated form field is this name in dashed case format.
// i.e.: field 'FamilyName' > mapKey 'family-name'
func (f *Form) ValidateAllowed(field string, allowed []string, errMsg ...string) {
	mapKey := pascalToDashed(field)
	val := f.Get(mapKey)

	if val == "" {
		return
	}

	msg := notAllowedErrMsg
	if len(errMsg) > 0 {
		msg = errMsg[0]
	}

	for _, a := range allowed {
		if val == a {
			return
		}
	}

	f.Errors.Add(field, msg)
}

// ValidateEmail value.
// Field name is a Pascal case string.
// Associated form field is this name in dashed case format.
// i.e.: field 'FamilyName' > mapKey 'family-name'
func (f *Form) ValidateEmail(field string, errMsg ...string) {
	mapKey := pascalToDashed(field)
	val := f.Get(mapKey)

	msg := notEmailErrMsg
	if len(errMsg) > 0 {
		msg = errMsg[0]
	}

	if len(val) > 254 || !emailRegex.MatchString(val) {
		f.Errors.Add(field, msg)
	}
}

// ValidateConfirmation value.
// Field name is a Pascal case string.
// Associated form field is this name in dashed case format.
// i.e.: field 'FamilyName' > mapKey 'family-name'
func (f *Form) ValidateConfirmation(field, confField string, errMsg ...string) {
	mapKey := pascalToDashed(field)
	val := f.Get(mapKey)

	mapConfKey := pascalToDashed(field)
	conf := f.Get(mapConfKey)

	msg := notMatchErrMsg
	if len(errMsg) > 0 {
		msg = errMsg[0]
	}

	if val != conf {
		f.Errors.Add(field, msg)
		f.Errors.Add(confField, msg)
	}
}

func (f *Form) Get(field string) string {
	return f.Changes[field][0]
}

func (f *Form) HasErrors() bool {
	return !f.IsValid()
}

func (f *Form) IsValid() bool {
	return len(f.Errors) == 0
}

func (es ErrorSet) Add(field, msg string) {
	es[field] = append(es[field], msg)
}

func (es ErrorSet) FieldErrors(field string) []string {
	return es[field]
}

func pascalToDashed(str string) string {
	str = uppercaseFirst(str)
	dashed := matchFirstCap.ReplaceAllString(str, "${1}-${2}")
	dashed = matchAllCap.ReplaceAllString(dashed, "${1}-${2}")
	return strings.ToLower(dashed)
}

func uppercaseFirst(str string) string {
	temp := []rune(str)
	temp[0] = unicode.ToUpper(temp[0])
	return string(temp)
}
