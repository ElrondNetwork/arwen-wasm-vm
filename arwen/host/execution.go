package host

import (
	"bytes"
	"errors"
	"math/big"

	"github.com/ElrondNetwork/arwen-wasm-vm/arwen"
	"github.com/ElrondNetwork/arwen-wasm-vm/math"
	"github.com/ElrondNetwork/elrond-go-logger/check"
	"github.com/ElrondNetwork/elrond-go/core"
	"github.com/ElrondNetwork/elrond-go/core/vmcommon"
)

func (host *vmHost) doRunSmartContractCreate(input *vmcommon.ContractCreateInput) *vmcommon.VMOutput {
	host.InitState()
	defer host.Clean()

	_, blockchain, metering, output, runtime, _, storage := host.GetContexts()

	address, err := blockchain.NewAddress(input.CallerAddr)
	if err != nil {
		return output.CreateVMOutputInCaseOfError(err)
	}

	runtime.SetVMInput(&input.VMInput)
	runtime.SetSCAddress(address)
	metering.InitStateFromContractCallInput(&input.VMInput)

	output.AddTxValueToAccount(address, input.CallValue)
	storage.SetAddress(runtime.GetSCAddress())

	codeDeployInput := arwen.CodeDeployInput{
		ContractCode:         input.ContractCode,
		ContractCodeMetadata: input.ContractCodeMetadata,
		ContractAddress:      address,
		CodeDeployerAddress:  input.CallerAddr,
	}

	vmOutput, err := host.performCodeDeployment(codeDeployInput)
	if err != nil {
		return output.CreateVMOutputInCaseOfError(err)
	}

	return vmOutput
}

func (host *vmHost) performCodeDeployment(input arwen.CodeDeployInput) (*vmcommon.VMOutput, error) {
	log.Trace("performCodeDeployment", "address", input.ContractAddress, "len(code)", len(input.ContractCode), "metadata", input.ContractCodeMetadata)

	_, _, metering, output, runtime, _, _ := host.GetContexts()

	err := metering.DeductInitialGasForDirectDeployment(input)
	if err != nil {
		output.SetReturnCode(vmcommon.OutOfGas)
		return nil, err
	}

	runtime.MustVerifyNextContractCode()

	err = runtime.StartWasmerInstance(input.ContractCode, metering.GetGasForExecution(), true)
	if err != nil {
		log.Debug("performCodeDeployment/StartWasmerInstance", "err", err)
		return nil, arwen.ErrContractInvalid
	}

	err = host.callInitFunction()
	if err != nil {
		return nil, err
	}

	output.DeployCode(input)
	vmOutput := output.GetVMOutput()
	runtime.CleanWasmerInstance()
	return vmOutput, nil
}

// doRunSmartContractUpgrade upgrades a contract directly
func (host *vmHost) doRunSmartContractUpgrade(input *vmcommon.ContractCallInput) *vmcommon.VMOutput {
	host.InitState()
	defer host.Clean()

	_, _, metering, output, runtime, _, storage := host.GetContexts()

	runtime.InitStateFromContractCallInput(input)
	metering.InitStateFromContractCallInput(&input.VMInput)
	output.AddTxValueToAccount(input.RecipientAddr, input.CallValue)
	storage.SetAddress(runtime.GetSCAddress())

	code, codeMetadata, err := runtime.ExtractCodeUpgradeFromArgs()
	if err != nil {
		return output.CreateVMOutputInCaseOfError(arwen.ErrInvalidUpgradeArguments)
	}

	codeDeployInput := arwen.CodeDeployInput{
		ContractCode:         code,
		ContractCodeMetadata: codeMetadata,
		ContractAddress:      input.RecipientAddr,
		CodeDeployerAddress:  input.CallerAddr,
	}

	vmOutput, err := host.performCodeDeployment(codeDeployInput)
	if err != nil {
		return output.CreateVMOutputInCaseOfError(err)
	}

	return vmOutput
}

