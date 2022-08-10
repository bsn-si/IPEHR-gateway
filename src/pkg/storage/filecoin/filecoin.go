package filecoin

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/filecoin-project/go-address"
	"github.com/filecoin-project/go-fil-markets/storagemarket"
	jsonrpc "github.com/filecoin-project/go-jsonrpc"
	lotusapi "github.com/filecoin-project/lotus/api"
	"github.com/ipfs/go-cid"

	"hms/gateway/pkg/errors"
)

type DealStatus = storagemarket.StorageDealStatus

type Client struct {
	rpcEndpoint   string
	authToken     string
	dealsMaxPrice uint64
	api           *lotusapi.FullNodeStruct
	closer        jsonrpc.ClientCloser
	httpClient    *http.Client
}

type Config struct {
	LotusRpcEndpoint string
	AuthToken        string
	DealsMaxPrice    uint64
}

type filrepMinersResult struct {
	Miners []struct {
		Id              uint
		Address         string
		Status          bool
		Reachability    string
		UptimeAverage   float64
		Price           string
		VerifiedPrice   string
		MinPieceSize    string
		MaxPieceSize    string
		RawPower        string
		QualityAdjPower string
		IsoCode         string
		Region          string
		Score           string
		Scores          struct {
			Total                  string
			Uptime                 string
			StorageDeals           string
			CommittedSectorsProofs string
		}
		FreeSpace    string
		StorageDeals struct {
			Total           uint
			NoPenalties     uint
			SuccessRate     string
			AveragePrice    string
			DataStored      string
			Slashed         uint
			Terminated      uint
			FaultTerminated uint
			Recent30days    uint
		}
		GoldenPath struct {
			StorageDealSuccessRate   bool
			RetrievalDealSuccessRate bool
			FastRetrieval            *bool
			VerifiedDealNoPrice      bool
			FaultsBelowThreshold     bool
		}
		Rank       string
		RegionRank string
	}
}

func NewClient(cfg *Config) (*Client, error) {
	c := &Client{
		rpcEndpoint:   cfg.LotusRpcEndpoint,
		authToken:     cfg.AuthToken,
		dealsMaxPrice: cfg.DealsMaxPrice,
		api:           &lotusapi.FullNodeStruct{},
		httpClient:    http.DefaultClient,
	}

	var err error
	c.closer, err = jsonrpc.NewMergeClient(
		context.Background(),
		c.rpcEndpoint,
		"Filecoin",
		[]interface{}{&c.api.Internal, &c.api.CommonStruct.Internal},
		http.Header{"Authorization": []string{"Bearer " + c.authToken}},
	)
	if err != nil {
		return nil, fmt.Errorf("Connecting with lotus failed: %w", err)
	}

	return c, nil
}

func (c *Client) StartDeal(ctx context.Context, CID *cid.Cid, dataSizeBytes uint64) (*cid.Cid, []byte, error) {
	walletAddr, err := c.api.WalletDefaultAddress(ctx)
	if err != nil {
		return nil, nil, fmt.Errorf("Lotus WalletDefaultAddress error: %w", err)
	}

	// MinerAddress
	minerAddr, err := c.FindMiner(dataSizeBytes)
	if err != nil {
		return nil, nil, fmt.Errorf("Miner address parsing error: %w", err)
	}

	// ClientStartDeal
	deal, err := c.api.ClientStartDeal(ctx, &lotusapi.StartDealParams{
		Data: &storagemarket.DataRef{
			TransferType: storagemarket.TTGraphsync,
			Root:         *CID,
		},
		Wallet: walletAddr,
		Miner:  *minerAddr,
		//EpochPrice:        big.NewInt(500000000), // TODO What should the price be?
		//MinBlocksDuration: 640000, // TODO what is it?
		//DealStartEpoch:    200,
		VerifiedDeal:  false, // TODO make verified
		FastRetrieval: true,
		//ProviderCollateral big.Int
	})

	return deal, minerAddr.Bytes(), nil
}

func (c *Client) GetDealStatus(ctx context.Context, CID *cid.Cid) (storagemarket.StorageDealStatus, error) {
	dealInfo, err := c.api.ClientGetDealInfo(ctx, *CID)
	if err != nil {
		return 0, fmt.Errorf("Lotus ClientGetDealInfo error: %w CID %s", err, CID.String())
	}

	return dealInfo.State, nil
}

func (c *Client) Close() {
	c.closer()
}

func (c *Client) FindMiner(dataSizeBytes uint64) (*address.Address, error) {
	url := "https://api.filrep.io/api/v1/miners?"
	url += "sortBy=score"
	url += "&limit=100"

	resp, err := c.httpClient.Get(url)
	if err != nil {
		return nil, fmt.Errorf("httpClient.Get error: %w url %s", err, url)
	}
	defer resp.Body.Close()

	var result filrepMinersResult

	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		return nil, fmt.Errorf("Miners decode error: %w", err)
	}

	for _, m := range result.Miners {
		price, err := strconv.ParseUint(m.Price, 10, 64)
		if err != nil {
			continue
		}

		minPieceSize, err := strconv.ParseUint(m.MinPieceSize, 10, 64)
		if err != nil {
			continue
		}

		maxPieceSize, err := strconv.ParseUint(m.MaxPieceSize, 10, 64)
		if err != nil {
			continue
		}

		freeSpace, err := strconv.ParseUint(m.FreeSpace, 10, 64)
		if err != nil {
			continue
		}

		dataStored, err := strconv.ParseUint(m.FreeSpace, 10, 64)
		if err != nil {
			continue
		}

		switch {
		case !m.Status:
			continue
		case m.Reachability != "reachable":
			continue
		case m.UptimeAverage < 0.99:
			continue
		case price > c.dealsMaxPrice:
			continue
		case dataSizeBytes < minPieceSize || dataSizeBytes > maxPieceSize:
			continue
		case m.Score != "100":
			continue
		case freeSpace < (1 << 40): // 1TB
			continue
		case m.StorageDeals.Total < 100:
			continue
		case m.StorageDeals.Total != m.StorageDeals.NoPenalties:
			continue
		case dataStored < (1 << 40): // 1TB
			continue
		}

		addr, err := address.NewFromString(m.Address)
		if err != nil {
			log.Println("Filecoin findMiner address parse error:", err)
			continue
		}

		return &addr, nil
	}

	return nil, fmt.Errorf("%w: No eligible miner was found", errors.ErrNotFound)
}
