package dex

import (
	"errors"
	"fmt"
	"math/rand"

	vmi "github.com/ElrondNetwork/elrond-vm-common"
)

func (pfe *fuzzDexExecutor) enterFarm(r *rand.Rand, statistics *eventsStatistics) error {
	amountMax := r.Intn(pfe.tokenDepositMaxValue) + 1
	user := pfe.userAddress(r.Intn(pfe.numUsers) + 1)
	amount := int64(r.Intn(pfe.enterFarmMaxValue) + 1)
	farm := pfe.farms[r.Intn(len(pfe.farms))]

	stakersLen := len(pfe.farmers)
	if stakersLen == 0 || r.Intn(2) == 0 {
	} else {
		nonce := rand.Intn(stakersLen) + 1

		if pfe.farmers[nonce].value != 0 {
			user = pfe.farmers[nonce].user
			amount = pfe.farmers[nonce].value
			farm = pfe.farmers[nonce].farm

			depositAmount := int64(amountMax)
			if int64(amountMax) > amount {
				depositAmount = amount
				delete(pfe.farmers, nonce)
			} else {
				depositAmount = int64(amountMax)
				pfe.farmers[nonce] = FarmerInfo{
					value: amount - depositAmount,
					user:  user,
					farm:  farm,
				}
			}

			_, err := pfe.executeTxStep(fmt.Sprintf(`
	{
		"step": "scCall",
		"txId": "stake",
		"tx": {
			"from": "%s",
			"to": "%s",
			"value": "0",
			"function": "depositFarmToken",
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
				depositAmount,
				nonce,
			))
			if err != nil {
				return err
			}
		}
	}

	output, err := pfe.executeTxStep(fmt.Sprintf(`
	{
		"step": "scCall",
		"txId": "stake",
		"tx": {
			"from": "%s",
			"to": "%s",
			"value": "0",
			"function": "enterFarm",
			"esdt": {
				"tokenIdentifier": "str:%s",
				"value": "%d"
			},
			"arguments": [],
			"gasLimit": "100,000,000",
			"gasPrice": "0"
		}
	}`,
		user,
		farm.address,
		farm.farmingToken,
		amount,
	))
	if err != nil {
		return err
	}

	if output.ReturnCode == vmi.Ok {
		statistics.enterFarmHits += 1

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
		statistics.enterFarmMisses += 1

		pfe.log("stake %s", farm.farmingToken)
		pfe.log("could enter farm because %s", output.ReturnMessage)

		return errors.New(output.ReturnMessage)
	}

	return nil
}
