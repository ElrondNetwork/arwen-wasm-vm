package arwen

import (
	"math/big"

	"github.com/ElrondNetwork/arwen-wasm-vm/v1_3/config"
	"github.com/ElrondNetwork/arwen-wasm-vm/v1_3/crypto"
	"github.com/ElrondNetwork/arwen-wasm-vm/v1_3/wasmer"
	vmcommon "github.com/ElrondNetwork/elrond-vm-common"
	"github.com/ElrondNetwork/elrond-vm-common/data/esdt"
)

// StateStack defines the functionality for working with a state stack
type StateStack interface {
	InitState()
	PushState()
	PopSetActiveState()
	PopDiscard()
	ClearStateStack()
}

// CallArgsParser defines the functionality to parse transaction data for a smart contract call
type CallArgsParser interface {
	ParseData(data string) (string, [][]byte, error)
	IsInterfaceNil() bool
}

// VMHost defines the functionality for working with the VM
type VMHost interface {
	vmcommon.VMExecutionHandler
	Crypto() crypto.VMCrypto
	Blockchain() BlockchainContext
	Runtime() RuntimeContext
	Async() AsyncContext
	BigInt() BigIntContext
	Output() OutputContext
	Metering() MeteringContext
	Storage() StorageContext
	IsArwenV2Enabled() bool
	IsAheadOfTimeCompileEnabled() bool
	IsDynamicGasLockingEnabled() bool
	IsArwenV3Enabled() bool
	IsESDTFunctionsEnabled() bool

	ExecuteESDTTransfer(destination []byte, sender []byte, tokenIdentifier []byte, nonce uint64, value *big.Int, callType vmcommon.CallType) (*vmcommon.VMOutput, uint64, error)
	CreateNewContract(input *vmcommon.ContractCreateInput) ([]byte, error)
	ExecuteOnSameContext(input *vmcommon.ContractCallInput) error
	ExecuteOnDestContext(input *vmcommon.ContractCallInput) (*vmcommon.VMOutput, error)
	GetAPIMethods() *wasmer.Imports
	GetProtocolBuiltinFunctions() vmcommon.FunctionNames
	SetProtocolBuiltinFunctions(vmcommon.FunctionNames)
	IsBuiltinFunctionName(functionName string) bool
	AreInSameShard(leftAddress []byte, rightAddress []byte) bool

	GetGasScheduleMap() config.GasScheduleMap
	GetContexts() (BigIntContext, BlockchainContext, MeteringContext, OutputContext, RuntimeContext, AsyncContext, StorageContext)
	SetRuntimeContext(runtime RuntimeContext)

	CallArgsParser() CallArgsParser

	InitState()
}

// BlockchainContext defines the functionality needed for interacting with the blockchain context
type BlockchainContext interface {
	StateStack

	NewAddress(creatorAddress []byte) ([]byte, error)
	AccountExists(addr []byte) bool
	GetBalance(addr []byte) []byte
	GetBalanceBigInt(addr []byte) *big.Int
	GetNonce(addr []byte) (uint64, error)
	CurrentEpoch() uint32
	GetStateRootHash() []byte
	LastTimeStamp() uint64
	LastNonce() uint64
	LastRound() uint64
	LastEpoch() uint32
	CurrentRound() uint64
	CurrentNonce() uint64
	CurrentTimeStamp() uint64
	CurrentRandomSeed() []byte
	LastRandomSeed() []byte
	IncreaseNonce(addr []byte)
	GetCodeHash(addr []byte) []byte
	GetCode(addr []byte) ([]byte, error)
	GetCodeSize(addr []byte) (int32, error)
	BlockHash(number uint64) []byte
	GetOwnerAddress() ([]byte, error)
	GetShardOfAddress(addr []byte) uint32
	IsSmartContract(addr []byte) bool
	IsPayable(address []byte) (bool, error)
	SaveCompiledCode(codeHash []byte, code []byte)
	GetCompiledCode(codeHash []byte) (bool, []byte)
	GetESDTToken(address []byte, tokenID []byte, nonce uint64) (*esdt.ESDigitalToken, error)
	GetUserAccount(address []byte) (vmcommon.UserAccountHandler, error)
	ProcessBuiltInFunction(input *vmcommon.ContractCallInput) (*vmcommon.VMOutput, error)
	GetSnapshot() int
	RevertToSnapshot(snapshot int)
}

