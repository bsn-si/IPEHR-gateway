package driver

import (
	"fmt"
	"reflect"

	aqlprocessor "github.com/bsn-si/IPEHR-gateway/src/pkg/aql/processor"
	"github.com/bsn-si/IPEHR-gateway/src/pkg/errors"
	"github.com/bsn-si/IPEHR-gateway/src/pkg/storage/treeindex"

	"github.com/google/uuid"
)

type dataCell struct {
	name  string
	alias string
	data  treeindex.Noder
}

func (dc dataCell) getName() string {
	if dc.alias != "" {
		return dc.alias
	}

	return dc.name
}

type dataRow struct {
	id    uuid.UUID
	cells map[string]dataCell
}

type dataRows []dataRow

func (exec *executer) findSources() (dataRows, error) {
	rows, err := exec.getDataRows(exec.query.From.ContainsExpr)
	if err != nil {
		return nil, errors.Wrap(err, "cannot find data rows")
	}

	return rows, nil
}

func (exec *executer) getDataRows(containsExpr aqlprocessor.ContainsExpr) (dataRows, error) {
	result, err := exec.processRowsContainsExpr(nil, &containsExpr)
	if err != nil {
		return nil, errors.Wrap(err, "cannot process rows")
	}

	return result, nil
}

func (exec *executer) processRowsContainsExpr(rootCell *dataCell, containsExpr *aqlprocessor.ContainsExpr) (dataRows, error) {
	result := dataRows{}

	switch operand := containsExpr.Operand.(type) {
	case aqlprocessor.ClassExpression:
		var node treeindex.Noder
		if rootCell != nil {
			node = rootCell.data
		}

		nodeDataCells, err := exec.getDataForClassExpr(node, operand)
		if err != nil {
			return nil, errors.Wrap(err, "cannot get data from node")
		}

		if len(containsExpr.Contains) > 0 {
			for i := range nodeDataCells {
				rowsSet := make([]dataRows, 0, len(containsExpr.Contains))

				for _, ce := range containsExpr.Contains {
					rows, err := exec.processRowsContainsExpr(&nodeDataCells[i], ce)
					if err != nil {
						return nil, errors.Wrap(err, "cannot process rows contains expr")
					}

					if rootCell != nil {
						for _, r := range rows {
							r.cells[rootCell.getName()] = *rootCell
						}
					}

					if len(rows) > 0 {
						rowsSet = append(rowsSet, rows)
					}
				}

				if len(rowsSet) == 1 {
					result = append(result, rowsSet[0]...)
				}
			}
		} else {
			for _, cell := range nodeDataCells {
				row := dataRow{
					id: uuid.New(),
					cells: map[string]dataCell{
						cell.getName(): cell,
					},
				}

				if rootCell != nil {
					row.cells[rootCell.getName()] = *rootCell
				}

				result = append(result, row)
			}
		}

		// result = append(result, nodeDataCells)
	default:
		return nil, fmt.Errorf("unexpected operand type: %T", operand) //nolint
	}

	return result, nil
}

func (exec *executer) getDataForClassExpr(node treeindex.Noder, operand aqlprocessor.ClassExpression) ([]dataCell, error) {
	var (
		result []dataCell
		err    error
	)

	if node == nil {
		result, err = exec.getDataForClassExpression(operand)
	} else {
		result, err = exec.getDataForClassExpressionnFromNode(node, operand)
	}

	if err != nil {
		return nil, err
	}

	return result, nil
}

func (exec *executer) getDataForClassExpression(operand aqlprocessor.ClassExpression) ([]dataCell, error) {
	cells := []dataCell{}

	switch name := operand.Identifiers[0]; name {
	case "EHR":
		ehrs, err := exec.index.GetEHRs("")
		if err != nil {
			return nil, errors.Wrap(err, "cannot get data source for EHRs")
		}

		for _, ehrNode := range ehrs {
			ok, err := exec.checkNodeByPathPredicate(ehrNode, operand.PathPredicate)
			if err != nil {
				return nil, err
			}

			if !ok {
				continue
			}

			dc := dataCell{
				name: operand.Identifiers[0],
				data: ehrNode,
			}

			if len(operand.Identifiers) > 1 {
				dc.alias = operand.Identifiers[1]
			}

			cells = append(cells, dc)
		}
	default:
		return nil, fmt.Errorf("unexpected data source type: %s", name) //nolint
	}

	return cells, nil
}

