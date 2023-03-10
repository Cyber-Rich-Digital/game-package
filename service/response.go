package service

import (
	"net/http"
)

type ResponseError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func (e ResponseError) Error() string {
	return e.Message
}

func badRequest(msg string) error {
	return ResponseError{
		Code:    http.StatusBadRequest,
		Message: msg,
	}
}

func notFound(msg string) error {
	return ResponseError{
		Code:    http.StatusNotFound,
		Message: msg,
	}
}

func internalServerError(msg string) error {
	return ResponseError{
		Code:    http.StatusInternalServerError,
		Message: msg,
	}
}
