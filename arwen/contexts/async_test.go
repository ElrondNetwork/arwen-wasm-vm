package contexts

import (
	"errors"
	"math/big"
	"testing"

	"github.com/ElrondNetwork/arwen-wasm-vm/v1_3/arwen"
	"github.com/ElrondNetwork/arwen-wasm-vm/v1_3/config"
	"github.com/ElrondNetwork/arwen-wasm-vm/v1_3/crypto/factory"
	contextmock "github.com/ElrondNetwork/arwen-wasm-vm/v1_3/mock/context"
	worldmock "github.com/ElrondNetwork/arwen-wasm-vm/v1_3/mock/world"
	"github.com/ElrondNetwork/arwen-wasm-vm/v1_3/wasmer"
	vmcommon "github.com/ElrondNetwork/elrond-vm-common"
	"github.com/stretchr/testify/require"
)

var mockWasmerInstance *wasmer.Instance
var OriginalCaller = []byte("address_original_caller")
var Alice = []byte("address_alice")
var Bob = []byte("address_bob")

const GasForAsyncStep = config.GasValueForTests

func InitializeArwenAndWasmer_AsyncContext() (*contextmock.VMHostMock, *worldmock.MockWorld) {
	imports := MakeAPIImports()
	_ = wasmer.SetImports(imports)

	gasSchedule := config.MakeGasMapForTests()
	gasCostConfig, _ := config.CreateGasConfig(gasSchedule)
	opcodeCosts := gasCostConfig.WASMOpcodeCost.ToOpcodeCostsArray()
	wasmer.SetOpcodeCosts(&opcodeCosts)

	host := &contextmock.VMHostMock{}
	host.SCAPIMethods = imports

	mockMetering := &contextmock.MeteringContextMock{}
	mockMetering.SetGasSchedule(gasSchedule)
	host.MeteringContext = mockMetering

	world := worldmock.NewMockWorld()
	host.BlockchainContext, _ = NewBlockchainContext(host, world)

	mockWasmerInstance = &wasmer.Instance{
		Exports: make(wasmer.ExportsMap),
	}
	runtimeContext, _ := NewRuntimeContext(host, []byte("vm"), false)
	runtimeContext.instance = mockWasmerInstance
	host.RuntimeContext = runtimeContext

	host.OutputContext, _ = NewOutputContext(host)
	host.CryptoHook = factory.NewVMCrypto()
	return host, world
}

func InitializeArwenAndWasmer_AsyncContext_AliceAndBob() (
	*contextmock.VMHostMock,
	*worldmock.MockWorld,
	*vmcommon.ContractCallInput,
) {
	host, world := InitializeArwenAndWasmer_AsyncContext()
	world.AcctMap.PutAccount(&worldmock.Account{
		Address: Alice,
		Balance: big.NewInt(88),
	})
	world.AcctMap.PutAccount(&worldmock.Account{
		Address: Bob,
		Balance: big.NewInt(12),
	})

	originalVMInput := &vmcommon.ContractCallInput{
		VMInput: vmcommon.VMInput{
			CallerAddr:     OriginalCaller,
			Arguments:      nil,
			CallType:       vmcommon.DirectCall,
			GasPrice:       1,
			CurrentTxHash:  []byte("txhash"),
			PrevTxHash:     []byte("txhash"),
			OriginalTxHash: []byte("txhash"),
		},
		RecipientAddr: Alice,
		Function:      "alice_function",
	}

	return host, world, originalVMInput
}

func TestNewAsyncContext(t *testing.T) {
	host, _ := InitializeArwenAndWasmer_AsyncContext()
	async := NewAsyncContext(host)
	require.NotNil(t, async)

	require.NotNil(t, async.host)
	require.Nil(t, async.stateStack)
	require.Nil(t, async.callerAddr)
	require.Nil(t, async.returnData)
	require.NotNil(t, async.asyncCallGroups)
	require.Empty(t, async.asyncCallGroups)
}

func TestAsyncContext_InitState(t *testing.T) {
	host, _ := InitializeArwenAndWasmer_AsyncContext()
	async := NewAsyncContext(host)
	require.NotNil(t, async)

	async.callerAddr = []byte("some address")
	async.gasPrice = 1000
	async.returnData = []byte("some return data")
	async.asyncCallGroups = nil

	async.InitState()

	require.NotNil(t, async.callerAddr)
	require.Empty(t, async.callerAddr)
	require.Zero(t, async.gasPrice)
	require.NotNil(t, async.returnData)
	require.Empty(t, async.returnData)
	require.NotNil(t, async.asyncCallGroups)
	require.Empty(t, async.asyncCallGroups)
}

func TestAsyncContext_InitStateFromContractCallInput(t *testing.T) {
	host, _ := InitializeArwenAndWasmer_AsyncContext()
	async := NewAsyncContext(host)
	require.NotNil(t, async)

	async.callerAddr = []byte("some address")
	async.gasPrice = 1000
	async.returnData = []byte("some return data")
	async.asyncCallGroups = nil

	input := &vmcommon.ContractCallInput{
		VMInput: vmcommon.VMInput{
			CallerAddr: []byte("real caller addr"),
			GasPrice:   42,
		},
	}

	async.InitStateFromInput(input)

	require.Equal(t, input.CallerAddr, async.callerAddr)
	require.Equal(t, uint64(42), async.gasPrice)
	require.NotNil(t, async.returnData)
	require.Empty(t, async.returnData)
	require.NotNil(t, async.asyncCallGroups)
	require.Empty(t, async.asyncCallGroups)
}

