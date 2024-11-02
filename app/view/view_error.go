package view

import (
	"fmt"
	"net/http"
)

type ViewError struct {
	Code    int64
	Message error
}

func (e ViewError) Error() string {
	return fmt.Sprintf("[%d]%s", e.Code, e.Message)
}

func NewBadRequestView(err error) ViewError {
	return ViewError{
		Code:    http.StatusBadRequest,
		Message: err,
	}
}

func NewUnauthorizedView(err error) ViewError {
	return ViewError{
		Code:    http.StatusUnauthorized,
		Message: err,
	}
}

func NewNotFoundView(err error) ViewError {
	return ViewError{
		Code:    http.StatusNotFound,
		Message: err,
	}
}

func NewInternalServerErrorView(err error) ViewError {
	return ViewError{
		Code:    http.StatusInternalServerError,
		Message: err,
	}
}
