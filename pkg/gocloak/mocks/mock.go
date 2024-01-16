// Code generated by MockGen. DO NOT EDIT.
// Source: keycloak.go

// Package mock_keycloak is a generated GoMock package.
package mock_keycloak

import (
	context "context"
	reflect "reflect"

	gocloak "github.com/Nerzal/gocloak/v8"
	gomock "github.com/golang/mock/gomock"
	keycloak "gitlab.com/a5805/ondeu/ondeu-back/pkg/gocloak"
)

// MockIClientAuth is a mock of IClientAuth interface.
type MockIClientAuth struct {
	ctrl     *gomock.Controller
	recorder *MockIClientAuthMockRecorder
}

// MockIClientAuthMockRecorder is the mock recorder for MockIClientAuth.
type MockIClientAuthMockRecorder struct {
	mock *MockIClientAuth
}

// NewMockIClientAuth creates a new mock instance.
func NewMockIClientAuth(ctrl *gomock.Controller) *MockIClientAuth {
	mock := &MockIClientAuth{ctrl: ctrl}
	mock.recorder = &MockIClientAuthMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockIClientAuth) EXPECT() *MockIClientAuthMockRecorder {
	return m.recorder
}

// Auth mocks base method.
func (m *MockIClientAuth) Auth(ctx context.Context) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "Auth", ctx)
}

// Auth indicates an expected call of Auth.
func (mr *MockIClientAuthMockRecorder) Auth(ctx interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Auth", reflect.TypeOf((*MockIClientAuth)(nil).Auth), ctx)
}

// Close mocks base method.
func (m *MockIClientAuth) Close() {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "Close")
}

// Close indicates an expected call of Close.
func (mr *MockIClientAuthMockRecorder) Close() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Close", reflect.TypeOf((*MockIClientAuth)(nil).Close))
}

// GetAccessToken mocks base method.
func (m *MockIClientAuth) GetAccessToken(clientId string) string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAccessToken", clientId)
	ret0, _ := ret[0].(string)
	return ret0
}

// GetAccessToken indicates an expected call of GetAccessToken.
func (mr *MockIClientAuthMockRecorder) GetAccessToken(clientId interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAccessToken", reflect.TypeOf((*MockIClientAuth)(nil).GetAccessToken), clientId)
}

// SetClient mocks base method.
func (m *MockIClientAuth) SetClient(clientId, clientSecret string) keycloak.IClientAuth {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SetClient", clientId, clientSecret)
	ret0, _ := ret[0].(keycloak.IClientAuth)
	return ret0
}

// SetClient indicates an expected call of SetClient.
func (mr *MockIClientAuthMockRecorder) SetClient(clientId, clientSecret interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetClient", reflect.TypeOf((*MockIClientAuth)(nil).SetClient), clientId, clientSecret)
}

// MockIKeycloak is a mock of IKeycloak interface.
type MockIKeycloak struct {
	ctrl     *gomock.Controller
	recorder *MockIKeycloakMockRecorder
}

// MockIKeycloakMockRecorder is the mock recorder for MockIKeycloak.
type MockIKeycloakMockRecorder struct {
	mock *MockIKeycloak
}

// NewMockIKeycloak creates a new mock instance.
func NewMockIKeycloak(ctrl *gomock.Controller) *MockIKeycloak {
	mock := &MockIKeycloak{ctrl: ctrl}
	mock.recorder = &MockIKeycloakMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockIKeycloak) EXPECT() *MockIKeycloakMockRecorder {
	return m.recorder
}

// CheckAccessToken mocks base method.
func (m *MockIKeycloak) CheckAccessToken(ctx context.Context, headers map[string][]string, realms []string, resources map[string][]string) (bool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CheckAccessToken", ctx, headers, realms, resources)
	ret0, _ := ret[0].(bool)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CheckAccessToken indicates an expected call of CheckAccessToken.
func (mr *MockIKeycloakMockRecorder) CheckAccessToken(ctx, headers, realms, resources interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CheckAccessToken", reflect.TypeOf((*MockIKeycloak)(nil).CheckAccessToken), ctx, headers, realms, resources)
}

// CheckRoles mocks base method.
func (m *MockIKeycloak) CheckRoles(ctx context.Context, claim map[string]interface{}, realms []string, resources map[string][]string) (bool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CheckRoles", ctx, claim, realms, resources)
	ret0, _ := ret[0].(bool)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CheckRoles indicates an expected call of CheckRoles.
func (mr *MockIKeycloakMockRecorder) CheckRoles(ctx, claim, realms, resources interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CheckRoles", reflect.TypeOf((*MockIKeycloak)(nil).CheckRoles), ctx, claim, realms, resources)
}

// GetRoles mocks base method.
func (m *MockIKeycloak) GetRoles(ctx context.Context, accessToken, clientID string) ([]*gocloak.Role, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetRoles", ctx, accessToken, clientID)
	ret0, _ := ret[0].([]*gocloak.Role)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetRoles indicates an expected call of GetRoles.
func (mr *MockIKeycloakMockRecorder) GetRoles(ctx, accessToken, clientID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetRoles", reflect.TypeOf((*MockIKeycloak)(nil).GetRoles), ctx, accessToken, clientID)
}

// GetUserInfoToken mocks base method.
func (m *MockIKeycloak) GetUserInfoToken(ctx context.Context, accessToken string) (keycloak.UserClaim, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetUserInfoToken", ctx, accessToken)
	ret0, _ := ret[0].(keycloak.UserClaim)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetUserInfoToken indicates an expected call of GetUserInfoToken.
func (mr *MockIKeycloakMockRecorder) GetUserInfoToken(ctx, accessToken interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUserInfoToken", reflect.TypeOf((*MockIKeycloak)(nil).GetUserInfoToken), ctx, accessToken)
}

// ValidateToken mocks base method.
func (m *MockIKeycloak) ValidateToken(ctx context.Context, headers map[string][]string) (bool, map[string]interface{}, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ValidateToken", ctx, headers)
	ret0, _ := ret[0].(bool)
	ret1, _ := ret[1].(map[string]interface{})
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// ValidateToken indicates an expected call of ValidateToken.
func (mr *MockIKeycloakMockRecorder) ValidateToken(ctx, headers interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ValidateToken", reflect.TypeOf((*MockIKeycloak)(nil).ValidateToken), ctx, headers)
}