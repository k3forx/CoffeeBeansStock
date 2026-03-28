package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/k3forx/CoffeeBeansStock/backend/internal/database"
)

// CoffeeBeanRepository defines the interface for coffee bean data access.
type CoffeeBeanRepository interface {
	GetByID(ctx context.Context, id uuid.UUID) (database.CoffeeBean, error)
	ListByUserID(ctx context.Context, userID uuid.UUID, limit, offset int32) ([]database.CoffeeBean, error)
	CountByUserID(ctx context.Context, userID uuid.UUID) (int64, error)
	Create(ctx context.Context, params CreateCoffeeBeanParams) (database.CoffeeBean, error)
	Update(ctx context.Context, params UpdateCoffeeBeanParams) (database.CoffeeBean, error)
	SoftDelete(ctx context.Context, id, userID uuid.UUID) error
}

// CreateCoffeeBeanParams holds parameters for creating a coffee bean.
type CreateCoffeeBeanParams struct {
	Origin       *string
	RoastLevel   *string
	Notes        *string
	Name         string
	UserID       uuid.UUID
	CurrentStock int32
}

// UpdateCoffeeBeanParams holds parameters for updating a coffee bean.
type UpdateCoffeeBeanParams struct {
	Name         *string
	Origin       *string
	RoastLevel   *string
	CurrentStock *int32
	Notes        *string
	ID           uuid.UUID
	UserID       uuid.UUID
}
