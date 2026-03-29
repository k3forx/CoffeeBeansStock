package repository

import (
	"context"

	"github.com/google/uuid"

	"github.com/k3forx/CoffeeBeansStock/backend/internal/database"
	"github.com/k3forx/CoffeeBeansStock/backend/internal/domain/user"
)

type userRepository struct {
	queries *database.Queries
}

func NewUserRepository(db database.DBTX) user.Repository {
	return &userRepository{queries: database.New(db)}
}

func (r *userRepository) GetByID(ctx context.Context, id uuid.UUID) (*user.User, error) {
	u, err := r.queries.GetUserByID(ctx, toUUID(id))
	if err != nil {
		return nil, notFoundOrErr(err)
	}
	return toDomainUser(u), nil
}

func (r *userRepository) Save(ctx context.Context, u *user.User) error {
	_, err := r.queries.CreateAnonymousUser(ctx, database.CreateAnonymousUserParams{
		ID:                  toUUID(u.ID()),
		LowStockThreshold:   u.LowStockThreshold(),
		NotificationEnabled: u.NotificationEnabled(),
	})
	return err
}

func toDomainUser(u database.User) *user.User {
	var id uuid.UUID
	if u.ID.Valid {
		id = uuid.UUID(u.ID.Bytes)
	}

	var email, passwordHash, name string
	if u.Email.Valid {
		email = u.Email.String
	}
	if u.PasswordHash.Valid {
		passwordHash = u.PasswordHash.String
	}
	if u.Name.Valid {
		name = u.Name.String
	}

	return user.Reconstruct(
		id, email, passwordHash, name,
		u.LowStockThreshold, u.NotificationEnabled,
		u.CreatedAt.Time, u.UpdatedAt.Time,
	)
}
