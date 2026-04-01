package repository

import (
	"testing"

	"github.com/google/uuid"

	domain "github.com/k3forx/CoffeeBeansStock/backend/internal/domain"
	"github.com/k3forx/CoffeeBeansStock/backend/internal/domain/user"
)

func TestUserRepository_Save(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		wantLowStockThreshold int32
	}{
		"creates anonymous user with default values": {
			wantLowStockThreshold: 100,
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			repo := NewUserRepository(newTestTx(t))

			u := user.NewAnonymousUser()
			err := repo.Save(t.Context(), u)
			if err != nil {
				t.Fatalf("Save() error: %v", err)
			}

			if u.ID() == uuid.Nil {
				t.Fatal("expected valid UUID, got nil")
			}
			if u.LowStockThreshold() != tt.wantLowStockThreshold {
				t.Errorf("LowStockThreshold = %d, want %d", u.LowStockThreshold(), tt.wantLowStockThreshold)
			}
		})
	}
}

func TestUserRepository_GetByID(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		setup   func(t *testing.T, repo user.Repository) uuid.UUID
		wantErr error
	}{
		"returns user by ID": {
			setup: func(t *testing.T, repo user.Repository) uuid.UUID {
				t.Helper()
				u := user.NewAnonymousUser()
				if err := repo.Save(t.Context(), u); err != nil {
					t.Fatalf("setup Save: %v", err)
				}
				return u.ID()
			},
			wantErr: nil,
		},
		"returns ErrNotFound for non-existent ID": {
			setup: func(t *testing.T, _ user.Repository) uuid.UUID {
				t.Helper()
				return uuid.New()
			},
			wantErr: domain.ErrNotFound,
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			repo := NewUserRepository(newTestTx(t))
			targetID := tt.setup(t, repo)

			got, err := repo.GetByID(t.Context(), targetID)

			if tt.wantErr != nil {
				if err == nil {
					t.Fatal("GetByID() expected error, got nil")
				}
				return
			}
			if err != nil {
				t.Fatalf("GetByID() unexpected error: %v", err)
			}

			if got.ID() != targetID {
				t.Errorf("ID = %v, want %v", got.ID(), targetID)
			}
			if got.LowStockThreshold() != 100 {
				t.Errorf("LowStockThreshold = %d, want 100", got.LowStockThreshold())
			}
		})
	}
}