func (host *vmHost) checkGasForGetCode(input *vmcommon.ContractCallInput, metering arwen.MeteringContext) error {
	if !host.IsArwenV2Enabled() {
		return nil
	}

	getCodeBaseCost := metering.GasSchedule().BaseOperationCost.GetCode
	if input.GasProvided < getCodeBaseCost {
		return arwen.ErrNotEnoughGas
	}

	return nil
}

func (host *vmHost) doRunSmartContractCall(input *vmcommon.ContractCallInput) *vmcommon.VMOutput {
	host.InitState()
	defer host.Clean()

	_, _, metering, output, runtime, async, storage := host.GetContexts()

	runtime.InitStateFromContractCallInput(input)
	async.InitStateFromInput(input)
	metering.InitStateFromContractCallInput(&input.VMInput)
	output.AddTxValueToAccount(input.RecipientAddr, input.CallValue)
	storage.SetAddress(runtime.GetSCAddress())

	err := host.checkGasForGetCode(input, metering)
	if err != nil {
		return output.CreateVMOutputInCaseOfError(arwen.ErrNotEnoughGas)
	}

	contract, err := runtime.GetSCCode()
	if err != nil {
		return output.CreateVMOutputInCaseOfError(arwen.ErrContractNotFound)
	}

	err = metering.DeductInitialGasForExecution(contract)
	if err != nil {
		return output.CreateVMOutputInCaseOfError(arwen.ErrNotEnoughGas)
	}

	err = runtime.StartWasmerInstance(contract, metering.GetGasForExecution(), false)
	if err != nil {
		return output.CreateVMOutputInCaseOfError(arwen.ErrContractInvalid)
	}

	err = host.callSCMethod()
	if err != nil {
		return output.CreateVMOutputInCaseOfError(err)
	}

	vmOutput := output.GetVMOutput()
	runtime.CleanWasmerInstance()
	return vmOutput
}

// ExecuteOnDestContext pushes each context to the corresponding stack
// and initializes new contexts for executing the contract call with the given input
func (host *vmHost) ExecuteOnDestContext(input *vmcommon.ContractCallInput) (*vmcommon.VMOutput, uint64, error) {
	log.Trace("ExecuteOnDestContext", "function", input.Function)

	bigInt, _, metering, output, runtime, async, storage := host.GetContexts()

	bigInt.PushState()
	bigInt.InitState()

	output.PushState()
	output.CensorVMOutput()

	runtime.PushState()
	runtime.InitStateFromContractCallInput(input)

	// TODO async.LoadOrInit(), not just Init; the contract invoked here likely has a
	// persisted AsyncContext of its own.
	async.PushState()
	async.InitStateFromInput(input)

	metering.PushState()
	metering.InitStateFromContractCallInput(&input.VMInput)

	storage.PushState()
	storage.SetAddress(runtime.GetSCAddress())

	gasUsedBeforeReset := uint64(0)

	// Perform a value transfer to the called SC. If the execution fails, this
	// transfer will not persist.
	err := output.TransferValueOnly(input.RecipientAddr, input.CallerAddr, input.CallValue)
	if err != nil {
		vmOutput := host.finishExecuteOnDestContext(err)
		return vmOutput, gasUsedBeforeReset, err
	}

	gasUsedBeforeReset, err = host.execute(input)
	if err != nil {
		vmOutput := host.finishExecuteOnDestContext(err)
		return vmOutput, gasUsedBeforeReset, err
	}

	err = async.Execute()
	vmOutput := host.finishExecuteOnDestContext(err)
	return vmOutput, gasUsedBeforeReset, err
}

