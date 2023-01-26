package api

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/bsn-si/IPEHR-gateway/src/pkg/api/mocks"
	"github.com/bsn-si/IPEHR-gateway/src/pkg/docs/model"
	"github.com/bsn-si/IPEHR-gateway/src/pkg/errors"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestQueryHandler_ExecStoredQuery(t *testing.T) {
	var (
		queryName = "some_aql_query"
		userID    = "5d44b88c-4199-4bad-97dc-d78268e01398"
		systemID  = "6d44b88c-4199-4bad-97dc-d78268e01398"
		ehrID     = uuid.MustParse("7d44b88c-4199-4bad-97dc-d78268e01398")
	)

	tests := []struct {
		name        string
		queryParams string
		prepare     func(svc *mocks.MockQueryService)
		wantStatus  int
		want        string
	}{
		{
			"1. invalid ehr_id",
			"ehr_id=invalid_ehr",
			func(svc *mocks.MockQueryService) {},
			400,
			`{"error":"ehr_id bad format"}`,
		},
		{
			"2. invalid offset",
			"offset=invalid_offset",
			func(svc *mocks.MockQueryService) {},
			400,
			`{"error":"offset bad format"}`,
		},
		{
			"3. invalid limit",
			"fetch=invalid",
			func(svc *mocks.MockQueryService) {},
			400,
			`{"error":"fetch bad format"}`,
		},
		{
			"4. error on get data",
			"ehr_id=7d44b88c-4199-4bad-97dc-d78268e01398&offset=1&fetch=10&some_val1=1&some_val2=some_str",
			func(svc *mocks.MockQueryService) {
				r := &model.QueryRequest{
					Offset: 1,
					Fetch:  10,
					QueryParameters: map[string]interface{}{
						"ehr_id":    ehrID,
						"some_val1": "1",
						"some_val2": "some_str",
					},
				}

				svc.EXPECT().ExecStoredQuery(gomock.Any(), userID, systemID, queryName, r).Return(nil, errors.New("some error"))
			},
			500,
			`{"error":"internal server error"}`,
		},
		{
			"5. success",
			"ehr_id=7d44b88c-4199-4bad-97dc-d78268e01398&offset=1&fetch=10&some_val1=1&some_val2=some_str",
			func(svc *mocks.MockQueryService) {
				r := &model.QueryRequest{
					Offset: 1,
					Fetch:  10,
					QueryParameters: map[string]interface{}{
						"ehr_id":    ehrID,
						"some_val1": "1",
						"some_val2": "some_str",
					},
				}
				resp := &model.QueryResponse{}

				svc.EXPECT().ExecStoredQuery(gomock.Any(), userID, systemID, queryName, r).Return(resp, nil)
			},
			200,
			`{"meta":{"_href":"","_type":"","_schema_version":"","_created":"","_generator":"","_executed_aql":""},"name":"","q":"","columns":null,"rows":null}`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			userSvc := mocks.NewMockUserService(ctrl)
			querySvc := mocks.NewMockQueryService(ctrl)

			userSvc.EXPECT().VerifyAccess(userID, "Bearer AccessKey").Return(nil)

			tt.prepare(querySvc)

			api := API{
				User:  NewUserHandler(userSvc),
				Query: NewQueryHandler(querySvc, "base_url"),
			}

			router := api.setupRouter(api.buildQueryAPI())

			req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/v1/query/%s?%s", queryName, tt.queryParams), nil)
			req.Header.Set("Authorization", "Bearer AccessKey")
			req.Header.Set("AuthUserId", userID)
			req.Header.Set("EhrSystemId", systemID)

			recorder := httptest.NewRecorder()
			router.ServeHTTP(recorder, req)

			resp := recorder.Result()
			defer resp.Body.Close()

			assert.Equal(t, tt.wantStatus, resp.StatusCode)

			respBody, _ := io.ReadAll(resp.Body)
			assert.Equal(t, tt.want, string(respBody))
		})
	}
}

func TestQueryHandler_PostExecStoredQuery(t *testing.T) {
	var (
		queryName = "some_aql_query"
		userID    = "5d44b88c-4199-4bad-97dc-d78268e01398"
		systemID  = "6d44b88c-4199-4bad-97dc-d78268e01398"
	)

	tests := []struct {
		name       string
		body       []byte
		prepare    func(svc *mocks.MockQueryService)
		wantStatus int
		want       string
	}{
		{
			"1. invalid offset",
			[]byte(`{"offset":"invalid_offset"}`),
			func(svc *mocks.MockQueryService) {},
			400,
			`{"error":"body bad format"}`,
		},
		{
			"2. invalid limit",
			[]byte(`{"fetch":"invalid"}`),
			func(svc *mocks.MockQueryService) {},
			400,
			`{"error":"body bad format"}`,
		},
		{
			"3. error on get data",
			[]byte(`{"offset":1,"fetch":10,"query_parameters":{"key":1}}`),
			func(svc *mocks.MockQueryService) {
				r := &model.QueryRequest{
					Offset: 1,
					Fetch:  10,
					QueryParameters: map[string]interface{}{
						"key": 1.0,
					},
				}

				svc.EXPECT().ExecStoredQuery(gomock.Any(), userID, systemID, queryName, r).Return(nil, errors.New("some error"))
			},
			500,
			`{"error":"internal server error"}`,
		},
		{
			"4. success",
			[]byte(`{"offset":1,"fetch":10,"query_parameters":{"key":1}}`),
			func(svc *mocks.MockQueryService) {
				r := &model.QueryRequest{
					Offset: 1,
					Fetch:  10,
					QueryParameters: map[string]interface{}{
						"key": 1.0,
					},
				}
				resp := &model.QueryResponse{}

				svc.EXPECT().ExecStoredQuery(gomock.Any(), userID, systemID, queryName, r).Return(resp, nil)
			},
			200,
			`{"meta":{"_href":"","_type":"","_schema_version":"","_created":"","_generator":"","_executed_aql":""},"name":"","q":"","columns":null,"rows":null}`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			userSvc := mocks.NewMockUserService(ctrl)
			querySvc := mocks.NewMockQueryService(ctrl)

			userSvc.EXPECT().VerifyAccess(userID, "Bearer AccessKey").Return(nil)

			tt.prepare(querySvc)

			api := API{
				User:  NewUserHandler(userSvc),
				Query: NewQueryHandler(querySvc, "base_url"),
			}

			router := api.setupRouter(api.buildQueryAPI())

			req := httptest.NewRequest(http.MethodPost, fmt.Sprintf("/v1/query/%s", queryName), bytes.NewBuffer(tt.body))
			req.Header.Set("Authorization", "Bearer AccessKey")
			req.Header.Set("AuthUserId", userID)
			req.Header.Set("EhrSystemId", systemID)

			recorder := httptest.NewRecorder()
			router.ServeHTTP(recorder, req)

			resp := recorder.Result()
			defer resp.Body.Close()

			assert.Equal(t, tt.wantStatus, resp.StatusCode)

			respBody, _ := io.ReadAll(resp.Body)
			assert.Equal(t, tt.want, string(respBody))
		})
	}
}

