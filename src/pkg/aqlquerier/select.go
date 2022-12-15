package aqlquerier

import (
	"database/sql/driver"
	"fmt"
	"hms/gateway/pkg/aqlprocessor"
	"hms/gateway/pkg/errors"
)

func (exec *executer) queryData(sources map[string]dataSource) (*Rows, error) {
	//TODO: add DISTINCT handling
	// exec.query.Select.Distinct

	// for _, source := range sources {
	// log.Println(source.name, source.alias, source.data)

	rows := []Row{}

	primitivesRow := Row{}

	for _, selectExpr := range exec.query.Select.SelectExprs {
		switch slct := selectExpr.Value.(type) {
		case *aqlprocessor.IdentifiedPathSelectValue:
			{
				columnValues, err := getDataByIdentifiedPath(slct.Val, sources)
				if err != nil {
					return nil, errors.Wrap(err, "cannot get data for identified data")
				}

				for _, val := range columnValues {
					row := Row{
						values: []interface{}{val},
					}

					rows = append(rows, row)
				}
			}
		case *aqlprocessor.PrimitiveSelectValue:
			{
				val := exec.getPrimitiveColumnValue(slct)
				primitivesRow.values = append(primitivesRow.values, val)
			}
		case *aqlprocessor.AggregateFunctionCallSelectValue:
			{
				return nil, errors.New("Aggregation function call not implemented")
			}
		case *aqlprocessor.FunctionCallSelectValue:
			{
				return nil, errors.New("Function call not implemented")
			}
		default:
			return nil, errors.New("Unexpected SelectExpr type")
		}
	}

	rows = append(rows, primitivesRow)

	result := &Rows{
		rows: rows,
	}

	return exec.fillColumns(result), nil
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
