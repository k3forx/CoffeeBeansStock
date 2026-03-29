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
	"github.com/k3forx/CoffeeBeansStock/backend/internal/usecase"
	"github.com/k3forx/CoffeeBeansStock/backend/internal/usecase/mock"
)

func TestCoffeeBeansService_List(t *testing.T) {
	t.Parallel()

	cases := map[string]struct {
		setup   func(ctrl *gomock.Controller) *usecase.CoffeeBeansService
		in      usecase.ListBeansInput
		out     *usecase.ListBeansResult
		wantErr bool
	}{
		"limit=0はデフォルト20に正規化される": {
			setup: func(ctrl *gomock.Controller) *usecase.CoffeeBeansService {
				r := mock.NewMockCoffeeBeanRepository(ctrl)
				r.EXPECT().CountByUserID(gomock.Any(), gomock.Any()).Return(int64(5), nil)
				r.EXPECT().ListByUserID(gomock.Any(), gomock.Any(), int32(20), int32(0)).Return([]*coffeebean.CoffeeBean{}, nil)
				return usecase.NewCoffeeBeansService(r, nil)
			},
			in: usecase.ListBeansInput{Limit: 0, Offset: 0},
			out: &usecase.ListBeansResult{
				Beans:      []*coffeebean.CoffeeBean{},
				Pagination: usecase.PaginationResponse{Total: 5, Limit: 20, Offset: 0, HasMore: false},
			},
		},
		"limit>100はデフォルト20に正規化される": {
			setup: func(ctrl *gomock.Controller) *usecase.CoffeeBeansService {
				r := mock.NewMockCoffeeBeanRepository(ctrl)
				r.EXPECT().CountByUserID(gomock.Any(), gomock.Any()).Return(int64(5), nil)
				r.EXPECT().ListByUserID(gomock.Any(), gomock.Any(), int32(20), int32(0)).Return([]*coffeebean.CoffeeBean{}, nil)
				return usecase.NewCoffeeBeansService(r, nil)
			},
			in: usecase.ListBeansInput{Limit: 101, Offset: 0},
			out: &usecase.ListBeansResult{
				Beans:      []*coffeebean.CoffeeBean{},
				Pagination: usecase.PaginationResponse{Total: 5, Limit: 20, Offset: 0, HasMore: false},
			},
		},
		"offset<0は0に正規化される": {
			setup: func(ctrl *gomock.Controller) *usecase.CoffeeBeansService {
				r := mock.NewMockCoffeeBeanRepository(ctrl)
				r.EXPECT().CountByUserID(gomock.Any(), gomock.Any()).Return(int64(3), nil)
				r.EXPECT().ListByUserID(gomock.Any(), gomock.Any(), int32(10), int32(0)).Return([]*coffeebean.CoffeeBean{}, nil)
				return usecase.NewCoffeeBeansService(r, nil)
			},
			in: usecase.ListBeansInput{Limit: 10, Offset: -1},
			out: &usecase.ListBeansResult{
				Beans:      []*coffeebean.CoffeeBean{},
				Pagination: usecase.PaginationResponse{Total: 3, Limit: 10, Offset: 0, HasMore: false},
			},
		},
		"hasMoreがtrueになる": {
			setup: func(ctrl *gomock.Controller) *usecase.CoffeeBeansService {
				r := mock.NewMockCoffeeBeanRepository(ctrl)
				r.EXPECT().CountByUserID(gomock.Any(), gomock.Any()).Return(int64(15), nil)
				r.EXPECT().ListByUserID(gomock.Any(), gomock.Any(), int32(10), int32(0)).Return([]*coffeebean.CoffeeBean{}, nil)
				return usecase.NewCoffeeBeansService(r, nil)
			},
			in: usecase.ListBeansInput{Limit: 10, Offset: 0},
			out: &usecase.ListBeansResult{
				Beans:      []*coffeebean.CoffeeBean{},
				Pagination: usecase.PaginationResponse{Total: 15, Limit: 10, Offset: 0, HasMore: true},
			},
		},
		"hasMoreがfalseになる（ちょうど）": {
			setup: func(ctrl *gomock.Controller) *usecase.CoffeeBeansService {
				r := mock.NewMockCoffeeBeanRepository(ctrl)
				r.EXPECT().CountByUserID(gomock.Any(), gomock.Any()).Return(int64(15), nil)
				r.EXPECT().ListByUserID(gomock.Any(), gomock.Any(), int32(10), int32(5)).Return([]*coffeebean.CoffeeBean{}, nil)
				return usecase.NewCoffeeBeansService(r, nil)
			},
			in: usecase.ListBeansInput{Limit: 10, Offset: 5},
			out: &usecase.ListBeansResult{
				Beans:      []*coffeebean.CoffeeBean{},
				Pagination: usecase.PaginationResponse{Total: 15, Limit: 10, Offset: 5, HasMore: false},
			},
		},
		"CountByUserIDがエラーを返す": {
			setup: func(ctrl *gomock.Controller) *usecase.CoffeeBeansService {
				r := mock.NewMockCoffeeBeanRepository(ctrl)
				r.EXPECT().CountByUserID(gomock.Any(), gomock.Any()).Return(int64(0), errors.New("db error"))
				return usecase.NewCoffeeBeansService(r, nil)
			},
			in:      usecase.ListBeansInput{Limit: 10, Offset: 0},
			wantErr: true,
		},
		"ListByUserIDがエラーを返す": {
			setup: func(ctrl *gomock.Controller) *usecase.CoffeeBeansService {
				r := mock.NewMockCoffeeBeanRepository(ctrl)
				r.EXPECT().CountByUserID(gomock.Any(), gomock.Any()).Return(int64(5), nil)
				r.EXPECT().ListByUserID(gomock.Any(), gomock.Any(), int32(10), int32(0)).Return(nil, errors.New("db error"))
				return usecase.NewCoffeeBeansService(r, nil)
			},
			in:      usecase.ListBeansInput{Limit: 10, Offset: 0},
			wantErr: true,
		},
	}

	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			svc := c.setup(ctrl)

			out, err := svc.List(t.Context(), c.in)

			if c.wantErr {
				if err == nil {
					t.Errorf("expected error, got nil")
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

func TestCoffeeBeansService_Create(t *testing.T) {
	t.Parallel()

	cases := map[string]struct {
		setup          func(ctrl *gomock.Controller) *usecase.CoffeeBeansService
		in             usecase.CreateBeanInput
		wantValidation bool
		wantErr        bool
	}{
		"無効なRoastLevelはValidationErrorを返す": {
			setup: func(ctrl *gomock.Controller) *usecase.CoffeeBeansService {
				r := mock.NewMockCoffeeBeanRepository(ctrl)
				return usecase.NewCoffeeBeansService(r, nil)
			},
			in: usecase.CreateBeanInput{
				Name:         "Test",
				RoastLevel:   "invalid",
				CurrentStock: 100,
			},
			wantValidation: true,
			wantErr:        true,
		},
		"無効なRoastDetailはValidationErrorを返す": {
			setup: func(ctrl *gomock.Controller) *usecase.CoffeeBeansService {
				r := mock.NewMockCoffeeBeanRepository(ctrl)
				return usecase.NewCoffeeBeansService(r, nil)
			},
			in: usecase.CreateBeanInput{
				Name:         "Test",
				RoastLevel:   "shallow",
				RoastDetail:  ptr("french"),
				CurrentStock: 100,
			},
			wantValidation: true,
			wantErr:        true,
		},
		"範囲外のStockはValidationErrorを返す": {
			setup: func(ctrl *gomock.Controller) *usecase.CoffeeBeansService {
				r := mock.NewMockCoffeeBeanRepository(ctrl)
				return usecase.NewCoffeeBeansService(r, nil)
			},
			in: usecase.CreateBeanInput{
				Name:         "Test",
				RoastLevel:   "shallow",
				CurrentStock: -1,
			},
			wantValidation: true,
			wantErr:        true,
		},
		"Repo.Saveのエラーが伝播する": {
			setup: func(ctrl *gomock.Controller) *usecase.CoffeeBeansService {
				r := mock.NewMockCoffeeBeanRepository(ctrl)
				r.EXPECT().Save(gomock.Any(), gomock.Any()).Return(errors.New("db error"))
				return usecase.NewCoffeeBeansService(r, nil)
			},
			in: usecase.CreateBeanInput{
				Name:         "Test",
				RoastLevel:   "shallow",
				CurrentStock: 100,
			},
			wantErr: true,
		},
		"RoastDetailなしで正常に作成できる": {
			setup: func(ctrl *gomock.Controller) *usecase.CoffeeBeansService {
				r := mock.NewMockCoffeeBeanRepository(ctrl)
				r.EXPECT().Save(gomock.Any(), gomock.Any()).Return(nil)
				return usecase.NewCoffeeBeansService(r, nil)
			},
			in: usecase.CreateBeanInput{
				Name:         "Sidama",
				RoastLevel:   "shallow",
				CurrentStock: 200,
			},
		},
		"RoastDetailありで正常に作成できる": {
			setup: func(ctrl *gomock.Controller) *usecase.CoffeeBeansService {
				r := mock.NewMockCoffeeBeanRepository(ctrl)
				r.EXPECT().Save(gomock.Any(), gomock.Any()).Return(nil)
				return usecase.NewCoffeeBeansService(r, nil)
			},
			in: usecase.CreateBeanInput{
				Name:         "Ethiopia",
				RoastLevel:   "shallow",
				RoastDetail:  ptr("light"),
				CurrentStock: 150,
			},
		},
	}

	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			svc := c.setup(ctrl)

			out, err := svc.Create(t.Context(), c.in)

			if c.wantErr {
				if err == nil {
					t.Errorf("expected error, got nil")
					return
				}
				if c.wantValidation {
					var ve *domain.ValidationError
					var ves domain.ValidationErrors
					if !errors.As(err, &ve) && !errors.As(err, &ves) {
						t.Errorf("expected ValidationError, got %T: %v", err, err)
					}
				}
				return
			}
			if err != nil {
				t.Errorf("unexpected error: %v", err)
				return
			}
			if out == nil || out.Bean == nil {
				t.Errorf("expected bean, got nil")
			}
		})
	}
}

