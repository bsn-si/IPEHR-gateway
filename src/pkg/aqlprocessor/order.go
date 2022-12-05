package aqlprocessor

import (
	"hms/gateway/pkg/aqlprocessor/aqlparser"
	"hms/gateway/pkg/errors"
)

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

func getOrder(ctx *aqlparser.OrderByClauseContext) (*Order, error) {
	result := Order{
		Orders: make([]OrderBy, 0, len(ctx.AllOrderByExpr())),
	}

	for _, orderByExpr := range ctx.AllOrderByExpr() {
		orderBy, err := getOrderBy(orderByExpr.(*aqlparser.OrderByExprContext))
		if err != nil {
			return nil, errors.Wrap(err, "cannot get Order.OrderBy")
		}

		result.Orders = append(result.Orders, orderBy)
	}

	return &result, nil
}

func getOrderBy(ctx *aqlparser.OrderByExprContext) (OrderBy, error) { //nolint
	orderBy := OrderBy{
		Identifier: ctx.IdentifiedPath().GetText(),
		Ordering:   AscendingOrdering,
	}

	if ctx.ASC() != nil || ctx.ASCENDING() != nil {
		orderBy.Ordering = AscendingOrdering
	} else if ctx.DESC() != nil || ctx.DESCENDING() != nil {
		orderBy.Ordering = DescendingOrdering
	}

	return orderBy, nil
}
