// Package group_access Stores access group data
// userID + GroupID -> model.GroupAccess
package groupAccess

import (
	"encoding/hex"

	"github.com/google/uuid"
	"github.com/vmihailenco/msgpack/v5"
	"golang.org/x/crypto/sha3"

	"hms/gateway/pkg/crypto/keybox"
	"hms/gateway/pkg/docs/model"
	"hms/gateway/pkg/indexer"
	"hms/gateway/pkg/keystore"
)

type Index struct {
	index    indexer.Indexer
	keystore *keystore.KeyStore
}

func New(ks *keystore.KeyStore) *Index {
	return &Index{
		index:    indexer.Init("group_access"),
		keystore: ks,
	}
}

func (i *Index) Add(userID string, groupAccess *model.GroupAccess) (err error) {
	userPubKey, _, err := i.keystore.Get(userID)
	if err != nil {
		return
	}

	groupAccessByte, err := msgpack.Marshal(groupAccess)
	if err != nil {
		return
	}

	groupAccessEncrypted, err := keybox.SealAnonymous(groupAccessByte, userPubKey)
	if err != nil {
		return
	}

	h := sha3.Sum256(append([]byte(userID), groupAccess.GroupUUID[:]...))
	indexKey := hex.EncodeToString(h[:])

	err = i.index.Add(indexKey, groupAccessEncrypted)

	return
}

func (i *Index) Get(userID string, groupAccessUUID *uuid.UUID) (groupAccess *model.GroupAccess, err error) {
	userPubKey, userPrivateKey, err := i.keystore.Get(userID)
	if err != nil {
		return
	}

	h := sha3.Sum256(append([]byte(userID), groupAccessUUID[:]...))
	indexKey := hex.EncodeToString(h[:])

	var groupAccessEncrypted []byte

	err = i.index.GetByID(indexKey, &groupAccessEncrypted)
	if err != nil {
		return
	}

	groupAccessByte, err := keybox.OpenAnonymous(groupAccessEncrypted, userPubKey, userPrivateKey)
	if err != nil {
		return
	}

	err = msgpack.Unmarshal(groupAccessByte, &groupAccess)

	return
}
