package contexts

import (
	"bytes"

	"github.com/ElrondNetwork/arwen-wasm-vm/arwen"
	"github.com/ElrondNetwork/arwen-wasm-vm/config"
	"github.com/ElrondNetwork/arwen-wasm-vm/math"
	"github.com/ElrondNetwork/elrond-go/core/vmcommon"
)

type meteringContext struct {
	host               arwen.VMHost
	stateStack         []*meteringContext
	gasSchedule        *config.GasCost
	blockGasLimit      uint64
	initialGasProvided uint64
	initialCost        uint64
	gasForExecution    uint64
	gasStates          map[string]*contractGasState
}

type contractGasState struct {
	forwarded uint64
	returned  uint64
}

// NewMeteringContext creates a new meteringContext
func NewMeteringContext(
	host arwen.VMHost,
	gasMap config.GasScheduleMap,
	blockGasLimit uint64,
) (*meteringContext, error) {

	gasSchedule, err := config.CreateGasConfig(gasMap)
	if err != nil {
		return nil, err
	}

	context := &meteringContext{
		host:          host,
		stateStack:    make([]*meteringContext, 0),
		gasStates:     make(map[string]*contractGasState),
		gasSchedule:   gasSchedule,
		blockGasLimit: blockGasLimit,
	}

	context.InitState()

	return context, nil
}

// InitState resets the internal state of the MeteringContext
func (context *meteringContext) InitState() {
	context.gasStates = make(map[string]*contractGasState)
	context.initialGasProvided = 0
	context.initialCost = 0
	context.gasForExecution = 0
}

// PushState pushes the current state of the MeteringContext on its internal state stack
func (context *meteringContext) PushState() {
	newState := &meteringContext{
		initialGasProvided: context.initialGasProvided,
		initialCost:        context.initialCost,
		gasForExecution:    context.gasForExecution,
	}

	context.stateStack = append(context.stateStack, newState)
}

// PopSetActiveState pops the state at the top of the internal state stack, and
// sets it as the current state
func (context *meteringContext) PopSetActiveState() {
	stateStackLen := len(context.stateStack)
	if stateStackLen == 0 {
		return
	}

	prevState := context.stateStack[stateStackLen-1]
	context.stateStack = context.stateStack[:stateStackLen-1]

	context.initialGasProvided = prevState.initialGasProvided
	context.initialCost = prevState.initialCost
	context.gasForExecution = prevState.gasForExecution
}

// PopDiscard pops the state at the top of the internal state stack, and discards it
func (context *meteringContext) PopDiscard() {
	stateStackLen := len(context.stateStack)
	if stateStackLen == 0 {
		return
	}

	context.stateStack = context.stateStack[:stateStackLen-1]
}

// ClearStateStack reinitializes the internal state stack to an empty stack
func (context *meteringContext) ClearStateStack() {
	context.stateStack = make([]*meteringContext, 0)
}

// InitStateFromContractCallInput initializes the internal state of the
// MeteringContext using values taken from the provided ContractCallInput
func (context *meteringContext) InitStateFromContractCallInput(input *vmcommon.VMInput) {
	context.unlockGasIfAsyncCallback(input)
	context.initialGasProvided = input.GasProvided
	context.gasForExecution = input.GasProvided
	context.initialCost = 0
}

// unlockGasIfAsyncCallback unlocks the locked gas if the call type is async callback
func (context *meteringContext) unlockGasIfAsyncCallback(input *vmcommon.VMInput) {
	if input.CallType != vmcommon.AsynchronousCallBack {
		return
	}

	gasProvided := math.AddUint64(input.GasProvided, input.GasLocked)

	context.gasForExecution = gasProvided
	input.GasProvided = gasProvided
	input.GasLocked = 0
}

// GasSchedule returns the current gas schedule
func (context *meteringContext) GasSchedule() *config.GasCost {
	return context.gasSchedule
}

// SetGasSchedule sets the gas schedule to the given gas map
func (context *meteringContext) SetGasSchedule(gasMap config.GasScheduleMap) {
	gasSchedule, err := config.CreateGasConfig(gasMap)
	if err != nil {
		log.Error("SetGasSchedule createGasConfig", "error", err)
		return
	}
	context.gasSchedule = gasSchedule
}

// UseGas consumes the specified amount of gas on the currently running Wasmer instance.
func (context *meteringContext) UseGas(gas uint64) {
	gasUsed := math.AddUint64(context.host.Runtime().GetPointsUsed(), gas)
	context.host.Runtime().SetPointsUsed(gasUsed)
}

// RestoreGas deducts the specified amount of gas from the gas currently spent on the running Wasmer instance.
func (context *meteringContext) RestoreGas(gas uint64) {
	gasUsed := context.host.Runtime().GetPointsUsed()
	if gas <= gasUsed {
		gasUsed -= gas
		context.host.Runtime().SetPointsUsed(gasUsed)
	}
}

