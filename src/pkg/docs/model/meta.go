package model

import (
	"hms/gateway/pkg/indexer/ehrIndexer"
)

type Attribute = uint8

const (
	AttributeCID             Attribute = 1
	AttributeCIDEncr         Attribute = 2
	AttributeKeyEncr         Attribute = 3
	AttributeDocBaseUIDHash  Attribute = 4
	AttributeDocUIDEncrypted Attribute = 5
	AttributeDealCid         Attribute = 6
	AttributeMinerAddress    Attribute = 7
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
