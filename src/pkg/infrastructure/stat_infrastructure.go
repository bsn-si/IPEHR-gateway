package infrastructure

import (
	"errors"
	"fmt"
	"log"

	"github.com/bsn-si/IPEHR-gateway/src/internal/repository"
	_ "github.com/bsn-si/IPEHR-gateway/src/pkg/aql/driver" //nolint
	"github.com/bsn-si/IPEHR-gateway/src/pkg/config"
	"github.com/bsn-si/IPEHR-gateway/src/pkg/service/stat"

	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/golang-migrate/migrate"
	"github.com/golang-migrate/migrate/database/sqlite3"
	_ "github.com/golang-migrate/migrate/source/file" //nolint
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3" //nolint
)

type StatInfra struct {
	DB        *sqlx.DB
	EthClient *ethclient.Client
	AqlDB     *sqlx.DB

	StatsRepo *repository.StatsStorage
	ChunkRepo *repository.IndexStorage
	Service   *stat.Service
}

func NewStatInfra(cfg *config.StatConfig) *StatInfra {
	ehtClient, err := ethclient.Dial(cfg.Sync.Endpoint)
	if err != nil {
		log.Fatal(err)
	}

	db, err := sqlx.Connect("sqlite3", cfg.LocalDB.Path)
	if err != nil {
		log.Fatal("sql.Open error: ", err)
	}

	if err := migrateDB(db, cfg.LocalDB.Migrations); err != nil {
		log.Fatal(err)
	}

	aqlDB, err := sqlx.Open("aql", "")
	if err != nil {
		log.Fatal(err)
	}

	statsRepo := repository.NetStatsSotrage(db)
	svc := stat.NewService(statsRepo)

	return &StatInfra{
		DB:        db,
		EthClient: ehtClient,
		AqlDB:     aqlDB,
		StatsRepo: statsRepo,
		ChunkRepo: repository.NewIndexStorage(db),
		Service:   svc,
	}
}

func (i *StatInfra) Close() {
	i.DB.Close()
	i.AqlDB.Close()
}

func migrateDB(db *sqlx.DB, migrations string) error {
	driver, err := sqlite3.WithInstance(db.DB, &sqlite3.Config{})
	if err != nil {
		return fmt.Errorf("sqlite3.WithInstance error: %w", err)
	}

	m, err := migrate.NewWithDatabaseInstance("file://"+migrations, "sqlite3", driver)
	if err != nil {
		return fmt.Errorf("migrate.NewWithDatabaseInstance error: %w", err)
	}

	if err := m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return fmt.Errorf("migrate.Up() error: %w", err)
	}

	return nil
}
