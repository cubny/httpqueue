// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/cubny/httpqueue/internal/app/timer (interfaces: Producer)

// Package mocks is a generated GoMock package.
package mocks

import (
	context "context"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"

	timer "github.com/cubny/httpqueue/internal/app/timer"
)

// Producer is a mock of Producer interface.
type Producer struct {
	ctrl     *gomock.Controller
	recorder *ProducerMockRecorder
}

// ProducerMockRecorder is the mock recorder for Producer.
type ProducerMockRecorder struct {
	mock *Producer
}

// NewProducer creates a new mock instance.
func NewProducer(ctrl *gomock.Controller) *Producer {
	mock := &Producer{ctrl: ctrl}
	mock.recorder = &ProducerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *Producer) EXPECT() *ProducerMockRecorder {
	return m.recorder
}

// Send mocks base method.
func (m *Producer) Send(arg0 context.Context, arg1 *timer.Timer) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Send", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// Send indicates an expected call of Send.
func (mr *ProducerMockRecorder) Send(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Send", reflect.TypeOf((*Producer)(nil).Send), arg0, arg1)
}
