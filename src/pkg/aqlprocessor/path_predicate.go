package aqlprocessor

import (
	"fmt"
	"hms/gateway/pkg/aqlprocessor/aqlparser"
	"hms/gateway/pkg/errors"

	"github.com/antlr/antlr4/runtime/Go/antlr/v4"
)

type PredicateType string

const (
	StandartPathPredicate   PredicateType = "STANDARD_PREDICATE"
	ArchetypedPathPredicate PredicateType = "ARCHETYPED_PREDICATE"
	NodePathPredicate       PredicateType = "NODE_PREDICATE"
)

type PathPredicate struct {
	Type PredicateType

	StandartPredicate *StandartPredicate
	NodePredicate     *NodePredicate
	Archetype         *ArchetypePathPredicate
}

type PathPredicateOperand struct {
	Primitive  *Primitive
	ObjectPath *ObjectPath
	Parameter  *string
	IDCode     *string
	AtCode     *string
}

// standardPredicate:
// objectPath COMPARISON_OPERATOR pathPredicateOperand;
type StandartPredicate struct {
	ObjectPath  *ObjectPath
	CMPOperator ComparisionSymbol
	Operand     *PathPredicateOperand
}

type ComparisionSymbol string

const (
	SymNone ComparisionSymbol = "NONE"
	SymLT   ComparisionSymbol = "<"
	SymGT   ComparisionSymbol = ">"
	SymLE   ComparisionSymbol = "<="
	SymGE   ComparisionSymbol = ">="
	SymNe   ComparisionSymbol = "!="
	SymEQ   ComparisionSymbol = "="
)

func getComparisionSimbol(ctx antlr.TerminalNode) (ComparisionSymbol, error) {
	switch ComparisionSymbol(ctx.GetText()) {
	case SymEQ:
		return SymEQ, nil
	case SymNe:
		return SymNe, nil
	case SymGE:
		return SymGE, nil
	case SymLE:
		return SymLE, nil
	case SymGT:
		return SymGT, nil
	case SymLT:
		return SymLT, nil
	default:
		return SymNone, fmt.Errorf("Unexpected comparison operator: %v", ctx.GetText()) //nolint
	}
}

type NodePredicate struct {
	Value string
}

type ArchetypePathPredicate struct {
	ArchetypeHRID *string
	Parameter     *string
}

func getPathPredicate(ctx *aqlparser.PathPredicateContext) (PathPredicate, error) {
	var (
		err    error
		result PathPredicate
	)

	if ctx.StandardPredicate() != nil {
		result, err = processStandartPredicate(ctx.StandardPredicate().(*aqlparser.StandardPredicateContext))
		if err != nil {
			return PathPredicate{}, errors.Wrap(err, "cannot process PathPredicate.StandartPredicate")
		}
	} else if ctx.ArchetypePredicate() != nil {
		result, err = processArchetypePredicate(ctx.ArchetypePredicate().(*aqlparser.ArchetypePredicateContext))
		if err != nil {
			return PathPredicate{}, errors.Wrap(err, "cannot process PathPredicate.ArchetypePredicate")
		}
	} else if ctx.NodePredicate() != nil {
		result, err = processNodePredicate(ctx.NodePredicate().(*aqlparser.NodePredicateContext))
		if err != nil {
			return PathPredicate{}, errors.Wrap(err, "cannot process PathPredicate.NodePredicate")
		}
	} else {
		return PathPredicate{}, fmt.Errorf("unknown path predicate type: %s", ctx.GetText()) //nolint
	}

	return result, nil
}

func processStandartPredicate(ctx *aqlparser.StandardPredicateContext) (PathPredicate, error) {
	pp := PathPredicate{
		Type:              StandartPathPredicate,
		StandartPredicate: &StandartPredicate{},
	}

	if ctx.ObjectPath() != nil {
		op, err := getObjectPath(ctx.ObjectPath().(*aqlparser.ObjectPathContext))
		if err != nil {
			return PathPredicate{}, errors.Wrap(err, "cannot get ObjectPath")
		}

		pp.StandartPredicate.ObjectPath = op
	}

	if ctx.COMPARISON_OPERATOR() != nil {
		symb, err := getComparisionSimbol(ctx.COMPARISON_OPERATOR())
		if err != nil {
			return PathPredicate{}, errors.Wrap(err, "cannot get ComparisionOperator")
		}

		pp.StandartPredicate.CMPOperator = symb
	}

	if ctx.PathPredicateOperand() != nil {
		operand, err := getPathPredicateOperand(ctx.PathPredicateOperand().(*aqlparser.PathPredicateOperandContext))
		if err != nil {
			return PathPredicate{}, err
		}

		pp.StandartPredicate.Operand = operand
	}

	return pp, nil
}

func processArchetypePredicate(ctx *aqlparser.ArchetypePredicateContext) (PathPredicate, error) {
	pp := PathPredicate{
		Type:      ArchetypedPathPredicate,
		Archetype: &ArchetypePathPredicate{},
	}

	if ctx.ARCHETYPE_HRID() != nil {
		pp.Archetype.ArchetypeHRID = toRef(ctx.ARCHETYPE_HRID().GetText())
	} else if ctx.PARAMETER() != nil {
		pp.Archetype.Parameter = toRef(ctx.PARAMETER().GetText())
	} else {
		return PathPredicate{}, fmt.Errorf("unexpected archetype predicate: %v", ctx.GetText()) //nolint
	}

	return pp, nil
}

// nodePredicate: (ID_CODE | AT_CODE) (
//
//		SYM_COMMA (
//			STRING
//			| PARAMETER
//			| TERM_CODE
//			| AT_CODE
//			| ID_CODE
//		)
//	)?
//	| ARCHETYPE_HRID (
//		SYM_COMMA (
//			STRING
//			| PARAMETER
//			| TERM_CODE
//			| AT_CODE
//			| ID_CODE
//		)
//	)?
//	| PARAMETER
//	| objectPath COMPARISON_OPERATOR pathPredicateOperand
//	| objectPath MATCHES CONTAINED_REGEX
//	| nodePredicate AND nodePredicate
//
// | nodePredicate OR nodePredicate;
func processNodePredicate(ctx *aqlparser.NodePredicateContext) (PathPredicate, error) {
	pp := PathPredicate{
		Type: NodePathPredicate,
		NodePredicate: &NodePredicate{
			Value: ctx.GetText(),
		},
	}
	//TODO add FULL realisation here

	return pp, nil
}

func getPathPredicateOperand(ctx *aqlparser.PathPredicateOperandContext) (*PathPredicateOperand, error) {
	result := PathPredicateOperand{}

	if ctx.Primitive() != nil {
		p, err := getPrimitive(ctx.Primitive().(*aqlparser.PrimitiveContext))
		if err != nil {
			return nil, errors.Wrap(err, "cannot get PathPredicateOverand.Primitive")
		}
		result.Primitive = &p
	} else if ctx.ObjectPath() != nil {
		op, err := getObjectPath(ctx.ObjectPath().(*aqlparser.ObjectPathContext))
		if err != nil {
			return nil, errors.Wrap(err, "cannot get PathPredicateOperand.ObjectPath")
		}
		result.ObjectPath = op
	} else if ctx.PARAMETER() != nil {
		result.Parameter = toRef(ctx.PARAMETER().GetText())
	} else if ctx.AT_CODE() != nil {
		result.AtCode = toRef(ctx.AT_CODE().GetText())
	} else if ctx.ID_CODE() != nil {
		result.IDCode = toRef(ctx.ID_CODE().GetText())
	}

	return &result, nil
}

func toRef[T any](v T) *T {
	return &v
}
