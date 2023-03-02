package driver

import (
	"database/sql/driver"
	"fmt"

	aqlprocessor "github.com/bsn-si/IPEHR-gateway/src/pkg/aql/processor"
	"github.com/bsn-si/IPEHR-gateway/src/pkg/errors"
)

func (exec *executer) queryData(sources dataRows) (*Rows, error) {
	if len(sources) == 0 {
		return &Rows{}, nil
	}

	result := &Rows{
		rows: []Row{},
	}

	//TODO: add DISTINCT handling
	// exec.query.Select.Distinct

	for _, dataRow := range sources {
		row := Row{
			values: []interface{}{},
		}

		for _, selectExpr := range exec.query.Select.SelectExprs {
			switch slct := selectExpr.Value.(type) {
			case *aqlprocessor.IdentifiedPathSelectValue:
				{
					var val any
					ip := slct.Val

					indexNode, ok := dataRow.cells[slct.Val.Identifier]
					if ok {
						if ip.ObjectPath != nil {
							val, _ = getValueForPath(ip.ObjectPath, indexNode.data)
						} else {
							return nil, errors.New("unsupported select expresion format")
						}
					}

					row.values = append(row.values, val)
				}
			case *aqlprocessor.PrimitiveSelectValue:
				{
					val := exec.getPrimitiveColumnValue(slct)
					row.values = append(row.values, val)
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

		result.rows = append(result.rows, row)
	}

	return exec.fillColumns(result), nil
}

func (exec *executer) getPrimitiveColumnValue(prim *aqlprocessor.PrimitiveSelectValue) driver.Value {
	if prim == nil {
		return nil
	}

	return prim.Val.Val
}

func (exec *executer) fillColumns(rows *Rows) *Rows {
	for i, se := range exec.query.Select.SelectExprs {
		c := Column{
			Path: se.Path,
			Name: se.AliasName,
		}

		if c.Name == "" {
			c.Name = fmt.Sprintf("#%d", i)
		}

		rows.columns = append(rows.columns, c)
	}

	return rows
}
