package mock

import (
	"github.com/ElrondNetwork/elrond-go/core/vmcommon"
)

var _ vmcommon.BlockchainHook = (*BlockchainHookStub)(nil)

// BlockchainHookStub is used in tests to check that interface methods were called
type BlockchainHookStub struct {
	NewAddressCalled              func(creatorAddress []byte, creatorNonce uint64, vmType []byte) ([]byte, error)
	GetStorageDataCalled          func(accountsAddress []byte, index []byte) ([]byte, error)
	GetBlockHashCalled            func(nonce uint64) ([]byte, error)
	LastNonceCalled               func() uint64
	LastRoundCalled               func() uint64
	LastTimeStampCalled           func() uint64
	LastRandomSeedCalled          func() []byte
	LastEpochCalled               func() uint32
	GetStateRootHashCalled        func() []byte
	CurrentNonceCalled            func() uint64
	CurrentRoundCalled            func() uint64
	CurrentTimeStampCalled        func() uint64
	CurrentRandomSeedCalled       func() []byte
	CurrentEpochCalled            func() uint32
	ProcessBuiltInFunctionCalled  func(input *vmcommon.ContractCallInput) (*vmcommon.VMOutput, error)
	GetBuiltinFunctionNamesCalled func() vmcommon.FunctionNames
	GetAllStateCalled             func(address []byte) (map[string][]byte, error)
	GetUserAccountCalled          func(address []byte) (vmcommon.UserAccountHandler, error)
	GetShardOfAddressCalled       func(address []byte) uint32
	IsSmartContractCalled         func(address []byte) bool
	IsPayableCalled               func(address []byte) (bool, error)
	GetCompiledCodeCalled         func(codeHash []byte) (bool, []byte)
	SaveCompiledCodeCalled        func(codeHash []byte, code []byte)
	GetCodeCalled                 func(account vmcommon.UserAccountHandler) []byte
}

// NewAddress mocked method
func (b *BlockchainHookStub) NewAddress(creatorAddress []byte, creatorNonce uint64, vmType []byte) ([]byte, error) {
	if b.NewAddressCalled != nil {
		return b.NewAddressCalled(creatorAddress, creatorNonce, vmType)
	}
	return []byte("newAddress"), nil
}

// GetStorageData mocked method
func (b *BlockchainHookStub) GetStorageData(accountAddress []byte, index []byte) ([]byte, error) {
	if b.GetStorageDataCalled != nil {
		return b.GetStorageDataCalled(accountAddress, index)
	}
	return nil, nil
}

// GetBlockhash mocked method
func (b *BlockchainHookStub) GetBlockhash(nonce uint64) ([]byte, error) {
	if b.GetBlockHashCalled != nil {
		return b.GetBlockHashCalled(nonce)
	}
	return []byte("roothash"), nil
}

// LastNonce mocked method
func (b *BlockchainHookStub) LastNonce() uint64 {
	if b.LastNonceCalled != nil {
		return b.LastNonceCalled()
	}
	return 0
}

// LastRound mocked method
func (b *BlockchainHookStub) LastRound() uint64 {
	if b.LastRoundCalled != nil {
		return b.LastRoundCalled()
	}
	return 0
}

// LastTimeStamp mocked method
func (b *BlockchainHookStub) LastTimeStamp() uint64 {
	if b.LastTimeStampCalled != nil {
		return b.LastTimeStampCalled()
	}
	return 0
}

// LastRandomSeed mocked method
func (b *BlockchainHookStub) LastRandomSeed() []byte {
	if b.LastRandomSeedCalled != nil {
		return b.LastRandomSeedCalled()
	}
	return []byte("seed")
}

// LastEpoch mocked method
func (b *BlockchainHookStub) LastEpoch() uint32 {
	if b.LastEpochCalled != nil {
		return b.LastEpochCalled()
	}
	return 0
}

// GetStateRootHash mocked method
func (b *BlockchainHookStub) GetStateRootHash() []byte {
	if b.GetStateRootHashCalled != nil {
		return b.GetStateRootHashCalled()
	}
	return []byte("roothash")
}

