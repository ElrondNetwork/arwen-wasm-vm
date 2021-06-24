package fuzzForwarder

import (
	"fmt"
	"strings"

	vmi "github.com/ElrondNetwork/elrond-go/core/vmcommon"
)

func (pfe *fuzzExecutor) executeCallPrintLogs(toIndex int) error {
	el, _ := pfe.getAllExpectedLogs(toIndex)
	fmt.Println(el)
	output, err := pfe.executeCall(toIndex, "")
	if err != nil {
		return err
	}
	for _, log := range output.Logs {
		fmt.Println(pfe.arwenTestExecutor.LogToString(log))
	}

	return nil
}

func (pfe *fuzzExecutor) executeCallCheckLogs(toIndex int) error {
	expectedLogs, err := pfe.getAllExpectedLogs(toIndex)
	if err != nil {
		return err
	}
	expectLogs := fmt.Sprintf(`,
		"expect": {
			"out": [],
			"status": "",
			"logs": [%s
			]
		}
	`,
		expectedLogs,
	)
	_, err = pfe.executeCall(toIndex, expectLogs)
	if err != nil {
		return err
	}

	return nil
}

func (pfe *fuzzExecutor) executeCall(toIndex int, mandosTxExpect string) (*vmi.VMOutput, error) {
	return pfe.executeTxStep(fmt.Sprintf(`
	{
		"step": "scCall",
		"txId": "%d",
		"tx": {
			"from": "%s",
			"to": "%s",
			"value": "0",
			"function": "forward_programmed_calls",
			"arguments": [ "%d" ],
			"gasLimit": "18,000,000,000,000,000",
			"gasPrice": "0"
		}
		%s
	}`,
		pfe.nextTxIndex(),
		pfe.data.mainCallerAddress,
		pfe.forwarderAddress(toIndex),
		pfe.data.maxCallDepth,
		mandosTxExpect,
	))
}

func (pfe *fuzzExecutor) popCall(forwarderIndex int) *programmedCall {
	if len(pfe.data.programmedCalls[forwarderIndex]) == 0 {
		return nil
	}
	result := pfe.data.programmedCalls[forwarderIndex][0]
	pfe.data.programmedCalls[forwarderIndex] = pfe.data.programmedCalls[forwarderIndex][1:]
	return result
}

func (pfe *fuzzExecutor) getAllExpectedLogs(toIndex int) (string, error) {
	var sb strings.Builder
	initialCall := &programmedCall{
		fromIndex: -1,
		toIndex:   toIndex,
		token:     "EGLD",
		nonce:     0,
		amount:    "0",
	}
	err := pfe.getExpectedLogsForCall(initialCall, pfe.data.maxCallDepth, &sb)
	if err != nil {
		return "", err
	}
	return sb.String(), nil
}

func (pfe *fuzzExecutor) getExpectedLogsForCall(call *programmedCall, maxDepth int, sb *strings.Builder) error {
	if sb.Len() > 0 {
		sb.WriteString(",")
	}
	sb.WriteString(fmt.Sprintf(`
				{
					"address": "%s",
					"identifier": "str:forward_programmed_calls",
					"topics": [
						"%s",
						"str:%s",
						"%d",
						"%s"
					],
					"data": ""
				}`,
		pfe.forwarderAddress(call.toIndex),
		pfe.forwarderAddress(call.fromIndex),
		call.token,
		call.nonce,
		call.amount,
	))
	pfe.log("%d calls %d with token %s, nonce %d, max depth: %d, call type: %s",
		call.fromIndex, call.toIndex, call.token, call.nonce, maxDepth, call.callType.String())

	if maxDepth == 0 {
		return nil
	}

	nextCall := pfe.popCall(call.toIndex)
	for nextCall != nil {
		err := pfe.getExpectedLogsForCall(nextCall, maxDepth-1, sb)
		if err != nil {
			return err
		}
		if nextCall.callType == asyncCall {
			return nil // exit early, just like the async call does
		}

		nextCall = pfe.popCall(call.toIndex)
	}

	return nil
}
