package filecoin_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/ipfs/go-cid"

	"hms/gateway/pkg/config"
	"hms/gateway/pkg/storage/filecoin"
)

func TestStartDeal(t *testing.T) {
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

	dataSize := uint64(1000000)

	minerAddress, err := filecoinClient.FindMiner(dataSize)
	if err != nil {
		t.Fatal(err)
	}

	t.Log("Miner address:", minerAddress)
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
