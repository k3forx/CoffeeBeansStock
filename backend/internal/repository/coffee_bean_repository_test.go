package repository

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/uuid"

	domain "github.com/k3forx/CoffeeBeansStock/backend/internal/domain"
	"github.com/k3forx/CoffeeBeansStock/backend/internal/domain/coffeebean"
	"github.com/k3forx/CoffeeBeansStock/backend/internal/domain/user"
)

// newTestBean creates a CoffeeBean domain object for testing.
func newTestBean(t *testing.T, userID uuid.UUID) *coffeebean.CoffeeBean {
	t.Helper()
	stock, _ := coffeebean.NewStock(500)
	bean, err := coffeebean.New(userID, "Test Bean", coffeebean.RoastMedium, nil, nil, nil, stock)
	if err != nil {
		t.Fatalf("setup: new bean: %v", err)
	}
	return bean
}

// newTestBeanFull creates a CoffeeBean with all optional fields set.
func newTestBeanFull(t *testing.T, userID uuid.UUID) *coffeebean.CoffeeBean {
	t.Helper()
	stock, _ := coffeebean.NewStock(1000)
	origin := "Ethiopia"
	notes := "fruity aroma"
	detail := coffeebean.RoastDetailCity
	bean, err := coffeebean.New(userID, "Full Bean", coffeebean.RoastMediumDeep, &detail, &origin, &notes, stock)
	if err != nil {
		t.Fatalf("setup: new full bean: %v", err)
	}
	return bean
}

// coffeeBeanCmpOpts returns cmp options for comparing CoffeeBean.
func coffeeBeanCmpOpts() []cmp.Option {
	return []cmp.Option{
		cmp.AllowUnexported(coffeebean.CoffeeBean{}, coffeebean.Stock{}, coffeebean.RoastLevel{}, coffeebean.RoastDetail{}),
		cmp.FilterPath(func(p cmp.Path) bool {
			switch p.String() {
			case "CreatedAt", "UpdatedAt", "createdAt", "updatedAt":
				return true
			}
			return false
		}, cmp.Ignore()),
	}
}

// setupBeanRepo creates a repo and saves a user, returning both.
func setupBeanRepo(t *testing.T) (coffeebean.Repository, uuid.UUID) {
	t.Helper()
	tx := newTestTx(t)
	userRepo := NewUserRepository(tx)
	u := user.NewAnonymousUser()
	if err := userRepo.Save(t.Context(), u); err != nil {
		t.Fatalf("setup: save user: %v", err)
	}
	return NewCoffeeBeanRepository(tx), u.ID()
}

func TestCoffeeBeanRepository_Save(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		makeBean func(t *testing.T, userID uuid.UUID) *coffeebean.CoffeeBean
	}{
		"saves bean with required fields only": {
			makeBean: newTestBean,
		},
		"saves bean with all optional fields": {
			makeBean: newTestBeanFull,
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			repo, userID := setupBeanRepo(t)
			bean := tt.makeBean(t, userID)

			err := repo.Save(t.Context(), bean)

			if err != nil {
				t.Errorf("Save() error: %v", err)
				return
			}

			got, err := repo.GetByID(t.Context(), bean.ID())
			if err != nil {
				t.Errorf("GetByID() after Save error: %v", err)
				return
			}
			if diff := cmp.Diff(bean, got, coffeeBeanCmpOpts()...); diff != "" {
				t.Errorf("mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func TestCoffeeBeanRepository_GetByID(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		setup   func(t *testing.T, repo coffeebean.Repository, userID uuid.UUID) uuid.UUID
		wantErr error
	}{
		"returns bean by ID": {
			setup: func(t *testing.T, repo coffeebean.Repository, userID uuid.UUID) uuid.UUID {
				t.Helper()
				bean := newTestBean(t, userID)
				if err := repo.Save(t.Context(), bean); err != nil {
					t.Fatalf("setup: save bean: %v", err)
				}
				return bean.ID()
			},
		},
		"returns ErrNotFound for non-existent ID": {
			setup: func(t *testing.T, _ coffeebean.Repository, _ uuid.UUID) uuid.UUID {
				t.Helper()
				return uuid.New()
			},
			wantErr: domain.ErrNotFound,
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			repo, userID := setupBeanRepo(t)
			targetID := tt.setup(t, repo, userID)

			got, err := repo.GetByID(t.Context(), targetID)

			if tt.wantErr != nil {
				if err == nil {
					t.Errorf("expected error %v, got nil", tt.wantErr)
				}
				return
			}
			if err != nil {
				t.Errorf("unexpected error: %v", err)
				return
			}
			if got.ID() != targetID {
				t.Errorf("ID = %v, want %v", got.ID(), targetID)
			}
		})
	}
}

func TestCoffeeBeanRepository_GetByIDForUpdate(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		setup   func(t *testing.T, repo coffeebean.Repository, userID uuid.UUID) uuid.UUID
		wantErr error
	}{
		"returns bean by ID with row lock": {
			setup: func(t *testing.T, repo coffeebean.Repository, userID uuid.UUID) uuid.UUID {
				t.Helper()
				bean := newTestBean(t, userID)
				if err := repo.Save(t.Context(), bean); err != nil {
					t.Fatalf("setup: save bean: %v", err)
				}
				return bean.ID()
			},
		},
		"returns ErrNotFound for non-existent ID": {
			setup: func(t *testing.T, _ coffeebean.Repository, _ uuid.UUID) uuid.UUID {
				t.Helper()
				return uuid.New()
			},
			wantErr: domain.ErrNotFound,
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			repo, userID := setupBeanRepo(t)
			targetID := tt.setup(t, repo, userID)

			got, err := repo.GetByIDForUpdate(t.Context(), targetID)

			if tt.wantErr != nil {
				if err == nil {
					t.Errorf("expected error %v, got nil", tt.wantErr)
				}
				return
			}
			if err != nil {
				t.Errorf("unexpected error: %v", err)
				return
			}
			if got.ID() != targetID {
				t.Errorf("ID = %v, want %v", got.ID(), targetID)
			}
		})
	}
}

