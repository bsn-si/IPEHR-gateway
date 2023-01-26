package treeindex

import (
	"encoding/json"
	"testing"

	"github.com/bsn-si/IPEHR-gateway/src/pkg/docs/model"
	"github.com/bsn-si/IPEHR-gateway/src/pkg/docs/model/base"

	"github.com/stretchr/testify/assert"
)

func Test_processEventContext(t *testing.T) {
	tests := []struct {
		name    string
		getCtx  func() (model.EventContext, error)
		want    Noder
		wantErr bool
	}{
		{
			"1. parse simple contex",
			getEventContext,
			&EventContextNode{
				BaseNode: BaseNode{
					NodeType: EventContextNodeType,
				},
				Attributes: Attributes{
					"start_time": newNode(&base.DvDateTime{
						DvTemporal: base.DvTemporal{
							DvValueBase: base.DvValueBase{
								Type: base.DvDateTimeItemType,
							},
						},
						Value: "2021-12-03T17:34:06.849379+01:00",
					}),
					"end_time": newNode(&base.DvDateTime{
						DvTemporal: base.DvTemporal{
							DvValueBase: base.DvValueBase{
								Type: base.DvDateTimeItemType,
							},
						},
						Value: "2021-12-03T17:34:06.849379+01:00",
					}),
					"location": newNode("some_text_here"),
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
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx, err := tt.getCtx()
			if err != nil {
				t.Fatal(err)
			}

			got, err := processEventContext(ctx)
			if (err != nil) != tt.wantErr {
				t.Errorf("processEventContext() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			assert.Equal(t, tt.want, got)
		})
	}
}

func getEventContext() (model.EventContext, error) {
	const jsonStr = `{
        "_type": "EVENT_CONTEXT",
        "start_time": {
            "_type": "DV_DATE_TIME",
            "value": "2021-12-03T17:34:06.849379+01:00"
        },
        "end_time": {
            "_type": "DV_DATE_TIME",
            "value": "2021-12-03T17:34:06.849379+01:00"
        },
		"location": "some_text_here",
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
        }
    }`

	ctx := model.EventContext{}
	if err := json.Unmarshal([]byte(jsonStr), &ctx); err != nil {
		return model.EventContext{}, err
	}

	return ctx, nil
}
