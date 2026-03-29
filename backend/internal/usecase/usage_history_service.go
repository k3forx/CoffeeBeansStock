package usecase

import (
	"context"
	"time"

	"github.com/google/uuid"

	domain "github.com/k3forx/CoffeeBeansStock/backend/internal/domain"
	"github.com/k3forx/CoffeeBeansStock/backend/internal/domain/coffeebean"
	"github.com/k3forx/CoffeeBeansStock/backend/internal/domain/unitofwork"
	"github.com/k3forx/CoffeeBeansStock/backend/internal/domain/usagehistory"
)

// UsageHistoryService handles usage history business logic.
type UsageHistoryService struct {
	usageRepo usagehistory.Repository
	beanRepo  coffeebean.Repository
	uow       unitofwork.UnitOfWork
}

// NewUsageHistoryService creates a new UsageHistoryService.
func NewUsageHistoryService(
	usageRepo usagehistory.Repository,
	beanRepo coffeebean.Repository,
	uow unitofwork.UnitOfWork,
) *UsageHistoryService {
	return &UsageHistoryService{
		usageRepo: usageRepo,
		beanRepo:  beanRepo,
		uow:       uow,
	}
}

// CreateUsageInput holds input for creating a usage history record.
type CreateUsageInput struct {
	UserID       uuid.UUID
	CoffeeBeanID uuid.UUID
	UsageDate    time.Time
	Quantity     int32
	UsageType    string
	Notes        *string
}

// CreateUsageOutput holds the result of creating a usage history record.
type CreateUsageOutput struct {
	Usage *usagehistory.UsageHistory
}

// DeleteUsageInput holds input for deleting a usage history record.
type DeleteUsageInput struct {
	UserID       uuid.UUID
	CoffeeBeanID uuid.UUID
	UsageID      uuid.UUID
}

// ListUsageInput holds input for listing usage history records.
type ListUsageInput struct {
	UserID       uuid.UUID
	CoffeeBeanID uuid.UUID
	Limit        int32
	Offset       int32
}

// ListUsageResult holds a paginated list of usage history records.
type ListUsageResult struct {
	Usages     []*usagehistory.UsageHistory
	Pagination PaginationResponse
}

// Create creates a new usage history record and decreases the coffee bean stock.
func (s *UsageHistoryService) Create(ctx context.Context, in CreateUsageInput) (*CreateUsageOutput, error) {
	qty, err := domain.NewQuantity(in.Quantity)
	if err != nil {
		return nil, err
	}

	usageType, err := usagehistory.NewUsageType(in.UsageType)
	if err != nil {
		return nil, err
	}

	var usage *usagehistory.UsageHistory
	err = s.uow.RunInTx(ctx, func(store unitofwork.Store) error {
		bean, txErr := store.CoffeeBeanRepo().GetByID(ctx, in.CoffeeBeanID)
		if txErr != nil {
			return txErr
		}
		if !bean.IsOwnedBy(in.UserID) {
			return domain.ErrForbidden
		}

		if !bean.CurrentStock().CanConsume(qty) {
			return domain.ErrInsufficientStock
		}

		newStock := bean.CurrentStock().Consume(qty)
		if updateErr := store.CoffeeBeanRepo().UpdateStock(ctx, bean.ID(), newStock); updateErr != nil {
			return updateErr
		}

		usage = usagehistory.New(in.CoffeeBeanID, in.UserID, in.UsageDate, qty, usageType, in.Notes)
		return store.UsageHistoryRepo().Save(ctx, usage)
	})
	if err != nil {
		return nil, err
	}

	return &CreateUsageOutput{Usage: usage}, nil
}

// Delete deletes a usage history record and restores the coffee bean stock.
func (s *UsageHistoryService) Delete(ctx context.Context, in DeleteUsageInput) error {
	usage, err := s.usageRepo.GetByID(ctx, in.UsageID)
	if err != nil {
		return err
	}
	if !usage.IsOwnedBy(in.UserID) {
		return domain.ErrForbidden
	}

	return s.uow.RunInTx(ctx, func(store unitofwork.Store) error {
		bean, txErr := store.CoffeeBeanRepo().GetByID(ctx, in.CoffeeBeanID)
		if txErr != nil {
			return txErr
		}
		if !bean.IsOwnedBy(in.UserID) {
			return domain.ErrForbidden
		}

		newStock := bean.CurrentStock().Add(usage.Quantity())
		if updateErr := store.CoffeeBeanRepo().UpdateStock(ctx, bean.ID(), newStock); updateErr != nil {
			return updateErr
		}

		return store.UsageHistoryRepo().Delete(ctx, in.UsageID, in.UserID)
	})
}

// ListByCoffeeBean returns a paginated list of usage history records for a coffee bean.
func (s *UsageHistoryService) ListByCoffeeBean(ctx context.Context, in ListUsageInput) (*ListUsageResult, error) {
	bean, err := s.beanRepo.GetByID(ctx, in.CoffeeBeanID)
	if err != nil {
		return nil, err
	}
	if !bean.IsOwnedBy(in.UserID) {
		return nil, domain.ErrForbidden
	}

	limit := in.Limit
	offset := in.Offset
	if limit <= 0 || limit > 100 {
		limit = 20
	}
	if offset < 0 {
		offset = 0
	}

	total, err := s.usageRepo.CountByCoffeeBeanID(ctx, in.CoffeeBeanID)
	if err != nil {
		return nil, err
	}

	usages, err := s.usageRepo.ListByCoffeeBeanID(ctx, in.CoffeeBeanID, limit, offset)
	if err != nil {
		return nil, err
	}

	return &ListUsageResult{
		Usages: usages,
		Pagination: PaginationResponse{
			Total:   total,
			Limit:   limit,
			Offset:  offset,
			HasMore: int64(offset+limit) < total,
		},
	}, nil
}
