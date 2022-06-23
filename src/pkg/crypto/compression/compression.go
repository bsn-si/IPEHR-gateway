package compression

import (
	"hms/gateway/pkg/compressor"
	"hms/gateway/pkg/config"
	"hms/gateway/pkg/crypto/common"
)

type Key struct {
	implementation common.KeyInterface
	compressor     *compressor.Compressor
}

func New(implementation common.KeyInterface) common.KeyInterface {
	globalConfig, err := config.New()
	if err != nil {
		return nil
	}
	return &Key{
		implementation: implementation,
		compressor:     compressor.New(globalConfig.CompressionLevel),
	}
}

func NewKeyFromBytes(implementation common.KeyInterface) common.KeyInterface {
	return &Key{
		implementation: implementation,
	}
}

func (k Key) Encrypt(msg []byte) ([]byte, error) {
	msgCompressed, err := k.compressor.Compress(&msg)
	if err != nil {
		return nil, err
	}
	return k.implementation.Encrypt(*msgCompressed)
}

func (k Key) Decrypt(encrypted []byte) ([]byte, error) {
	msgCompressed, err := k.implementation.Decrypt(encrypted)
	if err != nil {
		return nil, err
	}
	msg, err := k.compressor.Decompress(&msgCompressed)
	if err != nil {
		return nil, err
	}
	return *msg, nil
}

func (k Key) EncryptWithAuthData(msg, authData []byte) ([]byte, error) {
	msgCompressed, err := k.compressor.Compress(&msg)
	if err != nil {
		return nil, err
	}
	return k.implementation.EncryptWithAuthData(*msgCompressed, authData)
}

func (k Key) DecryptWithAuthData(encrypted, authData []byte) ([]byte, error) {
	msgCompressed, err := k.implementation.DecryptWithAuthData(encrypted, authData)
	if err != nil {
		return nil, err
	}
	msg, err := k.compressor.Decompress(&msgCompressed)
	if err != nil {
		return nil, err
	}
	return *msg, nil
}

func (k Key) String() string {
	return k.implementation.String()
}

func (k Key) Bytes() []byte {
	return k.implementation.Bytes()
}
