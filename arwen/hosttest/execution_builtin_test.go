package hosttest

import (
	"encoding/hex"
	"errors"
	"math/big"
	"testing"

	"github.com/ElrondNetwork/arwen-wasm-vm/v1_3/arwen"
	contextmock "github.com/ElrondNetwork/arwen-wasm-vm/v1_3/mock/context"
	worldmock "github.com/ElrondNetwork/arwen-wasm-vm/v1_3/mock/world"
	test "github.com/ElrondNetwork/arwen-wasm-vm/v1_3/testcommon"
	"github.com/ElrondNetwork/elrond-vm-common"
	"github.com/ElrondNetwork/elrond-vm-common/data/esdt"
	"github.com/stretchr/testify/require"
)

//TODO package contains snake case named files, rename those.

func TestExecution_ExecuteOnDestContext_ESDTTransferWithoutExecute(t *testing.T) {
	code := test.GetTestSCCodeModule("exec-dest-ctx-esdt/basic", "basic", "../../")
	scBalance := big.NewInt(1000)
	host, world := test.DefaultTestArwenForCallWithWorldMock(t, code, scBalance)

	tokenKey := worldmock.MakeTokenKey(test.ESDTTestTokenName, 0)
	err := world.BuiltinFuncs.SetTokenData(test.ParentAddress, tokenKey, &esdt.ESDigitalToken{
		Value: big.NewInt(100),
		Type:  uint32(vmcommon.Fungible),
	})
	require.Nil(t, err)

	input := test.DefaultTestContractCallInput()
	input.Function = "basic_transfer"
	input.GasProvided = 100000
	input.ESDTTokenName = test.ESDTTestTokenName
	input.ESDTValue = big.NewInt(16)

	vmOutput, err := host.RunSmartContractCall(input)

	verify := test.NewVMOutputVerifier(t, vmOutput, err)
	verify.
		Ok()
}

func TestExecution_ExecuteOnDestContext_MockBuiltinFunctions_Claim(t *testing.T) {
	parentGasUsed := uint64(1988)
	test.BuildInstanceCallTest(t).
		WithContracts(
			test.CreateInstanceContract(test.ParentAddress).
				WithCode(test.GetTestSCCode("exec-dest-ctx-builtin", "../../")).
				WithBalance(1000)).
		WithInput(test.CreateTestContractCallInputBuilder().
			WithRecipientAddr(test.ParentAddress).
			WithGasProvided(test.GasProvided).
			WithFunction("callBuiltinClaim").
			Build()).
		WithSetup(func(host arwen.VMHost, stubBlockchainHook *contextmock.BlockchainHookStub) {
			stubBlockchainHook.ProcessBuiltInFunctionCalled = dummyProcessBuiltInFunction
			host.SetProtocolBuiltinFunctions(getDummyBuiltinFunctionNames())
		}).
		AndAssertResults(func(host arwen.VMHost, stubBlockchainHook *contextmock.BlockchainHookStub, verify *test.VMOutputVerifier) {
			verify.
				Ok().
				BalanceDelta(test.ParentAddress, 42).
				GasUsed(test.ParentAddress, parentGasUsed).
				GasRemaining(test.GasProvided - parentGasUsed).
				ReturnData([]byte("succ"))
		})
}

func TestExecution_ExecuteOnDestContext_MockBuiltinFunctions_DoSomething(t *testing.T) {
	parentGasUsed := uint64(1992)
	test.BuildInstanceCallTest(t).
		WithContracts(
			test.CreateInstanceContract(test.ParentAddress).
				WithCode(test.GetTestSCCode("exec-dest-ctx-builtin", "../../")).
				WithBalance(1000)).
		WithInput(test.CreateTestContractCallInputBuilder().
			WithRecipientAddr(test.ParentAddress).
			WithGasProvided(test.GasProvided).
			WithFunction("callBuiltinDoSomething").
			Build()).
		WithSetup(func(host arwen.VMHost, stubBlockchainHook *contextmock.BlockchainHookStub) {
			stubBlockchainHook.ProcessBuiltInFunctionCalled = dummyProcessBuiltInFunction
			host.SetProtocolBuiltinFunctions(getDummyBuiltinFunctionNames())
		}).
		AndAssertResults(func(host arwen.VMHost, stubBlockchainHook *contextmock.BlockchainHookStub, verify *test.VMOutputVerifier) {
			verify.
				Ok().
				BalanceDelta(test.ParentAddress, big.NewInt(0).Sub(arwen.One, arwen.One).Int64()).
				GasUsed(test.ParentAddress, parentGasUsed).
				GasRemaining(test.GasProvided - parentGasUsed).
				ReturnData([]byte("succ"))
		})
}

