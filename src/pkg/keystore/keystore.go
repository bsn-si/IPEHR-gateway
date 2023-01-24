// Package keystore Storage for user public/private key pair
package keystore

import (
	cryptoRand "crypto/rand"
	"encoding/hex"
	"fmt"
	"log"

	"golang.org/x/crypto/nacl/box"
	"golang.org/x/crypto/sha3"

	"github.com/bsn-si/IPEHR-gateway/src/pkg/crypto/chachaPoly"
	"github.com/bsn-si/IPEHR-gateway/src/pkg/errors"
	"github.com/bsn-si/IPEHR-gateway/src/pkg/storage"
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
			log.Println("Generete new keys for userID", userID)

			publicKey, privateKey, err = k.generateAndStoreKeys(userID)
			if err != nil {
				return nil, nil, fmt.Errorf("generateAndStoreKeys error: %w", err)
			}

			return publicKey, privateKey, nil
		}
		return nil, nil, fmt.Errorf("storage.Get error: %w", err)
	}

	keysDecrypted, err := k.decryptUserKeys(keysEncrypted)
	if err != nil {
		return nil, nil, fmt.Errorf("decryptUserKeys error: %w", err)
	}

	publicKey = new([32]byte)
	privateKey = new([32]byte)

	copy(publicKey[:], keysDecrypted[0:32])
	copy(privateKey[:], keysDecrypted[32:64])

	return
}

// Generate and store new user key pair
func (k *KeyStore) generateAndStoreKeys(userID string) (*[32]byte, *[32]byte, error) {
	publicKey, privateKey, err := k.generateKeys()
	if err != nil {
		return nil, nil, fmt.Errorf("generateKeys error: %w", err)
	}

	err = k.storeKeys(userID, publicKey, privateKey)
	if err != nil {
		return nil, nil, fmt.Errorf("storeKeys error: %w", err)
	}

	return publicKey, privateKey, nil
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
		return fmt.Errorf("encryptUserKeys error: %w", err)
	}

	if err = k.storage.AddWithID(storeID, keysEncrypted); err != nil {
		return fmt.Errorf("storage.AddWithID error: %w", err)
	}

	return nil
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
