// Code generated from parser.g4 by ANTLR 4.12.0. DO NOT EDIT.

package parser // AqlParser
import "github.com/antlr/antlr4/runtime/Go/antlr/v4"

// BaseAqlParserListener is a complete listener for a parse tree produced by parser.
type BaseAqlParserListener struct{}

var _ AqlParserListener = &BaseAqlParserListener{}

// VisitTerminal is called when a terminal node is visited.
func (s *BaseAqlParserListener) VisitTerminal(node antlr.TerminalNode) {}

// VisitErrorNode is called when an error node is visited.
func (s *BaseAqlParserListener) VisitErrorNode(node antlr.ErrorNode) {}

// EnterEveryRule is called when any rule is entered.
func (s *BaseAqlParserListener) EnterEveryRule(ctx antlr.ParserRuleContext) {}

// ExitEveryRule is called when any rule is exited.
func (s *BaseAqlParserListener) ExitEveryRule(ctx antlr.ParserRuleContext) {}

// EnterSelectQuery is called when production selectQuery is entered.
func (s *BaseAqlParserListener) EnterSelectQuery(ctx *SelectQueryContext) {}

// ExitSelectQuery is called when production selectQuery is exited.
func (s *BaseAqlParserListener) ExitSelectQuery(ctx *SelectQueryContext) {}

// EnterSelectClause is called when production selectClause is entered.
func (s *BaseAqlParserListener) EnterSelectClause(ctx *SelectClauseContext) {}

// ExitSelectClause is called when production selectClause is exited.
func (s *BaseAqlParserListener) ExitSelectClause(ctx *SelectClauseContext) {}

// EnterFromClause is called when production fromClause is entered.
func (s *BaseAqlParserListener) EnterFromClause(ctx *FromClauseContext) {}

// ExitFromClause is called when production fromClause is exited.
func (s *BaseAqlParserListener) ExitFromClause(ctx *FromClauseContext) {}

// EnterWhereClause is called when production whereClause is entered.
func (s *BaseAqlParserListener) EnterWhereClause(ctx *WhereClauseContext) {}

// ExitWhereClause is called when production whereClause is exited.
func (s *BaseAqlParserListener) ExitWhereClause(ctx *WhereClauseContext) {}

// EnterOrderByClause is called when production orderByClause is entered.
func (s *BaseAqlParserListener) EnterOrderByClause(ctx *OrderByClauseContext) {}

// ExitOrderByClause is called when production orderByClause is exited.
func (s *BaseAqlParserListener) ExitOrderByClause(ctx *OrderByClauseContext) {}

// EnterLimitClause is called when production limitClause is entered.
func (s *BaseAqlParserListener) EnterLimitClause(ctx *LimitClauseContext) {}

// ExitLimitClause is called when production limitClause is exited.
func (s *BaseAqlParserListener) ExitLimitClause(ctx *LimitClauseContext) {}

// EnterSelectExpr is called when production selectExpr is entered.
func (s *BaseAqlParserListener) EnterSelectExpr(ctx *SelectExprContext) {}

// ExitSelectExpr is called when production selectExpr is exited.
func (s *BaseAqlParserListener) ExitSelectExpr(ctx *SelectExprContext) {}

// EnterFromExpr is called when production fromExpr is entered.
func (s *BaseAqlParserListener) EnterFromExpr(ctx *FromExprContext) {}

// ExitFromExpr is called when production fromExpr is exited.
func (s *BaseAqlParserListener) ExitFromExpr(ctx *FromExprContext) {}

// EnterWhereExpr is called when production whereExpr is entered.
func (s *BaseAqlParserListener) EnterWhereExpr(ctx *WhereExprContext) {}

// ExitWhereExpr is called when production whereExpr is exited.
func (s *BaseAqlParserListener) ExitWhereExpr(ctx *WhereExprContext) {}

// EnterOrderByExpr is called when production orderByExpr is entered.
func (s *BaseAqlParserListener) EnterOrderByExpr(ctx *OrderByExprContext) {}

// ExitOrderByExpr is called when production orderByExpr is exited.
func (s *BaseAqlParserListener) ExitOrderByExpr(ctx *OrderByExprContext) {}

