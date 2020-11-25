package contexts

import (
	"bytes"
	"fmt"
	"math/big"

	"github.com/ElrondNetwork/arwen-wasm-vm/arwen"
	"github.com/ElrondNetwork/arwen-wasm-vm/wasmer"
	logger "github.com/ElrondNetwork/elrond-go-logger"
	"github.com/ElrondNetwork/elrond-go/core/vmcommon"
)

var log = logger.GetOrCreate("arwen/runtime")

var _ arwen.RuntimeContext = (*runtimeContext)(nil)

type runtimeContext struct {
	host         arwen.VMHost
	instance     *wasmer.Instance
	vmInput      *vmcommon.VMInput
	scAddress    []byte
	codeSize     uint64
	callFunction string
	vmType       []byte
	readOnly     bool

	verifyCode bool

	stateStack    []*runtimeContext
	instanceStack []*wasmer.Instance

	maxWasmerInstances uint64

	validator *wasmValidator

	useWarmInstance     bool
	warmInstanceAddress []byte
	warmInstance        *wasmer.Instance
}

// NewRuntimeContext creates a new runtimeContext.
func NewRuntimeContext(host arwen.VMHost, vmType []byte, useWarmInstance bool) (*runtimeContext, error) {
	scAPINames := host.GetAPIMethods().Names()
	protocolBuiltinFunctions := host.GetProtocolBuiltinFunctions()

	context := &runtimeContext{
		host:                host,
		vmType:              vmType,
		stateStack:          make([]*runtimeContext, 0),
		instanceStack:       make([]*wasmer.Instance, 0),
		validator:           newWASMValidator(scAPINames, protocolBuiltinFunctions),
		useWarmInstance:     useWarmInstance,
		warmInstanceAddress: nil,
		warmInstance:        nil,
	}

	context.InitState()

	return context, nil
}

// InitState initializes all the contexts fields with default data.
func (context *runtimeContext) InitState() {
	context.vmInput = &vmcommon.VMInput{}
	context.scAddress = make([]byte, 0)
	context.callFunction = ""
	context.verifyCode = false
	context.readOnly = false
}

func (context *runtimeContext) setWarmInstanceWhenNeeded(gasLimit uint64) bool {
	scAddress := context.GetSCAddress()
	useWarm := context.useWarmInstance && context.warmInstanceAddress != nil && bytes.Equal(scAddress, context.warmInstanceAddress)
	if scAddress != nil && useWarm {
		log.Trace("Reusing the warm Wasmer instance")

		context.instance = context.warmInstance
		context.SetPointsUsed(0)
		context.instance.SetGasLimit(gasLimit)

		context.SetRuntimeBreakpointValue(arwen.BreakpointNone)
		return true
	}

	return false
}

func (context *runtimeContext) makeInstanceFromCompiledCode(codeHash []byte, gasLimit uint64, newCode bool) bool {
	if !context.host.IsAheadOfTimeCompileEnabled() {
		return false
	}

	if newCode || len(codeHash) == 0 {
		return false
	}

	blockchain := context.host.Blockchain()
	found, compiledCode := blockchain.GetCompiledCode(codeHash)
	if !found {
		log.Debug("compiled code was not found")
		return false
	}

	gasSchedule := context.host.Metering().GasSchedule()
	options := wasmer.CompilationOptions{
		GasLimit:           gasLimit,
		UnmeteredLocals:    uint64(gasSchedule.WASMOpcodeCost.LocalsUnmetered),
		OpcodeTrace:        false,
		Metering:           true,
		RuntimeBreakpoints: true,
	}
	newInstance, err := wasmer.NewInstanceFromCompiledCodeWithOptions(compiledCode, options)
	if err != nil {
		log.Warn("NewInstanceFromCompiledCodeWithOptions", "error", err)
		return false
	}

	context.instance = newInstance

	idContext := arwen.AddHostContext(context.host)
	context.instance.SetContextData(idContext)
	context.verifyCode = false

	return true
}

