package model_test

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/bsn-si/IPEHR-gateway/src/pkg/docs/model"
	"github.com/bsn-si/IPEHR-gateway/src/pkg/docs/model/base"

	"testing"

	"github.com/google/go-cmp/cmp"
)

type contributionTestData struct {
	c     model.Contribution
	cJSON []byte
}

type contributionVersionTestData struct {
	cv     base.Root
	cvJSON []byte
}

func TestContribution_UnmarshalJSON(t *testing.T) {
	tests := []struct {
		name    string
		data    contributionTestData
		wantErr bool
	}{
		{
			"1. error on unmarshal data",
			contributionTestData{},
			true,
		},
		{
			"2. empty contribution",
			newContribution(),
			false,
		}, {
			"3. contribution with composition",
			newContributionWithVersions([]contributionVersionTestData{
				{expectedComposition, []byte(compositionJSON)},
			}),
			false,
		}, {
			"4. contribution with composition and folder",
			newContributionWithVersions([]contributionVersionTestData{
				{expectedComposition, []byte(compositionJSON)},
				{model.Directory{
					Locatable: base.Locatable{
						Type:            base.FolderItemType,
						Name:            base.NewDvText("root"),
						ArchetypeNodeID: "openEHR-EHR-FOLDER.generic.v1"},
					FeederAudit: base.FeederAudit{},
					Folders:     nil,
					Details:     base.ItemStructure{},
					Items: []model.DirectoryItem{
						{
							ID: base.UIDBasedID{
								ObjectID: base.ObjectID{
									Type:  base.HierObjectIDItemType,
									Value: "replaceme",
								},
							},
							Type:      base.VersionCompositionItemType,
							Namespace: "my.system.id",
						},
					},
				}, []byte(`{
					  "_type": "FOLDER",
					  "name": {
						"_type": "DV_TEXT",
						"value": "root"
					  },
					  "archetype_node_id": "openEHR-EHR-FOLDER.generic.v1",
					  "items": [
						{
							"id": {
								"_type": "HIER_OBJECT_ID",
								"value": "replaceme"
							},
							"namespace": "my.system.id",
							"type": "VERSIONED_COMPOSITION"
						}
					  ]
					}`)},
			}),
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := model.Contribution{}
			if err := json.Unmarshal(tt.data.cJSON, &got); (err != nil) != tt.wantErr {
				t.Errorf("Contribution.UnmarshalJSON() error = %v, wantErr %v", err, tt.wantErr)
			}

			opts := cmp.AllowUnexported(
				base.ObjectVersionID{},
				base.PartyProxy{},
			)
			if diff := cmp.Diff(tt.data.c, got, opts); diff != "" {
				t.Errorf("Composition.UnmarshalJSON() mismatch{-want;+got}\n\t%s", diff)
			}
		})
	}
}

func newContribution() contributionTestData {
	c := model.Contribution{
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
								//"scheme": "<ID SCHEME NAME>" TODO no idea, in docs i didnt find information, but it has in tests in EHRBASE
							},
							Namespace: "DEMOGRAPHIC",
							Type:      "PERSON",
						},
					},
				},
			),
			// TODO ChangeType into struct?
			ChangeType: base.DvCodedText{
				DefiningCode: base.CodePhrase{
					TerminologyID: base.ObjectID{Value: "openehr"},
					CodeString:    "249",
				},
				DvText: base.DvText{Value: "creation"},
			},
			Description: base.DvText{Value: "<optional audit description>"},
		},

		Versions: []model.ContributionVersion{},
	}

	return contributionTestData{c, prepareContributionJSON("")}
}

func prepareContributionJSON(v string) []byte {
	if v == "" {
		v = "[]"
	}

	return []byte(fmt.Sprintf(`{
		"_type": "CONTRIBUTION",
		"versions": %s,
		"audit": {
			"_type": "AUDIT_DETAILS",
			"system_id": "test-system-id",
			"committer": {
				"_type": "PARTY_IDENTIFIED",
				"name": "<optional name of the committer>",
				"external_ref": {
					"id": {
						"_type": "GENERIC_ID",
						"value": "<OBJECT_ID>",
						"scheme": "<ID SCHEME NAME>"
					},
					"namespace": "DEMOGRAPHIC",
					"type": "PERSON"
				}
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
			"description": {
				"value": "<optional audit description>"
			}
		}
	}`, v))
}

func prepareContributionVersionJSON(v []byte) []byte {
	if v == nil {
		v = []byte("{}")
	}

	return []byte(fmt.Sprintf(`{
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
			},
			"description": {
				"value": "<optional audit description>"
			}
		},
		"uid": {
			"_type": "OBJECT_VERSION_ID",
			"value": "41f6fdb5-9ea5-4bb8-b2fa-21131543f82e::openEHRSys.example.com::1"
		},
		"data": %s,
		"lifecycle_state": {
			"_type": "DV_CODED_TEXT",
			"value": "complete",
			"defining_code": {
				"terminology_id": {
					"_type": "TERMINOLOGY_ID",
					"value": "openehr"
				},
				"code_string": "532"
			}
		}

	}`, v))
}

func newContributionWithVersions(data []contributionVersionTestData) contributionTestData {
	c := newContribution()
	cJSON := make([]string, 0, len(data))

	for _, d := range data {
		cVer := model.ContributionVersion{
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
				TimeCommitted: base.DvDateTime{
					Value: "2021-09-21T21:52:31.869-03:00",
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
					TerminologyID: base.ObjectID{
						Type:  "TERMINOLOGY_ID",
						Value: "openehr",
					},
					CodeString: "532",
				},
			),
			Data: d.cv,
		}

		c.c.Versions = append(c.c.Versions, cVer)
		cJSON = append(cJSON, string(prepareContributionVersionJSON(d.cvJSON)))
	}

	c.cJSON = prepareContributionJSON("[" + strings.Join(cJSON, ",") + "]")

	return c
}
