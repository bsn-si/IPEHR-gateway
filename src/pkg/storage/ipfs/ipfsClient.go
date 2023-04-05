package ipfs

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/ipfs/go-cid"

	"github.com/bsn-si/IPEHR-gateway/src/pkg/errors"
)

const (
	checkEndpointStatusInterval = 30 * time.Second
	active                      = "active"
	inactive                    = "inactive"
)

type (
	ipfsAddResult struct {
		Name string
		Hash string
		Size string
	}

	ipfsVersionResult struct {
		Commit  string
		Golang  string
		Repo    string
		System  string
		Version string
	}

	endpoint struct {
		APIURL    string
		Status    string
		LastCheck time.Time
	}

	Client struct {
		sync.Mutex
		endpoints  []*endpoint
		httpClient *http.Client
		done       chan bool
	}
)

func NewClient(apiURLs []string) (*Client, error) {
	client := &Client{
		httpClient: &http.Client{
			Timeout: time.Second * 10,
		},
	}

	for _, url := range apiURLs {
		if _, err := client.getVersion(url); err != nil {
			return nil, fmt.Errorf("IPFS endpoint get version error: %w", err)
		}

		client.endpoints = append(client.endpoints, &endpoint{
			APIURL:    url,
			Status:    active,
			LastCheck: time.Now(),
		})
	}

	ticker := time.NewTicker(checkEndpointStatusInterval)

	go func() {
		for {
			select {
			case <-client.done:
				return
			case <-ticker.C:
				client.checkEndpointStatus()
			}
		}
	}()

	return client, nil
}

func (i *Client) Close() {
	i.done <- true
}

// nolint
func (i *Client) getVersion(url string) (string, error) {
	url += "/version"

	resp, err := i.httpClient.Post(url, "", nil)
	if err != nil {
		return "", fmt.Errorf("IPFS get version error: %w URL %s", err, url)
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		resp.Body.Close()
		return "", fmt.Errorf("IPFS get version response error: %w url: %s", err, url)
	}

	resp.Body.Close()

	result := &ipfsVersionResult{}
	if err = json.Unmarshal(data, result); err != nil {
		return "", fmt.Errorf("IPFS version response decode error: %w URL: %s response: %s", err, url, data)
	}

	if result.Version == "" {
		return "", fmt.Errorf("%w IPFS get version error. URL: %s response: %s", errors.ErrCustom, url, result)
	}

	return result.Version, nil
}

// Add file to an IPFS node with CID version 0
// Returns CID or error
func (i *Client) Add(ctx context.Context, fileContent []byte) (*cid.Cid, error) {
	var (
		requestBody     bytes.Buffer
		multiPartWriter = multipart.NewWriter(&requestBody)
		fileWriter, _   = multiPartWriter.CreateFormFile("file", "file.txt")
	)

	i.Lock()
	defer i.Unlock()

	_, err := io.Copy(fileWriter, bytes.NewReader(fileContent))
	if err != nil {
		return nil, fmt.Errorf("io.Copy error: %w", err)
	}

	multiPartWriter.Close()

	for _, endpoint := range i.endpoints {
		if endpoint.Status != active {
			continue
		}

		url := endpoint.APIURL + "/add?cid-version=0"

		req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, &requestBody)
		if err != nil {
			return nil, fmt.Errorf("http.NewRequest add error: %w", err)
		}

		req.Header.Add("Content-Type", multiPartWriter.FormDataContentType())

		resp, err := i.httpClient.Do(req)
		if err != nil {
			log.Printf("[IPFS] add request error: %v URL: %s", err, url)

			endpoint.Status = inactive

			endpoint.LastCheck = time.Now()

			continue
		}

		if resp.StatusCode != http.StatusOK {
			log.Printf("[IPFS] add request status %s", resp.Status)

			endpoint.Status = inactive

			endpoint.LastCheck = time.Now()

			resp.Body.Close()

			continue
		}

		result := &ipfsAddResult{}

		if err = json.NewDecoder(resp.Body).Decode(result); err != nil {
			resp.Body.Close()
			return nil, fmt.Errorf("IPFS add response decode error: %w URL: %s", err, url)
		}

		CID, err := cid.Parse(result.Hash)
		if err != nil {
			resp.Body.Close()
			return nil, fmt.Errorf("IPFS add response CID parse error: %w URL: %s", err, url)
		}

		resp.Body.Close()

		return &CID, nil
	}

	return nil, fmt.Errorf("%w IPFS endpoints are not available", errors.ErrCustom)
}

// Get file from IPFS node by CID
// Returns ReadCloser or error
// Need to Close()
func (i *Client) Get(ctx context.Context, CID *cid.Cid) (io.ReadCloser, error) {
	for _, endpoint := range i.endpoints {
		if endpoint.Status != active {
			continue
		}

		url := endpoint.APIURL + "/cat?arg=" + CID.String()

		request, err := http.NewRequestWithContext(ctx, http.MethodPost, url, nil)
		if err != nil {
			return nil, fmt.Errorf("http.NewRequestWithContext error: %w", err)
		}

		resp, err := i.httpClient.Do(request)
		if err != nil {
			if !strings.Contains(err.Error(), "context deadline exceeded") {
				log.Printf("[IPFS] get request error: %v URL: %s", err, url)
			}

			continue
		}

		if resp.StatusCode != http.StatusOK {
			log.Printf("[IPFS] get request error: %v status %s", errors.ErrCustom, resp.Status)

			endpoint.Status = inactive

			endpoint.LastCheck = time.Now()

			resp.Body.Close()

			continue
		}

		return resp.Body, nil
	}

	return nil, errors.ErrNotFound
}

func (i *Client) checkEndpointStatus() {
	i.Lock()
	defer i.Unlock()

	for _, endpoint := range i.endpoints {
		if time.Since(endpoint.LastCheck) < checkEndpointStatusInterval {
			continue
		}

		_, err := i.getVersion(endpoint.APIURL)
		if err != nil {
			log.Printf("[IPFS] getVersion error: %v", err)

			endpoint.Status = inactive
		} else {
			endpoint.Status = active
		}

		endpoint.LastCheck = time.Now()
	}
}
