package dex

import (
	"errors"
	"fmt"
	worldmock "github.com/ElrondNetwork/arwen-wasm-vm/mock/world"
	"math/big"
)

func (pfe *fuzzDexExecutor) getTokensWithNonce(address []byte, toktik string, nonce int) (*big.Int, error) {
	token := worldmock.MakeTokenKey([]byte(toktik), uint64(nonce))
	return pfe.world.BuiltinFuncs.GetTokenBalance(address, token)
}

func (pfe *fuzzDexExecutor) getTokens(address []byte, toktik string) (*big.Int, error) {
	token := worldmock.MakeTokenKey([]byte(toktik), 0)
	return pfe.world.BuiltinFuncs.GetTokenBalance(address, token)
}

func (pfe *fuzzDexExecutor) checkTokens() error {
	expectedSumString := fmt.Sprintf("%d", pfe.numUsers)
	expectedSumString += "000000000000000000000000000000"

	sum, err := pfe.getSumForToken(pfe.wegldTokenId)
	if err != nil {
		return err
	}
	if sum != expectedSumString {
		return errors.New("sum differs")
	}

	for i := 1; i <= pfe.numTokens; i++ {
		sum, err = pfe.getSumForToken(pfe.tokenTicker(i))
		if err != nil {
			return err
		}
		if sum != expectedSumString {
			return errors.New("sum differs")
		}
	}

	return nil
}

func (pfe* fuzzDexExecutor) getSumForToken(tokenTicker string) (string, error) {
	totalSum := big.NewInt(0)

	for i := 1; i < pfe.numTokens; i++ {
		for j := i + 1; j <= pfe.numTokens; j++ {
			tokenA := pfe.tokenTicker(i)
			tokenB := pfe.tokenTicker(j)

			rawResponse, err := pfe.querySingleResult(pfe.ownerAddress, pfe.routerAddress,
				"getPair", fmt.Sprintf("\"str:%s\", \"str:%s\"", tokenA, tokenB))
			if err != nil {
				return "", err
			}

			result, err := pfe.getTokens(rawResponse[0], tokenTicker)
			if err != nil {
				return "", err
			}

			totalSum = big.NewInt(0).Add(totalSum, result)
		}
	}
	for i := 1; i <= pfe.numTokens; i++ {
		tokenA := pfe.wegldTokenId
		tokenB := pfe.tokenTicker(i)

		rawResponse, err := pfe.querySingleResult(pfe.ownerAddress, pfe.routerAddress,
			"getPair", fmt.Sprintf("\"str:%s\", \"str:%s\"", tokenA, tokenB))
		if err != nil {
			return "", err
		}

		result, err := pfe.getTokens(rawResponse[0], tokenTicker)
		if err != nil {
			return "", err
		}

		totalSum = big.NewInt(0).Add(totalSum, result)
	}
	for i := 1; i <= pfe.numTokens; i++ {
		tokenA := pfe.mexTokenId
		tokenB := pfe.tokenTicker(i)

		rawResponse, err := pfe.querySingleResult(pfe.ownerAddress, pfe.routerAddress,
			"getPair", fmt.Sprintf("\"str:%s\", \"str:%s\"", tokenA, tokenB))
		if err != nil {
			return "", err
		}

		result, err := pfe.getTokens(rawResponse[0], tokenTicker)
		if err != nil {
			return "", err
		}

		totalSum = big.NewInt(0).Add(totalSum, result)
	}
	tokenA := pfe.wegldTokenId
	tokenB := pfe.mexTokenId

	rawResponse, err := pfe.querySingleResult(pfe.ownerAddress, pfe.routerAddress,
		"getPair", fmt.Sprintf("\"str:%s\", \"str:%s\"", tokenA, tokenB))
	if err != nil {
		return "", err
	}

	result, err := pfe.getTokens(rawResponse[0], tokenTicker)
	if err != nil {
		return "", err
	}

	totalSum = big.NewInt(0).Add(totalSum, result)

	for i := 1; i <= pfe.numUsers; i++ {
		user := pfe.userAddress(i)
		result, err := pfe.getTokens(user, tokenTicker)
		if err != nil {
			return "", err
		}

		totalSum = big.NewInt(0).Add(totalSum, result)
	}

	//STAKING
	result, err = pfe.getTokens(pfe.wegldStakingAddress, tokenTicker)
	if err != nil {
		return "", err
	}

	totalSum = big.NewInt(0).Add(totalSum, result)
	totalSumString := totalSum.String()

	result, err = pfe.getTokens(pfe.mexStakingAddress, tokenTicker)
	if err != nil {
		return "", err
	}

	totalSum = big.NewInt(0).Add(totalSum, result)
	totalSumString = totalSum.String()

	return totalSumString, nil
}