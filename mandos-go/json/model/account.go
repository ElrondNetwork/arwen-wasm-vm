package mandosjsonmodel

import "bytes"

// Account is a json object representing an account.
type Account struct {
	Address         JSONBytesFromString
	Shard           JSONUint64
	IsSmartContract bool
	Comment         string
	Nonce           JSONUint64
	Balance         JSONBigInt
	Username        JSONBytesFromString
	Storage         []*StorageKeyValuePair
	Code            JSONBytesFromString
	Owner           JSONBytesFromString
	AsyncCallData   string
	ESDTData        []*ESDTData
}

// StorageKeyValuePair is a json key value pair in the storage map.
type StorageKeyValuePair struct {
	Key   JSONBytesFromString
	Value JSONBytesFromTree
}

// CheckAccount is a json object representing checks for an account.
type CheckAccount struct {
	Address               JSONBytesFromString
	Comment               string
	Nonce                 JSONCheckUint64
	Balance               JSONCheckBigInt
	Username              JSONCheckBytes
	CheckStorage          []*CheckStorageKeyValuePair
	IgnoreStorage         bool
	MoreStorageAllowed    bool
	Code                  JSONCheckBytes
	Owner                 JSONCheckBytes
	AsyncCallData         JSONCheckBytes
	CheckESDTData         []*CheckESDTData
	IgnoreESDT            bool
	MoreESDTTokensAllowed bool
}

// CheckStorageKeyValuePair checks a single entry in storage.
type CheckStorageKeyValuePair struct {
	Key        JSONBytesFromString
	CheckValue JSONCheckBytes
}

// CheckAccounts encodes rules to check mock accounts.
type CheckAccounts struct {
	Accounts            []*CheckAccount
	MoreAccountsAllowed bool
}

// FindAccount searches an account list by address.
func FindAccount(accounts []*Account, address []byte) *Account {
	for _, acct := range accounts {
		if bytes.Equal(acct.Address.Value, address) {
			return acct
		}
	}
	return nil
}

// FindCheckAccount searches a check account list by address.
func FindCheckAccount(accounts []*CheckAccount, address []byte) *CheckAccount {
	for _, acct := range accounts {
		if bytes.Equal(acct.Address.Value, address) {
			return acct
		}
	}
	return nil
}
