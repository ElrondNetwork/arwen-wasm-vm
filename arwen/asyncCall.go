package arwen

import (
	"github.com/ElrondNetwork/arwen-wasm-vm/math"
	"github.com/ElrondNetwork/elrond-go/core/vmcommon"
)

// AsyncCall holds the information about an individual async call
type AsyncCall struct {
	Status          AsyncCallStatus
	ExecutionMode   AsyncCallExecutionMode
	Destination     []byte
	Data            []byte
	GasLimit        uint64
	GasLocked       uint64
	ValueBytes      []byte
	SuccessCallback string
	ErrorCallback   string
}

// Clone creates a deep clone of the AsyncCall
func (ac *AsyncCall) Clone() *AsyncCall {
	clone := &AsyncCall{
		Status:          ac.Status,
		ExecutionMode:   ac.ExecutionMode,
		Destination:     make([]byte, len(ac.Destination)),
		Data:            make([]byte, len(ac.Data)),
		GasLimit:        ac.GasLimit,
		GasLocked:       ac.GasLocked,
		ValueBytes:      make([]byte, len(ac.ValueBytes)),
		SuccessCallback: ac.SuccessCallback,
		ErrorCallback:   ac.ErrorCallback,
	}

	copy(clone.Destination, ac.Destination)
	copy(clone.Data, ac.Data)
	copy(clone.ValueBytes, ac.ValueBytes)

	return clone
}

// GetDestination returns the destination of an async call
func (ac *AsyncCall) GetDestination() []byte {
	return ac.Destination
}

// GetData returns the transaction data of the async call
func (ac *AsyncCall) GetData() []byte {
	return ac.Data
}

// GetGasLimit returns the gas limit of the current async call
func (ac *AsyncCall) GetGasLimit() uint64 {
	return ac.GasLimit
}

// GetGasLocked returns the gas locked for the async callback
func (ac *AsyncCall) GetGasLocked() uint64 {
	return ac.GasLocked
}

// GetTotalGas returns the sum of the gas limit and gas locked
func (ac *AsyncCall) GetTotalGas() uint64 {
	return math.AddUint64(ac.GasLimit, ac.GasLocked)
}

// GetValue returns the byte representation of the value of the async call
func (ac *AsyncCall) GetValue() []byte {
	return ac.ValueBytes
}

// IsInterfaceNil returns true if there is no value under the interface
func (ac *AsyncCall) IsInterfaceNil() bool {
	return ac == nil
}

// UpdateStatus sets the status of the async call depending on the provided ReturnCode
func (ac *AsyncCall) UpdateStatus(returnCode vmcommon.ReturnCode) {
	ac.Status = AsyncCallResolved
	if returnCode != vmcommon.Ok {
		ac.Status = AsyncCallRejected
	}
}

// GetCallbackName returns the name of the callback to execute, depending on
// the status of the async call
func (ac *AsyncCall) GetCallbackName() string {
	if ac.Status == AsyncCallResolved {
		return ac.SuccessCallback
	}

	return ac.ErrorCallback
}
