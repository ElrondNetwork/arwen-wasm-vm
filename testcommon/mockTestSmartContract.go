package testcommon

import (
	"testing"

	"github.com/ElrondNetwork/arwen-wasm-vm/v1_3/arwen"
	mock "github.com/ElrondNetwork/arwen-wasm-vm/v1_3/mock/context"
)

// TestConfig is configuration for async call tests
type TestConfig struct {
	GasProvided        uint64
	GasProvidedToChild uint64
	GasUsedByParent    uint64
	GasUsedByChild     uint64
	GasUsedByCallback  uint64
	GasLockCost        uint64

	ParentBalance int64
	ChildBalance  int64

	TransferFromParentToChild int64
	TransferToThirdParty      int64
	TransferToVault           int64
	TransferFromChildToParent int64

	ESDTTokensToTransfer         uint64
	CallbackESDTTokensToTransfer uint64

	ChildCalls          int
	RecursiveChildCalls int
}

type testSmartContract struct {
	address []byte
	balance int64
	config  *TestConfig
	shardID uint32
}

// MockTestSmartContract represents the config data for the mock smart contract instance to be tested
type MockTestSmartContract struct {
	testSmartContract
	initMethods []func(*mock.InstanceMock, *TestConfig)
}

// CreateMockContract build a contract to be used in a test creted with BuildMockInstanceCallTest
func CreateMockContract(address []byte) *MockTestSmartContract {
	return CreateMockContractOnShard(address, 0)
}

// CreateMockContractOnShard build a contract to be used in a test creted with BuildMockInstanceCallTest
func CreateMockContractOnShard(address []byte, shardID uint32) *MockTestSmartContract {
	return &MockTestSmartContract{
		testSmartContract: testSmartContract{
			address: address,
			shardID: shardID,
		},
	}
}

// WithBalance provides the balance for the MockTestSmartContract
func (mockSC *MockTestSmartContract) WithBalance(balance int64) *MockTestSmartContract {
	mockSC.balance = balance
	return mockSC
}

// WithShardID provides the shardID for the MockTestSmartContract
func (mockSC *MockTestSmartContract) WithShardID(shardID uint32) *MockTestSmartContract {
	mockSC.shardID = shardID
	return mockSC
}

// WithConfig provides the config object for the MockTestSmartContract
func (mockSC *MockTestSmartContract) WithConfig(config *TestConfig) *MockTestSmartContract {
	mockSC.config = config
	return mockSC
}

// WithMethods provides the methods for the MockTestSmartContract
func (mockSC *MockTestSmartContract) WithMethods(initMethods ...func(*mock.InstanceMock, *TestConfig)) MockTestSmartContract {
	mockSC.initMethods = initMethods
	return *mockSC
}

func (mockSC *MockTestSmartContract) initialize(t testing.TB, host arwen.VMHost, imb *mock.InstanceBuilderMock) {
	instance := imb.CreateAndStoreInstanceMock(t, host, mockSC.address, mockSC.shardID, mockSC.balance)
	for _, initMethod := range mockSC.initMethods {
		initMethod(instance, mockSC.config)
	}
}
