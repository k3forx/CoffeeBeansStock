package unitofwork

import (
	"context"

	"github.com/k3forx/CoffeeBeansStock/backend/internal/domain/coffeebean"
	"github.com/k3forx/CoffeeBeansStock/backend/internal/domain/usagehistory"
	"github.com/k3forx/CoffeeBeansStock/backend/internal/domain/user"
)

type UnitOfWork interface {
	RunInTx(ctx context.Context, fn func(store Store) error) error
}

type Store interface {
	UserRepo() user.Repository
	CoffeeBeanRepo() coffeebean.Repository
	UsageHistoryRepo() usagehistory.Repository
}
