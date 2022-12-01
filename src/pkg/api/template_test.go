package api

import (
	"fmt"
	"hms/gateway/pkg/api/mocks"
	"hms/gateway/pkg/docs/model"
	"hms/gateway/pkg/errors"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/google/go-cmp/cmp"
	"github.com/google/uuid"
)

//
//go:generate mockgen -source template.go -package mocks -destination mocks/template_mock.go
//go:generate mockgen -source user.go -package mocks -destination mocks/user_mock.go
//

func TestTemplateHandler_GetByID(t *testing.T) {
	var (
		userID        = uuid.New().String()
		userAccessKey = "emptyJWTkey"
		systemID      = uuid.New().String()
	)

	templateID := "Vital Signs"

	template := &model.Template{
		TemplateID: templateID,
		Version:    "1",
		VerADL:     model.VerADL1_4,
		MimeType:   model.ADLTypeJSON,
		Body: []byte(`{
						  "templateId": "` + templateID + `",
						  "version": "1",
						  "defaultLanguage": "en",
						  "languages": [
							"en","es"
						  ],
						  "tree": {
							"id": "string",
							"name": "string",
							"localizedName": "string",
							"rmType": "string",
							"nodeId": "string",
							"min": 0,
							"max": 0,
							"localizedNames": {
							  "en": "Vital Signs"
							},
							"localizedDescriptions": {
							  "sl": "Generic encounter or progress note composition"
							},
							"aqlPath": "string"
						  }
						}`),
		Concept:     "",
		ArchetypeID: "",
		CreatedAt:   "",
	}

	tests := []struct {
		name       string
		accept     string
		templateID string
		adlVer     string
		prepare    func(gaSvc *mocks.MockTemplateService)
		wantStatus int
		wantResp   string
	}{
		{
			"1. empty result because template not found",
			"nonexistmime",
			"incorrect",
			model.VerADL1_4,
			func(gaSvc *mocks.MockTemplateService) {
				gaSvc.EXPECT().GetByID(gomock.Any(), userID, "incorrect").Return(nil, errors.ErrNotFound)
			},
			http.StatusNotFound,
			"",
		},
		{
			"2. empty result because request Accept header is not match with template",
			"notmatchmimetype",
			templateID,
			model.VerADL1_4,
			func(gaSvc *mocks.MockTemplateService) {
				gaSvc.EXPECT().GetByID(gomock.Any(), userID, gomock.Any()).Return(template, nil)
			},
			http.StatusNotAcceptable,
			"",
		},
		{
			"3. success result",
			template.MimeType,
			templateID,
			model.VerADL1_4,
			func(gaSvc *mocks.MockTemplateService) {
				gaSvc.EXPECT().GetByID(gomock.Any(), userID, gomock.Any()).Return(template, nil)
			},
			http.StatusOK,
			string(template.Body),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			tSvc := mocks.NewMockTemplateService(ctrl)
			tt.prepare(tSvc)

			// Mock for auth user service
			userSvc := mocks.NewMockUserHandlerService(ctrl)
			userSvc.EXPECT().VerifyAccess(gomock.Any(), gomock.Any()).Return(nil)

			api := API{
				Template: NewTemplateHandler(tSvc, ""),
				User:     NewUserHandler(userSvc),
			}

			router := api.setupRouter(api.buildDefinitionAPI())

			req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/v1/definition/template/%s/%s", tt.adlVer, url.QueryEscape(tt.templateID)), nil)
			req.Header.Set("Authorization", "Bearer "+userAccessKey)
			req.Header.Set("AuthUserId", userID)
			req.Header.Set("EhrSystemId", systemID)
			req.Header.Set("Accept", tt.accept)

			recorder := httptest.NewRecorder()
			router.ServeHTTP(recorder, req)

			resp := recorder.Result()
			defer resp.Body.Close()

			if diff := cmp.Diff(tt.wantStatus, resp.StatusCode); diff != "" {
				t.Errorf("TemplateHandler.GetByID() status code mismatch {-want;+got}\n\t%s", diff)
			}

			if tt.wantStatus != http.StatusOK {
				return
			}

			respBody, _ := io.ReadAll(resp.Body)
			if diff := cmp.Diff(tt.wantResp, string(respBody)); diff != "" {
				t.Errorf("StoredQueryHandler.GetStoredByVersion() status response {-want;+got}\n%s", diff)
			}
		})
	}
}
