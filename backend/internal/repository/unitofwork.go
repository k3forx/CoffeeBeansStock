package repository

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/k3forx/CoffeeBeansStock/backend/internal/domain/coffeebean"
	"github.com/k3forx/CoffeeBeansStock/backend/internal/domain/unitofwork"
	"github.com/k3forx/CoffeeBeansStock/backend/internal/domain/usagehistory"
	"github.com/k3forx/CoffeeBeansStock/backend/internal/domain/user"
)

type unitOfWorkImpl struct {
	pool *pgxpool.Pool
}

func NewUnitOfWork(pool *pgxpool.Pool) unitofwork.UnitOfWork {
	return &unitOfWorkImpl{pool: pool}
}

func (u *unitOfWorkImpl) RunInTx(ctx context.Context, fn func(store unitofwork.Store) error) error {
	tx, err := u.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer func() { _ = tx.Rollback(ctx) }()

	store := &txStore{
		userRepo:         NewUserRepository(tx),
		coffeeBeanRepo:   NewCoffeeBeanRepository(tx),
		usageHistoryRepo: NewUsageHistoryRepository(tx),
	}

	if err := fn(store); err != nil {
		return err
	}
	return tx.Commit(ctx)
}

type txStore struct {
	userRepo         user.Repository
	coffeeBeanRepo   coffeebean.Repository
	usageHistoryRepo usagehistory.Repository
}

func (s *txStore) UserRepo() user.Repository                 { return s.userRepo }
func (s *txStore) CoffeeBeanRepo() coffeebean.Repository     { return s.coffeeBeanRepo }
func (s *txStore) UsageHistoryRepo() usagehistory.Repository { return s.usageHistoryRepo }
