package user

import (
	"context"
	"encoding/json"
	e "forum/error"
	"net/http"
)

type UserService interface {
	CreateUser(ctx context.Context, arg CreateUserParams) (*User, error)
}

type UserController struct {
	s UserService
}

func NewUserController(s UserService) *UserController {
	return &UserController{s: s}
}

func (c *UserController) CreateUser(w http.ResponseWriter, r *http.Request) (User, error) {
	var arg CreateUserParams
	if err := json.NewDecoder(r.Body).Decode(&arg); err != nil {
		e.HandleBadRequest(w, r)
		return User{}, err
	}

	// TODO: have proper http status codes and error handling
	user, _ := c.s.CreateUser(r.Context(), arg)
	// handle errors

	return *user, nil
}
