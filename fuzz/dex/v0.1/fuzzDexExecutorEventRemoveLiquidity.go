package dex

import (
	"errors"
	"fmt"
	vmi "github.com/ElrondNetwork/elrond-vm-common"
	"math/big"
	"math/rand"
)

func (pfe *fuzzDexExecutor) removeLiquidity(r *rand.Rand, statistics *eventsStatistics) error {
	user := pfe.userAddress(r.Intn(pfe.numUsers) + 1)
	swapPair := pfe.swaps[r.Intn(len(pfe.swaps))]

	seed := r.Intn(pfe.removeLiquidityMaxValue) + 1
	amount := seed
	amountAmin := seed / 100
	amountBmin := seed / 100

	rawEquivalent, err := pfe.querySingleResultStringAddr(pfe.ownerAddress, swapPair.address,
		"getEquivalent", fmt.Sprintf("\"str:%s\", \"%d\"", swapPair.firstToken, 1000))
	if err != nil {
		return err
	}

	tokenABefore, err := pfe.getTokens(user, swapPair.firstToken)
	if err != nil {
		return err
	}
	tokenBBefore, err := pfe.getTokens(user, swapPair.secondToken)
	if err != nil {
		return err
	}
	tokenLpBefore, err := pfe.getTokens(user, swapPair.lpToken)
	if err != nil {
		return err
	}

	output, err := pfe.executeTxStep(fmt.Sprintf(`
	{
		"step": "scCall",
		"txId": "remove-liq",
		"tx": {
			"from": "%s",
			"to": "%s",
			"value": "0",
			"function": "removeLiquidity",
			"esdt": {
				"tokenIdentifier": "str:%s",
				"value": "%d"
			},
			"arguments": [
				"%d",
				"%d"
			],
			"gasLimit": "100,000,000",
			"gasPrice": "0"
		}
	}`,
		user,
		swapPair.address,
		swapPair.lpToken,
		amount,
		amountAmin,
		amountBmin,
	))
	if err != nil {
		return err
	}

	if output.ReturnCode == vmi.Ok {
		statistics.removeLiquidityHits += 1

		tokenAAfter, err := pfe.getTokens(user, swapPair.firstToken)
		if err != nil {
			return err
		}
		tokenBAfter, err := pfe.getTokens(user, swapPair.secondToken)
		if err != nil {
			return err
		}
		tokenLpAfter, err := pfe.getTokens(user, swapPair.lpToken)
		if err != nil {
			return err
		}

		rawOutput, err := pfe.querySingleResultStringAddr(pfe.ownerAddress, swapPair.address,
			"getEquivalent", fmt.Sprintf("\"str:%s\", \"%d\"", swapPair.firstToken, 1000))
		if err != nil {
			return err
		}

		if tokenABefore.Cmp(tokenAAfter) > -1 ||
			tokenBBefore.Cmp(tokenBAfter) > -1 ||
			tokenLpBefore.Cmp(tokenLpAfter) < 1 {
			return errors.New("FAILED remove liquidity balances on success")
		}

		before := big.NewInt(0).SetBytes(rawEquivalent[0])
		if big.NewInt(0).Cmp(before) != 0 {
			statistics.removeLiquidityPriceChecks += 1

			after := big.NewInt(0).SetBytes(rawOutput[0])
			difference := big.NewInt(0).Abs(big.NewInt(0).Sub(before, after))

			if big.NewInt(1).Cmp(difference) == -1 {
				return errors.New("PRICE CHANGED after add liquidity")
			}
		}
	} else {
		statistics.removeLiquidityMisses += 1

		pfe.log("remove liquidity %s -> %s", swapPair.firstToken, swapPair.secondToken)
		pfe.log("could not remove because %s", output.ReturnMessage)

		expectedErrors := map[string]bool{
			"Not enough LP token supply":    true,
			"Insufficient liquidity burned": true,
			"Not enough reserve":            true,
		}

		_, expected := expectedErrors[output.ReturnMessage]
		if !expected {
			return errors.New(output.ReturnMessage)
		}
	}

	return nil
}
