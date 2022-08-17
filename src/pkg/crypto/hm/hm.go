package hm

import (
	"encoding/binary"
	"fmt"

	"golang.org/x/crypto/chacha20poly1305"
	"golang.org/x/crypto/sha3"
)

var (
	ErrIncorrectKey   = fmt.Errorf("Key is incorrect")
	ErrIncorrectNonce = fmt.Errorf("Nonce is incorrect")
)

// key must not be nil
// key[0:4] != []byte{0,0,0,0}
func EncryptInt(x int64, key *[32]byte) int64 {
	if key == nil {
		panic(ErrIncorrectKey)
	}

	a := int64(binary.BigEndian.Uint32(key[0:4]))
	b := int64(binary.BigEndian.Uint32(key[4:8]))

	if a == 0 {
		panic(ErrIncorrectKey)
	}

	return x*a + b
}

// key must not be nil
// key[0:4] != []byte{0,0,0,0}
func EncryptFloat(x float64, key *[32]byte) float64 {
	if key == nil {
		panic(ErrIncorrectKey)
	}

	a := float64(binary.BigEndian.Uint32(key[0:4]))
	b := float64(binary.BigEndian.Uint32(key[4:8]))

	if a == 0 {
		panic(ErrIncorrectKey)
	}

	return x*a + b
}

// panic on empty key or nonce
func EncryptString(in string, key *[32]byte, nonce *[12]byte) (out []byte) {
	if key == nil {
		panic(ErrIncorrectKey)
	}

	if nonce == nil {
		panic(ErrIncorrectNonce)
	}

	aead, err := chacha20poly1305.New(key[:])
	if err != nil {
		panic(ErrIncorrectKey)
	}

	msg := sha3.Sum256([]byte(in))
	encrypted := aead.Seal(nonce[:], nonce[:], msg[:], nil)

	return encrypted[12:]
}
