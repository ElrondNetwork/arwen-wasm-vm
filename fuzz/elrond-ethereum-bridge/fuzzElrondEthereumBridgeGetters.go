package elrond_ethereum_bridge

import (
	"encoding/binary"
	"fmt"
	"math/big"
	"strconv"
)

func (fe *fuzzExecutor) interpretExpr(expression string) []byte {
	bytes, err := fe.mandosParser.ExprInterpreter.InterpretString(expression)
	if err != nil {
		panic(err)
	}
	return bytes
}

func (fe *fuzzExecutor) relayerAddress(index int) string {
	return fmt.Sprintf("address:relayer-%05d", index)
}

func (fe *fuzzExecutor) userAddress(index int) string {
	return fmt.Sprintf("address:user-%05d", index)
}

func (fe *fuzzExecutor) getNonce(address string) int {
	acct := fe.world.AcctMap.GetAccount(fe.interpretExpr(address))
	return int(acct.Nonce)
}

func (fe *fuzzExecutor) getBalance(address string) *big.Int {
	acct := fe.world.AcctMap.GetAccount(fe.interpretExpr(address))
	return acct.Balance
}

func (fe *fuzzExecutor) getEsdtBalance(address string, tokenId string) *big.Int {
	acct := fe.world.AcctMap.GetAccount(fe.interpretExpr(address))
	balance, err := acct.GetTokenBalanceByName(string(fe.interpretExpr(tokenId)))

	if err != nil {
		return big.NewInt(0)
	}
	return balance
}

func (fe *fuzzExecutor) generateValidRandomEsdtPayment(address string) (string, *big.Int, error) {
	acc := fe.world.AcctMap.GetAccount(fe.interpretExpr(address))
	allEsdts, err := acc.GetFullMockESDTData()
	if err != nil {
		return "", nil, err
	}

	// map key order is not guaranteed, so this is "random"
	for tokenId := range allEsdts {
		tokenIdMandosFormat := "str:" + tokenId
		balance := fe.getEsdtBalance(address, tokenIdMandosFormat)
		amount := fe.getRandomBigInt(balance)

		return tokenIdMandosFormat, amount, nil
	}

	return "", nil, fmt.Errorf("Account has no ESDT")
}

func (fe *fuzzExecutor) generateValidBridgedEsdtPayment() *SimpleTransfer {
	destAddress := fe.getRandomAccount()
	tokenId := fe.getRandomTokenIdFromWhitelist()
	amount := big.NewInt(fe.randSource.Int63n(1000) + 1)

	return &SimpleTransfer{
		to:      destAddress,
		tokenId: tokenId,
		amount:  amount,
	}
}

func (fe *fuzzExecutor) nextTxIndex() int {
	fe.txIndex++
	return fe.txIndex
}

func (fe *fuzzExecutor) nextEthereumBatchId() int {
	fe.ethereumBatchId++
	return fe.ethereumBatchId
}

// Do NOT use with bigInts that don't fit [0, max(int64))
// maxValue is inclusive
func (fe *fuzzExecutor) getRandomBigInt(maxValue *big.Int) *big.Int {
	return big.NewInt(int64(fe.randSource.Intn(int(maxValue.Int64())) + 1))
}

func (fe *fuzzExecutor) getRandomUser() string {
	index := fe.randSource.Intn(len(fe.data.actorAddresses.users))
	return fe.data.actorAddresses.users[index]
}

func (fe *fuzzExecutor) getRandomRelayer() string {
	index := fe.randSource.Intn(len(fe.data.actorAddresses.relayers))
	return fe.data.actorAddresses.relayers[index]
}

// 75% user account, 25% SC account
func (fe *fuzzExecutor) getRandomAccount() string {
	randNr := fe.randSource.Float32()
	if randNr < 0.75 {
		return fe.getRandomUser()
	} else {
		// doesn't matter which SC address is returned
		return fe.data.actorAddresses.esdtSafe
	}
}

func (fe *fuzzExecutor) getRandomTokenIdFromWhitelist() string {
	index := fe.randSource.Intn(len(fe.data.tokenWhitelist))
	return fe.data.tokenWhitelist[index]
}

func (fe *fuzzExecutor) getEthAddress() string {
	return "0x0102030405060708091011121314151617181920"
}

func (fe *fuzzExecutor) bytesToInt(bytes []byte) int {
	intBytes := bytes
	for len(intBytes) != 8 {
		intBytes = append([]byte{0}, intBytes...)
	}

	return int(binary.BigEndian.Uint64(intBytes))
}

func (fe *fuzzExecutor) GetExpectedBalancesAfterBridgeTransferToEthereum(
	transactions []*Transaction, statuses []TransactionStatus,
) (map[string]map[string]*big.Int, error) {

	if len(transactions) != len(statuses) {
		return nil, fmt.Errorf("Transactions and status array lengths do not match")
	}

	esdtSafeAddress := fe.data.actorAddresses.esdtSafe

	balances := make(map[string]map[string]*big.Int)
	balances[esdtSafeAddress] = make(map[string]*big.Int)

	for i, tx := range transactions {
		status := statuses[i]
		var userBalance *big.Int
		var esdtSafeBalance *big.Int

		// initialize each sub-map
		if _, ok := balances[tx.from]; !ok {
			balances[tx.from] = make(map[string]*big.Int)
		}

		// take entry from map if it already exists
		// this ensures consistency in cases where a single account performs multiple transfers
		if balance, ok := balances[tx.from][tx.tokenId]; ok {
			userBalance = balance
		} else {
			userBalance = fe.getEsdtBalance(tx.from, tx.tokenId)
		}

		if balance, ok := balances[esdtSafeAddress][tx.tokenId]; ok {
			esdtSafeBalance = balance
		} else {
			esdtSafeBalance = fe.getEsdtBalance(esdtSafeAddress, tx.tokenId)
		}

		// EsdtSafe SC either burns or returns the tokens, so balance decreases in both cases
		esdtSafeNewBalance := big.NewInt(0)
		esdtSafeNewBalance.Sub(esdtSafeBalance, tx.amount)
		balances[esdtSafeAddress][tx.tokenId] = esdtSafeNewBalance

		switch status {
		case Executed:
			// tokens are burned, so no change to user's balance
			balances[tx.from][tx.tokenId] = userBalance
		case Rejected:
			// tokens are returned to the user
			newUserBalance := big.NewInt(0)
			newUserBalance.Add(userBalance, tx.amount)

			balances[tx.from][tx.tokenId] = newUserBalance
		default:
			return nil, fmt.Errorf("Invalid status provided: %s", strconv.Itoa(int(status)))
		}
	}

	return balances, nil
}

// statuses are only available after the action is executed
// so we assume all of them executed successfuly and let the caller handle the error cases
func (fe *fuzzExecutor) GetExpectedBalancesAfterBridgeTransferToElrond(transfers []*SimpleTransfer,
) map[string]map[string]*big.Int {

	balances := make(map[string]map[string]*big.Int)

	for _, transf := range transfers {
		var userBalance *big.Int

		// initialize each sub-map
		if _, ok := balances[transf.to]; !ok {
			balances[transf.to] = make(map[string]*big.Int)
		}

		// take entry from map if it already exists
		// this ensures consistency in cases where a single account is the destination for multiple transfers
		if balance, ok := balances[transf.to][transf.tokenId]; ok {
			userBalance = balance
		} else {
			userBalance = fe.getEsdtBalance(transf.to, transf.tokenId)
		}

		newUserBalance := big.NewInt(0)
		newUserBalance.Add(userBalance, transf.amount)

		balances[transf.to][transf.tokenId] = newUserBalance
	}

	return balances
}
