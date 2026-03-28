package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/k3forx/CoffeeBeansStock/backend/internal/database"
)

type userRepository struct {
	queries *database.Queries
}

func NewUserRepository(queries *database.Queries) UserRepository {
	return &userRepository{queries: queries}
}

func (r *userRepository) GetByID(ctx context.Context, id uuid.UUID) (database.User, error) {
	return r.queries.GetUserByID(ctx, toUUID(id))
}

func (r *userRepository) GetByEmail(ctx context.Context, email string) (database.User, error) {
	return r.queries.GetUserByEmail(ctx, email)
}

func (r *userRepository) Create(ctx context.Context, params CreateUserParams) (database.User, error) {
	return r.queries.CreateUser(ctx, database.CreateUserParams{
		Email:               params.Email,
		PasswordHash:        params.PasswordHash,
		Name:                params.Name,
		LowStockThreshold:   params.LowStockThreshold,
		NotificationEnabled: params.NotificationEnabled,
	})
}
