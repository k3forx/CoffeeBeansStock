package usecase_test

import (
	"errors"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/google/uuid"
	"go.uber.org/mock/gomock"

	domain "github.com/k3forx/CoffeeBeansStock/backend/internal/domain"
	"github.com/k3forx/CoffeeBeansStock/backend/internal/domain/coffeebean"
	"github.com/k3forx/CoffeeBeansStock/backend/internal/domain/usagehistory"
	"github.com/k3forx/CoffeeBeansStock/backend/internal/usecase"
	"github.com/k3forx/CoffeeBeansStock/backend/internal/usecase/mock"
)

func TestUsageHistoryService_Create(t *testing.T) {
	t.Parallel()

	userID := uuid.New()
	otherUserID := uuid.New()
	beanID := uuid.New()
	usageDate := time.Date(2026, 3, 29, 0, 0, 0, 0, time.UTC)

	cases := map[string]struct {
		setup   func(ctrl *gomock.Controller) *usecase.UsageHistoryService
		in      usecase.CreateUsageInput
		wantErr error
		wantAny bool // expect any error (not a specific one)
	}{
		"正常に作成できる": {
			setup: func(ctrl *gomock.Controller) *usecase.UsageHistoryService {
				usageRepo := mock.NewMockUsageHistoryRepository(ctrl)
				beanRepo := mock.NewMockCoffeeBeanRepository(ctrl)
				uow := mock.NewMockUnitOfWork(ctrl)
				store := mock.NewMockStore(ctrl)

				bean := coffeebean.Reconstruct(
					beanID, userID, "Test", nil,
					coffeebean.RoastShallow, nil, coffeebean.ReconstructStock(100), nil,
					time.Now(), time.Now(),
				)

				storeBeanRepo := mock.NewMockCoffeeBeanRepository(ctrl)
				storeUsageRepo := mock.NewMockUsageHistoryRepository(ctrl)
				store.EXPECT().CoffeeBeanRepo().Return(storeBeanRepo).AnyTimes()
				store.EXPECT().UsageHistoryRepo().Return(storeUsageRepo).AnyTimes()
				storeBeanRepo.EXPECT().GetByID(gomock.Any(), beanID).Return(bean, nil)
				storeBeanRepo.EXPECT().UpdateStock(gomock.Any(), beanID, coffeebean.ReconstructStock(90)).Return(nil)
				storeUsageRepo.EXPECT().Save(gomock.Any(), gomock.Any()).Return(nil)

				setupUoWWithStore(uow, store)

				return usecase.NewUsageHistoryService(usageRepo, beanRepo, uow)
			},
			in: usecase.CreateUsageInput{
				UserID:       userID,
				CoffeeBeanID: beanID,
				UsageDate:    usageDate,
				Quantity:     10,
			},
		},
		"数量0はValidationErrorを返す": {
			setup: func(ctrl *gomock.Controller) *usecase.UsageHistoryService {
				usageRepo := mock.NewMockUsageHistoryRepository(ctrl)
				beanRepo := mock.NewMockCoffeeBeanRepository(ctrl)
				uow := mock.NewMockUnitOfWork(ctrl)
				return usecase.NewUsageHistoryService(usageRepo, beanRepo, uow)
			},
			in: usecase.CreateUsageInput{
				UserID:       userID,
				CoffeeBeanID: beanID,
				UsageDate:    usageDate,
				Quantity:     0,
			},
			wantAny: true,
		},
		"存在しないbeanはErrNotFoundを返す": {
			setup: func(ctrl *gomock.Controller) *usecase.UsageHistoryService {
				usageRepo := mock.NewMockUsageHistoryRepository(ctrl)
				beanRepo := mock.NewMockCoffeeBeanRepository(ctrl)
				uow := mock.NewMockUnitOfWork(ctrl)
				store := mock.NewMockStore(ctrl)

				storeBeanRepo := mock.NewMockCoffeeBeanRepository(ctrl)
				store.EXPECT().CoffeeBeanRepo().Return(storeBeanRepo).AnyTimes()
				storeBeanRepo.EXPECT().GetByID(gomock.Any(), beanID).Return(nil, domain.ErrNotFound)

				setupUoWWithStore(uow, store)

				return usecase.NewUsageHistoryService(usageRepo, beanRepo, uow)
			},
			in: usecase.CreateUsageInput{
				UserID:       userID,
				CoffeeBeanID: beanID,
				UsageDate:    usageDate,
				Quantity:     10,
			},
			wantErr: domain.ErrNotFound,
		},
		"別ユーザーのbeanはErrForbiddenを返す": {
			setup: func(ctrl *gomock.Controller) *usecase.UsageHistoryService {
				usageRepo := mock.NewMockUsageHistoryRepository(ctrl)
				beanRepo := mock.NewMockCoffeeBeanRepository(ctrl)
				uow := mock.NewMockUnitOfWork(ctrl)
				store := mock.NewMockStore(ctrl)

				bean := coffeebean.Reconstruct(
					beanID, userID, "Test", nil,
					coffeebean.RoastShallow, nil, coffeebean.ReconstructStock(100), nil,
					time.Now(), time.Now(),
				)

				storeBeanRepo := mock.NewMockCoffeeBeanRepository(ctrl)
				store.EXPECT().CoffeeBeanRepo().Return(storeBeanRepo).AnyTimes()
				storeBeanRepo.EXPECT().GetByID(gomock.Any(), beanID).Return(bean, nil)

				setupUoWWithStore(uow, store)

				return usecase.NewUsageHistoryService(usageRepo, beanRepo, uow)
			},
			in: usecase.CreateUsageInput{
				UserID:       otherUserID,
				CoffeeBeanID: beanID,
				UsageDate:    usageDate,
				Quantity:     10,
			},
			wantErr: domain.ErrForbidden,
		},
		"在庫不足はErrInsufficientStockを返す": {
			setup: func(ctrl *gomock.Controller) *usecase.UsageHistoryService {
				usageRepo := mock.NewMockUsageHistoryRepository(ctrl)
				beanRepo := mock.NewMockCoffeeBeanRepository(ctrl)
				uow := mock.NewMockUnitOfWork(ctrl)
				store := mock.NewMockStore(ctrl)

				bean := coffeebean.Reconstruct(
					beanID, userID, "Test", nil,
					coffeebean.RoastShallow, nil, coffeebean.ReconstructStock(5), nil,
					time.Now(), time.Now(),
				)

				storeBeanRepo := mock.NewMockCoffeeBeanRepository(ctrl)
				store.EXPECT().CoffeeBeanRepo().Return(storeBeanRepo).AnyTimes()
				storeBeanRepo.EXPECT().GetByID(gomock.Any(), beanID).Return(bean, nil)

				setupUoWWithStore(uow, store)

				return usecase.NewUsageHistoryService(usageRepo, beanRepo, uow)
			},
			in: usecase.CreateUsageInput{
				UserID:       userID,
				CoffeeBeanID: beanID,
				UsageDate:    usageDate,
				Quantity:     10,
			},
			wantErr: domain.ErrInsufficientStock,
		},
	}

	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			svc := c.setup(ctrl)

			out, err := svc.Create(t.Context(), c.in)

			if c.wantErr != nil {
				if !errors.Is(err, c.wantErr) {
					t.Errorf("error = %v, want %v", err, c.wantErr)
				}
				return
			}
			if c.wantAny {
				if err == nil {
					t.Errorf("expected error, got nil")
				}
				return
			}
			if err != nil {
				t.Errorf("unexpected error: %v", err)
				return
			}
			if out == nil || out.Usage == nil {
				t.Errorf("expected usage, got nil")
				return
			}
			if out.Usage.CoffeeBeanID() != c.in.CoffeeBeanID {
				t.Errorf("CoffeeBeanID = %v, want %v", out.Usage.CoffeeBeanID(), c.in.CoffeeBeanID)
			}
			if out.Usage.Quantity().Value() != c.in.Quantity {
				t.Errorf("Quantity = %d, want %d", out.Usage.Quantity().Value(), c.in.Quantity)
			}
		})
	}
}

