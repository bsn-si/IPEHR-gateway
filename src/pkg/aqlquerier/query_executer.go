package aqlquerier

import (
	"database/sql/driver"
	"fmt"
	"log"

	"hms/gateway/pkg/aqlprocessor"
	"hms/gateway/pkg/errors"
	"hms/gateway/pkg/storage/treeindex"
)

type executer struct {
	query  *aqlprocessor.Query
	params map[string]driver.Value

	index *treeindex.Tree
}

func (exec *executer) run() (*Rows, error) {
	// handle FROM block
	dataSources, err := exec.findSources()
	if err != nil {
		return nil, err
	}

	// handle WHERE block
	dataSources, err = exec.filterSources(dataSources)
	if err != nil {
		return nil, err
	}

	// handle SELECT block
	rows, err := exec.queryData(dataSources)
	if err != nil {
		return nil, err
	}

	// handle ORDER block
	//TODO: add order logic

	return rows, nil
}

func (exec *executer) findSources() (map[string]dataSource, error) {
	from := exec.query.From
	if len(from.Contains) > 0 || from.Operator != nil {
		return nil, errors.New("not implemented")
	}

	source := map[string]dataSource{}

	if from.Operand != nil {
		switch operand := from.Operand.(type) {
		case aqlprocessor.ClassExpression:
			if operand.PathPredicate != nil {
				return nil, errors.New("not implemented")
			}

			ds := dataSource{
				name: operand.Identifiers[0],
			}

			if len(operand.Identifiers) > 1 {
				ds.alias = operand.Identifiers[1]
			}

			data, err := exec.index.GetDataSourceByName(ds.name)
			if err != nil {
				return nil, errors.Wrap(err, "cannot get data source by name")
			}

			ds.data = data

			source[ds.getName()] = ds
		case aqlprocessor.VersionClassExpr:
			return nil, errors.New("not implemented")
		default:
			return nil, fmt.Errorf("unexpected FROM.Operand type: %T", operand) // nolint
		}
	}

	return source, nil
}

func (exec *executer) filterSources(sources map[string]dataSource) (map[string]dataSource, error) {
	if exec.query.Where == nil {
		return sources, nil
	}

	where := exec.query.Where

	return processWhere(where, sources)
}

