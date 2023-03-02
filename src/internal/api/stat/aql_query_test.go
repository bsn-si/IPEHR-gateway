package stat

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/bsn-si/IPEHR-gateway/src/internal/api/stat/mocks"
	"github.com/bsn-si/IPEHR-gateway/src/pkg/docs/model"
	"github.com/bsn-si/IPEHR-gateway/src/pkg/errors"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

//go:generate mockgen --package mocks --source aql_query.go --destination ./mocks/aql_query_mock.go

func TestAQLQueryAPI_QueryHandler(t *testing.T) {
	t.Parallel()

	queryReqStr := `{"q":"SELECT 1 FROM EHR e", "offset":1, "fetch":10, "query_parameters":{"key":1}}`
	queryReq := &model.QueryRequest{
		Query:  "SELECT 1 FROM EHR e",
		Offset: 1,
		Fetch:  10,
		QueryParameters: map[string]interface{}{
			"key": 1.0,
		},
	}

	tests := []struct {
		name     string
		data     string
		prepare  func(qm *mocks.MockAQLQuerier)
		wantCode int
		wantBody string
	}{
		{
			"1. invalid json",
			"invalid JSON string",
			func(qm *mocks.MockAQLQuerier) {},
			http.StatusBadRequest,
			`{"error":"invalid body"}`,
		},
		{
			"2. query validdation error",
			"{}",
			func(qm *mocks.MockAQLQuerier) {},
			http.StatusBadRequest,
			`{"error":"request validation error"}`,
		},
		{
			"3. error on execute query",
			queryReqStr,
			func(qm *mocks.MockAQLQuerier) {
				qm.EXPECT().ExecQuery(gomock.Any(), queryReq).Return(nil, errors.New("some error"))
			},
			http.StatusInternalServerError,
			`{"error":"internal server error"}`,
		},
		{
			"4. timeout error on execute query",
			queryReqStr,
			func(qm *mocks.MockAQLQuerier) {
				qm.EXPECT().ExecQuery(gomock.Any(), queryReq).Return(nil, errors.ErrTimeout)
			},
			http.StatusRequestTimeout,
			`{"error":"timeout exceeded"}`,
		},
		{
			"5. success",
			queryReqStr,
			func(qm *mocks.MockAQLQuerier) {
				resp := &model.QueryResponse{}
				qm.EXPECT().ExecQuery(gomock.Any(), queryReq).Return(resp, nil)
			},
			200,
			`{"meta":{"_href":"","_type":"","_schema_version":"","_created":"","_generator":"","_executed_aql":""},"name":"","q":"","columns":null,"rows":null}`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			queryMock := mocks.NewMockAQLQuerier(ctrl)
			tt.prepare(queryMock)

			api := &API{
				queryAPI: newAQLQueryAPI(queryMock),
			}

			w := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodPost, "/query/", bytes.NewBuffer([]byte(tt.data)))
			api.setupRouter(api.buildQueryAPI()).ServeHTTP(w, req)

			assert.Equal(t, tt.wantCode, w.Code)
			assert.Equal(t, tt.wantBody, w.Body.String())
		})
	}
}
