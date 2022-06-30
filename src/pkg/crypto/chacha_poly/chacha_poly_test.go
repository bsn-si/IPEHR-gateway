package chacha_poly

import (
	"hms/gateway/pkg/common"
	"hms/gateway/pkg/common/fake_data"
	"testing"
)

func TestEncryptWith(t *testing.T) {
	key := GenerateKey()
	msg, _ := fake_data.GetByteArray(20)

	t.Logf("Test message: %x", msg)

	encrypted, err := key.Encrypt(msg)
	if err != nil {
		panic(err)
	}

	t.Logf("Encrypted message: %x", encrypted)

	decrypted, err := key.Decrypt(encrypted)
	if err != nil {
		panic(err)
	}

	t.Logf("Decrypted message: %x", decrypted)

	if !common.SliceEqualBytes(msg, decrypted) {
		panic("Decryped message mismatch!")
	}
}

func TestEncryptWithAuthData(t *testing.T) {
	key := GenerateKey()
	msg, _ := fake_data.GetByteArray(20)

	t.Logf("Test message: %x", msg)

	authData := []byte("This is a additional data")
	encrypted, err := key.EncryptWithAuthData(msg, authData)
	if err != nil {
		panic(err)
	}

	t.Logf("Encrypted message: %x", encrypted)

	decrypted, err := key.DecryptWithAuthData(encrypted, authData)
	if err != nil {
		panic(err)
	}

	t.Logf("Decrypted message: %x", decrypted)

	if !common.SliceEqualBytes(msg, decrypted) {
		panic("Decryped message mismatch!")
	}
}
