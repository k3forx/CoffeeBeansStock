package usecase_test

import (
	"errors"
	"testing"

	"github.com/google/uuid"
	"go.uber.org/mock/gomock"

	domain "github.com/k3forx/CoffeeBeansStock/backend/internal/domain"
	domainauth "github.com/k3forx/CoffeeBeansStock/backend/internal/domain/auth"
	"github.com/k3forx/CoffeeBeansStock/backend/internal/usecase"
	"github.com/k3forx/CoffeeBeansStock/backend/internal/usecase/mock"
)

func TestAuthService_RegisterAnonymous(t *testing.T) {
	t.Parallel()

	cases := map[string]struct {
		setup   func(ctrl *gomock.Controller) *usecase.AuthService
		wantErr bool
	}{
		"userが保存されトークンペアが返る": {
			setup: func(ctrl *gomock.Controller) *usecase.AuthService {
				ur := mock.NewMockUserRepository(ctrl)
				tm := mock.NewMockTokenManager(ctrl)
				rt := mock.NewMockRefreshTokenRepository(ctrl)
				ur.EXPECT().Save(gomock.Any(), gomock.Any()).Return(nil)
				tm.EXPECT().GenerateTokenPair(gomock.Any()).Return(&domainauth.TokenPair{
					AccessToken:  "access-token",
					RefreshToken: "refresh-token",
				}, nil)
				rt.EXPECT().Store(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
				return usecase.NewAuthService(ur, tm, rt)
			},
		},
		"GenerateTokenPairのエラーが伝播する": {
			setup: func(ctrl *gomock.Controller) *usecase.AuthService {
				ur := mock.NewMockUserRepository(ctrl)
				tm := mock.NewMockTokenManager(ctrl)
				rt := mock.NewMockRefreshTokenRepository(ctrl)
				tm.EXPECT().GenerateTokenPair(gomock.Any()).Return(nil, errors.New("token error"))
				return usecase.NewAuthService(ur, tm, rt)
			},
			wantErr: true,
		},
		"userRepo.Saveのエラーが伝播する": {
			setup: func(ctrl *gomock.Controller) *usecase.AuthService {
				ur := mock.NewMockUserRepository(ctrl)
				tm := mock.NewMockTokenManager(ctrl)
				rt := mock.NewMockRefreshTokenRepository(ctrl)
				tm.EXPECT().GenerateTokenPair(gomock.Any()).Return(&domainauth.TokenPair{
					AccessToken:  "access-token",
					RefreshToken: "refresh-token",
				}, nil)
				ur.EXPECT().Save(gomock.Any(), gomock.Any()).Return(errors.New("db error"))
				return usecase.NewAuthService(ur, tm, rt)
			},
			wantErr: true,
		},
		"refreshTokens.Storeのエラーが伝播する": {
			setup: func(ctrl *gomock.Controller) *usecase.AuthService {
				ur := mock.NewMockUserRepository(ctrl)
				tm := mock.NewMockTokenManager(ctrl)
				rt := mock.NewMockRefreshTokenRepository(ctrl)
				tm.EXPECT().GenerateTokenPair(gomock.Any()).Return(&domainauth.TokenPair{
					AccessToken:  "access-token",
					RefreshToken: "refresh-token",
				}, nil)
				ur.EXPECT().Save(gomock.Any(), gomock.Any()).Return(nil)
				rt.EXPECT().Store(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(errors.New("store error"))
				return usecase.NewAuthService(ur, tm, rt)
			},
			wantErr: true,
		},
	}

	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			svc := c.setup(ctrl)

			result, err := svc.RegisterAnonymous(t.Context())

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
			if result == nil {
				t.Errorf("expected result, got nil")
				return
			}
			if result.AccessToken != "access-token" {
				t.Errorf("AccessToken = %q, want %q", result.AccessToken, "access-token")
			}
			if result.RefreshToken != "refresh-token" {
				t.Errorf("RefreshToken = %q, want %q", result.RefreshToken, "refresh-token")
			}
			if result.User == nil {
				t.Errorf("expected User in result, got nil")
			}
		})
	}
}

