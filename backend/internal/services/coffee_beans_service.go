package services

import (
	"context"

	"github.com/google/uuid"

	domain "github.com/k3forx/CoffeeBeansStock/backend/internal/domain"
	"github.com/k3forx/CoffeeBeansStock/backend/internal/domain/coffeebean"
	"github.com/k3forx/CoffeeBeansStock/backend/internal/domain/unitofwork"
)

// CoffeeBeansService handles coffee bean business logic.
type CoffeeBeansService struct {
	beanRepo coffeebean.Repository
	uow      unitofwork.UnitOfWork
}

// NewCoffeeBeansService creates a new CoffeeBeansService.
func NewCoffeeBeansService(beanRepo coffeebean.Repository, uow unitofwork.UnitOfWork) *CoffeeBeansService {
	return &CoffeeBeansService{beanRepo: beanRepo, uow: uow}
}

// ListBeansResult holds a paginated list of coffee beans.
type ListBeansResult struct {
	Beans      []*coffeebean.CoffeeBean
	Pagination PaginationResponse
}

// PaginationResponse holds pagination metadata.
type PaginationResponse struct {
	Total   int64 `json:"total"`
	Limit   int32 `json:"limit"`
	Offset  int32 `json:"offset"`
	HasMore bool  `json:"has_more"`
}

// CreateBeanInput holds input for creating a coffee bean.
type CreateBeanInput struct {
	Name         string
	Origin       *string
	RoastLevel   string
	Notes        *string
	CurrentStock int32
}

// UpdateBeanInput holds input for updating a coffee bean.
type UpdateBeanInput struct {
	Name         *string
	Origin       *string
	RoastLevel   *string
	Notes        *string
	CurrentStock *int32
}

// List returns a paginated list of coffee beans for the given user.
func (s *CoffeeBeansService) List(ctx context.Context, userID uuid.UUID, limit, offset int32) (*ListBeansResult, error) {
	if limit <= 0 || limit > 100 {
		limit = 20
	}
	if offset < 0 {
		offset = 0
	}

	total, err := s.beanRepo.CountByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	beans, err := s.beanRepo.ListByUserID(ctx, userID, limit, offset)
	if err != nil {
		return nil, err
	}

	return &ListBeansResult{
		Beans: beans,
		Pagination: PaginationResponse{
			Total:   total,
			Limit:   limit,
			Offset:  offset,
			HasMore: int64(offset+limit) < total,
		},
	}, nil
}

// Create creates a new coffee bean for the given user.
func (s *CoffeeBeansService) Create(ctx context.Context, userID uuid.UUID, input *CreateBeanInput) (*coffeebean.CoffeeBean, error) {
	roastLevel, err := coffeebean.NewRoastLevel(input.RoastLevel)
	if err != nil {
		return nil, err
	}
	stock, err := coffeebean.NewStock(input.CurrentStock)
	if err != nil {
		return nil, err
	}

	bean, err := coffeebean.New(userID, input.Name, roastLevel, input.Origin, input.Notes, stock)
	if err != nil {
		return nil, err
	}

	if err := s.beanRepo.Save(ctx, bean); err != nil {
		return nil, err
	}
	return bean, nil
}

// GetByID returns a coffee bean by ID, verifying ownership.
func (s *CoffeeBeansService) GetByID(ctx context.Context, userID, beanID uuid.UUID) (*coffeebean.CoffeeBean, error) {
	bean, err := s.beanRepo.GetByID(ctx, beanID)
	if err != nil {
		return nil, err
	}
	if !bean.IsOwnedBy(userID) {
		return nil, domain.ErrForbidden
	}
	return bean, nil
}

// Update updates a coffee bean, verifying ownership.
func (s *CoffeeBeansService) Update(ctx context.Context, userID, beanID uuid.UUID, input *UpdateBeanInput) (*coffeebean.CoffeeBean, error) {
	bean, err := s.beanRepo.GetByID(ctx, beanID)
	if err != nil {
		return nil, err
	}
	if !bean.IsOwnedBy(userID) {
		return nil, domain.ErrForbidden
	}

	var rl *coffeebean.RoastLevel
	if input.RoastLevel != nil {
		r, err := coffeebean.NewRoastLevel(*input.RoastLevel)
		if err != nil {
			return nil, err
		}
		rl = &r
	}
	var st *coffeebean.Stock
	if input.CurrentStock != nil {
		stk, err := coffeebean.NewStock(*input.CurrentStock)
		if err != nil {
			return nil, err
		}
		st = &stk
	}

	if err := bean.Update(input.Name, input.Origin, rl, input.Notes, st); err != nil {
		return nil, err
	}
	if err := s.beanRepo.Update(ctx, bean); err != nil {
		return nil, err
	}
	return bean, nil
}

// Delete soft-deletes a coffee bean, verifying ownership.
func (s *CoffeeBeansService) Delete(ctx context.Context, userID, beanID uuid.UUID) error {
	bean, err := s.beanRepo.GetByID(ctx, beanID)
	if err != nil {
		return err
	}
	if !bean.IsOwnedBy(userID) {
		return domain.ErrForbidden
	}
	return s.beanRepo.SoftDelete(ctx, beanID, userID)
}
