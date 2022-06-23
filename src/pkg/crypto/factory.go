package crypto

import (
	"hms/gateway/pkg/crypto/chacha_poly"
	"hms/gateway/pkg/crypto/common"
	"hms/gateway/pkg/crypto/compression"
)

/*
func GenerateKey() KeyInterface {
	return chacha_poly.GenerateKey()
}

func NewKeyFromBytes(keyBytes []byte) (KeyInterface, error) {
	return chacha_poly.NewKeyFromBytes(keyBytes)
}
*/

func GenerateKey() common.KeyInterface {
	return compression.New(chacha_poly.GenerateKey())
}

func NewKeyFromBytes(keyBytes []byte) (common.KeyInterface, error) {
	key, err := chacha_poly.NewKeyFromBytes(keyBytes)
	if err != nil {
		return nil, err
	}
	return compression.NewKeyFromBytes(key), nil
}
