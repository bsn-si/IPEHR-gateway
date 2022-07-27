package ipfs_test

import (
	"io/ioutil"
	"testing"

	"hms/gateway/pkg/config"
	"hms/gateway/pkg/storage/ipfs"
)

func TestAddFile(t *testing.T) {
	cfg, err := config.New()
	if err != nil {
		t.Fatal(err)
	}

	//expectedCid := "QmPYKPZhu6LdLrZJUbmUTPFCogwmmenaKMH5XMsrEBNG3m"
	fileContent := []byte("dfgg dtghreyh .sm,dfdsoiqwuefbw3586 (!!!) test one")

	testIpfsClient := ipfs.NewClient(cfg.Storage.Ipfs.EndpointURL)

	cid, err := testIpfsClient.Add(fileContent)
	if err != nil {
		t.Fatal(err)
	}

	fileCloser, err := testIpfsClient.Get(cid)
	if err != nil {
		t.Fatal(err)
	}
	defer fileCloser.Close()

	fileContent2, err := ioutil.ReadAll(fileCloser)
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
