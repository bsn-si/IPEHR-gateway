package gateway

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
	"time"

	"github.com/bsn-si/IPEHR-gateway/src/internal/api/gateway/mocks"
	"github.com/bsn-si/IPEHR-gateway/src/pkg/common"
	"github.com/bsn-si/IPEHR-gateway/src/pkg/docs/model"
	"github.com/bsn-si/IPEHR-gateway/src/pkg/docs/model/base"
	"github.com/bsn-si/IPEHR-gateway/src/pkg/docs/service/processing"
	"github.com/bsn-si/IPEHR-gateway/src/pkg/errors"
	"github.com/golang/mock/gomock"
	"github.com/google/go-cmp/cmp"
	"github.com/google/uuid"
)

//
//go:generate mockgen -source composition.go -package mocks -destination mocks/composition_mock.go
//go:generate mockgen -source directory.go -package mocks -destination mocks/directory_mock.go
//go:generate mockgen -source user.go -package mocks -destination mocks/user_mock.go
//go:generate mockgen -source ../../../pkg/helper/finder.go -package mocks -destination mocks/finder_mock.go
//

func TestDirectoryHandler_Create(t *testing.T) {
	var (
		userUUID    = uuid.New()
		ehrUUID     = uuid.New()
		directoryID = uuid.New()
		systemID    = common.EhrSystemID
	)

	directoryVersionUID, _ := base.NewObjectVersionID(directoryID.String(), systemID)
	nextDirectoryVersionUID, _ := base.NewObjectVersionID(directoryVersionUID.String(), systemID)
	_, _ = nextDirectoryVersionUID.IncreaseUIDVersion()

	tests := []struct {
		name       string
		payload    []byte
		prepare    func(dSvc *mocks.MockDirectoryService, iSvc *mocks.MockIndexer, uS *mocks.MockUserService)
		wantStatus int
		wantResp   *model.Directory
	}{
		{
			"1. failed because EhrID is not belong to current user",
			[]byte(``),
			func(_ *mocks.MockDirectoryService, i *mocks.MockIndexer, _ *mocks.MockUserService) {
				i.EXPECT().GetEhrUUIDByUserID(gomock.Any(), userUUID.String(), systemID).Return(nil, errors.ErrNotFound)
			},
			http.StatusNotFound,
			nil,
		},
		{
			"2. failed because body is empty",
			[]byte(``),
			func(_ *mocks.MockDirectoryService, i *mocks.MockIndexer, _ *mocks.MockUserService) {
				i.EXPECT().GetEhrUUIDByUserID(gomock.Any(), userUUID.String(), systemID).Return(&ehrUUID, errors.ErrNotFound)
			},
			http.StatusNotFound,
			nil,
		},
		{
			"4. Failed creation of DIRECTORY because doc already exist",
			[]byte(`{
					"_type": "FOLDER",
						"name": {
							"_type": "DV_TEXT",
							"value": "root"
					},
					"archetype_node_id": "openEHR-EHR-FOLDER.generic.v1"
				}
			`),
			func(dS *mocks.MockDirectoryService, i *mocks.MockIndexer, _ *mocks.MockUserService) {
				i.EXPECT().GetEhrUUIDByUserID(gomock.Any(), userUUID.String(), systemID).Return(&ehrUUID, nil)
				dS.EXPECT().GetByTimeOrLast(gomock.Any(), systemID, &ehrUUID, userUUID.String(), gomock.Any()).Return(&model.Directory{Locatable: base.Locatable{
					ObjectVersionID: *directoryVersionUID,
					Type:            base.FolderItemType,
					Name:            base.NewDvText("root"),
					ArchetypeNodeID: "openEHR-EHR-FOLDER.generic.v1",
				},
				}, nil)
			},
			http.StatusConflict,
			nil,
		},
		{
			"5. Success DIRECTORY create with increased version ID because last doc was deleted",
			[]byte(`{
					"_type": "FOLDER",
						"name": {
							"_type": "DV_TEXT",
							"value": "root"
					},
					"archetype_node_id": "openEHR-EHR-FOLDER.generic.v1"
				}
			`),
			func(dS *mocks.MockDirectoryService, i *mocks.MockIndexer, _ *mocks.MockUserService) {
				i.EXPECT().GetEhrUUIDByUserID(gomock.Any(), userUUID.String(), systemID).Return(&ehrUUID, nil)

				d := &model.Directory{Locatable: base.Locatable{
					ObjectVersionID: *directoryVersionUID,
					Type:            base.FolderItemType,
					Name:            base.NewDvText("root"),
					ArchetypeNodeID: "openEHR-EHR-FOLDER.generic.v1",
				}}

				dS.EXPECT().GetByTimeOrLast(gomock.Any(), systemID, &ehrUUID, userUUID.String(), gomock.Any()).Return(d, errors.ErrAlreadyDeleted)

				dS.EXPECT().IncreaseVersion(d, systemID).DoAndReturn(func(_ *model.Directory, _ string) (string, error) {
					nextDirectoryVersionUID.UID.Value = nextDirectoryVersionUID.String()
					return nextDirectoryVersionUID.String(), nil
				})

				dS.EXPECT().GetByID(gomock.Any(), userUUID.String(), systemID, &ehrUUID, nextDirectoryVersionUID).Return(nil, errors.ErrNotFound)

				dS.EXPECT().GetActiveProcRequest(userUUID.String(), processing.RequestDirectoryCreate).Return("", nil)

				dS.EXPECT().NewProcRequest(gomock.Any(), userUUID.String(), ehrUUID.String(), processing.RequestDirectoryCreate).Return(&MockRequest{}, nil)
				dS.EXPECT().Create(gomock.Any(), gomock.Any(), userUUID.String(), systemID, nextDirectoryVersionUID.String(), gomock.Any()).Return(nil)
			},
			http.StatusCreated,
			nil,
		},
		{
			"6. Filed creation DIRECTORY because DB contain prev unfinished request",
			[]byte(`{
					"_type": "FOLDER",
						"name": {
							"_type": "DV_TEXT",
							"value": "root"
					},
					"archetype_node_id": "openEHR-EHR-FOLDER.generic.v1"
				}
			`),
			func(dS *mocks.MockDirectoryService, i *mocks.MockIndexer, _ *mocks.MockUserService) {
				i.EXPECT().GetEhrUUIDByUserID(gomock.Any(), userUUID.String(), systemID).Return(&ehrUUID, nil)

				d := &model.Directory{Locatable: base.Locatable{
					ObjectVersionID: *directoryVersionUID,
					Type:            base.FolderItemType,
					Name:            base.NewDvText("root"),
					ArchetypeNodeID: "openEHR-EHR-FOLDER.generic.v1",
				}}

				dS.EXPECT().GetByTimeOrLast(gomock.Any(), systemID, &ehrUUID, userUUID.String(), gomock.Any()).Return(d, errors.ErrAlreadyDeleted)

				dS.EXPECT().IncreaseVersion(d, systemID).DoAndReturn(func(_ *model.Directory, _ string) (string, error) {
					nextDirectoryVersionUID.UID.Value = nextDirectoryVersionUID.String()
					return nextDirectoryVersionUID.String(), nil
				})

				dS.EXPECT().GetByID(gomock.Any(), userUUID.String(), systemID, &ehrUUID, nextDirectoryVersionUID).Return(nil, errors.ErrNotFound)

				dS.EXPECT().GetActiveProcRequest(userUUID.String(), processing.RequestDirectoryCreate).Return("not empty", nil)
			},
			http.StatusConflict,
			nil,
		},
		{
			"7. Success DIRECTORY create",
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
				i.EXPECT().GetEhrUUIDByUserID(gomock.Any(), userUUID.String(), systemID).Return(&ehrUUID, nil)

				dS.EXPECT().GetByTimeOrLast(gomock.Any(), systemID, &ehrUUID, userUUID.String(), gomock.Any()).Return(nil, nil)

				dS.EXPECT().GetByID(gomock.Any(), userUUID.String(), systemID, &ehrUUID, directoryVersionUID).Return(nil, errors.ErrNotFound)

				dS.EXPECT().GetActiveProcRequest(userUUID.String(), processing.RequestDirectoryCreate).Return("", nil)

				dS.EXPECT().NewProcRequest(gomock.Any(), userUUID.String(), ehrUUID.String(), processing.RequestDirectoryCreate).Return(&MockRequest{}, nil)
				dS.EXPECT().Create(gomock.Any(), gomock.Any(), userUUID.String(), systemID, directoryVersionUID.String(), gomock.Any()).Return(nil)
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

			api := API{
				Directory: NewDirectoryHandler(svc, userSvc, indexer, ""),
				User:      NewUserHandler(userSvc),
			}

			router := api.setupRouter(api.buildEhrDirectoryAPI())

			req := httptest.NewRequest(http.MethodPost, fmt.Sprintf("/v1/ehr/%s/directory/?&patient_id=%s", ehrUUID.String(), userUUID.String()), bytes.NewReader(tt.payload))
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

func TestDirectoryHandler_Update(t *testing.T) {
	var (
		userUUID    = uuid.New()
		ehrUUID     = uuid.New()
		directoryID = uuid.New()
		systemID    = common.EhrSystemID
	)

	directoryVersionUID, _ := base.NewObjectVersionID(directoryID.String(), systemID)
	nextDirectoryVersionUID, _ := base.NewObjectVersionID(directoryVersionUID.String(), systemID)
	_, _ = nextDirectoryVersionUID.IncreaseUIDVersion()

	tests := []struct {
		name          string
		directoryID   string
		payload       []byte
		prepare       func(dSvc *mocks.MockDirectoryService, iSvc *mocks.MockIndexer, uS *mocks.MockUserService)
		wantStatus    int
		wantResp      *model.Directory
		wantVersionID string
	}{
		{
			"1. failed because EhrID is not belong to current user",
			directoryVersionUID.String(),
			[]byte(``),
			func(_ *mocks.MockDirectoryService, i *mocks.MockIndexer, _ *mocks.MockUserService) {
				i.EXPECT().GetEhrUUIDByUserID(gomock.Any(), userUUID.String(), systemID).Return(nil, errors.ErrNotFound)
			},
			http.StatusNotFound,
			nil,
			"",
		},
		{
			"2. failed because body is empty",
			directoryVersionUID.String(),
			[]byte(``),
			func(_ *mocks.MockDirectoryService, i *mocks.MockIndexer, _ *mocks.MockUserService) {
				i.EXPECT().GetEhrUUIDByUserID(gomock.Any(), userUUID.String(), systemID).Return(&ehrUUID, errors.ErrNotFound)
			},
			http.StatusNotFound,
			nil,
			"",
		},
		{
			"3. failed because DIRECTORY with current version was not found",
			directoryVersionUID.String(),
			[]byte(`{
					"_type": "FOLDER",
						"name": {
							"_type": "DV_TEXT",
							"value": "root"
					},
					"uid": {
						"_type": "OBJECT_VERSION_ID",
						"value": "` + nextDirectoryVersionUID.String() + `"
		
					},
					"archetype_node_id": "openEHR-EHR-FOLDER.generic.v1"
				}
			`),
			func(dS *mocks.MockDirectoryService, i *mocks.MockIndexer, _ *mocks.MockUserService) {
				i.EXPECT().GetEhrUUIDByUserID(gomock.Any(), userUUID.String(), systemID).Return(&ehrUUID, nil)
				dS.EXPECT().GetByTimeOrLast(gomock.Any(), systemID, &ehrUUID, userUUID.String(), gomock.Any()).Return(nil, errors.ErrNotFound)
			},
			http.StatusNotFound,
			nil,
			"",
		},
		{
			"4. failed because DIRECTORY with current version ID is not match with version in DB",
			directoryVersionUID.String(),
			[]byte(`{
					"_type": "FOLDER",
						"name": {
							"_type": "DV_TEXT",
							"value": "root"
					},
					"uid": {
						"_type": "OBJECT_VERSION_ID",
						"value": "` + nextDirectoryVersionUID.String() + `"
		
					},
					"archetype_node_id": "openEHR-EHR-FOLDER.generic.v1"
				}
			`),
			func(dS *mocks.MockDirectoryService, i *mocks.MockIndexer, _ *mocks.MockUserService) {
				i.EXPECT().GetEhrUUIDByUserID(gomock.Any(), userUUID.String(), systemID).Return(&ehrUUID, nil)

				d := &model.Directory{
					Locatable: base.Locatable{
						ObjectVersionID: *directoryVersionUID,
					},
				}
				dS.EXPECT().GetByTimeOrLast(gomock.Any(), systemID, &ehrUUID, userUUID.String(), gomock.Any()).Return(d, nil)
			},
			http.StatusPreconditionFailed,
			nil,
			"",
		},
		{
			"5. Success DIRECTORY updated",
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
				i.EXPECT().GetEhrUUIDByUserID(gomock.Any(), userUUID.String(), systemID).Return(&ehrUUID, nil)
				d := &model.Directory{
					Locatable: base.Locatable{
						ObjectVersionID: *directoryVersionUID,
					},
				}

				dS.EXPECT().GetByTimeOrLast(gomock.Any(), systemID, &ehrUUID, userUUID.String(), gomock.Any()).Return(d, nil)
				dS.EXPECT().GetActiveProcRequest(userUUID.String(), processing.RequestDirectoryUpdate).Return("", nil)
				dS.EXPECT().NewProcRequest(gomock.Any(), userUUID.String(), ehrUUID.String(), processing.RequestDirectoryUpdate).Return(&MockRequest{}, nil)
				dS.EXPECT().Update(gomock.Any(), gomock.Any(), systemID, userUUID.String(), gomock.Any()).Return(nil)
			},
			http.StatusOK,
			nil,
			nextDirectoryVersionUID.String(),
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

			api := API{
				Directory: NewDirectoryHandler(svc, userSvc, indexer, ""),
				User:      NewUserHandler(userSvc),
			}

			router := api.setupRouter(api.buildEhrDirectoryAPI())

			req := httptest.NewRequest(http.MethodPut, fmt.Sprintf("/v1/ehr/%s/directory/?&patient_id=%s", ehrUUID.String(), userUUID.String()), bytes.NewReader(tt.payload))
			req.Header.Set("Authorization", "Bearer emptyJWTkey")
			req.Header.Set("AuthUserId", userUUID.String())
			req.Header.Set("EhrSystemId", systemID)
			req.Header.Set("If-Match", directoryVersionUID.String())
			req.Header.Set("Prefer", "return=representation")

			recorder := httptest.NewRecorder()
			router.ServeHTTP(recorder, req)

			resp := recorder.Result()
			defer resp.Body.Close()

			if diff := cmp.Diff(tt.wantStatus, resp.StatusCode); diff != "" {
				t.Errorf("DirectoryHandler.Update() status code mismatch {-want;+got}\n\t%s", diff)
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
				t.Errorf("DirectoryHandler.Update() body response {-want;+got}\n%s", diff)
			}

			if resp.StatusCode == http.StatusOK || resp.StatusCode == http.StatusNoContent {
				etag := req.Header.Get("ETag")
				if diff := cmp.Diff(etag, fmt.Sprintf("\"%s\"", tt.wantVersionID)); diff != "" {
					t.Errorf("DirectoryHandler.Update() etag {-want;+got}\n%s", diff)
				}
			}
		})
	}
}

func TestDirectoryHandler_Delete(t *testing.T) {
	var (
		userUUID    = uuid.New()
		ehrUUID     = uuid.New()
		directoryID = uuid.New()
		systemID    = common.EhrSystemID
	)

	directoryVersionUID, _ := base.NewObjectVersionID(directoryID.String(), systemID)
	nextDirectoryVersionUID, _ := base.NewObjectVersionID(directoryVersionUID.String(), systemID)
	_, _ = nextDirectoryVersionUID.IncreaseUIDVersion()

	tests := []struct {
		name          string
		directoryID   string
		prepare       func(dSvc *mocks.MockDirectoryService, iSvc *mocks.MockIndexer)
		wantStatus    int
		wantResp      *model.Directory
		wantVersionID string
	}{
		{
			"1. failed because EhrID is not belong to current user",
			directoryVersionUID.String(),
			func(_ *mocks.MockDirectoryService, i *mocks.MockIndexer) {
				i.EXPECT().GetEhrUUIDByUserID(gomock.Any(), userUUID.String(), systemID).Return(nil, errors.ErrNotFound)
			},
			http.StatusNotFound,
			nil,
			"",
		},
		{
			"2. failed because DIRECTORY with current version was not found",
			directoryVersionUID.String(),
			func(dS *mocks.MockDirectoryService, i *mocks.MockIndexer) {
				i.EXPECT().GetEhrUUIDByUserID(gomock.Any(), userUUID.String(), systemID).Return(&ehrUUID, nil)
				dS.EXPECT().GetByTimeOrLast(gomock.Any(), systemID, &ehrUUID, userUUID.String(), gomock.Any()).Return(nil, errors.ErrNotFound)
			},
			http.StatusNotFound,
			nil,
			"",
		},
		{
			"3. failed because DIRECTORY with current {version_id} is not match with {version_id} stored in DB",
			directoryVersionUID.String(),
			func(dS *mocks.MockDirectoryService, i *mocks.MockIndexer) {
				i.EXPECT().GetEhrUUIDByUserID(gomock.Any(), userUUID.String(), systemID).Return(&ehrUUID, nil)

				d := &model.Directory{
					Locatable: base.Locatable{
						ObjectVersionID: base.ObjectVersionID{
							UID: &base.UIDBasedID{
								ObjectID: base.ObjectID{
									Type:  "OBJECT_VERSION_ID",
									Value: nextDirectoryVersionUID.String(),
								},
							},
						},
					},
				}
				dS.EXPECT().GetByTimeOrLast(gomock.Any(), systemID, &ehrUUID, userUUID.String(), gomock.Any()).Return(d, nil)
			},
			http.StatusPreconditionFailed,
			nil,
			"",
		},
		{
			"4. Success DIRECTORY deleted",
			directoryID.String(),
			func(dS *mocks.MockDirectoryService, i *mocks.MockIndexer) {
				i.EXPECT().GetEhrUUIDByUserID(gomock.Any(), userUUID.String(), systemID).Return(&ehrUUID, nil)
				d := &model.Directory{
					Locatable: base.Locatable{
						ObjectVersionID: *directoryVersionUID,
					},
				}

				dS.EXPECT().GetByTimeOrLast(gomock.Any(), systemID, &ehrUUID, userUUID.String(), gomock.Any()).Return(d, nil)
				dS.EXPECT().GetActiveProcRequest(userUUID.String(), processing.RequestDirectoryDelete).Return("", nil)
				dS.EXPECT().NewProcRequest(gomock.Any(), userUUID.String(), ehrUUID.String(), processing.RequestDirectoryDelete).Return(&MockRequest{}, nil)
				dS.EXPECT().Delete(gomock.Any(), gomock.Any(), systemID, &ehrUUID, directoryVersionUID.String(), userUUID.String()).Return(d.String(), nil)
			},
			http.StatusNoContent,
			nil,
			nextDirectoryVersionUID.String(),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			svc := mocks.NewMockDirectoryService(ctrl)
			indexer := mocks.NewMockIndexer(ctrl)
			userSvc := mocks.NewMockUserService(ctrl)
			tt.prepare(svc, indexer)

			// Mock for auth user service
			userSvc.EXPECT().VerifyAccess(gomock.Any(), gomock.Any()).Return(nil)

			api := API{
				Directory: NewDirectoryHandler(svc, userSvc, indexer, ""),
				User:      NewUserHandler(userSvc),
			}

			router := api.setupRouter(api.buildEhrDirectoryAPI())

			req := httptest.NewRequest(http.MethodDelete, fmt.Sprintf("/v1/ehr/%s/directory/?&patient_id=%s", ehrUUID.String(), userUUID.String()), nil)
			req.Header.Set("Authorization", "Bearer emptyJWTkey")
			req.Header.Set("AuthUserId", userUUID.String())
			req.Header.Set("EhrSystemId", systemID)
			req.Header.Set("If-Match", directoryVersionUID.String())
			req.Header.Set("Prefer", "return=representation")

			recorder := httptest.NewRecorder()
			router.ServeHTTP(recorder, req)

			resp := recorder.Result()
			defer resp.Body.Close()

			if diff := cmp.Diff(tt.wantStatus, resp.StatusCode); diff != "" {
				t.Errorf("DirectoryHandler.Delete() status code mismatch {-want;+got}\n\t%s", diff)
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
				t.Errorf("DirectoryHandler.Delete() body response {-want;+got}\n%s", diff)
			}

			if resp.StatusCode == http.StatusNoContent {
				etag := req.Header.Get("ETag")
				if diff := cmp.Diff(etag, fmt.Sprintf("\"%s\"", tt.wantVersionID)); diff != "" {
					t.Errorf("DirectoryHandler.Delete() etag {-want;+got}\n%s", diff)
				}
			}
		})
	}
}

