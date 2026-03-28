package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/google/uuid"

	"github.com/k3forx/CoffeeBeansStock/backend/internal/api"
	"github.com/k3forx/CoffeeBeansStock/backend/internal/api/middleware"
	domain "github.com/k3forx/CoffeeBeansStock/backend/internal/domain"
	"github.com/k3forx/CoffeeBeansStock/backend/internal/domain/user"
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

type userResponse struct {
	ID                  uuid.UUID `json:"id"`
	Name                string    `json:"name,omitempty"`
	LowStockThreshold   int32     `json:"low_stock_threshold,omitempty"`
	NotificationEnabled bool      `json:"notification_enabled,omitempty"`
	CreatedAt           string    `json:"created_at"`
	UpdatedAt           string    `json:"updated_at,omitempty"`
}

func toUserResponse(u *user.User) userResponse {
	resp := userResponse{
		ID:                  u.ID(),
		Name:                u.Name(),
		LowStockThreshold:   u.LowStockThreshold(),
		NotificationEnabled: u.NotificationEnabled(),
	}
	if !u.CreatedAt().IsZero() {
		resp.CreatedAt = u.CreatedAt().Format(time.RFC3339)
	}
	if !u.UpdatedAt().IsZero() {
		resp.UpdatedAt = u.UpdatedAt().Format(time.RFC3339)
	}
	return resp
}

type authResponse struct {
	User         userResponse `json:"user"`
	AccessToken  string       `json:"access_token"`  //nolint:gosec
	RefreshToken string       `json:"refresh_token"` //nolint:gosec
}

// RegisterAnonymous handles anonymous user registration.
func (h *AuthHandler) RegisterAnonymous(w http.ResponseWriter, r *http.Request) {
	result, err := h.authService.RegisterAnonymous(r.Context())
	if err != nil {
		api.WriteError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "サーバーエラーが発生しました")
		return
	}

	api.WriteSuccess(w, http.StatusCreated, authResponse{
		User:         toUserResponse(result.User),
		AccessToken:  result.AccessToken,
		RefreshToken: result.RefreshToken,
	})
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
		if errors.Is(err, domain.ErrInvalidToken) || errors.Is(err, domain.ErrExpiredToken) {
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

	u, err := h.authService.GetMe(r.Context(), userID)
	if err != nil {
		api.WriteError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "サーバーエラーが発生しました")
		return
	}

	api.WriteSuccess(w, http.StatusOK, toUserResponse(u))
}
