package processor

import (
	"fmt"
	"io"

	"github.com/bsn-si/IPEHR-gateway/src/pkg/aql/parser"
	"github.com/bsn-si/IPEHR-gateway/src/pkg/errors"
)

type Order struct {
	Orders []OrderBy
}

func (o *Order) write(w io.Writer) {
	fmt.Fprint(w, "ORDER BY ")

	for i := range o.Orders {
		if i != 0 {
			fmt.Fprint(w, ",")
		}

		fmt.Fprint(w, " ")
		o.Orders[i].IdentifierPath.write(w)

		if o.Orders[i].Ordering == DescendingOrdering {
			fmt.Fprint(w, " DESC")
		} else if o.Orders[i].Ordering == AscendingOrdering {
			fmt.Fprint(w, " ASC")
		}
	}
}

type OrderBy struct {
	IdentifierPath IdentifiedPath
	Ordering       OrderingType
}

type OrderingType uint8

const (
	NoneOrdering OrderingType = iota
	DescendingOrdering
	AscendingOrdering
)

func getOrder(ctx *parser.OrderByClauseContext) (*Order, error) {
	result := Order{
		Orders: make([]OrderBy, 0, len(ctx.AllOrderByExpr())),
	}

	for _, orderByExpr := range ctx.AllOrderByExpr() {
		orderBy, err := getOrderBy(orderByExpr.(*parser.OrderByExprContext))
		if err != nil {
			return nil, errors.Wrap(err, "cannot get Order.OrderBy")
		}

		result.Orders = append(result.Orders, orderBy)
	}

	return &result, nil
}

func getOrderBy(ctx *parser.OrderByExprContext) (OrderBy, error) { //nolint
	ip, err := getIdentifiedPath(ctx.IdentifiedPath().(*parser.IdentifiedPathContext))
	if err != nil {
		return OrderBy{}, errors.Wrap(err, "cannot get OrderBy.IdentifiedPath")
	}

	orderBy := OrderBy{
		IdentifierPath: ip,
		Ordering:       NoneOrdering,
	}

	if ctx.ASC() != nil || ctx.ASCENDING() != nil {
		orderBy.Ordering = AscendingOrdering
	} else if ctx.DESC() != nil || ctx.DESCENDING() != nil {
		orderBy.Ordering = DescendingOrdering
	}

	return orderBy, nil
}
