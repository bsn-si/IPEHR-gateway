package crypto

import (
	"hms/gateway/pkg/config"
	"hms/gateway/pkg/crypto/chacha_poly"
	"hms/gateway/pkg/crypto/common"
	"hms/gateway/pkg/crypto/compression"
)

var compressionEnabled bool

func init() {
	globalConfig, err := config.New()
	if err == nil {
		compressionEnabled = globalConfig.CompressionEnabled
	} else {
		compressionEnabled = false
	}
}

func GenerateKey() common.KeyInterface {
	if compressionEnabled {
		return compression.New(chacha_poly.GenerateKey())
	}
	return chacha_poly.GenerateKey()
}

func NewKeyFromBytes(keyBytes []byte) (common.KeyInterface, error) {
	implementation, err := chacha_poly.NewKeyFromBytes(keyBytes)
	if err != nil {
		return nil, err
	}

	if compressionEnabled {
		return compression.NewKeyFromBytes(implementation), nil
	} else {
		return implementation, nil
	}
}
