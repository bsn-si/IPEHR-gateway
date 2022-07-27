package ipfs_test

import (
	"hms/gateway/pkg/config"
	"hms/gateway/pkg/storage/ipfs"
	"testing"
)

func TestAddFile(t *testing.T) {
	cfg, err := config.New()
	if err != nil {
		t.Fatal(err)
	}

	expectedCid := "QmPYKPZhu6LdLrZJUbmUTPFCogwmmenaKMH5XMsrEBNG3m"
	fileContent := []byte("dfgg dtghreyh .sm,dfdsoiqwuefbw3586 (!!!) test one")

	testIpfsClient := ipfs.NewClient(cfg.Storage.Ipfs.EndpointURL)

	cid, err := testIpfsClient.Add(fileContent)
	if err != nil {
		t.Fatal(err)
	}

	if cid != expectedCid {
		t.Fatalf("Received wrong CID (expected, received): %s, %s", expectedCid, cid)
	}
}
