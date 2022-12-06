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

		ip, err := getIdentifiedPath(ctx.IdentifiedPath().(*aqlparser.IdentifiedPathContext))
		if err != nil {
			return nil, errors.Wrap(err, "cannot get IdentifierExpr.IdentifierPath")
		}

		result.IdentifiedPath = &ip
	}

	if ctx.IdentifiedPath() != nil && ctx.COMPARISON_OPERATOR() != nil {
		ip, err := getIdentifiedPath(ctx.IdentifiedPath().(*aqlparser.IdentifiedPathContext))
		if err != nil {
			return nil, errors.Wrap(err, "cannot get IdentifierExpr.IdentifierPath")
		}

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
	Primitive      *Primitive
	Parameter      *Parameter
	IdentifiedPath *IdentifiedPath
	// FunctionCall *FunctionCall
}

func getTerminal(ctx *aqlparser.TerminalContext) (*Terminal, error) { //nolint
	t := &Terminal{}

	if ctx.Primitive() != nil {
		p, err := getPrimitive(ctx.Primitive().(*aqlparser.PrimitiveContext))
		if err != nil {
			return nil, errors.Wrap(err, "cannot get Terminal.Primitive")
		}

		t.Primitive = &p
	}

	if ctx.PARAMETER() != nil {
		p, err := getParameter(ctx.PARAMETER())
		if err != nil {
			return nil, errors.Wrap(err, "cannot get Terminal.PARAMETER")
		}

		t.Parameter = p
	}

	if ctx.IdentifiedPath() != nil {
		ip, err := getIdentifiedPath(ctx.IdentifiedPath().(*aqlparser.IdentifiedPathContext))
		if err != nil {
			return nil, errors.Wrap(err, "cannot get Terminal.IdentifiedPath")
		}

		t.IdentifiedPath = &ip
	}

	if ctx.FunctionCall() != nil {
		//TODO: implement
		return nil, errors.New("Terminal.FunctionCall not implemented")
	}

	return t, nil
}
