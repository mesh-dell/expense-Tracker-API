package handler

import (
	"errors"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"github.com/mesh-dell/expense-Tracker-API/internal/config"
	"github.com/mesh-dell/expense-Tracker-API/internal/custom"
	"github.com/mesh-dell/expense-Tracker-API/internal/users"
	"github.com/mesh-dell/expense-Tracker-API/internal/users/dtos"
	"github.com/mesh-dell/expense-Tracker-API/internal/users/service"
)

type UserHandler struct {
	userSvc    *service.UserService
	refreshSvc *service.RefreshTokenService
	cfg        config.Config
}

func NewUserHandler(u *service.UserService, cfg config.Config, r *service.RefreshTokenService) *UserHandler {
	return &UserHandler{
		userSvc:    u,
		cfg:        cfg,
		refreshSvc: r,
	}
}

func (h *UserHandler) Register(c *gin.Context) {
	var req dtos.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}
	user, err := h.userSvc.Register(c.Request.Context(), req)
	if err != nil {
		if errors.Is(err, custom.ErrEmailAlreadyExists) {
			c.IndentedJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		} else {
			c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	}
	t, err := users.IssueTokens(user.ID, h.cfg)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": "could not issue tokens"})
		return
	}
	err = h.refreshSvc.Save(c.Request.Context(), t)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": "could not save token"})
		return
	}
	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie("access_token", t.Access, int(time.Until(t.ExpAccess).Seconds()), "/", "", true, true)
	c.SetCookie("refresh_token", t.Refresh, int(time.Until(t.ExpRefresh).Seconds()), "/", "", true, true)
	c.IndentedJSON(http.StatusCreated, dtos.AuthResponse{
		User: struct {
			ID    uint   "json:\"id\""
			Email string "json:\"email\""
			Name  string "json:\"name\""
		}{
			ID:    user.ID,
			Email: user.Email,
			Name:  user.Name,
		},
	})
}

func (h *UserHandler) Login(c *gin.Context) {
	var req dtos.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}
	user, err := h.userSvc.Login(c.Request.Context(), req)
	if err != nil {
		if errors.Is(err, custom.ErrInvalidCredentials) {
			c.IndentedJSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			return
		} else {
			c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	}
	t, err := users.IssueTokens(user.ID, h.cfg)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": "could not issue tokens"})
		return
	}
	err = h.refreshSvc.Save(c.Request.Context(), t)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": "could not save token"})
		return
	}
	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie("access_token", t.Access, int(time.Until(t.ExpAccess).Seconds()), "/", "", true, true)
	c.SetCookie("refresh_token", t.Refresh, int(time.Until(t.ExpRefresh).Seconds()), "/", "", true, true)

	c.IndentedJSON(http.StatusOK, dtos.AuthResponse{
		User: struct {
			ID    uint   "json:\"id\""
			Email string "json:\"email\""
			Name  string "json:\"name\""
		}{
			ID:    user.ID,
			Email: user.Email,
			Name:  user.Name,
		},
	})
}

func (h *UserHandler) RefreshToken(c *gin.Context) {
	ref, err := c.Cookie("refresh_token")
	if err != nil {
		c.IndentedJSON(http.StatusUnauthorized, gin.H{"error": "missing refresh token"})
		return
	}
	claims := &jwt.StandardClaims{}
	parsed, err := jwt.ParseWithClaims(ref, claims, func(t *jwt.Token) (any, error) {
		return []byte(h.cfg.RefreshSecret), nil
	})
	if err != nil || !parsed.Valid {
		c.IndentedJSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
		return
	}
	stored, ok := h.refreshSvc.ValidateRefreshToken(c.Request.Context(), claims.Id)
	if !ok {
		c.IndentedJSON(http.StatusUnauthorized, gin.H{"error": "refresh revoked or expired"})
		return
	}
	tokens, err := users.IssueTokens(stored.UserID, h.cfg)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": "could not issue new tokens"})
		return
	}
	h.refreshSvc.RotateRefreshToken(c.Request.Context(), claims.Id, tokens)
	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie("access_token", tokens.Access, int(time.Until(tokens.ExpAccess).Seconds()), "/", "", true, true)
	c.SetCookie("refresh_token", tokens.Refresh, int(time.Until(tokens.ExpRefresh).Seconds()), "/", "", true, true)

	c.IndentedJSON(http.StatusOK, gin.H{"ok": "refreshed token successfully"})
}

func (h *UserHandler) GetMe(c *gin.Context) {
	userIDAny, exists := c.Get("userID")
	if !exists {
		c.IndentedJSON(http.StatusUnauthorized, gin.H{"user": nil})
		return
	}
	userID := userIDAny.(uint)

	user, err := h.userSvc.GetMe(c.Request.Context(), userID)
	if err != nil {
		c.IndentedJSON(http.StatusUnauthorized, gin.H{"user": nil})
		return
	}

	c.IndentedJSON(http.StatusOK, gin.H{
		"user": gin.H{
			"id":    user.ID,
			"email": user.Email,
			"name":  user.Name,
		},
	})
}
