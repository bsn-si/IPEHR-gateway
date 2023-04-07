package filecoin

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"time"

	"github.com/filecoin-project/go-address"
	"github.com/filecoin-project/go-fil-markets/retrievalmarket"
	"github.com/filecoin-project/go-fil-markets/storagemarket"
	jsonrpc "github.com/filecoin-project/go-jsonrpc"
	lotusapi "github.com/filecoin-project/lotus/api"
	"github.com/filecoin-project/lotus/chain/types"
	"github.com/ipfs/go-cid"

	"github.com/bsn-si/IPEHR-gateway/src/pkg/errors"
)

type DealStatus = storagemarket.StorageDealStatus

type Client struct {
	rpcEndpoint   string
	baseURL       string
	authToken     string
	dealsMaxPrice uint64
	miners        []string
	api           *lotusapi.FullNodeStruct
	closer        jsonrpc.ClientCloser
	httpClient    *http.Client
}

type Config struct {
	LotusRPCEndpoint string
	BaseURL          string
	AuthToken        string
	DealsMaxPrice    uint64
	Miners           []string
}

// nolint
type filrepMinersResult struct {
	Miners []struct {
		ID              uint
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
			Uptime                 interface{}
			StorageDeals           string
			CommittedSectorsProofs interface{}
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
		rpcEndpoint:   cfg.LotusRPCEndpoint,
		baseURL:       cfg.BaseURL,
		authToken:     cfg.AuthToken,
		dealsMaxPrice: cfg.DealsMaxPrice,
		miners:        cfg.Miners,
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

func (c *Client) StartDeal(ctx context.Context, CID *cid.Cid, dataSizeBytes uint64) (*cid.Cid, string, error) {
	walletAddr, err := c.api.WalletDefaultAddress(ctx)
	if err != nil {
		return nil, "", fmt.Errorf("Lotus WalletDefaultAddress error: %w", err)
	}

	// MinerAddress
	minerAddr, err := c.FindMiner(ctx)
	if err != nil {
		return nil, "", fmt.Errorf("Miner address parsing error: %w", err)
	}

	// ClientStartDeal
	deal, err := c.api.ClientStartDeal(ctx, &lotusapi.StartDealParams{
		Data: &storagemarket.DataRef{
			TransferType: storagemarket.TTGraphsync,
			Root:         *CID,
		},
		Wallet:            walletAddr,
		Miner:             *minerAddr,
		EpochPrice:        types.NewInt(c.dealsMaxPrice / 1e3), // TODO get from miner ask
		MinBlocksDuration: 518400,                              // epoch = 30 sec, 2880 per day, 180 days * 2880 = 518400
		//DealStartEpoch:    200,
		VerifiedDeal:  true,
		FastRetrieval: true,
		//ProviderCollateral big.Int
	})
	if err != nil {
		return nil, "", fmt.Errorf("Lotus.API.ClientStartDeal error: %w", err)
	}

	return deal, minerAddr.String(), nil
}

func (c *Client) GetDealStatus(ctx context.Context, CID *cid.Cid) (storagemarket.StorageDealStatus, uint64, error) {
	dealInfo, err := c.api.ClientGetDealInfo(ctx, *CID)
	if err != nil {
		return 0, 0, fmt.Errorf("Lotus ClientGetDealInfo error: %w CID %s", err, CID.String())
	}

	return dealInfo.State, uint64(dealInfo.DealID), nil
}

func (c *Client) StartRetrieve(ctx context.Context, CID *cid.Cid) (retrievalmarket.DealID, error) {
	offers, err := c.api.ClientFindData(ctx, *CID, nil)
	if err != nil {
		return 0, fmt.Errorf("Lotus ClientFindData error: %w CID: %s", err, CID.String())
	}

	if len(offers) == 0 {
		return 0, fmt.Errorf("%w ClientFindData offers is empty. dataCid: %s", errors.ErrCustom, CID.String())
	}

	walletAddr, err := c.api.WalletDefaultAddress(ctx)
	if err != nil {
		return 0, fmt.Errorf("Lotus WalletDefaultAddress error: %w", err)
	}

	// DEBUG
	log.Printf("offers: %+v", offers)

	// TODO why array?
	order := offers[0].Order(walletAddr)

	restrievalRes, err := c.api.ClientRetrieve(ctx, order)
	if err != nil {
		return 0, fmt.Errorf("Lotus ClientRetrieve error: %w", err)
	}

	retrieveStatus, err := c.GetRetrieveStatus(ctx, restrievalRes.DealID)
	if err != nil {
		return 0, fmt.Errorf("Lotus GetRetrieveStatus error: %w", err)
	}

	switch retrieveStatus {
	case
		retrievalmarket.DealStatusFailing,
		retrievalmarket.DealStatusRejected,
		retrievalmarket.DealStatusDealNotFound,
		retrievalmarket.DealStatusErrored,
		retrievalmarket.DealStatusInsufficientFunds:
		return 0, fmt.Errorf("%w Lotus Retrieve deal %d status not good: %s", errors.ErrCustom, restrievalRes.DealID, retrieveStatus)
	default: // ok
	}

	return restrievalRes.DealID, nil
}

func (c *Client) GetRetrieveStatus(ctx context.Context, dealID retrievalmarket.DealID) (retrievalmarket.DealStatus, error) {
	retrievalInfo, err := c.api.ClientListRetrievals(ctx)
	if err != nil {
		return 0, fmt.Errorf("Lotus ClientListRetrievals error: %w", err)
	}

	for _, ri := range retrievalInfo {
		if ri.ID == dealID {
			//log.Printf("retrievalInfo: %+v", ri)
			return ri.Status, nil
		}
	}

	return 0, errors.ErrNotFound
}

func (c *Client) SaveFile(ctx context.Context, CID *cid.Cid, dealID retrievalmarket.DealID) error {
	eref := lotusapi.ExportRef{
		Root:   *CID,
		DealID: dealID,
	}

	retRef := lotusapi.FileRef{
		Path: "/tmp/lotus/files/" + CID.String(), // TODO
	}

	err := c.api.ClientExport(ctx, eref, retRef)
	if err != nil {
		if err.Error() != "path already exists and overwriting is not allowed" {
			return fmt.Errorf("Lotus ClientExport error: %w", err)
		}
	}

	return nil
}

func (c *Client) Close() {
	c.closer()
}

/* Doesn't work
func (c *Client) FindMiner(ctx context.Context, dataSizeBytes uint64) (*address.Address, error) {
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
		case m.UptimeAverage < 0.9:
			continue
		case price > c.dealsMaxPrice:
			continue
		case dataSizeBytes < minPieceSize || dataSizeBytes > maxPieceSize:
			continue
		case m.Score != "100":
			continue
		case freeSpace < (1 << 40): // 1TB
			continue
		case m.StorageDeals.Total < 5:
			continue
		case m.StorageDeals.Total != m.StorageDeals.NoPenalties:
			continue
		case dataStored < (1 << 40): // 1TB
			continue
		}

		minerAddr, err := address.NewFromString(m.Address)
		if err != nil {
			log.Println("Filecoin findMiner address parse error:", err)
			continue
		}

		minerInfo, err := c.api.StateMinerInfo(ctx, minerAddr, types.EmptyTSK)
		if err != nil {
			//log.Println("Lotus api.StateMinerInfo error:", err)
			continue
		}

		cctx, cancel := context.WithTimeout(ctx, 3*time.Second)

		_, err = c.api.ClientQueryAsk(cctx, *minerInfo.PeerId, minerAddr)
		if err != nil {
			cancel()
			//log.Println("Lotus api.ClientQueryAsk error:", err)
			continue
		}

		cancel()

		return &minerAddr, nil
	}

	return nil, fmt.Errorf("%w: No eligible miner was found", errors.ErrNotFound)
}
*/

func (c *Client) FindMiner(ctx context.Context) (*address.Address, error) {
	blackList := []string{
		"",
		"f01482290",
		"f01497836",
		"f01543586",
		"f01602479",
		"f01606849",
		"f01606675",
		"f01386984",
		"f01579009",
		"f01611097",
		"f01423116",
		"f01310564",
		"f01394448",
		"f01672748",
		"f01443744",
	}
	_ = blackList

	rand.Seed(time.Now().UnixNano())

	for i := 0; i < 10; i++ {
		// nolint
		x := rand.Intn(len(c.miners))

		addr := c.miners[x]

		minerAddr, err := address.NewFromString(addr)
		if err != nil {
			log.Printf("Filecoin findMiner address parse error: %v addr: %s", err, addr)
			continue
		}

		minerInfo, err := c.api.StateMinerInfo(ctx, minerAddr, types.EmptyTSK)
		if err != nil {
			log.Println("Lotus api.StateMinerInfo error:", err)
			continue
		}

		cctx, cancel := context.WithTimeout(ctx, 10*time.Second)

		_, err = c.api.ClientQueryAsk(cctx, *minerInfo.PeerId, minerAddr)
		if err != nil {
			cancel()
			log.Println("Lotus api.ClientQueryAsk error:", err)

			continue
		}

		cancel()

		return &minerAddr, nil
	}

	return nil, fmt.Errorf("%w: No eligible miner was found", errors.ErrNotFound)
}

func (c *Client) BaseURL() string {
	return c.baseURL
}
