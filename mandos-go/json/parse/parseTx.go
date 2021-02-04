package mandosjsonparse

import (
	"errors"
	"fmt"

	mj "github.com/ElrondNetwork/arwen-wasm-vm/mandos-go/json/model"
	oj "github.com/ElrondNetwork/arwen-wasm-vm/mandos-go/orderedjson"
)

func (p *Parser) processTx(txType mj.TransactionType, blrRaw oj.OJsonObject) (*mj.Transaction, error) {
	bltMap, isMap := blrRaw.(*oj.OJsonMap)
	if !isMap {
		return nil, errors.New("unmarshalled transaction is not a map")
	}

	blt := mj.Transaction{
		Type:      txType,
		Value:     mj.JSONBigIntZero(),
		ESDTValue: mj.JSONBigIntZero(),
	}
	var err error
	for _, kvp := range bltMap.OrderedKV {

		switch kvp.Key {
		case "nonce":
			blt.Nonce, err = p.processUint64(kvp.Value)
			if err != nil {
				return nil, fmt.Errorf("invalid transaction nonce: %w", err)
			}
		case "from":
			if !txType.HasSender() {
				return nil, errors.New("`from` not allowed in transaction, it is always the zero address")
			}
			fromStr, err := p.parseString(kvp.Value)
			if err != nil {
				return nil, fmt.Errorf("invalid transaction from: %w", err)
			}
			var fromErr error
			blt.From, fromErr = p.parseAccountAddress(fromStr)
			if fromErr != nil {
				return nil, fromErr
			}
		case "to":
			toStr, err := p.parseString(kvp.Value)
			if err != nil {
				return nil, fmt.Errorf("invalid transaction to: %w", err)
			}

			if txType == mj.ScDeploy {
				if len(toStr) > 0 {
					return nil, errors.New("transaction to field not allowed for scDeploy transactions")
				}
			} else {
				blt.To, err = p.parseAccountAddress(toStr)
				if err != nil {
					return nil, err
				}
			}
		case "function":
			blt.Function, err = p.parseString(kvp.Value)
			if err != nil {
				return nil, fmt.Errorf("invalid transaction function: %w", err)
			}
			if txType == mj.ScDeploy && len(blt.Function) > 0 {
				return nil, errors.New("transaction function field not allowed for scDeploy transactions")
			}
			if txType == mj.Transfer && len(blt.Function) > 0 {
				return nil, errors.New("transaction function field not allowed for transfer transactions")
			}
		case "value":
			blt.Value, err = p.processBigInt(kvp.Value, bigIntUnsignedBytes)
			if err != nil {
				return nil, fmt.Errorf("invalid transaction value: %w", err)
			}
		case "esdtValue":
			blt.ESDTValue, err = p.processBigInt(kvp.Value, bigIntUnsignedBytes)
			if err != nil {
				return nil, fmt.Errorf("invalid transaction ESDT value: %w", err)
			}
		case "esdtTokenName":
			blt.ESDTTokenName, err = p.processStringAsByteArray(kvp.Value)
			if err != nil {
				return nil, fmt.Errorf("invalid transaction ESDT token name: %w", err)
			}
		case "arguments":
			blt.Arguments, err = p.parseSubTreeList(kvp.Value)
			if err != nil {
				return nil, fmt.Errorf("invalid transaction arguments: %w", err)
			}
			if txType == mj.Transfer && len(blt.Arguments) > 0 {
				return nil, errors.New("function arguments not allowed for transfer transactions")
			}
		case "contractCode":
			blt.Code, err = p.processStringAsByteArray(kvp.Value)
			if err != nil {
				return nil, fmt.Errorf("invalid transaction contract code: %w", err)
			}
			if txType != mj.ScDeploy && len(blt.Code.Value) > 0 {
				return nil, errors.New("transaction contractCode field only allowed int scDeploy transactions")
			}
		case "gasPrice":
			blt.GasPrice, err = p.processUint64(kvp.Value)
			if err != nil {
				return nil, fmt.Errorf("invalid transaction gasPrice: %w", err)
			}
		case "gasLimit":
			blt.GasLimit, err = p.processUint64(kvp.Value)
			if err != nil {
				return nil, fmt.Errorf("invalid transaction gasLimit: %w", err)
			}
		default:
			return nil, fmt.Errorf("unknown field in transaction: %w", err)
		}
	}

	return &blt, nil
}
