package aqlquerier

import (
	"database/sql/driver"
	"hms/gateway/pkg/aqlprocessor"
	"hms/gateway/pkg/errors"
)

func (exec *executer) queryData(sources map[string]dataSource) (*Rows, error) {
	//TODO: add DISTINCT handling
	// exec.query.Select.Distinct

	// for _, source := range sources {
	// log.Println(source.name, source.alias, source.data)

	rows := []Row{}

	for sourceName, source := range sources {
		for _, indexNodes := range source.data {
			for _, indexNode := range indexNodes {
				row := Row{
					values: []interface{}{},
				}

				anyNotNillValue := false

				for _, selectExpr := range exec.query.Select.SelectExprs {
					switch slct := selectExpr.Value.(type) {
					case *aqlprocessor.IdentifiedPathSelectValue:
						{
							ip := slct.Val
							var val any
							if sourceName == ip.Identifier {
								if ip.ObjectPath != nil {
									ok := false
									val, ok = getValueForPath(ip.ObjectPath, indexNode)
									if !ok {
										val = nil
									} else {
										anyNotNillValue = true
									}
								}
							}

							row.values = append(row.values, val)
						}
					case *aqlprocessor.PrimitiveSelectValue:
						{
							anyNotNillValue = true
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

				if anyNotNillValue {
					rows = append(rows, row)
				}
			}
		}
	}

	result := &Rows{
		rows: rows,
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
	for _, se := range exec.query.Select.SelectExprs {
		rows.columns = append(rows.columns, se.AliasName)
	}

	return rows
}
