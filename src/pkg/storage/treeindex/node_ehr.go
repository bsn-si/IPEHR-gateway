package treeindex

import "hms/gateway/pkg/docs/model"

func processEHR(ehr *model.EHR) (Noder, error) {
	return nil, nil
}

type EHRIndex struct {
	ehrs map[string]Noder
}

type EHRNode struct {
}

func newEHRNode(ehr model.EHR) Noder {
	return &EHRNode{}
}

func (ehr EHRNode) GetNodeType() NodeType {
	return EHRNodeType
}

func (ehr EHRNode) GetID() string {
	return ""
}

func (ehr EHRNode) TryGetChild(key string) Noder {
	return nil
}

func (ehr EHRNode) ForEach(func(name string, node Noder) bool) {
}

func (ehr EHRNode) addAttribute(key string, val Noder) {

}
