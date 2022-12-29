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
	"hms/gateway/pkg/common"
	"hms/gateway/pkg/docs/model"
	"hms/gateway/pkg/docs/model/base"
	"hms/gateway/pkg/docs/service/processing"
	"hms/gateway/pkg/errors"
	userModel "hms/gateway/pkg/user/model"
)

//
//go:generate mockgen -source directory.go -package mocks -destination mocks/directory_mock.go
//go:generate mockgen -source user.go -package mocks -destination mocks/user_mock.go
//go:generate mockgen -source ../helper/finder.go -package mocks -destination mocks/finder_mock.go
//

func TestDirectoryHandler_Create(t *testing.T) {
	var (
		userUUID    = uuid.New()
		ehrUUID     = uuid.New()
		directoryID = uuid.New()
		systemID    = common.EhrSystemID
	)

	directoryVersionUID, _ := base.NewObjectVersionID(directoryID.String(), systemID)

	tests := []struct {
		name        string
		directoryID string
		payload     []byte
		prepare     func(dSvc *mocks.MockDirectoryService, iSvc *mocks.MockIndexer, uS *mocks.MockUserService)
		wantStatus  int
		wantResp    *model.Directory
	}{
		{
			"1. failed because EhrID is not belong to current user",
			directoryVersionUID.String(),
			[]byte(``),
			func(_ *mocks.MockDirectoryService, i *mocks.MockIndexer, _ *mocks.MockUserService) {
				i.EXPECT().GetEhrUUIDByUserID(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, errors.ErrNotFound)
			},
			http.StatusNotFound,
			nil,
		},
		{
			"2. failed because body is empty",
			directoryVersionUID.String(),
			[]byte(``),
			func(_ *mocks.MockDirectoryService, i *mocks.MockIndexer, _ *mocks.MockUserService) {
				i.EXPECT().GetEhrUUIDByUserID(gomock.Any(), gomock.Any(), gomock.Any()).Return(&ehrUUID, errors.ErrNotFound)
			},
			http.StatusNotFound,
			nil,
		},
		{
			"3. failed because DIRECTORY with current ID already exist",
			directoryVersionUID.String(),
			[]byte(`{
					"_type": "FOLDER",
						"name": {
							"_type": "DV_TEXT",
							"value": "root"
					},
					"uid": {
						"_type": "OBJECT_VERSION_ID",
						"value": "` + directoryVersionUID.String() + `"
					
					},
					"archetype_node_id": "openEHR-EHR-FOLDER.generic.v1"
				}
			`),
			func(dS *mocks.MockDirectoryService, i *mocks.MockIndexer, _ *mocks.MockUserService) {
				i.EXPECT().GetEhrUUIDByUserID(gomock.Any(), gomock.Any(), gomock.Any()).Return(&ehrUUID, nil)
				dS.EXPECT().GetByID(gomock.Any(), userUUID.String(), directoryVersionUID.String()).Return(&model.Directory{}, nil)
			},
			http.StatusConflict,
			nil,
		},
		{
			"4. Success!!! NewProcRequest",
			directoryID.String(),
			[]byte(`{
					"_type": "FOLDER",
						"name": {
							"_type": "DV_TEXT",
							"value": "root"
					},
					"uid": {
						"_type": "OBJECT_VERSION_ID",
						"value": "` + directoryVersionUID.String() + `"
					
					},
					"archetype_node_id": "openEHR-EHR-FOLDER.generic.v1"
				}
			`),
			func(dS *mocks.MockDirectoryService, i *mocks.MockIndexer, uS *mocks.MockUserService) {
				i.EXPECT().GetEhrUUIDByUserID(gomock.Any(), gomock.Any(), gomock.Any()).Return(&ehrUUID, nil)
				dS.EXPECT().GetByID(gomock.Any(), userUUID.String(), directoryVersionUID.String()).Return(nil, nil)

				uS.EXPECT().Info(gomock.Any(), userUUID.String()).Return(&userModel.UserInfo{}, nil)

				dS.EXPECT().NewProcRequest(gomock.Any(), userUUID.String(), ehrUUID.String(), processing.RequestDirectoryCreate).Return(&MockRequest{}, nil)
				dS.EXPECT().Create(gomock.Any(), gomock.Any(), systemID, &ehrUUID, &userModel.UserInfo{}, gomock.Any()).Return(nil)
			},
			http.StatusCreated,
			nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			svc := mocks.NewMockDirectoryService(ctrl)
			indexer := mocks.NewMockIndexer(ctrl)
			userSvc := mocks.NewMockUserService(ctrl)
			tt.prepare(svc, indexer, userSvc)

			// Mock for auth user service
			userSvc.EXPECT().VerifyAccess(gomock.Any(), gomock.Any()).Return(nil)
			//userSvc.EXPECT().Info(gomock.Any(), userUUID.String()).Return(&userModel.UserInfo{}, nil)

			api := API{
				Directory: NewDirectoryHandler(svc, userSvc, indexer, ""),
				User:      NewUserHandler(userSvc),
			}

			router := api.setupRouter(api.buildEhrDirectoryAPI())

			req := httptest.NewRequest(http.MethodPost, fmt.Sprintf("/v1/ehr/%s/directory/", ehrUUID.String()), bytes.NewReader(tt.payload))
			req.Header.Set("Authorization", "Bearer emptyJWTkey")
			req.Header.Set("AuthUserId", userUUID.String())
			req.Header.Set("EhrSystemId", systemID)

			recorder := httptest.NewRecorder()
			router.ServeHTTP(recorder, req)

			resp := recorder.Result()
			defer resp.Body.Close()

			if diff := cmp.Diff(tt.wantStatus, resp.StatusCode); diff != "" {
				t.Errorf("DirectoryHandler.Create() status code mismatch {-want;+got}\n\t%s", diff)
				return
			}

			if tt.wantResp == nil {
				return
			}

			var got *model.Directory
			respBody, _ := io.ReadAll(resp.Body)

			err := json.Unmarshal(respBody, &got)
			if err != nil {
				t.Errorf("Cant unmarshal JSON: %v", err)
			}

			opts := cmp.AllowUnexported(
				base.ObjectVersionID{},
				base.PartyProxy{},
			)

			if diff := cmp.Diff(tt.wantResp, got, opts); diff != "" {
				t.Errorf("DirectoryHandler.Create() body response {-want;+got}\n%s", diff)
			}
		})
	}
}
