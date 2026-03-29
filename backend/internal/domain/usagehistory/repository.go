package usagehistory

import (
	"context"

	"github.com/google/uuid"
)

// Repository defines the interface for usage history persistence.
type Repository interface {
	GetByID(ctx context.Context, id uuid.UUID) (*UsageHistory, error)
	ListByCoffeeBeanID(ctx context.Context, coffeeBeanID uuid.UUID, limit, offset int32) ([]*UsageHistory, error)
	CountByCoffeeBeanID(ctx context.Context, coffeeBeanID uuid.UUID) (int64, error)
	Save(ctx context.Context, usage *UsageHistory) error
	Delete(ctx context.Context, id, userID uuid.UUID) error
}