func TestAsyncContext_GettersAndSetters(t *testing.T) {
	host, _ := InitializeArwenAndWasmer_AsyncContext()
	async := NewAsyncContext(host)
	require.NotNil(t, async)

	async.callerAddr = []byte("some address")
	async.gasPrice = 1000
	async.returnData = []byte("some return data")

	require.Equal(t, []byte("some address"), async.GetCallerAddress())
	require.Equal(t, uint64(1000), async.GetGasPrice())
	require.Equal(t, []byte("some return data"), async.GetReturnData())

	async.SetReturnData([]byte("rockets"))
	require.Equal(t, []byte("rockets"), async.GetReturnData())
}

func TestAsyncContext_RegisterAsyncCall_NewGroup_DeleteGroup(t *testing.T) {
	host, _ := InitializeArwenAndWasmer_AsyncContext()

	async := NewAsyncContext(host)
	require.NotNil(t, async)

	require.True(t, async.IsComplete())
	require.False(t, async.HasPendingCallGroups())

	group, exists := async.GetCallGroup("testGroup")
	require.Nil(t, group)
	require.False(t, exists)

	err := async.RegisterAsyncCall("testGroup", &arwen.AsyncCall{
		Destination: []byte("somewhere"),
		Data:        []byte("something"),
	})
	require.Nil(t, err)
	require.False(t, async.IsComplete())
	require.True(t, async.HasPendingCallGroups())
	require.Equal(t, 1, len(async.asyncCallGroups))

	group, exists = async.GetCallGroup("testGroup")
	require.NotNil(t, group)
	require.True(t, exists)
	require.True(t, group.HasPendingCalls())
	require.False(t, group.IsComplete())
	require.False(t, group.HasCallback())

	async.deleteCallGroupByID("testGroup")
	group, exists = async.GetCallGroup("testGroup")
	require.Nil(t, group)
	require.False(t, exists)
}

func TestAsyncContext_RegisterAsyncCall_ExistingGroup(t *testing.T) {
	host, _ := InitializeArwenAndWasmer_AsyncContext()

	async := NewAsyncContext(host)
	require.NotNil(t, async)
	require.Equal(t, 0, len(async.asyncCallGroups))

	err := async.AddCallGroup(arwen.NewAsyncCallGroup("testGroup"))
	require.Nil(t, err)
	require.Equal(t, 1, len(async.asyncCallGroups))
	require.False(t, async.IsComplete())

	group, exists := async.GetCallGroup("testGroup")
	require.NotNil(t, group)
	require.True(t, exists)

	err = async.RegisterAsyncCall("testGroup", &arwen.AsyncCall{
		Destination: []byte("somewhere"),
		Data:        []byte("something"),
	})
	require.Nil(t, err)
	require.Equal(t, 1, len(async.asyncCallGroups))
	require.False(t, async.IsComplete())
}

func TestAsyncContext_RegisterAsyncCall_ValidationAndFields(t *testing.T) {
	// TODO execution mode
	// TODO non-nil destination
	// TODO locked gas
}

func TestAsyncContext_SetGroupCallback_GroupDoesntExist(t *testing.T) {
	host, _ := InitializeArwenAndWasmer_AsyncContext()
	async := NewAsyncContext(host)

	mockWasmerInstance.Exports["callbackFunction"] = nil
	err := async.SetGroupCallback("testGroup", "callbackFunction", []byte{}, 0)
	require.True(t, errors.Is(err, arwen.ErrAsyncCallGroupDoesNotExist))
}

func TestAsyncContext_SetGroupCallback_OutOfGas(t *testing.T) {
	host, _ := InitializeArwenAndWasmer_AsyncContext()
	async := NewAsyncContext(host)

	err := async.RegisterAsyncCall("testGroup", &arwen.AsyncCall{
		Destination: []byte("somewhere"),
		Data:        []byte("something"),
	})
	require.Nil(t, err)

	mockMetering := host.Metering().(*contextmock.MeteringContextMock)
	mockMetering.Err = arwen.ErrNotEnoughGas

	mockWasmerInstance.Exports["callbackFunction"] = nil
	err = async.SetGroupCallback("testGroup", "callbackFunction", []byte{}, 0)
	require.True(t, errors.Is(err, arwen.ErrNotEnoughGas))
}

func TestAsyncContext_SetGroupCallback_Success(t *testing.T) {
	host, _ := InitializeArwenAndWasmer_AsyncContext()
	async := NewAsyncContext(host)

	mockMetering := host.Metering().(*contextmock.MeteringContextMock)
	mockMetering.GasComputedToLock = 42

	err := async.RegisterAsyncCall("testGroup", &arwen.AsyncCall{
		Destination: []byte("somewhere"),
		Data:        []byte("something"),
	})
	require.Nil(t, err)

	mockWasmerInstance.Exports["callbackFunction"] = nil
	err = async.SetGroupCallback("testGroup", "callbackFunction", []byte{}, 0)
	require.Nil(t, err)

	group, exists := async.GetCallGroup("testGroup")
	require.NotNil(t, group)
	require.True(t, exists)
	require.Equal(t, "callbackFunction", group.Callback)
	require.Equal(t, []byte{}, group.CallbackData)
	require.Equal(t, uint64(42), group.GasLocked)
}

