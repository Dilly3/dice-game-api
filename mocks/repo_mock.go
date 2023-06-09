// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/dilly3/dice-game-api/db/sqlc (interfaces: IGameRepo)

// Package mocks is a generated GoMock package.
package mocks

import (
	context "context"
	reflect "reflect"

	db "github.com/dilly3/dice-game-api/db/sqlc"
	gomock "github.com/golang/mock/gomock"
)

// MockIGameRepo is a mock of IGameRepo interface.
type MockIGameRepo struct {
	ctrl     *gomock.Controller
	recorder *MockIGameRepoMockRecorder
}

// MockIGameRepoMockRecorder is the mock recorder for MockIGameRepo.
type MockIGameRepoMockRecorder struct {
	mock *MockIGameRepo
}

// NewMockIGameRepo creates a new mock instance.
func NewMockIGameRepo(ctrl *gomock.Controller) *MockIGameRepo {
	mock := &MockIGameRepo{ctrl: ctrl}
	mock.recorder = &MockIGameRepoMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockIGameRepo) EXPECT() *MockIGameRepoMockRecorder {
	return m.recorder
}

// CreateTransaction mocks base method.
func (m *MockIGameRepo) CreateTransaction(arg0 context.Context, arg1 db.CreateTransactionParams) (db.Transaction, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateTransaction", arg0, arg1)
	ret0, _ := ret[0].(db.Transaction)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateTransaction indicates an expected call of CreateTransaction.
func (mr *MockIGameRepoMockRecorder) CreateTransaction(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateTransaction", reflect.TypeOf((*MockIGameRepo)(nil).CreateTransaction), arg0, arg1)
}

// CreateUser mocks base method.
func (m *MockIGameRepo) CreateUser(arg0 context.Context, arg1 db.CreateUserParams) (db.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateUser", arg0, arg1)
	ret0, _ := ret[0].(db.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateUser indicates an expected call of CreateUser.
func (mr *MockIGameRepoMockRecorder) CreateUser(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateUser", reflect.TypeOf((*MockIGameRepo)(nil).CreateUser), arg0, arg1)
}

// CreateUserTX mocks base method.
func (m *MockIGameRepo) CreateUserTX(arg0 context.Context, arg1 db.CreateUserParams) (db.User, db.Wallet, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateUserTX", arg0, arg1)
	ret0, _ := ret[0].(db.User)
	ret1, _ := ret[1].(db.Wallet)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// CreateUserTX indicates an expected call of CreateUserTX.
func (mr *MockIGameRepoMockRecorder) CreateUserTX(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateUserTX", reflect.TypeOf((*MockIGameRepo)(nil).CreateUserTX), arg0, arg1)
}

// CreateWallet mocks base method.
func (m *MockIGameRepo) CreateWallet(arg0 context.Context, arg1 db.CreateWalletParams) (db.Wallet, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateWallet", arg0, arg1)
	ret0, _ := ret[0].(db.Wallet)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateWallet indicates an expected call of CreateWallet.
func (mr *MockIGameRepoMockRecorder) CreateWallet(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateWallet", reflect.TypeOf((*MockIGameRepo)(nil).CreateWallet), arg0, arg1)
}

// CreditWallet mocks base method.
func (m *MockIGameRepo) CreditWallet(arg0 context.Context, arg1 db.UpdateWalletParams, arg2 bool) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreditWallet", arg0, arg1, arg2)
	ret0, _ := ret[0].(error)
	return ret0
}

// CreditWallet indicates an expected call of CreditWallet.
func (mr *MockIGameRepoMockRecorder) CreditWallet(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreditWallet", reflect.TypeOf((*MockIGameRepo)(nil).CreditWallet), arg0, arg1, arg2)
}

// DebitWallet mocks base method.
func (m *MockIGameRepo) DebitWallet(arg0 context.Context, arg1 db.UpdateWalletParams) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DebitWallet", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// DebitWallet indicates an expected call of DebitWallet.
func (mr *MockIGameRepoMockRecorder) DebitWallet(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DebitWallet", reflect.TypeOf((*MockIGameRepo)(nil).DebitWallet), arg0, arg1)
}

// DeleteUser mocks base method.
func (m *MockIGameRepo) DeleteUser(arg0 context.Context, arg1 string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteUser", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteUser indicates an expected call of DeleteUser.
func (mr *MockIGameRepoMockRecorder) DeleteUser(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteUser", reflect.TypeOf((*MockIGameRepo)(nil).DeleteUser), arg0, arg1)
}

// DeleteWallet mocks base method.
func (m *MockIGameRepo) DeleteWallet(arg0 context.Context, arg1 string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteWallet", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteWallet indicates an expected call of DeleteWallet.
func (mr *MockIGameRepoMockRecorder) DeleteWallet(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteWallet", reflect.TypeOf((*MockIGameRepo)(nil).DeleteWallet), arg0, arg1)
}

// GetTransaction mocks base method.
func (m *MockIGameRepo) GetTransaction(arg0 context.Context, arg1 db.GetTransactionParams) (db.Transaction, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetTransaction", arg0, arg1)
	ret0, _ := ret[0].(db.Transaction)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetTransaction indicates an expected call of GetTransaction.
func (mr *MockIGameRepoMockRecorder) GetTransaction(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetTransaction", reflect.TypeOf((*MockIGameRepo)(nil).GetTransaction), arg0, arg1)
}

// GetTransactionsByUsername mocks base method.
func (m *MockIGameRepo) GetTransactionsByUsername(arg0 context.Context, arg1 db.GetTransactionsByUsernameParams) ([]db.Transaction, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetTransactionsByUsername", arg0, arg1)
	ret0, _ := ret[0].([]db.Transaction)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetTransactionsByUsername indicates an expected call of GetTransactionsByUsername.
func (mr *MockIGameRepoMockRecorder) GetTransactionsByUsername(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetTransactionsByUsername", reflect.TypeOf((*MockIGameRepo)(nil).GetTransactionsByUsername), arg0, arg1)
}

// GetUser mocks base method.
func (m *MockIGameRepo) GetUser(arg0 context.Context, arg1 string) (db.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetUser", arg0, arg1)
	ret0, _ := ret[0].(db.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetUser indicates an expected call of GetUser.
func (mr *MockIGameRepoMockRecorder) GetUser(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUser", reflect.TypeOf((*MockIGameRepo)(nil).GetUser), arg0, arg1)
}

// GetUserByUsername mocks base method.
func (m *MockIGameRepo) GetUserByUsername(arg0 context.Context, arg1 string) (db.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetUserByUsername", arg0, arg1)
	ret0, _ := ret[0].(db.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetUserByUsername indicates an expected call of GetUserByUsername.
func (mr *MockIGameRepoMockRecorder) GetUserByUsername(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUserByUsername", reflect.TypeOf((*MockIGameRepo)(nil).GetUserByUsername), arg0, arg1)
}

// GetUserForUpdate mocks base method.
func (m *MockIGameRepo) GetUserForUpdate(arg0 context.Context, arg1 string) (db.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetUserForUpdate", arg0, arg1)
	ret0, _ := ret[0].(db.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetUserForUpdate indicates an expected call of GetUserForUpdate.
func (mr *MockIGameRepoMockRecorder) GetUserForUpdate(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUserForUpdate", reflect.TypeOf((*MockIGameRepo)(nil).GetUserForUpdate), arg0, arg1)
}

// GetWalletByUsername mocks base method.
func (m *MockIGameRepo) GetWalletByUsername(arg0 context.Context, arg1 string) (db.Wallet, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetWalletByUsername", arg0, arg1)
	ret0, _ := ret[0].(db.Wallet)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetWalletByUsername indicates an expected call of GetWalletByUsername.
func (mr *MockIGameRepoMockRecorder) GetWalletByUsername(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetWalletByUsername", reflect.TypeOf((*MockIGameRepo)(nil).GetWalletByUsername), arg0, arg1)
}

// GetWalletByUsernameForUpdate mocks base method.
func (m *MockIGameRepo) GetWalletByUsernameForUpdate(arg0 context.Context, arg1 string) (db.Wallet, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetWalletByUsernameForUpdate", arg0, arg1)
	ret0, _ := ret[0].(db.Wallet)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetWalletByUsernameForUpdate indicates an expected call of GetWalletByUsernameForUpdate.
func (mr *MockIGameRepoMockRecorder) GetWalletByUsernameForUpdate(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetWalletByUsernameForUpdate", reflect.TypeOf((*MockIGameRepo)(nil).GetWalletByUsernameForUpdate), arg0, arg1)
}

// ListUsers mocks base method.
func (m *MockIGameRepo) ListUsers(arg0 context.Context, arg1 db.ListUsersParams) ([]db.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ListUsers", arg0, arg1)
	ret0, _ := ret[0].([]db.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ListUsers indicates an expected call of ListUsers.
func (mr *MockIGameRepoMockRecorder) ListUsers(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListUsers", reflect.TypeOf((*MockIGameRepo)(nil).ListUsers), arg0, arg1)
}

// UpdateTransaction mocks base method.
func (m *MockIGameRepo) UpdateTransaction(arg0 context.Context, arg1 db.UpdateTransactionParams) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateTransaction", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateTransaction indicates an expected call of UpdateTransaction.
func (mr *MockIGameRepoMockRecorder) UpdateTransaction(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateTransaction", reflect.TypeOf((*MockIGameRepo)(nil).UpdateTransaction), arg0, arg1)
}

// UpdateUserGameMode mocks base method.
func (m *MockIGameRepo) UpdateUserGameMode(arg0 context.Context, arg1 db.UpdateUserGameModeParams) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateUserGameMode", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateUserGameMode indicates an expected call of UpdateUserGameMode.
func (mr *MockIGameRepoMockRecorder) UpdateUserGameMode(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateUserGameMode", reflect.TypeOf((*MockIGameRepo)(nil).UpdateUserGameMode), arg0, arg1)
}

// UpdateWallet mocks base method.
func (m *MockIGameRepo) UpdateWallet(arg0 context.Context, arg1 db.UpdateWalletParams) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateWallet", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateWallet indicates an expected call of UpdateWallet.
func (mr *MockIGameRepoMockRecorder) UpdateWallet(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateWallet", reflect.TypeOf((*MockIGameRepo)(nil).UpdateWallet), arg0, arg1)
}
