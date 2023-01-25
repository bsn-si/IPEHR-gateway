package ipfs_test

import (
	"context"
	"io"
	"testing"
	"time"

	"github.com/bsn-si/IPEHR-gateway/src/pkg/config"
	"github.com/bsn-si/IPEHR-gateway/src/pkg/storage/ipfs"

	"github.com/ipfs/go-cid"
)

func TestAddFile(t *testing.T) {
	cfg, err := config.New()
	if err != nil {
		t.Fatal(err)
	}

	//expectedCid := "QmPYKPZhu6LdLrZJUbmUTPFCogwmmenaKMH5XMsrEBNG3m"
	fileContent := []byte("dfgg dtghreyh .sm,dfdsoiqwuefbw3586 (!!!) test one")

	testIpfsClient, err := ipfs.NewClient(cfg.Storage.Ipfs.EndpointURLs)
	if err != nil {
		t.Fatal(err)
	}

	c, err := testIpfsClient.Add(context.Background(), fileContent)
	if err != nil {
		t.Fatal(err)
	}

	fileCloser, err := testIpfsClient.Get(context.Background(), c)
	if err != nil {
		t.Fatal(err)
	}
	defer fileCloser.Close()

	fileContent2, err := io.ReadAll(fileCloser)
	if err != nil {
		t.Fatal(err)
	}

	if len(fileContent) != len(fileContent2) {
		t.Fatalf("Received file length is wrong (expected, received): %d, %d", len(fileContent), len(fileContent2))
	}

	if string(fileContent) != string(fileContent2) {
		t.Fatalf("Received wrong file (expected, received): %s, %s", fileContent, fileContent2)
	}
}

func TestGetIncorrectCid(t *testing.T) {
	cfg, err := config.New()
	if err != nil {
		t.Fatal(err)
	}

	testIpfsClient, err := ipfs.NewClient(cfg.Storage.Ipfs.EndpointURLs)
	if err != nil {
		t.Fatal(err)
	}

	CID, err := cid.Parse("QmPYKPZhu6LdLrZJUbmUTPFCogwmmenaKMH5XMsrEBNG3n")
	if err != nil {
		t.Fatal(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	fileCloser, err := testIpfsClient.Get(ctx, &CID)
	if err == nil {
		if fileCloser != nil {
			fileCloser.Close()
		}

		t.Fatal("Expected error, received file")
	}
}
