package aqlprocessor

type Query struct {
	Select Select
	From   From
	Where  *Where
	Order  *Order
	Limit  *Limit

	parameters map[string]*Parameter
}

func (q *Query) addParameter(p *Parameter) {
	if q.parameters == nil {
		q.parameters = map[string]*Parameter{}
	}

	q.parameters[string(*p)] = p
}

func (q *Query) ParametersCount() int {
	return len(q.parameters)
}
