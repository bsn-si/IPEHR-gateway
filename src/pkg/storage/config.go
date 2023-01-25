package storage

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/bsn-si/IPEHR-gateway/src/pkg/errors"
)

type Config struct {
	processPath string
	path        string
}

func NewConfig(path string) (config *Config) {
	config = &Config{}
	path = config.prepare(path)

	if err := config.valid(path); err != nil {
		log.Panicln(err, path)
	}

	config.setPath(path)

	return config
}

// Absolute path
func (c *Config) Path() string {
	return c.path
}

// Absolute path of current executable process file
func (c *Config) ProcessPath() string {
	return c.processPath
}

func (c *Config) setPath(path string) {
	c.path = path
}

func (c *Config) setProcessPath(path string) {
	c.processPath = path
}

func (c *Config) prepare(path string) string {
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

func (c *Config) valid(path string) (err error) {
	if len(path) == 0 {
		return fmt.Errorf("Storage path is empty: %w", errors.ErrIsEmpty)
	}

	if path == "/" {
		return fmt.Errorf("%w: Can not use root folder as a storage", errors.ErrCustom)
	}

	if !strings.HasPrefix(path, c.ProcessPath()) {
		log.Printf("Notice: storage path '%s' is out of base path: %s", path, c.ProcessPath())
	}

	return nil
}
