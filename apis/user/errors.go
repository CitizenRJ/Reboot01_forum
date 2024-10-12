package user

import (
	"errors"
	"net/http"
)

type UserError struct {
	Status int   `json:"status"`
	Error  error `json:"error"`
}

var (
	UserNotFoundErr = UserError{http.StatusNotFound, errors.New("user not found")}
)

type errorData struct {
	Error string
}
