package adl14

import (
	"testing"

	"github.com/bsn-si/IPEHR-gateway/src/pkg/docs/model"
)

func TestParser_Validate(t *testing.T) {
	tests := []struct {
		name   string
		xmlStr string
		want   bool
	}{
		{
			"1. Invalid XML",
			`<xml>invalid/`,
			false,
		},
		{
			"2. Valid XML",
			`<xml>val</xml>`,
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			parser := NewParser()
			if parser.Validate([]byte(tt.xmlStr), model.ADLTypeXML) != tt.want {
				t.Error("adl14.Validate() error")
			}
		})
	}
}
