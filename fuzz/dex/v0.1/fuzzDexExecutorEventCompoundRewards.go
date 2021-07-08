package dex

import (
	"errors"
	"fmt"
	"math/rand"

	vmi "github.com/ElrondNetwork/elrond-vm-common"
)

func (pfe *fuzzDexExecutor) compoundRewards(r *rand.Rand, statistics *eventsStatistics) error {
	amountMax := r.Intn(pfe.compoundRewardsMaxValue) + 1

	stakersLen := len(pfe.farmers)
	if stakersLen == 0 {
		return nil
	}

	nonce := rand.Intn(stakersLen) + 1
	user := pfe.farmers[nonce].user
	amount := pfe.farmers[nonce].value
	if pfe.farmers[nonce].value == 0 {
		return nil
	}

	farm := pfe.farmers[nonce].farm
	claimAmount := int64(amountMax)
	if int64(amountMax) > amount {
		claimAmount = amount
		delete(pfe.farmers, nonce)
	} else {
		claimAmount = int64(amountMax)
		pfe.farmers[nonce] = FarmerInfo{
			value: amount - claimAmount,
			user:  user,
			farm:  farm,
		}
	}

	if farm.farmingToken != farm.rewardToken {
		return nil
	}

	output, err := pfe.executeTxStep(fmt.Sprintf(`
	{
		"step": "scCall",
		"txId": "claimRewards",
		"tx": {
			"from": "%s",
			"to": "%s",
			"value": "0",
			"function": "compoundRewards",
			"esdt": {
				"tokenIdentifier": "str:%s",
				"value": "%d",
				"nonce": "%d"
			},
			"arguments": [],
			"gasLimit": "100,000,000",
			"gasPrice": "0"
		}
	}`,
		user,
		farm.address,
		farm.farmToken,
		claimAmount,
		nonce,
	))
	if err != nil {
		return err
	}

	if output.ReturnCode == vmi.Ok {
		statistics.compoundRewardsHits += 1

		pfe.currentFarmTokenNonce[farm.address] += 1
		nonce := pfe.currentFarmTokenNonce[farm.address]
		bigint, errGet := pfe.getTokensWithNonce(user, farm.farmToken, nonce)
		if errGet != nil {
			return errGet
		}

		pfe.farmers[nonce] = FarmerInfo{
			user:  user,
			value: bigint.Int64(),
			farm:  farm,
		}
	} else {
		statistics.compoundRewardsMisses += 1

		expectedErrors := map[string]bool{
			"Farming token differ from reward token": true,
			"Farming token amount is zero":           true,
		}

		_, expected := expectedErrors[output.ReturnMessage]
		if !expected {
			return errors.New(output.ReturnMessage)
		}
	}

	return nil
}
