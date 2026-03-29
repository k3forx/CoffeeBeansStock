package repository

import (
	"context"

	"github.com/google/uuid"

	"github.com/k3forx/CoffeeBeansStock/backend/internal/database"
	"github.com/k3forx/CoffeeBeansStock/backend/internal/domain/coffeebean"
)

type coffeeBeanRepository struct {
	queries *database.Queries
}

func NewCoffeeBeanRepository(db database.DBTX) coffeebean.Repository {
	return &coffeeBeanRepository{queries: database.New(db)}
}

func (r *coffeeBeanRepository) GetByID(ctx context.Context, id uuid.UUID) (*coffeebean.CoffeeBean, error) {
	b, err := r.queries.GetCoffeeBeanByID(ctx, toUUID(id))
	if err != nil {
		return nil, notFoundOrErr(err)
	}
	return toDomainCoffeeBean(b), nil
}

func (r *coffeeBeanRepository) GetByIDForUpdate(ctx context.Context, id uuid.UUID) (*coffeebean.CoffeeBean, error) {
	b, err := r.queries.GetCoffeeBeanByIDForUpdate(ctx, toUUID(id))
	if err != nil {
		return nil, notFoundOrErr(err)
	}
	return toDomainCoffeeBean(b), nil
}

func (r *coffeeBeanRepository) ListByUserID(ctx context.Context, userID uuid.UUID, limit, offset int32) ([]*coffeebean.CoffeeBean, error) {
	rows, err := r.queries.ListCoffeeBeansByUserID(ctx, database.ListCoffeeBeansByUserIDParams{
		UserID: toUUID(userID),
		Limit:  limit,
		Offset: offset,
	})
	if err != nil {
		return nil, err
	}
	beans := make([]*coffeebean.CoffeeBean, len(rows))
	for i, row := range rows {
		beans[i] = toDomainCoffeeBean(row)
	}
	return beans, nil
}

func (r *coffeeBeanRepository) CountByUserID(ctx context.Context, userID uuid.UUID) (int64, error) {
	return r.queries.CountCoffeeBeansByUserID(ctx, toUUID(userID))
}

func (r *coffeeBeanRepository) Save(ctx context.Context, bean *coffeebean.CoffeeBean) error {
	params := database.CreateCoffeeBeanParams{
		ID:           toUUID(bean.ID()),
		UserID:       toUUID(bean.UserID()),
		Name:         bean.Name(),
		Origin:       toPgText(bean.Origin()),
		RoastLevel:   bean.RoastLevel().String(),
		CurrentStock: bean.CurrentStock().Value(),
		Notes:        toPgText(bean.Notes()),
	}
	if bean.RoastDetail() != nil {
		s := bean.RoastDetail().String()
		params.RoastDetail = toPgText(&s)
	}
	_, err := r.queries.CreateCoffeeBean(ctx, params)
	return err
}

func (r *coffeeBeanRepository) Update(ctx context.Context, bean *coffeebean.CoffeeBean) error {
	name := bean.Name()
	roastLevel := bean.RoastLevel().String()
	stock := bean.CurrentStock().Value()
	params := database.UpdateCoffeeBeanParams{
		ID:           toUUID(bean.ID()),
		UserID:       toUUID(bean.UserID()),
		Name:         toPgText(&name),
		Origin:       toPgText(bean.Origin()),
		RoastLevel:   toPgText(&roastLevel),
		CurrentStock: toPgInt4(&stock),
		Notes:        toPgText(bean.Notes()),
	}
	if bean.RoastDetail() != nil {
		s := bean.RoastDetail().String()
		params.RoastDetail = toPgText(&s)
	}
	_, err := r.queries.UpdateCoffeeBean(ctx, params)
	return notFoundOrErr(err)
}

func (r *coffeeBeanRepository) UpdateStock(ctx context.Context, id uuid.UUID, stock coffeebean.Stock) error {
	_, err := r.queries.UpdateCoffeeBeanStock(ctx, database.UpdateCoffeeBeanStockParams{
		ID:           toUUID(id),
		CurrentStock: stock.Value(),
	})
	return err
}

func (r *coffeeBeanRepository) SoftDelete(ctx context.Context, id, userID uuid.UUID) error {
	return r.queries.SoftDeleteCoffeeBean(ctx, database.SoftDeleteCoffeeBeanParams{
		ID:     toUUID(id),
		UserID: toUUID(userID),
	})
}

func toDomainCoffeeBean(b database.CoffeeBean) *coffeebean.CoffeeBean {
	var id, userID uuid.UUID
	if b.ID.Valid {
		id = uuid.UUID(b.ID.Bytes)
	}
	if b.UserID.Valid {
		userID = uuid.UUID(b.UserID.Bytes)
	}

	var origin *string
	if b.Origin.Valid {
		origin = &b.Origin.String
	}
	var notes *string
	if b.Notes.Valid {
		notes = &b.Notes.String
	}

	var roastDetail *coffeebean.RoastDetail
	if b.RoastDetail.Valid {
		rd := coffeebean.ReconstructRoastDetail(b.RoastDetail.String)
		roastDetail = &rd
	}

	return coffeebean.Reconstruct(
		id, userID, b.Name, origin,
		coffeebean.ReconstructRoastLevel(b.RoastLevel),
		roastDetail,
		coffeebean.ReconstructStock(b.CurrentStock),
		notes, b.CreatedAt.Time, b.UpdatedAt.Time,
	)
}
