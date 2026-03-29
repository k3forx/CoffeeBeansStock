package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/oapi-codegen/runtime/types"

	"github.com/k3forx/CoffeeBeansStock/backend/internal/api"
	"github.com/k3forx/CoffeeBeansStock/backend/internal/api/gen"
	"github.com/k3forx/CoffeeBeansStock/backend/internal/api/middleware"
	"github.com/k3forx/CoffeeBeansStock/backend/internal/domain/usagehistory"
	"github.com/k3forx/CoffeeBeansStock/backend/internal/usecase"
)

// UsageHistoryHandler handles usage history HTTP requests.
type UsageHistoryHandler struct {
	service *usecase.UsageHistoryService
}

// NewUsageHistoryHandler creates a new UsageHistoryHandler.
func NewUsageHistoryHandler(service *usecase.UsageHistoryService) *UsageHistoryHandler {
	return &UsageHistoryHandler{service: service}
}

// Create handles usage history creation.
func (h *UsageHistoryHandler) Create(w http.ResponseWriter, r *http.Request) {
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

	var req gen.CreateUsageHistoryRequest
	if decodeErr := json.NewDecoder(r.Body).Decode(&req); decodeErr != nil {
		api.WriteError(w, http.StatusBadRequest, "BAD_REQUEST", "リクエストの形式が不正です")
		return
	}

	output, err := h.service.Create(r.Context(), usecase.CreateUsageInput{
		UserID:       userID,
		CoffeeBeanID: beanID,
		UsageDate:    req.UsageDate.Time,
		Quantity:     req.Quantity,
		UsageType:    string(req.UsageType),
		Notes:        req.Notes,
	})
	if err != nil {
		handleDomainError(w, err, "リソースが見つかりません")
		return
	}

	api.WriteSuccess(w, http.StatusCreated, toUsageHistoryResponse(output.Usage))
}

// List handles listing usage history records.
func (h *UsageHistoryHandler) List(w http.ResponseWriter, r *http.Request) {
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

	limit := int32(parseQueryInt(r, "limit", 20))
	offset := int32(parseQueryInt(r, "offset", 0))

	result, err := h.service.ListByCoffeeBean(r.Context(), usecase.ListUsageInput{
		UserID:       userID,
		CoffeeBeanID: beanID,
		Limit:        limit,
		Offset:       offset,
	})
	if err != nil {
		handleDomainError(w, err, "コーヒー豆が見つかりません")
		return
	}

	usages := make([]gen.UsageHistoryResponse, len(result.Usages))
	for i, u := range result.Usages {
		usages[i] = toUsageHistoryResponse(u)
	}

	api.WriteSuccess(w, http.StatusOK, gen.ListUsageHistoryResponse{
		Usages: usages,
		Pagination: gen.PaginationResponse{
			Total:   result.Pagination.Total,
			Limit:   result.Pagination.Limit,
			Offset:  result.Pagination.Offset,
			HasMore: result.Pagination.HasMore,
		},
	})
}

// Delete handles usage history deletion.
func (h *UsageHistoryHandler) Delete(w http.ResponseWriter, r *http.Request) {
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

	usageID, err := uuid.Parse(chi.URLParam(r, "usageId"))
	if err != nil {
		api.WriteError(w, http.StatusBadRequest, "BAD_REQUEST", "不正な使用記録IDです")
		return
	}

	if deleteErr := h.service.Delete(r.Context(), usecase.DeleteUsageInput{
		UserID:       userID,
		CoffeeBeanID: beanID,
		UsageID:      usageID,
	}); deleteErr != nil {
		handleDomainError(w, deleteErr, "使用記録が見つかりません")
		return
	}

	api.WriteSuccess(w, http.StatusOK, gen.DeleteResponse{Message: "使用記録を削除しました"})
}

func toUsageHistoryResponse(u *usagehistory.UsageHistory) gen.UsageHistoryResponse {
	return gen.UsageHistoryResponse{
		Id:           u.ID(),
		CoffeeBeanId: u.CoffeeBeanID(),
		UsageDate:    types.Date{Time: u.UsageDate()},
		Quantity:     u.Quantity().Value(),
		UsageType:    gen.UsageType(u.UsageType().String()),
		Notes:        u.Notes(),
		CreatedAt:    u.CreatedAt(),
	}
}
