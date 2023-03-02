package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/bsn-si/IPEHR-gateway/src/pkg/common"
	"github.com/bsn-si/IPEHR-gateway/src/pkg/docs/model"
	"github.com/bsn-si/IPEHR-gateway/src/pkg/errors"
)

const executeQueryPath = `/query/`

type AQLQueryServiceClient struct {
	statsServiceHost string
	client           *http.Client
}

func NewAQLQueryServiceClient(statsHost string) *AQLQueryServiceClient {
	cli := &http.Client{
		Timeout: common.WebRequestTimeout,
	}

	return &AQLQueryServiceClient{
		statsServiceHost: statsHost,
		client:           cli,
	}
}

func (cli *AQLQueryServiceClient) ExecQuery(ctx context.Context, query *model.QueryRequest) (*model.QueryResponse, error) {
	reqData, err := json.Marshal(query)
	if err != nil {
		return nil, errors.Wrap(err, "cannot marshal request body")
	}

	url := fmt.Sprintf("%s%s", cli.statsServiceHost, executeQueryPath)

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewBuffer(reqData))
	if err != nil {
		return nil, errors.Wrap(err, "cannot create new request")
	}

	resp, err := cli.client.Do(req)
	if err != nil {
		return nil, errors.Wrap(err, "cannot do http request")
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		errResp := map[string]string{}
		if err := json.NewDecoder(resp.Body).Decode(&errResp); err != nil {
			log.Printf("[ERROR] cannot unmarshal execute AQL error body: %v", err)
		} else {
			log.Printf("[ERROR] execute AQL query response error: %v", errResp["error"])
		}

		switch resp.StatusCode {
		case http.StatusRequestTimeout:
			return nil, errors.ErrTimeout
		default:
			return nil, errors.ErrInternalServerError
		}
	}

	result := &model.QueryResponse{}
	if err := json.NewDecoder(resp.Body).Decode(result); err != nil {
		return nil, errors.Wrap(err, "cannot unmarshal respose body")
	}

	return result, nil
}