func TestCoffeeBeansService_GetByID(t *testing.T) {
	t.Parallel()

	userID := uuid.New()
	otherUserID := uuid.New()
	beanID := uuid.New()
	testBean := coffeebean.Reconstruct(
		beanID, userID, "Test Bean", nil,
		coffeebean.RoastShallow, nil, coffeebean.ReconstructStock(100), nil,
		time.Now(), time.Now(),
	)

	cases := map[string]struct {
		setup   func(ctrl *gomock.Controller) *usecase.CoffeeBeansService
		in      usecase.GetBeanByIDInput
		out     *usecase.GetBeanByIDOutput
		wantErr error
	}{
		"自分のbeanを取得できる": {
			setup: func(ctrl *gomock.Controller) *usecase.CoffeeBeansService {
				r := mock.NewMockCoffeeBeanRepository(ctrl)
				r.EXPECT().GetByID(gomock.Any(), beanID).Return(testBean, nil)
				return usecase.NewCoffeeBeansService(r, nil)
			},
			in:  usecase.GetBeanByIDInput{UserID: userID, BeanID: beanID},
			out: &usecase.GetBeanByIDOutput{Bean: testBean},
		},
		"ErrNotFoundが伝播する": {
			setup: func(ctrl *gomock.Controller) *usecase.CoffeeBeansService {
				r := mock.NewMockCoffeeBeanRepository(ctrl)
				r.EXPECT().GetByID(gomock.Any(), beanID).Return(nil, domain.ErrNotFound)
				return usecase.NewCoffeeBeansService(r, nil)
			},
			in:      usecase.GetBeanByIDInput{UserID: userID, BeanID: beanID},
			wantErr: domain.ErrNotFound,
		},
		"別ユーザーのbeanはErrForbiddenを返す": {
			setup: func(ctrl *gomock.Controller) *usecase.CoffeeBeansService {
				r := mock.NewMockCoffeeBeanRepository(ctrl)
				r.EXPECT().GetByID(gomock.Any(), beanID).Return(testBean, nil)
				return usecase.NewCoffeeBeansService(r, nil)
			},
			in:      usecase.GetBeanByIDInput{UserID: otherUserID, BeanID: beanID},
			wantErr: domain.ErrForbidden,
		},
	}

	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			svc := c.setup(ctrl)

			out, err := svc.GetByID(t.Context(), c.in)

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
			if diff := cmp.Diff(c.out, out, cmp.AllowUnexported(coffeebean.CoffeeBean{}, coffeebean.RoastLevel{}, coffeebean.Stock{})); diff != "" {
				t.Errorf("mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func TestCoffeeBeansService_Update(t *testing.T) {
	t.Parallel()

	userID := uuid.New()
	otherUserID := uuid.New()
	beanID := uuid.New()

	cases := map[string]struct {
		setup          func(ctrl *gomock.Controller) *usecase.CoffeeBeansService
		in             usecase.UpdateBeanInput
		wantValidation bool
		wantErr        error
	}{
		"無効なRoastLevelはValidationErrorを返す（トランザクション前に失敗）": {
			setup: func(ctrl *gomock.Controller) *usecase.CoffeeBeansService {
				r := mock.NewMockCoffeeBeanRepository(ctrl)
				return usecase.NewCoffeeBeansService(r, nil)
			},
			in:             usecase.UpdateBeanInput{UserID: userID, BeanID: beanID, RoastLevel: ptr("invalid")},
			wantValidation: true,
		},
		"ErrNotFoundが伝播する": {
			setup: func(ctrl *gomock.Controller) *usecase.CoffeeBeansService {
				r := mock.NewMockCoffeeBeanRepository(ctrl)
				store := mock.NewMockStore(ctrl)
				store.EXPECT().CoffeeBeanRepo().Return(r).AnyTimes()
				r.EXPECT().GetByIDForUpdate(gomock.Any(), beanID).Return(nil, domain.ErrNotFound)
				uow := fakeRunInTx(ctrl, store)
				return usecase.NewCoffeeBeansService(r, uow)
			},
			in:      usecase.UpdateBeanInput{UserID: userID, BeanID: beanID},
			wantErr: domain.ErrNotFound,
		},
		"別ユーザーのbeanはErrForbiddenを返す": {
			setup: func(ctrl *gomock.Controller) *usecase.CoffeeBeansService {
				r := mock.NewMockCoffeeBeanRepository(ctrl)
				testBean := coffeebean.Reconstruct(
					beanID, userID, "Test Bean", nil,
					coffeebean.RoastShallow, nil, coffeebean.ReconstructStock(100), nil,
					time.Now(), time.Now(),
				)
				store := mock.NewMockStore(ctrl)
				store.EXPECT().CoffeeBeanRepo().Return(r).AnyTimes()
				r.EXPECT().GetByIDForUpdate(gomock.Any(), beanID).Return(testBean, nil)
				uow := fakeRunInTx(ctrl, store)
				return usecase.NewCoffeeBeansService(r, uow)
			},
			in:      usecase.UpdateBeanInput{UserID: otherUserID, BeanID: beanID},
			wantErr: domain.ErrForbidden,
		},
		"Repo.Updateのエラーが伝播する": {
			setup: func(ctrl *gomock.Controller) *usecase.CoffeeBeansService {
				r := mock.NewMockCoffeeBeanRepository(ctrl)
				testBean := coffeebean.Reconstruct(
					beanID, userID, "Test Bean", nil,
					coffeebean.RoastShallow, nil, coffeebean.ReconstructStock(100), nil,
					time.Now(), time.Now(),
				)
				store := mock.NewMockStore(ctrl)
				store.EXPECT().CoffeeBeanRepo().Return(r).AnyTimes()
				r.EXPECT().GetByIDForUpdate(gomock.Any(), beanID).Return(testBean, nil)
				r.EXPECT().Update(gomock.Any(), gomock.Any()).Return(errors.New("db error"))
				uow := fakeRunInTx(ctrl, store)
				return usecase.NewCoffeeBeansService(r, uow)
			},
			in: usecase.UpdateBeanInput{UserID: userID, BeanID: beanID, Name: ptr("Updated")},
		},
		"正常に更新できる（部分更新）": {
			setup: func(ctrl *gomock.Controller) *usecase.CoffeeBeansService {
				r := mock.NewMockCoffeeBeanRepository(ctrl)
				testBean := coffeebean.Reconstruct(
					beanID, userID, "Test Bean", nil,
					coffeebean.RoastShallow, nil, coffeebean.ReconstructStock(100), nil,
					time.Now(), time.Now(),
				)
				store := mock.NewMockStore(ctrl)
				store.EXPECT().CoffeeBeanRepo().Return(r).AnyTimes()
				r.EXPECT().GetByIDForUpdate(gomock.Any(), beanID).Return(testBean, nil)
				r.EXPECT().Update(gomock.Any(), gomock.Any()).Return(nil)
				uow := fakeRunInTx(ctrl, store)
				return usecase.NewCoffeeBeansService(r, uow)
			},
			in: usecase.UpdateBeanInput{UserID: userID, BeanID: beanID, Name: ptr("New Name")},
		},
	}

	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			svc := c.setup(ctrl)

			_, err := svc.Update(t.Context(), c.in)

			if c.wantValidation {
				if err == nil {
					t.Errorf("expected validation error, got nil")
					return
				}
				var ve *domain.ValidationError
				if !errors.As(err, &ve) {
					t.Errorf("expected ValidationError, got %T: %v", err, err)
				}
				return
			}
			if c.wantErr != nil {
				if !errors.Is(err, c.wantErr) {
					t.Errorf("error = %v, want %v", err, c.wantErr)
				}
				return
			}
			// "Repo.Updateのエラーが伝播する" — expects non-nil error
			if name == "Repo.Updateのエラーが伝播する" {
				if err == nil {
					t.Errorf("expected error, got nil")
				}
				return
			}
			if err != nil {
				t.Errorf("unexpected error: %v", err)
			}
		})
	}
}

