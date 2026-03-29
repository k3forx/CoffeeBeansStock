package domain_test

import (
	"errors"
	"testing"

	domain "github.com/k3forx/CoffeeBeansStock/backend/internal/domain"
)

func TestNewQuantity(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		input   int32
		want    int32
		wantErr bool
	}{
		"正の値で作成できる": {input: 10, want: 10},
		"1で作成できる":     {input: 1, want: 1},
		"0はエラーになる":    {input: 0, wantErr: true},
		"負の値はエラーになる":  {input: -1, wantErr: true},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			got, err := domain.NewQuantity(tt.input)

			if tt.wantErr {
				if err == nil {
					t.Errorf("expected error, got nil")
					return
				}
				var ve *domain.ValidationError
				if !errors.As(err, &ve) {
					t.Errorf("expected ValidationError, got %T: %v", err, err)
				}
				return
			}
			if err != nil {
				t.Errorf("unexpected error: %v", err)
				return
			}
			if got.Value() != tt.want {
				t.Errorf("NewQuantity(%d).Value() = %d, want %d", tt.input, got.Value(), tt.want)
			}
		})
	}
}

func TestReconstructQuantity(t *testing.T) {
	t.Parallel()

	got := domain.ReconstructQuantity(42)
	if got.Value() != 42 {
		t.Errorf("ReconstructQuantity(42).Value() = %d, want 42", got.Value())
	}
}