func (host *vmHost) finishExecuteOnDestContext(executeErr error) *vmcommon.VMOutput {
	bigInt, _, metering, output, runtime, async, storage := host.GetContexts()

	var vmOutput *vmcommon.VMOutput
	if executeErr != nil {
		// Execution failed: restore contexts as if the execution didn't happen,
		// but first create a vmOutput to capture the error.
		vmOutput = output.CreateVMOutputInCaseOfError(executeErr)
	} else {
		// Retrieve the VMOutput before popping the Runtime state and the previous
		// instance, to ensure accurate GasRemaining
		vmOutput = output.GetVMOutput()
	}

	childContract := runtime.GetSCAddress()
	gasSpentByChildContract := metering.GasSpentByContract()

	if vmOutput.ReturnCode != vmcommon.Ok {
		gasSpentByChildContract = 0
	}

	// Restore the previous context states, except Output, which will be merged
	// into the initial state (VMOutput), but only if it the child execution
	// returned vmcommon.Ok.
	bigInt.PopSetActiveState()
	metering.PopSetActiveState()
	runtime.PopSetActiveState()
	storage.PopSetActiveState()
	async.PopSetActiveState()

	// Restore remaining gas to the caller Wasmer instance
	metering.RestoreGas(vmOutput.GasRemaining)
	metering.ForwardGas(runtime.GetSCAddress(), childContract, gasSpentByChildContract)

	if vmOutput.ReturnCode == vmcommon.Ok {
		output.PopMergeActiveState()
	} else {
		output.PopSetActiveState()
	}

	return vmOutput
}

// ExecuteOnSameContext executes the contract call with the given input
// on the same runtime context. Some other contexts are backed up.
func (host *vmHost) ExecuteOnSameContext(input *vmcommon.ContractCallInput) error {
	log.Trace("ExecuteOnSameContext", "function", input.Function)

	if host.IsBuiltinFunctionName(input.Function) {
		return arwen.ErrBuiltinCallOnSameContextDisallowed
	}

	bigInt, _, metering, output, runtime, _, _ := host.GetContexts()

	// Back up the states of the contexts (except Storage and Async, which aren't
	// affected by ExecuteOnSameContext())
	bigInt.PushState()
	output.PushState()

	runtime.PushState()
	runtime.InitStateFromContractCallInput(input)

	metering.PushState()
	metering.InitStateFromContractCallInput(&input.VMInput)

	// Perform a value transfer to the called SC. If the execution fails, this
	// transfer will not persist.
	err := output.TransferValueOnly(input.RecipientAddr, input.CallerAddr, input.CallValue)
	if err != nil {
		host.finishExecuteOnSameContext(err)
		return err
	}

	_, err = host.execute(input)
	if err != nil {
		host.finishExecuteOnSameContext(err)
		return err
	}

	host.finishExecuteOnSameContext(nil)
	return nil
}

func (host *vmHost) finishExecuteOnSameContext(executeErr error) {
	bigInt, _, metering, output, runtime, _, _ := host.GetContexts()

	if output.ReturnCode() != vmcommon.Ok || executeErr != nil {
		// Execution failed: restore contexts as if the execution didn't happen.
		bigInt.PopSetActiveState()
		metering.PopSetActiveState()
		output.PopSetActiveState()
		runtime.PopSetActiveState()

		return
	}

	// Retrieve the VMOutput before popping the Runtime state and the previous
	// instance, to ensure accurate GasRemaining
	vmOutput := output.GetVMOutput()
	childContract := runtime.GetSCAddress()
	gasSpentByContract := metering.GasSpentByContract()
	if vmOutput.ReturnCode != vmcommon.Ok {
		gasSpentByContract = 0
	}

	// Execution successful: discard the backups made at the beginning and
	// resume from the new state. However, output.PopDiscard() will ensure that
	// all GasUsed records will be restored, undoing the action of output.ResetGas()
	bigInt.PopDiscard()
	output.PopDiscard()
	metering.PopSetActiveState()
	runtime.PopSetActiveState()

	// Restore remaining gas to the caller Wasmer instance
	metering.RestoreGas(vmOutput.GasRemaining)
	metering.ForwardGas(runtime.GetSCAddress(), childContract, gasSpentByContract)
}

func (host *vmHost) isInitFunctionBeingCalled() bool {
	functionName := host.Runtime().Function()
	return functionName == arwen.InitFunctionName || functionName == arwen.InitFunctionNameEth
}

func (host *vmHost) isBuiltinFunctionBeingCalled() bool {
	functionName := host.Runtime().Function()
	return host.IsBuiltinFunctionName(functionName)
}

// IsBuiltinFunctionName returns true if the given function name is the same as any protocol builtin function
func (host *vmHost) IsBuiltinFunctionName(functionName string) bool {
	_, ok := host.protocolBuiltinFunctions[functionName]
	return ok
}

