package model_test

import (
	"encoding/json"
	"hms/gateway/pkg/docs/model"
	"hms/gateway/pkg/docs/model/base"
	"testing"

	"github.com/google/go-cmp/cmp"
)

type directoryTestData struct {
	d    model.Directory
	JSON []byte
}

func TestDirectory_UnmarshalJSON(t *testing.T) {
	tests := []struct {
		name    string
		data    directoryTestData
		wantErr bool
	}{
		{
			"1. error on unmarshal data",
			directoryTestData{},
			true,
		},
		{
			"2. empty directory",
			directoryTestData{
				d: model.Directory{
					Locatable: base.Locatable{Type: base.FolderItemType, Name: base.NewDvText("root"),
						ArchetypeNodeID: "openEHR-EHR-FOLDER.generic.v1"},
					FeederAudit: base.FeederAudit{},
					Folders:     nil,
					Details:     base.ItemStructure{},
					Items:       nil,
				},
				JSON: []byte(`
					{
					  "_type": "FOLDER",
					  "name": {
						"_type": "DV_TEXT",
						"value": "root"
					  },
					  "archetype_node_id": "openEHR-EHR-FOLDER.generic.v1"
					}
				`),
			},
			false,
		},
		{
			"3. empty directory with items",
			directoryTestData{
				d: model.Directory{
					Locatable: base.Locatable{Type: base.FolderItemType, Name: base.NewDvText("root"),
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
				},
				JSON: []byte(`
					{
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
					}
				`),
			},
			false,
		},
		//{
		//	"4. empty directory with details",
		//	directoryTestData{
		//		d: model.Directory{
		//			Locatable: base.Locatable{Type: base.FolderItemType, Name: base.NewDvText("root"),
		//				ArchetypeNodeID: "openEHR-EHR-FOLDER.generic.v1"},
		//			FeederAudit: base.FeederAudit{},
		//			Folders:     nil,
		//			Details:     base.ItemStructure{},
		//			Items: []model.DirectoryItem{
		//				{
		//					ID: base.UIDBasedID{
		//						ObjectID: base.ObjectID{
		//							Type:  base.HierObjectIDItemType,
		//							Value: "replaceme",
		//						},
		//					},
		//					Type:      base.VersionCompositionItemType,
		//					Namespace: "my.system.id",
		//				},
		//			},
		//		},
		//		JSON: []byte(`
		//			{
		//			  "_type": "FOLDER",
		//			  "name": {
		//				"_type": "DV_TEXT",
		//				"value": "root"
		//			  },
		//			  "archetype_node_id": "openEHR-EHR-FOLDER.generic.v1",
		//			  "items": [
		//				{
		//					"id": {
		//						"_type": "HIER_OBJECT_ID",
		//						"value": "replaceme"
		//					},
		//					"namespace": "my.system.id",
		//					"type": "VERSIONED_COMPOSITION"
		//				}
		//			  ]
		//			}
		//		`),
		//	},
		//	false,
		//},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := model.Directory{}
			if err := json.Unmarshal(tt.data.JSON, &got); (err != nil) != tt.wantErr {
				t.Errorf("Directory.UnmarshalJSON() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			opts := cmp.AllowUnexported(
				base.ObjectVersionID{},
				base.PartyProxy{},
			)
			if diff := cmp.Diff(tt.data.d, got, opts); diff != "" {
				t.Errorf("Directory.UnmarshalJSON() mismatch{-want;+got}\n\t%s", diff)
			}
		})
	}
}
