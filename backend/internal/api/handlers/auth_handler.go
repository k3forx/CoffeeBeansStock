package handlers

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/k3forx/CoffeeBeansStock/backend/internal/api"
	"github.com/k3forx/CoffeeBeansStock/backend/internal/api/gen"
	"github.com/k3forx/CoffeeBeansStock/backend/internal/api/middleware"
	domain "github.com/k3forx/CoffeeBeansStock/backend/internal/domain"
	"github.com/k3forx/CoffeeBeansStock/backend/internal/domain/user"
	"github.com/k3forx/CoffeeBeansStock/backend/internal/usecase"
)

// AuthHandler handles authentication-related HTTP requests.
type AuthHandler struct {
	authService *usecase.AuthService
}

// NewAuthHandler creates a new AuthHandler.
func NewAuthHandler(authService *usecase.AuthService) *AuthHandler {
	return &AuthHandler{authService: authService}
}

func toUserResponse(u *user.User) gen.UserResponse {
	resp := gen.UserResponse{
		Id:        u.ID(),
		CreatedAt: u.CreatedAt(),
	}
	if name := u.Name(); name != "" {
		resp.Name = &name
	}
	if threshold := u.LowStockThreshold(); threshold != 0 {
		resp.LowStockThreshold = &threshold
	}
	if enabled := u.NotificationEnabled(); enabled {
		resp.NotificationEnabled = &enabled
	}
	if !u.UpdatedAt().IsZero() {
		updatedAt := u.UpdatedAt()
		resp.UpdatedAt = &updatedAt
	}
	return resp
}

// RegisterAnonymous handles anonymous user registration.
func (h *AuthHandler) RegisterAnonymous(w http.ResponseWriter, r *http.Request) {
	result, err := h.authService.RegisterAnonymous(r.Context())
	if err != nil {
		api.WriteError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "サーバーエラーが発生しました")
		return
	}

	api.WriteSuccess(w, http.StatusCreated, gen.AuthResponse{
		User:         toUserResponse(result.User),
		AccessToken:  result.AccessToken,
		RefreshToken: result.RefreshToken,
	})
}

// Refresh handles token refresh requests.
func (h *AuthHandler) Refresh(w http.ResponseWriter, r *http.Request) {
	var req gen.RefreshRequest
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

	result, err := h.authService.Refresh(r.Context(), usecase.RefreshInput{
		RefreshToken: req.RefreshToken,
	})
	if err != nil {
		if errors.Is(err, domain.ErrInvalidToken) || errors.Is(err, domain.ErrExpiredToken) {
			api.WriteError(w, http.StatusUnauthorized, "UNAUTHORIZED", "リフレッシュトークンが無効です")
			return
		}
		api.WriteError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "サーバーエラーが発生しました")
		return
	}

	api.WriteSuccess(w, http.StatusOK, gen.RefreshResult{
		AccessToken:  result.AccessToken,
		RefreshToken: result.RefreshToken,
	})
}

// GetMe returns the authenticated user's profile.
func (h *AuthHandler) GetMe(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserID(r.Context())
	if !ok {
		api.WriteError(w, http.StatusUnauthorized, "UNAUTHORIZED", "認証が必要です")
		return
	}

	output, err := h.authService.GetMe(r.Context(), usecase.GetMeInput{UserID: userID})
	if err != nil {
		api.WriteError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "サーバーエラーが発生しました")
		return
	}

	api.WriteSuccess(w, http.StatusOK, toUserResponse(output.User))
}
