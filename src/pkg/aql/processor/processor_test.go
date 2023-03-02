package processor

import (
	"testing"

	"github.com/bsn-si/IPEHR-gateway/src/pkg/aql/parser"
	"github.com/google/go-cmp/cmp"
)

func TestProcessorTestProcessor(t *testing.T) {
	tests := []struct {
		name    string
		query   string
		want    *Query
		wantErr bool
	}{
		{
			"1. invalid query",
			"SELECT invalid",
			nil,
			true,
		},
		{
			"2. valid query",
			`SELECT 
    o/data[at0002]/events[at0003]/data[at0001]/items[at0004]/value/magnitude AS temperature, 
    o/data[at0002]/events[at0003]/data[at0001]/items[at0004]/value/units AS unit 
FROM 
    EHR e[ehr_id/value='554f896d-faca-4513-bddf-664541146308d']
        CONTAINS Observation o[openEHR-EHR-OBSERVATION.body_temperature-zn.v1] 
WHERE 
    o/data[at0002]/events[at0003]/data[at0001]/items[at0004]/value/magnitude > $temperature 
    AND o/data[at0002]/events[at0003]/data[at0001]/items[at0.63 and name/value='Symptoms']/value/defining_code/code_string=$chills 
ORDER BY temperature DESC
LIMIT 3`,
			&Query{
				Parameters: map[string]*Parameter{
					"temperature": toRef(Parameter("temperature")),
					"chills":      toRef(Parameter("chills")),
				},
				Select: Select{
					SelectExprs: []SelectExpr{
						{
							Path:      "o/data[at0002]/events[at0003]/data[at0001]/items[at0004]/value/magnitude",
							AliasName: "temperature",
							Value: &IdentifiedPathSelectValue{
								Val: IdentifiedPath{
									Identifier: "o",
									ObjectPath: &ObjectPath{
										Paths: []PartPath{
											{
												Identifier: "data", PathPredicate: &PathPredicate{
													Type: NodePathPredicate,
													NodePredicate: &NodePredicate{
														Operator:          NoneOperator,
														ComparisionSymbol: SymNone,
														AtCode:            toRef(AtCode("0002")),
													},
												},
											},
											{
												Identifier: "events", PathPredicate: &PathPredicate{
													Type: NodePathPredicate,
													NodePredicate: &NodePredicate{
														Operator:          NoneOperator,
														ComparisionSymbol: SymNone,
														AtCode:            toRef(AtCode("0003")),
													},
												},
											},
											{
												Identifier: "data", PathPredicate: &PathPredicate{
													Type: NodePathPredicate,
													NodePredicate: &NodePredicate{
														Operator:          NoneOperator,
														ComparisionSymbol: SymNone,
														AtCode:            toRef(AtCode("0001")),
													},
												},
											},
											{
												Identifier: "items", PathPredicate: &PathPredicate{
													Type: NodePathPredicate,
													NodePredicate: &NodePredicate{
														Operator:          NoneOperator,
														ComparisionSymbol: SymNone,
														AtCode:            toRef(AtCode("0004")),
													},
												},
											},
											{Identifier: "value"},
											{Identifier: "magnitude"},
										},
									},
								},
							},
						},
						{
							Path:      "o/data[at0002]/events[at0003]/data[at0001]/items[at0004]/value/units",
							AliasName: "unit",
							Value: &IdentifiedPathSelectValue{
								Val: IdentifiedPath{
									Identifier: "o",
									ObjectPath: &ObjectPath{
										Paths: []PartPath{
											{
												Identifier: "data", PathPredicate: &PathPredicate{
													Type: NodePathPredicate,
													NodePredicate: &NodePredicate{
														Operator:          NoneOperator,
														ComparisionSymbol: SymNone,
														AtCode:            toRef(AtCode("0002")),
													},
												},
											},
											{
												Identifier: "events", PathPredicate: &PathPredicate{
													Type: NodePathPredicate,
													NodePredicate: &NodePredicate{
														Operator:          NoneOperator,
														ComparisionSymbol: SymNone,
														AtCode:            toRef(AtCode("0003")),
													},
												},
											},
											{
												Identifier: "data", PathPredicate: &PathPredicate{
													Type: NodePathPredicate,
													NodePredicate: &NodePredicate{
														Operator:          NoneOperator,
														ComparisionSymbol: SymNone,
														AtCode:            toRef(AtCode("0001")),
													},
												},
											},
											{
												Identifier: "items", PathPredicate: &PathPredicate{
													Type: NodePathPredicate,
													NodePredicate: &NodePredicate{
														Operator:          NoneOperator,
														ComparisionSymbol: SymNone,
														AtCode:            toRef(AtCode("0004")),
													},
												},
											},
											{Identifier: "value"},
											{Identifier: "units"},
										},
									},
								},
							},
						},
					},
				},
				From: From{
					ContainsExpr{
						Operand: ClassExpression{
							Identifiers: []string{"EHR", "e"},
							PathPredicate: &PathPredicate{
								Type: StandartPathPredicate,
								StandartPredicate: &StandartPredicate{
									CMPOperator: SymEQ,
									ObjectPath: &ObjectPath{
										Paths: []PartPath{
											{Identifier: "ehr_id"},
											{Identifier: "value"},
										},
									},
									Operand: &PathPredicateOperand{
										Primitive: &Primitive{Val: "554f896d-faca-4513-bddf-664541146308d", Type: parser.AqlLexerSTRING},
									},
								},
							},
						},
						Contains: []*ContainsExpr{
							{
								Operand: ClassExpression{
									Identifiers: []string{"Observation", "o"},
									PathPredicate: &PathPredicate{
										Type: ArchetypedPathPredicate,
										Archetype: &ArchetypePathPredicate{
											ArchetypeHRID: toRef("openEHR-EHR-OBSERVATION.body_temperature-zn.v1"),
										},
									},
								},
							},
						},
					},
				},
				Where: &Where{
					OperatorType: ANDOperator,
					Next: []*Where{
						{IdentifiedExpr: &IdentifiedExpr{
							IdentifiedPath: &IdentifiedPath{
								Identifier: "o",
								ObjectPath: &ObjectPath{
									Paths: []PartPath{
										{
											Identifier: "data",
											PathPredicate: &PathPredicate{
												Type:          NodePathPredicate,
												NodePredicate: &NodePredicate{Operator: NoneOperator, ComparisionSymbol: SymNone, AtCode: toRef(AtCode("0002"))},
											},
										},
										{
											Identifier: "events",
											PathPredicate: &PathPredicate{
												Type:          NodePathPredicate,
												NodePredicate: &NodePredicate{Operator: NoneOperator, ComparisionSymbol: SymNone, AtCode: toRef(AtCode("0003"))},
											},
										},
										{
											Identifier: "data",
											PathPredicate: &PathPredicate{
												Type:          NodePathPredicate,
												NodePredicate: &NodePredicate{Operator: NoneOperator, ComparisionSymbol: SymNone, AtCode: toRef(AtCode("0001"))},
											},
										},
										{
											Identifier: "items",
											PathPredicate: &PathPredicate{
												Type:          NodePathPredicate,
												NodePredicate: &NodePredicate{Operator: NoneOperator, ComparisionSymbol: SymNone, AtCode: toRef(AtCode("0004"))},
											},
										},
										{Identifier: "value"},
										{Identifier: "magnitude"},
									},
								},
							},
							Terminal:           &Terminal{Parameter: toRef(Parameter("temperature"))},
							ComparisonOperator: toRef(SymGT),
						}},
						{IdentifiedExpr: &IdentifiedExpr{
							IdentifiedPath: &IdentifiedPath{
								Identifier: "o",
								ObjectPath: &ObjectPath{
									Paths: []PartPath{
										{
											Identifier: "data",
											PathPredicate: &PathPredicate{
												Type:          NodePathPredicate,
												NodePredicate: &NodePredicate{Operator: NoneOperator, ComparisionSymbol: SymNone, AtCode: toRef(AtCode("0002"))},
											},
										},
										{
											Identifier: "events",
											PathPredicate: &PathPredicate{
												Type:          NodePathPredicate,
												NodePredicate: &NodePredicate{Operator: NoneOperator, ComparisionSymbol: SymNone, AtCode: toRef(AtCode("0003"))},
											},
										},
										{
											Identifier: "data",
											PathPredicate: &PathPredicate{
												Type:          NodePathPredicate,
												NodePredicate: &NodePredicate{Operator: NoneOperator, ComparisionSymbol: SymNone, AtCode: toRef(AtCode("0001"))},
											},
										},
										{
											Identifier: "items",
											PathPredicate: &PathPredicate{
												Type: NodePathPredicate,
												NodePredicate: &NodePredicate{
													Operator:          ANDOperator,
													ComparisionSymbol: SymNone,
													Next: []*NodePredicate{
														{Operator: NoneOperator, ComparisionSymbol: SymNone, AtCode: toRef(AtCode("0.63"))},
														{
															Operator:          NoneOperator,
															ComparisionSymbol: SymEQ,
															ObjectPath: &ObjectPath{
																Paths: []PartPath{
																	{Identifier: "name"},
																	{Identifier: "value"},
																},
															},
															PathPredicateOperand: &PathPredicateOperand{
																Primitive: &Primitive{Val: "Symptoms", Type: parser.AqlLexerSTRING},
															},
														},
													},
												},
											},
										},
										{Identifier: "value"},
										{Identifier: "defining_code"},
										{Identifier: "code_string"},
									},
								},
							},
							Terminal:           &Terminal{Parameter: toRef(Parameter("chills"))},
							ComparisonOperator: toRef(SymEQ),
						}},
					},
				},
				Order: &Order{
					Orders: []OrderBy{
						{
							IdentifierPath: IdentifiedPath{
								Identifier: "temperature",
							},
							Ordering: DescendingOrdering,
						},
					},
				},
				Limit: &Limit{
					Limit: 3,
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

			opts := cmp.AllowUnexported(
				Query{},
			)
			if diff := cmp.Diff(tt.want, got, opts); diff != "" {
				t.Errorf("Mismatch {+want;-got}:\n\t%s", diff)
			}
		})
	}
}

type Attribute = uint8

const (
	AttributeID Attribute = iota + 1
	AttributeIDEncr
	AttributeKeyEncr
	AttributeDocBaseUIDHash
	AttributeDocUIDEncr
	AttributeDealCid
	AttributeMinerAddress
	AttributeContent
	AttributeContentEncr
	AttributeDescriptionEncr
	AttributePasswordHash
	AttributeTimestamp
)
