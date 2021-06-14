package delegation

import (
	"fmt"
	"math/big"

	vmi "github.com/ElrondNetwork/elrond-go/core/vmcommon"
)

func (pfe *fuzzDelegationExecutor) stake(delegIndex int, amount *big.Int) error {
	// keep track of stake added
	pfe.totalStakeAdded.Add(pfe.totalStakeAdded, amount)

	// get the stake from the big sack
	_, err := pfe.executeTxStep(fmt.Sprintf(`
	{
		"step": "transfer",
		"txId": "%d",
		"tx": {
			"from": "%s",
			"to": "%s",
			"value": "%d"
		}
	}`,
		pfe.nextTxIndex(),
		pfe.faucetAddress,
		pfe.delegatorAddress(delegIndex),
		amount,
	))
	if err != nil {
		return err
	}

	// actual staking
	_, err = pfe.executeTxStep(fmt.Sprintf(`
	{
		"step": "scCall",
		"txId": "%d",
		"tx": {
			"from": "%s",
			"to": "%s",
			"value": "%d",
			"function": "stake",
			"arguments": [],
			"gasLimit": "100,000,000",
			"gasPrice": "0"
		}
	}`,
		pfe.nextTxIndex(),
		pfe.delegatorAddress(delegIndex),
		pfe.delegationContractAddress,
		amount,
	))
	pfe.log("stake, delegator: %d, amount: %d", delegIndex, amount)
	pfe.printUserStakeByType(delegIndex)
	return err
}

func (pfe *fuzzDelegationExecutor) unStake(delegatorIndex int, stake *big.Int) error {
	output, err := pfe.executeTxStep(fmt.Sprintf(`
	{
		"step": "scCall",
		"txId": "%d",
		"tx": {
			"from": "%s",
			"to": "%s",
			"value": "0",
			"function": "unStake",
			"arguments": ["%d"],
			"gasLimit": "500,000,000",
			"gasPrice": "0"
		}
	}`,
		pfe.nextTxIndex(),
		pfe.delegatorAddress(delegatorIndex),
		pfe.delegationContractAddress,
		stake,
	))
	if err != nil {
		return err
	}
	if output.ReturnCode == vmi.Ok {
		pfe.log("unStake, delegator: %d", delegatorIndex)
	} else if output.ReturnCode == vmi.OutOfGas {
		panic(fmt.Sprintf("unStake, delegator: %d, out of gas, message: %s", delegatorIndex, output.ReturnMessage))
	} else {
		pfe.log("unStake, delegator: %d, fail, %s", delegatorIndex, output.ReturnMessage)
	}

	pfe.printUserStakeByType(delegatorIndex)
	return nil
}

func (pfe *fuzzDelegationExecutor) unBond(delegatorIndex int) error {
	unbondable, err := pfe.getUserStakeOfType(delegatorIndex, UserUnbondable)
	if err != nil {
		return err
	}
	unbondedSum := big.NewInt(0)

	callIndex := 0
	for unbondable.Sign() > 0 {
		callIndex++
		output, err := pfe.executeTxStep(fmt.Sprintf(`
		{
			"step": "scCall",
			"txId": "%d",
			"tx": {
				"from": "%s",
				"to": "%s",
				"value": "0",
				"function": "unBond",
				"arguments": [],
				"gasLimit": "80,000,000",
				"gasPrice": "0"
			},
			"expect": {
				"out": [ "*" ],
				"status": "",
				"refund": "*"
			}
		}`,
			pfe.nextTxIndex(),
			pfe.delegatorAddress(delegatorIndex),
			pfe.delegationContractAddress,
		))
		if err != nil {
			return err
		}
		unbondedSum = unbondedSum.Add(unbondedSum, big.NewInt(0).SetBytes(output.ReturnData[0]))
		unbondable, err = pfe.getUserStakeOfType(delegatorIndex, UserUnbondable)
		if err != nil {
			return err
		}
		withdrawOnly, err := pfe.getUserStakeOfType(delegatorIndex, UserWithdrawOnly)
		if err != nil {
			return err
		}
		deferredPayment, err := pfe.getUserStakeOfType(delegatorIndex, UserDeferredPayment)
		if err != nil {
			return err
		}

		pfe.log("unBond, call #%d, delegator: %d, Still unbondable: %d, WithdrawOnly stake: %d. DeferredPayment: %d",
			callIndex, delegatorIndex, unbondable, withdrawOnly, deferredPayment)
	}

	pfe.log("stake withdrawn %d", unbondedSum)
	pfe.totalStakeWithdrawn.Add(pfe.totalStakeWithdrawn, unbondedSum)

	if unbondedSum.Cmp(big.NewInt(0)) > 0 {
		_, err = pfe.executeTxStep(fmt.Sprintf(`
		{
			"step": "transfer",
			"txId": "%d",
			"tx": {
				"from": "%s",
				"to": "%s",
				"value": "%d"
			}
		}`,
			pfe.nextTxIndex(),
			pfe.delegatorAddress(delegatorIndex),
			pfe.withdrawTargetAddress,
			unbondedSum,
		))
		if err != nil {
			return err
		}
	}

	pfe.printUserStakeByType(delegatorIndex)
	return nil
}

