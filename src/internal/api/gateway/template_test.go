package gateway

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/google/go-cmp/cmp"
	"github.com/google/uuid"

	"github.com/bsn-si/IPEHR-gateway/src/internal/api/gateway/mocks"
	"github.com/bsn-si/IPEHR-gateway/src/pkg/docs/model"
	"github.com/bsn-si/IPEHR-gateway/src/pkg/docs/parser/adl14"
	"github.com/bsn-si/IPEHR-gateway/src/pkg/errors"
)

//
//go:generate mockgen -source template.go -package mocks -destination mocks/template_mock.go
//go:generate mockgen -source user.go -package mocks -destination mocks/user_mock.go

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
				gaSvc.EXPECT().GetByID(gomock.Any(), userID, gomock.Any(), "incorrect").Return(nil, errors.ErrNotFound)
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
				gaSvc.EXPECT().GetByID(gomock.Any(), userID, gomock.Any(), gomock.Any()).Return(template, nil)
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
				gaSvc.EXPECT().GetByID(gomock.Any(), userID, gomock.Any(), gomock.Any()).Return(template, nil)
			},
			http.StatusOK,
			string(template.Body),
		},
	}

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	tSvc := mocks.NewMockTemplateService(ctrl)
	userSvc := mocks.NewMockUserService(ctrl)

	api := API{
		Template: NewTemplateHandler(tSvc, ""),
		User:     NewUserHandler(userSvc),
	}

	router := api.setupRouter(api.buildDefinitionAPI())

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.prepare(tSvc)

			// Mock for auth user service
			userSvc.EXPECT().VerifyAccess(gomock.Any(), gomock.Any()).Return(nil)

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
				t.Errorf("TemplateHandler.GetByID() status response {-want;+got}\n%s", diff)
			}
		})
	}
}

func TestTemplateHandler_ListStored(t *testing.T) {
	var (
		userID   = uuid.New().String()
		systemID = uuid.New().String()
		tM       = make([]*model.Template, 1)
	)

	tM[0] = &model.Template{
		TemplateID:  "1",
		Version:     "1.0.1",
		Concept:     "concept",
		ArchetypeID: "openEHR-EHR-COMPOSITION.encounter.v1",
		CreatedAt:   "2017-07-16T19:20:30.450+01:00",
	}

	tJSON, _ := json.Marshal(tM)

	tests := []struct {
		name       string
		adlVer     string
		prepare    func(gaSvc *mocks.MockTemplateService)
		wantStatus int
		wantResp   string
	}{
		{
			"1. empty result because data was not found",
			model.VerADL1_4,
			func(gaSvc *mocks.MockTemplateService) {
				gaSvc.EXPECT().GetList(gomock.Any(), gomock.Any(), gomock.Any())
			},
			http.StatusOK,
			`[]`,
		},
		{
			"2. success result",
			model.VerADL1_4,
			func(gaSvc *mocks.MockTemplateService) {
				gaSvc.EXPECT().GetList(gomock.Any(), gomock.Any(), gomock.Any()).Return(tM, nil)
			},
			http.StatusOK,
			string(tJSON),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			tSvc := mocks.NewMockTemplateService(ctrl)
			tt.prepare(tSvc)

			// Mock for auth user service
			userSvc := mocks.NewMockUserService(ctrl)
			userSvc.EXPECT().VerifyAccess(gomock.Any(), gomock.Any()).Return(nil)

			api := API{
				Template: NewTemplateHandler(tSvc, ""),
				User:     NewUserHandler(userSvc),
			}

			router := api.setupRouter(api.buildDefinitionAPI())

			req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/v1/definition/template/%s", tt.adlVer), nil)
			req.Header.Set("Authorization", "Bearer emptyJWTkey")
			req.Header.Set("AuthUserId", userID)
			req.Header.Set("EhrSystemId", systemID)

			recorder := httptest.NewRecorder()
			router.ServeHTTP(recorder, req)

			resp := recorder.Result()
			defer resp.Body.Close()

			if diff := cmp.Diff(tt.wantStatus, resp.StatusCode); diff != "" {
				t.Errorf("TemplateHandler.GetList() status code mismatch {-want;+got}\n\t%s", diff)
			}

			if tt.wantStatus == http.StatusNotFound {
				return
			}

			respBody, _ := io.ReadAll(resp.Body)
			if diff := cmp.Diff(tt.wantResp, string(respBody)); diff != "" {
				t.Errorf("TemplateHandler.GetList() status response {-want;+got}\n%s", diff)
			}
		})
	}
}

