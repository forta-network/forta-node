// Code generated by MockGen. DO NOT EDIT.
// Source: clients/ratelimiter/rate_limiter.go

// Package mock_ratelimiter is a generated GoMock package.
package mock_ratelimiter

import (
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// MockRateLimiter is a mock of RateLimiter interface.
type MockRateLimiter struct {
	ctrl     *gomock.Controller
	recorder *MockRateLimiterMockRecorder
}

// MockRateLimiterMockRecorder is the mock recorder for MockRateLimiter.
type MockRateLimiterMockRecorder struct {
	mock *MockRateLimiter
}

// NewMockRateLimiter creates a new mock instance.
func NewMockRateLimiter(ctrl *gomock.Controller) *MockRateLimiter {
	mock := &MockRateLimiter{ctrl: ctrl}
	mock.recorder = &MockRateLimiterMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockRateLimiter) EXPECT() *MockRateLimiterMockRecorder {
	return m.recorder
}

// ExceedsLimit mocks base method.
func (m *MockRateLimiter) ExceedsLimit(clientID string) bool {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ExceedsLimit", clientID)
	ret0, _ := ret[0].(bool)
	return ret0
}

// ExceedsLimit indicates an expected call of ExceedsLimit.
func (mr *MockRateLimiterMockRecorder) ExceedsLimit(clientID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ExceedsLimit", reflect.TypeOf((*MockRateLimiter)(nil).ExceedsLimit), clientID)
}