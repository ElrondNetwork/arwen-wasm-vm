package host

import (
	"bytes"
	"errors"
	"fmt"
	"math/big"
	"testing"

	"github.com/ElrondNetwork/arwen-wasm-vm/arwen"
	"github.com/ElrondNetwork/arwen-wasm-vm/config"
	contextmock "github.com/ElrondNetwork/arwen-wasm-vm/mock/context"
	"github.com/ElrondNetwork/arwen-wasm-vm/wasmer"
	"github.com/ElrondNetwork/elrond-go/core/vmcommon"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var counterKey = []byte("COUNTER")
var WASMLocalsLimit = uint64(4000)

const (
	increment               = "increment"
	callRecursive           = "callRecursive"
	parentCallsChild        = "parentCallsChild"
	parentPerformAsyncCall  = "parentPerformAsyncCall"
	parentFunctionChildCall = "parentFunctionChildCall"
)

func TestSCMem(t *testing.T) {
	code := arwen.GetTestSCCode("misc", "../../")
	host, _ := defaultTestArwenForCall(t, code, nil)
	input := DefaultTestContractCallInput()
	input.GasProvided = 100000
	input.Function = "iterate_over_byte_array"
	vmOutput, err := host.RunSmartContractCall(input)
	require.Nil(t, err)
	require.Equal(t, vmcommon.Ok, vmOutput.ReturnCode)

	testString := "this is some random string of bytes"
	expectedData := [][]byte{
		[]byte(testString),
		{35},
	}
	for _, c := range testString {
		expectedData = append(expectedData, []byte{byte(c)})
	}
	require.Equal(t, expectedData, vmOutput.ReturnData)
}

func TestExecution_DeployNewAddressErr(t *testing.T) {
	stubBlockchainHook := &contextmock.BlockchainHookStub{}

	errNewAddress := errors.New("new address error")

	host := defaultTestArwen(t, stubBlockchainHook)
	input := DefaultTestContractCreateInput()
	stubBlockchainHook.GetUserAccountCalled = func(address []byte) (vmcommon.UserAccountHandler, error) {
		require.Equal(t, input.CallerAddr, address)
		return &contextmock.StubAccount{}, nil
	}
	stubBlockchainHook.NewAddressCalled = func(creatorAddress []byte, nonce uint64, vmType []byte) ([]byte, error) {
		require.Equal(t, input.CallerAddr, creatorAddress)
		require.Equal(t, uint64(0), nonce)
		require.Equal(t, defaultVMType, vmType)
		return nil, errNewAddress
	}

	vmOutput, err := host.RunSmartContractCreate(input)
	require.Nil(t, err)
	require.NotNil(t, vmOutput)
	require.Equal(t, vmcommon.ExecutionFailed, vmOutput.ReturnCode)
	require.Equal(t, errNewAddress.Error(), vmOutput.ReturnMessage)
}

func TestExecution_DeployOutOfGas(t *testing.T) {
	newAddress := []byte("new smartcontract")
	host := defaultTestArwenForDeployment(t, 24, newAddress)
	input := DefaultTestContractCreateInput()
	input.GasProvided = 8 // default deployment requires 9 units of Gas
	vmOutput, err := host.RunSmartContractCreate(input)
	require.Nil(t, err)
	require.NotNil(t, vmOutput)
	require.Equal(t, vmcommon.OutOfGas, vmOutput.ReturnCode)
	require.Equal(t, arwen.ErrNotEnoughGas.Error(), vmOutput.ReturnMessage)
}

func TestExecution_DeployNotWASM(t *testing.T) {
	newAddress := []byte("new smartcontract")
	host := defaultTestArwenForDeployment(t, 24, newAddress)
	input := DefaultTestContractCreateInput()
	input.GasProvided = 9
	input.ContractCode = []byte("not WASM")
	vmOutput, err := host.RunSmartContractCreate(input)
	require.Nil(t, err)
	require.NotNil(t, vmOutput)
	require.Equal(t, vmcommon.ContractInvalid, vmOutput.ReturnCode)
}

func TestExecution_DeployWASM_WithoutMemory(t *testing.T) {
	newAddress := []byte("new smartcontract")
	host := defaultTestArwenForDeployment(t, 24, newAddress)
	input := DefaultTestContractCreateInput()
	input.GasProvided = 1000
	input.ContractCode = arwen.GetTestSCCode("memoryless", "../../")
	vmOutput, err := host.RunSmartContractCreate(input)
	require.Nil(t, err)
	require.NotNil(t, vmOutput)
	require.Equal(t, vmcommon.ContractInvalid, vmOutput.ReturnCode)
}

func TestExecution_DeployWASM_WrongInit(t *testing.T) {
	newAddress := []byte("new smartcontract")
	host := defaultTestArwenForDeployment(t, 24, newAddress)
	input := DefaultTestContractCreateInput()
	input.GasProvided = 1000
	input.ContractCode = arwen.GetTestSCCode("init-wrong", "../../")
	vmOutput, err := host.RunSmartContractCreate(input)
	require.Nil(t, err)
	require.NotNil(t, vmOutput)
	require.Equal(t, vmcommon.ContractInvalid, vmOutput.ReturnCode)
}

func TestExecution_DeployWASM_WrongMethods(t *testing.T) {
	newAddress := []byte("new smartcontract")
	host := defaultTestArwenForDeployment(t, 24, newAddress)
	input := DefaultTestContractCreateInput()
	input.GasProvided = 1000
	input.ContractCode = arwen.GetTestSCCode("signatures", "../../")
	vmOutput, err := host.RunSmartContractCreate(input)
	require.Nil(t, err)
	require.NotNil(t, vmOutput)
	require.Equal(t, vmcommon.ContractInvalid, vmOutput.ReturnCode)
}

func TestExecution_DeployWASM_Successful(t *testing.T) {
	newAddress := []byte("new smartcontract")
	host := defaultTestArwenForDeployment(t, 24, newAddress)
	input := DefaultTestContractCreateInput()
	input.CallValue = big.NewInt(88)
	input.GasProvided = 1000
	input.ContractCode = arwen.GetTestSCCode("init-correct", "../../")
	input.Arguments = [][]byte{{0}}

	vmOutput, err := host.RunSmartContractCreate(input)
	require.Nil(t, err)
	require.NotNil(t, vmOutput)
	require.Equal(t, vmcommon.Ok, vmOutput.ReturnCode)
	require.Len(t, vmOutput.ReturnData, 1)
	require.Equal(t, []byte("init successful"), vmOutput.ReturnData[0])
	require.Equal(t, uint64(528), vmOutput.GasRemaining)
	require.Len(t, vmOutput.OutputAccounts, 2)
	require.Equal(t, uint64(24), vmOutput.OutputAccounts["caller"].Nonce)
	require.Equal(t, input.ContractCode, vmOutput.OutputAccounts[string(newAddress)].Code)
	require.Equal(t, big.NewInt(88), vmOutput.OutputAccounts[string(newAddress)].BalanceDelta)
}

func TestExecution_DeployWASM_Popcnt(t *testing.T) {
	newAddress := []byte("new smartcontract")
	host := defaultTestArwenForDeployment(t, 24, newAddress)
	input := DefaultTestContractCreateInput()
	input.CallValue = big.NewInt(88)
	input.GasProvided = 1000
	input.ContractCode = arwen.GetTestSCCode("init-simple-popcnt", "../../")
	input.Arguments = [][]byte{}

	vmOutput, err := host.RunSmartContractCreate(input)
	require.Nil(t, err)
	require.NotNil(t, vmOutput)
	require.Equal(t, vmcommon.Ok, vmOutput.ReturnCode)
	require.Len(t, vmOutput.ReturnData, 1)
	require.Equal(t, []byte{3}, vmOutput.ReturnData[0])
}

func TestExecution_DeployWASM_AtMaximumLocals(t *testing.T) {
	newAddress := []byte("new smartcontract")
	host := defaultTestArwenForDeployment(t, 24, newAddress)
	input := DefaultTestContractCreateInput()
	input.CallValue = big.NewInt(88)
	input.GasProvided = 1000
	input.ContractCode = makeBytecodeWithLocals(WASMLocalsLimit)
	input.Arguments = [][]byte{{0}}

	vmOutput, err := host.RunSmartContractCreate(input)
	require.Nil(t, err)
	require.NotNil(t, vmOutput)
	require.Equal(t, vmcommon.Ok, vmOutput.ReturnCode)
}

func TestExecution_DeployWASM_MoreThanMaximumLocals(t *testing.T) {
	newAddress := []byte("new smartcontract")
	host := defaultTestArwenForDeployment(t, 24, newAddress)
	input := DefaultTestContractCreateInput()
	input.CallValue = big.NewInt(88)
	input.GasProvided = 1000
	input.ContractCode = makeBytecodeWithLocals(WASMLocalsLimit + 1)
	input.Arguments = [][]byte{{0}}

	vmOutput, err := host.RunSmartContractCreate(input)
	require.Nil(t, err)
	require.NotNil(t, vmOutput)
	require.Equal(t, vmcommon.ContractInvalid, vmOutput.ReturnCode)
}

func TestExecution_DeployWASM_Init_Errors(t *testing.T) {
	newAddress := []byte("new smartcontract")
	host := defaultTestArwenForDeployment(t, 24, newAddress)
	input := DefaultTestContractCreateInput()
	input.CallValue = big.NewInt(88)
	input.GasProvided = 1000
	input.ContractCode = arwen.GetTestSCCode("init-correct", "../../")

	// init() calls signalError()
	input.Arguments = [][]byte{{1}}
	vmOutput, err := host.RunSmartContractCreate(input)
	require.Nil(t, err)
	require.NotNil(t, vmOutput)
	require.Equal(t, vmcommon.UserError, vmOutput.ReturnCode)

	// init() starts an infinite loop
	input.Arguments = [][]byte{{2}}
	vmOutput, err = host.RunSmartContractCreate(input)
	require.Nil(t, err)
	require.NotNil(t, vmOutput)
	require.Equal(t, vmcommon.OutOfGas, vmOutput.ReturnCode)
}

func TestExecution_ManyDeployments(t *testing.T) {
	if testing.Short() {
		t.Skip("not a short test")
	}

	ownerNonce := uint64(23)
	newAddress := "new smartcontract"
	stubBlockchainHook := &contextmock.BlockchainHookStub{}
	stubBlockchainHook.GetUserAccountCalled = func(address []byte) (vmcommon.UserAccountHandler, error) {
		return &contextmock.StubAccount{Nonce: ownerNonce}, nil
	}
	stubBlockchainHook.NewAddressCalled = func(creatorAddress []byte, nonce uint64, vmType []byte) ([]byte, error) {
		ownerNonce++
		return []byte(newAddress + " " + fmt.Sprint(ownerNonce)), nil
	}

	host := defaultTestArwen(t, stubBlockchainHook)
	input := DefaultTestContractCreateInput()
	input.CallerAddr = []byte("owner")
	input.Arguments = make([][]byte, 0)
	input.CallValue = big.NewInt(88)
	input.ContractCode = arwen.GetTestSCCode("init-simple", "../../")

	numDeployments := 1000
	for i := 0; i < numDeployments; i++ {
		input.GasProvided = 100000
		vmOutput, err := host.RunSmartContractCreate(input)
		require.Nil(t, err)
		require.NotNil(t, vmOutput)
		if vmOutput.ReturnCode != vmcommon.Ok {
			fmt.Printf("Deployed %d SCs\n", i)
			fmt.Printf(vmOutput.ReturnMessage)
		}
		require.Equal(t, vmcommon.Ok, vmOutput.ReturnCode)
	}
}

func TestExecution_Deploy_DisallowFloatingPoint(t *testing.T) {
	newAddress := []byte("new smartcontract")
	host := defaultTestArwenForDeployment(t, 24, newAddress)
	input := DefaultTestContractCreateInput()
	input.CallValue = big.NewInt(88)
	input.GasProvided = 1000
	input.ContractCode = arwen.GetTestSCCode("num-with-fp", "../../")

	vmOutput, err := host.RunSmartContractCreate(input)
	require.Nil(t, err)
	require.NotNil(t, vmOutput)
	require.Equal(t, vmcommon.ContractInvalid, vmOutput.ReturnCode)
}

func TestExecution_CallGetUserAccountErr(t *testing.T) {
	stubBlockchainHook := &contextmock.BlockchainHookStub{}

	errGetAccount := errors.New("get code error")

	host := defaultTestArwen(t, stubBlockchainHook)
	input := DefaultTestContractCallInput()
	stubBlockchainHook.GetUserAccountCalled = func(address []byte) (vmcommon.UserAccountHandler, error) {
		return nil, errGetAccount
	}

	input.GasProvided = 100
	vmOutput, err := host.RunSmartContractCall(input)
	require.Nil(t, err)
	require.NotNil(t, vmOutput)
	require.Equal(t, vmcommon.ContractNotFound, vmOutput.ReturnCode)
	require.Equal(t, arwen.ErrContractNotFound.Error(), vmOutput.ReturnMessage)
}

func TestExecution_NotEnoughGasForGetCode(t *testing.T) {
	stubBlockchainHook := &contextmock.BlockchainHookStub{}

	host := defaultTestArwen(t, stubBlockchainHook)
	input := DefaultTestContractCallInput()

	input.GasProvided = 0
	vmOutput, err := host.RunSmartContractCall(input)
	require.Nil(t, err)
	require.NotNil(t, vmOutput)
	require.Equal(t, vmcommon.OutOfGas, vmOutput.ReturnCode)
}

func TestExecution_CallOutOfGas(t *testing.T) {
	code := arwen.GetTestSCCode("counter", "../../")
	host, _ := defaultTestArwenForCall(t, code, nil)
	input := DefaultTestContractCallInput()
	input.Function = increment

	vmOutput, err := host.RunSmartContractCall(input)
	require.Nil(t, err)
	require.NotNil(t, vmOutput)
	require.Equal(t, vmcommon.OutOfGas, vmOutput.ReturnCode)
	require.Equal(t, arwen.ErrNotEnoughGas.Error(), vmOutput.ReturnMessage)
}

func TestExecution_CallWasmerError(t *testing.T) {
	code := []byte("not WASM")
	host, _ := defaultTestArwenForCall(t, code, nil)
	input := DefaultTestContractCallInput()
	input.GasProvided = 100000
	input.Function = increment

	vmOutput, err := host.RunSmartContractCall(input)
	require.Nil(t, err)
	require.NotNil(t, vmOutput)
	require.Equal(t, vmcommon.ContractInvalid, vmOutput.ReturnCode)
}

func TestExecution_CallSCMethod(t *testing.T) {
	code := arwen.GetTestSCCode("counter", "../../")
	host, _ := defaultTestArwenForCall(t, code, nil)
	input := DefaultTestContractCallInput()
	input.GasProvided = 100000

	// Calling init() directly is forbidden
	input.Function = "init"
	vmOutput, err := host.RunSmartContractCall(input)
	require.Nil(t, err)
	require.NotNil(t, vmOutput)
	require.Equal(t, vmcommon.UserError, vmOutput.ReturnCode)
	require.Equal(t, arwen.ErrInitFuncCalledInRun.Error(), vmOutput.ReturnMessage)

	// Calling callBack() directly is forbidden
	input.Function = "callBack"
	vmOutput, err = host.RunSmartContractCall(input)
	require.Nil(t, err)
	require.NotNil(t, vmOutput)
	require.Equal(t, vmcommon.UserError, vmOutput.ReturnCode)
	require.Equal(t, arwen.ErrCallBackFuncCalledInRun.Error(), vmOutput.ReturnMessage)

	// Handle calling a missing function
	input.Function = "wrong"
	vmOutput, err = host.RunSmartContractCall(input)
	require.Nil(t, err)
	require.NotNil(t, vmOutput)
	require.Equal(t, vmcommon.FunctionNotFound, vmOutput.ReturnCode)
}

func TestExecution_Call_Successful(t *testing.T) {
	code := arwen.GetTestSCCode("counter", "../../")
	host, stubBlockchainHook := defaultTestArwenForCall(t, code, nil)
	stubBlockchainHook.GetStorageDataCalled = func(scAddress []byte, key []byte) ([]byte, error) {
		return big.NewInt(1001).Bytes(), nil
	}
	input := DefaultTestContractCallInput()
	input.GasProvided = 100000
	input.Function = increment

	vmOutput, err := host.RunSmartContractCall(input)
	require.Nil(t, err)
	require.NotNil(t, vmOutput)
	require.Equal(t, vmcommon.Ok, vmOutput.ReturnCode)
	require.Len(t, vmOutput.OutputAccounts, 1)
	require.Len(t, vmOutput.OutputAccounts[string(parentAddress)].StorageUpdates, 1)

	storedBytes := vmOutput.OutputAccounts[string(parentAddress)].StorageUpdates[string(counterKey)].Data
	require.Equal(t, big.NewInt(1002).Bytes(), storedBytes)
}

func TestExecution_RawRustContract(t *testing.T) {
	code := arwen.GetTestSCCode("rawrs", "../../")
	host, _ := defaultTestArwenForCall(t, code, nil)

	input := DefaultTestContractCallInput()
	input.GasProvided = 1000
	input.Function = "method"

	vmOutput, err := host.RunSmartContractCall(input)
	require.Nil(t, err)
	require.NotNil(t, vmOutput)
	require.Equal(t, vmcommon.Ok, vmOutput.ReturnCode)
	require.Equal(t, [][]byte{{4}}, vmOutput.ReturnData)
}

func TestExecution_Call_GasConsumptionOnLocals(t *testing.T) {
	gasWithZeroLocals, gasSchedule := callCustomSCAndGetGasUsed(t, 0)
	costPerLocal := uint64(gasSchedule.WASMOpcodeCost.LocalAllocate)

	UnmeteredLocals := uint64(gasSchedule.WASMOpcodeCost.LocalsUnmetered)

	// Any number of local variables below `UnmeteredLocals` must be instantiated
	// without metering, i.e. gas-free.
	for _, locals := range []uint64{1, UnmeteredLocals / 2, UnmeteredLocals} {
		gasUsed, _ := callCustomSCAndGetGasUsed(t, locals)
		require.Equal(t, gasWithZeroLocals, gasUsed)
	}

	// Any number of local variables above `UnmeteredLocals` must be instantiated
	// with metering, i.e. will cost gas.
	for _, locals := range []uint64{UnmeteredLocals + 1, UnmeteredLocals * 2, UnmeteredLocals * 4} {
		gasUsed, _ := callCustomSCAndGetGasUsed(t, locals)
		meteredLocals := locals - UnmeteredLocals
		costOfLocals := costPerLocal * meteredLocals
		expectedGasUsed := gasWithZeroLocals + costOfLocals
		require.Equal(t, expectedGasUsed, gasUsed)
	}
}

func callCustomSCAndGetGasUsed(t *testing.T, locals uint64) (uint64, *config.GasCost) {
	code := makeBytecodeWithLocals(locals)
	host, _ := defaultTestArwenForCall(t, code, nil)
	gasSchedule := host.Metering().GasSchedule()

	gasLimit := uint64(100000)
	input := DefaultTestContractCallInput()
	input.GasProvided = gasLimit
	input.Function = "answer"

	vmOutput, err := host.RunSmartContractCall(input)
	require.Nil(t, err)
	require.NotNil(t, vmOutput)

	compilationCost := uint64(len(code)) * gasSchedule.BaseOperationCost.CompilePerByte
	return gasLimit - vmOutput.GasRemaining - compilationCost, gasSchedule
}

func TestExecution_ExecuteOnSameContext_Simple(t *testing.T) {
	parentCode := arwen.GetTestSCCode("exec-same-ctx-simple-parent", "../../")
	childCode := arwen.GetTestSCCode("exec-same-ctx-simple-child", "../../")

	host, _ := defaultTestArwenForTwoSCs(t, parentCode, childCode, nil, nil)
	input := DefaultTestContractCallInput()
	input.RecipientAddr = parentAddress
	input.Function = parentFunctionChildCall
	input.GasProvided = gasProvided

	vmOutput, err := host.RunSmartContractCall(input)
	require.Nil(t, err)
	require.NotNil(t, vmOutput)
	require.Equal(t, vmcommon.Ok, vmOutput.ReturnCode)

	expectedVMOutput := expectedVMOutputSameCtxSimple(parentCode, childCode)
	require.Equal(t, expectedVMOutput, vmOutput)
}

func TestExecution_Call_Breakpoints(t *testing.T) {
	t.Parallel()

	code := arwen.GetTestSCCode("breakpoint", "../../")
	host, _ := defaultTestArwenForCall(t, code, nil)
	input := DefaultTestContractCallInput()
	input.GasProvided = 100000
	input.Function = "testFunc"

	// Send the number 15 to the SC, causing it to finish with the number 100
	input.Arguments = [][]byte{{15}}
	vmOutput, err := host.RunSmartContractCall(input)
	require.Nil(t, err)
	require.NotNil(t, vmOutput)
	require.Equal(t, vmcommon.Ok, vmOutput.ReturnCode)
	require.Equal(t, [][]byte{{100}}, vmOutput.ReturnData)

	// Send the number 1 to the SC, causing it to exit with ReturnMessage "exit
	// here" if the breakpoint mechanism works properly, or with the
	// ReturnMessage "exit later" if the breakpoint mechanism fails to stop the
	// execution.
	input.Arguments = [][]byte{{1}}
	vmOutput, err = host.RunSmartContractCall(input)
	require.Nil(t, err)
	require.NotNil(t, vmOutput)
	require.Equal(t, vmcommon.UserError, vmOutput.ReturnCode)
	require.Len(t, vmOutput.ReturnData, 0)
	require.Equal(t, "exit here", vmOutput.ReturnMessage)
}

func TestExecution_ExecuteOnSameContext_Prepare(t *testing.T) {
	parentCode := arwen.GetTestSCCode("exec-same-ctx-parent", "../../")
	parentSCBalance := big.NewInt(1000)

	// Execute the parent SC method "parentFunctionPrepare", which sets storage,
	// finish data and performs a transfer. This step validates the test to the
	// actual call to ExecuteOnSameContext().
	host, _ := defaultTestArwenForCall(t, parentCode, parentSCBalance)
	input := DefaultTestContractCallInput()
	input.RecipientAddr = parentAddress
	input.Function = "parentFunctionPrepare"
	input.GasProvided = gasProvided

	vmOutput, err := host.RunSmartContractCall(input)
	require.Nil(t, err)
	require.Equal(t, vmcommon.Ok, vmOutput.ReturnCode)

	expectedVMOutput := expectedVMOutputSameCtxPrepare(parentCode)
	require.Equal(t, expectedVMOutput, vmOutput)
}

func TestExecution_ExecuteOnSameContext_Wrong(t *testing.T) {
	parentCode := arwen.GetTestSCCode("exec-same-ctx-parent", "../../")
	parentSCBalance := big.NewInt(1000)

	// Call parentFunctionWrongCall() of the parent SC, which will try to call a
	// non-existing SC.
	host, _ := defaultTestArwenForCall(t, parentCode, parentSCBalance)
	input := DefaultTestContractCallInput()
	input.RecipientAddr = parentAddress
	input.Function = "parentFunctionWrongCall"
	input.GasProvided = gasProvided

	vmOutput, err := host.RunSmartContractCall(input)
	require.Nil(t, err)
	require.NotNil(t, vmOutput)

	if host.Runtime().ElrondSyncExecAPIErrorShouldFailExecution() == false {
		require.Equal(t, vmcommon.Ok, vmOutput.ReturnCode)
		expectedVMOutput := expectedVMOutputSameCtxWrongContractCalled(parentCode)
		require.Equal(t, expectedVMOutput, vmOutput)
	} else {
		require.Equal(t, vmcommon.ExecutionFailed, vmOutput.ReturnCode)
		require.Equal(t, "account not found", vmOutput.ReturnMessage)
		require.Zero(t, vmOutput.GasRemaining)
	}
}

func TestExecution_ExecuteOnSameContext_OutOfGas(t *testing.T) {
	// Scenario:
	// Parent sets data into the storage, finishes data and creates a bigint
	// Parent calls executeOnSameContext, sending some value as well
	// Parent provides insufficient gas to executeOnSameContext (enoguh to start the SC though)
	// Child SC starts executing: sets data into the storage, finishes data and changes the bigint
	// Child starts an infinite loop, which must surely end with OutOfGas
	// Execution returns to parent, which finishes with the result of executeOnSameContext
	// Assertions: modifications made by the child are did not take effect
	// Assertions: the value sent by the parent to the child was returned to the parent
	// Assertions: the parent lost all the gas provided to executeOnSameContext
	parentCode := arwen.GetTestSCCode("exec-same-ctx-parent", "../../")
	childCode := arwen.GetTestSCCode("exec-same-ctx-child", "../../")

	// Call parentFunctionChildCall_OutOfGas() of the parent SC, which will call
	// the child SC using executeOnSameContext() with sufficient gas for
	// compilation and starting, but the child starts an infinite loop which will
	// end in OutOfGas.
	host, _ := defaultTestArwenForTwoSCs(t, parentCode, childCode, nil, nil)
	input := DefaultTestContractCallInput()
	input.RecipientAddr = parentAddress
	input.Function = "parentFunctionChildCall_OutOfGas"
	input.GasProvided = gasProvided

	vmOutput, err := host.RunSmartContractCall(input)
	require.Nil(t, err)
	require.NotNil(t, vmOutput)

	if host.Runtime().ElrondSyncExecAPIErrorShouldFailExecution() == false {
		require.Equal(t, vmcommon.Ok, vmOutput.ReturnCode)
		expectedVMOutput := expectedVMOutputSameCtxOutOfGas(parentCode, childCode)
		require.Equal(t, expectedVMOutput, vmOutput)
	} else {
		require.Equal(t, vmcommon.ExecutionFailed, vmOutput.ReturnCode)
		require.Equal(t, arwen.ErrNotEnoughGas.Error(), vmOutput.ReturnMessage)
		require.Zero(t, vmOutput.GasRemaining)
	}
}

func TestExecution_ExecuteOnSameContext_Successful(t *testing.T) {
	parentCode := arwen.GetTestSCCode("exec-same-ctx-parent", "../../")
	childCode := arwen.GetTestSCCode("exec-same-ctx-child", "../../")

	// Call parentFunctionChildCall() of the parent SC, which will call the child
	// SC and pass some arguments using executeOnSameContext().
	host, _ := defaultTestArwenForTwoSCs(t, parentCode, childCode, nil, nil)
	input := DefaultTestContractCallInput()
	input.RecipientAddr = parentAddress
	input.Function = parentFunctionChildCall
	input.GasProvided = gasProvided

	vmOutput, err := host.RunSmartContractCall(input)
	require.Nil(t, err)
	require.Equal(t, vmcommon.Ok, vmOutput.ReturnCode)
	expectedVMOutput := expectedVMOutputSameCtxSuccessfulChildCall(parentCode, childCode)
	require.Equal(t, expectedVMOutput, vmOutput)
}

func TestExecution_ExecuteOnSameContext_Successful_BigInts(t *testing.T) {
	parentCode := arwen.GetTestSCCode("exec-same-ctx-parent", "../../")
	childCode := arwen.GetTestSCCode("exec-same-ctx-child", "../../")

	// Call parentFunctionChildCall_BigInts() of the parent SC, which will call a
	// method of the child SC that takes some big Int references as arguments and
	// produce a new big Int out of the arguments.
	host, _ := defaultTestArwenForTwoSCs(t, parentCode, childCode, nil, nil)
	input := DefaultTestContractCallInput()
	input.RecipientAddr = parentAddress
	input.Function = "parentFunctionChildCall_BigInts"
	input.GasProvided = gasProvided

	vmOutput, err := host.RunSmartContractCall(input)
	require.Nil(t, err)
	require.Equal(t, vmcommon.Ok, vmOutput.ReturnCode)
	expectedVMOutput := expectedVMOutputSameCtxSuccessfulChildCallBigInts(parentCode, childCode)
	require.Equal(t, expectedVMOutput, vmOutput)
}

func TestExecution_ExecuteOnSameContext_Recursive_Direct(t *testing.T) {
	// Scenario:
	// SC has a method "callRecursive" which takes a byte as argument (number of recursive calls)
	// callRecursive() saves to storage "keyNNN" → "valueNNN", where NNN is the argument
	// callRecursive() saves to storage a counter starting at 1, increased by every recursive call
	// callRecursive() creates a bigInt and increments it with every iteration
	// callRecursive() finishes "finishNNN" in each iteration
	// callRecursive() calls itself using executeOnSameContext(), with the argument decremented
	// callRecursive() handles argument == 0 as follows: saves to storage the
	//		value of the bigInt counter, then exits without recursive call
	// Assertions: the VMOutput must contain as many StorageUpdates as the argument requires
	// Assertions: the VMOutput must contain as many finished values as the argument requires
	// Assertions: there must be a StorageUpdate with the value of the bigInt counter
	code := arwen.GetTestSCCode("exec-same-ctx-recursive", "../../")
	scBalance := big.NewInt(1000)

	host, _ := defaultTestArwenForCall(t, code, scBalance)
	input := DefaultTestContractCallInput()
	input.RecipientAddr = parentAddress
	input.Function = callRecursive
	input.GasProvided = gasProvided

	recursiveCalls := byte(5)
	input.Arguments = [][]byte{
		{recursiveCalls},
	}

	vmOutput, err := host.RunSmartContractCall(input)
	require.Nil(t, err)
	require.NotNil(t, vmOutput)
	require.Equal(t, vmcommon.Ok, vmOutput.ReturnCode)

	// TODO set proper gas calculation in the expectedVMOutput, like the other
	// tests
	expectedVMOutput := expectedVMOutputSameCtxRecursiveDirect(code, int(recursiveCalls))
	expectedVMOutput.GasRemaining = vmOutput.GasRemaining
	require.Equal(t, expectedVMOutput, vmOutput)
	require.Equal(t, int64(recursiveCalls+1), host.BigInt().GetOne(16).Int64())
}

func TestExecution_ExecuteOnSameContext_Recursive_Direct_ErrMaxInstances(t *testing.T) {
	code := arwen.GetTestSCCode("exec-same-ctx-recursive", "../../")
	scBalance := big.NewInt(1000)

	host, _ := defaultTestArwenForCall(t, code, scBalance)
	input := DefaultTestContractCallInput()
	input.RecipientAddr = parentAddress
	input.Function = callRecursive
	input.GasProvided = gasProvided

	recursiveCalls := byte(11)
	input.Arguments = [][]byte{
		{recursiveCalls},
	}

	vmOutput, err := host.RunSmartContractCall(input)
	require.Nil(t, err)
	require.NotNil(t, vmOutput)
	if host.Runtime().ElrondSyncExecAPIErrorShouldFailExecution() == false {
		require.Equal(t, vmcommon.Ok, vmOutput.ReturnCode)
		expectedVMOutput := expectedVMOutputSameCtxRecursiveDirectErrMaxInstances(code, int(recursiveCalls))
		expectedVMOutput.GasRemaining = vmOutput.GasRemaining
		require.Equal(t, expectedVMOutput, vmOutput)
		require.Equal(t, int64(1), host.BigInt().GetOne(16).Int64())
	} else {
		require.Equal(t, vmcommon.ExecutionFailed, vmOutput.ReturnCode)
		require.Equal(t, arwen.ErrExecutionFailed.Error(), vmOutput.ReturnMessage)
		require.Zero(t, vmOutput.GasRemaining)
	}
}

func TestExecution_ExecuteOnSameContext_Recursive_Mutual_Methods(t *testing.T) {
	// Scenario:
	// SC has a method "callRecursiveMutualMethods" which takes a byte as
	//		argument (number of recursive calls)
	// callRecursiveMutualMethods() sets the finish value "start recursive mutual calls"
	// callRecursiveMutualMethods() calls recursiveMethodA() on the same context,
	//		passing the argument

	// recursiveMethodA() saves to storage "AkeyNNN" → "AvalueNNN", where NNN is the argument
	// recursiveMethodA() saves to storage a counter starting at 1, increased by every recursive call
	// recursiveMethodA() creates a bigInt and increments it with every iteration
	// recursiveMethodA() finishes "AfinishNNN" in each iteration
	// recursiveMethodA() calls recursiveMethodB() with the argument decremented
	// recursiveMethodB() is a copy of recursiveMethodA()
	// when argument == 0, either of them will save to storage the
	//		value of the bigInt counter, then exits without recursive call
	// callRecursiveMutualMethods() sets the finish value "end recursive mutual calls" and exits
	// Assertions: the VMOutput must contain as many StorageUpdates as the argument requires
	// Assertions: the VMOutput must contain as many finished values as the argument requires
	// Assertions: there must be a StorageUpdate with the value of the bigInt counter
	code := arwen.GetTestSCCode("exec-same-ctx-recursive", "../../")
	scBalance := big.NewInt(1000)

	host, _ := defaultTestArwenForCall(t, code, scBalance)
	input := DefaultTestContractCallInput()
	input.RecipientAddr = parentAddress
	input.Function = "callRecursiveMutualMethods"
	input.GasProvided = gasProvided

	recursiveCalls := byte(5)
	input.Arguments = [][]byte{
		{recursiveCalls},
	}

	vmOutput, err := host.RunSmartContractCall(input)
	require.Nil(t, err)
	require.NotNil(t, vmOutput)
	require.Equal(t, vmcommon.Ok, vmOutput.ReturnCode)

	// TODO set proper gas calculation in the expectedVMOutput, like the other
	// tests
	expectedVMOutput := expectedVMOutputSameCtxRecursiveMutualMethods(code, int(recursiveCalls))
	expectedVMOutput.GasRemaining = vmOutput.GasRemaining
	require.Equal(t, expectedVMOutput, vmOutput)
	require.Equal(t, int64(recursiveCalls+1), host.BigInt().GetOne(16).Int64())
}

func TestExecution_ExecuteOnSameContext_Recursive_Mutual_SCs(t *testing.T) {
	// Scenario:
	// Parent has method parentCallChild()
	// Child has method childCallParent()
	// The two methods are identical, just named differently
	// The methods do the following:
	//		parent: save to storage "PkeyNNN" → "PvalueNNN"
	//		parent:	finish "PfinishNNN"
	//		child:	save to storage "CkeyNNN" → "CvalueNNN"
	//		child:	finish "CfinishNNN"
	//		both:		increment a shared bigInt counter
	//		both:		whoever exits must save the shared bigInt counter to storage
	parentCode := arwen.GetTestSCCode("exec-same-ctx-recursive-parent", "../../")
	childCode := arwen.GetTestSCCode("exec-same-ctx-recursive-child", "../../")

	// Call parentFunctionChildCall() of the parent SC, which will call the child
	// SC and pass some arguments using executeOnDestContext().
	host, _ := defaultTestArwenForTwoSCs(t, parentCode, childCode, nil, nil)
	input := DefaultTestContractCallInput()
	input.RecipientAddr = parentAddress
	input.Function = parentCallsChild
	input.GasProvided = gasProvided

	recursiveCalls := byte(4)
	input.Arguments = [][]byte{
		{recursiveCalls},
	}

	vmOutput, err := host.RunSmartContractCall(input)
	require.Nil(t, err)
	require.NotNil(t, vmOutput)
	require.Equal(t, vmcommon.Ok, vmOutput.ReturnCode)

	// TODO set proper gas calculation in the expectedVMOutput, like the other
	// tests
	expectedVMOutput := expectedVMOutputSameCtxRecursiveMutualSCs(parentCode, childCode, int(recursiveCalls))
	expectedVMOutput.GasRemaining = vmOutput.GasRemaining
	require.Equal(t, expectedVMOutput, vmOutput)
	require.Equal(t, int64(recursiveCalls+1), host.BigInt().GetOne(88).Int64())
}

func TestExecution_ExecuteOnSameContext_Recursive_Mutual_SCs_OutOfGas(t *testing.T) {
	parentCode := arwen.GetTestSCCode("exec-same-ctx-recursive-parent", "../../")
	childCode := arwen.GetTestSCCode("exec-same-ctx-recursive-child", "../../")

	// Call parentFunctionChildCall() of the parent SC, which will call the child
	// SC and pass some arguments using executeOnDestContext().
	host, _ := defaultTestArwenForTwoSCs(t, parentCode, childCode, nil, nil)
	input := DefaultTestContractCallInput()
	input.RecipientAddr = parentAddress
	input.Function = parentCallsChild
	input.GasProvided = 10000

	recursiveCalls := byte(5)
	input.Arguments = [][]byte{
		{recursiveCalls},
	}

	vmOutput, err := host.RunSmartContractCall(input)
	require.Nil(t, err)
	require.NotNil(t, vmOutput)

	if host.Runtime().ElrondSyncExecAPIErrorShouldFailExecution() == false {
		require.Equal(t, vmcommon.OutOfGas, vmOutput.ReturnCode)
		require.Equal(t, arwen.ErrNotEnoughGas.Error(), vmOutput.ReturnMessage)
	} else {
		require.Equal(t, vmcommon.ExecutionFailed, vmOutput.ReturnCode)
		require.Equal(t, arwen.ErrExecutionFailed.Error(), vmOutput.ReturnMessage)
	}
}

func TestExecution_ExecuteOnDestContext_Prepare(t *testing.T) {
	parentCode := arwen.GetTestSCCode("exec-dest-ctx-parent", "../../")
	parentSCBalance := big.NewInt(1000)

	// Execute the parent SC method "parentFunctionPrepare", which sets storage,
	// finish data and performs a transfer. This step validates the test to the
	// actual call to ExecuteOnSameContext().
	host, _ := defaultTestArwenForCall(t, parentCode, parentSCBalance)
	input := DefaultTestContractCallInput()
	input.RecipientAddr = parentAddress
	input.Function = "parentFunctionPrepare"
	input.GasProvided = gasProvided

	vmOutput, err := host.RunSmartContractCall(input)
	require.Nil(t, err)
	require.NotNil(t, vmOutput)
	require.Equal(t, vmcommon.Ok, vmOutput.ReturnCode)

	expectedVMOutput := expectedVMOutputDestCtxPrepare(parentCode)
	require.Equal(t, expectedVMOutput, vmOutput)
}

func TestExecution_ExecuteOnDestContext_Wrong(t *testing.T) {
	parentCode := arwen.GetTestSCCode("exec-dest-ctx-parent", "../../")
	parentSCBalance := big.NewInt(1000)

	// Call parentFunctionWrongCall() of the parent SC, which will try to call a
	// non-existing SC.
	host, _ := defaultTestArwenForCall(t, parentCode, parentSCBalance)
	input := DefaultTestContractCallInput()
	input.RecipientAddr = parentAddress
	input.Function = "parentFunctionWrongCall"
	input.GasProvided = gasProvided

	vmOutput, err := host.RunSmartContractCall(input)
	require.Nil(t, err)
	require.NotNil(t, vmOutput)

	if host.Runtime().ElrondSyncExecAPIErrorShouldFailExecution() == false {
		require.Equal(t, vmcommon.Ok, vmOutput.ReturnCode)
		expectedVMOutput := expectedVMOutputDestCtxWrongContractCalled(parentCode)
		require.Equal(t, expectedVMOutput, vmOutput)
	} else {
		require.Equal(t, vmcommon.ExecutionFailed, vmOutput.ReturnCode)
		require.Equal(t, "account not found", vmOutput.ReturnMessage)
		require.Zero(t, vmOutput.GasRemaining)
	}
}

func TestExecution_ExecuteOnDestContext_OutOfGas(t *testing.T) {
	// Scenario:
	// Parent sets data into the storage, finishes data and creates a bigint
	// Parent calls executeOnDestContext, sending some value as well
	// Parent provides insufficient gas to executeOnDestContext (enoguh to start the SC though)
	// Child SC starts executing: sets data into the storage, finishes data and changes the bigint
	// Child starts an infinite loop, which must surely end with OutOfGas
	// Execution returns to parent, which finishes with the result of executeOnDestContext
	// Assertions: modifications made by the child are did not take effect (no OutputAccount is created)
	// Assertions: the value sent by the parent to the child was returned to the parent
	// Assertions: the parent lost all the gas provided to executeOnDestContext
	parentCode := arwen.GetTestSCCode("exec-dest-ctx-parent", "../../")
	childCode := arwen.GetTestSCCode("exec-dest-ctx-child", "../../")

	// Call parentFunctionChildCall_OutOfGas() of the parent SC, which will call
	// the child SC using executeOnDestContext() with sufficient gas for
	// compilation and starting, but the child starts an infinite loop which will
	// end in OutOfGas.
	host, _ := defaultTestArwenForTwoSCs(t, parentCode, childCode, nil, nil)
	input := DefaultTestContractCallInput()
	input.RecipientAddr = parentAddress
	input.Function = "parentFunctionChildCall_OutOfGas"
	input.GasProvided = gasProvided

	vmOutput, err := host.RunSmartContractCall(input)
	require.Nil(t, err)
	require.NotNil(t, vmOutput)

	if host.Runtime().ElrondSyncExecAPIErrorShouldFailExecution() == false {
		require.Equal(t, vmcommon.Ok, vmOutput.ReturnCode)
		expectedVMOutput := expectedVMOutputDestCtxOutOfGas(parentCode)
		require.Equal(t, expectedVMOutput, vmOutput)
		require.Equal(t, int64(42), host.BigInt().GetOne(12).Int64())
	} else {
		require.Equal(t, vmcommon.ExecutionFailed, vmOutput.ReturnCode)
		require.Equal(t, arwen.ErrNotEnoughGas.Error(), vmOutput.ReturnMessage)
		require.Zero(t, vmOutput.GasRemaining)
	}
}

func TestExecution_ExecuteOnDestContext_Successful(t *testing.T) {
	parentCode := arwen.GetTestSCCode("exec-dest-ctx-parent", "../../")
	childCode := arwen.GetTestSCCode("exec-dest-ctx-child", "../../")

	// Call parentFunctionChildCall() of the parent SC, which will call the child
	// SC and pass some arguments using executeOnDestContext().
	host, _ := defaultTestArwenForTwoSCs(t, parentCode, childCode, nil, nil)
	input := DefaultTestContractCallInput()
	input.RecipientAddr = parentAddress
	input.Function = parentFunctionChildCall
	input.GasProvided = gasProvided

	vmOutput, err := host.RunSmartContractCall(input)
	require.Nil(t, err)
	require.NotNil(t, vmOutput)
	require.Equal(t, vmcommon.Ok, vmOutput.ReturnCode)

	expectedVMOutput := expectedVMOutputDestCtxSuccessfulChildCall(parentCode, childCode)
	expectedVMOutput.OutputAccounts[string(parentAddress)].StorageUpdates[string(childKey)] = &vmcommon.StorageUpdate{Offset: childKey, Data: nil}
	assert.Equal(t, expectedVMOutput, vmOutput)
}

func TestExecution_ExecuteOnDestContext_GasRemaining(t *testing.T) {
	// This test ensures that host.ExecuteOnDestContext() calls
	// metering.GasLeft() on the Wasmer instance of the child, and not of the
	// parent.

	parentCode := arwen.GetTestSCCode("exec-dest-ctx-parent", "../../")
	childCode := arwen.GetTestSCCode("exec-dest-ctx-child", "../../")

	// Pretend that the execution of the parent SC was requested, with the
	// following ContractCallInput:
	input := DefaultTestContractCallInput()
	input.RecipientAddr = parentAddress
	input.Function = "parentFunctionChildCall_BigInts"
	input.GasProvided = gasProvided

	// Initialize the VM with the parent SC and child SC, but without really
	// executing the parent. The initialization emulates the behavior of
	// host.doRunSmartContractCall(). Gas cost for compilation is skipped.
	host, _ := defaultTestArwenForTwoSCs(t, parentCode, childCode, nil, nil)
	host.InitState()

	_, _, metering, output, runtime, async, storage := host.GetContexts()
	runtime.InitStateFromContractCallInput(input)
	async.InitStateFromInput(input)
	output.AddTxValueToAccount(input.RecipientAddr, input.CallValue)
	storage.SetAddress(runtime.GetSCAddress())
	_ = metering.DeductInitialGasForExecution([]byte{})

	contract, err := runtime.GetSCCode()
	require.Nil(t, err)

	vmInput := runtime.GetVMInput()
	err = runtime.StartWasmerInstance(contract, vmInput.GasProvided, false)
	require.Nil(t, err)

	// Use a lot of gas on the parent contract
	metering.UseGas(500000)
	require.Equal(t, input.GasProvided-500001, metering.GasLeft())

	// Create a second ContractCallInput, used to call the child SC using
	// host.ExecuteOnDestContext().
	childInput := arwen.MakeContractCallInput(parentAddress, childAddress, "childFunction", 99)
	arwen.AddArgument(childInput, []byte("some data"))
	arwen.AddArgument(childInput, []byte("argument"))
	arwen.AddArgument(childInput, []byte("another argument"))
	childInput.GasProvided = 10000

	childOutput, _, err := host.ExecuteOnDestContext(childInput)
	require.Nil(t, err)
	require.NotNil(t, childOutput)
	require.Equal(t, uint64(7752), childOutput.GasRemaining)

	host.Clean()
}

func TestExecution_ExecuteOnDestContext_Successful_BigInts(t *testing.T) {
	parentCode := arwen.GetTestSCCode("exec-dest-ctx-parent", "../../")
	childCode := arwen.GetTestSCCode("exec-dest-ctx-child", "../../")

	// Call parentFunctionChildCall_BigInts() of the parent SC, which will call a
	// method of the child SC that takes some big Int references as arguments and
	// produce a new big Int out of the arguments.
	host, _ := defaultTestArwenForTwoSCs(t, parentCode, childCode, nil, nil)
	input := DefaultTestContractCallInput()
	input.RecipientAddr = parentAddress
	input.Function = "parentFunctionChildCall_BigInts"
	input.GasProvided = gasProvided

	vmOutput, err := host.RunSmartContractCall(input)
	require.Nil(t, err)
	require.NotNil(t, vmOutput)
	require.Equal(t, vmcommon.Ok, vmOutput.ReturnCode)

	expectedVMOutput := expectedVMOutputDestCtxSuccessfulChildCallBigInts(parentCode, childCode)
	require.Equal(t, expectedVMOutput, vmOutput)
}

func TestExecution_ExecuteOnDestContext_Recursive_Direct(t *testing.T) {
	code := arwen.GetTestSCCode("exec-dest-ctx-recursive", "../../")
	scBalance := big.NewInt(1000)

	host, _ := defaultTestArwenForCall(t, code, scBalance)
	input := DefaultTestContractCallInput()
	input.RecipientAddr = parentAddress
	input.Function = callRecursive
	input.GasProvided = gasProvided

	recursiveCalls := byte(6)
	input.Arguments = [][]byte{
		{recursiveCalls},
	}

	vmOutput, err := host.RunSmartContractCall(input)
	require.Nil(t, err)
	require.NotNil(t, vmOutput)
	require.Equal(t, vmcommon.Ok, vmOutput.ReturnCode)

	// TODO set proper gas calculation in the expectedVMOutput, like the other
	// tests
	expectedVMOutput := expectedVMOutputDestCtxRecursiveDirect(code, int(recursiveCalls))
	expectedVMOutput.GasRemaining = vmOutput.GasRemaining
	require.Equal(t, expectedVMOutput, vmOutput)
	require.Equal(t, int64(1), host.BigInt().GetOne(16).Int64())
}

func TestExecution_ExecuteOnDestContext_Recursive_Mutual_Methods(t *testing.T) {
	code := arwen.GetTestSCCode("exec-dest-ctx-recursive", "../../")
	scBalance := big.NewInt(1000)

	host, _ := defaultTestArwenForCall(t, code, scBalance)
	input := DefaultTestContractCallInput()
	input.RecipientAddr = parentAddress
	input.Function = "callRecursiveMutualMethods"
	input.GasProvided = gasProvided

	recursiveCalls := byte(7)
	input.Arguments = [][]byte{
		{recursiveCalls},
	}

	vmOutput, err := host.RunSmartContractCall(input)
	require.Nil(t, err)
	require.NotNil(t, vmOutput)
	require.Equal(t, vmcommon.Ok, vmOutput.ReturnCode)

	// TODO set proper gas calculation in the expectedVMOutput, like the other
	// tests
	expectedVMOutput := expectedVMOutputDestCtxRecursiveMutualMethods(code, int(recursiveCalls))
	expectedVMOutput.GasRemaining = vmOutput.GasRemaining
	require.Equal(t, expectedVMOutput, vmOutput)
	require.Equal(t, int64(0), host.BigInt().GetOne(16).Int64())
}

func TestExecution_ExecuteOnDestContext_Recursive_Mutual_SCs(t *testing.T) {
	parentCode := arwen.GetTestSCCode("exec-dest-ctx-recursive-parent", "../../")
	childCode := arwen.GetTestSCCode("exec-dest-ctx-recursive-child", "../../")

	// Call parentFunctionChildCall() of the parent SC, which will call the child
	// SC and pass some arguments using executeOnDestContext().
	host, _ := defaultTestArwenForTwoSCs(t, parentCode, childCode, nil, nil)
	input := DefaultTestContractCallInput()
	input.RecipientAddr = parentAddress
	input.Function = parentCallsChild
	input.GasProvided = gasProvided

	recursiveCalls := byte(6)
	input.Arguments = [][]byte{
		{recursiveCalls},
	}

	vmOutput, err := host.RunSmartContractCall(input)
	require.Nil(t, err)
	require.NotNil(t, vmOutput)
	require.Equal(t, vmcommon.Ok, vmOutput.ReturnCode)

	// TODO set proper gas calculation in the expectedVMOutput, like the other
	// tests
	expectedVMOutput := expectedVMOutputDestCtxRecursiveMutualSCs(parentCode, childCode, int(recursiveCalls))
	expectedVMOutput.GasRemaining = vmOutput.GasRemaining
	require.Equal(t, expectedVMOutput, vmOutput)
	require.Equal(t, int64(1), host.BigInt().GetOne(88).Int64())
}

func TestExecution_ExecuteOnDestContext_Recursive_Mutual_SCs_OutOfGas(t *testing.T) {
	parentCode := arwen.GetTestSCCode("exec-dest-ctx-recursive-parent", "../../")
	childCode := arwen.GetTestSCCode("exec-dest-ctx-recursive-child", "../../")

	// Call parentFunctionChildCall() of the parent SC, which will call the child
	// SC and pass some arguments using executeOnDestContext().
	host, _ := defaultTestArwenForTwoSCs(t, parentCode, childCode, nil, nil)
	input := DefaultTestContractCallInput()
	input.RecipientAddr = parentAddress
	input.Function = parentCallsChild
	input.GasProvided = 10000

	recursiveCalls := byte(5)
	input.Arguments = [][]byte{
		{recursiveCalls},
	}

	vmOutput, err := host.RunSmartContractCall(input)
	require.Nil(t, err)
	require.NotNil(t, vmOutput)

	if host.Runtime().ElrondSyncExecAPIErrorShouldFailExecution() == false {
		require.Equal(t, vmcommon.OutOfGas, vmOutput.ReturnCode)
		require.Equal(t, arwen.ErrNotEnoughGas.Error(), vmOutput.ReturnMessage)
	} else {
		require.Equal(t, vmcommon.ExecutionFailed, vmOutput.ReturnCode)
		require.Equal(t, arwen.ErrExecutionFailed.Error(), vmOutput.ReturnMessage)
		require.Zero(t, vmOutput.GasRemaining)
	}
}

func TestExecution_ExecuteOnDestContextByCaller_SimpleTransfer(t *testing.T) {
	// The child contract is designed to send some tokens back to its caller, as
	// many as requested. The parent calls the child using
	// executeOnDestContextByCaller(), which means that the child will not see
	// the parent as its caller, but the original caller of the transaction
	// instead. Thus the original caller (the user address) will receive 42
	// tokens, and not the parent, even if the parent is the one making the call
	// to the child.
	parentCode := arwen.GetTestSCCodeModule("exec-dest-ctx-by-caller/parent", "parent", "../../")
	childCode := arwen.GetTestSCCodeModule("exec-dest-ctx-by-caller/child", "child", "../../")

	host, _ := defaultTestArwenForTwoSCs(t, parentCode, childCode, nil, nil)
	input := DefaultTestContractCallInput()
	input.Function = "call_child"
	input.GasProvided = 2000

	vmOutput, err := host.RunSmartContractCall(input)
	require.Nil(t, err)
	require.NotNil(t, vmOutput)
	require.Equal(t, vmcommon.Ok, vmOutput.ReturnCode)

	// TODO calculate expected remaining gas properly, instead of copying it from
	// the actual vmOutput.
	expectedVMOutput := expectedVMOutputDestCtxByCallerSimpleTransfer(42)
	expectedVMOutput.GasRemaining = vmOutput.GasRemaining
	require.Equal(t, expectedVMOutput, vmOutput)
}

func TestExecution_AsyncCall_GasLimitConsumed(t *testing.T) {
	parentCode := arwen.GetTestSCCode("async-call-parent", "../../")
	childCode := arwen.GetTestSCCode("async-call-child", "../../")
	host, stubBlockchainHook := defaultTestArwenForTwoSCs(t, parentCode, childCode, nil, nil)
	stubBlockchainHook.GetUserAccountCalled = func(scAddress []byte) (vmcommon.UserAccountHandler, error) {
		if bytes.Equal(scAddress, parentAddress) {
			return &contextmock.StubAccount{
				Address: parentAddress,
				Balance: big.NewInt(1000),
			}, nil
		}
		return nil, errAccountNotFound
	}
	stubBlockchainHook.GetCodeCalled = func(account vmcommon.UserAccountHandler) []byte {
		if bytes.Equal(parentAddress, account.AddressBytes()) {
			return parentCode
		}
		return nil
	}
	stubBlockchainHook.GetShardOfAddressCalled = func(address []byte) uint32 {
		if bytes.Equal(address, parentAddress) {
			return 0
		}
		return 1
	}

	input := DefaultTestContractCallInput()
	input.RecipientAddr = parentAddress
	input.Function = parentPerformAsyncCall
	input.GasProvided = 1000000
	input.Arguments = [][]byte{{0}}

	vmOutput, err := host.RunSmartContractCall(input)
	require.Nil(t, err)
	require.NotNil(t, vmOutput)
	require.Zero(t, vmOutput.GasRemaining)
}

func TestExecution_AsyncCall(t *testing.T) {
	// Scenario
	// Parent SC calls Child SC
	// Before asyncCall, Parent sets storage, makes a value transfer to ThirdParty and finishes some data
	// Parent performs asyncCall to Child with a sufficient amount of ERD, with arguments:
	//	* the address of ThirdParty
	//	* number of ERD the Child should send to ThirdParty
	//  * a string, to be set as the data on the transfer to ThirdParty
	// Child stores the received arguments to storage
	// Child performs two transfers:
	//	* to ThirdParty, sending the amount of ERD specified as argument in asyncCall
	//	* to the Vault, a fixed address known by the Child, sending exactly 4 ERD with the data provided by Parent
	// Child finishes with "thirdparty" if the transfer to ThirdParty was successful
	// Child finishes with "vault" if the transfer to Vault was successful
	// Parent callBack() verifies its arguments and expects both "thirdparty" and "vault"
	// Assertions: OutputAccounts for
	//		* Parent: negative balance delta (payment for child + thirdparty + vault => 2), storage
	//		* Child: zero balance delta, storage
	//		* ThirdParty: positive balance delta
	//		* Vault
	parentCode := arwen.GetTestSCCode("async-call-parent", "../../")
	childCode := arwen.GetTestSCCode("async-call-child", "../../")

	// Call parentFunctionChildCall() of the parent SC, which will call the child
	// SC and pass some arguments using asyncCall().
	host, _ := defaultTestArwenForTwoSCs(t, parentCode, childCode, nil, nil)
	input := DefaultTestContractCallInput()
	input.RecipientAddr = parentAddress
	input.Function = parentPerformAsyncCall
	input.GasProvided = 116000
	input.Arguments = [][]byte{{0}}

	vmOutput, err := host.RunSmartContractCall(input)
	require.Nil(t, err)
	require.NotNil(t, vmOutput)
	require.Equal(t, vmcommon.Ok, vmOutput.ReturnCode)

	// TODO calculate expected remaining gas properly, instead of copying it from
	// the actual vmOutput.
	expectedVMOutput := expectedVMOutputAsyncCall(parentCode, childCode)
	expectedVMOutput.GasRemaining = vmOutput.GasRemaining
	require.Equal(t, expectedVMOutput, vmOutput)
}

func TestExecution_AsyncCall_ChildFails(t *testing.T) {
	// Scenario
	// Identical to TestExecution_AsyncCall(), except that the child is
	// instructed to call signalError().
	// Because "vault" was not received by the callBack(), the Parent sends 4 ERD
	// to the Vault directly.
	parentCode := arwen.GetTestSCCode("async-call-parent", "../../")
	childCode := arwen.GetTestSCCode("async-call-child", "../../")

	// Call parentFunctionChildCall() of the parent SC, which will call the child
	// SC and pass some arguments using asyncCall().
	host, _ := defaultTestArwenForTwoSCs(t, parentCode, childCode, nil, nil)
	host.Metering().GasSchedule().ElrondAPICost.AsyncCallbackGasLock = 3000

	input := DefaultTestContractCallInput()
	input.RecipientAddr = parentAddress
	input.Function = parentPerformAsyncCall
	input.GasProvided = 1000000
	input.Arguments = [][]byte{{1}}
	input.CurrentTxHash = []byte("txhash")

	vmOutput, err := host.RunSmartContractCall(input)
	require.Nil(t, err)
	require.NotNil(t, vmOutput)
	require.Equal(t, vmcommon.Ok, vmOutput.ReturnCode)

	// TODO calculate expected remaining gas properly, instead of copying it from
	// the actual vmOutput.
	expectedVMOutput := expectedVMOutputAsyncCallChildFails(parentCode, childCode)
	expectedVMOutput.GasRemaining = vmOutput.GasRemaining
	require.Equal(t, expectedVMOutput, vmOutput)
}

func TestExecution_AsyncCall_CallBackFails(t *testing.T) {
	// Scenario
	// Identical to TestExecution_AsyncCall(), except that the callback is
	// instructed to call signalError().
	parentCode := arwen.GetTestSCCode("async-call-parent", "../../")
	childCode := arwen.GetTestSCCode("async-call-child", "../../")

	// Call parentFunctionChildCall() of the parent SC, which will call the child
	// SC and pass some arguments using asyncCall().
	host, _ := defaultTestArwenForTwoSCs(t, parentCode, childCode, nil, nil)
	input := DefaultTestContractCallInput()
	input.RecipientAddr = parentAddress
	input.Function = parentPerformAsyncCall
	input.GasProvided = 200000
	input.Arguments = [][]byte{{0, 3}}
	input.CurrentTxHash = []byte("txhash")

	vmOutput, err := host.RunSmartContractCall(input)
	require.Nil(t, err)
	require.NotNil(t, vmOutput)

	// TODO calculate expected remaining gas properly, instead of copying it from
	// the actual vmOutput.
	expectedVMOutput := expectedVMOutputAsyncCallCallBackFails(parentCode, childCode)
	expectedVMOutput.GasRemaining = vmOutput.GasRemaining
	require.Equal(t, expectedVMOutput, vmOutput)
}

func TestExecution_CreateNewContract_Success(t *testing.T) {
	parentCode := arwen.GetTestSCCode("deployer", "../../")
	childCode := arwen.GetTestSCCode("init-correct", "../../")
	parentBalance := big.NewInt(1000)

	host, stubBlockchainHook := defaultTestArwenForCall(t, parentCode, parentBalance)
	stubBlockchainHook.GetStorageDataCalled = func(address []byte, key []byte) ([]byte, error) {
		if bytes.Equal(address, parentAddress) {
			if bytes.Equal(key, []byte{'A'}) {
				return childCode, nil
			}
			return nil, nil
		}
		return nil, arwen.ErrInvalidAccount
	}

	input := DefaultTestContractCallInput()
	input.Function = "deployChildContract"
	input.Arguments = [][]byte{{'A'}, {0}}
	input.GasProvided = 1_000_000

	vmOutput, err := host.RunSmartContractCall(input)
	require.Nil(t, err)
	require.NotNil(t, vmOutput)
	require.Equal(t, vmcommon.Ok, vmOutput.ReturnCode)

	expectedVMOutput := expectedVMOutputCreateNewContractSuccess(parentCode, childCode)

	// TODO calculate expected remaining gas properly, instead of copying it from
	// the actual vmOutput.
	expectedVMOutput.GasRemaining = vmOutput.GasRemaining
	require.Equal(t, expectedVMOutput, vmOutput)
}

func TestExecution_CreateNewContract_Fail(t *testing.T) {
	parentCode := arwen.GetTestSCCode("deployer", "../../")
	childCode := arwen.GetTestSCCode("init-correct", "../../")
	parentBalance := big.NewInt(1000)

	host, stubBlockchainHook := defaultTestArwenForCall(t, parentCode, parentBalance)
	stubBlockchainHook.GetStorageDataCalled = func(address []byte, key []byte) ([]byte, error) {
		if bytes.Equal(address, parentAddress) {
			if bytes.Equal(key, []byte{'A'}) {
				return childCode, nil
			}
			return nil, nil
		}
		return nil, arwen.ErrInvalidAccount
	}

	input := DefaultTestContractCallInput()
	input.Function = "deployChildContract"
	input.Arguments = [][]byte{{'A'}, {1}}
	input.GasProvided = 1_000_000

	vmOutput, err := host.RunSmartContractCall(input)
	require.Nil(t, err)
	require.NotNil(t, vmOutput)
	require.Equal(t, vmcommon.Ok, vmOutput.ReturnCode)

	expectedVMOutput := expectedVMOutputCreateNewContractFail(parentCode, childCode)

	// TODO calculate expected remaining gas properly, instead of copying it from
	// the actual vmOutput.
	expectedVMOutput.GasRemaining = vmOutput.GasRemaining
	require.Equal(t, expectedVMOutput, vmOutput)
}

func TestExecution_Mocked_Wasmer_Instances(t *testing.T) {
	host, _, ibm := defaultTestArwenForCallWithInstanceMocks(t)

	parentInstance := ibm.CreateAndStoreInstanceMock(parentAddress)
	parentInstance.GetExports()["callChild"] = mockMethod(func() {
		host.Output().Finish([]byte("parent returns this"))
		host.Metering().UseGas(500)
		_, err := host.Storage().SetStorage([]byte("parent"), []byte("parent storage"))
		require.Nil(t, err)
		childInput := DefaultTestContractCallInput()
		childInput.CallerAddr = parentAddress
		childInput.RecipientAddr = childAddress
		childInput.CallValue = big.NewInt(4)
		childInput.Function = "doSomething"
		childInput.GasProvided = 1000
		_, _, err = host.ExecuteOnDestContext(childInput)
		require.Nil(t, err)
	})

	childInstance := ibm.CreateAndStoreInstanceMock(childAddress)
	childInstance.GetExports()["doSomething"] = mockMethod(func() {
		host.Output().Finish([]byte("child returns this"))
		host.Metering().UseGas(100)
		_, err := host.Storage().SetStorage([]byte("child"), []byte("child storage"))
		require.Nil(t, err)
	})

	input := DefaultTestContractCallInput()
	input.Function = "callChild"
	input.GasProvided = 1000

	vmOutput, err := host.RunSmartContractCall(input)
	require.Nil(t, err)
	require.NotNil(t, vmOutput)
	require.Equal(t, vmcommon.Ok, vmOutput.ReturnCode)

	expectedVMOutput := expectedVMOutputMockedWasmerInstances()
	expectedVMOutput.GasRemaining = 307
	require.Equal(t, expectedVMOutput, vmOutput)
}

func TestExecution_GasUsed_SingleContract(t *testing.T) {
	host, _, ibm := defaultTestArwenForCallWithInstanceMocks(t)
	host.Metering().GasSchedule().BaseOperationCost.CompilePerByte = 0
	host.Metering().GasSchedule().BaseOperationCost.AoTPreparePerByte = 0

	gasProvided := uint64(1000)
	gasUsedByParent := uint64(401)

	parentInstance := ibm.CreateAndStoreInstanceMock(parentAddress)
	parentInstance.GetExports()["doSomething"] = mockMethod(func() {
		host.Metering().UseGas(gasUsedByParent)
	})

	input := DefaultTestContractCallInput()
	input.Function = "doSomething"
	input.GasProvided = gasProvided

	vmOutput, err := host.RunSmartContractCall(input)
	require.Nil(t, err)
	require.NotNil(t, vmOutput)
	require.Equal(t, gasProvided-gasUsedByParent-1, vmOutput.GasRemaining)

	parentAccount := vmOutput.OutputAccounts[string(parentAddress)]
	require.Equal(t, gasUsedByParent+1, parentAccount.GasUsed)
}

func TestExecution_GasUsed_ExecuteOnSameCtx(t *testing.T) {
	host, _, ibm := defaultTestArwenForCallWithInstanceMocks(t)
	host.Metering().GasSchedule().BaseOperationCost.CompilePerByte = 1
	host.Metering().GasSchedule().BaseOperationCost.AoTPreparePerByte = 1

	gasProvided := uint64(1000)
	contractCompilationCost := uint64(32)
	gasUsedByParentExec := uint64(400)
	gasUsedByChildExec := uint64(200)
	gasUsedByParent := contractCompilationCost + gasUsedByParentExec + 1
	gasUsedByChild := contractCompilationCost + gasUsedByChildExec

	parentInstance := ibm.CreateAndStoreInstanceMock(parentAddress)
	parentInstance.GetExports()["function"] = mockMethod(func() {
		host.Metering().UseGas(gasUsedByParentExec)
		childInput := DefaultTestContractCallInput()
		childInput.CallerAddr = parentAddress
		childInput.RecipientAddr = childAddress
		childInput.GasProvided = 300
		err := host.ExecuteOnSameContext(childInput)
		require.Nil(t, err)
	})

	childInstance := ibm.CreateAndStoreInstanceMock(childAddress)
	childInstance.GetExports()["function"] = mockMethod(func() {
		host.Metering().UseGas(gasUsedByChildExec)
	})

	input := DefaultTestContractCallInput()
	input.Function = "function"
	input.GasProvided = gasProvided

	expectedGasRemaining := gasProvided - gasUsedByParent - gasUsedByChild - 1
	vmOutput, err := host.RunSmartContractCall(input)
	require.Nil(t, err)
	require.NotNil(t, vmOutput)
	require.Equal(t, expectedGasRemaining, vmOutput.GasRemaining)

	parentAccount := vmOutput.OutputAccounts[string(parentAddress)]
	require.Equal(t, gasUsedByParent, parentAccount.GasUsed)

	childAccount := vmOutput.OutputAccounts[string(childAddress)]
	require.Equal(t, gasUsedByChild+1, childAccount.GasUsed)
}

func mockMethod(method func()) wasmer.ExportedFunctionCallback {
	wrappedMethod := func(...interface{}) (wasmer.Value, error) {
		method()
		return wasmer.Void(), nil
	}
	return wrappedMethod
}

// makeBytecodeWithLocals rewrites the bytecode of "answer" to change the
// number of i64 locals it instantiates
func makeBytecodeWithLocals(numLocals uint64) []byte {
	originalCode := arwen.GetTestSCCode("answer", "../../")
	firstSlice := originalCode[:0x5B]
	secondSlice := originalCode[0x5C:]

	encodedNumLocals := arwen.U64ToLEB128(numLocals)
	extraBytes := len(encodedNumLocals) - 1

	result := make([]byte, 0)
	result = append(result, firstSlice...)
	result = append(result, encodedNumLocals...)
	result = append(result, secondSlice...)

	result[0x57] = byte(int(result[0x57]) + extraBytes)
	result[0x59] = byte(int(result[0x59]) + extraBytes)

	return result
}
