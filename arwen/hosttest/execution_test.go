package hosttest

import (
	"bytes"
	"errors"
	"fmt"
	"math"
	"math/big"
	"testing"

	"github.com/ElrondNetwork/arwen-wasm-vm/v1_3/arwen"
	"github.com/ElrondNetwork/arwen-wasm-vm/v1_3/config"
	contextmock "github.com/ElrondNetwork/arwen-wasm-vm/v1_3/mock/context"
	mock "github.com/ElrondNetwork/arwen-wasm-vm/v1_3/mock/context"
	worldmock "github.com/ElrondNetwork/arwen-wasm-vm/v1_3/mock/world"
	test "github.com/ElrondNetwork/arwen-wasm-vm/v1_3/testcommon"
	testcommon "github.com/ElrondNetwork/arwen-wasm-vm/v1_3/testcommon"
	"github.com/ElrondNetwork/arwen-wasm-vm/v1_3/wasmer"
	vmcommon "github.com/ElrondNetwork/elrond-vm-common"
	"github.com/stretchr/testify/require"
)

var counterKey = []byte("COUNTER")
var WASMLocalsLimit = uint64(4000)
var maxUint8AsInt = int(math.MaxUint8)
var newAddress = testcommon.MakeTestSCAddress("new smartcontract")

const (
	get                     = "get"
	increment               = "increment"
	callRecursive           = "callRecursive"
	parentCallsChild        = "parentCallsChild"
	parentPerformAsyncCall  = "parentPerformAsyncCall"
	parentFunctionChildCall = "parentFunctionChildCall"
)

func init() {
	test.ParentCompilationCostSameCtx = uint64(len(test.GetTestSCCode("exec-same-ctx-parent", "../../", "../../../")))
	test.ChildCompilationCostSameCtx = uint64(len(test.GetTestSCCode("exec-same-ctx-child", "../../", "../../../")))

	test.ParentCompilationCostDestCtx = uint64(len(test.GetTestSCCode("exec-dest-ctx-parent", "../../", "../../../")))
	test.ChildCompilationCostDestCtx = uint64(len(test.GetTestSCCode("exec-dest-ctx-child", "../../", "../../../")))
}

func TestSCMem(t *testing.T) {
	testString := "this is some random string of bytes"
	returnData := [][]byte{
		[]byte(testString),
		{35},
	}
	for _, c := range testString {
		returnData = append(returnData, []byte{byte(c)})
	}

	test.BuildInstanceCallTest(t).
		WithContracts(
			test.CreateInstanceContract(test.ParentAddress).
				WithCode(test.GetTestSCCode("misc", "../../"))).
		WithInput(test.CreateTestContractCallInputBuilder().
			WithRecipientAddr(test.ParentAddress).
			WithGasProvided(100000).
			WithFunction("iterate_over_byte_array").
			Build()).
		AndAssertResults(func(host arwen.VMHost, blockchainHook *contextmock.BlockchainHookStub, verify *test.VMOutputVerifier) {
			verify.
				Ok().
				ReturnData(returnData...)
		})
}

func TestExecution_DeployNewAddressErr(t *testing.T) {
	errNewAddress := errors.New("new address error")

	input := test.CreateTestContractCreateInputBuilder().
		WithGasProvided(1000).
		WithContractCode([]byte("contract")).
		Build()

	test.BuildInstanceCreatorTest(t).
		WithInput(input).
		WithSetup(func(host arwen.VMHost, stubBlockchainHook *contextmock.BlockchainHookStub) {
			stubBlockchainHook.GetUserAccountCalled = func(address []byte) (vmcommon.UserAccountHandler, error) {
				require.Equal(t, input.CallerAddr, address)
				return &contextmock.StubAccount{}, nil
			}
			stubBlockchainHook.NewAddressCalled = func(creatorAddress []byte, nonce uint64, vmType []byte) ([]byte, error) {
				require.Equal(t, input.CallerAddr, creatorAddress)
				require.Equal(t, uint64(0), nonce)
				require.Equal(t, test.DefaultVMType, vmType)
				return nil, errNewAddress
			}
		}).
		AndAssertResults(func(blockchainHook *contextmock.BlockchainHookStub, verify *test.VMOutputVerifier) {
			verify.
				ReturnCode(vmcommon.ExecutionFailed).
				ReturnMessage(errNewAddress.Error())
		})
}

func TestExecution_DeployOutOfGas(t *testing.T) {
	test.BuildInstanceCreatorTest(t).
		WithInput(test.CreateTestContractCreateInputBuilder().
			WithGasProvided(8).
			Build()).
		WithAddress(newAddress).
		AndAssertResults(func(blockchainHook *contextmock.BlockchainHookStub, verify *test.VMOutputVerifier) {
			verify.
				ReturnCode(vmcommon.OutOfGas).
				ReturnMessage(arwen.ErrNotEnoughGas.Error())
		})
}

func TestExecution_DeployNotWASM(t *testing.T) {
	test.BuildInstanceCreatorTest(t).
		WithInput(test.CreateTestContractCreateInputBuilder().
			WithGasProvided(9).
			WithContractCode([]byte("not WASM")).
			Build()).
		WithAddress(newAddress).
		AndAssertResults(func(blockchainHook *contextmock.BlockchainHookStub, verify *test.VMOutputVerifier) {
			verify.
				ReturnCode(vmcommon.ContractInvalid)
		})
}

func TestExecution_DeployWASM_WithoutMemory(t *testing.T) {
	test.BuildInstanceCreatorTest(t).
		WithInput(test.CreateTestContractCreateInputBuilder().
			WithGasProvided(1000).
			WithContractCode(test.GetTestSCCode("memoryless", "../../")).
			Build()).
		WithAddress(newAddress).
		AndAssertResults(func(blockchainHook *contextmock.BlockchainHookStub, verify *test.VMOutputVerifier) {
			verify.
				ReturnCode(vmcommon.ContractInvalid)
		})
}

func TestExecution_DeployWASM_WrongInit(t *testing.T) {
	test.BuildInstanceCreatorTest(t).
		WithInput(test.CreateTestContractCreateInputBuilder().
			WithGasProvided(1000).
			WithContractCode(test.GetTestSCCode("init-wrong", "../../")).
			Build()).
		WithAddress(newAddress).
		AndAssertResults(func(blockchainHook *contextmock.BlockchainHookStub, verify *test.VMOutputVerifier) {
			verify.
				ReturnCode(vmcommon.ContractInvalid)
		})
}

func TestExecution_DeployWASM_WrongMethods(t *testing.T) {
	test.BuildInstanceCreatorTest(t).
		WithInput(test.CreateTestContractCreateInputBuilder().
			WithGasProvided(1000).
			WithContractCode(test.GetTestSCCode("signatures", "../../")).
			Build()).
		WithAddress(newAddress).
		AndAssertResults(func(blockchainHook *contextmock.BlockchainHookStub, verify *test.VMOutputVerifier) {
			verify.
				ReturnCode(vmcommon.ContractInvalid)
		})
}

func TestExecution_DeployWASM_Successful(t *testing.T) {
	input := test.CreateTestContractCreateInputBuilder().
		WithGasProvided(1000).
		WithContractCode(test.GetTestSCCode("init-correct", "../../")).
		WithCallValue(88).
		WithArguments([]byte{0}).
		Build()
	test.BuildInstanceCreatorTest(t).
		WithInput(input).
		WithAddress(newAddress).
		AndAssertResults(func(blockchainHook *contextmock.BlockchainHookStub, verify *test.VMOutputVerifier) {
			verify.
				Ok().
				ReturnData([]byte("init successful")).
				GasRemaining(528).
				Nonce([]byte("caller"), 24).
				Code(newAddress, input.ContractCode).
				BalanceDelta(newAddress, 88)
		})
}

func TestExecution_DeployWASM_Popcnt(t *testing.T) {
	test.BuildInstanceCreatorTest(t).
		WithInput(test.CreateTestContractCreateInputBuilder().
			WithGasProvided(1000).
			WithCallValue(88).
			WithArguments().
			WithContractCode(test.GetTestSCCode("init-simple-popcnt", "../../")).
			Build()).
		WithAddress(newAddress).
		AndAssertResults(func(blockchainHook *contextmock.BlockchainHookStub, verify *test.VMOutputVerifier) {
			verify.
				Ok().
				ReturnData([]byte{3})
		})
}

func TestExecution_DeployWASM_AtMaximumLocals(t *testing.T) {
	test.BuildInstanceCreatorTest(t).
		WithInput(test.CreateTestContractCreateInputBuilder().
			WithGasProvided(1000).
			WithCallValue(88).
			WithContractCode(makeBytecodeWithLocals(WASMLocalsLimit)).
			Build()).
		WithAddress(newAddress).
		AndAssertResults(func(blockchainHook *contextmock.BlockchainHookStub, verify *test.VMOutputVerifier) {
			verify.
				Ok()
		})
}

func TestExecution_DeployWASM_MoreThanMaximumLocals(t *testing.T) {
	test.BuildInstanceCreatorTest(t).
		WithInput(test.CreateTestContractCreateInputBuilder().
			WithGasProvided(1000).
			WithCallValue(88).
			WithContractCode(makeBytecodeWithLocals(WASMLocalsLimit + 1)).
			Build()).
		WithAddress(newAddress).
		AndAssertResults(func(blockchainHook *contextmock.BlockchainHookStub, verify *test.VMOutputVerifier) {
			verify.
				ReturnCode(vmcommon.ContractInvalid)
		})
}

func TestExecution_DeployWASM_Init_Errors(t *testing.T) {
	test.BuildInstanceCreatorTest(t).
		WithInput(test.CreateTestContractCreateInputBuilder().
			WithGasProvided(1000).
			WithCallValue(88).
			WithArguments([]byte{1}).
			WithContractCode(test.GetTestSCCode("init-correct", "../../")).
			Build()).
		WithAddress(newAddress).
		AndAssertResults(func(blockchainHook *contextmock.BlockchainHookStub, verify *test.VMOutputVerifier) {
			verify.
				ReturnCode(vmcommon.UserError)
		})
}

func TestExecution_DeployWASM_Init_InfiniteLoop_Errors(t *testing.T) {
	test.BuildInstanceCreatorTest(t).
		WithInput(test.CreateTestContractCreateInputBuilder().
			WithGasProvided(1000).
			WithCallValue(88).
			WithArguments([]byte{2}).
			WithContractCode(test.GetTestSCCode("init-correct", "../../")).
			Build()).
		WithAddress(newAddress).
		AndAssertResults(func(blockchainHook *contextmock.BlockchainHookStub, verify *test.VMOutputVerifier) {
			verify.
				ReturnCode(vmcommon.OutOfGas)
		})
}

func TestExecution_ManyDeployments(t *testing.T) {
	if testing.Short() {
		t.Skip("not a short test")
	}

	ownerNonce := uint64(23)
	numDeployments := 1000

	for i := 0; i < numDeployments; i++ {
		test.BuildInstanceCreatorTest(t).
			WithInput(test.CreateTestContractCreateInputBuilder().
				WithGasProvided(100000).
				WithCallValue(88).
				WithCallerAddr([]byte("owner")).
				WithContractCode(test.GetTestSCCode("init-simple", "../../")).
				Build()).
			WithAddress(newAddress).
			WithSetup(func(host arwen.VMHost, stubBlockchainHook *contextmock.BlockchainHookStub) {
				stubBlockchainHook.GetUserAccountCalled = func(address []byte) (vmcommon.UserAccountHandler, error) {
					return &contextmock.StubAccount{Nonce: ownerNonce}, nil
				}
				stubBlockchainHook.NewAddressCalled = func(creatorAddress []byte, nonce uint64, vmType []byte) ([]byte, error) {
					ownerNonce++
					return []byte(string(newAddress) + " " + fmt.Sprint(ownerNonce)), nil
				}
			}).
			AndAssertResults(func(blockchainHook *contextmock.BlockchainHookStub, verify *test.VMOutputVerifier) {
				verify.
					Ok()
			})
	}
}

