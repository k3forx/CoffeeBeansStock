package domain

type Quantity struct {
	value int32
}

// NewQuantity creates a new Quantity with validation (v > 0).
func NewQuantity(v int32) (Quantity, error) {
	if v <= 0 {
		return Quantity{}, &ValidationError{
			Field:   "quantity",
			Message: "数量は1以上で指定してください",
		}
	}
	return Quantity{value: v}, nil
}

// ReconstructQuantity restores a Quantity from persisted data without validation.
func ReconstructQuantity(v int32) Quantity { return Quantity{value: v} }

func (q Quantity) Value() int32 { return q.value }
