package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/k3forx/CoffeeBeansStock/backend/internal/database"
)

type userRepository struct {
	queries *database.Queries
}

// NewUserRepository creates a new UserRepository backed by sqlc queries.
func NewUserRepository(queries *database.Queries) UserRepository {
	return &userRepository{queries: queries}
}

func (r *userRepository) GetByID(ctx context.Context, id uuid.UUID) (database.User, error) {
	return r.queries.GetUserByID(ctx, toUUID(id))
}

func (r *userRepository) CreateAnonymous(ctx context.Context) (database.User, error) {
	return r.queries.CreateAnonymousUser(ctx, database.CreateAnonymousUserParams{
		LowStockThreshold:   100,
		NotificationEnabled: true,
	})
}
