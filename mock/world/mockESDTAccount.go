package worldmock

import (
	"bytes"
	"math/big"

	"github.com/ElrondNetwork/elrond-go/core"
	"github.com/ElrondNetwork/elrond-go/data"
	"github.com/ElrondNetwork/elrond-go/data/esdt"
)

// MakeTokenKey creates the storage key corresponding to the given tokenName.
func MakeTokenKey(tokenName []byte, nonce uint64) []byte {
	nonceBytes := big.NewInt(0).SetUint64(nonce).Bytes()
	tokenKey := append(ESDTTokenKeyPrefix, tokenName...)
	tokenKey = append(tokenKey, nonceBytes...)
	return tokenKey
}

// MakeTokenRolesKey creates the storage key corresponding to the roles for the
// given tokenName.
func MakeTokenRolesKey(tokenName []byte) []byte {
	tokenRolesKey := append(ESDTRoleKeyPrefix, tokenName...)
	return tokenRolesKey
}

// MakeLastNonceKey creates the storage key corresponding to the last nonce of
// the given tokenName.
func MakeLastNonceKey(tokenName []byte) []byte {
	tokenNonceKey := append(ESDTNonceKeyPrefix, tokenName...)
	return tokenNonceKey
}

// IsESDTKey returns true if the given storage key is ESDT-related
func IsESDTKey(key []byte) bool {
	return IsTokenKey(key) || IsRoleKey(key) || IsNonceKey(key)
}

// IsTokenKey returns true if the given storage key belongs to an ESDT token.
func IsTokenKey(key []byte) bool {
	return bytes.HasPrefix(key, ESDTTokenKeyPrefix)
}

// IsRoleKey returns true if the given storage key belongs to an ESDT role.
func IsRoleKey(key []byte) bool {
	return bytes.HasPrefix(key, ESDTRoleKeyPrefix)
}

// IsNonceKey returns true if the given storage key belongs to an ESDT nonce.
func IsNonceKey(key []byte) bool {
	return bytes.HasPrefix(key, ESDTNonceKeyPrefix)
}

// GetTokenNameFromKey extracts the token name from the given storage key; it
// does not check whether the key is indeed a token key or not.
func GetTokenNameFromKey(key []byte) []byte {
	return key[len(ESDTTokenKeyPrefix):]
}

// GetTokenBalanceByName returns the ESDT balance of the account, specified by
// the token name.
func (a *Account) GetTokenBalanceByName(tokenName string) (*big.Int, error) {
	tokenKey := MakeTokenKey([]byte(tokenName), 0)
	return a.GetTokenBalance(tokenKey)
}

// GetTokenBalance returns the ESDT balance of the account, specified by the
// token key.
func (a *Account) GetTokenBalance(tokenKey []byte) (*big.Int, error) {
	tokenData, err := a.GetTokenData(tokenKey)
	if err != nil {
		return nil, err
	}

	return tokenData.Value, nil
}

// SetTokenBalance sets the ESDT balance of the account, specified by the token
// key.
func (a *Account) SetTokenBalance(tokenKey []byte, balance *big.Int) error {
	tokenData, err := a.GetTokenData(tokenKey)
	if err != nil {
		return err
	}

	if balance.Sign() < 0 {
		return data.ErrNegativeValue
	}

	tokenData.Value = balance
	return a.SetTokenData(tokenKey, tokenData)
}

// GetTokenData gets the ESDT information related to a token from the storage of the account.
func (a *Account) GetTokenData(tokenKey []byte) (*esdt.ESDigitalToken, error) {
	esdtData := &esdt.ESDigitalToken{
		Value: big.NewInt(0),
		Type:  uint32(core.Fungible),
		TokenMetaData: &esdt.MetaData{
			Name:  GetTokenNameFromKey(tokenKey),
			Nonce: 0,
		},
	}

	marshaledData, err := a.DataTrieTracker().RetrieveValue(tokenKey)
	if err != nil || len(marshaledData) == 0 {
		return esdtData, nil
	}

	err = WorldMarshalizer.Unmarshal(esdtData, marshaledData)
	if err != nil {
		return nil, err
	}

	return esdtData, nil
}

