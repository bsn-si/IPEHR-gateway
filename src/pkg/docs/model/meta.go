package model

import (
	"hms/gateway/pkg/docs/status"
	"hms/gateway/pkg/docs/types"
)

type DocumentMeta struct {
	TypeCode       types.DocumentType
	StorageId      *[32]byte
	DocIdEncrypted []byte
	Timestamp      uint64
	Status         status.DocumentStatus
}
