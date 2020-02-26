// Code generated by mockery v1.0.0. DO NOT EDIT.

package mocks

import (
	assets "chainlink/core/assets"

	accounts "github.com/ethereum/go-ethereum/accounts"

	big "math/big"

	common "github.com/ethereum/go-ethereum/common"

	decimal "github.com/shopspring/decimal"

	eth "chainlink/core/eth"

	ethereum "github.com/ethereum/go-ethereum"

	mock "github.com/stretchr/testify/mock"

	models "chainlink/core/store/models"

	null "gopkg.in/guregu/null.v3"

	store "chainlink/core/store"
)

// TxManager is an autogenerated mock type for the TxManager type
type TxManager struct {
	mock.Mock
}

// BumpGasUntilSafe provides a mock function with given fields: hash
func (_m *TxManager) BumpGasUntilSafe(hash common.Hash) (*eth.TxReceipt, store.AttemptState, error) {
	ret := _m.Called(hash)

	var r0 *eth.TxReceipt
	if rf, ok := ret.Get(0).(func(common.Hash) *eth.TxReceipt); ok {
		r0 = rf(hash)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*eth.TxReceipt)
		}
	}

	var r1 store.AttemptState
	if rf, ok := ret.Get(1).(func(common.Hash) store.AttemptState); ok {
		r1 = rf(hash)
	} else {
		r1 = ret.Get(1).(store.AttemptState)
	}

	var r2 error
	if rf, ok := ret.Get(2).(func(common.Hash) error); ok {
		r2 = rf(hash)
	} else {
		r2 = ret.Error(2)
	}

	return r0, r1, r2
}

// CheckAttempt provides a mock function with given fields: txAttempt, blockHeight
func (_m *TxManager) CheckAttempt(txAttempt *models.TxAttempt, blockHeight uint64) (*eth.TxReceipt, store.AttemptState, error) {
	ret := _m.Called(txAttempt, blockHeight)

	var r0 *eth.TxReceipt
	if rf, ok := ret.Get(0).(func(*models.TxAttempt, uint64) *eth.TxReceipt); ok {
		r0 = rf(txAttempt, blockHeight)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*eth.TxReceipt)
		}
	}

	var r1 store.AttemptState
	if rf, ok := ret.Get(1).(func(*models.TxAttempt, uint64) store.AttemptState); ok {
		r1 = rf(txAttempt, blockHeight)
	} else {
		r1 = ret.Get(1).(store.AttemptState)
	}

	var r2 error
	if rf, ok := ret.Get(2).(func(*models.TxAttempt, uint64) error); ok {
		r2 = rf(txAttempt, blockHeight)
	} else {
		r2 = ret.Error(2)
	}

	return r0, r1, r2
}