// CreateNewContract creates a new contract indirectly (from another Smart Contract)
func (host *vmHost) CreateNewContract(input *vmcommon.ContractCreateInput) ([]byte, error) {
	_, blockchain, metering, output, runtime, _, _ := host.GetContexts()

	codeDeployInput := arwen.CodeDeployInput{
		ContractCode:         input.ContractCode,
		ContractCodeMetadata: input.ContractCodeMetadata,
		ContractAddress:      nil,
		CodeDeployerAddress:  input.CallerAddr,
	}
	err := metering.DeductInitialGasForIndirectDeployment(codeDeployInput)
	if err != nil {
		return nil, err
	}

	if runtime.ReadOnly() {
		err = arwen.ErrInvalidCallOnReadOnlyMode
		return nil, err
	}

	newContractAddress, err := blockchain.NewAddress(input.CallerAddr)
	if err != nil {
		return nil, err
	}

	if blockchain.AccountExists(newContractAddress) {
		err = arwen.ErrDeploymentOverExistingAccount
		return nil, err
	}

	codeDeployInput.ContractAddress = newContractAddress
	output.DeployCode(codeDeployInput)

	defer func() {
		if err != nil {
			output.DeleteOutputAccount(newContractAddress)
		}
	}()

	runtime.MustVerifyNextContractCode()

	initCallInput := &vmcommon.ContractCallInput{
		RecipientAddr:     newContractAddress,
		Function:          arwen.InitFunctionName,
		AllowInitFunction: true,
		VMInput:           input.VMInput,
	}

	_, _, err = host.ExecuteOnDestContext(initCallInput)
	if err != nil {
		return nil, err
	}

	blockchain.IncreaseNonce(input.CallerAddr)

	return newContractAddress, nil
}

func (host *vmHost) checkUpgradePermission(vmInput *vmcommon.ContractCallInput) error {
	contract, err := host.blockChainHook.GetUserAccount(vmInput.RecipientAddr)
	if err != nil {
		return err
	}
	if check.IfNilReflect(contract) {
		return arwen.ErrNilContract
	}

	codeMetadata := vmcommon.CodeMetadataFromBytes(contract.GetCodeMetadata())
	isUpgradeable := codeMetadata.Upgradeable
	callerAddress := vmInput.CallerAddr
	ownerAddress := contract.GetOwnerAddress()
	isCallerOwner := bytes.Equal(callerAddress, ownerAddress)

	if isUpgradeable && isCallerOwner {
		return nil
	}

	return arwen.ErrUpgradeNotAllowed
}

// executeUpgrade upgrades a contract indirectly (from another contract). This
// function follows the convention of executeSmartContractCall().
func (host *vmHost) executeUpgrade(input *vmcommon.ContractCallInput) error {
	_, _, metering, output, runtime, _, _ := host.GetContexts()

	err := host.checkUpgradePermission(input)
	if err != nil {
		return err
	}

	code, codeMetadata, err := runtime.ExtractCodeUpgradeFromArgs()
	if err != nil {
		return arwen.ErrInvalidUpgradeArguments
	}

	codeDeployInput := arwen.CodeDeployInput{
		ContractCode:         code,
		ContractCodeMetadata: codeMetadata,
		ContractAddress:      input.RecipientAddr,
		CodeDeployerAddress:  input.CallerAddr,
	}

	err = metering.DeductInitialGasForDirectDeployment(codeDeployInput)
	if err != nil {
		output.SetReturnCode(vmcommon.OutOfGas)
		return err
	}

	runtime.MustVerifyNextContractCode()

	err = runtime.StartWasmerInstance(codeDeployInput.ContractCode, metering.GetGasForExecution(), true)
	if err != nil {
		log.Debug("performCodeDeployment/StartWasmerInstance", "err", err)
		return arwen.ErrContractInvalid
	}

	err = host.callInitFunction()
	if err != nil {
		return err
	}

	output.DeployCode(codeDeployInput)
	if output.ReturnCode() != vmcommon.Ok {
		return arwen.ErrReturnCodeNotOk
	}

	return nil
}

