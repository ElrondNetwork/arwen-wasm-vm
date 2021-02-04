package nodepart

import (
	"github.com/ElrondNetwork/arwen-wasm-vm/arwen"
	"github.com/ElrondNetwork/arwen-wasm-vm/ipc/common"
)

func (part *NodePart) replyToBlockchainNewAddress(request common.MessageHandler) common.MessageHandler {
	typedRequest := request.(*common.MessageBlockchainNewAddressRequest)
	result, err := part.blockchain.NewAddress(typedRequest.CreatorAddress, typedRequest.CreatorNonce, typedRequest.VMType)
	response := common.NewMessageBlockchainNewAddressResponse(result, err)
	return response
}

func (part *NodePart) replyToBlockchainGetStorageData(request common.MessageHandler) common.MessageHandler {
	typedRequest := request.(*common.MessageBlockchainGetStorageDataRequest)
	data, err := part.blockchain.GetStorageData(typedRequest.Address, typedRequest.Index)
	response := common.NewMessageBlockchainGetStorageDataResponse(data, err)
	return response
}

func (part *NodePart) replyToBlockchainGetBlockhash(request common.MessageHandler) common.MessageHandler {
	typedRequest := request.(*common.MessageBlockchainGetBlockhashRequest)
	result, err := part.blockchain.GetBlockhash(typedRequest.Nonce)
	response := common.NewMessageBlockchainGetBlockhashResponse(result, err)
	return response
}

func (part *NodePart) replyToBlockchainLastNonce(_ common.MessageHandler) common.MessageHandler {
	result := part.blockchain.LastNonce()
	response := common.NewMessageBlockchainLastNonceResponse(result)
	return response
}

func (part *NodePart) replyToBlockchainLastRound(_ common.MessageHandler) common.MessageHandler {
	result := part.blockchain.LastRound()
	response := common.NewMessageBlockchainLastRoundResponse(result)
	return response
}

func (part *NodePart) replyToBlockchainLastTimeStamp(_ common.MessageHandler) common.MessageHandler {
	result := part.blockchain.LastTimeStamp()
	response := common.NewMessageBlockchainLastTimeStampResponse(result)
	return response
}

func (part *NodePart) replyToBlockchainLastRandomSeed(_ common.MessageHandler) common.MessageHandler {
	result := part.blockchain.LastRandomSeed()
	response := common.NewMessageBlockchainLastRandomSeedResponse(result)
	return response
}

func (part *NodePart) replyToBlockchainLastEpoch(_ common.MessageHandler) common.MessageHandler {
	result := part.blockchain.LastEpoch()
	response := common.NewMessageBlockchainLastEpochResponse(result)
	return response
}

func (part *NodePart) replyToBlockchainGetStateRootHash(_ common.MessageHandler) common.MessageHandler {
	result := part.blockchain.GetStateRootHash()
	response := common.NewMessageBlockchainGetStateRootHashResponse(result)
	return response
}

func (part *NodePart) replyToBlockchainCurrentNonce(_ common.MessageHandler) common.MessageHandler {
	result := part.blockchain.CurrentNonce()
	response := common.NewMessageBlockchainCurrentNonceResponse(result)
	return response
}

func (part *NodePart) replyToBlockchainCurrentRound(_ common.MessageHandler) common.MessageHandler {
	result := part.blockchain.CurrentRound()
	response := common.NewMessageBlockchainCurrentRoundResponse(result)
	return response
}

func (part *NodePart) replyToBlockchainCurrentTimeStamp(_ common.MessageHandler) common.MessageHandler {
	result := part.blockchain.CurrentTimeStamp()
	response := common.NewMessageBlockchainCurrentTimeStampResponse(result)
	return response
}

func (part *NodePart) replyToBlockchainCurrentRandomSeed(_ common.MessageHandler) common.MessageHandler {
	result := part.blockchain.CurrentRandomSeed()
	response := common.NewMessageBlockchainCurrentRandomSeedResponse(result)
	return response
}

func (part *NodePart) replyToBlockchainCurrentEpoch(_ common.MessageHandler) common.MessageHandler {
	result := part.blockchain.CurrentEpoch()
	response := common.NewMessageBlockchainCurrentEpochResponse(result)
	return response
}

