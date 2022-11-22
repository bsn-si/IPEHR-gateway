// Code generated by MockGen. DO NOT EDIT.
// Source: stored_query.go

// Package mocks is a generated GoMock package.
package mocks

import (
	context "context"
	model "hms/gateway/pkg/docs/model"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// MockStoredQueryService is a mock of StoredQueryService interface.
type MockStoredQueryService struct {
	ctrl     *gomock.Controller
	recorder *MockStoredQueryServiceMockRecorder
}

// MockStoredQueryServiceMockRecorder is the mock recorder for MockStoredQueryService.
type MockStoredQueryServiceMockRecorder struct {
	mock *MockStoredQueryService
}

// NewMockStoredQueryService creates a new mock instance.
func NewMockStoredQueryService(ctrl *gomock.Controller) *MockStoredQueryService {
	mock := &MockStoredQueryService{ctrl: ctrl}
	mock.recorder = &MockStoredQueryServiceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockStoredQueryService) EXPECT() *MockStoredQueryServiceMockRecorder {
	return m.recorder
}

// Get mocks base method.
func (m *MockStoredQueryService) Get(ctx context.Context, userID, qualifiedQueryName string) ([]model.StoredQuery, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Get", ctx, userID, qualifiedQueryName)
	ret0, _ := ret[0].([]model.StoredQuery)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Get indicates an expected call of Get.
func (mr *MockStoredQueryServiceMockRecorder) Get(ctx, userID, qualifiedQueryName interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Get", reflect.TypeOf((*MockStoredQueryService)(nil).Get), ctx, userID, qualifiedQueryName)
}