package contexts

import (
	"errors"
	"fmt"
	"math/big"
	"testing"

	"github.com/ElrondNetwork/arwen-wasm-vm/arwen"
	"github.com/ElrondNetwork/arwen-wasm-vm/arwen/cryptoapi"
	"github.com/ElrondNetwork/arwen-wasm-vm/arwen/elrondapi"
	"github.com/ElrondNetwork/arwen-wasm-vm/config"
	"github.com/ElrondNetwork/arwen-wasm-vm/crypto"
	contextmock "github.com/ElrondNetwork/arwen-wasm-vm/mock/context"
	worldmock "github.com/ElrondNetwork/arwen-wasm-vm/mock/world"
	"github.com/ElrondNetwork/arwen-wasm-vm/wasmer"
	"github.com/ElrondNetwork/elrond-go/core/vmcommon"
	"github.com/stretchr/testify/require"
)

const WASMPageSize = 65536
const counterWasmCode = "./../../test/contracts/counter/output/counter.wasm"

func MakeAPIImports() *wasmer.Imports {
	imports, _ := elrondapi.ElrondEIImports()
	imports, _ = elrondapi.BigIntImports(imports)
	imports, _ = elrondapi.SmallIntImports(imports)
	imports, _ = cryptoapi.CryptoImports(imports)
	return imports
}

func InitializeArwenAndWasmer() *contextmock.VMHostMock {
	imports := MakeAPIImports()
	_ = wasmer.SetImports(imports, 0)

	gasSchedule := config.MakeGasMapForTests()
	gasCostConfig, _ := config.CreateGasConfig(gasSchedule)
	opcodeCosts := gasCostConfig.WASMOpcodeCost.ToOpcodeCostsArray()
	wasmer.SetOpcodeCosts(&opcodeCosts)

	host := &contextmock.VMHostMock{}
	host.SCAPIMethods = imports

	mockMetering := &contextmock.MeteringContextMock{}
	mockMetering.SetGasSchedule(gasSchedule)
	host.MeteringContext = mockMetering
	host.BlockchainContext, _ = NewBlockchainContext(host, worldmock.NewMockWorld())
	host.OutputContext, _ = NewOutputContext(host)
	host.CryptoHook = crypto.NewVMCrypto()
	return host
}

func TestNewRuntimeContext(t *testing.T) {
	host := InitializeArwenAndWasmer()

	vmType := []byte("type")

	runtimeContext, err := NewRuntimeContext(host, vmType, false)
	require.Nil(t, err)
	require.NotNil(t, runtimeContext)

	require.Equal(t, &vmcommon.VMInput{}, runtimeContext.vmInput)
	require.Equal(t, []byte{}, runtimeContext.scAddress)
	require.Equal(t, "", runtimeContext.callFunction)
	require.Equal(t, false, runtimeContext.readOnly)
	require.Nil(t, runtimeContext.asyncCallInfo)
}

func TestRuntimeContext_InitState(t *testing.T) {
	host := InitializeArwenAndWasmer()

	vmType := []byte("type")

	runtimeContext, err := NewRuntimeContext(host, vmType, false)
	require.Nil(t, err)
	require.NotNil(t, runtimeContext)

	runtimeContext.vmInput = nil
	runtimeContext.scAddress = []byte("some address")
	runtimeContext.callFunction = "a function"
	runtimeContext.readOnly = true
	runtimeContext.asyncCallInfo = &arwen.AsyncCallInfo{}

	runtimeContext.InitState()

	require.Equal(t, &vmcommon.VMInput{}, runtimeContext.vmInput)
	require.Equal(t, []byte{}, runtimeContext.scAddress)
	require.Equal(t, "", runtimeContext.callFunction)
	require.Equal(t, false, runtimeContext.readOnly)
	require.Nil(t, runtimeContext.asyncCallInfo)
}