func TestAsyncContext_DetermineExecutionMode(t *testing.T) {
	leftAddress := []byte("left")
	leftAccount := &worldmock.Account{
		Address: leftAddress,
		Code:    []byte("left code"),
		ShardID: 0,
	}

	rightAddress := []byte("right")
	rightAccount := &worldmock.Account{
		Address: rightAddress,
		Code:    []byte("right code"),
		ShardID: 0,
	}

	host, world := InitializeArwenAndWasmer_AsyncContext()
	world.AcctMap.PutAccount(leftAccount)
	world.AcctMap.PutAccount(rightAccount)
	runtime := host.Runtime()

	async := NewAsyncContext(host)

	runtime.SetSCAddress(leftAddress)
	execMode, err := async.determineExecutionMode(rightAddress, []byte("func"))
	require.Nil(t, err)
	require.Equal(t, arwen.SyncExecution, execMode)

	execMode, err = async.determineExecutionMode(rightAddress, []byte(""))
	require.NotNil(t, err)
	require.Equal(t, arwen.AsyncUnknown, execMode)

	execMode, err = async.determineExecutionMode(rightAddress, []byte(""))
	require.NotNil(t, err)
	require.Equal(t, arwen.AsyncUnknown, execMode)

	host.IsBuiltinFunc = true
	runtime.SetSCAddress(leftAddress)
	execMode, err = async.determineExecutionMode(rightAddress, []byte("func"))
	require.Nil(t, err)
	require.Equal(t, arwen.AsyncBuiltinFuncIntraShard, execMode)

	host.IsBuiltinFunc = false
	rightAccount.Code = []byte{}
	rightAccount.ShardID = 1

	// Erase the code of the rightAccount from the Output context cache
	outputAccount, _ := host.Output().GetOutputAccount(rightAddress)
	outputAccount.Code = []byte{}

	runtime.SetSCAddress(leftAddress)
	execMode, err = async.determineExecutionMode(rightAddress, []byte("func"))
	require.Nil(t, err)
	require.Equal(t, arwen.AsyncUnknown, execMode)

	host.IsBuiltinFunc = true
	rightAccount.Code = []byte{}
	rightAccount.ShardID = 1

	// Erase the code of the rightAccount from the Output context cache
	outputAccount, _ = host.Output().GetOutputAccount(rightAddress)
	outputAccount.Code = []byte{}

	runtime.SetSCAddress(leftAddress)
	execMode, err = async.determineExecutionMode(rightAddress, []byte("func"))
	require.Nil(t, err)
	require.Equal(t, arwen.AsyncBuiltinFuncCrossShard, execMode)
}

func TestAsyncContext_IsValidCallbackName(t *testing.T) {
	host, _ := InitializeArwenAndWasmer_AsyncContext()
	async := NewAsyncContext(host)

	mockWasmerInstance.Exports["a"] = nil
	mockWasmerInstance.Exports["my_contract_method_22"] = nil
	mockWasmerInstance.Exports["not_builtin"] = nil
	mockWasmerInstance.Exports["callBack"] = nil
	mockWasmerInstance.Exports["callback"] = nil
	mockWasmerInstance.Exports["function_do"] = nil

	require.True(t, async.isValidCallbackName("a"))
	require.True(t, async.isValidCallbackName("my_contract_method_22"))

	require.True(t, async.isValidCallbackName("not_builtin"))
	host.IsBuiltinFunc = true
	require.False(t, async.isValidCallbackName("builtin"))
	host.IsBuiltinFunc = false

	require.True(t, async.isValidCallbackName("callBack"))
	require.True(t, async.isValidCallbackName("callback"))
	require.True(t, async.isValidCallbackName("function_do"))

	require.False(t, async.isValidCallbackName("function-do"))
	require.False(t, async.isValidCallbackName("3_my_contract_method_22"))
	require.False(t, async.isValidCallbackName("init"))
	require.False(t, async.isValidCallbackName("München"))
	require.False(t, async.isValidCallbackName("Göteborg"))
	require.False(t, async.isValidCallbackName("東京"))
	require.False(t, async.isValidCallbackName("function.org"))
	require.False(t, async.isValidCallbackName("Ainulindalë"))
}