func TestUsageHistoryService_Delete(t *testing.T) {
	t.Parallel()

	userID := uuid.New()
	otherUserID := uuid.New()
	beanID := uuid.New()
	usageID := uuid.New()

	testUsage := usagehistory.Reconstruct(
		usageID, beanID, userID,
		time.Date(2026, 3, 29, 0, 0, 0, 0, time.UTC),
		domain.ReconstructQuantity(10),
		nil, time.Now(),
	)

	cases := map[string]struct {
		setup   func(ctrl *gomock.Controller) *usecase.UsageHistoryService
		in      usecase.DeleteUsageInput
		wantErr error
	}{
		"正常に削除できる": {
			setup: func(ctrl *gomock.Controller) *usecase.UsageHistoryService {
				usageRepo := mock.NewMockUsageHistoryRepository(ctrl)
				beanRepo := mock.NewMockCoffeeBeanRepository(ctrl)
				uow := mock.NewMockUnitOfWork(ctrl)
				store := mock.NewMockStore(ctrl)

				usageRepo.EXPECT().GetByID(gomock.Any(), usageID).Return(testUsage, nil)

				bean := coffeebean.Reconstruct(
					beanID, userID, "Test", nil,
					coffeebean.RoastShallow, nil, coffeebean.ReconstructStock(90), nil,
					time.Now(), time.Now(),
				)

				storeBeanRepo := mock.NewMockCoffeeBeanRepository(ctrl)
				storeUsageRepo := mock.NewMockUsageHistoryRepository(ctrl)
				store.EXPECT().CoffeeBeanRepo().Return(storeBeanRepo).AnyTimes()
				store.EXPECT().UsageHistoryRepo().Return(storeUsageRepo).AnyTimes()
				storeBeanRepo.EXPECT().GetByID(gomock.Any(), beanID).Return(bean, nil)
				storeBeanRepo.EXPECT().UpdateStock(gomock.Any(), beanID, coffeebean.ReconstructStock(100)).Return(nil)
				storeUsageRepo.EXPECT().Delete(gomock.Any(), usageID, userID).Return(nil)

				setupUoWWithStore(uow, store)

				return usecase.NewUsageHistoryService(usageRepo, beanRepo, uow)
			},
			in: usecase.DeleteUsageInput{UserID: userID, CoffeeBeanID: beanID, UsageID: usageID},
		},
		"使用記録が見つからない場合はErrNotFoundを返す": {
			setup: func(ctrl *gomock.Controller) *usecase.UsageHistoryService {
				usageRepo := mock.NewMockUsageHistoryRepository(ctrl)
				beanRepo := mock.NewMockCoffeeBeanRepository(ctrl)
				uow := mock.NewMockUnitOfWork(ctrl)

				usageRepo.EXPECT().GetByID(gomock.Any(), usageID).Return(nil, domain.ErrNotFound)

				return usecase.NewUsageHistoryService(usageRepo, beanRepo, uow)
			},
			in:      usecase.DeleteUsageInput{UserID: userID, CoffeeBeanID: beanID, UsageID: usageID},
			wantErr: domain.ErrNotFound,
		},
		"別ユーザーの使用記録はErrForbiddenを返す": {
			setup: func(ctrl *gomock.Controller) *usecase.UsageHistoryService {
				usageRepo := mock.NewMockUsageHistoryRepository(ctrl)
				beanRepo := mock.NewMockCoffeeBeanRepository(ctrl)
				uow := mock.NewMockUnitOfWork(ctrl)

				usageRepo.EXPECT().GetByID(gomock.Any(), usageID).Return(testUsage, nil)

				return usecase.NewUsageHistoryService(usageRepo, beanRepo, uow)
			},
			in:      usecase.DeleteUsageInput{UserID: otherUserID, CoffeeBeanID: beanID, UsageID: usageID},
			wantErr: domain.ErrForbidden,
		},
	}

	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			svc := c.setup(ctrl)

			err := svc.Delete(t.Context(), c.in)

			if c.wantErr != nil {
				if !errors.Is(err, c.wantErr) {
					t.Errorf("error = %v, want %v", err, c.wantErr)
				}
				return
			}
			if err != nil {
				t.Errorf("unexpected error: %v", err)
			}
		})
	}
}

