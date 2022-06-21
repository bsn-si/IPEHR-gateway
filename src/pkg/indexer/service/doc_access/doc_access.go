// Package doc_access User document keys index
package doc_access

import (
	"encoding/hex"

	"github.com/google/uuid"
	"golang.org/x/crypto/sha3"

	"hms/gateway/pkg/crypto/chacha_poly"
	"hms/gateway/pkg/crypto/keybox"
	"hms/gateway/pkg/indexer"
	"hms/gateway/pkg/keystore"
)

type DocAccessIndex struct {
	index    indexer.Indexer
	keystore *keystore.KeyStore
}

func New(ks *keystore.KeyStore) *DocAccessIndex {
	return &DocAccessIndex{
		index:    indexer.Init("doc_access"),
		keystore: ks,
	}
}

// Add user's document key
func (u *DocAccessIndex) Add(userId string, docStorageId *[32]byte, docKey []byte) error {
	userUUID, err := uuid.Parse(userId)
	if err != nil {
		return err
	}

	// Getting user publicKey
	userPubKey, _, err := u.keystore.Get(userId)
	if err != nil {
		return err
	}

	// Document key encryption
	keyEncrypted, err := keybox.SealAnonymous(docKey, userPubKey)
	if err != nil {
		return err
	}

	// Index doc_id -> encrypted_doc_key
	indexKey := sha3.Sum256(append(docStorageId[:], userUUID[:]...))
	indexKeyStr := hex.EncodeToString(indexKey[:])

	if err = u.index.Add(indexKeyStr, keyEncrypted); err != nil {
		return err
	}

	return nil
}

// Get user key
func (u *DocAccessIndex) Get(userId string) ([]byte, error) {
	var keyEncrypted []byte
	err := u.index.GetById(userId, &keyEncrypted)
	return keyEncrypted, err
}

func (u *DocAccessIndex) GetDocumentKey(userId string, docStorageId *[32]byte) (docKey *chacha_poly.Key, err error) {
	userUUID, err := uuid.Parse(userId)
	if err != nil {
		return
	}

	userPubKey, userPrivateKey, err := u.keystore.Get(userId)
	if err != nil {
		return
	}

	indexKey := sha3.Sum256(append(docStorageId[:], userUUID[:]...))
	indexKeyStr := hex.EncodeToString(indexKey[:])

	var keyEncrypted []byte
	err = u.index.GetById(indexKeyStr, &keyEncrypted)
	if err != nil {
		return
	}

	docKeyBytes, err := keybox.OpenAnonymous(keyEncrypted, userPubKey, userPrivateKey)
	if err != nil {
		return
	}

	docKey, err = chacha_poly.NewKeyFromBytes(docKeyBytes)
	if err != nil {
		return
	}

	return
}
