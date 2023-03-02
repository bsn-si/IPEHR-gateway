package processor

import (
	"fmt"
	"io"

	"github.com/bsn-si/IPEHR-gateway/src/pkg/aql/parser"
	"github.com/bsn-si/IPEHR-gateway/src/pkg/errors"
)

type From struct {
	ContainsExpr
}

func (f *From) write(w io.Writer) {
	fmt.Fprint(w, "FROM\n")
	f.ContainsExpr.write(w)
}

type ContainsExpr struct {
	Operand  Operand
	Contains []*ContainsExpr
	Operator *OperatorType
	Brackets bool
}

type Operand interface {
	write(w io.Writer)
}

func (cw *ContainsExpr) write(w io.Writer) {
	if cw.Operand != nil {
		cw.Operand.write(w)
	}

	if len(cw.Contains) > 0 {
		if cw.Operand != nil {
			if cw.Operator != nil && *cw.Operator == NOTOperator {
				fmt.Fprintf(w, "NOT ")
			}

			fmt.Fprintf(w, "CONTAINS ")
		}

		if cw.Brackets {
			fmt.Fprintf(w, "(")
		}

		if len(cw.Contains) == 1 {
			cw.Contains[0].write(w)
		}

		if len(cw.Contains) == 2 {
			cw.Contains[0].write(w)

			if cw.Operator != nil {
				fmt.Fprintf(w, " %s ", *cw.Operator)
			}

			cw.Contains[1].write(w)
		}

		if cw.Brackets {
			fmt.Fprintf(w, ")")
		}
	}
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

func (ce ClassExpression) write(w io.Writer) {
	for i := range ce.Identifiers {
		if i != 0 {
			fmt.Fprint(w, " ")
		}

		fmt.Fprint(w, ce.Identifiers[i])
	}

	if ce.PathPredicate != nil {
		fmt.Fprint(w, "[")
		ce.PathPredicate.write(w)
		fmt.Fprint(w, "]")
	}
}

type VersionClassExpr struct {
	Version          string
	Variable         *string
	VersionPredicate *VersionPredicate
}

func (vce VersionClassExpr) write(w io.Writer) {
	fmt.Fprintf(w, "%s ", vce.Version)

	if vce.Variable != nil {
		fmt.Fprintf(w, "%s", *vce.Variable)
	}

	if vce.VersionPredicate != nil {
		fmt.Fprint(w, "[")
		vce.VersionPredicate.write(w)
		fmt.Fprint(w, "]")
	}
}

type VersionPredicate struct {
	LatestVersion     *string
	AllVersions       *string
	StandartPredicate *StandartPredicate
}

func (vp VersionPredicate) write(w io.Writer) {
	if vp.LatestVersion != nil {
		fmt.Fprint(w, vp.LatestVersion)
		return
	}

	if vp.AllVersions != nil {
		fmt.Fprint(w, vp.AllVersions)
		return
	}

	if vp.StandartPredicate != nil {
		vp.StandartPredicate.write(w)
	}
}

func getFrom(ctx *parser.FromExprContext) (From, error) {
	f := From{}

	if ctx.ContainsExpr() != nil {
		cExpr, err := getContainsExpr(ctx.ContainsExpr().(*parser.ContainsExprContext))
		if err != nil {
			return From{}, errors.Wrap(err, "cannot process From.ContainsExpr")
		}

		f.ContainsExpr = *cExpr
	} else {
		return From{}, errors.New("empty From.ContainsExpr")
	}

	return f, nil
}

func getContainsExpr(ctx *parser.ContainsExprContext) (*ContainsExpr, error) {
	result := ContainsExpr{}

	if ctx.ClassExprOperand() != nil {
		switch ctx := ctx.ClassExprOperand().(type) {
		case *parser.ClassExpressionContext:
			{
				ce := ClassExpression{}
				for _, id := range ctx.AllIDENTIFIER() {
					ce.Identifiers = append(ce.Identifiers, id.GetText())
				}

				if ctx.PathPredicate() != nil {
					p, err := getPathPredicate(ctx.PathPredicate().(*parser.PathPredicateContext))
					if err != nil {
						return nil, err
					}

					ce.PathPredicate = &p
				}

				result.Operand = ce
			}
		case *parser.VersionClassExprContext:
			{
				vce := VersionClassExpr{}
				vce.Version = ctx.VERSION().GetText()
				vce.Variable = toRef(ctx.IDENTIFIER().GetText())
				if ctx.VersionPredicate() != nil {
					vp, err := getVersionPredicate(ctx.VersionPredicate().(*parser.VersionPredicateContext))
					if err != nil {
						return nil, err
					}

					vce.VersionPredicate = &vp
				}

				result.Operand = vce
			}
		default:
			return nil, fmt.Errorf("unexpected ContainsExp operand class: %T", ctx) //nolint
		}
	}

	if len(ctx.AllContainsExpr()) > 0 {
		for _, ce := range ctx.AllContainsExpr() {
			cExp, err := getContainsExpr(ce.(*parser.ContainsExprContext))
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

	if ctx.SYM_LEFT_PAREN() != nil && ctx.SYM_RIGHT_PAREN() != nil {
		result.Brackets = true
	}

	return &result, nil
}

func getVersionPredicate(ctx *parser.VersionPredicateContext) (VersionPredicate, error) {
	vp := VersionPredicate{}
	if ctx.LATEST_VERSION() != nil {
		vp.LatestVersion = toRef(ctx.LATEST_VERSION().GetText())
	} else if ctx.ALL_VERSIONS() != nil {
		vp.AllVersions = toRef(ctx.ALL_VERSIONS().GetText())
	} else if ctx.StandardPredicate() != nil {
		sp, err := getStandartPredicate(ctx.StandardPredicate().(*parser.StandardPredicateContext))
		if err != nil {
			return VersionPredicate{}, errors.Wrap(err, "cannot get VersionPredicate.StandardPredicate")
		}

		vp.StandartPredicate = sp
	}

	return vp, nil
}