func TestExecution_MultipleArwens_OverlappingContractInstanceData(t *testing.T) {
	code := test.GetTestSCCode("counter", "../../")

	input := test.DefaultTestContractCallInput()
	input.GasProvided = 1000000
	input.Function = get

	host1, instanceRecorder1 := test.DefaultTestArwenForCallWithInstanceRecorderMock(t, code, nil)
	_, _, _, _, runtimeContext1, _, _ := host1.GetContexts()
	runtimeContextMock := contextmock.NewRuntimeContextWrapper(&runtimeContext1)
	runtimeContextMock.CleanWasmerInstanceFunc = func() {}
	host1.SetRuntimeContext(runtimeContextMock)

	for i := 0; i < 5; i++ {
		vmOutput, err := host1.RunSmartContractCall(input)
		verify := test.NewVMOutputVerifier(t, vmOutput, err)
		verify.Ok()
	}

	var host1InstancesData = make(map[interface{}]bool)
	for _, instance := range instanceRecorder1.GetContractInstances(code) {
		host1InstancesData[instance.GetData()] = true
	}

	host2, instanceRecorder2 := test.DefaultTestArwenForCallWithInstanceRecorderMock(t, code, nil)
	_, _, _, _, runtimeContext2, _, _ := host2.GetContexts()
	runtimeContextMock = contextmock.NewRuntimeContextWrapper(&runtimeContext2)
	runtimeContextMock.CleanWasmerInstanceFunc = func() {}
	runtimeContextMock.GetSCCodeFunc = func() ([]byte, error) {
		return code, nil
	}
	host2.SetRuntimeContext(runtimeContextMock)

	for i := 0; i < maxUint8AsInt+1; i++ {
		vmOutput, err := host2.RunSmartContractCall(input)
		verify := test.NewVMOutputVerifier(t, vmOutput, err)
		verify.Ok()
	}

	for _, instance := range instanceRecorder2.GetContractInstances(code) {
		_, found := host1InstancesData[instance.GetData()]
		require.False(t, found)
	}
}

func TestExecution_MultipleArwens_CleanInstanceWhileOthersAreRunning(t *testing.T) {

	code := test.GetTestSCCode("counter", "../../")

	input := test.DefaultTestContractCallInput()
	input.GasProvided = 1000000
	input.Function = get

	interHostsChan := make(chan string)
	host1Chan := make(chan string)

	host1, _ := test.DefaultTestArwenForCall(t, code, nil)
	_, _, _, _, runtimeContext1, _, _ := host1.GetContexts()
	runtimeContextMock := contextmock.NewRuntimeContextWrapper(&runtimeContext1)
	runtimeContextMock.GetFunctionToCallFunc = func() (wasmer.ExportedFunctionCallback, error) {
		interHostsChan <- "waitForHost2"
		return runtimeContextMock.GetWrappedRuntimeContext().GetFunctionToCall()
	}
	host1.SetRuntimeContext(runtimeContextMock)

	var vmOutput1 *vmcommon.VMOutput
	var err1 error
	go func() {
		vmOutput1, err1 = host1.RunSmartContractCall(input)
		interHostsChan <- "finish"
		host1Chan <- "finish"
	}()

	host2, _ := test.DefaultTestArwenForCall(t, code, nil)
	_, _, _, _, runtimeContext2, _, _ := host2.GetContexts()
	runtimeContextMock = contextmock.NewRuntimeContextWrapper(&runtimeContext2)
	runtimeContextMock.GetFunctionToCallFunc = func() (wasmer.ExportedFunctionCallback, error) {
		// wait to make sure host1 is running also
		<-interHostsChan
		// wait for host1 to finish
		<-interHostsChan
		return runtimeContextMock.GetWrappedRuntimeContext().GetFunctionToCall()
	}
	host2.SetRuntimeContext(runtimeContextMock)

	vmOutput2, err2 := host2.RunSmartContractCall(input)

	<-host1Chan

	verify1 := test.NewVMOutputVerifier(t, vmOutput1, err1)
	verify1.Ok()

	verify2 := test.NewVMOutputVerifier(t, vmOutput2, err2)
	verify2.Ok()
}

func TestExecution_Deploy_DisallowFloatingPoint(t *testing.T) {
	test.BuildInstanceCreatorTest(t).
		WithInput(test.CreateTestContractCreateInputBuilder().
			WithGasProvided(1000).
			WithCallValue(88).
			WithArguments([]byte{2}).
			WithContractCode(test.GetTestSCCode("num-with-fp", "../../")).
			Build()).
		WithAddress(newAddress).
		AndAssertResults(func(blockchainHook *contextmock.BlockchainHookStub, verify *test.VMOutputVerifier) {
			verify.
				ReturnCode(vmcommon.ContractInvalid)
		})
}

func TestExecution_CallGetUserAccountErr(t *testing.T) {
	errGetAccount := errors.New("get code error")
	test.BuildInstanceCallTest(t).
		WithInput(test.CreateTestContractCallInputBuilder().
			WithGasProvided(100).
			Build()).
		WithSetup(func(host arwen.VMHost, stubBlockchainHook *contextmock.BlockchainHookStub) {
			stubBlockchainHook.GetUserAccountCalled = func(address []byte) (vmcommon.UserAccountHandler, error) {
				return nil, errGetAccount
			}
		}).
		AndAssertResults(func(host arwen.VMHost, stubBlockchainHook *contextmock.BlockchainHookStub, verify *test.VMOutputVerifier) {
			verify.
				ReturnCode(vmcommon.ContractNotFound).
				ReturnMessage(arwen.ErrContractNotFound.Error())
		})
}

func TestExecution_NotEnoughGasForGetCode(t *testing.T) {
	test.BuildInstanceCallTest(t).
		WithInput(test.CreateTestContractCallInputBuilder().
			WithGasProvided(0).
			Build()).
		AndAssertResults(func(host arwen.VMHost, stubBlockchainHook *contextmock.BlockchainHookStub, verify *test.VMOutputVerifier) {
			verify.
				ReturnCode(vmcommon.OutOfGas).
				ReturnMessage(arwen.ErrNotEnoughGas.Error())
		})
}

func TestExecution_CallOutOfGas(t *testing.T) {
	test.BuildInstanceCallTest(t).
		WithContracts(
			test.CreateInstanceContract(test.ParentAddress).
				WithCode(test.GetTestSCCode("counter", "../../"))).
		WithInput(test.CreateTestContractCallInputBuilder().
			WithGasProvided(0).
			Build()).
		AndAssertResults(func(host arwen.VMHost, stubBlockchainHook *contextmock.BlockchainHookStub, verify *test.VMOutputVerifier) {
			verify.
				ReturnCode(vmcommon.OutOfGas).
				ReturnMessage(arwen.ErrNotEnoughGas.Error())
		})
}

func TestExecution_CallWasmerError(t *testing.T) {
	test.BuildInstanceCallTest(t).
		WithContracts(
			test.CreateInstanceContract(test.ParentAddress).
				WithCode([]byte("not WASM"))).
		WithInput(test.CreateTestContractCallInputBuilder().
			WithGasProvided(100000).
			WithFunction(increment).
			Build()).
		AndAssertResults(func(host arwen.VMHost, stubBlockchainHook *contextmock.BlockchainHookStub, verify *test.VMOutputVerifier) {
			verify.
				ReturnCode(vmcommon.ContractInvalid)
		})
}

func TestExecution_CallSCMethod_Init(t *testing.T) {
	test.BuildInstanceCallTest(t).
		WithContracts(
			test.CreateInstanceContract(test.ParentAddress).
				WithCode(test.GetTestSCCode("counter", "../../"))).
		WithInput(test.CreateTestContractCallInputBuilder().
			WithGasProvided(100000).
			WithFunction("init").
			Build()).
		AndAssertResults(func(host arwen.VMHost, stubBlockchainHook *contextmock.BlockchainHookStub, verify *test.VMOutputVerifier) {
			verify.
				ReturnCode(vmcommon.UserError).
				ReturnMessage(arwen.ErrInitFuncCalledInRun.Error())
		})
}

func TestExecution_CallSCMethod_Callback(t *testing.T) {
	test.BuildInstanceCallTest(t).
		WithContracts(
			test.CreateInstanceContract(test.ParentAddress).
				WithCode(test.GetTestSCCode("counter", "../../"))).
		WithInput(test.CreateTestContractCallInputBuilder().
			WithGasProvided(100000).
			WithFunction("callBack").
			Build()).
		AndAssertResults(func(host arwen.VMHost, stubBlockchainHook *contextmock.BlockchainHookStub, verify *test.VMOutputVerifier) {
			verify.
				ReturnCode(vmcommon.UserError).
				ReturnMessage(arwen.ErrCallBackFuncCalledInRun.Error())
		})
}

func TestExecution_CallSCMethod_MissingFunction(t *testing.T) {
	test.BuildInstanceCallTest(t).
		WithContracts(
			test.CreateInstanceContract(test.ParentAddress).
				WithCode(test.GetTestSCCode("counter", "../../"))).
		WithInput(test.CreateTestContractCallInputBuilder().
			WithGasProvided(100000).
			WithFunction("wrong").
			Build()).
		AndAssertResults(func(host arwen.VMHost, stubBlockchainHook *contextmock.BlockchainHookStub, verify *test.VMOutputVerifier) {
			verify.
				ReturnCode(vmcommon.FunctionNotFound)
		})
}

func TestExecution_Call_Successful(t *testing.T) {
	test.BuildInstanceCallTest(t).
		WithContracts(
			test.CreateInstanceContract(test.ParentAddress).
				WithCode(test.GetTestSCCode("counter", "../../"))).
		WithInput(test.CreateTestContractCallInputBuilder().
			WithGasProvided(100000).
			WithFunction(increment).
			Build()).
		WithSetup(func(host arwen.VMHost, stubBlockchainHook *contextmock.BlockchainHookStub) {
			stubBlockchainHook.GetStorageDataCalled = func(scAddress []byte, key []byte) ([]byte, error) {
				return big.NewInt(1001).Bytes(), nil
			}
		}).
		AndAssertResults(func(host arwen.VMHost, stubBlockchainHook *contextmock.BlockchainHookStub, verify *test.VMOutputVerifier) {
			verify.
				Ok().
				Storage(
					test.CreateStoreEntry(test.ParentAddress).WithKey(counterKey).WithValue(big.NewInt(1002).Bytes()),
				)
		})
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
	var gasSchedule *config.GasCost
	var gasUsed uint64

	gasLimit := uint64(100000)
	code := makeBytecodeWithLocals(locals)

	test.BuildInstanceCallTest(t).
		WithContracts(
			test.CreateInstanceContract(test.ParentAddress).
				WithCode(code)).
		WithInput(test.CreateTestContractCallInputBuilder().
			WithGasProvided(gasLimit).
			WithFunction("answer").
			Build()).
		WithSetup(func(host arwen.VMHost, stubBlockchainHook *contextmock.BlockchainHookStub) {
			gasSchedule = host.Metering().GasSchedule()
		}).
		AndAssertResults(func(host arwen.VMHost, stubBlockchainHook *contextmock.BlockchainHookStub, verify *test.VMOutputVerifier) {
			compilationCost := uint64(len(code)) * gasSchedule.BaseOperationCost.CompilePerByte
			gasUsed = gasLimit - verify.VmOutput.GasRemaining - compilationCost
			verify.
				Ok()
		})

	return gasUsed, gasSchedule
}

func TestExecution_ExecuteOnSameContext_Simple(t *testing.T) {
	parentGasUsed := uint64(521)
	childGasUsed := uint64(6870)
	executionCost := parentGasUsed + childGasUsed

	var returnData [][]byte

	returnData = append(returnData, []byte("child"))
	returnData = append(returnData, []byte{})
	for i := 1; i < 100; i++ {
		returnData = append(returnData, []byte{byte(i)})
	}
	returnData = append(returnData, []byte{})
	returnData = append(returnData, []byte("child"))
	returnData = append(returnData, []byte{})
	for i := 1; i < 100; i++ {
		returnData = append(returnData, []byte{byte(i)})
	}
	returnData = append(returnData, []byte{})
	returnData = append(returnData, []byte("parent"))

	test.BuildInstanceCallTest(t).
		WithContracts(
			test.CreateInstanceContract(test.ParentAddress).
				WithCode(test.GetTestSCCode("exec-same-ctx-simple-parent", "../../")).
				WithBalance(1000),
			test.CreateInstanceContract(test.ChildAddress).
				WithCode(test.GetTestSCCode("exec-same-ctx-simple-child", "../../")).
				WithBalance(1000)).
		WithInput(test.CreateTestContractCallInputBuilder().
			WithRecipientAddr(test.ParentAddress).
			WithFunction(parentFunctionChildCall).
			WithGasProvided(test.GasProvided).
			Build()).
		AndAssertResults(func(host arwen.VMHost, stubBlockchainHook *contextmock.BlockchainHookStub, verify *test.VMOutputVerifier) {
			verify.
				Ok().
				// test.ParentAddress
				BalanceDelta(test.ParentAddress, -198).
				GasUsed(test.ParentAddress, parentGasUsed).
				// test.ChildAddress
				BalanceDelta(test.ChildAddress, 198).
				GasUsed(test.ChildAddress, childGasUsed).
				// other
				GasRemaining(test.GasProvided - executionCost).
				ReturnData(returnData...)
		})
}

func TestExecution_Call_Breakpoints(t *testing.T) {
	t.Parallel()

	test.BuildInstanceCallTest(t).
		WithContracts(
			test.CreateInstanceContract(test.ParentAddress).
				WithCode(test.GetTestSCCode("breakpoint", "../../"))).
		WithInput(test.CreateTestContractCallInputBuilder().
			WithGasProvided(100000).
			WithFunction("testFunc").
			WithArguments([]byte{15}).
			Build()).
		AndAssertResults(func(host arwen.VMHost, stubBlockchainHook *contextmock.BlockchainHookStub, verify *test.VMOutputVerifier) {
			verify.
				Ok().
				ReturnData([]byte{100})
		})
}

