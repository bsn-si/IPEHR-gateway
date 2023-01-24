package chachaPoly_test

import (
	"bytes"
	"testing"

	"github.com/bsn-si/IPEHR-gateway/src/pkg/common/fakeData"
	"github.com/bsn-si/IPEHR-gateway/src/pkg/crypto/chachaPoly"
)

func TestEncryptWith(t *testing.T) {
	key := chachaPoly.GenerateKey()
	msg, _ := fakeData.GetByteArray(20)

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

	if !bytes.Equal(msg, decrypted) {
		panic("Decryped message mismatch!")
	}
}

func TestEncryptWithAuthData(t *testing.T) {
	key := chachaPoly.GenerateKey()
	msg, _ := fakeData.GetByteArray(20)

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

	if !bytes.Equal(msg, decrypted) {
		panic("Decryped message mismatch!")
	}
}
