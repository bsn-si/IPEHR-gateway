package aqlquerier

import (
	"encoding/json"
	"fmt"
	"hms/gateway/pkg/aqlprocessor"
	"hms/gateway/pkg/errors"
	"hms/gateway/pkg/storage/treeindex"
	"log"
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

	data, _ := json.MarshalIndent(from, "", "\t")
	log.Println(string(data))

	result := make(map[string]dataSource)
	result, err := exec.getSources(exec.query.From.ContainsExpr, result)

	log.Println("resulted data sources", result)

	return result, err
}

func (exec *executer) getSources(from aqlprocessor.ContainsExpr, sources map[string]dataSource) (map[string]dataSource, error) {
	if from.Operand == nil {
		return sources, nil
	}

	switch operand := from.Operand.(type) {
	case aqlprocessor.ClassExpression:
		ds, err := exec.getDataSourceForClassExpression(operand)
		if err != nil {
			return nil, errors.Wrap(err, "cannot get data source for class expression")
		}

		sources[ds.getName()] = ds
	case aqlprocessor.VersionClassExpr:
		return nil, errors.New("not implemented")
	default:
		return nil, fmt.Errorf("unexpected FROM.Operand type: %T", operand) // nolint
	}

	if len(from.Contains) > 0 {
		sourceMaps := make([]map[string]dataSource, 0, len(from.Contains))
		sourceMaps = append(sourceMaps, sources)

		for _, ce := range from.Contains {
			if ce == nil {
				continue
			}

			s, err := exec.getSources(*ce, sources)
			if err != nil {
				return nil, errors.Wrap(err, "cannot get contains data sources")
			}

			sourceMaps = append(sourceMaps, s)
		}

		if from.Operator != nil {
			return nil, errors.New("not implemented")
		}

		return mergeDataSourceMaps(sourceMaps), nil
	}

	return sources, nil
}

func mergeDataSourceMaps(dsMaps []map[string]dataSource) map[string]dataSource {
	result := map[string]dataSource{}

	for _, m := range dsMaps {
		for key, ds := range m {
			if originDS, ok := result[key]; ok {
				for key, contanier := range ds.data {
					originDS.data[key] = contanier
				}

				continue
			}

			result[key] = ds
		}
	}

	return result
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

	switch ds.name {
	case "EHR":
		ehrs, err := exec.index.GetEHRs("")
		if err != nil {
			return dataSource{}, errors.Wrap(err, "cannot get data source for EHRs")
		}

		ds.data = make(treeindex.Container, len(ehrs))
		for _, ehrNode := range ehrs {
			ds.data[ehrNode.GetID()] = append(ds.data[ehrNode.GetID()], ehrNode)
		}
	case "COMPOSITION":
	default:
		return dataSource{}, fmt.Errorf("unexpeted data source type: %s", ds.name) //nolint
	}

	// data, err := exec.index.GetDataSourceByName(ds.name)
	// if err != nil {
	// 	return dataSource{}, errors.Wrap(err, "cannot get data source by name")
	// }

	// ds.data = data

	return ds, nil
}
