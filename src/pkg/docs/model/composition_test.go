package model_test

import (
	"encoding/json"
	"os"
	"testing"

	"github.com/bsn-si/IPEHR-gateway/src/pkg/docs/model"
	"github.com/bsn-si/IPEHR-gateway/src/pkg/docs/model/base"

	"github.com/google/go-cmp/cmp"
	"github.com/stretchr/testify/assert"
)

func TestComposition_UnmarshalJSON(t *testing.T) {
	tests := []struct {
		name    string
		data    []byte
		wantErr bool
		want    model.Composition
	}{
		{
			"1. error on unmarshal data",
			[]byte(""),
			true,
			model.Composition{},
		},
		{
			"2. composition without content",
			[]byte(`{
				  "_type": "COMPOSITION",
				  "name": {
				    "_type": "DV_TEXT",
				    "value": "International Patient Summary"
				  },
				  "uid": {
				    "_type": "OBJECT_VERSION_ID",
				    "value": "41f6fdb5-9ea5-4bb8-b2fa-21131543f82e::openEHRSys.example.com::1"
				  },
				  "archetype_details": {
				    "_type": "ARCHETYPED",
				    "archetype_id": {
				      "_type": "ARCHETYPE_ID",
				      "value": "openEHR-EHR-COMPOSITION.health_summary.v1"
				    },
				    "template_id": {
				      "_type": "TEMPLATE_ID",
				      "value": "International Patient Summary"
				    },
				    "rm_version": "1.0.4"
				  }
			}`),
			false,
			model.Composition{
				Locatable: base.Locatable{
					Type: "COMPOSITION",
					Name: base.NewDvText("International Patient Summary"),
					ObjectVersionID: base.ObjectVersionID{
						UID: &base.UIDBasedID{
							ObjectID: base.ObjectID{
								Type:  "OBJECT_VERSION_ID",
								Value: "41f6fdb5-9ea5-4bb8-b2fa-21131543f82e::openEHRSys.example.com::1",
							},
						},
					},
					ArchetypeDetails: &base.Archetyped{
						Type: "ARCHETYPED",
						ArchetypeID: base.ObjectID{
							Type:  "ARCHETYPE_ID",
							Value: "openEHR-EHR-COMPOSITION.health_summary.v1",
						},
						TemplateID: &base.ObjectID{
							Type:  "TEMPLATE_ID",
							Value: "International Patient Summary",
						},
						RmVersion: "1.0.4",
					},
				},
			},
		},
		{
			"3. parse json",
			[]byte(compositionJSON),
			false,
			expectedComposition,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := model.Composition{}
			if err := json.Unmarshal(tt.data, &got); (err != nil) != tt.wantErr {
				t.Errorf("Composition.UnmarshalJSON() error = %v, wantErr %v", err, tt.wantErr)
			}

			opts := cmp.AllowUnexported(
				base.ObjectVersionID{},
				base.PartyProxy{},
			)
			if diff := cmp.Diff(tt.want, got, opts); diff != "" {
				t.Errorf("Composition.UnmarshalJSON() mismatch{-want;+got}\n\t%s", diff)
			}
		})
	}
}

func TestParseComposition(t *testing.T) {
	filePath := "./test_fixtures/composition.json"

	inJSON, err := os.ReadFile(filePath)
	if err != nil {
		t.Fatal("Can't open composition.json file", filePath)
	}

	res := model.Composition{}

	if err := json.Unmarshal(inJSON, &res); err != nil {
		t.Error(err)
		return
	}

	if res.UID.Value == "" {
		t.Error("Composition is not parsed correctly")
	}
}

func TestMarshalAndUnmarshalComposition(t *testing.T) {
	filePath := "./test_fixtures/composition.json"

	inJSON, err := os.ReadFile(filePath)
	if err != nil {
		t.Fatal("Can't open composition.json file", filePath)
	}

	composition := model.Composition{}

	err = json.Unmarshal(inJSON, &composition)
	assert.Nil(t, err)

	data, err := json.Marshal(composition)
	assert.Nil(t, err)

	newComposition := model.Composition{}

	err = json.Unmarshal(data, &newComposition)
	if !assert.NoError(t, err) {
		return
	}

	assert.Equal(t, composition, newComposition)
}

func toRef[T any](v T) *T {
	return &v
}

