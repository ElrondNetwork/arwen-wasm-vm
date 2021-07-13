package host

import (
	"bytes"
	"encoding/hex"
	"errors"
	"fmt"
	"math/big"

	"github.com/ElrondNetwork/arwen-wasm-vm/v1_3/arwen"
	"github.com/ElrondNetwork/arwen-wasm-vm/v1_3/math"
	"github.com/ElrondNetwork/elrond-go-logger/check"
	vmcommon "github.com/ElrondNetwork/elrond-vm-common"
)

func (host *vmHost) doRunSmartContractCreate(input *vmcommon.ContractCreateInput) *vmcommon.VMOutput {
	host.InitState()
	defer func() {
		errs := host.GetRuntimeErrors()
		if errs != nil {
			log.Trace(fmt.Sprintf("doRunSmartContractCreate full error list"), "error", errs)
		}
		host.Clean()
	}()

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
		log.Trace("doRunSmartContractCreate", "error", err)
		return output.CreateVMOutputInCaseOfError(err)
	}

	log.Trace("doRunSmartContractCreate",
		"retCode", vmOutput.ReturnCode,
		"message", vmOutput.ReturnMessage,
		"data", vmOutput.ReturnData)

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
		log.Trace("performCodeDeployment/StartWasmerInstance", "err", err)
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
	defer func() {
		errs := host.GetRuntimeErrors()
		if errs != nil {
			log.Trace(fmt.Sprintf("doRunSmartContractUpgrade full error list"), "error", errs)
		}
		host.Clean()
	}()

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
		log.Trace("doRunSmartContractUpgrade", "error", err)
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

func (host *vmHost) doRunSmartContractCall(input *vmcommon.ContractCallInput) (vmOutput *vmcommon.VMOutput) {
	host.InitState()
	defer func() {
		errs := host.GetRuntimeErrors()
		if errs != nil {
			log.Trace(fmt.Sprintf("doRunSmartContractCall full error list for %s", input.Function), "error", errs)
		}
		host.Clean()
	}()

	_, _, metering, output, runtime, _, storage := host.GetContexts()

	runtime.InitStateFromContractCallInput(input)
	metering.InitStateFromContractCallInput(&input.VMInput)
	output.AddTxValueToAccount(input.RecipientAddr, input.CallValue)
	storage.SetAddress(runtime.GetSCAddress())

	err := host.checkGasForGetCode(input, metering)
	if err != nil {
		log.Trace("doRunSmartContractCall get code", "error", arwen.ErrNotEnoughGas)
		return output.CreateVMOutputInCaseOfError(arwen.ErrNotEnoughGas)
	}

	contract, err := runtime.GetSCCode()
	if err != nil {
		log.Trace("doRunSmartContractCall get code", "error", arwen.ErrContractNotFound)
		return output.CreateVMOutputInCaseOfError(arwen.ErrContractNotFound)
	}

	err = metering.DeductInitialGasForExecution(contract)
	if err != nil {
		log.Trace("doRunSmartContractCall initial gas", "error", arwen.ErrNotEnoughGas)
		return output.CreateVMOutputInCaseOfError(arwen.ErrNotEnoughGas)
	}

	err = runtime.StartWasmerInstance(contract, metering.GetGasForExecution(), false)
	if err != nil {
		return output.CreateVMOutputInCaseOfError(arwen.ErrContractInvalid)
	}

	err = host.callSCMethod()
	if err != nil {
		log.Trace("doRunSmartContractCall", "error", err)
		return output.CreateVMOutputInCaseOfError(err)
	}

	vmOutput = output.GetVMOutput()

	log.Trace("doRunSmartContractCall finished",
		"retCode", vmOutput.ReturnCode,
		"message", vmOutput.ReturnMessage,
		"data", vmOutput.ReturnData)

	runtime.CleanWasmerInstance()
	return
}

func copyTxHashesFromContext(copyEnabled bool, runtime arwen.RuntimeContext, input *vmcommon.ContractCallInput) {
	if !copyEnabled {
		return
	}
	currentVMInput := runtime.GetVMInput()
	if len(currentVMInput.OriginalTxHash) > 0 {
		input.OriginalTxHash = currentVMInput.OriginalTxHash
	}
	if len(currentVMInput.CurrentTxHash) > 0 {
		input.CurrentTxHash = currentVMInput.CurrentTxHash
	}
	if len(currentVMInput.PrevTxHash) > 0 {
		input.PrevTxHash = currentVMInput.PrevTxHash
	}

}

