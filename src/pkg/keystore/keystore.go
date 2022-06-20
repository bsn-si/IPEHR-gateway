// Package keystore Storage for user public/private key pair
package keystore

import (
	cryptoRand "crypto/rand"
	"hms/gateway/pkg/config"
	"hms/gateway/pkg/crypto/chacha_poly"
	"hms/gateway/pkg/errors"
	"hms/gateway/pkg/storage"

	"golang.org/x/crypto/nacl/box"
	"golang.org/x/crypto/sha3"
)

type KeyStore struct {
	storage storage.Storager
	cfg     *config.Config
}

func New() *KeyStore {
	cfg, err := config.New()
	if err != nil {
		return nil
	}

	return &KeyStore{
		storage: storage.Init(),
		cfg:     cfg,
	}
}

// Get user key pair
func (k *KeyStore) Get(userId string) (publicKey, privateKey *[32]byte, err error) {
	storeId := k.storeId(userId)

	keysEncrypted, err := k.storage.Get(storeId)
	if err != nil {
		if errors.Is(err, errors.IsNotExist) {
			publicKey, privateKey, err = k.generateAndStoreKeys(userId)
		}
		return
	}

	keysDecrypted, err := k.decryptUserKeys(&keysEncrypted)
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
	keysDecrypted := append(publicKey[:], privateKey[:]...)

	keysEncrypted, err := k.encryptUserKeys(&keysDecrypted)
	if err != nil {
		return err
	}

	return k.storage.AddWithId(storeId, keysEncrypted)
}

// Get store file ID where the user keys is
func (k *KeyStore) storeId(userId string) *[32]byte {
	id := sha3.Sum256([]byte(userId + "keys"))
	return &id
}

func (k *KeyStore) encryptUserKeys(keysDecrypted *[]byte) (keysEncrypted []byte, err error) {
	key, err := chacha_poly.NewKeyFromBytes(k.cfg.KeystoreKey)
	if err != nil {
		return
	}

	keysEncrypted, err = key.Encrypt(*keysDecrypted)

	return
}

func (k *KeyStore) decryptUserKeys(keysEncrypted *[]byte) (keysDecrypted []byte, err error) {
	key, err := chacha_poly.NewKeyFromBytes(k.cfg.KeystoreKey)
	if err != nil {
		return
	}

	keysDecrypted, err = key.Decrypt(*keysEncrypted)
	return
}
