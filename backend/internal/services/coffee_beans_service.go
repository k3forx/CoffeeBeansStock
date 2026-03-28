package services

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/k3forx/CoffeeBeansStock/backend/internal/api"
	"github.com/k3forx/CoffeeBeansStock/backend/internal/database"
	"github.com/k3forx/CoffeeBeansStock/backend/internal/repository"
)

var (
	ErrBeanNotFound = errors.New("coffee bean not found")
	ErrForbidden    = errors.New("forbidden")
)

type CoffeeBeansService struct {
	beanRepo repository.CoffeeBeanRepository
}

func NewCoffeeBeansService(beanRepo repository.CoffeeBeanRepository) *CoffeeBeansService {
	return &CoffeeBeansService{beanRepo: beanRepo}
}

type CoffeeBeanResponse struct {
	ID           uuid.UUID `json:"id"`
	Name         string    `json:"name"`
	Origin       *string   `json:"origin,omitempty"`
	RoastLevel   *string   `json:"roast_level,omitempty"`
	CurrentStock int32     `json:"current_stock"`
	Notes        *string   `json:"notes,omitempty"`
	CreatedAt    string    `json:"created_at"`
	UpdatedAt    string    `json:"updated_at"`
}

type ListBeansResult struct {
	Beans      []CoffeeBeanResponse `json:"beans"`
	Pagination PaginationResponse   `json:"pagination"`
}

type PaginationResponse struct {
	Total   int64 `json:"total"`
	Limit   int32 `json:"limit"`
	Offset  int32 `json:"offset"`
	HasMore bool  `json:"has_more"`
}

type CreateBeanInput struct {
	Name         string
	Origin       *string
	RoastLevel   *string
	CurrentStock int32
	Notes        *string
}

type UpdateBeanInput struct {
	Name         *string
	Origin       *string
	RoastLevel   *string
	CurrentStock *int32
	Notes        *string
}

func (s *CoffeeBeansService) ValidateCreateInput(input *CreateBeanInput) []api.FieldError {
	var errs []api.FieldError

	if input.Name == "" {
		errs = append(errs, api.FieldError{Field: "name", Message: "名前は必須です"})
	} else if len(input.Name) > 200 {
		errs = append(errs, api.FieldError{Field: "name", Message: "名前は200文字以内で入力してください"})
	}

	if input.CurrentStock < 0 {
		errs = append(errs, api.FieldError{Field: "current_stock", Message: "在庫数は0以上で入力してください"})
	} else if input.CurrentStock > 50000 {
		errs = append(errs, api.FieldError{Field: "current_stock", Message: "在庫数は50000以下で入力してください"})
	}

	return errs
}

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

	result := make([]CoffeeBeanResponse, len(beans))
	for i, b := range beans {
		result[i] = toCoffeeBeanResponse(b)
	}

	return &ListBeansResult{
		Beans: result,
		Pagination: PaginationResponse{
			Total:   total,
			Limit:   limit,
			Offset:  offset,
			HasMore: int64(offset+limit) < total,
		},
	}, nil
}

func (s *CoffeeBeansService) Create(ctx context.Context, userID uuid.UUID, input *CreateBeanInput) (*CoffeeBeanResponse, error) {
	bean, err := s.beanRepo.Create(ctx, repository.CreateCoffeeBeanParams{
		UserID:       userID,
		Name:         input.Name,
		Origin:       input.Origin,
		RoastLevel:   input.RoastLevel,
		CurrentStock: input.CurrentStock,
		Notes:        input.Notes,
	})
	if err != nil {
		return nil, err
	}

	resp := toCoffeeBeanResponse(bean)
	return &resp, nil
}

func (s *CoffeeBeansService) GetByID(ctx context.Context, userID uuid.UUID, beanID uuid.UUID) (*CoffeeBeanResponse, error) {
	bean, err := s.beanRepo.GetByID(ctx, beanID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrBeanNotFound
		}
		return nil, err
	}

	if bean.UserID.Valid && uuid.UUID(bean.UserID.Bytes) != userID {
		return nil, ErrForbidden
	}

	resp := toCoffeeBeanResponse(bean)
	return &resp, nil
}

func (s *CoffeeBeansService) Update(ctx context.Context, userID uuid.UUID, beanID uuid.UUID, input *UpdateBeanInput) (*CoffeeBeanResponse, error) {
	_, err := s.GetByID(ctx, userID, beanID)
	if err != nil {
		return nil, err
	}

	bean, err := s.beanRepo.Update(ctx, repository.UpdateCoffeeBeanParams{
		ID:           beanID,
		UserID:       userID,
		Name:         input.Name,
		Origin:       input.Origin,
		RoastLevel:   input.RoastLevel,
		CurrentStock: input.CurrentStock,
		Notes:        input.Notes,
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrBeanNotFound
		}
		return nil, err
	}

	resp := toCoffeeBeanResponse(bean)
	return &resp, nil
}

func (s *CoffeeBeansService) Delete(ctx context.Context, userID uuid.UUID, beanID uuid.UUID) error {
	_, err := s.GetByID(ctx, userID, beanID)
	if err != nil {
		return err
	}

	return s.beanRepo.SoftDelete(ctx, beanID, userID)
}

func toCoffeeBeanResponse(b database.CoffeeBean) CoffeeBeanResponse {
	resp := CoffeeBeanResponse{
		Name:         b.Name,
		CurrentStock: b.CurrentStock,
	}
	if b.ID.Valid {
		resp.ID = uuid.UUID(b.ID.Bytes)
	}
	if b.Origin.Valid {
		resp.Origin = &b.Origin.String
	}
	if b.RoastLevel.Valid {
		resp.RoastLevel = &b.RoastLevel.String
	}
	if b.Notes.Valid {
		resp.Notes = &b.Notes.String
	}
	if b.CreatedAt.Valid {
		resp.CreatedAt = b.CreatedAt.Time.Format("2006-01-02T15:04:05Z")
	}
	if b.UpdatedAt.Valid {
		resp.UpdatedAt = b.UpdatedAt.Time.Format("2006-01-02T15:04:05Z")
	}
	return resp
}
