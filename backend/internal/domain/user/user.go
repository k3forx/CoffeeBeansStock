package user

import (
	"errors"
	"time"

	"github.com/google/uuid"

	domain "github.com/k3forx/CoffeeBeansStock/backend/internal/domain"
)

type User struct {
	createdAt         time.Time
	updatedAt         time.Time
	email             string
	passwordHash      string
	name              string
	lowStockThreshold int32
	id                uuid.UUID
	gramsPerCup       GramsPerCup
}

// NewAnonymousUser creates a new anonymous user with default settings.
func NewAnonymousUser() *User {
	now := time.Now().UTC()
	return &User{
		id:                uuid.New(),
		email:             "",
		passwordHash:      "",
		name:              "",
		lowStockThreshold: 100,
		gramsPerCup:       DefaultGramsPerCup(),
		createdAt:         now,
		updatedAt:         now,
	}
}

// Reconstruct restores a User from persisted data without validation.
func Reconstruct(
	id uuid.UUID, email, passwordHash, name string,
	lowStockThreshold int32,
	gramsPerCup GramsPerCup,
	createdAt, updatedAt time.Time,
) *User {
	return &User{
		id:                id,
		email:             email,
		passwordHash:      passwordHash,
		name:              name,
		lowStockThreshold: lowStockThreshold,
		gramsPerCup:       gramsPerCup,
		createdAt:         createdAt,
		updatedAt:         updatedAt,
	}
}

// Update modifies the user's settings with validation.
func (u *User) Update(gramsPerCup, lowStockThreshold *int32) error {
	var errs domain.ValidationErrors

	if gramsPerCup != nil {
		gpc, err := NewGramsPerCup(*gramsPerCup)
		if err == nil {
			u.gramsPerCup = gpc
		} else if ve, ok := errors.AsType[*domain.ValidationError](err); ok {
			errs = append(errs, ve)
		} else {
			return err
		}
	}
	if lowStockThreshold != nil {
		u.lowStockThreshold = *lowStockThreshold
	}
	if len(errs) > 0 {
		return errs
	}
	u.updatedAt = time.Now().UTC()
	return nil
}

func (u *User) ID() uuid.UUID            { return u.id }
func (u *User) Email() string            { return u.email }
func (u *User) PasswordHash() string     { return u.passwordHash }
func (u *User) Name() string             { return u.name }
func (u *User) LowStockThreshold() int32 { return u.lowStockThreshold }
func (u *User) GramsPerCup() GramsPerCup { return u.gramsPerCup }
func (u *User) CreatedAt() time.Time     { return u.createdAt }
func (u *User) UpdatedAt() time.Time     { return u.updatedAt }
