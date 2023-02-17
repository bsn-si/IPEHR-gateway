package aqlprocessor

import (
	"fmt"

	"github.com/bsn-si/IPEHR-gateway/src/pkg/aqlprocessor/aqlparser"
	"github.com/bsn-si/IPEHR-gateway/src/pkg/errors"
	"github.com/vmihailenco/msgpack/v5"
)

type Select struct {
	Distinct    bool
	SelectExprs []SelectExpr
}

type SelectExpr struct {
	Value     SelectValuer
	Path      string
	AliasName string
}

type SelectValueWrap struct {
	Type      uint8
	Value     any
	Path      string
	AliasName string
}

type SelectValueType = uint8

const (
	SelectValueTypeIdentifiedPath SelectValueType = iota
	SelectValueTypePrimitive
	SelectValueAggregateFunction
	SelectValueFunctionCall
)

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

func (se SelectExpr) EncodeMsgpack(enc *msgpack.Encoder) error {
	var _type SelectValueType

	switch se.Value.(type) {
	case *IdentifiedPathSelectValue:
		_type = SelectValueTypeIdentifiedPath
	case *PrimitiveSelectValue:
		_type = SelectValueTypePrimitive
	case *AggregateFunctionCallSelectValue:
		//TODO
	case *FunctionCallSelectValue:
		//TODO
	default:
		return errors.Errorf("Unsupported SelectExpr.Value type: %T", se.Value)
	}

	return enc.Encode(SelectValueWrap{
		Type:      _type,
		Path:      se.Path,
		AliasName: se.AliasName,
		Value:     se.Value,
	})
}

func (se *SelectExpr) UnmarshalMsgpack(data []byte) error {
	var tmp map[string]any

	err := msgpack.Unmarshal(data, &tmp)
	if err != nil {
		panic(err)
	}

	_type, ok := tmp["Type"]
	if !ok {
		return errors.New("Unknown type of SelectExpr")
	}

	wrap := SelectValueWrap{}

	switch _type {
	case SelectValueTypeIdentifiedPath:
		wrap.Value = &IdentifiedPathSelectValue{}
	case SelectValueTypePrimitive:
		wrap.Value = &PrimitiveSelectValue{}
	case SelectValueAggregateFunction:
		wrap.Value = &AggregateFunctionCallSelectValue{}
	case SelectValueFunctionCall:
		wrap.Value = &FunctionCallSelectValue{}
	default:
		return errors.Errorf("Unsupported SelectValueType: %d", _type)
	}

	err = msgpack.Unmarshal(data, &wrap)
	if err != nil {
		return errors.New("SelectExpr unmarshaling error")
	}

	se.Path = wrap.Path
	se.AliasName = wrap.AliasName
	se.Value = wrap.Value

	return nil
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
		selectExpr.Path = ctx.ColumnExpr().GetText()

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
	case *aqlparser.AggregateFunctionCallContext: // nolint
		// selectValue = &AggregateFunctionCallSelectValue{}

		return nil, errors.New("column expr Aggregate Func Call Not implemented")
	case *aqlparser.FunctionCallContext: // nolint
		// selectValue = &FunctionCallSelectValue{}

		return nil, errors.New("column expr Func Call Not implemented")
	default:
		return nil, fmt.Errorf("unexpected column expresion type: %T", val) // nolint
	}
}
