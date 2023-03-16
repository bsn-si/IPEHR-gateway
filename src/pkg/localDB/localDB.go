package localDB

import (
	"errors"
	"log"
	"os"
	"strings"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func New(filepath string) (*gorm.DB, error) {
	if _, err := os.Stat(filepath); errors.Is(err, os.ErrNotExist) {
		s := strings.Split(filepath, "/")
		if len(s) > 1 {
			s = s[:len(s)-1]

			err = os.MkdirAll(strings.Join(s, "/"), 0777)
			if err != nil {
				log.Fatal(err)
			}

			filepath = strings.Join(s[:len(s)+1], "/")
		}

		f, err := os.OpenFile(filepath, os.O_CREATE, 0600)
		if err != nil {
			log.Fatal(err)
		}

		f.Close()
	}

	db, err := gorm.Open(sqlite.Open(filepath), &gorm.Config{
		//Logger: logger.Default.LogMode(logger.Error),
		Logger: logger.Default.LogMode(logger.Warn),
	})
	if err != nil {
		return nil, err
	}

	return db, nil
}
