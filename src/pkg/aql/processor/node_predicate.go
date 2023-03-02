package processor

import (
	"fmt"
	"io"
	"strings"

	"github.com/bsn-si/IPEHR-gateway/src/pkg/aql/parser"
	"github.com/bsn-si/IPEHR-gateway/src/pkg/errors"

	"github.com/antlr/antlr4/runtime/Go/antlr/v4"
)

type NodePredicate struct {
	Operator          OperatorType
	ComparisionSymbol ComparisionSymbol

	AtCode         *AtCode
	IDCode         *IDCode
	ArchetypeHRID  *string
	AdditionalData *NodePredicateAdditionalData

	Parameter            *Parameter
	Next                 []*NodePredicate
	ObjectPath           *ObjectPath
	PathPredicateOperand *PathPredicateOperand
	IsMatches            bool
	ContainedRegex       *string
}

func (np *NodePredicate) write(w io.Writer) {
	fmt.Fprint(w, "nodePridecate")
}

type NodePredicateAdditionalData struct {
	String    *string
	Parameter *Parameter
	TermCode  *string
	AtCode    *AtCode
	IDCode    *IDCode
}

type AtCode string

func (code AtCode) ToString() string {
	return fmt.Sprintf("at%v", code)
}

func getAtCode(tn antlr.TerminalNode) (AtCode, error) {
	return getCode[AtCode](tn, "at")
}

type IDCode string

func getIDCode(tn antlr.TerminalNode) (IDCode, error) {
	return getCode[IDCode](tn, "id")
}

func getCode[T ~string](tn antlr.TerminalNode, prefix string) (T, error) {
	str := tn.GetText()
	if !strings.Contains(str, prefix) {
		return "", fmt.Errorf("string don't contains prefix '%s'", prefix) // nolint
	}

	return T(strings.TrimLeft(str, prefix)), nil
}

func getNodePredicate(ctx *parser.NodePredicateContext) (*NodePredicate, error) {
	if ctx.AT_CODE() != nil || ctx.ID_CODE() != nil {
		return getNodePredicateWithATOrIDCode(ctx)
	}

	if ctx.ARCHETYPE_HRID() != nil {
		return getNodePredicateWithArchetypeHIRD(ctx)
	}

	if ctx.ObjectPath() != nil && ctx.COMPARISON_OPERATOR() != nil && ctx.PathPredicateOperand() != nil {
		return getNodePredicateWithComparisionOperator(ctx)
	}

	if ctx.ObjectPath() != nil && ctx.MATCHES() != nil && ctx.CONTAINED_REGEX() != nil {
		return getNodePredicateWithMatches(ctx)
	}

	if ctx.PARAMETER() != nil {
		return getNodePredicateWithParameter(ctx)
	}

	if ctx.AND() != nil {
		return getNodePredicateWithAndOperator(ctx)
	}

	if ctx.OR() != nil {
		return getNodePredicateWithOrOperator(ctx)
	}

	return nil, errors.New("unexpected NodePredicate state")
}

func getNodePredicateWithATOrIDCode(ctx *parser.NodePredicateContext) (*NodePredicate, error) {
	np := &NodePredicate{
		Operator:          NoneOperator,
		ComparisionSymbol: SymNone,
	}

	if ctx.AT_CODE() != nil {
		atCode, err := getAtCode(ctx.AT_CODE())
		if err != nil {
			return nil, errors.Wrap(err, "cannot get NodePredicate.AT_CODE")
		}

		np.AtCode = &atCode
	} else if ctx.ID_CODE() != nil {
		idCode, err := getIDCode(ctx.ID_CODE())
		if err != nil {
			return nil, errors.Wrap(err, "cannot get NodePredicate.ID_CODE")
		}

		np.IDCode = &idCode
	}

	if ctx.SYM_COMMA() != nil && ctx.NodePredicateAdditionalData() != nil {
		ad, err := getNodePredicateAdditionalData(ctx.NodePredicateAdditionalData().(*parser.NodePredicateAdditionalDataContext))
		if err != nil {
			return nil, errors.Wrap(err, "cannot get NodePredicate.SYM_COMMA")
		}

		np.AdditionalData = ad
	}

	return np, nil
}

func getNodePredicateWithArchetypeHIRD(ctx *parser.NodePredicateContext) (*NodePredicate, error) {
	np := &NodePredicate{
		Operator:          NoneOperator,
		ComparisionSymbol: SymNone,
		ArchetypeHRID:     toRef(ctx.ARCHETYPE_HRID().GetText()),
	}

	if ctx.SYM_COMMA() != nil && ctx.NodePredicateAdditionalData() != nil {
		ad, err := getNodePredicateAdditionalData(ctx.NodePredicateAdditionalData().(*parser.NodePredicateAdditionalDataContext))
		if err != nil {
			return nil, errors.Wrap(err, "cannot get NodePredicate.SYM_COMMA")
		}

		np.AdditionalData = ad
	}

	return np, nil
}

