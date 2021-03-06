// Code generated by MockGen. DO NOT EDIT.
// Source: resourcemanager.go

// Package resourcemanager is a generated GoMock package.
package resourcemanager

import (
	gomock "github.com/golang/mock/gomock"
	v1alpha "github.com/nokamoto/demo20-apis/cloud/resourcemanager/v1alpha"
	reflect "reflect"
)

// Mockresourcemanager is a mock of resourcemanager interface
type Mockresourcemanager struct {
	ctrl     *gomock.Controller
	recorder *MockresourcemanagerMockRecorder
}

// MockresourcemanagerMockRecorder is the mock recorder for Mockresourcemanager
type MockresourcemanagerMockRecorder struct {
	mock *Mockresourcemanager
}

// NewMockresourcemanager creates a new mock instance
func NewMockresourcemanager(ctrl *gomock.Controller) *Mockresourcemanager {
	mock := &Mockresourcemanager{ctrl: ctrl}
	mock.recorder = &MockresourcemanagerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *Mockresourcemanager) EXPECT() *MockresourcemanagerMockRecorder {
	return m.recorder
}

// Get mocks base method
func (m *Mockresourcemanager) Get(arg0 string) (*v1alpha.Project, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Get", arg0)
	ret0, _ := ret[0].(*v1alpha.Project)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Get indicates an expected call of Get
func (mr *MockresourcemanagerMockRecorder) Get(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Get", reflect.TypeOf((*Mockresourcemanager)(nil).Get), arg0)
}

// Create mocks base method
func (m *Mockresourcemanager) Create(id string, project *v1alpha.Project) (*v1alpha.Project, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Create", id, project)
	ret0, _ := ret[0].(*v1alpha.Project)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Create indicates an expected call of Create
func (mr *MockresourcemanagerMockRecorder) Create(id, project interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Create", reflect.TypeOf((*Mockresourcemanager)(nil).Create), id, project)
}
