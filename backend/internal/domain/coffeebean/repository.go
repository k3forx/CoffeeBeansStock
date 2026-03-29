package coffeebean

import (
	"context"

	"github.com/google/uuid"
)

type Repository interface {
	GetByID(ctx context.Context, id uuid.UUID) (*CoffeeBean, error)
	ListByUserID(ctx context.Context, userID uuid.UUID, limit, offset int32) ([]*CoffeeBean, error)
	CountByUserID(ctx context.Context, userID uuid.UUID) (int64, error)
	Save(ctx context.Context, bean *CoffeeBean) error
	Update(ctx context.Context, bean *CoffeeBean) error
	UpdateStock(ctx context.Context, id uuid.UUID, stock Stock) error
	SoftDelete(ctx context.Context, id, userID uuid.UUID) error
}