func TestExecution_Call_Breakpoints_UserError(t *testing.T) {
	t.Parallel()
	test.BuildInstanceCallTest(t).
		WithContracts(
			test.CreateInstanceContract(test.ParentAddress).
				WithCode(test.GetTestSCCode("breakpoint", "../../"))).
		WithInput(test.CreateTestContractCallInputBuilder().
			WithGasProvided(100000).
			WithFunction("testFunc").
			WithArguments([]byte{1}).
			Build()).
		AndAssertResults(func(host arwen.VMHost, stubBlockchainHook *contextmock.BlockchainHookStub, verify *test.VMOutputVerifier) {
			verify.
				ReturnData().
				ReturnCode(vmcommon.UserError).
				ReturnMessage("exit here")
		})
}

func TestExecution_ExecuteOnSameContext_Prepare(t *testing.T) {
	expectedExecutionCost := uint64(138)
	test.BuildInstanceCallTest(t).
		WithContracts(
			test.CreateInstanceContract(test.ParentAddress).
				WithCode(test.GetTestSCCode("exec-same-ctx-parent", "../../")).
				WithBalance(1000)).
		WithInput(test.CreateTestContractCallInputBuilder().
			WithRecipientAddr(test.ParentAddress).
			WithFunction("parentFunctionPrepare").
			WithGasProvided(test.GasProvided).
			Build()).
		AndAssertResults(func(host arwen.VMHost, stubBlockchainHook *contextmock.BlockchainHookStub, verify *test.VMOutputVerifier) {
			verify.
				Ok().
				GasUsed(test.ParentAddress, 3405).
				Balance(test.ParentAddress, 1000).
				BalanceDelta(test.ParentAddress, -test.ParentTransferValue).
				BalanceDelta(test.ParentTransferReceiver, test.ParentTransferValue).
				GasRemaining(test.GasProvided-
					test.ParentCompilationCostSameCtx-
					expectedExecutionCost).
				ReturnData(test.ParentFinishA, test.ParentFinishB, []byte("succ")).
				Storage(
					test.CreateStoreEntry(test.ParentAddress).WithKey(test.ParentKeyA).WithValue(test.ParentDataA),
					test.CreateStoreEntry(test.ParentAddress).WithKey(test.ParentKeyB).WithValue(test.ParentDataB),
				).
				Transfers(
					test.CreateTransferEntry(test.ParentAddress, test.ParentTransferReceiver).
						WithData(test.ParentTransferData).
						WithValue(big.NewInt(test.ParentTransferValue)),
				)
		})
}

func TestExecution_ExecuteOnSameContext_Wrong(t *testing.T) {
	executionCostBeforeExecuteAPI := uint64(156)
	executeAPICost := uint64(39)
	gasLostOnFailure := uint64(50000)
	finalCost := uint64(44)

	test.BuildInstanceCallTest(t).
		WithContracts(
			test.CreateInstanceContract(test.ParentAddress).
				WithCode(test.GetTestSCCode("exec-same-ctx-parent", "../../")).
				WithBalance(1000)).
		WithInput(test.CreateTestContractCallInputBuilder().
			WithRecipientAddr(test.ParentAddress).
			WithFunction("parentFunctionWrongCall").
			WithGasProvided(test.GasProvided).
			Build()).
		AndAssertResults(func(host arwen.VMHost, stubBlockchainHook *contextmock.BlockchainHookStub, verify *test.VMOutputVerifier) {
			if !host.Runtime().ElrondSyncExecAPIErrorShouldFailExecution() {
				verify.
					Ok().
					GasUsed(test.ParentAddress, 3405).
					Balance(test.ParentAddress, 1000).
					BalanceDelta(test.ParentAddress, -test.ParentTransferValue).
					BalanceDelta(test.ParentTransferReceiver, test.ParentTransferValue).
					GasRemaining(test.GasProvided-
						test.ParentCompilationCostSameCtx-
						executionCostBeforeExecuteAPI-
						executeAPICost-
						gasLostOnFailure-
						finalCost).
					ReturnData(test.ParentFinishA, test.ParentFinishB, []byte("succ"), []byte("fail")).
					Storage(
						test.CreateStoreEntry(test.ParentAddress).WithKey(test.ParentKeyA).WithValue(test.ParentDataA),
						test.CreateStoreEntry(test.ParentAddress).WithKey(test.ParentKeyB).WithValue(test.ParentDataB),
						test.CreateStoreEntry(test.ChildAddress).WithKey(test.ChildKey).WithValue(test.ChildData),
					).
					Transfers(
						test.CreateTransferEntry(test.ParentAddress, test.ParentTransferReceiver).
							WithData(test.ParentTransferData).
							WithValue(big.NewInt(test.ParentTransferValue)),
					)
			} else {
				verify.
					ReturnCode(vmcommon.ExecutionFailed).
					ReturnMessage("account not found").
					GasRemaining(0)
			}
		})
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

	// Call parentFunctionChildCall_OutOfGas() of the parent SC, which will call
	// the child SC using executeOnSameContext() with sufficient gas for
	// compilation and starting, but the child starts an infinite loop which will
	// end in OutOfGas.

	executionCostBeforeExecuteAPI := uint64(90)
	executeAPICost := uint64(1)
	gasLostOnFailure := uint64(3500)
	finalCost := uint64(54)

	test.BuildInstanceCallTest(t).
		WithContracts(
			test.CreateInstanceContract(test.ParentAddress).
				WithCode(test.GetTestSCCode("exec-same-ctx-parent", "../../")).
				WithBalance(1000),
			test.CreateInstanceContract(test.ChildAddress).
				WithCode(test.GetTestSCCode("exec-same-ctx-child", "../../")).
				WithBalance(1000)).
		WithInput(test.CreateTestContractCallInputBuilder().
			WithRecipientAddr(test.ParentAddress).
			WithFunction("parentFunctionChildCall_OutOfGas").
			WithGasProvided(test.GasProvided).
			Build()).
		AndAssertResults(func(host arwen.VMHost, stubBlockchainHook *contextmock.BlockchainHookStub, verify *test.VMOutputVerifier) {
			if !host.Runtime().ElrondSyncExecAPIErrorShouldFailExecution() {
				verify.
					Ok().
					Balance(test.ParentAddress, 1000).
					BalanceDelta(test.ParentAddress, 0).
					GasRemaining(test.GasProvided-
						test.ParentCompilationCostSameCtx-
						executionCostBeforeExecuteAPI-
						executeAPICost-
						gasLostOnFailure-
						finalCost).
					ReturnData(test.ParentFinishA, []byte("fail")).
					Storage(
						test.CreateStoreEntry(test.ParentAddress).WithKey(test.ParentKeyA).WithValue(test.ParentDataA),
						test.CreateStoreEntry(test.ParentAddress).WithKey(test.ParentKeyB).WithValue(test.ParentDataB),
					)
			} else {
				verify.
					ReturnCode(vmcommon.ExecutionFailed).
					ReturnMessage(arwen.ErrNotEnoughGas.Error()).
					HasRuntimeErrors(arwen.ErrNotEnoughGas.Error(), arwen.ErrExecutionFailed.Error()).
					GasRemaining(0)
			}
		})
}

func TestExecution_ExecuteOnSameContext_Successful(t *testing.T) {
	executeAPICost := uint64(39)
	childExecutionCost := uint64(437)
	parentGasBeforeExecuteAPI := uint64(172)
	finalCost := uint64(134)

	parentAccountBalance := int64(1000)

	returnData := [][]byte{test.ParentFinishA, test.ParentFinishB, []byte("succ")}
	returnData = append(returnData, test.ChildFinish, test.ParentDataA)
	for _, c := range test.ParentDataA {
		returnData = append(returnData, []byte{c})
	}
	returnData = append(returnData, test.ParentDataB)
	for _, c := range test.ParentDataB {
		returnData = append(returnData, []byte{c})
	}
	returnData = append(returnData, []byte("child ok"), []byte("succ"), []byte("succ"))

	// Call parentFunctionChildCall() of the parent SC, which will call the child
	// SC and pass some arguments using executeOnSameContext().

	test.BuildInstanceCallTest(t).
		WithContracts(
			test.CreateInstanceContract(test.ParentAddress).
				WithCode(test.GetTestSCCode("exec-same-ctx-parent", "../../")).
				WithBalance(parentAccountBalance),
			test.CreateInstanceContract(test.ChildAddress).
				WithCode(test.GetTestSCCode("exec-same-ctx-child", "../../")).
				WithBalance(1000)).
		WithInput(test.CreateTestContractCallInputBuilder().
			WithRecipientAddr(test.ParentAddress).
			WithFunction(parentFunctionChildCall).
			WithGasProvided(test.GasProvided).
			Build()).
		AndAssertResults(func(host arwen.VMHost, stubBlockchainHook *contextmock.BlockchainHookStub, verify *test.VMOutputVerifier) {
			verify.
				Ok().
				// test.ParentAddress
				Balance(test.ParentAddress, parentAccountBalance).
				BalanceDelta(test.ParentAddress, -141).
				GasUsed(test.ParentAddress, 3612).
				// test.ChildAddress
				Balance(test.ChildAddress, 1000).
				BalanceDelta(test.ChildAddress, 3).
				GasUsed(test.ChildAddress, test.ChildCompilationCostSameCtx+childExecutionCost).
				// others
				BalanceDelta(test.ChildTransferReceiver, 96).
				BalanceDelta(test.ParentTransferReceiver, test.ParentTransferValue).
				GasRemaining(test.GasProvided-
					test.ParentCompilationCostSameCtx-
					parentGasBeforeExecuteAPI-
					executeAPICost-
					test.ChildCompilationCostSameCtx-
					childExecutionCost-
					finalCost).
				ReturnData(returnData...).
				Storage(
					test.CreateStoreEntry(test.ParentAddress).WithKey(test.ParentKeyA).WithValue(test.ParentDataA),
					test.CreateStoreEntry(test.ParentAddress).WithKey(test.ParentKeyB).WithValue(test.ParentDataB),
					test.CreateStoreEntry(test.ParentAddress).WithKey(test.ChildKey).WithValue(test.ChildData),
				).
				Transfers(
					test.CreateTransferEntry(test.ParentAddress, test.ParentTransferReceiver).
						WithData(test.ParentTransferData).
						WithValue(big.NewInt(test.ParentTransferValue)),
					test.CreateTransferEntry(test.ChildAddress, test.ChildTransferReceiver).
						WithData([]byte("qwerty")).
						WithValue(big.NewInt(96)),
				)
		})
}

func TestExecution_ExecuteOnSameContext_Successful_BigInts(t *testing.T) {
	// Call parentFunctionChildCall_BigInts() of the parent SC, which will call a
	// method of the child SC that takes some big Int references as arguments and
	// produce a new big Int out of the arguments.

	childExecutionCost := uint64(108)
	parentGasBeforeExecuteAPI := uint64(114)
	executeAPICost := uint64(13)
	finalCost := uint64(67)

	test.BuildInstanceCallTest(t).
		WithContracts(
			test.CreateInstanceContract(test.ParentAddress).
				WithCode(test.GetTestSCCode("exec-same-ctx-parent", "../../")).
				WithBalance(1000),
			test.CreateInstanceContract(test.ChildAddress).
				WithCode(test.GetTestSCCode("exec-same-ctx-child", "../../")).
				WithBalance(1000)).
		WithInput(test.CreateTestContractCallInputBuilder().
			WithRecipientAddr(test.ParentAddress).
			WithFunction("parentFunctionChildCall_BigInts").
			WithGasProvided(test.GasProvided).
			Build()).
		AndAssertResults(func(host arwen.VMHost, stubBlockchainHook *contextmock.BlockchainHookStub, verify *test.VMOutputVerifier) {
			verify.
				Ok().
				// test.ParentAddress
				Balance(test.ParentAddress, 1000).
				BalanceDelta(test.ParentAddress, -99).
				GasUsed(test.ParentAddress, 3461).
				// test.ChildAddress
				BalanceDelta(test.ChildAddress, 99).
				GasUsed(test.ChildAddress, test.ChildCompilationCostSameCtx+childExecutionCost).
				// others
				GasRemaining(test.GasProvided-
					test.ParentCompilationCostSameCtx-
					parentGasBeforeExecuteAPI-
					executeAPICost-
					test.ChildCompilationCostSameCtx-
					childExecutionCost-
					finalCost).
				ReturnData([]byte("child ok"), []byte("succ"), []byte("succ"))
		})
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

	recursiveCalls := 5
	var returnData [][]byte

	for i := recursiveCalls; i >= 0; i-- {
		finishString := fmt.Sprintf("Rfinish%03d", i)
		returnData = append(returnData, []byte(finishString))
	}
	for i := recursiveCalls - 1; i >= 0; i-- {
		returnData = append(returnData, []byte("succ"))
	}

	var storeEntries []test.StoreEntry

	for i := 0; i <= recursiveCalls; i++ {
		key := fmt.Sprintf("Rkey%03d.........................", i)
		value := fmt.Sprintf("Rvalue%03d", i)
		storeEntries = append(storeEntries, test.CreateStoreEntry(test.ParentAddress).WithKey([]byte(key)).WithValue([]byte(value)))
	}

	storeEntries = append(storeEntries, test.CreateStoreEntry(test.ParentAddress).WithKey(test.RecursiveIterationCounterKey).WithValue([]byte{byte(recursiveCalls + 1)}))
	storeEntries = append(storeEntries, test.CreateStoreEntry(test.ParentAddress).WithKey(test.RecursiveIterationBigCounterKey).WithValue(big.NewInt(int64(recursiveCalls+1)).Bytes()))

	test.BuildInstanceCallTest(t).
		WithContracts(
			test.CreateInstanceContract(test.ParentAddress).
				WithCode(test.GetTestSCCode("exec-same-ctx-recursive", "../../")).
				WithBalance(1000)).
		WithInput(test.CreateTestContractCallInputBuilder().
			WithRecipientAddr(test.ParentAddress).
			WithFunction(callRecursive).
			WithGasProvided(test.GasProvided).
			WithArguments([]byte{byte(recursiveCalls)}).
			Build()).
		AndAssertResults(func(host arwen.VMHost, stubBlockchainHook *contextmock.BlockchainHookStub, verify *test.VMOutputVerifier) {
			verify.
				Ok().
				Balance(test.ParentAddress, 1000).
				BalanceDelta(test.ParentAddress, 0).
				GasUsed(test.ParentAddress, 25424).
				ReturnData(returnData...).
				Storage(storeEntries...)

			require.Equal(t, int64(recursiveCalls+1), host.BigInt().GetOne(16).Int64())
		})
}

