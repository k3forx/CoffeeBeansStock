package usagehistory

import (
	"time"

	"github.com/google/uuid"

	domain "github.com/k3forx/CoffeeBeansStock/backend/internal/domain"
)

// UsageHistory represents a single usage record of a coffee bean.
type UsageHistory struct {
	id           uuid.UUID
	coffeeBeanID uuid.UUID
	userID       uuid.UUID
	usageDate    time.Time
	quantity     domain.Quantity
	usageType    UsageType
	notes        *string
	createdAt    time.Time
}

// New creates a new UsageHistory with a generated ID and current timestamp.
func New(
	coffeeBeanID, userID uuid.UUID,
	usageDate time.Time,
	quantity domain.Quantity,
	usageType UsageType,
	notes *string,
) *UsageHistory {
	return &UsageHistory{
		id:           uuid.New(),
		coffeeBeanID: coffeeBeanID,
		userID:       userID,
		usageDate:    usageDate,
		quantity:     quantity,
		usageType:    usageType,
		notes:        notes,
		createdAt:    time.Now().UTC(),
	}
}

// Reconstruct restores a UsageHistory from persisted data without validation.
func Reconstruct(
	id, coffeeBeanID, userID uuid.UUID,
	usageDate time.Time,
	quantity domain.Quantity,
	usageType UsageType,
	notes *string,
	createdAt time.Time,
) *UsageHistory {
	return &UsageHistory{
		id:           id,
		coffeeBeanID: coffeeBeanID,
		userID:       userID,
		usageDate:    usageDate,
		quantity:     quantity,
		usageType:    usageType,
		notes:        notes,
		createdAt:    createdAt,
	}
}

func (h *UsageHistory) ID() uuid.UUID             { return h.id }
func (h *UsageHistory) CoffeeBeanID() uuid.UUID   { return h.coffeeBeanID }
func (h *UsageHistory) UserID() uuid.UUID         { return h.userID }
func (h *UsageHistory) UsageDate() time.Time      { return h.usageDate }
func (h *UsageHistory) Quantity() domain.Quantity { return h.quantity }
func (h *UsageHistory) UsageType() UsageType      { return h.usageType }
func (h *UsageHistory) Notes() *string            { return h.notes }
func (h *UsageHistory) CreatedAt() time.Time      { return h.createdAt }

// IsOwnedBy checks if the usage history belongs to the given user.
func (h *UsageHistory) IsOwnedBy(userID uuid.UUID) bool {
	return h.userID == userID
}
