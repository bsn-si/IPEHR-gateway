package api

import (
	"bytes"
	"encoding/json"
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
		urlPath  = "/v1/user/group"
	)

	ug := model.UserGroup{
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
					gomock.Any(),
					gomock.Any(),
					gomock.Any(),
					gomock.Any(),
					gomock.Any(),
				).Return(nil)
			},
			http.StatusCreated,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			userSvc := mocks.NewMockUserService(ctrl)
			userSvc.EXPECT().VerifyAccess(gomock.Any(), gomock.Any()).Return(nil)
			tt.prepare(userSvc)

			api := API{
				User: NewUserHandler(userSvc),
			}

			router := api.setupRouter(api.buildUserAPI())

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
