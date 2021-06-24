package fuzzForwarder

import (
	"fmt"
)

func (pfe *fuzzExecutor) programCall(
	callType programmedCallType,
	fromIndex, toIndex int,
	tokenName string, nonce int, amount string) error {

	pfe.log("%d will call %d with token %s, nonce %d, call type: %s", fromIndex, toIndex, tokenName, nonce, callType.String())

	pfe.data.programmedCalls[fromIndex] = append(pfe.data.programmedCalls[fromIndex], &programmedCall{
		callType:  callType,
		fromIndex: fromIndex,
		toIndex:   toIndex,
		token:     tokenName,
		nonce:     nonce,
		amount:    amount,
	})

	_, err := pfe.executeTxStep(fmt.Sprintf(`
	{
		"step": "scCall",
		"txId": "%d",
		"tx": {
			"from": "%s",
			"to": "%s",
			"value": "0",
			"function": "add_programmed_call",
			"arguments": [
				"%d",
				"%s",
				"str:%s",
				"%d",
				"%s"
			],
			"gasLimit": "50,000,000",
			"gasPrice": "0"
		},
		"expect": {
			"out": [],
			"status": ""
		}
	}`,
		pfe.nextTxIndex(),
		pfe.data.mainCallerAddress,
		pfe.forwarderAddress(fromIndex),
		callType,
		pfe.forwarderAddress(toIndex),
		tokenName,
		nonce,
		amount,
	))
	if err != nil {
		return err
	}

	return nil
}
