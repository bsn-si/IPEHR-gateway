package api

import (
	"bytes"
	"encoding/json"
	"errors"
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

//go:generate mockgen -source group_access.go -package mocks -destination mocks/group_access_mock.go

func TestGroupAccessHandler_Create(t *testing.T) {
	const baseURL = "/"
	userID := uuid.New().String()
	ehrSystemID := "some_ehr_system_id"
	groupUUID := uuid.New()
	ga := &model.GroupAccess{
		Description: "desc",
		GroupUUID:   &groupUUID,
	}
	gaJson, _ := json.Marshal(ga)

	tests := []struct {
		name       string
		prepare    func(gaSvc *mocks.MockGroupAccessService)
		data       []byte
		wantStatus int
		wantResp   string
	}{
		{
			"1. invalid request body",
			func(gaSvc *mocks.MockGroupAccessService) {},
			nil,
			http.StatusBadRequest,
			`{"error":"Request decoding error"}`,
		},
		{
			"2. request validation error",
			func(gaSvc *mocks.MockGroupAccessService) {
				gaSvc.EXPECT().Create(gomock.Any(), userID, &model.GroupAccessCreateRequest{}).
					Return(nil, errors.New("some error"))
			},
			[]byte("{}"),
			http.StatusInternalServerError,
			`{"error":"Group Access creating error"}`,
		},
		{
			"3. success",
			func(gaSvc *mocks.MockGroupAccessService) {

				gaSvc.EXPECT().Create(gomock.Any(), userID, &model.GroupAccessCreateRequest{}).
					Return(ga, nil)
			},
			[]byte("{}"),
			http.StatusCreated,
			string(gaJson),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			gaSvc := mocks.NewMockGroupAccessService(ctrl)
			tt.prepare(gaSvc)

			api := API{
				GroupAccess: NewGroupAccessHandler(gaSvc, baseURL),
				testMode:    true,
			}
			router := api.setupRouter(api.buildGroupAccessAPI())

			req := httptest.NewRequest(http.MethodPost, "/v1/access/group", bytes.NewBuffer(tt.data))
			req.Header.Set("AuthUserId", userID)
			req.Header.Set("EhrSystemId", ehrSystemID)

			recorder := httptest.NewRecorder()
			router.ServeHTTP(recorder, req)

			resp := recorder.Result()
			if diff := cmp.Diff(tt.wantStatus, resp.StatusCode); diff != "" {
				t.Errorf("GroupAccessHandler.Create() status code mismatch {-want;+got}\n\t%s", diff)
			}

			respBody, _ := io.ReadAll(resp.Body)
			if diff := cmp.Diff(tt.wantResp, string(respBody)); diff != "" {
				t.Errorf("GroupAccessHandler.Create() status response {-want;+got}\n\t%s", diff)
			}
		})
	}
}

func TestGroupAccessHandler_Get(t *testing.T) {
	const baseURL = "/"
	var (
		userID  = uuid.New().String()
		groupID = uuid.New()

		ehrSystemID = "some_ehr_system_id"
	)
	ga := &model.GroupAccess{
		Description: "desc",
		GroupUUID:   &groupID,
	}
	gaJson, _ := json.Marshal(ga)

	tests := []struct {
		name       string
		groupID    string
		prepare    func(gaSvc *mocks.MockGroupAccessService)
		wantStatus int
		wantResp   string
	}{
		{
			"1. error on parse group_id",
			"invalid_uuid",
			func(gaSvc *mocks.MockGroupAccessService) {},
			http.StatusBadRequest,
			`{"error":"groupId is incorrect"}`,
		},
		{
			"2. error on get access group",
			groupID.String(),
			func(gaSvc *mocks.MockGroupAccessService) {
				gaSvc.EXPECT().Get(gomock.Any(), userID, &groupID).Return(nil, errors.New("some error"))
			},
			http.StatusNotFound,
			`{"error":"Group access not found"}`,
		},
		{
			"3. success get access group",
			groupID.String(),
			func(gaSvc *mocks.MockGroupAccessService) {
				gaSvc.EXPECT().Get(gomock.Any(), userID, &groupID).Return(ga, nil)
			},
			http.StatusOK,
			string(gaJson),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			gaSvc := mocks.NewMockGroupAccessService(ctrl)
			tt.prepare(gaSvc)

			api := API{
				GroupAccess: NewGroupAccessHandler(gaSvc, baseURL),
				testMode:    true,
			}
			router := api.setupRouter(api.buildGroupAccessAPI())

			req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/v1/access/group/%s", tt.groupID), nil)
			req.Header.Set("AuthUserId", userID)
			req.Header.Set("EhrSystemId", ehrSystemID)

			recorder := httptest.NewRecorder()
			router.ServeHTTP(recorder, req)

			resp := recorder.Result()
			if diff := cmp.Diff(tt.wantStatus, resp.StatusCode); diff != "" {
				t.Errorf("GroupAccessHandler.Get() status code mismatch {-want;+got}\n\t%s", diff)
			}

			respBody, _ := io.ReadAll(resp.Body)
			if diff := cmp.Diff(tt.wantResp, string(respBody)); diff != "" {
				t.Errorf("GroupAccessHandler.Get() status response {-want;+got}\n\t%s", diff)
			}
		})
	}
}
