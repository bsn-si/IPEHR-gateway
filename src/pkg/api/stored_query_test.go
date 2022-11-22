package api

import (
	"encoding/json"
	"fmt"
	"hms/gateway/pkg/api/mocks"
	"hms/gateway/pkg/docs/model"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/google/go-cmp/cmp"
	"github.com/google/uuid"
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
			"1. empty result because no qualifiedQueryName",
			"",
			func(gaSvc *mocks.MockQueryService) {},
			http.StatusNotFound,
			"",
		},
		{
			"2. empty result because qualifiedQueryName was not found",
			"notexist",
			func(gaSvc *mocks.MockQueryService) {
				gaSvc.EXPECT().Get(gomock.Any(), gomock.Any(), gomock.Any())
			},
			http.StatusOK,
			`[]`,
		},
		{
			"3. success result",
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
				Query: NewQueryHandler(sqSvc),
				User:  NewUserHandler(userSvc),
			}

			router := api.setupRouter(api.buildStoredQueryAPI())

			req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/v1/definition/query/%s", tt.qualifiedQueryName), nil)
			req.Header.Set("Authorization", "Bearer emptyJWTkey")
			req.Header.Set("AuthUserId", userID)

			recorder := httptest.NewRecorder()
			router.ServeHTTP(recorder, req)

			resp := recorder.Result()
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
