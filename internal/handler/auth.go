package handler

import (
	"encoding/json"
	"net/http"
	"time"

	"test-backend/internal/response"
	"test-backend/internal/service"
	apperr "test-backend/pkg/errors"

	"github.com/golang-jwt/jwt/v5"
)

// AuthHandler handles login and JWT issuance.
type AuthHandler struct {
	userSvc       *service.UserService
	jwtSecret     string
	jwtExpiryHours int
}

// NewAuthHandler returns a new AuthHandler.
func NewAuthHandler(userSvc *service.UserService, jwtSecret string, jwtExpiryHours int) *AuthHandler {
	return &AuthHandler{userSvc: userSvc, jwtSecret: jwtSecret, jwtExpiryHours: jwtExpiryHours}
}

// LoginRequest is the body for POST /api/v1/login.
type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

// Login handles POST /api/v1/login.
func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, apperr.NewBadRequest("Invalid request body. Expected valid JSON."))
		return
	}
	if req.Email == "" || req.Password == "" {
		response.Error(w, apperr.NewBadRequest("Email and password are required."))
		return
	}
	user, err := h.userSvc.Login(req.Email, req.Password)
	if err != nil {
		response.Error(w, err)
		return
	}
	if h.jwtSecret == "" {
		response.Error(w, apperr.NewUnauthorized("Authentication is not configured."))
		return
	}
	exp := time.Now().Add(time.Duration(h.jwtExpiryHours) * time.Hour)
	claims := jwt.RegisteredClaims{
		Subject:   user.ID,
		ExpiresAt: jwt.NewNumericDate(exp),
		IssuedAt:  jwt.NewNumericDate(time.Now()),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenStr, err := token.SignedString([]byte(h.jwtSecret))
	if err != nil {
		response.Error(w, apperr.NewUnauthorized("Failed to issue token."))
		return
	}
	expiresIn := int(exp.Sub(time.Now()).Seconds())
	data := map[string]any{
		"token":           tokenStr,
		"user":            user,
		"expires_in":      expiresIn,
	}
	response.Success(w, "Login successful.", data)
}
