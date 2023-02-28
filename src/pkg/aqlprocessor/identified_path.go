package aqlprocessor

import (
	"fmt"
	"io"

	"github.com/bsn-si/IPEHR-gateway/src/pkg/aqlprocessor/aqlparser"
	"github.com/bsn-si/IPEHR-gateway/src/pkg/errors"
)

type IdentifiedPath struct {
	Identifier    string
	PathPredicate *PathPredicate
	ObjectPath    *ObjectPath
}

func (ip *IdentifiedPath) write(w io.Writer) {
	fmt.Fprintf(w, "%s ", ip.Identifier)
	if ip.PathPredicate != nil {
		ip.PathPredicate.write(w)
	}
	if ip.ObjectPath != nil {
		fmt.Fprintln(w, "/")
		ip.ObjectPath.write(w)
	}
}

type ObjectPath struct {
	Paths []PartPath
}

func (op *ObjectPath) write(w io.Writer) {
	for i := range op.Paths {
		op.Paths[i].write(w)
		if i < len(op.Paths)-1 {
			fmt.Fprint(w, "/")
		}
	}
}

type PartPath struct {
	Identifier    string
	PathPredicate *PathPredicate
}

func (pp PartPath) write(w io.Writer) {
	fmt.Fprint(w, pp.Identifier)
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
