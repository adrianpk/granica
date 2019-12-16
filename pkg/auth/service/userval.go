package service

import (
	"errors"

	"gitlab.com/mikrowezel/backend/granica/internal/model"
	"gitlab.com/mikrowezel/backend/service"
)

type (
	UserValidator struct {
		Model interface{}
		service.Validator
	}
)

func NewUserValidator(u model.User) UserValidator {
	return UserValidator{
		Model:     u,
		Validator: service.NewValidator(),
	}
}

func (uv UserValidator) ValidateForCreate() error {
	ok1 := uv.ValidateMinUsername(4, "min_length_err_msg")
	ok2 := uv.ValidateMaxUsername(16, "max_length_err_msg")
	if ok1 && ok2 {
		return nil
	}

	return errors.New("user has errors")
}

// NOTE: Update validations shoud be different
// than the ones executed on creation.
func (uv UserValidator) ValidateForUpdate() error {
	ok1 := uv.ValidateMinUsername(4, "min_length_err_msg")
	ok2 := uv.ValidateMaxUsername(16, "max_length_err_msg")
	if ok1 && ok2 {
		return nil
	}

	return errors.New("user has errors")
}

func (uv UserValidator) ValidateMinUsername(min int, errMsg ...string) (ok bool) {
	u, ok := uv.Model.(model.User)
	if !ok {
		return true
	}

	ok = uv.ValidateMinLength(u.Username.String, min)
	if ok {
		return true
	}

	msg := service.MinLengthErrMsg
	if len(errMsg) > 0 {
		msg = errMsg[0]
	}

	uv.Errors["Username"] = append(uv.Errors["Username"], msg)
	return false
}

func (uv UserValidator) ValidateMaxUsername(max int, errMsg ...string) (ok bool) {
	u, ok := uv.Model.(model.User)
	if !ok {
		return true
	}

	ok = uv.ValidateMaxLength(u.Username.String, max)
	if ok {
		return true
	}

	msg := service.MaxLengthErrMsg
	if len(errMsg) > 0 {
		msg = errMsg[0]
	}

	uv.Errors["Username"] = append(uv.Errors["Username"], msg)
	return false
}
