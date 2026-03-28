package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/k3forx/CoffeeBeansStock/backend/internal/database"
)

type UserRepository interface {
	GetByID(ctx context.Context, id uuid.UUID) (database.User, error)
	CreateAnonymous(ctx context.Context) (database.User, error)
}
