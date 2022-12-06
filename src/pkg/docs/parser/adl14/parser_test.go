package adl14

import (
	"hms/gateway/pkg/docs/model"
	"testing"
)

func TestParser_Validate(t *testing.T) {
	p := NewParser()

	t.Run("Invalid XML", func(t *testing.T) {
		r := p.Validate([]byte(`<xml>invalid>/`), model.ADLTypeXML)

		if r {
			t.Errorf("adl14.Validate() xml is valid")
			return
		}
	})

	t.Run("Valid XML", func(t *testing.T) {
		r := p.Validate([]byte(`<xml></xml>`), model.ADLTypeXML)

		if !r {
			t.Errorf("adl14.Validate() xml is invalid")
			return
		}
	})
}
