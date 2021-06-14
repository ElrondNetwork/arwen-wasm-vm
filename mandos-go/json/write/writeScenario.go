package mandosjsonwrite

import (
	mj "github.com/ElrondNetwork/arwen-wasm-vm/mandos-go/json/model"
	oj "github.com/ElrondNetwork/arwen-wasm-vm/mandos-go/orderedjson"
)

// ScenarioToJSONString converts a scenario object to its JSON representation.
func ScenarioToJSONString(scenario *mj.Scenario) string {
	jobj := ScenarioToOrderedJSON(scenario)
	return oj.JSONString(jobj) + "\n"
}

// ScenarioToOrderedJSON converts a scenario object to an ordered JSON object.
func ScenarioToOrderedJSON(scenario *mj.Scenario) oj.OJsonObject {
	scenarioOJ := oj.NewMap()

	if len(scenario.Name) > 0 {
		scenarioOJ.Put("name", stringToOJ(scenario.Name))
	}

	if len(scenario.Comment) > 0 {
		scenarioOJ.Put("comment", stringToOJ(scenario.Comment))
	}

	if !scenario.CheckGas {
		ojFalse := oj.OJsonBool(false)
		scenarioOJ.Put("checkGas", &ojFalse)
	}

	scenarioOJ.Put("gasSchedule", gasScheduleToOJ(scenario.GasSchedule))

	var stepOJList []oj.OJsonObject

	for _, generalStep := range scenario.Steps {
		stepOJ := oj.NewMap()
		stepOJ.Put("step", stringToOJ(generalStep.StepTypeName()))
		switch step := generalStep.(type) {
		case *mj.ExternalStepsStep:
			if len(step.Comment) > 0 {
				stepOJ.Put("comment", stringToOJ(step.Comment))
			}
			stepOJ.Put("path", stringToOJ(step.Path))
		case *mj.SetStateStep:
			if len(step.Comment) > 0 {
				stepOJ.Put("comment", stringToOJ(step.Comment))
			}
			if len(step.Accounts) > 0 {
				stepOJ.Put("accounts", AccountsToOJ(step.Accounts))
			}
			if len(step.NewAddressMocks) > 0 {
				stepOJ.Put("newAddresses", newAddressMocksToOJ(step.NewAddressMocks))
			}
			if step.PreviousBlockInfo != nil {
				stepOJ.Put("previousBlockInfo", blockInfoToOJ(step.PreviousBlockInfo))
			}
			if step.CurrentBlockInfo != nil {
				stepOJ.Put("currentBlockInfo", blockInfoToOJ(step.CurrentBlockInfo))
			}
			if len(step.BlockHashes) > 0 {
				stepOJ.Put("blockHashes", blockHashesToOJ(step.BlockHashes))
			}
		case *mj.CheckStateStep:
			if len(step.Comment) > 0 {
				stepOJ.Put("comment", stringToOJ(step.Comment))
			}
			stepOJ.Put("accounts", checkAccountsToOJ(step.CheckAccounts))
		case *mj.DumpStateStep:
			if len(step.Comment) > 0 {
				stepOJ.Put("comment", stringToOJ(step.Comment))
			}
		case *mj.TxStep:
			if len(step.TxIdent) > 0 {
				stepOJ.Put("txId", stringToOJ(step.TxIdent))
			}
			if len(step.Comment) > 0 {
				stepOJ.Put("comment", stringToOJ(step.Comment))
			}
			stepOJ.Put("tx", transactionToScenarioOJ(step.Tx))
			if step.Tx.Type.IsSmartContractTx() && step.ExpectedResult != nil {
				stepOJ.Put("expect", resultToOJ(step.ExpectedResult))
			}
		}

		stepOJList = append(stepOJList, stepOJ)
	}

	stepsOJ := oj.OJsonList(stepOJList)
	scenarioOJ.Put("steps", &stepsOJ)

	return scenarioOJ
}

