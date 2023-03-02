package processor

import (
	"testing"

	"github.com/bsn-si/IPEHR-gateway/src/pkg/aql/parser"
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
						Parameter: toRef(Parameter("nameValue")),
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
							Terminal: &Terminal{
								Parameter: toRef(Parameter("nameValue")),
							},
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
							Terminal: &Terminal{
								Parameter: toRef(Parameter("templateId")),
							},
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
				o/data[at0001,id123]/events[at0006,'str']/data[at0003,at0002]/items[at0004,$some_parameter]/value/value >= 140`,
			&Where{
				OperatorType: ANDOperator,
				Next: []*Where{
					{
						Brackets: true,
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
											Terminal: &Terminal{
												Parameter: toRef(Parameter("nameValue")),
											},
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
											Terminal: &Terminal{
												Parameter: toRef(Parameter("templateId")),
											},
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
												Type: NodePathPredicate,
												NodePredicate: &NodePredicate{
													AtCode:            toRef(AtCode("0001")),
													Operator:          NoneOperator,
													ComparisionSymbol: SymNone,
													AdditionalData: &NodePredicateAdditionalData{
														IDCode: toRef(IDCode("123")),
													},
												},
											},
										},
										{
											Identifier: "events",
											PathPredicate: &PathPredicate{
												Type: NodePathPredicate,
												NodePredicate: &NodePredicate{
													AtCode:            toRef(AtCode("0006")),
													Operator:          NoneOperator,
													ComparisionSymbol: SymNone,
													AdditionalData: &NodePredicateAdditionalData{
														String: toRef("str"),
													},
												},
											},
										},
										{
											Identifier: "data",
											PathPredicate: &PathPredicate{
												Type: NodePathPredicate,
												NodePredicate: &NodePredicate{
													AtCode:            toRef(AtCode("0003")),
													Operator:          NoneOperator,
													ComparisionSymbol: SymNone,
													AdditionalData: &NodePredicateAdditionalData{
														AtCode: toRef(AtCode("0002")),
													},
												},
											},
										},
										{
											Identifier: "items",
											PathPredicate: &PathPredicate{
												Type: NodePathPredicate,
												NodePredicate: &NodePredicate{
													AtCode:            toRef(AtCode("0004")),
													Operator:          NoneOperator,
													ComparisionSymbol: SymNone,
													AdditionalData: &NodePredicateAdditionalData{
														Parameter: toRef(Parameter("some_parameter")),
													},
												},
											},
										},
										{Identifier: "value"},
										{Identifier: "value"},
									},
								},
							},
							ComparisonOperator: toRef(SymGE),
							Terminal: &Terminal{
								Primitive: &Primitive{
									Val:  140,
									Type: parser.AqlLexerINTEGER,
								},
							},
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
