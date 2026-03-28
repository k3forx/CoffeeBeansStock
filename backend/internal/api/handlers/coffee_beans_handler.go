package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"

	"github.com/k3forx/CoffeeBeansStock/backend/internal/api"
	"github.com/k3forx/CoffeeBeansStock/backend/internal/api/gen"
	"github.com/k3forx/CoffeeBeansStock/backend/internal/api/middleware"
	domain "github.com/k3forx/CoffeeBeansStock/backend/internal/domain"
	"github.com/k3forx/CoffeeBeansStock/backend/internal/domain/coffeebean"
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

func toCoffeeBeanResponse(b *coffeebean.CoffeeBean) gen.CoffeeBeanResponse {
	resp := gen.CoffeeBeanResponse{
		Id:           b.ID(),
		Name:         b.Name(),
		Origin:       b.Origin(),
		RoastLevel:   gen.RoastLevel(b.RoastLevel().String()),
		CurrentStock: b.CurrentStock().Value(),
		Notes:        b.Notes(),
		CreatedAt:    b.CreatedAt(),
		UpdatedAt:    b.UpdatedAt(),
	}
	return resp
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

	beans := make([]gen.CoffeeBeanResponse, len(result.Beans))
	for i, b := range result.Beans {
		beans[i] = toCoffeeBeanResponse(b)
	}

	api.WriteSuccess(w, http.StatusOK, gen.ListBeansResponse{
		Beans: beans,
		Pagination: gen.PaginationResponse{
			Total:   result.Pagination.Total,
			Limit:   result.Pagination.Limit,
			Offset:  result.Pagination.Offset,
			HasMore: result.Pagination.HasMore,
		},
	})
}

// Create handles coffee bean creation.
func (h *CoffeeBeansHandler) Create(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserID(r.Context())
	if !ok {
		api.WriteError(w, http.StatusUnauthorized, "UNAUTHORIZED", "認証が必要です")
		return
	}

	var req gen.CreateBeanRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		api.WriteError(w, http.StatusBadRequest, "BAD_REQUEST", "リクエストの形式が不正です")
		return
	}

	bean, err := h.service.Create(r.Context(), userID, &services.CreateBeanInput{
		Name:         req.Name,
		Origin:       req.Origin,
		RoastLevel:   string(req.RoastLevel),
		Notes:        req.Notes,
		CurrentStock: req.CurrentStock,
	})
	if err != nil {
		handleBeanError(w, err)
		return
	}

	api.WriteSuccess(w, http.StatusCreated, toCoffeeBeanResponse(bean))
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

	api.WriteSuccess(w, http.StatusOK, toCoffeeBeanResponse(bean))
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

	var req gen.UpdateBeanRequest
	if err = json.NewDecoder(r.Body).Decode(&req); err != nil {
		api.WriteError(w, http.StatusBadRequest, "BAD_REQUEST", "リクエストの形式が不正です")
		return
	}

	var roastLevel *string
	if req.RoastLevel != nil {
		s := string(*req.RoastLevel)
		roastLevel = &s
	}

	bean, err := h.service.Update(r.Context(), userID, beanID, &services.UpdateBeanInput{
		Name:         req.Name,
		Origin:       req.Origin,
		RoastLevel:   roastLevel,
		CurrentStock: req.CurrentStock,
		Notes:        req.Notes,
	})
	if err != nil {
		handleBeanError(w, err)
		return
	}

	api.WriteSuccess(w, http.StatusOK, toCoffeeBeanResponse(bean))
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

	api.WriteSuccess(w, http.StatusOK, gen.DeleteResponse{Message: "コーヒー豆を削除しました"})
}

func handleBeanError(w http.ResponseWriter, err error) {
	var validationErrs domain.ValidationErrors
	var validationErr *domain.ValidationError
	switch {
	case errors.As(err, &validationErrs):
		details := make([]api.FieldError, len(validationErrs))
		for i, ve := range validationErrs {
			details[i] = api.FieldError{Field: ve.Field, Message: ve.Message}
		}
		api.WriteValidationError(w, details)
	case errors.As(err, &validationErr):
		api.WriteValidationError(w, []api.FieldError{
			{Field: validationErr.Field, Message: validationErr.Message},
		})
	case errors.Is(err, domain.ErrNotFound):
		api.WriteError(w, http.StatusNotFound, "NOT_FOUND", "コーヒー豆が見つかりません")
	case errors.Is(err, domain.ErrForbidden):
		api.WriteError(w, http.StatusForbidden, "FORBIDDEN", "このリソースにアクセスする権限がありません")
	case errors.Is(err, domain.ErrInsufficientStock):
		api.WriteError(w, http.StatusConflict, "INSUFFICIENT_STOCK", "在庫が不足しています")
	default:
		api.WriteError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "サーバーエラーが発生しました")
	}
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
