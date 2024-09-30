// Code generated by MockGen. DO NOT EDIT.
// Source: C:\Users\akawadia\Downloads\CryptoTracker\internal\services\crypto_services.go

// Package mock_services is a generated GoMock package.
package mock_services

import (
	models "cryptotracker/models"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// MockCryptoService is a mock of CryptoService interface.
type MockCryptoService struct {
	ctrl     *gomock.Controller
	recorder *MockCryptoServiceMockRecorder
}

// MockCryptoServiceMockRecorder is the mock recorder for MockCryptoService.
type MockCryptoServiceMockRecorder struct {
	mock *MockCryptoService
}

// NewMockCryptoService creates a new mock instance.
func NewMockCryptoService(ctrl *gomock.Controller) *MockCryptoService {
	mock := &MockCryptoService{ctrl: ctrl}
	mock.recorder = &MockCryptoServiceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockCryptoService) EXPECT() *MockCryptoServiceMockRecorder {
	return m.recorder
}

// DisplayTopCryptocurrencies mocks base method.
func (m *MockCryptoService) DisplayTopCryptocurrencies(count int) ([]interface{}, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DisplayTopCryptocurrencies", count)
	ret0, _ := ret[0].([]interface{})
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// DisplayTopCryptocurrencies indicates an expected call of DisplayTopCryptocurrencies.
func (mr *MockCryptoServiceMockRecorder) DisplayTopCryptocurrencies(count interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DisplayTopCryptocurrencies", reflect.TypeOf((*MockCryptoService)(nil).DisplayTopCryptocurrencies), count)
}

// SearchCryptocurrency mocks base method.
func (m *MockCryptoService) SearchCryptocurrency(user *models.User, cryptoSymbol string) (float64, string, string, *models.Cryptocurrency, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SearchCryptocurrency", user, cryptoSymbol)
	ret0, _ := ret[0].(float64)
	ret1, _ := ret[1].(string)
	ret2, _ := ret[2].(string)
	ret3, _ := ret[3].(*models.Cryptocurrency)
	ret4, _ := ret[4].(error)
	return ret0, ret1, ret2, ret3, ret4
}

// SearchCryptocurrency indicates an expected call of SearchCryptocurrency.
func (mr *MockCryptoServiceMockRecorder) SearchCryptocurrency(user, cryptoSymbol interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SearchCryptocurrency", reflect.TypeOf((*MockCryptoService)(nil).SearchCryptocurrency), user, cryptoSymbol)
}

// SetPriceAlert mocks base method.
func (m *MockCryptoService) SetPriceAlert(user *models.User, symbol string, targetPrice float64) (float64, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SetPriceAlert", user, symbol, targetPrice)
	ret0, _ := ret[0].(float64)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// SetPriceAlert indicates an expected call of SetPriceAlert.
func (mr *MockCryptoServiceMockRecorder) SetPriceAlert(user, symbol, targetPrice interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetPriceAlert", reflect.TypeOf((*MockCryptoService)(nil).SetPriceAlert), user, symbol, targetPrice)
}