// ExecuteOnDestContext pushes each context to the corresponding stack
// and initializes new contexts for executing the contract call with the given input
func (host *vmHost) ExecuteOnDestContext(input *vmcommon.ContractCallInput) (vmOutput *vmcommon.VMOutput, err error) {
	log.Trace("ExecuteOnDestContext", "caller", input.CallerAddr, "dest", input.RecipientAddr, "function", input.Function)

	scExecutionInput := input

	blockchain := host.Blockchain()
	blockchain.PushState()

	if host.IsBuiltinFunctionName(input.Function) {
		scExecutionInput, vmOutput, err = host.handleBuiltinFunctionCall(input)
		if err != nil {
			blockchain.PopSetActiveState()
			host.Runtime().AddError(err, input.Function)
			vmOutput = host.Output().CreateVMOutputInCaseOfError(err)
			return
		}
	}

	if scExecutionInput != nil {
		vmOutput, err = host.executeOnDestContextNoBuiltinFunction(scExecutionInput)
	}

	if err != nil {
		blockchain.PopSetActiveState()
	} else {
		blockchain.PopDiscard()
	}

	return
}

func (host *vmHost) handleBuiltinFunctionCall(input *vmcommon.ContractCallInput) (*vmcommon.ContractCallInput, *vmcommon.VMOutput, error) {
	output := host.Output()
	postBuiltinInput, builtinOutput, err := host.callBuiltinFunction(input)
	if err != nil {
		log.Trace("ExecuteOnDestContext builtin function", "error", err)
		return nil, nil, err
	}

	output.AddToActiveState(builtinOutput)

	return postBuiltinInput, builtinOutput, nil
}

func (host *vmHost) executeOnDestContextNoBuiltinFunction(input *vmcommon.ContractCallInput) (*vmcommon.VMOutput, error) {
	bigInt, _, metering, output, runtime, async, storage := host.GetContexts()
	bigInt.PushState()
	bigInt.InitState()

	output.PushState()
	output.CensorVMOutput()

	copyTxHashesFromContext(host.IsESDTFunctionsEnabled(), runtime, input)
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

	var err error
	var vmOutput *vmcommon.VMOutput

	// Perform a value transfer to the called SC. If the execution fails, this
	// transfer will not persist.
	if input.CallType != vmcommon.AsynchronousCallBack || input.CallValue.Cmp(arwen.Zero) == 0 {
		err = output.TransferValueOnly(input.RecipientAddr, input.CallerAddr, input.CallValue, false)
		if err != nil {
			vmOutput = host.finishExecuteOnDestContext(err)
			if err == nil && vmOutput.ReturnCode != vmcommon.Ok {
				err = arwen.ErrExecutionFailed
			}
			log.Trace("ExecuteOnDestContext transfer", "error", err)
			return vmOutput, err
		}
	}

	err = host.execute(input)
	if err != nil {
		log.Trace("ExecuteOnDestContext execution", "error", err)
		vmOutput = host.finishExecuteOnDestContext(err)
		return vmOutput, err
	}

	err = async.Execute()
	vmOutput = host.finishExecuteOnDestContext(err)
	return vmOutput, err
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

	gasSpentByChildContract := metering.GasSpentByContract()

	// Restore the previous context states
	bigInt.PopSetActiveState()
	storage.PopSetActiveState()

	if vmOutput.ReturnCode == vmcommon.Ok {
		metering.PopMergeActiveState()
		output.PopMergeActiveState()
	} else {
		metering.PopSetActiveState()
		output.PopSetActiveState()
	}

	// Return to the caller context completely
	runtime.PopSetActiveState()
	async.PopSetActiveState()

	// Restore remaining gas to the caller Wasmer instance
	metering.RestoreGas(vmOutput.GasRemaining)

	log.Trace("ExecuteOnDestContext finished", "gas spent", gasSpentByChildContract)

	return vmOutput
}

// ExecuteOnSameContext executes the contract call with the given input
// on the same runtime context. Some other contexts are backed up.
func (host *vmHost) ExecuteOnSameContext(input *vmcommon.ContractCallInput) (err error) {
	log.Trace("ExecuteOnSameContext", "function", input.Function)

	if host.IsBuiltinFunctionName(input.Function) {
		return arwen.ErrBuiltinCallOnSameContextDisallowed
	}

	bigInt, blockchain, metering, output, runtime, _, _ := host.GetContexts()

	// Back up the states of the contexts (except Storage and Async, which aren't
	// affected by ExecuteOnSameContext())
	bigInt.PushState()
	output.PushState()

	copyTxHashesFromContext(host.IsESDTFunctionsEnabled(), runtime, input)
	runtime.PushState()
	runtime.InitStateFromContractCallInput(input)

	metering.PushState()
	metering.InitStateFromContractCallInput(&input.VMInput)

	blockchain.PushState()

	// Perform a value transfer to the called SC. If the execution fails, this
	// transfer will not persist.
	err = output.TransferValueOnly(input.RecipientAddr, input.CallerAddr, input.CallValue, false)
	if err != nil {
		runtime.AddError(err, input.Function)
		host.finishExecuteOnSameContext(err)
		return
	}

	err = host.execute(input)
	runtime.AddError(err, input.Function)
	host.finishExecuteOnSameContext(err)
	return err
}

