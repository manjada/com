package dto

import (
	"fmt"
	"github.com/manjada/com/web"
	"net/http"
)

const (
	generalError = "An error occurred, please try again later or contact the administrator"
)
const (
	dbErrorCode       = 1001
	readFileErrorCode = 1002
)

var ERR_INVALID_EMAIL_FORMAT = ErrCodeUser(1003, "Invalid Email Format")
var ERR_PARSE_JSON = ErrCodeUser(1004, "Invalid json provided")
var ERR_VALIDATE_REQUIRED = ErrCodeUser(1005, "Field Required")
var ERR_USER_NOT_FOUND = ErrCodeUser(1005, "user or password is wrong")
var ERR_SYSTEM = ErrCodeUser(1006, "An error occurred in the system, please try again or contact the administrator")
var ERR_OTP_MAX = ErrCodeUser(1007, "Percobaan otp sudah lebih dari batas maksimal silahkan request otp kembali")
var ERR_OTP_INVALID = ErrCodeUser(1007, "Otp tidak sesuai")
var ERR_DATA_NOT_FOUND = ErrCodeUser(1008, "Error Data Not Found")
var ERR_DATA_EXISTS = ErrCodeUser(1009, "Data Exists")
var ERR_TOKEN_EXPIRED = ErrCodeUser(1010, "Invalid or expired token")
var ERR_UNAUTHORIZED = ErrCodeUser(1011, "Unauthorized access")

type ErrorCustom struct {
	CodeError int
	Desc      string
}

type SystemError struct {
	Wrap      string
	Message   string
	CodeError int
}

type CustomErrorInterface interface {
	Code() int
	InvalidResponse(c web.Context) error
}

func (o *SystemError) InvalidResponse(c web.Context) error {
	if o != nil {
		return c.JSON(http.StatusBadRequest, Response{
			Status:    http.StatusBadRequest,
			ErrorCode: o.Code(),
			Message:   o.Error(),
		})
	}
	return nil
}

func (o *ErrorCustom) InvalidResponse(c web.Context) error {
	if o != nil {
		return c.JSON(http.StatusBadRequest, Response{
			Status:    http.StatusBadRequest,
			ErrorCode: o.Code(),
			Message:   o.Error(),
		})
	}
	return nil
}

func (o *ErrorCustom) Code() int {
	return o.CodeError
}

func (o *SystemError) Code() int {
	return o.CodeError
}

func (e ErrorCustom) Error() string {
	return fmt.Sprintf("%s", e.Desc)
}

func (e SystemError) Error() string {
	return fmt.Sprintf("%s (%d): %s ", e.Wrap, e.CodeError, e.Message)
}

func ErrorDb(err error) *SystemError {
	return &SystemError{"Error Db process", err.Error(), dbErrorCode}
}

func ErrorParse(err error) *SystemError {
	return &SystemError{"Error parse", err.Error(), dbErrorCode}
}

func ErrorReadFile(err error) *SystemError {
	return &SystemError{"Failed to read file", err.Error(), readFileErrorCode}
}

func ErrorSystem(err error) *SystemError {
	return &SystemError{"Error from system", err.Error(), dbErrorCode}
}

func ErrorUser(data *ErrorCustom, s string) *ErrorCustom {
	if s != "" {
		data.Desc = s
	}
	return data
}

func ErrCodeUser(code int, desc string) *ErrorCustom {
	return &ErrorCustom{
		CodeError: code,
		Desc:      desc,
	}
}
