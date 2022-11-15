package model

import (
	"hms/gateway/pkg/indexer/ehrIndexer"

	"github.com/ipfs/go-cid"
)

type DocumentMeta ehrIndexer.DocsDocumentMeta

func (m *DocumentMeta) Cid() *cid.Cid {
	CID, err := cid.Parse(m.CID)
	if err != nil {
		panic(err)
	}

	return &CID
}