// RuntimeContext defines the functionality needed for interacting with the runtime context
type RuntimeContext interface {
	StateStack

	InitStateFromContractCallInput(input *vmcommon.ContractCallInput)
	SetCustomCallFunction(callFunction string)
	GetVMInput() *vmcommon.VMInput
	SetVMInput(vmInput *vmcommon.VMInput)
	GetSCAddress() []byte
	SetSCAddress(scAddress []byte)
	GetSCCode() ([]byte, error)
	GetSCCodeSize() uint64
	GetVMType() []byte
	Function() string
	Arguments() [][]byte
	GetCurrentTxHash() []byte
	GetOriginalTxHash() []byte
	ExtractCodeUpgradeFromArgs() ([]byte, []byte, error)
	SignalUserError(message string)
	FailExecution(err error)
	MustVerifyNextContractCode()
	SetRuntimeBreakpointValue(value BreakpointValue)
	GetRuntimeBreakpointValue() BreakpointValue
	IsContractOnTheStack(address []byte) bool
	RunningInstancesCount() uint64
	IsFunctionImported(name string) bool
	IsWarmInstance() bool
	ResetWarmInstance()
	ReadOnly() bool
	SetReadOnly(readOnly bool)
	StartWasmerInstance(contract []byte, gasLimit uint64, newCode bool) error
	CleanWasmerInstance()
	SetMaxInstanceCount(uint64)
	VerifyContractCode() error
	GetInstance() wasmer.InstanceHandler
	GetInstanceExports() wasmer.ExportsMap
	GetInitFunction() wasmer.ExportedFunctionCallback
	GetFunctionToCall() (wasmer.ExportedFunctionCallback, error)
	GetPointsUsed() uint64
	SetPointsUsed(gasPoints uint64)
	MemStore(offset int32, data []byte) error
	MemLoad(offset int32, length int32) ([]byte, error)
	MemLoadMultiple(offset int32, lengths []int32) ([][]byte, error)
	ElrondAPIErrorShouldFailExecution() bool
	ElrondSyncExecAPIErrorShouldFailExecution() bool
	CryptoAPIErrorShouldFailExecution() bool
	BigIntAPIErrorShouldFailExecution() bool

	AddError(err error, otherInfo ...string)
	GetAllErrors() error

	ValidateCallbackName(callbackName string) error
	HasFunction(functionName string) bool
	GetPrevTxHash() []byte

	// TODO remove after implementing proper mocking of Wasmer instances; this is
	// used for tests only
	ReplaceInstanceBuilder(builder InstanceBuilder)
}

// BigIntContext defines the functionality needed for interacting with the big int context
type BigIntContext interface {
	StateStack

	Put(value int64) int32
	GetOne(id int32) *big.Int
	GetTwo(id1, id2 int32) (*big.Int, *big.Int)
	GetThree(id1, id2, id3 int32) (*big.Int, *big.Int, *big.Int)
}

// OutputContext defines the functionality needed for interacting with the output context
type OutputContext interface {
	StateStack
	PopMergeActiveState()
	CensorVMOutput()
	AddToActiveState(rightOutput *vmcommon.VMOutput)

	GetOutputAccount(address []byte) (*vmcommon.OutputAccount, bool)
	GetOutputAccounts() map[string]*vmcommon.OutputAccount
	DeleteOutputAccount(address []byte)
	WriteLog(address []byte, topics [][]byte, data []byte)
	TransferValueOnly(destination []byte, sender []byte, value *big.Int, checkPayable bool) error
	Transfer(destination []byte, sender []byte, gasLimit uint64, gasLocked uint64, value *big.Int, input []byte, callType vmcommon.CallType) error
	TransferESDT(destination []byte, sender []byte, tokenIdentifier []byte, nonce uint64, value *big.Int, callInput *vmcommon.ContractCallInput) (uint64, error)
	SelfDestruct(address []byte, beneficiary []byte)
	GetRefund() uint64
	SetRefund(refund uint64)
	ReturnCode() vmcommon.ReturnCode
	SetReturnCode(returnCode vmcommon.ReturnCode)
	ReturnMessage() string
	SetReturnMessage(message string)
	ReturnData() [][]byte
	ClearReturnData()
	Finish(data []byte)
	PrependFinish(data []byte)
	GetVMOutput() *vmcommon.VMOutput
	AddTxValueToAccount(address []byte, value *big.Int)
	DeployCode(input CodeDeployInput)
	CreateVMOutputInCaseOfError(err error) *vmcommon.VMOutput
}

