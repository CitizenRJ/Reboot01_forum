package user

import "context"

type userService struct {
	repo Repository
}

type Repository interface {
	CreateUser(ctx context.Context, arg CreateUserParams) (*User, error)
}

func NewUserService(r Repository) *userService {
	return &userService{repo: r}
}

func (s *userService) CreateUser(ctx context.Context, arg CreateUserParams) (*User, error) {
	// do some business logic and validation
	return s.repo.CreateUser(ctx, arg)
}