// StartWasmerInstance initializes a Wasmer instance, either from the provided
// WASM bytecode or from cached precompiled code, if the maxWasmerInstances has
// not been reached.
func (context *runtimeContext) StartWasmerInstance(contract []byte, gasLimit uint64, newCode bool) error {
	if context.RunningInstancesCount() >= context.maxWasmerInstances {
		context.instance = nil
		return arwen.ErrMaxInstancesReached
	}

	warmInstanceUsed := context.setWarmInstanceWhenNeeded(gasLimit)
	if warmInstanceUsed {
		return nil
	}

	blockchain := context.host.Blockchain()
	codeHash := blockchain.GetCodeHash(context.GetSCAddress())
	compiledCodeUsed := context.makeInstanceFromCompiledCode(codeHash, gasLimit, newCode)
	if compiledCodeUsed {
		return nil
	}

	return context.makeInstanceFromContractByteCode(contract, codeHash, gasLimit, newCode)
}

func (context *runtimeContext) makeInstanceFromContractByteCode(contract []byte, codeHash []byte, gasLimit uint64, newCode bool) error {
	log.Trace("Creating a new Wasmer instance")

	gasSchedule := context.host.Metering().GasSchedule()
	options := wasmer.CompilationOptions{
		GasLimit:           gasLimit,
		UnmeteredLocals:    uint64(gasSchedule.WASMOpcodeCost.LocalsUnmetered),
		OpcodeTrace:        false,
		Metering:           true,
		RuntimeBreakpoints: true,
	}
	newInstance, err := wasmer.NewInstanceWithOptions(contract, options)
	if err != nil {
		context.instance = nil
		return err
	}

	context.instance = newInstance

	if newCode || len(codeHash) == 0 {
		codeHash, err = context.host.Crypto().Sha256(contract)
		if err != nil {
			context.CleanWasmerInstance()
			return err
		}
	}

	context.saveCompiledCode(codeHash)

	idContext := arwen.AddHostContext(context.host)
	context.instance.SetContextData(idContext)

	if newCode {
		err = context.VerifyContractCode()
		if err != nil {
			context.CleanWasmerInstance()
			return err
		}
	}

	if context.useWarmInstance {
		context.warmInstanceAddress = context.GetSCAddress()
		context.warmInstance = context.instance
	}

	return nil
}

// GetSCCode returns the WASM bytecode of the current contract, while also caching its size.
func (context *runtimeContext) GetSCCode() ([]byte, error) {
	blockchain := context.host.Blockchain()
	code, err := blockchain.GetCode(context.scAddress)
	if err != nil {
		return nil, err
	}

	context.codeSize = uint64(len(code))
	return code, nil
}

// GetSCCodeSize returns the cached size of the current SC code.
func (context *runtimeContext) GetSCCodeSize() uint64 {
	return context.codeSize
}

func (context *runtimeContext) saveCompiledCode(codeHash []byte) {
	compiledCode, err := context.instance.Cache()
	if err != nil {
		log.Error("getCompiledCode from instance", "error", err)
	}

	blockchain := context.host.Blockchain()
	blockchain.SaveCompiledCode(codeHash, compiledCode)
}

// IsWarmInstance verifies whether the currently running Wasmer instance is the 'warm' instance or not.
func (context *runtimeContext) IsWarmInstance() bool {
	if context.instance != nil && context.instance == context.warmInstance {
		return true
	}

	return false
}

// ResetWarmInstance cleans the current 'warm' instance and closes it, also clearing its related context fields.
func (context *runtimeContext) ResetWarmInstance() {
	if context.instance == nil {
		return
	}

	arwen.RemoveHostContext(*context.instance.Data)
	context.instance.Clean()

	context.instance = nil
	context.warmInstanceAddress = nil
	context.warmInstance = nil
}