func TestDirectoryHandler_GetByTimeOrLast(t *testing.T) {
	var (
		userUUID    = uuid.New()
		ehrUUID     = uuid.New()
		directoryID = uuid.New()
		systemID    = common.EhrSystemID
	)

	directoryVersionUID, _ := base.NewObjectVersionID(directoryID.String(), systemID)
	d := newDirectoryWithFolders(directoryVersionUID)

	tests := []struct {
		name         string
		timeVersion  string
		path         string
		prepare      func(dSvc *mocks.MockDirectoryService, iSvc *mocks.MockIndexer)
		wantStatus   int
		wantResp     bool
		wantPathName string
	}{
		{
			"1. failed because EhrID is not belong to current user",
			"",
			"",
			func(_ *mocks.MockDirectoryService, i *mocks.MockIndexer) {
				i.EXPECT().GetEhrUUIDByUserID(gomock.Any(), userUUID.String(), systemID).Return(nil, errors.ErrNotFound)
			},
			http.StatusNotFound,
			false,
			"",
		},
		{
			"2. failed because time has incorrect format",
			"incorrect_time_format",
			"",
			func(dS *mocks.MockDirectoryService, i *mocks.MockIndexer) {
				i.EXPECT().GetEhrUUIDByUserID(gomock.Any(), userUUID.String(), systemID).Return(&ehrUUID, nil)
			},
			http.StatusBadRequest,
			false,
			"",
		},
		{
			"3. failed because DIRECTORY by exact time was not found",
			time.Now().Format(common.OpenEhrTimeFormat),
			"",
			func(dS *mocks.MockDirectoryService, i *mocks.MockIndexer) {
				i.EXPECT().GetEhrUUIDByUserID(gomock.Any(), userUUID.String(), systemID).Return(&ehrUUID, nil)
				dS.EXPECT().GetByTimeOrLast(gomock.Any(), systemID, &ehrUUID, userUUID.String(), gomock.Any()).Return(nil, errors.ErrNotFound)
			},
			http.StatusNotFound,
			false,
			"",
		},
		{
			"4. failed because DIRECTORY by time was deleted",
			time.Now().Format(common.OpenEhrTimeFormat),
			"",
			func(dS *mocks.MockDirectoryService, i *mocks.MockIndexer) {
				i.EXPECT().GetEhrUUIDByUserID(gomock.Any(), userUUID.String(), systemID).Return(&ehrUUID, nil)
				dS.EXPECT().GetByTimeOrLast(gomock.Any(), systemID, &ehrUUID, userUUID.String(), gomock.Any()).Return(nil, errors.ErrAlreadyDeleted)
			},
			http.StatusNoContent,
			false,
			"",
		},
		{
			"5. failed because DIRECTORY by {path} was not found",
			time.Now().Format(common.OpenEhrTimeFormat),
			"not_exist_path",
			func(dS *mocks.MockDirectoryService, i *mocks.MockIndexer) {
				i.EXPECT().GetEhrUUIDByUserID(gomock.Any(), userUUID.String(), systemID).Return(&ehrUUID, nil)
				dS.EXPECT().GetByTimeOrLast(gomock.Any(), systemID, &ehrUUID, userUUID.String(), gomock.Any()).Return(d, nil)
			},
			http.StatusNotFound,
			false,
			"",
		},
		{
			"6. success return root folder from DIRECTORY by {version_at_time} and {path}",
			time.Now().Format(common.OpenEhrTimeFormat),
			"root",
			func(dS *mocks.MockDirectoryService, i *mocks.MockIndexer) {
				i.EXPECT().GetEhrUUIDByUserID(gomock.Any(), userUUID.String(), systemID).Return(&ehrUUID, nil)
				dS.EXPECT().GetByTimeOrLast(gomock.Any(), systemID, &ehrUUID, userUUID.String(), gomock.Any()).Return(d, nil)
			},
			http.StatusOK,
			true,
			"root",
		},
		{
			"7. success return sub folder from DIRECTORY by {version_at_time} and {path}",
			time.Now().Format(common.OpenEhrTimeFormat),
			"root/1",
			func(dS *mocks.MockDirectoryService, i *mocks.MockIndexer) {
				i.EXPECT().GetEhrUUIDByUserID(gomock.Any(), userUUID.String(), systemID).Return(&ehrUUID, nil)
				dS.EXPECT().GetByTimeOrLast(gomock.Any(), systemID, &ehrUUID, userUUID.String(), gomock.Any()).Return(d, nil)
			},
			http.StatusOK,
			true,
			"1",
		},
		{
			"8. success return last sub folder from DIRECTORY by {version_at_time} and {path}",
			time.Now().Format(common.OpenEhrTimeFormat),
			"root/1/1-1/1-1-2",
			func(dS *mocks.MockDirectoryService, i *mocks.MockIndexer) {
				i.EXPECT().GetEhrUUIDByUserID(gomock.Any(), userUUID.String(), systemID).Return(&ehrUUID, nil)
				dS.EXPECT().GetByTimeOrLast(gomock.Any(), systemID, &ehrUUID, userUUID.String(), gomock.Any()).Return(d, nil)
			},
			http.StatusOK,
			true,
			"1-1-2",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			svc := mocks.NewMockDirectoryService(ctrl)
			indexer := mocks.NewMockIndexer(ctrl)
			userSvc := mocks.NewMockUserService(ctrl)
			tt.prepare(svc, indexer)

			// Mock for auth user service
			userSvc.EXPECT().VerifyAccess(gomock.Any(), gomock.Any()).Return(nil)

			api := API{
				Directory: NewDirectoryHandler(svc, userSvc, indexer, ""),
				User:      NewUserHandler(userSvc),
			}

			router := api.setupRouter(api.buildEhrDirectoryAPI())

			req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/v1/ehr/%s/directory/?version_at_time=%s&path=%s&patient_id=%s", ehrUUID.String(), url.QueryEscape(tt.timeVersion), url.QueryEscape(tt.path), userUUID.String()), nil)
			req.Header.Set("Authorization", "Bearer emptyJWTkey")
			req.Header.Set("AuthUserId", userUUID.String())
			req.Header.Set("EhrSystemId", systemID)

			recorder := httptest.NewRecorder()
			router.ServeHTTP(recorder, req)

			resp := recorder.Result()
			defer resp.Body.Close()

			if diff := cmp.Diff(tt.wantStatus, resp.StatusCode); diff != "" {
				t.Errorf("DirectoryHandler.GetByTimeOrLast() status code mismatch {-want;+got}\n\t%s", diff)
				return
			}

			if tt.wantResp == false {
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

			if diff := cmp.Diff(tt.wantPathName, got.Name.Value, opts); diff != "" {
				t.Errorf("DirectoryHandler.GetByTimeOrLast() body response {-want;+got}\n%s", diff)
			}
		})
	}
}

