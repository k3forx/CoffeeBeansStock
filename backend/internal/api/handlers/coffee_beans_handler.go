package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/k3forx/CoffeeBeansStock/backend/internal/api"
	"github.com/k3forx/CoffeeBeansStock/backend/internal/api/middleware"
	"github.com/k3forx/CoffeeBeansStock/backend/internal/services"
)

// CoffeeBeansHandler handles coffee beans CRUD HTTP requests.
type CoffeeBeansHandler struct {
	service *services.CoffeeBeansService
}

// NewCoffeeBeansHandler creates a new CoffeeBeansHandler.
func NewCoffeeBeansHandler(service *services.CoffeeBeansService) *CoffeeBeansHandler {
	return &CoffeeBeansHandler{service: service}
}

type createBeanRequest struct {
	Origin       *string `json:"origin"`
	RoastLevel   *string `json:"roast_level"`
	Notes        *string `json:"notes"`
	Name         string  `json:"name"`
	CurrentStock int32   `json:"current_stock"`
}

type updateBeanRequest struct {
	Name         *string `json:"name"`
	Origin       *string `json:"origin"`
	RoastLevel   *string `json:"roast_level"`
	CurrentStock *int32  `json:"current_stock"`
	Notes        *string `json:"notes"`
}

// List returns a paginated list of coffee beans.
func (h *CoffeeBeansHandler) List(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserID(r.Context())
	if !ok {
		api.WriteError(w, http.StatusUnauthorized, "UNAUTHORIZED", "認証が必要です")
		return
	}

	limit := int32(parseQueryInt(r, "limit", 20))
	offset := int32(parseQueryInt(r, "offset", 0))

	result, err := h.service.List(r.Context(), userID, limit, offset)
	if err != nil {
		api.WriteError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "サーバーエラーが発生しました")
		return
	}

	api.WriteSuccess(w, http.StatusOK, result)
}

// Create handles coffee bean creation.
func (h *CoffeeBeansHandler) Create(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserID(r.Context())
	if !ok {
		api.WriteError(w, http.StatusUnauthorized, "UNAUTHORIZED", "認証が必要です")
		return
	}

	var req createBeanRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		api.WriteError(w, http.StatusBadRequest, "BAD_REQUEST", "リクエストの形式が不正です")
		return
	}

	input := &services.CreateBeanInput{
		Name:         req.Name,
		Origin:       req.Origin,
		RoastLevel:   req.RoastLevel,
		CurrentStock: req.CurrentStock,
		Notes:        req.Notes,
	}

	if errs := h.service.ValidateCreateInput(input); len(errs) > 0 {
		api.WriteValidationError(w, errs)
		return
	}

	bean, err := h.service.Create(r.Context(), userID, input)
	if err != nil {
		api.WriteError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "サーバーエラーが発生しました")
		return
	}

	api.WriteSuccess(w, http.StatusCreated, bean)
}

// Get returns a single coffee bean by ID.
func (h *CoffeeBeansHandler) Get(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserID(r.Context())
	if !ok {
		api.WriteError(w, http.StatusUnauthorized, "UNAUTHORIZED", "認証が必要です")
		return
	}

	beanID, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		api.WriteError(w, http.StatusBadRequest, "BAD_REQUEST", "不正なIDです")
		return
	}

	bean, err := h.service.GetByID(r.Context(), userID, beanID)
	if err != nil {
		handleBeanError(w, err)
		return
	}

	api.WriteSuccess(w, http.StatusOK, bean)
}

// Update handles coffee bean updates.
func (h *CoffeeBeansHandler) Update(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserID(r.Context())
	if !ok {
		api.WriteError(w, http.StatusUnauthorized, "UNAUTHORIZED", "認証が必要です")
		return
	}

	beanID, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		api.WriteError(w, http.StatusBadRequest, "BAD_REQUEST", "不正なIDです")
		return
	}

	var req updateBeanRequest
	if err = json.NewDecoder(r.Body).Decode(&req); err != nil {
		api.WriteError(w, http.StatusBadRequest, "BAD_REQUEST", "リクエストの形式が不正です")
		return
	}

	input := &services.UpdateBeanInput{
		Name:         req.Name,
		Origin:       req.Origin,
		RoastLevel:   req.RoastLevel,
		CurrentStock: req.CurrentStock,
		Notes:        req.Notes,
	}

	bean, err := h.service.Update(r.Context(), userID, beanID, input)
	if err != nil {
		handleBeanError(w, err)
		return
	}

	api.WriteSuccess(w, http.StatusOK, bean)
}

// Delete handles coffee bean deletion.
func (h *CoffeeBeansHandler) Delete(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserID(r.Context())
	if !ok {
		api.WriteError(w, http.StatusUnauthorized, "UNAUTHORIZED", "認証が必要です")
		return
	}

	beanID, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		api.WriteError(w, http.StatusBadRequest, "BAD_REQUEST", "不正なIDです")
		return
	}

	if err := h.service.Delete(r.Context(), userID, beanID); err != nil {
		handleBeanError(w, err)
		return
	}

	api.WriteSuccess(w, http.StatusOK, map[string]string{"message": "コーヒー豆を削除しました"})
}

func handleBeanError(w http.ResponseWriter, err error) {
	if errors.Is(err, services.ErrBeanNotFound) {
		api.WriteError(w, http.StatusNotFound, "BEAN_NOT_FOUND", "コーヒー豆が見つかりません")
		return
	}
	if errors.Is(err, services.ErrForbidden) {
		api.WriteError(w, http.StatusForbidden, "FORBIDDEN", "このリソースにアクセスする権限がありません")
		return
	}
	api.WriteError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "サーバーエラーが発生しました")
}

func parseQueryInt(r *http.Request, key string, defaultVal int) int {
	v := r.URL.Query().Get(key)
	if v == "" {
		return defaultVal
	}
	n, err := strconv.Atoi(v)
	if err != nil {
		return defaultVal
	}
	return n
}
