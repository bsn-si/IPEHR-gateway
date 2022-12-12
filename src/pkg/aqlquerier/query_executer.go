package aqlquerier

import (
	"database/sql/driver"
	"fmt"
	"log"
	"os"
	"strings"

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

func (exec *executer) queryData(sources map[string]dataSource) (*Rows, error) {
	rows := &Rows{}

	//TODO: add DISTINCT handling
	// exec.query.Select.Distinct

	// for _, source := range sources {
	// log.Println(source.name, source.alias, source.data)

	// row := Row{}

	primitives := Row{}

	for _, selectExpr := range exec.query.Select.SelectExprs {
		switch slct := selectExpr.Value.(type) {
		case *aqlprocessor.IdentifiedPathSelectValue:
			columnValues, err := exec.getDataByIdentifiedPath(slct, sources)
			if err != nil {
				return nil, errors.Wrap(err, "cannot get data fo~r identified data")
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

			primitives.values = append(primitives.values, val)
		case *aqlprocessor.AggregateFunctionCallSelectValue:
			return nil, errors.New("Aggregation function call not implemented")
		case *aqlprocessor.FunctionCallSelectValue:
			return nil, errors.New("Function call not implemented")
		default:
			return nil, errors.New("Unexpected SelectExpr type")
		}
		// }
	}

	rows.rows = append(rows.rows, primitives)

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

func (exec *executer) getDataByIdentifiedPath(slctExpr *aqlprocessor.IdentifiedPathSelectValue, sources map[string]dataSource) ([]any, error) {
	result := []any{}
	logger := log.New(os.Stderr, "\t[getDataByIdentifiedPath]\t", log.LstdFlags)
	logger.Println()

	selectPath := slctExpr.Val

	source, ok := sources[selectPath.Identifier]
	if !ok {
		return nil, errors.New("unexpected identifier " + selectPath.Identifier)
	}

	// logger.Println(source)
	// logger.Println(source.data.Len())
	// data, _ := json.MarshalIndent(val, "", "\t")
	// logger.Println(string(data))

	for _, indexNodes := range source.data {
		// data, _ := json.MarshalIndent(nodes, "", "\t")
		// log.Println("data", string(data))

		for _, indexNode := range indexNodes {

			if selectPath.ObjectPath != nil {
				// logger.Println("NODE:", name, name == selectPath.ObjectPath.Paths[0].Identifier)
				if resultData, ok := getValueForPath(selectPath.ObjectPath, indexNode); ok {
					logger.Printf("%v = %T", resultData, resultData)
					result = append(result, resultData)
				}
			}
			// indexNode.ForEach(func(name string, node treeindex.Noder) bool {
			//
			//
			// return true
			// })
		}
	}

	return result, nil
}

func getValueForPath(path *aqlprocessor.ObjectPath, node treeindex.Noder) (any, bool) {
	var result any
	found := false

	logger := log.New(os.Stdout, "\t[getValueForPath]\t", log.LstdFlags)

	// logger.Println("NODE_ID", node.GetID())
	// logger.Println(path.Paths)

	var walkFunc func(name string, node treeindex.Noder) bool

	offset := 0
	idx := 0
	walkFunc = func(name string, node treeindex.Noder) bool {
		offset++
		strOffset := strings.Repeat("\t", offset)

		if idx >= len(path.Paths) {
			return false
		}

		p := path.Paths[idx]

		// logger.Println(strOffset, "p.Identifier", p.Identifier, p.Identifier == name, p.PathPredicate != nil)

		if p.PathPredicate != nil {
			switch p.PathPredicate.Type {
			case aqlprocessor.StandartPathPredicate:
				logger.Println(strOffset, "standart")
			case aqlprocessor.ArchetypedPathPredicate:
				logger.Println(strOffset, "archetype")
			case aqlprocessor.NodePathPredicate:
				np := p.PathPredicate.NodePredicate

				if np.AtCode != nil {
					logger.Println(strOffset, "node.at_code", np.AtCode.ToString(), "node_id", node.GetID())

					if name == p.Identifier && np.AtCode.ToString() == node.GetID() {
						idx++
					}
				}
			default:
				return false
			}
		} else if name == p.Identifier {
			idx++
		}

		switch node := node.(type) {
		case *treeindex.ValueNode:
			if name == p.Identifier {
				result = node.GetData()
				found = true
				// logger.Println(strOffset, name, "=", node.GetData())

				return false
			}
		case *treeindex.SliceNode:
			// logger.Println(strOffset, "slice_node")
			if p.PathPredicate != nil {
				switch p.PathPredicate.Type {
				case aqlprocessor.StandartPathPredicate:
					logger.Println(strOffset, "standart")
				case aqlprocessor.ArchetypedPathPredicate:
					logger.Println(strOffset, "archetype")
				case aqlprocessor.NodePathPredicate:
					np := p.PathPredicate.NodePredicate

					if np.AtCode != nil {
						newNode := node.TryGetChild(np.AtCode.ToString())
						if newNode != nil {
							idx++
							newNode.ForEach(walkFunc)
						}
					}
				default:
					return false
				}
			}
		default:
			logger.Println(strOffset, name, node.GetID())
			node.ForEach(walkFunc)
		}

		offset--
		return true
	}

	node.ForEach(walkFunc)

	logger.Println("RESULT = ", result)
	return result, found
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
