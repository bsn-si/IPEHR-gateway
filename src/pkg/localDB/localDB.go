package localDB

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func New(filepath string) (*gorm.DB, error) {
	db, err := gorm.Open(sqlite.Open(filepath), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	return db, nil
}
