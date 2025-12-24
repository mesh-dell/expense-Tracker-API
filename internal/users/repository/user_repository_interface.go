package repository

import (
	"context"

	"github.com/mesh-dell/expense-Tracker-API/internal/users"
)

type IUserRepository interface {
	FindUserByEmail(ctx context.Context, email string) (*users.User, error)
	CreateUser(ctx context.Context, user *users.User) error
	FindUserByID(ctx context.Context, userID uint) (users.User, error)
}