// MeteringContext defines the functionality needed for interacting with the metering context
type MeteringContext interface {
	StateStack
	PopMergeActiveState()

	InitStateFromContractCallInput(input *vmcommon.VMInput)
	SetGasSchedule(gasMap config.GasScheduleMap)
	GasSchedule() *config.GasCost
	UseGas(gas uint64)
	FreeGas(gas uint64)
	RestoreGas(gas uint64)
	GasLeft() uint64
	GasUsedForExecution() uint64
	GasSpentByContract() uint64
	GetGasForExecution() uint64
	GetGasProvided() uint64
	GetSCPrepareInitialCost() uint64
	BoundGasLimit(value int64) uint64
	BlockGasLimit() uint64
	DeductInitialGasForExecution(contract []byte) error
	DeductInitialGasForDirectDeployment(input CodeDeployInput) error
	DeductInitialGasForIndirectDeployment(input CodeDeployInput) error
	ComputeGasLockedForAsync() uint64
	UseGasForAsyncStep() error
	UseGasBounded(gasToUse uint64) error
	GetGasLocked() uint64
	UpdateGasStateOnSuccess(vmOutput *vmcommon.VMOutput) error
	UpdateGasStateOnFailure(vmOutput *vmcommon.VMOutput)
	TrackGasUsedByBuiltinFunction(builtinInput *vmcommon.ContractCallInput, builtinOutput *vmcommon.VMOutput, postBuiltinInput *vmcommon.ContractCallInput)
}

// StorageStatus defines the states the storage can be in
type StorageStatus int

const (
	// StorageUnchanged signals that the storage was not changed
	StorageUnchanged StorageStatus = iota

	// StorageModified signals that the storage has been modified
	StorageModified

	// StorageAdded signals that something was added to storage
	StorageAdded

	// StorageDeleted signals that something was removed from storage
	StorageDeleted
)

// StorageContext defines the functionality needed for interacting with the storage context
type StorageContext interface {
	StateStack

	SetAddress(address []byte)
	GetStorageUpdates(address []byte) map[string]*vmcommon.StorageUpdate
	GetStorageFromAddress(address []byte, key []byte) []byte
	GetStorage(key []byte) []byte
	GetStorageUnmetered(key []byte) []byte
	SetStorage(key []byte, value []byte) (StorageStatus, error)
	SetProtectedStorage(key []byte, value []byte) (StorageStatus, error)
}

// AsyncCallInfoHandler defines the functionality for working with AsyncCallInfo
type AsyncCallInfoHandler interface {
	GetDestination() []byte
	GetData() []byte
	GetGasLimit() uint64
	GetGasLocked() uint64
	GetValueBytes() []byte
}

// InstanceBuilder defines the functionality needed to create Wasmer instances
type InstanceBuilder interface {
	NewInstanceWithOptions(contractCode []byte, options wasmer.CompilationOptions) (wasmer.InstanceHandler, error)
	NewInstanceFromCompiledCodeWithOptions(compiledCode []byte, options wasmer.CompilationOptions) (wasmer.InstanceHandler, error)
}

type AsyncContext interface {
	StateStack

	InitStateFromInput(input *vmcommon.ContractCallInput)
	HasPendingCallGroups() bool
	IsComplete() bool
	GetCallGroup(groupID string) (*AsyncCallGroup, bool)
	SetGroupCallback(groupID string, callbackName string, data []byte, gas uint64) error
	SetContextCallback(callbackName string, data []byte, gas uint64) error
	HasCallback() bool
	PostprocessCrossShardCallback() error
	GetCallerAddress() []byte
	GetReturnData() []byte
	SetReturnData(data []byte)
	GetGasPrice() uint64

	Execute() error
	RegisterAsyncCall(groupID string, call *AsyncCall) error
	RegisterLegacyAsyncCall(address []byte, data []byte, value []byte) error
	UpdateCurrentCallStatus() (*AsyncCall, error)

	Load() error
	Save() error
	Delete() error
}
