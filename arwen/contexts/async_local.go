package contexts

import (
	"math/big"

	"github.com/ElrondNetwork/arwen-wasm-vm/v1_3/arwen"
	"github.com/ElrondNetwork/arwen-wasm-vm/v1_3/math"
	vmcommon "github.com/ElrondNetwork/elrond-vm-common"
)

func (context *asyncContext) executeAsyncLocalCalls() error {
	for _, group := range context.asyncCallGroups {
		for _, call := range group.AsyncCalls {
			if (call.ExecutionMode != arwen.SyncExecution) && (call.ExecutionMode != arwen.AsyncBuiltinFuncIntraShard) {
				continue
			}

			err := context.executeAsyncLocalCall(call)
			if err != nil {
				return err
			}
		}

		group.DeleteCompletedAsyncCalls()

		// If all the AsyncCalls in the AsyncCallGroup were executed synchronously,
		// then the AsyncCallGroup can have its callback executed.
		if group.IsComplete() {
			context.executeCallGroupCallback(group)
		}
	}

	context.DeleteCompletedGroups()

	if !context.HasPendingCallGroups() {
		context.executeContextCallback()
	}

	return nil
}

func (context *asyncContext) executeAsyncLocalCall(asyncCall *arwen.AsyncCall) error {
	// Briefly restore the AsyncCall GasLimit, after it was consumed in its
	// entirety by addAsyncCall(); this is required, because ExecuteOnDestContext()
	// must also consume the GasLimit in its entirety, before starting execution,
	// but will restore any GasRemaining to the current instance.
	metering := context.host.Metering()
	metering.RestoreGas(asyncCall.GetGasLimit())

	destinationCallInput, err := context.createContractCallInput(asyncCall)
	if err != nil {
		return err
	}

	vmOutput, err := context.host.ExecuteOnDestContext(destinationCallInput)

	// The vmOutput instance returned by host.ExecuteOnDestContext() is never nil,
	// by design. Using it without checking for err is safe here.
	asyncCall.UpdateStatus(vmOutput.ReturnCode)

	if asyncCall.HasCallback() {
		callbackVMOutput, callbackErr := context.executeSyncCallback(asyncCall, vmOutput, err)
		context.finishAsyncLocalExecution(callbackVMOutput, callbackErr)
	}

	// TODO accumulate remaining gas from the callback into the AsyncContext,
	// after fixing the bug caught by TestExecution_ExecuteOnDestContext_GasRemaining().

	return nil
}

func (context *asyncContext) executeSyncCallback(
	asyncCall *arwen.AsyncCall,
	destinationVMOutput *vmcommon.VMOutput,
	destinationErr error,
) (*vmcommon.VMOutput, error) {
	metering := context.host.Metering()

	callbackInput, err := context.createCallbackInput(asyncCall, destinationVMOutput, destinationErr)
	if err != nil {
		return nil, err
	}

	// Restore gas locked while still on the caller instance; otherwise, the
	// locked gas will appear to have been used twice by the caller instance.
	metering.RestoreGas(asyncCall.GetGasLocked())
	callbackVMOutput, callBackErr := context.host.ExecuteOnDestContext(callbackInput)

	return callbackVMOutput, callBackErr
}

// executeCallGroupCallback synchronously executes the designated callback of
// the AsyncCallGroup, as it was set with SetGroupCallback().
//
// Gas for the execution has been already paid for when SetGroupCallback() was
// set. The remaining gas is refunded to context.callerAddr, which initiated
// the call and paid for the gas in the first place.
func (context *asyncContext) executeCallGroupCallback(group *arwen.AsyncCallGroup) {
	if !group.HasCallback() {
		return
	}

	input := context.createGroupCallbackInput(group)
	vmOutput, err := context.host.ExecuteOnDestContext(input)
	context.finishAsyncLocalExecution(vmOutput, err)
}

