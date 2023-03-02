package gateway

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/bsn-si/IPEHR-gateway/src/internal/api/gateway/mocks"
	"github.com/bsn-si/IPEHR-gateway/src/pkg/docs/model"
	"github.com/bsn-si/IPEHR-gateway/src/pkg/docs/model/base"
	"github.com/bsn-si/IPEHR-gateway/src/pkg/errors"
	"github.com/golang/mock/gomock"
	"github.com/google/go-cmp/cmp"
	"github.com/google/uuid"
)

//
//go:generate mockgen -source query.go -package mocks -destination mocks/query_mock.go
//go:generate mockgen -source user.go -package mocks -destination mocks/user_mock.go
//

func TestStoredQueryHandler_List(t *testing.T) {
	var (
		userID   = uuid.New().String()
		systemID = uuid.New().String()
		sqM      = make([]*model.StoredQuery, 1)
	)

	sqM[0] = &model.StoredQuery{
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
			"1. empty result because qualified_query_name was not found",
			"notexist",
			func(gaSvc *mocks.MockQueryService) {
				gaSvc.EXPECT().List(gomock.Any(), userID, systemID, "notexist")
			},
			http.StatusOK,
			`[]`,
		},
		{
			"2. success result",
			"exist",
			func(gaSvc *mocks.MockQueryService) {
				gaSvc.EXPECT().List(gomock.Any(), userID, systemID, "exist").Return(sqM, nil)
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
			userSvc := mocks.NewMockUserService(ctrl)
			userSvc.EXPECT().VerifyAccess(gomock.Any(), gomock.Any()).Return(nil)

			api := API{
				Query: NewQueryHandler(sqSvc, ""),
				User:  NewUserHandler(userSvc),
			}

			router := api.setupRouter(api.buildDefinitionAPI())

			req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/v1/definition/query/%s", tt.qualifiedQueryName), nil)
			req.Header.Set("Authorization", "Bearer emptyJWTkey")
			req.Header.Set("AuthUserId", userID)
			req.Header.Set("EhrSystemId", systemID)

			recorder := httptest.NewRecorder()
			router.ServeHTTP(recorder, req)

			resp := recorder.Result()
			defer resp.Body.Close()

			if diff := cmp.Diff(tt.wantStatus, resp.StatusCode); diff != "" {
				t.Errorf("StoredQueryHandler.List() status code mismatch {-want;+got}\n\t%s", diff)
			}

			if tt.wantStatus == http.StatusNotFound {
				return
			}

			respBody, _ := io.ReadAll(resp.Body)
			if diff := cmp.Diff(tt.wantResp, string(respBody)); diff != "" {
				t.Errorf("StoredQueryHandler.List() status response {-want;+got}\n%s", diff)
			}
		})
	}
}

