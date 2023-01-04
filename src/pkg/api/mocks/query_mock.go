// Code generated by MockGen. DO NOT EDIT.
// Source: query.go

// Package mocks is a generated GoMock package.
package mocks

import (
	context "context"
	model "hms/gateway/pkg/docs/model"
	base "hms/gateway/pkg/docs/model/base"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// MockQueryService is a mock of QueryService interface.
type MockQueryService struct {
	ctrl     *gomock.Controller
	recorder *MockQueryServiceMockRecorder
}

// MockQueryServiceMockRecorder is the mock recorder for MockQueryService.
type MockQueryServiceMockRecorder struct {
	mock *MockQueryService
}

// NewMockQueryService creates a new mock instance.
func NewMockQueryService(ctrl *gomock.Controller) *MockQueryService {
	mock := &MockQueryService{ctrl: ctrl}
	mock.recorder = &MockQueryServiceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockQueryService) EXPECT() *MockQueryServiceMockRecorder {
	return m.recorder
}

// ExecStoredQuery mocks base method.
func (m *MockQueryService) ExecStoredQuery(ctx context.Context, qualifiedQueryName string, query *model.QueryRequest) (*model.QueryResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ExecStoredQuery", ctx, qualifiedQueryName, query)
	ret0, _ := ret[0].(*model.QueryResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ExecStoredQuery indicates an expected call of ExecStoredQuery.
func (mr *MockQueryServiceMockRecorder) ExecStoredQuery(ctx, qualifiedQueryName, query interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ExecStoredQuery", reflect.TypeOf((*MockQueryService)(nil).ExecStoredQuery), ctx, qualifiedQueryName, query)
}

// GetByVersion mocks base method.
func (m *MockQueryService) GetByVersion(ctx context.Context, userID, systemID, name string, version *base.VersionTreeID) (*model.StoredQuery, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetByVersion", ctx, userID, systemID, name, version)
	ret0, _ := ret[0].(*model.StoredQuery)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetByVersion indicates an expected call of GetByVersion.
func (mr *MockQueryServiceMockRecorder) GetByVersion(ctx, userID, systemID, name, version interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetByVersion", reflect.TypeOf((*MockQueryService)(nil).GetByVersion), ctx, userID, systemID, name, version)
}

// List mocks base method.
func (m *MockQueryService) List(ctx context.Context, userID, systemID, qualifiedQueryName string) ([]*model.StoredQuery, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "List", ctx, userID, systemID, qualifiedQueryName)
	ret0, _ := ret[0].([]*model.StoredQuery)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// List indicates an expected call of List.
func (mr *MockQueryServiceMockRecorder) List(ctx, userID, systemID, qualifiedQueryName interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "List", reflect.TypeOf((*MockQueryService)(nil).List), ctx, userID, systemID, qualifiedQueryName)
}

// Store mocks base method.
func (m *MockQueryService) Store(ctx context.Context, userID, systemID, reqID, qType, name, q string) (*model.StoredQuery, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Store", ctx, userID, systemID, reqID, qType, name, q)
	ret0, _ := ret[0].(*model.StoredQuery)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Store indicates an expected call of Store.
func (mr *MockQueryServiceMockRecorder) Store(ctx, userID, systemID, reqID, qType, name, q interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Store", reflect.TypeOf((*MockQueryService)(nil).Store), ctx, userID, systemID, reqID, qType, name, q)
}

// StoreVersion mocks base method.
func (m *MockQueryService) StoreVersion(ctx context.Context, userID, systemID, reqID, qType, name string, version *base.VersionTreeID, q string) (*model.StoredQuery, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "StoreVersion", ctx, userID, systemID, reqID, qType, name, version, q)
	ret0, _ := ret[0].(*model.StoredQuery)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// StoreVersion indicates an expected call of StoreVersion.
func (mr *MockQueryServiceMockRecorder) StoreVersion(ctx, userID, systemID, reqID, qType, name, version, q interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "StoreVersion", reflect.TypeOf((*MockQueryService)(nil).StoreVersion), ctx, userID, systemID, reqID, qType, name, version, q)
}

// Validate mocks base method.
func (m *MockQueryService) Validate(data []byte) bool {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Validate", data)
	ret0, _ := ret[0].(bool)
	return ret0
}

// Validate indicates an expected call of Validate.
func (mr *MockQueryServiceMockRecorder) Validate(data interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Validate", reflect.TypeOf((*MockQueryService)(nil).Validate), data)
}
