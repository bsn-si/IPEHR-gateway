package localfile

import (
	"encoding/hex"
	"fmt"
	"log"
	"os"

	"golang.org/x/crypto/sha3"

	"github.com/bsn-si/IPEHR-gateway/src/pkg/errors"
)

type Config struct {
	BasePath string
	Depth    uint8
}

type Storage struct {
	basePath string
	depth    uint8
}

func Init(config *Config) (*Storage, error) {
	if len(config.BasePath) == 0 {
		return nil, fmt.Errorf("%w: BasePath", errors.ErrIsEmpty)
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

	return &Storage{
		basePath: config.BasePath,
		depth:    config.Depth,
	}, nil
}

func (s *Storage) Add(data []byte) (id *[32]byte, err error) {
	id = s.idByContent(&data)

	err = s.writeFile(id, &data)

	return
}

func (s *Storage) idByContent(data *[]byte) *[32]byte {
	h := sha3.Sum256(*data)
	return &h
}

func (s *Storage) ReplaceWithID(id *[32]byte, data []byte) (err error) {
	return s.AddWithID(id, data)
}

func (s *Storage) AddWithID(id *[32]byte, data []byte) (err error) {
	err = s.writeFile(id, &data)
	return
}

func (s *Storage) Get(id *[32]byte) (data []byte, err error) {
	idStr := hex.EncodeToString(id[:])

	path := s.filepath(idStr)
	if _, err = os.Stat(path); os.IsNotExist(err) {
		return nil, errors.ErrIsNotExist
	}

	data, err = os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("os.ReadFile error: %w", err)
	}

	return
}

func (s *Storage) writeFile(id *[32]byte, data *[]byte) (err error) {
	idStr := hex.EncodeToString(id[:])

	path := s.dirpath(idStr)
	if _, err = os.Stat(path); os.IsNotExist(err) {
		if err = os.MkdirAll(path, os.ModePerm); err != nil {
			return err
		}
	}

	filepath := s.filepath(idStr)

	err = os.WriteFile(filepath, *data, 0600)
	if err != nil {
		return err
	}

	return nil
}

func (s *Storage) dirpath(id string) (path string) {
	path = s.basePath
	for i := 0; i < int(s.depth)*2; i = i + 2 {
		path += id[i:i+2] + "/"
	}

	return path
}

func (s *Storage) filepath(id string) (path string) {
	path = s.basePath
	for i := 0; i < int(s.depth)*2; i = i + 2 {
		path += id[i:i+2] + "/"
	}

	return path + id
}

func (s *Storage) Clean() (err error) {
	if s.basePath == "/" {
		log.Panicln("Can not clean the base folder is root!")
	}

	if _, err = os.Stat(s.basePath); err != nil {
		return err
	}

	if err = os.RemoveAll(s.basePath); err != nil {
		return err
	}

	return nil
}
