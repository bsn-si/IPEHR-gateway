package compressor_test

import (
	"bytes"
	"testing"

	"github.com/bsn-si/IPEHR-gateway/src/pkg/common/fakeData"
	"github.com/bsn-si/IPEHR-gateway/src/pkg/compressor"
)

func TestCompression(t *testing.T) {
	compressor := compressor.New(5)

	testData, err := fakeData.GetByteArray(1000)
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