func TestExecution_ExecuteOnSameContext_Recursive_Direct_ErrMaxInstances(t *testing.T) {
	recursiveCalls := byte(11)
	test.BuildInstanceCallTest(t).
		WithContracts(
			test.CreateInstanceContract(test.ParentAddress).
				WithCode(test.GetTestSCCode("exec-same-ctx-recursive", "../../")).
				WithBalance(1000)).
		WithInput(test.CreateTestContractCallInputBuilder().
			WithRecipientAddr(test.ParentAddress).
			WithFunction(callRecursive).
			WithGasProvided(test.GasProvided).
			WithArguments([]byte{byte(recursiveCalls)}).
			Build()).
		AndAssertResults(func(host arwen.VMHost, stubBlockchainHook *contextmock.BlockchainHookStub, verify *test.VMOutputVerifier) {
			if host.Runtime().ElrondSyncExecAPIErrorShouldFailExecution() == false {
				verify.
					Ok().
					Balance(test.ParentAddress, 1000).
					BalanceDelta(test.ParentAddress, 0).
					ReturnData(
						[]byte(fmt.Sprintf("Rfinish%03d", recursiveCalls)),
						[]byte("fail"),
					).
					Storage(
						test.CreateStoreEntry(test.ParentAddress).
							WithKey([]byte(fmt.Sprintf("Rkey%03d.........................", recursiveCalls))).
							WithValue([]byte(fmt.Sprintf("Rvalue%03d", recursiveCalls))),
					)
				require.Equal(t, int64(1), host.BigInt().GetOne(16).Int64())
			} else {
				verify.
					ReturnCode(vmcommon.ExecutionFailed).
					ReturnMessage(arwen.ErrExecutionFailed.Error()).
					HasRuntimeErrors(arwen.ErrMaxInstancesReached.Error(), arwen.ErrExecutionFailed.Error()).
					GasRemaining(0)
			}
		})
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

	recursiveCalls := 5

	var returnData [][]byte
	var storeEntries []test.StoreEntry

	storeEntries = append(storeEntries, test.CreateStoreEntry(test.ParentAddress).WithKey(test.RecursiveIterationCounterKey).WithValue([]byte{byte(recursiveCalls + 1)}))
	storeEntries = append(storeEntries, test.CreateStoreEntry(test.ParentAddress).WithKey(test.RecursiveIterationBigCounterKey).WithValue(big.NewInt(int64(recursiveCalls+1)).Bytes()))

	returnData = append(returnData, []byte("start recursive mutual calls"))

	for i := 0; i <= recursiveCalls; i++ {
		var finishData string
		var key string
		var value string
		iteration := recursiveCalls - i
		if i%2 == 0 {
			finishData = fmt.Sprintf("Afinish%03d", iteration)
			key = fmt.Sprintf("Akey%03d.........................", iteration)
			value = fmt.Sprintf("Avalue%03d", iteration)
		} else {
			finishData = fmt.Sprintf("Bfinish%03d", iteration)
			key = fmt.Sprintf("Bkey%03d.........................", iteration)
			value = fmt.Sprintf("Bvalue%03d", iteration)
		}
		storeEntries = append(storeEntries, test.CreateStoreEntry(test.ParentAddress).WithKey([]byte(key)).WithValue([]byte(value)))
		returnData = append(returnData, []byte(finishData))
	}

	for i := recursiveCalls; i >= 0; i-- {
		returnData = append(returnData, []byte("succ"))
	}

	returnData = append(returnData, []byte("end recursive mutual calls"))

	test.BuildInstanceCallTest(t).
		WithContracts(
			test.CreateInstanceContract(test.ParentAddress).
				WithCode(test.GetTestSCCode("exec-same-ctx-recursive", "../../")).
				WithBalance(1000)).
		WithInput(test.CreateTestContractCallInputBuilder().
			WithRecipientAddr(test.ParentAddress).
			WithFunction("callRecursiveMutualMethods").
			WithGasProvided(test.GasProvided).
			WithArguments([]byte{byte(recursiveCalls)}).
			Build()).
		AndAssertResults(func(host arwen.VMHost, stubBlockchainHook *contextmock.BlockchainHookStub, verify *test.VMOutputVerifier) {
			verify.
				Ok().
				Balance(test.ParentAddress, 1000).
				BalanceDelta(test.ParentAddress, (big.NewInt(0).Sub(big.NewInt(1), big.NewInt(1))).Int64()).
				GasUsed(test.ParentAddress, 29593).
				ReturnData(returnData...).
				Storage(storeEntries...)

			require.Equal(t, int64(recursiveCalls+1), host.BigInt().GetOne(16).Int64())
		})
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

	// Call parentFunctionChildCall() of the parent SC, which will call the child
	// SC and pass some arguments using executeOnDestContext().

	recursiveCalls := 4

	var expectedParentBalanceDelta, expectedChildBalanceDelta int64
	if recursiveCalls%2 == 1 {
		expectedParentBalanceDelta = -5
		expectedChildBalanceDelta = 5
	} else {
		expectedParentBalanceDelta = big.NewInt(0).Sub(big.NewInt(1), big.NewInt(1)).Int64()
		expectedChildBalanceDelta = big.NewInt(0).Sub(big.NewInt(1), big.NewInt(1)).Int64()
	}

	var returnData [][]byte
	var storeEntries []test.StoreEntry

	for i := 0; i <= recursiveCalls; i++ {
		var finishData string
		var key string
		var value string
		iteration := recursiveCalls - i
		if i%2 == 0 {
			finishData = fmt.Sprintf("Pfinish%03d", iteration)
			key = fmt.Sprintf("Pkey%03d.........................", iteration)
			value = fmt.Sprintf("Pvalue%03d", iteration)
		} else {
			finishData = fmt.Sprintf("Cfinish%03d", iteration)
			key = fmt.Sprintf("Ckey%03d.........................", iteration)
			value = fmt.Sprintf("Cvalue%03d", iteration)
		}
		storeEntries = append(storeEntries, test.CreateStoreEntry(test.ParentAddress).WithKey([]byte(key)).WithValue([]byte(value)))
		returnData = append(returnData, []byte(finishData))
	}

	for i := recursiveCalls - 1; i >= 0; i-- {
		returnData = append(returnData, []byte("succ"))
	}

	storeEntries = append(storeEntries, test.CreateStoreEntry(test.ParentAddress).WithKey(test.RecursiveIterationCounterKey).WithValue([]byte{byte(recursiveCalls + 1)}))
	storeEntries = append(storeEntries, test.CreateStoreEntry(test.ParentAddress).WithKey(test.RecursiveIterationBigCounterKey).WithValue(big.NewInt(int64(recursiveCalls+1)).Bytes()))

	test.BuildInstanceCallTest(t).
		WithContracts(
			test.CreateInstanceContract(test.ParentAddress).
				WithCode(test.GetTestSCCode("exec-same-ctx-recursive-parent", "../../")).
				WithBalance(1000),
			test.CreateInstanceContract(test.ChildAddress).
				WithCode(test.GetTestSCCode("exec-same-ctx-recursive-child", "../../")).
				WithBalance(1000)).
		WithInput(test.CreateTestContractCallInputBuilder().
			WithRecipientAddr(test.ParentAddress).
			WithFunction(parentCallsChild).
			WithGasProvided(test.GasProvided).
			WithArguments([]byte{byte(recursiveCalls)}).
			Build()).
		AndAssertResults(func(host arwen.VMHost, stubBlockchainHook *contextmock.BlockchainHookStub, verify *test.VMOutputVerifier) {
			verify.
				Ok().
				// test.ParentAddress
				Balance(test.ParentAddress, 1000).
				BalanceDelta(test.ParentAddress, expectedParentBalanceDelta).
				GasUsed(test.ParentAddress, 5426).
				// test.ChildAddress
				Balance(test.ChildAddress, 1000).
				BalanceDelta(test.ChildAddress, expectedChildBalanceDelta).
				GasUsed(test.ChildAddress, 3652).
				// other
				ReturnData(returnData...).
				Storage(storeEntries...)

			require.Equal(t, int64(recursiveCalls+1), host.BigInt().GetOne(88).Int64())
		})
}

func TestExecution_ExecuteOnSameContext_Recursive_Mutual_SCs_OutOfGas(t *testing.T) {
	// Call parentFunctionChildCall() of the parent SC, which will call the child
	// SC and pass some arguments using executeOnDestContext().
	recursiveCalls := byte(5)

	test.BuildInstanceCallTest(t).
		WithContracts(
			test.CreateInstanceContract(test.ParentAddress).
				WithCode(test.GetTestSCCode("exec-same-ctx-recursive-parent", "../../")).
				WithBalance(1000),
			test.CreateInstanceContract(test.ChildAddress).
				WithCode(test.GetTestSCCode("exec-same-ctx-recursive-child", "../../")).
				WithBalance(1000)).
		WithInput(test.CreateTestContractCallInputBuilder().
			WithRecipientAddr(test.ParentAddress).
			WithFunction(parentCallsChild).
			WithGasProvided(10000).
			WithArguments([]byte{byte(recursiveCalls)}).
			Build()).
		AndAssertResults(func(host arwen.VMHost, stubBlockchainHook *contextmock.BlockchainHookStub, verify *test.VMOutputVerifier) {
			if host.Runtime().ElrondSyncExecAPIErrorShouldFailExecution() == false {
				verify.
					ReturnCode(vmcommon.OutOfGas).
					ReturnMessage(arwen.ErrNotEnoughGas.Error()).
					GasRemaining(0)
			} else {
				verify.
					ReturnCode(vmcommon.ExecutionFailed).
					ReturnMessage(arwen.ErrExecutionFailed.Error()).
					HasRuntimeErrors(arwen.ErrNotEnoughGas.Error(), arwen.ErrExecutionFailed.Error()).
					GasRemaining(0)
			}
		})
}

func TestExecution_ExecuteOnDestContext_Prepare(t *testing.T) {
	// Execute the parent SC method "parentFunctionPrepare", which sets storage,
	// finish data and performs a transfer. This step validates the test to the
	// actual call to ExecuteOnSameContext().

	expectedExecutionCost := uint64(138)

	test.BuildInstanceCallTest(t).
		WithContracts(
			test.CreateInstanceContract(test.ParentAddress).
				WithCode(test.GetTestSCCode("exec-dest-ctx-parent", "../../")).
				WithBalance(1000)).
		WithInput(test.CreateTestContractCallInputBuilder().
			WithRecipientAddr(test.ParentAddress).
			WithFunction("parentFunctionPrepare").
			WithGasProvided(test.GasProvided).
			Build()).
		AndAssertResults(func(host arwen.VMHost, stubBlockchainHook *contextmock.BlockchainHookStub, verify *test.VMOutputVerifier) {
			verify.
				Ok().
				Balance(test.ParentAddress, 1000).
				BalanceDelta(test.ParentAddress, -test.ParentTransferValue).
				GasUsed(test.ParentAddress, 4309).
				BalanceDelta(test.ParentTransferReceiver, test.ParentTransferValue).
				GasRemaining(test.GasProvided-
					test.ParentCompilationCostDestCtx-
					expectedExecutionCost).
				ReturnData(test.ParentFinishA, test.ParentFinishB, []byte("succ")).
				Storage(
					test.CreateStoreEntry(test.ParentAddress).WithKey(test.ParentKeyA).WithValue(test.ParentDataA),
					test.CreateStoreEntry(test.ParentAddress).WithKey(test.ParentKeyB).WithValue(test.ParentDataB),
				).
				Transfers(
					test.CreateTransferEntry(test.ParentAddress, test.ParentTransferReceiver).
						WithData(test.ParentTransferData).
						WithValue(big.NewInt(test.ParentTransferValue)),
				)
		})
}