func TestCoffeeBeanRepository_ListByUserID(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		seedCount int
		limit     int32
		offset    int32
		wantCount int
	}{
		"returns all beans for user": {
			seedCount: 3,
			limit:     10,
			offset:    0,
			wantCount: 3,
		},
		"returns empty list when no beans": {
			seedCount: 0,
			limit:     10,
			offset:    0,
			wantCount: 0,
		},
		"respects limit": {
			seedCount: 5,
			limit:     2,
			offset:    0,
			wantCount: 2,
		},
		"respects offset": {
			seedCount: 5,
			limit:     10,
			offset:    3,
			wantCount: 2,
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			repo, userID := setupBeanRepo(t)
			for i := 0; i < tt.seedCount; i++ {
				stock, _ := coffeebean.NewStock(int32(100 * (i + 1)))
				bean, err := coffeebean.New(userID, "Bean "+string(rune('A'+i)), coffeebean.RoastMedium, nil, nil, nil, stock)
				if err != nil {
					t.Fatalf("setup: new bean: %v", err)
				}
				if err := repo.Save(t.Context(), bean); err != nil {
					t.Fatalf("setup: save bean: %v", err)
				}
			}

			got, err := repo.ListByUserID(t.Context(), userID, tt.limit, tt.offset)

			if err != nil {
				t.Errorf("unexpected error: %v", err)
				return
			}
			if len(got) != tt.wantCount {
				t.Errorf("count = %d, want %d", len(got), tt.wantCount)
			}
		})
	}
}

func TestCoffeeBeanRepository_CountByUserID(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		seedCount int
		wantCount int64
	}{
		"counts beans for user": {
			seedCount: 3,
			wantCount: 3,
		},
		"returns zero when no beans": {
			seedCount: 0,
			wantCount: 0,
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			repo, userID := setupBeanRepo(t)
			for i := 0; i < tt.seedCount; i++ {
				stock, _ := coffeebean.NewStock(100)
				bean, err := coffeebean.New(userID, "Bean "+string(rune('A'+i)), coffeebean.RoastMedium, nil, nil, nil, stock)
				if err != nil {
					t.Fatalf("setup: new bean: %v", err)
				}
				if err := repo.Save(t.Context(), bean); err != nil {
					t.Fatalf("setup: save bean: %v", err)
				}
			}

			got, err := repo.CountByUserID(t.Context(), userID)

			if err != nil {
				t.Errorf("unexpected error: %v", err)
				return
			}
			if got != tt.wantCount {
				t.Errorf("count = %d, want %d", got, tt.wantCount)
			}
		})
	}
}

