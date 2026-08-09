package handlers

import (
	"math"
	"net/http"

	"github.com/k3forx/CoffeeBeansStock/backend/internal/api"
	"github.com/k3forx/CoffeeBeansStock/backend/internal/api/gen"
	"github.com/k3forx/CoffeeBeansStock/backend/internal/domain/coffeebean"
	"github.com/k3forx/CoffeeBeansStock/backend/internal/usecase"
)

// CoffeeBeansHandler handles coffee beans CRUD HTTP requests.
type CoffeeBeansHandler struct {
	service *usecase.CoffeeBeansService
}

// NewCoffeeBeansHandler creates a new CoffeeBeansHandler.
func NewCoffeeBeansHandler(service *usecase.CoffeeBeansService) *CoffeeBeansHandler {
	return &CoffeeBeansHandler{service: service}
}

func toCoffeeBeanResponse(b *coffeebean.CoffeeBean, cr coffeebean.ConsumptionRate) gen.CoffeeBeanResponse {
	resp := gen.CoffeeBeanResponse{
		Id:           b.ID(),
		Name:         b.Name(),
		Origin:       b.Origin(),
		RoastLevel:   gen.RoastLevel(b.RoastLevel().String()),
		CurrentStock: b.CurrentStock().Value(),
		Notes:        b.Notes(),
		CreatedAt:    b.CreatedAt(),
		UpdatedAt:    b.UpdatedAt(),
		ConsumptionRate: gen.ConsumptionRate{
			RemainingCups: cr.RemainingCups(),
		},
	}
	if cr.RemainingDays() != nil {
		resp.ConsumptionRate.RemainingDays = cr.RemainingDays()
	}
	if cr.DailyConsumption() != nil {
		v := float32(*cr.DailyConsumption())
		resp.ConsumptionRate.DailyConsumptionGrams = &v
	}
	if cr.WeeklyTotal() != nil {
		resp.ConsumptionRate.WeeklyTotalGrams = cr.WeeklyTotal()
	}
	if al := cr.AlertLevel(); al != "" {
		s := gen.ConsumptionRateAlertLevel(al)
		resp.ConsumptionRate.AlertLevel = &s
	}
	if b.RoastDetail() != nil {
		rd := gen.RoastDetail(b.RoastDetail().String())
		resp.RoastDetail = &rd
	}
	return resp
}

// List returns a paginated list of coffee beans.
func (h *CoffeeBeansHandler) List(w http.ResponseWriter, r *http.Request) {
	userID, ok := requireUserID(w, r)
	if !ok {
		return
	}

	limit := int32(parseQueryInt(r, "limit", 20, 1, 100))
	offset := int32(parseQueryInt(r, "offset", 0, 0, math.MaxInt32))

	result, err := h.service.List(r.Context(), usecase.ListBeansInput{
		UserID: userID,
		Limit:  limit,
		Offset: offset,
	})
	if err != nil {
		handleDomainError(w, r, err, "")
		return
	}

	beans := make([]gen.CoffeeBeanResponse, len(result.Beans))
	for i, b := range result.Beans {
		cr := result.ConsumptionRates[b.ID()]
		beans[i] = toCoffeeBeanResponse(b, cr)
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
	userID, ok := requireUserID(w, r)
	if !ok {
		return
	}

	var req gen.CreateBeanRequest
	if !decodeJSON(w, r, &req) {
		return
	}

	output, err := h.service.Create(r.Context(), usecase.CreateBeanInput{
		UserID:       userID,
		Name:         req.Name,
		Origin:       req.Origin,
		RoastLevel:   string(req.RoastLevel),
		RoastDetail:  toStringPtr(req.RoastDetail),
		Notes:        req.Notes,
		CurrentStock: req.CurrentStock,
	})
	if err != nil {
		handleDomainError(w, r, err, "コーヒー豆が見つかりません")
		return
	}

	api.WriteSuccess(w, http.StatusCreated, toCoffeeBeanResponse(output.Bean, output.ConsumptionRate))
}

// Get returns a single coffee bean by ID.
func (h *CoffeeBeansHandler) Get(w http.ResponseWriter, r *http.Request) {
	userID, ok := requireUserID(w, r)
	if !ok {
		return
	}

	beanID, ok := parseUUIDParam(w, r, "id")
	if !ok {
		return
	}

	output, err := h.service.GetByID(r.Context(), usecase.GetBeanByIDInput{
		UserID: userID,
		BeanID: beanID,
	})
	if err != nil {
		handleDomainError(w, r, err, "コーヒー豆が見つかりません")
		return
	}

	api.WriteSuccess(w, http.StatusOK, toCoffeeBeanResponse(output.Bean, output.ConsumptionRate))
}

// Update handles coffee bean updates.
func (h *CoffeeBeansHandler) Update(w http.ResponseWriter, r *http.Request) {
	userID, ok := requireUserID(w, r)
	if !ok {
		return
	}

	beanID, ok := parseUUIDParam(w, r, "id")
	if !ok {
		return
	}

	var req gen.UpdateBeanRequest
	if !decodeJSON(w, r, &req) {
		return
	}

	output, err := h.service.Update(r.Context(), usecase.UpdateBeanInput{
		UserID:       userID,
		BeanID:       beanID,
		Name:         req.Name,
		Origin:       req.Origin,
		RoastLevel:   toStringPtr(req.RoastLevel),
		RoastDetail:  toStringPtr(req.RoastDetail),
		CurrentStock: req.CurrentStock,
		Notes:        req.Notes,
	})
	if err != nil {
		handleDomainError(w, r, err, "コーヒー豆が見つかりません")
		return
	}

	api.WriteSuccess(w, http.StatusOK, toCoffeeBeanResponse(output.Bean, output.ConsumptionRate))
}

// Delete handles coffee bean deletion.
func (h *CoffeeBeansHandler) Delete(w http.ResponseWriter, r *http.Request) {
	userID, ok := requireUserID(w, r)
	if !ok {
		return
	}

	beanID, ok := parseUUIDParam(w, r, "id")
	if !ok {
		return
	}

	if err := h.service.Delete(r.Context(), usecase.DeleteBeanInput{
		UserID: userID,
		BeanID: beanID,
	}); err != nil {
		handleDomainError(w, r, err, "コーヒー豆が見つかりません")
		return
	}

	api.WriteSuccess(w, http.StatusOK, gen.DeleteResponse{Message: "コーヒー豆を削除しました"})
}
