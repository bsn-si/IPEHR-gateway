package model

import (
	"hms/gateway/pkg/docs/status"
	"hms/gateway/pkg/docs/types"
)

type DocumentMeta struct {
	TypeCode       types.DocumentType
	StorageID      *[32]byte
	DocIDEncrypted []byte
	Timestamp      uint64
	Status         status.DocumentStatus
}
