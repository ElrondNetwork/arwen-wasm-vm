package delegation

import (
	"errors"
	"fmt"
	"math/big"
)

func (pfe *fuzzDelegationExecutor) addRewards(amount *big.Int) error {
	// keep track of added rewards
	pfe.totalRewards.Add(pfe.totalRewards, amount)

	// simulate system rewards to delegation contract
	_, err := pfe.executeTxStep(fmt.Sprintf(`
	{
		"step": "validatorReward",
		"txId": "%d",
		"tx": {
			"to": "%s",
			"value": "%d"
		}
	}`,
		pfe.nextTxIndex(),
		pfe.delegationContractAddress,
		amount,
	))
	pfe.log("reward: %d", amount)
	return err
}

func (pfe *fuzzDelegationExecutor) claimRewards(delegatorIndex int) error {
	claimableRewardsBefore, err := pfe.getClaimableRewards(delegatorIndex)
	if err != nil {
		return err
	}
	if claimableRewardsBefore.Sign() == 0 {
		pfe.log("no rewards, delegator: %d", delegatorIndex)
		return nil
	}

	_, err = pfe.executeTxStep(fmt.Sprintf(`
	{
		"step": "scCall",
		"txId": "%d",
		"tx": {
			"from": "%s",
			"to": "%s",
			"value": "0",
			"function": "claimRewards",
			"arguments": [],
			"gasLimit": "500,000,000",
			"gasPrice": "0"
		},
		"expect": {
			"out": [],
			"status": "",
			"logs": "*",
			"gas": "*",
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

	claimableRewardsAfter, err := pfe.getClaimableRewards(delegatorIndex)
	if err != nil {
		return err
	}
	if claimableRewardsAfter.Sign() != 0 {
		return errors.New("getClaimableRewards should always yield 0 after claimRewards")
	}

	pfe.log("claim rewards, delegator: %d, amount: %d", delegatorIndex, claimableRewardsBefore)
	return nil
}

func (pfe *fuzzDelegationExecutor) getClaimableRewards(delegatorIndex int) (*big.Int, error) {
	output, err := pfe.executeTxStep(fmt.Sprintf(`
	{
		"step": "scQuery",
		"txId": "%d",
		"tx": {
			"to": "%s",
			"function": "getClaimableRewards",
			"arguments": [
				"%s"
			]
		},
		"expect": {
			"out": [ "*" ],
			"status": ""
		}
	}`,
		pfe.nextTxIndex(),
		pfe.delegationContractAddress,
		pfe.delegatorAddress(delegatorIndex),
	))
	if err != nil {
		return nil, err
	}
	return big.NewInt(0).SetBytes(output.ReturnData[0]), nil
}
