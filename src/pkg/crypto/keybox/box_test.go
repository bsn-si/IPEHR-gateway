package keybox_test

import (
	cryptoRand "crypto/rand"
	"testing"

	"github.com/bsn-si/IPEHR-gateway/src/pkg/common/fakeData"
	"github.com/bsn-si/IPEHR-gateway/src/pkg/crypto/keybox"

	"golang.org/x/crypto/nacl/box"
)

func TestCrypt(t *testing.T) {
	senderPublicKey, senderPrivateKey, err := box.GenerateKey(cryptoRand.Reader)
	if err != nil {
		panic(err)
	}

	recipientPublicKey, recipientPrivateKey, err := box.GenerateKey(cryptoRand.Reader)
	if err != nil {
		panic(err)
	}

	msg, err := fakeData.GetByteArray(20)
	if err != nil {
		t.Fatalf("%s", err)
	}

	t.Logf("Test message: %x", msg)

	encrypted, _ := keybox.Seal(msg, recipientPublicKey, senderPrivateKey)

	t.Logf("Encrypted message: %x", encrypted)

	decrypted, _ := keybox.Open(encrypted, senderPublicKey, recipientPrivateKey)

	t.Logf("Decrypted message: %x", decrypted)

	if string(msg) != string(decrypted) {
		t.Error("Crypting error")
	}
}

func TestPrecomputedCrypt(t *testing.T) {
	senderPublicKey, senderPrivateKey, err := box.GenerateKey(cryptoRand.Reader)
	if err != nil {
		panic(err)
	}

	recipientPublicKey, recipientPrivateKey, err := box.GenerateKey(cryptoRand.Reader)
	if err != nil {
		panic(err)
	}

	var sharedEncryptKey, sharedDecryptKey [keybox.KeyLength]byte

	keybox.Precompute(&sharedEncryptKey, recipientPublicKey, senderPrivateKey)
	keybox.Precompute(&sharedDecryptKey, senderPublicKey, recipientPrivateKey)

	t.Logf("Precomputed shared key: %x", sharedDecryptKey)

	msg, err := fakeData.GetByteArray(20)

	t.Logf("Test message: %x", msg)

	if err != nil {
		t.Fatalf("%s", err)
	}

	encrypted, _ := keybox.Seal(msg, recipientPublicKey, senderPrivateKey)

	t.Logf("Encrypted message: %x", encrypted)

	decrypted, _ := keybox.OpenAfterPrecomputation(encrypted, &sharedDecryptKey)

	t.Logf("Decrypted message: %x", decrypted)

	if string(msg) != string(decrypted) {
		t.Error("Decrypted string is not equal with original string")
	}
}

func TestAnonymousCrypt(t *testing.T) {
	publicKey, privateKey, err := box.GenerateKey(cryptoRand.Reader)
	if err != nil {
		panic(err)
	}

	msg, err := fakeData.GetByteArray(20)

	t.Logf("Test message: %x", msg)

	if err != nil {
		t.Fatalf("%s", err)
	}

	encrypted, _ := keybox.SealAnonymous(msg, publicKey)

	t.Logf("Encrypted message: %x", encrypted)

	decrypted, _ := keybox.OpenAnonymous(encrypted, publicKey, privateKey)

	t.Logf("Decrypted message: %x", decrypted)

	if string(msg) != string(decrypted) {
		t.Error("Anonymous encryption error")
	}
}
