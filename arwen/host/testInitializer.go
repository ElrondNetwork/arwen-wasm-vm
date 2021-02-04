package host

import (
	"bytes"
	"errors"
	"math/big"
	"testing"

	"github.com/ElrondNetwork/arwen-wasm-vm/arwen"
	"github.com/ElrondNetwork/arwen-wasm-vm/config"
	contextmock "github.com/ElrondNetwork/arwen-wasm-vm/mock/context"
	"github.com/ElrondNetwork/elrond-go/core/vmcommon"
	"github.com/stretchr/testify/require"
)

var defaultVMType = []byte{0xF, 0xF}
var errAccountNotFound = errors.New("account not found")

var userAddress = []byte("userAccount.....................")
var parentAddress = []byte("\x00\x00\x00\x00\x00\x00\x00\x00\x0f\x0fparentSC..............")
var childAddress = []byte("\x00\x00\x00\x00\x00\x00\x00\x00\x0f\x0fchildSC...............")

var customGasSchedule = config.GasScheduleMap(nil)

// DefaultTestArwenForDeployment creates an Arwen vmHost configured for testing deployments
func defaultTestArwenForDeployment(t *testing.T, _ uint64, newAddress []byte) *vmHost {
	stubBlockchainHook := &contextmock.BlockchainHookStub{}
	stubBlockchainHook.GetUserAccountCalled = func(address []byte) (vmcommon.UserAccountHandler, error) {
		return &contextmock.StubAccount{
			Nonce: 24,
		}, nil
	}
	stubBlockchainHook.NewAddressCalled = func(creatorAddress []byte, nonce uint64, vmType []byte) ([]byte, error) {
		return newAddress, nil
	}

	host := defaultTestArwen(t, stubBlockchainHook)
	return host
}

func defaultTestArwenForCall(tb testing.TB, code []byte, balance *big.Int) (*vmHost, *contextmock.BlockchainHookStub) {
	stubBlockchainHook := &contextmock.BlockchainHookStub{}
	stubBlockchainHook.GetUserAccountCalled = func(scAddress []byte) (vmcommon.UserAccountHandler, error) {
		if bytes.Equal(scAddress, parentAddress) {
			return &contextmock.StubAccount{
				Balance: balance,
			}, nil
		}
		return nil, errAccountNotFound
	}
	stubBlockchainHook.GetCodeCalled = func(account vmcommon.UserAccountHandler) []byte {
		return code
	}

	host := defaultTestArwen(tb, stubBlockchainHook)
	return host, stubBlockchainHook
}

func defaultTestArwenForCallWithInstanceMocks(tb testing.TB) (*vmHost, *contextmock.BlockchainHookStub, *contextmock.InstanceBuilderMock) {
	stubBlockchainHook := &contextmock.BlockchainHookStub{}
	stubBlockchainHook.GetUserAccountCalled = func(scAddress []byte) (vmcommon.UserAccountHandler, error) {
		return &contextmock.StubAccount{
			Address: scAddress,
			Balance: big.NewInt(1000),
		}, nil
	}
	stubBlockchainHook.GetCodeCalled = func(account vmcommon.UserAccountHandler) []byte {
		return account.AddressBytes()
	}

	host := defaultTestArwen(tb, stubBlockchainHook)

	code := arwen.GetTestSCCode("counter", "../../")
	instanceBuilderMock := contextmock.NewInstanceBuilderMock(tb, code)
	host.Runtime().ReplaceInstanceBuilder(instanceBuilderMock)

	return host, stubBlockchainHook, instanceBuilderMock
}

// defaultTestArwenForTwoSCs creates an Arwen vmHost configured for testing calls between 2 SmartContracts
func defaultTestArwenForTwoSCs(
	t *testing.T,
	parentCode []byte,
	childCode []byte,
	parentSCBalance *big.Int,
	childSCBalance *big.Int,
) (*vmHost, *contextmock.BlockchainHookStub) {
	stubBlockchainHook := &contextmock.BlockchainHookStub{}

	if parentSCBalance == nil {
		parentSCBalance = big.NewInt(1000)
	}

	if childSCBalance == nil {
		childSCBalance = big.NewInt(1000)
	}

	stubBlockchainHook.GetUserAccountCalled = func(scAddress []byte) (vmcommon.UserAccountHandler, error) {
		if bytes.Equal(scAddress, parentAddress) {
			return &contextmock.StubAccount{
				Address: parentAddress,
				Balance: parentSCBalance,
			}, nil
		}
		if bytes.Equal(scAddress, childAddress) {
			return &contextmock.StubAccount{
				Address: childAddress,
				Balance: childSCBalance,
			}, nil
		}

		return nil, errAccountNotFound
	}
	stubBlockchainHook.GetCodeCalled = func(account vmcommon.UserAccountHandler) []byte {
		if bytes.Equal(account.AddressBytes(), parentAddress) {
			return parentCode
		}
		if bytes.Equal(account.AddressBytes(), childAddress) {
			return childCode
		}
		return nil
	}

	host := defaultTestArwen(t, stubBlockchainHook)
	return host, stubBlockchainHook
}

func defaultTestArwen(tb testing.TB, blockchain vmcommon.BlockchainHook) *vmHost {
	gasSchedule := customGasSchedule
	if gasSchedule == nil {
		gasSchedule = config.MakeGasMapForTests()
	}

	host, err := NewArwenVM(blockchain, &arwen.VMHostParameters{
		VMType:                   defaultVMType,
		BlockGasLimit:            uint64(1000),
		GasSchedule:              gasSchedule,
		ProtocolBuiltinFunctions: make(vmcommon.FunctionNames),
		ElrondProtectedKeyPrefix: []byte("ELROND"),
		UseWarmInstance:          false,
		DynGasLockEnableEpoch:    0,
	})
	require.Nil(tb, err)
	require.NotNil(tb, host)
	return host
}

// DefaultTestContractCreateInput creates a vmcommon.ContractCreateInput struct with default values
func DefaultTestContractCreateInput() *vmcommon.ContractCreateInput {
	return &vmcommon.ContractCreateInput{
		VMInput: vmcommon.VMInput{
			CallerAddr: []byte("caller"),
			Arguments: [][]byte{
				[]byte("argument 1"),
				[]byte("argument 2"),
			},
			CallValue:   big.NewInt(0),
			CallType:    vmcommon.DirectCall,
			GasPrice:    0,
			GasProvided: 0,
		},
		ContractCode: []byte("contract"),
	}
}

// DefaultTestContractCallInput creates a vmcommon.ContractCallInput struct with default values
func DefaultTestContractCallInput() *vmcommon.ContractCallInput {
	return &vmcommon.ContractCallInput{
		VMInput: vmcommon.VMInput{
			CallerAddr:  userAddress,
			Arguments:   make([][]byte, 0),
			CallValue:   big.NewInt(0),
			CallType:    vmcommon.DirectCall,
			GasPrice:    0,
			GasProvided: 0,
		},
		RecipientAddr: parentAddress,
		Function:      "function",
	}
}

func LoadGasScheduleConfig(filepath string) (config.GasScheduleMap, error) {
	gasScheduleConfig, err := arwen.LoadTomlFileToMap(filepath)
	if err != nil {
		return nil, err
	}

	flattenedGasSchedule := make(config.GasScheduleMap)
	for libType, costs := range gasScheduleConfig {
		flattenedGasSchedule[libType] = make(map[string]uint64)
		costsMap := costs.(map[string]interface{})
		for operationName, cost := range costsMap {
			flattenedGasSchedule[libType][operationName] = uint64(cost.(int64))
		}
	}

	return flattenedGasSchedule, nil
}
