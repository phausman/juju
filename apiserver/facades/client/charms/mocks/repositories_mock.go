// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/juju/juju/apiserver/facades/client/charms (interfaces: CSRepository,CharmHubClient,Strategy,StoreCharm)

// Package mocks is a generated GoMock package.
package mocks

import (
	context "context"
	gomock "github.com/golang/mock/gomock"
	charm "github.com/juju/charm/v8"
	params "github.com/juju/charmrepo/v6/csclient/params"
	charmhub "github.com/juju/juju/charmhub"
	transport "github.com/juju/juju/charmhub/transport"
	charm0 "github.com/juju/juju/core/charm"
	reflect "reflect"
)

// MockCSRepository is a mock of CSRepository interface
type MockCSRepository struct {
	ctrl     *gomock.Controller
	recorder *MockCSRepositoryMockRecorder
}

// MockCSRepositoryMockRecorder is the mock recorder for MockCSRepository
type MockCSRepositoryMockRecorder struct {
	mock *MockCSRepository
}

// NewMockCSRepository creates a new mock instance
func NewMockCSRepository(ctrl *gomock.Controller) *MockCSRepository {
	mock := &MockCSRepository{ctrl: ctrl}
	mock.recorder = &MockCSRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockCSRepository) EXPECT() *MockCSRepositoryMockRecorder {
	return m.recorder
}

// Get mocks base method
func (m *MockCSRepository) Get(arg0 *charm.URL, arg1 string) (*charm.CharmArchive, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Get", arg0, arg1)
	ret0, _ := ret[0].(*charm.CharmArchive)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Get indicates an expected call of Get
func (mr *MockCSRepositoryMockRecorder) Get(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Get", reflect.TypeOf((*MockCSRepository)(nil).Get), arg0, arg1)
}

// ResolveWithPreferredChannel mocks base method
func (m *MockCSRepository) ResolveWithPreferredChannel(arg0 *charm.URL, arg1 params.Channel) (*charm.URL, params.Channel, []string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ResolveWithPreferredChannel", arg0, arg1)
	ret0, _ := ret[0].(*charm.URL)
	ret1, _ := ret[1].(params.Channel)
	ret2, _ := ret[2].([]string)
	ret3, _ := ret[3].(error)
	return ret0, ret1, ret2, ret3
}

// ResolveWithPreferredChannel indicates an expected call of ResolveWithPreferredChannel
func (mr *MockCSRepositoryMockRecorder) ResolveWithPreferredChannel(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ResolveWithPreferredChannel", reflect.TypeOf((*MockCSRepository)(nil).ResolveWithPreferredChannel), arg0, arg1)
}

// MockCharmHubClient is a mock of CharmHubClient interface
type MockCharmHubClient struct {
	ctrl     *gomock.Controller
	recorder *MockCharmHubClientMockRecorder
}

// MockCharmHubClientMockRecorder is the mock recorder for MockCharmHubClient
type MockCharmHubClientMockRecorder struct {
	mock *MockCharmHubClient
}

// NewMockCharmHubClient creates a new mock instance
func NewMockCharmHubClient(ctrl *gomock.Controller) *MockCharmHubClient {
	mock := &MockCharmHubClient{ctrl: ctrl}
	mock.recorder = &MockCharmHubClientMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockCharmHubClient) EXPECT() *MockCharmHubClientMockRecorder {
	return m.recorder
}

// Info mocks base method
func (m *MockCharmHubClient) Info(arg0 context.Context, arg1 string) (transport.InfoResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Info", arg0, arg1)
	ret0, _ := ret[0].(transport.InfoResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Info indicates an expected call of Info
func (mr *MockCharmHubClientMockRecorder) Info(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Info", reflect.TypeOf((*MockCharmHubClient)(nil).Info), arg0, arg1)
}

// Refresh mocks base method
func (m *MockCharmHubClient) Refresh(arg0 context.Context, arg1 charmhub.RefreshConfig) ([]transport.RefreshResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Refresh", arg0, arg1)
	ret0, _ := ret[0].([]transport.RefreshResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Refresh indicates an expected call of Refresh
func (mr *MockCharmHubClientMockRecorder) Refresh(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Refresh", reflect.TypeOf((*MockCharmHubClient)(nil).Refresh), arg0, arg1)
}

// MockStrategy is a mock of Strategy interface
type MockStrategy struct {
	ctrl     *gomock.Controller
	recorder *MockStrategyMockRecorder
}

// MockStrategyMockRecorder is the mock recorder for MockStrategy
type MockStrategyMockRecorder struct {
	mock *MockStrategy
}

// NewMockStrategy creates a new mock instance
func NewMockStrategy(ctrl *gomock.Controller) *MockStrategy {
	mock := &MockStrategy{ctrl: ctrl}
	mock.recorder = &MockStrategyMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockStrategy) EXPECT() *MockStrategyMockRecorder {
	return m.recorder
}

// CharmURL mocks base method
func (m *MockStrategy) CharmURL() *charm.URL {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CharmURL")
	ret0, _ := ret[0].(*charm.URL)
	return ret0
}

// CharmURL indicates an expected call of CharmURL
func (mr *MockStrategyMockRecorder) CharmURL() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CharmURL", reflect.TypeOf((*MockStrategy)(nil).CharmURL))
}

// Finish mocks base method
func (m *MockStrategy) Finish() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Finish")
	ret0, _ := ret[0].(error)
	return ret0
}

// Finish indicates an expected call of Finish
func (mr *MockStrategyMockRecorder) Finish() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Finish", reflect.TypeOf((*MockStrategy)(nil).Finish))
}

// Run mocks base method
func (m *MockStrategy) Run(arg0 charm0.State, arg1 charm0.JujuVersionValidator) (charm0.DownloadResult, bool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Run", arg0, arg1)
	ret0, _ := ret[0].(charm0.DownloadResult)
	ret1, _ := ret[1].(bool)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// Run indicates an expected call of Run
func (mr *MockStrategyMockRecorder) Run(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Run", reflect.TypeOf((*MockStrategy)(nil).Run), arg0, arg1)
}

// Validate mocks base method
func (m *MockStrategy) Validate() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Validate")
	ret0, _ := ret[0].(error)
	return ret0
}

// Validate indicates an expected call of Validate
func (mr *MockStrategyMockRecorder) Validate() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Validate", reflect.TypeOf((*MockStrategy)(nil).Validate))
}

// MockStoreCharm is a mock of StoreCharm interface
type MockStoreCharm struct {
	ctrl     *gomock.Controller
	recorder *MockStoreCharmMockRecorder
}

// MockStoreCharmMockRecorder is the mock recorder for MockStoreCharm
type MockStoreCharmMockRecorder struct {
	mock *MockStoreCharm
}

// NewMockStoreCharm creates a new mock instance
func NewMockStoreCharm(ctrl *gomock.Controller) *MockStoreCharm {
	mock := &MockStoreCharm{ctrl: ctrl}
	mock.recorder = &MockStoreCharmMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockStoreCharm) EXPECT() *MockStoreCharmMockRecorder {
	return m.recorder
}

// Actions mocks base method
func (m *MockStoreCharm) Actions() *charm.Actions {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Actions")
	ret0, _ := ret[0].(*charm.Actions)
	return ret0
}

// Actions indicates an expected call of Actions
func (mr *MockStoreCharmMockRecorder) Actions() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Actions", reflect.TypeOf((*MockStoreCharm)(nil).Actions))
}

