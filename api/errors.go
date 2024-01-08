package api

import (
	"github.com/gofiber/fiber/v2"
	"net/http"
)

type Error struct {
	Code int    `json:"code"`
	Err  string `json:"error"`
}

func ErrorHandler(c *fiber.Ctx, err error) error {
	if apiError, ok := err.(Error); ok {
		return c.Status(apiError.Code).JSON(apiError)
	}
	apiError := NewError(http.StatusInternalServerError, err.Error())
	return c.Status(apiError.Code).JSON(apiError)
}

// Error implements the Error interface
func (e Error) Error() string {
	return e.Err
}
func NewError(code int, msg string) Error {
	return Error{
		Code: code,
		Err:  msg,
	}
}

func ErrInvalidId() Error {
	return Error{
		Code: http.StatusBadRequest,
		Err:  "Invalid id",
	}
}

func ErrUnauthorized() Error {
	return Error{
		Code: http.StatusUnauthorized,
		Err:  "Unauthorized request",
	}
}

func ErrBadRequest() Error {
	return Error{
		Code: http.StatusBadRequest,
		Err:  "Invalid JSON request",
	}
}

func ErrResourceNotFound(res string) Error {
	return Error{
		Code: http.StatusNotFound,
		Err:  res + "Resource not found",
	}
}
