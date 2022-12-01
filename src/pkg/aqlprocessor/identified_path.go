package aqlprocessor

import (
	"hms/gateway/pkg/aqlprocessor/aqlparser"
	"hms/gateway/pkg/errors"
	"log"
)

type IdentifiedPath struct {
	Identifier    string
	PathPredicate *PathPredicate
	ObjectPath    *ObjectPath
}

type ObjectPath struct {
	Paths []PartPath
}

type PartPath struct {
	Identifier    string
	PathPredicate *PathPredicate
}

func NewIdentifiedPath(ctx *aqlparser.IdentifiedPathContext) IdentifiedPath {
	ip := IdentifiedPath{
		Identifier: ctx.IDENTIFIER().GetText(),
	}

	if pp := ctx.PathPredicate(); pp != nil {
		predicate, err := processPathPredicate(pp.(*aqlparser.PathPredicateContext))
		if err != nil {
			log.Fatal(err)
		}

		ip.PathPredicate = &predicate
	}

	if slash := ctx.SYM_SLASH(); slash != nil && ctx.ObjectPath() != nil {
		op, err := newObjectPath(ctx.ObjectPath().(*aqlparser.ObjectPathContext))
		if err != nil {
			log.Fatal(err)
		}

		ip.ObjectPath = op
	}

	return ip
}

func newObjectPath(ctx *aqlparser.ObjectPathContext) (*ObjectPath, error) {
	result := ObjectPath{
		Paths: make([]PartPath, 0, len(ctx.AllPathPart())),
	}

	for _, pp := range ctx.AllPathPart() {
		val, err := processPathPart(pp.(*aqlparser.PathPartContext))
		if err != nil {
			return nil, errors.Wrap(err, "cannot process ObjectPath.PathPart")
		}

		result.Paths = append(result.Paths, val)
	}

	return &result, nil
}

func processPathPart(ctx *aqlparser.PathPartContext) (PartPath, error) {
	op := PartPath{
		Identifier: ctx.IDENTIFIER().GetText(),
	}

	if ctx.PathPredicate() != nil {
		pp, err := processPathPredicate(ctx.PathPredicate().(*aqlparser.PathPredicateContext))
		if err != nil {
			return op, err
		}

		op.PathPredicate = &pp
	}

	return op, nil
}
