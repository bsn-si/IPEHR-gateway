package common

// KeyInterface Key implementations interface
type KeyInterface interface {
	Encrypt(msg []byte) ([]byte, error)
	Decrypt(encrypted []byte) ([]byte, error)
	EncryptWithAuthData(msg, authData []byte) ([]byte, error)
	DecryptWithAuthData(encrypted, authData []byte) ([]byte, error)
	String() string
	Bytes() []byte
}
