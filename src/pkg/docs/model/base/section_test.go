package base_test

import (
	"encoding/json"
	"testing"

	"github.com/bsn-si/IPEHR-gateway/src/pkg/docs/model/base"

	"github.com/google/go-cmp/cmp"
)

func TestSection_UnmarshalJSON(t *testing.T) {
	tests := []struct {
		name    string
		data    []byte
		want    base.Section
		wantErr bool
	}{
		{
			"1. empty json",
			nil,
			base.Section{},
			true,
		},
		{
			"2. selection object",
			[]byte(selectionJSON),
			base.Section{
				Locatable: base.Locatable{
					Type: "SECTION",
					Name: base.NewDvText("Medication Summary"),
					ArchetypeDetails: &base.Archetyped{
						Type: "ARCHETYPED",
						ArchetypeID: base.ObjectID{
							Type:  "ARCHETYPE_ID",
							Value: "openEHR-EHR-SECTION.adhoc.v1",
						},
						RmVersion: "1.0.4",
					},
					ArchetypeNodeID: "openEHR-EHR-SECTION.adhoc.v1",
				},
				Items: []base.Root{
					&base.Action{
						CareEntry: base.CareEntry{
							Entry: base.Entry{
								ContentItem: base.ContentItem{base.Locatable{
									Type: base.ActionItemType,
								}},
							},
						},
					},
					&base.Evaluation{
						CareEntry: base.CareEntry{
							Entry: base.Entry{
								ContentItem: base.ContentItem{base.Locatable{
									Type: base.EvaluationItemType,
								}},
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
			got := base.Section{}
			if err := json.Unmarshal(tt.data, &got); (err != nil) != tt.wantErr {
				t.Errorf("Section.UnmarshalJSON() error = %v, wantErr %v", err, tt.wantErr)
			}

			opts := cmp.AllowUnexported(
				base.ObjectVersionID{},
				base.PartyProxy{},
			)
			if diff := cmp.Diff(tt.want, got, opts); diff != "" {
				t.Errorf("Section.UnmarshalJSON() mismatch {-want;+got}\n\t%s", diff)
			}
		})
	}
}

const selectionJSON = `{
    "_type": "SECTION",
    "name": {
        "_type": "DV_TEXT",
        "value": "Medication Summary"
    },
    "archetype_details": {
        "_type": "ARCHETYPED",
        "archetype_id": {
            "_type": "ARCHETYPE_ID",
            "value": "openEHR-EHR-SECTION.adhoc.v1"
        },
        "rm_version": "1.0.4"
    },
    "archetype_node_id": "openEHR-EHR-SECTION.adhoc.v1",
    "items": [
        {
            "_type": "ACTION"
        },
        {
            "_type": "EVALUATION"
        }
    ]
}`
