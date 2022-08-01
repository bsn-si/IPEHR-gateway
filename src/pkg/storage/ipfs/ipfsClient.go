package ipfs

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"

	"github.com/ipfs/go-cid"

	"hms/gateway/pkg/errors"
)

type ipfsAddResult struct {
	Name string
	Hash string
	Size string
}

type ipfsVersionResult struct {
	Commit  string
	Golang  string
	Repo    string
	System  string
	Version string
}

type Client struct {
	apiURL     string
	httpClient *http.Client
	version    string
}

func NewClient(apiURL string) (*Client, error) {
	client := &Client{
		apiURL:     apiURL,
		httpClient: http.DefaultClient,
	}

	url := apiURL + "/version"
	resp, err := client.httpClient.Post(url, "", nil)
	if err != nil {
		return nil, fmt.Errorf("IPFS NewClient error: %w apiURL %s", err, apiURL)
	}
	defer resp.Body.Close()

	result := &ipfsVersionResult{}
	err = json.NewDecoder(resp.Body).Decode(result)
	if err != nil {
		return nil, fmt.Errorf("IPFS version response decode error: %w", err)
	}

	if result.Version == "" {
		return nil, fmt.Errorf("IPFS version error: %w", errors.ErrCustom)
	}

	client.version = result.Version

	return client, nil
}

// Add file to an IPFS node with CID version 0
// Returns CID or error
func (i *Client) Add(fileContent []byte) (*cid.Cid, error) {
	var (
		url             = i.apiURL + "/add?cid-version=0"
		requestBody     bytes.Buffer
		multiPartWriter = multipart.NewWriter(&requestBody)
		fileWriter, _   = multiPartWriter.CreateFormFile("file", "file.txt")
	)

	_, err := io.Copy(fileWriter, bytes.NewReader(fileContent))
	if err != nil {
		return nil, fmt.Errorf("io.Copy error: %w", err)
	}

	multiPartWriter.Close()

	req, err := http.NewRequest(http.MethodPost, url, &requestBody)
	if err != nil {
		return nil, fmt.Errorf("http.NewRequest add error: %w", err)
	}

	req.Header.Add("Content-Type", multiPartWriter.FormDataContentType())

	resp, err := i.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("IPFS add request error: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("IPFS add request error: %w status %s", errors.ErrCustom, resp.Status)
	}

	result := &ipfsAddResult{}
	err = json.NewDecoder(resp.Body).Decode(result)
	if err != nil {
		return nil, fmt.Errorf("IPFS add response decode error: %w", err)
	}

	CID, err := cid.Parse(result.Hash)
	if err != nil {
		return nil, fmt.Errorf("IPFS add response CID parse error: %w", err)
	}

	return &CID, nil
}

// Get file from IPFS node by CID
// Returns ReadCloser or error
// Need to Close()
func (i *Client) Get(CID *cid.Cid) (io.ReadCloser, error) {
	url := i.apiURL + "/cat?arg=" + CID.String()

	resp, err := i.httpClient.Post(url, "", nil)
	if err != nil {
		return nil, fmt.Errorf("IPFS get request error: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("IPFS get request error: %w status %s", errors.ErrCustom, resp.Status)
	}

	return resp.Body, nil
}
