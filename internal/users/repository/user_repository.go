package repository

import (
	"context"
	"errors"

	"github.com/mesh-dell/expense-Tracker-API/internal/users"
	"gorm.io/gorm"
)

type UserRepository struct {
	db *gorm.DB
}

// CreateUser implements [IUserRepository].
func (r *UserRepository) CreateUser(ctx context.Context, user *users.User) error {
	return gorm.G[users.User](r.db).Create(ctx, user)
}

// FindUserByEmail implements [IUserRepository].
func (r *UserRepository) FindUserByEmail(ctx context.Context, email string) (*users.User, error) {
	user, err := gorm.G[users.User](r.db).Where("email = ?", email).First(ctx)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		} else {
			return nil, err
		}
	}
	return &user, nil
}

func NewUserRepository(db *gorm.DB) IUserRepository {
	return &UserRepository{
		db: db,
	}
}
