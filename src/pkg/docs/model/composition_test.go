package model_test

import (
	"encoding/json"
	"hms/gateway/pkg/docs/model"
	"hms/gateway/pkg/docs/model/base"
	"testing"

	"github.com/google/go-cmp/cmp"
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
					Name: base.DvText{
						Type:  "DV_TEXT",
						Value: "International Patient Summary",
					},
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
			"3. read data from file",
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

			if diff := cmp.Diff(tt.want, got, cmp.AllowUnexported(base.ObjectVersionID{})); diff != "" {
				t.Errorf("Composition.UnmarshalJSON() mismatch{-want;+got}\n\t%s", diff)
			}
		})
	}
}

var expectedComposition = model.Composition{
	Language: base.CodePhrase{
		TerminologyID: base.ObjectID{
			Type:  "TERMINOLOGY_ID",
			Value: "ISO_639-1",
		},
		CodeString: "en",
	},
	Territory: base.CodePhrase{
		TerminologyID: base.ObjectID{
			Type:  "TERMINOLOGY_ID",
			Value: "ISO_3166-1",
		},
		CodeString: "US",
	},
	Category: base.DvCodedText{
		DefiningCode: base.CodePhrase{
			TerminologyID: base.ObjectID{
				Type:  "TERMINOLOGY_ID",
				Value: "openehr",
			},
			CodeString: "433",
		},
		DvText: base.DvText{
			Type:  "DV_CODED_TEXT",
			Value: "event",
		},
	},
	Context: &model.EventContext{
		StartTime: base.DvDateTime{
			Value: "2021-12-03T17:34:06.849379+01:00",
		},
		Setting: base.DvCodedText{
			DefiningCode: base.CodePhrase{
				TerminologyID: base.ObjectID{
					Type:  "TERMINOLOGY_ID",
					Value: "openehr",
				},
				CodeString: "238",
			},
			DvText: base.DvText{
				Type:  "DV_CODED_TEXT",
				Value: "other care",
			},
		},
		HealthCareFacility: &base.PartyIdentified{
			Name: "Hospital",
			PartyProxy: base.PartyProxy{
				ExternalRef: base.ObjectRef{
					ID: base.ObjectID{
						Type:  "GENERIC_ID",
						Value: "9091",
					},
					Namespace: "HOSPITAL-NS",
					Type:      "PARTY",
				},
			},
		},
		Participations: &[]base.Participation{
			{
				Function: base.DvText{
					Type:  "DV_TEXT",
					Value: "requester",
				},
				Mode: &base.DvCodedText{
					DefiningCode: base.CodePhrase{
						TerminologyID: base.ObjectID{
							Type:  "TERMINOLOGY_ID",
							Value: "openehr",
						},
						CodeString: "216",
					},
					DvText: base.DvText{
						Type:  "DV_CODED_TEXT",
						Value: "face-to-face communication",
					},
				},
				Performer: base.PartyProxy{
					ExternalRef: base.ObjectRef{
						ID: base.ObjectID{
							Type:  "GENERIC_ID",
							Value: "199",
						},
						Namespace: "HOSPITAL-NS",
						Type:      "PERSON",
					},
				},
			},
			{
				Function: base.DvText{
					Type:  "DV_TEXT",
					Value: "performer",
				},
				Mode: &base.DvCodedText{
					DefiningCode: base.CodePhrase{
						TerminologyID: base.ObjectID{
							Type:  "TERMINOLOGY_ID",
							Value: "openehr",
						},
						CodeString: "193",
					},
					DvText: base.DvText{
						Type:  "DV_CODED_TEXT",
						Value: "not specified",
					},
				},
				Performer: base.PartyProxy{
					ExternalRef: base.ObjectRef{
						ID: base.ObjectID{
							Type:  "GENERIC_ID",
							Value: "198",
						},
						Namespace: "HOSPITAL-NS",
						Type:      "PERSON",
					},
				},
			},
		},
	},
	Locatable: base.Locatable{
		Type: "COMPOSITION",
		Name: base.DvText{
			Type:  "DV_TEXT",
			Value: "International Patient Summary",
		},
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
	Content: []base.Section{
		{
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
			Items: []base.ContentItem{},
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
        "_type": "PARTY_REF",
        "id": {
          "_type": "GENERIC_ID",
          "value": "9091",
          "scheme": "HOSPITAL-NS"
        },
        "namespace": "HOSPITAL-NS",
        "type": "PARTY"
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
            "_type": "PARTY_REF",
            "id": {
              "_type": "GENERIC_ID",
              "value": "199",
              "scheme": "HOSPITAL-NS"
            },
            "namespace": "HOSPITAL-NS",
            "type": "PERSON"
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
            "_type": "PARTY_REF",
            "id": {
              "_type": "GENERIC_ID",
              "value": "198",
              "scheme": "HOSPITAL-NS"
            },
            "namespace": "HOSPITAL-NS",
            "type": "PERSON"
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
