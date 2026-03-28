package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/k3forx/CoffeeBeansStock/backend/internal/database"
)

type coffeeBeanRepository struct {
	queries *database.Queries
}

// NewCoffeeBeanRepository creates a new CoffeeBeanRepository backed by sqlc queries.
func NewCoffeeBeanRepository(queries *database.Queries) CoffeeBeanRepository {
	return &coffeeBeanRepository{queries: queries}
}

func (r *coffeeBeanRepository) GetByID(ctx context.Context, id uuid.UUID) (database.CoffeeBean, error) {
	return r.queries.GetCoffeeBeanByID(ctx, toUUID(id))
}

func (r *coffeeBeanRepository) ListByUserID(ctx context.Context, userID uuid.UUID, limit, offset int32) ([]database.CoffeeBean, error) {
	return r.queries.ListCoffeeBeansByUserID(ctx, database.ListCoffeeBeansByUserIDParams{
		UserID: toUUID(userID),
		Limit:  limit,
		Offset: offset,
	})
}

func (r *coffeeBeanRepository) CountByUserID(ctx context.Context, userID uuid.UUID) (int64, error) {
	return r.queries.CountCoffeeBeansByUserID(ctx, toUUID(userID))
}

func (r *coffeeBeanRepository) Create(ctx context.Context, params CreateCoffeeBeanParams) (database.CoffeeBean, error) {
	return r.queries.CreateCoffeeBean(ctx, database.CreateCoffeeBeanParams{
		UserID:       toUUID(params.UserID),
		Name:         params.Name,
		Origin:       toPgText(params.Origin),
		RoastLevel:   toPgText(params.RoastLevel),
		CurrentStock: params.CurrentStock,
		Notes:        toPgText(params.Notes),
	})
}

func (r *coffeeBeanRepository) Update(ctx context.Context, params UpdateCoffeeBeanParams) (database.CoffeeBean, error) {
	return r.queries.UpdateCoffeeBean(ctx, database.UpdateCoffeeBeanParams{
		ID:           toUUID(params.ID),
		UserID:       toUUID(params.UserID),
		Name:         toPgText(params.Name),
		Origin:       toPgText(params.Origin),
		RoastLevel:   toPgText(params.RoastLevel),
		CurrentStock: toPgInt4(params.CurrentStock),
		Notes:        toPgText(params.Notes),
	})
}

func (r *coffeeBeanRepository) SoftDelete(ctx context.Context, id, userID uuid.UUID) error {
	return r.queries.SoftDeleteCoffeeBean(ctx, database.SoftDeleteCoffeeBeanParams{
		ID:     toUUID(id),
		UserID: toUUID(userID),
	})
}