// MustVerifyNextContractCode will cause the validation of the WASM bytecode before the next Wasmer instance is started.
func (context *runtimeContext) MustVerifyNextContractCode() {
	context.verifyCode = true
}

// SetMaxInstanceCount sets the maximum number of allowed Wasmer instances on
// the instance stack, for recursivity.
func (context *runtimeContext) SetMaxInstanceCount(maxInstances uint64) {
	context.maxWasmerInstances = maxInstances
}

// InitStateFromContractCallInput initializes the state of the runtime context
// (and the async context) from the provided ContractCallInput.
func (context *runtimeContext) InitStateFromInput(input *vmcommon.ContractCallInput) {
	context.SetVMInput(&input.VMInput)
	context.scAddress = input.RecipientAddr
	context.callFunction = input.Function
}

// SetCustomCallFunction sets a custom function to be called next, instead of
// the one specified by the current ContractCallInput.
func (context *runtimeContext) SetCustomCallFunction(callFunction string) {
	context.callFunction = callFunction
}

// PushState pushes the current state of the runtime context onto its state stack.
func (context *runtimeContext) PushState() {
	newState := &runtimeContext{
		vmInput:      context.vmInput,
		scAddress:    context.scAddress,
		callFunction: context.callFunction,
		readOnly:     context.readOnly,
	}

	context.stateStack = append(context.stateStack, newState)
}

// PopSetActiveState pops the state at the top of the state stack and sets it as the current state.
func (context *runtimeContext) PopSetActiveState() {
	stateStackLen := len(context.stateStack)
	if stateStackLen == 0 {
		return
	}

	prevState := context.stateStack[stateStackLen-1]
	context.stateStack = context.stateStack[:stateStackLen-1]

	context.vmInput = prevState.vmInput
	context.scAddress = prevState.scAddress
	context.callFunction = prevState.callFunction
	context.readOnly = prevState.readOnly
}

// PopDiscard pops the state at the top of the state stack and discards it.
func (context *runtimeContext) PopDiscard() {
	stateStackLen := len(context.stateStack)
	if stateStackLen == 0 {
		return
	}

	context.stateStack = context.stateStack[:stateStackLen-1]
}

// ClearStateStack discards the entire state state stack and initializes it anew.
func (context *runtimeContext) ClearStateStack() {
	context.stateStack = make([]*runtimeContext, 0)
}

// PushInstance pushes the current Wasmer instance on the instance stack (separate from the state stack).
func (context *runtimeContext) PushInstance() {
	context.instanceStack = append(context.instanceStack, context.instance)
}

// PopInstance pops the Wasmer instance off the top of the instance stack, and sets it as the current Wasmer instance.
func (context *runtimeContext) PopInstance() {
	instanceStackLen := len(context.instanceStack)
	if instanceStackLen == 0 {
		return
	}

	prevInstance := context.instanceStack[instanceStackLen-1]
	context.instanceStack = context.instanceStack[:instanceStackLen-1]

	context.CleanWasmerInstance()
	context.instance = prevInstance
}

// RunningInstanceCount returns the number of the currently running Wasmer instances.
func (context *runtimeContext) RunningInstancesCount() uint64 {
	return uint64(len(context.instanceStack))
}

// ClearInstanceStack closes and removes all Wasmer instances from the instance stack
func (context *runtimeContext) ClearInstanceStack() {
	for _, instance := range context.instanceStack {
		instance.Clean()
	}
	context.instanceStack = make([]*wasmer.Instance, 0)
}

// GetVMType returns the bytes that identify the Arwen VM
func (context *runtimeContext) GetVMType() []byte {
	return context.vmType
}

// GetVMInput returns the current VMInput
func (context *runtimeContext) GetVMInput() *vmcommon.VMInput {
	return context.vmInput
}

