package auth

import (
	"fmt"
	"github.com/go-playground/validator"
	"github.com/hashicorp/go-multierror"
)

type Validator struct {
	valid *validator.Validate
}

type ValidatorInterface interface {
	Validate(i interface{}) error
}

func NewValidate() ValidatorInterface {
	return Validator{valid: validator.New()}
}

func (r Validator) Validate(i interface{}) error {
	var result error
	if err := r.valid.Struct(i); err != nil {
		if _, ok := err.(*validator.InvalidValidationError); ok {
			result = multierror.Append(result, err)
			return result
		}

		for _, err := range err.(validator.ValidationErrors) {
			switch err.Tag() {
			case "required":
				errRequired := fmt.Errorf("%s is required", err.Field())
				result = multierror.Append(result, errRequired)
			default:
				errUndefined := fmt.Errorf("error is undefined")
				result = multierror.Append(result, errUndefined)
			}
		}
		return result
	}
	return nil
}
