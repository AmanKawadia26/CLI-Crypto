// Code generated by MockGen. DO NOT EDIT.
// Source: C:\Users\akawadia\Downloads\CryptoTracker\internal\repositories\notification_repository.go

// Package mock_repositories is a generated GoMock package.
package mock_repositories

import (
	models "cryptotracker/models"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	pgx "github.com/jackc/pgx/v4"
)

// MockNotificationRepository is a mock of NotificationRepository interface.
type MockNotificationRepository struct {
	ctrl     *gomock.Controller
	recorder *MockNotificationRepositoryMockRecorder
}

// MockNotificationRepositoryMockRecorder is the mock recorder for MockNotificationRepository.
type MockNotificationRepositoryMockRecorder struct {
	mock *MockNotificationRepository
}

// NewMockNotificationRepository creates a new mock instance.
func NewMockNotificationRepository(ctrl *gomock.Controller) *MockNotificationRepository {
	mock := &MockNotificationRepository{ctrl: ctrl}
	mock.recorder = &MockNotificationRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockNotificationRepository) EXPECT() *MockNotificationRepositoryMockRecorder {
	return m.recorder
}

// CheckPriceAlertsRepo mocks base method.
func (m *MockNotificationRepository) CheckPriceAlertsRepo(username string) ([]models.PriceNotification, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CheckPriceAlertsRepo", username)
	ret0, _ := ret[0].([]models.PriceNotification)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CheckPriceAlertsRepo indicates an expected call of CheckPriceAlertsRepo.
func (mr *MockNotificationRepositoryMockRecorder) CheckPriceAlertsRepo(username interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CheckPriceAlertsRepo", reflect.TypeOf((*MockNotificationRepository)(nil).CheckPriceAlertsRepo), username)
}

// CheckUnavailableCryptoRequestsRepo mocks base method.
func (m *MockNotificationRepository) CheckUnavailableCryptoRequestsRepo(username string) (pgx.Rows, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CheckUnavailableCryptoRequestsRepo", username)
	ret0, _ := ret[0].(pgx.Rows)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CheckUnavailableCryptoRequestsRepo indicates an expected call of CheckUnavailableCryptoRequestsRepo.
func (mr *MockNotificationRepositoryMockRecorder) CheckUnavailableCryptoRequestsRepo(username interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CheckUnavailableCryptoRequestsRepo", reflect.TypeOf((*MockNotificationRepository)(nil).CheckUnavailableCryptoRequestsRepo), username)
}

// UpdatePriceNotificationStatusRepo mocks base method.
func (m *MockNotificationRepository) UpdatePriceNotificationStatusRepo(notification *models.PriceNotification, username string, currentPrice float64) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdatePriceNotificationStatusRepo", notification, username, currentPrice)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdatePriceNotificationStatusRepo indicates an expected call of UpdatePriceNotificationStatusRepo.
func (mr *MockNotificationRepositoryMockRecorder) UpdatePriceNotificationStatusRepo(notification, username, currentPrice interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdatePriceNotificationStatusRepo", reflect.TypeOf((*MockNotificationRepository)(nil).UpdatePriceNotificationStatusRepo), notification, username, currentPrice)
}