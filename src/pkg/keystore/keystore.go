// Package keystore Storage for user public/private key pair
package keystore

import (
	cryptoRand "crypto/rand"
	"encoding/hex"

	"golang.org/x/crypto/nacl/box"
	"golang.org/x/crypto/sha3"

	"hms/gateway/pkg/crypto/chachaPoly"
	"hms/gateway/pkg/errors"
	"hms/gateway/pkg/storage"
)

type KeyStore struct {
	storage     storage.Storager
	keystoreKey []byte
}

func New(key string) *KeyStore {
	if key == "" {
		panic("Keystore key is empty. Check the config.")
	}

	keyByte, err := hex.DecodeString(key)
	if err != nil {
		return nil
	}

	return &KeyStore{
		storage:     storage.Storage(),
		keystoreKey: keyByte,
	}
}

// Get user key pair
func (k *KeyStore) Get(userID string) (publicKey, privateKey *[32]byte, err error) {
	storeID := k.storeID(userID)

	keysEncrypted, err := k.storage.Get(storeID)
	if err != nil {
		if errors.Is(err, errors.ErrIsNotExist) {
			publicKey, privateKey, err = k.generateAndStoreKeys(userID)
		}

		return
	}

	keysDecrypted, err := k.decryptUserKeys(keysEncrypted)
	if err != nil {
		return
	}

	publicKey = new([32]byte)
	privateKey = new([32]byte)

	copy(publicKey[:], keysDecrypted[0:32])
	copy(privateKey[:], keysDecrypted[32:64])

	return
}

// Generate and store new user key pair
func (k *KeyStore) generateAndStoreKeys(userID string) (publicKey, privateKey *[32]byte, err error) {
	publicKey, privateKey, err = k.generateKeys()
	if err != nil {
		return
	}

	err = k.storeKeys(userID, publicKey, privateKey)

	return
}

// Generate new user key pair
func (k *KeyStore) generateKeys() (publicKey, privateKey *[32]byte, err error) {
	publicKey, privateKey, err = box.GenerateKey(cryptoRand.Reader)
	return
}

// Store user key pair
func (k *KeyStore) storeKeys(userID string, publicKey, privateKey *[32]byte) error {
	storeID := k.storeID(userID)

	keysDecrypted := append(publicKey[:], privateKey[:]...)

	keysEncrypted, err := k.encryptUserKeys(keysDecrypted)
	if err != nil {
		return err
	}

	return k.storage.AddWithID(storeID, keysEncrypted)
}

// Get store file ID where the user keys is
func (k *KeyStore) storeID(userID string) *[32]byte {
	id := sha3.Sum256([]byte(userID + "keys"))
	return &id
}

func (k *KeyStore) encryptUserKeys(keysDecrypted []byte) (keysEncrypted []byte, err error) {
	key, err := chachaPoly.NewKeyFromBytes(k.keystoreKey)
	if err != nil {
		return
	}

	keysEncrypted, err = key.Encrypt(keysDecrypted)

	return
}

func (k *KeyStore) decryptUserKeys(keysEncrypted []byte) (keysDecrypted []byte, err error) {
	key, err := chachaPoly.NewKeyFromBytes(k.keystoreKey)
	if err != nil {
		return
	}

	keysDecrypted, err = key.Decrypt(keysEncrypted)

	return
}
