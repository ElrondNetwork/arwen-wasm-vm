package worldmock

import (
	"bytes"
	"encoding/hex"
	"fmt"
)

// SCAddressNumLeadingZeros is the number of zero bytes every smart contract address begins with.
const SCAddressNumLeadingZeros = 8

// IsSmartContractAddress verifies the address format.
// In Elrond, smart contract addresses start with 8 bytes of 0.
func IsSmartContractAddress(address []byte) bool {
	leadingZeros := make([]byte, SCAddressNumLeadingZeros)
	return bytes.Equal(address[:SCAddressNumLeadingZeros], leadingZeros)
}

func (a *Account) Validate() error {
	if len(a.Address) != 32 {
		return fmt.Errorf(
			"account address should be 32 bytes long: 0x%s",
			hex.EncodeToString(a.Address))
	}

	scAddress := IsSmartContractAddress(a.Address)
	if len(a.Code) > 0 {
		if !scAddress {
			return fmt.Errorf(
				"account has a smart contract address, but has no code: 0x%s",
				hex.EncodeToString(a.Address))
		}
	} else {
		if scAddress {
			return fmt.Errorf(
				"account has code but not a smart contract address: %s",
				hex.EncodeToString(a.Address))
		}
	}

	return nil
}
