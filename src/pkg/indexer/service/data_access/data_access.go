// Package data_access keys index
package data_access

import (
	"encoding/hex"

	"github.com/google/uuid"
	"golang.org/x/crypto/sha3"

	"hms/gateway/pkg/crypto/keybox"
	"hms/gateway/pkg/indexer"
	"hms/gateway/pkg/keystore"
)

type DataAccessIndex struct {
	index    indexer.Indexer
	keystore *keystore.KeyStore
}

func New() *DataAccessIndex {
	return &DataAccessIndex{
		index:    indexer.Init("data_access"),
		keystore: keystore.New(),
	}
}

// Add DataAccessIndex key
func (d *DataAccessIndex) Add(userId, accessGroupId string, accessGroupKey []byte) error {
	userUUID, err := uuid.Parse(userId)
	if err != nil {
		return err
	}

	accessGroupUUID, err := uuid.Parse(accessGroupId)
	if err != nil {
		return err
	}

	// Getting user publicKey
	userPubKey, _, err := d.keystore.Get(userId)
	if err != nil {
		return err
	}

	// access group key encryption
	keyEncrypted, err := keybox.SealAnonymous(accessGroupKey, userPubKey)
	if err != nil {
		return err
	}

	// Index doc_id -> encrypted_doc_key
	indexKey := sha3.Sum256(append(userUUID[:], accessGroupUUID[:]...))
	indexKeyStr := hex.EncodeToString(indexKey[:])
	if err = d.index.Add(indexKeyStr, keyEncrypted); err != nil {
		return err
	}

	return nil
}

// Get DataAccessIndex key
func (d *DataAccessIndex) Get(userId, accessGroupId string) (groupAccessKeyEncrypted []byte, err error) {
	userUUID, err := uuid.Parse(userId)
	if err != nil {
		return
	}

	accessGroupUUID, err := uuid.Parse(accessGroupId)
	if err != nil {
		return
	}

	indexKey := sha3.Sum256(append(userUUID[:], accessGroupUUID[:]...))
	indexKeyStr := hex.EncodeToString(indexKey[:])
	err = d.index.GetById(indexKeyStr, &groupAccessKeyEncrypted)

	return groupAccessKeyEncrypted, err
}
