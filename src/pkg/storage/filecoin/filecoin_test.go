package filecoin_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/ipfs/go-cid"

	"github.com/bsn-si/IPEHR-gateway/src/pkg/config"
	"github.com/bsn-si/IPEHR-gateway/src/pkg/storage/filecoin"
)

func TestStartDeal(t *testing.T) {
	t.Skip()

	filecoinClient, err := prepare(t)
	if err != nil {
		t.Fatal(err)
	}

	defer clean(t, filecoinClient)

	ctx := context.Background()

	CID, err := cid.Decode("QmPYKPZhu6LdLrZJUbmUTPFCogwmmenaKMH5XMsrEBNG3m")
	if err != nil {
		t.Fatal(err)
	}

	dataSize := uint64(50)

	dealCID, minerAddr, err := filecoinClient.StartDeal(ctx, &CID, dataSize)
	if err != nil {
		t.Fatal(err)
	}

	t.Log("dealCid:", dealCID, "minerAddress:", minerAddr)
}

func TestFindMiner(t *testing.T) {
	filecoinClient, err := prepare(t)
	if err != nil {
		t.Fatal(err)
	}

	defer clean(t, filecoinClient)

	minerAddress, err := filecoinClient.FindMiner(context.Background())
	if err != nil {
		t.Fatal(err)
	}

	t.Log("Miner address:", minerAddress)
}

func TestRetrieve(t *testing.T) {
	t.Skip()

	filecoinClient, err := prepare(t)
	if err != nil {
		t.Fatal(err)
	}

	defer clean(t, filecoinClient)

	ctx := context.Background()

	CID, err := cid.Decode("QmRYm3NtD5uDH1msh9YAXCYeuPQKFFJkr4cXtgWDfiuczM")
	if err != nil {
		t.Fatal(err)
	}

	dealID, err := filecoinClient.StartRetrieve(ctx, &CID)
	if err != nil {
		t.Fatal(err)
	}

	t.Log("Retrieval dealID:", dealID)

	retrieveStatus, err := filecoinClient.GetRetrieveStatus(ctx, dealID)
	if err != nil {
		t.Fatal(err)
	}

	t.Log("retrieveStatus:", retrieveStatus)
}

func prepare(t *testing.T) (*filecoin.Client, error) {
	t.Helper()

	cfg, err := config.New()
	if err != nil {
		return nil, err
	}

	filecoinCfg := (filecoin.Config)(cfg.Storage.Filecoin)

	client, err := filecoin.NewClient(&filecoinCfg)
	if err != nil {
		return nil, fmt.Errorf("filecoin.New error: %w", err)
	}

	return client, nil
}

func clean(t *testing.T, s *filecoin.Client) {
	t.Helper()

	s.Close()
}
