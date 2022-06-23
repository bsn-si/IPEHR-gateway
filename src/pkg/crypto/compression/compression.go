package compression

import (
	"hms/gateway/pkg/compressor"
	"hms/gateway/pkg/config"
	"hms/gateway/pkg/crypto/common"
)

type Compression struct {
	implementation common.KeyInterface
	compressor     compressor.CompressorInterface
}

func New(implementation common.KeyInterface) common.KeyInterface {
	globalConfig, err := config.New()
	if err != nil {
		return nil
	}
	return &Compression{
		implementation: implementation,
		compressor:     compressor.New(globalConfig.CompressionLevel),
	}
}

func NewKeyFromBytes(implementation common.KeyInterface) common.KeyInterface {
	return New(implementation)
}

func (c *Compression) Encrypt(msg []byte) ([]byte, error) {
	msgCompressed, err := c.compressor.Compress(&msg)
	if err != nil {
		return nil, err
	}
	return c.implementation.Encrypt(*msgCompressed)
}

func (c *Compression) Decrypt(encrypted []byte) ([]byte, error) {
	msgCompressed, err := c.implementation.Decrypt(encrypted)
	if err != nil {
		return nil, err
	}
	msg, err := c.compressor.Decompress(&msgCompressed)
	if err != nil {
		return nil, err
	}
	return *msg, nil
}

func (c *Compression) EncryptWithAuthData(msg, authData []byte) ([]byte, error) {
	msgCompressed, err := c.compressor.Compress(&msg)
	if err != nil {
		return nil, err
	}
	return c.implementation.EncryptWithAuthData(*msgCompressed, authData)
}

func (c *Compression) DecryptWithAuthData(encrypted, authData []byte) ([]byte, error) {
	msgCompressed, err := c.implementation.DecryptWithAuthData(encrypted, authData)
	if err != nil {
		return nil, err
	}
	msg, err := c.compressor.Decompress(&msgCompressed)
	if err != nil {
		return nil, err
	}
	return *msg, nil
}

func (c *Compression) String() string {
	return c.implementation.String()
}

func (c *Compression) Bytes() []byte {
	return c.implementation.Bytes()
}