func TestTemplateHandler_Store(t *testing.T) {
	var (
		userID        = uuid.New().String()
		userAccessKey = "emptyJWTkey"
		systemID      = uuid.New().String()
	)

	template := &model.Template{
		TemplateID:  "Vital Signs",
		Version:     "1",
		VerADL:      model.VerADL1_4,
		MimeType:    model.ADLTypeXML,
		Body:        nil,
		Concept:     "some concept",
		ArchetypeID: "openEHR-EHR-COMPOSITION.encounter.v1",
	}

	template.Body = []byte(`<template xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance" xmlns="http://schemas.openehr.org/v1">
								<language>
									<terminology_id>
										<value>ISO_639-1</value>
									</terminology_id>
									<code_string>en</code_string>
								</language>
								<description>
									<original_author id="Original Author">Not Specified</original_author>
									<lifecycle_state>Initial</lifecycle_state>
									<other_details id="MetaDataSet:Sample Set ">Template metadata sample set</other_details>
									<details>
										<language>
											<terminology_id>
												<value>ISO_639-1</value>
											</terminology_id>
											<code_string>en</code_string>
										</language>
										<purpose>Not Specified</purpose>
									</details>
								</description>
								<uid>
									<value>b4d7f203-b329-4e89-a58a-c605b19e94de</value>
								</uid>
								<template_id>
									<value>` + template.TemplateID + `</value>
								</template_id>
								<concept>` + template.Concept + `</concept>
								<definition archetype_id="` + template.ArchetypeID + `"
									concept_name="Encounter" name="vital_signs2">
									<rm_type_name>COMPOSITION</rm_type_name>
									<occurrences>
										<lower_included>true</lower_included>
										<upper_included>true</upper_included>
										<lower_unbounded>false</lower_unbounded>
										<upper_unbounded>false</upper_unbounded>
										<lower>1</lower>
										<upper>1</upper>
									</occurrences>
									<node_id>at0000</node_id>
								</definition>
							</template>
		`)

	tests := []struct {
		name       string
		body       string
		adlVer     string
		prepare    func(gaSvc *mocks.MockTemplateService)
		wantStatus int
	}{
		{
			"1. empty result because body is empty",
			"",
			model.VerADL1_4,
			func(gaSvc *mocks.MockTemplateService) {},
			http.StatusBadRequest,
		},
		{
			"2. empty result because body is not valid",
			"<xml>...ups",
			model.VerADL1_4,
			func(gaSvc *mocks.MockTemplateService) {
				gaSvc.EXPECT().Parser(model.VerADL1_4).Return(adl14.NewParser(), nil)
			},
			http.StatusBadRequest,
		},
		{
			"3. success",
			string(template.Body),
			model.VerADL1_4,
			func(gaSvc *mocks.MockTemplateService) {
				gaSvc.EXPECT().Parser(model.VerADL1_4).Return(adl14.NewParser(), nil)
				gaSvc.EXPECT().Store(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
			},
			http.StatusCreated,
		},
	}

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	tSvc := mocks.NewMockTemplateService(ctrl)
	userSvc := mocks.NewMockUserService(ctrl)

	api := API{
		Template: NewTemplateHandler(tSvc, ""),
		User:     NewUserHandler(userSvc),
	}

	router := api.setupRouter(api.buildDefinitionAPI())

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.prepare(tSvc)

			// Mock for auth user service
			userSvc.EXPECT().VerifyAccess(gomock.Any(), gomock.Any()).Return(nil)

			reqBody := strings.NewReader(tt.body)

			req := httptest.NewRequest(http.MethodPost, fmt.Sprintf("/v1/definition/template/%s", tt.adlVer), reqBody)
			req.Header.Set("Authorization", "Bearer "+userAccessKey)
			req.Header.Set("AuthUserId", userID)
			req.Header.Set("EhrSystemId", systemID)

			recorder := httptest.NewRecorder()
			router.ServeHTTP(recorder, req)

			resp := recorder.Result()
			defer resp.Body.Close()

			if diff := cmp.Diff(tt.wantStatus, resp.StatusCode); diff != "" {
				t.Errorf("TemplateHandler.GetByID() status code mismatch {-want;+got}\n\t%s", diff)
			}

			if tt.wantStatus != http.StatusCreated {
				return
			}

			if diff := cmp.Diff(resp.Header.Get("Location"), "/definition/template/"+tt.adlVer+"/"+url.QueryEscape(template.TemplateID)); diff != "" {
				t.Errorf("TemplateHandler.store() location {-want;+got}\n%s", diff)
			}
		})
	}
}
