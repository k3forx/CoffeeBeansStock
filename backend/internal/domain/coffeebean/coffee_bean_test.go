package coffeebean

import (
	"errors"
	"testing"

	"github.com/google/uuid"

	domain "github.com/k3forx/CoffeeBeansStock/backend/internal/domain"
)

func TestNew_RoastDetailConsistency(t *testing.T) {
	t.Parallel()

	userID := uuid.New()
	stock, _ := NewStock(100)
	detailLight := RoastDetailLight
	detailItalian := RoastDetailItalian

	tests := map[string]struct {
		roastLevel  RoastLevel
		roastDetail *RoastDetail
		wantErr     bool
		wantDetail  *RoastDetail
	}{
		"category_only": {
			roastLevel: RoastShallow,
			wantDetail: nil,
		},
		"consistent_category_and_detail": {
			roastLevel:  RoastShallow,
			roastDetail: &detailLight,
			wantDetail:  &detailLight,
		},
		"inconsistent_category_and_detail": {
			roastLevel:  RoastShallow,
			roastDetail: &detailItalian,
			wantErr:     true,
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			bean, err := New(userID, "Test Bean", tt.roastLevel, tt.roastDetail, nil, nil, stock)

			if tt.wantErr {
				if err == nil {
					t.Errorf("expected error, got nil")
				}
				var ve domain.ValidationErrors
				if err != nil && !errors.As(err, &ve) {
					t.Errorf("expected ValidationErrors, got %T", err)
				}
				return
			}
			if err != nil {
				t.Errorf("unexpected error: %v", err)
				return
			}
			if bean.RoastLevel() != tt.roastLevel {
				t.Errorf("RoastLevel() = %v, want %v", bean.RoastLevel(), tt.roastLevel)
			}
			if tt.wantDetail == nil {
				if bean.RoastDetail() != nil {
					t.Errorf("RoastDetail() = %v, want nil", bean.RoastDetail())
				}
			} else {
				if bean.RoastDetail() == nil || *bean.RoastDetail() != *tt.wantDetail {
					t.Errorf("RoastDetail() = %v, want %v", bean.RoastDetail(), tt.wantDetail)
				}
			}
		})
	}
}

func TestUpdate_RoastDetailConsistency(t *testing.T) {
	t.Parallel()

	userID := uuid.New()
	stock, _ := NewStock(100)
	detailLight := RoastDetailLight
	detailCinnamon := RoastDetailCinnamon
	detailFrench := RoastDetailFrench
	levelDeep := RoastDeep

	tests := map[string]struct {
		initLevel    RoastLevel
		initDetail   *RoastDetail
		updateLevel  *RoastLevel
		updateDetail *RoastDetail
		wantErr      bool
		wantDetail   *RoastDetail
	}{
		"update_detail_consistent": {
			initLevel:    RoastShallow,
			updateDetail: &detailCinnamon,
			wantDetail:   &detailCinnamon,
		},
		"update_detail_inconsistent": {
			initLevel:    RoastShallow,
			updateDetail: &detailFrench,
			wantErr:      true,
		},
		"update_level_makes_existing_detail_inconsistent": {
			initLevel:   RoastShallow,
			initDetail:  &detailLight,
			updateLevel: &levelDeep,
			wantErr:     true,
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			bean, _ := New(userID, "Test Bean", tt.initLevel, tt.initDetail, nil, nil, stock)
			err := bean.Update(nil, nil, tt.updateLevel, tt.updateDetail, nil, nil)

			if tt.wantErr {
				if err == nil {
					t.Errorf("expected error, got nil")
				}
				return
			}
			if err != nil {
				t.Errorf("unexpected error: %v", err)
				return
			}
			if tt.wantDetail == nil {
				if bean.RoastDetail() != nil {
					t.Errorf("RoastDetail() = %v, want nil", bean.RoastDetail())
				}
			} else {
				if bean.RoastDetail() == nil || *bean.RoastDetail() != *tt.wantDetail {
					t.Errorf("RoastDetail() = %v, want %v", bean.RoastDetail(), tt.wantDetail)
				}
			}
		})
	}
}
