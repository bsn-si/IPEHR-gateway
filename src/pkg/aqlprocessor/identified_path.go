package aqlprocessor

import (
	"hms/gateway/pkg/aqlprocessor/aqlparser"
	"hms/gateway/pkg/errors"
	"log"
)

type IdentifiedPath struct {
	Name          string
	PathPredicate string
	Paths         []ObjectPath
}

type ObjectPath struct {
	Name          string
	PathPredicate string
}

func NewIdentifiedPath(ctx *aqlparser.IdentifiedPathContext) IdentifiedPath {
	ip := IdentifiedPath{
		Name: ctx.IDENTIFIER().GetText(),
	}

	if pp := ctx.PathPredicate(); pp != nil {
		predicate, err := processPathPredicate(pp.(*aqlparser.PathPredicateContext))
		if err != nil {
			log.Fatal(err)
		}

		ip.PathPredicate = predicate
	}

	if slash := ctx.SYM_SLASH(); slash != nil {
		if ctx.ObjectPath() != nil {
			paths, err := processObjectPath(ctx.ObjectPath().(*aqlparser.ObjectPathContext))
			if err != nil {
				log.Fatal(err)
			}

			ip.Paths = paths
		}
	}

	return ip
}

func processObjectPath(ctx *aqlparser.ObjectPathContext) ([]ObjectPath, error) {
	result := make([]ObjectPath, 0, len(ctx.AllPathPart()))

	for _, pp := range ctx.AllPathPart() {
		val, err := processPathPart(pp.(*aqlparser.PathPartContext))
		if err != nil {
			return nil, errors.Wrap(err, "cannot process ObjectPath.PathPart")
		}

		result = append(result, val)
	}

	return result, nil
}

func processPathPart(ctx *aqlparser.PathPartContext) (ObjectPath, error) {
	op := ObjectPath{
		Name: ctx.IDENTIFIER().GetText(),
	}

	if ctx.PathPredicate() != nil {
		pp, err := processPathPredicate(ctx.PathPredicate().(*aqlparser.PathPredicateContext))
		if err != nil {
			return op, err
		}

		op.PathPredicate = pp
	}

	return op, nil
}

func processPathPredicate(ctx *aqlparser.PathPredicateContext) (string, error) {
	log.Println("path predicate childs", ctx.GetChildCount())

	if ctx.StandardPredicate() != nil {
		sp := ctx.StandardPredicate().(*aqlparser.StandardPredicateContext)

		return sp.GetText(), nil
	} else if ctx.ArchetypePredicate() != nil {
		return "", errors.New("archetype predicate handler not implemented")
	} else if ctx.NodePredicate() != nil {
		str, err := processNodePredicate(ctx.NodePredicate().(*aqlparser.NodePredicateContext))
		if err != nil {
			return "", errors.Wrap(err, "cannot process PathPredicate.NodePredicate")
		}

		return str, nil
	}

	return ctx.GetText(), nil
}

func processNodePredicate(np *aqlparser.NodePredicateContext) (string, error) {
	return np.GetText(), nil
}
