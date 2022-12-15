package infrastructure

import (
	"log"

	"github.com/ethereum/go-ethereum/ethclient"

	"ipehr/stat/pkg/config"
	"ipehr/stat/pkg/localDB"
)

type Infra struct {
	DB        *localDB.DB
	EthClient *ethclient.Client
}

func New(cfg *config.Config) *Infra {
	ehtClient, err := ethclient.Dial(cfg.Sync.Endpoint)
	if err != nil {
		log.Fatal(err)
	}

	db := localDB.New(cfg.LocalDB.Path)

	err = db.Migrate(cfg.LocalDB.Migrations)
	if err != nil {
		log.Fatal(err)
	}

	return &Infra{
		DB:        db,
		EthClient: ehtClient,
	}
}

func (i *Infra) Close() {
	i.DB.Close()
}
