package coffeebean

import (
	"errors"
	"testing"

	domain "github.com/k3forx/CoffeeBeansStock/backend/internal/domain"
)

func TestNewRoastLevel(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		input   string
		want    RoastLevel
		wantErr bool
	}{
		"shallow":     {input: "shallow", want: RoastShallow},
		"medium":      {input: "medium", want: RoastMedium},
		"medium_deep": {input: "medium_deep", want: RoastMediumDeep},
		"deep":        {input: "deep", want: RoastDeep},
		"invalid":     {input: "light", wantErr: true},
		"empty":       {input: "", wantErr: true},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			got, err := NewRoastLevel(tt.input)

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
				t.Errorf("NewRoastLevel() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRoastLevel_ValidDetailFor(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		level  RoastLevel
		detail RoastDetail
		want   bool
	}{
		"shallow_light":         {level: RoastShallow, detail: RoastDetailLight, want: true},
		"shallow_cinnamon":      {level: RoastShallow, detail: RoastDetailCinnamon, want: true},
		"shallow_medium":        {level: RoastShallow, detail: RoastDetailMedium, want: false},
		"medium_medium":         {level: RoastMedium, detail: RoastDetailMedium, want: true},
		"medium_high":           {level: RoastMedium, detail: RoastDetailHigh, want: true},
		"medium_city":           {level: RoastMedium, detail: RoastDetailCity, want: false},
		"medium_deep_city":      {level: RoastMediumDeep, detail: RoastDetailCity, want: true},
		"medium_deep_full_city": {level: RoastMediumDeep, detail: RoastDetailFullCity, want: true},
		"medium_deep_french":    {level: RoastMediumDeep, detail: RoastDetailFrench, want: false},
		"deep_french":           {level: RoastDeep, detail: RoastDetailFrench, want: true},
		"deep_italian":          {level: RoastDeep, detail: RoastDetailItalian, want: true},
		"deep_light":            {level: RoastDeep, detail: RoastDetailLight, want: false},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			got := tt.level.ValidDetailFor(tt.detail)
			if got != tt.want {
				t.Errorf("ValidDetailFor(%s) = %v, want %v", tt.detail, got, tt.want)
			}
		})
	}
}
