package worldmock

import (
	"math/big"
)

// AccountMap is a map from address to account
type AccountMap map[string]*Account

// Account holds the account info
type Account struct {
	Exists          bool
	Address         []byte
	Nonce           uint64
	Balance         *big.Int
	BalanceDelta    *big.Int
	Storage         map[string][]byte
	Code            []byte
	CodeHash        []byte
	CodeMetadata    []byte
	AsyncCallData   string
	OwnerAddress    []byte
	Username        []byte
	ShardID         uint32
	IsSmartContract bool
	ESDTData        map[string]*ESDTData
}

// ESDTData models an account holding an ESDT token
type ESDTData struct {
	Balance      *big.Int
	BalanceDelta *big.Int
	Frozen       bool
}

var storageDefaultValue = []byte{}

// NewAccountMap creates a new AccountMap instance.
func NewAccountMap() AccountMap {
	return AccountMap(make(map[string]*Account))
}

// CreateAccount instantiates an empty account for the given address.
func (am AccountMap) CreateAccount(address []byte) *Account {
	newAccount := &Account{
		Nonce:           0,
		Balance:         big.NewInt(0),
		BalanceDelta:    big.NewInt(0),
		Storage:         make(map[string][]byte),
		IsSmartContract: false,
	}

	newAccount.Address = make([]byte, len(address))
	copy(newAccount.Address, address)
	am.PutAccount(newAccount)

	return newAccount
}

// PutAccount inserts account based on address.
func (am AccountMap) PutAccount(account *Account) {
	am[addressKey(account.Address)] = account
}

// PutAccounts inserts multiple accounts based on address.
func (am AccountMap) PutAccounts(accounts []*Account) {
	for _, account := range accounts {
		am.PutAccount(account)
	}
}

// GetAccount retrieves account based on address
func (am AccountMap) GetAccount(address []byte) *Account {
	return am[addressKey(address)]
}

// DeleteAccount removes account based on address
func (am AccountMap) DeleteAccount(address []byte) {
	delete(am, addressKey(address))
}

func addressKey(address []byte) string {
	return string(address)
}

// StorageValue yields the storage value for key, default 0
func (a *Account) StorageValue(key string) []byte {
	value, found := a.Storage[key]
	if !found {
		return storageDefaultValue
	}
	return value
}

// AddressBytes -
func (a *Account) AddressBytes() []byte {
	return a.Address
}

// GetNonce -
func (a *Account) GetNonce() uint64 {
	return a.Nonce
}

// GetCode -
func (a *Account) GetCode() []byte {
	return a.Code
}

// GetCodeMetadata -
func (a *Account) GetCodeMetadata() []byte {
	return a.CodeMetadata
}

// GetCodeHash -
func (a *Account) GetCodeHash() []byte {
	return a.CodeHash
}

// GetRootHash -
func (a *Account) GetRootHash() []byte {
	return []byte{}
}

// GetBalance -
func (a *Account) GetBalance() *big.Int {
	return a.Balance
}

// SetBalance -
func (a *Account) SetBalance(balance int64) {
	a.Balance = big.NewInt(balance)
}

// GetDeveloperReward -
func (a *Account) GetDeveloperReward() *big.Int {
	return big.NewInt(0)
}

// GetOwnerAddress -
func (a *Account) GetOwnerAddress() []byte {
	return a.OwnerAddress
}

// GetUserName -
func (a *Account) GetUserName() []byte {
	return a.Username
}

// IsInterfaceNil -
func (a *Account) IsInterfaceNil() bool {
	return a == nil
}