func TestExecution_ExecuteOnDestContext_MockBuiltinFunctions_Nonexistent(t *testing.T) {
	test.BuildInstanceCallTest(t).
		WithContracts(
			test.CreateInstanceContract(test.ParentAddress).
				WithCode(test.GetTestSCCode("exec-dest-ctx-builtin", "../../")).
				WithBalance(1000)).
		WithInput(test.CreateTestContractCallInputBuilder().
			WithRecipientAddr(test.ParentAddress).
			WithGasProvided(test.GasProvided).
			WithFunction("callNonexistingBuiltin").
			Build()).
		WithSetup(func(host arwen.VMHost, stubBlockchainHook *contextmock.BlockchainHookStub) {
			stubBlockchainHook.ProcessBuiltInFunctionCalled = dummyProcessBuiltInFunction
			host.SetProtocolBuiltinFunctions(getDummyBuiltinFunctionNames())
		}).
		AndAssertResults(func(host arwen.VMHost, stubBlockchainHook *contextmock.BlockchainHookStub, verify *test.VMOutputVerifier) {
			verify.
				ReturnCode(vmcommon.ExecutionFailed).
				ReturnMessage(arwen.ErrFuncNotFound.Error()).
				GasRemaining(0)
		})
}

func TestExecution_ExecuteOnDestContext_MockBuiltinFunctions_Fail(t *testing.T) {
	test.BuildInstanceCallTest(t).
		WithContracts(
			test.CreateInstanceContract(test.ParentAddress).
				WithCode(test.GetTestSCCode("exec-dest-ctx-builtin", "../../")).
				WithBalance(1000)).
		WithInput(test.CreateTestContractCallInputBuilder().
			WithRecipientAddr(test.ParentAddress).
			WithGasProvided(test.GasProvided).
			WithFunction("callBuiltinFail").
			Build()).
		WithSetup(func(host arwen.VMHost, stubBlockchainHook *contextmock.BlockchainHookStub) {
			stubBlockchainHook.ProcessBuiltInFunctionCalled = dummyProcessBuiltInFunction
			host.SetProtocolBuiltinFunctions(getDummyBuiltinFunctionNames())
		}).
		AndAssertResults(func(host arwen.VMHost, stubBlockchainHook *contextmock.BlockchainHookStub, verify *test.VMOutputVerifier) {
			verify.
				ReturnCode(vmcommon.ExecutionFailed).
				ReturnMessage("whatdidyoudo").
				GasRemaining(0)
		})
}

func TestExecution_AsyncCall_MockBuiltinFails(t *testing.T) {
	test.BuildInstanceCallTest(t).
		WithContracts(
			test.CreateInstanceContract(test.ParentAddress).
				WithCode(test.GetTestSCCode("async-call-builtin", "../../")).
				WithBalance(1000)).
		WithInput(test.CreateTestContractCallInputBuilder().
			WithRecipientAddr(test.ParentAddress).
			WithGasProvided(test.GasProvided).
			WithFunction("performAsyncCallToBuiltin").
			WithArguments([]byte{1}).
			Build()).
		WithSetup(func(host arwen.VMHost, stubBlockchainHook *contextmock.BlockchainHookStub) {
			stubBlockchainHook.ProcessBuiltInFunctionCalled = dummyProcessBuiltInFunction
			host.SetProtocolBuiltinFunctions(getDummyBuiltinFunctionNames())
		}).
		AndAssertResults(func(host arwen.VMHost, stubBlockchainHook *contextmock.BlockchainHookStub, verify *test.VMOutputVerifier) {
			verify.
				Ok().
				ReturnData([]byte("hello"), []byte{10})
		})
}

func TestESDT_GettersAPI(t *testing.T) {
	test.BuildInstanceCallTest(t).
		WithContracts(
			test.CreateInstanceContract(test.ParentAddress).
				WithCode(test.GetTestSCCode("exchange", "../../")).
				WithBalance(1000)).
		WithInput(test.CreateTestContractCallInputBuilder().
			WithRecipientAddr(test.ParentAddress).
			WithGasProvided(test.GasProvided).
			WithFunction("validateGetters").
			WithESDTValue(big.NewInt(5)).
			WithESDTTokenName(test.ESDTTestTokenName).
			Build()).
		WithSetup(func(host arwen.VMHost, stubBlockchainHook *contextmock.BlockchainHookStub) {
			stubBlockchainHook.ProcessBuiltInFunctionCalled = dummyProcessBuiltInFunction
			host.SetProtocolBuiltinFunctions(getDummyBuiltinFunctionNames())
		}).
		AndAssertResults(func(host arwen.VMHost, stubBlockchainHook *contextmock.BlockchainHookStub, verify *test.VMOutputVerifier) {
			verify.
				Ok()
		})
}

