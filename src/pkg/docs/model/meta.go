package model

import (
	"hms/gateway/pkg/docs/status"
	"hms/gateway/pkg/docs/types"
)

type DocumentMeta struct {
	TypeCode        types.DocumentType
	CID             *[32]byte
	DocUIDEncrypted []byte
	DocBaseUIDHash  *[32]byte
	IsLastVersion   bool
	Version         string
	Status          status.DocumentStatus
	Timestamp       uint64
}
