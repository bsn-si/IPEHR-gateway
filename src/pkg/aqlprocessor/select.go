package aqlprocessor

type Select struct {
	Distinct    bool
	SelectExprs []SelectExpr
}

func (s *Select) SetCurrentColumnValue(val SelectValuer) {
	s.SelectExprs[len(s.SelectExprs)-1].Value = val
}

type SelectExpr struct {
	Value     SelectValuer
	AliasName string
}

type SelectValuer interface{}

type BaseSelectValue struct {
}

type IdentifiedPathSelectValue struct {
	BaseSelectValue
}

type PrimitiveSelectValue struct {
	BaseSelectValue
	Val Primitive
}

type AggregateFunctionCallSelectValue struct {
	BaseSelectValue
}

type FunctionCallSelectValue struct {
	BaseSelectValue
}