func TestAsyncContext_FindCall(t *testing.T) {
	host, _ := InitializeArwenAndWasmer_AsyncContext()
	async := NewAsyncContext(host)

	groupID, index, err := async.findCall([]byte("somewhere"))
	require.Equal(t, "", groupID)
	require.Equal(t, -1, index)
	require.True(t, errors.Is(err, arwen.ErrAsyncCallNotFound))

	err = async.RegisterAsyncCall("testGroup", &arwen.AsyncCall{
		Destination: []byte("somewhere"),
		Data:        []byte("something"),
	})
	require.Nil(t, err)

	groupID, index, err = async.findCall([]byte("somewhere"))
	require.Nil(t, err)
	require.Equal(t, "testGroup", groupID)
	require.Equal(t, 0, index)

	err = async.RegisterAsyncCall("testGroup", &arwen.AsyncCall{
		Destination: []byte("somewhere_else"),
		Data:        []byte("something"),
	})
	require.Nil(t, err)

	groupID, index, err = async.findCall([]byte("somewhere_else"))
	require.Nil(t, err)
	require.Equal(t, "testGroup", groupID)
	require.Equal(t, 1, index)

	err = async.RegisterAsyncCall("another_testGroup", &arwen.AsyncCall{
		Destination: []byte("somewhere_else_entirely"),
		Data:        []byte("something"),
	})
	require.Nil(t, err)

	groupID, index, err = async.findCall([]byte("somewhere_else_entirely"))
	require.Nil(t, err)
	require.Equal(t, "another_testGroup", groupID)
	require.Equal(t, 0, index)
}

func TestAsyncContext_UpdateCurrentCallStatus(t *testing.T) {
	vmInput := &vmcommon.ContractCallInput{
		VMInput: vmcommon.VMInput{
			CallerAddr: []byte("caller"),
			Arguments:  [][]byte{{0}},
			CallType:   vmcommon.DirectCall,
		},
	}

	host, _ := InitializeArwenAndWasmer_AsyncContext()
	async := NewAsyncContext(host)

	// CallType == DirectCall, async.UpdateCurrentCallStatus() does nothing
	host.Runtime().InitStateFromContractCallInput(vmInput)
	asyncCall, err := async.UpdateCurrentCallStatus()
	require.Nil(t, asyncCall)
	require.Nil(t, err)

	// CallType == AsynchronousCall, async.UpdateCurrentCallStatus() does nothing
	vmInput.CallType = vmcommon.AsynchronousCall
	host.Runtime().InitStateFromContractCallInput(vmInput)
	asyncCall, err = async.UpdateCurrentCallStatus()
	require.Nil(t, asyncCall)
	require.Nil(t, err)

	// CallType == AsynchronousCallback, but no AsyncCalls registered in the
	// AsyncContext, so async.UpdateCurrentCallStatus() returns an error
	vmInput.CallType = vmcommon.AsynchronousCallBack
	vmInput.Arguments = nil
	host.Runtime().InitStateFromContractCallInput(vmInput)
	asyncCall, err = async.UpdateCurrentCallStatus()
	require.Nil(t, asyncCall)
	require.True(t, errors.Is(err, arwen.ErrCannotInterpretCallbackArgs))

	// CallType == AsynchronousCallback, but no AsyncCalls registered in the
	// AsyncContext, so async.UpdateCurrentCallStatus() returns an error
	vmInput.CallType = vmcommon.AsynchronousCallBack
	vmInput.Arguments = [][]byte{{0}}
	host.Runtime().InitStateFromContractCallInput(vmInput)
	asyncCall, err = async.UpdateCurrentCallStatus()
	require.Nil(t, asyncCall)
	require.True(t, errors.Is(err, arwen.ErrAsyncCallNotFound))

	// CallType == AsynchronousCallback, and there is an AsyncCall registered,
	// but it's not the expected one.
	err = async.RegisterAsyncCall("testGroup", &arwen.AsyncCall{
		Destination: []byte("some_address"),
		Data:        []byte("function"),
	})
	require.Nil(t, err)

	vmInput.CallType = vmcommon.AsynchronousCallBack
	vmInput.Arguments = [][]byte{{0}}
	host.Runtime().InitStateFromContractCallInput(vmInput)
	asyncCall, err = async.UpdateCurrentCallStatus()
	require.Nil(t, asyncCall)
	require.True(t, errors.Is(err, arwen.ErrAsyncCallNotFound))

	// CallType == AsynchronousCallback, but this time there is a corresponding AsyncCall
	// registered, causing async.UpdateCurrentCallStatus() to find and update the AsyncCall
	err = async.RegisterAsyncCall("testGroup", &arwen.AsyncCall{
		Destination: vmInput.CallerAddr,
		Data:        []byte("function"),
	})
	require.Nil(t, err)

	vmInput.CallType = vmcommon.AsynchronousCallBack
	vmInput.Arguments = [][]byte{{0}}
	host.Runtime().InitStateFromContractCallInput(vmInput)
	asyncCall, err = async.UpdateCurrentCallStatus()
	require.Nil(t, err)
	require.NotNil(t, asyncCall)
	require.Equal(t, arwen.AsyncCallResolved, asyncCall.Status)

	// CallType == AsynchronousCallback, there is a corresponding AsyncCall
	// registered, causing async.UpdateCurrentCallStatus() to find and update the
	// AsyncCall, but with AsyncCallRejected
	vmInput.CallType = vmcommon.AsynchronousCallBack
	vmInput.Arguments = [][]byte{{1}}
	host.Runtime().InitStateFromContractCallInput(vmInput)
	asyncCall, err = async.UpdateCurrentCallStatus()
	require.Nil(t, err)
	require.NotNil(t, asyncCall)
	require.Equal(t, arwen.AsyncCallRejected, asyncCall.Status)
}

