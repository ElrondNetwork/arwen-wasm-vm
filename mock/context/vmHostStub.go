package mock

import (
	"math/big"

	"github.com/ElrondNetwork/arwen-wasm-vm/v1_3/arwen"
	"github.com/ElrondNetwork/arwen-wasm-vm/v1_3/config"
	"github.com/ElrondNetwork/arwen-wasm-vm/v1_3/crypto"
	"github.com/ElrondNetwork/arwen-wasm-vm/v1_3/wasmer"
	"github.com/ElrondNetwork/elrond-vm-common"
)

var _ arwen.VMHost = (*VMHostStub)(nil)

// VMHostStub is used in tests to check the VMHost interface method calls
type VMHostStub struct {
	InitStateCalled       func()
	PushStateCalled       func()
	PopStateCalled        func()
	ClearStateStackCalled func()
	GetVersionCalled      func() string

	CryptoCalled                      func() crypto.VMCrypto
	BlockchainCalled                  func() arwen.BlockchainContext
	RuntimeCalled                     func() arwen.RuntimeContext
	ManagedTypesCalled                func() arwen.ManagedTypesContext
	OutputCalled                      func() arwen.OutputContext
	MeteringCalled                    func() arwen.MeteringContext
	StorageCalled                     func() arwen.StorageContext
	ExecuteESDTTransferCalled         func(destination []byte, sender []byte, tokenIdentifier []byte, nonce uint64, value *big.Int, callType vmcommon.CallType) (*vmcommon.VMOutput, uint64, error)
	CreateNewContractCalled           func(input *vmcommon.ContractCreateInput) ([]byte, error)
	ExecuteOnSameContextCalled        func(input *vmcommon.ContractCallInput) (*arwen.AsyncContextInfo, error)
	ExecuteOnDestContextCalled        func(input *vmcommon.ContractCallInput) (*vmcommon.VMOutput, *arwen.AsyncContextInfo, error)
	GetAPIMethodsCalled               func() *wasmer.Imports
	GetProtocolBuiltinFunctionsCalled func() vmcommon.FunctionNames
	SetProtocolBuiltinFunctionsCalled func(vmcommon.FunctionNames)
	IsBuiltinFunctionNameCalled       func(functionName string) bool
	AreInSameShardCalled              func(left []byte, right []byte) bool

	RunSmartContractCallCalled   func(input *vmcommon.ContractCallInput) (vmOutput *vmcommon.VMOutput, err error)
	RunSmartContractCreateCalled func(input *vmcommon.ContractCreateInput) (vmOutput *vmcommon.VMOutput, err error)
	GetGasScheduleMapCalled      func() config.GasScheduleMap
	GasScheduleChangeCalled      func(newGasSchedule config.GasScheduleMap)
	IsInterfaceNilCalled         func() bool

	SetRuntimeContextCalled func(runtime arwen.RuntimeContext)
	GetContextsCalled       func() (arwen.ManagedTypesContext, arwen.BlockchainContext, arwen.MeteringContext, arwen.OutputContext, arwen.RuntimeContext, arwen.StorageContext)
}

// GetVersion mocked method
func (vhs *VMHostStub) GetVersion() string {
	if vhs.GetVersionCalled != nil {
		return vhs.GetVersionCalled()
	}

	return "stub"
}

// InitState mocked method
func (vhs *VMHostStub) InitState() {
	if vhs.InitStateCalled != nil {
		vhs.InitStateCalled()
	}
}

// PushState mocked method
func (vhs *VMHostStub) PushState() {
	if vhs.PushStateCalled != nil {
		vhs.PushStateCalled()
	}
}

// PopState mocked method
func (vhs *VMHostStub) PopState() {
	if vhs.PopStateCalled != nil {
		vhs.PopStateCalled()
	}
}

// ClearStateStack mocked method
func (vhs *VMHostStub) ClearStateStack() {
	if vhs.ClearStateStackCalled != nil {
		vhs.ClearStateStackCalled()
	}
}

// Crypto mocked method
func (vhs *VMHostStub) Crypto() crypto.VMCrypto {
	if vhs.CryptoCalled != nil {
		return vhs.CryptoCalled()
	}
	return nil
}

// Blockchain mocked method
func (vhs *VMHostStub) Blockchain() arwen.BlockchainContext {
	if vhs.BlockchainCalled != nil {
		return vhs.BlockchainCalled()
	}
	return nil
}

// Runtime mocked method
func (vhs *VMHostStub) Runtime() arwen.RuntimeContext {
	if vhs.RuntimeCalled != nil {
		return vhs.RuntimeCalled()
	}
	return nil
}

// BigInt mocked method
func (vhs *VMHostStub) ManagedTypes() arwen.ManagedTypesContext {
	if vhs.ManagedTypesCalled != nil {
		return vhs.ManagedTypesCalled()
	}
	return nil
}

// IsArwenV2Enabled mocked method
func (vhs *VMHostStub) IsArwenV2Enabled() bool {
	return true
}

// IsArwenV3Enabled mocked method
func (vhs *VMHostStub) IsArwenV3Enabled() bool {
	return true
}

// IsAheadOfTimeCompileEnabled mocked method
func (vhs *VMHostStub) IsAheadOfTimeCompileEnabled() bool {
	return true
}

// IsDynamicGasLockingEnabled mocked method
func (vhs *VMHostStub) IsDynamicGasLockingEnabled() bool {
	return true
}

// IsESDTFunctionsEnabled mocked method
func (vhs *VMHostStub) IsESDTFunctionsEnabled() bool {
	return true
}

// Output mocked method
func (vhs *VMHostStub) Output() arwen.OutputContext {
	if vhs.OutputCalled != nil {
		return vhs.OutputCalled()
	}
	return nil
}

