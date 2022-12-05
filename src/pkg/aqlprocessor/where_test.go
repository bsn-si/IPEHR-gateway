package aqlprocessor

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestProcessor_Where(t *testing.T) {
	tests := []struct {
		name    string
		query   string
		want    *Where
		wantErr bool
	}{
		{
			"1. empty where",
			`SELECT val FROM EHR`,
			nil,
			false,
		},
		{
			"2. simple where",
			`SELECT val FROM EHR e WHERE c/name/value=$nameValue`,
			&Where{
				IdentifiedExpr: &IdentifiedExpr{
					IdentifiedPath: &IdentifiedPath{
						Identifier: "c",
						ObjectPath: &ObjectPath{
							Paths: []PartPath{
								{Identifier: "name"},
								{Identifier: "value"},
							},
						},
					},
					ComparisonOperator: toRef(SymEQ),
					Terminal: &Terminal{
						Text: "$nameValue",
					},
				},
			},
			false,
		},
		{
			"3. Where with AND",
			`SELECT val FROM EHR
			WHERE
				c/name/value=$nameValue AND c/archetype_details/template_id/value>=$templateId`,
			&Where{
				OperatorType: ANDOperator,
				Next: []*Where{
					{
						IdentifiedExpr: &IdentifiedExpr{
							IdentifiedPath: &IdentifiedPath{
								Identifier: "c",
								ObjectPath: &ObjectPath{
									Paths: []PartPath{
										{Identifier: "name"},
										{Identifier: "value"},
									},
								},
							},
							ComparisonOperator: toRef(SymEQ),
							Terminal:           &Terminal{Text: "$nameValue"},
						},
					},
					{
						IdentifiedExpr: &IdentifiedExpr{
							IdentifiedPath: &IdentifiedPath{
								Identifier: "c",
								ObjectPath: &ObjectPath{
									Paths: []PartPath{
										{Identifier: "archetype_details"},
										{Identifier: "template_id"},
										{Identifier: "value"},
									},
								},
							},
							ComparisonOperator: toRef(SymGE),
							Terminal:           &Terminal{Text: "$templateId"},
						},
					},
				},
			},
			false,
		},
		{
			"4. Where with AND and OR",
			`SELECT val FROM EHR
			WHERE
				(c/name/value = $nameValue OR c/archetype_details/template_id/value = $templateId) AND
				o/data[at0001]/events[at0006]/data[at0003]/items[at0004]/value/value >= 140`,
			&Where{
				OperatorType: ANDOperator,
				Next: []*Where{
					{
						Next: []*Where{
							{
								OperatorType: OROperator,
								Next: []*Where{
									{
										IdentifiedExpr: &IdentifiedExpr{
											IdentifiedPath: &IdentifiedPath{
												Identifier: "c",
												ObjectPath: &ObjectPath{
													Paths: []PartPath{
														{Identifier: "name"},
														{Identifier: "value"},
													},
												},
											},
											ComparisonOperator: toRef(SymEQ),
											Terminal:           &Terminal{Text: "$nameValue"},
										},
									},
									{
										IdentifiedExpr: &IdentifiedExpr{
											IdentifiedPath: &IdentifiedPath{
												Identifier: "c",
												ObjectPath: &ObjectPath{
													Paths: []PartPath{
														{Identifier: "archetype_details"},
														{Identifier: "template_id"},
														{Identifier: "value"},
													},
												},
											},
											Terminal:           &Terminal{Text: "$templateId"},
											ComparisonOperator: toRef(SymEQ),
										},
									},
								},
							},
						},
					},
					{
						IdentifiedExpr: &IdentifiedExpr{
							IdentifiedPath: &IdentifiedPath{
								Identifier: "o",
								ObjectPath: &ObjectPath{
									Paths: []PartPath{
										{
											Identifier: "data", PathPredicate: &PathPredicate{
												Type:          NodePathPredicate,
												NodePredicate: &NodePredicate{Value: "at0001"},
											},
										},
										{
											Identifier: "events",
											PathPredicate: &PathPredicate{
												Type:          NodePathPredicate,
												NodePredicate: &NodePredicate{Value: "at0006"},
											},
										},
										{
											Identifier: "data",
											PathPredicate: &PathPredicate{
												Type:          NodePathPredicate,
												NodePredicate: &NodePredicate{Value: "at0003"},
											},
										},
										{
											Identifier: "items",
											PathPredicate: &PathPredicate{
												Type:          NodePathPredicate,
												NodePredicate: &NodePredicate{Value: "at0004"},
											},
										},
										{Identifier: "value"},
										{Identifier: "value"},
									},
								},
							},
							ComparisonOperator: toRef(SymGE),
							Terminal:           &Terminal{Text: "140"},
						},
					},
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

			if diff := cmp.Diff(tt.want, got.Where); diff != "" {
				t.Errorf("Mismatch {+want;-got}:\n\t%s", diff)
			}
		})
	}
}
