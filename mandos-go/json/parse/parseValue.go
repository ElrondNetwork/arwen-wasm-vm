package mandosjsonparse

import (
	"errors"
	"math/big"

	mj "github.com/ElrondNetwork/arwen-wasm-vm/mandos-go/json/model"
	oj "github.com/ElrondNetwork/arwen-wasm-vm/mandos-go/orderedjson"
	twos "github.com/ElrondNetwork/big-int-util/twos-complement"
)

type bigIntParseFormat int

const (
	bigIntSignedBytes bigIntParseFormat = iota
	bigIntUnsignedBytes
)

func (p *Parser) processCheckBigInt(obj oj.OJsonObject, format bigIntParseFormat) (mj.JSONCheckBigInt, error) {
	if IsStar(obj) {
		// "*" means any value, skip checking it
		return mj.JSONCheckBigInt{
			Value:    nil,
			IsStar:   true,
			Original: "*"}, nil
	}

	jbi, err := p.processBigInt(obj, format)
	if err != nil {
		return mj.JSONCheckBigInt{}, err
	}
	return mj.JSONCheckBigInt{
		Value:    jbi.Value,
		IsStar:   false,
		Original: jbi.Original,
	}, nil
}

func (p *Parser) processBigInt(obj oj.OJsonObject, format bigIntParseFormat) (mj.JSONBigInt, error) {
	strVal, err := p.parseString(obj)
	if err != nil {
		return mj.JSONBigInt{}, err
	}

	bi, err := p.parseBigInt(strVal, format)
	return mj.JSONBigInt{
		Value:    bi,
		Original: strVal,
	}, err
}

func (p *Parser) parseBigInt(strRaw string, format bigIntParseFormat) (*big.Int, error) {
	bytes, err := p.ExprInterpreter.InterpretString(strRaw)
	if err != nil {
		return nil, err
	}
	switch format {
	case bigIntSignedBytes:
		return twos.FromBytes(bytes), nil
	case bigIntUnsignedBytes:
		return big.NewInt(0).SetBytes(bytes), nil
	default:
		return nil, errors.New("unknown format requested")
	}
}

func (p *Parser) processCheckUint64(obj oj.OJsonObject) (mj.JSONCheckUint64, error) {
	if IsStar(obj) {
		// "*" means any value, skip checking it
		return mj.JSONCheckUint64{
			Value:    0,
			IsStar:   true,
			Original: "*"}, nil
	}

	ju, err := p.processUint64(obj)
	if err != nil {
		return mj.JSONCheckUint64{}, err
	}
	return mj.JSONCheckUint64{
		Value:    ju.Value,
		IsStar:   false,
		Original: ju.Original}, nil

}

func (p *Parser) processUint64(obj oj.OJsonObject) (mj.JSONUint64, error) {
	bi, err := p.processBigInt(obj, bigIntUnsignedBytes)
	if err != nil {
		return mj.JSONUint64{}, err
	}

	if bi.Value == nil || !bi.Value.IsUint64() {
		return mj.JSONUint64{}, errors.New("value is not uint64")
	}

	return mj.JSONUint64{
		Value:    bi.Value.Uint64(),
		Original: bi.Original}, nil
}

func (p *Parser) parseCheckBytes(obj oj.OJsonObject) (mj.JSONCheckBytes, error) {
	if IsStar(obj) {
		// "*" means any value, skip checking it
		return mj.JSONCheckBytesStar(), nil
	}

	jb, err := p.processSubTreeAsByteArray(obj)
	if err != nil {
		return mj.JSONCheckBytes{}, err
	}
	return mj.JSONCheckBytes{
		Value:    jb.Value,
		IsStar:   false,
		Original: jb.Original,
	}, nil
}

func (p *Parser) processStringAsByteArray(obj oj.OJsonObject) (mj.JSONBytesFromString, error) {
	strVal, err := p.parseString(obj)
	if err != nil {
		return mj.JSONBytesFromString{}, err
	}
	result, err := p.ExprInterpreter.InterpretString(strVal)
	return mj.NewJSONBytesFromString(result, strVal), err
}

func (p *Parser) processSubTreeAsByteArray(obj oj.OJsonObject) (mj.JSONBytesFromTree, error) {
	value, err := p.ExprInterpreter.InterpretSubTree(obj)
	return mj.JSONBytesFromTree{
		Value:    value,
		Original: obj,
	}, err
}

func (p *Parser) parseString(obj oj.OJsonObject) (string, error) {
	str, isStr := obj.(*oj.OJsonString)
	if !isStr {
		return "", errors.New("not a string value")
	}
	return str.Value, nil
}

// IsStar returns whether check object is othe form "*".
func IsStar(obj oj.OJsonObject) bool {
	str, isStr := obj.(*oj.OJsonString)
	if !isStr {
		return false
	}
	return str.Value == "*"
}
