package driver

import (
	"fmt"

	aqlprocessor "github.com/bsn-si/IPEHR-gateway/src/pkg/aql/processor"
	"github.com/bsn-si/IPEHR-gateway/src/pkg/errors"
)

func processWhere(where *aqlprocessor.Where, sources dataRows) (dataRows, error) {
	if ie := where.IdentifiedExpr; ie != nil {
		return getDataSourceForIdentifierExpr(ie, sources)
	}

	if where.OperatorType != aqlprocessor.NoneOperator && where.OperatorType != "" {
		if len(where.Next) == 0 {
			return nil, errors.New("unexpected where conditions count")
		}

		results := make([]dataRows, len(where.Next))

		for i, where := range where.Next {
			processedDataSources, err := processWhere(where, sources)
			if err != nil {
				return nil, errors.Wrap(err, "cannot filter inner WHERE conditions")
			}

			results[i] = processedDataSources
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

func getDataSourceForIdentifierExpr(ie *aqlprocessor.IdentifiedExpr, rows dataRows) (dataRows, error) {
	result := dataRows{}

	if ie.ComparisonOperator == nil || ie.IdentifiedPath == nil {
		return result, nil
	}

	for _, row := range rows {
		ip := *ie.IdentifiedPath

		if cell, ok := row.cells[ip.Identifier]; ok {
			containsValues := false

			indexNode := cell.data
			if ip.ObjectPath != nil {
				value, ok := getValueForPath(ip.ObjectPath, indexNode)
				if !ok {
					continue
				}

				containsValues = compare(ie.Terminal, value, *ie.ComparisonOperator)
			}

			if containsValues {
				result = append(result, row)
			}
		} else {
			result = append(result, row)
		}
	}

	return result, nil
}

func mergeDataSourcesNOT(origin, exclude dataRows) dataRows {
	result := dataRows{}

	m := make(map[string]bool, len(exclude))
	for _, r := range exclude {
		m[r.id.String()] = true
	}

	for _, row := range origin {
		if m[row.id.String()] {
			continue
		}

		result = append(result, row)
	}

	return result
}

func mergeDataSourcesAND(left, right dataRows) dataRows {
	result := dataRows{}

	m := make(map[string]bool, len(right))
	for _, r := range right {
		m[r.id.String()] = true
	}

	for _, row := range left {
		if !m[row.id.String()] {
			continue
		}

		result = append(result, row)
	}

	return result
}

func mergeDataSourcesOR(left, right dataRows) dataRows {
	result := left

	m := make(map[string]bool, len(left))
	for _, r := range left {
		m[r.id.String()] = true
	}

	for _, r := range right {
		if m[r.id.String()] {
			continue
		}

		result = append(result, r)
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
