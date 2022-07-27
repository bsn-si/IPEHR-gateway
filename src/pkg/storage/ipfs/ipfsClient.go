package ipfs

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"

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

// Add file to an IPFS node
// Returns CID or error
func (i *Client) Add(fileContent []byte) (string, error) {
	var (
		url             = i.apiURL + "/add"
		requestBody     bytes.Buffer
		multiPartWriter = multipart.NewWriter(&requestBody)
		fileWriter, _   = multiPartWriter.CreateFormFile("file", "file.txt")
	)

	_, err := io.Copy(fileWriter, bytes.NewReader(fileContent))
	if err != nil {
		return "", fmt.Errorf("io.Copy error: %w", err)
	}

	multiPartWriter.Close()

	req, err := http.NewRequest(http.MethodPost, url, &requestBody)
	if err != nil {
		return "", fmt.Errorf("http.NewRequest error: %w", err)
	}

	req.Header.Add("Content-Type", multiPartWriter.FormDataContentType())

	resp, err := i.httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("IPFS add request error: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("IPFS add request error: %w status %s", errors.ErrCustom, resp.Status)
	}

	result := ipfsAddResult{}

	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		return "", fmt.Errorf("IPFS add response decode error: %w", err)
	}

	resp.Body.Close()

	return result.Hash, nil
}