func (exec *executer) getDataForClassExpressionnFromNode(node treeindex.Noder, from aqlprocessor.ClassExpression) ([]dataCell, error) {
	result := []dataCell{}

	name := from.Identifiers[0]

	alias := ""
	if len(from.Identifiers) > 1 {
		alias = from.Identifiers[1]
	}

	var container treeindex.Container

	switch node := node.(type) {
	case *treeindex.EHRNode:
		container = node.GetCompositions()
	case *treeindex.CompositionNode:
		sources, err := node.GetDataSourceByName(name)
		if err != nil {
			return nil, errors.Wrap(err, "cannot get sources from compositions node")
		}

		container = sources
	default:
		return nil, fmt.Errorf("not imlemented error for: %T", node) //nolint
	}

	for _, nodes := range container {
		for _, node := range nodes {
			ok, err := exec.checkNodeByPathPredicate(node, from.PathPredicate)
			if err != nil {
				return nil, err
			}

			if !ok {
				continue
			}

			dc := dataCell{
				name:  name,
				alias: alias,
				data:  node,
			}

			result = append(result, dc)
		}
	}

	return result, nil
}

func (exec *executer) checkNodeByPathPredicate(node treeindex.Noder, pathPredicate *aqlprocessor.PathPredicate) (bool, error) {
	if pathPredicate == nil {
		return true, nil
	}

	switch pathPredicate.Type {
	case aqlprocessor.StandartPathPredicate:
		return exec.checkNodeByStandartPathPredicate(node, pathPredicate.StandartPredicate)
	case aqlprocessor.ArchetypedPathPredicate:
		return exec.checkNodeByArchetypePredicate(node, pathPredicate.Archetype)
	case aqlprocessor.NodePathPredicate:
	default:
		return false, fmt.Errorf("unexpected PathPredicate Type: %v", pathPredicate.Type) //nolint
	}

	return false, errors.New("not implemented")
}

func (exec *executer) checkNodeByStandartPathPredicate(node treeindex.Noder, predicate *aqlprocessor.StandartPredicate) (bool, error) {
	val, ok := getValueForPath(predicate.ObjectPath, node)
	if !ok {
		return false, nil
	}

	if predicate.Operand != nil {
		if param := predicate.Operand.Parameter; param != nil {
			paramVal, ok := exec.params[string(*param)]
			if !ok {
				return false, nil
			}

			if reflect.TypeOf(paramVal) == reflect.TypeOf(val) {
				switch pv := paramVal.(type) {
				case string:
					return pv == val.(string), nil
				default:
					return false, fmt.Errorf("unexpected type: %T", paramVal) //nolint
				}
			}
		}

		return false, errors.New("standart predicate operand operations are not implemented")
	}

	return false, errors.New("unexpected standart predicate state")
}

func (exec *executer) checkNodeByArchetypePredicate(node treeindex.Noder, predicate *aqlprocessor.ArchetypePathPredicate) (bool, error) {
	targetArchetypeID := ""
	if predicate.ArchetypeHRID != nil {
		targetArchetypeID = *predicate.ArchetypeHRID
	} else if predicate.Parameter != nil {
		paramVal, ok := exec.params[string(*predicate.Parameter)]
		if !ok {
			return false, nil
		}

		targetArchetypeID, ok = paramVal.(string)
		if !ok {
			return false, nil
		}
	} else {
		return false, errors.New("unexpected archetype predicate state")
	}

	nodeArchetypeID := node.TryGetChild("archetype_node_id")
	if nodeArchetypeID == nil {
		return false, nil
	}

	valueNode, ok := nodeArchetypeID.(*treeindex.ValueNode)
	if !ok {
		return false, errors.New("invalid  archetype_node_id type")
	}

	return valueNode.GetData().(string) == targetArchetypeID, nil
}
