package coffeebean

import domain "github.com/k3forx/CoffeeBeansStock/backend/internal/domain"

const MaxStock = 50000

type Stock struct {
	value int32
}

func NewStock(v int32) (Stock, error) {
	if v < 0 || v > MaxStock {
		return Stock{}, &domain.ValidationError{
			Field:   "current_stock",
			Message: "在庫量は0〜50000の範囲で指定してください",
		}
	}
	return Stock{value: v}, nil
}

func ReconstructStock(v int32) Stock { return Stock{value: v} }

func (s Stock) Value() int32 { return s.value }

func (s Stock) CanConsume(qty domain.Quantity) bool {
	return s.value >= qty.Value()
}

func (s Stock) Consume(qty domain.Quantity) Stock {
	return Stock{value: s.value - qty.Value()}
}

func (s Stock) Add(qty domain.Quantity) Stock {
	return Stock{value: s.value + qty.Value()}
}