func TestExecution_ExecuteOnDestContext_Wrong(t *testing.T) {
	// Call parentFunctionWrongCall() of the parent SC, which will try to call a
	// non-existing SC.

	executionCostBeforeExecuteAPI := uint64(156)
	executeAPICost := uint64(42)
	gasLostOnFailure := uint64(10000)
	finalCost := uint64(44)

	test.BuildInstanceCallTest(t).
		WithContracts(
			test.CreateInstanceContract(test.ParentAddress).
				WithCode(test.GetTestSCCode("exec-dest-ctx-parent", "../../")).
				WithBalance(1000)).
		WithInput(test.CreateTestContractCallInputBuilder().
			WithRecipientAddr(test.ParentAddress).
			WithFunction("parentFunctionWrongCall").
			WithGasProvided(test.GasProvided).
			Build()).
		AndAssertResults(func(host arwen.VMHost, stubBlockchainHook *contextmock.BlockchainHookStub, verify *test.VMOutputVerifier) {
			if host.Runtime().ElrondSyncExecAPIErrorShouldFailExecution() == false {
				verify.
					Ok().
					Balance(test.ParentAddress, 1000).
					BalanceDelta(test.ParentAddress, -42).
					GasUsed(test.ParentAddress, 3612).
					BalanceDelta(test.ChildTransferReceiver, 96).
					BalanceDelta(test.ParentTransferReceiver, test.ParentTransferValue).
					GasRemaining(test.GasProvided-
						test.ParentCompilationCostDestCtx-
						executionCostBeforeExecuteAPI-
						executeAPICost-
						gasLostOnFailure-
						finalCost).
					ReturnData(test.ParentFinishA, test.ParentFinishB, []byte("succ"), []byte("fail")).
					Storage(
						test.CreateStoreEntry(test.ParentAddress).WithKey(test.ParentKeyA).WithValue(test.ParentDataA),
						test.CreateStoreEntry(test.ParentAddress).WithKey(test.ParentKeyB).WithValue(test.ParentDataB),
						test.CreateStoreEntry(test.ParentAddress).WithKey(test.ChildKey).WithValue(test.ChildData),
					).
					Transfers(
						test.CreateTransferEntry(test.ChildAddress, test.ChildTransferReceiver).
							WithData([]byte("qwerty")).
							WithValue(big.NewInt(96)),
						test.CreateTransferEntry(test.ParentAddress, test.ParentTransferReceiver).
							WithData(test.ParentTransferData).
							WithValue(big.NewInt(test.ParentTransferValue)),
					)
			} else {
				verify.
					ReturnCode(vmcommon.ExecutionFailed).
					ReturnMessage("account not found").
					GasRemaining(0)
			}
		})
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

	// Call parentFunctionChildCall_OutOfGas() of the parent SC, which will call
	// the child SC using executeOnDestContext() with sufficient gas for
	// compilation and starting, but the child starts an infinite loop which will
	// end in OutOfGas.

	executionCostBeforeExecuteAPI := uint64(90)
	executeAPICost := uint64(1)
	gasLostOnFailure := uint64(3500)
	finalCost := uint64(54)

	test.BuildInstanceCallTest(t).
		WithContracts(
			test.CreateInstanceContract(test.ParentAddress).
				WithCode(test.GetTestSCCode("exec-dest-ctx-parent", "../../")).
				WithBalance(1000),
			test.CreateInstanceContract(test.ChildAddress).
				WithCode(test.GetTestSCCode("exec-dest-ctx-child", "../../")).
				WithBalance(1000),
		).
		WithInput(test.CreateTestContractCallInputBuilder().
			WithRecipientAddr(test.ParentAddress).
			WithFunction("parentFunctionChildCall_OutOfGas").
			WithGasProvided(test.GasProvided).
			Build()).
		AndAssertResults(func(host arwen.VMHost, stubBlockchainHook *contextmock.BlockchainHookStub, verify *test.VMOutputVerifier) {
			if host.Runtime().ElrondSyncExecAPIErrorShouldFailExecution() == false {
				verify.
					Ok().
					Balance(test.ParentAddress, 1000).
					GasRemaining(test.GasProvided-
						test.ParentCompilationCostDestCtx-
						executionCostBeforeExecuteAPI-
						executeAPICost-
						gasLostOnFailure-
						finalCost).
					ReturnData(test.ParentFinishA, []byte("fail")).
					Storage(
						test.CreateStoreEntry(test.ParentAddress).WithKey(test.ParentKeyA).WithValue(test.ParentDataA),
						test.CreateStoreEntry(test.ParentAddress).WithKey(test.ParentKeyB).WithValue(test.ParentDataB),
					)
				require.Equal(t, int64(42), host.BigInt().GetOne(12).Int64())
			} else {
				verify.
					ReturnCode(vmcommon.ExecutionFailed).
					ReturnMessage(arwen.ErrNotEnoughGas.Error()).
					HasRuntimeErrors(arwen.ErrNotEnoughGas.Error(), arwen.ErrExecutionFailed.Error()).
					GasRemaining(0)
			}
		})
}

func TestExecution_ExecuteOnDestContext_Successful(t *testing.T) {
	// Call parentFunctionChildCall() of the parent SC, which will call the child
	// SC and pass some arguments using executeOnDestContext().

	parentGasBeforeExecuteAPI := uint64(168)
	executeAPICost := uint64(42)
	childExecutionCost := uint64(91)
	finalCost := uint64(65)
	childTransferValue := int64(12)

	test.BuildInstanceCallTest(t).
		WithContracts(
			test.CreateInstanceContract(test.ParentAddress).
				WithCode(test.GetTestSCCode("exec-dest-ctx-parent", "../../")).
				WithBalance(1000),
			test.CreateInstanceContract(test.ChildAddress).
				WithCode(test.GetTestSCCode("exec-dest-ctx-child", "../../")).
				WithBalance(1000),
		).
		WithInput(test.CreateTestContractCallInputBuilder().
			WithRecipientAddr(test.ParentAddress).
			WithFunction(parentFunctionChildCall).
			WithGasProvided(test.GasProvided).
			Build()).
		AndAssertResults(func(host arwen.VMHost, stubBlockchainHook *contextmock.BlockchainHookStub, verify *test.VMOutputVerifier) {
			verify.
				Ok().
				// test.ParentAddress
				Balance(test.ParentAddress, 1000).
				BalanceDelta(test.ParentAddress, -141).
				GasUsed(test.ParentAddress, 4444).
				/// test.ChildAddress
				Balance(test.ChildAddress, 1000).
				BalanceDelta(test.ChildAddress, 99-childTransferValue).
				GasUsed(test.ChildAddress, 2250).
				// other
				BalanceDelta(test.ChildTransferReceiver, childTransferValue).
				GasRemaining(test.GasProvided-
					test.ParentCompilationCostDestCtx-
					parentGasBeforeExecuteAPI-
					executeAPICost-
					test.ChildCompilationCostDestCtx-
					childExecutionCost-
					finalCost).
				ReturnData(test.ParentFinishA, test.ParentFinishB, []byte("succ"), test.ChildFinish, []byte("succ"), []byte("succ")).
				Storage(
					test.CreateStoreEntry(test.ParentAddress).WithKey(test.ParentKeyA).WithValue(test.ParentDataA),
					test.CreateStoreEntry(test.ParentAddress).WithKey(test.ParentKeyB).WithValue(test.ParentDataB),
					test.CreateStoreEntry(test.ParentAddress).WithKey(test.ChildKey).WithValue(nil),
					test.CreateStoreEntry(test.ChildAddress).WithKey(test.ChildKey).WithValue(test.ChildData),
				).
				Transfers(
					test.CreateTransferEntry(test.ChildAddress, test.ChildTransferReceiver).
						WithData([]byte("Second sentence.")).
						WithValue(big.NewInt(childTransferValue)),
					test.CreateTransferEntry(test.ParentAddress, test.ParentTransferReceiver).
						WithData(test.ParentTransferData).
						WithValue(big.NewInt(test.ParentTransferValue)),
				)
		})
}

func TestExecution_ExecuteOnDestContext_Successful_ChildReturns(t *testing.T) {
	// Call parentFunctionChildCall() of the parent SC, which will call the child
	// SC and pass some arguments using executeOnDestContext().

	parentGasBeforeExecuteAPI := uint64(168)
	executeAPICost := uint64(42)
	childExecutionCost := uint64(91)
	parentGasAfterExecuteAPI := uint64(273)
	childTransferValue := int64(12)

	test.BuildInstanceCallTest(t).
		WithContracts(
			test.CreateInstanceContract(test.ParentAddress).
				WithCode(test.GetTestSCCode("exec-dest-ctx-parent", "../../")).
				WithBalance(1000),
			test.CreateInstanceContract(test.ChildAddress).
				WithCode(test.GetTestSCCode("exec-dest-ctx-child", "../../")).
				WithBalance(1000),
		).
		WithInput(test.CreateTestContractCallInputBuilder().
			WithRecipientAddr(test.ParentAddress).
			WithFunction("parentFunctionChildCall_ReturnedData").
			WithGasProvided(test.GasProvided).
			Build()).
		AndAssertResults(func(host arwen.VMHost, stubBlockchainHook *contextmock.BlockchainHookStub, verify *test.VMOutputVerifier) {
			verify.
				Ok().
				// test.ParentAddress
				Balance(test.ParentAddress, 1000).
				BalanceDelta(test.ParentAddress, -141).
				GasUsed(test.ParentAddress, 4652).
				/// test.ChildAddress
				Balance(test.ChildAddress, 1000).
				BalanceDelta(test.ChildAddress, 99-childTransferValue).
				GasUsed(test.ChildAddress, 2250).
				// other
				BalanceDelta(test.ChildTransferReceiver, childTransferValue).
				GasRemaining(test.GasProvided-
					test.ParentCompilationCostDestCtx-
					parentGasBeforeExecuteAPI-
					executeAPICost-
					test.ChildCompilationCostDestCtx-
					childExecutionCost-
					parentGasAfterExecuteAPI).
				ReturnData(test.ParentFinishA, test.ParentFinishB, []byte("succ"), test.ChildFinish, []byte("succ")).
				Storage(
					test.CreateStoreEntry(test.ParentAddress).WithKey(test.ParentKeyA).WithValue(test.ParentDataA),
					test.CreateStoreEntry(test.ParentAddress).WithKey(test.ParentKeyB).WithValue(test.ParentDataB),
					test.CreateStoreEntry(test.ChildAddress).WithKey(test.ChildKey).WithValue(test.ChildData),
				).
				Transfers(
					test.CreateTransferEntry(test.ChildAddress, test.ChildTransferReceiver).
						WithData([]byte("Second sentence.")).
						WithValue(big.NewInt(childTransferValue)),
					test.CreateTransferEntry(test.ParentAddress, test.ParentTransferReceiver).
						WithData(test.ParentTransferData).
						WithValue(big.NewInt(test.ParentTransferValue)),
				)
		})
}

func TestExecution_ExecuteOnDestContext_GasRemaining(t *testing.T) {
	// This test ensures that host.ExecuteOnDestContext() calls
	// metering.GasLeft() on the Wasmer instance of the child, and not of the
	// parent.

	parentCode := test.GetTestSCCode("exec-dest-ctx-parent", "../../")
	childCode := test.GetTestSCCode("exec-dest-ctx-child", "../../")

	// Pretend that the execution of the parent SC was requested, with the
	// following ContractCallInput:
	input := test.DefaultTestContractCallInput()
	input.RecipientAddr = test.ParentAddress
	input.Function = "parentFunctionChildCall"
	input.GasProvided = test.GasProvided

	// Initialize the VM with the parent SC and child SC, but without really
	// executing the parent. The initialization emulates the behavior of
	// host.doRunSmartContractCall(). Gas cost for compilation is skipped.
	host, _ := test.DefaultTestArwenForTwoSCs(t, parentCode, childCode, nil, nil)
	host.InitState()

	_, _, metering, output, runtime, _, storage := host.GetContexts()
	runtime.InitStateFromContractCallInput(input)
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
	childInput := test.DefaultTestContractCallInput()
	childInput.CallerAddr = test.ParentAddress
	childInput.CallValue = big.NewInt(99)
	childInput.Function = "childFunction"
	childInput.RecipientAddr = test.ChildAddress
	childInput.Arguments = [][]byte{
		[]byte("some data"),
		[]byte("argument"),
		[]byte("another argument"),
	}
	childInput.GasProvided = 10000

	childOutput, err := host.ExecuteOnDestContext(childInput)
	verify := test.NewVMOutputVerifier(t, childOutput, err)
	verify.
		Ok().
		GasRemaining(7758)
}

