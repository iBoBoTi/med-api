// Code generated by MockGen. DO NOT EDIT.
// Source: db/auth_repository.go

// Package mocks is a generated GoMock package.
package mocks

import (
	reflect "reflect"

	models "github.com/decagonhq/meddle-api/models"
	gomock "github.com/golang/mock/gomock"
)

// MockAuthRepository is a mock of AuthRepository interface.
type MockAuthRepository struct {
	ctrl     *gomock.Controller
	recorder *MockAuthRepositoryMockRecorder
}

// MockAuthRepositoryMockRecorder is the mock recorder for MockAuthRepository.
type MockAuthRepositoryMockRecorder struct {
	mock *MockAuthRepository
}

// NewMockAuthRepository creates a new mock instance.
func NewMockAuthRepository(ctrl *gomock.Controller) *MockAuthRepository {
	mock := &MockAuthRepository{ctrl: ctrl}
	mock.recorder = &MockAuthRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockAuthRepository) EXPECT() *MockAuthRepositoryMockRecorder {
	return m.recorder
}

// AddToBlackList mocks base method.
func (m *MockAuthRepository) AddToBlackList(blacklist *models.BlackList) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AddToBlackList", blacklist)
	ret0, _ := ret[0].(error)
	return ret0
}

// AddToBlackList indicates an expected call of AddToBlackList.
func (mr *MockAuthRepositoryMockRecorder) AddToBlackList(blacklist interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddToBlackList", reflect.TypeOf((*MockAuthRepository)(nil).AddToBlackList), blacklist)
}

// CreateUser mocks base method.
func (m *MockAuthRepository) CreateUser(user *models.User) (*models.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateUser", user)
	ret0, _ := ret[0].(*models.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateUser indicates an expected call of CreateUser.
func (mr *MockAuthRepositoryMockRecorder) CreateUser(user interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateUser", reflect.TypeOf((*MockAuthRepository)(nil).CreateUser), user)
}

// FindUserByEmailOrPhoneNumber mocks base method.
func (m *MockAuthRepository) FindUserByEmailOrPhoneNumber(email, phoneNumber string) (*models.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FindUserByEmailOrPhoneNumber", email, phoneNumber)
	ret0, _ := ret[0].(*models.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FindUserByEmailOrPhoneNumber indicates an expected call of FindUserByEmailOrPhoneNumber.
func (mr *MockAuthRepositoryMockRecorder) FindUserByEmailOrPhoneNumber(email, phoneNumber interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindUserByEmailOrPhoneNumber", reflect.TypeOf((*MockAuthRepository)(nil).FindUserByEmailOrPhoneNumber), email, phoneNumber)
}

// FindUserByUsername mocks base method.
func (m *MockAuthRepository) FindUserByUsername(username string) (*models.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FindUserByUsername", username)
	ret0, _ := ret[0].(*models.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FindUserByUsername indicates an expected call of FindUserByUsername.
func (mr *MockAuthRepositoryMockRecorder) FindUserByUsername(username interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindUserByUsername", reflect.TypeOf((*MockAuthRepository)(nil).FindUserByUsername), username)
}

// IsEmailExist mocks base method.
func (m *MockAuthRepository) IsEmailExist(email string) (bool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "IsEmailExist", email)
	ret0, _ := ret[0].(bool)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// IsEmailExist indicates an expected call of IsEmailExist.
func (mr *MockAuthRepositoryMockRecorder) IsEmailExist(email interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "IsEmailExist", reflect.TypeOf((*MockAuthRepository)(nil).IsEmailExist), email)
}

// IsPhoneExist mocks base method.
func (m *MockAuthRepository) IsPhoneExist(email string) (bool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "IsPhoneExist", email)
	ret0, _ := ret[0].(bool)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// IsPhoneExist indicates an expected call of IsPhoneExist.
func (mr *MockAuthRepositoryMockRecorder) IsPhoneExist(email interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "IsPhoneExist", reflect.TypeOf((*MockAuthRepository)(nil).IsPhoneExist), email)
}

// TokenInBlacklist mocks base method.
func (m *MockAuthRepository) TokenInBlacklist(token *string) bool {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "TokenInBlacklist", token)
	ret0, _ := ret[0].(bool)
	return ret0
}

// TokenInBlacklist indicates an expected call of TokenInBlacklist.
func (mr *MockAuthRepositoryMockRecorder) TokenInBlacklist(token interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "TokenInBlacklist", reflect.TypeOf((*MockAuthRepository)(nil).TokenInBlacklist), token)
}

// UpdateUser mocks base method.
func (m *MockAuthRepository) UpdateUser(user *models.User) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateUser", user)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateUser indicates an expected call of UpdateUser.
func (mr *MockAuthRepositoryMockRecorder) UpdateUser(user interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateUser", reflect.TypeOf((*MockAuthRepository)(nil).UpdateUser), user)
}
