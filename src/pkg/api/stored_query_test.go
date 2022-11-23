package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/google/go-cmp/cmp"
	"github.com/google/uuid"

	"hms/gateway/pkg/api/mocks"
	"hms/gateway/pkg/docs/model"
)

//
//go:generate mockgen -source query.go -package mocks -destination mocks/query_mock.go
//go:generate mockgen -source user.go -package mocks -destination mocks/user_mock.go
//

func TestStoredQueryHandler_Get(t *testing.T) {
	var (
		userID = uuid.New().String()
		sqM    = make([]model.StoredQuery, 1)
	)

	sqM[0] = model.StoredQuery{
		Name:        "org.openehr::compositions",
		Type:        "aql",
		Version:     "1.0.1",
		TimeCreated: "2017-07-16T19:20:30.450+01:00",
		Query:       "SELECT 1",
	}

	sqJSON, _ := json.Marshal(sqM)

	tests := []struct {
		name               string
		qualifiedQueryName string
		prepare            func(gaSvc *mocks.MockQueryService)
		wantStatus         int
		wantResp           string
	}{
		{
			"1. empty result because qualifiedQueryName was not found",
			"notexist",
			func(gaSvc *mocks.MockQueryService) {
				gaSvc.EXPECT().Get(gomock.Any(), gomock.Any(), gomock.Any())
			},
			http.StatusOK,
			`[]`,
		},
		{
			"2. success result",
			"exist",
			func(gaSvc *mocks.MockQueryService) {
				gaSvc.EXPECT().Get(gomock.Any(), gomock.Any(), gomock.Any()).Return(sqM, nil)
			},
			http.StatusOK,
			string(sqJSON),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			sqSvc := mocks.NewMockQueryService(ctrl)
			tt.prepare(sqSvc)

			// Mock for auth user service
			userSvc := mocks.NewMockUserHandlerService(ctrl)
			userSvc.EXPECT().VerifyAccess(gomock.Any(), gomock.Any()).Return(nil)

			api := API{
				Query: NewQueryHandler(sqSvc, ""),
				User:  NewUserHandler(userSvc),
			}

			router := api.setupRouter(api.buildDefinitionAPI())

			req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/v1/definition/query/%s", tt.qualifiedQueryName), nil)
			req.Header.Set("Authorization", "Bearer emptyJWTkey")
			req.Header.Set("AuthUserId", userID)

			recorder := httptest.NewRecorder()
			router.ServeHTTP(recorder, req)

			resp := recorder.Result()
			defer resp.Body.Close()

			if diff := cmp.Diff(tt.wantStatus, resp.StatusCode); diff != "" {
				t.Errorf("StoredQueryHandler.Get() status code mismatch {-want;+got}\n\t%s", diff)
			}

			if tt.wantStatus == http.StatusNotFound {
				return
			}

			respBody, _ := io.ReadAll(resp.Body)
			if diff := cmp.Diff(tt.wantResp, string(respBody)); diff != "" {
				t.Errorf("StoredQueryHandler.Get() status response {-want;+got}\n%s", diff)
			}
		})
	}
}

func TestStoredQueryHandler_Put(t *testing.T) {
	var (
		userID  = uuid.New().String()
		urlPath = "/v1/definition/query"
	)

	sqM := model.StoredQuery{
		Name:        "org.openehr::compositions",
		Type:        "aql",
		Version:     "1.0.1",
		TimeCreated: "2017-07-16T19:20:30.450+01:00",
		Query:       "SELECT 1",
	}

	tests := []struct {
		name               string
		qualifiedQueryName string
		body               string
		prepare            func(gaSvc *mocks.MockQueryService)
		wantStatus         int
		wantLocation       string
	}{
		{
			"1. bad request because request body is empty",
			sqM.Name.String(),
			"",
			func(gaSvc *mocks.MockQueryService) {},
			http.StatusBadRequest,
			"",
		},
		{
			"2. success result",
			sqM.Name.String(),
			sqM.Query,
			func(gaSvc *mocks.MockQueryService) {
				gaSvc.EXPECT().Validate(gomock.Any()).Return(true)
				gaSvc.EXPECT().Store(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(sqM, nil)
			},
			http.StatusOK,
			urlPath + "/" + sqM.Name.String() + "/" + sqM.Version,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			sqSvc := mocks.NewMockQueryService(ctrl)
			tt.prepare(sqSvc)

			// Mock for auth user service
			userSvc := mocks.NewMockUserHandlerService(ctrl)
			userSvc.EXPECT().VerifyAccess(gomock.Any(), gomock.Any()).Return(nil)

			api := API{
				Query: NewQueryHandler(sqSvc, ""),
				User:  NewUserHandler(userSvc),
			}

			router := api.setupRouter(api.buildDefinitionAPI())

			reqBody := strings.NewReader(tt.body)

			req := httptest.NewRequest(http.MethodPut, fmt.Sprintf("%s/%s", urlPath, tt.qualifiedQueryName), reqBody)
			req.Header.Set("Authorization", "Bearer AccessKey")
			req.Header.Set("AuthUserId", userID)

			recorder := httptest.NewRecorder()
			router.ServeHTTP(recorder, req)

			resp := recorder.Result()
			defer resp.Body.Close()

			if diff := cmp.Diff(tt.wantStatus, resp.StatusCode); diff != "" {
				t.Errorf("StoredQueryHandler.Put() status code mismatch {-want;+got}\n\t%s", diff)
			}

			respBody, _ := io.ReadAll(resp.Body)
			if tt.wantStatus == http.StatusBadRequest {
				t.Logf(string(respBody))
				return
			}

			if diff := cmp.Diff(tt.wantLocation, resp.Header.Get("Location")); diff != "" {
				t.Errorf("StoredQueryHandler.Put() status response {-want;+got}\n%s", diff)
			}
		})
	}
}
