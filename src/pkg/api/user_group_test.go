package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/google/go-cmp/cmp"
	"github.com/google/uuid"

	"hms/gateway/pkg/api/mocks"
	"hms/gateway/pkg/user/model"
)

func TestUserHandler_GroupCreate(t *testing.T) {
	var (
		userID   = uuid.New().String()
		systemID = uuid.New().String()
		groupID  = uuid.New()
		urlPath  = "/v1/user/group"
	)

	ug := model.UserGroup{
		GroupID:     &groupID,
		Name:        "testGroupName",
		Description: "testDescription",
	}

	ugBytes, _ := json.Marshal(ug)

	tests := []struct {
		name       string
		groupName  string
		body       []byte
		prepare    func(userSvc *mocks.MockUserService)
		wantStatus int
	}{
		{
			"1. bad request because request body is empty",
			ug.Name,
			nil,
			func(svc *mocks.MockUserService) {},
			http.StatusBadRequest,
		},
		{
			"2. bad request because name is empty",
			ug.Name,
			nil,
			func(svc *mocks.MockUserService) {},
			http.StatusBadRequest,
		},
		{
			"3. success result",
			"",
			ugBytes,
			func(svc *mocks.MockUserService) {
				svc.EXPECT().GroupCreate(
					gomock.Any(),
					userID,
					systemID,
					gomock.Any(),
					ug.Name,
					ug.Description,
				).Return(ug.GroupID, nil)
			},
			http.StatusCreated,
		},
	}

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	userSvc := mocks.NewMockUserService(ctrl)

	api := API{
		User: NewUserHandler(userSvc),
	}

	router := api.setupRouter(api.buildUserAPI())

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			userSvc.EXPECT().VerifyAccess(gomock.Any(), gomock.Any()).Return(nil)
			tt.prepare(userSvc)

			reqBody := bytes.NewReader(tt.body)
			req := httptest.NewRequest(http.MethodPost, urlPath, reqBody)
			req.Header.Set("Authorization", "Bearer AccessKey")
			req.Header.Set("AuthUserId", userID)
			req.Header.Set("EhrSystemId", systemID)

			recorder := httptest.NewRecorder()
			router.ServeHTTP(recorder, req)

			resp := recorder.Result()
			defer resp.Body.Close()

			if diff := cmp.Diff(tt.wantStatus, resp.StatusCode); diff != "" {
				t.Errorf("UserHandler.Post() status code mismatch {-want;+got}\n\t%s", diff)
			}

			respBody, _ := io.ReadAll(resp.Body)
			if tt.wantStatus == http.StatusBadRequest {
				t.Logf(string(respBody))
				return
			}
		})
	}
}

func TestUserHandler_GroupGetByID(t *testing.T) {
	var (
		systemID = uuid.New().String()
		userID   = uuid.New().String()
		groupID  = uuid.New()
	)

	ug := &model.UserGroup{
		GroupID:     &groupID,
		Name:        "testName",
		Description: "testDescription",
		Members:     []string{},
	}

	ugJSON, _ := json.Marshal(ug)

	tests := []struct {
		name       string
		groupID    string
		prepare    func(userSvc *mocks.MockUserService)
		wantStatus int
		wantResp   string
	}{
		{
			"1. error because {group_id} is incorrect",
			"incorrect",
			func(svc *mocks.MockUserService) {},
			http.StatusBadRequest,
			`{"error":"group_id must be UUID"}`,
		},
		{
			"2. success result",
			ug.GroupID.String(),
			func(svc *mocks.MockUserService) {
				svc.EXPECT().GroupGetByID(gomock.Any(), userID, &groupID).Return(ug, nil)
			},
			http.StatusOK,
			string(ugJSON),
		},
	}

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	userSvc := mocks.NewMockUserService(ctrl)

	api := API{
		User: NewUserHandler(userSvc),
	}

	router := api.setupRouter(api.buildUserAPI())

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			userSvc.EXPECT().VerifyAccess(userID, "Bearer emptyJWTkey").Return(nil)
			tt.prepare(userSvc)

			req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/v1/user/group/%s", tt.groupID), nil)
			req.Header.Set("Authorization", "Bearer emptyJWTkey")
			req.Header.Set("AuthUserId", userID)
			req.Header.Set("EhrSystemId", systemID)

			recorder := httptest.NewRecorder()
			router.ServeHTTP(recorder, req)

			resp := recorder.Result()
			defer resp.Body.Close()

			if diff := cmp.Diff(tt.wantStatus, resp.StatusCode); diff != "" {
				t.Errorf("UserHandler.GroupGetByID() status code mismatch {-want;+got}\n\t%s", diff)
			}

			if tt.wantStatus == http.StatusNotFound {
				return
			}

			respBody, _ := io.ReadAll(resp.Body)
			if diff := cmp.Diff(tt.wantResp, string(respBody)); diff != "" {
				t.Errorf("UserHandler.GroupGetByID() status response {-want;+got}\n%s", diff)
			}
		})
	}
}