func TestCoffeeBeanRepository_Update(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		setup   func(t *testing.T, repo coffeebean.Repository, userID uuid.UUID) *coffeebean.CoffeeBean
		wantErr error
	}{
		"updates bean fields": {
			setup: func(t *testing.T, repo coffeebean.Repository, userID uuid.UUID) *coffeebean.CoffeeBean {
				t.Helper()
				bean := newTestBean(t, userID)
				if err := repo.Save(t.Context(), bean); err != nil {
					t.Fatalf("setup: save bean: %v", err)
				}
				newName := "Updated Bean"
				newOrigin := "Colombia"
				newLevel := coffeebean.RoastDeep
				newDetail := coffeebean.RoastDetailFrench
				newNotes := "bold flavor"
				newStock, _ := coffeebean.NewStock(999)
				if err := bean.Update(&newName, &newOrigin, &newLevel, &newDetail, &newNotes, &newStock); err != nil {
					t.Fatalf("setup: update domain: %v", err)
				}
				return bean
			},
		},
		"returns ErrNotFound for non-existent bean": {
			setup: func(t *testing.T, _ coffeebean.Repository, userID uuid.UUID) *coffeebean.CoffeeBean {
				t.Helper()
				return newTestBean(t, userID)
			},
			wantErr: domain.ErrNotFound,
		},
		"returns ErrNotFound when user does not own the bean": {
			setup: func(t *testing.T, repo coffeebean.Repository, userID uuid.UUID) *coffeebean.CoffeeBean {
				t.Helper()
				bean := newTestBean(t, userID)
				if err := repo.Save(t.Context(), bean); err != nil {
					t.Fatalf("setup: save bean: %v", err)
				}
				// Reconstruct with a different userID to simulate wrong owner
				otherUserID := uuid.New()
				return coffeebean.Reconstruct(
					bean.ID(), otherUserID, bean.Name(), bean.Origin(),
					bean.RoastLevel(), bean.RoastDetail(), bean.CurrentStock(),
					bean.Notes(), bean.CreatedAt(), bean.UpdatedAt(),
				)
			},
			wantErr: domain.ErrNotFound,
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			repo, userID := setupBeanRepo(t)
			bean := tt.setup(t, repo, userID)

			err := repo.Update(t.Context(), bean)

			if tt.wantErr != nil {
				if err == nil {
					t.Errorf("expected error %v, got nil", tt.wantErr)
				}
				return
			}
			if err != nil {
				t.Errorf("unexpected error: %v", err)
				return
			}

			got, err := repo.GetByID(t.Context(), bean.ID())
			if err != nil {
				t.Errorf("GetByID() after Update error: %v", err)
				return
			}
			if diff := cmp.Diff(bean, got, coffeeBeanCmpOpts()...); diff != "" {
				t.Errorf("mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func TestCoffeeBeanRepository_SoftDelete(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		setup   func(t *testing.T, repo coffeebean.Repository, userID uuid.UUID) (beanID uuid.UUID, deleteUserID uuid.UUID)
		wantErr error
	}{
		"soft deletes a bean": {
			setup: func(t *testing.T, repo coffeebean.Repository, userID uuid.UUID) (uuid.UUID, uuid.UUID) {
				t.Helper()
				bean := newTestBean(t, userID)
				if err := repo.Save(t.Context(), bean); err != nil {
					t.Fatalf("setup: save bean: %v", err)
				}
				return bean.ID(), userID
			},
		},
		"does not error when bean does not exist": {
			setup: func(t *testing.T, _ coffeebean.Repository, userID uuid.UUID) (uuid.UUID, uuid.UUID) {
				t.Helper()
				return uuid.New(), userID
			},
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			repo, userID := setupBeanRepo(t)
			beanID, deleteUserID := tt.setup(t, repo, userID)

			err := repo.SoftDelete(t.Context(), beanID, deleteUserID)

			if tt.wantErr != nil {
				if err == nil {
					t.Errorf("expected error %v, got nil", tt.wantErr)
				}
				return
			}
			if err != nil {
				t.Errorf("unexpected error: %v", err)
				return
			}

			// Verify the bean is no longer found
			_, err = repo.GetByID(t.Context(), beanID)
			if err != domain.ErrNotFound {
				t.Errorf("GetByID() after SoftDelete: expected ErrNotFound, got %v", err)
			}
		})
	}
}

func TestCoffeeBeanRepository_CountByUserID_excludes_deleted(t *testing.T) {
	t.Parallel()

	repo, userID := setupBeanRepo(t)

	// Create 3 beans, soft-delete 1
	var deleteID uuid.UUID
	for i := range 3 {
		stock, _ := coffeebean.NewStock(100)
		bean, err := coffeebean.New(userID, "Bean "+string(rune('A'+i)), coffeebean.RoastMedium, nil, nil, nil, stock)
		if err != nil {
			t.Fatalf("setup: new bean: %v", err)
		}
		if err := repo.Save(t.Context(), bean); err != nil {
			t.Fatalf("setup: save bean: %v", err)
		}
		if i == 0 {
			deleteID = bean.ID()
		}
	}
	if err := repo.SoftDelete(t.Context(), deleteID, userID); err != nil {
		t.Fatalf("setup: soft delete: %v", err)
	}

	got, err := repo.CountByUserID(t.Context(), userID)

	if err != nil {
		t.Errorf("unexpected error: %v", err)
		return
	}
	if got != 2 {
		t.Errorf("count = %d, want 2 (one bean was soft-deleted)", got)
	}
}
