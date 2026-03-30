package auth

import (
	"context"
	"time"

	"github.com/google/uuid"
)

// RefreshTokenRepository manages persisted refresh token hashes.
type RefreshTokenRepository interface {
	Store(ctx context.Context, userID uuid.UUID, tokenHash string, expiresAt time.Time) error
	ExistsByHash(ctx context.Context, tokenHash string) (bool, error)
	DeleteByHash(ctx context.Context, tokenHash string) error
	DeleteByUserID(ctx context.Context, userID uuid.UUID) error
}
