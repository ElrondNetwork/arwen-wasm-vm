package arwen

import (
	"github.com/ElrondNetwork/arwen-wasm-vm/config"
	"github.com/ElrondNetwork/elrond-go/core/vmcommon"
)

// CallbackFunctionName is the name of the default asynchronous callback
// function of a smart contract
const CallbackFunctionName = "callBack"

// ProtectedStoragePrefix is the storage key prefix that will be protected by
// Arwen explicitly, and implicitly by the Elrond node due to '@'; the
// protection can be disabled temporarily by the StorageContext
const ProtectedStoragePrefix = "ARWEN@"

// TimeLockKeyPrefix is the storage key prefix used for timelock-related storage.
const TimeLockKeyPrefix = ProtectedStoragePrefix + "TIMELOCK"

// AsyncDataPrefix is the storage key prefix used for AsyncContext-related storage.
const AsyncDataPrefix = ProtectedStoragePrefix + "ASYNC"

// LegacyAsyncCallGroupID is the AsyncCallGroup identifier reserved for the
// implementation of the legacy asyncCall() EEI function
const LegacyAsyncCallGroupID = "LegacyAsync"

// BreakpointValue encodes Wasmer runtime breakpoint types
type BreakpointValue uint64

const (
	// BreakpointNone signifies the lack of a breakpoint
	BreakpointNone BreakpointValue = iota

	// BreakpointExecutionFailed means that Wasmer must stop immediately due to failure indicated by Arwen
	BreakpointExecutionFailed

	// BreakpointAsyncCall means that Wasmer must stop immediately so Arwen can execute an AsyncCall
	BreakpointAsyncCall

	// BreakpointSignalError means that Wasmer must stop immediately due to a contract-signalled error
	BreakpointSignalError

	// BreakpointOutOfGas means that Wasmer must stop immediately due to gas being exhausted
	BreakpointOutOfGas
)

// AsyncCallExecutionMode encodes the execution modes of an AsyncCall
type AsyncCallExecutionMode uint

const (
	// SyncExecution indicates that the async call can be executed synchronously,
	// with its corresponding callback
	SyncExecution AsyncCallExecutionMode = iota

	// AsyncBuiltinFuncIntraShard indicates that the async call is an intra-shard built in function call
	AsyncBuiltinFuncIntraShard

	// AsyncBuiltinFuncCrossShard indicates that the async call is a cross-shard call to a
	// built-in function, which is executed half in-shard, half cross-shard
	AsyncBuiltinFuncCrossShard

	// AsyncUnknown indicates that the async call cannot be executed locally, and
	// must be forwarded to the destination account
	AsyncUnknown
)

// AsyncCallStatus represents the different status an async call can have
type AsyncCallStatus uint8

const (
	// AsyncCallPending is the status of an async call that awaits complete execution
	AsyncCallPending AsyncCallStatus = iota

	// AsyncCallResolved is the status of an async call that was executed completely and successfully
	AsyncCallResolved

	// AsyncCallRejected is the status of an async call that was executed completely but unsuccessfully
	AsyncCallRejected
)

// CodeDeployInput contains code deploy state, whether it comes from a ContractCreateInput or a ContractCallInput
type CodeDeployInput struct {
	ContractCode         []byte
	ContractCodeMetadata []byte
	ContractAddress      []byte
	CodeDeployerAddress  []byte
}

// VMHostParameters represents the parameters to be passed to VMHost
type VMHostParameters struct {
	VMType                   []byte
	BlockGasLimit            uint64
	GasSchedule              config.GasScheduleMap
	ProtocolBuiltinFunctions vmcommon.FunctionNames
	ElrondProtectedKeyPrefix []byte
	ArwenV2EnableEpoch       uint32
	AheadOfTimeEnableEpoch   uint32
	DynGasLockEnableEpoch    uint32
	UseWarmInstance          bool
}