// FreeGas refunds the specified amount of gas to the caller.
func (context *meteringContext) FreeGas(gas uint64) {
	refund := math.AddUint64(context.host.Output().GetRefund(), gas)
	context.host.Output().SetRefund(refund)
}

// GasLeft computes the amount of gas left on the currently running Wasmer instance.
func (context *meteringContext) GasLeft() uint64 {
	gasProvided := context.gasForExecution
	gasUsed := context.host.Runtime().GetPointsUsed()

	if gasProvided < gasUsed {
		return 0
	}

	return gasProvided - gasUsed
}

// ForwardGas accumulates the gas forwarded by the current contract for the execution of other contracts
func (context *meteringContext) ForwardGas(sourceAddress []byte, destAddress []byte, gas uint64) {
	// Gas forwarded to any contract (including self-forwarding) is recorded for
	// the current contract.
	context.addForwardedGas(sourceAddress, gas)

	// If the address to which the gas is being forwarded already exists on the
	// execution stack, but is not directly below the current contract, it means
	// that any gas that has been forwarded to it is in fact returning via the
	// current contract.
	selfCall := bytes.Equal(sourceAddress, destAddress)
	alreadyOnExecutionStack := context.host.Runtime().IsContractOnTheStack(destAddress)
	if selfCall || alreadyOnExecutionStack {
		context.addReturnedGas(destAddress, gas)
	}
}

func (context *meteringContext) addForwardedGas(address []byte, gas uint64) {
	state := context.getContractGasState(address)
	state.forwarded = math.AddUint64(state.forwarded, gas)
}

func (context *meteringContext) addReturnedGas(address []byte, gas uint64) {
	state := context.getContractGasState(address)
	state.returned = math.AddUint64(state.returned, gas)
}

func (context *meteringContext) getTotalForwardedGas(address []byte) uint64 {
	state := context.getContractGasState(address)
	total, _ := math.SubUint64(state.forwarded, state.returned)
	return total
}

func (context *meteringContext) getContractGasState(address []byte) *contractGasState {
	key := string(address)
	state, exists := context.gasStates[key]
	if !exists {
		state = &contractGasState{
			forwarded: 0,
			returned:  0,
		}
		context.gasStates[key] = state
	}

	return state
}

// GasUsedByContract returns the gas used by the current contract.
func (context *meteringContext) GasUsedByContract() (uint64, uint64) {
	runtime := context.host.Runtime()
	executionGasUsed := runtime.GetPointsUsed()

	gasUsed := uint64(0)
	if context.host.IsArwenV2Enabled() {
		gasUsed = context.initialCost
	}

	gasUsed = math.AddUint64(gasUsed, executionGasUsed)

	totalGasForwarded := context.getTotalForwardedGas(runtime.GetSCAddress())
	remainedFromForwarded := uint64(0)
	gasUsed, remainedFromForwarded = math.SubUint64(gasUsed, totalGasForwarded)

	return gasUsed, remainedFromForwarded
}

// GasUsedForExecution returns the actual gas used for execution for the contract which needs to be restored
func (context *meteringContext) GasUsedForExecution() uint64 {
	gasUsed, _ := context.GasUsedByContract()
	gasUsed, _ = math.SubUint64(gasUsed, context.initialCost)
	return gasUsed
}

// GasSpentByContract calculates the entire gas consumption of the contract,
// without any gas forwarding.
func (context *meteringContext) GasSpentByContract() uint64 {
	runtime := context.host.Runtime()
	executionGasUsed := runtime.GetPointsUsed()

	gasSpent := uint64(0)
	if context.host.IsArwenV2Enabled() {
		gasSpent = context.initialCost
	}

	gasSpent = math.AddUint64(gasSpent, executionGasUsed)

	return gasSpent
}

// GetGasForExecution returns the gas left after the deduction of the initial gas from the provided gas
func (context *meteringContext) GetGasForExecution() uint64 {
	return context.gasForExecution
}

// GetGasProvided returns the fully provided gas for the sc execution
func (context *meteringContext) GetGasProvided() uint64 {
	return context.initialGasProvided
}

// GetSCPrepareInitialCost return the initial prepare cost for the sc execution
func (context *meteringContext) GetSCPrepareInitialCost() uint64 {
	return context.initialCost
}

// BoundGasLimit returns the maximum between the provided amount and the gas
// left on the currently running Wasmer instance.
func (context *meteringContext) BoundGasLimit(limit uint64) uint64 {
	gasLeft := context.GasLeft()

	if gasLeft < limit {
		return gasLeft
	}
	return limit
}

