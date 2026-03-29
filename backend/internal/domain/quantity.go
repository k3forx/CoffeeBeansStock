package domain

type Quantity struct {
	value int32
}

func (q Quantity) Value() int32 { return q.value }