// executeSyncHalfOfBuiltinFunction will synchronously call the requested
// built-in function. This is required for all cross-shard calls to built-in
// functions, because they will handle cross-shard calls themselves, by
// generating entries in vmOutput.OutputAccounts, and they need to be executed
// synchronously to do that. As a consequence, it is not necessary to call
// sendAsyncCallCrossShard(). The vmOutput produced by the built-in function,
// containing the cross-shard call, has ALREADY been merged into the main
// output by the inner call to host.ExecuteOnDestContext(). Moreover, the
// status of the AsyncCall is not updated here - it will be updated by
// PostprocessCrossShardCallback(), when the cross-shard call returns.
func (context *asyncContext) executeSyncHalfOfBuiltinFunction(asyncCall *arwen.AsyncCall) error {
	destinationCallInput, err := context.createContractCallInput(asyncCall)
	if err != nil {
		return err
	}

	vmOutput, err := context.host.ExecuteOnDestContext(destinationCallInput)
	if err != nil {
		return err
	}

	// If the in-shard half of the built-in function call has failed, go no
	// further and execute the error callback of this AsyncCall.
	if vmOutput.ReturnCode != vmcommon.Ok {
		asyncCall.Reject()
		callbackVMOutput, callbackErr := context.executeSyncCallback(asyncCall, vmOutput, err)
		context.finishAsyncLocalExecution(callbackVMOutput, callbackErr)
	}

	// The gas that remains after executing the in-shard half of the built-in
	// function is provided to the cross-shard half.
	asyncCall.GasLimit = vmOutput.GasRemaining

	return nil
}

// executeSyncContextCallback will execute the callback of the original caller
// synchronously, already assuming the original caller is in the same shard
func (context *asyncContext) executeSyncContextCallback() {
	callbackCallInput := context.createContextCallbackInput()
	callbackVMOutput, callBackErr := context.host.ExecuteOnDestContext(callbackCallInput)
	context.finishAsyncLocalExecution(callbackVMOutput, callBackErr)
}

// TODO return values are never used by code that calls finishAsyncLocalExecution
func (context *asyncContext) finishAsyncLocalExecution(vmOutput *vmcommon.VMOutput, err error) {
	if err == nil {
		return
	}

	runtime := context.host.Runtime()
	output := context.host.Output()

	runtime.GetVMInput().GasProvided = 0

	if vmOutput == nil {
		vmOutput = output.CreateVMOutputInCaseOfError(err)
	}

	// TODO Discuss consistency between in-shard and cross-shard results
	// TODO of the callback, and how they're accessible to the caller / user.
	// TODO Currently, a failed callback in-shard leaves the ReturnCode to
	// TODO vmcommon.Ok, unless the following line is uncommented.
	// output.SetReturnCode(vmOutput.ReturnCode)

	output.SetReturnMessage(vmOutput.ReturnMessage)
	output.Finish([]byte(vmOutput.ReturnCode.String()))
	output.Finish(runtime.GetCurrentTxHash())
}

func (context *asyncContext) createContractCallInput(asyncCall *arwen.AsyncCall) (*vmcommon.ContractCallInput, error) {
	host := context.host
	runtime := host.Runtime()
	sender := runtime.GetSCAddress()

	function, arguments, err := host.CallArgsParser().ParseData(string(asyncCall.GetData()))
	if err != nil {
		return nil, err
	}

	gasLimit := asyncCall.GetGasLimit()
	gasToUse := host.Metering().GasSchedule().ElrondAPICost.AsyncCallStep
	if gasLimit <= gasToUse {
		return nil, arwen.ErrNotEnoughGas
	}
	gasLimit -= gasToUse

	contractCallInput := &vmcommon.ContractCallInput{
		VMInput: vmcommon.VMInput{
			CallerAddr:     sender,
			Arguments:      arguments,
			CallValue:      big.NewInt(0).SetBytes(asyncCall.GetValue()),
			CallType:       vmcommon.AsynchronousCall,
			GasPrice:       runtime.GetVMInput().GasPrice,
			GasProvided:    gasLimit,
			GasLocked:      asyncCall.GetGasLocked(),
			CurrentTxHash:  runtime.GetCurrentTxHash(),
			OriginalTxHash: runtime.GetOriginalTxHash(),
			PrevTxHash:     runtime.GetPrevTxHash(),
		},
		RecipientAddr: asyncCall.GetDestination(),
		Function:      function,
	}

	return contractCallInput, nil
}

