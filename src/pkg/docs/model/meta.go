package model

import (
	"hms/gateway/pkg/docs/status"
	"hms/gateway/pkg/docs/types"

	"github.com/ipfs/go-cid"
)

type DocumentMeta struct {
	TypeCode        types.DocumentType
	CID             []byte
	DealCID         []byte
	MinerAddress    []byte
	DocUIDEncrypted []byte
	DocBaseUIDHash  *[32]byte
	IsLastVersion   bool
	Version         string
	Status          status.DocumentStatus
	Timestamp       uint64
}

func (m *DocumentMeta) Cid() *cid.Cid {
	CID, err := cid.Parse(m.CID)
	if err != nil {
		panic(err)
	}

	return &CID
}
