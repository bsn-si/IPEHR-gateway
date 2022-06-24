package compressor

import (
	"bytes"
	"compress/gzip"
	"io/ioutil"
)

type Compressor struct {
	compressionLevel int
}

func New(compressionLevel int) *Compressor {
	return &Compressor{
		compressionLevel: compressionLevel,
	}
}

func (c *Compressor) Compress(data []byte) (compressedData []byte, err error) {
	var buf bytes.Buffer
	zw, err := gzip.NewWriterLevel(&buf, c.compressionLevel)
	if err != nil {
		return
	}
	defer zw.Close()

	if _, err = zw.Write(data); err != nil {
		return
	}

	if err = zw.Close(); err != nil {
		return
	}

	return buf.Bytes(), nil
}

func (c *Compressor) Decompress(data []byte) (decompressedData []byte, err error) {
	buf := bytes.NewReader(data)
	zr, err := gzip.NewReader(buf)
	if err != nil {
		return
	}
	defer zr.Close()

	decompressed, err := ioutil.ReadAll(zr)
	if err != nil {
		return
	}

	return decompressed, nil
}
