package repository

import (
	"context"
	"errors"

	"github.com/mesh-dell/expense-Tracker-API/internal/custom"
	"github.com/mesh-dell/expense-Tracker-API/internal/users"
	"gorm.io/gorm"
)

type RefreshTokenRepository struct {
	db *gorm.DB
}

func NewRefreshTokenRepository(db *gorm.DB) *RefreshTokenRepository {
	return &RefreshTokenRepository{
		db: db,
	}
}

func (r *RefreshTokenRepository) Save(ctx context.Context, tokens *users.Tokens) error {
	ref := &users.RefreshToken{
		UserID:    tokens.UserID,
		JTI:       tokens.JTIRefresh,
		ExpiresAt: tokens.ExpRefresh,
	}
	return gorm.G[users.RefreshToken](r.db).Create(ctx, ref)
}

func (r *RefreshTokenRepository) Find(ctx context.Context, jti string) (*users.RefreshToken, error) {
	ref, err := gorm.G[users.RefreshToken](r.db).Where("jti = ?", jti).First(ctx)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, custom.ErrRefreshTokenNotFound
	}
	if err != nil {
		return nil, err
	}
	return &ref, nil
}

func (r *RefreshTokenRepository) Delete(ctx context.Context, jti string) error {
	_, err := gorm.G[users.RefreshToken](r.db).Where("jti = ?", jti).Delete(ctx)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil
		} else {
			return err
		}
	}
	return nil
}

func (r *RefreshTokenRepository) DeleteAllForUser(ctx context.Context, userID uint) error {
	_, err := gorm.G[users.RefreshToken](r.db).Where("user_id = ?").Delete(ctx)
	return err
}