// EnterColumnExpr is called when production columnExpr is entered.
func (s *BaseAqlParserListener) EnterColumnExpr(ctx *ColumnExprContext) {}

// ExitColumnExpr is called when production columnExpr is exited.
func (s *BaseAqlParserListener) ExitColumnExpr(ctx *ColumnExprContext) {}

// EnterContainsExpr is called when production containsExpr is entered.
func (s *BaseAqlParserListener) EnterContainsExpr(ctx *ContainsExprContext) {}

// ExitContainsExpr is called when production containsExpr is exited.
func (s *BaseAqlParserListener) ExitContainsExpr(ctx *ContainsExprContext) {}

// EnterIdentifiedExpr is called when production identifiedExpr is entered.
func (s *BaseAqlParserListener) EnterIdentifiedExpr(ctx *IdentifiedExprContext) {}

// ExitIdentifiedExpr is called when production identifiedExpr is exited.
func (s *BaseAqlParserListener) ExitIdentifiedExpr(ctx *IdentifiedExprContext) {}

// EnterClassExpression is called when production classExpression is entered.
func (s *BaseAqlParserListener) EnterClassExpression(ctx *ClassExpressionContext) {}

// ExitClassExpression is called when production classExpression is exited.
func (s *BaseAqlParserListener) ExitClassExpression(ctx *ClassExpressionContext) {}

// EnterVersionClassExpr is called when production versionClassExpr is entered.
func (s *BaseAqlParserListener) EnterVersionClassExpr(ctx *VersionClassExprContext) {}

// ExitVersionClassExpr is called when production versionClassExpr is exited.
func (s *BaseAqlParserListener) ExitVersionClassExpr(ctx *VersionClassExprContext) {}

// EnterTerminal is called when production terminal is entered.
func (s *BaseAqlParserListener) EnterTerminal(ctx *TerminalContext) {}

// ExitTerminal is called when production terminal is exited.
func (s *BaseAqlParserListener) ExitTerminal(ctx *TerminalContext) {}

// EnterIdentifiedPath is called when production identifiedPath is entered.
func (s *BaseAqlParserListener) EnterIdentifiedPath(ctx *IdentifiedPathContext) {}

// ExitIdentifiedPath is called when production identifiedPath is exited.
func (s *BaseAqlParserListener) ExitIdentifiedPath(ctx *IdentifiedPathContext) {}

// EnterPathPredicate is called when production pathPredicate is entered.
func (s *BaseAqlParserListener) EnterPathPredicate(ctx *PathPredicateContext) {}

// ExitPathPredicate is called when production pathPredicate is exited.
func (s *BaseAqlParserListener) ExitPathPredicate(ctx *PathPredicateContext) {}

// EnterStandardPredicate is called when production standardPredicate is entered.
func (s *BaseAqlParserListener) EnterStandardPredicate(ctx *StandardPredicateContext) {}

// ExitStandardPredicate is called when production standardPredicate is exited.
func (s *BaseAqlParserListener) ExitStandardPredicate(ctx *StandardPredicateContext) {}

// EnterArchetypePredicate is called when production archetypePredicate is entered.
func (s *BaseAqlParserListener) EnterArchetypePredicate(ctx *ArchetypePredicateContext) {}

// ExitArchetypePredicate is called when production archetypePredicate is exited.
func (s *BaseAqlParserListener) ExitArchetypePredicate(ctx *ArchetypePredicateContext) {}

// EnterNodePredicate is called when production nodePredicate is entered.
func (s *BaseAqlParserListener) EnterNodePredicate(ctx *NodePredicateContext) {}

// ExitNodePredicate is called when production nodePredicate is exited.
func (s *BaseAqlParserListener) ExitNodePredicate(ctx *NodePredicateContext) {}

// EnterNodePredicateAdditionalData is called when production nodePredicateAdditionalData is entered.
func (s *BaseAqlParserListener) EnterNodePredicateAdditionalData(ctx *NodePredicateAdditionalDataContext) {
}

// ExitNodePredicateAdditionalData is called when production nodePredicateAdditionalData is exited.
func (s *BaseAqlParserListener) ExitNodePredicateAdditionalData(ctx *NodePredicateAdditionalDataContext) {
}

// EnterVersionPredicate is called when production versionPredicate is entered.
func (s *BaseAqlParserListener) EnterVersionPredicate(ctx *VersionPredicateContext) {}