func TestRuntimeContext_NewWasmerInstance(t *testing.T) {
	host := InitializeArwenAndWasmer()

	vmType := []byte("type")

	runtimeContext, err := NewRuntimeContext(host, vmType, false)
	require.Nil(t, err)

	runtimeContext.SetMaxInstanceCount(1)

	gasLimit := uint64(100000000)
	var dummy []byte
	err = runtimeContext.StartWasmerInstance(dummy, gasLimit, false)
	require.NotNil(t, err)
	require.True(t, errors.Is(err, wasmer.ErrInvalidBytecode))

	gasLimit = uint64(100000000)
	dummy = []byte("contract")
	err = runtimeContext.StartWasmerInstance(dummy, gasLimit, false)
	require.NotNil(t, err)

	path := counterWasmCode
	contractCode := arwen.GetSCCode(path)
	err = runtimeContext.StartWasmerInstance(contractCode, gasLimit, false)
	require.Nil(t, err)
	require.Equal(t, arwen.BreakpointNone, runtimeContext.GetRuntimeBreakpointValue())
}

func TestRuntimeContext_StateSettersAndGetters(t *testing.T) {
	imports := MakeAPIImports()
	host := &contextmock.VMHostMock{}
	host.SCAPIMethods = imports

	vmType := []byte("type")
	runtimeContext, _ := NewRuntimeContext(host, vmType, false)

	arguments := [][]byte{[]byte("argument 1"), []byte("argument 2")}
	vmInput := vmcommon.VMInput{
		CallerAddr: []byte("caller"),
		Arguments:  arguments,
		CallValue:  big.NewInt(0),
	}
	callInput := &vmcommon.ContractCallInput{
		VMInput:       vmInput,
		RecipientAddr: []byte("recipient"),
		Function:      "test function",
	}

	runtimeContext.InitStateFromContractCallInput(callInput)
	require.Equal(t, []byte("caller"), runtimeContext.GetVMInput().CallerAddr)
	require.Equal(t, []byte("recipient"), runtimeContext.GetSCAddress())
	require.Equal(t, "test function", runtimeContext.Function())
	require.Equal(t, vmType, runtimeContext.GetVMType())
	require.Equal(t, arguments, runtimeContext.Arguments())

	vmInput2 := vmcommon.VMInput{
		CallerAddr: []byte("caller2"),
		Arguments:  arguments,
		CallValue:  big.NewInt(0),
	}
	runtimeContext.SetVMInput(&vmInput2)
	require.Equal(t, []byte("caller2"), runtimeContext.GetVMInput().CallerAddr)

	runtimeContext.SetSCAddress([]byte("smartcontract"))
	require.Equal(t, []byte("smartcontract"), runtimeContext.GetSCAddress())
}

func TestRuntimeContext_PushPopInstance(t *testing.T) {
	host := InitializeArwenAndWasmer()

	vmType := []byte("type")
	runtimeContext, _ := NewRuntimeContext(host, vmType, false)
	runtimeContext.SetMaxInstanceCount(1)

	gasLimit := uint64(100000000)
	path := counterWasmCode
	contractCode := arwen.GetSCCode(path)
	err := runtimeContext.StartWasmerInstance(contractCode, gasLimit, false)
	require.Nil(t, err)

	instance := runtimeContext.instance

	runtimeContext.pushInstance()
	runtimeContext.instance = nil
	require.Equal(t, 1, len(runtimeContext.instanceStack))

	runtimeContext.popInstance()
	require.NotNil(t, runtimeContext.instance)
	require.Equal(t, instance, runtimeContext.instance)
	require.Equal(t, 0, len(runtimeContext.instanceStack))

	runtimeContext.pushInstance()
	require.Equal(t, 1, len(runtimeContext.instanceStack))
}

