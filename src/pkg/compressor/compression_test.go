package compressor

import (
	"bytes"
	"testing"

	"hms/gateway/pkg/common/fake_data"
)

func TestCompression(t *testing.T) {
	compressor := New(5)

	testData, err := fake_data.GetByteArray(1000)
	if err != nil {
		t.Fatal(err)
	}

	compressedData, err := compressor.Compress(testData)
	if err != nil {
		t.Fatal(err)
	}

	decompressedData, err := compressor.Decompress(compressedData)
	if err != nil {
		t.Fatal(err)
	}

	if !bytes.Equal(decompressedData, testData) {
		t.Fatal("Source and decompressed data is not equal")
	}
}
