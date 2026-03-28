package services

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/k3forx/CoffeeBeansStock/backend/internal/database"
)

var (
	ErrBeanNotFound = errors.New("coffee bean not found")
	ErrForbidden    = errors.New("forbidden")
)

type CoffeeBeansService struct {
	queries *database.Queries
}

func NewCoffeeBeansService(queries *database.Queries) *CoffeeBeansService {
	return &CoffeeBeansService{queries: queries}
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

func (s *CoffeeBeansService) ValidateCreateInput(input *CreateBeanInput) []FieldError {
	var errs []FieldError

	if input.Name == "" {
		errs = append(errs, FieldError{Field: "name", Message: "名前は必須です"})
	} else if len(input.Name) > 200 {
		errs = append(errs, FieldError{Field: "name", Message: "名前は200文字以内で入力してください"})
	}

	if input.CurrentStock < 0 {
		errs = append(errs, FieldError{Field: "current_stock", Message: "在庫数は0以上で入力してください"})
	} else if input.CurrentStock > 50000 {
		errs = append(errs, FieldError{Field: "current_stock", Message: "在庫数は50000以下で入力してください"})
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

	pgUserID := pgtype.UUID{Bytes: userID, Valid: true}

	total, err := s.queries.CountCoffeeBeansByUserID(ctx, pgUserID)
	if err != nil {
		return nil, err
	}

	beans, err := s.queries.ListCoffeeBeansByUserID(ctx, database.ListCoffeeBeansByUserIDParams{
		UserID: pgUserID,
		Limit:  limit,
		Offset: offset,
	})
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
	bean, err := s.queries.CreateCoffeeBean(ctx, database.CreateCoffeeBeanParams{
		UserID:       pgtype.UUID{Bytes: userID, Valid: true},
		Name:         input.Name,
		Origin:       toPgText(input.Origin),
		RoastLevel:   toPgText(input.RoastLevel),
		CurrentStock: input.CurrentStock,
		Notes:        toPgText(input.Notes),
	})
	if err != nil {
		return nil, err
	}

	resp := toCoffeeBeanResponse(bean)
	return &resp, nil
}

func (s *CoffeeBeansService) GetByID(ctx context.Context, userID uuid.UUID, beanID uuid.UUID) (*CoffeeBeanResponse, error) {
	bean, err := s.queries.GetCoffeeBeanByID(ctx, pgtype.UUID{Bytes: beanID, Valid: true})
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
	// Verify ownership first
	_, err := s.GetByID(ctx, userID, beanID)
	if err != nil {
		return nil, err
	}

	bean, err := s.queries.UpdateCoffeeBean(ctx, database.UpdateCoffeeBeanParams{
		ID:           pgtype.UUID{Bytes: beanID, Valid: true},
		UserID:       pgtype.UUID{Bytes: userID, Valid: true},
		Name:         toPgText(input.Name),
		Origin:       toPgText(input.Origin),
		RoastLevel:   toPgText(input.RoastLevel),
		CurrentStock: toPgInt4(input.CurrentStock),
		Notes:        toPgText(input.Notes),
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
	// Verify ownership first
	_, err := s.GetByID(ctx, userID, beanID)
	if err != nil {
		return err
	}

	return s.queries.SoftDeleteCoffeeBean(ctx, database.SoftDeleteCoffeeBeanParams{
		ID:     pgtype.UUID{Bytes: beanID, Valid: true},
		UserID: pgtype.UUID{Bytes: userID, Valid: true},
	})
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

func toPgText(s *string) pgtype.Text {
	if s == nil {
		return pgtype.Text{Valid: false}
	}
	return pgtype.Text{String: *s, Valid: true}
}

func toPgInt4(i *int32) pgtype.Int4 {
	if i == nil {
		return pgtype.Int4{Valid: false}
	}
	return pgtype.Int4{Int32: *i, Valid: true}
}
