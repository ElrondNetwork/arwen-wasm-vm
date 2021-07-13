package elrond_ethereum_bridge

import (
	"math/big"
	"strconv"
)

// Note: proposals that affect a specific child contract will not be here
// they will be in the child SC's file

func (fe *fuzzExecutor) sign(relayerAddress string, actionId int) error {
	_, err := fe.performSmartContractCall(
		relayerAddress,
		fe.data.actorAddresses.multisig,
		big.NewInt(0),
		"sign",
		[]string{strconv.Itoa(actionId)},
	)
	if err != nil {
		return err
	}

	if !fe.hasSignedAlready(relayerAddress, actionId) {
		fe.data.multisigState.signatures[actionId] = append(fe.data.multisigState.signatures[actionId], relayerAddress)
	}

	return nil
}

func (fe *fuzzExecutor) performAction(relayerAddress string, actionId int) ([][]byte, error) {
	output, err := fe.performSmartContractCall(
		relayerAddress,
		fe.data.actorAddresses.multisig,
		big.NewInt(0),
		"performAction",
		[]string{strconv.Itoa(actionId)},
	)
	if err != nil {
		return nil, err
	}

	delete(fe.data.multisigState.actions, actionId)
	delete(fe.data.multisigState.signatures, actionId)

	return output, nil
}

func (fe *fuzzExecutor) hasSignedAlready(relayerAddress string, actionId int) bool {
	for _, signer := range fe.data.multisigState.signatures[actionId] {
		if signer == relayerAddress {
			return true
		}
	}

	return false
}

func (fe *fuzzExecutor) getActionIdForSetCurrentTransactionBatchStatus(
	esdtSafeBatchId int, statuses ...TransactionStatus) (int, error) {

	args := []string{strconv.Itoa(esdtSafeBatchId)}
	for _, status := range statuses {
		args = append(args, strconv.Itoa(int(status)))
	}

	output, err := fe.performSmartContractCall(
		fe.getRandomRelayer(),
		fe.data.actorAddresses.multisig,
		big.NewInt(0),
		"getActionIdForSetCurrentTransactionBatchStatus",
		args,
	)
	if err != nil {
		return 0, err
	}

	actionId := fe.bytesToInt(output[0])

	return actionId, nil
}

func (fe *fuzzExecutor) getActionIdForTransferBatch(batchId int, transfers []*SimpleTransfer) (int, error) {
	args := []string{strconv.Itoa(batchId)}
	for _, transfer := range transfers {
		args = append(args, transfer.to, transfer.tokenId, transfer.amount.String())
	}

	output, err := fe.performSmartContractCall(
		fe.getRandomRelayer(),
		fe.data.actorAddresses.multisig,
		big.NewInt(0),
		"getActionIdForTransferBatch",
		args,
	)
	if err != nil {
		return 0, err
	}

	actionId := fe.bytesToInt(output[0])

	return actionId, nil
}