func TestAuthService_Refresh(t *testing.T) {
	t.Parallel()

	validUserID := uuid.New()
	testUser := newTestUser(validUserID)
	newTokenPair := &domainauth.TokenPair{
		AccessToken:  "new-access-token",
		RefreshToken: "new-refresh-token",
	}

	cases := map[string]struct {
		setup   func(ctrl *gomock.Controller) *usecase.AuthService
		in      usecase.RefreshInput
		out     *usecase.RefreshResult
		wantErr error
	}{
		"新しいトークンペアが返る": {
			setup: func(ctrl *gomock.Controller) *usecase.AuthService {
				ur := mock.NewMockUserRepository(ctrl)
				tm := mock.NewMockTokenManager(ctrl)
				rt := mock.NewMockRefreshTokenRepository(ctrl)
				claims := &domainauth.TokenClaims{UserID: validUserID.String()}
				tm.EXPECT().ValidateToken("valid-token").Return(claims, nil)
				ur.EXPECT().GetByID(gomock.Any(), validUserID).Return(testUser, nil)
				rt.EXPECT().ExistsByHash(gomock.Any(), gomock.Any()).Return(true, nil)
				rt.EXPECT().DeleteByHash(gomock.Any(), gomock.Any()).Return(nil)
				tm.EXPECT().GenerateTokenPair(validUserID).Return(newTokenPair, nil)
				rt.EXPECT().Store(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
				return usecase.NewAuthService(ur, tm, rt)
			},
			in: usecase.RefreshInput{RefreshToken: "valid-token"},
			out: &usecase.RefreshResult{
				AccessToken:  "new-access-token",
				RefreshToken: "new-refresh-token",
			},
		},
		"ValidateTokenのエラーが伝播する": {
			setup: func(ctrl *gomock.Controller) *usecase.AuthService {
				ur := mock.NewMockUserRepository(ctrl)
				tm := mock.NewMockTokenManager(ctrl)
				rt := mock.NewMockRefreshTokenRepository(ctrl)
				tm.EXPECT().ValidateToken("bad-token").Return(nil, domain.ErrInvalidToken)
				return usecase.NewAuthService(ur, tm, rt)
			},
			in:      usecase.RefreshInput{RefreshToken: "bad-token"},
			wantErr: domain.ErrInvalidToken,
		},
		"UserIDが無効なUUIDの場合はErrInvalidTokenを返す": {
			setup: func(ctrl *gomock.Controller) *usecase.AuthService {
				ur := mock.NewMockUserRepository(ctrl)
				tm := mock.NewMockTokenManager(ctrl)
				rt := mock.NewMockRefreshTokenRepository(ctrl)
				claims := &domainauth.TokenClaims{UserID: "not-a-uuid"}
				tm.EXPECT().ValidateToken("token").Return(claims, nil)
				return usecase.NewAuthService(ur, tm, rt)
			},
			in:      usecase.RefreshInput{RefreshToken: "token"},
			wantErr: domain.ErrInvalidToken,
		},
		"ErrNotFoundはErrInvalidTokenに変換される": {
			setup: func(ctrl *gomock.Controller) *usecase.AuthService {
				ur := mock.NewMockUserRepository(ctrl)
				tm := mock.NewMockTokenManager(ctrl)
				rt := mock.NewMockRefreshTokenRepository(ctrl)
				claims := &domainauth.TokenClaims{UserID: validUserID.String()}
				tm.EXPECT().ValidateToken("token").Return(claims, nil)
				ur.EXPECT().GetByID(gomock.Any(), validUserID).Return(nil, domain.ErrNotFound)
				return usecase.NewAuthService(ur, tm, rt)
			},
			in:      usecase.RefreshInput{RefreshToken: "token"},
			wantErr: domain.ErrInvalidToken,
		},
		"その他のRepoエラーはそのまま伝播する": {
			setup: func(ctrl *gomock.Controller) *usecase.AuthService {
				ur := mock.NewMockUserRepository(ctrl)
				tm := mock.NewMockTokenManager(ctrl)
				rt := mock.NewMockRefreshTokenRepository(ctrl)
				claims := &domainauth.TokenClaims{UserID: validUserID.String()}
				tm.EXPECT().ValidateToken("token").Return(claims, nil)
				ur.EXPECT().GetByID(gomock.Any(), validUserID).Return(nil, errors.New("db error"))
				return usecase.NewAuthService(ur, tm, rt)
			},
			in:      usecase.RefreshInput{RefreshToken: "token"},
			wantErr: domain.ErrInvalidToken, // placeholder: checked via wantErrFn below
		},
		"GenerateTokenPairのエラーが伝播する": {
			setup: func(ctrl *gomock.Controller) *usecase.AuthService {
				ur := mock.NewMockUserRepository(ctrl)
				tm := mock.NewMockTokenManager(ctrl)
				rt := mock.NewMockRefreshTokenRepository(ctrl)
				claims := &domainauth.TokenClaims{UserID: validUserID.String()}
				tm.EXPECT().ValidateToken("token").Return(claims, nil)
				ur.EXPECT().GetByID(gomock.Any(), validUserID).Return(testUser, nil)
				rt.EXPECT().ExistsByHash(gomock.Any(), gomock.Any()).Return(true, nil)
				rt.EXPECT().DeleteByHash(gomock.Any(), gomock.Any()).Return(nil)
				tm.EXPECT().GenerateTokenPair(validUserID).Return(nil, errors.New("token error"))
				return usecase.NewAuthService(ur, tm, rt)
			},
			in: usecase.RefreshInput{RefreshToken: "token"},
		},
		"リフレッシュトークンがDBに存在しない場合はErrInvalidTokenを返す": {
			setup: func(ctrl *gomock.Controller) *usecase.AuthService {
				ur := mock.NewMockUserRepository(ctrl)
				tm := mock.NewMockTokenManager(ctrl)
				rt := mock.NewMockRefreshTokenRepository(ctrl)
				claims := &domainauth.TokenClaims{UserID: validUserID.String()}
				tm.EXPECT().ValidateToken("revoked-token").Return(claims, nil)
				ur.EXPECT().GetByID(gomock.Any(), validUserID).Return(testUser, nil)
				rt.EXPECT().ExistsByHash(gomock.Any(), gomock.Any()).Return(false, nil)
				return usecase.NewAuthService(ur, tm, rt)
			},
			in:      usecase.RefreshInput{RefreshToken: "revoked-token"},
			wantErr: domain.ErrInvalidToken,
		},
	}

	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			svc := c.setup(ctrl)

			out, err := svc.Refresh(t.Context(), c.in)

			// Special case: non-sentinel error check
			if name == "その他のRepoエラーはそのまま伝播する" {
				if err == nil {
					t.Errorf("expected error, got nil")
					return
				}
				if errors.Is(err, domain.ErrInvalidToken) {
					t.Errorf("expected non-ErrInvalidToken error, got ErrInvalidToken")
				}
				return
			}
			if name == "GenerateTokenPairのエラーが伝播する" {
				if err == nil {
					t.Errorf("expected error, got nil")
				}
				return
			}

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
			if out.AccessToken != c.out.AccessToken {
				t.Errorf("AccessToken = %q, want %q", out.AccessToken, c.out.AccessToken)
			}
			if out.RefreshToken != c.out.RefreshToken {
				t.Errorf("RefreshToken = %q, want %q", out.RefreshToken, c.out.RefreshToken)
			}
		})
	}
}

