package repository

import (
	"context"
	"time"

	"github.com/google/uuid"

	"github.com/k3forx/CoffeeBeansStock/backend/internal/database"
	domain "github.com/k3forx/CoffeeBeansStock/backend/internal/domain"
	"github.com/k3forx/CoffeeBeansStock/backend/internal/domain/usagehistory"
)

type usageHistoryRepository struct {
	queries *database.Queries
}

// NewUsageHistoryRepository creates a new usagehistory.Repository backed by PostgreSQL.
func NewUsageHistoryRepository(db database.DBTX) usagehistory.Repository {
	return &usageHistoryRepository{queries: database.New(db)}
}

func (r *usageHistoryRepository) GetByID(ctx context.Context, id uuid.UUID) (*usagehistory.UsageHistory, error) {
	row, err := r.queries.GetUsageHistoryByID(ctx, toUUID(id))
	if err != nil {
		return nil, notFoundOrErr(err)
	}
	return toDomainUsageHistory(row), nil
}

func (r *usageHistoryRepository) ListByCoffeeBeanID(ctx context.Context, coffeeBeanID uuid.UUID, limit, offset int32) ([]*usagehistory.UsageHistory, error) {
	rows, err := r.queries.ListUsageHistoriesByCoffeeBeanID(ctx, database.ListUsageHistoriesByCoffeeBeanIDParams{
		CoffeeBeanID: toUUID(coffeeBeanID),
		Limit:        limit,
		Offset:       offset,
	})
	if err != nil {
		return nil, err
	}
	usages := make([]*usagehistory.UsageHistory, len(rows))
	for i, row := range rows {
		usages[i] = toDomainUsageHistory(row)
	}
	return usages, nil
}

func (r *usageHistoryRepository) CountByCoffeeBeanID(ctx context.Context, coffeeBeanID uuid.UUID) (int64, error) {
	return r.queries.CountUsageHistoriesByCoffeeBeanID(ctx, toUUID(coffeeBeanID))
}

func (r *usageHistoryRepository) Save(ctx context.Context, usage *usagehistory.UsageHistory) error {
	_, err := r.queries.CreateUsageHistory(ctx, database.CreateUsageHistoryParams{
		CoffeeBeanID: toUUID(usage.CoffeeBeanID()),
		UserID:       toUUID(usage.UserID()),
		UsageDate:    toPgDate(usage.UsageDate()),
		Quantity:     usage.Quantity().Value(),
		Notes:        toPgText(usage.Notes()),
	})
	return err
}

func (r *usageHistoryRepository) GetRecentUsageSummary(ctx context.Context, coffeeBeanID uuid.UUID, since time.Time) (int32, error) {
	rows, err := r.queries.GetRecentUsageHistoryForConsumptionRate(ctx, database.GetRecentUsageHistoryForConsumptionRateParams{
		CoffeeBeanID: toUUID(coffeeBeanID),
		UsageDate:    toPgDate(since),
	})
	if err != nil {
		return 0, err
	}
	var total int32
	for _, row := range rows {
		total += row.Quantity
	}
	return total, nil
}

func (r *usageHistoryRepository) GetRecentUsageSummaryByUserID(ctx context.Context, userID uuid.UUID, since time.Time) (map[uuid.UUID]int32, error) {
	rows, err := r.queries.GetRecentUsageSummaryByUserID(ctx, database.GetRecentUsageSummaryByUserIDParams{
		UserID:    toUUID(userID),
		UsageDate: toPgDate(since),
	})
	if err != nil {
		return nil, err
	}
	result := make(map[uuid.UUID]int32, len(rows))
	for _, row := range rows {
		if row.CoffeeBeanID.Valid {
			result[uuid.UUID(row.CoffeeBeanID.Bytes)] = row.TotalQuantity
		}
	}
	return result, nil
}

func (r *usageHistoryRepository) Delete(ctx context.Context, id, userID uuid.UUID) error {
	return r.queries.DeleteUsageHistory(ctx, database.DeleteUsageHistoryParams{
		ID:     toUUID(id),
		UserID: toUUID(userID),
	})
}

func toDomainUsageHistory(row database.UsageHistory) *usagehistory.UsageHistory {
	var id, coffeeBeanID, userID uuid.UUID
	if row.ID.Valid {
		id = uuid.UUID(row.ID.Bytes)
	}
	if row.CoffeeBeanID.Valid {
		coffeeBeanID = uuid.UUID(row.CoffeeBeanID.Bytes)
	}
	if row.UserID.Valid {
		userID = uuid.UUID(row.UserID.Bytes)
	}

	var notes *string
	if row.Notes.Valid {
		notes = &row.Notes.String
	}

	return usagehistory.Reconstruct(
		id, coffeeBeanID, userID,
		row.UsageDate.Time,
		domain.ReconstructQuantity(row.Quantity),
		notes,
		row.CreatedAt.Time,
	)
}
