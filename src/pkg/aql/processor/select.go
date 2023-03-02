package processor

import (
	"fmt"
	"io"

	"github.com/bsn-si/IPEHR-gateway/src/pkg/aql/parser"
	"github.com/bsn-si/IPEHR-gateway/src/pkg/errors"
)

type Select struct {
	Distinct    bool
	SelectExprs []SelectExpr
}

func (s *Select) write(w io.Writer) {
	if s.Distinct {
		fmt.Fprintln(w, "SELECT DISTINCT")
	} else {
		fmt.Fprintln(w, "SELECT")
	}

	for i, se := range s.SelectExprs {
		fmt.Fprintf(w, "\t%s", se.Path)

		if se.AliasName != "" {
			fmt.Fprintf(w, " AS %s", se.AliasName)
		}

		if i < len(s.SelectExprs)-1 {
			fmt.Fprint(w, ",")
		}

		fmt.Fprintln(w)
	}
}

type SelectExpr struct {
	Value     SelectValuer
	Path      string
	AliasName string
}

type SelectValuer interface{}

type IdentifiedPathSelectValue struct {
	Val IdentifiedPath
}

type PrimitiveSelectValue struct {
	Val Primitive
}

type AggregateFunctionCallSelectValue struct {
}

type FunctionCallSelectValue struct {
}

func getSelect(ctx *parser.SelectClauseContext) (*Select, error) {
	result := Select{}

	if distinct := ctx.DISTINCT(); distinct != nil {
		result.Distinct = true
	}

	// if top := ctx.Top(); top != nil {
	// deprecated
	// }

	result.SelectExprs = make([]SelectExpr, 0, len(ctx.AllSelectExpr()))

	for _, se := range ctx.AllSelectExpr() {
		selectExpr, err := getSelectExpr(se.(*parser.SelectExprContext))
		if err != nil {
			return nil, errors.Wrap(err, "cannot get Select.SelectExpr")
		}

		result.SelectExprs = append(result.SelectExprs, selectExpr)
	}

	return &result, nil
}

func getSelectExpr(ctx *parser.SelectExprContext) (SelectExpr, error) {
	selectExpr := SelectExpr{}

	if ctx.ColumnExpr() != nil {
		selectExpr.Path = ctx.ColumnExpr().GetText()

		columVal, err := getColumnExpr(ctx.ColumnExpr().(*parser.ColumnExprContext))
		if err != nil {
			return SelectExpr{}, errors.Wrap(err, "cannot get SelectExpr.ColumnExpr")
		}

		selectExpr.Value = columVal
	}

	if alias := ctx.GetAliasName(); alias != nil {
		selectExpr.AliasName = alias.GetText()
	}

	return selectExpr, nil
}

func getColumnExpr(ctx *parser.ColumnExprContext) (SelectValuer, error) {
	switch val := ctx.GetChild(0).(type) {
	case *parser.IdentifiedPathContext:
		ip, err := getIdentifiedPath(val)
		if err != nil {
			return nil, errors.Wrap(err, "cannot get ColumnExpr.IdentifierPath")
		}

		ifsv := &IdentifiedPathSelectValue{
			Val: ip,
		}

		return ifsv, nil
	case *parser.PrimitiveContext:
		p, err := getPrimitive(val)
		if err != nil {
			return nil, errors.Wrap(err, "cannot get ColumnExpr.Primitive")
		}

		psv := &PrimitiveSelectValue{
			Val: p,
		}

		return psv, nil
	case *parser.AggregateFunctionCallContext: // nolint
		// selectValue = &AggregateFunctionCallSelectValue{}

		return nil, errors.New("column expr Aggregate Func Call Not implemented")
	case *parser.FunctionCallContext: // nolint
		// selectValue = &FunctionCallSelectValue{}

		return nil, errors.New("column expr Func Call Not implemented")
	default:
		return nil, fmt.Errorf("unexpected column expresion type: %T", val) // nolint
	}
}
