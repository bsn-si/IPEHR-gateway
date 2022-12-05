package aqlprocessor

import (
	"testing"

	"github.com/google/go-cmp/cmp"
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
								Identifier: "field",
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
								Identifier: "field",
								PathPredicate: &PathPredicate{
									Type: NodePathPredicate,
									NodePredicate: &NodePredicate{
										Value: "at0003",
									},
								},
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
								Identifier: "field",
								PathPredicate: &PathPredicate{
									Type: NodePathPredicate,
									NodePredicate: &NodePredicate{
										Value: "at0003",
									},
								},
								ObjectPath: &ObjectPath{
									Paths: []PartPath{
										{
											Identifier: "value",
										},
									},
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
								Identifier: "o",
								ObjectPath: &ObjectPath{
									Paths: []PartPath{
										{
											Identifier: "field",
											PathPredicate: &PathPredicate{
												Type: NodePathPredicate,
												NodePredicate: &NodePredicate{
													Value: "at0003",
												},
											},
										},
										{Identifier: "value1"},
										{Identifier: "value2"},
									},
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

			tt.want.From = got.From
			if diff := cmp.Diff(tt.want, got); diff != "" {
				t.Errorf("mismatch {+want;-got}:\n\t%s", diff)
			}
		})
	}
}