func TestRuntimeContext_PushPopState(t *testing.T) {
	imports := MakeAPIImports()
	host := &contextmock.VMHostMock{}
	host.SCAPIMethods = imports

	vmType := []byte("type")
	runtimeContext, _ := NewRuntimeContext(host, vmType, false)
	runtimeContext.SetMaxInstanceCount(1)

	vmInput := vmcommon.VMInput{
		CallerAddr:  []byte("caller"),
		GasProvided: 1000,
		CallValue:   big.NewInt(0),
		ESDTValue:   big.NewInt(0),
	}

	funcName := "test_func"
	scAddress := []byte("smartcontract")
	input := &vmcommon.ContractCallInput{
		VMInput:       vmInput,
		RecipientAddr: scAddress,
		Function:      funcName,
	}
	runtimeContext.InitStateFromContractCallInput(input)

	runtimeContext.PushState()
	require.Equal(t, 1, len(runtimeContext.stateStack))

	// change state
	runtimeContext.SetSCAddress([]byte("dummy"))
	runtimeContext.SetVMInput(nil)
	runtimeContext.SetReadOnly(true)

	require.Equal(t, []byte("dummy"), runtimeContext.GetSCAddress())
	require.Nil(t, runtimeContext.GetVMInput())
	require.True(t, runtimeContext.ReadOnly())

	runtimeContext.PopSetActiveState()

	//check state was restored correctly
	require.Equal(t, scAddress, runtimeContext.GetSCAddress())
	require.Equal(t, funcName, runtimeContext.Function())
	require.Equal(t, &vmInput, runtimeContext.GetVMInput())
	require.False(t, runtimeContext.ReadOnly())
	require.Nil(t, runtimeContext.Arguments())

	runtimeContext.PushState()
	require.Equal(t, 1, len(runtimeContext.stateStack))

	runtimeContext.PushState()
	require.Equal(t, 2, len(runtimeContext.stateStack))

	runtimeContext.PopDiscard()
	require.Equal(t, 1, len(runtimeContext.stateStack))

	runtimeContext.ClearStateStack()
	require.Equal(t, 0, len(runtimeContext.stateStack))
}

func TestRuntimeContext_Instance(t *testing.T) {
	host := InitializeArwenAndWasmer()

	vmType := []byte("type")
	runtimeContext, _ := NewRuntimeContext(host, vmType, false)
	runtimeContext.SetMaxInstanceCount(1)

	gasLimit := uint64(100000000)
	path := counterWasmCode
	contractCode := arwen.GetSCCode(path)
	err := runtimeContext.StartWasmerInstance(contractCode, gasLimit, false)
	require.Nil(t, err)

	gasPoints := uint64(100)
	runtimeContext.SetPointsUsed(gasPoints)
	require.Equal(t, gasPoints, runtimeContext.GetPointsUsed())

	funcName := "increment"
	input := &vmcommon.ContractCallInput{
		VMInput: vmcommon.VMInput{
			CallValue: big.NewInt(0),
		},
		RecipientAddr: []byte("addr"),
		Function:      funcName,
	}
	runtimeContext.InitStateFromContractCallInput(input)

	f, err := runtimeContext.GetFunctionToCall()
	require.Nil(t, err)
	require.NotNil(t, f)

	input.Function = "func"
	runtimeContext.InitStateFromContractCallInput(input)
	f, err = runtimeContext.GetFunctionToCall()
	require.Equal(t, arwen.ErrFuncNotFound, err)
	require.Nil(t, f)

	initFunc := runtimeContext.GetInitFunction()
	require.NotNil(t, initFunc)

	runtimeContext.CleanWasmerInstance()
	require.Nil(t, runtimeContext.instance)
}