var expectedComposition = model.Composition{
	Language: base.CodePhrase{
		Type: base.CodePhraseItemType,
		TerminologyID: base.ObjectID{
			Type:  "TERMINOLOGY_ID",
			Value: "ISO_639-1",
		},
		CodeString: "en",
	},
	Territory: base.CodePhrase{
		Type: base.CodePhraseItemType,
		TerminologyID: base.ObjectID{
			Type:  "TERMINOLOGY_ID",
			Value: "ISO_3166-1",
		},
		CodeString: "US",
	},
	Composer: base.NewPartyProxy(
		&base.PartyIdentified{
			Name: "Silvia Blake",
			PartyProxyBase: base.PartyProxyBase{
				Type: base.PartyIdentifiedItemType,
			},
		},
	),
	Category: base.NewDvCodedText(
		"event",
		base.CodePhrase{
			Type: base.CodePhraseItemType,
			TerminologyID: base.ObjectID{
				Type:  "TERMINOLOGY_ID",
				Value: "openehr",
			},
			CodeString: "433",
		},
	),
	Context: &model.EventContext{
		StartTime: base.DvDateTime{
			DvTemporal: base.DvTemporal{
				DvValueBase: base.DvValueBase{
					Type: base.DvDateTimeItemType,
				},
			},
			Value: "2021-12-03T17:34:06.849379+01:00",
		},
		Setting: base.NewDvCodedText(
			"other care",
			base.CodePhrase{
				Type: base.CodePhraseItemType,
				TerminologyID: base.ObjectID{
					Type:  "TERMINOLOGY_ID",
					Value: "openehr",
				},
				CodeString: "238",
			},
		),
		HealthCareFacility: &base.PartyIdentified{
			Name: "Hospital",
			PartyProxyBase: base.PartyProxyBase{
				Type: base.PartyIdentifiedItemType,
				ExternalRef: &base.ObjectRef{
					ID: base.ObjectID{
						Type:  "GENERIC_ID",
						Value: "9091",
					},
					Namespace: "HOSPITAL-NS",
					Type:      "PARTY_REF",
				},
			},
		},
		Participations: []base.Participation{
			{
				Function: base.NewDvText("requester"),
				Mode: toRef(base.NewDvCodedText(
					"face-to-face communication",
					base.CodePhrase{
						Type: base.CodePhraseItemType,
						TerminologyID: base.ObjectID{
							Type:  "TERMINOLOGY_ID",
							Value: "openehr",
						},
						CodeString: "216",
					},
				)),
				Performer: base.NewPartyProxy(
					&base.PartyIdentified{
						Name: "Dr. Marcus Johnson",
						PartyProxyBase: base.PartyProxyBase{
							Type: base.PartyIdentifiedItemType,
							ExternalRef: &base.ObjectRef{
								ID: base.ObjectID{
									Type:  "GENERIC_ID",
									Value: "199",
								},
								Namespace: "HOSPITAL-NS",
								Type:      "PARTY_REF",
							},
						},
					},
				),
			},
			{
				Function: base.NewDvText("performer"),
				Mode: toRef(base.NewDvCodedText(
					"not specified",
					base.CodePhrase{
						Type: base.CodePhraseItemType,
						TerminologyID: base.ObjectID{
							Type:  "TERMINOLOGY_ID",
							Value: "openehr",
						},
						CodeString: "193",
					},
				)),
				Performer: base.NewPartyProxy(
					&base.PartyIdentified{
						Name: "Lara Markham",
						PartyProxyBase: base.PartyProxyBase{
							Type: base.PartyIdentifiedItemType,
							ExternalRef: &base.ObjectRef{
								ID: base.ObjectID{
									Type:  "GENERIC_ID",
									Value: "198",
								},
								Namespace: "HOSPITAL-NS",
								Type:      "PARTY_REF",
							},
						},
					},
				),
			},
		},
	},
	Locatable: base.Locatable{
		Type:            "COMPOSITION",
		Name:            base.NewDvText("International Patient Summary"),
		ArchetypeNodeID: "openEHR-EHR-COMPOSITION.health_summary.v1",
		ObjectVersionID: base.ObjectVersionID{
			UID: &base.UIDBasedID{
				ObjectID: base.ObjectID{
					Type:  "OBJECT_VERSION_ID",
					Value: "41f6fdb5-9ea5-4bb8-b2fa-21131543f82e::openEHRSys.example.com::1",
				},
			},
		},
		ArchetypeDetails: &base.Archetyped{
			Type: "ARCHETYPED",
			ArchetypeID: base.ObjectID{
				Type: "ARCHETYPE_ID", Value: "openEHR-EHR-COMPOSITION.health_summary.v1",
			},
			TemplateID: &base.ObjectID{
				Type: "TEMPLATE_ID", Value: "International Patient Summary",
			},
			RmVersion: "1.0.4",
		},
	},
	Content: []base.Root{
		&base.Section{
			Locatable: base.Locatable{
				Type:            "SECTION",
				Name:            base.NewDvText("Medication Summary"),
				ArchetypeNodeID: "openEHR-EHR-SECTION.adhoc.v1",
				ArchetypeDetails: &base.Archetyped{
					Type: "ARCHETYPED",
					ArchetypeID: base.ObjectID{
						Type: "ARCHETYPE_ID", Value: "openEHR-EHR-SECTION.adhoc.v1",
					},
					RmVersion: "1.0.4",
				},
			},
			Items: []base.Root{},
		},
	},
}

