package aqlprocessor

type Query struct {
	Select Select
	From   From
	Where  *Where
	Order  *Order
	Limit  *Limit

	Parameters map[string]*Parameter
}

func (q *Query) addParameter(p *Parameter) {
	if q.Parameters == nil {
		q.Parameters = map[string]*Parameter{}
	}

	q.Parameters[string(*p)] = p
}

func (q *Query) ParametersCount() int {
	return len(q.Parameters)
}
