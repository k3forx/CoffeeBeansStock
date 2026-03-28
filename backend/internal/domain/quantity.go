package domain

type Quantity struct {
	value int32
}

func NewQuantity(v int32) (Quantity, error) {
	if v <= 0 {
		return Quantity{}, &ValidationError{Field: "quantity", Message: "数量は1以上にしてください"}
	}
	return Quantity{value: v}, nil
}

func (q Quantity) Value() int32 { return q.value }
