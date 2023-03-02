package processor

import (
	"fmt"
	"io"

	"github.com/bsn-si/IPEHR-gateway/src/pkg/aql/parser"
	"github.com/bsn-si/IPEHR-gateway/src/pkg/errors"
)

type Where struct {
	IdentifiedExpr *IdentifiedExpr
	Next           []*Where
	OperatorType   OperatorType
	Brackets       bool
}

func (wh *Where) write(w io.Writer) {
	if wh.IdentifiedExpr != nil {
		wh.IdentifiedExpr.write(w)
		return
	}

	if wh.Brackets && len(wh.Next) == 1 {
		fmt.Fprint(w, "(")
		wh.Next[0].write(w)
		fmt.Fprint(w, ")")
		return
	}

	if wh.OperatorType == NOTOperator && len(wh.Next) == 1 {
		fmt.Fprint(w, "NOT ")
		wh.Next[0].write(w)
		return
	}

	if (wh.OperatorType == OROperator || wh.OperatorType == ANDOperator) && len(wh.Next) == 2 {
		wh.Next[0].write(w)
		fmt.Fprintf(w, " %s ", wh.OperatorType)
		wh.Next[1].write(w)
		return
	}

	fmt.Fprintf(w, "%+v ", *wh)
}

type IdentifiedExpr struct {
	Next               *IdentifiedExpr
	IsExists           bool
	IdentifiedPath     *IdentifiedPath
	Terminal           *Terminal
	ComparisonOperator *ComparisionSymbol

	Brackets bool
}

func (ie *IdentifiedExpr) write(w io.Writer) {
	if ie.IsExists && ie.IdentifiedPath != nil {
		fmt.Fprintf(w, "EXISTS ")
		ie.IdentifiedPath.write(w)
		return
	}

	if ie.Brackets && ie.Next != nil {
		fmt.Fprintf(w, "(")
		ie.Next.write(w)
		fmt.Fprintf(w, ")")
		return
	}

	if ie.IdentifiedPath != nil && ie.ComparisonOperator != nil && ie.Terminal != nil {
		ie.IdentifiedPath.write(w)
		fmt.Fprintf(w, " %s ", *ie.ComparisonOperator)
		ie.Terminal.write(w)
		return
	}

	fmt.Fprintf(w, "%+v ", ie)
}

func getWhere(ctx *parser.WhereExprContext) (*Where, error) {
	result := Where{}

	if ctx.IdentifiedExpr() != nil {
		ie, err := getIdentifiedExpr(ctx.IdentifiedExpr().(*parser.IdentifiedExprContext))
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

	result.Brackets = ctx.SYM_LEFT_PAREN() != nil && ctx.SYM_RIGHT_PAREN() != nil

	for _, whereExpr := range ctx.AllWhereExpr() {
		ww, err := getWhere(whereExpr.(*parser.WhereExprContext))
		if err != nil {
			return nil, errors.Wrap(err, "cannot get Where.InnerWhere")
		}

		result.Next = append(result.Next, ww)
	}

	return &result, nil
}

func getIdentifiedExpr(ctx *parser.IdentifiedExprContext) (*IdentifiedExpr, error) {
	result := IdentifiedExpr{}

	if ctx.IdentifiedExpr() != nil && ctx.SYM_LEFT_PAREN() != nil && ctx.SYM_RIGHT_PAREN() != nil {
		result.Brackets = true
		next, err := getIdentifiedExpr(ctx.IdentifiedExpr().(*parser.IdentifiedExprContext))

		if err != nil {
			return nil, errors.Wrap(err, "cannot get IdentifiedExpr.IdentifiedExpr")
		}

		result.Next = next
	}

	if ctx.EXISTS() != nil && ctx.IdentifiedPath() != nil {
		result.IsExists = true

		ip, err := getIdentifiedPath(ctx.IdentifiedPath().(*parser.IdentifiedPathContext))
		if err != nil {
			return nil, errors.Wrap(err, "cannot get IdentifierExpr.IdentifierPath")
		}

		result.IdentifiedPath = &ip
	}

	if ctx.IdentifiedPath() != nil && ctx.COMPARISON_OPERATOR() != nil {
		ip, err := getIdentifiedPath(ctx.IdentifiedPath().(*parser.IdentifiedPathContext))
		if err != nil {
			return nil, errors.Wrap(err, "cannot get IdentifierExpr.IdentifierPath")
		}

		result.IdentifiedPath = &ip

		co, err := getComparisionSimbol(ctx.COMPARISON_OPERATOR())
		if err != nil {
			return nil, errors.Wrap(err, "cannot get IdentifiedExpr.ComparisonOperator")
		}

		result.ComparisonOperator = &co

		terminal, err := getTerminal(ctx.Terminal().(*parser.TerminalContext))
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

func (t *Terminal) write(w io.Writer) {
	if t.Primitive != nil {
		t.Primitive.write(w)
		return
	}

	if t.Parameter != nil {
		fmt.Fprintf(w, "$%s", *t.Parameter)
		return
	}

	if t.IdentifiedPath != nil {
		t.IdentifiedPath.write(w)
		return
	}
}

func getTerminal(ctx *parser.TerminalContext) (*Terminal, error) { //nolint
	t := &Terminal{}

	if ctx.Primitive() != nil {
		p, err := getPrimitive(ctx.Primitive().(*parser.PrimitiveContext))
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
		ip, err := getIdentifiedPath(ctx.IdentifiedPath().(*parser.IdentifiedPathContext))
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
