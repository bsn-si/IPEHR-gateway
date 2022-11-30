package model

import (
	"hms/gateway/pkg/indexer/ehrIndexer"
)

type Attribute = uint8

const (
	AttributeID             Attribute = 1
	AttributeIDEncr         Attribute = 2
	AttributeKeyEncr        Attribute = 3
	AttributeDocBaseUIDHash Attribute = 4
	AttributeDocUIDEncr     Attribute = 5
	AttributeDealCid        Attribute = 6
	AttributeMinerAddress   Attribute = 7
	AttributeContentEncr    Attribute = 8
)

type DocumentMeta ehrIndexer.DocsDocumentMeta

func (dm *DocumentMeta) GetAttr(attr Attribute) []byte {
	for _, a := range dm.Attrs {
		if a.Code == attr {
			return a.Value
		}
	}

	return nil
}
