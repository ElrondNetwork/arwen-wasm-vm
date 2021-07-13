package elrond_ethereum_bridge

import (
	"math/big"
	"strconv"
)

func (fe *fuzzExecutor) proposeMultiTransferEsdtBatch(relayerAddress string,
	batchId int, transfers []*SimpleTransfer) (int, error) {

	args := []string{strconv.Itoa(batchId)}
	for _, transfer := range transfers {
		args = append(args, transfer.to, transfer.tokenId, transfer.amount.String())
	}

	output, err := fe.performSmartContractCall(
		relayerAddress,
		fe.data.actorAddresses.multisig,
		big.NewInt(0),
		"proposeMultiTransferEsdtBatch",
		args,
	)
	if err != nil {
		return 0, err
	}

	actionId := fe.bytesToInt(output[0])
	fe.data.multisigState.actions[actionId] = Action{
		actionType: ActionBatchTransferEsdtToken,
		data: BatchTransferEsdtTokenActionData{
			batchId:   batchId,
			transfers: transfers,
		},
	}
	fe.data.multisigState.signatures[actionId] = []string{relayerAddress}

	return actionId, nil
}
