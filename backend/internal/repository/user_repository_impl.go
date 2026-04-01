package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"

	"github.com/k3forx/CoffeeBeansStock/backend/internal/apperrors"
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
		ID:                toUUID(u.ID()),
		LowStockThreshold: u.LowStockThreshold(),
		GramsPerCup:       u.GramsPerCup().Value(),
	})
	return apperrors.Wrap(err)
}

func (r *userRepository) Update(ctx context.Context, u *user.User) error {
	threshold := u.LowStockThreshold()
	gpc := u.GramsPerCup().Value()

	_, err := r.queries.UpdateUser(ctx, database.UpdateUserParams{
		ID:                toUUID(u.ID()),
		LowStockThreshold: pgtype.Int4{Int32: threshold, Valid: true},
		GramsPerCup:       pgtype.Int4{Int32: gpc, Valid: true},
	})
	return apperrors.Wrap(err)
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
		u.LowStockThreshold,
		user.ReconstructGramsPerCup(u.GramsPerCup),
		u.CreatedAt.Time, u.UpdatedAt.Time,
	)
}
