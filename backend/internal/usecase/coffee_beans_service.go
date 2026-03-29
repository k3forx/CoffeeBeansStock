package usecase

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

// ListBeansInput holds input for listing coffee beans.
type ListBeansInput struct {
	UserID uuid.UUID
	Limit  int32
	Offset int32
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
	UserID       uuid.UUID
	Origin       *string
	RoastDetail  *string
	Notes        *string
	Name         string
	RoastLevel   string
	CurrentStock int32
}

// CreateBeanOutput holds the result of creating a coffee bean.
type CreateBeanOutput struct {
	Bean *coffeebean.CoffeeBean
}

// GetBeanByIDInput holds input for getting a coffee bean by ID.
type GetBeanByIDInput struct {
	UserID uuid.UUID
	BeanID uuid.UUID
}

// GetBeanByIDOutput holds the result of getting a coffee bean by ID.
type GetBeanByIDOutput struct {
	Bean *coffeebean.CoffeeBean
}

// UpdateBeanInput holds input for updating a coffee bean.
type UpdateBeanInput struct {
	UserID       uuid.UUID
	BeanID       uuid.UUID
	Name         *string
	Origin       *string
	RoastLevel   *string
	RoastDetail  *string
	Notes        *string
	CurrentStock *int32
}

// UpdateBeanOutput holds the result of updating a coffee bean.
type UpdateBeanOutput struct {
	Bean *coffeebean.CoffeeBean
}

// DeleteBeanInput holds input for deleting a coffee bean.
type DeleteBeanInput struct {
	UserID uuid.UUID
	BeanID uuid.UUID
}

// List returns a paginated list of coffee beans for the given user.
func (s *CoffeeBeansService) List(ctx context.Context, in ListBeansInput) (*ListBeansResult, error) {
	limit := in.Limit
	offset := in.Offset
	if limit <= 0 || limit > 100 {
		limit = 20
	}
	if offset < 0 {
		offset = 0
	}

	total, err := s.beanRepo.CountByUserID(ctx, in.UserID)
	if err != nil {
		return nil, err
	}

	beans, err := s.beanRepo.ListByUserID(ctx, in.UserID, limit, offset)
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
func (s *CoffeeBeansService) Create(ctx context.Context, in CreateBeanInput) (*CreateBeanOutput, error) {
	roastLevel, err := coffeebean.NewRoastLevel(in.RoastLevel)
	if err != nil {
		return nil, err
	}
	var roastDetail *coffeebean.RoastDetail
	if in.RoastDetail != nil {
		rd, rdErr := coffeebean.NewRoastDetail(*in.RoastDetail)
		if rdErr != nil {
			return nil, rdErr
		}
		roastDetail = &rd
	}
	stock, err := coffeebean.NewStock(in.CurrentStock)
	if err != nil {
		return nil, err
	}

	bean, err := coffeebean.New(in.UserID, in.Name, roastLevel, roastDetail, in.Origin, in.Notes, stock)
	if err != nil {
		return nil, err
	}

	if err := s.beanRepo.Save(ctx, bean); err != nil {
		return nil, err
	}
	return &CreateBeanOutput{Bean: bean}, nil
}

// GetByID returns a coffee bean by ID, verifying ownership.
func (s *CoffeeBeansService) GetByID(ctx context.Context, in GetBeanByIDInput) (*GetBeanByIDOutput, error) {
	bean, err := s.beanRepo.GetByID(ctx, in.BeanID)
	if err != nil {
		return nil, err
	}
	if !bean.IsOwnedBy(in.UserID) {
		return nil, domain.ErrForbidden
	}
	return &GetBeanByIDOutput{Bean: bean}, nil
}

// Update updates a coffee bean, verifying ownership.
func (s *CoffeeBeansService) Update(ctx context.Context, in UpdateBeanInput) (*UpdateBeanOutput, error) {
	bean, err := s.beanRepo.GetByID(ctx, in.BeanID)
	if err != nil {
		return nil, err
	}
	if !bean.IsOwnedBy(in.UserID) {
		return nil, domain.ErrForbidden
	}

	var rl *coffeebean.RoastLevel
	if in.RoastLevel != nil {
		r, err := coffeebean.NewRoastLevel(*in.RoastLevel)
		if err != nil {
			return nil, err
		}
		rl = &r
	}
	var rd *coffeebean.RoastDetail
	if in.RoastDetail != nil {
		d, err := coffeebean.NewRoastDetail(*in.RoastDetail)
		if err != nil {
			return nil, err
		}
		rd = &d
	}
	var st *coffeebean.Stock
	if in.CurrentStock != nil {
		stk, err := coffeebean.NewStock(*in.CurrentStock)
		if err != nil {
			return nil, err
		}
		st = &stk
	}

	if err := bean.Update(in.Name, in.Origin, rl, rd, in.Notes, st); err != nil {
		return nil, err
	}
	if err := s.beanRepo.Update(ctx, bean); err != nil {
		return nil, err
	}
	return &UpdateBeanOutput{Bean: bean}, nil
}

// Delete soft-deletes a coffee bean, verifying ownership.
func (s *CoffeeBeansService) Delete(ctx context.Context, in DeleteBeanInput) error {
	bean, err := s.beanRepo.GetByID(ctx, in.BeanID)
	if err != nil {
		return err
	}
	if !bean.IsOwnedBy(in.UserID) {
		return domain.ErrForbidden
	}
	return s.beanRepo.SoftDelete(ctx, in.BeanID, in.UserID)
}
