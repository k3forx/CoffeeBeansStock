package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/k3forx/CoffeeBeansStock/backend/internal/database"
)

// UserRepository defines the interface for user data access.
type UserRepository interface {
	GetByID(ctx context.Context, id uuid.UUID) (database.User, error)
	CreateAnonymous(ctx context.Context) (database.User, error)
}
