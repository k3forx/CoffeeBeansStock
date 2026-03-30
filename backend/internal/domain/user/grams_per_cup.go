package user

import domain "github.com/k3forx/CoffeeBeansStock/backend/internal/domain"

const (
	minGramsPerCup     = 1
	maxGramsPerCup     = 100
	defaultGramsPerCup = 15
)

// GramsPerCup represents the amount of coffee grounds used per cup.
type GramsPerCup struct {
	value int32
}

// NewGramsPerCup creates a GramsPerCup with validation.
func NewGramsPerCup(v int32) (GramsPerCup, error) {
	if v < minGramsPerCup || v > maxGramsPerCup {
		return GramsPerCup{}, &domain.ValidationError{
			Field:   "grams_per_cup",
			Message: "1〜100の範囲で指定してください",
		}
	}
	return GramsPerCup{value: v}, nil
}

// DefaultGramsPerCup returns the default value (15g).
func DefaultGramsPerCup() GramsPerCup {
	return GramsPerCup{value: defaultGramsPerCup}
}

// ReconstructGramsPerCup restores a GramsPerCup from persisted data without validation.
func ReconstructGramsPerCup(v int32) GramsPerCup {
	return GramsPerCup{value: v}
}

func (g GramsPerCup) Value() int32 { return g.value }
