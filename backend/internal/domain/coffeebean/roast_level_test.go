package coffeebean

import (
	"errors"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"

	domain "github.com/k3forx/CoffeeBeansStock/backend/internal/domain"
)

func TestNewRoastLevel(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		input   string
		want    RoastLevel
		wantErr bool
	}{
		{name: "shallow", input: "shallow", want: RoastShallow},
		{name: "medium", input: "medium", want: RoastMedium},
		{name: "medium_deep", input: "medium_deep", want: RoastMediumDeep},
		{name: "deep", input: "deep", want: RoastDeep},
		{name: "invalid", input: "light", wantErr: true},
		{name: "empty", input: "", wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got, err := NewRoastLevel(tt.input)
			if tt.wantErr {
				var ve *domain.ValidationError
				if err == nil {
					t.Fatal("expected error, got nil")
				}
				if !errors.As(err, &ve) {
					t.Fatalf("expected ValidationError, got %T", err)
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if diff := cmp.Diff(tt.want, got, cmpopts.IgnoreUnexported(RoastLevel{})); diff != "" {
				t.Errorf("mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func TestRoastLevel_ValidDetailFor(t *testing.T) {
	t.Parallel()

	tests := []struct {
		level  RoastLevel
		detail RoastDetail
		want   bool
	}{
		// shallow
		{RoastShallow, RoastDetailLight, true},
		{RoastShallow, RoastDetailCinnamon, true},
		{RoastShallow, RoastDetailMedium, false},
		// medium
		{RoastMedium, RoastDetailMedium, true},
		{RoastMedium, RoastDetailHigh, true},
		{RoastMedium, RoastDetailCity, false},
		// medium_deep
		{RoastMediumDeep, RoastDetailCity, true},
		{RoastMediumDeep, RoastDetailFullCity, true},
		{RoastMediumDeep, RoastDetailFrench, false},
		// deep
		{RoastDeep, RoastDetailFrench, true},
		{RoastDeep, RoastDetailItalian, true},
		{RoastDeep, RoastDetailLight, false},
	}
	for _, tt := range tests {
		t.Run(tt.level.String()+"_"+tt.detail.String(), func(t *testing.T) {
			t.Parallel()
			got := tt.level.ValidDetailFor(tt.detail)
			if got != tt.want {
				t.Errorf("ValidDetailFor(%s) = %v, want %v", tt.detail, got, tt.want)
			}
		})
	}
}
