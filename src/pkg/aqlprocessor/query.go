package aqlprocessor

type Query struct {
	Select Select
	From   From
	Where  *Where
	Order  *Order
	Limit  *Limit
}

type Order struct {
	Orders []OrderBy
}

type OrderBy struct {
	Identifier string
	Ordering   OrderingType
}

type OrderingType uint8

const (
	DescendingOrdering OrderingType = iota
	AscendingOrdering
)

type Limit struct {
	Limit  int
	Offset int
}
