package handlers

import (
	"net/http"

	"github.com/k3forx/CoffeeBeansStock/backend/internal/api"
	"github.com/k3forx/CoffeeBeansStock/backend/internal/api/gen"
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
	if gpc := u.GramsPerCup().Value(); gpc != 0 {
		resp.GramsPerCup = &gpc
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
	if !decodeJSON(w, r, &req) {
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
		handleDomainError(w, err, "リフレッシュトークンが無効です")
		return
	}

	api.WriteSuccess(w, http.StatusOK, gen.RefreshResult{
		AccessToken:  result.AccessToken,
		RefreshToken: result.RefreshToken,
	})
}

// GetMe returns the authenticated user's profile.
func (h *AuthHandler) GetMe(w http.ResponseWriter, r *http.Request) {
	userID, ok := requireUserID(w, r)
	if !ok {
		return
	}

	output, err := h.authService.GetMe(r.Context(), usecase.GetMeInput{UserID: userID})
	if err != nil {
		api.WriteError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "サーバーエラーが発生しました")
		return
	}

	api.WriteSuccess(w, http.StatusOK, toUserResponse(output.User))
}

// UpdateMe updates the authenticated user's settings.
func (h *AuthHandler) UpdateMe(w http.ResponseWriter, r *http.Request) {
	userID, ok := requireUserID(w, r)
	if !ok {
		return
	}

	var req gen.UpdateUserRequest
	if !decodeJSON(w, r, &req) {
		return
	}

	output, err := h.authService.UpdateMe(r.Context(), usecase.UpdateMeInput{
		UserID:              userID,
		GramsPerCup:         req.GramsPerCup,
		LowStockThreshold:   req.LowStockThreshold,
		NotificationEnabled: req.NotificationEnabled,
	})
	if err != nil {
		handleDomainError(w, err, "ユーザーが見つかりません")
		return
	}

	api.WriteSuccess(w, http.StatusOK, toUserResponse(output.User))
}