func TestDirectoryHandler_GetByVersion(t *testing.T) {
	var (
		userUUID    = uuid.New()
		ehrUUID     = uuid.New()
		directoryID = uuid.New()
		systemID    = common.EhrSystemID
	)

	directoryVersionUID, _ := base.NewObjectVersionID(directoryID.String(), systemID)
	d := newDirectoryWithFolders(directoryVersionUID)

	nextDirectoryVersionUID, _ := base.NewObjectVersionID(directoryVersionUID.String(), systemID)
	_, _ = nextDirectoryVersionUID.IncreaseUIDVersion()

	tests := []struct {
		name         string
		version      string
		path         string
		prepare      func(dSvc *mocks.MockDirectoryService, iSvc *mocks.MockIndexer)
		wantStatus   int
		wantResp     bool
		wantPathName string
	}{
		{
			"1. failed because EhrID is not belong to current user",
			directoryVersionUID.String(),
			"",
			func(_ *mocks.MockDirectoryService, i *mocks.MockIndexer) {
				i.EXPECT().GetEhrUUIDByUserID(gomock.Any(), userUUID.String(), systemID).Return(nil, errors.ErrNotFound)
			},
			http.StatusNotFound,
			false,
			"",
		},
		{
			"2. failed because time has incorrect format",
			"incorrect_time_format",
			"",
			func(dS *mocks.MockDirectoryService, i *mocks.MockIndexer) {
				i.EXPECT().GetEhrUUIDByUserID(gomock.Any(), userUUID.String(), systemID).Return(&ehrUUID, nil)
			},
			http.StatusBadRequest,
			false,
			"",
		},
		{
			"3. failed because DIRECTORY by {version_uid} was not found",
			nextDirectoryVersionUID.String(),
			"",
			func(dS *mocks.MockDirectoryService, i *mocks.MockIndexer) {
				i.EXPECT().GetEhrUUIDByUserID(gomock.Any(), userUUID.String(), systemID).Return(&ehrUUID, nil)
				//dS.EXPECT().GetByID(gomock.Any(), userUUID.String(), systemID, &ehrUUID, nextDirectoryVersionUID).Return(nil, errors.ErrNotFound)
				dS.EXPECT().GetByID(gomock.Any(), userUUID.String(), systemID, &ehrUUID, gomock.Any()).Return(nil, errors.ErrNotFound)
			},
			http.StatusNotFound,
			false,
			"",
		},
		{
			"4. failed because DIRECTORY by {version_uid} was deleted",
			nextDirectoryVersionUID.String(),
			"",
			func(dS *mocks.MockDirectoryService, i *mocks.MockIndexer) {
				i.EXPECT().GetEhrUUIDByUserID(gomock.Any(), userUUID.String(), systemID).Return(&ehrUUID, nil)
				dS.EXPECT().GetByID(gomock.Any(), userUUID.String(), systemID, &ehrUUID, gomock.Any()).Return(nil, errors.ErrAlreadyDeleted)
			},
			http.StatusNoContent,
			false,
			"",
		},
		{
			"5. failed because DIRECTORY by {path} was not found",
			d.String(),
			"not_exist_path",
			func(dS *mocks.MockDirectoryService, i *mocks.MockIndexer) {
				i.EXPECT().GetEhrUUIDByUserID(gomock.Any(), userUUID.String(), systemID).Return(&ehrUUID, nil)
				dS.EXPECT().GetByID(gomock.Any(), userUUID.String(), systemID, &ehrUUID, gomock.Any()).Return(d, nil)
			},
			http.StatusNotFound,
			false,
			"",
		},
		{
			"6. success return root folder from DIRECTORY by {version_at_time} and {path}",
			d.String(),
			"root",
			func(dS *mocks.MockDirectoryService, i *mocks.MockIndexer) {
				i.EXPECT().GetEhrUUIDByUserID(gomock.Any(), userUUID.String(), systemID).Return(&ehrUUID, nil)
				dS.EXPECT().GetByID(gomock.Any(), userUUID.String(), systemID, &ehrUUID, gomock.Any()).Return(d, nil)
			},
			http.StatusOK,
			true,
			"root",
		},
		{
			"7. success return sub folder from DIRECTORY by {version_at_time} and {path}",
			d.String(),
			"root/1",
			func(dS *mocks.MockDirectoryService, i *mocks.MockIndexer) {
				i.EXPECT().GetEhrUUIDByUserID(gomock.Any(), userUUID.String(), systemID).Return(&ehrUUID, nil)
				dS.EXPECT().GetByID(gomock.Any(), userUUID.String(), systemID, &ehrUUID, gomock.Any()).Return(d, nil)
			},
			http.StatusOK,
			true,
			"1",
		},
		{
			"8. success return last sub folder from DIRECTORY by {version_at_time} and {path}",
			d.String(),
			"root/1/1-1/1-1-2",
			func(dS *mocks.MockDirectoryService, i *mocks.MockIndexer) {
				i.EXPECT().GetEhrUUIDByUserID(gomock.Any(), userUUID.String(), systemID).Return(&ehrUUID, nil)
				dS.EXPECT().GetByID(gomock.Any(), userUUID.String(), systemID, &ehrUUID, gomock.Any()).Return(d, nil)
			},
			http.StatusOK,
			true,
			"1-1-2",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			svc := mocks.NewMockDirectoryService(ctrl)
			indexer := mocks.NewMockIndexer(ctrl)
			userSvc := mocks.NewMockUserService(ctrl)
			tt.prepare(svc, indexer)

			// Mock for auth user service
			userSvc.EXPECT().VerifyAccess(gomock.Any(), gomock.Any()).Return(nil)

			api := API{
				Directory: NewDirectoryHandler(svc, userSvc, indexer, ""),
				User:      NewUserHandler(userSvc),
			}

			router := api.setupRouter(api.buildEhrDirectoryAPI())

			req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/v1/ehr/%s/directory/%s?path=%s&patient_id=%s", ehrUUID.String(), url.QueryEscape(tt.version), url.QueryEscape(tt.path), userUUID.String()), nil)
			req.Header.Set("Authorization", "Bearer emptyJWTkey")
			req.Header.Set("AuthUserId", userUUID.String())
			req.Header.Set("EhrSystemId", systemID)

			recorder := httptest.NewRecorder()
			router.ServeHTTP(recorder, req)

			resp := recorder.Result()
			defer resp.Body.Close()

			if diff := cmp.Diff(tt.wantStatus, resp.StatusCode); diff != "" {
				t.Errorf("DirectoryHandler.GetByTimeOrLast() status code mismatch {-want;+got}\n\t%s", diff)
				return
			}

			if tt.wantResp == false {
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

			if diff := cmp.Diff(tt.wantPathName, got.Name.Value, opts); diff != "" {
				t.Errorf("DirectoryHandler.GetByTimeOrLast() body response {-want;+got}\n%s", diff)
			}
		})
	}
}

