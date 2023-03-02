package gateway

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

	"github.com/bsn-si/IPEHR-gateway/src/internal/api/gateway/mocks"
	"github.com/bsn-si/IPEHR-gateway/src/pkg/common"
	"github.com/bsn-si/IPEHR-gateway/src/pkg/docs/model"
	"github.com/bsn-si/IPEHR-gateway/src/pkg/docs/model/base"
	"github.com/bsn-si/IPEHR-gateway/src/pkg/docs/service/processing"
	"github.com/bsn-si/IPEHR-gateway/src/pkg/errors"
	userModel "github.com/bsn-si/IPEHR-gateway/src/pkg/user/model"
)

//
//go:generate mockgen -source composition.go -package mocks -destination mocks/composition_mock.go
//go:generate mockgen -source contribution.go -package mocks -destination mocks/contribution_mock.go
//go:generate mockgen -source user.go -package mocks -destination mocks/user_mock.go
//

type contributionResponseTestData struct {
	c     *model.ContributionResponse
	cJSON []byte
}

type contributionTestData struct {
	c     *model.Contribution
	cJSON []byte
}

func TestContributionHandler_GetByID(t *testing.T) {
	var (
		userID   = uuid.New().String()
		systemID = common.EhrSystemID
	)

	cR := newContributionResponse()
	contributionID := cR.c.UID.Value

	tests := []struct {
		name            string
		contributionUID string
		prepare         func(gaSvc *mocks.MockContributionService)
		wantStatus      int
		wantResp        *model.ContributionResponse
	}{
		{
			"1. empty result because doc with {contribution_uid} is not exist",
			contributionID,
			func(gaSvc *mocks.MockContributionService) {
				gaSvc.EXPECT().GetByID(gomock.Any(), userID, contributionID).Return(nil, errors.ErrNotFound)
			},
			http.StatusNotFound,
			nil,
		},
		{
			"2. success result",
			contributionID,
			func(gaSvc *mocks.MockContributionService) {
				gaSvc.EXPECT().GetByID(gomock.Any(), userID, contributionID).Return(cR.c, nil)
			},
			http.StatusOK,
			cR.c,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			svc := mocks.NewMockContributionService(ctrl)
			tt.prepare(svc)

			// Mock for auth user service
			userSvc := mocks.NewMockUserService(ctrl)
			userSvc.EXPECT().VerifyAccess(gomock.Any(), gomock.Any()).Return(nil)

			tplSvc := mocks.NewMockTemplateService(ctrl)
			comSvc := mocks.NewMockCompositionService(ctrl)

			api := API{
				Contribution: NewContributionHandler(svc, userSvc, tplSvc, comSvc, ""),
				User:         NewUserHandler(userSvc),
			}

			router := api.setupRouter(api.buildEhrContributionAPI())

			req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/v1/ehr/7d44b88c-4199-4bad-97dc-d78268e01398/contribution/%s", tt.contributionUID), nil)
			req.Header.Set("Authorization", "Bearer emptyJWTkey")
			req.Header.Set("AuthUserId", userID)
			req.Header.Set("EhrSystemId", systemID)

			recorder := httptest.NewRecorder()
			router.ServeHTTP(recorder, req)

			resp := recorder.Result()
			defer resp.Body.Close()

			if diff := cmp.Diff(tt.wantStatus, resp.StatusCode); diff != "" {
				t.Errorf("ContributionHandler.GetByID() status code mismatch {-want;+got}\n\t%s", diff)
				return
			}

			if tt.wantStatus != http.StatusOK {
				return
			}

			var got *model.ContributionResponse
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
				t.Errorf("ContributionHandler.GetByID() status response {-want;+got}\n%s", diff)
			}
		})
	}
}

