package service

import (
	"context"
	"fmt"

	"github.com/mesh-dell/expense-Tracker-API/internal/custom"
	"github.com/mesh-dell/expense-Tracker-API/internal/users"
	"github.com/mesh-dell/expense-Tracker-API/internal/users/dtos"
	"github.com/mesh-dell/expense-Tracker-API/internal/users/repository"
	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	repo repository.IUserRepository
}

func NewUserService(r repository.IUserRepository) *UserService {
	return &UserService{
		repo: r,
	}
}

func (svc *UserService) Login(ctx context.Context, req dtos.LoginRequest) (*users.User, error) {
	user, err := svc.repo.FindUserByEmail(ctx, req.Email)
	if err != nil {
		return nil, fmt.Errorf("get user error: %v", err)
	}
	if user == nil {
		return nil, custom.ErrInvalidCredentials
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		return nil, custom.ErrInvalidCredentials
	}
	return user, nil
}

func (svc *UserService) Register(ctx context.Context, req dtos.RegisterRequest) (*users.User, error) {
	if exists, _ := svc.repo.FindUserByEmail(ctx, req.Email); exists != nil {
		return nil, custom.ErrEmailAlreadyExists
	}
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("error hashing password: %v", err)
	}
	user := &users.User{
		Name:         req.Name,
		Email:        req.Email,
		PasswordHash: string(passwordHash),
	}
	err = svc.repo.CreateUser(ctx, user)
	if err != nil {
		return nil, fmt.Errorf("error creating user: %v", err)
	}
	return user, nil
}
