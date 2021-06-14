package arwenmandos

import (
	"fmt"
	"math/big"

	mj "github.com/ElrondNetwork/arwen-wasm-vm/mandos-go/json/model"
	mjwrite "github.com/ElrondNetwork/arwen-wasm-vm/mandos-go/json/write"
	oj "github.com/ElrondNetwork/arwen-wasm-vm/mandos-go/orderedjson"
	vmi "github.com/ElrondNetwork/elrond-go/core/vmcommon"
)

func checkTxResults(
	txIndex string,
	blResult *mj.TransactionResult,
	checkGas bool,
	output *vmi.VMOutput,
) error {

	if !blResult.Status.Check(big.NewInt(int64(output.ReturnCode))) {
		return fmt.Errorf("result code mismatch. Tx %s. Want: %s. Have: %d (%s). Message: %s",
			txIndex, blResult.Status.Original, int(output.ReturnCode), output.ReturnCode.String(), output.ReturnMessage)
	}

	if !blResult.Message.Check([]byte(output.ReturnMessage)) {
		return fmt.Errorf("result message mismatch. Tx %s. Want: %s. Have: %s",
			txIndex, blResult.Message.Original, output.ReturnMessage)
	}

	// check result
	if len(output.ReturnData) != len(blResult.Out) {
		return fmt.Errorf("result length mismatch. Tx %s. Want: %s. Have: %s",
			txIndex,
			checkBytesListPretty(blResult.Out),
			mj.ResultAsString(output.ReturnData))
	}
	for i, expected := range blResult.Out {
		if !expected.Check(output.ReturnData[i]) {
			return fmt.Errorf("result mismatch. Tx %s. Want: %s. Have: %s",
				txIndex,
				checkBytesListPretty(blResult.Out),
				mj.ResultAsString(output.ReturnData))
		}
	}

	// check refund
	if !blResult.Refund.Check(output.GasRefund) {
		return fmt.Errorf("result gas refund mismatch. Tx %s. Want: %s. Have: 0x%x",
			txIndex, blResult.Refund.Original, output.GasRefund)
	}

	// check gas
	// unlike other checks, if unspecified the remaining gas check is ignored
	if checkGas && !blResult.Gas.IsUnspecified() && !blResult.Gas.Check(output.GasRemaining) {
		return fmt.Errorf("result gas mismatch. Tx %s. Want: %s. Got: %d (0x%x)",
			txIndex,
			blResult.Gas.Original,
			output.GasRemaining,
			output.GasRemaining)
	}

	// "logs": "*" means any value is accepted, log check ignored
	if blResult.LogsStar {
		return nil
	}

	// this is the real log check
	if len(blResult.Logs) != len(output.Logs) {
		return fmt.Errorf("wrong number of logs. Tx %s. Want:%d. Got:%d",
			txIndex,
			len(blResult.Logs),
			len(output.Logs))
	}
	for i, outLog := range output.Logs {
		testLog := blResult.Logs[i]
		if !testLog.Address.Check(outLog.Address) {
			return fmt.Errorf("bad log address. Tx %s. Want:\n%s\nGot:\n%s",
				txIndex,
				mjwrite.LogToString(testLog),
				mjwrite.LogToString(convertLogToTestFormat(outLog)))
		}
		if !testLog.Identifier.Check(outLog.Identifier) {
			return fmt.Errorf("bad log identifier. Tx %s. Want:\n%s\nGot:\n%s",
				txIndex,
				mjwrite.LogToString(testLog),
				mjwrite.LogToString(convertLogToTestFormat(outLog)))
		}
		if len(outLog.Topics) != len(testLog.Topics) {
			return fmt.Errorf("wrong number of log topics. Tx %s. Want:\n%s\nGot:\n%s",
				txIndex,
				mjwrite.LogToString(testLog),
				mjwrite.LogToString(convertLogToTestFormat(outLog)))
		}
		for ti := range outLog.Topics {
			if !testLog.Topics[ti].Check(outLog.Topics[ti]) {
				return fmt.Errorf("bad log topic. Tx %s. Want:\n%s\nGot:\n%s",
					txIndex,
					mjwrite.LogToString(testLog),
					mjwrite.LogToString(convertLogToTestFormat(outLog)))
			}
		}
		if !testLog.Data.Check(outLog.Data) {
			return fmt.Errorf("bad log data. Tx %s. Want:\n%s\nGot:\n%s",
				txIndex,
				mjwrite.LogToString(testLog),
				mjwrite.LogToString(convertLogToTestFormat(outLog)))
		}
	}

	return nil
}

// JSONCheckBytesString formats a list of JSONCheckBytes for printing to console.
// TODO: move somewhere else
func checkBytesListPretty(jcbs []mj.JSONCheckBytes) string {
	str := "["
	for i, jcb := range jcbs {
		if i > 0 {
			str += ", "
		}

		str += oj.JSONString(jcb.Original)
	}
	return str + "]"
}
