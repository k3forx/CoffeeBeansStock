package coffeebean

import (
	"errors"
	"testing"

	domain "github.com/k3forx/CoffeeBeansStock/backend/internal/domain"
)

func TestNewStock(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		input   int32
		want    int32
		wantErr bool
	}{
		"zero":        {input: 0, want: 0},
		"normal":      {input: 100, want: 100},
		"max":         {input: MaxStock, want: MaxStock},
		"negative":    {input: -1, wantErr: true},
		"over_max":    {input: MaxStock + 1, wantErr: true},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			got, err := NewStock(tt.input)

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
			if got.Value() != tt.want {
				t.Errorf("NewStock() = %v, want %v", got.Value(), tt.want)
			}
		})
	}
}