func TestESDT_GettersAPI_ExecuteAfterBuiltinCall(t *testing.T) {
	host, world := test.DefaultTestArwenWithWorldMock(t)

	initialESDTTokenBalance := uint64(1000)

	// Deploy the "parent" contract, which will call the exchange; the actual
	// code of the contract is not important, because the exchange will be called
	// by the "parent" using a manual call to host.ExecuteOnDestContext().
	dummyCode := test.GetTestSCCode("init-simple", "../../")
	parentAccount := world.AcctMap.CreateSmartContractAccount(test.UserAddress, test.ParentAddress, dummyCode, world)
	_ = parentAccount.SetTokenBalanceUint64(test.ESDTTestTokenKey, initialESDTTokenBalance)

	// Deploy the exchange contract, which will receive ESDT and verify that it
	// can see the received token amount and token name.
	exchangeAddress := test.MakeTestSCAddress("exchange")
	exchangeCode := test.GetTestSCCode("exchange", "../../")
	exchange := world.AcctMap.CreateSmartContractAccount(test.UserAddress, exchangeAddress, exchangeCode, world)
	exchange.Balance = big.NewInt(1000)

	// Prepare Arwen to appear as if the parent contract is being executed
	input := test.DefaultTestContractCallInput()
	host.Runtime().InitStateFromContractCallInput(input)
	_ = host.Runtime().StartWasmerInstance(dummyCode, input.GasProvided, true)
	err := host.Output().TransferValueOnly(input.RecipientAddr, input.CallerAddr, input.CallValue, false)
	require.Nil(t, err)

	// Transfer ESDT to the exchange and call its "validateGetters" method
	esdtValue := int64(5)
	input.CallerAddr = test.ParentAddress
	input.RecipientAddr = exchangeAddress
	input.Function = vmcommon.BuiltInFunctionESDTTransfer
	input.GasProvided = 10000
	input.Arguments = [][]byte{
		test.ESDTTestTokenName,
		big.NewInt(esdtValue).Bytes(),
		[]byte("validateGetters"),
	}

	vmOutput, err := host.ExecuteOnDestContext(input)

	verify := test.NewVMOutputVerifier(t, vmOutput, err)
	verify.
		Ok()

	parentESDTBalance, _ := parentAccount.GetTokenBalanceUint64(test.ESDTTestTokenKey)
	require.Equal(t, initialESDTTokenBalance-uint64(esdtValue), parentESDTBalance)
}

func dummyProcessBuiltInFunction(input *vmcommon.ContractCallInput) (*vmcommon.VMOutput, error) {
	outputAccounts := make(map[string]*vmcommon.OutputAccount)
	outputAccounts[string(test.ParentAddress)] = &vmcommon.OutputAccount{
		BalanceDelta: big.NewInt(0),
		Address:      test.ParentAddress}

	if input.Function == "builtinClaim" {
		outputAccounts[string(test.ParentAddress)].BalanceDelta = big.NewInt(42)
		return &vmcommon.VMOutput{
			GasRemaining:   400 + input.GasLocked,
			OutputAccounts: outputAccounts,
		}, nil
	}
	if input.Function == "builtinDoSomething" {
		return &vmcommon.VMOutput{
			GasRemaining:   400 + input.GasLocked,
			OutputAccounts: outputAccounts,
		}, nil
	}
	if input.Function == "builtinFail" {
		return nil, errors.New("whatdidyoudo")
	}
	if input.Function == vmcommon.BuiltInFunctionESDTTransfer {
		vmOutput := &vmcommon.VMOutput{
			GasRemaining: 0,
		}
		function := string(input.Arguments[2])
		esdtTransferTxData := function
		for _, arg := range input.Arguments[3:] {
			esdtTransferTxData += "@" + hex.EncodeToString(arg)
		}
		outTransfer := vmcommon.OutputTransfer{
			Value:         big.NewInt(0),
			GasLimit:      input.GasProvided - test.ESDTTransferGasCost + input.GasLocked,
			Data:          []byte(esdtTransferTxData),
			CallType:      vmcommon.AsynchronousCall,
			SenderAddress: input.CallerAddr,
		}
		vmOutput.OutputAccounts = make(map[string]*vmcommon.OutputAccount)
		vmOutput.OutputAccounts[string(input.RecipientAddr)] = &vmcommon.OutputAccount{
			Address:         input.RecipientAddr,
			OutputTransfers: []vmcommon.OutputTransfer{outTransfer},
		}
		// TODO when ESDT token balance querying is implemented, ensure the
		// transfers that happen here are persisted in the mock accounts
		return vmOutput, nil
	}

	return nil, arwen.ErrFuncNotFound
}

func getDummyBuiltinFunctionNames() vmcommon.FunctionNames {
	names := make(vmcommon.FunctionNames)

	var empty struct{}
	names["builtinClaim"] = empty
	names["builtinDoSomething"] = empty
	names["builtinFail"] = empty
	names[vmcommon.BuiltInFunctionESDTTransfer] = empty

	return names
}
