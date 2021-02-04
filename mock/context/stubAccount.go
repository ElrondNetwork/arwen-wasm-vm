package mock

import (
	"math/big"
)

// StubAccount is used with the blockchain stub in arwen context tests.
// It has minimal functionality.
type StubAccount struct {
	Address      []byte
	Nonce        uint64
	Balance      *big.Int
	CodeHash     []byte
	CodeMetadata []byte
	OwnerAddress []byte
	UserName     []byte

	ShardID uint32
	Storage map[string][]byte
	Err     error
}

// AddressBytes -
func (a *StubAccount) AddressBytes() []byte {
	return a.Address
}

// GetNonce -
func (a *StubAccount) GetNonce() uint64 {
	return a.Nonce
}

// GetCodeMetadata -
func (a *StubAccount) GetCodeMetadata() []byte {
	return a.CodeMetadata
}

// GetCodeHash -
func (a *StubAccount) GetCodeHash() []byte {
	return a.CodeHash
}

// GetRootHash -
func (a *StubAccount) GetRootHash() []byte {
	return []byte{}
}

// GetBalance -
func (a *StubAccount) GetBalance() *big.Int {
	if a.Balance == nil {
		return big.NewInt(0)
	}
	return a.Balance
}

// GetDeveloperReward -
func (a *StubAccount) GetDeveloperReward() *big.Int {
	return big.NewInt(0)
}

// GetOwnerAddress -
func (a *StubAccount) GetOwnerAddress() []byte {
	return a.OwnerAddress
}

// GetUserName -
func (a *StubAccount) GetUserName() []byte {
	return a.UserName
}

// IsInterfaceNil -
func (a *StubAccount) IsInterfaceNil() bool {
	return a == nil
}
