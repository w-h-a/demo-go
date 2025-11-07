package userrepo

import (
	"context"

	"github.com/w-h-a/demo-go/api/user"
)

// UserRepo is the interface for our user data store.
type UserRepo interface {
	Create(ctx context.Context, dto user.CreateUserDTO) (user.User, error)
	GetByID(ctx context.Context, id string) (user.User, error)
	GetByEmail(ctx context.Context, email string) (user.User, error)
	GetAll(ctx context.Context, opts ...GetAllOption) ([]user.User, error)
}
