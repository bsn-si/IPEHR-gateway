package aqlprocessor

import (
	"hms/gateway/pkg/aqlprocessor/aqlparser"
	"hms/gateway/pkg/errors"
)

type Where struct {
	IdentifiedExpr *IdentifiedExpr
	Next           []*Where
	OperatorType   OperatorType
}

type IdentifiedExpr struct {
	Next               *IdentifiedExpr
	IsExists           *bool
	IdentifiedPath     *IdentifiedPath
	Terminal           *Terminal
	ComparisonOperator *ComparisionSymbol
}

func getWhere(ctx *aqlparser.WhereExprContext) (*Where, error) {
	result := Where{}

	if ctx.IdentifiedExpr() != nil {
		ie, err := getIdentifiedExpr(ctx.IdentifiedExpr().(*aqlparser.IdentifiedExprContext))
		if err != nil {
			return nil, errors.Wrap(err, "cannot get Where.IdentifiedExpr")
		}

		result.IdentifiedExpr = ie
	}

	if ctx.NOT() != nil {
		result.OperatorType = NOTOperator
	}

	if ctx.AND() != nil {
		result.OperatorType = ANDOperator
	}

	if ctx.OR() != nil {
		result.OperatorType = OROperator
	}

	for _, whereExpr := range ctx.AllWhereExpr() {
		ww, err := getWhere(whereExpr.(*aqlparser.WhereExprContext))
		if err != nil {
			return nil, errors.Wrap(err, "cannot get Where.InnerWhere")
		}

		result.Next = append(result.Next, ww)
	}

	return &result, nil
}

func getIdentifiedExpr(ctx *aqlparser.IdentifiedExprContext) (*IdentifiedExpr, error) {
	result := IdentifiedExpr{}

	if ctx.IdentifiedExpr() != nil {
		next, err := getIdentifiedExpr(ctx.IdentifiedExpr().(*aqlparser.IdentifiedExprContext))
		if err != nil {
			return nil, errors.Wrap(err, "cannot get IdentifiedExpr.IdentifiedExpr")
		}

		result.Next = next
	}

	if ctx.EXISTS() != nil && ctx.IdentifiedPath() != nil {
		result.IsExists = toRef(true)
		ip := NewIdentifiedPath(ctx.IdentifiedPath().(*aqlparser.IdentifiedPathContext))
		result.IdentifiedPath = &ip
	}

	if ctx.IdentifiedPath() != nil && ctx.COMPARISON_OPERATOR() != nil {
		ip := NewIdentifiedPath(ctx.IdentifiedPath().(*aqlparser.IdentifiedPathContext))
		result.IdentifiedPath = &ip

		co, err := getComparisionSimbol(ctx.COMPARISON_OPERATOR())
		if err != nil {
			return nil, errors.Wrap(err, "cannot get IdentifiedExpr.ComparisonOperator")
		}

		result.ComparisonOperator = &co

		terminal, err := getTerminal(ctx.Terminal().(*aqlparser.TerminalContext))
		if err != nil {
			return nil, errors.Wrap(err, "cannot get IdentifiedExpr.Terminal value")
		}

		result.Terminal = terminal
	}

	return &result, nil
}

type Terminal struct {
	Text string
}

func getTerminal(ctx *aqlparser.TerminalContext) (*Terminal, error) { //nolint
	return &Terminal{
		Text: ctx.GetText(),
	}, nil
}