func (pfe *fuzzDelegationExecutor) getUserStakeOfType(delegatorIndex int, fundType string) (*big.Int, error) {
	output, err := pfe.executeTxStep(fmt.Sprintf(`
	{
		"step": "scCall",
		"txId": "%d",
		"tx": {
			"from": "%s",
			"to": "%s",
			"value": "0",
			"function": "%s",
			"arguments": ["%s"],
			"gasLimit": "100,000,000",
			"gasPrice": "0"
		}
	}`,
		pfe.nextTxIndex(),
		pfe.delegatorAddress(delegatorIndex),
		pfe.delegationContractAddress,
		fundType,
		pfe.delegatorAddress(delegatorIndex),
	))
	if err != nil {
		return nil, err
	}

	if len(output.ReturnData) != 0 {
		result := big.NewInt(0).SetBytes(output.ReturnData[0])
		return result, nil
	}

	return big.NewInt(0), nil
}

func (pfe *fuzzDelegationExecutor) printUserStakeByType(delegatorIndex int) {
	output, err := pfe.executeTxStep(fmt.Sprintf(`
	{
		"step": "scCall",
		"txId": "%d",
		"tx": {
			"from": "%s",
			"to": "%s",
			"value": "0",
			"function": "getUserStakeByType",
			"arguments": ["%s"],
			"gasLimit": "100,000,000",
			"gasPrice": "0"
		}
	}`,
		pfe.nextTxIndex(),
		pfe.ownerAddress,
		pfe.delegationContractAddress,
		pfe.delegatorAddress(delegatorIndex),
	))
	if err != nil {
		pfe.log("getUserStakeByType error")
		return
	}

	pfe.log("user %d stake by type:", delegatorIndex)
	pfe.printFundsInEachBucket(output.ReturnData)
}

func (pfe *fuzzDelegationExecutor) printTotalStakeByType() {
	output, err := pfe.executeTxStep(fmt.Sprintf(`
	{
		"step": "scCall",
		"txId": "%d",
		"tx": {
			"from": "%s",
			"to": "%s",
			"value": "0",
			"function": "getTotalStakeByType",
			"arguments": [],
			"gasLimit": "100,000,000",
			"gasPrice": "0"
		}
	}`,
		pfe.nextTxIndex(),
		pfe.ownerAddress,
		pfe.delegationContractAddress,
	))
	if err != nil {
		pfe.log("getTotalStakeByType error")
		return
	}

	pfe.log("total stake by type:")
	pfe.printFundsInEachBucket(output.ReturnData)
}

func (pfe *fuzzDelegationExecutor) printFundsInEachBucket(returnedData [][]byte) {
	if len(returnedData) == 5 {
		pfe.log("funds in contract: "+
			"WithdrawOnly: %d "+
			"Waiting: %d "+
			"Active: %d "+
			"UnStaked %d "+
			"DeferredPayment %d",
			big.NewInt(0).SetBytes(returnedData[0]),
			big.NewInt(0).SetBytes(returnedData[1]),
			big.NewInt(0).SetBytes(returnedData[2]),
			big.NewInt(0).SetBytes(returnedData[3]),
			big.NewInt(0).SetBytes(returnedData[4]),
		)
	}
}
