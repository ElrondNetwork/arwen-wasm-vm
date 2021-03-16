package contexts

type bytesMap map[int32][]byte

type bytesHeapContext struct {
	values     bytesMap
	stateStack []bytesMap
}

// NewBytesHeapContext creates a new bytesHeapContext
func NewBytesHeapContext() (*bytesHeapContext, error) {
	context := &bytesHeapContext{
		values:     make(bytesMap),
		stateStack: make([]bytesMap, 0),
	}

	return context, nil
}

// InitState initializes the underlying values map
func (context *bytesHeapContext) InitState() {
	context.values = make(bytesMap)
}

// PushState appends the values map to the state stack
func (context *bytesHeapContext) PushState() {
	newState := context.clone()
	context.stateStack = append(context.stateStack, newState)
}

// PopSetActiveState removes the latest entry from the state stack and sets it as the current values map
func (context *bytesHeapContext) PopSetActiveState() {
	stateStackLen := len(context.stateStack)
	if stateStackLen == 0 {
		return
	}

	prevValues := context.stateStack[stateStackLen-1]
	context.stateStack = context.stateStack[:stateStackLen-1]

	context.values = prevValues
}

// PopDiscard removes the latest entry from the state stack
func (context *bytesHeapContext) PopDiscard() {
	stateStackLen := len(context.stateStack)
	if stateStackLen == 0 {
		return
	}

	context.stateStack = context.stateStack[:stateStackLen-1]
}

// ClearStateStack initializes the state stack
func (context *bytesHeapContext) ClearStateStack() {
	context.stateStack = make([]bytesMap, 0)
}

func (context *bytesHeapContext) clone() bytesMap {
	newState := make(bytesMap, len(context.values))
	for handle, bytes := range context.values {
		newState[handle] = make([]byte, len(bytes))
		copy(newState[handle], bytes)
	}
	return newState
}

// NewByteBuffer adds the given value to the current values map and returns the handle
func (context *bytesHeapContext) NewByteBuffer(value []byte) int32 {
	newHandle := int32(len(context.values))
	for {
		if _, ok := context.values[newHandle]; !ok {
			break
		}
		newHandle++
	}

	context.values[newHandle] = value

	return newHandle
}

// GetByteBuffer returns the value at the given handle. If there is no value under that handle, it will return an empty buffer
func (context *bytesHeapContext) GetByteBuffer(handle int32) []byte {
	if _, ok := context.values[handle]; !ok {
		context.values[handle] = make([]byte, 0)
	}

	return context.values[handle]
}

// SetByteBuffer writes the given bytes at the given handle
func (context *bytesHeapContext) SetByteBuffer(handle int32, value []byte) {
	context.values[handle] = value
}

// IsInterfaceNil returns true if there is no value under the interface
func (context *bytesHeapContext) IsInterfaceNil() bool {
	return context == nil
}
