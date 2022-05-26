package localfile

import (
	"encoding/hex"
	"fmt"
	"os"

	"golang.org/x/crypto/sha3"

	"hms/gateway/pkg/errors"
)

var (
	ErrIdIsTooShort = fmt.Errorf("id is too short")
)

type Config struct {
	BasePath string
	Depth    uint8
}

type LocalFileStorage struct {
	basePath string
	depth    uint8
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
	return &LocalFileStorage{
		basePath: config.BasePath,
		depth:    config.Depth,
	}, nil
}

func (s *LocalFileStorage) Add(data []byte) (id *[32]byte, err error) {
	h := sha3.Sum256(data)
	idStr := hex.EncodeToString(h[:])

	path := s.dirpath(idStr)

	if _, err = os.Stat(path); os.IsNotExist(err) {
		if err = os.MkdirAll(path, os.ModePerm); err != nil {
			return nil, err
		}
	}

	filepath := s.filepath(idStr)

	err = os.WriteFile(filepath, data, 0644)
	if err != nil {
		return nil, err
	}
	return &h, nil
}

func (s *LocalFileStorage) AddWithId(id *[32]byte, data []byte) (err error) {
	if len(*id) < int(s.depth*2) {
		return ErrIdIsTooShort
	}

	idStr := hex.EncodeToString(id[:])

	path := s.dirpath(idStr)
	if _, err = os.Stat(path); os.IsNotExist(err) {
		if err = os.MkdirAll(path, os.ModePerm); err != nil {
			return err
		}
	}

	filepath := s.filepath(idStr)
	if _, err = os.Stat(filepath); err == nil {
		return errors.AlreadyExist
	}

	err = os.WriteFile(filepath, data, 0644)
	if err != nil {
		return err
	}
	return nil
}

func (s *LocalFileStorage) ReplaceWithId(id *[32]byte, data []byte) (err error) {
	if len(*id) < int(s.depth*2) {
		return ErrIdIsTooShort
	}

	idStr := hex.EncodeToString(id[:])

	path := s.dirpath(idStr)
	if _, err = os.Stat(path); os.IsNotExist(err) {
		if err = os.MkdirAll(path, os.ModePerm); err != nil {
			return err
		}
	}

	filepath := s.filepath(idStr)
	err = os.WriteFile(filepath, data, 0644)
	if err != nil {
		return err
	}
	return nil
}

func (s *LocalFileStorage) Get(id *[32]byte) (data []byte, err error) {
	if len(*id) < int(s.depth*2) {
		return nil, ErrIdIsTooShort
	}

	idStr := hex.EncodeToString(id[:])

	path := s.filepath(idStr)
	if _, err = os.Stat(path); os.IsNotExist(err) {
		return nil, errors.IsNotExist
	}
	return os.ReadFile(path)
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