func TestAsyncContext_SendAsyncCallCrossShard(t *testing.T) {
	host, world := InitializeArwenAndWasmer_AsyncContext()
	world.AcctMap.PutAccount(&worldmock.Account{
		Address: []byte("smartcontract"),
		Balance: big.NewInt(88),
	})

	host.Runtime().SetSCAddress([]byte("smartcontract"))
	async := NewAsyncContext(host)

	asyncCall := &arwen.AsyncCall{
		Destination: []byte("destination"),
		GasLimit:    42,
		GasLocked:   98,
		ValueBytes:  big.NewInt(88).Bytes(),
		Data:        []byte("some_data"),
	}

	host.Runtime().GetVMInput().GasProvided = 200
	err := async.sendAsyncCallCrossShard(asyncCall)
	require.Nil(t, err)

	mockMetering := host.Metering().(*contextmock.MeteringContextMock)
	mockMetering.GasProvidedMock = 200
	mockMetering.GasLeftMock = 60

	vmOutput := host.Output().GetVMOutput()
	require.NotNil(t, vmOutput)
	require.Equal(t, vmcommon.Ok, vmOutput.ReturnCode)

	smartcontract, ok := vmOutput.OutputAccounts["smartcontract"]
	require.True(t, ok)
	require.Equal(t, big.NewInt(-88), smartcontract.BalanceDelta)
	require.Empty(t, smartcontract.OutputTransfers)

	destination, ok := vmOutput.OutputAccounts["destination"]
	require.True(t, ok)
	require.Equal(t, big.NewInt(88), destination.BalanceDelta)
	require.Len(t, destination.OutputTransfers, 1)

	asyncTransfer := destination.OutputTransfers[0]
	require.Equal(t, big.NewInt(88), asyncTransfer.Value)
	require.Equal(t, uint64(42), asyncTransfer.GasLimit)
	require.Equal(t, uint64(98), asyncTransfer.GasLocked)
	require.Equal(t, []byte("some_data"), asyncTransfer.Data)
	require.Equal(t, vmcommon.AsynchronousCall, asyncTransfer.CallType)
}

func TestAsyncContext_ExecuteSyncCall_EarlyOutOfGas(t *testing.T) {
	// Scenario 1
	// Assert error propagation in async.executeSyncCall() from
	// async.createContractCallInput()
	host, _, originalVMInput := InitializeArwenAndWasmer_AsyncContext_AliceAndBob()
	host.Runtime().InitStateFromContractCallInput(originalVMInput)
	async := NewAsyncContext(host)

	asyncCall := defaultAsyncCall_AliceToBob()
	asyncCall.Data = []byte("function")
	asyncCall.GasLimit = 1
	err := async.executeAsyncLocalCall(asyncCall)
	require.True(t, errors.Is(err, arwen.ErrNotEnoughGas))
}

func TestAsyncContext_ExecuteSyncCall_NoDynamicGasLocking_Simulation(t *testing.T) {
	// Scenario 2
	// Successful execution at destination, but not enough gas for callback execution
	// (this situation should not happen in practice, due to dynamic gas locking)
	host, _, originalVMInput := InitializeArwenAndWasmer_AsyncContext_AliceAndBob()
	host.Runtime().InitStateFromContractCallInput(originalVMInput)
	async := NewAsyncContext(host)

	asyncCall := defaultAsyncCall_AliceToBob()
	asyncCall.GasLimit = 10

	gasConsumedByDestination := uint64(3)
	destOutput := &vmcommon.VMOutput{
		ReturnCode:   vmcommon.Ok,
		GasRemaining: asyncCall.GasLimit - gasConsumedByDestination,
	}
	host.EnqueueVMOutput(destOutput)

	err := async.executeAsyncLocalCall(asyncCall)
	require.Nil(t, err)
	require.Equal(t, arwen.AsyncCallResolved, asyncCall.Status)

	// Only one ContractCallInput was stored by the VmHostMock: constructing the
	// ContractCallInput for the callback has failed with insufficient gas before
	// reaching host.ExecutionOnDestContext()
	require.Len(t, host.StoredInputs, 1)

	// The ContractCallInput generated to execute the destination call synchronously
	destInput := defaultCallInput_AliceToBob(originalVMInput)
	destInput.GasProvided = asyncCall.GasLimit - GasForAsyncStep
	require.Equal(t, destInput, host.StoredInputs[0])

	// Verify the final VMOutput, containing the failure.
	expectedOutput := arwen.MakeVMOutput()
	expectedOutput.ReturnCode = vmcommon.Ok
	expectedOutput.ReturnMessage = "not enough gas"
	expectedOutput.GasRemaining = 0
	arwen.AddFinishData(expectedOutput, []byte("out of gas"))
	arwen.AddFinishData(expectedOutput, originalVMInput.CurrentTxHash)

	// The expectedOutput must also contain an OutputAccount corresponding to
	// Alice, because of a call to host.Output().GetOutputAccount() in
	// host.Output().GetVMOutput(), which creates and caches an empty account for
	// her.
	arwen.AddNewOutputAccount(expectedOutput, Alice, 0, nil)

	host.Output().GetOutputAccount(Alice) // TODO matei-p keep?
	vmOutput := host.Output().GetVMOutput()
	require.Equal(t, expectedOutput, vmOutput)
}