// Connect provides a mock function with given fields: _a0
func (_m *TxManager) Connect(_a0 *models.Head) error {
	ret := _m.Called(_a0)

	var r0 error
	if rf, ok := ret.Get(0).(func(*models.Head) error); ok {
		r0 = rf(_a0)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Connected provides a mock function with given fields:
func (_m *TxManager) Connected() bool {
	ret := _m.Called()

	var r0 bool
	if rf, ok := ret.Get(0).(func() bool); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(bool)
	}

	return r0
}

// ContractLINKBalance provides a mock function with given fields: wr
func (_m *TxManager) ContractLINKBalance(wr models.WithdrawalRequest) (assets.Link, error) {
	ret := _m.Called(wr)

	var r0 assets.Link
	if rf, ok := ret.Get(0).(func(models.WithdrawalRequest) assets.Link); ok {
		r0 = rf(wr)
	} else {
		r0 = ret.Get(0).(assets.Link)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(models.WithdrawalRequest) error); ok {
		r1 = rf(wr)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// CreateTx provides a mock function with given fields: to, data
func (_m *TxManager) CreateTx(to common.Address, data []byte) (*models.Tx, error) {
	ret := _m.Called(to, data)

	var r0 *models.Tx
	if rf, ok := ret.Get(0).(func(common.Address, []byte) *models.Tx); ok {
		r0 = rf(to, data)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*models.Tx)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(common.Address, []byte) error); ok {
		r1 = rf(to, data)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// CreateTxWithEth provides a mock function with given fields: from, to, value
func (_m *TxManager) CreateTxWithEth(from common.Address, to common.Address, value *assets.Eth) (*models.Tx, error) {
	ret := _m.Called(from, to, value)

	var r0 *models.Tx
	if rf, ok := ret.Get(0).(func(common.Address, common.Address, *assets.Eth) *models.Tx); ok {
		r0 = rf(from, to, value)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*models.Tx)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(common.Address, common.Address, *assets.Eth) error); ok {
		r1 = rf(from, to, value)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// CreateTxWithGas provides a mock function with given fields: surrogateID, to, data, gasPriceWei, gasLimit
func (_m *TxManager) CreateTxWithGas(surrogateID null.String, to common.Address, data []byte, gasPriceWei *big.Int, gasLimit uint64) (*models.Tx, error) {
	ret := _m.Called(surrogateID, to, data, gasPriceWei, gasLimit)

	var r0 *models.Tx
	if rf, ok := ret.Get(0).(func(null.String, common.Address, []byte, *big.Int, uint64) *models.Tx); ok {
		r0 = rf(surrogateID, to, data, gasPriceWei, gasLimit)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*models.Tx)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(null.String, common.Address, []byte, *big.Int, uint64) error); ok {
		r1 = rf(surrogateID, to, data, gasPriceWei, gasLimit)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Disconnect provides a mock function with given fields:
func (_m *TxManager) Disconnect() {
	_m.Called()
}

// GetAggregatorLatestRound provides a mock function with given fields: address
func (_m *TxManager) GetAggregatorLatestRound(address common.Address) (*big.Int, error) {
	ret := _m.Called(address)

	var r0 *big.Int
	if rf, ok := ret.Get(0).(func(common.Address) *big.Int); ok {
		r0 = rf(address)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*big.Int)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(common.Address) error); ok {
		r1 = rf(address)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetAggregatorLatestSubmission provides a mock function with given fields: aggregatorAddress, oracleAddress
func (_m *TxManager) GetAggregatorLatestSubmission(aggregatorAddress common.Address, oracleAddress common.Address) (*big.Int, *big.Int, error) {
	ret := _m.Called(aggregatorAddress, oracleAddress)

	var r0 *big.Int
	if rf, ok := ret.Get(0).(func(common.Address, common.Address) *big.Int); ok {
		r0 = rf(aggregatorAddress, oracleAddress)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*big.Int)
		}
	}

	var r1 *big.Int
	if rf, ok := ret.Get(1).(func(common.Address, common.Address) *big.Int); ok {
		r1 = rf(aggregatorAddress, oracleAddress)
	} else {
		if ret.Get(1) != nil {
			r1 = ret.Get(1).(*big.Int)
		}
	}

	var r2 error
	if rf, ok := ret.Get(2).(func(common.Address, common.Address) error); ok {
		r2 = rf(aggregatorAddress, oracleAddress)
	} else {
		r2 = ret.Error(2)
	}

	return r0, r1, r2
}

// GetAggregatorPrice provides a mock function with given fields: address, precision
func (_m *TxManager) GetAggregatorPrice(address common.Address, precision int32) (decimal.Decimal, error) {
	ret := _m.Called(address, precision)

	var r0 decimal.Decimal
	if rf, ok := ret.Get(0).(func(common.Address, int32) decimal.Decimal); ok {
		r0 = rf(address, precision)
	} else {
		r0 = ret.Get(0).(decimal.Decimal)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(common.Address, int32) error); ok {
		r1 = rf(address, precision)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetAggregatorReportingRound provides a mock function with given fields: address
func (_m *TxManager) GetAggregatorReportingRound(address common.Address) (*big.Int, error) {
	ret := _m.Called(address)

	var r0 *big.Int
	if rf, ok := ret.Get(0).(func(common.Address) *big.Int); ok {
		r0 = rf(address)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*big.Int)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(common.Address) error); ok {
		r1 = rf(address)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetAggregatorTimedOutStatus provides a mock function with given fields: address, round
func (_m *TxManager) GetAggregatorTimedOutStatus(address common.Address, round *big.Int) (bool, error) {
	ret := _m.Called(address, round)

	var r0 bool
	if rf, ok := ret.Get(0).(func(common.Address, *big.Int) bool); ok {
		r0 = rf(address, round)
	} else {
		r0 = ret.Get(0).(bool)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(common.Address, *big.Int) error); ok {
		r1 = rf(address, round)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetBlockByNumber provides a mock function with given fields: hex
func (_m *TxManager) GetBlockByNumber(hex string) (eth.BlockHeader, error) {
	ret := _m.Called(hex)

	var r0 eth.BlockHeader
	if rf, ok := ret.Get(0).(func(string) eth.BlockHeader); ok {
		r0 = rf(hex)
	} else {
		r0 = ret.Get(0).(eth.BlockHeader)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(hex)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetChainID provides a mock function with given fields:
func (_m *TxManager) GetChainID() (*big.Int, error) {
	ret := _m.Called()

	var r0 *big.Int
	if rf, ok := ret.Get(0).(func() *big.Int); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*big.Int)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func() error); ok {
		r1 = rf()
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetERC20Balance provides a mock function with given fields: address, contractAddress
func (_m *TxManager) GetERC20Balance(address common.Address, contractAddress common.Address) (*big.Int, error) {
	ret := _m.Called(address, contractAddress)

	var r0 *big.Int
	if rf, ok := ret.Get(0).(func(common.Address, common.Address) *big.Int); ok {
		r0 = rf(address, contractAddress)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*big.Int)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(common.Address, common.Address) error); ok {
		r1 = rf(address, contractAddress)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetEthBalance provides a mock function with given fields: address
func (_m *TxManager) GetEthBalance(address common.Address) (*assets.Eth, error) {
	ret := _m.Called(address)

	var r0 *assets.Eth
	if rf, ok := ret.Get(0).(func(common.Address) *assets.Eth); ok {
		r0 = rf(address)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*assets.Eth)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(common.Address) error); ok {
		r1 = rf(address)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetLINKBalance provides a mock function with given fields: address
func (_m *TxManager) GetLINKBalance(address common.Address) (*assets.Link, error) {
	ret := _m.Called(address)

	var r0 *assets.Link
	if rf, ok := ret.Get(0).(func(common.Address) *assets.Link); ok {
		r0 = rf(address)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*assets.Link)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(common.Address) error); ok {
		r1 = rf(address)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetLogs provides a mock function with given fields: q
func (_m *TxManager) GetLogs(q ethereum.FilterQuery) ([]eth.Log, error) {
	ret := _m.Called(q)

	var r0 []eth.Log
	if rf, ok := ret.Get(0).(func(ethereum.FilterQuery) []eth.Log); ok {
		r0 = rf(q)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]eth.Log)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(ethereum.FilterQuery) error); ok {
		r1 = rf(q)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetNonce provides a mock function with given fields: address
func (_m *TxManager) GetNonce(address common.Address) (uint64, error) {
	ret := _m.Called(address)

	var r0 uint64
	if rf, ok := ret.Get(0).(func(common.Address) uint64); ok {
		r0 = rf(address)
	} else {
		r0 = ret.Get(0).(uint64)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(common.Address) error); ok {
		r1 = rf(address)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetTxReceipt provides a mock function with given fields: hash
func (_m *TxManager) GetTxReceipt(hash common.Hash) (*eth.TxReceipt, error) {
	ret := _m.Called(hash)

	var r0 *eth.TxReceipt
	if rf, ok := ret.Get(0).(func(common.Hash) *eth.TxReceipt); ok {
		r0 = rf(hash)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*eth.TxReceipt)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(common.Hash) error); ok {
		r1 = rf(hash)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// NextActiveAccount provides a mock function with given fields:
func (_m *TxManager) NextActiveAccount() *store.ManagedAccount {
	ret := _m.Called()

	var r0 *store.ManagedAccount
	if rf, ok := ret.Get(0).(func() *store.ManagedAccount); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*store.ManagedAccount)
		}
	}

	return r0
}

// OnNewHead provides a mock function with given fields: _a0
func (_m *TxManager) OnNewHead(_a0 *models.Head) {
	_m.Called(_a0)
}

// Register provides a mock function with given fields: _a0
func (_m *TxManager) Register(_a0 []accounts.Account) {
	_m.Called(_a0)
}

// SendRawTx provides a mock function with given fields: hex
func (_m *TxManager) SendRawTx(hex string) (common.Hash, error) {
	ret := _m.Called(hex)

	var r0 common.Hash
	if rf, ok := ret.Get(0).(func(string) common.Hash); ok {
		r0 = rf(hex)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(common.Hash)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(hex)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// SubscribeToLogs provides a mock function with given fields: channel, q
func (_m *TxManager) SubscribeToLogs(channel chan<- eth.Log, q ethereum.FilterQuery) (eth.Subscription, error) {
	ret := _m.Called(channel, q)

	var r0 eth.Subscription
	if rf, ok := ret.Get(0).(func(chan<- eth.Log, ethereum.FilterQuery) eth.Subscription); ok {
		r0 = rf(channel, q)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(eth.Subscription)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(chan<- eth.Log, ethereum.FilterQuery) error); ok {
		r1 = rf(channel, q)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// SubscribeToNewHeads provides a mock function with given fields: channel
func (_m *TxManager) SubscribeToNewHeads(channel chan<- eth.BlockHeader) (eth.Subscription, error) {
	ret := _m.Called(channel)

	var r0 eth.Subscription
	if rf, ok := ret.Get(0).(func(chan<- eth.BlockHeader) eth.Subscription); ok {
		r0 = rf(channel)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(eth.Subscription)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(chan<- eth.BlockHeader) error); ok {
		r1 = rf(channel)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// WithdrawLINK provides a mock function with given fields: wr
func (_m *TxManager) WithdrawLINK(wr models.WithdrawalRequest) (common.Hash, error) {
	ret := _m.Called(wr)

	var r0 common.Hash
	if rf, ok := ret.Get(0).(func(models.WithdrawalRequest) common.Hash); ok {
		r0 = rf(wr)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(common.Hash)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(models.WithdrawalRequest) error); ok {
		r1 = rf(wr)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}
