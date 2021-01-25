package contexts

import (
	"strings"
	"testing"

	"github.com/ElrondNetwork/arwen-wasm-vm/arwen"
	"github.com/ElrondNetwork/arwen-wasm-vm/wasmer"
	"github.com/ElrondNetwork/elrond-go/core/vmcommon"
	"github.com/stretchr/testify/require"
)

func TestFunctionsGuard_isValidFunctionName(t *testing.T) {
	imports := MakeAPIImports()

	protocolBuiltinFunctions := vmcommon.FunctionNames{
		"fromProtocolFoo": {},
		"fromProtocolBar": {},
	}

	validator := newWASMValidator(imports.Names(), protocolBuiltinFunctions)

	require.Nil(t, validator.verifyValidFunctionName("foo"))
	require.Nil(t, validator.verifyValidFunctionName("_"))
	require.Nil(t, validator.verifyValidFunctionName("a"))
	require.Nil(t, validator.verifyValidFunctionName("i"))

	require.NotNil(t, validator.verifyValidFunctionName(""))
	require.NotNil(t, validator.verifyValidFunctionName("â"))
	require.NotNil(t, validator.verifyValidFunctionName("ș"))
	require.NotNil(t, validator.verifyValidFunctionName("Ä"))

	require.NotNil(t, validator.verifyValidFunctionName("fromProtocolFoo"))
	require.NotNil(t, validator.verifyValidFunctionName("fromProtocolBar"))

	require.Nil(t, validator.verifyValidFunctionName(strings.Repeat("_", 255)))
	require.NotNil(t, validator.verifyValidFunctionName(strings.Repeat("_", 256)))

	require.NotNil(t, validator.verifyValidFunctionName("getArgument"))
	require.NotNil(t, validator.verifyValidFunctionName("asyncCall"))
	require.Nil(t, validator.verifyValidFunctionName("getArgument55"))
}

func TestFunctionsGuard_Arity(t *testing.T) {
	host := InitializeArwenAndWasmer()
	imports := host.SCAPIMethods

	validator := newWASMValidator(imports.Names(), make(vmcommon.FunctionNames))

	gasLimit := uint64(100000000)
	path := "./../../test/contracts/signatures/output/signatures.wasm"
	contractCode := arwen.GetSCCode(path)
	options := wasmer.CompilationOptions{
		GasLimit:           gasLimit,
		OpcodeTrace:        false,
		Metering:           true,
		RuntimeBreakpoints: true,
		ImportObjectIndex:  0,
	}
	instance, err := wasmer.NewInstanceWithOptions(contractCode, options)
	require.Nil(t, err)

	inArity, _ := validator.getInputArity(instance, "goodFunction")
	require.Equal(t, 0, inArity)

	outArity, _ := validator.getOutputArity(instance, "goodFunction")
	require.Equal(t, 0, outArity)

	inArity, _ = validator.getInputArity(instance, "wrongReturn")
	require.Equal(t, 0, inArity)

	outArity, _ = validator.getOutputArity(instance, "wrongReturn")
	require.Equal(t, 1, outArity)

	inArity, _ = validator.getInputArity(instance, "wrongParams")
	require.Equal(t, 1, inArity)

	outArity, _ = validator.getOutputArity(instance, "wrongParams")
	require.Equal(t, 0, outArity)

	inArity, _ = validator.getInputArity(instance, "wrongParamsAndReturn")
	require.Equal(t, 2, inArity)

	outArity, _ = validator.getOutputArity(instance, "wrongParamsAndReturn")
	require.Equal(t, 1, outArity)

	err = validator.verifyVoidFunction(instance, "goodFunction")
	require.Nil(t, err)

	err = validator.verifyVoidFunction(instance, "wrongReturn")
	require.NotNil(t, err)

	err = validator.verifyVoidFunction(instance, "wrongParams")
	require.NotNil(t, err)

	err = validator.verifyVoidFunction(instance, "wrongParamsAndReturn")
	require.NotNil(t, err)
}
