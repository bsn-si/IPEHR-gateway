package aqlprocessor

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestProcessor_SelectIdentifiedPath(t *testing.T) {
	tests := []struct {
		name    string
		query   string
		want    Query
		wantErr bool
	}{
		{
			"1. select field",
			`SELECT field FROM EHR`,
			Query{
				Select: Select{
					SelectExprs: []SelectExpr{{
						Value: &IdentifiedPathSelectValue{
							Val: IdentifiedPath{
								Name: "field",
							},
						},
					}},
				},
			},
			false,
		},
		{
			"2. select field with path_predicate",
			`SELECT field[at0003] FROM EHR`,
			Query{
				Select: Select{
					SelectExprs: []SelectExpr{{
						Value: &IdentifiedPathSelectValue{
							Val: IdentifiedPath{
								Name:          "field",
								PathPredicate: "at0003",
							},
						},
					}},
				},
			},
			false,
		},
		{
			"3. select field[at0003]/value",
			`SELECT field[at0003]/value FROM EHR`,
			Query{
				Select: Select{
					SelectExprs: []SelectExpr{{
						Value: &IdentifiedPathSelectValue{
							Val: IdentifiedPath{
								Name:          "field",
								PathPredicate: "at0003",
								Paths: []ObjectPath{
									{Name: "value"},
								},
							},
						},
					}},
				},
			},
			false,
		},
		{
			"4. select o/field[at0003]/value1/value2 with alisas",
			`SELECT o/field[at0003]/value1/value2 AS new_name FROM EHR`,
			Query{
				Select: Select{
					SelectExprs: []SelectExpr{{
						Value: &IdentifiedPathSelectValue{
							Val: IdentifiedPath{
								Name: "o",
								Paths: []ObjectPath{
									{Name: "field", PathPredicate: "at0003"},
									{Name: "value1"},
									{Name: "value2"},
								},
							},
						},
						AliasName: "new_name",
					}},
				},
			},
			false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := NewAqlProcessor(tt.query)
			got, err := p.Process()
			if (err != nil) != tt.wantErr {
				t.Errorf("Process Query err: '%v', want: %v", err, tt.wantErr)
			}

			if !tt.wantErr && assert.NoError(t, err) {
				assert.Equal(t, tt.want, got)
			}
		})
	}
}
