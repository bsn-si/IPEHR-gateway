package aqlprocessor

import (
	"github.com/bsn-si/IPEHR-gateway/src/pkg/aqlprocessor/aqlparser"
	"github.com/bsn-si/IPEHR-gateway/src/pkg/errors"
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

func getIdentifiedPath(ctx *aqlparser.IdentifiedPathContext) (IdentifiedPath, error) {
	ip := IdentifiedPath{
		Identifier: ctx.IDENTIFIER().GetText(),
	}

	if pp := ctx.PathPredicate(); pp != nil {
		predicate, err := getPathPredicate(pp.(*aqlparser.PathPredicateContext))
		if err != nil {
			return ip, errors.Wrap(err, "cannot process IdetiedPath.PathPredicate")
		}

		ip.PathPredicate = &predicate
	}

	if slash := ctx.SYM_SLASH(); slash != nil && ctx.ObjectPath() != nil {
		op, err := getObjectPath(ctx.ObjectPath().(*aqlparser.ObjectPathContext))
		if err != nil {
			return ip, errors.Wrap(err, "cannot get ObjectPath")
		}

		ip.ObjectPath = op
	}

	return ip, nil
}

func getObjectPath(ctx *aqlparser.ObjectPathContext) (*ObjectPath, error) {
	result := ObjectPath{
		Paths: make([]PartPath, 0, len(ctx.AllPathPart())),
	}

	for _, pp := range ctx.AllPathPart() {
		val, err := getPathPart(pp.(*aqlparser.PathPartContext))
		if err != nil {
			return nil, errors.Wrap(err, "cannot process ObjectPath.PathPart")
		}

		result.Paths = append(result.Paths, val)
	}

	return &result, nil
}

func getPathPart(ctx *aqlparser.PathPartContext) (PartPath, error) {
	op := PartPath{
		Identifier: ctx.IDENTIFIER().GetText(),
	}

	if ctx.PathPredicate() != nil {
		pp, err := getPathPredicate(ctx.PathPredicate().(*aqlparser.PathPredicateContext))
		if err != nil {
			return op, errors.Wrap(err, "cannot get PathPredicate")
		}

		op.PathPredicate = &pp
	}

	return op, nil
}
