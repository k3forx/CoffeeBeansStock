package coffeebean

import (
	"errors"
	"testing"

	domain "github.com/k3forx/CoffeeBeansStock/backend/internal/domain"
)

func TestNewRoastDetail(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		input   string
		want    RoastDetail
		wantErr bool
	}{
		"valid_light":     {input: "light", want: RoastDetailLight},
		"valid_cinnamon":  {input: "cinnamon", want: RoastDetailCinnamon},
		"valid_medium":    {input: "medium", want: RoastDetailMedium},
		"valid_high":      {input: "high", want: RoastDetailHigh},
		"valid_city":      {input: "city", want: RoastDetailCity},
		"valid_full_city": {input: "full_city", want: RoastDetailFullCity},
		"valid_french":    {input: "french", want: RoastDetailFrench},
		"valid_italian":   {input: "italian", want: RoastDetailItalian},
		"invalid":         {input: "invalid", wantErr: true},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			got, err := NewRoastDetail(tt.input)

			if tt.wantErr {
				if err == nil {
					t.Errorf("expected error, got nil")
				}
				var ve *domain.ValidationError
				if err != nil && !errors.As(err, &ve) {
					t.Errorf("expected ValidationError, got %T", err)
				}
				return
			}
			if err != nil {
				t.Errorf("unexpected error: %v", err)
				return
			}
			if got != tt.want {
				t.Errorf("NewRoastDetail() = %v, want %v", got, tt.want)
			}
		})
	}
}