func TestContributionHandler_Create(t *testing.T) {
	var (
		userID   = uuid.New().String()
		systemID = common.EhrSystemID
		ehrUUID  = uuid.New().String()
	)

	c := newContribution()
	contributionID := c.c.UID.Value
	cR := newContributionResponse()

	tests := []struct {
		name       string
		prepare    func(gaSvc *mocks.MockContributionService, userSvc *mocks.MockUserService)
		payload    []byte
		wantStatus int
		wantResp   *model.ContributionResponse
	}{
		{
			"0. empty result because body is empty",
			func(_ *mocks.MockContributionService, _ *mocks.MockUserService) {},
			nil,
			http.StatusBadRequest,
			nil,
		},
		{
			"1. empty result because doc with {contribution_uid} is already exist",
			func(gaSvc *mocks.MockContributionService, _ *mocks.MockUserService) {
				gaSvc.EXPECT().GetByID(gomock.Any(), userID, contributionID).Return(cR.c, nil)
			},
			c.cJSON,
			http.StatusConflict,
			nil,
		},
		{
			"2. failed because data is not valid",
			func(gaSvc *mocks.MockContributionService, _ *mocks.MockUserService) {
				gaSvc.EXPECT().GetByID(gomock.Any(), userID, contributionID).Return(nil, nil)
				gaSvc.EXPECT().Validate(gomock.Any(), gomock.Any(), gomock.Any()).Return(false, errors.ErrIsNotValid)
			},
			c.cJSON,
			http.StatusBadRequest,
			nil,
		},
		{
			"3. failed because cant save versions",
			func(gaSvc *mocks.MockContributionService, _ *mocks.MockUserService) {
				gaSvc.EXPECT().GetByID(gomock.Any(), userID, contributionID).Return(nil, nil)
				gaSvc.EXPECT().Validate(gomock.Any(), gomock.Any(), gomock.Any()).Return(true, nil)
				gaSvc.EXPECT().NewProcRequest(gomock.Any(), userID, ehrUUID, processing.RequestContributionCreate)
				gaSvc.EXPECT().Execute(gomock.Any(), gomock.Any(), userID, ehrUUID, gomock.Any(), gomock.Any()).Return(errors.ErrNotImplemented)
			},
			c.cJSON,
			http.StatusBadRequest,
			nil,
		},
		{
			"4. failed because cant store CONTRIBUTION",
			func(gaSvc *mocks.MockContributionService, userSvc *mocks.MockUserService) {
				gaSvc.EXPECT().GetByID(gomock.Any(), userID, contributionID).Return(nil, nil)
				gaSvc.EXPECT().Validate(gomock.Any(), gomock.Any(), gomock.Any()).Return(true, nil)
				gaSvc.EXPECT().NewProcRequest(gomock.Any(), userID, ehrUUID, processing.RequestContributionCreate)
				gaSvc.EXPECT().Execute(gomock.Any(), gomock.Any(), userID, ehrUUID, gomock.Any(), gomock.Any()).Return(nil)
				u := &userModel.UserInfo{}
				userSvc.EXPECT().Info(gomock.Any(), userID, systemID).Return(u, nil)
				gaSvc.EXPECT().Store(gomock.Any(), gomock.Any(), systemID, u, gomock.Any()).Return(errors.ErrCustom)
			},
			c.cJSON,
			http.StatusInternalServerError,
			nil,
		},
		{
			"5. failed because can create response for contribution",
			func(gaSvc *mocks.MockContributionService, userSvc *mocks.MockUserService) {
				gaSvc.EXPECT().GetByID(gomock.Any(), userID, contributionID).Return(nil, nil)
				gaSvc.EXPECT().Validate(gomock.Any(), gomock.Any(), gomock.Any()).Return(true, nil)
				gaSvc.EXPECT().NewProcRequest(gomock.Any(), userID, ehrUUID, processing.RequestContributionCreate).Return(&MockRequest{}, nil)
				gaSvc.EXPECT().Execute(gomock.Any(), gomock.Any(), userID, ehrUUID, gomock.Any(), gomock.Any()).Return(nil)

				u := &userModel.UserInfo{}
				userSvc.EXPECT().Info(gomock.Any(), userID, systemID).Return(u, nil)
				gaSvc.EXPECT().Store(gomock.Any(), gomock.Any(), systemID, u, gomock.Any()).Return(nil)
				gaSvc.EXPECT().PrepareResponse(gomock.Any(), systemID, gomock.Any()).Return(nil, errors.ErrNotImplemented)
			},
			c.cJSON,
			http.StatusInternalServerError,
			nil,
		}, {
			"6. successfully created",
			func(gaSvc *mocks.MockContributionService, userSvc *mocks.MockUserService) {
				gaSvc.EXPECT().GetByID(gomock.Any(), userID, contributionID).Return(nil, nil)
				gaSvc.EXPECT().Validate(gomock.Any(), gomock.Any(), gomock.Any()).Return(true, nil)
				gaSvc.EXPECT().NewProcRequest(gomock.Any(), userID, ehrUUID, processing.RequestContributionCreate).Return(&MockRequest{}, nil)
				gaSvc.EXPECT().Execute(gomock.Any(), gomock.Any(), userID, ehrUUID, gomock.Any(), gomock.Any()).Return(nil)

				u := &userModel.UserInfo{}
				userSvc.EXPECT().Info(gomock.Any(), userID, systemID).Return(u, nil)
				gaSvc.EXPECT().Store(gomock.Any(), gomock.Any(), systemID, u, gomock.Any()).Return(nil)
				gaSvc.EXPECT().PrepareResponse(gomock.Any(), systemID, gomock.Any()).Return(cR.c, nil)
			},
			c.cJSON,
			http.StatusCreated,
			cR.c,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			svc := mocks.NewMockContributionService(ctrl)
			userSvc := mocks.NewMockUserService(ctrl)
			tt.prepare(svc, userSvc)

			// Mock for auth user service
			userSvc.EXPECT().VerifyAccess(gomock.Any(), gomock.Any()).Return(nil)

			tplSvc := mocks.NewMockTemplateService(ctrl)
			comSvc := mocks.NewMockCompositionService(ctrl)

			api := API{
				Contribution: NewContributionHandler(svc, userSvc, tplSvc, comSvc, ""),
				User:         NewUserHandler(userSvc),
			}

			router := api.setupRouter(api.buildEhrContributionAPI())

			req := httptest.NewRequest(http.MethodPost, fmt.Sprintf("/v1/ehr/%s/contribution/", ehrUUID), bytes.NewReader(tt.payload))
			req.Header.Set("Authorization", "Bearer emptyJWTkey")
			req.Header.Set("AuthUserId", userID)
			req.Header.Set("EhrSystemId", systemID)
			req.Header.Set("Prefer", "return=representation")

			recorder := httptest.NewRecorder()
			router.ServeHTTP(recorder, req)

			resp := recorder.Result()
			defer resp.Body.Close()

			if diff := cmp.Diff(tt.wantStatus, resp.StatusCode); diff != "" {
				t.Errorf("ContributionHandler.GetByID() status code mismatch {-want;+got}\n\t%s response: %s", diff, resp.Body)

				return
			}

			if tt.wantStatus != http.StatusCreated {
				return
			}

			var got *model.ContributionResponse
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
				t.Errorf("ContributionHandler.GetByID() status response {-want;+got}\n%s", diff)
			}
		})
	}
}

