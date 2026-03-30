package usecase_test

import (
	"context"
	"time"

	"github.com/google/uuid"
	"go.uber.org/mock/gomock"

	"github.com/k3forx/CoffeeBeansStock/backend/internal/domain/coffeebean"
	"github.com/k3forx/CoffeeBeansStock/backend/internal/domain/unitofwork"
	"github.com/k3forx/CoffeeBeansStock/backend/internal/domain/user"
	"github.com/k3forx/CoffeeBeansStock/backend/internal/usecase/mock"
)

//go:fix inline
func ptr[T any](v T) *T { return new(v) }

func newTestUser(id uuid.UUID) *user.User {
	now := time.Now()
	return user.Reconstruct(id, "", "", "", 100, true, user.DefaultGramsPerCup(), now, now)
}

func newTestBean(userID uuid.UUID) *coffeebean.CoffeeBean {
	return coffeebean.Reconstruct(
		uuid.New(), userID, "Test Bean", nil,
		coffeebean.RoastShallow, nil, coffeebean.ReconstructStock(100), nil,
		time.Now(), time.Now(),
	)
}

// fakeRunInTx returns a mock UoW whose RunInTx executes the callback with the given store.
func fakeRunInTx(ctrl *gomock.Controller, store *mock.MockStore) *mock.MockUnitOfWork {
	uow := mock.NewMockUnitOfWork(ctrl)
	uow.EXPECT().RunInTx(gomock.Any(), gomock.Any()).DoAndReturn(
		func(ctx context.Context, fn func(unitofwork.Store) error) error {
			return fn(store)
		},
	)
	return uow
}

// setupUoWWithStore creates a mock UoW that executes the fn with the given store.
func setupUoWWithStore(uow *mock.MockUnitOfWork, store *mock.MockStore) {
	uow.EXPECT().RunInTx(gomock.Any(), gomock.Any()).DoAndReturn(
		func(_ any, fn func(unitofwork.Store) error) error {
			return fn(store)
		},
	)
}
