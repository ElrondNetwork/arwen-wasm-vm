package elrond_ethereum_bridge

import "math/big"

type TransactionStatus int

const (
	None TransactionStatus = iota
	Pending
	InProgress
	Executed
	Rejected
)

type ActorAddresses struct {
	owner             string
	relayers          []string
	users             []string
	multisig          string
	priceAggregator   string
	egldEsdtSwap      string
	esdtSafe          string
	ethereumFeePrepay string
	multiTransferEsdt string
}

type Transaction struct {
	from    string
	tokenId string
	amount  *big.Int
}

type SimpleTransfer struct {
	to      string
	tokenId string
	amount  *big.Int
}

type MultisigState struct {
	requiredStake *big.Int
	quorum        int
	actions       map[int]Action   // action ID -> action data
	signatures    map[int][]string // action ID -> signer address list

	allEsdtSafeTransactions []*Transaction

	currentEsdtSafeBatchId          int
	currentEsdtSafeTransactionBatch []*Transaction

	currentEthereumBatchId int
}