// executeSmartContractCall executes an indirect call to a smart contract,
// assuming there is an already-running Wasmer instance with another contract
// that has requested the indirect call. This method creates a new Wasmer
// instance and pushes the previous one onto the Runtime instance stack, but it
// will not pop the previous instance back - that remains the responsibility of
// the calling code. Also, this method does not restore the gas remaining after
// the indirect call, it does not push the states of any Host Context onto
// their respective stacks, nor does it pop any state stack. Handling the state
// stacks and the remaining gas are responsibilities of the calling code, which
// must push and pop as required, before and after calling this method, and
// handle the remaining gas. These principles also apply to indirect contract
// upgrading (via host.executeUpgrade(), which also does not pop the previous
// instance from the Runtime instance stack, nor does it restore the remaining
// gas).
func (host *vmHost) executeSmartContractCall(
	input *vmcommon.ContractCallInput,
	metering arwen.MeteringContext,
	runtime arwen.RuntimeContext,
	output arwen.OutputContext,
	withInitialGasDeduct bool,
) error {
	if host.isInitFunctionBeingCalled() && !input.AllowInitFunction {
		return arwen.ErrInitFuncCalledInRun
	}

	// Use all gas initially, on the Wasmer instance of the caller. In case of
	// successful execution, the unused gas will be restored.
	metering.UseGas(input.GasProvided)

	isUpgrade := input.Function == arwen.UpgradeFunctionName
	if isUpgrade {
		return host.executeUpgrade(input)
	}

	contract, err := runtime.GetSCCode()
	if err != nil {
		return err
	}

	if withInitialGasDeduct {
		err = metering.DeductInitialGasForExecution(contract)
		if err != nil {
			return err
		}
	}

	// Replace the current Wasmer instance of the Runtime with a new one; this
	// assumes that the instance was preserved on the Runtime instance stack
	// before calling executeSmartContractCall().
	err = runtime.StartWasmerInstance(contract, metering.GetGasForExecution(), false)
	if err != nil {
		return err
	}

	err = host.callSCMethodIndirect()
	if err != nil {
		return err
	}

	if output.ReturnCode() != vmcommon.Ok {
		return arwen.ErrReturnCodeNotOk
	}

	return nil
}

func (host *vmHost) execute(input *vmcommon.ContractCallInput) (uint64, error) {
	_, _, metering, output, runtime, _, storage := host.GetContexts()

	if host.isBuiltinFunctionBeingCalled() {
		newVMInput, gasUsedBeforeReset, err := host.callBuiltinFunction(input)
		if err != nil {
			return gasUsedBeforeReset, err
		}

		if newVMInput != nil {
			runtime.InitStateFromContractCallInput(newVMInput)
			metering.InitStateFromContractCallInput(&newVMInput.VMInput)
			storage.SetAddress(runtime.GetSCAddress())
			err = host.executeSmartContractCall(newVMInput, metering, runtime, output, false)
			if err != nil {
				host.revertESDTTransfer(input)
			}

			return gasUsedBeforeReset, err
		}

		return gasUsedBeforeReset, nil
	}

	return 0, host.executeSmartContractCall(input, metering, runtime, output, true)
}

func (host *vmHost) callSCMethodIndirect() error {
	function, err := host.Runtime().GetFunctionToCall()
	if err != nil {
		if errors.Is(err, arwen.ErrNilCallbackFunction) {
			return nil
		}
		return err
	}

	_, err = function()
	if err != nil {
		err = host.handleBreakpointIfAny(err)
	}

	return err
}

