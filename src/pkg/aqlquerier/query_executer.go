package aqlquerier

import (
	"database/sql/driver"
	"fmt"

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
		return nil, errors.Wrap(err, "cannot find data sources")
	}

	// handle WHERE block
	dataSources, err = exec.filterSources(dataSources)
	if err != nil {
		return nil, errors.Wrap(err, "cannot filter data sources")
	}

	// handle SELECT block
	rows, err := exec.queryData(dataSources)
	if err != nil {
		return nil, errors.Wrap(err, "cannot query rows from data sources")
	}

	rows, err = exec.orderRows(rows)
	if err != nil {
		return nil, errors.Wrap(err, "cannot order rows")
	}

	return exec.limitRows(rows), nil
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

	return processWhere(exec.query.Where, sources)
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
			val := exec.getPrimitiveColumnValue(slct)
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

func (exec *executer) orderRows(rows *Rows) (*Rows, error) {
	// handle ORDER block
	//TODO: add order logic
	return rows, nil
}

func (exec *executer) limitRows(rows *Rows) *Rows {
	if exec.query.Limit == nil {
		return rows
	}

	limit := exec.query.Limit.Limit
	offset := exec.query.Limit.Offset

	if offset >= 0 {
		if offset > len(rows.rows) {
			offset = len(rows.rows)
		}

		rows.rows = rows.rows[offset:]
	}

	if limit < len(rows.rows) {
		rows.rows = rows.rows[:limit]
	}

	return rows
}

func (exec *executer) getPrimitiveColumnValue(prim *aqlprocessor.PrimitiveSelectValue) driver.Value {
	if prim == nil {
		return nil
	}

	return prim.Val.Val
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