func TestAsyncContext_ExecuteSyncCall_Successful(t *testing.T) {
	// Scenario 3
	// Successful execution at destination, and successful callback execution;
	// the AsyncCall contains sufficient gas this time.
	host, _, originalVMInput := InitializeArwenAndWasmer_AsyncContext_AliceAndBob()
	host.Runtime().InitStateFromContractCallInput(originalVMInput)
	async := NewAsyncContext(host)
	asyncCall := defaultAsyncCall_AliceToBob()
	asyncCall.GasLimit = 100
	asyncCall.GasLocked = 90
	gasConsumedByDestination := uint64(23)
	gasConsumedByCallback := uint64(22)

	// The expected input passed to host.ExecuteOnDestContext() to call Bob as destination
	destInput := defaultCallInput_AliceToBob(originalVMInput)
	destInput.GasProvided = asyncCall.GasLimit - GasForAsyncStep
	destInput.GasLocked = asyncCall.GasLocked

	// Prepare the output of Bob (the destination call)
	destOutput := defaultDestOutput_Ok()
	destOutput.GasRemaining = destInput.GasProvided - gasConsumedByDestination

	// Prepare the input to Alice's callback
	callbackInput := defaultCallbackInput_BobToAlice(originalVMInput)
	callbackInput.GasProvided = destOutput.GasRemaining + asyncCall.GasLocked
	callbackInput.GasProvided -= defaultOutputDataLengthAsArgs(asyncCall, destOutput)
	callbackInput.GasProvided -= GasForAsyncStep
	callbackInput.GasLocked = 0

	// Prepare the output of Alice's callback
	callbackOutput := defaultCallbackOutput_Ok()
	callbackOutput.GasRemaining = callbackInput.GasProvided - gasConsumedByCallback

	// Enqueue the prepared VMOutputs
	host.EnqueueVMOutput(destOutput)
	host.EnqueueVMOutput(callbackOutput)

	err := async.executeAsyncLocalCall(asyncCall)
	require.Nil(t, err)
	require.Equal(t, arwen.AsyncCallResolved, asyncCall.Status)

	// There were two calls to host.ExecuteOnDestContext()
	require.Len(t, host.StoredInputs, 2)
	require.Equal(t, destInput, host.StoredInputs[0])
	require.Equal(t, callbackInput, host.StoredInputs[1])

	// Verify the final output of the execution; GasRemaining is set to 0 because
	// the test uses a mocked host.ExecuteOnDestContext(), which does not know to
	// manipulate the state stack of the OutputContext, therefore VMOutputs are
	// not merged between executions.
	expectedOutput := arwen.MakeVMOutput()
	expectedOutput.ReturnCode = vmcommon.Ok
	expectedOutput.GasRemaining = 0

	// The expectedOutput must also contain an OutputAccount corresponding to
	// Alice, because of a call to host.Output().GetOutputAccount() in
	// host.Output().GetVMOutput(), which creates and caches an empty account for
	// her.
	arwen.AddNewOutputAccount(expectedOutput, Alice, 0, nil)

	host.Output().GetOutputAccount(Alice) // TODO matei-p keep?
	actualOutput := host.Output().GetVMOutput()
	require.Equal(t, expectedOutput, actualOutput)
}

func TestAsyncContext_CreateContractCallInput(t *testing.T) {
	host, _, originalVMInput := InitializeArwenAndWasmer_AsyncContext_AliceAndBob()
	host.Runtime().InitStateFromContractCallInput(originalVMInput)
	async := NewAsyncContext(host)
	asyncCall := &arwen.AsyncCall{
		Destination: Bob,
		ValueBytes:  big.NewInt(88).Bytes(),
	}

	asyncCall.Data = []byte{}
	input, err := async.createContractCallInput(asyncCall)
	require.Nil(t, input)
	require.Error(t, err)

	asyncCall.Data = []byte("function")
	asyncCall.GasLimit = 1
	input, err = async.createContractCallInput(asyncCall)
	require.Nil(t, input)
	require.True(t, errors.Is(err, arwen.ErrNotEnoughGas))

	asyncCall.Data = []byte("function@0A0B0C@03")
	asyncCall.GasLimit = 2
	input, err = async.createContractCallInput(asyncCall)
	require.Nil(t, err)
	require.NotNil(t, input)

	expectedInput := defaultCallInput_AliceToBob(originalVMInput)
	expectedInput.GasProvided = 1
	require.Equal(t, expectedInput, input)
}