func getNodePredicateWithComparisionOperator(ctx *parser.NodePredicateContext) (*NodePredicate, error) {
	cs, err := getComparisionSimbol(ctx.COMPARISON_OPERATOR())
	if err != nil {
		return nil, errors.Wrap(err, "cannot get NodePredicate.ComparisionOperator")
	}

	op, err := getObjectPath(ctx.ObjectPath().(*parser.ObjectPathContext))
	if err != nil {
		return nil, errors.Wrap(err, "cannot get NodePredicate.ObjectPath")
	}

	ppo, err := getPathPredicateOperand(ctx.PathPredicateOperand().(*parser.PathPredicateOperandContext))
	if err != nil {
		return nil, errors.Wrap(err, "cannot get NodePredicate.PathPredicateOperand")
	}

	np := &NodePredicate{
		Operator:             NoneOperator,
		ComparisionSymbol:    cs,
		ObjectPath:           op,
		PathPredicateOperand: ppo,
	}

	return np, nil
}

func getNodePredicateWithMatches(ctx *parser.NodePredicateContext) (*NodePredicate, error) {
	op, err := getObjectPath(ctx.ObjectPath().(*parser.ObjectPathContext))
	if err != nil {
		return nil, errors.Wrap(err, "cannot get NodePredicate.ObjectPath")
	}

	np := &NodePredicate{
		Operator:          NoneOperator,
		ComparisionSymbol: SymNone,
		IsMatches:         true,
		ObjectPath:        op,
		ContainedRegex:    toRef(ctx.CONTAINED_REGEX().GetText()),
	}

	return np, nil
}

func getNodePredicateWithParameter(ctx *parser.NodePredicateContext) (*NodePredicate, error) {
	p, err := getParameter(ctx.PARAMETER())
	if err != nil {
		return nil, errors.Wrap(err, "cannot get NodePredicate.PARAMETER")
	}

	np := &NodePredicate{
		Operator:          NoneOperator,
		ComparisionSymbol: SymNone,
		Parameter:         p,
	}

	return np, nil
}

func getNodePredicateWithAndOperator(ctx *parser.NodePredicateContext) (*NodePredicate, error) {
	np := &NodePredicate{
		Operator:          ANDOperator,
		ComparisionSymbol: SymNone,
		Next:              make([]*NodePredicate, 0, len(ctx.AllNodePredicate())),
	}

	for _, npCtx := range ctx.AllNodePredicate() {
		nextNp, err := getNodePredicate(npCtx.(*parser.NodePredicateContext))
		if err != nil {
			return nil, errors.Wrap(err, "cannot get NodePredicate.NodePredicate")
		}

		np.Next = append(np.Next, nextNp)
	}

	return np, nil
}

func getNodePredicateWithOrOperator(ctx *parser.NodePredicateContext) (*NodePredicate, error) {
	np := &NodePredicate{
		Operator:          OROperator,
		ComparisionSymbol: SymNone,
		Next:              make([]*NodePredicate, 0, len(ctx.AllNodePredicate())),
	}

	for _, npCtx := range ctx.AllNodePredicate() {
		nextNp, err := getNodePredicate(npCtx.(*parser.NodePredicateContext))
		if err != nil {
			return nil, errors.Wrap(err, "cannot get NodePredicate.NodePredicate")
		}

		np.Next = append(np.Next, nextNp)
	}

	return np, nil
}

func getNodePredicateAdditionalData(ctx *parser.NodePredicateAdditionalDataContext) (*NodePredicateAdditionalData, error) {
	result := NodePredicateAdditionalData{}

	if ctx.STRING() != nil {
		result.String = toRef(trimString(ctx.STRING().GetText()))
	} else if ctx.PARAMETER() != nil {
		p, err := getParameter(ctx.PARAMETER())
		if err != nil {
			return nil, errors.Wrap(err, "cannot get NodePredicate.SYM_COMMA.PARAMETER")
		}

		result.Parameter = p
	} else if ctx.TERM_CODE() != nil {
		result.TermCode = toRef(ctx.TERM_CODE().GetText())
	} else if ctx.AT_CODE() != nil {
		atCode, err := getAtCode(ctx.AT_CODE())
		if err != nil {
			return nil, errors.Wrap(err, "cannot get NodePredicate.SYM_COMMA.AT_CODE")
		}

		result.AtCode = &atCode
	} else if ctx.ID_CODE() != nil {
		idCode, err := getIDCode(ctx.ID_CODE())
		if err != nil {
			return nil, errors.Wrap(err, "cannot get NodePredicate.SYM_COMMA.ID_CODE")
		}

		result.IDCode = &idCode
	} else {
		return nil, errors.New("unexpected NodePredicate.SYM_COMMA object state")
	}

	return &result, nil
}
