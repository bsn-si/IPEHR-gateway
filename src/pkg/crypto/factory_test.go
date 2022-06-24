package crypto

import (
	"hms/gateway/pkg/crypto/chacha_poly"
	"hms/gateway/pkg/crypto/compression"
	"testing"
)

func TestSwitch(t *testing.T) {
	compressionEnabled = true

	implementation := GenerateKey()

	switch v := implementation.(type) {
	case *compression.Compression:
	default:
		t.Fatalf("Incorrect type returned for enabled compression %T", v)
	}

	compressionEnabled = false

	implementation = GenerateKey()

	switch v := implementation.(type) {
	case *chacha_poly.Key:
	default:
		t.Fatalf("Incorrect type returned for enabled compression %T", v)
	}
}
