package arwenmandos

import (
	"errors"
	"fmt"
	"math/big"

	mj "github.com/ElrondNetwork/arwen-wasm-vm/mandos-go/json/model"
	worldmock "github.com/ElrondNetwork/arwen-wasm-vm/mock/world"
	"github.com/ElrondNetwork/elrond-go/core"
	"github.com/ElrondNetwork/elrond-go/core/vmcommon"
	"github.com/ElrondNetwork/elrond-go/data/esdt"
	"github.com/ElrondNetwork/elrond-go/process/smartContract/builtInFunctions"
)

func convertAccount(testAcct *mj.Account) (*worldmock.Account, error) {
	storage := make(map[string][]byte)
	for _, stkvp := range testAcct.Storage {
		key := string(stkvp.Key.Value)
		storage[key] = stkvp.Value.Value
	}

	if len(testAcct.Address.Value) != 32 {
		return nil, errors.New("bad test: account address should be 32 bytes long")
	}

	account := &worldmock.Account{
		Address:         testAcct.Address.Value,
		Nonce:           testAcct.Nonce.Value,
		Balance:         big.NewInt(0).Set(testAcct.Balance.Value),
		BalanceDelta:    big.NewInt(0),
		DeveloperReward: big.NewInt(0),
		Username:        testAcct.Username.Value,
		Storage:         storage,
		Code:            testAcct.Code.Value,
		OwnerAddress:    testAcct.Owner.Value,
		AsyncCallData:   testAcct.AsyncCallData,
		ShardID:         uint32(testAcct.Shard.Value),
		IsSmartContract: len(testAcct.Code.Value) > 0,
		CodeMetadata: (&vmcommon.CodeMetadata{
			Payable:     true,
			Upgradeable: true,
			Readable:    true,
		}).ToBytes(), // TODO: add explicit fields in mandos json
	}

	for _, mandosESDTData := range testAcct.ESDTData {
		tokenName := mandosESDTData.TokenIdentifier.Value
		isFrozen := mandosESDTData.Frozen.Value > 0
		for _, instance := range mandosESDTData.Instances {
			tokenNonce := instance.Nonce.Value
			tokenKey := worldmock.MakeTokenKey(tokenName, tokenNonce)
			tokenBalance := instance.Balance.Value
			tokenData := &esdt.ESDigitalToken{
				Value:      tokenBalance,
				Type:       uint32(core.Fungible),
				Properties: makeESDTUserMetadataBytes(isFrozen),
				TokenMetaData: &esdt.MetaData{
					Name:       tokenName,
					Nonce:      tokenNonce,
					Creator:    instance.Creator.Value,
					Royalties:  uint32(instance.Royalties.Value),
					Hash:       instance.Hash.Value,
					URIs:       [][]byte{instance.Uri.Value},
					Attributes: instance.Attributes.Value,
				},
			}
			err := account.SetTokenData(tokenKey, tokenData)
			if err != nil {
				return nil, err
			}
			err = account.SetLastNonce(tokenName, mandosESDTData.LastNonce.Value)
			if err != nil {
				return nil, err
			}
		}
		err := account.SetTokenRolesAsStrings(tokenName, mandosESDTData.Roles)
		if err != nil {
			return nil, err
		}
	}

	return account, nil
}

func validateSetStateAccount(mandosAccount *mj.Account, converted *worldmock.Account) error {
	err := converted.Validate()
	if err != nil {
		return fmt.Errorf(
			`"setState" step validation failed for account "%s": %w`,
			mandosAccount.Address.Original,
			err)
	}
	return nil
}

func makeESDTUserMetadataBytes(frozen bool) []byte {
	metadata := &builtInFunctions.ESDTUserMetadata{
		Frozen: frozen,
	}

	return metadata.ToBytes()
}

func validateNewAddressMocks(testNAMs []*mj.NewAddressMock) error {
	for _, testNAM := range testNAMs {
		if !worldmock.IsSmartContractAddress(testNAM.NewAddress.Value) {
			return fmt.Errorf(
				`address in "setState" "newAddresses" field should have SC format: %s`,
				testNAM.NewAddress.Original)
		}
	}
	return nil
}

func convertNewAddressMocks(testNAMs []*mj.NewAddressMock) []*worldmock.NewAddressMock {
	var result []*worldmock.NewAddressMock
	for _, testNAM := range testNAMs {
		result = append(result, &worldmock.NewAddressMock{
			CreatorAddress: testNAM.CreatorAddress.Value,
			CreatorNonce:   testNAM.CreatorNonce.Value,
			NewAddress:     testNAM.NewAddress.Value,
		})
	}
	return result
}

func convertBlockInfo(testBlockInfo *mj.BlockInfo) *worldmock.BlockInfo {
	if testBlockInfo == nil {
		return nil
	}

	var randomsSeed [48]byte
	if testBlockInfo.BlockRandomSeed != nil {
		copy(randomsSeed[:], testBlockInfo.BlockRandomSeed.Value)
	}

	result := &worldmock.BlockInfo{
		BlockTimestamp: testBlockInfo.BlockTimestamp.Value,
		BlockNonce:     testBlockInfo.BlockNonce.Value,
		BlockRound:     testBlockInfo.BlockRound.Value,
		BlockEpoch:     uint32(testBlockInfo.BlockEpoch.Value),
		RandomSeed:     &randomsSeed,
	}

	return result
}

// this is a small hack, so we can reuse mandos's JSON printing in error messages
func convertLogToTestFormat(outputLog *vmcommon.LogEntry) *mj.LogEntry {
	testLog := mj.LogEntry{
		Address:    mj.JSONCheckBytesReconstructed(outputLog.Address),
		Identifier: mj.JSONCheckBytesReconstructed(outputLog.Identifier),
		Data:       mj.JSONCheckBytesReconstructed(outputLog.Data),
		Topics:     make([]mj.JSONCheckBytes, len(outputLog.Topics)),
	}
	for i, topic := range outputLog.Topics {
		testLog.Topics[i] = mj.JSONCheckBytesReconstructed(topic)
	}

	return &testLog
}

func generateTxHash(txIndex string) []byte {
	txIndexBytes := []byte(txIndex)
	if len(txIndexBytes) > 32 {
		return txIndexBytes[:32]
	}
	for i := len(txIndexBytes); i < 32; i++ {
		txIndexBytes = append(txIndexBytes, '.')
	}
	return txIndexBytes
}

func addESDTToVMInput(esdtData *mj.ESDTTxData, vmInput *vmcommon.VMInput) {
	if esdtData != nil {
		vmInput.ESDTTokenName = esdtData.TokenIdentifier.Value
		vmInput.ESDTValue = esdtData.Value.Value
		vmInput.ESDTTokenNonce = esdtData.Nonce.Value
		if vmInput.ESDTTokenNonce != 0 {
			vmInput.ESDTTokenType = uint32(core.NonFungible)
		} else {
			vmInput.ESDTTokenType = uint32(core.Fungible)
		}
	}
}