func TestExecution_ExecuteOnDestContext_Successful_BigInts(t *testing.T) {
	// Call parentFunctionChildCall_BigInts() of the parent SC, which will call a
	// method of the child SC that takes some big Int references as arguments and
	// produce a new big Int out of the arguments.

	parentGasBeforeExecuteAPI := uint64(115)
	executeAPICost := uint64(13)
	childExecutionCost := uint64(101)
	finalCost := uint64(68)

	test.BuildInstanceCallTest(t).
		WithContracts(
			test.CreateInstanceContract(test.ParentAddress).
				WithCode(test.GetTestSCCode("exec-dest-ctx-parent", "../../")).
				WithBalance(1000),
			test.CreateInstanceContract(test.ChildAddress).
				WithCode(test.GetTestSCCode("exec-dest-ctx-child", "../../")).
				WithBalance(1000),
		).
		WithInput(test.CreateTestContractCallInputBuilder().
			WithRecipientAddr(test.ParentAddress).
			WithFunction("parentFunctionChildCall_BigInts").
			WithGasProvided(test.GasProvided).
			Build()).
		AndAssertResults(func(host arwen.VMHost, stubBlockchainHook *contextmock.BlockchainHookStub, verify *test.VMOutputVerifier) {
			verify.
				Ok().
				// test.ParentAddress
				Balance(test.ParentAddress, 1000).
				BalanceDelta(test.ParentAddress, -99).
				GasUsed(test.ParentAddress, 4366).
				/// test.ChildAddress
				BalanceDelta(test.ChildAddress, 99).
				GasUsed(test.ChildAddress, 2259).
				// other
				GasRemaining(test.GasProvided-
					test.ParentCompilationCostDestCtx-
					parentGasBeforeExecuteAPI-
					executeAPICost-
					test.ChildCompilationCostDestCtx-
					childExecutionCost-
					finalCost).
				ReturnData([]byte("child ok"), []byte("succ"), []byte("succ"))
		})
}

func TestExecution_ExecuteOnDestContext_Recursive_Direct(t *testing.T) {
	recursiveCalls := 6

	var returnData [][]byte
	var storeEntries []test.StoreEntry

	for i := recursiveCalls; i >= 0; i-- {
		finishString := fmt.Sprintf("Rfinish%03d", i)
		returnData = append(returnData, []byte(finishString))
	}

	for i := recursiveCalls - 1; i >= 0; i-- {
		returnData = append(returnData, []byte("succ"))
	}

	for i := 0; i <= recursiveCalls; i++ {
		key := fmt.Sprintf("Rkey%03d.........................", i)
		value := fmt.Sprintf("Rvalue%03d", i)
		storeEntries = append(storeEntries, test.CreateStoreEntry(test.ParentAddress).WithKey([]byte(key)).WithValue([]byte(value)))
	}

	storeEntries = append(storeEntries,
		test.CreateStoreEntry(test.ParentAddress).WithKey(test.RecursiveIterationCounterKey).WithValue([]byte{byte(recursiveCalls + 1)}),
		test.CreateStoreEntry(test.ParentAddress).WithKey(test.RecursiveIterationBigCounterKey).WithValue(big.NewInt(int64(1)).Bytes()))

	test.BuildInstanceCallTest(t).
		WithContracts(
			test.CreateInstanceContract(test.ParentAddress).
				WithCode(test.GetTestSCCode("exec-dest-ctx-recursive", "../../")).
				WithBalance(1000),
		).
		WithInput(test.CreateTestContractCallInputBuilder().
			WithRecipientAddr(test.ParentAddress).
			WithFunction(callRecursive).
			WithGasProvided(test.GasProvided).
			WithArguments([]byte{byte(recursiveCalls)}).
			Build()).
		AndAssertResults(func(host arwen.VMHost, stubBlockchainHook *contextmock.BlockchainHookStub, verify *test.VMOutputVerifier) {
			verify.
				Ok().
				Balance(test.ParentAddress, 1000).
				BalanceDelta(test.ParentAddress, big.NewInt(0).Sub(big.NewInt(1), big.NewInt(1)).Int64()).
				GasUsed(test.ParentAddress, 29670).
				ReturnData(returnData...).
				Storage(storeEntries...)

			require.Equal(t, int64(1), host.BigInt().GetOne(16).Int64())
		})
}

func TestExecution_ExecuteOnDestContext_Recursive_Mutual_Methods(t *testing.T) {
	recursiveCalls := 7

	var returnData [][]byte
	var storeEntries []test.StoreEntry

	storeEntries = append(storeEntries, test.CreateStoreEntry(test.ParentAddress).WithKey(test.RecursiveIterationCounterKey).WithValue([]byte{byte(recursiveCalls + 1)}))
	storeEntries = append(storeEntries, test.CreateStoreEntry(test.ParentAddress).WithKey(test.RecursiveIterationBigCounterKey).WithValue(big.NewInt(int64(1)).Bytes()))

	returnData = append(returnData, []byte("start recursive mutual calls"))

	for i := 0; i <= recursiveCalls; i++ {
		var finishData string
		var key string
		var value string
		iteration := recursiveCalls - i
		if i%2 == 0 {
			finishData = fmt.Sprintf("Afinish%03d", iteration)
			key = fmt.Sprintf("Akey%03d.........................", iteration)
			value = fmt.Sprintf("Avalue%03d", iteration)
		} else {
			finishData = fmt.Sprintf("Bfinish%03d", iteration)
			key = fmt.Sprintf("Bkey%03d.........................", iteration)
			value = fmt.Sprintf("Bvalue%03d", iteration)
		}
		storeEntries = append(storeEntries, test.CreateStoreEntry(test.ParentAddress).WithKey([]byte(key)).WithValue([]byte(value)))
		returnData = append(returnData, []byte(finishData))
	}

	for i := recursiveCalls; i >= 0; i-- {
		returnData = append(returnData, []byte("succ"))
	}

	returnData = append(returnData, []byte("end recursive mutual calls"))

	test.BuildInstanceCallTest(t).
		WithContracts(
			test.CreateInstanceContract(test.ParentAddress).
				WithCode(test.GetTestSCCode("exec-dest-ctx-recursive", "../../")).
				WithBalance(1000),
		).
		WithInput(test.CreateTestContractCallInputBuilder().
			WithRecipientAddr(test.ParentAddress).
			WithFunction("callRecursiveMutualMethods").
			WithGasProvided(test.GasProvided).
			WithArguments([]byte{byte(recursiveCalls)}).
			Build()).
		AndAssertResults(func(host arwen.VMHost, stubBlockchainHook *contextmock.BlockchainHookStub, verify *test.VMOutputVerifier) {
			verify.
				Ok().
				Balance(test.ParentAddress, 1000).
				BalanceDelta(test.ParentAddress, big.NewInt(0).Sub(big.NewInt(1), big.NewInt(1)).Int64()).
				GasUsed(test.ParentAddress, 38083).
				ReturnData(returnData...).
				Storage(storeEntries...)

			require.Equal(t, int64(0), host.BigInt().GetOne(16).Int64())
		})
}

func TestExecution_ExecuteOnDestContext_Recursive_Mutual_SCs(t *testing.T) {
	// Call parentFunctionChildCall() of the parent SC, which will call the child
	// SC and pass some arguments using executeOnDestContext().

	recursiveCalls := 6

	parentIterations := (recursiveCalls / 2) + (recursiveCalls % 2)
	childIterations := recursiveCalls - parentIterations
	balanceDelta := int64(5*parentIterations - 3*childIterations)

	var returnData [][]byte
	var storeEntries []test.StoreEntry

	for i := 0; i <= recursiveCalls; i++ {
		var finishData string
		var key string
		var value string
		iteration := recursiveCalls - i
		if i%2 == 0 {
			finishData = fmt.Sprintf("Pfinish%03d", iteration)
			key = fmt.Sprintf("Pkey%03d.........................", iteration)
			value = fmt.Sprintf("Pvalue%03d", iteration)
			storeEntries = append(storeEntries, test.CreateStoreEntry(test.ParentAddress).WithKey([]byte(key)).WithValue([]byte(value)))
		} else {
			finishData = fmt.Sprintf("Cfinish%03d", iteration)
			key = fmt.Sprintf("Ckey%03d.........................", iteration)
			value = fmt.Sprintf("Cvalue%03d", iteration)
			storeEntries = append(storeEntries, test.CreateStoreEntry(test.ChildAddress).WithKey([]byte(key)).WithValue([]byte(value)))
		}
		returnData = append(returnData, []byte(finishData))
	}

	for i := recursiveCalls - 1; i >= 0; i-- {
		returnData = append(returnData, []byte("succ"))
	}

	counterValue := (recursiveCalls + recursiveCalls%2) / 2

	storeEntries = append(storeEntries, test.CreateStoreEntry(test.ParentAddress).WithKey(test.RecursiveIterationCounterKey).WithValue([]byte{byte(counterValue + 1)}))
	storeEntries = append(storeEntries, test.CreateStoreEntry(test.ChildAddress).WithKey(test.RecursiveIterationCounterKey).WithValue(big.NewInt(int64(counterValue)).Bytes()))

	if recursiveCalls%2 == 0 {
		storeEntries = append(storeEntries, test.CreateStoreEntry(test.ParentAddress).WithKey(test.RecursiveIterationBigCounterKey).WithValue(big.NewInt(int64(1)).Bytes()))
	} else {
		storeEntries = append(storeEntries, test.CreateStoreEntry(test.ChildAddress).WithKey(test.RecursiveIterationBigCounterKey).WithValue(big.NewInt(int64(1)).Bytes()))
	}

	test.BuildInstanceCallTest(t).
		WithContracts(
			test.CreateInstanceContract(test.ParentAddress).
				WithCode(test.GetTestSCCode("exec-dest-ctx-recursive-parent", "../../")).
				WithBalance(1000),
			test.CreateInstanceContract(test.ChildAddress).
				WithCode(test.GetTestSCCode("exec-dest-ctx-recursive-child", "../../")).
				WithBalance(1000),
		).
		WithInput(test.CreateTestContractCallInputBuilder().
			WithRecipientAddr(test.ParentAddress).
			WithFunction(parentCallsChild).
			WithGasProvided(test.GasProvided).
			WithArguments([]byte{byte(recursiveCalls)}).
			Build()).
		AndAssertResults(func(host arwen.VMHost, stubBlockchainHook *contextmock.BlockchainHookStub, verify *test.VMOutputVerifier) {
			verify.
				Ok().
				// test.ParentAddress
				Balance(test.ParentAddress, 1000).
				BalanceDelta(test.ParentAddress, -balanceDelta).
				GasUsed(test.ParentAddress, 7252).
				// test.ChildAddress
				Balance(test.ChildAddress, 1000).
				BalanceDelta(test.ChildAddress, balanceDelta).
				GasUsed(test.ChildAddress, 5464).
				// others
				ReturnData(returnData...).
				Storage(storeEntries...)

			require.Equal(t, int64(1), host.BigInt().GetOne(88).Int64())
		})
}

func TestExecution_ExecuteOnDestContext_Recursive_Mutual_SCs_OutOfGas(t *testing.T) {
	// Call parentFunctionChildCall() of the parent SC, which will call the child
	// SC and pass some arguments using executeOnDestContext().

	recursiveCalls := byte(5)

	test.BuildInstanceCallTest(t).
		WithContracts(
			test.CreateInstanceContract(test.ParentAddress).
				WithCode(test.GetTestSCCode("exec-dest-ctx-recursive-parent", "../../")).
				WithBalance(1000),
			test.CreateInstanceContract(test.ChildAddress).
				WithCode(test.GetTestSCCode("exec-dest-ctx-recursive-child", "../../")).
				WithBalance(1000),
		).
		WithInput(test.CreateTestContractCallInputBuilder().
			WithRecipientAddr(test.ParentAddress).
			WithFunction(parentCallsChild).
			WithGasProvided(10000).
			WithArguments([]byte{byte(recursiveCalls)}).
			Build()).
		AndAssertResults(func(host arwen.VMHost, stubBlockchainHook *contextmock.BlockchainHookStub, verify *test.VMOutputVerifier) {
			if host.Runtime().ElrondSyncExecAPIErrorShouldFailExecution() == false {
				verify.
					ReturnCode(vmcommon.OutOfGas).
					ReturnMessage(arwen.ErrNotEnoughGas.Error())
			} else {
				verify.
					ReturnCode(vmcommon.ExecutionFailed).
					ReturnMessage(arwen.ErrExecutionFailed.Error()).
					HasRuntimeErrors(arwen.ErrNotEnoughGas.Error(), arwen.ErrExecutionFailed.Error()).
					GasRemaining(0)
			}
		})
}

