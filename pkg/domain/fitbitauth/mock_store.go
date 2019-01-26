// Code generated by MockGen. DO NOT EDIT.
// Source: store.go

// Package fitbitauth is a generated GoMock package.
package fitbitauth

import (
	gomock "github.com/golang/mock/gomock"
	oauth2 "golang.org/x/oauth2"
	reflect "reflect"
)

// MockStore is a mock of Store interface
type MockStore struct {
	ctrl     *gomock.Controller
	recorder *MockStoreMockRecorder
}

// MockStoreMockRecorder is the mock recorder for MockStore
type MockStoreMockRecorder struct {
	mock *MockStore
}

// NewMockStore creates a new mock instance
func NewMockStore(ctrl *gomock.Controller) *MockStore {
	mock := &MockStore{ctrl: ctrl}
	mock.recorder = &MockStoreMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockStore) EXPECT() *MockStoreMockRecorder {
	return m.recorder
}

// WriteFitbitToken mocks base method
func (m *MockStore) WriteFitbitToken(token *oauth2.Token) error {
	ret := m.ctrl.Call(m, "WriteFitbitToken", token)
	ret0, _ := ret[0].(error)
	return ret0
}

// WriteFitbitToken indicates an expected call of WriteFitbitToken
func (mr *MockStoreMockRecorder) WriteFitbitToken(token interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "WriteFitbitToken", reflect.TypeOf((*MockStore)(nil).WriteFitbitToken), token)
}

// FetchFitbitToken mocks base method
func (m *MockStore) FetchFitbitToken() (*oauth2.Token, error) {
	ret := m.ctrl.Call(m, "FetchFitbitToken")
	ret0, _ := ret[0].(*oauth2.Token)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FetchFitbitToken indicates an expected call of FetchFitbitToken
func (mr *MockStoreMockRecorder) FetchFitbitToken() *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FetchFitbitToken", reflect.TypeOf((*MockStore)(nil).FetchFitbitToken))
}