func newContribution() contributionTestData {
	c := model.Contribution{
		UID: base.UIDBasedID{
			ObjectID: base.ObjectID{
				Type:  base.UIDBasedIDItemType,
				Value: "0826851c-c4c2-4d61-92b9-410fb8275ff0",
			},
		},
		Audit: model.AuditDetails{
			Type:     base.AuditDetailsType,
			SystemID: "test-system-id",
			TimeCommitted: base.DvDateTime{
				DvTemporal: base.DvTemporal{
					DvValueBase: base.DvValueBase{Type: base.DvDateTimeItemType},
				},
				Value: "2021-12-03T16:05:19.513939+01:00",
			},
			Committer: base.NewPartyProxy(
				&base.PartyIdentified{
					Name: "<optional name of the committer>",
					PartyProxyBase: base.PartyProxyBase{
						Type: base.PartyIdentifiedItemType,
						ExternalRef: &base.ObjectRef{
							ID: base.ObjectID{
								Type:  "GENERIC_ID",
								Value: "<OBJECT_ID>",
							},
							Namespace: "DEMOGRAPHIC",
							Type:      "PERSON",
						},
					},
				},
			),
			ChangeType: base.DvCodedText{
				DefiningCode: base.CodePhrase{
					TerminologyID: base.ObjectID{Value: "openehr"},
					CodeString:    "249",
				},
				DvText: base.DvText{Value: "creation"},
			},
			Description: base.DvText{Value: "<optional audit description>"},
		},

		Versions: []model.ContributionVersion{},
	}

	return contributionTestData{&c, prepareContributionJSON("")}
}

