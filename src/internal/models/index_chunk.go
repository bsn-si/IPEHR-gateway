package models

import (
	"crypto/sha256"
	"fmt"
	"time"

	"github.com/google/uuid"
)

type IndexChunk struct {
	Key       string    `db:"key"`
	CreatedAt time.Time `db:"created_at"`
	GroupID   string    `db:"group_id"`
	DataID    string    `db:"data_id"`
	EhrID     string    `db:"ehr_id"`
	Data      []byte    `db:"data"`
	Hash      string    `db:"hash"`
}

func NewIndexChunk(groupID, dataID, ehrID string, data []byte) IndexChunk {
	hash := sha256.New()
	hash.Write([]byte(fmt.Sprintf("%s.%s.%s", groupID, dataID, ehrID)))
	hash.Write(data)

	idxChunck := IndexChunk{
		Key:     uuid.NewString(),
		GroupID: groupID,
		DataID:  dataID,
		EhrID:   ehrID,
		Data:    data,
		Hash:    string(hash.Sum(nil)),
	}

	return idxChunck
}

func (idxChunk IndexChunk) Validate() bool {
	hash := sha256.New()
	hash.Write([]byte(fmt.Sprintf("%s.%s.%s", idxChunk.GroupID, idxChunk.DataID, idxChunk.EhrID)))
	hash.Write(idxChunk.Data)

	return idxChunk.Hash == string(hash.Sum(nil))
}