func TestStoredQueryHandler_GetByVersion(t *testing.T) {
	var (
		userID   = uuid.New().String()
		systemID = uuid.New().String()
	)

	sqM := &model.StoredQuery{
		Name:        "org.openehr::compositions",
		Type:        "aql",
		Version:     "1.0.1",
		TimeCreated: "2017-07-16T19:20:30.450+01:00",
		Query:       "SELECT 1",
	}

	qV, err := base.NewVersionTreeID(sqM.Version)
	if err != nil {
		t.Fatal(err)
	}

	sqJSON, _ := json.Marshal(sqM)

	tests := []struct {
		name               string
		qualifiedQueryName string
		version            string
		prepare            func(gaSvc *mocks.MockQueryService)
		wantStatus         int
		wantResp           string
	}{
		{
			"1. empty result because {qualified_query_name} is incorrect",
			"",
			"incorrect",
			func(gaSvc *mocks.MockQueryService) {},
			http.StatusBadRequest,
			"",
		},
		{
			"2. empty result because {version} is incorrect",
			sqM.Name,
			"incorrect",
			func(gaSvc *mocks.MockQueryService) {},
			http.StatusBadRequest,
			"",
		},
		{
			"3. empty result because {qualified_query_name} was not found",
			"notexist",
			sqM.Version,
			func(gaSvc *mocks.MockQueryService) {
				gaSvc.EXPECT().GetByVersion(gomock.Any(), userID, systemID, "notexist", qV).Return(nil, errors.ErrNotFound)
			},
			http.StatusNotFound,
			"",
		},
		{
			"4. empty result because document with that version is not exist",
			sqM.Name,
			"999.999.999",
			func(gaSvc *mocks.MockQueryService) {
				gaSvc.EXPECT().GetByVersion(gomock.Any(), userID, systemID, sqM.Name, gomock.Any()).Return(nil, errors.ErrNotFound)
			},
			http.StatusNotFound,
			"",
		},
		{
			"5. success result",
			sqM.Name,
			sqM.Version,
			func(gaSvc *mocks.MockQueryService) {
				gaSvc.EXPECT().GetByVersion(gomock.Any(), userID, systemID, sqM.Name, qV).Return(sqM, nil)
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
			userSvc := mocks.NewMockUserService(ctrl)
			userSvc.EXPECT().VerifyAccess(gomock.Any(), gomock.Any()).Return(nil)

			api := API{
				Query: NewQueryHandler(sqSvc, ""),
				User:  NewUserHandler(userSvc),
			}

			router := api.setupRouter(api.buildDefinitionAPI())

			req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/v1/definition/query/%s/%s", tt.qualifiedQueryName, tt.version), nil)
			req.Header.Set("Authorization", "Bearer emptyJWTkey")
			req.Header.Set("AuthUserId", userID)
			req.Header.Set("EhrSystemId", systemID)

			recorder := httptest.NewRecorder()
			router.ServeHTTP(recorder, req)

			resp := recorder.Result()
			defer resp.Body.Close()

			if diff := cmp.Diff(tt.wantStatus, resp.StatusCode); diff != "" {
				t.Errorf("StoredQueryHandler.List() status code mismatch {-want;+got}\n\t%s", diff)
			}

			if tt.wantStatus == http.StatusNotFound {
				return
			}

			respBody, _ := io.ReadAll(resp.Body)
			if diff := cmp.Diff(tt.wantResp, string(respBody)); diff != "" {
				t.Errorf("StoredQueryHandler.GetStoredByVersion() status response {-want;+got}\n%s", diff)
			}
		})
	}
}

