// Package doc_access User document keys index
package docAccess

import (
	"encoding/hex"

	"github.com/google/uuid"
	"golang.org/x/crypto/sha3"

	"hms/gateway/pkg/crypto/chachaPoly"
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
		index:    indexer.Init("doc_access"),
		keystore: ks,
	}
}

// Add user's document key
func (u *Index) Add(userID string, docStorageID *[32]byte, docKey []byte) error {
	userUUID, err := uuid.Parse(userID)
	if err != nil {
		return err
	}

	// Getting user publicKey
	userPubKey, _, err := u.keystore.Get(userID)
	if err != nil {
		return err
	}

	// Document key encryption
	keyEncrypted, err := keybox.SealAnonymous(docKey, userPubKey)
	if err != nil {
		return err
	}

	// Index doc_id -> encrypted_doc_key
	indexKey := sha3.Sum256(append(docStorageID[:], userUUID[:]...))
	indexKeyStr := hex.EncodeToString(indexKey[:])

	if err = u.index.Add(indexKeyStr, keyEncrypted); err != nil {
		return err
	}

	return nil
}

// Get user key
func (u *Index) Get(userID string) ([]byte, error) {
	var keyEncrypted []byte

	err := u.index.GetByID(userID, &keyEncrypted)

	return keyEncrypted, err
}

func (u *Index) GetDocumentKey(userID string, docStorageID *[32]byte) (docKey *chachaPoly.Key, err error) {
	userUUID, err := uuid.Parse(userID)
	if err != nil {
		return
	}

	userPubKey, userPrivateKey, err := u.keystore.Get(userID)
	if err != nil {
		return
	}

	indexKey := sha3.Sum256(append(docStorageID[:], userUUID[:]...))
	indexKeyStr := hex.EncodeToString(indexKey[:])

	var keyEncrypted []byte

	err = u.index.GetByID(indexKeyStr, &keyEncrypted)
	if err != nil {
		return
	}

	docKeyBytes, err := keybox.OpenAnonymous(keyEncrypted, userPubKey, userPrivateKey)
	if err != nil {
		return
	}

	docKey, err = chachaPoly.NewKeyFromBytes(docKeyBytes)
	if err != nil {
		return
	}

	return docKey, nil
}
