// Code generated by MockGen. DO NOT EDIT.
// Source: ./internal/wallet_repository.go

// Package mock_internal is a generated GoMock package.
package mock

import (
        context "context"
        sql "database/sql"
        reflect "reflect"

        gomock "github.com/golang/mock/gomock"
        internal "github.com/hokdre/mini-ewallet/internal"
        model "github.com/hokdre/mini-ewallet/internal/model"
)

// MockWalletRepository is a mock of WalletRepository interface.
type MockWalletRepository struct {
        ctrl     *gomock.Controller
        recorder *MockWalletRepositoryMockRecorder
}

// MockWalletRepositoryMockRecorder is the mock recorder for MockWalletRepository.
type MockWalletRepositoryMockRecorder struct {
        mock *MockWalletRepository
}

// NewMockWalletRepository creates a new mock instance.
func NewMockWalletRepository(ctrl *gomock.Controller) *MockWalletRepository {
        mock := &MockWalletRepository{ctrl: ctrl}
        mock.recorder = &MockWalletRepositoryMockRecorder{mock}
        return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockWalletRepository) EXPECT() *MockWalletRepositoryMockRecorder {
        return m.recorder
}

// CreateTx mocks base method.
func (m *MockWalletRepository) CreateTx(ctx context.Context, tx *sql.Tx, newWallet model.Wallet) error {
        m.ctrl.T.Helper()
        ret := m.ctrl.Call(m, "CreateTx", ctx, tx, newWallet)
        ret0, _ := ret[0].(error)
        return ret0
}

// CreateTx indicates an expected call of CreateTx.
func (mr *MockWalletRepositoryMockRecorder) CreateTx(ctx, tx, newWallet interface{}) *gomock.Call {
        mr.mock.ctrl.T.Helper()
        return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateTx", reflect.TypeOf((*MockWalletRepository)(nil).CreateTx), ctx, tx, newWallet)
}

// Decrement mocks base method.
func (m *MockWalletRepository) Decrement(ctx context.Context, tx *sql.Tx, wallet model.Wallet, amount int64) (int64, error) {
        m.ctrl.T.Helper()
        ret := m.ctrl.Call(m, "Decrement", ctx, tx, wallet, amount)
        ret0, _ := ret[0].(int64)
        ret1, _ := ret[1].(error)
        return ret0, ret1
}

// Decrement indicates an expected call of Decrement.
func (mr *MockWalletRepositoryMockRecorder) Decrement(ctx, tx, wallet, amount interface{}) *gomock.Call {
        mr.mock.ctrl.T.Helper()
        return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Decrement", reflect.TypeOf((*MockWalletRepository)(nil).Decrement), ctx, tx, wallet, amount)
}

// GetOne mocks base method.
func (m *MockWalletRepository) GetOne(ctx context.Context, filter internal.WalletFilter) (model.Wallet, error) {
        m.ctrl.T.Helper()
        ret := m.ctrl.Call(m, "GetOne", ctx, filter)
        ret0, _ := ret[0].(model.Wallet)
        ret1, _ := ret[1].(error)
        return ret0, ret1
}

// GetOne indicates an expected call of GetOne.
func (mr *MockWalletRepositoryMockRecorder) GetOne(ctx, filter interface{}) *gomock.Call {
        mr.mock.ctrl.T.Helper()
        return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetOne", reflect.TypeOf((*MockWalletRepository)(nil).GetOne), ctx, filter)
}

// Increment mocks base method.
func (m *MockWalletRepository) Increment(ctx context.Context, tx *sql.Tx, wallet model.Wallet, amount int64) (int64, error) {
        m.ctrl.T.Helper()
        ret := m.ctrl.Call(m, "Increment", ctx, tx, wallet, amount)
        ret0, _ := ret[0].(int64)
        ret1, _ := ret[1].(error)
        return ret0, ret1
}

// Increment indicates an expected call of Increment.
func (mr *MockWalletRepositoryMockRecorder) Increment(ctx, tx, wallet, amount interface{}) *gomock.Call {
        mr.mock.ctrl.T.Helper()
        return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Increment", reflect.TypeOf((*MockWalletRepository)(nil).Increment), ctx, tx, wallet, amount)
}

// Update mocks base method.
func (m *MockWalletRepository) Update(ctx context.Context, wallet model.Wallet) error {
        m.ctrl.T.Helper()
        ret := m.ctrl.Call(m, "Update", ctx, wallet)
        ret0, _ := ret[0].(error)
        return ret0
}

// Update indicates an expected call of Update.
func (mr *MockWalletRepositoryMockRecorder) Update(ctx, wallet interface{}) *gomock.Call {
        mr.mock.ctrl.T.Helper()
        return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Update", reflect.TypeOf((*MockWalletRepository)(nil).Update), ctx, wallet)
}