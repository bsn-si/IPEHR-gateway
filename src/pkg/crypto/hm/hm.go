package hm

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"strconv"

	"golang.org/x/crypto/chacha20poly1305"
	"golang.org/x/crypto/sha3"
)

var (
	ErrIncorrectKey   = fmt.Errorf("Key is incorrect")
	ErrIncorrectNonce = fmt.Errorf("Nonce is incorrect")
)

func EncryptInt(in interface{}, key *[32]byte, errors []error) (out int64) {
	if key == nil {
		errors = append(errors, ErrIncorrectKey)
		return
	}

	a := int64(binary.BigEndian.Uint32(key[0:4]))
	b := int64(binary.BigEndian.Uint32(key[4:8]))

	if a == 0 || b == 0 {
		errors = append(errors, ErrIncorrectKey)
		return
	}

	var (
		num int64
		err error
	)
	switch in.(type) {
	case string:
		num, err = strconv.ParseInt(in.(string), 10, 64)
		if err != nil {
			errors = append(errors, err)
			return
		}
	case int:
		num = int64(in.(int))
	case int64:
		num = int64(in.(int64))
	case uint64:
		num = int64(in.(uint64))
	default:
		errors = append(errors, fmt.Errorf("incorrect type of input data %T", in))
		return 0
	}

	return num*a + b
}

func EncryptFloat(in interface{}, key *[32]byte, errors []error) (out float64) {
	if key == nil {
		errors = append(errors, ErrIncorrectKey)
		return
	}

	a := float64(binary.BigEndian.Uint32(key[0:4]))
	b := float64(binary.BigEndian.Uint32(key[4:8]))

	if a == 0 || b == 0 {
		errors = append(errors, ErrIncorrectKey)
		return
	}

	var (
		num float64
		err error
	)
	switch in.(type) {
	case string:
		num, err = strconv.ParseFloat(in.(string), 64)
		if err != nil {
			errors = append(errors, err)
			return
		}
	case float64:
		num = in.(float64)
	default:
		errors = append(errors, fmt.Errorf("incorrect type of input data %T", in))
		return 0
	}

	return num*a + b
}

func EncryptString(in string, key *[32]byte, nonce *[12]byte, errors []error) (out []byte) {
	if key == nil || bytes.Equal(key[:], make([]byte, 32)) {
		errors = append(errors, ErrIncorrectKey)
		return
	}

	if nonce == nil || bytes.Equal(key[:], make([]byte, 12)) {
		errors = append(errors, ErrIncorrectNonce)
		return
	}

	aead, err := chacha20poly1305.New(key[:])
	if err != nil {
		errors = append(errors, err)
		return
	}

	msg := sha3.Sum256([]byte(in))
	encrypted := aead.Seal(nonce[:], nonce[:], msg[:], nil)

	return encrypted[12:]
}
