package mock

import (
	"github.com/ElrondNetwork/arwen-wasm-vm/arwen"
	"github.com/ElrondNetwork/arwen-wasm-vm/wasmer"
	"github.com/ElrondNetwork/elrond-go/core/vmcommon"
)

var _ arwen.RuntimeContext = (*RuntimeContextMock)(nil)

// RuntimeContextMock is used in tests to check the RuntimeContextMock interface method calls
type RuntimeContextMock struct {
	Err                    error
	VMInput                *vmcommon.VMInput
	SCAddress              []byte
	SCCode                 []byte
	SCCodeSize             uint64
	CallFunction           string
	VMType                 []byte
	IsContractOnStack      bool
	ReadOnlyFlag           bool
	VerifyCode             bool
	CurrentBreakpointValue arwen.BreakpointValue
	PointsUsed             uint64
	InstanceCtxID          int
	MemLoadResult          []byte
	MemLoadMultipleResult  [][]byte
	FailCryptoAPI          bool
	FailElrondAPI          bool
	FailElrondSyncExecAPI  bool
	FailBigIntAPI          bool
	DefaultAsyncCall       *arwen.AsyncCall
	RunningInstances       uint64
	CurrentTxHash          []byte
	OriginalTxHash         []byte
	PrevTxHash             []byte
	HasFunctionResult      bool
}

// InitState mocked method
func (r *RuntimeContextMock) InitState() {
}

// ReplaceInstanceBuilder mocked method()
func (r *RuntimeContextMock) ReplaceInstanceBuilder(_ arwen.InstanceBuilder) {
}

// StartWasmerInstance mocked method
func (r *RuntimeContextMock) StartWasmerInstance(_ []byte, _ uint64, _ bool) error {
	if r.Err != nil {
		return r.Err
	}
	return nil
}

// InitStateFromContractCallInput mocked method
func (r *RuntimeContextMock) InitStateFromContractCallInput(_ *vmcommon.ContractCallInput) {
}

// PushState mocked method
func (r *RuntimeContextMock) PushState() {
}

// PopSetActiveState mocked method
func (r *RuntimeContextMock) PopSetActiveState() {
}

// PopDiscard mocked method
func (r *RuntimeContextMock) PopDiscard() {
}

// MustVerifyNextContractCode mocked method
func (r *RuntimeContextMock) MustVerifyNextContractCode() {
}

// ClearStateStack mocked method
func (r *RuntimeContextMock) ClearStateStack() {
}

// PushInstance mocked method
func (r *RuntimeContextMock) PushInstance() {
}

// PopInstance mocked method
func (r *RuntimeContextMock) PopInstance() {
}

// IsWarmInstance mocked method
func (r *RuntimeContextMock) IsWarmInstance() bool {
	return false
}

// ResetWarmInstance mocked method
func (r *RuntimeContextMock) ResetWarmInstance() {
}

// RunningInstancesCount mocked method
func (r *RuntimeContextMock) RunningInstancesCount() uint64 {
	return r.RunningInstances
}

// SetMaxInstanceCount mocked method
func (r *RuntimeContextMock) SetMaxInstanceCount(uint64) {
}

// ClearInstanceStack mocked method
func (r *RuntimeContextMock) ClearInstanceStack() {
}

func (r *RuntimeContextMock) ValidateCallbackName(callbackName string) error {
	return nil
}

// GetVMType mocked method
func (r *RuntimeContextMock) GetVMType() []byte {
	return r.VMType
}

// GetVMInput mocked method
func (r *RuntimeContextMock) GetVMInput() *vmcommon.VMInput {
	return r.VMInput
}

// SetVMInput mocked method
func (r *RuntimeContextMock) SetVMInput(vmInput *vmcommon.VMInput) {
	r.VMInput = vmInput
}

// IsContractOnTheStack mocked method
func (r *RuntimeContextMock) IsContractOnTheStack(_ []byte) bool {
	return r.IsContractOnStack
}

// GetSCAddress mocked method
func (r *RuntimeContextMock) GetSCAddress() []byte {
	return r.SCAddress
}

// SetSCAddress mocked method
func (r *RuntimeContextMock) SetSCAddress(scAddress []byte) {
	r.SCAddress = scAddress
}

// GetSCCode mocked method
func (r *RuntimeContextMock) GetSCCode() ([]byte, error) {
	return r.SCCode, r.Err
}

// GetSCCodeSize mocked method
func (r *RuntimeContextMock) GetSCCodeSize() uint64 {
	return r.SCCodeSize
}

// Function mocked method
func (r *RuntimeContextMock) Function() string {
	return r.CallFunction
}

// Arguments mocked method
func (r *RuntimeContextMock) Arguments() [][]byte {
	return r.VMInput.Arguments
}

