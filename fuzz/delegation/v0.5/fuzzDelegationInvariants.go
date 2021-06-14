package delegation

import (
	"fmt"
)

func (pfe *fuzzDelegationExecutor) validateOwnerStakeShare() error {
	_, err := pfe.executeTxStep(fmt.Sprintf(`
	{
		"step": "scCall",
		"txId": "%d",
		"tx": {
			"from": "%s",
			"to": "%s",
			"value": "0",
			"function": "validateOwnerStakeShare",
			"arguments": [],
			"gasLimit": "50,000,000",
			"gasPrice": "0"
		},
		"expect": {
			"out": [],
			"status": "",
			"logs": [],
			"gas": "*",
			"refund": "*"
		}
	}`,
		pfe.nextTxIndex(),
		pfe.ownerAddress,
		pfe.delegationContractAddress,
	))
	if err != nil {
		return err
	}

	return nil
}

func (pfe *fuzzDelegationExecutor) validateDelegationCapInvariant() error {
	_, err := pfe.executeTxStep(fmt.Sprintf(`
	{
		"step": "scCall",
		"txId": "%d",
		"tx": {
			"from": "%s",
			"to": "%s",
			"value": "0",
			"function": "validateDelegationCapInvariant",
			"arguments": [],
			"gasLimit": "50,000,000",
			"gasPrice": "0"
		},
		"expect": {
			"out": [],
			"status": "",
			"logs": [],
			"gas": "*",
			"refund": "*"
		}
	}`,
		pfe.nextTxIndex(),
		pfe.ownerAddress,
		pfe.delegationContractAddress,
	))
	if err != nil {
		return err
	}

	return nil
}
