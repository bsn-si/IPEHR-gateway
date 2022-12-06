package aqlprocessor

type Query struct {
	Select Select
	From   From
	Where  *Where
	Order  *Order
	Limit  *Limit
}
