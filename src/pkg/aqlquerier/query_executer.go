package aqlquerier

import (
	"database/sql/driver"
	"encoding/json"
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

	log.Println("data sources: ", dataSources)

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
	data, _ := json.Marshal(exec.query.From)
	log.Println(string(data))

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

	for _, _ = range sources {
		row := Row{}
		for _, selectExpr := range exec.query.Select.SelectExprs {
			switch slct := selectExpr.Value.(type) {
			case *aqlprocessor.IdentifiedPathSelectValue:
				return nil, errors.New("Identified path not implemented")
			case *aqlprocessor.PrimitiveSelectValue:
				data, _ := json.Marshal(slct)
				val, err := exec.getPrimitiveColumnValue(slct)
				if err != nil {
					return nil, errors.Wrap(err, "cannot get primitive select value")
				}

				log.Printf("PRIMITIVE: %v, %T, %v", string(data), slct.Val.Val, val)
				row.values = append(row.values, val)
			case *aqlprocessor.AggregateFunctionCallSelectValue:
				return nil, errors.New("Aggregation function call not implemented")
			case *aqlprocessor.FunctionCallSelectValue:
				return nil, errors.New("Function call not implemented")
			default:
				return nil, errors.New("Unexpected SelectExpr type")
			}
		}

		log.Println(row)
		rows.rows = append(rows.rows, row)
	}

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