func TestExecution_ExecuteOnSameContext_MultipleChildren(t *testing.T) {
	world := worldmock.NewMockWorld()
	host := test.DefaultTestArwen(t, world)

	alphaCode := test.GetTestSCCodeModule("exec-sync-ctx-multiple/alpha", "alpha", "../../")
	alpha := test.AddTestSmartContractToWorld(world, "alphaSC", alphaCode)
	alpha.Balance = big.NewInt(100)

	betaCode := test.GetTestSCCodeModule("exec-sync-ctx-multiple/beta", "beta", "../../")
	gammaCode := test.GetTestSCCodeModule("exec-sync-ctx-multiple/gamma", "gamma", "../../")
	deltaCode := test.GetTestSCCodeModule("exec-sync-ctx-multiple/delta", "delta", "../../")

	_ = test.AddTestSmartContractToWorld(world, "betaSC", betaCode)
	_ = test.AddTestSmartContractToWorld(world, "gammaSC", gammaCode)
	_ = test.AddTestSmartContractToWorld(world, "deltaSC", deltaCode)

	expectedReturnData := [][]byte{
		[]byte("arg1"),
		[]byte("succ"),
		[]byte("arg2"),
		[]byte("succ"),
		[]byte("arg3"),
		[]byte("succ"),
	}

	// Alpha uses executeOnSameContext() to call beta, gamma and delta one after
	// the other, in the same transaction.
	input := test.DefaultTestContractCallInput()
	input.Function = "callChildrenDirectly_SameCtx"
	input.GasProvided = 1000000
	input.RecipientAddr = alpha.Address

	vmOutput, err := host.RunSmartContractCall(input)

	verify := test.NewVMOutputVerifier(t, vmOutput, err)
	verify.
		Ok().
		ReturnData(expectedReturnData...)
}

func TestExecution_ExecuteOnDestContext_MultipleChildren(t *testing.T) {
	world := worldmock.NewMockWorld()
	host := test.DefaultTestArwen(t, world)

	alphaCode := test.GetTestSCCodeModule("exec-sync-ctx-multiple/alpha", "alpha", "../../")
	alpha := test.AddTestSmartContractToWorld(world, "alphaSC", alphaCode)
	alpha.Balance = big.NewInt(100)

	betaCode := test.GetTestSCCodeModule("exec-sync-ctx-multiple/beta", "beta", "../../")
	gammaCode := test.GetTestSCCodeModule("exec-sync-ctx-multiple/gamma", "gamma", "../../")
	deltaCode := test.GetTestSCCodeModule("exec-sync-ctx-multiple/delta", "delta", "../../")

	_ = test.AddTestSmartContractToWorld(world, "betaSC", betaCode)
	_ = test.AddTestSmartContractToWorld(world, "gammaSC", gammaCode)
	_ = test.AddTestSmartContractToWorld(world, "deltaSC", deltaCode)

	expectedReturnData := [][]byte{
		[]byte("arg1"),
		[]byte("succ"),
		[]byte("arg2"),
		[]byte("succ"),
		[]byte("arg3"),
		[]byte("succ"),
	}

	// Alpha uses executeOnDestContext() to call beta, gamma and delta one after
	// the other, in the same transaction.
	input := test.DefaultTestContractCallInput()
	input.Function = "callChildrenDirectly_DestCtx"
	input.GasProvided = 1000000
	input.RecipientAddr = alpha.Address

	vmOutput, err := host.RunSmartContractCall(input)

	verify := test.NewVMOutputVerifier(t, vmOutput, err)
	verify.
		Ok().
		ReturnData(expectedReturnData...)
}

func TestExecution_ExecuteOnDestContextByCaller_SimpleTransfer(t *testing.T) {
	// The child contract is designed to send some tokens back to its caller, as
	// many as requested. The parent calls the child using
	// executeOnDestContextByCaller(), which means that the child will not see
	// the parent as its caller, but the original caller of the transaction
	// instead. Thus the original caller (the user address) will receive 42
	// tokens, and not the parent, even if the parent is the one making the call
	// to the child.

	transferValue := int64(42)

	test.BuildInstanceCallTest(t).
		WithContracts(
			test.CreateInstanceContract(test.ParentAddress).
				WithCode(test.GetTestSCCodeModule("exec-dest-ctx-by-caller/parent", "parent", "../../")).
				WithBalance(1000),
			test.CreateInstanceContract(test.ChildAddress).
				WithCode(test.GetTestSCCodeModule("exec-dest-ctx-by-caller/child", "child", "../../")).
				WithBalance(1000),
		).
		WithInput(test.CreateTestContractCallInputBuilder().
			WithRecipientAddr(test.ParentAddress).
			WithFunction("call_child").
			WithGasProvided(2000).
			Build()).
		AndAssertResults(func(host arwen.VMHost, stubBlockchainHook *contextmock.BlockchainHookStub, verify *test.VMOutputVerifier) {
			verify.
				Ok().
				// test.ParentAddress
				GasUsed(test.ParentAddress, 762).
				/// test.ChildAddress
				Balance(test.ChildAddress, 1000).
				BalanceDelta(test.ChildAddress, -transferValue).
				GasUsed(test.ChildAddress, 667).
				// test.UserAddress
				BalanceDelta(test.UserAddress, transferValue).
				// other
				ReturnData([]byte("sent"), []byte("child called")).
				Transfers(
					test.CreateTransferEntry(test.ChildAddress, test.UserAddress).
						WithData([]byte{}).
						WithValue(big.NewInt(transferValue)),
				)
		})
}

func TestExecution_AsyncCall_GasLimitConsumed(t *testing.T) {
	parentCode := test.GetTestSCCode("async-call-parent", "../../")
	childCode := test.GetTestSCCode("async-call-child", "../../")

	test.BuildInstanceCallTest(t).
		WithContracts(
			test.CreateInstanceContract(test.ParentAddress).
				WithCode(parentCode).
				WithBalance(1000),
			test.CreateInstanceContract(test.ChildAddress).
				WithCode(childCode).
				WithBalance(1000),
		).
		WithInput(test.CreateTestContractCallInputBuilder().
			WithRecipientAddr(test.ParentAddress).
			WithFunction(parentPerformAsyncCall).
			WithGasProvided(1000000).
			WithArguments([]byte{0}).
			Build()).
		WithSetup(func(host arwen.VMHost, stubBlockchainHook *contextmock.BlockchainHookStub) {
			stubBlockchainHook.GetUserAccountCalled = func(scAddress []byte) (vmcommon.UserAccountHandler, error) {
				if bytes.Equal(scAddress, test.ParentAddress) {
					return &contextmock.StubAccount{
						Address: test.ParentAddress,
						Balance: big.NewInt(1000),
					}, nil
				}
				return nil, test.ErrAccountNotFound
			}
			stubBlockchainHook.GetCodeCalled = func(account vmcommon.UserAccountHandler) []byte {
				if bytes.Equal(test.ParentAddress, account.AddressBytes()) {
					return parentCode
				}
				return nil
			}
			stubBlockchainHook.GetShardOfAddressCalled = func(address []byte) uint32 {
				if bytes.Equal(address, test.ParentAddress) {
					return 0
				}
				return 1
			}
		}).
		AndAssertResults(func(host arwen.VMHost, stubBlockchainHook *contextmock.BlockchainHookStub, verify *test.VMOutputVerifier) {
			verify.
				Ok().
				GasRemaining(0)
		})
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

	// Call parentFunctionChildCall() of the parent SC, which will call the child
	// SC and pass some arguments using asyncCall().

	test.BuildInstanceCallTest(t).
		WithContracts(
			test.CreateInstanceContract(test.ParentAddress).
				WithCode(test.GetTestSCCode("async-call-parent", "../../")).
				WithBalance(1000),
			test.CreateInstanceContract(test.ChildAddress).
				WithCode(test.GetTestSCCode("async-call-child", "../../")).
				WithBalance(1000),
		).
		WithInput(test.CreateTestContractCallInputBuilder().
			WithRecipientAddr(test.ParentAddress).
			WithFunction(parentPerformAsyncCall).
			WithGasProvided(116000).
			WithArguments([]byte{0}).
			Build()).
		AndAssertResults(func(host arwen.VMHost, stubBlockchainHook *contextmock.BlockchainHookStub, verify *test.VMOutputVerifier) {
			verify.
				Ok().
				// TODO matei-p restore correct gas assertions
				// GasUsed(test.ParentAddress, 9114).
				// GasUsed(test.ChildAddress, 2534).
				// GasRemaining(104352).
				Balance(test.ParentAddress, 1000).
				Balance(test.ChildAddress, 1000).
				BalanceDelta(test.ThirdPartyAddress, 6).
				ReturnData(test.ParentFinishA, test.ParentFinishB, []byte{0}, []byte("thirdparty"), []byte("vault"), []byte{0}, []byte("succ")).
				Storage(
					test.CreateStoreEntry(test.ParentAddress).WithKey(test.ParentKeyA).WithValue(test.ParentDataA),
					test.CreateStoreEntry(test.ParentAddress).WithKey(test.ParentKeyB).WithValue(test.ParentDataB),
					test.CreateStoreEntry(test.ChildAddress).WithKey(test.ChildKey).WithValue(test.ChildData),
				).
				Transfers(
					test.CreateTransferEntry(test.ParentAddress, test.ThirdPartyAddress).
						WithData([]byte("hello")).
						WithValue(big.NewInt(3)),
					test.CreateTransferEntry(test.ChildAddress, test.ThirdPartyAddress).
						WithData([]byte(" there")).
						WithValue(big.NewInt(3)),
					test.CreateTransferEntry(test.ChildAddress, test.VaultAddress).
						WithData([]byte{}).
						WithValue(big.NewInt(4)),
				)
		})
}

func TestExecution_AsyncCall_ChildFails(t *testing.T) {
	// Scenario
	// Identical to TestExecution_AsyncCall(), except that the child is
	// instructed to call signalError().
	// Because "vault" was not received by the callBack(), the Parent sends 4 ERD
	// to the Vault directly.

	// Call parentFunctionChildCall() of the parent SC, which will call the child
	// SC and pass some arguments using asyncCall().

	test.BuildInstanceCallTest(t).
		WithContracts(
			test.CreateInstanceContract(test.ParentAddress).
				WithCode(test.GetTestSCCode("async-call-parent", "../../")).
				WithBalance(1000),
			test.CreateInstanceContract(test.ChildAddress).
				WithCode(test.GetTestSCCode("async-call-child", "../../")).
				WithBalance(1000),
		).
		WithInput(test.CreateTestContractCallInputBuilder().
			WithRecipientAddr(test.ParentAddress).
			WithFunction(parentPerformAsyncCall).
			WithGasProvided(1000000).
			WithArguments([]byte{1}).
			WithCurrentTxHash([]byte("txhash")).
			Build()).
		WithSetup(func(host arwen.VMHost, stubBlockchainHook *contextmock.BlockchainHookStub) {
			host.Metering().GasSchedule().ElrondAPICost.AsyncCallbackGasLock = 3000
		}).
		AndAssertResults(func(host arwen.VMHost, stubBlockchainHook *contextmock.BlockchainHookStub, verify *test.VMOutputVerifier) {
			verify.
				Ok().
				// TODO matei-p enable gas checks
				// GasUsed(test.ParentAddress, 998352).
				// GasRemaining(1648).
				ReturnData(test.ParentFinishA, test.ParentFinishB, []byte("succ")).
				Storage(
					test.CreateStoreEntry(test.ParentAddress).WithKey(test.ParentKeyA).WithValue(test.ParentDataA),
					test.CreateStoreEntry(test.ParentAddress).WithKey(test.ParentKeyB).WithValue(test.ParentDataB),
				)
		})
}

