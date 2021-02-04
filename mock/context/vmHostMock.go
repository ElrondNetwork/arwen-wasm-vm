package mock

import (
	"github.com/ElrondNetwork/arwen-wasm-vm/arwen"
	"github.com/ElrondNetwork/arwen-wasm-vm/crypto"
	"github.com/ElrondNetwork/arwen-wasm-vm/wasmer"
	"github.com/ElrondNetwork/elrond-go/core/parsers"
	"github.com/ElrondNetwork/elrond-go/core/vmcommon"
)

var _ arwen.VMHost = (*VMHostMock)(nil)

// VMHostMock is used in tests to check the VMHost interface method calls
type VMHostMock struct {
	BlockChainHook vmcommon.BlockchainHook
	CryptoHook     crypto.VMCrypto

	EthInput []byte

	BlockchainContext arwen.BlockchainContext
	RuntimeContext    arwen.RuntimeContext
	AsyncContext      arwen.AsyncContext
	OutputContext     arwen.OutputContext
	MeteringContext   arwen.MeteringContext
	StorageContext    arwen.StorageContext
	BigIntContext     arwen.BigIntContext

	SCAPIMethods  *wasmer.Imports
	IsBuiltinFunc bool

	StoredInputs []*vmcommon.ContractCallInput

	VMOutputQueue    []*vmcommon.VMOutput
	VMOutputToReturn int
	Err              error
}

// Crypto mocked method
func (host *VMHostMock) Crypto() crypto.VMCrypto {
	return host.CryptoHook
}

// Blockchain mocked method
func (host *VMHostMock) Blockchain() arwen.BlockchainContext {
	return host.BlockchainContext
}

// Runtime mocked method
func (host *VMHostMock) Runtime() arwen.RuntimeContext {
	return host.RuntimeContext
}

// Async mocked method
func (host *VMHostMock) Async() arwen.AsyncContext {
	return host.AsyncContext
}

// Output mocked method
func (host *VMHostMock) Output() arwen.OutputContext {
	return host.OutputContext
}

// Metering mocked method
func (host *VMHostMock) Metering() arwen.MeteringContext {
	return host.MeteringContext
}

// Storage mocked method
func (host *VMHostMock) Storage() arwen.StorageContext {
	return host.StorageContext
}

// BigInt mocked method
func (host *VMHostMock) BigInt() arwen.BigIntContext {
	return host.BigIntContext
}

// CallArgsParser mocked method
func (host *VMHostMock) CallArgsParser() arwen.CallArgsParser {
	return parsers.NewCallArgsParser()
}

// IsArwenV2Enabled mocked method
func (host *VMHostMock) IsArwenV2Enabled() bool {
	return true
}

// IsAheadOfTimeCompileEnabled mocked method
func (host *VMHostMock) IsAheadOfTimeCompileEnabled() bool {
	return true
}

// IsDynamicGasLockingEnabled mocked method
func (host *VMHostMock) IsDynamicGasLockingEnabled() bool {
	return true
}

// AreInSameShard mocked method
func (host *VMHostMock) AreInSameShard(left []byte, right []byte) bool {
	leftShard := host.BlockchainContext.GetShardOfAddress(left)
	rightShard := host.BlockchainContext.GetShardOfAddress(right)
	return leftShard == rightShard
}

// CreateNewContract mocked method
func (host *VMHostMock) CreateNewContract(_ *vmcommon.ContractCreateInput) ([]byte, error) {
	return nil, nil
}

// ExecuteOnSameContext mocked method
func (host *VMHostMock) ExecuteOnSameContext(_ *vmcommon.ContractCallInput) error {
	return nil
}

// ExecuteOnDestContext mocked method
func (host *VMHostMock) ExecuteOnDestContext(input *vmcommon.ContractCallInput) (*vmcommon.VMOutput, uint64, error) {
	if host.Err != nil {
		return nil, 0, host.Err
	}
	host.StoreInput(input)
	return host.GetNextVMOutput(), 0, nil
}

// InitState mocked method
func (host *VMHostMock) InitState() {
}

// PushState mocked method
func (host *VMHostMock) PushState() {
}

// PopState mocked method
func (host *VMHostMock) PopState() {
}

// ClearStateStack mocked method
func (host *VMHostMock) ClearStateStack() {
}

// GetAPIMethods mocked method
func (host *VMHostMock) GetAPIMethods() *wasmer.Imports {
	return host.SCAPIMethods
}

// GetProtocolBuiltinFunctions mocked method
func (host *VMHostMock) GetProtocolBuiltinFunctions() vmcommon.FunctionNames {
	return make(vmcommon.FunctionNames)
}

// IsBuiltinFunctionName mocked method
func (host *VMHostMock) IsBuiltinFunctionName(_ string) bool {
	return host.IsBuiltinFunc
}

func (host *VMHostMock) StoreInput(input *vmcommon.ContractCallInput) {
	if host.StoredInputs == nil {
		host.StoredInputs = make([]*vmcommon.ContractCallInput, 0)
	}
	host.StoredInputs = append(host.StoredInputs, input)
}

func (host *VMHostMock) EnqueueVMOutput(vmOutput *vmcommon.VMOutput) {
	if host.VMOutputQueue == nil {
		host.VMOutputQueue = make([]*vmcommon.VMOutput, 1)
		host.VMOutputQueue[0] = vmOutput
		host.VMOutputToReturn = 0
		return
	}

	host.VMOutputQueue = append(host.VMOutputQueue, vmOutput)
}

func (host *VMHostMock) GetNextVMOutput() *vmcommon.VMOutput {
	if host.VMOutputToReturn >= len(host.VMOutputQueue) {
		return nil
	}

	vmOutput := host.VMOutputQueue[host.VMOutputToReturn]
	host.VMOutputToReturn += 1
	return vmOutput
}
