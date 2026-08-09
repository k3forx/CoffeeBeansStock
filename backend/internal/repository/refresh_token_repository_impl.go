package repository

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"

	"github.com/k3forx/CoffeeBeansStock/backend/internal/apperrors"
	"github.com/k3forx/CoffeeBeansStock/backend/internal/database"
	domainauth "github.com/k3forx/CoffeeBeansStock/backend/internal/domain/auth"
)

type refreshTokenRepository struct {
	queries *database.Queries
}

func NewRefreshTokenRepository(db database.DBTX) domainauth.RefreshTokenRepository {
	return &refreshTokenRepository{queries: database.New(db)}
}

func (r *refreshTokenRepository) Store(ctx context.Context, userID uuid.UUID, tokenHash string, expiresAt time.Time) error {
	_, err := r.queries.CreateRefreshToken(ctx, database.CreateRefreshTokenParams{
		UserID:    toUUID(userID),
		TokenHash: tokenHash,
		ExpiresAt: pgtype.Timestamptz{Time: expiresAt, Valid: true},
	})
	return apperrors.Wrap(err)
}

func (r *refreshTokenRepository) ExistsByHash(ctx context.Context, tokenHash string) (bool, error) {
	_, err := r.queries.GetRefreshTokenByHash(ctx, tokenHash)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

func (r *refreshTokenRepository) DeleteByHash(ctx context.Context, tokenHash string) error {
	err := r.queries.DeleteRefreshTokenByHash(ctx, tokenHash)
	return apperrors.Wrap(err)
}

func (r *refreshTokenRepository) DeleteByUserID(ctx context.Context, userID uuid.UUID) error {
	err := r.queries.DeleteRefreshTokensByUserID(ctx, toUUID(userID))
	return apperrors.Wrap(err)
}
