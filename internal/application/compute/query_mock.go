// Code generated by MockGen. DO NOT EDIT.
// Source: query.go

// Package compute is a generated GoMock package.
package compute

import (
	gomock "github.com/golang/mock/gomock"
	gorm "github.com/jinzhu/gorm"
	compute "github.com/nokamoto/demo20-apps/internal/mysql/compute"
	reflect "reflect"
)

// MockinstanceQuery is a mock of instanceQuery interface
type MockinstanceQuery struct {
	ctrl     *gomock.Controller
	recorder *MockinstanceQueryMockRecorder
}

// MockinstanceQueryMockRecorder is the mock recorder for MockinstanceQuery
type MockinstanceQueryMockRecorder struct {
	mock *MockinstanceQuery
}

// NewMockinstanceQuery creates a new mock instance
func NewMockinstanceQuery(ctrl *gomock.Controller) *MockinstanceQuery {
	mock := &MockinstanceQuery{ctrl: ctrl}
	mock.recorder = &MockinstanceQueryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockinstanceQuery) EXPECT() *MockinstanceQueryMockRecorder {
	return m.recorder
}

// Create mocks base method
func (m *MockinstanceQuery) Create(arg0 *gorm.DB, arg1 *compute.Instance) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Create", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// Create indicates an expected call of Create
func (mr *MockinstanceQueryMockRecorder) Create(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Create", reflect.TypeOf((*MockinstanceQuery)(nil).Create), arg0, arg1)
}