const compositionJSON = `{
  "_type": "COMPOSITION",
  "name": {
    "_type": "DV_TEXT",
    "value": "International Patient Summary"
  },
  "uid": {
    "_type": "OBJECT_VERSION_ID",
    "value": "41f6fdb5-9ea5-4bb8-b2fa-21131543f82e::openEHRSys.example.com::1"
  },
  "archetype_details": {
    "_type": "ARCHETYPED",
    "archetype_id": {
      "_type": "ARCHETYPE_ID",
      "value": "openEHR-EHR-COMPOSITION.health_summary.v1"
    },
    "template_id": {
      "_type": "TEMPLATE_ID",
      "value": "International Patient Summary"
    },
    "rm_version": "1.0.4"
  },
  "archetype_node_id": "openEHR-EHR-COMPOSITION.health_summary.v1",
  "language": {
    "_type": "CODE_PHRASE",
    "terminology_id": {
      "_type": "TERMINOLOGY_ID",
      "value": "ISO_639-1"
    },
    "code_string": "en"
  },
  "territory": {
    "_type": "CODE_PHRASE",
    "terminology_id": {
      "_type": "TERMINOLOGY_ID",
      "value": "ISO_3166-1"
    },
    "code_string": "US"
  },
  "category": {
    "_type": "DV_CODED_TEXT",
    "value": "event",
    "defining_code": {
      "_type": "CODE_PHRASE",
      "terminology_id": {
        "_type": "TERMINOLOGY_ID",
        "value": "openehr"
      },
      "code_string": "433"
    }
  },
  "composer": {
    "_type": "PARTY_IDENTIFIED",
    "name": "Silvia Blake"
  },
  "context": {
    "_type": "EVENT_CONTEXT",
    "start_time": {
      "_type": "DV_DATE_TIME",
      "value": "2021-12-03T17:34:06.849379+01:00"
    },
    "setting": {
      "_type": "DV_CODED_TEXT",
      "value": "other care",
      "defining_code": {
        "_type": "CODE_PHRASE",
        "terminology_id": {
          "_type": "TERMINOLOGY_ID",
          "value": "openehr"
        },
        "code_string": "238"
      }
    },
    "health_care_facility": {
      "_type": "PARTY_IDENTIFIED",
      "external_ref": {
        "type": "PARTY_REF",
        "id": {
          "_type": "GENERIC_ID",
          "value": "9091",
          "scheme": "HOSPITAL-NS"
        },
        "namespace": "HOSPITAL-NS"
      },
      "name": "Hospital"
    },
    "participations": [
      {
        "_type": "PARTICIPATION",
        "function": {
          "_type": "DV_TEXT",
          "value": "requester"
        },
        "performer": {
          "_type": "PARTY_IDENTIFIED",
          "external_ref": {
            "type": "PARTY_REF",
            "id": {
              "_type": "GENERIC_ID",
              "value": "199",
              "scheme": "HOSPITAL-NS"
            },
            "namespace": "HOSPITAL-NS"
          },
          "name": "Dr. Marcus Johnson"
        },
        "mode": {
          "_type": "DV_CODED_TEXT",
          "value": "face-to-face communication",
          "defining_code": {
            "_type": "CODE_PHRASE",
            "terminology_id": {
              "_type": "TERMINOLOGY_ID",
              "value": "openehr"
            },
            "code_string": "216"
          }
        }
      },
      {
        "_type": "PARTICIPATION",
        "function": {
          "_type": "DV_TEXT",
          "value": "performer"
        },
        "performer": {
          "_type": "PARTY_IDENTIFIED",
          "external_ref": {
            "type": "PARTY_REF",
            "id": {
              "_type": "GENERIC_ID",
              "value": "198",
              "scheme": "HOSPITAL-NS"
            },
            "namespace": "HOSPITAL-NS"
          },
          "name": "Lara Markham"
        },
        "mode": {
          "_type": "DV_CODED_TEXT",
          "value": "not specified",
          "defining_code": {
            "_type": "CODE_PHRASE",
            "terminology_id": {
              "_type": "TERMINOLOGY_ID",
              "value": "openehr"
            },
            "code_string": "193"
          }
        }
      }
    ]
  },
  "content": [
	{
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
      "items": []
    }
  ]
}`
