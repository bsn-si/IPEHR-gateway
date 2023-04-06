package treeindex

import (
	"encoding/json"
	"os"
	"testing"

	"github.com/bsn-si/IPEHR-gateway/src/pkg/docs/model"
	"github.com/bsn-si/IPEHR-gateway/src/pkg/docs/model/base"
	"github.com/bsn-si/IPEHR-gateway/src/pkg/errors"
	"github.com/vmihailenco/msgpack/v5"

	"github.com/stretchr/testify/assert"
)

func TestEHRIndex_AddEHR(t *testing.T) {
	tests := []struct {
		name    string
		getEHR  func() (model.EHR, error)
		want    map[string]*EHRNode
		wantErr bool
	}{
		{
			"1. success add new EHR object",
			func() (model.EHR, error) {
				ehr, err := loadEHRFromFile("./test_fixtures/ehr.json")
				if err != nil {
					return model.EHR{}, err
				}

				return ehr, nil
			},
			map[string]*EHRNode{
				"7d44b88c-4199-4bad-97dc-d78268e01398": {
					BaseNode: BaseNode{
						ID:       "7d44b88c-4199-4bad-97dc-d78268e01398",
						Type:     base.EHRItemType,
						NodeType: EHRNodeType,
					},
					Attributes: Attributes{
						"system_id": newValueNode("d60e2348-b083-48ce-93b9-916cef1d3a5a"),
						"ehr_id":    newValueNode("7d44b88c-4199-4bad-97dc-d78268e01398"),
					},
					Compositions: Container{},
				},
			},
			false,
		},
		{
			"2. EHR with simple Composition",
			func() (model.EHR, error) {
				ehr, err := loadEHRFromFile("./test_fixtures/ehr_with_composition.json")
				if err != nil {
					return model.EHR{}, err
				}

				return ehr, nil
			},
			map[string]*EHRNode{
				"7d44b88c-4199-4bad-97dc-d78268e01398": {
					BaseNode: BaseNode{
						ID:       "7d44b88c-4199-4bad-97dc-d78268e01398",
						Type:     base.EHRItemType,
						NodeType: EHRNodeType,
					},
					Attributes: Attributes{
						"system_id": newValueNode("d60e2348-b083-48ce-93b9-916cef1d3a5a"),
						"ehr_id":    newValueNode("7d44b88c-4199-4bad-97dc-d78268e01398"),
					},
					Compositions: Container{
						"openEHR-EHR-COMPOSITION.health_summary.v1": []Noder{
							&CompositionNode{
								BaseNode: BaseNode{
									ID:       "openEHR-EHR-COMPOSITION.health_summary.v1",
									Type:     base.CompositionItemType,
									Name:     "International Patient Summary",
									NodeType: CompostionNodeType,
								},
								Tree: *NewTree(),
								Attributes: Attributes{
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
								},
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
			idx := NewEHRIndex()
			ehr, err := tt.getEHR()
			if err != nil {
				t.Errorf("EHRIndex.AddEHR(), cannot get EHR: %v", err)
				return
			}

			if err := idx.AddEHR(&ehr); (err != nil) != tt.wantErr {
				t.Errorf("EHRIndex.AddEHR() error = %v, wantErr %v", err, tt.wantErr)
			}

			assert.Equal(t, tt.want, idx.Ehrs)
		})
	}
}

func TestEHRIndex_MessagePack(t *testing.T) {
	tests := []struct {
		name    string
		getEHR  func() (model.EHR, error)
		wantErr bool
	}{
		{
			"1. success add new EHR object",
			func() (model.EHR, error) {
				ehr, err := loadEHRFromFile("./test_fixtures/ehr.json")
				if err != nil {
					return model.EHR{}, err
				}

				return ehr, nil
			},
			false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			idx := NewEHRIndex()
			ehr, err := tt.getEHR()
			if err != nil {
				t.Errorf("EHRIndex.AddEHR(), cannot get EHR: %v", err)
				return
			}

			if err := idx.AddEHR(&ehr); (err != nil) != tt.wantErr {
				t.Errorf("EHRIndex.AddEHR() error = %v, wantErr %v", err, tt.wantErr)
			}

			data, err := msgpack.Marshal(idx)
			if err != nil {
				t.Errorf("EHRIndex.Marshal() error = %v, wantErr %v", err, tt.wantErr)
			}

			gotIdx := NewEHRIndex()
			if err := msgpack.Unmarshal(data, &gotIdx); err != nil {
				t.Errorf("EHRIndex.Unmarshal() error = %v, wantErr %v", err, tt.wantErr)
			}

			assert.Equal(t, idx, gotIdx)
		})
	}
}

func loadEHRFromFile(name string) (model.EHR, error) {
	data, err := os.ReadFile(name)
	if err != nil {
		return model.EHR{}, errors.Wrap(err, "cannot read file")
	}

	ehr := model.EHR{}
	if err := json.Unmarshal(data, &ehr); err != nil {
		return model.EHR{}, errors.Wrap(err, "cannot unmarshal EHR")
	}

	return ehr, nil
}
