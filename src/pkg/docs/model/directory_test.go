package model_test

import (
	"encoding/json"
	"testing"

	"github.com/bsn-si/IPEHR-gateway/src/pkg/docs/model"
	"github.com/bsn-si/IPEHR-gateway/src/pkg/docs/model/base"

	"github.com/google/go-cmp/cmp"
)

type directoryTestData struct {
	d    model.Directory
	JSON []byte
}

func TestDirectory_UnmarshalJSON(t *testing.T) {
	var tests = []struct {
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
		{
			"4. directory with subfolders",
			directoryTestData{
				d: model.Directory{
					Locatable: base.Locatable{
						Type:            base.FolderItemType,
						Name:            base.NewDvText("root"),
						ArchetypeNodeID: "openEHR-EHR-FOLDER.generic.v1",
					},
					FeederAudit: base.FeederAudit{},
					Details:     base.ItemStructure{},
					Folders: []*model.Directory{
						{
							Locatable: base.Locatable{
								Type:            base.FolderItemType,
								Name:            base.NewDvText("emergency"),
								ArchetypeNodeID: "openEHR-EHR-FOLDER.generic.v1",
							},
							FeederAudit: base.FeederAudit{},
							Details:     base.ItemStructure{},
							Folders: []*model.Directory{
								{
									Locatable: base.Locatable{
										Type:            base.FolderItemType,
										Name:            base.NewDvText("episode_x"),
										ArchetypeNodeID: "openEHR-EHR-FOLDER.generic.v1",
									},
									FeederAudit: base.FeederAudit{},
									Details:     base.ItemStructure{},
									Folders: []*model.Directory{
										{
											Locatable: base.Locatable{
												Type:            base.FolderItemType,
												Name:            base.NewDvText("summary_compo_x"),
												ArchetypeNodeID: "openEHR-EHR-FOLDER.generic.v1",
											},
											FeederAudit: base.FeederAudit{},
											Details:     base.ItemStructure{},
											Folders:     nil,
										},
									},
								},
								{
									Locatable: base.Locatable{
										Type:            base.FolderItemType,
										Name:            base.NewDvText("episode_y"),
										ArchetypeNodeID: "openEHR-EHR-FOLDER.generic.v1",
									},
									FeederAudit: base.FeederAudit{},
									Details:     base.ItemStructure{},
									Folders: []*model.Directory{
										{
											Locatable: base.Locatable{
												Type:            base.FolderItemType,
												Name:            base.NewDvText("summary_compo_y"),
												ArchetypeNodeID: "openEHR-EHR-FOLDER.generic.v1",
											},
											FeederAudit: base.FeederAudit{},
											Details:     base.ItemStructure{},
											Folders:     nil,
										},
									},
								},
							},
						},
						{
							Locatable: base.Locatable{
								Type:            base.FolderItemType,
								Name:            base.NewDvText("hospitalization"),
								ArchetypeNodeID: "openEHR-EHR-FOLDER.generic.v1",
							},
							FeederAudit: base.FeederAudit{},
							Details:     base.ItemStructure{},
							Folders: []*model.Directory{
								{
									Locatable: base.Locatable{
										Type:            base.FolderItemType,
										Name:            base.NewDvText("summary_compo_z"),
										ArchetypeNodeID: "openEHR-EHR-FOLDER.generic.v1",
									},
									FeederAudit: base.FeederAudit{},
									Details:     base.ItemStructure{},
									Folders:     nil,
								},
							},
						},
						{
							Locatable: base.Locatable{
								Type:            base.FolderItemType,
								Name:            base.NewDvText("foldername-w-special-chars"),
								ArchetypeNodeID: "openEHR-EHR-FOLDER.generic.v1",
							},
							FeederAudit: base.FeederAudit{},
							Details:     base.ItemStructure{},
							Folders:     nil,
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
					  "folders": [
						{
						  "_type": "FOLDER",
						  "name": {
							"_type": "DV_TEXT",
							"value": "emergency"
						  },
						  "archetype_node_id": "openEHR-EHR-FOLDER.generic.v1",
						  "folders": [
							{
							  "_type": "FOLDER",
							  "name": {
								"_type": "DV_TEXT",
								"value": "episode_x"
							  },
							  "archetype_node_id": "openEHR-EHR-FOLDER.generic.v1",
							  "folders": [
								{
								  "_type": "FOLDER",
								  "name": {
									"_type": "DV_TEXT",
									"value": "summary_compo_x"
								  },
								  "archetype_node_id": "openEHR-EHR-FOLDER.generic.v1"
								}
							  ]
							},
							{
							  "_type": "FOLDER",
							  "name": {
								"_type": "DV_TEXT",
								"value": "episode_y"
							  },
							  "archetype_node_id": "openEHR-EHR-FOLDER.generic.v1",
							  "folders": [
								{
								  "_type": "FOLDER",
								  "name": {
									"_type": "DV_TEXT",
									"value": "summary_compo_y"
								  },
								  "archetype_node_id": "openEHR-EHR-FOLDER.generic.v1"
								}
							  ]
							}
						  ]
						},
						{
						  "_type": "FOLDER",
						  "name": {
							"_type": "DV_TEXT",
							"value": "hospitalization"
						  },
						  "archetype_node_id": "openEHR-EHR-FOLDER.generic.v1",
						  "folders": [
							{
							  "_type": "FOLDER",
							  "name": {
								"_type": "DV_TEXT",
								"value": "summary_compo_z"
							  },
							  "archetype_node_id": "openEHR-EHR-FOLDER.generic.v1"
							}
						  ]
						},
						{
						  "_type": "FOLDER",
						  "name": {
							"_type": "DV_TEXT",
							"value": "foldername-w-special-chars"
						  },
						  "archetype_node_id": "openEHR-EHR-FOLDER.generic.v1"
						}
					  ]
					}
				`),
			},
			false,
		},
		{
			name: "5. empty directory with details",
			data: directoryTestData{
				d: model.Directory{
					Locatable: base.Locatable{
						Type:            base.FolderItemType,
						Name:            base.NewDvText("root"),
						ArchetypeNodeID: "openEHR-EHR-FOLDER.generic.v1",
					},
					FeederAudit: base.FeederAudit{},
					Folders:     nil,
					Details: base.ItemStructure{Data: &base.ItemTree{
						DataStructure: base.DataStructure{
							Locatable: base.Locatable{
								Type:            base.ItemTreeItemType,
								Name:            base.NewDvText("Tree"),
								ArchetypeNodeID: "at0003",
							},
						},
						Items: base.Items{
							&base.Element{
								Item: base.Item{
									Type:            base.ElementItemType,
									Name:            base.NewDvText("text"),
									ArchetypeNodeID: "at0004",
								},
								Value: toRef(base.NewDvText("Lorem ipsum dolor sit amet")),
							},
						}}},
				},
				JSON: []byte(`
					{
					  "_type": "FOLDER",
					  "name": {
						"_type": "DV_TEXT",
						"value": "root"
					  },
					  "archetype_node_id": "openEHR-EHR-FOLDER.generic.v1",
					  "details": {
						"_type": "ITEM_TREE",
						"name": {
						  "_type": "DV_TEXT",
						  "value": "Tree"
						},
						"archetype_node_id": "at0003",
						"items": [
						  {
							"_type": "ELEMENT",
							"name": {
							  "_type": "DV_TEXT",
							  "value": "text"
							},
							"archetype_node_id": "at0004",
							"value": {
							  "_type": "DV_TEXT",
							  "value": "Lorem ipsum dolor sit amet"
							}
						  }
						]
					  }
					}
				`),
			},
		},
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

func TestDirectory_GetByPath(t *testing.T) {
	d := &model.Directory{
		Locatable: base.Locatable{
			Type:            base.FolderItemType,
			Name:            base.NewDvText("root"),
			ArchetypeNodeID: "openEHR-EHR-FOLDER.generic.v1",
		},
		Folders: []*model.Directory{
			{
				Locatable: base.Locatable{
					Type:            base.FolderItemType,
					Name:            base.NewDvText("1"),
					ArchetypeNodeID: "openEHR-EHR-FOLDER.generic.v1",
				},
				Folders: []*model.Directory{
					{
						Locatable: base.Locatable{
							Type:            base.FolderItemType,
							Name:            base.NewDvText("1-1"),
							ArchetypeNodeID: "openEHR-EHR-FOLDER.generic.v1",
						},
						Folders: []*model.Directory{
							{
								Locatable: base.Locatable{
									Type:            base.FolderItemType,
									Name:            base.NewDvText("1-1-1"),
									ArchetypeNodeID: "openEHR-EHR-FOLDER.generic.v1",
								},
							},
							{
								Locatable: base.Locatable{
									Type:            base.FolderItemType,
									Name:            base.NewDvText("1-1-2"),
									ArchetypeNodeID: "openEHR-EHR-FOLDER.generic.v1",
								},
							},
						},
					},
					{
						Locatable: base.Locatable{
							Type:            base.FolderItemType,
							Name:            base.NewDvText("1-2"),
							ArchetypeNodeID: "openEHR-EHR-FOLDER.generic.v1",
						},
						FeederAudit: base.FeederAudit{},
						Details:     base.ItemStructure{},
						Folders: []*model.Directory{
							{
								Locatable: base.Locatable{
									Type:            base.FolderItemType,
									Name:            base.NewDvText("1-2-1"),
									ArchetypeNodeID: "openEHR-EHR-FOLDER.generic.v1",
								},
								FeederAudit: base.FeederAudit{},
								Details:     base.ItemStructure{},
								Folders:     nil,
							},
						},
					},
				},
			},
			{
				Locatable: base.Locatable{
					Type:            base.FolderItemType,
					Name:            base.NewDvText("2"),
					ArchetypeNodeID: "openEHR-EHR-FOLDER.generic.v1",
				},
				FeederAudit: base.FeederAudit{},
				Details:     base.ItemStructure{},
				Folders: []*model.Directory{
					{
						Locatable: base.Locatable{
							Type:            base.FolderItemType,
							Name:            base.NewDvText("2-1"),
							ArchetypeNodeID: "openEHR-EHR-FOLDER.generic.v1",
						},
						FeederAudit: base.FeederAudit{},
						Details:     base.ItemStructure{},
						Folders:     nil,
					},
				},
			},
			{
				Locatable: base.Locatable{
					Type:            base.FolderItemType,
					Name:            base.NewDvText("3"),
					ArchetypeNodeID: "openEHR-EHR-FOLDER.generic.v1",
				},
				FeederAudit: base.FeederAudit{},
				Details:     base.ItemStructure{},
				Folders:     nil,
			},
		},
	}

	var tests = []struct {
		name     string
		path     string
		wantName string
		wantErr  bool
	}{
		{
			name:     "1. Error because path is empty",
			path:     "",
			wantName: "root",
			wantErr:  false,
		}, {
			name:     "2. Error because path not found",
			path:     "unknown folder name",
			wantName: "",
			wantErr:  true,
		}, {
			name:     "3. Successfully find root",
			path:     "root",
			wantName: "root",
			wantErr:  false,
		}, {
			name:     "4. Successfully find first sub folder",
			path:     "root/1",
			wantName: "1",
			wantErr:  false,
		}, {
			name:     "5. Successfully find last folder",
			path:     "root/1/1-1/1-1-2",
			wantName: "1-1-2",
			wantErr:  false,
		}, {
			name:     "6. Successfully find folder in sub folder",
			path:     "root/2/2-1",
			wantName: "2-1",
			wantErr:  false,
		}, {
			name:     "7. Successfully find last folder",
			path:     "root/3",
			wantName: "3",
			wantErr:  false,
		}, {
			name:     "8. Successfully trim slashes and find folder",
			path:     "//////////root///3////////",
			wantName: "3",
			wantErr:  false,
		}, {
			name:     "9. Fail because havent found folder",
			path:     "root/1/1-1/not exist folder",
			wantName: "not exist folder",
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotD, err := d.GetByPath(tt.path)
			if err != nil {
				if tt.wantErr {
					return
				}

				t.Errorf("Directory.GetByPath() have error %v", err)
				return
			}

			if diff := cmp.Diff(gotD.Name.Value, tt.wantName); diff != "" {
				t.Errorf("Directory.GetByPath() mismatch{-want;+got}\n\t%s", diff)
			}
		})
	}
}
