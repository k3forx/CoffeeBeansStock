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

type AuthHandler struct {
	authService *services.AuthService
}

func NewAuthHandler(authService *services.AuthService) *AuthHandler {
	return &AuthHandler{authService: authService}
}

type signupRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Name     string `json:"name"`
}

type loginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type refreshRequest struct {
	RefreshToken string `json:"refresh_token"`
}

func (h *AuthHandler) Signup(w http.ResponseWriter, r *http.Request) {
	var req signupRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		api.WriteError(w, http.StatusBadRequest, "BAD_REQUEST", "リクエストの形式が不正です")
		return
	}

	input := &services.SignupInput{
		Email:    req.Email,
		Password: req.Password,
		Name:     req.Name,
	}

	if errs := h.authService.ValidateSignupInput(input); len(errs) > 0 {
		api.WriteValidationError(w, errs)
		return
	}

	result, err := h.authService.Signup(r.Context(), input)
	if err != nil {
		if errors.Is(err, services.ErrEmailAlreadyExists) {
			api.WriteError(w, http.StatusConflict, "CONFLICT", "このメールアドレスは既に登録されています")
			return
		}
		api.WriteError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "サーバーエラーが発生しました")
		return
	}

	api.WriteSuccess(w, http.StatusCreated, result)
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req loginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		api.WriteError(w, http.StatusBadRequest, "BAD_REQUEST", "リクエストの形式が不正です")
		return
	}

	input := &services.LoginInput{
		Email:    req.Email,
		Password: req.Password,
	}

	if errs := h.authService.ValidateLoginInput(input); len(errs) > 0 {
		api.WriteValidationError(w, errs)
		return
	}

	result, err := h.authService.Login(r.Context(), input)
	if err != nil {
		if errors.Is(err, services.ErrInvalidCredentials) {
			api.WriteError(w, http.StatusUnauthorized, "UNAUTHORIZED", "メールアドレスまたはパスワードが正しくありません")
			return
		}
		api.WriteError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "サーバーエラーが発生しました")
		return
	}

	api.WriteSuccess(w, http.StatusOK, result)
}

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

