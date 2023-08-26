package infrastructure

import (
	"context"
	"log"
	"net/http"

	"github.com/ethereum/go-ethereum/ethclient"
	"gorm.io/gorm"

	"github.com/bsn-si/IPEHR-gateway/src/pkg/compressor"
	"github.com/bsn-si/IPEHR-gateway/src/pkg/config"
	"github.com/bsn-si/IPEHR-gateway/src/pkg/docs/service/processing"
	"github.com/bsn-si/IPEHR-gateway/src/pkg/indexer"
	"github.com/bsn-si/IPEHR-gateway/src/pkg/keystore"
	"github.com/bsn-si/IPEHR-gateway/src/pkg/localDB"
	"github.com/bsn-si/IPEHR-gateway/src/pkg/storage"
	"github.com/bsn-si/IPEHR-gateway/src/pkg/storage/filecoin"
	"github.com/bsn-si/IPEHR-gateway/src/pkg/storage/ipfs"
)

type Infra struct {
	LocalDB            *gorm.DB
	Keystore           *keystore.KeyStore
	HTTPClient         *http.Client
	EthClient          *ethclient.Client
	IpfsClient         *ipfs.Client
	FilecoinClient     *filecoin.Client
	Index              *indexer.Index
	LocalStorage       storage.Storager
	Compressor         compressor.Interface
	CompressionEnabled bool
}

func New(ctx context.Context, cfg *config.Config) *Infra {
	sc := storage.NewConfig(cfg.Storage.Localfile.Path)
	storage.Init(sc)

	var (
		db  *gorm.DB
		err error
	)

	if cfg.DB.UsePostgres {
		db, err = localDB.NewForPostgres(
			cfg.DB.Postgres.Host,
			cfg.DB.Postgres.Port,
			cfg.DB.Postgres.User,
			cfg.DB.Postgres.Password,
			cfg.DB.Postgres.DBName,
		)
		if err != nil {
			log.Fatal(err, "Postgres Connecton err:", cfg.DB.FilePath)
		}
	} else {
		db, err = localDB.New(cfg.DB.FilePath)
		if err != nil {
			log.Fatal(err, "DB path:", cfg.DB.FilePath)
		}
	}

	if err = db.AutoMigrate(&processing.Request{}); err != nil {
		log.Fatal(err)
	}

	if err = db.AutoMigrate(&processing.Retrieve{}); err != nil {
		log.Fatal(err)
	}

	if err = db.AutoMigrate(&processing.EthereumTx{}); err != nil {
		log.Fatal(err)
	}

	if err = db.AutoMigrate(&processing.FileCoinTx{}); err != nil {
		log.Fatal(err)
	}

	ethClient, err := ethclient.Dial(cfg.Blockchain.Endpoint)
	if err != nil {
		log.Fatal(err)
	}

	ipfsClient, err := ipfs.NewClient(cfg.Storage.Ipfs.EndpointURLs)
	if err != nil {
		log.Fatal(err)
	}

	filecoinCfg := filecoin.Config(cfg.Storage.Filecoin)

	filecoinClient, err := filecoin.NewClient(&filecoinCfg)
	if err != nil {
		log.Fatal(err)
	}

	return &Infra{
		LocalDB:            db,
		Keystore:           keystore.New(cfg.KeystoreKey),
		HTTPClient:         http.DefaultClient,
		EthClient:          ethClient,
		IpfsClient:         ipfsClient,
		FilecoinClient:     filecoinClient,
		Index:              indexer.New(cfg.Blockchain, ethClient),
		LocalStorage:       storage.Storage(),
		Compressor:         compressor.New(cfg.CompressionLevel),
		CompressionEnabled: cfg.CompressionEnabled,
	}
}