func TestRuntimeContext_Breakpoints(t *testing.T) {
	host := InitializeArwenAndWasmer()

	mockOutput := &contextmock.OutputContextMock{
		OutputAccountMock: NewVMOutputAccount([]byte("address")),
	}
	mockOutput.OutputAccountMock.Code = []byte("code")
	mockOutput.SetReturnMessage("")

	host.OutputContext = mockOutput

	vmType := []byte("type")
	runtimeContext, _ := NewRuntimeContext(host, vmType, false)
	runtimeContext.SetMaxInstanceCount(1)

	gasLimit := uint64(100000000)
	path := counterWasmCode
	contractCode := arwen.GetSCCode(path)
	err := runtimeContext.StartWasmerInstance(contractCode, gasLimit, false)
	require.Nil(t, err)

	// Set and get curent breakpoint value
	require.Equal(t, arwen.BreakpointNone, runtimeContext.GetRuntimeBreakpointValue())
	runtimeContext.SetRuntimeBreakpointValue(arwen.BreakpointOutOfGas)
	require.Equal(t, arwen.BreakpointOutOfGas, runtimeContext.GetRuntimeBreakpointValue())

	runtimeContext.SetRuntimeBreakpointValue(arwen.BreakpointNone)
	require.Equal(t, arwen.BreakpointNone, runtimeContext.GetRuntimeBreakpointValue())

	// Signal user error
	mockOutput.SetReturnCode(vmcommon.Ok)
	mockOutput.SetReturnMessage("")
	runtimeContext.SetRuntimeBreakpointValue(arwen.BreakpointNone)

	runtimeContext.SignalUserError("something happened")
	require.Equal(t, arwen.BreakpointSignalError, runtimeContext.GetRuntimeBreakpointValue())
	require.Equal(t, vmcommon.UserError, mockOutput.ReturnCode())
	require.Equal(t, "something happened", mockOutput.ReturnMessage())

	// Fail execution
	mockOutput.SetReturnCode(vmcommon.Ok)
	mockOutput.SetReturnMessage("")
	runtimeContext.SetRuntimeBreakpointValue(arwen.BreakpointNone)

	runtimeContext.FailExecution(nil)
	require.Equal(t, arwen.BreakpointExecutionFailed, runtimeContext.GetRuntimeBreakpointValue())
	require.Equal(t, vmcommon.ExecutionFailed, mockOutput.ReturnCode())
	require.Equal(t, "execution failed", mockOutput.ReturnMessage())

	mockOutput.SetReturnCode(vmcommon.Ok)
	mockOutput.SetReturnMessage("")
	runtimeContext.SetRuntimeBreakpointValue(arwen.BreakpointNone)
	require.Equal(t, arwen.BreakpointNone, runtimeContext.GetRuntimeBreakpointValue())

	runtimeError := errors.New("runtime error")
	runtimeContext.FailExecution(runtimeError)
	require.Equal(t, arwen.BreakpointExecutionFailed, runtimeContext.GetRuntimeBreakpointValue())
	require.Equal(t, vmcommon.ExecutionFailed, mockOutput.ReturnCode())
	require.Equal(t, runtimeError.Error(), mockOutput.ReturnMessage())
}

func TestRuntimeContext_MemLoadStoreOk(t *testing.T) {
	host := InitializeArwenAndWasmer()

	vmType := []byte("type")
	runtimeContext, _ := NewRuntimeContext(host, vmType, false)
	runtimeContext.SetMaxInstanceCount(1)

	gasLimit := uint64(100000000)
	path := counterWasmCode
	contractCode := arwen.GetSCCode(path)
	err := runtimeContext.StartWasmerInstance(contractCode, gasLimit, false)
	require.Nil(t, err)

	memory := runtimeContext.instance.GetMemory()

	memContents, err := runtimeContext.MemLoad(10, 10)
	require.Nil(t, err)
	require.Equal(t, []byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0}, memContents)

	pageSize := uint32(65536)
	require.Equal(t, 2*pageSize, memory.Length())

	memContents = []byte("test data")
	err = runtimeContext.MemStore(10, memContents)
	require.Nil(t, err)
	require.Equal(t, 2*pageSize, memory.Length())

	memContents, err = runtimeContext.MemLoad(10, 10)
	require.Nil(t, err)
	require.Equal(t, []byte{'t', 'e', 's', 't', ' ', 'd', 'a', 't', 'a', 0}, memContents)
}

func TestRuntimeContext_MemoryIsBlank(t *testing.T) {
	host := InitializeArwenAndWasmer()

	vmType := []byte("type")
	runtimeContext, _ := NewRuntimeContext(host, vmType, false)
	runtimeContext.SetMaxInstanceCount(1)

	gasLimit := uint64(100000000)
	path := "./../../test/contracts/init-simple/output/init-simple.wasm"
	contractCode := arwen.GetSCCode(path)
	err := runtimeContext.StartWasmerInstance(contractCode, gasLimit, false)
	require.Nil(t, err)

	memory := runtimeContext.instance.GetMemory()
	err = memory.Grow(30)
	require.Nil(t, err)

	totalPages := 32
	memoryContents := memory.Data()
	require.Equal(t, memory.Length(), uint32(len(memoryContents)))
	require.Equal(t, totalPages*WASMPageSize, len(memoryContents))

	for i, value := range memoryContents {
		if value != byte(0) {
			msg := fmt.Sprintf("Non-zero value found at %d in Wasmer memory: 0x%X", i, value)
			require.Fail(t, msg)
		}
	}
}