func TestAsyncContext_CreateCallbackInput_DestinationCallSuccessful(t *testing.T) {
	host, _, originalVMInput := InitializeArwenAndWasmer_AsyncContext_AliceAndBob()
	host.Runtime().InitStateFromContractCallInput(originalVMInput)
	async := NewAsyncContext(host)

	asyncCall := defaultAsyncCall_AliceToBob()
	asyncCall.Status = arwen.AsyncCallResolved
	asyncCall.GasLocked = 82

	vmOutput := defaultDestOutput_Ok()
	vmOutput.GasRemaining = 12

	destinationErr := error(nil)
	callbackInput, err := async.createCallbackInput(asyncCall, vmOutput, destinationErr)
	require.Nil(t, err)

	expectedGasProvided := asyncCall.GasLocked + vmOutput.GasRemaining
	expectedGasProvided -= defaultOutputDataLengthAsArgs(asyncCall, vmOutput)
	expectedGasProvided -= host.Metering().GasSchedule().ElrondAPICost.AsyncCallStep

	expectedInput := defaultCallbackInput_BobToAlice(originalVMInput)
	expectedInput.GasProvided = expectedGasProvided
	require.Equal(t, expectedInput, callbackInput)
}

func TestAsyncContext_CreateCallbackInput_DestinationCallFailed(t *testing.T) {
	host, _, originalVMInput := InitializeArwenAndWasmer_AsyncContext_AliceAndBob()
	host.Runtime().InitStateFromContractCallInput(originalVMInput)
	async := NewAsyncContext(host)

	asyncCall := defaultAsyncCall_AliceToBob()
	asyncCall.Status = arwen.AsyncCallRejected
	asyncCall.GasLocked = 82

	vmOutput := defaultDestOutput_UserError()
	destinationErr := arwen.ErrSignalError
	callbackInput, err := async.createCallbackInput(asyncCall, vmOutput, destinationErr)
	require.Nil(t, err)

	expectedGasProvided := asyncCall.GasLocked + vmOutput.GasRemaining
	expectedGasProvided -= defaultOutputDataLengthAsArgs(asyncCall, vmOutput)
	expectedGasProvided -= host.Metering().GasSchedule().ElrondAPICost.AsyncCallStep

	expectedInput := arwen.MakeContractCallInput(Bob, Alice, "errorCallback", 0)
	arwen.AddArgument(expectedInput, []byte{byte(vmcommon.UserError)})
	arwen.AddArgument(expectedInput, []byte(vmOutput.ReturnMessage))
	arwen.CopyTxHashes(expectedInput, originalVMInput)
	expectedInput.GasProvided = expectedGasProvided
	expectedInput.CallType = vmcommon.AsynchronousCallBack
	expectedInput.ReturnCallAfterError = true
	require.Equal(t, expectedInput, callbackInput)
}

func TestAsyncContext_CreateCallbackInput_NotEnoughGas(t *testing.T) {
	// Due to dynamic gas locking, this situation should never happen
	host, _, originalVMInput := InitializeArwenAndWasmer_AsyncContext_AliceAndBob()
	host.Runtime().InitStateFromContractCallInput(originalVMInput)
	async := NewAsyncContext(host)

	asyncCall := defaultAsyncCall_AliceToBob()
	asyncCall.Status = arwen.AsyncCallRejected

	vmOutput := &vmcommon.VMOutput{
		ReturnCode:    vmcommon.UserError,
		ReturnData:    [][]byte{},
		ReturnMessage: "there was a user error",
		GasRemaining:  0,
	}

	destinationErr := arwen.ErrSignalError
	callbackInput, err := async.createCallbackInput(asyncCall, vmOutput, destinationErr)
	require.Nil(t, callbackInput)
	require.True(t, errors.Is(err, arwen.ErrNotEnoughGas))
}

func TestAsyncContext_FinishSyncExecution_NilError_NilVMOutput(t *testing.T) {
	host, _, originalVMInput := InitializeArwenAndWasmer_AsyncContext_AliceAndBob()
	host.Runtime().InitStateFromContractCallInput(originalVMInput)
	async := NewAsyncContext(host)
	async.finishAsyncLocalExecution(nil, nil)

	// The expectedOutput must also contain an OutputAccount corresponding to
	// Alice, because of a call to host.Output().GetOutputAccount() in
	// host.Output().GetVMOutput(), which creates and caches an empty account for
	// her.
	expectedOutput := arwen.MakeVMOutput()
	arwen.AddNewOutputAccount(expectedOutput, Alice, 0, nil)

	host.Output().GetOutputAccount(Alice) // TODO matei-p keep?
	require.Equal(t, expectedOutput, host.Output().GetVMOutput())
}

func TestAsyncContext_FinishSyncExecution_Error_NilVMOutput(t *testing.T) {
	host, _, originalVMInput := InitializeArwenAndWasmer_AsyncContext_AliceAndBob()
	host.Runtime().InitStateFromContractCallInput(originalVMInput)
	async := NewAsyncContext(host)

	syncExecErr := arwen.ErrNotEnoughGas
	async.finishAsyncLocalExecution(nil, syncExecErr)

	expectedOutput := arwen.MakeVMOutput()
	expectedOutput.ReturnCode = vmcommon.Ok
	expectedOutput.ReturnMessage = syncExecErr.Error()
	arwen.AddFinishData(expectedOutput, []byte(vmcommon.OutOfGas.String()))
	arwen.AddFinishData(expectedOutput, originalVMInput.CurrentTxHash)

	// The expectedOutput must also contain an OutputAccount corresponding to
	// Alice, because of a call to host.Output().GetOutputAccount() in
	// host.Output().GetVMOutput(), which creates and caches an empty account for
	// her.
	arwen.AddNewOutputAccount(expectedOutput, Alice, 0, nil)

	host.Output().GetOutputAccount(Alice) // TODO matei-p keep?
	require.Equal(t, expectedOutput, host.Output().GetVMOutput())
}

