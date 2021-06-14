package delegation

import (
	"fmt"
	"math/big"
)

func (pfe *fuzzDelegationExecutor) interpretExpr(expression string) []byte {
	bytes, err := pfe.mandosParser.ExprInterpreter.InterpretString(expression)
	if err != nil {
		panic(err)
	}
	return bytes
}

func (pfe *fuzzDelegationExecutor) delegatorAddress(delegIndex int) string {
	if delegIndex == 0 {
		return pfe.ownerAddress
	}

	return fmt.Sprintf("address:delegator-%05d", delegIndex)
}

func (pfe *fuzzDelegationExecutor) getAuctionBalance() *big.Int {
	// pfe.mandosParser.ExprInterpreter.InterpretString()
	acct := pfe.world.AcctMap.GetAccount(pfe.interpretExpr(pfe.auctionMockAddress))
	return acct.Balance
}

func (pfe *fuzzDelegationExecutor) getWithdrawTargetBalance() *big.Int {
	acct := pfe.world.AcctMap.GetAccount(pfe.interpretExpr(pfe.withdrawTargetAddress))
	return acct.Balance
}

func (pfe *fuzzDelegationExecutor) getContractBalance() *big.Int {
	acct := pfe.world.AcctMap.GetAccount(pfe.interpretExpr(pfe.delegationContractAddress))
	return acct.Balance
}
