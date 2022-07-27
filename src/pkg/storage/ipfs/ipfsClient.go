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

type Client struct {
	apiURL     string
	httpClient *http.Client
}

func NewClient(apiURL string) *Client {
	return &Client{
		apiURL:     apiURL,
		httpClient: http.DefaultClient,
	}
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