func TestAuthService_GetMe(t *testing.T) {
	t.Parallel()

	userID := uuid.New()
	testUser := newTestUser(userID)

	cases := map[string]struct {
		setup   func(ctrl *gomock.Controller) *usecase.AuthService
		in      usecase.GetMeInput
		out     *usecase.GetMeOutput
		wantErr error
	}{
		"ユーザーを返す": {
			setup: func(ctrl *gomock.Controller) *usecase.AuthService {
				ur := mock.NewMockUserRepository(ctrl)
				tm := mock.NewMockTokenManager(ctrl)
				rt := mock.NewMockRefreshTokenRepository(ctrl)
				ur.EXPECT().GetByID(gomock.Any(), userID).Return(testUser, nil)
				return usecase.NewAuthService(ur, tm, rt)
			},
			in:  usecase.GetMeInput{UserID: userID},
			out: &usecase.GetMeOutput{User: testUser},
		},
		"ErrNotFoundが伝播する": {
			setup: func(ctrl *gomock.Controller) *usecase.AuthService {
				ur := mock.NewMockUserRepository(ctrl)
				tm := mock.NewMockTokenManager(ctrl)
				rt := mock.NewMockRefreshTokenRepository(ctrl)
				ur.EXPECT().GetByID(gomock.Any(), userID).Return(nil, domain.ErrNotFound)
				return usecase.NewAuthService(ur, tm, rt)
			},
			in:      usecase.GetMeInput{UserID: userID},
			wantErr: domain.ErrNotFound,
		},
	}

	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			svc := c.setup(ctrl)

			out, err := svc.GetMe(t.Context(), c.in)

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
			if out.User != c.out.User {
				t.Errorf("returned user does not match expected")
			}
		})
	}
}

