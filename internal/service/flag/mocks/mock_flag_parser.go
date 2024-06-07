// Code generated by MockGen. DO NOT EDIT.
// Source: flag_parser.go
//
// Generated by this command:
//
//	mockgen -destination=mocks/mock_flag_parser.go -package=mocks -source=flag_parser.go
//
// Package mocks is a generated GoMock package.
package mocks

import (
	config "project-helper/internal/config"
	reflect "reflect"

	gomock "go.uber.org/mock/gomock"
)

// MockConfigService is a mock of ConfigService interface.
type MockConfigService struct {
	ctrl     *gomock.Controller
	recorder *MockConfigServiceMockRecorder
}

// MockConfigServiceMockRecorder is the mock recorder for MockConfigService.
type MockConfigServiceMockRecorder struct {
	mock *MockConfigService
}

// NewMockConfigService creates a new mock instance.
func NewMockConfigService(ctrl *gomock.Controller) *MockConfigService {
	mock := &MockConfigService{ctrl: ctrl}
	mock.recorder = &MockConfigServiceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockConfigService) EXPECT() *MockConfigServiceMockRecorder {
	return m.recorder
}

// GetConfig mocks base method.
func (m *MockConfigService) GetConfig() *config.Application {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetConfig")
	ret0, _ := ret[0].(*config.Application)
	return ret0
}

// GetConfig indicates an expected call of GetConfig.
func (mr *MockConfigServiceMockRecorder) GetConfig() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetConfig", reflect.TypeOf((*MockConfigService)(nil).GetConfig))
}