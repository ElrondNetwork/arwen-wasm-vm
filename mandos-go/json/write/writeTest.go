package mandosjsonwrite

import (
	mj "github.com/ElrondNetwork/arwen-wasm-vm/mandos-go/json/model"
	oj "github.com/ElrondNetwork/arwen-wasm-vm/mandos-go/orderedjson"
)

// TestToJSONString converts a test object to its JSON representation.
func TestToJSONString(testTopLevel []*mj.Test) string {
	jobj := TestToOrderedJSON(testTopLevel)
	return oj.JSONString(jobj) + "\n"
}

// TestToOrderedJSON converts a test object to an ordered JSON object.
func TestToOrderedJSON(testTopLevel []*mj.Test) oj.OJsonObject {
	result := oj.NewMap()
	for _, test := range testTopLevel {
		result.Put(test.TestName, testToOJ(test))
	}

	return result
}

func testToOJ(test *mj.Test) oj.OJsonObject {
	testOJ := oj.NewMap()

	if !test.CheckGas {
		ojFalse := oj.OJsonBool(false)
		testOJ.Put("checkGas", &ojFalse)
	}

	testOJ.Put("pre", AccountsToOJ(test.Pre))

	var blockList []oj.OJsonObject
	for _, block := range test.Blocks {
		blockList = append(blockList, blockToOJ(block))
	}
	blocksOJ := oj.OJsonList(blockList)
	testOJ.Put("blocks", &blocksOJ)
	testOJ.Put("network", stringToOJ(test.Network))
	testOJ.Put("blockHashes", blockHashesToOJ(test.BlockHashes))
	testOJ.Put("postState", checkAccountsToOJ(test.PostState))
	return testOJ
}

func transactionToTestOJ(tx *mj.Transaction) oj.OJsonObject {
	transactionOJ := oj.NewMap()
	transactionOJ.Put("nonce", uint64ToOJ(tx.Nonce))
	transactionOJ.Put("function", stringToOJ(tx.Function))
	transactionOJ.Put("gasLimit", uint64ToOJ(tx.GasLimit))
	transactionOJ.Put("value", bigIntToOJ(tx.Value))
	transactionOJ.Put("to", bytesFromStringToOJ(tx.To))

	var argList []oj.OJsonObject
	for _, arg := range tx.Arguments {
		argList = append(argList, bytesFromTreeToOJ(arg))
	}
	argOJ := oj.OJsonList(argList)
	transactionOJ.Put("arguments", &argOJ)

	if len(tx.Code.Original) > 0 {
		transactionOJ.Put("contractCode", bytesFromStringToOJ(tx.Code))
	}
	transactionOJ.Put("gasPrice", uint64ToOJ(tx.GasPrice))
	transactionOJ.Put("from", bytesFromStringToOJ(tx.From))

	return transactionOJ
}

func blockToOJ(block *mj.Block) oj.OJsonObject {
	blockOJ := oj.NewMap()

	var resultList []oj.OJsonObject
	for _, blr := range block.Results {
		resultList = append(resultList, resultToOJ(blr))
	}
	resultsOJ := oj.OJsonList(resultList)
	blockOJ.Put("results", &resultsOJ)

	var txList []oj.OJsonObject
	for _, tx := range block.Transactions {
		txList = append(txList, transactionToTestOJ(tx))
	}
	txsOJ := oj.OJsonList(txList)
	blockOJ.Put("transactions", &txsOJ)

	blockHeaderOJ := oj.NewMap()
	blockHeaderOJ.Put("gasLimit", bigIntToOJ(block.BlockHeader.GasLimit))
	blockHeaderOJ.Put("number", bigIntToOJ(block.BlockHeader.Number))
	blockHeaderOJ.Put("difficulty", bigIntToOJ(block.BlockHeader.Difficulty))
	blockHeaderOJ.Put("timestamp", uint64ToOJ(block.BlockHeader.Timestamp))
	blockHeaderOJ.Put("coinbase", bigIntToOJ(block.BlockHeader.Beneficiary))
	blockOJ.Put("blockHeader", blockHeaderOJ)

	return blockOJ
}
