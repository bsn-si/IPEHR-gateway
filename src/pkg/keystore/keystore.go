// Package keystore Storage for user public/private key pair
package keystore

import (
	cryptoRand "crypto/rand"
	"hms/gateway/pkg/errors"
	"hms/gateway/pkg/storage"

	"golang.org/x/crypto/nacl/box"
	"golang.org/x/crypto/sha3"
)

type KeyStore struct {
	storage storage.Storager
}

func New() *KeyStore {
	return &KeyStore{
		storage: *storage.Storage(),
	}
}

// Get user key pair
func (k *KeyStore) Get(userId string) (publicKey, privateKey *[32]byte, err error) {
	storeId := k.storeId(userId)

	keys, err := k.storage.Get(storeId)
	if err != nil {
		if errors.Is(err, errors.IsNotExist) {
			publicKey, privateKey, err = k.generateAndStoreKeys(userId)
		}
		return
	}

	publicKey = new([32]byte)
	privateKey = new([32]byte)

	copy(publicKey[:], keys[0:32])
	copy(privateKey[:], keys[32:64])

	return
}

// Generate and store new user key pair
func (k *KeyStore) generateAndStoreKeys(userId string) (publicKey, privateKey *[32]byte, err error) {
	publicKey, privateKey, err = k.generateKeys()
	if err != nil {
		return
	}
	err = k.storeKeys(userId, publicKey, privateKey)
	return
}

// Generate new user key pair
func (k *KeyStore) generateKeys() (publicKey, privateKey *[32]byte, err error) {
	publicKey, privateKey, err = box.GenerateKey(cryptoRand.Reader)
	return
}

// Store user key pair
func (k *KeyStore) storeKeys(userId string, publicKey, privateKey *[32]byte) error {
	storeId := k.storeId(userId)
	keys := append(publicKey[:], privateKey[:]...)
	return k.storage.AddWithId(storeId, keys)
}

// Get store file ID where the user keys is
func (k *KeyStore) storeId(userId string) *[32]byte {
	id := sha3.Sum256([]byte(userId + "keys"))
	return &id
}
