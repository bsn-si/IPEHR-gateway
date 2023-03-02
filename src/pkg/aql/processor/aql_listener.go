package processor

import (
	"log"

	"github.com/antlr/antlr4/runtime/Go/antlr/v4"
	"github.com/bsn-si/IPEHR-gateway/src/pkg/aql/parser"
)

type AQLListener struct {
	*parser.BaseAqlParserListener
	query Query
}

func NewAQLListener() *AQLListener {
	return &AQLListener{
		query: Query{},
	}
}

// VisitTerminal is called when a terminal node is visited.
func (aql *AQLListener) VisitTerminal(node antlr.TerminalNode) {
	if p, err := getParameter(node); err == nil {
		aql.query.addParameter(p)
	}
}

// EnterSelectClause is called when production selectClause is entered.
func (aql *AQLListener) EnterSelectClause(ctx *parser.SelectClauseContext) {
	slct, err := getSelect(ctx)
	if err != nil {
		handleError(ctx.GetParser(), ctx.GetStart(), err)
		log.Printf("get Select err: %v", err)
	}

	aql.query.Select = *slct
}

// EnterFromClause is called when production fromClause is entered.
func (aql *AQLListener) EnterFromClause(ctx *parser.FromClauseContext) {
	if ctx.FromExpr() == nil {
		return
	}

	from, err := getFrom(ctx.FromExpr().(*parser.FromExprContext))
	if err != nil {
		handleError(ctx.GetParser(), ctx.GetStart(), err)
		log.Printf("get From err: %v", err)
	}

	aql.query.From = from
}

// EnterWhereClause is called when production whereClause is entered.
func (aql *AQLListener) EnterWhereClause(ctx *parser.WhereClauseContext) {
	if ctx.WhereExpr() == nil {
		return
	}

	where, err := getWhere(ctx.WhereExpr().(*parser.WhereExprContext))
	if err != nil {
		handleError(ctx.GetParser(), ctx.GetStart(), err)
		log.Printf("get Where err: %v", err)
	}

	aql.query.Where = where
}

// EnterOrderByClause is called when production orderByClause is entered.
func (aql *AQLListener) EnterOrderByClause(ctx *parser.OrderByClauseContext) {
	order, err := getOrder(ctx)
	if err != nil {
		handleError(ctx.GetParser(), ctx.GetStart(), err)
		log.Printf("get Order err: %v", err)
	}

	aql.query.Order = order
}

// EnterLimitClause is called when production limitClause is entered.
func (aql *AQLListener) EnterLimitClause(ctx *parser.LimitClauseContext) {
	limit, err := getLimit(ctx)
	if err != nil {
		handleError(ctx.GetParser(), ctx.GetStart(), err)
		log.Printf("get Limit err: %v", err)
	}

	aql.query.Limit = limit
}