func TestAsyncContext_FinishSyncExecution_ErrorAndVMOutput(t *testing.T) {
	host, _, originalVMInput := InitializeArwenAndWasmer_AsyncContext_AliceAndBob()
	host.Runtime().InitStateFromContractCallInput(originalVMInput)
	async := NewAsyncContext(host)

	syncExecOutput := arwen.MakeVMOutput()
	syncExecOutput.ReturnCode = vmcommon.UserError
	syncExecOutput.ReturnMessage = "user made an error"
	syncExecErr := arwen.ErrSignalError
	async.finishAsyncLocalExecution(syncExecOutput, syncExecErr)

	expectedOutput := arwen.MakeVMOutput()
	expectedOutput.ReturnCode = vmcommon.Ok
	expectedOutput.ReturnMessage = "user made an error"
	arwen.AddFinishData(expectedOutput, []byte(vmcommon.UserError.String()))
	arwen.AddFinishData(expectedOutput, originalVMInput.CurrentTxHash)

	// The expectedOutput must also contain an OutputAccount corresponding to
	// Alice, because of a call to host.Output().GetOutputAccount() in
	// host.Output().GetVMOutput(), which creates and caches an empty account for
	// her.
	arwen.AddNewOutputAccount(expectedOutput, Alice, 0, nil)

	host.Output().GetOutputAccount(Alice) // TODO matei-p keep?
	require.Equal(t, expectedOutput, host.Output().GetVMOutput())
}

func defaultAsyncCall_AliceToBob() *arwen.AsyncCall {
	return &arwen.AsyncCall{
		Destination:     Bob,
		Data:            []byte("function@0A0B0C@03"),
		GasLimit:        0,
		GasLocked:       0,
		ValueBytes:      big.NewInt(88).Bytes(),
		SuccessCallback: "successCallback",
		ErrorCallback:   "errorCallback",
		Status:          arwen.AsyncCallPending,
	}
}

func defaultCallInput_AliceToBob(originalVMInput *vmcommon.ContractCallInput) *vmcommon.ContractCallInput {
	destInput := arwen.MakeContractCallInput(Alice, Bob, "function", 88)
	arwen.CopyTxHashes(destInput, originalVMInput)
	arwen.AddArgument(destInput, []byte{10, 11, 12})
	arwen.AddArgument(destInput, []byte{3})
	destInput.CallType = vmcommon.AsynchronousCall

	return destInput
}

func defaultDestOutput_UserError() *vmcommon.VMOutput {
	return &vmcommon.VMOutput{
		ReturnCode:    vmcommon.UserError,
		ReturnData:    [][]byte{},
		ReturnMessage: "user error occurred",
		GasRemaining:  0,
	}
}

func defaultDestOutput_Ok() *vmcommon.VMOutput {
	return &vmcommon.VMOutput{
		ReturnCode: vmcommon.Ok,
		ReturnData: [][]byte{
			[]byte("first"),
			[]byte("second"),
			{},
			[]byte("third"),
		},
		ReturnMessage: "a message",
		GasRemaining:  0,
	}
}

func defaultCallbackInput_BobToAlice(originalVMInput *vmcommon.ContractCallInput) *vmcommon.ContractCallInput {
	input := arwen.MakeContractCallInput(Bob, Alice, "successCallback", 0)
	arwen.AddArgument(input, big.NewInt(int64(vmcommon.Ok)).Bytes())
	arwen.AddArgument(input, []byte("first"))
	arwen.AddArgument(input, []byte("second"))
	arwen.AddArgument(input, []byte{})
	arwen.AddArgument(input, []byte("third"))
	arwen.CopyTxHashes(input, originalVMInput)
	input.CallType = vmcommon.AsynchronousCallBack
	return input
}

func defaultCallbackOutput_Ok() *vmcommon.VMOutput {
	vmOutput := arwen.MakeVMOutput()
	arwen.AddFinishData(vmOutput, []byte("cbFirst"))
	arwen.AddFinishData(vmOutput, []byte("cbSecond"))

	return vmOutput
}

func defaultOutputDataLengthAsArgs(asyncCall *arwen.AsyncCall, vmOutput *vmcommon.VMOutput) uint64 {
	separator := 1
	hexSize := 2
	returnCode := 1 * hexSize

	dataLength := 0
	if vmOutput.ReturnCode == vmcommon.Ok {
		dataLength += len(asyncCall.SuccessCallback) + separator
		for _, data := range vmOutput.ReturnData {
			dataLength += separator
			dataLength += len(data) * hexSize
		}
	} else {
		dataLength += len(asyncCall.ErrorCallback) + separator + returnCode
		dataLength += separator + len(vmOutput.ReturnMessage)*hexSize
	}

	return uint64(dataLength)
}
