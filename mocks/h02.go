// Code generated by MockGen. DO NOT EDIT.
// Source: reciever-ms/protocol/h02/decoder (interfaces: IDecoder)

// Package mocks is a generated GoMock package.
package mocks

import (
	context "context"
	protocol "reciever-ms/protocol"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// MockIDecoder is a mock of IDecoder interface.
type MockIDecoder struct {
	ctrl     *gomock.Controller
	recorder *MockIDecoderMockRecorder
}

// MockIDecoderMockRecorder is the mock recorder for MockIDecoder.
type MockIDecoderMockRecorder struct {
	mock *MockIDecoder
}

// NewMockIDecoder creates a new mock instance.
func NewMockIDecoder(ctrl *gomock.Controller) *MockIDecoder {
	mock := &MockIDecoder{ctrl: ctrl}
	mock.recorder = &MockIDecoderMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockIDecoder) EXPECT() *MockIDecoderMockRecorder {
	return m.recorder
}

// Decode mocks base method.
func (m *MockIDecoder) Decode(arg0 context.Context, arg1 []byte) (*protocol.DecodeResult, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Decode", arg0, arg1)
	ret0, _ := ret[0].(*protocol.DecodeResult)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Decode indicates an expected call of Decode.
func (mr *MockIDecoderMockRecorder) Decode(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Decode", reflect.TypeOf((*MockIDecoder)(nil).Decode), arg0, arg1)
}
