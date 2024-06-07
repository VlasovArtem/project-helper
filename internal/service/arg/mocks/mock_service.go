// Code generated by MockGen. DO NOT EDIT.
// Source: service.go
//
// Generated by this command:
//
//	mockgen -destination=mocks/mock_service.go -package=mocks -source=service.go
//
// Package mocks is a generated GoMock package.
package mocks

import (
	config "project-helper/internal/config"
	dto "project-helper/internal/domain/dto"
	entity "project-helper/internal/domain/entity"
	reflect "reflect"

	gomock "go.uber.org/mock/gomock"
)

// MockFlagService is a mock of FlagService interface.
type MockFlagService struct {
	ctrl     *gomock.Controller
	recorder *MockFlagServiceMockRecorder
}

// MockFlagServiceMockRecorder is the mock recorder for MockFlagService.
type MockFlagServiceMockRecorder struct {
	mock *MockFlagService
}

// NewMockFlagService creates a new mock instance.
func NewMockFlagService(ctrl *gomock.Controller) *MockFlagService {
	mock := &MockFlagService{ctrl: ctrl}
	mock.recorder = &MockFlagServiceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockFlagService) EXPECT() *MockFlagServiceMockRecorder {
	return m.recorder
}

// GetOperationFlags mocks base method.
func (m *MockFlagService) GetOperationFlags(operation config.Operation) *entity.Flags {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetOperationFlags", operation)
	ret0, _ := ret[0].(*entity.Flags)
	return ret0
}

// GetOperationFlags indicates an expected call of GetOperationFlags.
func (mr *MockFlagServiceMockRecorder) GetOperationFlags(operation any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetOperationFlags", reflect.TypeOf((*MockFlagService)(nil).GetOperationFlags), operation)
}

// MockEnhanceArgService is a mock of EnhanceArgService interface.
type MockEnhanceArgService struct {
	ctrl     *gomock.Controller
	recorder *MockEnhanceArgServiceMockRecorder
}

// MockEnhanceArgServiceMockRecorder is the mock recorder for MockEnhanceArgService.
type MockEnhanceArgServiceMockRecorder struct {
	mock *MockEnhanceArgService
}

// NewMockEnhanceArgService creates a new mock instance.
func NewMockEnhanceArgService(ctrl *gomock.Controller) *MockEnhanceArgService {
	mock := &MockEnhanceArgService{ctrl: ctrl}
	mock.recorder = &MockEnhanceArgServiceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockEnhanceArgService) EXPECT() *MockEnhanceArgServiceMockRecorder {
	return m.recorder
}

// EnhanceArgs mocks base method.
func (m *MockEnhanceArgService) EnhanceArgs(request *dto.EnhanceArgsRequest) ([]string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "EnhanceArgs", request)
	ret0, _ := ret[0].([]string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// EnhanceArgs indicates an expected call of EnhanceArgs.
func (mr *MockEnhanceArgServiceMockRecorder) EnhanceArgs(request any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "EnhanceArgs", reflect.TypeOf((*MockEnhanceArgService)(nil).EnhanceArgs), request)
}

// GetEnhancedOperationArgs mocks base method.
func (m *MockEnhanceArgService) GetEnhancedOperationArgs(request *dto.GetEnhancedOperationArgs) ([]string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetEnhancedOperationArgs", request)
	ret0, _ := ret[0].([]string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetEnhancedOperationArgs indicates an expected call of GetEnhancedOperationArgs.
func (mr *MockEnhanceArgServiceMockRecorder) GetEnhancedOperationArgs(request any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetEnhancedOperationArgs", reflect.TypeOf((*MockEnhanceArgService)(nil).GetEnhancedOperationArgs), request)
}

// MockPredefinedArgService is a mock of PredefinedArgService interface.
type MockPredefinedArgService struct {
	ctrl     *gomock.Controller
	recorder *MockPredefinedArgServiceMockRecorder
}

// MockPredefinedArgServiceMockRecorder is the mock recorder for MockPredefinedArgService.
type MockPredefinedArgServiceMockRecorder struct {
	mock *MockPredefinedArgService
}

// NewMockPredefinedArgService creates a new mock instance.
func NewMockPredefinedArgService(ctrl *gomock.Controller) *MockPredefinedArgService {
	mock := &MockPredefinedArgService{ctrl: ctrl}
	mock.recorder = &MockPredefinedArgServiceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockPredefinedArgService) EXPECT() *MockPredefinedArgServiceMockRecorder {
	return m.recorder
}

// GetPredefinedArgValues mocks base method.
func (m *MockPredefinedArgService) GetPredefinedArgValues(request *dto.GetPredefinedArgsRequest) ([]string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetPredefinedArgValues", request)
	ret0, _ := ret[0].([]string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetPredefinedArgValues indicates an expected call of GetPredefinedArgValues.
func (mr *MockPredefinedArgServiceMockRecorder) GetPredefinedArgValues(request any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetPredefinedArgValues", reflect.TypeOf((*MockPredefinedArgService)(nil).GetPredefinedArgValues), request)
}