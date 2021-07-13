package arwen

import "bytes"

// AsyncCallGroup is a structure containing a group of async calls and a callback
// that should be called when all these async calls are resolved
type AsyncCallGroup struct {
	Callback     string
	GasLocked    uint64
	CallbackData []byte
	Identifier   string
	AsyncCalls   []*AsyncCall
}

// NewAsyncCallGroup creates a new instance of AsyncCallGroup
func NewAsyncCallGroup(identifier string) *AsyncCallGroup {
	return &AsyncCallGroup{
		Callback:     "",
		GasLocked:    0,
		CallbackData: make([]byte, 0),
		Identifier:   identifier,
		AsyncCalls:   make([]*AsyncCall, 0),
	}
}

// Clone creates a deep clone of the AsyncCallGroup
func (acg *AsyncCallGroup) Clone() *AsyncCallGroup {
	callCount := len(acg.AsyncCalls)
	clone := &AsyncCallGroup{
		Callback:   acg.Callback,
		GasLocked:  acg.GasLocked,
		Identifier: acg.Identifier,
		AsyncCalls: make([]*AsyncCall, callCount),
	}

	copy(clone.CallbackData, acg.CallbackData)

	for i := 0; i < callCount; i++ {
		clone.AsyncCalls[i] = acg.AsyncCalls[i].Clone()
	}

	return clone
}

// AddAsyncCall adds a given AsyncCall to the AsyncCallGroup
func (acg *AsyncCallGroup) AddAsyncCall(call *AsyncCall) {
	acg.AsyncCalls = append(acg.AsyncCalls, call)
}

// HasPendingCalls verifies whether the AsyncCallGroup has any AsyncCalls left
// to return from the destination call
func (acg *AsyncCallGroup) HasPendingCalls() bool {
	return len(acg.AsyncCalls) > 0
}

// IsComplete verifies whether all AsyncCalls have been completed
func (acg *AsyncCallGroup) IsComplete() bool {
	return len(acg.AsyncCalls) == 0
}

// HasCallback verifies whether a callback function has been set for this AsyncCallGroup
func (acg *AsyncCallGroup) HasCallback() bool {
	return acg.Callback != ""
}

// FindByDestination returns the index of an AsyncCall in this AsyncCallGroup
// that matches the provided destination
func (acg *AsyncCallGroup) FindByDestination(destination []byte) (int, bool) {
	for index, call := range acg.AsyncCalls {
		if bytes.Equal(destination, call.Destination) {
			return index, true
		}
	}
	return -1, false
}

// DeleteAsyncCall removes an AsyncCall from this AsyncCallGroup, given its index
func (acg *AsyncCallGroup) DeleteAsyncCall(index int) {
	asyncCalls := acg.AsyncCalls
	if len(asyncCalls) == 0 {
		return
	}

	last := len(asyncCalls) - 1
	if index < 0 || index > last {
		return
	}

	asyncCalls[index] = asyncCalls[last]
	asyncCalls = asyncCalls[:last]
	acg.AsyncCalls = asyncCalls
}

// DeleteCompletedAsyncCalls removes all completed AsyncCalls, keeping only
// those with status AsyncCallPending
func (acg *AsyncCallGroup) DeleteCompletedAsyncCalls() {
	remainingAsyncCalls := make([]*AsyncCall, 0)
	for _, asyncCall := range acg.AsyncCalls {
		if asyncCall.Status == AsyncCallPending {
			remainingAsyncCalls = append(remainingAsyncCalls, asyncCall)
		}
	}

	acg.AsyncCalls = remainingAsyncCalls
}
