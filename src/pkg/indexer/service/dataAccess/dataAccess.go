// Package dataAccess keys index
package dataAccess

import (
	"encoding/hex"

	"github.com/google/uuid"
	"golang.org/x/crypto/sha3"

	"hms/gateway/pkg/crypto/keybox"
	"hms/gateway/pkg/indexer"
	"hms/gateway/pkg/keystore"
)

type Index struct {
	index    indexer.Indexer
	keystore *keystore.KeyStore
}

func New(ks *keystore.KeyStore) *Index {
	return &Index{
		index:    indexer.Init("data_access"),
		keystore: ks,
	}
}

// Add Index key
func (d *Index) Add(userID, accessGroupID string, accessGroupKey []byte) error {
	userUUID, err := uuid.Parse(userID)
	if err != nil {
		return err
	}

	accessGroupUUID, err := uuid.Parse(accessGroupID)
	if err != nil {
		return err
	}

	// Getting user publicKey
	userPubKey, _, err := d.keystore.Get(userID)
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

// Get Index key
func (d *Index) Get(userID, accessGroupID string) (groupAccessKeyEncrypted []byte, err error) {
	userUUID, err := uuid.Parse(userID)
	if err != nil {
		return
	}

	accessGroupUUID, err := uuid.Parse(accessGroupID)
	if err != nil {
		return
	}

	indexKey := sha3.Sum256(append(userUUID[:], accessGroupUUID[:]...))
	indexKeyStr := hex.EncodeToString(indexKey[:])
	err = d.index.GetByID(indexKeyStr, &groupAccessKeyEncrypted)

	return groupAccessKeyEncrypted, err
}
