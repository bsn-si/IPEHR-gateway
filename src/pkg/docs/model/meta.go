package model

import "github.com/bsn-si/IPEHR-gateway/src/pkg/indexer/ehrIndexer"

type DocumentMeta ehrIndexer.DocsDocumentMeta

func (dm *DocumentMeta) GetAttr(code Attribute) []byte {
	for _, attr := range dm.Attrs {
		if attr.Code == code {
			return attr.Value
		}
	}

	return nil
}
