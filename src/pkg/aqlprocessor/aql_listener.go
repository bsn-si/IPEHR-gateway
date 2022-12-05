package aqlprocessor

import (
	"log"

	"hms/gateway/pkg/aqlprocessor/aqlparser"
)

type AQLListener struct {
	*aqlparser.BaseAqlParserListener
	query Query
}

func NewAQLListener() *AQLListener {
	return &AQLListener{
		query: Query{},
	}
}

// EnterSelectClause is called when production selectClause is entered.
func (aql *AQLListener) EnterSelectClause(ctx *aqlparser.SelectClauseContext) {
	slct, err := getSelect(ctx)
	if err != nil {
		log.Fatalf("get Select err: %v", err)
	}

	aql.query.Select = *slct
}

// EnterFromClause is called when production fromClause is entered.
func (aql *AQLListener) EnterFromClause(ctx *aqlparser.FromClauseContext) {
	if ctx.FromExpr() == nil {
		return
	}

	from, err := getFrom(ctx.FromExpr().(*aqlparser.FromExprContext))
	if err != nil {
		log.Fatalf("get From err: %v", err)
	}

	aql.query.From = from
}

// EnterWhereClause is called when production whereClause is entered.
func (aql *AQLListener) EnterWhereClause(ctx *aqlparser.WhereClauseContext) {
	if ctx.WhereExpr() == nil {
		return
	}

	where, err := getWhere(ctx.WhereExpr().(*aqlparser.WhereExprContext))
	if err != nil {
		log.Fatalf("get Where err: %v", err)
	}

	aql.query.Where = where
}

// EnterOrderByClause is called when production orderByClause is entered.
func (aql *AQLListener) EnterOrderByClause(ctx *aqlparser.OrderByClauseContext) {
	order, err := getOrder(ctx)
	if err != nil {
		log.Fatalf("get Order err: %v", err)
	}

	aql.query.Order = order
}

// EnterLimitClause is called when production limitClause is entered.
func (aql *AQLListener) EnterLimitClause(ctx *aqlparser.LimitClauseContext) {
	limit, err := getLimit(ctx)
	if err != nil {
		log.Fatalf("get Limit err: %v", err)
	}

	aql.query.Limit = limit
}
