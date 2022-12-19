package api

import (
	"fmt"
	"hms/gateway/pkg/api/mocks"
	"hms/gateway/pkg/docs/model"
	"hms/gateway/pkg/docs/model/base"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/google/go-cmp/cmp"
	"github.com/google/uuid"

	"hms/gateway/pkg/errors"
)

//
//go:generate mockgen -source composition.go -package mocks -destination mocks/composition_mock.go
//go:generate mockgen -source contribution.go -package mocks -destination mocks/contribution_mock.go
//go:generate mockgen -source user.go -package mocks -destination mocks/user_mock.go
//

type contributionTestData struct {
	c     model.ContributionResponse
	cJSON []byte
}

func TestContributionHandler_GetByID(t *testing.T) {
	var (
		userID   = uuid.New().String()
		systemID = uuid.New().String()
	)

	cM := newContributionResponse()

	tests := []struct {
		name            string
		contributionUID string
		prepare         func(gaSvc *mocks.MockContributionService)
		wantStatus      int
		wantResp        string
	}{
		{
			"1. empty result because doc with {contribution_uid} is not exist",
			"123",
			func(gaSvc *mocks.MockContributionService) {
				gaSvc.EXPECT().GetByID(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, errors.ErrNotFound)
			},
			http.StatusNotFound,
			"",
		},
		{
			"2. success result",
			"123",
			func(gaSvc *mocks.MockContributionService) {
				gaSvc.EXPECT().GetByID(gomock.Any(), gomock.Any(), gomock.Any()).Return(&cM.c, nil)
			},
			http.StatusOK,
			string(cM.cJSON),
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
			}

			if tt.wantStatus != http.StatusOK {
				return
			}

			// TODO Dont know how to compare 2 jsons
			//respBody, _ := io.ReadAll(resp.Body)
			//if diff := cmp.Diff(tt.wantResp, string(respBody)); diff != "" {
			//	t.Errorf("ContributionHandler.GetByID() status response {-want;+got}\n%s", diff)
			//}
		})
	}
}

func newContributionResponse() contributionTestData {
	c := model.ContributionResponse{
		UID: base.UIDBasedID{
			ObjectID: base.ObjectID{
				Value: "0826851c-c4c2-4d61-92b9-410fb8275ff0",
			},
		},
		Audit: model.AuditDetails{
			Type:     base.AuditDetailsType,
			SystemID: "test-system-id",
			TimeCommited: base.DvDateTime{
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

	return contributionTestData{c, prepareContributionJSON("")}
}

func prepareContributionJSON(v string) []byte {
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
			"time_commited": "2021-12-03T16:05:19.513939+01:00",
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
