package aqlprocessor

import (
	"hms/gateway/pkg/aqlprocessor/aqlparser"
	"log"
	"strconv"

	"github.com/antlr/antlr4/runtime/Go/antlr/v4"
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

// VisitTerminal is called when a terminal node is visited.
func (aql *AQLListener) VisitTerminal(node antlr.TerminalNode) {}

// VisitErrorNode is called when an error node is visited.
func (aql *AQLListener) VisitErrorNode(node antlr.ErrorNode) {}

// EnterEveryRule is called when any rule is entered.
func (aql *AQLListener) EnterEveryRule(ctx antlr.ParserRuleContext) {}

// ExitEveryRule is called when any rule is exited.
func (aql *AQLListener) ExitEveryRule(ctx antlr.ParserRuleContext) {}

// EnterSelectQuery is called when production selectQuery is entered.
func (aql *AQLListener) EnterSelectQuery(ctx *aqlparser.SelectQueryContext) {
	aql.query.Select = Select{}
}

// ExitSelectQuery is called when production selectQuery is exited.
func (aql *AQLListener) ExitSelectQuery(ctx *aqlparser.SelectQueryContext) {
}

// EnterSelectClause is called when production selectClause is entered.
func (aql *AQLListener) EnterSelectClause(ctx *aqlparser.SelectClauseContext) {
}

// ExitSelectClause is called when production selectClause is exited.
func (aql *AQLListener) ExitSelectClause(ctx *aqlparser.SelectClauseContext) {
}

// EnterFromClause is called when production fromClause is entered.
func (aql *AQLListener) EnterFromClause(ctx *aqlparser.FromClauseContext) {
	aql.query.From = From{}
}

// ExitFromClause is called when production fromClause is exited.
func (aql *AQLListener) ExitFromClause(ctx *aqlparser.FromClauseContext) {
}

// EnterWhereClause is called when production whereClause is entered.
func (aql *AQLListener) EnterWhereClause(ctx *aqlparser.WhereClauseContext) {
	aql.query.Where = &Where{}
}

// ExitWhereClause is called when production whereClause is exited.
func (aql *AQLListener) ExitWhereClause(ctx *aqlparser.WhereClauseContext) {}

// EnterOrderByClause is called when production orderByClause is entered.
func (aql *AQLListener) EnterOrderByClause(ctx *aqlparser.OrderByClauseContext) {
	if ctx.IsEmpty() {
		return
	}

	for _, expr := range ctx.AllOrderByExpr() {
		log.Println(expr.GetText())
	}

	aql.query.Order = &Order{}
}

// ExitOrderByClause is called when production orderByClause is exited.
func (aql *AQLListener) ExitOrderByClause(ctx *aqlparser.OrderByClauseContext) {}

// EnterLimitClause is called when production limitClause is entered.
func (aql *AQLListener) EnterLimitClause(ctx *aqlparser.LimitClauseContext) {
	if ctx.IsEmpty() {
		return
	}

	aql.query.Limit = &Limit{}

	if limit := ctx.GetLimit(); limit != nil && limit.GetTokenType() == aqlparser.AqlLexerINTEGER {
		limit, err := strconv.Atoi(limit.GetText())
		if err != nil {
			ctx.SetException(antlr.InputMisMatchException{
				BaseRecognitionException: &antlr.BaseRecognitionException{},
			})
		}

		aql.query.Limit.Limit = limit
	}

	if offset := ctx.GetOffset(); offset != nil && offset.GetTokenType() == aqlparser.AqlLexerINTEGER {
		offset, err := strconv.Atoi(offset.GetText())
		if err != nil {
			ctx.SetException(antlr.InputMisMatchException{
				BaseRecognitionException: &antlr.BaseRecognitionException{},
			})
		}

		aql.query.Limit.Offset = offset
	}
}

// ExitLimitClause is called when production limitClause is exited.
func (aql *AQLListener) ExitLimitClause(ctx *aqlparser.LimitClauseContext) {}

// EnterSelectExpr is called when production selectExpr is entered.
func (aql *AQLListener) EnterSelectExpr(ctx *aqlparser.SelectExprContext) {
	// fmt.Println("select expr")
	// fmt.Println(ctx.GetText())
}

// ExitSelectExpr is called when production selectExpr is exited.
func (aql *AQLListener) ExitSelectExpr(ctx *aqlparser.SelectExprContext) {
}

// EnterFromExpr is called when production fromExpr is entered.
func (aql *AQLListener) EnterFromExpr(ctx *aqlparser.FromExprContext) {}

// ExitFromExpr is called when production fromExpr is exited.
func (aql *AQLListener) ExitFromExpr(ctx *aqlparser.FromExprContext) {}

// EnterWhereExpr is called when production whereExpr is entered.
func (aql *AQLListener) EnterWhereExpr(ctx *aqlparser.WhereExprContext) {}

// ExitWhereExpr is called when production whereExpr is exited.
func (aql *AQLListener) ExitWhereExpr(ctx *aqlparser.WhereExprContext) {}

// EnterOrderByExpr is called when production orderByExpr is entered.
func (aql *AQLListener) EnterOrderByExpr(ctx *aqlparser.OrderByExprContext) {
	order := OrderBy{
		Identifier: ctx.IdentifiedPath().GetText(),
		Ordering:   AscendingOrdering,
	}

	if ctx.ASC() != nil || ctx.ASCENDING() != nil {
		order.Ordering = AscendingOrdering
	} else if ctx.DESC() != nil || ctx.DESCENDING() != nil {
		order.Ordering = DescendingOrdering
	}

	aql.query.Order.Orders = append(aql.query.Order.Orders, order)
}

// ExitOrderByExpr is called when production orderByExpr is exited.
func (aql *AQLListener) ExitOrderByExpr(ctx *aqlparser.OrderByExprContext) {}

// EnterColumnExpr is called when production columnExpr is entered.
func (aql *AQLListener) EnterColumnExpr(ctx *aqlparser.ColumnExprContext) {
}

// ExitColumnExpr is called when production columnExpr is exited.
func (aql *AQLListener) ExitColumnExpr(ctx *aqlparser.ColumnExprContext) {}

// EnterContainsExpr is called when production containsExpr is entered.
func (aql *AQLListener) EnterContainsExpr(ctx *aqlparser.ContainsExprContext) {}

// ExitContainsExpr is called when production containsExpr is exited.
func (aql *AQLListener) ExitContainsExpr(ctx *aqlparser.ContainsExprContext) {}

// EnterIdentifiedExpr is called when production identifiedExpr is entered.
func (aql *AQLListener) EnterIdentifiedExpr(ctx *aqlparser.IdentifiedExprContext) {
}

// ExitIdentifiedExpr is called when production identifiedExpr is exited.
func (aql *AQLListener) ExitIdentifiedExpr(ctx *aqlparser.IdentifiedExprContext) {}

// EnterClassExpression is called when production classExpression is entered.
func (aql *AQLListener) EnterClassExpression(ctx *aqlparser.ClassExpressionContext) {}

// ExitClassExpression is called when production classExpression is exited.
func (aql *AQLListener) ExitClassExpression(ctx *aqlparser.ClassExpressionContext) {}

// EnterVersionClassExpr is called when production versionClassExpr is entered.
func (aql *AQLListener) EnterVersionClassExpr(ctx *aqlparser.VersionClassExprContext) {}

// ExitVersionClassExpr is called when production versionClassExpr is exited.
func (aql *AQLListener) ExitVersionClassExpr(ctx *aqlparser.VersionClassExprContext) {}

// EnterTerminal is called when production terminal is entered.
func (aql *AQLListener) EnterTerminal(ctx *aqlparser.TerminalContext) {}

// ExitTerminal is called when production terminal is exited.
func (aql *AQLListener) ExitTerminal(ctx *aqlparser.TerminalContext) {}

// EnterIdentifiedPath is called when production identifiedPath is entered.
func (aql *AQLListener) EnterIdentifiedPath(ctx *aqlparser.IdentifiedPathContext) {
}

// ExitIdentifiedPath is called when production identifiedPath is exited.
func (aql *AQLListener) ExitIdentifiedPath(ctx *aqlparser.IdentifiedPathContext) {}

// EnterPathPredicate is called when production pathPredicate is entered.
func (aql *AQLListener) EnterPathPredicate(ctx *aqlparser.PathPredicateContext) {}

// ExitPathPredicate is called when production pathPredicate is exited.
func (aql *AQLListener) ExitPathPredicate(ctx *aqlparser.PathPredicateContext) {}

// EnterStandardPredicate is called when production standardPredicate is entered.
func (aql *AQLListener) EnterStandardPredicate(ctx *aqlparser.StandardPredicateContext) {}

// ExitStandardPredicate is called when production standardPredicate is exited.
func (aql *AQLListener) ExitStandardPredicate(ctx *aqlparser.StandardPredicateContext) {}

// EnterArchetypePredicate is called when production archetypePredicate is entered.
func (aql *AQLListener) EnterArchetypePredicate(ctx *aqlparser.ArchetypePredicateContext) {}

// ExitArchetypePredicate is called when production archetypePredicate is exited.
func (aql *AQLListener) ExitArchetypePredicate(ctx *aqlparser.ArchetypePredicateContext) {}

// EnterNodePredicate is called when production nodePredicate is entered.
func (aql *AQLListener) EnterNodePredicate(ctx *aqlparser.NodePredicateContext) {}

// ExitNodePredicate is called when production nodePredicate is exited.
func (aql *AQLListener) ExitNodePredicate(ctx *aqlparser.NodePredicateContext) {}

// EnterVersionPredicate is called when production versionPredicate is entered.
func (aql *AQLListener) EnterVersionPredicate(ctx *aqlparser.VersionPredicateContext) {}

// ExitVersionPredicate is called when production versionPredicate is exited.
func (aql *AQLListener) ExitVersionPredicate(ctx *aqlparser.VersionPredicateContext) {}

// EnterPathPredicateOperand is called when production pathPredicateOperand is entered.
func (aql *AQLListener) EnterPathPredicateOperand(ctx *aqlparser.PathPredicateOperandContext) {}

// ExitPathPredicateOperand is called when production pathPredicateOperand is exited.
func (aql *AQLListener) ExitPathPredicateOperand(ctx *aqlparser.PathPredicateOperandContext) {}

// EnterObjectPath is called when production objectPath is entered.
func (aql *AQLListener) EnterObjectPath(ctx *aqlparser.ObjectPathContext) {}

// ExitObjectPath is called when production objectPath is exited.
func (aql *AQLListener) ExitObjectPath(ctx *aqlparser.ObjectPathContext) {}

// EnterPathPart is called when production pathPart is entered.
func (aql *AQLListener) EnterPathPart(ctx *aqlparser.PathPartContext) {}

// ExitPathPart is called when production pathPart is exited.
func (aql *AQLListener) ExitPathPart(ctx *aqlparser.PathPartContext) {}

// EnterLikeOperand is called when production likeOperand is entered.
func (aql *AQLListener) EnterLikeOperand(ctx *aqlparser.LikeOperandContext) {}

// ExitLikeOperand is called when production likeOperand is exited.
func (aql *AQLListener) ExitLikeOperand(ctx *aqlparser.LikeOperandContext) {}

// EnterMatchesOperand is called when production matchesOperand is entered.
func (aql *AQLListener) EnterMatchesOperand(ctx *aqlparser.MatchesOperandContext) {}

// ExitMatchesOperand is called when production matchesOperand is exited.
func (aql *AQLListener) ExitMatchesOperand(ctx *aqlparser.MatchesOperandContext) {}

// EnterValueListItem is called when production valueListItem is entered.
func (aql *AQLListener) EnterValueListItem(ctx *aqlparser.ValueListItemContext) {}

// ExitValueListItem is called when production valueListItem is exited.
func (aql *AQLListener) ExitValueListItem(ctx *aqlparser.ValueListItemContext) {}

// EnterPrimitive is called when production primitive is entered.
func (aql *AQLListener) EnterPrimitive(ctx *aqlparser.PrimitiveContext) {}

// ExitPrimitive is called when production primitive is exited.
func (aql *AQLListener) ExitPrimitive(ctx *aqlparser.PrimitiveContext) {}

// EnterNumericPrimitive is called when production numericPrimitive is entered.
func (aql *AQLListener) EnterNumericPrimitive(ctx *aqlparser.NumericPrimitiveContext) {}

// ExitNumericPrimitive is called when production numericPrimitive is exited.
func (aql *AQLListener) ExitNumericPrimitive(ctx *aqlparser.NumericPrimitiveContext) {}

// EnterFunctionCall is called when production functionCall is entered.
func (aql *AQLListener) EnterFunctionCall(ctx *aqlparser.FunctionCallContext) {}

// ExitFunctionCall is called when production functionCall is exited.
func (aql *AQLListener) ExitFunctionCall(ctx *aqlparser.FunctionCallContext) {}

// EnterAggregateFunctionCall is called when production aggregateFunctionCall is entered.
func (aql *AQLListener) EnterAggregateFunctionCall(ctx *aqlparser.AggregateFunctionCallContext) {}

// ExitAggregateFunctionCall is called when production aggregateFunctionCall is exited.
func (aql *AQLListener) ExitAggregateFunctionCall(ctx *aqlparser.AggregateFunctionCallContext) {}

// EnterTerminologyFunction is called when production terminologyFunction is entered.
func (aql *AQLListener) EnterTerminologyFunction(ctx *aqlparser.TerminologyFunctionContext) {}

// ExitTerminologyFunction is called when production terminologyFunction is exited.
func (aql *AQLListener) ExitTerminologyFunction(ctx *aqlparser.TerminologyFunctionContext) {}

// EnterTop is called when production top is entered.
func (aql *AQLListener) EnterTop(ctx *aqlparser.TopContext) {}

// ExitTop is called when production top is exited.
func (aql *AQLListener) ExitTop(ctx *aqlparser.TopContext) {}
