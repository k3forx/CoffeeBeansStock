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

func TestStock_CanConsume(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		stock int32
		qty   int32
		want  bool
	}{
		"sufficient":  {stock: 100, qty: 50, want: true},
		"exact":       {stock: 100, qty: 100, want: true},
		"over":        {stock: 100, qty: 101, want: false},
		"zero_stock":  {stock: 0, qty: 1, want: false},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			s := ReconstructStock(tt.stock)
			qty, _ := domain.NewQuantity(tt.qty)

			got := s.CanConsume(qty)
			if got != tt.want {
				t.Errorf("CanConsume(%d) = %v, want %v", tt.qty, got, tt.want)
			}
		})
	}
}

func TestStock_Consume(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		stock int32
		qty   int32
		want  int32
	}{
		"partial": {stock: 100, qty: 30, want: 70},
		"full":    {stock: 100, qty: 100, want: 0},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			s := ReconstructStock(tt.stock)
			qty, _ := domain.NewQuantity(tt.qty)

			got := s.Consume(qty)
			if got.Value() != tt.want {
				t.Errorf("Consume(%d) = %v, want %v", tt.qty, got.Value(), tt.want)
			}
		})
	}
}

func TestStock_Add(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		stock int32
		qty   int32
		want  int32
	}{
		"add": {stock: 100, qty: 50, want: 150},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			s := ReconstructStock(tt.stock)
			qty, _ := domain.NewQuantity(tt.qty)

			got := s.Add(qty)
			if got.Value() != tt.want {
				t.Errorf("Add(%d) = %v, want %v", tt.qty, got.Value(), tt.want)
			}
		})
	}
}
