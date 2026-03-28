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

	t.Run("category_only", func(t *testing.T) {
		t.Parallel()
		bean, err := New(userID, "Test Bean", RoastShallow, nil, nil, nil, stock)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if bean.RoastLevel() != RoastShallow {
			t.Errorf("RoastLevel = %v, want shallow", bean.RoastLevel())
		}
		if bean.RoastDetail() != nil {
			t.Errorf("RoastDetail = %v, want nil", bean.RoastDetail())
		}
	})

	t.Run("consistent_category_and_detail", func(t *testing.T) {
		t.Parallel()
		detail := RoastDetailLight
		bean, err := New(userID, "Test Bean", RoastShallow, &detail, nil, nil, stock)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if bean.RoastDetail() == nil || *bean.RoastDetail() != RoastDetailLight {
			t.Errorf("RoastDetail = %v, want light", bean.RoastDetail())
		}
	})

	t.Run("inconsistent_category_and_detail", func(t *testing.T) {
		t.Parallel()
		detail := RoastDetailItalian
		_, err := New(userID, "Test Bean", RoastShallow, &detail, nil, nil, stock)
		if err == nil {
			t.Fatal("expected error for inconsistent roast level and detail")
		}
		var ve domain.ValidationErrors
		if !errors.As(err, &ve) {
			t.Fatalf("expected ValidationErrors, got %T", err)
		}
	})
}

func TestUpdate_RoastDetailConsistency(t *testing.T) {
	t.Parallel()

	userID := uuid.New()
	stock, _ := NewStock(100)

	t.Run("update_detail_consistent", func(t *testing.T) {
		t.Parallel()
		bean, _ := New(userID, "Test Bean", RoastShallow, nil, nil, nil, stock)
		detail := RoastDetailCinnamon
		err := bean.Update(nil, nil, nil, &detail, nil, nil)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if bean.RoastDetail() == nil || *bean.RoastDetail() != RoastDetailCinnamon {
			t.Errorf("RoastDetail = %v, want cinnamon", bean.RoastDetail())
		}
	})

	t.Run("update_detail_inconsistent", func(t *testing.T) {
		t.Parallel()
		bean, _ := New(userID, "Test Bean", RoastShallow, nil, nil, nil, stock)
		detail := RoastDetailFrench
		err := bean.Update(nil, nil, nil, &detail, nil, nil)
		if err == nil {
			t.Fatal("expected error for inconsistent update")
		}
	})

	t.Run("update_level_makes_existing_detail_inconsistent", func(t *testing.T) {
		t.Parallel()
		detail := RoastDetailLight
		bean, _ := New(userID, "Test Bean", RoastShallow, &detail, nil, nil, stock)
		newLevel := RoastDeep
		err := bean.Update(nil, nil, &newLevel, nil, nil, nil)
		if err == nil {
			t.Fatal("expected error when changing level makes existing detail inconsistent")
		}
	})
}
