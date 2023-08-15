package localDB

import (
	"errors"
	"fmt"
	"log"
	"os"
	"strings"

	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func New(filepath string) (*gorm.DB, error) {
	log.Println("try to connect to sqLITE")
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

func NewForPostgres(host string, port int, user, password, dbname string) (*gorm.DB, error) {
	// "user=gorm password=gorm dbname=gorm port=9920 sslmode=disable TimeZone=Asia/Shanghai",

	log.Println("try to connect to postgresql")
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=disable",
		host, user, password, dbname, port)
	// https://github.com/go-gorm/postgres
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	d, err := db.DB()
	if err != nil {
		return nil, err
	}

	if err := d.Ping(); err != nil {
		return nil, err
	}

	log.Println("connected to postgresql")

	return db, nil
}