func (host *vmHost) finishExecuteOnSameContext(executeErr error) {
	bigInt, blockchain, metering, output, runtime, _, _ := host.GetContexts()

	if output.ReturnCode() != vmcommon.Ok || executeErr != nil {
		// Execution failed: restore contexts as if the execution didn't happen.
		bigInt.PopSetActiveState()
		metering.PopSetActiveState()
		output.PopSetActiveState()
		runtime.PopSetActiveState()
		blockchain.PopSetActiveState()
		return
	}

	// Execution successful; retrieve the VMOutput before popping the Runtime
	// state and the previous instance, to ensure accurate GasRemaining and
	// GasUsed for all accounts.
	vmOutput := output.GetVMOutput()

	metering.PopMergeActiveState()
	output.PopDiscard()
	bigInt.PopDiscard()
	blockchain.PopDiscard()
	runtime.PopSetActiveState()

	// Restore remaining gas to the caller (parent) Wasmer instance
	metering.RestoreGas(vmOutput.GasRemaining)
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
func (host *vmHost) CreateNewContract(input *vmcommon.ContractCreateInput) (newContractAddress []byte, err error) {
	newContractAddress = nil
	err = nil

	defer func() {
		if err != nil {
			newContractAddress = nil
		}
	}()

	_, blockchain, metering, output, runtime, _, _ := host.GetContexts()

	codeDeployInput := arwen.CodeDeployInput{
		ContractCode:         input.ContractCode,
		ContractCodeMetadata: input.ContractCodeMetadata,
		ContractAddress:      nil,
		CodeDeployerAddress:  input.CallerAddr,
	}
	err = metering.DeductInitialGasForIndirectDeployment(codeDeployInput)
	if err != nil {
		return
	}

	if runtime.ReadOnly() {
		err = arwen.ErrInvalidCallOnReadOnlyMode
		return
	}

	newContractAddress, err = blockchain.NewAddress(input.CallerAddr)
	if err != nil {
		return
	}

	if blockchain.AccountExists(newContractAddress) {
		err = arwen.ErrDeploymentOverExistingAccount
		return
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
	_, err = host.ExecuteOnDestContext(initCallInput)
	if err != nil {
		return
	}

	blockchain.IncreaseNonce(input.CallerAddr)

	return
}

func (host *vmHost) checkUpgradePermission(vmInput *vmcommon.ContractCallInput) error {
	contract, err := host.Blockchain().GetUserAccount(vmInput.RecipientAddr)
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
		log.Trace("performCodeDeployment/StartWasmerInstance", "err", err)
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

// execute executes an indirect call to a smart contract, assuming there is an
// already-running Wasmer instance of another contract that has requested the
// indirect call. This method creates a new Wasmer instance and pushes the
// previous one onto the Runtime instance stack, but it will not pop the
// previous instance back - that remains the responsibility of the calling
// code. Also, this method does not restore the gas remaining after the
// indirect call, it does not push the states of any Host Context onto their
// respective stacks, nor does it pop any state stack. Handling the state
// stacks and the remaining gas are responsibilities of the calling code, which
// must push and pop as required, before and after calling this method, and
// handle the remaining gas. These principles also apply to indirect contract
// upgrading (via host.executeUpgrade(), which also does not pop the previous
// instance from the Runtime instance stack, nor does it restore the remaining
// gas).
func (host *vmHost) execute(input *vmcommon.ContractCallInput) error {
	_, _, metering, output, runtime, _, _ := host.GetContexts()

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

	err = metering.DeductInitialGasForExecution(contract)
	if err != nil {
		return err
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

// ExecuteESDTTransfer calls the process built in function with the given transfer for ESDT/ESDTNFT if nonce > 0
// there are no NFTs with nonce == 0
func (host *vmHost) ExecuteESDTTransfer(destination []byte, sender []byte, tokenIdentifier []byte, nonce uint64, value *big.Int, callType vmcommon.CallType) (*vmcommon.VMOutput, uint64, error) {
	_, _, metering, _, runtime, _, _ := host.GetContexts()

	esdtTransferInput := &vmcommon.ContractCallInput{
		VMInput: vmcommon.VMInput{
			CallerAddr:  sender,
			Arguments:   make([][]byte, 0),
			CallValue:   big.NewInt(0),
			CallType:    callType,
			GasPrice:    runtime.GetVMInput().GasPrice,
			GasProvided: metering.GasLeft(),
			GasLocked:   0,
		},
		RecipientAddr:     destination,
		Function:          vmcommon.BuiltInFunctionESDTTransfer,
		AllowInitFunction: false,
	}

	if nonce > 0 {
		esdtTransferInput.Function = vmcommon.BuiltInFunctionESDTNFTTransfer
		esdtTransferInput.RecipientAddr = esdtTransferInput.CallerAddr
		nonceAsBytes := big.NewInt(0).SetUint64(nonce).Bytes()
		esdtTransferInput.Arguments = append(esdtTransferInput.Arguments, tokenIdentifier, nonceAsBytes, value.Bytes(), destination)
	} else {
		esdtTransferInput.Arguments = append(esdtTransferInput.Arguments, tokenIdentifier, value.Bytes())
	}

	vmOutput, err := host.Blockchain().ProcessBuiltInFunction(esdtTransferInput)
	log.Trace("ESDT transfer", "sender", sender, "dest", destination)
	log.Trace("ESDT transfer", "token", tokenIdentifier, "value", value)
	if err != nil {
		log.Trace("ESDT transfer", "error", err)
		return vmOutput, esdtTransferInput.GasProvided, err
	}
	if vmOutput.ReturnCode != vmcommon.Ok {
		log.Trace("ESDT transfer", "error", err, "retcode", vmOutput.ReturnCode, "message", vmOutput.ReturnMessage)
		return vmOutput, esdtTransferInput.GasProvided, arwen.ErrExecutionFailed
	}

	gasConsumed := math.SubUint64(esdtTransferInput.GasProvided, vmOutput.GasRemaining)
	for _, outAcc := range vmOutput.OutputAccounts {
		for _, transfer := range outAcc.OutputTransfers {
			gasConsumed = math.SubUint64(gasConsumed, transfer.GasLimit)
		}
	}
	if callType != vmcommon.AsynchronousCallBack {
		if metering.GasLeft() < gasConsumed {
			log.Trace("ESDT transfer", "error", arwen.ErrNotEnoughGas)
			return vmOutput, esdtTransferInput.GasProvided, arwen.ErrNotEnoughGas
		}
		metering.UseGas(gasConsumed)
	}

	return vmOutput, gasConsumed, nil
}

func (host *vmHost) callBuiltinFunction(input *vmcommon.ContractCallInput) (*vmcommon.ContractCallInput, *vmcommon.VMOutput, error) {
	metering := host.Metering()

	vmOutput, err := host.Blockchain().ProcessBuiltInFunction(input)
	if err != nil {
		metering.UseGas(input.GasProvided)
		return nil, nil, err
	}

	newVMInput, err := host.isSCExecutionAfterBuiltInFunc(input, vmOutput)
	if err != nil {
		metering.UseGas(input.GasProvided)
		return nil, nil, err
	}

	if newVMInput != nil {
		for _, outAcc := range vmOutput.OutputAccounts {
			outAcc.OutputTransfers = make([]vmcommon.OutputTransfer, 0)
		}
	}

	metering.TrackGasUsedByBuiltinFunction(input, vmOutput, newVMInput)

	host.addESDTTransferToVMOutputSCIntraShardCall(input, vmOutput)

	return newVMInput, vmOutput, nil
}

// add output transfer of esdt transfer when sc calling another sc intra shard to log the transfer information
func (host *vmHost) addESDTTransferToVMOutputSCIntraShardCall(
	input *vmcommon.ContractCallInput,
	output *vmcommon.VMOutput,
) {
	if output.ReturnCode != vmcommon.Ok {
		return
	}
	if !host.AreInSameShard(input.RecipientAddr, input.CallerAddr) {
		return
	}
	isESDTTransfer := input.Function == vmcommon.BuiltInFunctionESDTTransfer || input.Function == vmcommon.BuiltInFunctionESDTNFTTransfer
	if !isESDTTransfer {
		return
	}

	recipientAddr := input.RecipientAddr
	if input.Function == vmcommon.BuiltInFunctionESDTNFTTransfer {
		if len(input.Arguments) != 4 {
			return
		}
		recipientAddr = input.Arguments[3]
	}
	addOutputTransferToVMOutput(input.Function, input.Arguments, input.CallerAddr, recipientAddr, input.CallType, output)
}

func addOutputTransferToVMOutput(
	function string,
	arguments [][]byte,
	sender []byte,
	recipient []byte,
	callType vmcommon.CallType,
	vmOutput *vmcommon.VMOutput,
) {
	esdtTransferTxData := function
	for _, arg := range arguments {
		esdtTransferTxData += "@" + hex.EncodeToString(arg)
	}
	outTransfer := vmcommon.OutputTransfer{
		Value:         big.NewInt(0),
		Data:          []byte(esdtTransferTxData),
		CallType:      callType,
		SenderAddress: sender,
	}

	if len(vmOutput.OutputAccounts) == 0 {
		vmOutput.OutputAccounts = make(map[string]*vmcommon.OutputAccount)
	}
	outAcc, ok := vmOutput.OutputAccounts[string(recipient)]
	if !ok {
		outAcc = &vmcommon.OutputAccount{
			Address:         recipient,
			OutputTransfers: make([]vmcommon.OutputTransfer, 0),
		}
	}
	outAcc.OutputTransfers = append(outAcc.OutputTransfers, outTransfer)
	vmOutput.OutputAccounts[string(recipient)] = outAcc
}

func (host *vmHost) checkFinalGasAfterExit() error {
	if !host.IsArwenV2Enabled() {
		return nil
	}

	totalUsedPoints := host.Runtime().GetPointsUsed()
	if totalUsedPoints > host.Metering().GetGasForExecution() {
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
			log.Trace("UpdateCurrentCallStatus failed", "error", err)
			err = async.PostprocessCrossShardCallback()
			if err != nil {
				log.Trace("call SC method failed", "error", err)
			}
			return err
		}

		runtime.SetCustomCallFunction(asyncCall.GetCallbackName())
	}

	// TODO refactor this, and apply this condition in other places where a
	// function is called
	var err error
	if runtime.Function() != "" {
		err = host.verifyAllowedFunctionCall()
		if err != nil {
			log.Trace("call SC method failed", "error", err)
			return err
		}

		function, err := runtime.GetFunctionToCall()
		if err != nil {
			log.Trace("call SC method failed", "error", err)
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
			log.Trace("call SC method failed", "error", err)
			return err
		}

		err = async.Execute()
		if err != nil {
			log.Trace("call SC method failed", "error", err)
			return err
		}
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

	if err != nil {
		log.Trace("call SC method failed", "error", err)
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
	recipient := vmInput.RecipientAddr
	if vmInput.Function == vmcommon.BuiltInFunctionESDTNFTTransfer && bytes.Equal(vmInput.CallerAddr, vmInput.RecipientAddr) {
		recipient = vmInput.Arguments[3]
	}
	if !host.AreInSameShard(vmInput.CallerAddr, recipient) {
		return nil, nil
	}
	if !host.Blockchain().IsSmartContract(recipient) {
		return nil, nil
	}

	outAcc, ok := vmOutput.OutputAccounts[string(recipient)]
	if !ok {
		return nil, nil
	}
	if len(outAcc.OutputTransfers) != 1 {
		return nil, nil
	}

	callType := vmInput.CallType
	scCallOutTransfer := outAcc.OutputTransfers[0]

	argParser := host.CallArgsParser()
	function, arguments, err := argParser.ParseData(string(scCallOutTransfer.Data))
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
		RecipientAddr:     recipient,
		Function:          function,
		AllowInitFunction: false,
	}

	fillWithESDTValue(vmInput, newVMInput)

	return newVMInput, nil
}

func fillWithESDTValue(fullVMInput *vmcommon.ContractCallInput, newVMInput *vmcommon.ContractCallInput) {
	isESDTTransfer := fullVMInput.Function == vmcommon.BuiltInFunctionESDTTransfer || fullVMInput.Function == vmcommon.BuiltInFunctionESDTNFTTransfer
	if !isESDTTransfer {
		return
	}

	newVMInput.ESDTTokenName = fullVMInput.Arguments[0]
	newVMInput.ESDTValue = big.NewInt(0).SetBytes(fullVMInput.Arguments[1])

	if fullVMInput.Function == vmcommon.BuiltInFunctionESDTNFTTransfer {
		newVMInput.ESDTTokenNonce = big.NewInt(0).SetBytes(fullVMInput.Arguments[1]).Uint64()
		newVMInput.ESDTValue = big.NewInt(0).SetBytes(fullVMInput.Arguments[2])
		newVMInput.ESDTTokenType = uint32(vmcommon.NonFungible)
	}
}
