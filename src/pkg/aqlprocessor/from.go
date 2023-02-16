package aqlprocessor

import (
	"fmt"

	"github.com/bsn-si/IPEHR-gateway/src/pkg/aqlprocessor/aqlparser"
	"github.com/bsn-si/IPEHR-gateway/src/pkg/errors"
)

type From struct {
	ContainsExpr
}

type ContainsExpr struct {
	Operand  *ClassExpression
	Contains []*ContainsExpr
	Operator *OperatorType
}

type OperatorType string

const (
	NoneOperator OperatorType = "NONE"
	ANDOperator  OperatorType = "AND"
	OROperator   OperatorType = "OR"
	NOTOperator  OperatorType = "NOT"
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
		cExpr, err := getContainsExpr(ctx.ContainsExpr().(*aqlparser.ContainsExprContext))
		if err != nil {
			return From{}, errors.Wrap(err, "cannot process From.ContainsExpr")
		}

		f.ContainsExpr = *cExpr
	} else {
		return From{}, errors.New("empty From.ContainsExpr")
	}

	return f, nil
}

func getContainsExpr(ctx *aqlparser.ContainsExprContext) (*ContainsExpr, error) {
	result := ContainsExpr{}

	if ctx.ClassExprOperand() != nil {
		switch ctx := ctx.ClassExprOperand().(type) {
		case *aqlparser.ClassExpressionContext:
			{
				ce := ClassExpression{}
				for _, id := range ctx.AllIDENTIFIER() {
					ce.Identifiers = append(ce.Identifiers, id.GetText())
				}

				if ctx.PathPredicate() != nil {
					p, err := getPathPredicate(ctx.PathPredicate().(*aqlparser.PathPredicateContext))
					if err != nil {
						return nil, err
					}

					ce.PathPredicate = &p
				}

				result.Operand = &ce
			}
			/*
				case *aqlparser.VersionClassExprContext:
					{
						vce := VersionClassExpr{}
						vce.Version = ctx.VERSION().GetText()
						vce.Variable = toRef(ctx.IDENTIFIER().GetText())
						if ctx.VersionPredicate() != nil {
							pp, err := getVersionPredicate(ctx.VersionPredicate().(*aqlparser.VersionPredicateContext))
							if err != nil {
								return nil, err
							}

							vce.VersionPredicate = &pp
						}

						result.Operand = vce
					}
			*/
		default:
			return nil, fmt.Errorf("unexpected ContainsExp operand class: %T", ctx) //nolint
		}
	}

	if len(ctx.AllContainsExpr()) > 0 {
		for _, ce := range ctx.AllContainsExpr() {
			cExp, err := getContainsExpr(ce.(*aqlparser.ContainsExprContext))
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

/*
func getVersionPredicate(ctx *aqlparser.VersionPredicateContext) (PathPredicate, error) {
	sp, err := getStandartPredicate(ctx.StandardPredicate().(*aqlparser.StandardPredicateContext))
	if err != nil {
		return PathPredicate{}, errors.Wrap(err, "cannot get VersionPredicate.StandardPredicate")
	}

	return PathPredicate{
		Type:              StandartPathPredicate,
		StandartPredicate: sp,
	}, nil
}
*/