// UseGasForAsyncStep consumes the AsyncCallStep gas cost on the currently
// running Wasmer instance
func (context *meteringContext) UseGasForAsyncStep() error {
	gasSchedule := context.GasSchedule().ElrondAPICost
	gasToDeduct := gasSchedule.AsyncCallStep
	return context.UseGasBounded(gasToDeduct)
}

// UseGasBounded consumes the specified amount of gas on the currently running
// Wasmer instance, but returns an error if there is not enough gas left.
func (context *meteringContext) UseGasBounded(gasToUse uint64) error {
	gasLeft := context.GasLeft()
	if gasLeft < gasToUse {
		return arwen.ErrNotEnoughGas
	}
	context.UseGas(gasToUse)
	return nil
}

// ComputeGasLockedForAsync calculates the minimum amount of gas to lock for async callbacks
func (context *meteringContext) ComputeGasLockedForAsync() uint64 {
	baseGasSchedule := context.GasSchedule().BaseOperationCost
	apiGasSchedule := context.GasSchedule().ElrondAPICost
	codeSize := context.host.Runtime().GetSCCodeSize()

	costPerByte := baseGasSchedule.CompilePerByte
	if context.host.IsAheadOfTimeCompileEnabled() {
		costPerByte = baseGasSchedule.AoTPreparePerByte
	}

	// Exact amount of gas required to compile this SC again, to execute the callback
	compilationGasLock := uint64(0)
	if context.host.IsDynamicGasLockingEnabled() {
		compilationGasLock = math.MulUint64(codeSize, costPerByte)
	}

	// Minimum amount required to execute the callback
	executionGasLock := math.AddUint64(apiGasSchedule.AsyncCallStep, apiGasSchedule.AsyncCallbackGasLock)
	gasLockedForAsync := math.AddUint64(compilationGasLock, executionGasLock)

	return gasLockedForAsync
}

// UnlockGasIfAsyncCallback adds the locked gas to the gas provided for execution, before execution starts.
func (context *meteringContext) UnlockGasIfAsyncCallback() {
	input := context.host.Runtime().GetVMInput()
	if input.CallType != vmcommon.AsynchronousCallBack {
		return
	}

	gasProvided := math.AddUint64(input.GasProvided, input.GasLocked)

	input.GasProvided = gasProvided
	input.GasLocked = 0
}

// GetGasLocked returns the amount of gas locked during the current execution, as specified by the VMInput.
func (context *meteringContext) GetGasLocked() uint64 {
	input := context.host.Runtime().GetVMInput()
	return input.GasLocked
}

// BlockGasLimit returns the maximum amount of gas allowed to be consumed in a block.
func (context *meteringContext) BlockGasLimit() uint64 {
	return context.blockGasLimit
}

// DeductInitialGasForExecution deducts gas for compilation and locks gas if the execution is an asynchronous call
func (context *meteringContext) DeductInitialGasForExecution(contract []byte) error {
	costPerByte := context.gasSchedule.BaseOperationCost.CompilePerByte
	baseCost := uint64(0)
	if context.host.IsAheadOfTimeCompileEnabled() {
		costPerByte = context.gasSchedule.BaseOperationCost.AoTPreparePerByte
		baseCost = context.gasSchedule.BaseOperationCost.GetCode
	}
	err := context.deductInitialGas(contract, baseCost, costPerByte)
	if err != nil {
		return err
	}

	return nil
}

// DeductInitialGasForDirectDeployment deducts gas for the deployment of a contract initiated by a Transaction
func (context *meteringContext) DeductInitialGasForDirectDeployment(input arwen.CodeDeployInput) error {
	return context.deductInitialGas(
		input.ContractCode,
		context.gasSchedule.ElrondAPICost.CreateContract,
		context.gasSchedule.BaseOperationCost.CompilePerByte,
	)
}

// DeductInitialGasForIndirectDeployment deducts gas for the deployment of a contract initiated by another SmartContract
func (context *meteringContext) DeductInitialGasForIndirectDeployment(input arwen.CodeDeployInput) error {
	return context.deductInitialGas(
		input.ContractCode,
		0,
		context.gasSchedule.BaseOperationCost.CompilePerByte,
	)
}

func (context *meteringContext) deductInitialGas(
	code []byte,
	baseCost uint64,
	costPerByte uint64,
) error {
	input := context.host.Runtime().GetVMInput()
	codeLength := uint64(len(code))
	codeCost := math.MulUint64(codeLength, costPerByte)
	initialCost := math.AddUint64(baseCost, codeCost)

	if initialCost > input.GasProvided {
		return arwen.ErrNotEnoughGas
	}

	context.initialCost = initialCost
	context.gasForExecution = input.GasProvided - initialCost
	return nil
}
