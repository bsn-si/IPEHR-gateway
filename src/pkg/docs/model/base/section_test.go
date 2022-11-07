package base

import (
	"encoding/json"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestSection_UnmarshalJSON(t *testing.T) {
	tests := []struct {
		name    string
		data    []byte
		want    Section
		wantErr bool
	}{
		{
			"1. empty json",
			nil,
			Section{},
			true,
		},
		{
			"2. selection object",
			[]byte(selectionJSON),
			Section{
				Locatable: Locatable{
					Type: "SECTION",
					Name: NewDvText("Medication Summary"),
					ArchetypeDetails: &Archetyped{
						Type: "ARCHETYPED",
						ArchetypeID: ObjectID{
							Type:  "ARCHETYPE_ID",
							Value: "openEHR-EHR-SECTION.adhoc.v1",
						},
						RmVersion: "1.0.4",
					},
					ArchetypeNodeID: "openEHR-EHR-SECTION.adhoc.v1",
				},
				Items: []ContentItem{
					&Action{
						Time: DvDateTime{Value: "2021-12-03T16:05:19.513939+01:00"},
						IsmTransition: IsmTransition{
							CurrentState: DvCodedText{
								DefiningCode: CodePhrase{
									TerminologyID: ObjectID{Type: TerminologyIDContentItemType, Value: "openehr"},
									CodeString:    "245",
								},
								DvText: DvText{Type: DvCodedTextContentItemType, Value: "active"},
							},
						},
						Description: ItemStructure{DataStructure{Locatable{
							Type:            ItemTreeContentItemType,
							Name:            NewDvText("Tree"),
							ArchetypeNodeID: "at0017",
						}}},
						CareEntry: CareEntry{
							Protocol: ItemStructure{DataStructure{Locatable{
								Type:            ItemTreeContentItemType,
								Name:            NewDvText("Tree"),
								ArchetypeNodeID: "at0030",
							}}},
							Entry: Entry{
								Language: CodePhrase{
									TerminologyID: ObjectID{Type: TerminologyIDContentItemType, Value: "ISO_639-1"},
									CodeString:    "en",
								},
								Encoding: CodePhrase{
									TerminologyID: ObjectID{Type: TerminologyIDContentItemType, Value: "IANA_character-sets"},
									CodeString:    "UTF-8",
								},
								OtherParticipations: &[]Participation{},
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
			got := Section{}
			if err := json.Unmarshal(tt.data, &got); (err != nil) != tt.wantErr {
				t.Errorf("Section.UnmarshalJSON() error = %v, wantErr %v", err, tt.wantErr)
			}

			if diff := cmp.Diff(tt.want, got, cmp.AllowUnexported(ObjectVersionID{})); diff != "" {
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
            "_type": "ACTION",
            "name": {
                "_type": "DV_TEXT",
                "value": "Medication statement"
            },
            "archetype_details": {
                "_type": "ARCHETYPED",
                "archetype_id": {
                    "_type": "ARCHETYPE_ID",
                    "value": "openEHR-EHR-ACTION.medication.v1"
                },
                "rm_version": "1.0.4"
            },
            "archetype_node_id": "openEHR-EHR-ACTION.medication.v1",
            "language": {
                "_type": "CODE_PHRASE",
                "terminology_id": {
                    "_type": "TERMINOLOGY_ID",
                    "value": "ISO_639-1"
                },
                "code_string": "en"
            },
            "encoding": {
                "_type": "CODE_PHRASE",
                "terminology_id": {
                    "_type": "TERMINOLOGY_ID",
                    "value": "IANA_character-sets"
                },
                "code_string": "UTF-8"
            },
            "subject": {
                "_type": "PARTY_SELF"
            },
            "other_participations": [],
            "protocol": {
                "_type": "ITEM_TREE",
                "name": {
                    "_type": "DV_TEXT",
                    "value": "Tree"
                },
                "archetype_node_id": "at0030",
                "items": [
                    {
                        "_type": "ELEMENT",
                        "name": {
                            "_type": "DV_TEXT",
                            "value": "Order ID"
                        },
                        "archetype_node_id": "at0103",
                        "value": {
                            "_type": "DV_IDENTIFIER",
                            "issuer": "Issuer",
                            "assigner": "Assigner",
                            "id": "9a0e5173-07c8-443d-b414-24432b9d95ca",
                            "type": "Prescription"
                        }
                    }
                ]
            },
            "time": {
                "_type": "DV_DATE_TIME",
                "value": "2021-12-03T16:05:19.513939+01:00"
            },
            "description": {
                "_type": "ITEM_TREE",
                "name": {
                    "_type": "DV_TEXT",
                    "value": "Tree"
                },
                "archetype_node_id": "at0017",
                "items": []
            },
            "ism_transition": {
                "_type": "ISM_TRANSITION",
                "current_state": {
                    "_type": "DV_CODED_TEXT",
                    "value": "active",
                    "defining_code": {
                        "_type": "CODE_PHRASE",
                        "terminology_id": {
                            "_type": "TERMINOLOGY_ID",
                            "value": "openehr"
                        },
                        "code_string": "245"
                    }
                }
            }
        }
    ]
}`