func prepareContributionJSON(v string) []byte {
	if v == "" {
		v = "[]"
	}

	return []byte(fmt.Sprintf(`{
		"uid":{
			"_type": "UID_BASED_ID",
			"value": "0826851c-c4c2-4d61-92b9-410fb8275ff0"
		},
		"versions": %s,
		"audit": {
			"committer": {
				"_type": "PARTY_IDENTIFIED",
				"name": "<optional name of the committer>"
			},
			"change_type": {
				"value": "creation",
				"defining_code": {
					"terminology_id": {
						"value": "openehr"
					},
					"code_string": "249"
				}
			},
			"description": {
				"value": "<optional audit description>"
			}
		}
	}`, v))
}

func newContributionResponse() contributionResponseTestData {
	c := model.ContributionResponse{
		UID: base.UIDBasedID{
			ObjectID: base.ObjectID{
				Type:  base.UIDBasedIDItemType,
				Value: "0826851c-c4c2-4d61-92b9-410fb8275ff0",
			},
		},
		Audit: model.AuditDetails{
			Type:     base.AuditDetailsType,
			SystemID: "test-system-id",
			TimeCommitted: base.DvDateTime{
				DvTemporal: base.DvTemporal{
					DvValueBase: base.DvValueBase{Type: base.DvDateTimeItemType},
				},
				Value: "2021-12-03T16:05:19.513939+01:00",
			},
			Committer: base.NewPartyProxy(
				&base.PartyIdentified{
					Name: "<optional name of the committer>",
					PartyProxyBase: base.PartyProxyBase{
						Type: base.PartyIdentifiedItemType,
						ExternalRef: &base.ObjectRef{
							ID: base.ObjectID{
								Type:  "GENERIC_ID",
								Value: "<OBJECT_ID>",
							},
							Namespace: "DEMOGRAPHIC",
							Type:      "PERSON",
						},
					},
				},
			),
			ChangeType: base.DvCodedText{
				DefiningCode: base.CodePhrase{
					TerminologyID: base.ObjectID{Value: "openehr"},
					CodeString:    "249",
				},
				DvText: base.DvText{Value: "creation"},
			},
			Description: base.DvText{Value: "<optional audit description>"},
		},

		Versions: []model.ContributionVersionResponse{},
	}

	return contributionResponseTestData{&c, prepareContributionResponseJSON("")}
}

func prepareContributionResponseJSON(v string) []byte {
	if v == "" {
		v = "[]"
	}

	return []byte(fmt.Sprintf(`{
		"uid":{
			"value": "0826851c-c4c2-4d61-92b9-410fb8275ff0"
		},
		"versions": %s,
		"audit": {
			"system_id": "test-system-id",
			"time_committed": "2021-12-03T16:05:19.513939+01:00",
			"committer": {
				"_type": "PARTY_IDENTIFIED",
				"name": "<optional name of the committer>",
				"external_ref": {
					"id": {
						"_type": "GENERIC_ID",
						"value": "<OBJECT_ID>",
					},
					"namespace": "DEMOGRAPHIC",
					"type": "PERSON"
				}
			},
			"change_type": {
				"value": "creation",
				"defining_code": {
					"terminology_id": {
						"value": "openehr"
					},
					"code_string": "249"
				}
			},
			"description": {
				"value": "<optional audit description>"
			}
		}
	}`, v))
}
