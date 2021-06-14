package mandosjsonparse

import (
	"errors"
	"fmt"

	mj "github.com/ElrondNetwork/arwen-wasm-vm/mandos-go/json/model"
	oj "github.com/ElrondNetwork/arwen-wasm-vm/mandos-go/orderedjson"
)

func (p *Parser) processCheckESDTData(
	tokenName mj.JSONBytesFromString,
	esdtDataRaw oj.OJsonObject) (*mj.CheckESDTData, error) {

	switch data := esdtDataRaw.(type) {
	case *oj.OJsonString:
		// simple string representing balance "400,000,000,000"
		esdtData := mj.CheckESDTData{
			TokenIdentifier: tokenName,
		}
		balance, err := p.processCheckBigInt(esdtDataRaw, bigIntUnsignedBytes)
		if err != nil {
			return nil, fmt.Errorf("invalid ESDT balance: %w", err)
		}
		esdtData.Instances = []*mj.CheckESDTInstance{
			{
				Nonce:   mj.JSONCheckUint64{Value: 0, Original: ""},
				Balance: balance,
			},
		}
		return &esdtData, nil
	case *oj.OJsonMap:
		return p.processCheckESDTDataMap(tokenName, data)
	default:
		return nil, errors.New("invalid JSON object for ESDT")
	}
}

// map containing other fields too, e.g.:
// {
// 	"balance": "400,000,000,000",
// 	"frozen": "true"
// }
func (p *Parser) processCheckESDTDataMap(tokenName mj.JSONBytesFromString, esdtDataMap *oj.OJsonMap) (*mj.CheckESDTData, error) {
	esdtData := mj.CheckESDTData{
		TokenIdentifier: tokenName,
	}
	// var err error
	firstInstance := &mj.CheckESDTInstance{
		Nonce:      mj.JSONCheckUint64Unspecified(),
		Balance:    mj.JSONCheckBigIntUnspecified(),
		Creator:    mj.JSONCheckBytesUnspecified(),
		Royalties:  mj.JSONCheckUint64Unspecified(),
		Hash:       mj.JSONCheckBytesUnspecified(),
		Uri:        mj.JSONCheckBytesUnspecified(),
		Attributes: mj.JSONCheckBytesUnspecified(),
	}
	firstInstanceLoaded := false
	var explicitInstances []*mj.CheckESDTInstance

	for _, kvp := range esdtDataMap.OrderedKV {
		// it is allowed to load the instance directly, fields set to the first instance
		instanceFieldLoaded, err := p.tryProcessCheckESDTInstanceField(kvp, firstInstance)
		if err != nil {
			return nil, fmt.Errorf("invalid account ESDT instance field: %w", err)
		}
		if instanceFieldLoaded {
			firstInstanceLoaded = true
		} else {
			switch kvp.Key {
			case "instances":
				explicitInstances, err = p.processCheckESDTInstances(kvp.Value)
				if err != nil {
					return nil, fmt.Errorf("invalid account ESDT instances: %w", err)
				}
			case "lastNonce":
				esdtData.LastNonce, err = p.processCheckUint64(kvp.Value)
				if err != nil {
					return nil, fmt.Errorf("invalid account ESDT lastNonce: %w", err)
				}
			case "roles":
				esdtData.Roles, err = p.processStringList(kvp.Value)
				if err != nil {
					return nil, fmt.Errorf("invalid account ESDT roles: %w", err)
				}
			case "frozen":
				esdtData.Frozen, err = p.processCheckUint64(kvp.Value)
				if err != nil {
					return nil, fmt.Errorf("invalid ESDT frozen flag: %w", err)
				}
			default:
				return nil, fmt.Errorf("unknown ESDT data field: %s", kvp.Key)
			}
		}
	}

	if firstInstanceLoaded {
		esdtData.Instances = []*mj.CheckESDTInstance{firstInstance}
	}
	esdtData.Instances = append(esdtData.Instances, explicitInstances...)

	return &esdtData, nil
}

func (p *Parser) tryProcessCheckESDTInstanceField(kvp *oj.OJsonKeyValuePair, targetInstance *mj.CheckESDTInstance) (bool, error) {
	var err error
	switch kvp.Key {
	case "nonce":
		targetInstance.Nonce, err = p.processCheckUint64(kvp.Value)
		if err != nil {
			return false, fmt.Errorf("invalid account nonce: %w", err)
		}
	case "balance":
		targetInstance.Balance, err = p.processCheckBigInt(kvp.Value, bigIntUnsignedBytes)
		if err != nil {
			return false, fmt.Errorf("invalid ESDT balance: %w", err)
		}
	case "creator":
		targetInstance.Creator, err = p.parseCheckBytes(kvp.Value)
		if err != nil {
			return false, fmt.Errorf("invalid ESDT NFT creator address: %w", err)
		}
	case "royalties":
		targetInstance.Royalties, err = p.processCheckUint64(kvp.Value)
		if err != nil {
			return false, fmt.Errorf("invalid ESDT NFT royalties: %w", err)
		}
		if targetInstance.Royalties.Value > 10000 {
			return false, errors.New("invalid ESDT NFT royalties: value exceeds maximum allowed 10000")
		}
	case "hash":
		targetInstance.Hash, err = p.parseCheckBytes(kvp.Value)
		if err != nil {
			return false, fmt.Errorf("invalid ESDT NFT hash: %w", err)
		}
	case "uri":
		targetInstance.Uri, err = p.parseCheckBytes(kvp.Value)
		if err != nil {
			return false, fmt.Errorf("invalid ESDT NFT URI: %w", err)
		}
	case "attributes":
		targetInstance.Attributes, err = p.parseCheckBytes(kvp.Value)
		if err != nil {
			return false, fmt.Errorf("invalid ESDT NFT attributes: %w", err)
		}
	default:
		return false, nil
	}
	return true, nil
}

func (p *Parser) processCheckESDTInstances(esdtInstancesRaw oj.OJsonObject) ([]*mj.CheckESDTInstance, error) {
	var instancesResult []*mj.CheckESDTInstance
	esdtInstancesList, isList := esdtInstancesRaw.(*oj.OJsonList)
	if !isList {
		return nil, errors.New("esdt instances object is not a list")
	}
	for _, instanceItem := range esdtInstancesList.AsList() {
		instanceAsMap, isMap := instanceItem.(*oj.OJsonMap)
		if !isMap {
			return nil, errors.New("JSON map expected as esdt instances list item")
		}

		instance := mj.NewCheckESDTInstance()

		for _, kvp := range instanceAsMap.OrderedKV {
			instanceFieldLoaded, err := p.tryProcessCheckESDTInstanceField(kvp, instance)
			if err != nil {
				return nil, fmt.Errorf("invalid account ESDT instance field in instances list: %w", err)
			}
			if !instanceFieldLoaded {
				return nil, fmt.Errorf("invalid account ESDT instance field in instances list: `%s`", kvp.Key)
			}
		}

		instancesResult = append(instancesResult, instance)

	}

	return instancesResult, nil
}