func newDirectoryWithFolders(id *base.ObjectVersionID) *model.Directory {
	return &model.Directory{
		Locatable: base.Locatable{
			ObjectVersionID: *id,
			Type:            base.FolderItemType,
			Name:            base.NewDvText("root"),
			ArchetypeNodeID: "openEHR-EHR-FOLDER.generic.v1",
		},
		Folders: []*model.Directory{
			{
				Locatable: base.Locatable{
					Type:            base.FolderItemType,
					Name:            base.NewDvText("1"),
					ArchetypeNodeID: "openEHR-EHR-FOLDER.generic.v1",
				},
				Folders: []*model.Directory{
					{
						Locatable: base.Locatable{
							Type:            base.FolderItemType,
							Name:            base.NewDvText("1-1"),
							ArchetypeNodeID: "openEHR-EHR-FOLDER.generic.v1",
						},
						Folders: []*model.Directory{
							{
								Locatable: base.Locatable{
									Type:            base.FolderItemType,
									Name:            base.NewDvText("1-1-1"),
									ArchetypeNodeID: "openEHR-EHR-FOLDER.generic.v1",
								},
							},
							{
								Locatable: base.Locatable{
									Type:            base.FolderItemType,
									Name:            base.NewDvText("1-1-2"),
									ArchetypeNodeID: "openEHR-EHR-FOLDER.generic.v1",
								},
							},
						},
					},
					{
						Locatable: base.Locatable{
							Type:            base.FolderItemType,
							Name:            base.NewDvText("1-2"),
							ArchetypeNodeID: "openEHR-EHR-FOLDER.generic.v1",
						},
						FeederAudit: base.FeederAudit{},
						Details:     base.ItemStructure{},
						Folders: []*model.Directory{
							{
								Locatable: base.Locatable{
									Type:            base.FolderItemType,
									Name:            base.NewDvText("1-2-1"),
									ArchetypeNodeID: "openEHR-EHR-FOLDER.generic.v1",
								},
								FeederAudit: base.FeederAudit{},
								Details:     base.ItemStructure{},
								Folders:     nil,
							},
						},
					},
				},
			},
			{
				Locatable: base.Locatable{
					Type:            base.FolderItemType,
					Name:            base.NewDvText("2"),
					ArchetypeNodeID: "openEHR-EHR-FOLDER.generic.v1",
				},
				FeederAudit: base.FeederAudit{},
				Details:     base.ItemStructure{},
				Folders: []*model.Directory{
					{
						Locatable: base.Locatable{
							Type:            base.FolderItemType,
							Name:            base.NewDvText("2-1"),
							ArchetypeNodeID: "openEHR-EHR-FOLDER.generic.v1",
						},
						FeederAudit: base.FeederAudit{},
						Details:     base.ItemStructure{},
						Folders:     nil,
					},
				},
			},
			{
				Locatable: base.Locatable{
					Type:            base.FolderItemType,
					Name:            base.NewDvText("3"),
					ArchetypeNodeID: "openEHR-EHR-FOLDER.generic.v1",
				},
				FeederAudit: base.FeederAudit{},
				Details:     base.ItemStructure{},
				Folders:     nil,
			},
		},
	}
}
