package mock

import (
	"testing"

	"github.com/ElrondNetwork/arwen-wasm-vm/wasmer"
	"github.com/stretchr/testify/require"
)

// InstanceBuilderMock can be passed to RuntimeContext as an InstanceBuilder to
// create mocked Wasmer instances.
type InstanceBuilderMock struct {
	tb          testing.TB
	InstanceMap map[string]*wasmer.Instance
	DefaultCode []byte
}

// NewInstanceBuilderMock constructs a new InstanceBuilderMock
func NewInstanceBuilderMock(tb testing.TB, defaultCode []byte) *InstanceBuilderMock {
	return &InstanceBuilderMock{
		tb:          tb,
		InstanceMap: make(map[string]*wasmer.Instance),
		DefaultCode: defaultCode,
	}
}

// CreateAndStoreInstanceMock creates a real Wasmer instance using the
// DefaultCode and returns it, so that a test may alter its Exports map to
// inject new contract methods; afterwards, the RuntimeContext will call
// NewInstanceWithOptions() and obtain the instance with the injected methods.
//
// It is necessary to call CreateAndStoreInstanceMock() for any contract that is
// to be called, or at least manually populate the InstanceMap appropriately
// (real WASM contracts may be used to populate it, as well).
func (builder *InstanceBuilderMock) CreateAndStoreInstanceMock(code []byte) wasmer.InstanceHandler {
	options := wasmer.CompilationOptions{
		GasLimit:           1000000000,
		UnmeteredLocals:    4000,
		OpcodeTrace:        false,
		Metering:           true,
		RuntimeBreakpoints: true,
	}

	instance, err := wasmer.NewInstanceWithOptions(builder.DefaultCode, options)
	require.NotNil(builder.tb, instance)
	require.Nil(builder.tb, err)
	builder.InstanceMap[string(code)] = instance

	return instance
}

// GetStoredInstanceMock retrieves and initializes a stored Wasmer instance, or
// nil if it doesn't exist
func (builder *InstanceBuilderMock) GetStoredInstanceMock(code []byte, gasLimit uint64) (wasmer.InstanceHandler, bool) {
	instance, ok := builder.InstanceMap[string(code)]
	if ok {
		instance.SetPointsUsed(0)
		instance.SetGasLimit(gasLimit)
		return instance, true
	}
	return nil, false
}

// NewInstanceWithOptions attempts to load a prepared instance using
// GetStoredInstanceMock; if it doesn't exist, it creates a true Wasmer
// instance with the provided contract code.
func (builder *InstanceBuilderMock) NewInstanceWithOptions(
	contractCode []byte,
	options wasmer.CompilationOptions,
) (wasmer.InstanceHandler, error) {

	instance, ok := builder.GetStoredInstanceMock(contractCode, options.GasLimit)
	if ok {
		return instance, nil
	}
	return wasmer.NewInstanceWithOptions(contractCode, options)
}

// NewInstanceFromCompiledCodeWithOptions attempts to load a prepared instance
// using GetStoredInstanceMock; if it doesn't exist, it creates a true Wasmer
// instance with the provided precompiled code.
func (builder *InstanceBuilderMock) NewInstanceFromCompiledCodeWithOptions(
	compiledCode []byte,
	options wasmer.CompilationOptions,
) (wasmer.InstanceHandler, error) {
	instance, ok := builder.GetStoredInstanceMock(compiledCode, options.GasLimit)
	if ok {
		return instance, nil
	}
	return wasmer.NewInstanceFromCompiledCodeWithOptions(compiledCode, options)
}
