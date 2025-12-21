package users

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	"github.com/mesh-dell/expense-Tracker-API/internal/config"
)

type Tokens struct {
	Access     string
	Refresh    string
	JTIAccess  string
	JTIRefresh string
	ExpAccess  time.Time
	ExpRefresh time.Time
	UserID     uint
	Issuer     string
	Audience   string
}

func IssueTokens(userID uint, cfg config.Config) (*Tokens, error) {
	t := &Tokens{
		UserID:     userID,
		JTIAccess:  uuid.NewString(),
		JTIRefresh: uuid.NewString(),
		ExpAccess:  time.Now().Add(time.Duration(cfg.AccessExpiry) * time.Second),
		ExpRefresh: time.Now().Add(time.Duration(cfg.RefreshExpiry) * time.Second),
		Issuer:     "expense-tracker-app",
		Audience:   "expense-tracker-client",
	}

	acc := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.StandardClaims{
		Subject:   fmt.Sprint(userID),
		Id:        t.JTIAccess,
		Issuer:    t.Issuer,
		Audience:  t.Audience,
		IssuedAt:  time.Now().Unix(),
		ExpiresAt: t.ExpAccess.Unix(),
	})

	ref := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.StandardClaims{
		Subject:   fmt.Sprint(userID),
		Id:        t.JTIRefresh,
		Issuer:    t.Issuer,
		Audience:  t.Audience,
		IssuedAt:  time.Now().Unix(),
		ExpiresAt: t.ExpRefresh.Unix(),
	})

	var err error
	t.Access, err = acc.SignedString([]byte(cfg.AccessSecret))
	if err != nil {
		return nil, err
	}
	t.Refresh, err = ref.SignedString([]byte(cfg.RefreshSecret))
	if err != nil {
		return nil, err
	}
	return t, err
}
