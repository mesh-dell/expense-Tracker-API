package service

import (
	"context"
	"time"

	"github.com/mesh-dell/expense-Tracker-API/internal/users"
	"github.com/mesh-dell/expense-Tracker-API/internal/users/repository"
)

type RefreshTokenService struct {
	repo repository.RefreshTokenRepository
}

func NewRefreshTokenService(r repository.RefreshTokenRepository) *RefreshTokenService {
	return &RefreshTokenService{
		repo: r,
	}
}

func (svc *RefreshTokenService) Save(ctx context.Context, tokens *users.Tokens) error {
	return svc.repo.Save(ctx, tokens)
}

func (svc *RefreshTokenService) RotateRefreshToken(ctx context.Context, oldJTI string, newTokens *users.Tokens) error {
	_ = svc.repo.Delete(ctx, oldJTI)
	return svc.Save(ctx, newTokens)
}

func (svc *RefreshTokenService) ValidateRefreshToken(ctx context.Context, jti string) (*users.RefreshToken, bool) {
	stored, _ := svc.repo.Find(ctx, jti)
	if stored == nil {
		return nil, false
	}
	if time.Now().After(stored.ExpiresAt) {
		svc.repo.Delete(ctx, stored.JTI)
		return nil, false
	}
	return stored, true
}
