package handlers

import (
	"net/http"

	"github.com/oapi-codegen/runtime/types"

	"github.com/k3forx/CoffeeBeansStock/backend/internal/api"
	"github.com/k3forx/CoffeeBeansStock/backend/internal/api/gen"
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
	userID, ok := requireUserID(w, r)
	if !ok {
		return
	}

	beanID, ok := parseUUIDParam(w, r, "id")
	if !ok {
		return
	}

	var req gen.CreateUsageHistoryRequest
	if !decodeJSON(w, r, &req) {
		return
	}

	output, err := h.service.Create(r.Context(), usecase.CreateUsageInput{
		UserID:       userID,
		CoffeeBeanID: beanID,
		UsageDate:    req.UsageDate.Time,
		Quantity:     req.Quantity,
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
	userID, ok := requireUserID(w, r)
	if !ok {
		return
	}

	beanID, ok := parseUUIDParam(w, r, "id")
	if !ok {
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
	userID, ok := requireUserID(w, r)
	if !ok {
		return
	}

	beanID, ok := parseUUIDParam(w, r, "id")
	if !ok {
		return
	}

	usageID, ok := parseUUIDParam(w, r, "usageId")
	if !ok {
		return
	}

	if err := h.service.Delete(r.Context(), usecase.DeleteUsageInput{
		UserID:       userID,
		CoffeeBeanID: beanID,
		UsageID:      usageID,
	}); err != nil {
		handleDomainError(w, err, "使用記録が見つかりません")
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
		Notes:        u.Notes(),
		CreatedAt:    u.CreatedAt(),
	}
}