// GetCurrentTxHash mocked method
func (r *RuntimeContextMock) GetCurrentTxHash() []byte {
	return r.CurrentTxHash
}

// GetOriginalTxHash mocked method
func (r *RuntimeContextMock) GetOriginalTxHash() []byte {
	return r.OriginalTxHash
}

func (r *RuntimeContextMock) GetPrevTxHash() []byte {
	return r.PrevTxHash
}

// ExtractCodeUpgradeFromArgs mocked method
func (r *RuntimeContextMock) ExtractCodeUpgradeFromArgs() ([]byte, []byte, error) {
	arguments := r.VMInput.Arguments
	if len(arguments) < 2 {
		panic("ExtractCodeUpgradeFromArgs: bad test setup")
	}

	return r.VMInput.Arguments[0], r.VMInput.Arguments[1], nil
}

// SignalExit mocked method
func (r *RuntimeContextMock) SignalExit(_ int) {
}

// SignalUserError mocked method
func (r *RuntimeContextMock) SignalUserError(_ string) {
}

// SetRuntimeBreakpointValue mocked method
func (r *RuntimeContextMock) SetRuntimeBreakpointValue(_ arwen.BreakpointValue) {
}

// GetRuntimeBreakpointValue mocked method
func (r *RuntimeContextMock) GetRuntimeBreakpointValue() arwen.BreakpointValue {
	return r.CurrentBreakpointValue
}

// PrepareLegacyAsyncCall mocked method
func (r *RuntimeContextMock) PrepareLegacyAsyncCall(address []byte, data []byte, value []byte) error {
	return r.Err
}

// VerifyContractCode mocked method
func (r *RuntimeContextMock) VerifyContractCode() error {
	return r.Err
}

// GetPointsUsed mocked method
func (r *RuntimeContextMock) GetPointsUsed() uint64 {
	return r.PointsUsed
}

// SetPointsUsed mocked method
func (r *RuntimeContextMock) SetPointsUsed(gasPoints uint64) {
	r.PointsUsed = gasPoints
}

// ReadOnly mocked method
func (r *RuntimeContextMock) ReadOnly() bool {
	return r.ReadOnlyFlag
}

// SetReadOnly mocked method
func (r *RuntimeContextMock) SetReadOnly(readOnly bool) {
	r.ReadOnlyFlag = readOnly
}

// GetInstanceExports mocked method
func (r *RuntimeContextMock) GetInstanceExports() wasmer.ExportsMap {
	return nil
}

// CleanWasmerInstance mocked method
func (r *RuntimeContextMock) CleanWasmerInstance() {
}

// GetFunctionToCall mocked method
func (r *RuntimeContextMock) GetFunctionToCall() (wasmer.ExportedFunctionCallback, error) {
	if r.Err != nil {
		return nil, r.Err
	}
	return nil, nil
}

// GetInitFunction mocked method
func (r *RuntimeContextMock) GetInitFunction() wasmer.ExportedFunctionCallback {
	return nil
}

// MemLoad mocked method
func (r *RuntimeContextMock) MemLoad(_ int32, _ int32) ([]byte, error) {
	if r.Err != nil {
		return nil, r.Err
	}

	return r.MemLoadResult, nil
}

// MemLoadMultiple mocked method
func (r *RuntimeContextMock) MemLoadMultiple(_ int32, _ []int32) ([][]byte, error) {
	if r.Err != nil {
		return nil, r.Err
	}

	return r.MemLoadMultipleResult, nil
}

// MemStore mocked method
func (r *RuntimeContextMock) MemStore(_ int32, _ []byte) error {
	return r.Err
}

// ElrondAPIErrorShouldFailExecution mocked method
func (r *RuntimeContextMock) ElrondAPIErrorShouldFailExecution() bool {
	return r.FailElrondAPI
}

// ElrondSyncExecAPIErrorShouldFailExecution mocked method
func (r *RuntimeContextMock) ElrondSyncExecAPIErrorShouldFailExecution() bool {
	return r.FailElrondSyncExecAPI
}

// CryptoAPIErrorShouldFailExecution mocked method
func (r *RuntimeContextMock) CryptoAPIErrorShouldFailExecution() bool {
	return r.FailCryptoAPI
}

// BigIntAPIErrorShouldFailExecution mocked method
func (r *RuntimeContextMock) BigIntAPIErrorShouldFailExecution() bool {
	return r.FailBigIntAPI
}

// FailExecution mocked method
func (r *RuntimeContextMock) FailExecution(_ error) {
}

// SetCustomCallFunction mocked method
func (r *RuntimeContextMock) SetCustomCallFunction(_ string) {
}

// HasCallbackMethod mocked method
func (r *RuntimeContextMock) HasFunction(functionName string) bool {
	return r.HasFunctionResult
}
