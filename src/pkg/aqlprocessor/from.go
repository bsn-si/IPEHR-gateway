package aqlprocessor

import (
	"fmt"
	"hms/gateway/pkg/aqlprocessor/aqlparser"
	"hms/gateway/pkg/errors"
	"log"
)

type From struct {
	ContainsExpr
}

type ContainsExpr struct {
	Operand  any
	Contains []*ContainsExpr
	Operator *OperatorType
}

type OperatorType string

const (
	ANDOperator OperatorType = "AND"
	OROperator  OperatorType = "OR"
	NOTOperator OperatorType = "NOT"
)

type ClassExpression struct {
	Identifiers   []string
	PathPredicate *PathPredicate
}

type VersionClassExpr struct {
	Version          string
	Variable         *string
	VersionPredicate *PathPredicate
}

func getFrom(ctx *aqlparser.FromExprContext) (From, error) {
	f := From{}

	if ctx.ContainsExpr() != nil {
		cExpr, err := newContainsExpr(ctx.ContainsExpr().(*aqlparser.ContainsExprContext))
		if err != nil {
			return From{}, errors.Wrap(err, "cannot process From.ContainsExpr")
		}

		f.ContainsExpr = *cExpr
	} else {
		return From{}, errors.New("empty From.ContainsExpr")
	}

	return f, nil
}

func newContainsExpr(ctx *aqlparser.ContainsExprContext) (*ContainsExpr, error) {
	result := ContainsExpr{}

	if ctx.ClassExprOperand() != nil {
		switch ctx := ctx.ClassExprOperand().(type) {
		case *aqlparser.ClassExpressionContext:
			{
				log.Println("CLASS EXPR OPERAND", ctx.GetText())

				ce := ClassExpression{}
				for _, id := range ctx.AllIDENTIFIER() {
					log.Println("\t", id.GetText())
					ce.Identifiers = append(ce.Identifiers, id.GetText())
				}

				if ctx.PathPredicate() != nil {
					p, err := processPathPredicate(ctx.PathPredicate().(*aqlparser.PathPredicateContext))
					if err != nil {
						return nil, err
					}

					log.Println("\tPATH_PREDICATE", p)
					ce.PathPredicate = &p
				}

				result.Operand = ce
			}
		case *aqlparser.VersionClassExprContext:
			{
				vce := VersionClassExpr{}
				vce.Version = ctx.VERSION().GetText()
				vce.Variable = toRef(ctx.IDENTIFIER().GetText())
				if ctx.VersionPredicate() != nil {
					pp, err := processVersionPredicate(ctx.VersionPredicate().(*aqlparser.VersionPredicateContext))
					if err != nil {
						return nil, err
					}

					vce.VersionPredicate = &pp
				}

				result.Operand = vce
			}
		default:
			return nil, fmt.Errorf("unexpected ContainsExp operand class: %T", ctx) //nolint
		}
	}

	if len(ctx.AllContainsExpr()) > 0 {
		for _, ce := range ctx.AllContainsExpr() {
			cExp, err := newContainsExpr(ce.(*aqlparser.ContainsExprContext))
			if err != nil {
				return nil, err
			}

			result.Contains = append(result.Contains, cExp)
		}
	}

	if ctx.AND() != nil {
		result.Operator = toRef(ANDOperator)
	}

	if ctx.OR() != nil {
		result.Operator = toRef(OROperator)
	}

	if ctx.NOT() != nil {
		result.Operator = toRef(NOTOperator)
	}

	return &result, nil
}

func processVersionPredicate(ctx *aqlparser.VersionPredicateContext) (PathPredicate, error) {
	return processStandartPredicate(ctx.StandardPredicate().(*aqlparser.StandardPredicateContext))
}
