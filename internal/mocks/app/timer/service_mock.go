// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/cubny/httpqueue/internal/app/timer (interfaces: Service)

// Package mocks is a generated GoMock package.
package mocks

import (
	context "context"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"

	timer "github.com/cubny/httpqueue/internal/app/timer"
)

// Service is a mock of Service interface.
type Service struct {
	ctrl     *gomock.Controller
	recorder *ServiceMockRecorder
}

// ServiceMockRecorder is the mock recorder for Service.
type ServiceMockRecorder struct {
	mock *Service
}

// NewService creates a new mock instance.
func NewService(ctrl *gomock.Controller) *Service {
	mock := &Service{ctrl: ctrl}
	mock.recorder = &ServiceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *Service) EXPECT() *ServiceMockRecorder {
	return m.recorder
}

// ArchiveTimer mocks base method.
func (m *Service) ArchiveTimer(arg0 context.Context, arg1 string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ArchiveTimer", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// ArchiveTimer indicates an expected call of ArchiveTimer.
func (mr *ServiceMockRecorder) ArchiveTimer(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ArchiveTimer", reflect.TypeOf((*Service)(nil).ArchiveTimer), arg0, arg1)
}

// CreateTimer mocks base method.
func (m *Service) CreateTimer(arg0 context.Context, arg1 timer.SetTimerCommand) (*timer.Timer, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateTimer", arg0, arg1)
	ret0, _ := ret[0].(*timer.Timer)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateTimer indicates an expected call of CreateTimer.
func (mr *ServiceMockRecorder) CreateTimer(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateTimer", reflect.TypeOf((*Service)(nil).CreateTimer), arg0, arg1)
}

// GetTimer mocks base method.
func (m *Service) GetTimer(arg0 context.Context, arg1 string) (*timer.Timer, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetTimer", arg0, arg1)
	ret0, _ := ret[0].(*timer.Timer)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetTimer indicates an expected call of GetTimer.
func (mr *ServiceMockRecorder) GetTimer(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetTimer", reflect.TypeOf((*Service)(nil).GetTimer), arg0, arg1)
}