// Config mocks base method
func (m *MockStoreCharm) Config() *charm.Config {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Config")
	ret0, _ := ret[0].(*charm.Config)
	return ret0
}

// Config indicates an expected call of Config
func (mr *MockStoreCharmMockRecorder) Config() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Config", reflect.TypeOf((*MockStoreCharm)(nil).Config))
}

// LXDProfile mocks base method
func (m *MockStoreCharm) LXDProfile() *charm.LXDProfile {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "LXDProfile")
	ret0, _ := ret[0].(*charm.LXDProfile)
	return ret0
}

// LXDProfile indicates an expected call of LXDProfile
func (mr *MockStoreCharmMockRecorder) LXDProfile() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "LXDProfile", reflect.TypeOf((*MockStoreCharm)(nil).LXDProfile))
}

// Meta mocks base method
func (m *MockStoreCharm) Meta() *charm.Meta {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Meta")
	ret0, _ := ret[0].(*charm.Meta)
	return ret0
}

// Meta indicates an expected call of Meta
func (mr *MockStoreCharmMockRecorder) Meta() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Meta", reflect.TypeOf((*MockStoreCharm)(nil).Meta))
}

// Metrics mocks base method
func (m *MockStoreCharm) Metrics() *charm.Metrics {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Metrics")
	ret0, _ := ret[0].(*charm.Metrics)
	return ret0
}

// Metrics indicates an expected call of Metrics
func (mr *MockStoreCharmMockRecorder) Metrics() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Metrics", reflect.TypeOf((*MockStoreCharm)(nil).Metrics))
}

// Revision mocks base method
func (m *MockStoreCharm) Revision() int {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Revision")
	ret0, _ := ret[0].(int)
	return ret0
}

// Revision indicates an expected call of Revision
func (mr *MockStoreCharmMockRecorder) Revision() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Revision", reflect.TypeOf((*MockStoreCharm)(nil).Revision))
}

// Version mocks base method
func (m *MockStoreCharm) Version() string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Version")
	ret0, _ := ret[0].(string)
	return ret0
}

// Version indicates an expected call of Version
func (mr *MockStoreCharmMockRecorder) Version() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Version", reflect.TypeOf((*MockStoreCharm)(nil).Version))
}
