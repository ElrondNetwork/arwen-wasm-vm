package dex

import (
	"errors"
	"fmt"
	"math/big"
	"math/rand"
	"strconv"

	vmi "github.com/ElrondNetwork/elrond-vm-common"
)

func (pfe *fuzzDexExecutor) claimRewards(r *rand.Rand, statistics *eventsStatistics) error {
	amountMax := r.Intn(pfe.claimRewardsMaxValue) + 1

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

	mexBefore, err := pfe.getTokens(user, pfe.mexTokenId)
	if err != nil {
		return err
	}

	tokenData, err := pfe.getTokenData(user, farm.farmToken, nonce)
	if err != nil {
		return err
	}

	sizeByte := make([]byte, 4)
	copy(sizeByte, tokenData.TokenMetaData.Attributes[0:4])
	size, _ := strconv.Atoi(string(sizeByte))
	rps := make([]byte, size)
	copy(rps, tokenData.TokenMetaData.Attributes[4:(4+size)])

	output, err := pfe.executeTxStep(fmt.Sprintf(`
	{
		"step": "scCall",
		"txId": "claimRewards",
		"tx": {
			"from": "%s",
			"to": "%s",
			"value": "0",
			"function": "claimRewards",
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
		statistics.claimRewardsHits += 1

		mexAfter, err := pfe.getTokens(user, pfe.mexTokenId)
		if err != nil {
			return err
		}

		rpsAfter, err := pfe.querySingleResult(user, farm.address, "getRewardPerShare", "")
		if err != nil {
			return err
		}

		if mexAfter.Cmp(mexBefore) == 1 {
			statistics.claimRewardsWithRewards += 1

			rpsAttrs := big.NewInt(0).SetBytes(rps)
			rpsDifference := big.NewInt(0).Sub(big.NewInt(0).SetBytes(rpsAfter[0]), rpsAttrs)
			shouldEarn := big.NewInt(0).Div(big.NewInt(0).Mul(big.NewInt(claimAmount), rpsDifference), big.NewInt(1000000000000))
			earned := big.NewInt(0).Sub(mexAfter, mexBefore)

			if rpsAttrs.Cmp(big.NewInt(0)) != 0 && earned.Cmp(shouldEarn) != 0 {
				return errors.New("reward is not good")
			}

		} else if mexAfter.Cmp(mexBefore) == -1 {
			return errors.New("LOST mex while claimRewards")
		}

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
		statistics.claimRewardsMisses += 1

		expectedErrors := map[string]bool{
			"Farming token amount is zero": true,
		}

		_, expected := expectedErrors[output.ReturnMessage]
		if !expected {
			return errors.New(output.ReturnMessage)
		}
	}

	return nil
}