// SetTokenData sets the ESDT information related to a token into the storage of the account.
func (a *Account) SetTokenData(tokenKey []byte, tokenData *esdt.ESDigitalToken) error {
	marshaledData, err := WorldMarshalizer.Marshal(tokenData)
	if err != nil {
		return err
	}

	return a.DataTrieTracker().SaveKeyValue(tokenKey, marshaledData)
}

// SetTokenRoles sets the specified roles to the account, corresponding to the given tokenName.
func (a *Account) SetTokenRoles(tokenName []byte, roles [][]byte) error {
	tokenRolesKey := MakeTokenRolesKey(tokenName)
	tokenRolesData := &esdt.ESDTRoles{
		Roles: roles,
	}

	marshaledData, err := WorldMarshalizer.Marshal(tokenRolesData)
	if err != nil {
		return err
	}

	return a.DataTrieTracker().SaveKeyValue(tokenRolesKey, marshaledData)
}

// SetTokenRolesAsStrings sets the specified roles to the account, corresponding to the given tokenName.
func (a *Account) SetTokenRolesAsStrings(tokenName []byte, rolesAsStrings []string) error {
	roles := make([][]byte, len(rolesAsStrings))
	for i := 0; i < len(roles); i++ {
		roles[i] = []byte(rolesAsStrings[i])
	}

	return a.SetTokenRoles(tokenName, roles)
}

// SetLastNonces writes the last nonces of each specified ESDT into the storage.
func (a *Account) SetLastNonces(lastNonces map[string]uint64) error {
	for tokenName, nonce := range lastNonces {
		tokenNonceKey := MakeLastNonceKey([]byte(tokenName))
		nonceBytes := big.NewInt(0).SetUint64(nonce).Bytes()

		err := a.DataTrieTracker().SaveKeyValue(tokenNonceKey, nonceBytes)
		if err != nil {
			return err
		}
	}

	return nil
}

// GetTokenRoles returns the roles of the account for the specified tokenName.
func (a *Account) GetTokenRoles(tokenName []byte) ([][]byte, error) {
	tokenRolesKey := MakeTokenRolesKey(tokenName)
	tokenRolesData := &esdt.ESDTRoles{
		Roles: make([][]byte, 0),
	}

	marshaledData, err := a.DataTrieTracker().RetrieveValue(tokenRolesKey)
	if err != nil || len(marshaledData) == 0 {
		return tokenRolesData.Roles, nil
	}

	err = WorldMarshalizer.Unmarshal(tokenRolesData, marshaledData)
	if err != nil {
		return nil, err
	}

	return tokenRolesData.Roles, nil

}

// GetAllTokenData returns the information about all the ESDT tokens held by the account.
func (a *Account) GetAllTokenData() (map[string]*esdt.ESDigitalToken, error) {
	tokenDataMap := make(map[string]*esdt.ESDigitalToken)
	tokenKeys := a.GetTokenKeys()
	for _, tokenKey := range tokenKeys {
		tokenData, err := a.GetTokenData(tokenKey)
		if err != nil {
			return nil, err
		}

		if tokenData.TokenMetaData == nil {
			tokenData.TokenMetaData = &esdt.MetaData{
				Nonce: 0,
			}
		}

		tokenKeyStr := string(tokenKey)
		tokenDataMap[tokenKeyStr] = tokenData
	}

	return tokenDataMap, nil
}

// GetTokenKeys returns the storage keys of all the ESDT tokens owned by the account.
func (a *Account) GetTokenKeys() [][]byte {
	tokenKeys := make([][]byte, 0)
	for key := range a.Storage {
		if IsTokenKey([]byte(key)) {
			tokenKeys = append(tokenKeys, []byte(key))
		}
	}

	return tokenKeys
}

// GetTokenNames returns the names of all the tokens owned by the account.
func (a *Account) GetTokenNames() ([][]byte, error) {
	tokenKeys := a.GetTokenKeys()
	tokenNames := make([][]byte, len(tokenKeys))

	for i := 0; i < len(tokenKeys); i++ {
		tokenKey := tokenKeys[i]
		tokenData, err := a.GetTokenData(tokenKey)
		if err != nil {
			return nil, err
		}

		tokenNames = append(tokenNames, tokenData.TokenMetaData.Name)
	}

	return tokenNames, nil
}
