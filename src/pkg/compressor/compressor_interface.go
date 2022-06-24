package compressor

type Interface interface {
	Compress(decompressedData []byte) (compressedData []byte, err error)
	Decompress(compressedData []byte) (decompressedData []byte, err error)
}
