package custom

import "errors"

var (
	ErrInvalidCredentials   = errors.New("invalid email or password")
	ErrEmailAlreadyExists   = errors.New("email already exists")
	ErrRefreshTokenNotFound = errors.New("refresh token not found")
	ErrExpenseNotFound      = errors.New("expense not found")
	ErrInvalidCategory      = errors.New("invalid expense category")
)
