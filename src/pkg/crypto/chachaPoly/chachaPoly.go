package chachaPoly

import (
	crypto_rand "crypto/rand"
	"encoding/hex"
	"fmt"

	"golang.org/x/crypto/chacha20poly1305"

	"github.com/bsn-si/IPEHR-gateway/src/pkg/errors"
)

const (
	// KeyLength is the size of the key used by this AEAD, in bytes.
	KeyLength = 32

	// NonceLength is the size of the nonce used with the standard variant of this
	// AEAD, in bytes.
	//
	// Note that this is too short to be safely generated at random if the same
	// key is reused more than 2³² times.
	NonceLength = 12

	// Overhead is the size of the Poly1305 authentication tag, and the
	// difference between a ciphertext length and its plaintext.
	Overhead = 16
)

type Key [KeyLength]byte

func GenerateKey() *Key {
	key := new(Key)
	if _, err := crypto_rand.Read(key[:]); err != nil {
		panic(err)
	}

	return key
}

func NewKeyFromBytes(keyBytes []byte) (*Key, error) {
	if len(keyBytes) != KeyLength {
		return nil, fmt.Errorf("%w: Key length is incorrect", errors.ErrEncryption)
	}

	key := new(Key)

	copy(key[:], keyBytes)

	return key, nil
}

func (k Key) Encrypt(msg []byte) ([]byte, error) {
	if len(k) != KeyLength {
		return nil, fmt.Errorf("%w: Key length is incorrect", errors.ErrEncryption)
	}

	aead, err := chacha20poly1305.New(k[:])
	if err != nil {
		return nil, fmt.Errorf("key init error: %w", err)
	}

	nonce := make([]byte, NonceLength, NonceLength+len(msg)+Overhead)
	if _, err := crypto_rand.Read(nonce); err != nil {
		return nil, fmt.Errorf("nonce creating error: %w", err)
	}

	encrypted := aead.Seal(nonce, nonce, msg, nil)

	return encrypted, nil
}

func (k Key) EncryptWithAuthData(msg, authData []byte) ([]byte, error) {
	if len(k) != KeyLength {
		return nil, fmt.Errorf("%w: Key length is incorrect", errors.ErrEncryption)
	}

	aead, err := chacha20poly1305.New(k[:])
	if err != nil {
		return nil, fmt.Errorf("key init error: %w", err)
	}

	nonce := make([]byte, NonceLength, NonceLength+len(msg)+Overhead)
	if _, err := crypto_rand.Read(nonce); err != nil {
		return nil, fmt.Errorf("nonce creating error: %w", err)
	}

	encrypted := aead.Seal(nonce, nonce, msg, authData)

	return encrypted, nil
}

func (k Key) Decrypt(encrypted []byte) ([]byte, error) {
	if len(k) != KeyLength {
		return nil, fmt.Errorf("%w: Key length is incorrect", errors.ErrEncryption)
	}

	aead, err := chacha20poly1305.New(k[:])
	if err != nil {
		return nil, fmt.Errorf("key init error: %w", err)
	}

	if len(encrypted) < NonceLength {
		return nil, fmt.Errorf("%w: Ciphertext too short", errors.ErrEncryption)
	}

	// Split nonce and ciphertext.
	nonce, ciphertext := encrypted[:NonceLength], encrypted[NonceLength:]

	// Decrypt the message and check it wasn't tampered with.
	msg, err := aead.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return nil, fmt.Errorf("ciphertext open error: %w", err)
	}

	return msg, nil
}

func (k Key) DecryptWithAuthData(encrypted, authData []byte) ([]byte, error) {
	if len(k) != KeyLength {
		return nil, fmt.Errorf("%w: Key length is incorrect", errors.ErrEncryption)
	}

	aead, err := chacha20poly1305.New(k[:])
	if err != nil {
		return nil, fmt.Errorf("key init error: %w", err)
	}

	if len(encrypted) < NonceLength {
		return nil, fmt.Errorf("%w: Ciphertext too short", errors.ErrEncryption)
	}

	// Split nonce and ciphertext.
	nonce, ciphertext := encrypted[:NonceLength], encrypted[NonceLength:]

	// Decrypt the message and check it wasn't tampered with.
	msg, err := aead.Open(nil, nonce, ciphertext, authData)
	if err != nil {
		return nil, fmt.Errorf("ciphertext open error: %w", err)
	}

	return msg, nil
}

func (k Key) String() string {
	return hex.EncodeToString(k[:])
}

func (k Key) Bytes() []byte {
	return k[:]
}
