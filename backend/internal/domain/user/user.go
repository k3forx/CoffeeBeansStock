package user

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	createdAt           time.Time
	updatedAt           time.Time
	email               string
	passwordHash        string
	name                string
	lowStockThreshold   int32
	id                  uuid.UUID
	notificationEnabled bool
}

// NewAnonymousUser creates a new anonymous user with default settings.
func NewAnonymousUser() *User {
	now := time.Now()
	return &User{
		id:                  uuid.New(),
		email:               "",
		passwordHash:        "",
		name:                "",
		lowStockThreshold:   100,
		notificationEnabled: true,
		createdAt:           now,
		updatedAt:           now,
	}
}

// Reconstruct restores a User from persisted data without validation.
func Reconstruct(
	id uuid.UUID, email, passwordHash, name string,
	lowStockThreshold int32, notificationEnabled bool,
	createdAt, updatedAt time.Time,
) *User {
	return &User{
		id:                  id,
		email:               email,
		passwordHash:        passwordHash,
		name:                name,
		lowStockThreshold:   lowStockThreshold,
		notificationEnabled: notificationEnabled,
		createdAt:           createdAt,
		updatedAt:           updatedAt,
	}
}

func (u *User) ID() uuid.UUID             { return u.id }
func (u *User) Email() string             { return u.email }
func (u *User) PasswordHash() string      { return u.passwordHash }
func (u *User) Name() string              { return u.name }
func (u *User) LowStockThreshold() int32  { return u.lowStockThreshold }
func (u *User) NotificationEnabled() bool { return u.notificationEnabled }
func (u *User) CreatedAt() time.Time      { return u.createdAt }
func (u *User) UpdatedAt() time.Time      { return u.updatedAt }