func (host *vmHost) revertESDTTransfer(input *vmcommon.ContractCallInput) {
	if input.Function != core.BuiltInFunctionESDTTransfer {
		return
	}
	if len(input.Arguments) < 2 {
		return
	}

	revertInput := &vmcommon.ContractCallInput{
		VMInput: vmcommon.VMInput{
			CallerAddr:     input.RecipientAddr,
			Arguments:      input.Arguments[:2],
			CallValue:      big.NewInt(0),
			CallType:       vmcommon.DirectCall,
			GasPrice:       input.GasPrice,
			GasProvided:    input.GasProvided,
			GasLocked:      0,
			OriginalTxHash: input.OriginalTxHash,
			CurrentTxHash:  input.CurrentTxHash,
			ESDTValue:      big.NewInt(0),
			ESDTTokenName:  nil,
		},
		RecipientAddr:     input.CallerAddr,
		Function:          input.Function,
		AllowInitFunction: false,
	}

	vmOutput, err := host.blockChainHook.ProcessBuiltInFunction(revertInput)
	if err != nil {
		log.Error("revertESDTTransfer failed", "error", err)
	}
	if vmOutput.ReturnCode != vmcommon.Ok {
		log.Error("revertESDTTransfer failed", "returnCode", vmOutput.ReturnCode, "returnMessage", vmOutput.ReturnMessage)
	}
}

func (host *vmHost) callBuiltinFunction(input *vmcommon.ContractCallInput) (*vmcommon.ContractCallInput, uint64, error) {
	_, _, metering, output, runtime, _, _ := host.GetContexts()

	gasConsumedForExecution := host.computeGasUsedInExecutionBeforeReset(input)
	runtime.SetPointsUsed(0)
	vmOutput, err := host.blockChainHook.ProcessBuiltInFunction(input)
	if err != nil {
		metering.UseGas(input.GasProvided)
		return nil, gasConsumedForExecution, err
	}

	gasConsumed, _ := math.SubUint64(input.GasProvided, vmOutput.GasRemaining)
	for _, outAcc := range vmOutput.OutputAccounts {
		for _, outTransfer := range outAcc.OutputTransfers {
			if outTransfer.GasLimit > 0 || outTransfer.GasLocked > 0 {
				gasForwarded := math.AddUint64(outTransfer.GasLocked, outTransfer.GasLimit)
				metering.ForwardGas(runtime.GetSCAddress(), nil, gasForwarded)
				gasConsumed = math.AddUint64(gasConsumed, outTransfer.GasLocked)
			}
		}
	}

	if vmOutput.GasRemaining < input.GasProvided {
		metering.UseGas(gasConsumed)
	}

	newVMInput, err := host.isSCExecutionAfterBuiltInFunc(input, vmOutput)
	if err != nil {
		return nil, gasConsumedForExecution, err
	}

	if newVMInput != nil {
		for _, outAcc := range vmOutput.OutputAccounts {
			outAcc.OutputTransfers = make([]vmcommon.OutputTransfer, 0)
		}
	}

	output.AddToActiveState(vmOutput)

	return newVMInput, gasConsumedForExecution, nil
}

func (host *vmHost) computeGasUsedInExecutionBeforeReset(vmInput *vmcommon.ContractCallInput) uint64 {
	metering := host.Metering()
	gasUsedForExecution, _ := math.SubUint64(metering.GasUsedForExecution(), vmInput.GasLocked)
	return gasUsedForExecution
}

func (host *vmHost) checkFinalGasAfterExit() error {
	if !host.IsArwenV2Enabled() {
		return nil
	}

	if host.Runtime().GetPointsUsed() > host.Metering().GetGasForExecution() {
		return arwen.ErrNotEnoughGas
	}

	return nil
}

func (host *vmHost) callInitFunction() error {
	runtime := host.Runtime()
	init := runtime.GetInitFunction()
	if init == nil {
		return nil
	}

	_, err := init()
	if err != nil {
		err = host.handleBreakpointIfAny(err)
	}

	if err == nil {
		err = host.checkFinalGasAfterExit()
	}

	return err
}

