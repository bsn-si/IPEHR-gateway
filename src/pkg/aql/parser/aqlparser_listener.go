// Code generated from parser.g4 by ANTLR 4.12.0. DO NOT EDIT.

package parser // AqlParser
import "github.com/antlr/antlr4/runtime/Go/antlr/v4"

// AqlParserListener is a complete listener for a parse tree produced by parser.
type AqlParserListener interface {
	antlr.ParseTreeListener

	// EnterSelectQuery is called when entering the selectQuery production.
	EnterSelectQuery(c *SelectQueryContext)

	// EnterSelectClause is called when entering the selectClause production.
	EnterSelectClause(c *SelectClauseContext)

	// EnterFromClause is called when entering the fromClause production.
	EnterFromClause(c *FromClauseContext)

	// EnterWhereClause is called when entering the whereClause production.
	EnterWhereClause(c *WhereClauseContext)

	// EnterOrderByClause is called when entering the orderByClause production.
	EnterOrderByClause(c *OrderByClauseContext)

	// EnterLimitClause is called when entering the limitClause production.
	EnterLimitClause(c *LimitClauseContext)

	// EnterSelectExpr is called when entering the selectExpr production.
	EnterSelectExpr(c *SelectExprContext)

	// EnterFromExpr is called when entering the fromExpr production.
	EnterFromExpr(c *FromExprContext)

	// EnterWhereExpr is called when entering the whereExpr production.
	EnterWhereExpr(c *WhereExprContext)

	// EnterOrderByExpr is called when entering the orderByExpr production.
	EnterOrderByExpr(c *OrderByExprContext)

	// EnterColumnExpr is called when entering the columnExpr production.
	EnterColumnExpr(c *ColumnExprContext)

	// EnterContainsExpr is called when entering the containsExpr production.
	EnterContainsExpr(c *ContainsExprContext)

	// EnterIdentifiedExpr is called when entering the identifiedExpr production.
	EnterIdentifiedExpr(c *IdentifiedExprContext)

	// EnterClassExpression is called when entering the classExpression production.
	EnterClassExpression(c *ClassExpressionContext)

	// EnterVersionClassExpr is called when entering the versionClassExpr production.
	EnterVersionClassExpr(c *VersionClassExprContext)

	// EnterTerminal is called when entering the terminal production.
	EnterTerminal(c *TerminalContext)

	// EnterIdentifiedPath is called when entering the identifiedPath production.
	EnterIdentifiedPath(c *IdentifiedPathContext)

	// EnterPathPredicate is called when entering the pathPredicate production.
	EnterPathPredicate(c *PathPredicateContext)

	// EnterStandardPredicate is called when entering the standardPredicate production.
	EnterStandardPredicate(c *StandardPredicateContext)

	// EnterArchetypePredicate is called when entering the archetypePredicate production.
	EnterArchetypePredicate(c *ArchetypePredicateContext)

	// EnterNodePredicate is called when entering the nodePredicate production.
	EnterNodePredicate(c *NodePredicateContext)

	// EnterNodePredicateAdditionalData is called when entering the nodePredicateAdditionalData production.
	EnterNodePredicateAdditionalData(c *NodePredicateAdditionalDataContext)

	// EnterVersionPredicate is called when entering the versionPredicate production.
	EnterVersionPredicate(c *VersionPredicateContext)

	// EnterPathPredicateOperand is called when entering the pathPredicateOperand production.
	EnterPathPredicateOperand(c *PathPredicateOperandContext)

	// EnterObjectPath is called when entering the objectPath production.
	EnterObjectPath(c *ObjectPathContext)

	// EnterPathPart is called when entering the pathPart production.
	EnterPathPart(c *PathPartContext)

	// EnterLikeOperand is called when entering the likeOperand production.
	EnterLikeOperand(c *LikeOperandContext)

	// EnterMatchesOperand is called when entering the matchesOperand production.
	EnterMatchesOperand(c *MatchesOperandContext)

	// EnterValueListItem is called when entering the valueListItem production.
	EnterValueListItem(c *ValueListItemContext)

	// EnterPrimitive is called when entering the primitive production.
	EnterPrimitive(c *PrimitiveContext)

	// EnterNumericPrimitive is called when entering the numericPrimitive production.
	EnterNumericPrimitive(c *NumericPrimitiveContext)

	// EnterFunctionCall is called when entering the functionCall production.
	EnterFunctionCall(c *FunctionCallContext)

	// EnterAggregateFunctionCall is called when entering the aggregateFunctionCall production.
	EnterAggregateFunctionCall(c *AggregateFunctionCallContext)

	// EnterTerminologyFunction is called when entering the terminologyFunction production.
	EnterTerminologyFunction(c *TerminologyFunctionContext)

	// EnterTop is called when entering the top production.
	EnterTop(c *TopContext)

	// ExitSelectQuery is called when exiting the selectQuery production.
	ExitSelectQuery(c *SelectQueryContext)

	// ExitSelectClause is called when exiting the selectClause production.
	ExitSelectClause(c *SelectClauseContext)

	// ExitFromClause is called when exiting the fromClause production.
	ExitFromClause(c *FromClauseContext)

	// ExitWhereClause is called when exiting the whereClause production.
	ExitWhereClause(c *WhereClauseContext)

	// ExitOrderByClause is called when exiting the orderByClause production.
	ExitOrderByClause(c *OrderByClauseContext)

	// ExitLimitClause is called when exiting the limitClause production.
	ExitLimitClause(c *LimitClauseContext)

	// ExitSelectExpr is called when exiting the selectExpr production.
	ExitSelectExpr(c *SelectExprContext)

	// ExitFromExpr is called when exiting the fromExpr production.
	ExitFromExpr(c *FromExprContext)

	// ExitWhereExpr is called when exiting the whereExpr production.
	ExitWhereExpr(c *WhereExprContext)

	// ExitOrderByExpr is called when exiting the orderByExpr production.
	ExitOrderByExpr(c *OrderByExprContext)

	// ExitColumnExpr is called when exiting the columnExpr production.
	ExitColumnExpr(c *ColumnExprContext)

	// ExitContainsExpr is called when exiting the containsExpr production.
	ExitContainsExpr(c *ContainsExprContext)

	// ExitIdentifiedExpr is called when exiting the identifiedExpr production.
	ExitIdentifiedExpr(c *IdentifiedExprContext)

	// ExitClassExpression is called when exiting the classExpression production.
	ExitClassExpression(c *ClassExpressionContext)

	// ExitVersionClassExpr is called when exiting the versionClassExpr production.
	ExitVersionClassExpr(c *VersionClassExprContext)

	// ExitTerminal is called when exiting the terminal production.
	ExitTerminal(c *TerminalContext)

	// ExitIdentifiedPath is called when exiting the identifiedPath production.
	ExitIdentifiedPath(c *IdentifiedPathContext)

	// ExitPathPredicate is called when exiting the pathPredicate production.
	ExitPathPredicate(c *PathPredicateContext)

	// ExitStandardPredicate is called when exiting the standardPredicate production.
	ExitStandardPredicate(c *StandardPredicateContext)

	// ExitArchetypePredicate is called when exiting the archetypePredicate production.
	ExitArchetypePredicate(c *ArchetypePredicateContext)

	// ExitNodePredicate is called when exiting the nodePredicate production.
	ExitNodePredicate(c *NodePredicateContext)

	// ExitNodePredicateAdditionalData is called when exiting the nodePredicateAdditionalData production.
	ExitNodePredicateAdditionalData(c *NodePredicateAdditionalDataContext)

	// ExitVersionPredicate is called when exiting the versionPredicate production.
	ExitVersionPredicate(c *VersionPredicateContext)

	// ExitPathPredicateOperand is called when exiting the pathPredicateOperand production.
	ExitPathPredicateOperand(c *PathPredicateOperandContext)

	// ExitObjectPath is called when exiting the objectPath production.
	ExitObjectPath(c *ObjectPathContext)

	// ExitPathPart is called when exiting the pathPart production.
	ExitPathPart(c *PathPartContext)

	// ExitLikeOperand is called when exiting the likeOperand production.
	ExitLikeOperand(c *LikeOperandContext)

	// ExitMatchesOperand is called when exiting the matchesOperand production.
	ExitMatchesOperand(c *MatchesOperandContext)

	// ExitValueListItem is called when exiting the valueListItem production.
	ExitValueListItem(c *ValueListItemContext)

	// ExitPrimitive is called when exiting the primitive production.
	ExitPrimitive(c *PrimitiveContext)

	// ExitNumericPrimitive is called when exiting the numericPrimitive production.
	ExitNumericPrimitive(c *NumericPrimitiveContext)

	// ExitFunctionCall is called when exiting the functionCall production.
	ExitFunctionCall(c *FunctionCallContext)

	// ExitAggregateFunctionCall is called when exiting the aggregateFunctionCall production.
	ExitAggregateFunctionCall(c *AggregateFunctionCallContext)

	// ExitTerminologyFunction is called when exiting the terminologyFunction production.
	ExitTerminologyFunction(c *TerminologyFunctionContext)

	// ExitTop is called when exiting the top production.
	ExitTop(c *TopContext)
}
