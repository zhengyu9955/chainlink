package mocks

import (
	"errors"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/smartcontractkit/aggregator-monitor/client"
	"math/big"
	"math/rand"
	"time"
)

type AggregatorMockInit func() *AggregatorMock

func NewAggregatorOracle() *client.AggregatorOracle {
	return &client.AggregatorOracle{
		Name:    "Unknown",
		Address: common.BigToAddress(big.NewInt(rand.Int63())),
	}
}

type AggregatorMock struct {
	address common.Address
	oracles client.OracleMapping
}

func NewAggregatorMock() *AggregatorMock {
	return &AggregatorMock{
		address: common.BigToAddress(big.NewInt(rand.Int63())),
		oracles: GenerateOracles(),
	}
}

func (a *AggregatorMock) Name() string {
	return "Mock"
}

func (a *AggregatorMock) Address() common.Address {
	return a.address
}

func (a *AggregatorMock) LatestAnswer() (*big.Int, error) {
	return big.NewInt(50), nil
}

func (a *AggregatorMock) LatestRound() (*big.Int, error) {
	return big.NewInt(100), nil
}

func (a *AggregatorMock) Oracles() (client.OracleMapping, error) {
	return a.oracles, nil
}

func (a *AggregatorMock) LogToTransaction(log types.Log) (*types.Transaction, error) {
	return types.NewTransaction(0, common.Address{}, big.NewInt(0), 0, big.NewInt(4000000), []byte{}), nil
}

func (a *AggregatorMock) SubscribeToNewRound(logChan chan<- types.Log) (client.Subscription, error) {
	rlog := types.Log{
		Address: a.address,
		Topics: []common.Hash{
			common.HexToHash("0xc3c45d1924f55369653f407ee9f095309d1e687b2c0011b1f709042d4f457e17"),
			common.BigToHash(big.NewInt(101)),
			common.BigToHash(big.NewInt(time.Now().Unix())),
		},
		Data:        nil,
		BlockNumber: 9083920,
		TxHash:      common.HexToHash("0x525013b79eb6d1948160ca3a0640c84e2ea080f852e1455b401ee8234ed6fb8d"),
		BlockHash:   common.HexToHash("0xfd53f5db740996d4b0f036a7ae435b6e38faf6bbe636b7e249c0586c0c1e443f"),
	}
	go func() { logChan <- rlog }()
	return &SubscriptionMock{}, nil
}

func (a *AggregatorMock) SubscribeToOracleAnswer(
	roundId *big.Int,
	oracle common.Address,
	logChan chan<- types.Log,
) (client.Subscription, error) {
	alog := types.Log{
		Address: a.address,
		Topics: []common.Hash{
			common.HexToHash("0xb51168059c83c860caf5b830c5d2e64c2172c6fb2fe9f25447d9838e18d93b60"),
			common.BigToHash(big.NewInt(50)),
			common.BigToHash(roundId),
			oracle.Hash(),
		},
		Data:        nil,
		BlockNumber: 9083921,
		TxHash:      common.HexToHash("0xb51168059c83c860caf5b830c5d2e64c2172c6fb2fe9f25447d9838e18d93b60"),
		BlockHash:   common.HexToHash("0x6f8135b9e3b688ebea07d37d940f3c7b43a9fd6144c6e47ef7703f87c3b4e450"),
	}
	go func() { logChan <- alog }()
	return &SubscriptionMock{}, nil
}

func (a *AggregatorMock) UnmarshalNewRoundEvent(log types.Log) (*client.NewRoundEvent, error) {
	nr := &client.NewRoundEvent{}
	if len(log.Topics) == 3 {
		nr.RoundID = log.Topics[1].Big()
		nr.StartedBy = common.BytesToAddress(log.Topics[2].Bytes())
	} else {
		return nr, errors.New("invalid log type")
	}
	return nr, nil
}

func (a *AggregatorMock) SetAddress(address common.Address) {
	a.address = address
}

func (a *AggregatorMock) SetOracles(oracles client.OracleMapping) {
	a.oracles = oracles
}

type AggregatorMockTimeout struct {
	AggregatorMock
}

func NewAggregatorMockTimeout() client.Aggregator {
	return &AggregatorMockTimeout{
		AggregatorMock{
			address: common.BigToAddress(big.NewInt(rand.Int63())),
			oracles: GenerateOracles(),
		},
	}
}

func (a *AggregatorMockTimeout) SubscribeToOracleAnswer(
	roundId *big.Int,
	oracle common.Address,
	logChan chan<- types.Log,
) (client.Subscription, error) {
	return &SubscriptionMock{}, nil
}

type AggregatorMockError struct {
	AggregatorMock
}

func NewAggregatorMockError() client.Aggregator {
	return &AggregatorMockError{
		AggregatorMock{
			address: common.BigToAddress(big.NewInt(rand.Int63())),
			oracles: GenerateOracles(),
		},
	}
}

func (a *AggregatorMockError) SubscribeToOracleAnswer(
	roundId *big.Int,
	oracle common.Address,
	logChan chan<- types.Log,
) (client.Subscription, error) {
	return &SubscriptionMockError{}, nil
}

type AggregatorMockInvalidLog struct {
	AggregatorMock
}

func NewAggregatorMockInvalidLog() client.Aggregator {
	return &AggregatorMockInvalidLog{
		AggregatorMock{
			address: common.BigToAddress(big.NewInt(rand.Int63())),
			oracles: GenerateOracles(),
		},
	}
}

func (a *AggregatorMockInvalidLog) SubscribeToNewRound(logChan chan<- types.Log) (client.Subscription, error) {
	go func() { logChan <- types.Log{} }()
	return &SubscriptionMock{}, nil
}

func GenerateOracles() client.OracleMapping {
	oracles := client.OracleMapping{}
	for i := 0; i < 10; i++ {
		oracles[common.BigToAddress(big.NewInt(rand.Int63()))] = "Unknown"
	}
	return oracles
}
