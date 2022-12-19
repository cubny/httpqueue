// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/hibiken/asynq (interfaces: Handler)

// Package mocks is a generated GoMock package.
package mocks

import (
	context "context"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	asynq "github.com/hibiken/asynq"
)

// Handler is a mock of Handler interface.
type Handler struct {
	ctrl     *gomock.Controller
	recorder *HandlerMockRecorder
}

// HandlerMockRecorder is the mock recorder for Handler.
type HandlerMockRecorder struct {
	mock *Handler
}

// NewHandler creates a new mock instance.
func NewHandler(ctrl *gomock.Controller) *Handler {
	mock := &Handler{ctrl: ctrl}
	mock.recorder = &HandlerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *Handler) EXPECT() *HandlerMockRecorder {
	return m.recorder
}

// ProcessTask mocks base method.
func (m *Handler) ProcessTask(arg0 context.Context, arg1 *asynq.Task) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ProcessTask", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// ProcessTask indicates an expected call of ProcessTask.
func (mr *HandlerMockRecorder) ProcessTask(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ProcessTask", reflect.TypeOf((*Handler)(nil).ProcessTask), arg0, arg1)
}