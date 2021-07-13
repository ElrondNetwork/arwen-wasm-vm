package elrond_ethereum_bridge

type MultisigAction int

const (
	ActionNone MultisigAction = iota
	ActionAddBoardMember
	ActionAddProposer
	ActionRemoveUser
	ActionSlashUser
	ActionChangeQuorum
	ActionSetCurrentTransactionBatchStatus
	ActionBatchTransferEsdtToken
)

type Action struct {
	actionType MultisigAction
	data       interface{}
}

type SetCurrentTransactionBatchStatusActionData struct {
	relayerRewardAddress   string
	esdtSafeBatchId        int
	transactionBatchStatus []TransactionStatus
}

type BatchTransferEsdtTokenActionData struct {
	batchId   int
	transfers []*SimpleTransfer
}
