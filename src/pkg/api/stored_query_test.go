package api

import (
	"fmt"
	"hms/gateway/pkg/api/mocks"
	"hms/gateway/pkg/docs/model"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/google/go-cmp/cmp"
)

//
//go:generate mockgen -source stored_query.go -package mocks -destination mocks/stored_query_mock.go
//

func TestStoredQueryHandler_Get(t *testing.T) {
	var (
		//userID = uuid.New().String()
		sqM = make([]model.StoredQuery, 1)
	)

	sqM[0] = model.StoredQuery{
		Name:        "",
		Type:        "",
		Version:     "",
		TimeCreated: "",
		Query:       "",
	}

	//sqJson, _ := json.Marshal(sqM)

	tests := []struct {
		name               string
		qualifiedQueryName string
		prepare            func(gaSvc *mocks.MockStoredQueryService)
		wantStatus         int
		wantResp           string
	}{
		// TODO repair tests
		{
			"1. empty result because no qualifiedQueryName",
			"",
			func(gaSvc *mocks.MockStoredQueryService) {},
			http.StatusNotFound,
			"",
		},
		//{
		//	"2. error on get access group",
		//	groupID.String(),
		//	func(gaSvc *mocks.MockGroupAccessService) {
		//		gaSvc.EXPECT().Get(gomock.Any(), userID, &groupID).Return(nil, errors.New("some error"))
		//	},
		//	http.StatusNotFound,
		//	`{"error":"Group access not found"}`,
		//},
		//{
		//	"3. success get access group",
		//	groupID.String(),
		//	func(gaSvc *mocks.MockGroupAccessService) {
		//		gaSvc.EXPECT().Get(gomock.Any(), userID, &groupID).Return(ga, nil)
		//	},
		//	http.StatusOK,
		//	string(sqJson),
		//},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			sqSvc := mocks.NewMockStoredQueryService(ctrl)
			tt.prepare(sqSvc)

			api := API{
				StoredQuery: NewStoredQueryHandler(sqSvc),
			}

			// TODO add auth mock
			router := api.setupRouter(api.buildStoredQueryAPI())

			req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/v1/definition/query/%s", tt.qualifiedQueryName), nil)
			//req.Header.Set("AuthUserId", userID) // TODO ???

			recorder := httptest.NewRecorder()
			router.ServeHTTP(recorder, req)

			resp := recorder.Result()
			if diff := cmp.Diff(tt.wantStatus, resp.StatusCode); diff != "" {
				t.Errorf("StoredQueryHandler.Get() status code mismatch {-want;+got}\n\t%s", diff)
			}

			respBody, _ := io.ReadAll(resp.Body)
			if diff := cmp.Diff(tt.wantResp, string(respBody)); diff != "" {
				t.Errorf("StoredQueryHandler.Get() status response {-want;+got}\n%s", diff)
			}
		})
	}
}