// TODO function too large; refactor needed
func (context *asyncContext) createCallbackInput(
	asyncCall *arwen.AsyncCall,
	vmOutput *vmcommon.VMOutput,
	destinationErr error,
) (*vmcommon.ContractCallInput, error) {
	metering := context.host.Metering()
	runtime := context.host.Runtime()

	// always provide return code as the first argument to callback function
	arguments := [][]byte{
		big.NewInt(int64(vmOutput.ReturnCode)).Bytes(),
	}

	esdtFunction := ""
	isESDTOnCallBack := false
	esdtArgs := make([][]byte, 0)
	returnWithError := false
	if destinationErr == nil && vmOutput.ReturnCode == vmcommon.Ok {
		// when execution went Ok, callBack arguments are:
		// [0, result1, result2, ....]
		isESDTOnCallBack, esdtFunction, esdtArgs = context.isESDTTransferOnReturnDataWithNoAdditionalData(vmOutput.ReturnData)
		arguments = append(arguments, vmOutput.ReturnData...)
	} else {
		// when execution returned error, callBack arguments are:
		// [error code, error message]
		arguments = append(arguments, []byte(vmOutput.ReturnMessage))
		returnWithError = true
	}

	callbackFunction := asyncCall.GetCallbackName()

	gasLimit := math.AddUint64(vmOutput.GasRemaining, asyncCall.GetGasLocked())
	dataLength := computeDataLengthFromArguments(callbackFunction, arguments)

	gasToUse := metering.GasSchedule().ElrondAPICost.AsyncCallStep
	copyPerByte := metering.GasSchedule().BaseOperationCost.DataCopyPerByte
	gas := math.MulUint64(copyPerByte, uint64(dataLength))
	gasToUse = math.AddUint64(gasToUse, gas)
	if gasLimit <= gasToUse {
		return nil, arwen.ErrNotEnoughGas
	}
	gasLimit -= gasToUse

	// Return to the sender SC, calling its specified callback method.
	contractCallInput := &vmcommon.ContractCallInput{
		VMInput: vmcommon.VMInput{
			CallerAddr:           asyncCall.Destination,
			Arguments:            arguments,
			CallValue:            context.computeCallValueFromVMOutput(vmOutput),
			CallType:             vmcommon.AsynchronousCallBack,
			GasPrice:             runtime.GetVMInput().GasPrice,
			GasProvided:          gasLimit,
			GasLocked:            0,
			CurrentTxHash:        runtime.GetCurrentTxHash(),
			OriginalTxHash:       runtime.GetOriginalTxHash(),
			PrevTxHash:           runtime.GetPrevTxHash(),
			ReturnCallAfterError: returnWithError,
		},
		RecipientAddr: runtime.GetSCAddress(),
		Function:      callbackFunction,
	}

	if isESDTOnCallBack {
		contractCallInput.Function = esdtFunction
		contractCallInput.Arguments = make([][]byte, 0, len(arguments))
		contractCallInput.Arguments = append(contractCallInput.Arguments, esdtArgs[0], esdtArgs[1])
		if esdtFunction == vmcommon.BuiltInFunctionESDTNFTTransfer {
			contractCallInput.Arguments = append(contractCallInput.Arguments, esdtArgs[2], esdtArgs[3])
		}
		contractCallInput.Arguments = append(contractCallInput.Arguments, []byte(callbackFunction))
		contractCallInput.Arguments = append(contractCallInput.Arguments, big.NewInt(int64(vmOutput.ReturnCode)).Bytes())
		if len(vmOutput.ReturnData) > 1 {
			contractCallInput.Arguments = append(contractCallInput.Arguments, vmOutput.ReturnData[1:]...)
		}
	}

	return contractCallInput, nil
}

