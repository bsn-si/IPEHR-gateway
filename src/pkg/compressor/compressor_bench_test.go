package compressor_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/bsn-si/IPEHR-gateway/src/pkg/compressor"
	"github.com/bsn-si/IPEHR-gateway/src/pkg/crypto/chachaPoly"
)

func BenchmarkCompression(b *testing.B) {
	data, err := testData()
	if err != nil {
		b.Fatal(err)
	}

	for l := 0; l <= 9; l++ {
		b.Run("Compression "+fmt.Sprintf("%d", l), func(b *testing.B) {
			compressor := compressor.New(l)
			b.ResetTimer()

			for i := 0; i < b.N; i++ {
				compressed, _ := compressor.Compress(data)
				_, _ = compressor.Decompress(compressed)
			}
		})
	}
}

// If someone wants to check compression ratio
func TestCompressionRatio(t *testing.T) {
	t.Skip()

	data, err := testData()
	if err != nil {
		t.Fatal(err)
	}

	key := chachaPoly.GenerateKey()

	encryptedData, err := key.Encrypt(data)
	if err != nil {
		t.Fatal(err)
	}

	dataSize := float32(len(data))
	dataSizeEncrypted := float32(len(encryptedData))

	// raw data
	for l := 0; l <= 9; l++ {
		compressor := compressor.New(l)
		compressed, _ := compressor.Compress(data)
		ratio := dataSize / float32(len(compressed))
		t.Logf("Raw data. Level: %d Ratio: %.1f times", l, ratio)
	}

	// encrypted
	for l := 0; l <= 9; l++ {
		compressor := compressor.New(l)
		compressed, _ := compressor.Compress(encryptedData)
		ratio := dataSizeEncrypted / float32(len(compressed))
		t.Logf("Encrypted. Level: %d Ratio: %.1f times", l, ratio)
	}
}

func testData() (data []byte, err error) {
	filePath := "./test_fixtures/composition.json"

	data, err = os.ReadFile(filePath)

	return
}
