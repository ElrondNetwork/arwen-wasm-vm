package elrond_ethereum_bridge

import (
	"fmt"
	"math/big"
)

func (fe *fuzzExecutor) wrapEgld(userAddress string, amount *big.Int) error {
	egldBalanceBefore := fe.getBalance(userAddress)
	esdtBalanceBefore := fe.getEsdtBalance(userAddress, fe.data.wrappedEgldTokenId)

	_, err := fe.performSmartContractCall(
		userAddress,
		fe.data.actorAddresses.egldEsdtSwap,
		amount,
		"wrapEgld",
		[]string{},
	)
	if err != nil {
		return err
	}

	actualEgldBalanceAfter := fe.getBalance(userAddress)
	actualEsdtBalanceAfter := fe.getEsdtBalance(userAddress, fe.data.wrappedEgldTokenId)

	expectedEgldBalanceAfter := big.NewInt(0)
	expectedEgldBalanceAfter.Sub(egldBalanceBefore, amount)

	expectedEsdtBalanceAfter := big.NewInt(0)
	expectedEsdtBalanceAfter.Add(esdtBalanceBefore, amount)

	if actualEgldBalanceAfter.Cmp(expectedEgldBalanceAfter) != 0 {
		return fmt.Errorf("Wrong EGLD balance after wrapping, Expected: %s, Have: %s",
			expectedEgldBalanceAfter.String(),
			actualEgldBalanceAfter.String())
	}
	if actualEsdtBalanceAfter.Cmp(expectedEsdtBalanceAfter) != 0 {
		return fmt.Errorf("Wrong ESDT balance after wrapping, Expected: %s, Have: %s",
			expectedEsdtBalanceAfter.String(),
			actualEsdtBalanceAfter.String())
	}

	return nil
}

func (fe *fuzzExecutor) unwrapEgld(userAddress string, amount *big.Int) error {
	egldBalanceBefore := fe.getBalance(userAddress)
	esdtBalanceBefore := fe.getEsdtBalance(userAddress, fe.data.wrappedEgldTokenId)

	_, err := fe.performEsdtTransferSmartContractCall(
		userAddress,
		fe.data.actorAddresses.egldEsdtSwap,
		fe.data.wrappedEgldTokenId,
		amount,
		"unwrapEgld",
		[]string{},
	)
	if err != nil {
		return err
	}

	actualEgldBalanceAfter := fe.getBalance(userAddress)
	actualEsdtBalanceAfter := fe.getEsdtBalance(userAddress, fe.data.wrappedEgldTokenId)

	expectedEgldBalanceAfter := big.NewInt(0)
	expectedEgldBalanceAfter.Add(egldBalanceBefore, amount)

	expectedEsdtBalanceAfter := big.NewInt(0)
	expectedEsdtBalanceAfter.Sub(esdtBalanceBefore, amount)

	if actualEgldBalanceAfter.Cmp(expectedEgldBalanceAfter) != 0 {
		return fmt.Errorf("Wrong EGLD balance after unwrapping, Expected: %s, Have: %s",
			expectedEgldBalanceAfter.String(),
			actualEgldBalanceAfter.String())
	}
	if actualEsdtBalanceAfter.Cmp(expectedEsdtBalanceAfter) != 0 {
		return fmt.Errorf("Wrong ESDT balance after unwrapping, Expected: %s, Have: %s",
			expectedEsdtBalanceAfter.String(),
			actualEsdtBalanceAfter.String())
	}

	return nil
}