func TestQueryHandler_ExecPostQuery(t *testing.T) {
	var (
		userID   = "5d44b88c-4199-4bad-97dc-d78268e01398"
		systemID = "6d44b88c-4199-4bad-97dc-d78268e01398"
	)

	tests := []struct {
		name       string
		body       []byte
		prepare    func(svc *mocks.MockQueryService)
		wantStatus int
		want       string
	}{
		{
			"1. invalid offset",
			[]byte(`{"offset":"invalid_offset"}`),
			func(svc *mocks.MockQueryService) {},
			http.StatusBadRequest,
			`{"error":"body bad format"}`,
		},
		{
			"2. invalid limit",
			[]byte(`{"fetch":"invalid"}`),
			func(svc *mocks.MockQueryService) {},
			http.StatusBadRequest,
			`{"error":"body bad format"}`,
		},
		{
			"3. false because query is empty",
			[]byte(`{"offset":1,"fetch":10,"query_parameters":{"key":1}}`),
			func(svc *mocks.MockQueryService) {},
			http.StatusBadRequest,
			`{"error":"Request validation error"}`,
		},
		{
			"4. error on get data",
			[]byte(`{"q":"SELECT 1", "offset":1, "fetch":10, "query_parameters":{"key":1}}`),
			func(svc *mocks.MockQueryService) {
				r := &model.QueryRequest{
					Query:  "SELECT 1",
					Offset: 1,
					Fetch:  10,
					QueryParameters: map[string]interface{}{
						"key": 1.0,
					},
				}

				svc.EXPECT().ExecQuery(gomock.Any(), r).Return(nil, errors.New("some error"))
			},
			http.StatusInternalServerError,
			`{"error":"internal server error"}`,
		},
		{
			"5. error on get data because long request",
			[]byte(`{"q":"SELECT 1", "offset":1, "fetch":10, "query_parameters":{"key":1}}`),
			func(svc *mocks.MockQueryService) {
				r := &model.QueryRequest{
					Query:  "SELECT 1",
					Offset: 1,
					Fetch:  10,
					QueryParameters: map[string]interface{}{
						"key": 1.0,
					},
				}

				svc.EXPECT().ExecQuery(gomock.Any(), r).Return(nil, errors.ErrTimeout)
			},
			http.StatusRequestTimeout,
			`{"error":"timeout exceeded"}`,
		},
		{
			"5. success",
			[]byte(`{"q":"SELECT 1", "offset":1, "fetch":10, "query_parameters":{"key":1}}`),
			func(svc *mocks.MockQueryService) {
				r := &model.QueryRequest{
					Query:  "SELECT 1",
					Offset: 1,
					Fetch:  10,
					QueryParameters: map[string]interface{}{
						"key": 1.0,
					},
				}
				resp := &model.QueryResponse{}

				svc.EXPECT().ExecQuery(gomock.Any(), r).Return(resp, nil)
			},
			200,
			`{"meta":{"_href":"","_type":"","_schema_version":"","_created":"","_generator":"","_executed_aql":""},"name":"","q":"","columns":null,"rows":null}`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			userSvc := mocks.NewMockUserService(ctrl)
			querySvc := mocks.NewMockQueryService(ctrl)

			userSvc.EXPECT().VerifyAccess(userID, "Bearer AccessKey").Return(nil)

			tt.prepare(querySvc)

			api := API{
				User:  NewUserHandler(userSvc),
				Query: NewQueryHandler(querySvc, "base_url"),
			}

			router := api.setupRouter(api.buildQueryAPI())

			req := httptest.NewRequest(http.MethodPost, "/v1/query/aql", bytes.NewBuffer(tt.body))
			req.Header.Set("Authorization", "Bearer AccessKey")
			req.Header.Set("AuthUserId", userID)
			req.Header.Set("EhrSystemId", systemID)

			recorder := httptest.NewRecorder()
			router.ServeHTTP(recorder, req)

			resp := recorder.Result()
			defer resp.Body.Close()

			assert.Equal(t, tt.wantStatus, resp.StatusCode)

			respBody, _ := io.ReadAll(resp.Body)
			assert.Equal(t, tt.want, string(respBody))
		})
	}
}
