package mandosjsonmodel

import (
	"math/big"

	oj "github.com/ElrondNetwork/arwen-wasm-vm/mandos-go/orderedjson"
)

// JSONBytesFromString stores a byte slice
type JSONBytesFromString struct {
	Value    []byte
	Original string
}

// NewJSONBytesFromString creates a new JSONBytesFromString instance.
func NewJSONBytesFromString(value []byte, originalStr string) JSONBytesFromString {
	return JSONBytesFromString{
		Value:    value,
		Original: originalStr,
	}
}

// JSONBytesFromTree stores a parsed byte slice, either from a string, or from a list of strings.
// The list of strings representation can be used in storage, arguments or results,
// and it is designed to make it easier to express serialized objects.
// The strings in the list get simply concatenated to produce a value.
type JSONBytesFromTree struct {
	Value    []byte
	Original oj.OJsonObject
}

// OriginalEmpty returns true if the object originates from "".
func (jb JSONBytesFromTree) OriginalEmpty() bool {
	if str, isStr := jb.Original.(*oj.OJsonString); isStr {
		return len(str.Value) == 0
	}
	return false
}

// JSONBigInt stores the parsed big int value but also the original parsed string
type JSONBigInt struct {
	Value    *big.Int
	Original string
}

// JSONBigIntZero provides an unitialized zero value.
func JSONBigIntZero() JSONBigInt {
	return JSONBigInt{
		Value:    big.NewInt(0),
		Original: "",
	}
}

// JSONUint64 stores the parsed uint64 value but also the original parsed string
type JSONUint64 struct {
	Value    uint64
	Original string
}

// JSONUint64Zero provides an unitialized zero value.
func JSONUint64Zero() JSONUint64 {
	return JSONUint64{
		Value:    0,
		Original: "",
	}
}
