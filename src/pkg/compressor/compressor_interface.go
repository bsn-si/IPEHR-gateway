package compressor

type CompressorInterface interface {
	Compress(decompressedData *[]byte) (compressedData *[]byte, err error)
	Decompress(compressedData *[]byte) (decompressedData *[]byte, err error)
}
