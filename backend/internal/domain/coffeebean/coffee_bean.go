package coffeebean

import (
	"time"

	"github.com/google/uuid"

	domain "github.com/k3forx/CoffeeBeansStock/backend/internal/domain"
)

type CoffeeBean struct {
	id           uuid.UUID
	userID       uuid.UUID
	name         string
	origin       *string
	roastLevel   RoastLevel
	currentStock Stock
	notes        *string
	createdAt    time.Time
	updatedAt    time.Time
}

// New creates a new CoffeeBean with validation.
func New(
	userID uuid.UUID, name string, roastLevel RoastLevel,
	origin, notes *string, stock Stock,
) (*CoffeeBean, error) {
	var errs domain.ValidationErrors
	if name == "" {
		errs = append(errs, &domain.ValidationError{Field: "name", Message: "名前は必須です"})
	} else if len(name) > 200 {
		errs = append(errs, &domain.ValidationError{Field: "name", Message: "名前は200文字以内にしてください"})
	}
	if len(errs) > 0 {
		return nil, errs
	}

	now := time.Now()
	return &CoffeeBean{
		id:           uuid.New(),
		userID:       userID,
		name:         name,
		origin:       origin,
		roastLevel:   roastLevel,
		currentStock: stock,
		notes:        notes,
		createdAt:    now,
		updatedAt:    now,
	}, nil
}

// Reconstruct restores a CoffeeBean from persisted data without validation.
func Reconstruct(
	id, userID uuid.UUID, name string, origin *string,
	roastLevel RoastLevel, stock Stock, notes *string,
	createdAt, updatedAt time.Time,
) *CoffeeBean {
	return &CoffeeBean{
		id:           id,
		userID:       userID,
		name:         name,
		origin:       origin,
		roastLevel:   roastLevel,
		currentStock: stock,
		notes:        notes,
		createdAt:    createdAt,
		updatedAt:    updatedAt,
	}
}

func (b *CoffeeBean) ID() uuid.UUID          { return b.id }
func (b *CoffeeBean) UserID() uuid.UUID      { return b.userID }
func (b *CoffeeBean) Name() string           { return b.name }
func (b *CoffeeBean) Origin() *string        { return b.origin }
func (b *CoffeeBean) RoastLevel() RoastLevel { return b.roastLevel }
func (b *CoffeeBean) CurrentStock() Stock    { return b.currentStock }
func (b *CoffeeBean) Notes() *string         { return b.notes }
func (b *CoffeeBean) CreatedAt() time.Time   { return b.createdAt }
func (b *CoffeeBean) UpdatedAt() time.Time   { return b.updatedAt }

// IsOwnedBy checks if the bean belongs to the given user.
func (b *CoffeeBean) IsOwnedBy(userID uuid.UUID) bool {
	return b.userID == userID
}

// Update modifies the bean's master data with validation.
func (b *CoffeeBean) Update(
	name *string, origin *string, roastLevel *RoastLevel,
	notes *string, stock *Stock,
) error {
	var errs domain.ValidationErrors
	if name != nil {
		if *name == "" {
			errs = append(errs, &domain.ValidationError{Field: "name", Message: "名前は必須です"})
		} else if len(*name) > 200 {
			errs = append(errs, &domain.ValidationError{Field: "name", Message: "名前は200文字以内にしてください"})
		}
	}
	if len(errs) > 0 {
		return errs
	}

	if name != nil {
		b.name = *name
	}
	if origin != nil {
		b.origin = origin
	}
	if roastLevel != nil {
		b.roastLevel = *roastLevel
	}
	if stock != nil {
		b.currentStock = *stock
	}
	if notes != nil {
		b.notes = notes
	}
	b.updatedAt = time.Now()
	return nil
}

// ConsumeStock reduces stock by the given quantity.
func (b *CoffeeBean) ConsumeStock(qty domain.Quantity) error {
	if !b.currentStock.CanConsume(qty) {
		return domain.ErrInsufficientStock
	}
	b.currentStock = b.currentStock.Consume(qty)
	b.updatedAt = time.Now()
	return nil
}

// AddStock increases stock by the given quantity.
func (b *CoffeeBean) AddStock(qty domain.Quantity) error {
	b.currentStock = b.currentStock.Add(qty)
	b.updatedAt = time.Now()
	return nil
}
