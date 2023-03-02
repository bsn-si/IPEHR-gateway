package client

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/bsn-si/IPEHR-gateway/src/pkg/docs/model"
	"github.com/stretchr/testify/assert"
)

func TestAQLQueryServiceClient_ExecuteQuery(t *testing.T) {
	query := &model.QueryRequest{}
	resp := &model.QueryResponse{
		Name:  "name",
		Query: "query",
		Columns: []model.QueryColumn{
			{
				Name: "name", Path: "path",
			},
		},
	}

	tests := []struct {
		name    string
		handler func(w http.ResponseWriter, r *http.Request)
		want    *model.QueryResponse
		wantErr bool
	}{
		{
			"1. path not found error",
			http.NotFound,
			nil,
			true,
		},
		{
			"2. timeout",
			func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusRequestTimeout)
				_ = json.NewEncoder(w).Encode(map[string]string{"error": "request timeout error"})
			},
			nil,
			true,
		},
		{
			"3. success",
			func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
				_ = json.NewEncoder(w).Encode(resp)
			},
			resp,
			false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			srv := httptest.NewServer(http.HandlerFunc(tt.handler))
			defer srv.Close()

			cli := NewAQLQueryServiceClient(srv.URL)
			got, err := cli.ExecQuery(context.Background(), query)
			if (err != nil) != tt.wantErr {
				t.Errorf("AQLQueryServiceClient.ExecuteQuery() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			assert.Equal(t, tt.want, got)
		})
	}
}
