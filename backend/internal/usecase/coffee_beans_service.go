package usecase

import (
	"context"
	"time"

	"github.com/google/uuid"

	domain "github.com/k3forx/CoffeeBeansStock/backend/internal/domain"
	"github.com/k3forx/CoffeeBeansStock/backend/internal/domain/coffeebean"
	"github.com/k3forx/CoffeeBeansStock/backend/internal/domain/unitofwork"
	"github.com/k3forx/CoffeeBeansStock/backend/internal/domain/usagehistory"
	"github.com/k3forx/CoffeeBeansStock/backend/internal/domain/user"
)

// CoffeeBeansService handles coffee bean business logic.
type CoffeeBeansService struct {
	beanRepo  coffeebean.Repository
	usageRepo usagehistory.Repository
	userRepo  user.Repository
	uow       unitofwork.UnitOfWork
}

// NewCoffeeBeansService creates a new CoffeeBeansService.
func NewCoffeeBeansService(beanRepo coffeebean.Repository, usageRepo usagehistory.Repository, userRepo user.Repository, uow unitofwork.UnitOfWork) *CoffeeBeansService {
	return &CoffeeBeansService{beanRepo: beanRepo, usageRepo: usageRepo, userRepo: userRepo, uow: uow}
}

// ListBeansInput holds input for listing coffee beans.
type ListBeansInput struct {
	UserID uuid.UUID
	Limit  int32
	Offset int32
}

// ListBeansResult holds a paginated list of coffee beans.
type ListBeansResult struct {
	Beans            []*coffeebean.CoffeeBean
	ConsumptionRates map[uuid.UUID]coffeebean.ConsumptionRate
	Pagination       PaginationResponse
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
	Bean            *coffeebean.CoffeeBean
	ConsumptionRate coffeebean.ConsumptionRate
}

// GetBeanByIDInput holds input for getting a coffee bean by ID.
type GetBeanByIDInput struct {
	UserID uuid.UUID
	BeanID uuid.UUID
}

// GetBeanByIDOutput holds the result of getting a coffee bean by ID.
type GetBeanByIDOutput struct {
	Bean            *coffeebean.CoffeeBean
	ConsumptionRate coffeebean.ConsumptionRate
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
	Bean            *coffeebean.CoffeeBean
	ConsumptionRate coffeebean.ConsumptionRate
}

// DeleteBeanInput holds input for deleting a coffee bean.
type DeleteBeanInput struct {
	UserID uuid.UUID
	BeanID uuid.UUID
}

// List returns a paginated list of coffee beans for the given user.
func (s *CoffeeBeansService) List(ctx context.Context, in ListBeansInput) (*ListBeansResult, error) {
	limit, offset := normalizePagination(in.Limit, in.Offset)

	total, err := s.beanRepo.CountByUserID(ctx, in.UserID)
	if err != nil {
		return nil, err
	}

	beans, err := s.beanRepo.ListByUserID(ctx, in.UserID, limit, offset)
	if err != nil {
		return nil, err
	}

	// Fetch user for grams_per_cup
	u, err := s.userRepo.GetByID(ctx, in.UserID)
	if err != nil {
		return nil, err
	}
	gramsPerCup := u.GramsPerCup().Value()

	// Fetch usage summary in batch
	since := time.Now().AddDate(0, 0, -7)
	usageSummary, err := s.usageRepo.GetRecentUsageSummaryByUserID(ctx, in.UserID, since)
	if err != nil {
		return nil, err
	}

	// Build consumption rates map
	consumptionRates := make(map[uuid.UUID]coffeebean.ConsumptionRate, len(beans))
	for _, b := range beans {
		var weeklyUsage *int32
		if t, ok := usageSummary[b.ID()]; ok && t > 0 {
			weeklyUsage = &t
		}
		consumptionRates[b.ID()] = coffeebean.NewConsumptionRate(b.CurrentStock().Value(), gramsPerCup, weeklyUsage)
	}

	return &ListBeansResult{
		Beans:            beans,
		ConsumptionRates: consumptionRates,
		Pagination:       newPaginationResponse(total, limit, offset),
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

	// New bean has no usage history, so weeklyUsage is nil
	u, err := s.userRepo.GetByID(ctx, in.UserID)
	if err != nil {
		return nil, err
	}
	gramsPerCup := u.GramsPerCup().Value()
	cr := coffeebean.NewConsumptionRate(bean.CurrentStock().Value(), gramsPerCup, nil)

	return &CreateBeanOutput{Bean: bean, ConsumptionRate: cr}, nil
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

	cr, err := s.computeConsumptionRate(ctx, in.BeanID, in.UserID, bean.CurrentStock().Value())
	if err != nil {
		return nil, err
	}

	return &GetBeanByIDOutput{Bean: bean, ConsumptionRate: cr}, nil
}

// Update updates a coffee bean, verifying ownership.
func (s *CoffeeBeansService) Update(ctx context.Context, in UpdateBeanInput) (*UpdateBeanOutput, error) {
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

	var result *UpdateBeanOutput
	err := s.uow.RunInTx(ctx, func(store unitofwork.Store) error {
		bean, err := store.CoffeeBeanRepo().GetByIDForUpdate(ctx, in.BeanID)
		if err != nil {
			return err
		}
		if !bean.IsOwnedBy(in.UserID) {
			return domain.ErrForbidden
		}
		if err := bean.Update(in.Name, in.Origin, rl, rd, in.Notes, st); err != nil {
			return err
		}
		if err := store.CoffeeBeanRepo().Update(ctx, bean); err != nil {
			return err
		}
		result = &UpdateBeanOutput{Bean: bean}
		return nil
	})
	if err != nil {
		return nil, err
	}

	// Compute consumption rate after the transaction
	cr, err := s.computeConsumptionRate(ctx, in.BeanID, in.UserID, result.Bean.CurrentStock().Value())
	if err != nil {
		return nil, err
	}
	result.ConsumptionRate = cr

	return result, nil
}

// Delete soft-deletes a coffee bean, verifying ownership.
func (s *CoffeeBeansService) Delete(ctx context.Context, in DeleteBeanInput) error {
	return s.uow.RunInTx(ctx, func(store unitofwork.Store) error {
		bean, err := store.CoffeeBeanRepo().GetByIDForUpdate(ctx, in.BeanID)
		if err != nil {
			return err
		}
		if !bean.IsOwnedBy(in.UserID) {
			return domain.ErrForbidden
		}
		return store.CoffeeBeanRepo().SoftDelete(ctx, in.BeanID, in.UserID)
	})
}

func (s *CoffeeBeansService) computeConsumptionRate(ctx context.Context, beanID, userID uuid.UUID, currentStock int32) (coffeebean.ConsumptionRate, error) {
	u, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return coffeebean.ConsumptionRate{}, err
	}
	gramsPerCup := u.GramsPerCup().Value()

	since := time.Now().AddDate(0, 0, -7)
	total, err := s.usageRepo.GetRecentUsageSummary(ctx, beanID, since)
	if err != nil {
		return coffeebean.ConsumptionRate{}, err
	}

	var weeklyUsage *int32
	if total > 0 {
		weeklyUsage = &total
	}

	return coffeebean.NewConsumptionRate(currentStock, gramsPerCup, weeklyUsage), nil
}
