package handlers

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/k3forx/CoffeeBeansStock/backend/internal/api"
	"github.com/k3forx/CoffeeBeansStock/backend/internal/api/middleware"
	"github.com/k3forx/CoffeeBeansStock/backend/internal/auth"
	"github.com/k3forx/CoffeeBeansStock/backend/internal/services"
)

// AuthHandler handles authentication-related HTTP requests.
type AuthHandler struct {
	authService *services.AuthService
}

// NewAuthHandler creates a new AuthHandler.
func NewAuthHandler(authService *services.AuthService) *AuthHandler {
	return &AuthHandler{authService: authService}
}

type refreshRequest struct {
	RefreshToken string `json:"refresh_token"` //nolint:gosec // request field, not a hardcoded secret
}

// RegisterAnonymous handles anonymous user registration.
func (h *AuthHandler) RegisterAnonymous(w http.ResponseWriter, r *http.Request) {
	result, err := h.authService.RegisterAnonymous(r.Context())
	if err != nil {
		api.WriteError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "サーバーエラーが発生しました")
		return
	}

	api.WriteSuccess(w, http.StatusCreated, result)
}

// Refresh handles token refresh requests.
func (h *AuthHandler) Refresh(w http.ResponseWriter, r *http.Request) {
	var req refreshRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		api.WriteError(w, http.StatusBadRequest, "BAD_REQUEST", "リクエストの形式が不正です")
		return
	}

	if req.RefreshToken == "" {
		api.WriteValidationError(w, []api.FieldError{
			{Field: "refresh_token", Message: "リフレッシュトークンは必須です"},
		})
		return
	}

	result, err := h.authService.Refresh(r.Context(), req.RefreshToken)
	if err != nil {
		if errors.Is(err, auth.ErrInvalidToken) || errors.Is(err, auth.ErrExpiredToken) {
			api.WriteError(w, http.StatusUnauthorized, "UNAUTHORIZED", "リフレッシュトークンが無効です")
			return
		}
		api.WriteError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "サーバーエラーが発生しました")
		return
	}

	api.WriteSuccess(w, http.StatusOK, result)
}

// GetMe returns the authenticated user's profile.
func (h *AuthHandler) GetMe(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserID(r.Context())
	if !ok {
		api.WriteError(w, http.StatusUnauthorized, "UNAUTHORIZED", "認証が必要です")
		return
	}

	user, err := h.authService.GetMe(r.Context(), userID)
	if err != nil {
		api.WriteError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "サーバーエラーが発生しました")
		return
	}

	api.WriteSuccess(w, http.StatusOK, user)
}
