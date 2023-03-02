package processor

import (
	"bytes"
	"fmt"
)

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

func (q *Query) String() string {
	buffer := &bytes.Buffer{}
	q.Select.write(buffer)
	q.From.write(buffer)

	if q.Where != nil {
		fmt.Fprintf(buffer, "\nWHERE ")
		q.Where.write(buffer)
	}

	if q.Order != nil {
		fmt.Fprintln(buffer)
		q.Order.write(buffer)
	}

	if q.Limit != nil {
		fmt.Fprintln(buffer)
		q.Limit.write(buffer)
	}

	return buffer.String()
}