func (part *NodePart) replyToBlockchainProcessBuiltinFunction(request common.MessageHandler) common.MessageHandler {
	typedRequest := request.(*common.MessageBlockchainProcessBuiltinFunctionRequest)
	vmOutput, err := part.blockchain.ProcessBuiltInFunction(&typedRequest.CallInput)
	response := common.NewMessageBlockchainProcessBuiltinFunctionResponse(vmOutput, err)
	return response
}

func (part *NodePart) replyToBlockchainGetBuiltinFunctionNames(_ common.MessageHandler) common.MessageHandler {
	functionNames := part.blockchain.GetBuiltinFunctionNames()
	response := common.NewMessageBlockchainGetBuiltinFunctionNamesResponse(functionNames)
	return response
}

func (part *NodePart) replyToBlockchainGetAllState(request common.MessageHandler) common.MessageHandler {
	typedRequest := request.(*common.MessageBlockchainGetAllStateRequest)
	state, err := part.blockchain.GetAllState(typedRequest.Address)
	response := common.NewMessageBlockchainGetAllStateResponse(state, err)
	return response
}

func (part *NodePart) replyToBlockchainGetUserAccount(request common.MessageHandler) common.MessageHandler {
	typedRequest := request.(*common.MessageBlockchainGetUserAccountRequest)
	account, err := part.blockchain.GetUserAccount(typedRequest.Address)

	if err != nil || arwen.IfNil(account) {
		return common.NewMessageBlockchainGetUserAccountResponse(nil, err)
	}

	return common.NewMessageBlockchainGetUserAccountResponse(&common.Account{
		Nonce:           account.GetNonce(),
		Address:         account.AddressBytes(),
		Balance:         account.GetBalance(),
		CodeMetadata:    account.GetCodeMetadata(),
		CodeHash:        account.GetCodeHash(),
		RootHash:        account.GetRootHash(),
		DeveloperReward: account.GetDeveloperReward(),
		OwnerAddress:    account.GetOwnerAddress(),
		UserName:        account.GetUserName(),
	}, err)
}

func (part *NodePart) replyToBlockchainGetCode(request common.MessageHandler) common.MessageHandler {
	typedRequest := request.(*common.MessageBlockchainGetCodeRequest)
	code := part.blockchain.GetCode(typedRequest.Account)

	return common.NewMessageBlockchainGetCodeResponse(code)
}

func (part *NodePart) replyToBlockchainGetShardOfAddress(request common.MessageHandler) common.MessageHandler {
	typedRequest := request.(*common.MessageBlockchainGetShardOfAddressRequest)
	result := part.blockchain.GetShardOfAddress(typedRequest.Address)
	response := common.NewMessageBlockchainGetShardOfAddressResponse(result)
	return response
}

func (part *NodePart) replyToBlockchainIsSmartContract(request common.MessageHandler) common.MessageHandler {
	typedRequest := request.(*common.MessageBlockchainIsSmartContractRequest)
	result := part.blockchain.IsSmartContract(typedRequest.Address)
	response := common.NewMessageBlockchainIsSmartContractResponse(result)
	return response
}

func (part *NodePart) replyToBlockchainIsPayable(request common.MessageHandler) common.MessageHandler {
	typedRequest := request.(*common.MessageBlockchainIsPayableRequest)
	result, err := part.blockchain.IsPayable(typedRequest.Address)
	response := common.NewMessageBlockchainIsPayableResponse(result, err)
	return response
}

func (part *NodePart) replyToBlockchainSaveCompiledCode(request common.MessageHandler) common.MessageHandler {
	typedRequest := request.(*common.MessageBlockchainSaveCompiledCodeRequest)
	part.blockchain.SaveCompiledCode(typedRequest.CodeHash, typedRequest.Code)
	response := common.NewMessageBlockchainSaveCompiledCodeResponse()
	return response
}

func (part *NodePart) replyToBlockchainGetCompiledCode(request common.MessageHandler) common.MessageHandler {
	typedRequest := request.(*common.MessageBlockchainGetCompiledCodeRequest)
	found, code := part.blockchain.GetCompiledCode(typedRequest.CodeHash)
	response := common.NewMessageBlockchainGetCompiledCodeResponse(found, code)
	return response
}
