// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/cubny/httpqueue/internal/app/timer (interfaces: Outbox)

// Package mocks is a generated GoMock package.
package mocks

import (
	context "context"
	reflect "reflect"

	timer "github.com/cubny/httpqueue/internal/app/timer"
	gomock "github.com/golang/mock/gomock"
)

// Outbox is a mock of Outbox interface.
type Outbox struct {
	ctrl     *gomock.Controller
	recorder *OutboxMockRecorder
}

// OutboxMockRecorder is the mock recorder for Outbox.
type OutboxMockRecorder struct {
	mock *Outbox
}

// NewOutbox creates a new mock instance.
func NewOutbox(ctrl *gomock.Controller) *Outbox {
	mock := &Outbox{ctrl: ctrl}
	mock.recorder = &OutboxMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *Outbox) EXPECT() *OutboxMockRecorder {
	return m.recorder
}

// DequeueOutbox mocks base method.
func (m *Outbox) DequeueOutbox(arg0 context.Context, arg1 int) ([]*timer.Timer, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DequeueOutbox", arg0, arg1)
	ret0, _ := ret[0].([]*timer.Timer)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// DequeueOutbox indicates an expected call of DequeueOutbox.
func (mr *OutboxMockRecorder) DequeueOutbox(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DequeueOutbox", reflect.TypeOf((*Outbox)(nil).DequeueOutbox), arg0, arg1)
}