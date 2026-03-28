package user

import (
	"context"

	"github.com/google/uuid"
)

type Repository interface {
	GetByID(ctx context.Context, id uuid.UUID) (*User, error)
	Save(ctx context.Context, user *User) error
}
