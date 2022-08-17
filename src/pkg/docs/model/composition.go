package model

import (
	"bytes"
	"encoding/json"
	"hms/gateway/pkg/docs/model/base"
	"hms/gateway/pkg/docs/types"
	"io/ioutil"
)

// Composition Content of one version in a VERSIONED_COMPOSITION. A Composition is considered the unit
// of modification of the record, the unit of transmission in record Extracts, and the unit of
// attestation by authorising clinicians. In this latter sense, it may be considered equivalent to a
// signed document.
// https://specifications.openehr.org/releases/RM/latest/ehr.html#_composition_class
type Composition struct {
	Language  base.CodePhrase  `json:"language"`
	Territory base.CodePhrase  `json:"territory"`
	Category  base.DvCodedText `json:"category"`
	Context   *EventContext    `json:"context,omitempty"`
	Composer  base.PartyProxy  `json:"composer"`
	Content   interface{}      `json:"content,omitempty"`
	base.Locatable
}

func (c *Composition) Validate() bool {
	validation := true
	if c.Type != types.Composition.String() {
		validation = false
	}

	return validation
}

func (c *Composition) FromJSON(reader *bytes.Reader) (err error) {
	data, err := ioutil.ReadAll(reader)
	if err == nil {
		err = json.Unmarshal(data, &c)
	}

	//c.prepare()
	return
}

//func (c *Composition) prepare() {
// TODO we can move logic here like initialization, e.g.:
//	c.ObjectVersionID.New(c.UID.Value, cfg.CreatingSystemID)
// TODO but in what case we need global variables and its look like bad arch
//}