// SetVMInput sets the current VMInput to the one provided (cloned)
func (context *runtimeContext) SetVMInput(vmInput *vmcommon.VMInput) {
	if !context.host.IsArwenV2Enabled() || vmInput == nil {
		context.vmInput = vmInput
		return
	}

	context.vmInput = &vmcommon.VMInput{
		CallType:      vmInput.CallType,
		GasPrice:      vmInput.GasPrice,
		GasProvided:   vmInput.GasProvided,
		GasLocked:     vmInput.GasLocked,
		CallValue:     big.NewInt(0),
		ESDTValue:     big.NewInt(0),
		ESDTTokenName: nil,
	}

	if vmInput.CallValue != nil {
		context.vmInput.CallValue.Set(vmInput.CallValue)
	}

	if len(vmInput.CallerAddr) > 0 {
		context.vmInput.CallerAddr = make([]byte, len(vmInput.CallerAddr))
		copy(context.vmInput.CallerAddr, vmInput.CallerAddr)
	}

	if vmInput.ESDTValue != nil {
		context.vmInput.ESDTValue.Set(vmInput.ESDTValue)
	}

	if len(vmInput.ESDTTokenName) > 0 {
		context.vmInput.ESDTTokenName = make([]byte, len(vmInput.ESDTTokenName))
		copy(context.vmInput.ESDTTokenName, vmInput.ESDTTokenName)
	}

	if len(vmInput.OriginalTxHash) > 0 {
		context.vmInput.OriginalTxHash = make([]byte, len(vmInput.OriginalTxHash))
		copy(context.vmInput.OriginalTxHash, vmInput.OriginalTxHash)
	}

	if len(vmInput.CurrentTxHash) > 0 {
		context.vmInput.CurrentTxHash = make([]byte, len(vmInput.CurrentTxHash))
		copy(context.vmInput.CurrentTxHash, vmInput.CurrentTxHash)
	}

	if len(vmInput.PrevTxHash) > 0 {
		context.vmInput.PrevTxHash = make([]byte, len(vmInput.PrevTxHash))
		copy(context.vmInput.PrevTxHash, vmInput.PrevTxHash)
	}

	if len(vmInput.Arguments) > 0 {
		context.vmInput.Arguments = make([][]byte, len(vmInput.Arguments))
		for i, arg := range vmInput.Arguments {
			context.vmInput.Arguments[i] = make([]byte, len(arg))
			copy(context.vmInput.Arguments[i], arg)
		}
	}
}

// GetSCAddress returns the address of the contract currently executed
func (context *runtimeContext) GetSCAddress() []byte {
	return context.scAddress
}

// SetSCAddress sets the address of the contract currently executed
func (context *runtimeContext) SetSCAddress(scAddress []byte) {
	context.scAddress = scAddress
}

// GetCurrentTxHash returns the hash of the current transaction, as specified by the current VMInput.
func (context *runtimeContext) GetCurrentTxHash() []byte {
	return context.vmInput.CurrentTxHash
}

// GetCurrentTxHash returns the hash of the original transaction, in the case of async calls, as specified by the current VMInput.
func (context *runtimeContext) GetOriginalTxHash() []byte {
	return context.vmInput.OriginalTxHash
}

// GetPrevTxHash returns the hash of the previous transaction, in the case of async calls, as specified by the current VMInput.
func (context *runtimeContext) GetPrevTxHash() []byte {
	return context.vmInput.PrevTxHash
}

// Function returns the name of the contract function to be called next
func (context *runtimeContext) Function() string {
	return context.callFunction
}

// Arguments returns the binary arguments that will be passed to the contract to be executed, as specified by the current VMInput.
func (context *runtimeContext) Arguments() [][]byte {
	return context.vmInput.Arguments
}

// ExtractCodeUpgradeFromArgs extracts the code and code metadata from the
// current VMInput.Arguments, assuming a contract code upgrade has been requested.
func (context *runtimeContext) ExtractCodeUpgradeFromArgs() ([]byte, []byte, error) {
	const numMinUpgradeArguments = 2

	arguments := context.vmInput.Arguments
	if len(arguments) < numMinUpgradeArguments {
		return nil, nil, arwen.ErrInvalidUpgradeArguments
	}

	code := arguments[0]
	codeMetadata := arguments[1]
	context.vmInput.Arguments = context.vmInput.Arguments[numMinUpgradeArguments:]
	return code, codeMetadata, nil
}

