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
