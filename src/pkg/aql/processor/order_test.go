package processor

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestParser_Order(t *testing.T) {
	tests := []struct {
		name    string
		query   string
		want    *Order
		wantErr bool
	}{
		{
			"1. empty order block",
			`SELECT val FROM ehr`,
			nil,
			false,
		},
		{
			"2. simple order",
			`SELECT val FROM ehr ORDER BY c/name/value`,
			&Order{
				Orders: []OrderBy{
					{
						IdentifierPath: IdentifiedPath{
							Identifier: "c",
							ObjectPath: &ObjectPath{
								Paths: []PartPath{
									{Identifier: "name"},
									{Identifier: "value"},
								},
							},
						},
						Ordering: NoneOrdering,
					},
				},
			},
			false,
		},
		{
			"3. ASCENDING ordering",
			`SELECT val FROM ehr ORDER BY c/name/value ASCENDING`,
			&Order{
				Orders: []OrderBy{
					{
						IdentifierPath: IdentifiedPath{
							Identifier: "c",
							ObjectPath: &ObjectPath{
								Paths: []PartPath{
									{Identifier: "name"},
									{Identifier: "value"},
								},
							},
						},
						Ordering: AscendingOrdering,
					},
				},
			},
			false,
		},
		{
			"4. ASC ordering",
			`SELECT val FROM ehr ORDER BY c/name/value ASC`,
			&Order{
				Orders: []OrderBy{
					{
						IdentifierPath: IdentifiedPath{
							Identifier: "c",
							ObjectPath: &ObjectPath{
								Paths: []PartPath{
									{Identifier: "name"},
									{Identifier: "value"},
								},
							},
						},
						Ordering: AscendingOrdering,
					},
				},
			},
			false,
		},

		{
			"5. ASCENDING ordering",
			`SELECT val FROM ehr ORDER BY c/name/value DESCENDING`,
			&Order{
				Orders: []OrderBy{
					{
						IdentifierPath: IdentifiedPath{
							Identifier: "c",
							ObjectPath: &ObjectPath{
								Paths: []PartPath{
									{Identifier: "name"},
									{Identifier: "value"},
								},
							},
						},
						Ordering: DescendingOrdering,
					},
				},
			},
			false,
		},
		{
			"6. DESC ordering",
			`SELECT val FROM ehr ORDER BY c/name/value DESC`,
			&Order{
				Orders: []OrderBy{
					{
						IdentifierPath: IdentifiedPath{
							Identifier: "c",
							ObjectPath: &ObjectPath{
								Paths: []PartPath{
									{Identifier: "name"},
									{Identifier: "value"},
								},
							},
						},
						Ordering: DescendingOrdering,
					},
				},
			},
			false,
		},
		{
			"7. Multi column ordering",
			`SELECT val FROM ehr ORDER BY c/name/value DESC, c/name/value2`,
			&Order{
				Orders: []OrderBy{
					{
						IdentifierPath: IdentifiedPath{
							Identifier: "c",
							ObjectPath: &ObjectPath{
								Paths: []PartPath{
									{Identifier: "name"},
									{Identifier: "value"},
								},
							},
						},
						Ordering: DescendingOrdering,
					},
					{
						IdentifierPath: IdentifiedPath{
							Identifier: "c",
							ObjectPath: &ObjectPath{
								Paths: []PartPath{
									{Identifier: "name"},
									{Identifier: "value2"},
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

			if diff := cmp.Diff(tt.want, got.Order); diff != "" {
				t.Errorf("mismatch {+want;-got}:\n\t%s", diff)
			}
		})
	}
}