// FailExecution informs Wasmer to immediately stop the execution of the contract
// with BreakpointExecutionFailed and sets the corresponding VMOutput fields accordingly
func (context *runtimeContext) FailExecution(err error) {
	context.host.Output().SetReturnCode(vmcommon.ExecutionFailed)

	var message string
	if err != nil {
		message = err.Error()
	} else {
		message = "execution failed"
	}

	context.host.Output().SetReturnMessage(message)
	context.SetRuntimeBreakpointValue(arwen.BreakpointExecutionFailed)
}

// SignalUserError informs Wasmer to immediately stop the execution of the contract
// with BreakpointSignalError and sets the corresponding VMOutput fields accordingly
func (context *runtimeContext) SignalUserError(message string) {
	context.host.Output().SetReturnCode(vmcommon.UserError)
	context.host.Output().SetReturnMessage(message)
	context.SetRuntimeBreakpointValue(arwen.BreakpointSignalError)
}

// SetRuntimeBreakpointValue sets the specified runtime breakpoint in Wasmer,
// immediately stopping the contract execution.
func (context *runtimeContext) SetRuntimeBreakpointValue(value arwen.BreakpointValue) {
	context.instance.SetBreakpointValue(uint64(value))
}

// GetRuntimeBreakpointValue retrieves the value of the breakpoint that has
// stopped the execution of the contract.
func (context *runtimeContext) GetRuntimeBreakpointValue() arwen.BreakpointValue {
	return arwen.BreakpointValue(context.instance.GetBreakpointValue())
}

// VerifyContractCode performs validation on the WASM bytecode (declaration of memory and legal functions).
func (context *runtimeContext) VerifyContractCode() error {
	if !context.verifyCode {
		return nil
	}

	context.verifyCode = false

	err := context.validator.verifyMemoryDeclaration(context.instance)
	if err != nil {
		return err
	}

	err = context.validator.verifyFunctions(context.instance)
	if err != nil {
		return err
	}

	return nil
}

func (context *runtimeContext) ValidateCallbackName(callbackName string) error {
	if callbackName == arwen.InitFunctionName {
		return arwen.ErrInvalidFunctionName
	}

	return context.validator.verifyValidFunctionName(callbackName)
}

// ELrondAPIErrorShouldFailExecution specifies whether an error in the EEI should abort contract execution.
func (context *runtimeContext) ElrondAPIErrorShouldFailExecution() bool {
	return true
}

// ElrondSyncExecAPIErrorShouldFailExecution specifies whether an error in the
// EEI functions for synchronous execution should abort contract execution.
func (context *runtimeContext) ElrondSyncExecAPIErrorShouldFailExecution() bool {
	return true
}

// BigIntAPIErrorShouldFailExecution specifies whether an error in the EEI
// functions for BigInt operations should abort contract execution.
func (context *runtimeContext) BigIntAPIErrorShouldFailExecution() bool {
	return true
}

// CryptoAPIErrorShouldFailExecution specifies whether an error in the EEI
// functions for crypto operations should abort contract execution.
func (context *runtimeContext) CryptoAPIErrorShouldFailExecution() bool {
	return true
}

// GetPointsUsed returns the gas amount spent by the currently running Wasmer instance.
func (context *runtimeContext) GetPointsUsed() uint64 {
	return context.instance.GetPointsUsed()
}

// SetPointsUsed directly sets the gas amount already spent by the currently running Wasmer instance.
func (context *runtimeContext) SetPointsUsed(gasPoints uint64) {
	context.instance.SetPointsUsed(gasPoints)
}

