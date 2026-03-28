package coffeebean

import (
	"errors"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"

	domain "github.com/k3forx/CoffeeBeansStock/backend/internal/domain"
)

func TestNewRoastDetail(t *testing.T) {
	t.Parallel()

	validValues := []string{"light", "cinnamon", "medium", "high", "city", "full_city", "french", "italian"}
	for _, v := range validValues {
		t.Run("valid_"+v, func(t *testing.T) {
			t.Parallel()
			got, err := NewRoastDetail(v)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			want := RoastDetail{value: v}
			if diff := cmp.Diff(want, got, cmpopts.IgnoreUnexported(RoastDetail{})); diff != "" {
				t.Errorf("mismatch (-want +got):\n%s", diff)
			}
		})
	}

	t.Run("invalid", func(t *testing.T) {
		t.Parallel()
		_, err := NewRoastDetail("invalid")
		var ve *domain.ValidationError
		if err == nil {
			t.Fatal("expected error, got nil")
		}
		if !errors.As(err, &ve) {
			t.Fatalf("expected ValidationError, got %T", err)
		}
	})
}
