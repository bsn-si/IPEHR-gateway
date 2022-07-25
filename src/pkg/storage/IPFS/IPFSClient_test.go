package IPFS_test

import (
	"hms/gateway/pkg/config"
	"hms/gateway/pkg/storage/IPFS"
	"hms/gateway/pkg/storage/IPFS/HTTPClientMock"
	"testing"
)

func TestAddFile(t *testing.T) {
	cfg, err := config.New()
	if err != nil {
		t.Fatal(err)
	}

	expectedCid := "QmPYKPZhu6LdLrZJUbmUTPFCogwmmenaKMH5XMsrEBNG3m"

	fileContent := []byte("dfgg dtghreyh .sm,dfdsoiqwuefbw3586 (!!!) test one")

	// Uncomment the line below if You want to test with real IPFS node
	//testHttpClient := httpClient.New()
	testHTTPClient := HTTPClientMock.New()
	testHTTPClient.SetPostRes([]byte(`{"Name":"file.txt","Hash":"` + expectedCid + `","Size":"58"}`))
	testIpfsClient := IPFS.New(cfg.IpfsNodeAPI, testHTTPClient)

	cid, err := testIpfsClient.Add(&fileContent)
	if err != nil {
		t.Fatal(err)
	}

	if cid != expectedCid {
		t.Fatalf("Received wrong CID (expected, received): %s, %s", expectedCid, cid)
	}
}
