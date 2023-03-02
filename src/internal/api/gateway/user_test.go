package gateway

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/bsn-si/IPEHR-gateway/src/internal/api/gateway/mocks"
	"github.com/bsn-si/IPEHR-gateway/src/pkg/errors"
	"github.com/bsn-si/IPEHR-gateway/src/pkg/user/model"
	"github.com/bsn-si/IPEHR-gateway/src/pkg/user/roles"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestUserHandler_Info(t *testing.T) {
	var (
		userID   = uuid.New().String()
		systemID = uuid.New().String()
	)

	tests := []struct {
		name       string
		userID     string
		prepare    func(svc *mocks.MockUserService)
		wantStatus int
		want       string
	}{
		{
			"1. error on get user info",
			userID,
			func(svc *mocks.MockUserService) {
				svc.EXPECT().Info(gomock.Any(), userID, systemID).Return(nil, errors.New("some error"))
			},
			http.StatusInternalServerError,
			"",
		},
		{
			"2. error not found",
			userID,
			func(svc *mocks.MockUserService) {
				svc.EXPECT().Info(gomock.Any(), userID, systemID).Return(nil, errors.ErrNotFound)
			},
			http.StatusNotFound,
			"",
		},
		{
			"3. success return doctor",
			userID,
			func(svc *mocks.MockUserService) {
				doctor := &model.UserInfo{
					Role:        roles.Doctor.String(),
					TimeCreated: "123",
					Name:        "some_name",
				}
				svc.EXPECT().Info(gomock.Any(), userID, systemID).Return(doctor, nil)
			},
			http.StatusOK,
			`{"role":"Doctor","name":"some_name","timeCreated":"123"}`,
		},
		{
			"4. success return patient",
			userID,
			func(svc *mocks.MockUserService) {
				ehrID := uuid.MustParse("09b8e497-c9d8-4562-99a7-a2f614037971")
				user := &model.UserInfo{
					Role:        roles.Patient.String(),
					TimeCreated: "123",
					Name:        "some_name",
					EhrID:       &ehrID,
				}
				svc.EXPECT().Info(gomock.Any(), userID, systemID).Return(user, nil)
			},
			http.StatusOK,
			`{"role":"Patient","timeCreated":"123","ehrID":"09b8e497-c9d8-4562-99a7-a2f614037971"}`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			userSvc := mocks.NewMockUserService(ctrl)
			userSvc.EXPECT().VerifyAccess(userID, "Bearer AccessKey").Return(nil)
			tt.prepare(userSvc)

			api := API{
				User: NewUserHandler(userSvc),
			}

			router := api.setupRouter(api.buildUserAPI())

			req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/v1/user/%s", tt.userID), nil)
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
