package processor

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestProcessor_From(t *testing.T) {
	tests := []struct {
		name    string
		query   string
		want    From
		wantErr bool
	}{
		{
			"1. invalid FROM",
			"SELECT invalid FROM ",
			From{},
			true,
		},
		{
			"2. Simple FROM",
			"SELECT val FROM EHR",
			From{
				ContainsExpr: ContainsExpr{
					Operand: ClassExpression{
						Identifiers: []string{"EHR"},
					},
				},
			},
			false,
		},
		{
			"3. FROM with vartiabel",
			"SELECT val FROM EHR e",
			From{
				ContainsExpr: ContainsExpr{
					Operand: ClassExpression{
						Identifiers: []string{"EHR", "e"},
					},
				},
			},

			false,
		},
		{
			"4. FROM with path predicate",
			"SELECT val FROM EHR e [ehr_id/value=$ehrUid]",
			From{
				ContainsExpr: ContainsExpr{
					Operand: ClassExpression{
						Identifiers: []string{"EHR", "e"},
						PathPredicate: &PathPredicate{
							Type: StandartPathPredicate,
							StandartPredicate: &StandartPredicate{
								ObjectPath: &ObjectPath{
									Paths: []PartPath{
										{Identifier: "ehr_id"},
										{Identifier: "value"},
									},
								},
								CMPOperator: SymEQ,
								Operand: &PathPredicateOperand{
									Parameter: toRef(Parameter("ehrUid")),
								},
							},
						},
					},
				},
			},
			false,
		},
		{
			"5. FROM with CONTANIS",
			`SELECT val FROM EHR e CONTAINS COMPOSITION c[openEHR-EHR-COMPOSITION.report.v1]`,
			From{
				ContainsExpr: ContainsExpr{
					Operand: ClassExpression{
						Identifiers:   []string{"EHR", "e"},
						PathPredicate: nil,
					},
					Contains: []*ContainsExpr{
						{
							Operand: ClassExpression{
								Identifiers: []string{"COMPOSITION", "c"},
								PathPredicate: &PathPredicate{
									Type: ArchetypedPathPredicate,
									Archetype: &ArchetypePathPredicate{
										ArchetypeHRID: toRef("openEHR-EHR-COMPOSITION.report.v1"),
									},
								},
							},
						},
					},
				},
			},
			false,
		},
		{
			"6, FROM with CONTAINS and AND",
			`SELECT val FROM EHR e
				CONTAINS COMPOSITION c [openEHR-EHR-COMPOSITION.referral.v1] AND COMPOSITION c1 [openEHR-EHR-COMPOSITION.report.v1]`,
			From{
				ContainsExpr: ContainsExpr{
					Operand: ClassExpression{
						Identifiers: []string{"EHR", "e"},
					},
					Contains: []*ContainsExpr{
						{
							Operator: toRef(ANDOperator),
							Contains: []*ContainsExpr{
								{
									Operand: ClassExpression{
										Identifiers: []string{"COMPOSITION", "c"},
										PathPredicate: &PathPredicate{
											Type: ArchetypedPathPredicate,
											Archetype: &ArchetypePathPredicate{
												ArchetypeHRID: toRef("openEHR-EHR-COMPOSITION.referral.v1"),
											},
										},
									},
								},
								{
									Operand: ClassExpression{
										Identifiers: []string{"COMPOSITION", "c1"},
										PathPredicate: &PathPredicate{
											Type: ArchetypedPathPredicate,
											Archetype: &ArchetypePathPredicate{
												ArchetypeHRID: toRef("openEHR-EHR-COMPOSITION.report.v1"),
											},
										},
									},
								},
							},
						},
					},
				},
			},
			false,
		},
		{
			"7. FROM with OR operator",
			`SELECT val
			 FROM EHR e
				CONTAINS COMPOSITION c [openEHR-EHR-COMPOSITION.referral.v1]
					CONTAINS (OBSERVATION o [openEHR-EHR-OBSERVATION.laboratory-hba1c.v1] OR OBSERVATION o1 [openEHR-EHR-OBSERVATION.laboratory-glucose.v1])`,
			From{
				ContainsExpr: ContainsExpr{
					Operand: ClassExpression{
						Identifiers: []string{"EHR", "e"},
					},
					Contains: []*ContainsExpr{
						{
							Operand: ClassExpression{
								Identifiers: []string{"COMPOSITION", "c"},
								PathPredicate: &PathPredicate{
									Type: ArchetypedPathPredicate,
									Archetype: &ArchetypePathPredicate{
										ArchetypeHRID: toRef("openEHR-EHR-COMPOSITION.referral.v1"),
									},
								},
							},
							Contains: []*ContainsExpr{
								{
									Brackets: true,
									Contains: []*ContainsExpr{
										{
											Operator: toRef(OROperator),
											Contains: []*ContainsExpr{
												{
													Operand: ClassExpression{
														Identifiers: []string{"OBSERVATION", "o"},
														PathPredicate: &PathPredicate{
															Type: ArchetypedPathPredicate,
															Archetype: &ArchetypePathPredicate{
																ArchetypeHRID: toRef("openEHR-EHR-OBSERVATION.laboratory-hba1c.v1"),
															},
														},
													},
												},
												{
													Operand: ClassExpression{
														Identifiers: []string{"OBSERVATION", "o1"},
														PathPredicate: &PathPredicate{
															Type: ArchetypedPathPredicate,
															Archetype: &ArchetypePathPredicate{
																ArchetypeHRID: toRef("openEHR-EHR-OBSERVATION.laboratory-glucose.v1"),
															},
														},
													},
												},
											},
										},
									},
								},
							},
						},
					},
				},
			},
			false,
		},
		{
			"8. FROM with NOT operator",
			`SELECT val
			FROM EHR e
				CONTAINS COMPOSITION c [openEHR-EHR-COMPOSITION.referral.v1]
					NOT CONTAINS OBSERVATION o [openEHR-EHR-OBSERVATION.laboratory_test_result.v1]`,
			From{
				ContainsExpr: ContainsExpr{
					Operand: ClassExpression{
						Identifiers: []string{"EHR", "e"},
					},
					Contains: []*ContainsExpr{
						{
							Operator: toRef(NOTOperator),
							Operand: ClassExpression{
								Identifiers: []string{"COMPOSITION", "c"},
								PathPredicate: &PathPredicate{
									Type: ArchetypedPathPredicate,
									Archetype: &ArchetypePathPredicate{
										ArchetypeHRID: toRef("openEHR-EHR-COMPOSITION.referral.v1"),
									},
								},
							},
							Contains: []*ContainsExpr{
								{
									Operand: ClassExpression{
										Identifiers: []string{"OBSERVATION", "o"},
										PathPredicate: &PathPredicate{
											Type: ArchetypedPathPredicate,
											Archetype: &ArchetypePathPredicate{
												ArchetypeHRID: toRef("openEHR-EHR-OBSERVATION.laboratory_test_result.v1"),
											},
										},
									},
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

			if got == nil {
				return
			}

			if diff := cmp.Diff(tt.want, got.From); diff != "" {
				t.Errorf("Mismatch {+want;-got}:\n\t%s", diff)
			}
		})
	}
}
