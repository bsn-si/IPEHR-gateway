package aqlquerier

import (
	"fmt"
	"hms/gateway/pkg/aqlprocessor"
	"hms/gateway/pkg/errors"
	"hms/gateway/pkg/storage/treeindex"
)

func processWhere(where *aqlprocessor.Where, sources map[string]dataSource) (map[string]dataSource, error) {
	if ie := where.IdentifiedExpr; ie != nil {
		return getDataSourceForIdentifierExpr(ie, sources)
	}

	if where.OperatorType != aqlprocessor.NoneOperator && where.OperatorType != "" {
		if len(where.Next) == 0 {
			return nil, errors.New("unexpected where conditions count")
		}

		results := make([]map[string]dataSource, len(where.Next))

		for i, where := range where.Next {
			s, err := processWhere(where, sources)
			if err != nil {
				return nil, errors.Wrap(err, "cannot filter inner WHERE conditions")
			}

			results[i] = s
		}

		switch where.OperatorType {
		case aqlprocessor.NOTOperator:
			if len(results) != 1 {
				return sources, nil
			}

			return mergeDataSourcesNOT(sources, results[0]), nil
		case aqlprocessor.ANDOperator:
			if len(results) != 2 {
				return nil, errors.New("invalid data sources count")
			}

			return mergeDataSourcesAND(results[0], results[1]), nil
		case aqlprocessor.OROperator:
			if len(results) != 2 {
				return nil, errors.New("invalid data sources count")
			}

			return mergeDataSourcesOR(results[0], results[1]), nil
		default:
			return nil, fmt.Errorf("unexpected operator type: %v", where.OperatorType) //nolint
		}
	} else if len(where.Next) == 1 {
		return processWhere(where.Next[0], sources)
	}

	return nil, errors.New("unexpected WHERE object state")
}

func getDataSourceForIdentifierExpr(ie *aqlprocessor.IdentifiedExpr, sources map[string]dataSource) (map[string]dataSource, error) {
	result := map[string]dataSource{}

	if ie.ComparisonOperator == nil || ie.IdentifiedPath == nil {
		return result, nil
	}

	for key, source := range sources {
		ip := *ie.IdentifiedPath
		if key == ip.Identifier {
			newSource := dataSource{
				name:  source.name,
				alias: source.alias,
				data:  treeindex.Container{},
			}

			for key, indexNodes := range source.data {
				for _, indexNode := range indexNodes {
					containsValues := false

					if ip.ObjectPath != nil {
						value, ok := getValueForPath(ip.ObjectPath, indexNode)
						if !ok {
							continue
						}

						containsValues = compare(ie.Terminal, value, *ie.ComparisonOperator)
					}

					if containsValues {
						newSource.data[key] = append(newSource.data[key], indexNode)
					}
				}
			}

			result[key] = newSource
		} else {
			result[key] = source
		}
	}

	return result, nil
}

func mergeDataSourcesNOT(origin, exclude map[string]dataSource) map[string]dataSource {
	result := map[string]dataSource{}

	for key, source := range origin {
		excludeSource, ok := exclude[key]
		if !ok {
			result[key] = source
		}

		newSource := dataSource{
			name:  source.name,
			alias: source.alias,
			data:  treeindex.Container{},
		}

		for key, val := range source.data {
			if _, ok := excludeSource.data[key]; !ok {
				newSource.data[key] = val
				continue
			}
		}

		if len(newSource.data) > 0 {
			result[key] = newSource
		}
	}

	return result
}

func mergeDataSourcesAND(left, right map[string]dataSource) map[string]dataSource {
	result := map[string]dataSource{}

	for k, val1 := range left {
		val2, ok := right[k]
		if !ok {
			continue
		}

		commonDataSource := dataSource{
			name:  val1.name,
			alias: val1.alias,
			data:  treeindex.Container{},
		}

		for observKey, nodesA := range val1.data {
			if _, ok := val2.data[observKey]; !ok {
				continue
			}

			commonDataSource.data[observKey] = nodesA
		}

		if len(commonDataSource.data) > 0 {
			result[k] = commonDataSource
		}
	}

	return result
}

func mergeDataSourcesOR(left, right map[string]dataSource) map[string]dataSource {
	result := left

	for key, rightSource := range right {
		leftSource, ok := result[key]
		if !ok {
			result[key] = rightSource
			continue
		}

		for key, val := range rightSource.data {
			if _, ok := leftSource.data[key]; ok {
				continue
			}

			leftSource.data[key] = val
		}
	}

	return result
}

func compare(term *aqlprocessor.Terminal, val any, cmpOperator aqlprocessor.ComparisionSymbol) bool {
	if term.Primitive != nil {
		return term.Primitive.Compare(val, cmpOperator)
	}

	//TODO: add logic for other conditions

	return false
}
