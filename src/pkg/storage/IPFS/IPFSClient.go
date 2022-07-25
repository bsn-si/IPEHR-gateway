package IPFS

import (
	"encoding/json"
	"hms/gateway/pkg/storage/IPFS/HTTPClient"
)

type ipfsAddResult struct {
	Name string
	Hash string
	Size string
}

type Client struct {
	apiURL     string
	httpClient HTTPClient.Interface
}

func New(apiURL string, httpClient HTTPClient.Interface) *Client {
	return &Client{
		apiURL:     apiURL,
		httpClient: httpClient,
	}
}

// Add file to an IPFS node
func (i *Client) Add(fileContent *[]byte) (cid string, err error) {
	res, err := i.httpClient.PostFile(i.apiURL+"/add", fileContent)
	if err != nil {
		return
	}

	var resJSON ipfsAddResult

	err = json.Unmarshal(res, &resJSON)
	if err != nil {
		return
	}

	cid = resJSON.Hash

	return
}