func processWhere(where *aqlprocessor.Where, sources map[string]dataSource) (map[string]dataSource, error) {
	if ie := where.IdentifiedExpr; ie != nil {
		result := map[string]dataSource{}

		if ie.ComparisonOperator != nil && ie.IdentifiedPath != nil {
			for key, source := range sources {
				ip := *ie.IdentifiedPath
				// val, ok := getValueForPath(*ie.IdentifiedPath, source)
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
		}

		return result, nil
	} else if where.OperatorType != aqlprocessor.NoneOperator {
		if len(where.Next) != 2 {
			return nil, errors.New("unexpected where conditions count")
		}

		result := map[string]dataSource{}
		results := make([]map[string]dataSource, 2)
		for i, where := range where.Next {
			s, err := processWhere(where, sources)
			if err != nil {
				return nil, errors.Wrap(err, "cannot filter inner WHERE conditions")
			}

			results[i] = s
		}

		switch where.OperatorType {
		case aqlprocessor.ANDOperator:
			result = mergeDataSourcesAND(results[0], results[1])
		case aqlprocessor.OROperator:
			result = mergeDataSourcesOR(results[0], results[1])
		default:
			log.Println("UNEXPECTED ", where.OperatorType)
		}

		return result, nil
	}

	return sources, nil
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

func (exec *executer) queryData(sources map[string]dataSource) (*Rows, error) {
	rows := &Rows{}

	//TODO: add DISTINCT handling
	// exec.query.Select.Distinct

	// for _, source := range sources {
	// log.Println(source.name, source.alias, source.data)

	// row := Row{}

	primitivesRow := Row{}

	for _, selectExpr := range exec.query.Select.SelectExprs {
		switch slct := selectExpr.Value.(type) {
		case *aqlprocessor.IdentifiedPathSelectValue:
			columnValues, err := getDataByIdentifiedPath(slct.Val, sources)
			if err != nil {
				return nil, errors.Wrap(err, "cannot get data for identified data")
			}

			for _, val := range columnValues {
				row := Row{
					values: []interface{}{val},
				}
				rows.rows = append(rows.rows, row)
			}
		case *aqlprocessor.PrimitiveSelectValue:
			val, err := exec.getPrimitiveColumnValue(slct)
			if err != nil {
				return nil, errors.Wrap(err, "cannot get primitive select value")
			}

			primitivesRow.values = append(primitivesRow.values, val)
		case *aqlprocessor.AggregateFunctionCallSelectValue:
			return nil, errors.New("Aggregation function call not implemented")
		case *aqlprocessor.FunctionCallSelectValue:
			return nil, errors.New("Function call not implemented")
		default:
			return nil, errors.New("Unexpected SelectExpr type")
		}
	}

	rows.rows = append(rows.rows, primitivesRow)

	return exec.fillColumns(rows), nil
}

func (exec *executer) getPrimitiveColumnValue(prim *aqlprocessor.PrimitiveSelectValue) (driver.Value, error) {
	if prim == nil {
		return nil, nil
	}

	return prim.Val.Val, nil
}

func (exec *executer) fillColumns(rows *Rows) *Rows {
	for _, se := range exec.query.Select.SelectExprs {
		rows.columns = append(rows.columns, se.AliasName)
	}

	return rows
}

func getDataByIdentifiedPath(ip aqlprocessor.IdentifiedPath, sources map[string]dataSource) ([]any, error) {
	result := []any{}

	source, ok := sources[ip.Identifier]
	if !ok {
		return nil, fmt.Errorf("unexpected identifier %v", ip.Identifier) //nolint
	}

	for _, indexNodes := range source.data {
		for _, indexNode := range indexNodes {
			if ip.ObjectPath != nil {
				if resultData, ok := getValueForPath(ip.ObjectPath, indexNode); ok {
					result = append(result, resultData)
				}
			}
		}
	}

	return result, nil
}

func getValueForPath(path *aqlprocessor.ObjectPath, node treeindex.Noder) (any, bool) {
	index := 0
	queue := []treeindex.Noder{node}

	for len(queue) > 0 {
		if index >= len(path.Paths) {
			return nil, false
		}

		path := path.Paths[index]

		node := queue[0]
		queue = queue[1:]

		switch node := node.(type) {
		case *treeindex.ObjectNode:
			{
				nextNode := node.TryGetChild(path.Identifier)
				if nextNode == nil {
					continue
				}

				switch nextNode.GetNodeType() {
				case treeindex.ObjectNodeType:
					if path.PathPredicate != nil && path.PathPredicate.Type == aqlprocessor.NodePathPredicate {
						if np := path.PathPredicate.NodePredicate; np.AtCode != nil && nextNode.GetID() == np.AtCode.ToString() {
							index++
						}
					}
				case treeindex.DataValueNodeType:
					index++
				}

				queue = append(queue, nextNode)
			}
		case *treeindex.SliceNode:
			if path.PathPredicate != nil && path.PathPredicate.Type == aqlprocessor.NodePathPredicate {
				np := path.PathPredicate.NodePredicate

				if np.AtCode != nil {
					nextNode := node.TryGetChild(np.AtCode.ToString())
					if nextNode != nil {
						queue = append(queue, nextNode)
						index++
					}
				}
			}
		case *treeindex.DataValueNode:
			if valueNode := node.TryGetChild(path.Identifier); valueNode != nil {
				queue = append(queue, valueNode)
			}
		case *treeindex.ValueNode:
			return node.GetData(), true
		default:
		}
	}

	return nil, false
}

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
