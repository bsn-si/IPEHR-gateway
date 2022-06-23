package compressor

import (
	"fmt"
	"hms/gateway/pkg/common/utils"
	"hms/gateway/pkg/crypto/chacha_poly"
	"os"
	"testing"
)

func BenchmarkCompression(b *testing.B) {
	data, err := testData()
	if err != nil {
		b.Fatal(err)
	}

	for l := 0; l <= 9; l++ {
		b.Run("Compression "+fmt.Sprintf("%d", l), func(b *testing.B) {
			compressor := New(l)
			b.ResetTimer()

			for i := 0; i < b.N; i++ {
				compressed, _ := compressor.Compress(&data)
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

	key := chacha_poly.GenerateKey()
	encryptedData, err := key.Encrypt(data)
	if err != nil {
		t.Fatal(err)
	}

	dataSize := float32(len(data))
	dataSizeEncrypted := float32(len(encryptedData))

	// raw data
	for l := 0; l <= 9; l++ {
		compressor := New(l)
		compressed, _ := compressor.Compress(&data)
		ratio := dataSize / float32(len(*compressed))
		t.Logf("Raw data. Level: %d Ratio: %.1f times", l, ratio)
	}

	// encrypted
	for l := 0; l <= 9; l++ {
		compressor := New(l)
		compressed, _ := compressor.Compress(&encryptedData)
		ratio := dataSizeEncrypted / float32(len(*compressed))
		t.Logf("Encrypted. Level: %d Ratio: %.1f times", l, ratio)
	}
}

func testData() (data []byte, err error) {
	rootDir, err := utils.ProjectRootDir()
	if err != nil {
		return
	}
	filePath := rootDir + "/data/mock/ehr/composition.json"

	data, err = os.ReadFile(filePath)

	return
}
