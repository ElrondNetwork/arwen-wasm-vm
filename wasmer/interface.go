package wasmer

// InstanceHandler defines the functionality for working with WASM Instances
type InstanceHandler interface {
	HasMemory() bool
	SetContextData(data int)
	GetPointsUsed() uint64
	SetPointsUsed(points uint64)
	SetGasLimit(gasLimit uint64)
	SetBreakpointValue(value uint64)
	GetBreakpointValue() uint64
	Cache() ([]byte, error)
	Clean()
	GetExports() ExportsMap
	GetSignature(functionName string) (*ExportedFunctionSignature, bool)
	GetData() *int
	GetInstanceCtxMemory() *Memory
	GetMemory() *Memory
}
