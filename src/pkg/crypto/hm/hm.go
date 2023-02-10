package hm

import (
	"encoding/binary"
	"errors"
	"math/big"

	"golang.org/x/crypto/chacha20poly1305"
)

type (
	Key   [32]byte
	Nonce [12]byte
)

var (
	ErrIncorrectKey   = errors.New("Key is incorrect")
	ErrIncorrectNonce = errors.New("Nonce is incorrect")
	ErrOverflow       = errors.New("Calculation overflow")
)

// Key must not be nil
// key[0:4] != []byte{0,0,0,0}
// Limitations: -2147483648 <= x <= 2147483647 with key 0xffffffffffffffff...
// Returns x*a + b
func EncryptInt64(x int64, key *Key) (int64, error) {
	if key == nil {
		return 0, ErrIncorrectKey
	}

	a := big.NewInt(int64(binary.BigEndian.Uint32(key[0:4])))
	b := big.NewInt(int64(binary.BigEndian.Uint32(key[4:8])))

	if a.Int64() == 0 {
		return 0, ErrIncorrectKey
	}

	xBig := a.Mul(a, big.NewInt(x))
	xBig.Add(xBig, b)

	if !xBig.IsInt64() {
		return 0, ErrOverflow
	}

	return xBig.Int64(), nil
}

// Key must not be nil
// Key[0:4] != []byte{0,0,0,0}
// Returs (x - b)/a
func DecryptInt64(x int64, key *[32]byte) (int64, error) {
	if key == nil {
		return 0, ErrIncorrectKey
	}

	a := int64(binary.BigEndian.Uint32(key[0:4]))
	b := int64(binary.BigEndian.Uint32(key[4:8]))

	if a == 0 {
		return 0, ErrIncorrectKey
	}

	return (x - b) / a, nil
}

// Key must not be nil
// key[0:4] != []byte{0,0,0,0}
// Returns x*a + b
func EncryptFloat64(x float64, key *Key) (float64, error) {
	if key == nil {
		return 0, ErrIncorrectKey
	}

	if binary.BigEndian.Uint32(key[0:4]) == 0 {
		return 0, ErrIncorrectKey
	}

	a := big.NewFloat(float64(binary.BigEndian.Uint32(key[0:4])))
	b := big.NewFloat(float64(binary.BigEndian.Uint32(key[4:8])))

	xBig := a.Mul(a, big.NewFloat(x))
	xBig.Add(xBig, b)

	x, acc := xBig.Float64()
	if acc != 0 {
		return 0, ErrOverflow
	}

	return x, nil
}

// key must not be nil
// key[0:4] != []byte{0,0,0,0}
// returs (x - b)/a
func DecryptFloat64(x float64, key *Key) (float64, error) {
	if key == nil {
		return 0, ErrIncorrectKey
	}

	a := float64(binary.BigEndian.Uint32(key[0:4]))
	b := float64(binary.BigEndian.Uint32(key[4:8]))

	if a == 0 {
		return 0, ErrIncorrectKey
	}

	return (x - b) / a, nil
}

// Key or nonce must not be nil
// key[0:4] != []byte{0,0,0,0}
func EncryptString(in string, key *Key, nonce *Nonce) (out []byte) {
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

	encrypted := aead.Seal(nonce[:], nonce[:], []byte(in), nil)

	return encrypted[12:]
}
