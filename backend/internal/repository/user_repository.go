package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/k3forx/CoffeeBeansStock/backend/internal/database"
)

type UserRepository interface {
	GetByID(ctx context.Context, id uuid.UUID) (database.User, error)
	GetByEmail(ctx context.Context, email string) (database.User, error)
	Create(ctx context.Context, params CreateUserParams) (database.User, error)
}

type CreateUserParams struct {
	Email               string
	PasswordHash        string
	Name                string
	LowStockThreshold   int32
	NotificationEnabled bool
}
