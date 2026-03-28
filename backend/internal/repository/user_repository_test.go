package repository

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/k3forx/CoffeeBeansStock/backend/internal/database"
)

func TestUserRepository_CreateAnonymous(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		want database.User
	}{
		"creates anonymous user with default values": {
			want: database.User{
				LowStockThreshold:   100,
				NotificationEnabled: true,
			},
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			ctx := t.Context()
			tx, err := testPool.Begin(ctx)
			if err != nil {
				t.Fatalf("begin tx: %v", err)
			}
			defer func() { _ = tx.Rollback(ctx) }()

			repo := NewUserRepository(testQueries.WithTx(tx))

			got, err := repo.CreateAnonymous(ctx)
			if err != nil {
				t.Fatalf("CreateAnonymous() error: %v", err)
			}

			if !got.ID.Valid {
				t.Fatal("expected valid UUID, got invalid")
			}
			if got.LowStockThreshold != tt.want.LowStockThreshold {
				t.Errorf("LowStockThreshold = %d, want %d", got.LowStockThreshold, tt.want.LowStockThreshold)
			}
			if got.NotificationEnabled != tt.want.NotificationEnabled {
				t.Errorf("NotificationEnabled = %v, want %v", got.NotificationEnabled, tt.want.NotificationEnabled)
			}
			if got.Email.Valid {
				t.Errorf("Email should be null for anonymous user, got %q", got.Email.String)
			}
			if got.Name.Valid {
				t.Errorf("Name should be null for anonymous user, got %q", got.Name.String)
			}
		})
	}
}

func TestUserRepository_GetByID(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		setup   func(t *testing.T, repo UserRepository) uuid.UUID
		wantErr bool
	}{
		"returns user by ID": {
			setup: func(t *testing.T, repo UserRepository) uuid.UUID {
				t.Helper()
				user, err := repo.CreateAnonymous(t.Context())
				if err != nil {
					t.Fatalf("setup CreateAnonymous: %v", err)
				}
				id, err := uuid.FromBytes(user.ID.Bytes[:])
				if err != nil {
					t.Fatalf("setup parse UUID: %v", err)
				}
				return id
			},
			wantErr: false,
		},
		"returns error for non-existent ID": {
			setup: func(t *testing.T, _ UserRepository) uuid.UUID {
				t.Helper()
				return uuid.New()
			},
			wantErr: true,
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			ctx := t.Context()
			tx, err := testPool.Begin(ctx)
			if err != nil {
				t.Fatalf("begin tx: %v", err)
			}
			defer func() { _ = tx.Rollback(ctx) }()

			repo := NewUserRepository(testQueries.WithTx(tx))
			targetID := tt.setup(t, repo)

			got, err := repo.GetByID(ctx, targetID)

			if tt.wantErr {
				if err == nil {
					t.Fatal("GetByID() expected error, got nil")
				}
				return
			}
			if err != nil {
				t.Fatalf("GetByID() unexpected error: %v", err)
			}

			want := database.User{
				ID:                  pgtype.UUID{Bytes: targetID, Valid: true},
				LowStockThreshold:   100,
				NotificationEnabled: true,
			}

			if diff := cmp.Diff(want, got, userCmpOpts()...); diff != "" {
				t.Errorf("GetByID() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func userCmpOpts() cmp.Options {
	ignoreTimestamps := cmp.FilterPath(func(p cmp.Path) bool {
		switch p.String() {
		case "CreatedAt", "UpdatedAt":
			return true
		}
		return false
	}, cmp.Ignore())

	return cmp.Options{ignoreTimestamps}
}