func (context *asyncContext) createGroupCallbackInput(group *arwen.AsyncCallGroup) *vmcommon.ContractCallInput {
	runtime := context.host.Runtime()
	input := &vmcommon.ContractCallInput{
		VMInput: vmcommon.VMInput{
			CallerAddr:     context.callerAddr,
			Arguments:      [][]byte{group.CallbackData},
			CallValue:      big.NewInt(0),
			GasPrice:       context.gasPrice,
			GasProvided:    group.GasLocked,
			CurrentTxHash:  runtime.GetCurrentTxHash(),
			OriginalTxHash: runtime.GetOriginalTxHash(),
			PrevTxHash:     runtime.GetPrevTxHash(),
		},
		RecipientAddr: runtime.GetSCAddress(),
		Function:      group.Callback,
	}

	return input
}

func (context *asyncContext) createContextCallbackInput() *vmcommon.ContractCallInput {
	host := context.host
	runtime := host.Runtime()

	_, arguments, err := host.CallArgsParser().ParseData(string(context.returnData))
	if err != nil {
		arguments = [][]byte{context.returnData}
	}

	// TODO ensure a new value for VMInput.CurrentTxHash
	input := &vmcommon.ContractCallInput{
		VMInput: vmcommon.VMInput{
			CallerAddr:     context.callerAddr,
			Arguments:      arguments,
			CallValue:      runtime.GetVMInput().CallValue,
			CallType:       vmcommon.AsynchronousCallBack,
			GasPrice:       runtime.GetVMInput().GasPrice,
			GasProvided:    context.gasRemaining,
			CurrentTxHash:  runtime.GetCurrentTxHash(),
			OriginalTxHash: runtime.GetOriginalTxHash(),
			PrevTxHash:     runtime.GetPrevTxHash(),
		},
		RecipientAddr: runtime.GetSCAddress(),
		Function:      context.callback,
	}
	return input
}

func (context *asyncContext) isESDTTransferOnReturnDataWithNoAdditionalData(data [][]byte) (bool, string, [][]byte) {
	if len(data) == 0 {
		return false, "", nil
	}

	argParser := context.host.CallArgsParser()
	functionName, args, err := argParser.ParseData(string(data[0]))
	if err != nil {
		return false, "", nil
	}

	return isESDTTransferOnReturnDataFromFunctionAndArgs(functionName, args)
}

func isESDTTransferOnReturnDataFromFunctionAndArgs(functionName string, args [][]byte) (bool, string, [][]byte) {
	if functionName == vmcommon.BuiltInFunctionESDTTransfer && len(args) == 2 {
		return true, functionName, args
	}

	if functionName == vmcommon.BuiltInFunctionESDTNFTTransfer && len(args) == 4 {
		return true, functionName, args
	}

	return false, functionName, args
}

func (context *asyncContext) computeCallValueFromVMOutput(destinationVMOutput *vmcommon.VMOutput) *big.Int {
	if len(destinationVMOutput.ReturnData) > 0 {
		return big.NewInt(0)
	}

	returnTransfer := big.NewInt(0)
	callBackReceiver := context.host.Runtime().GetSCAddress()
	outAcc, ok := destinationVMOutput.OutputAccounts[string(callBackReceiver)]
	if !ok {
		return returnTransfer
	}

	if len(outAcc.OutputTransfers) == 0 {
		return returnTransfer
	}

	lastOutTransfer := outAcc.OutputTransfers[len(outAcc.OutputTransfers)-1]
	if len(lastOutTransfer.Data) == 0 {
		returnTransfer.Set(lastOutTransfer.Value)
	}

	return returnTransfer
}
