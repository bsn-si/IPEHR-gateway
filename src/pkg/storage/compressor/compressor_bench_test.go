package compressor

import (
	"fmt"
	"hms/gateway/pkg/common/utils"
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

func testData() (data []byte, err error) {
	rootDir, err := utils.ProjectRootDir()
	if err != nil {
		return
	}
	filePath := rootDir + "/data/mock/ehr/composition.json"

	data, err = os.ReadFile(filePath)

	return
}
