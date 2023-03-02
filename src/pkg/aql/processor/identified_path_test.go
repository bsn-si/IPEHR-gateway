package processor

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestProcessor_SelectIdentifiedPath(t *testing.T) {
	tests := []struct {
		name    string
		query   string
		want    Select
		wantErr bool
	}{
		{
			"1. select field",
			`SELECT field FROM EHR`,
			Select{
				SelectExprs: []SelectExpr{{
					Path: "field",
					Value: &IdentifiedPathSelectValue{
						Val: IdentifiedPath{
							Identifier: "field",
						},
					},
				}},
			},
			false,
		},
		{
			"2. select field with path_predicate",
			`SELECT field[at0003] FROM EHR`,
			Select{
				SelectExprs: []SelectExpr{{
					Path: "field[at0003]",
					Value: &IdentifiedPathSelectValue{
						Val: IdentifiedPath{
							Identifier: "field",
							PathPredicate: &PathPredicate{
								Type: NodePathPredicate,
								NodePredicate: &NodePredicate{
									AtCode:            toRef(AtCode("0003")),
									Operator:          NoneOperator,
									ComparisionSymbol: SymNone,
								},
							},
						},
					},
				}},
			},
			false,
		},
		{
			"3. select field[at0003]/value",
			`SELECT field[at0003]/value FROM EHR`,
			Select{
				SelectExprs: []SelectExpr{{
					Value: &IdentifiedPathSelectValue{
						Val: IdentifiedPath{
							Identifier: "field",
							PathPredicate: &PathPredicate{
								Type: NodePathPredicate,
								NodePredicate: &NodePredicate{
									AtCode:            toRef(AtCode("0003")),
									Operator:          NoneOperator,
									ComparisionSymbol: SymNone,
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
					Path: "field[at0003]/value",
				}},
			},
			false,
		},
		{
			"4. select o/field[at0003]/value1/value2 with alisas",
			`SELECT o/field[at0003]/value1/value2 AS new_name FROM EHR`,
			Select{
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
												AtCode:            toRef(AtCode("0003")),
												Operator:          NoneOperator,
												ComparisionSymbol: SymNone,
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
					Path:      "o/field[at0003]/value1/value2",
				}},
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

			if got == nil {
				return
			}

			if diff := cmp.Diff(tt.want, got.Select); diff != "" {
				t.Errorf("mismatch {+want;-got}:\n\t%s", diff)
			}
		})
	}
}
