package aqlprocessor

import (
	"fmt"
	"hms/gateway/pkg/aqlprocessor/aqlparser"
	"hms/gateway/pkg/errors"
)

type Select struct {
	Distinct    bool
	SelectExprs []SelectExpr
}

type SelectExpr struct {
	Value     SelectValuer
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

func getSelect(ctx *aqlparser.SelectClauseContext) (*Select, error) {
	result := Select{}

	if distinct := ctx.DISTINCT(); distinct != nil {
		result.Distinct = true
	}

	// if top := ctx.Top(); top != nil {
	// deprecated
	// }

	result.SelectExprs = make([]SelectExpr, 0, len(ctx.AllSelectExpr()))

	for _, se := range ctx.AllSelectExpr() {
		selectExpr, err := getSelectExpr(se.(*aqlparser.SelectExprContext))
		if err != nil {
			return nil, errors.Wrap(err, "cannot get Select.SelectExpr")
		}

		result.SelectExprs = append(result.SelectExprs, selectExpr)
	}

	return &result, nil
}

func getSelectExpr(ctx *aqlparser.SelectExprContext) (SelectExpr, error) {
	selectExpr := SelectExpr{}

	if ctx.ColumnExpr() != nil {
		columVal, err := getColumnExpr(ctx.ColumnExpr().(*aqlparser.ColumnExprContext))
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

func getColumnExpr(ctx *aqlparser.ColumnExprContext) (SelectValuer, error) {
	switch val := ctx.GetChild(0).(type) {
	case *aqlparser.IdentifiedPathContext:
		ip, err := getIdentifiedPath(val)
		if err != nil {
			return nil, errors.Wrap(err, "cannot get ColumnExpr.IdentifierPath")
		}

		ifsv := &IdentifiedPathSelectValue{
			Val: ip,
		}

		return ifsv, nil
	case *aqlparser.PrimitiveContext:
		p, err := getPrimitive(val)
		if err != nil {
			return nil, errors.Wrap(err, "cannot get ColumnExpr.Primitive")
		}
		psv := &PrimitiveSelectValue{
			Val: p,
		}

		return psv, nil
	case *aqlparser.AggregateFunctionCallContext:
		// selectValue = &AggregateFunctionCallSelectValue{}

		return nil, errors.New("column expr Aggregate Func Call Not implemented")
	case *aqlparser.FunctionCallContext:
		// selectValue = &FunctionCallSelectValue{}

		return nil, errors.New("column expr Func Call Not implemented")
	default:
		return nil, fmt.Errorf("unexpected column expresion type: %T", val) // nolint
	}
}
