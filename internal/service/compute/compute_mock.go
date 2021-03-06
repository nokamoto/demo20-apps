// Code generated by MockGen. DO NOT EDIT.
// Source: compute.go

// Package compute is a generated GoMock package.
package compute

import (
	gomock "github.com/golang/mock/gomock"
	v1alpha "github.com/nokamoto/demo20-apis/cloud/compute/v1alpha"
	reflect "reflect"
)

// Mockcompute is a mock of compute interface
type Mockcompute struct {
	ctrl     *gomock.Controller
	recorder *MockcomputeMockRecorder
}

// MockcomputeMockRecorder is the mock recorder for Mockcompute
type MockcomputeMockRecorder struct {
	mock *Mockcompute
}

// NewMockcompute creates a new mock instance
func NewMockcompute(ctrl *gomock.Controller) *Mockcompute {
	mock := &Mockcompute{ctrl: ctrl}
	mock.recorder = &MockcomputeMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *Mockcompute) EXPECT() *MockcomputeMockRecorder {
	return m.recorder
}

// Create mocks base method
func (m *Mockcompute) Create(arg0, arg1 string, arg2 *v1alpha.Instance) (*v1alpha.Instance, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Create", arg0, arg1, arg2)
	ret0, _ := ret[0].(*v1alpha.Instance)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Create indicates an expected call of Create
func (mr *MockcomputeMockRecorder) Create(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Create", reflect.TypeOf((*Mockcompute)(nil).Create), arg0, arg1, arg2)
}

// RandomName mocks base method
func (m *Mockcompute) RandomName(arg0 string) string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RandomName", arg0)
	ret0, _ := ret[0].(string)
	return ret0
}

// RandomName indicates an expected call of RandomName
func (mr *MockcomputeMockRecorder) RandomName(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RandomName", reflect.TypeOf((*Mockcompute)(nil).RandomName), arg0)
}