func (host *vmHost) callSCMethod() error {
	runtime := host.Runtime()
	vmInput := runtime.GetVMInput()
	async := host.Async()
	callType := vmInput.CallType

	if callType == vmcommon.AsynchronousCallBack {
		async.Load()
		asyncCall, err := async.UpdateCurrentCallStatus()
		if err != nil {
			return async.PostprocessCrossShardCallback()
		}

		runtime.SetCustomCallFunction(asyncCall.GetCallbackName())
	}

	err := host.verifyAllowedFunctionCall()
	if err != nil {
		return err
	}

	function, err := runtime.GetFunctionToCall()
	if err != nil {
		return err
	}

	_, err = function()
	if err != nil {
		err = host.handleBreakpointIfAny(err)
	}
	if err == nil {
		err = host.checkFinalGasAfterExit()
	}
	if err != nil {
		return err
	}

	err = async.Execute()
	if err != nil {
		return err
	}

	switch callType {
	case vmcommon.DirectCall:
		break
	case vmcommon.AsynchronousCall:
		err = host.sendAsyncCallbackToCaller()
	case vmcommon.AsynchronousCallBack:
		err = async.PostprocessCrossShardCallback()
	default:
		err = arwen.ErrUnknownCallType
	}

	return err
}

func (host *vmHost) verifyAllowedFunctionCall() error {
	runtime := host.Runtime()
	functionName := runtime.Function()

	isInit := functionName == arwen.InitFunctionName || functionName == arwen.InitFunctionNameEth
	if isInit {
		return arwen.ErrInitFuncCalledInRun
	}

	isCallBack := functionName == arwen.CallbackFunctionName
	isInAsyncCallBack := runtime.GetVMInput().CallType == vmcommon.AsynchronousCallBack
	if isCallBack && !isInAsyncCallBack {
		return arwen.ErrCallBackFuncCalledInRun
	}

	return nil
}

func (host *vmHost) isSCExecutionAfterBuiltInFunc(
	vmInput *vmcommon.ContractCallInput,
	vmOutput *vmcommon.VMOutput,
) (*vmcommon.ContractCallInput, error) {
	if vmOutput.ReturnCode != vmcommon.Ok {
		return nil, nil
	}

	if !core.IsSmartContractAddress(vmInput.RecipientAddr) {
		return nil, nil
	}

	if !host.AreInSameShard(vmInput.CallerAddr, vmInput.RecipientAddr) {
		return nil, nil
	}

	outAcc, ok := vmOutput.OutputAccounts[string(vmInput.RecipientAddr)]
	if !ok {
		return nil, nil
	}
	if len(outAcc.OutputTransfers) != 1 {
		return nil, nil
	}

	callType := vmInput.CallType
	scCallOutTransfer := outAcc.OutputTransfers[0]
	txData := prependCallbackToTxDataIfAsyncCall(scCallOutTransfer.Data, callType)

	function, arguments, err := host.CallArgsParser().ParseData(txData)
	if err != nil {
		return nil, err
	}

	// TODO CurrentTxHash might equal PrevTxHash, because Arwen does not generate
	// SCRs between async steps; analyze and fix this (e.g. create dummy SCRs
	// just for hashing, or hash the VMInput itself)
	newVMInput := &vmcommon.ContractCallInput{
		VMInput: vmcommon.VMInput{
			CallerAddr:     vmInput.CallerAddr,
			Arguments:      arguments,
			CallValue:      big.NewInt(0),
			CallType:       callType,
			GasPrice:       vmInput.GasPrice,
			GasProvided:    scCallOutTransfer.GasLimit,
			GasLocked:      scCallOutTransfer.GasLocked,
			OriginalTxHash: vmInput.OriginalTxHash,
			CurrentTxHash:  vmInput.CurrentTxHash,
			PrevTxHash:     vmInput.PrevTxHash,
		},
		RecipientAddr:     vmInput.RecipientAddr,
		Function:          function,
		AllowInitFunction: false,
	}

	fillWithESDTValue(vmInput, newVMInput)

	return newVMInput, nil
}

func fillWithESDTValue(fullVMInput *vmcommon.ContractCallInput, newVMInput *vmcommon.ContractCallInput) {
	if fullVMInput.Function != core.BuiltInFunctionESDTTransfer {
		return
	}

	newVMInput.ESDTTokenName = fullVMInput.Arguments[0]
	newVMInput.ESDTValue = big.NewInt(0).SetBytes(fullVMInput.Arguments[1])
}

func prependCallbackToTxDataIfAsyncCall(txData []byte, callType vmcommon.CallType) string {
	if callType == vmcommon.AsynchronousCallBack {
		return string(append([]byte(arwen.CallbackFunctionName), txData...))
	}

	return string(txData)
}
