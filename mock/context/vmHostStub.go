package mock

import (
	"math/big"

	"github.com/ElrondNetwork/arwen-wasm-vm/arwen"
	"github.com/ElrondNetwork/arwen-wasm-vm/crypto"
	"github.com/ElrondNetwork/arwen-wasm-vm/wasmer"
	"github.com/ElrondNetwork/elrond-go/core/vmcommon"
)

var _ arwen.VMHost = (*VMHostStub)(nil)

// VMHostStub is used in tests to check the VMHost interface method calls
type VMHostStub struct {
	InitStateCalled       func()
	PushStateCalled       func()
	PopStateCalled        func()
	ClearStateStackCalled func()

	CryptoCalled                      func() crypto.VMCrypto
	BlockchainCalled                  func() arwen.BlockchainContext
	RuntimeCalled                     func() arwen.RuntimeContext
	BigIntCalled                      func() arwen.BigIntContext
	BytesHeapCalled                   func() arwen.BytesHeapContext
	OutputCalled                      func() arwen.OutputContext
	MeteringCalled                    func() arwen.MeteringContext
	StorageCalled                     func() arwen.StorageContext
	RevertESDTTransferCalled          func(input *vmcommon.ContractCallInput)
	ExecuteESDTTransferCalled         func(destination []byte, sender []byte, tokenIdentifier []byte, nonce uint64, value *big.Int) (*vmcommon.VMOutput, uint64, error)
	CreateNewContractCalled           func(input *vmcommon.ContractCreateInput) ([]byte, error)
	ExecuteOnSameContextCalled        func(input *vmcommon.ContractCallInput) (*arwen.AsyncContextInfo, error)
	ExecuteOnDestContextCalled        func(input *vmcommon.ContractCallInput) (*vmcommon.VMOutput, *arwen.AsyncContextInfo, uint64, error)
	GetAPIMethodsCalled               func() *wasmer.Imports
	GetProtocolBuiltinFunctionsCalled func() vmcommon.FunctionNames
	IsBuiltinFunctionNameCalled       func(functionName string) bool
	AreInSameShardCalled              func(left []byte, right []byte) bool
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
func (vhs *VMHostStub) BigInt() arwen.BigIntContext {
	if vhs.BigIntCalled != nil {
		return vhs.BigIntCalled()
	}
	return nil
}

// BytesHeap mocked method
func (vhs *VMHostStub) BytesHeap() arwen.BytesHeapContext {
	if vhs.BytesHeapCalled != nil {
		return vhs.BytesHeapCalled()
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

// RevertESDTTransfer mocked method
func (vhs *VMHostStub) RevertESDTTransfer(input *vmcommon.ContractCallInput) {
	if vhs.RevertESDTTransferCalled != nil {
		vhs.RevertESDTTransferCalled(input)
	}
}

// ExecuteESDTTransfer mocked method
func (vhs *VMHostStub) ExecuteESDTTransfer(destination []byte, sender []byte, tokenIdentifier []byte, nonce uint64, value *big.Int) (*vmcommon.VMOutput, uint64, error) {
	if vhs.ExecuteESDTTransferCalled != nil {
		return vhs.ExecuteESDTTransferCalled(destination, sender, tokenIdentifier, nonce, value)
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
func (vhs *VMHostStub) ExecuteOnDestContext(input *vmcommon.ContractCallInput) (*vmcommon.VMOutput, *arwen.AsyncContextInfo, uint64, error) {
	if vhs.ExecuteOnDestContextCalled != nil {
		return vhs.ExecuteOnDestContextCalled(input)
	}
	return nil, nil, 0, nil
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

// IsBuiltinFunctionName mocked method
func (vhs *VMHostStub) IsBuiltinFunctionName(functionName string) bool {
	if vhs.IsBuiltinFunctionNameCalled != nil {
		return vhs.IsBuiltinFunctionNameCalled(functionName)
	}
	return false
}
