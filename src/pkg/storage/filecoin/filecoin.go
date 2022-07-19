package filecoin

import (
	"context"
	"encoding/hex"
	"fmt"
	"net/http"
	"os"

	"github.com/filecoin-project/go-address"
	"github.com/filecoin-project/go-fil-markets/storagemarket"
	jsonrpc "github.com/filecoin-project/go-jsonrpc"
	lotusapi "github.com/filecoin-project/lotus/api"
	"github.com/ipfs/go-cid"
	"golang.org/x/crypto/sha3"

	"hms/gateway/pkg/errors"
)

type Storage struct {
	filesPath   string
	rpcEndpoint string
	authToken   string
	api         *lotusapi.FullNodeStruct
	closer      jsonrpc.ClientCloser
}

type Config struct {
	LotusRpcEndpoint string
	AuthToken        string
	FilesPath        string
}

func New(cfg *Config) (*Storage, error) {
	_, err := os.Stat(cfg.FilesPath)
	if os.IsNotExist(err) {
		return nil, fmt.Errorf("FilesPath is not exist: %w", err)
	}

	s := &Storage{
		filesPath:   cfg.FilesPath,
		rpcEndpoint: cfg.LotusRpcEndpoint,
		authToken:   cfg.AuthToken,
		api:         &lotusapi.FullNodeStruct{},
	}

	s.closer, err = jsonrpc.NewMergeClient(
		context.Background(),
		s.rpcEndpoint,
		"Filecoin",
		[]interface{}{&s.api.Internal, &s.api.CommonStruct.Internal},
		http.Header{"Authorization": []string{"Bearer " + s.authToken}},
	)
	if err != nil {
		return nil, fmt.Errorf("Connecting with lotus failed: %w", err)
	}

	return s, nil
}

func (s *Storage) Add(ctx context.Context, data []byte) (*[32]byte, error) {
	return nil, errors.ErrIsUnsupported
}

func (s *Storage) StartDeal(ctx context.Context, data []byte) (*cid.Cid, error) {
	// Считаем хэш
	h := sha3.Sum256(data)
	filename := hex.EncodeToString(h[:])
	filepath := s.filesPath + "/" + filename

	// Сохраняем во временный файл
	err := os.WriteFile(filepath, data, 0600)
	if err != nil {
		return nil, fmt.Errorf("WriteFile error: %w path %s", err, filepath)
	}

	// ClientImport
	importRes, err := s.api.ClientImport(ctx, lotusapi.FileRef{Path: filepath, IsCAR: false})
	if err != nil {
		return nil, fmt.Errorf("Lotus ClientImport error: %w filepath %s", err, filepath)
	}

	walletAddr, err := s.api.WalletDefaultAddress(ctx)
	if err != nil {
		return nil, fmt.Errorf("Lotus WalletDefaultAddress error: %w", err)
	}

	// MinerAddress
	// TODO нужна система выбора майнера
	minerAddr, err := address.NewFromString("f01662887")
	if err != nil {
		return nil, fmt.Errorf("Miner address parsing error: %w", err)
	}

	// ClientStartDeal
	deal, err := s.api.ClientStartDeal(ctx, &lotusapi.StartDealParams{
		Data: &storagemarket.DataRef{
			TransferType: storagemarket.TTGraphsync,
			Root:         importRes.Root,
		},
		Wallet: walletAddr,
		Miner:  minerAddr,
		//EpochPrice:        big.NewInt(500000000), // TODO What should the price be?
		MinBlocksDuration: 640000, // TODO what is it?
		DealStartEpoch:    200,    // TODO what is it?
		VerifiedDeal:      false,  // TODO make verified
		FastRetrieval:     true,
		//ProviderCollateral big.Int
	})

	return deal, nil
}

func (s *Storage) Close() {
	s.closer()
}

func (s *Storage) FilesPath() string {
	return s.filesPath
}