func TestCoffeeBeansService_Delete(t *testing.T) {
	t.Parallel()

	userID := uuid.New()
	otherUserID := uuid.New()
	beanID := uuid.New()
	testBean := coffeebean.Reconstruct(
		beanID, userID, "Test Bean", nil,
		coffeebean.RoastShallow, nil, coffeebean.ReconstructStock(100), nil,
		time.Now(), time.Now(),
	)

	cases := map[string]struct {
		setup   func(ctrl *gomock.Controller) *usecase.CoffeeBeansService
		in      usecase.DeleteBeanInput
		wantErr error
	}{
		"正常に削除できる": {
			setup: func(ctrl *gomock.Controller) *usecase.CoffeeBeansService {
				r := mock.NewMockCoffeeBeanRepository(ctrl)
				store := mock.NewMockStore(ctrl)
				store.EXPECT().CoffeeBeanRepo().Return(r).AnyTimes()
				r.EXPECT().GetByIDForUpdate(gomock.Any(), beanID).Return(testBean, nil)
				r.EXPECT().SoftDelete(gomock.Any(), beanID, userID).Return(nil)
				uow := fakeRunInTx(ctrl, store)
				return usecase.NewCoffeeBeansService(r, uow)
			},
			in: usecase.DeleteBeanInput{UserID: userID, BeanID: beanID},
		},
		"ErrNotFoundが伝播する": {
			setup: func(ctrl *gomock.Controller) *usecase.CoffeeBeansService {
				r := mock.NewMockCoffeeBeanRepository(ctrl)
				store := mock.NewMockStore(ctrl)
				store.EXPECT().CoffeeBeanRepo().Return(r).AnyTimes()
				r.EXPECT().GetByIDForUpdate(gomock.Any(), beanID).Return(nil, domain.ErrNotFound)
				uow := fakeRunInTx(ctrl, store)
				return usecase.NewCoffeeBeansService(r, uow)
			},
			in:      usecase.DeleteBeanInput{UserID: userID, BeanID: beanID},
			wantErr: domain.ErrNotFound,
		},
		"別ユーザーのbeanはErrForbiddenを返す": {
			setup: func(ctrl *gomock.Controller) *usecase.CoffeeBeansService {
				r := mock.NewMockCoffeeBeanRepository(ctrl)
				store := mock.NewMockStore(ctrl)
				store.EXPECT().CoffeeBeanRepo().Return(r).AnyTimes()
				r.EXPECT().GetByIDForUpdate(gomock.Any(), beanID).Return(testBean, nil)
				uow := fakeRunInTx(ctrl, store)
				return usecase.NewCoffeeBeansService(r, uow)
			},
			in:      usecase.DeleteBeanInput{UserID: otherUserID, BeanID: beanID},
			wantErr: domain.ErrForbidden,
		},
		"SoftDeleteのエラーが伝播する": {
			setup: func(ctrl *gomock.Controller) *usecase.CoffeeBeansService {
				r := mock.NewMockCoffeeBeanRepository(ctrl)
				store := mock.NewMockStore(ctrl)
				store.EXPECT().CoffeeBeanRepo().Return(r).AnyTimes()
				r.EXPECT().GetByIDForUpdate(gomock.Any(), beanID).Return(testBean, nil)
				r.EXPECT().SoftDelete(gomock.Any(), beanID, userID).Return(errors.New("db error"))
				uow := fakeRunInTx(ctrl, store)
				return usecase.NewCoffeeBeansService(r, uow)
			},
			in: usecase.DeleteBeanInput{UserID: userID, BeanID: beanID},
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
			// "SoftDeleteのエラーが伝播する"
			if name == "SoftDeleteのエラーが伝播する" {
				if err == nil {
					t.Errorf("expected error, got nil")
				}
				return
			}
			if err != nil {
				t.Errorf("unexpected error: %v", err)
			}
		})
	}
}
