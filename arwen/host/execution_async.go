package host

import (
	"encoding/hex"

	"github.com/ElrondNetwork/arwen-wasm-vm/arwen"
	"github.com/ElrondNetwork/elrond-go/core/vmcommon"
)

func (host *vmHost) handleAsyncCallBreakpoint() error {
	runtime := host.Runtime()
	async := host.Async()
	runtime.SetRuntimeBreakpointValue(arwen.BreakpointNone)

	legacyGroupID := arwen.LegacyAsyncCallGroupID
	legacyGroup, exists := async.GetCallGroup(legacyGroupID)
	if !exists {
		return arwen.ErrLegacyAsyncCallNotFound

	}

	if legacyGroup.IsComplete() {
		return arwen.ErrLegacyAsyncCallInvalid
	}

	return nil
}

func (host *vmHost) sendAsyncCallbackToCaller() error {
	if !host.Async().IsComplete() {
		return nil
	}

	runtime := host.Runtime()
	output := host.Output()
	metering := host.Metering()
	currentCall := runtime.GetVMInput()

	retData := []byte("@" + hex.EncodeToString([]byte(output.ReturnCode().String())))
	for _, data := range output.ReturnData() {
		retData = append(retData, []byte("@"+hex.EncodeToString(data))...)
	}

	err := output.Transfer(
		currentCall.CallerAddr,
		runtime.GetSCAddress(),
		metering.GasLeft(),
		0,
		currentCall.CallValue,
		retData,
		vmcommon.AsynchronousCallBack,
	)
	if err != nil {
		metering.UseGas(metering.GasLeft())
		runtime.FailExecution(err)
		return err
	}

	return nil
}
