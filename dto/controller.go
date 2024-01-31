package dto

import (
	"fmt"
	"github.com/go-playground/validator"
	"github.com/hashicorp/go-multierror"
	"github.com/labstack/echo/v4"
	"net/http"
)

const Success string = "Success"

type Validator struct {
	valid *validator.Validate
}

func NewValidate() Validator {
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

func ResponseSuccess(c echo.Context, data ...interface{}) error {
	return c.JSON(http.StatusOK, ResponseData{
		Response: Response{
			Status:  http.StatusOK,
			Message: Success,
		},
		Data: data,
	})
}