func TestStoredQueryHandler_Put(t *testing.T) {
	var (
		userID   = uuid.New().String()
		systemID = uuid.New().String()
		urlPath  = "/v1/definition/query"
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
			sqM.Name,
			"",
			func(gaSvc *mocks.MockQueryService) {},
			http.StatusBadRequest,
			"",
		},
		{
			"2. failed because query is incorrect",
			sqM.Name,
			sqM.Query,
			func(gaSvc *mocks.MockQueryService) {
				gaSvc.EXPECT().Validate([]byte(sqM.Query)).Return(false)
			},
			http.StatusBadRequest,
			urlPath + "/" + sqM.Name + "/" + sqM.Version,
		},
		{
			"3. success result",
			sqM.Name,
			sqM.Query,
			func(gaSvc *mocks.MockQueryService) {
				gaSvc.EXPECT().Validate([]byte(sqM.Query)).Return(true)
				gaSvc.EXPECT().Store(
					gomock.Any(),
					userID,
					systemID,
					gomock.Any(),
					model.QueryTypeAQL,
					sqM.Name,
					sqM.Query,
				).Return(&sqM, nil)
			},
			http.StatusOK,
			urlPath + "/" + sqM.Name + "/" + sqM.Version,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			sqSvc := mocks.NewMockQueryService(ctrl)
			tt.prepare(sqSvc)

			// Mock for auth user service
			userSvc := mocks.NewMockUserService(ctrl)
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
			req.Header.Set("EhrSystemId", systemID)

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

func TestStoredQueryHandler_PutByVer(t *testing.T) {
	var (
		userID   = uuid.New().String()
		systemID = uuid.New().String()
		urlPath  = "/v1/definition/query"
	)

	sqM := model.StoredQuery{
		Name:        "org.openehr::compositions",
		Type:        "aql",
		Version:     "1.0.1",
		TimeCreated: "2017-07-16T19:20:30.450+01:00",
		Query:       "SELECT 1",
	}

	qV, err := base.NewVersionTreeID(sqM.Version)
	if err != nil {
		t.Fatal(err)
	}

	tests := []struct {
		name               string
		qualifiedQueryName string
		body               string
		prepare            func(gaSvc *mocks.MockQueryService)
		wantVersion        string
		wantStatus         int
		wantLocation       string
	}{
		{
			"1. bad request because request body is empty",
			sqM.Name,
			"",
			func(gaSvc *mocks.MockQueryService) {
				gaSvc.EXPECT().GetByVersion(
					gomock.Any(),
					userID,
					systemID,
					sqM.Name,
					qV,
				).Return(nil, nil)
			},
			sqM.Version,
			http.StatusBadRequest,
			"",
		},
		{
			"2. bad request because version is not valid",
			sqM.Name,
			sqM.Query,
			func(gaSvc *mocks.MockQueryService) {},
			"notvalidversion.12.4",
			http.StatusBadRequest,
			"",
		},
		{
			"3. bad request because version of query already exist",
			sqM.Name,
			sqM.Query,
			func(gaSvc *mocks.MockQueryService) {
				gaSvc.EXPECT().GetByVersion(
					gomock.Any(),
					userID,
					systemID,
					sqM.Name,
					qV,
				).Return(&sqM, nil)
			},
			sqM.Version,
			http.StatusConflict,
			urlPath + "/" + sqM.Name + "/" + sqM.Version,
		},
		{
			"4. failed because query is incorrect",
			sqM.Name,
			sqM.Query,
			func(gaSvc *mocks.MockQueryService) {
				gaSvc.EXPECT().GetByVersion(
					gomock.Any(),
					userID,
					systemID,
					sqM.Name,
					qV,
				).Return(nil, nil)
				gaSvc.EXPECT().Validate([]byte(sqM.Query)).Return(false)
			},
			sqM.Version,
			http.StatusBadRequest,
			urlPath + "/" + sqM.Name + "/" + sqM.Version,
		},
		{
			"5. success result",
			sqM.Name,
			sqM.Query,
			func(gaSvc *mocks.MockQueryService) {
				gaSvc.EXPECT().GetByVersion(
					gomock.Any(),
					userID,
					systemID,
					sqM.Name,
					qV,
				).Return(nil, nil)

				gaSvc.EXPECT().Validate([]byte(sqM.Query)).Return(true)

				gaSvc.EXPECT().StoreVersion(
					gomock.Any(),
					userID,
					systemID,
					gomock.Any(),
					model.QueryTypeAQL,
					sqM.Name,
					qV,
					sqM.Query,
				).Return(&sqM, nil)
			},
			sqM.Version,
			http.StatusOK,
			urlPath + "/" + sqM.Name + "/" + sqM.Version,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			sqSvc := mocks.NewMockQueryService(ctrl)
			tt.prepare(sqSvc)

			// Mock for auth user service
			userSvc := mocks.NewMockUserService(ctrl)
			userSvc.EXPECT().VerifyAccess(gomock.Any(), gomock.Any()).Return(nil)

			api := API{
				Query: NewQueryHandler(sqSvc, ""),
				User:  NewUserHandler(userSvc),
			}

			router := api.setupRouter(api.buildDefinitionAPI())

			reqBody := strings.NewReader(tt.body)

			req := httptest.NewRequest(http.MethodPut, fmt.Sprintf("%s/%s/%s", urlPath, tt.qualifiedQueryName, tt.wantVersion), reqBody)
			req.Header.Set("Authorization", "Bearer AccessKey")
			req.Header.Set("AuthUserId", userID)
			req.Header.Set("EhrSystemId", systemID)

			recorder := httptest.NewRecorder()
			router.ServeHTTP(recorder, req)

			resp := recorder.Result()
			defer resp.Body.Close()

			if diff := cmp.Diff(tt.wantStatus, resp.StatusCode); diff != "" {
				t.Errorf("StoredQueryHandler.Put() status code mismatch {-want;+got}\n\t%s", diff)
			}

			if resp.StatusCode != http.StatusOK {
				return
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
