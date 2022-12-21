package treeindex

import (
	"encoding/json"
	"hms/gateway/pkg/docs/model"
	"hms/gateway/pkg/docs/model/base"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_processComposition(t *testing.T) {
	tests := []struct {
		name    string
		getCmp  func() (model.Composition, error)
		want    Noder
		wantErr bool
	}{
		{
			"1. simple composition",
			func() (model.Composition, error) {
				return loadComposition("./test_fixtures/simple_composition.json")
			},
			&CompositionNode{
				baseNode: baseNode{
					ID:   "openEHR-EHR-COMPOSITION.health_summary.v1",
					Type: base.CompositionItemType,
					Name: "International Patient Summary",
				},
				Tree: Tree{
					actions:       Container{},
					evaluations:   Container{},
					instructions:  Container{},
					obeservations: Container{},
				},
				attributes: map[string]Noder{
					"language": newNode(&base.CodePhrase{
						Type: base.CodePhraseItemType,
						TerminologyID: base.ObjectID{
							Type:  base.TerminologyIDItemType,
							Value: "ISO_639-1",
						},
						CodeString: "en",
					}),
					"territory": newNode(&base.CodePhrase{
						Type: base.CodePhraseItemType,
						TerminologyID: base.ObjectID{
							Type:  base.TerminologyIDItemType,
							Value: "ISO_3166-1",
						},
						CodeString: "US",
					}),
					"category": newNode(base.NewDvCodedText("event", base.CodePhrase{
						Type: base.CodePhraseItemType,
						TerminologyID: base.ObjectID{
							Type:  base.TerminologyIDItemType,
							Value: "openehr",
						},
						CodeString: "443",
					})),
					"context": &EventContextNode{
						attributes: map[string]Noder{
							"start_time": newNode(&base.DvDateTime{
								DvTemporal: base.DvTemporal{
									DvValueBase: base.DvValueBase{
										Type: base.DvDateTimeItemType,
									},
								},
								Value: "2021-12-03T17:34:06.849379+01:00",
							}),
							"setting": newNode(base.NewDvCodedText("other care", base.CodePhrase{
								Type: base.CodePhraseItemType,
								TerminologyID: base.ObjectID{
									Type:  base.TerminologyIDItemType,
									Value: "openehr",
								},
								CodeString: "238",
							})),
						},
					},
				},
			},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmp, err := tt.getCmp()
			if err != nil {
				t.Errorf("processComposition() load Composition error: %v", err)
				return
			}

			got, err := processComposition(cmp)
			if (err != nil) != tt.wantErr {
				t.Errorf("processComposition() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			assert.Equal(t, tt.want, got)
		})
	}
}

func loadComposition(name string) (model.Composition, error) {
	data, err := os.ReadFile(name)
	if err != nil {
		return model.Composition{}, err
	}

	cmp := model.Composition{}
	if err := json.Unmarshal(data, &cmp); err != nil {
		return model.Composition{}, err
	}

	return cmp, nil
}
