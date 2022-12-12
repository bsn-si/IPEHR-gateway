package model_test

import (
	"encoding/json"
	"hms/gateway/pkg/docs/model"
	"hms/gateway/pkg/docs/model/base"
	//"os"
	"testing"

	"github.com/google/go-cmp/cmp"
	//"github.com/stretchr/testify/assert"
)

func TestContribution_UnmarshalJSON(t *testing.T) {
	tests := []struct {
		name    string
		data    []byte
		wantErr bool
		want    model.Contribution
	}{
		//{
		//	"1. error on unmarshal data",
		//	[]byte(""),
		//	true,
		//	model.ContributionVersion{},
		//},
		{
			"2. contribution with composition",
			[]byte(`{
		"_type": "ORIGINAL_VERSION",
		"contribution": {
			"id": {
				"_type": "HIER_OBJECT_ID",
				"value": "720ed9fa-4bba-4817-9173-4c05b00acf6a"
			},
			"namespace": "EHR::COMMON",
			"type": "CONTRIBUTION"
		},
		"commit_audit": {
			"_type": "AUDIT_DETAILS",
			"system_id": "CABOLABS_EHRSERVER",
			"time_committed": {
				"value": "2021-09-21T21:52:31.869-03:00"
			},
			"change_type": {
				"value": "creation",
				"defining_code": {
					"terminology_id": {
						"value": "openehr"
					},
					"code_string": "249"
				}
			},
			"committer": {
				"_type": "PARTY_IDENTIFIED",
				"external_ref": {
					"id": {
						"_type": "HIER_OBJECT_ID",
						"value": "f7e48c23-21b2-4b58-b9e0-a3ccece1bcf1"
					},
					"namespace": "DEMOGRAPHIC",
					"type": "PERSON"
				},
				"name": "Dr. Yamamoto"
			}
		},
		"uid": {
			"_type": "OBJECT_VERSION_ID",
			"value": "d11739a8-545d-4137-9bcd-9e5617252a0b::EMR_APP::1"
		},
		"data": {
			"_type": "COMPOSITION",
			"name": {
				"_type": "DV_TEXT",
				"value": "Minimal"
			},
			"archetype_details": {
				"archetype_id": {
					"value": "openEHR-EHR-COMPOSITION.minimal.v1"
				},
				"template_id": {
					"value": "minimal_evaluation.en.v1"
				},
				"rm_version": "1.0.2"
			},
			"archetype_node_id": "openEHR-EHR-COMPOSITION.minimal.v1",
			"language": {
				"terminology_id": {
					"value": "ISO_639-1"
				},
				"code_string": "en"
			},
			"territory": {
				"terminology_id": {
					"value": "ISO_3166-1"
				},
				"code_string": "UY"
			},
			"category": {
				"value": "event",
				"defining_code": {
					"terminology_id": {
						"value": "openehr"
					},
					"code_string": "433"
				}
			},
			"composer": {
				"_type": "PARTY_IDENTIFIED",
				"external_ref": {
					"id": {
						"_type": "HIER_OBJECT_ID",
						"value": "fc376c46-29c1-4090-bc01-5cc046af7f26"
					},
					"namespace": "DEMOGRAPHIC",
					"type": "PERSON"
				},
				"name": "Dr. House"
			},
			"context": {
				"start_time": {
					"value": "2021-09-21T21:52:31.927-03:00"
				},
				"setting": {
					"value": "primary medical care",
					"defining_code": {
						"terminology_id": {
							"value": "openehr"
						},
						"code_string": "228"
					}
				},
				"participations": [{
					"function": {
						"value": "legal guardian consent author"
					},
					"performer": {
						"_type": "PARTY_RELATED",
						"name": "Alexandra Alamo",
						"relationship": {
							"value": "mother",
							"defining_code": {
								"terminology_id": {
									"value": "openehr"
								},
								"code_string": "10"
							}
						}
					},
					"mode": {
						"value": "not specified",
						"defining_code": {
							"terminology_id": {
								"value": "openehr"
							},
							"code_string": "193"
						}
					}
				}]
			},
			"content": [{
				"_type": "EVALUATION",
				"name": {
					"_type": "DV_TEXT",
					"value": "Minimal"
				},
				"archetype_details": {
					"archetype_id": {
						"value": "openEHR-EHR-EVALUATION.minimal.v1"
					},
					"template_id": {
						"value": "minimal_evaluation.en.v1"
					},
					"rm_version": "1.0.2"
				},
				"archetype_node_id": "openEHR-EHR-EVALUATION.minimal.v1",
				"language": {
					"terminology_id": {
						"value": "ISO_639-1"
					},
					"code_string": "en"
				},
				"encoding": {
					"terminology_id": {
						"value": "IANA_character-sets"
					},
					"code_string": "UTF-8"
				},
				"subject": {
					"_type": "PARTY_SELF"
				},
				"data": {
					"_type": "ITEM_TREE",
					"name": {
						"_type": "DV_TEXT",
						"value": "Arbol"
					},
					"archetype_node_id": "at0001",
					"items": [{
						"_type": "ELEMENT",
						"name": {
							"_type": "DV_TEXT",
							"value": "quantity"
						},
						"archetype_node_id": "at0002",
						"value": {
							"_type": "DV_QUANTITY",
							"magnitude": 974.0,
							"units": "kg"
						}
					}]
				}
			}]
		},
		"lifecycle_state": {
			"value": "complete",
			"defining_code": {
				"terminology_id": {
					"value": "openehr"
				},
				"code_string": "532"
			}
		}
	}`),
			false,
			model.Contribution{
				Audit: model.AuditDetails{
					Type:     base.AuditDetailsType,
					SystemID: "test-system-id",
					Committer: base.NewPartyProxy(
						&base.PartyIdentified{
							Name: "<optional name of the committer>",
							PartyProxyBase: base.PartyProxyBase{
								Type: base.PartyIdentifiedItemType,
								ExternalRef: &base.ObjectRef{
									ID: base.ObjectID{
										Type:  "GENERIC_ID",
										Value: "<OBJECT_ID>",
										//"scheme": "<ID SCHEME NAME>" TODO mmm no idea, in docs i didnt find information
									},
									Namespace: "DEMOGRAPHIC",
									Type:      "PERSON",
								},
							},
						},
					),
					ChangeType: base.DvCodedText{
						DefiningCode: base.CodePhrase{
							TerminologyID: base.ObjectID{Value: "openehr"},
							CodeString:    "249",
						},
						DvText: base.DvText{Value: "creation"},
					},

					Description: base.DvText{Value: "<optional audit description>"},
				},

				Versions: []model.ContributionVersion{
					{
						Type: "ORIGINAL_VERSION",
						Contribution: base.ObjectRef{
							ID: base.ObjectID{
								Type:  "HIER_OBJECT_ID",
								Value: "720ed9fa-4bba-4817-9173-4c05b00acf6a",
							},
							Namespace: "EHR::COMMON",
							Type:      "CONTRIBUTION",
						},
						CommitAudit: model.AuditDetails{
							Type:     base.AuditDetailsType,
							SystemID: "CABOLABS_EHRSERVER",
							TimeCommited: base.DvDateTime{
								Value: "2021-12-03T16:05:19.513939+01:00",
							},
							ChangeType: base.DvCodedText{
								DefiningCode: base.CodePhrase{
									TerminologyID: base.ObjectID{Value: "openehr"},
									CodeString:    "249",
								},
								DvText: base.DvText{Value: "creation"},
							},
							Committer: base.NewPartyProxy(
								&base.PartyIdentified{
									Name: "Dr. Yamamoto",
									PartyProxyBase: base.PartyProxyBase{
										Type: base.PartyIdentifiedItemType,
										ExternalRef: &base.ObjectRef{
											ID: base.ObjectID{
												Type:  "HIER_OBJECT_ID",
												Value: "f7e48c23-21b2-4b58-b9e0-a3ccece1bcf1",
											},
											Namespace: "DEMOGRAPHIC",
											Type:      "PERSON",
										},
									},
								},
							),
							Description: base.DvText{Value: "<optional audit description>"},
						},
						UID: base.UIDBasedID{
							ObjectID: base.ObjectID{
								Type:  "OBJECT_VERSION_ID",
								Value: "41f6fdb5-9ea5-4bb8-b2fa-21131543f82e::openEHRSys.example.com::1",
							},
						},

						LifecycleState: base.NewDvCodedText(
							"complete",
							base.CodePhrase{
								Type: base.TerminologyIDItemType,
								TerminologyID: base.ObjectID{
									Type:  "TERMINOLOGY_ID",
									Value: "openehr",
								},
								CodeString: "532",
							},
						),

						Data: expectedComposition,
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := model.Contribution{}
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

//	func TestParseComposition(t *testing.T) {
//		wd, _ := os.Getwd()
//		filePath := wd + "/../../../../data/mock/ehr/composition.json"
//
//		inJSON, err := os.ReadFile(filePath)
//		if err != nil {
//			t.Fatal("Can't open composition.json file", filePath)
//		}
//
//		res := model.Composition{}
//
//		if err := json.Unmarshal(inJSON, &res); err != nil {
//			t.Error(err)
//			return
//		}
//
//		if res.UID.Value == "" {
//			t.Error("Composition is not parsed correctly")
//		}
//	}
//
//	func TestMarshalAndUnmarshalComposition(t *testing.T) {
//		wd, _ := os.Getwd()
//		filePath := wd + "/../../../../data/mock/ehr/composition.json"
//
//		inJSON, err := os.ReadFile(filePath)
//		if err != nil {
//			t.Fatal("Can't open composition.json file", filePath)
//		}
//
//		composition := model.Composition{}
//
//		err = json.Unmarshal(inJSON, &composition)
//		assert.Nil(t, err)
//
//		data, err := json.Marshal(composition)
//		assert.Nil(t, err)
//
//		newComposition := model.Composition{}
//
//		err = json.Unmarshal(data, &newComposition)
//		if !assert.NoError(t, err) {
//			return
//		}
//
//		assert.Equal(t, composition, newComposition)
//	}
//
//	func toRef[T any](v T) *T {
//		return &v
//	}
var expectedComposition2 = model.Composition{
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
		Participations: &[]base.Participation{
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

//
//const compositionJSON = `{
//  "_type": "COMPOSITION",
//  "name": {
//    "_type": "DV_TEXT",
//    "value": "International Patient Summary"
//  },
//  "uid": {
//    "_type": "OBJECT_VERSION_ID",
//    "value": "41f6fdb5-9ea5-4bb8-b2fa-21131543f82e::openEHRSys.example.com::1"
//  },
//  "archetype_details": {
//    "_type": "ARCHETYPED",
//    "archetype_id": {
//      "_type": "ARCHETYPE_ID",
//      "value": "openEHR-EHR-COMPOSITION.health_summary.v1"
//    },
//    "template_id": {
//      "_type": "TEMPLATE_ID",
//      "value": "International Patient Summary"
//    },
//    "rm_version": "1.0.4"
//  },
//  "archetype_node_id": "openEHR-EHR-COMPOSITION.health_summary.v1",
//  "language": {
//    "_type": "CODE_PHRASE",
//    "terminology_id": {
//      "_type": "TERMINOLOGY_ID",
//      "value": "ISO_639-1"
//    },
//    "code_string": "en"
//  },
//  "territory": {
//    "_type": "CODE_PHRASE",
//    "terminology_id": {
//      "_type": "TERMINOLOGY_ID",
//      "value": "ISO_3166-1"
//    },
//    "code_string": "US"
//  },
//  "category": {
//    "_type": "DV_CODED_TEXT",
//    "value": "event",
//    "defining_code": {
//      "_type": "CODE_PHRASE",
//      "terminology_id": {
//        "_type": "TERMINOLOGY_ID",
//        "value": "openehr"
//      },
//      "code_string": "433"
//    }
//  },
//  "composer": {
//    "_type": "PARTY_IDENTIFIED",
//    "name": "Silvia Blake"
//  },
//  "context": {
//    "_type": "EVENT_CONTEXT",
//    "start_time": {
//      "_type": "DV_DATE_TIME",
//      "value": "2021-12-03T17:34:06.849379+01:00"
//    },
//    "setting": {
//      "_type": "DV_CODED_TEXT",
//      "value": "other care",
//      "defining_code": {
//        "_type": "CODE_PHRASE",
//        "terminology_id": {
//          "_type": "TERMINOLOGY_ID",
//          "value": "openehr"
//        },
//        "code_string": "238"
//      }
//    },
//    "health_care_facility": {
//      "_type": "PARTY_IDENTIFIED",
//      "external_ref": {
//        "_type": "PARTY_REF",
//        "id": {
//          "_type": "GENERIC_ID",
//          "value": "9091",
//          "scheme": "HOSPITAL-NS"
//        },
//        "namespace": "HOSPITAL-NS",
//        "type": "PARTY"
//      },
//      "name": "Hospital"
//    },
//    "participations": [
//      {
//        "_type": "PARTICIPATION",
//        "function": {
//          "_type": "DV_TEXT",
//          "value": "requester"
//        },
//        "performer": {
//          "_type": "PARTY_IDENTIFIED",
//          "external_ref": {
//            "_type": "PARTY_REF",
//            "id": {
//              "_type": "GENERIC_ID",
//              "value": "199",
//              "scheme": "HOSPITAL-NS"
//            },
//            "namespace": "HOSPITAL-NS",
//            "type": "PERSON"
//          },
//          "name": "Dr. Marcus Johnson"
//        },
//        "mode": {
//          "_type": "DV_CODED_TEXT",
//          "value": "face-to-face communication",
//          "defining_code": {
//            "_type": "CODE_PHRASE",
//            "terminology_id": {
//              "_type": "TERMINOLOGY_ID",
//              "value": "openehr"
//            },
//            "code_string": "216"
//          }
//        }
//      },
//      {
//        "_type": "PARTICIPATION",
//        "function": {
//          "_type": "DV_TEXT",
//          "value": "performer"
//        },
//        "performer": {
//          "_type": "PARTY_IDENTIFIED",
//          "external_ref": {
//            "_type": "PARTY_REF",
//            "id": {
//              "_type": "GENERIC_ID",
//              "value": "198",
//              "scheme": "HOSPITAL-NS"
//            },
//            "namespace": "HOSPITAL-NS",
//            "type": "PERSON"
//          },
//          "name": "Lara Markham"
//        },
//        "mode": {
//          "_type": "DV_CODED_TEXT",
//          "value": "not specified",
//          "defining_code": {
//            "_type": "CODE_PHRASE",
//            "terminology_id": {
//              "_type": "TERMINOLOGY_ID",
//              "value": "openehr"
//            },
//            "code_string": "193"
//          }
//        }
//      }
//    ]
//  },
//  "content": [
//	{
//      "_type": "SECTION",
//      "name": {
//        "_type": "DV_TEXT",
//        "value": "Medication Summary"
//      },
//      "archetype_details": {
//        "_type": "ARCHETYPED",
//        "archetype_id": {
//          "_type": "ARCHETYPE_ID",
//          "value": "openEHR-EHR-SECTION.adhoc.v1"
//        },
//        "rm_version": "1.0.4"
//      },
//      "archetype_node_id": "openEHR-EHR-SECTION.adhoc.v1",
//      "items": []
//    }
//  ]
//}`