func TestUsageHistoryService_ListByCoffeeBean(t *testing.T) {
	t.Parallel()

	userID := uuid.New()
	otherUserID := uuid.New()
	beanID := uuid.New()

	testBean := coffeebean.Reconstruct(
		beanID, userID, "Test", nil,
		coffeebean.RoastShallow, nil, coffeebean.ReconstructStock(100), nil,
		time.Now(), time.Now(),
	)

	cases := map[string]struct {
		setup   func(ctrl *gomock.Controller) *usecase.UsageHistoryService
		in      usecase.ListUsageInput
		out     *usecase.ListUsageResult
		wantErr error
	}{
		"正常に一覧取得できる": {
			setup: func(ctrl *gomock.Controller) *usecase.UsageHistoryService {
				usageRepo := mock.NewMockUsageHistoryRepository(ctrl)
				beanRepo := mock.NewMockCoffeeBeanRepository(ctrl)
				uow := mock.NewMockUnitOfWork(ctrl)

				beanRepo.EXPECT().GetByID(gomock.Any(), beanID).Return(testBean, nil)
				usageRepo.EXPECT().CountByCoffeeBeanID(gomock.Any(), beanID).Return(int64(5), nil)
				usageRepo.EXPECT().ListByCoffeeBeanID(gomock.Any(), beanID, int32(20), int32(0)).Return([]*usagehistory.UsageHistory{}, nil)

				return usecase.NewUsageHistoryService(usageRepo, beanRepo, uow)
			},
			in: usecase.ListUsageInput{UserID: userID, CoffeeBeanID: beanID, Limit: 20, Offset: 0},
			out: &usecase.ListUsageResult{
				Usages:     []*usagehistory.UsageHistory{},
				Pagination: usecase.PaginationResponse{Total: 5, Limit: 20, Offset: 0, HasMore: false},
			},
		},
		"limit=0はデフォルト20に正規化される": {
			setup: func(ctrl *gomock.Controller) *usecase.UsageHistoryService {
				usageRepo := mock.NewMockUsageHistoryRepository(ctrl)
				beanRepo := mock.NewMockCoffeeBeanRepository(ctrl)
				uow := mock.NewMockUnitOfWork(ctrl)

				beanRepo.EXPECT().GetByID(gomock.Any(), beanID).Return(testBean, nil)
				usageRepo.EXPECT().CountByCoffeeBeanID(gomock.Any(), beanID).Return(int64(0), nil)
				usageRepo.EXPECT().ListByCoffeeBeanID(gomock.Any(), beanID, int32(20), int32(0)).Return([]*usagehistory.UsageHistory{}, nil)

				return usecase.NewUsageHistoryService(usageRepo, beanRepo, uow)
			},
			in: usecase.ListUsageInput{UserID: userID, CoffeeBeanID: beanID, Limit: 0, Offset: 0},
			out: &usecase.ListUsageResult{
				Usages:     []*usagehistory.UsageHistory{},
				Pagination: usecase.PaginationResponse{Total: 0, Limit: 20, Offset: 0, HasMore: false},
			},
		},
		"beanが見つからない場合はErrNotFoundを返す": {
			setup: func(ctrl *gomock.Controller) *usecase.UsageHistoryService {
				usageRepo := mock.NewMockUsageHistoryRepository(ctrl)
				beanRepo := mock.NewMockCoffeeBeanRepository(ctrl)
				uow := mock.NewMockUnitOfWork(ctrl)

				beanRepo.EXPECT().GetByID(gomock.Any(), beanID).Return(nil, domain.ErrNotFound)

				return usecase.NewUsageHistoryService(usageRepo, beanRepo, uow)
			},
			in:      usecase.ListUsageInput{UserID: userID, CoffeeBeanID: beanID, Limit: 20, Offset: 0},
			wantErr: domain.ErrNotFound,
		},
		"別ユーザーのbeanはErrForbiddenを返す": {
			setup: func(ctrl *gomock.Controller) *usecase.UsageHistoryService {
				usageRepo := mock.NewMockUsageHistoryRepository(ctrl)
				beanRepo := mock.NewMockCoffeeBeanRepository(ctrl)
				uow := mock.NewMockUnitOfWork(ctrl)

				beanRepo.EXPECT().GetByID(gomock.Any(), beanID).Return(testBean, nil)

				return usecase.NewUsageHistoryService(usageRepo, beanRepo, uow)
			},
			in:      usecase.ListUsageInput{UserID: otherUserID, CoffeeBeanID: beanID, Limit: 20, Offset: 0},
			wantErr: domain.ErrForbidden,
		},
	}

	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			svc := c.setup(ctrl)

			out, err := svc.ListByCoffeeBean(t.Context(), c.in)

			if c.wantErr != nil {
				if !errors.Is(err, c.wantErr) {
					t.Errorf("error = %v, want %v", err, c.wantErr)
				}
				return
			}
			if err != nil {
				t.Errorf("unexpected error: %v", err)
				return
			}
			if diff := cmp.Diff(c.out.Pagination, out.Pagination); diff != "" {
				t.Errorf("Pagination mismatch (-want +got):\n%s", diff)
			}
		})
	}
}