// ExitVersionPredicate is called when production versionPredicate is exited.
func (s *BaseAqlParserListener) ExitVersionPredicate(ctx *VersionPredicateContext) {}

// EnterPathPredicateOperand is called when production pathPredicateOperand is entered.
func (s *BaseAqlParserListener) EnterPathPredicateOperand(ctx *PathPredicateOperandContext) {}

// ExitPathPredicateOperand is called when production pathPredicateOperand is exited.
func (s *BaseAqlParserListener) ExitPathPredicateOperand(ctx *PathPredicateOperandContext) {}

// EnterObjectPath is called when production objectPath is entered.
func (s *BaseAqlParserListener) EnterObjectPath(ctx *ObjectPathContext) {}

// ExitObjectPath is called when production objectPath is exited.
func (s *BaseAqlParserListener) ExitObjectPath(ctx *ObjectPathContext) {}

// EnterPathPart is called when production pathPart is entered.
func (s *BaseAqlParserListener) EnterPathPart(ctx *PathPartContext) {}

// ExitPathPart is called when production pathPart is exited.
func (s *BaseAqlParserListener) ExitPathPart(ctx *PathPartContext) {}

// EnterLikeOperand is called when production likeOperand is entered.
func (s *BaseAqlParserListener) EnterLikeOperand(ctx *LikeOperandContext) {}

// ExitLikeOperand is called when production likeOperand is exited.
func (s *BaseAqlParserListener) ExitLikeOperand(ctx *LikeOperandContext) {}

// EnterMatchesOperand is called when production matchesOperand is entered.
func (s *BaseAqlParserListener) EnterMatchesOperand(ctx *MatchesOperandContext) {}

// ExitMatchesOperand is called when production matchesOperand is exited.
func (s *BaseAqlParserListener) ExitMatchesOperand(ctx *MatchesOperandContext) {}

// EnterValueListItem is called when production valueListItem is entered.
func (s *BaseAqlParserListener) EnterValueListItem(ctx *ValueListItemContext) {}

// ExitValueListItem is called when production valueListItem is exited.
func (s *BaseAqlParserListener) ExitValueListItem(ctx *ValueListItemContext) {}

// EnterPrimitive is called when production primitive is entered.
func (s *BaseAqlParserListener) EnterPrimitive(ctx *PrimitiveContext) {}

// ExitPrimitive is called when production primitive is exited.
func (s *BaseAqlParserListener) ExitPrimitive(ctx *PrimitiveContext) {}

// EnterNumericPrimitive is called when production numericPrimitive is entered.
func (s *BaseAqlParserListener) EnterNumericPrimitive(ctx *NumericPrimitiveContext) {}

// ExitNumericPrimitive is called when production numericPrimitive is exited.
func (s *BaseAqlParserListener) ExitNumericPrimitive(ctx *NumericPrimitiveContext) {}

// EnterFunctionCall is called when production functionCall is entered.
func (s *BaseAqlParserListener) EnterFunctionCall(ctx *FunctionCallContext) {}

// ExitFunctionCall is called when production functionCall is exited.
func (s *BaseAqlParserListener) ExitFunctionCall(ctx *FunctionCallContext) {}

// EnterAggregateFunctionCall is called when production aggregateFunctionCall is entered.
func (s *BaseAqlParserListener) EnterAggregateFunctionCall(ctx *AggregateFunctionCallContext) {}

// ExitAggregateFunctionCall is called when production aggregateFunctionCall is exited.
func (s *BaseAqlParserListener) ExitAggregateFunctionCall(ctx *AggregateFunctionCallContext) {}

// EnterTerminologyFunction is called when production terminologyFunction is entered.
func (s *BaseAqlParserListener) EnterTerminologyFunction(ctx *TerminologyFunctionContext) {}

// ExitTerminologyFunction is called when production terminologyFunction is exited.
func (s *BaseAqlParserListener) ExitTerminologyFunction(ctx *TerminologyFunctionContext) {}

// EnterTop is called when production top is entered.
func (s *BaseAqlParserListener) EnterTop(ctx *TopContext) {}

// ExitTop is called when production top is exited.
func (s *BaseAqlParserListener) ExitTop(ctx *TopContext) {}
