package contexts

import (
	"github.com/ElrondNetwork/arwen-wasm-vm/wasmer"
)

type wasmerInstanceBuilder struct {
}

// NewInstanceWithOptions creates a new Wasmer instance from WASM bytecode,
// respecting the provided options
func (builder *wasmerInstanceBuilder) NewInstanceWithOptions(
	contractCode []byte,
	options wasmer.CompilationOptions,
) (wasmer.InstanceHandler, error) {
	return wasmer.NewInstanceWithOptions(contractCode, options)
}

// NewInstanceFromCompiledCodeWithOptions creates a new Wasmer instance from
// precompiled machine code, respecting the provided options
func (builder *wasmerInstanceBuilder) NewInstanceFromCompiledCodeWithOptions(
	compiledCode []byte,
	options wasmer.CompilationOptions,
) (wasmer.InstanceHandler, error) {
	return wasmer.NewInstanceFromCompiledCodeWithOptions(compiledCode, options)
}