// CurrentNonce mocked method
func (b *BlockchainHookStub) CurrentNonce() uint64 {
	if b.CurrentNonceCalled != nil {
		return b.CurrentNonceCalled()
	}
	return 0
}

// CurrentRound mocked method
func (b *BlockchainHookStub) CurrentRound() uint64 {
	if b.CurrentRoundCalled != nil {
		return b.CurrentRoundCalled()
	}
	return 0
}

// CurrentTimeStamp mocked method
func (b *BlockchainHookStub) CurrentTimeStamp() uint64 {
	if b.CurrentTimeStampCalled != nil {
		return b.CurrentTimeStampCalled()
	}
	return 0
}

// CurrentRandomSeed mocked method
func (b *BlockchainHookStub) CurrentRandomSeed() []byte {
	if b.CurrentRandomSeedCalled != nil {
		return b.CurrentRandomSeedCalled()
	}
	return []byte("seed")
}

// CurrentEpoch mocked method
func (b *BlockchainHookStub) CurrentEpoch() uint32 {
	if b.CurrentEpochCalled != nil {
		return b.CurrentEpochCalled()
	}
	return 0
}

// ProcessBuiltInFunction mocked method
func (b *BlockchainHookStub) ProcessBuiltInFunction(input *vmcommon.ContractCallInput) (*vmcommon.VMOutput, error) {
	if b.ProcessBuiltInFunctionCalled != nil {
		return b.ProcessBuiltInFunctionCalled(input)
	}
	return &vmcommon.VMOutput{}, nil
}

// GetBuiltinFunctionNames mocked method
func (b *BlockchainHookStub) GetBuiltinFunctionNames() vmcommon.FunctionNames {
	if b.GetBuiltinFunctionNamesCalled != nil {
		return b.GetBuiltinFunctionNamesCalled()
	}
	return make(vmcommon.FunctionNames)
}

// GetAllState mocked method
func (b *BlockchainHookStub) GetAllState(address []byte) (map[string][]byte, error) {
	if b.GetAllStateCalled != nil {
		return b.GetAllStateCalled(address)
	}
	return nil, nil
}

// GetUserAccount mocked method
func (b *BlockchainHookStub) GetUserAccount(address []byte) (vmcommon.UserAccountHandler, error) {
	if b.GetUserAccountCalled != nil {
		return b.GetUserAccountCalled(address)
	}
	return nil, nil
}

// GetCode mocked method
func (b *BlockchainHookStub) GetCode(account vmcommon.UserAccountHandler) []byte {
	if b.GetCodeCalled != nil {
		return b.GetCodeCalled(account)
	}
	return nil
}

// GetShardOfAddress mocked method
func (b *BlockchainHookStub) GetShardOfAddress(address []byte) uint32 {
	if b.GetShardOfAddressCalled != nil {
		return b.GetShardOfAddressCalled(address)
	}
	return 0
}

// IsSmartContract mocked method
func (b *BlockchainHookStub) IsSmartContract(address []byte) bool {
	if b.IsSmartContractCalled != nil {
		return b.IsSmartContractCalled(address)
	}
	return false
}

// IsPayable mocked method
func (b *BlockchainHookStub) IsPayable(address []byte) (bool, error) {
	if b.IsPayableCalled != nil {
		return b.IsPayableCalled(address)
	}
	return true, nil
}

// SaveCompiledCode mocked method
func (b *BlockchainHookStub) SaveCompiledCode(codeHash []byte, code []byte) {
	if b.SaveCompiledCodeCalled != nil {
		b.SaveCompiledCodeCalled(codeHash, code)
	}
}

// GetCompiledCode mocked method
func (b *BlockchainHookStub) GetCompiledCode(codeHash []byte) (bool, []byte) {
	if b.GetCompiledCodeCalled != nil {
		return b.GetCompiledCodeCalled(codeHash)
	}
	return false, nil
}

// ClearCompiledCodes mocked method
func (b *BlockchainHookStub) ClearCompiledCodes() {
}

// IsInterfaceNil mocked method
func (b *BlockchainHookStub) IsInterfaceNil() bool {
	return b == nil
}
