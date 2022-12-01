package aqlprocessor

import (
	"fmt"
	"hms/gateway/pkg/aqlprocessor/aqlparser"
	"hms/gateway/pkg/errors"
	"log"
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
	ObjectPath  ObjectPath
	CMPOperator ComparisionSymbol
	Operand     *PathPredicateOperand
}

type ComparisionSymbol string

const (
	SymLT ComparisionSymbol = "<"
	SymGT ComparisionSymbol = ">"
	SymLE ComparisionSymbol = "<="
	SymGE ComparisionSymbol = ">="
	SymNe ComparisionSymbol = "!="
	SymEQ ComparisionSymbol = "="
)

type ArchetypePathPredicate struct {
	ArchetypeHRID *string
	Parameter     *string
}

func processPathPredicate(ctx *aqlparser.PathPredicateContext) (PathPredicate, error) {
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
	log.Println("STANDART PREDICATE")

	pp := PathPredicate{
		Type:              StandartPathPredicate,
		StandartPredicate: &StandartPredicate{},
	}

	if ctx.ObjectPath() != nil {
		op := ctx.ObjectPath().(*aqlparser.ObjectPathContext)

		log.Println("\t\tOBJECT_PATH:", op.GetText())
	}

	if ctx.COMPARISON_OPERATOR() != nil {
		switch ComparisionSymbol(ctx.COMPARISON_OPERATOR().GetText()) {
		case SymEQ:
			pp.StandartPredicate.CMPOperator = SymEQ
		case SymNe:
			pp.StandartPredicate.CMPOperator = SymNe
		case SymGE:
			pp.StandartPredicate.CMPOperator = SymGE
		case SymLE:
			pp.StandartPredicate.CMPOperator = SymLE
		case SymGT:
			pp.StandartPredicate.CMPOperator = SymGT
		case SymLT:
			pp.StandartPredicate.CMPOperator = SymLT
		default:
			return PathPredicate{}, fmt.Errorf("Unexpected comparison operator: %v", ctx.COMPARISON_OPERATOR().GetText()) //nolint
		}
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
	}
	//TODO add FULL realisation here

	return pp, nil
}

func getPathPredicateOperand(ctx *aqlparser.PathPredicateOperandContext) (*PathPredicateOperand, error) {
	result := PathPredicateOperand{}

	log.Println("\t\tPATH_PREDICATE_OPERAND", ctx.GetText())

	if ctx.Primitive() != nil {
		log.Println("\t\t\tPRIMITIVE", ctx.Primitive().GetText())
	} else if ctx.ObjectPath() != nil {
		log.Println("\t\t\tOBJECT PATH", ctx.ObjectPath().GetText())
	} else if ctx.PARAMETER() != nil {
		log.Println("\t\t\tPARAMETER", ctx.PARAMETER().GetText())
		result.Parameter = toRef(ctx.PARAMETER().GetText())
	} else {
		log.Printf("\t\t\tUNEXPECTED TYPE: %T", ctx)
	}

	log.Println("\tPATH OPERAND", ctx.GetChildCount())

	return &result, nil
}

func toRef[T any](v T) *T {
	return &v
}
