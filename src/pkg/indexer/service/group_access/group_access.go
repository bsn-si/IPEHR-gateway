// Package group_access Stores access group data
// userId + GroupId -> model.GroupAccess
package group_access

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

type GroupAccessIndex struct {
	index    indexer.Indexer
	keystore *keystore.KeyStore
}

func New(ks *keystore.KeyStore) *GroupAccessIndex {
	return &GroupAccessIndex{
		index:    indexer.Init("group_access"),
		keystore: ks,
	}
}

func (g *GroupAccessIndex) Add(userId string, groupAccess *model.GroupAccess) (err error) {
	userPubKey, _, err := g.keystore.Get(userId)
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

	h := sha3.Sum256(append([]byte(userId), groupAccess.GroupUUID[:]...))
	indexKey := hex.EncodeToString(h[:])

	err = g.index.Add(indexKey, groupAccessEncrypted)

	return
}

func (g *GroupAccessIndex) Get(userId string, groupAccessUUID *uuid.UUID) (groupAccess *model.GroupAccess, err error) {
	userPubKey, userPrivateKey, err := g.keystore.Get(userId)
	if err != nil {
		return
	}

	h := sha3.Sum256(append([]byte(userId), groupAccessUUID[:]...))
	indexKey := hex.EncodeToString(h[:])

	var groupAccessEncrypted []byte
	err = g.index.GetById(indexKey, &groupAccessEncrypted)
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
