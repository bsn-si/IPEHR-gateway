package treeindex

import (
	"encoding/json"
	"hms/gateway/pkg/docs/model"
	"hms/gateway/pkg/errors"
)

var DefaultEHRIndex = NewEHRIndex()

type EHRIndex struct {
	ehrs Container
}

func NewEHRIndex() *EHRIndex {
	idx := EHRIndex{
		ehrs: Container{},
	}

	return &idx
}

func (idx *EHRIndex) AddEHR(ehr model.EHR) error {
	node, err := walk(ehr)

	if err != nil {
		return errors.Wrap(err, "cannot add EHR object")
	}

	idx.ehrs[node.GetID()] = append(idx.ehrs[node.GetID()], node)

	return nil
}

func (idx EHRIndex) MarshalJSON() ([]byte, error) {
	return json.Marshal(idx.ehrs)
}
