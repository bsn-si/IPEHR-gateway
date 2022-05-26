package keybox

import (
	cryptoRand "crypto/rand"
	"testing"

	"golang.org/x/crypto/nacl/box"
	"golang.org/x/crypto/sha3"

	"hms/gateway/pkg/common/fake_data"
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

	testStrings, err := fake_data.GetByteArrays(b.N, KeyLength)

	if err != nil {
		b.Fatalf("%s", err)
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		encrypted, _ := Seal(testStrings[i], recipientPublicKey, senderPrivateKey)
		Open(encrypted, senderPublicKey, recipientPrivateKey)
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

	testStrings, err := fake_data.GetByteArrays(b.N, KeyLength)

	if err != nil {
		b.Fatalf("%s", err)
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		Seal(testStrings[i], recipientPublicKey, senderPrivateKey)
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

	var sharedEncryptKey, sharedDecryptKey [KeyLength]byte
	Precompute(&sharedEncryptKey, recipientPublicKey, senderPrivateKey)
	Precompute(&sharedDecryptKey, senderPublicKey, recipientPrivateKey)

	testStrings, err := fake_data.GetByteArrays(b.N, KeyLength)

	if err != nil {
		b.Fatalf("%s", err)
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		encrypted, _ := SealAfterPrecomputation(testStrings[i], &sharedEncryptKey)
		OpenAfterPrecomputation(encrypted, &sharedDecryptKey)
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

	var sharedEncryptKey [KeyLength]byte
	Precompute(&sharedEncryptKey, recipientPublicKey, senderPrivateKey)

	testStrings, err := fake_data.GetByteArrays(b.N, KeyLength)

	if err != nil {
		b.Fatalf("%s", err)
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		SealAfterPrecomputation(testStrings[i], &sharedEncryptKey)
	}
}

func BenchmarkAnonymousCrypt(b *testing.B) {
	publicKey, privateKey, err := box.GenerateKey(cryptoRand.Reader)
	if err != nil {
		panic(err)
	}

	testStrings, err := fake_data.GetByteArrays(b.N, KeyLength)

	if err != nil {
		b.Fatalf("%s", err)
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		encrypted, _ := SealAnonymous(testStrings[i], publicKey)
		OpenAnonymous(encrypted, publicKey, privateKey)
	}
}

// Bench Seal without Opens for recognize which part took more time
func BenchmarkAnonymousCryptSealOnly(b *testing.B) {
	publicKey, _, err := box.GenerateKey(cryptoRand.Reader)
	if err != nil {
		panic(err)
	}

	testStrings, err := fake_data.GetByteArrays(b.N, KeyLength)

	if err != nil {
		b.Fatalf("%s", err)
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		SealAnonymous(testStrings[i], publicKey)
	}
}

func BenchmarkSha3_256(b *testing.B) {
	testStrings, err := fake_data.GetByteArrays(b.N, 64)

	if err != nil {
		b.Fatalf("%s", err)
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		sha3.Sum256(testStrings[i])
	}
}
