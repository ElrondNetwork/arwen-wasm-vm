package arwendebug

import (
	"math/big"
)

// AccountsMap is a map from address to account
type AccountsMap map[string]*Account

// Account is a debug account
type Account struct {
	AddressHex      string
	Nonce           uint64
	Balance         *big.Int
	BalanceString   string
	CodeHex         string
	CodeMetadataHex string
	OwnerAddressHex string
	UserNameHex     string
	Storage         map[string]string
	ShardID         uint32
}

// NewAccount creates a new debug account
func NewAccount(address []byte, nonce uint64, balance *big.Int) *Account {
	if balance == nil {
		balance = big.NewInt(0)
	}

	return &Account{
		AddressHex:    toHex(address),
		Nonce:         nonce,
		Balance:       balance,
		BalanceString: balance.String(),
		CodeHex:       "",
		Storage:       make(map[string]string),
	}
}

// AddressBytes -
func (a *Account) AddressBytes() []byte {
	return fromHexNoError(a.AddressHex)
}

// GetNonce -
func (a *Account) GetNonce() uint64 {
	return a.Nonce
}

// GetCode -
func (a *Account) GetCode() []byte {
	return fromHexNoError(a.CodeHex)
}

// GetCodeMetadata -
func (a *Account) GetCodeMetadata() []byte {
	return fromHexNoError(a.CodeMetadataHex)
}

// GetCodeHash -
func (a *Account) GetCodeHash() []byte {
	return []byte{}
}

// GetRootHash -
func (a *Account) GetRootHash() []byte {
	return []byte{}
}

// GetBalance -
func (a *Account) GetBalance() *big.Int {
	if a.Balance == nil {
		return big.NewInt(0)
	}
	return a.Balance
}

// GetDeveloperReward -
func (a *Account) GetDeveloperReward() *big.Int {
	return big.NewInt(0)
}

// GetOwnerAddress -
func (a *Account) GetOwnerAddress() []byte {
	return fromHexNoError(a.OwnerAddressHex)
}

// GetUserName -
func (a *Account) GetUserName() []byte {
	return fromHexNoError(a.UserNameHex)
}

// IsInterfaceNil -
func (a *Account) IsInterfaceNil() bool {
	return a == nil
}
