package storage

import (
	"errors"
	"log"
	"os"
	"path/filepath"
	"strings"
)

type StorageConfig struct {
	processPath string
	path        string
}

func NewConfig(path string) (config *StorageConfig) {
	config = &StorageConfig{}
	path = config.prepare(path)
	err := config.valid(path)
	if err != nil {
		log.Panicln(err, path)
	}

	config.setPath(path)
	return config
}

// Absolute path
func (c *StorageConfig) Path() string {
	return c.path
}

// Absolute path of current executable process file
func (c *StorageConfig) ProcessPath() string {
	return c.processPath
}

func (c *StorageConfig) setPath(path string) {
	c.path = path
}

func (c *StorageConfig) setProcessPath(path string) {
	c.processPath = path
}

func (c *StorageConfig) prepare(path string) string {
	processPath, err := os.Executable()
	if err != nil {
		log.Panicln(err)
	}
	c.setProcessPath(filepath.Dir(processPath))

	path = strings.TrimSpace(path)

	if len(path) == 0 {
		log.Panicln("Storage path is empty")
	}

	if path[0:1] != "/" {
		path = c.ProcessPath() + "/" + path
	}

	path = filepath.Clean(path)

	return path
}

func (c *StorageConfig) valid(path string) (err error) {
	if len(path) == 0 {
		return errors.New("Storage path is empty")
	}

	if path == "/" {
		return errors.New("Can not use root folder as a storage")
	}

	if !strings.HasPrefix(path, c.ProcessPath()) {
		log.Printf("Notice: storage path '%s' is out of base path: %s", path, c.ProcessPath())
	}

	return nil
}