func transactionToScenarioOJ(tx *mj.Transaction) oj.OJsonObject {
	transactionOJ := oj.NewMap()
	if tx.Type.HasSender() {
		transactionOJ.Put("from", bytesFromStringToOJ(tx.From))
	}
	if tx.Type.HasReceiver() {
		transactionOJ.Put("to", bytesFromStringToOJ(tx.To))
	}
	if tx.Type.HasValue() {
		transactionOJ.Put("value", bigIntToOJ(tx.Value))
	}
	if tx.ESDTValue != nil {
		esdtItemOJ := esdtTxDataToOJ(tx.ESDTValue)
		transactionOJ.Put("esdt", esdtItemOJ)
	}
	if tx.Type.HasFunction() {
		transactionOJ.Put("function", stringToOJ(tx.Function))
	}
	if tx.Type == mj.ScDeploy {
		transactionOJ.Put("contractCode", bytesFromStringToOJ(tx.Code))
	}

	if tx.Type.HasFunction() || tx.Type == mj.ScDeploy {
		var argList []oj.OJsonObject
		for _, arg := range tx.Arguments {
			argList = append(argList, bytesFromTreeToOJ(arg))
		}
		argOJ := oj.OJsonList(argList)
		transactionOJ.Put("arguments", &argOJ)
	}

	if tx.Type.HasGas() {
		transactionOJ.Put("gasLimit", uint64ToOJ(tx.GasLimit))
		transactionOJ.Put("gasPrice", uint64ToOJ(tx.GasPrice))
	}

	return transactionOJ
}

func newAddressMocksToOJ(newAddressMocks []*mj.NewAddressMock) oj.OJsonObject {
	var namList []oj.OJsonObject
	for _, namEntry := range newAddressMocks {
		namOJ := oj.NewMap()
		namOJ.Put("creatorAddress", bytesFromStringToOJ(namEntry.CreatorAddress))
		namOJ.Put("creatorNonce", uint64ToOJ(namEntry.CreatorNonce))
		namOJ.Put("newAddress", bytesFromStringToOJ(namEntry.NewAddress))
		namList = append(namList, namOJ)
	}
	namOJList := oj.OJsonList(namList)
	return &namOJList
}

func blockInfoToOJ(blockInfo *mj.BlockInfo) oj.OJsonObject {
	blockInfoOJ := oj.NewMap()
	if len(blockInfo.BlockTimestamp.Original) > 0 {
		blockInfoOJ.Put("blockTimestamp", uint64ToOJ(blockInfo.BlockTimestamp))
	}
	if len(blockInfo.BlockNonce.Original) > 0 {
		blockInfoOJ.Put("blockNonce", uint64ToOJ(blockInfo.BlockNonce))
	}
	if len(blockInfo.BlockRound.Original) > 0 {
		blockInfoOJ.Put("blockRound", uint64ToOJ(blockInfo.BlockRound))
	}
	if len(blockInfo.BlockEpoch.Original) > 0 {
		blockInfoOJ.Put("blockEpoch", uint64ToOJ(blockInfo.BlockEpoch))
	}
	if blockInfo.BlockRandomSeed != nil {
		blockInfoOJ.Put("blockRandomSeed", bytesFromTreeToOJ(*blockInfo.BlockRandomSeed))
	}

	return blockInfoOJ
}

func gasScheduleToOJ(gasSchedule mj.GasSchedule) oj.OJsonObject {
	switch gasSchedule {
	case mj.GasScheduleDefault:
		return stringToOJ("default")
	case mj.GasScheduleDummy:
		return stringToOJ("dummy")
	case mj.GasScheduleV1:
		return stringToOJ("v1")
	case mj.GasScheduleV2:
		return stringToOJ("v2")
	case mj.GasScheduleV3:
		return stringToOJ("v3")
	default:
		return stringToOJ("")
	}
}
