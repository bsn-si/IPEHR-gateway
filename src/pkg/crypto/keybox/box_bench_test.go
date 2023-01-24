package keybox_test

import (
	cryptoRand "crypto/rand"
	"testing"

	"golang.org/x/crypto/nacl/box"
	"golang.org/x/crypto/sha3"

	"github.com/bsn-si/IPEHR-gateway/src/pkg/common/fakeData"
	"github.com/bsn-si/IPEHR-gateway/src/pkg/crypto/keybox"
)

func BenchmarkCrypt(b *testing.B) {
	senderPublicKey, senderPrivateKey, err := box.GenerateKey(cryptoRand.Reader)
	if err != nil {
		panic(err)
	}

	recipientPublicKey, recipientPrivateKey, err := box.GenerateKey(cryptoRand.Reader)
	if err != nil {
		panic(err)
	}

	testStrings, err := fakeData.GetByteArrays(b.N, keybox.KeyLength)

	if err != nil {
		b.Fatalf("%s", err)
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		encrypted, _ := keybox.Seal(testStrings[i], recipientPublicKey, senderPrivateKey)

		if _, err = keybox.Open(encrypted, senderPublicKey, recipientPrivateKey); err != nil {
			b.Fatal(err)
		}
	}
}

// Bench Seal without Opens for recognize which part took more time
func BenchmarkCryptSealOnly(b *testing.B) {
	_, senderPrivateKey, err := box.GenerateKey(cryptoRand.Reader)
	if err != nil {
		panic(err)
	}

	recipientPublicKey, _, err := box.GenerateKey(cryptoRand.Reader)
	if err != nil {
		panic(err)
	}

	testStrings, err := fakeData.GetByteArrays(b.N, keybox.KeyLength)
	if err != nil {
		b.Fatalf("%s", err)
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_, err = keybox.Seal(testStrings[i], recipientPublicKey, senderPrivateKey)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkPrecomputedCrypt(b *testing.B) {
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

	testStrings, err := fakeData.GetByteArrays(b.N, keybox.KeyLength)
	if err != nil {
		b.Fatalf("%s", err)
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		encrypted, _ := keybox.SealAfterPrecomputation(testStrings[i], &sharedEncryptKey)

		_, err = keybox.OpenAfterPrecomputation(encrypted, &sharedDecryptKey)
		if err != nil {
			b.Fatal(err)
		}
	}
}

// Bench Seal without Opens for recognize which part took more time
func BenchmarkPrecomputedCryptSealOnly(b *testing.B) {
	_, senderPrivateKey, err := box.GenerateKey(cryptoRand.Reader)
	if err != nil {
		panic(err)
	}

	recipientPublicKey, _, err := box.GenerateKey(cryptoRand.Reader)
	if err != nil {
		panic(err)
	}

	var sharedEncryptKey [keybox.KeyLength]byte

	keybox.Precompute(&sharedEncryptKey, recipientPublicKey, senderPrivateKey)

	testStrings, err := fakeData.GetByteArrays(b.N, keybox.KeyLength)
	if err != nil {
		b.Fatalf("%s", err)
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_, err = keybox.SealAfterPrecomputation(testStrings[i], &sharedEncryptKey)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkAnonymousCrypt(b *testing.B) {
	publicKey, privateKey, err := box.GenerateKey(cryptoRand.Reader)
	if err != nil {
		panic(err)
	}

	testStrings, err := fakeData.GetByteArrays(b.N, keybox.KeyLength)

	if err != nil {
		b.Fatalf("%s", err)
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		encrypted, _ := keybox.SealAnonymous(testStrings[i], publicKey)

		if _, err = keybox.OpenAnonymous(encrypted, publicKey, privateKey); err != nil {
			b.Fatal(err)
		}
	}
}

// Bench Seal without Opens for recognize which part took more time
func BenchmarkAnonymousCryptSealOnly(b *testing.B) {
	publicKey, _, err := box.GenerateKey(cryptoRand.Reader)
	if err != nil {
		panic(err)
	}

	testStrings, err := fakeData.GetByteArrays(b.N, keybox.KeyLength)

	if err != nil {
		b.Fatalf("%s", err)
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		if _, err = keybox.SealAnonymous(testStrings[i], publicKey); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkSha3_256(b *testing.B) {
	testStrings, err := fakeData.GetByteArrays(b.N, 64)

	if err != nil {
		b.Fatalf("%s", err)
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		sha3.Sum256(testStrings[i])
	}
}
