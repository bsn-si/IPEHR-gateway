package infrastructure

import (
	"net/http"

	"github.com/ethereum/go-ethereum/ethclient"
	"gorm.io/gorm"

	"hms/gateway/pkg/compressor"
	"hms/gateway/pkg/config"
	"hms/gateway/pkg/docs/service/processing"
	"hms/gateway/pkg/indexer"
	"hms/gateway/pkg/keystore"
	"hms/gateway/pkg/localDB"
	log "hms/gateway/pkg/logging"
	"hms/gateway/pkg/storage"
	"hms/gateway/pkg/storage/filecoin"
	"hms/gateway/pkg/storage/ipfs"
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

func New(cfg *config.Config) *Infra {
	sc := storage.NewConfig(cfg.Storage.Localfile.Path)
	storage.Init(sc)

	db, err := localDB.New(cfg.DB.FilePath)
	if err != nil {
		log.Fatal(err, "DB path:", cfg.DB.FilePath)
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

	ks := keystore.New(cfg.KeystoreKey)

	ehtClient, err := ethclient.Dial(cfg.Contract.Endpoint)
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
		LocalDB:        db,
		Keystore:       ks,
		HTTPClient:     http.DefaultClient,
		EthClient:      ehtClient,
		IpfsClient:     ipfsClient,
		FilecoinClient: filecoinClient,
		Index: indexer.New(
			cfg.Contract.Address,
			cfg.Contract.PrivKeyPath,
			ehtClient,
			cfg.Contract.GasTipCap,
		),
		LocalStorage:       storage.Storage(),
		Compressor:         compressor.New(cfg.CompressionLevel),
		CompressionEnabled: cfg.CompressionEnabled,
	}
}
