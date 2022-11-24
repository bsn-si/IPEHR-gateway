package compressor

import (
	"bytes"
	"compress/gzip"
	"io"
)

const (
	NoCompression      = gzip.NoCompression
	BestSpeed          = gzip.BestSpeed
	BestCompression    = gzip.BestCompression
	DefaultCompression = gzip.DefaultCompression
	HuffmanOnly        = gzip.HuffmanOnly
)

type Compressor struct {
	level int
}

func New(level int) *Compressor {
	return &Compressor{
		level: level,
	}
}

func (c *Compressor) Compress(data []byte) (compressedData []byte, err error) {
	var buf bytes.Buffer

	zw, err := gzip.NewWriterLevel(&buf, c.level)
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

	decompressed, err := io.ReadAll(zr)
	if err != nil {
		return
	}

	return decompressed, nil
}
