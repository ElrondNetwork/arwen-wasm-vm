package mandosexpressionreconstructor

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"math/big"
	"strconv"
	"strings"

	ei "github.com/ElrondNetwork/arwen-wasm-vm/mandos-go/expression/interpreter"
)

type ExprReconstructorHint uint64

const (
	// NoHint indicates that the type if not known
	NoHint ExprReconstructorHint = iota

	// NumberHint hints that value should be a number
	NumberHint

	// AddressHint hints that value should be an address
	AddressHint

	// StrHint hints that value should be a string expression, e.g. a username, "str:..."
	StrHint

	// CodeHint hints that value should be a smart contract code, normally loaded from a file
	CodeHint
)

const maxBytesInterpretedAsNumber = 15

// ExprReconstructor is a component that attempts to convert raw bytes to a human-readable format.
type ExprReconstructor struct{}

func (er *ExprReconstructor) Reconstruct(value []byte, hint ExprReconstructorHint) string {
	switch hint {
	case NumberHint:
		return fmt.Sprintf("%d", big.NewInt(0).SetBytes(value))
	case StrHint:
		return fmt.Sprintf("str:%s", string(value))
	case AddressHint:
		return addressPretty(value)
	case CodeHint:
		return codePretty(value)
	default:
		return unknownByteArrayPretty(value)
	}
}

func (er *ExprReconstructor) ReconstructFromBigInt(value *big.Int) string {
	return er.Reconstruct(value.Bytes(), NumberHint)
}

func (er *ExprReconstructor) ReconstructFromUint64(value uint64) string {
	return er.Reconstruct(big.NewInt(0).SetUint64(value).Bytes(), NumberHint)
}

func unknownByteArrayPretty(bytes []byte) string {
	if len(bytes) == 0 {
		return ""
	}

	// fully interpret as string
	if canInterpretAsString(bytes) {
		return fmt.Sprintf("0x%s (str:%s)", hex.EncodeToString(bytes), string(bytes))
	}

	// interpret as number
	if len(bytes) < maxBytesInterpretedAsNumber {
		asInt := big.NewInt(0).SetBytes(bytes)
		return fmt.Sprintf("0x%s (%d)", hex.EncodeToString(bytes), asInt)
	}

	// default interpret as string with escaped bytes
	return fmt.Sprintf("0x%s (str:%s)", hex.EncodeToString(bytes), strconv.Quote(string(bytes)))
}

func addressPretty(value []byte) string {
	if len(value) != 32 {
		return unknownByteArrayPretty(value)
	}

	// smart contract addresses
	leadingZeros := make([]byte, ei.SCAddressNumLeadingZeros)
	if bytes.Equal(value[:ei.SCAddressNumLeadingZeros], leadingZeros) {
		if value[31] == byte('_') {
			addrStr := string(value[ei.SCAddressNumLeadingZeros:])
			addrStr = strings.TrimRight(addrStr, "_")
			return fmt.Sprintf("sc:%s", addrStr)
		} else {
			// last byte is the shard id and is explicit
			addrStr := string(value[ei.SCAddressNumLeadingZeros:31])
			addrStr = strings.TrimRight(addrStr, "_")
			shard_id := value[31]
			return fmt.Sprintf("sc:%s#%x", addrStr, shard_id)
		}
	}

	// regular addresses
	if value[31] == byte('_') {
		addrStr := string(value)
		addrStr = strings.TrimRight(addrStr, "_")
		return fmt.Sprintf("address:%s", addrStr)
	} else {
		// last byte is the shard id and is explicit
		addrStr := string(value[:31])
		addrStr = strings.TrimRight(addrStr, "_")
		shard_id := value[31]
		return fmt.Sprintf("address:%s#%02x", addrStr, shard_id)
	}
}

func canInterpretAsString(bytes []byte) bool {
	if len(bytes) == 0 {
		return false
	}
	for _, b := range bytes {
		if b < 32 || b > 126 {
			return false
		}
	}
	return true
}

func codePretty(bytes []byte) string {
	if len(bytes) == 0 {
		return ""
	}
	encoded := hex.EncodeToString(bytes)

	if len(encoded) > 20 {
		return fmt.Sprintf("0x%s...", encoded[:20])
	}

	return fmt.Sprintf("0x%s", encoded)
}