func TestRuntimeContext_MemLoadCases(t *testing.T) {
	host := InitializeArwenAndWasmer()

	vmType := []byte("type")
	runtimeContext, _ := NewRuntimeContext(host, vmType, false)
	runtimeContext.SetMaxInstanceCount(1)

	gasLimit := uint64(100000000)
	path := counterWasmCode
	contractCode := arwen.GetSCCode(path)
	err := runtimeContext.StartWasmerInstance(contractCode, gasLimit, false)
	require.Nil(t, err)

	memory := runtimeContext.instance.GetMemory()

	var offset int32
	var length int32
	// Offset too small
	offset = -3
	length = 10
	memContents, err := runtimeContext.MemLoad(offset, length)
	require.True(t, errors.Is(err, arwen.ErrBadBounds))
	require.Nil(t, memContents)

	// Offset too larget
	offset = int32(memory.Length() + 1)
	length = 10
	memContents, err = runtimeContext.MemLoad(offset, length)
	require.True(t, errors.Is(err, arwen.ErrBadBounds))
	require.Nil(t, memContents)

	// Negative length
	offset = 10
	length = -2
	memContents, err = runtimeContext.MemLoad(offset, length)
	require.True(t, errors.Is(err, arwen.ErrNegativeLength))
	require.Nil(t, memContents)

	// Requested end too large
	memContents = []byte("test data")
	offset = int32(memory.Length() - 9)
	err = runtimeContext.MemStore(offset, memContents)
	require.Nil(t, err)

	offset = int32(memory.Length() - 9)
	length = 9
	memContents, err = runtimeContext.MemLoad(offset, length)
	require.Nil(t, err)
	require.Equal(t, []byte("test data"), memContents)

	offset = int32(memory.Length() - 8)
	length = 9
	memContents, err = runtimeContext.MemLoad(offset, length)
	require.Nil(t, err)
	require.Equal(t, []byte{'e', 's', 't', ' ', 'd', 'a', 't', 'a', 0}, memContents)

	// Zero length
	offset = int32(memory.Length() - 8)
	length = 0
	memContents, err = runtimeContext.MemLoad(offset, length)
	require.Nil(t, err)
	require.Equal(t, []byte{}, memContents)
}

func TestRuntimeContext_MemStoreCases(t *testing.T) {
	host := InitializeArwenAndWasmer()

	vmType := []byte("type")
	runtimeContext, _ := NewRuntimeContext(host, vmType, false)
	runtimeContext.SetMaxInstanceCount(1)

	gasLimit := uint64(100000000)
	path := counterWasmCode
	contractCode := arwen.GetSCCode(path)
	err := runtimeContext.StartWasmerInstance(contractCode, gasLimit, false)
	require.Nil(t, err)

	pageSize := uint32(65536)
	memory := runtimeContext.instance.GetMemory()
	require.Equal(t, 2*pageSize, memory.Length())

	// Bad lower bounds
	memContents := []byte("test data")
	offset := int32(-2)
	err = runtimeContext.MemStore(offset, memContents)
	require.True(t, errors.Is(err, arwen.ErrBadLowerBounds))

	// Memory growth
	require.Equal(t, 2*pageSize, memory.Length())
	offset = int32(memory.Length() - 4)
	err = runtimeContext.MemStore(offset, memContents)
	require.Nil(t, err)
	require.Equal(t, 3*pageSize, memory.Length())

	// Bad upper bounds - forcing the Wasmer memory to grow more than a page at a
	// time is not allowed
	memContents = make([]byte, pageSize+100)
	offset = int32(memory.Length() - 50)
	err = runtimeContext.MemStore(offset, memContents)
	require.True(t, errors.Is(err, arwen.ErrBadUpperBounds))
	require.Equal(t, 4*pageSize, memory.Length())

	// Write something, then overwrite, then overwrite with empty byte slice
	memContents = []byte("this is a message")
	offset = int32(memory.Length() - 100)
	err = runtimeContext.MemStore(offset, memContents)
	require.Nil(t, err)

	memContents, err = runtimeContext.MemLoad(offset, 17)
	require.Nil(t, err)
	require.Equal(t, []byte("this is a message"), memContents)

	memContents = []byte("this is something")
	err = runtimeContext.MemStore(offset, memContents)
	require.Nil(t, err)

	memContents, err = runtimeContext.MemLoad(offset, 17)
	require.Nil(t, err)
	require.Equal(t, []byte("this is something"), memContents)

	memContents = []byte{}
	err = runtimeContext.MemStore(offset, memContents)
	require.Nil(t, err)

	memContents, err = runtimeContext.MemLoad(offset, 17)
	require.Nil(t, err)
	require.Equal(t, []byte("this is something"), memContents)
}

