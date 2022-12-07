// Code generated by MockGen. DO NOT EDIT.
// Source: template.go

// Package mocks is a generated GoMock package.
package mocks

import (
	context "context"
	model "hms/gateway/pkg/docs/model"
	template "hms/gateway/pkg/docs/service/template"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// MockTemplateService is a mock of TemplateService interface.
type MockTemplateService struct {
	ctrl     *gomock.Controller
	recorder *MockTemplateServiceMockRecorder
}

// MockTemplateServiceMockRecorder is the mock recorder for MockTemplateService.
type MockTemplateServiceMockRecorder struct {
	mock *MockTemplateService
}

// NewMockTemplateService creates a new mock instance.
func NewMockTemplateService(ctrl *gomock.Controller) *MockTemplateService {
	mock := &MockTemplateService{ctrl: ctrl}
	mock.recorder = &MockTemplateServiceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockTemplateService) EXPECT() *MockTemplateServiceMockRecorder {
	return m.recorder
}

// GetByID mocks base method.
func (m *MockTemplateService) GetByID(ctx context.Context, userID, templateID string) (*model.Template, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetByID", ctx, userID, templateID)
	ret0, _ := ret[0].(*model.Template)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetByID indicates an expected call of GetByID.
func (mr *MockTemplateServiceMockRecorder) GetByID(ctx, userID, templateID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetByID", reflect.TypeOf((*MockTemplateService)(nil).GetByID), ctx, userID, templateID)
}

// GetList mocks base method.
func (m *MockTemplateService) GetList(ctx context.Context, userID, systemID string) ([]*model.TemplateResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetList", ctx, userID, systemID)
	ret0, _ := ret[0].([]*model.TemplateResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetList indicates an expected call of GetList.
func (mr *MockTemplateServiceMockRecorder) GetList(ctx, userID, systemID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetList", reflect.TypeOf((*MockTemplateService)(nil).GetList), ctx, userID, systemID)
}

// Parser mocks base method.
func (m *MockTemplateService) Parser(version model.ADLVer) (template.ADLParser, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Parser", version)
	ret0, _ := ret[0].(template.ADLParser)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Parser indicates an expected call of Parser.
func (mr *MockTemplateServiceMockRecorder) Parser(version interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Parser", reflect.TypeOf((*MockTemplateService)(nil).Parser), version)
}

// Store mocks base method.
func (m_2 *MockTemplateService) Store(ctx context.Context, userID, systemID, reqID string, m *model.Template) error {
	m_2.ctrl.T.Helper()
	ret := m_2.ctrl.Call(m_2, "Store", ctx, userID, systemID, reqID, m)
	ret0, _ := ret[0].(error)
	return ret0
}

// Store indicates an expected call of Store.
func (mr *MockTemplateServiceMockRecorder) Store(ctx, userID, systemID, reqID, m interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Store", reflect.TypeOf((*MockTemplateService)(nil).Store), ctx, userID, systemID, reqID, m)
}
