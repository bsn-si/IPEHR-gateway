package keybox

import (
	cryptoRand "crypto/rand"
	"testing"

	"golang.org/x/crypto/nacl/box"

	"hms/gateway/pkg/common/fake_data"
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

	msg, err := fake_data.GetByteArray(20)
	if err != nil {
		t.Fatalf("%s", err)
	}

	t.Logf("Test message: %x", msg)

	encrypted, _ := Seal(msg, recipientPublicKey, senderPrivateKey)

	t.Logf("Encrypted message: %x", encrypted)

	decrypted, _ := Open(encrypted, senderPublicKey, recipientPrivateKey)

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

	var sharedEncryptKey, sharedDecryptKey [KeyLength]byte
	Precompute(&sharedEncryptKey, recipientPublicKey, senderPrivateKey)
	Precompute(&sharedDecryptKey, senderPublicKey, recipientPrivateKey)

	t.Logf("Precomputed shared key: %x", sharedDecryptKey)

	msg, err := fake_data.GetByteArray(20)

	t.Logf("Test message: %x", msg)

	if err != nil {
		t.Fatalf("%s", err)
	}

	encrypted, _ := Seal(msg, recipientPublicKey, senderPrivateKey)

	t.Logf("Encrypted message: %x", encrypted)

	decrypted, _ := OpenAfterPrecomputation(encrypted, &sharedDecryptKey)

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

	msg, err := fake_data.GetByteArray(20)

	t.Logf("Test message: %x", msg)

	if err != nil {
		t.Fatalf("%s", err)
	}

	encrypted, _ := SealAnonymous(msg, publicKey)

	t.Logf("Encrypted message: %x", encrypted)

	decrypted, _ := OpenAnonymous(encrypted, publicKey, privateKey)

	t.Logf("Decrypted message: %x", decrypted)

	if string(msg) != string(decrypted) {
		t.Error("Anonymous encryption error")
	}
}
