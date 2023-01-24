package base_test

import (
	"encoding/json"
	"testing"

	"github.com/bsn-si/IPEHR-gateway/src/pkg/docs/model/base"

	"github.com/google/go-cmp/cmp"
)

func TestAction_UnmarshalJSON(t *testing.T) {
	tests := []struct {
		name    string
		data    []byte
		want    base.Action
		wantErr bool
	}{
		{
			"1. invalid json",
			[]byte("invalid_json"),
			base.Action{},
			true,
		},
		{
			"2. valid json",
			[]byte(actionJSON),
			base.Action{
				Time: base.DvDateTime{
					DvTemporal: base.DvTemporal{
						DvValueBase: base.DvValueBase{Type: base.DvDateTimeItemType},
					},
					Value: "2021-12-03T16:05:19.513939+01:00",
				},
				IsmTransition: base.IsmTransition{
					CurrentState: base.DvCodedText{
						DefiningCode: base.CodePhrase{
							Type:          base.CodePhraseItemType,
							TerminologyID: base.ObjectID{Type: base.TerminologyIDItemType, Value: "openehr"},
							CodeString:    "245",
						},
						DvText: base.DvText{DvValueBase: base.DvValueBase{Type: base.DvCodedTextItemType}, Value: "active"},
					},
				},
				Description: base.ItemTree{
					DataStructure: base.DataStructure{
						Locatable: base.Locatable{
							Type:            base.ItemTreeItemType,
							Name:            base.NewDvText("Tree"),
							ArchetypeNodeID: "at0017",
						},
					},
					Items: base.Items{},
				},
				CareEntry: base.CareEntry{
					Protocol: &base.ItemStructure{
						Data: &base.ItemTree{
							DataStructure: base.DataStructure{
								Locatable: base.Locatable{
									Type:            base.ItemTreeItemType,
									Name:            base.NewDvText("Tree"),
									ArchetypeNodeID: "at0030",
								},
							},
							Items: base.Items{},
						},
					},
					Entry: base.Entry{
						ContentItem: base.ContentItem{
							Locatable: base.Locatable{
								Type:            base.ActionItemType,
								Name:            base.NewDvText("Medication statement"),
								ArchetypeNodeID: "openEHR-EHR-ACTION.medication.v1",
								ArchetypeDetails: &base.Archetyped{
									Type: base.ArchetypedItemType,
									ArchetypeID: base.ObjectID{
										Type:  base.ArchetypeIDItemType,
										Value: "openEHR-EHR-ACTION.medication.v1",
									},
									RmVersion: "1.0.4",
								},
							},
						},
						Language: base.CodePhrase{
							Type:          base.CodePhraseItemType,
							TerminologyID: base.ObjectID{Type: base.TerminologyIDItemType, Value: "ISO_639-1"},
							CodeString:    "en",
						},
						Encoding: base.CodePhrase{
							Type:          base.CodePhraseItemType,
							TerminologyID: base.ObjectID{Type: base.TerminologyIDItemType, Value: "IANA_character-sets"},
							CodeString:    "UTF-8",
						},
						OtherParticipations: []base.Participation{},
						Subject: base.NewPartyProxy(
							&base.PartySelf{
								base.PartyProxyBase{
									Type: base.PartySelfItemType,
								},
							},
						),
					},
				},
			},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := base.Action{}
			if err := json.Unmarshal(tt.data, &got); (err != nil) != tt.wantErr {
				t.Errorf("Action.UnmarshalJSON() error = %v, wantErr %v", err, tt.wantErr)
			}

			opts := cmp.AllowUnexported(
				base.ObjectVersionID{},
				base.PartyProxy{},
			)
			if diff := cmp.Diff(tt.want, got, opts); diff != "" {
				t.Errorf("Action.UnmarshalJSON() mismatch {-want;+got}\n\t%s", diff)
			}
		})
	}
}

const actionJSON = `{
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
        "items": []
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
}`