func TestExecution_AsyncCall_CallBackFails(t *testing.T) {
	// Scenario
	// Identical to TestExecution_AsyncCall(), except that the callback is
	// instructed to call signalError().

	// Call parentFunctionChildCall() of the parent SC, which will call the child
	// SC and pass some arguments using asyncCall().

	test.BuildInstanceCallTest(t).
		WithContracts(
			test.CreateInstanceContract(test.ParentAddress).
				WithCode(test.GetTestSCCode("async-call-parent", "../../")).
				WithBalance(1000),
			test.CreateInstanceContract(test.ChildAddress).
				WithCode(test.GetTestSCCode("async-call-child", "../../")).
				WithBalance(1000),
		).
		WithInput(test.CreateTestContractCallInputBuilder().
			WithRecipientAddr(test.ParentAddress).
			WithFunction(parentPerformAsyncCall).
			WithGasProvided(200000).
			WithArguments([]byte{0, 3}).
			WithCurrentTxHash([]byte("txhash")).
			Build()).
		AndAssertResults(func(host arwen.VMHost, stubBlockchainHook *contextmock.BlockchainHookStub, verify *test.VMOutputVerifier) {
			verify.
				Ok().
				ReturnMessage("callBack error").
				// TODO matei-p enable gas checks
				// GasUsed(test.ParentAddress, 197437).
				// GasUsed(test.ChildAddress, 2534).
				// TODO Why is there a minuscule amount of gas remaining after the callback
				// fails? This is supposed to be 0.
				// GasRemaining(29).
				BalanceDelta(test.ThirdPartyAddress, 6).
				BalanceDelta(test.ChildAddress, big.NewInt(0).Sub(big.NewInt(1), big.NewInt(1)).Int64()).
				ReturnData(test.ParentFinishA, test.ParentFinishB, []byte{3}, []byte("thirdparty"), []byte("vault"), []byte("user error"), []byte("txhash")).
				Storage(
					test.CreateStoreEntry(test.ParentAddress).WithKey(test.ParentKeyA).WithValue(test.ParentDataA),
					test.CreateStoreEntry(test.ParentAddress).WithKey(test.ParentKeyB).WithValue(test.ParentDataB),
					test.CreateStoreEntry(test.ChildAddress).WithKey(test.ChildKey).WithValue(test.ChildData),
				).
				Transfers(
					test.CreateTransferEntry(test.ParentAddress, test.ThirdPartyAddress).
						WithData([]byte("hello")).
						WithValue(big.NewInt(3)),
					test.CreateTransferEntry(test.ChildAddress, test.ThirdPartyAddress).
						WithData([]byte(" there")).
						WithValue(big.NewInt(3)),
					test.CreateTransferEntry(test.ChildAddress, test.VaultAddress).
						WithData([]byte{}).
						WithValue(big.NewInt(4)),
				)
		})
}

func TestExecution_CreateNewContract_Success(t *testing.T) {
	childCode := test.GetTestSCCode("init-correct", "../../")
	childAddress := []byte("newAddress")
	l := len(childCode)

	test.BuildInstanceCallTest(t).
		WithContracts(
			test.CreateInstanceContract(test.ParentAddress).
				WithCode(test.GetTestSCCode("deployer", "../../")).
				WithBalance(1000),
		).
		WithInput(test.CreateTestContractCallInputBuilder().
			WithRecipientAddr(test.ParentAddress).
			WithFunction("deployChildContract").
			WithGasProvided(1_000_000).
			WithArguments([]byte{'A'}, []byte{0}).
			WithCurrentTxHash([]byte("txhash")).
			Build()).
		WithSetup(func(host arwen.VMHost, stubBlockchainHook *contextmock.BlockchainHookStub) {
			stubBlockchainHook.GetStorageDataCalled = func(address []byte, key []byte) ([]byte, error) {
				if bytes.Equal(address, test.ParentAddress) {
					if bytes.Equal(key, []byte{'A'}) {
						return childCode, nil
					}
					return nil, nil
				}
				return nil, arwen.ErrInvalidAccount
			}
		}).
		AndAssertResults(func(host arwen.VMHost, stubBlockchainHook *contextmock.BlockchainHookStub, verify *test.VMOutputVerifier) {
			verify.
				Ok().
				// test.ParentAddress
				Balance(test.ParentAddress, 1000).
				GasUsed(test.ParentAddress, 885).
				/// test.ChildAddress
				BalanceDelta(childAddress, 42).
				Code(childAddress, childCode).
				CodeMetadata(childAddress, []byte{1, 0}).
				CodeDeployerAddress(childAddress, test.ParentAddress).
				GasUsed(childAddress, 472).
				// other
				ReturnData([]byte{byte(l / 256), byte(l % 256)}, []byte("init successful"), []byte("succ")).
				Storage(
					test.CreateStoreEntry(test.ParentAddress).WithKey([]byte{'A'}).WithValue(childCode))
		})
}

func TestExecution_DeployNewContractFromExistingCode_Success(t *testing.T) {
	sourceAddress := testcommon.MakeTestSCAddress("sourceAddress")
	sourceCode := test.GetTestSCCode("init-correct", "../../")
	generatedNewAddress := []byte("newAddress")

	test.BuildInstanceCallTest(t).
		WithContracts(
			test.CreateInstanceContract(sourceAddress).
				WithCode(sourceCode).
				WithBalance(1000),
			test.CreateInstanceContract(test.ParentAddress).
				WithCode(test.GetTestSCCode("deployer-fromanother-contract", "../../")).
				WithBalance(1000),
		).
		WithInput(test.CreateTestContractCallInputBuilder().
			WithRecipientAddr(test.ParentAddress).
			WithFunction("deployCodeFromAnotherContract").
			WithArguments(sourceAddress).
			WithGasProvided(1_000_000).
			Build()).
		AndAssertResults(func(host arwen.VMHost, stubBlockchainHook *contextmock.BlockchainHookStub, verify *test.VMOutputVerifier) {
			verify.
				Ok().
				Code(generatedNewAddress, sourceCode).
				CodeMetadata(generatedNewAddress, testcommon.DefaultCodeMetadata).
				ReturnData(
					// returned by the new deployed contract from the existing source code
					[]byte("init successful"),
					// returned by the deployer contract
					[]byte("succ"),
				)
		})
}

func TestExecution_UpgradeContractFromExistingCode_Success(t *testing.T) {
	initialAddress := testcommon.MakeTestSCAddress("destAddress")
	initialCode := test.GetTestSCCode("init-simple", "../../")
	sourceAddress := testcommon.MakeTestSCAddress("sourceAddress")
	sourceCode := test.GetTestSCCode("init-correct", "../../")

	test.BuildInstanceCallTest(t).
		WithContracts(
			test.CreateInstanceContract(sourceAddress).
				WithCode(sourceCode).
				WithBalance(1000),
			test.CreateInstanceContract(initialAddress).
				WithCode(initialCode).
				WithBalance(1000),
			test.CreateInstanceContract(test.ParentAddress).
				WithCode(test.GetTestSCCode("upgrader-fromanother-contract", "../../")).
				WithBalance(1000),
		).
		WithInput(test.CreateTestContractCallInputBuilder().
			WithRecipientAddr(test.ParentAddress).
			WithFunction("upgradeCodeFromAnotherContract").
			WithArguments(initialAddress, sourceAddress).
			WithGasProvided(1_000_000).
			Build()).
		AndAssertResults(func(host arwen.VMHost, stubBlockchainHook *contextmock.BlockchainHookStub, verify *test.VMOutputVerifier) {
			verify.
				Ok().
				ReturnData(
					// returned by the replaced contract code
					[]byte("init successful"),
				)
		})
}

func TestExecution_CreateNewContract_Fail(t *testing.T) {
	childCode := test.GetTestSCCode("init-correct", "../../")
	l := len(childCode)

	test.BuildInstanceCallTest(t).
		WithContracts(
			test.CreateInstanceContract(test.ParentAddress).
				WithCode(test.GetTestSCCode("deployer", "../../")).
				WithBalance(1000),
		).
		WithInput(test.CreateTestContractCallInputBuilder().
			WithRecipientAddr(test.ParentAddress).
			WithFunction("deployChildContract").
			WithGasProvided(1_000_000).
			WithArguments([]byte{'A'}, []byte{1}).
			Build()).
		WithSetup(func(host arwen.VMHost, stubBlockchainHook *contextmock.BlockchainHookStub) {
			stubBlockchainHook.GetStorageDataCalled = func(address []byte, key []byte) ([]byte, error) {
				if bytes.Equal(address, test.ParentAddress) {
					if bytes.Equal(key, []byte{'A'}) {
						return childCode, nil
					}
					return nil, nil
				}
				return nil, arwen.ErrInvalidAccount
			}
		}).
		AndAssertResults(func(host arwen.VMHost, stubBlockchainHook *contextmock.BlockchainHookStub, verify *test.VMOutputVerifier) {
			verify.
				Ok().
				GasUsed(test.ParentAddress, 2885).
				ReturnData([]byte{byte(l / 256), byte(l % 256)}, []byte("fail")).
				Storage(test.CreateStoreEntry(test.ParentAddress).WithKey([]byte{'A'}).WithValue(childCode))
		})
}

func TestExecution_CreateNewContract_IsSmartContract(t *testing.T) {

	childCode := test.GetTestSCCode("deployer-child", "../../")

	newAddr := "newAddr_"
	ownerNonce := uint64(23)
	parentAddress := testcommon.MakeTestSCAddress(fmt.Sprintf("%s_%d", newAddr, 24))
	childAddress := testcommon.MakeTestSCAddress(fmt.Sprintf("%s_%d", newAddr, 25))

	input := test.CreateTestContractCreateInputBuilder().
		WithCallValue(1000).
		WithGasProvided(100_000).
		WithContractCode(test.GetTestSCCode("deployer-parent", "../../")).
		WithArguments(parentAddress, childCode).
		Build()

	test.BuildInstanceCreatorTest(t).
		WithInput(input).
		WithSetup(func(host arwen.VMHost, stubBlockchainHook *contextmock.BlockchainHookStub) {
			stubBlockchainHook.GetUserAccountCalled = func(address []byte) (vmcommon.UserAccountHandler, error) {
				strAddress := string(address)
				if strAddress == string(childAddress) {
					return nil, errors.New("not found")
				}
				return &contextmock.StubAccount{
					Nonce: 24,
				}, nil
			}
			stubBlockchainHook.NewAddressCalled = func(creatorAddress []byte, nonce uint64, vmType []byte) ([]byte, error) {
				ownerNonce++
				return testcommon.MakeTestSCAddress(fmt.Sprintf("%s_%d", newAddr, ownerNonce)), nil
			}
			stubBlockchainHook.IsSmartContractCalled = func(address []byte) bool {
				outputAccounts := host.Output().GetOutputAccounts()
				_, isSmartContract := outputAccounts[string(address)]
				return isSmartContract
			}
		}).
		AndAssertResults(func(blockchainHook *contextmock.BlockchainHookStub, verify *test.VMOutputVerifier) {
			verify.
				Ok().
				ReturnData([]byte("succ")) /* returned from child contract init */
		})
}

func TestExecution_Mocked_Wasmer_Instances(t *testing.T) {
	test.BuildMockInstanceCallTest(t).
		WithContracts(
			test.CreateMockContract(test.ParentAddress).
				WithBalance(1000).
				WithMethods(func(parentInstance *mock.InstanceMock, testConfig *testcommon.TestConfig) {
					parentInstance.AddMockMethod("callChild", func() *mock.InstanceMock {
						host := parentInstance.Host
						host.Output().Finish([]byte("parent returns this"))
						host.Metering().UseGas(500)
						_, err := host.Storage().SetStorage([]byte("parent"), []byte("parent storage"))
						require.Nil(t, err)
						childInput := test.DefaultTestContractCallInput()
						childInput.CallerAddr = test.ParentAddress
						childInput.RecipientAddr = test.ChildAddress
						childInput.CallValue = big.NewInt(4)
						childInput.Function = "doSomething"
						childInput.GasProvided = 1000
						_, err = host.ExecuteOnDestContext(childInput)
						require.Nil(t, err)
						return parentInstance
					})
				}),
			test.CreateMockContract(test.ChildAddress).
				WithBalance(0).
				WithMethods(func(childInstance *mock.InstanceMock, testConfig *testcommon.TestConfig) {
					childInstance.AddMockMethod("doSomething", func() *mock.InstanceMock {
						host := childInstance.Host
						host.Output().Finish([]byte("child returns this"))
						host.Metering().UseGas(100)
						_, err := host.Storage().SetStorage([]byte("child"), []byte("child storage"))
						require.Nil(t, err)
						return childInstance
					})
				}),
		).
		WithInput(test.CreateTestContractCallInputBuilder().
			WithRecipientAddr(test.ParentAddress).
			WithGasProvided(1000).
			WithFunction("callChild").
			Build()).
		AndAssertResults(func(world *worldmock.MockWorld, verify *test.VMOutputVerifier) {
			verify.
				Ok().
				// test.ParentAddress
				Balance(test.ParentAddress, 1000).
				BalanceDelta(test.ParentAddress, -4).
				GasUsed(test.ParentAddress, 547).
				/// test.ChildAddress
				BalanceDelta(test.ChildAddress, 4).
				GasUsed(test.ChildAddress, 146).
				// other
				GasRemaining(307).
				ReturnData([]byte("parent returns this"), []byte("child returns this")).
				Storage(
					test.CreateStoreEntry(test.ParentAddress).WithKey([]byte("parent")).WithValue([]byte("parent storage")),
					test.CreateStoreEntry(test.ChildAddress).WithKey([]byte("child")).WithValue([]byte("child storage")),
				)
		})
}

// makeBytecodeWithLocals rewrites the bytecode of "answer" to change the
// number of i64 locals it instantiates
func makeBytecodeWithLocals(numLocals uint64) []byte {
	originalCode := test.GetTestSCCode("answer", "../../")
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
