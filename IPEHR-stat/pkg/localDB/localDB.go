package localDB

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"sync"

	"github.com/golang-migrate/migrate"
	"github.com/golang-migrate/migrate/database/sqlite3"
	_ "github.com/golang-migrate/migrate/source/file"
	_ "github.com/mattn/go-sqlite3"
)

type DB struct {
	sync.Mutex
	db *sql.DB
}

func New(filePath string) *DB {
	db, err := sql.Open("sqlite3", filePath)
	if err != nil {
		log.Fatal("sql.Open error: ", err)
	}

	return &DB{sync.Mutex{}, db}
}

func (db *DB) Close() {
	db.Close()
}

func (db *DB) Migrate(migrations string) error {
	driver, err := sqlite3.WithInstance(db.db, &sqlite3.Config{})
	if err != nil {
		return fmt.Errorf("sqlite3.WithInstance error: %w", err)
	}

	m, err := migrate.NewWithDatabaseInstance("file://"+migrations, "sqlite3", driver)
	if err != nil {
		return fmt.Errorf("migrate.NewWithDatabaseInstance error: %w", err)
	}

	err = m.Up()
	if err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return fmt.Errorf("migrate.Up() error: %w", err)
	}

	return nil
}