func TestRuntimeContext_MemLoadStoreVsInstanceStack(t *testing.T) {
	host := InitializeArwenAndWasmer()

	vmType := []byte("type")
	runtimeContext, _ := NewRuntimeContext(host, vmType, false)
	runtimeContext.SetMaxInstanceCount(2)

	gasLimit := uint64(100000000)
	path := counterWasmCode
	contractCode := arwen.GetSCCode(path)
	err := runtimeContext.StartWasmerInstance(contractCode, gasLimit, false)
	require.Nil(t, err)

	// Write "test data1" to the WASM memory of the current instance
	memContents := []byte("test data1")
	err = runtimeContext.MemStore(10, memContents)
	require.Nil(t, err)

	memContents, err = runtimeContext.MemLoad(10, 10)
	require.Nil(t, err)
	require.Equal(t, []byte("test data1"), memContents)

	// Push the current instance down the instance stack
	runtimeContext.pushInstance()
	require.Equal(t, 1, len(runtimeContext.instanceStack))

	// Create a new Wasmer instance
	contractCode = arwen.GetSCCode(path)
	err = runtimeContext.StartWasmerInstance(contractCode, gasLimit, false)
	require.Nil(t, err)

	// Write "test data2" to the WASM memory of the new instance
	memContents = []byte("test data2")
	err = runtimeContext.MemStore(10, memContents)
	require.Nil(t, err)

	memContents, err = runtimeContext.MemLoad(10, 10)
	require.Nil(t, err)
	require.Equal(t, []byte("test data2"), memContents)

	// Pop the initial instance from the stack, making it the 'current instance'
	runtimeContext.popInstance()
	require.Equal(t, 0, len(runtimeContext.instanceStack))

	// Check whether the previously-written string "test data1" is still in the
	// memory of the initial instance
	memContents, err = runtimeContext.MemLoad(10, 10)
	require.Nil(t, err)
	require.Equal(t, []byte("test data1"), memContents)

	// Write "test data3" to the WASM memory of the initial instance (now current)
	memContents = []byte("test data3")
	err = runtimeContext.MemStore(10, memContents)
	require.Nil(t, err)

	memContents, err = runtimeContext.MemLoad(10, 10)
	require.Nil(t, err)
	require.Equal(t, []byte("test data3"), memContents)
}

func TestRuntimeContext_PopSetActiveStateIfStackIsEmptyShouldNotPanic(t *testing.T) {
	t.Parallel()

	host := InitializeArwenAndWasmer()

	vmType := []byte("type")
	runtimeContext, _ := NewRuntimeContext(host, vmType, false)
	runtimeContext.PopSetActiveState()

	require.Equal(t, 0, len(runtimeContext.stateStack))
}

func TestRuntimeContext_PopDiscardIfStackIsEmptyShouldNotPanic(t *testing.T) {
	t.Parallel()

	host := InitializeArwenAndWasmer()

	vmType := []byte("type")
	runtimeContext, _ := NewRuntimeContext(host, vmType, false)
	runtimeContext.PopDiscard()

	require.Equal(t, 0, len(runtimeContext.stateStack))
}

func TestRuntimeContext_PopInstanceIfStackIsEmptyShouldNotPanic(t *testing.T) {
	t.Parallel()

	host := InitializeArwenAndWasmer()

	vmType := []byte("type")
	runtimeContext, _ := NewRuntimeContext(host, vmType, false)
	runtimeContext.popInstance()

	require.Equal(t, 0, len(runtimeContext.stateStack))
}
