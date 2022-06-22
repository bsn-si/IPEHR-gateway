package localfile

import (
	"bytes"
	"compress/gzip"
	"encoding/hex"
	"fmt"
	"golang.org/x/crypto/sha3"
	config2 "hms/gateway/pkg/config"
	"hms/gateway/pkg/errors"
	"io/ioutil"
	"log"
	"os"
)

type Config struct {
	BasePath string
	Depth    uint8
}

type LocalFileStorage struct {
	basePath           string
	depth              uint8
	compressionEnabled bool
	compressionLevel   int
}

func Init(config *Config) (*LocalFileStorage, error) {
	if len(config.BasePath) == 0 {
		return nil, fmt.Errorf("BasePath is empty")
	}

	if config.Depth == 0 {
		config.Depth = 1
	}

	if config.BasePath[len(config.BasePath)-1] != '/' {
		config.BasePath += "/"
	}

	_, err := os.Stat(config.BasePath)
	if os.IsNotExist(err) {
		if err = os.MkdirAll(config.BasePath, os.ModePerm); err != nil {
			return nil, err
		}
	}

	globalConfig, err := config2.New()
	if err != nil {
		return nil, err
	}

	return &LocalFileStorage{
		basePath:           config.BasePath,
		depth:              config.Depth,
		compressionEnabled: globalConfig.CompressionEnabled,
		compressionLevel:   globalConfig.CompressionLevel,
	}, nil
}

func (s *LocalFileStorage) Add(data []byte) (id *[32]byte, err error) {
	h := sha3.Sum256(data)

	id = &h

	err = s.writeFile(id, &data)
	return
}

func (s *LocalFileStorage) ReplaceWithId(id *[32]byte, data []byte) (err error) {
	return s.AddWithId(id, data)
}

func (s *LocalFileStorage) AddWithId(id *[32]byte, data []byte) (err error) {
	err = s.writeFile(id, &data)
	return
}

func (s *LocalFileStorage) Get(id *[32]byte) (data []byte, err error) {
	idStr := hex.EncodeToString(id[:])

	path := s.filepath(idStr)
	if _, err = os.Stat(path); os.IsNotExist(err) {
		return nil, errors.IsNotExist
	}
	data, err = os.ReadFile(path)
	if err != nil {
		return
	}

	if s.compressionEnabled {
		dataDecompressed, err := s.decompress(&data)
		if err != nil {
			return nil, err
		}
		data = *dataDecompressed
	}

	return
}

func (s *LocalFileStorage) writeFile(id *[32]byte, data *[]byte) (err error) {
	idStr := hex.EncodeToString(id[:])

	path := s.dirpath(idStr)
	if _, err = os.Stat(path); os.IsNotExist(err) {
		if err = os.MkdirAll(path, os.ModePerm); err != nil {
			return err
		}
	}

	if s.compressionEnabled {
		data, err = s.compress(data)
		if err != nil {
			return
		}
	}

	filepath := s.filepath(idStr)
	err = os.WriteFile(filepath, *data, 0644)
	if err != nil {
		return err
	}
	return nil
}

func (s *LocalFileStorage) compress(data *[]byte) (compressedData *[]byte, err error) {
	var buf bytes.Buffer
	zw, err := gzip.NewWriterLevel(&buf, s.compressionLevel)
	if err != nil {
		return
	}

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

func (s *LocalFileStorage) decompress(data *[]byte) (decompressedData *[]byte, err error) {
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

func (s *LocalFileStorage) dirpath(id string) (path string) {
	path = s.basePath
	i := 0
	for i < int(s.depth)*2 {
		path += id[i:i+2] + "/"
		i += 2
	}
	return path
}

func (s *LocalFileStorage) filepath(id string) (path string) {
	path = s.basePath
	i := 0
	for i < int(s.depth)*2 {
		path += id[i:i+2] + "/"
		i += 2
	}
	return path + id
}

func (s *LocalFileStorage) Clean() (err error) {
	if s.basePath == "/" {
		log.Panicln("Can not clean the base folder is root!")
	}

	_, err = os.Stat(s.basePath)
	if err != nil {
		return nil
	}

	if err = os.RemoveAll(s.basePath); err != nil {
		return err
	}

	return nil
}
