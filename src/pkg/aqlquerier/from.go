package aqlquerier

import (
	"fmt"
	"hms/gateway/pkg/aqlprocessor"
	"hms/gateway/pkg/errors"
	"hms/gateway/pkg/storage/treeindex"
)

type dataSource struct {
	name  string
	alias string
	data  treeindex.Container
}

func (ds dataSource) getName() string {
	if ds.alias != "" {
		return ds.alias
	}

	return ds.name
}

func (exec *executer) findSources() (map[string]dataSource, error) {
	from := exec.query.From
	if len(from.Contains) > 0 || from.Operator != nil {
		return nil, errors.New("not implemented")
	}

	result := map[string]dataSource{}

	if from.Operand == nil {
		return result, nil
	}

	switch operand := from.Operand.(type) {
	case aqlprocessor.ClassExpression:
		ds, err := exec.getDataSourceForClassExpression(operand)
		if err != nil {
			return nil, errors.Wrap(err, "cannot get data source for class expression")
		}

		result[ds.getName()] = ds
	case aqlprocessor.VersionClassExpr:
		return nil, errors.New("not implemented")
	default:
		return nil, fmt.Errorf("unexpected FROM.Operand type: %T", operand) // nolint
	}

	return result, nil
}

func (exec *executer) getDataSourceForClassExpression(operand aqlprocessor.ClassExpression) (dataSource, error) {
	if operand.PathPredicate != nil {
		return dataSource{}, errors.New("not implemented")
	}

	ds := dataSource{
		name: operand.Identifiers[0],
	}

	if len(operand.Identifiers) > 1 {
		ds.alias = operand.Identifiers[1]
	}

	data, err := exec.index.GetDataSourceByName(ds.name)
	if err != nil {
		return dataSource{}, errors.Wrap(err, "cannot get data source by name")
	}

	ds.data = data

	return ds, nil
}
