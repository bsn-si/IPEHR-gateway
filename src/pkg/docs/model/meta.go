package model

import "hms/gateway/pkg/indexer/ehrIndexer"

type DocumentMeta ehrIndexer.DocsDocumentMeta

func (dm *DocumentMeta) GetAttr(code Attribute) []byte {
	for _, attr := range dm.Attrs {
		if attr.Code == code {
			return attr.Value
		}
	}

	return nil
}
