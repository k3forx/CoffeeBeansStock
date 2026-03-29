package usagehistory_test

import (
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/google/uuid"

	domain "github.com/k3forx/CoffeeBeansStock/backend/internal/domain"
	"github.com/k3forx/CoffeeBeansStock/backend/internal/domain/usagehistory"
)

func TestNew(t *testing.T) {
	t.Parallel()

	coffeeBeanID := uuid.New()
	userID := uuid.New()
	usageDate := time.Date(2026, 3, 29, 0, 0, 0, 0, time.UTC)
	qty := domain.ReconstructQuantity(10)
	notes := "test notes"

	got := usagehistory.New(coffeeBeanID, userID, usageDate, qty, &notes)

	if got.ID() == uuid.Nil {
		t.Errorf("ID should not be nil")
	}
	if got.CoffeeBeanID() != coffeeBeanID {
		t.Errorf("CoffeeBeanID() = %v, want %v", got.CoffeeBeanID(), coffeeBeanID)
	}
	if got.UserID() != userID {
		t.Errorf("UserID() = %v, want %v", got.UserID(), userID)
	}
	if got.UsageDate() != usageDate {
		t.Errorf("UsageDate() = %v, want %v", got.UsageDate(), usageDate)
	}
	if got.Quantity().Value() != 10 {
		t.Errorf("Quantity().Value() = %d, want 10", got.Quantity().Value())
	}
	if got.Notes() == nil || *got.Notes() != "test notes" {
		t.Errorf("Notes() = %v, want \"test notes\"", got.Notes())
	}
	if got.CreatedAt().IsZero() {
		t.Errorf("CreatedAt should not be zero")
	}
}

func TestReconstruct(t *testing.T) {
	t.Parallel()

	id := uuid.New()
	coffeeBeanID := uuid.New()
	userID := uuid.New()
	usageDate := time.Date(2026, 3, 29, 0, 0, 0, 0, time.UTC)
	qty := domain.ReconstructQuantity(5)
	createdAt := time.Date(2026, 3, 29, 12, 0, 0, 0, time.UTC)
	notes := "reconstructed"

	got := usagehistory.Reconstruct(id, coffeeBeanID, userID, usageDate, qty, &notes, createdAt)

	want := usagehistory.Reconstruct(id, coffeeBeanID, userID, usageDate, qty, &notes, createdAt)

	if diff := cmp.Diff(want, got,
		cmp.AllowUnexported(usagehistory.UsageHistory{}, domain.Quantity{}),
	); diff != "" {
		t.Errorf("Reconstruct mismatch (-want +got):\n%s", diff)
	}
}

func TestUsageHistory_IsOwnedBy(t *testing.T) {
	t.Parallel()

	userID := uuid.New()
	otherUserID := uuid.New()
	usage := usagehistory.Reconstruct(
		uuid.New(), uuid.New(), userID,
		time.Now(), domain.ReconstructQuantity(5),
		nil, time.Now(),
	)

	tests := map[string]struct {
		checkUserID uuid.UUID
		want        bool
	}{
		"自分のものはtrue":  {checkUserID: userID, want: true},
		"他人のものはfalse": {checkUserID: otherUserID, want: false},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			if got := usage.IsOwnedBy(tt.checkUserID); got != tt.want {
				t.Errorf("IsOwnedBy() = %v, want %v", got, tt.want)
			}
		})
	}
}
