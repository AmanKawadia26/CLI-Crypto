// Code generated by MockGen. DO NOT EDIT.
// Source: C:\Users\akawadia\Downloads\CryptoTracker\internal\services\user_services.go

// Package mock_services is a generated GoMock package.
package mock_services

import (
	models "cryptotracker/models"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// MockUserServices is a mock of UserServices interface.
type MockUserServices struct {
	ctrl     *gomock.Controller
	recorder *MockUserServicesMockRecorder
}

// MockUserServicesMockRecorder is the mock recorder for MockUserServices.
type MockUserServicesMockRecorder struct {
	mock *MockUserServices
}

// NewMockUserServices creates a new mock instance.
func NewMockUserServices(ctrl *gomock.Controller) *MockUserServices {
	mock := &MockUserServices{ctrl: ctrl}
	mock.recorder = &MockUserServicesMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockUserServices) EXPECT() *MockUserServicesMockRecorder {
	return m.recorder
}

// GetUserProfile mocks base method.
func (m *MockUserServices) GetUserProfile(username string) (*models.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetUserProfile", username)
	ret0, _ := ret[0].(*models.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetUserProfile indicates an expected call of GetUserProfile.
func (mr *MockUserServicesMockRecorder) GetUserProfile(username interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUserProfile", reflect.TypeOf((*MockUserServices)(nil).GetUserProfile), username)
}
