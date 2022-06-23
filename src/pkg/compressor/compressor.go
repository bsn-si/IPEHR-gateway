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

func (c *Compressor) Compress(data *[]byte) (compressedData *[]byte, err error) {
	var buf bytes.Buffer
	zw, err := gzip.NewWriterLevel(&buf, c.compressionLevel)
	if err != nil {
		return
	}

	defer func(zw *gzip.Writer) {
		_ = zw.Close()
	}(zw)

	_, err = zw.Write(*data)
	if err != nil {
		return
	}

	err = zw.Close()
	if err != nil {
		return
	}

	bufBytes := buf.Bytes()
	compressedData = &bufBytes

	return
}

func (c *Compressor) Decompress(data *[]byte) (decompressedData *[]byte, err error) {
	buf := bytes.NewReader(*data)
	zr, err := gzip.NewReader(buf)
	if err != nil {
		return
	}

	defer func(zr *gzip.Reader) {
		_ = zr.Close()
	}(zr)

	decompressed, err := ioutil.ReadAll(zr)
	if err != nil {
		return
	}

	decompressedData = &decompressed

	return
}