// ReadOnly verifies whether the read-only execution flag is set.
func (context *runtimeContext) ReadOnly() bool {
	return context.readOnly
}

// SetReadOnly sets the read-only execution flag.
func (context *runtimeContext) SetReadOnly(readOnly bool) {
	context.readOnly = readOnly
}

// GetInstanceExports returns the objects exported by the WASM bytecode after
// the current Wasmer instance was started.
func (context *runtimeContext) GetInstanceExports() wasmer.ExportsMap {
	return context.instance.Exports
}

// CleanWasmerInstance cleans the current Wasmer instance.
func (context *runtimeContext) CleanWasmerInstance() {
	if context.instance == nil || context.IsWarmInstance() {
		return
	}

	arwen.RemoveHostContext(*context.instance.Data)
	context.instance.Clean()
	context.instance = nil
}

// GetFunctionToCall returns the callable contract method to be executed, as exported by the Wasmer instance.
func (context *runtimeContext) GetFunctionToCall() (wasmer.ExportedFunctionCallback, error) {
	exports := context.instance.Exports
	if function, ok := exports[context.callFunction]; ok {
		return function, nil
	}

	if context.callFunction == arwen.CallbackFunctionName {
		return nil, arwen.ErrNilCallbackFunction
	}

	return nil, arwen.ErrFuncNotFound
}

// GetInitFunction returns the callable contract method which initializes the
// contract immediately after deployment.
func (context *runtimeContext) GetInitFunction() wasmer.ExportedFunctionCallback {
	exports := context.instance.Exports
	if init, ok := exports[arwen.InitFunctionName]; ok {
		return init
	}

	return nil
}

func (context *runtimeContext) HasCallbackMethod() bool {
	_, ok := context.instance.Exports[arwen.CallbackFunctionName]
	return ok
}

// MemLoad reads a specified number of bytes from the given offset from the
// WASM memory of the currently running Wasmer instance.
func (context *runtimeContext) MemLoad(offset int32, length int32) ([]byte, error) {
	if length == 0 {
		return []byte{}, nil
	}

	memory := context.instance.InstanceCtx.Memory()
	memoryView := memory.Data()
	memoryLength := memory.Length()
	requestedEnd := uint32(offset + length)
	isOffsetTooSmall := offset < 0
	isOffsetTooLarge := uint32(offset) > memoryLength
	isRequestedEndTooLarge := requestedEnd > memoryLength
	isLengthNegative := length < 0

	if isOffsetTooSmall || isOffsetTooLarge {
		return nil, fmt.Errorf("mem load: %w", arwen.ErrBadBounds)
	}
	if isLengthNegative {
		return nil, fmt.Errorf("mem load: %w", arwen.ErrNegativeLength)
	}

	result := make([]byte, length)
	if isRequestedEndTooLarge {
		copy(result, memoryView[offset:])
	} else {
		copy(result, memoryView[offset:requestedEnd])
	}

	return result, nil
}

// MemStore writes the specified data bytes at the given offset in the WASM
// memory of the currently running Wasmer instance.
func (context *runtimeContext) MemStore(offset int32, data []byte) error {
	dataLength := int32(len(data))
	if dataLength == 0 {
		return nil
	}

	memory := context.instance.InstanceCtx.Memory()
	memoryView := memory.Data()
	memoryLength := memory.Length()
	requestedEnd := uint32(offset + dataLength)
	isOffsetTooSmall := offset < 0
	isNewPageNecessary := requestedEnd > memoryLength

	if isOffsetTooSmall {
		return arwen.ErrBadLowerBounds
	}
	if isNewPageNecessary {
		err := memory.Grow(1)
		if err != nil {
			return err
		}

		memoryView = memory.Data()
		memoryLength = memory.Length()
	}

	isRequestedEndTooLarge := requestedEnd > memoryLength
	if isRequestedEndTooLarge {
		return arwen.ErrBadUpperBounds
	}

	copy(memoryView[offset:requestedEnd], data)
	return nil
}
