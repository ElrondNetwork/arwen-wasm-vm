package mandosjsonparse

import (
	"errors"
	"fmt"

	mj "github.com/ElrondNetwork/arwen-wasm-vm/mandos-go/json/model"
	oj "github.com/ElrondNetwork/arwen-wasm-vm/mandos-go/orderedjson"
)

func (p *Parser) processESDTData(
	tokenName mj.JSONBytesFromString,
	esdtDataRaw oj.OJsonObject) (*mj.ESDTData, error) {

	switch data := esdtDataRaw.(type) {
	case *oj.OJsonString:
		// simple string representing balance "400,000,000,000"
		esdtData := mj.ESDTData{
			TokenIdentifier: tokenName,
		}
		balance, err := p.processBigInt(esdtDataRaw, bigIntUnsignedBytes)
		if err != nil {
			return nil, fmt.Errorf("invalid ESDT balance: %w", err)
		}
		esdtData.Instances = []*mj.ESDTInstance{
			{
				Nonce:   mj.JSONUint64{Value: 0, Original: ""},
				Balance: balance,
			},
		}
		return &esdtData, nil
	case *oj.OJsonMap:
		return p.processESDTDataMap(tokenName, data)
	default:
		return nil, errors.New("invalid JSON object for ESDT")
	}
}

// map containing other fields too, e.g.:
// {
// 	"balance": "400,000,000,000",
// 	"frozen": "true"
// }
func (p *Parser) processESDTDataMap(tokenName mj.JSONBytesFromString, esdtDataMap *oj.OJsonMap) (*mj.ESDTData, error) {
	esdtData := mj.ESDTData{
		TokenIdentifier: tokenName,
	}
	// var err error
	firstInstance := &mj.ESDTInstance{}
	firstInstanceLoaded := false
	var explicitInstances []*mj.ESDTInstance

	for _, kvp := range esdtDataMap.OrderedKV {
		// it is allowed to load the instance directly, fields set to the first instance
		instanceFieldLoaded, err := p.tryProcessESDTInstanceField(kvp, firstInstance)
		if err != nil {
			return nil, fmt.Errorf("invalid account ESDT instance field: %w", err)
		}
		if instanceFieldLoaded {
			firstInstanceLoaded = true
		} else {
			switch kvp.Key {
			case "instances":
				explicitInstances, err = p.processESDTInstances(kvp.Value)
				if err != nil {
					return nil, fmt.Errorf("invalid account ESDT instances: %w", err)
				}
			case "lastNonce":
				esdtData.LastNonce, err = p.processUint64(kvp.Value)
				if err != nil {
					return nil, fmt.Errorf("invalid account ESDT lastNonce: %w", err)
				}
			case "roles":
				esdtData.Roles, err = p.processStringList(kvp.Value)
				if err != nil {
					return nil, fmt.Errorf("invalid account ESDT roles: %w", err)
				}
			case "frozen":
				esdtData.Frozen, err = p.processUint64(kvp.Value)
				if err != nil {
					return nil, fmt.Errorf("invalid ESDT frozen flag: %w", err)
				}
			default:
				return nil, fmt.Errorf("unknown ESDT data field: %s", kvp.Key)
			}
		}
	}

	if firstInstanceLoaded {
		esdtData.Instances = []*mj.ESDTInstance{firstInstance}
	}
	esdtData.Instances = append(esdtData.Instances, explicitInstances...)

	return &esdtData, nil
}

func (p *Parser) tryProcessESDTInstanceField(kvp *oj.OJsonKeyValuePair, targetInstance *mj.ESDTInstance) (bool, error) {
	var err error
	switch kvp.Key {
	case "nonce":
		targetInstance.Nonce, err = p.processUint64(kvp.Value)
		if err != nil {
			return false, fmt.Errorf("invalid account nonce: %w", err)
		}
	case "balance":
		targetInstance.Balance, err = p.processBigInt(kvp.Value, bigIntUnsignedBytes)
		if err != nil {
			return false, fmt.Errorf("invalid ESDT balance: %w", err)
		}
	case "creator":
		targetInstance.Creator, err = p.processStringAsByteArray(kvp.Value)
		if err != nil || len(targetInstance.Creator.Value) != 32 {
			return false, fmt.Errorf("invalid ESDT NFT creator address: %w", err)
		}
	case "royalties":
		targetInstance.Royalties, err = p.processUint64(kvp.Value)
		if err != nil || targetInstance.Royalties.Value > 10000 {
			return false, fmt.Errorf("invalid ESDT NFT royalties: %w", err)
		}
	case "hash":
		targetInstance.Hash, err = p.processStringAsByteArray(kvp.Value)
		if err != nil {
			return false, fmt.Errorf("invalid ESDT NFT hash: %w", err)
		}
	case "uri":
		targetInstance.Uri, err = p.processSubTreeAsByteArray(kvp.Value)
		if err != nil {
			return false, fmt.Errorf("invalid ESDT NFT URI: %w", err)
		}
	case "attributes":
		targetInstance.Attributes, err = p.processStringAsByteArray(kvp.Value)
		if err != nil {
			return false, fmt.Errorf("invalid ESDT NFT attributes: %w", err)
		}
	default:
		return false, nil
	}
	return true, nil
}

func (p *Parser) processESDTInstances(esdtInstancesRaw oj.OJsonObject) ([]*mj.ESDTInstance, error) {
	var instancesResult []*mj.ESDTInstance
	esdtInstancesList, isList := esdtInstancesRaw.(*oj.OJsonList)
	if !isList {
		return nil, errors.New("esdt instances object is not a list")
	}
	for _, instanceItem := range esdtInstancesList.AsList() {
		instanceAsMap, isMap := instanceItem.(*oj.OJsonMap)
		if !isMap {
			return nil, errors.New("JSON map expected as esdt instances list item")
		}

		instance := &mj.ESDTInstance{}

		for _, kvp := range instanceAsMap.OrderedKV {
			instanceFieldLoaded, err := p.tryProcessESDTInstanceField(kvp, instance)
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