// Metering mocked method
func (vhs *VMHostStub) Metering() arwen.MeteringContext {
	if vhs.MeteringCalled != nil {
		return vhs.MeteringCalled()
	}
	return nil
}

// Storage mocked method
func (vhs *VMHostStub) Storage() arwen.StorageContext {
	if vhs.StorageCalled != nil {
		return vhs.StorageCalled()
	}
	return nil
}

// ExecuteESDTTransfer mocked method
func (vhs *VMHostStub) ExecuteESDTTransfer(destination []byte, sender []byte, tokenIdentifier []byte, nonce uint64, value *big.Int, callType vmcommon.CallType) (*vmcommon.VMOutput, uint64, error) {
	if vhs.ExecuteESDTTransferCalled != nil {
		return vhs.ExecuteESDTTransferCalled(destination, sender, tokenIdentifier, nonce, value, callType)
	}
	return nil, 0, nil
}

// CreateNewContract mocked method
func (vhs *VMHostStub) CreateNewContract(input *vmcommon.ContractCreateInput) ([]byte, error) {
	if vhs.CreateNewContractCalled != nil {
		return vhs.CreateNewContractCalled(input)
	}
	return nil, nil
}

// ExecuteOnSameContext mocked method
func (vhs *VMHostStub) ExecuteOnSameContext(input *vmcommon.ContractCallInput) (*arwen.AsyncContextInfo, error) {
	if vhs.ExecuteOnSameContextCalled != nil {
		return vhs.ExecuteOnSameContextCalled(input)
	}
	return nil, nil
}

// ExecuteOnDestContext mocked method
func (vhs *VMHostStub) ExecuteOnDestContext(input *vmcommon.ContractCallInput) (*vmcommon.VMOutput, *arwen.AsyncContextInfo, error) {
	if vhs.ExecuteOnDestContextCalled != nil {
		return vhs.ExecuteOnDestContextCalled(input)
	}
	return nil, nil, nil
}

// AreInSameShard mocked method
func (vhs *VMHostStub) AreInSameShard(left []byte, right []byte) bool {
	if vhs.AreInSameShardCalled != nil {
		return vhs.AreInSameShardCalled(left, right)
	}
	return true
}

// GetAPIMethods mocked method
func (vhs *VMHostStub) GetAPIMethods() *wasmer.Imports {
	if vhs.GetAPIMethodsCalled != nil {
		return vhs.GetAPIMethodsCalled()
	}
	return nil
}

// GetProtocolBuiltinFunctions mocked method
func (vhs *VMHostStub) GetProtocolBuiltinFunctions() vmcommon.FunctionNames {
	if vhs.GetProtocolBuiltinFunctionsCalled != nil {
		return vhs.GetProtocolBuiltinFunctionsCalled()
	}
	return make(vmcommon.FunctionNames)
}

// SetProtocolBuiltinFunctions mocked method
func (vhs *VMHostStub) SetProtocolBuiltinFunctions(functionNames vmcommon.FunctionNames) {
	if vhs.SetProtocolBuiltinFunctionsCalled != nil {
		vhs.SetProtocolBuiltinFunctionsCalled(functionNames)
	}
}

// IsBuiltinFunctionName mocked method
func (vhs *VMHostStub) IsBuiltinFunctionName(functionName string) bool {
	if vhs.IsBuiltinFunctionNameCalled != nil {
		return vhs.IsBuiltinFunctionNameCalled(functionName)
	}
	return false
}

// GetGasScheduleMap returns the currently stored gas schedule
func (vhs *VMHostStub) GetGasScheduleMap() config.GasScheduleMap {
	if vhs.GetGasScheduleMapCalled != nil {
		return vhs.GetGasScheduleMapCalled()
	}
	return nil
}

// RunSmartContractCall mocked method
func (vhs *VMHostStub) RunSmartContractCall(input *vmcommon.ContractCallInput) (vmOutput *vmcommon.VMOutput, err error) {
	if vhs.RunSmartContractCallCalled != nil {
		return vhs.RunSmartContractCallCalled(input)
	}
	return nil, nil
}

// RunSmartContractCreate mocked method
func (vhs *VMHostStub) RunSmartContractCreate(input *vmcommon.ContractCreateInput) (vmOutput *vmcommon.VMOutput, err error) {
	if vhs.RunSmartContractCreateCalled != nil {
		return vhs.RunSmartContractCreateCalled(input)
	}
	return nil, nil
}

// GasScheduleChange mocked method
func (vhs *VMHostStub) GasScheduleChange(newGasSchedule config.GasScheduleMap) {
	if vhs.GasScheduleChangeCalled != nil {
		vhs.GasScheduleChangeCalled(newGasSchedule)
	}
}

// IsInterfaceNil mocked method
func (vhs *VMHostStub) IsInterfaceNil() bool {
	if vhs.IsInterfaceNilCalled != nil {
		return vhs.IsInterfaceNilCalled()
	}
	return false
}

// GetContexts mocked method
func (vhs *VMHostStub) GetContexts() (
	arwen.ManagedTypesContext,
	arwen.BlockchainContext,
	arwen.MeteringContext,
	arwen.OutputContext,
	arwen.RuntimeContext,
	arwen.StorageContext,
) {
	if vhs.GetContextsCalled != nil {
		return vhs.GetContextsCalled()
	}
	return nil, nil, nil, nil, nil, nil
}

// SetRuntimeContext mocked method
func (vhs *VMHostStub) SetRuntimeContext(runtime arwen.RuntimeContext) {
	if vhs.SetRuntimeContextCalled != nil {
		vhs.SetRuntimeContextCalled(runtime)
	}
}
