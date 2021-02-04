package arwenpart

import (
	"github.com/ElrondNetwork/arwen-wasm-vm/ipc/common"
	"github.com/ElrondNetwork/elrond-go/core/vmcommon"
)

var _ vmcommon.BlockchainHook = (*BlockchainHookGateway)(nil)

// BlockchainHookGateway forwards requests to the actual hook
type BlockchainHookGateway struct {
	messenger *ArwenMessenger
}

// NewBlockchainHookGateway creates a new gateway
func NewBlockchainHookGateway(messenger *ArwenMessenger) *BlockchainHookGateway {
	return &BlockchainHookGateway{messenger: messenger}
}

// NewAddress forwards a message to the actual hook
func (blockchain *BlockchainHookGateway) NewAddress(creatorAddress []byte, creatorNonce uint64, vmType []byte) ([]byte, error) {
	request := common.NewMessageBlockchainNewAddressRequest(creatorAddress, creatorNonce, vmType)
	rawResponse, err := blockchain.messenger.SendHookCallRequest(request)
	if err != nil {
		return nil, err
	}

	if rawResponse.GetKind() != common.BlockchainNewAddressResponse {
		return nil, common.ErrBadHookResponseFromNode
	}

	response := rawResponse.(*common.MessageBlockchainNewAddressResponse)
	return response.Result, response.GetError()
}

// GetStorageData forwards a message to the actual hook
func (blockchain *BlockchainHookGateway) GetStorageData(address []byte, index []byte) ([]byte, error) {
	request := common.NewMessageBlockchainGetStorageDataRequest(address, index)
	rawResponse, err := blockchain.messenger.SendHookCallRequest(request)
	if err != nil {
		return nil, err
	}

	if rawResponse.GetKind() != common.BlockchainGetStorageDataResponse {
		return nil, common.ErrBadHookResponseFromNode
	}

	response := rawResponse.(*common.MessageBlockchainGetStorageDataResponse)
	return response.Data, response.GetError()
}

// GetBlockhash forwards a message to the actual hook
func (blockchain *BlockchainHookGateway) GetBlockhash(nonce uint64) ([]byte, error) {
	request := common.NewMessageBlockchainGetBlockhashRequest(nonce)
	rawResponse, err := blockchain.messenger.SendHookCallRequest(request)
	if err != nil {
		return nil, err
	}

	if rawResponse.GetKind() != common.BlockchainGetBlockhashResponse {
		log.Error("GetBlockhash", "err", common.ErrBadHookResponseFromNode)
		return nil, common.ErrBadHookResponseFromNode
	}

	response := rawResponse.(*common.MessageBlockchainGetBlockhashResponse)
	return response.Result, response.GetError()
}

// LastNonce forwards a message to the actual hook
func (blockchain *BlockchainHookGateway) LastNonce() uint64 {
	request := common.NewMessageBlockchainLastNonceRequest()
	rawResponse, err := blockchain.messenger.SendHookCallRequest(request)
	if err != nil {
		return 0
	}

	if rawResponse.GetKind() != common.BlockchainLastNonceResponse {
		log.Error("LastNonce", "err", common.ErrBadHookResponseFromNode)
		return 0
	}

	response := rawResponse.(*common.MessageBlockchainLastNonceResponse)
	return response.Result
}

// LastRound forwards a message to the actual hook
func (blockchain *BlockchainHookGateway) LastRound() uint64 {
	request := common.NewMessageBlockchainLastRoundRequest()
	rawResponse, err := blockchain.messenger.SendHookCallRequest(request)
	if err != nil {
		return 0
	}

	if rawResponse.GetKind() != common.BlockchainLastRoundResponse {
		log.Error("LastRound", "err", common.ErrBadHookResponseFromNode)
		return 0
	}

	response := rawResponse.(*common.MessageBlockchainLastRoundResponse)
	return response.Result
}

// LastTimeStamp forwards a message to the actual hook
func (blockchain *BlockchainHookGateway) LastTimeStamp() uint64 {
	request := common.NewMessageBlockchainLastTimeStampRequest()
	rawResponse, err := blockchain.messenger.SendHookCallRequest(request)
	if err != nil {
		return 0
	}

	if rawResponse.GetKind() != common.BlockchainLastTimeStampResponse {
		log.Error("LastTimeStamp", "err", common.ErrBadHookResponseFromNode)
		return 0
	}

	response := rawResponse.(*common.MessageBlockchainLastTimeStampResponse)
	return response.Result
}

// LastRandomSeed forwards a message to the actual hook
func (blockchain *BlockchainHookGateway) LastRandomSeed() []byte {
	request := common.NewMessageBlockchainLastRandomSeedRequest()
	rawResponse, err := blockchain.messenger.SendHookCallRequest(request)
	if err != nil {
		return nil
	}

	if rawResponse.GetKind() != common.BlockchainLastRandomSeedResponse {
		log.Error("LastRandomSeed", "err", common.ErrBadHookResponseFromNode)
		return nil
	}

	response := rawResponse.(*common.MessageBlockchainLastRandomSeedResponse)
	return response.Result
}

// LastEpoch forwards a message to the actual hook
func (blockchain *BlockchainHookGateway) LastEpoch() uint32 {
	request := common.NewMessageBlockchainLastEpochRequest()
	rawResponse, err := blockchain.messenger.SendHookCallRequest(request)
	if err != nil {
		return 0
	}

	if rawResponse.GetKind() != common.BlockchainLastEpochResponse {
		log.Error("LastEpoch", "err", common.ErrBadHookResponseFromNode)
		return 0
	}

	response := rawResponse.(*common.MessageBlockchainLastEpochResponse)
	return response.Result
}

// GetStateRootHash forwards a message to the actual hook
func (blockchain *BlockchainHookGateway) GetStateRootHash() []byte {
	request := common.NewMessageBlockchainGetStateRootHashRequest()
	rawResponse, err := blockchain.messenger.SendHookCallRequest(request)
	if err != nil {
		return nil
	}

	if rawResponse.GetKind() != common.BlockchainGetStateRootHashResponse {
		log.Error("GetStateRootHash", "err", common.ErrBadHookResponseFromNode)
		return nil
	}

	response := rawResponse.(*common.MessageBlockchainGetStateRootHashResponse)
	return response.Result
}

// CurrentNonce forwards a message to the actual hook
func (blockchain *BlockchainHookGateway) CurrentNonce() uint64 {
	request := common.NewMessageBlockchainCurrentNonceRequest()
	rawResponse, err := blockchain.messenger.SendHookCallRequest(request)
	if err != nil {
		return 0
	}

	if rawResponse.GetKind() != common.BlockchainCurrentNonceResponse {
		log.Error("CurrentNonce", "err", common.ErrBadHookResponseFromNode)
		return 0
	}

	response := rawResponse.(*common.MessageBlockchainCurrentNonceResponse)
	return response.Result
}

// CurrentRound forwards a message to the actual hook
func (blockchain *BlockchainHookGateway) CurrentRound() uint64 {
	request := common.NewMessageBlockchainCurrentRoundRequest()
	rawResponse, err := blockchain.messenger.SendHookCallRequest(request)
	if err != nil {
		return 0
	}

	if rawResponse.GetKind() != common.BlockchainCurrentRoundResponse {
		log.Error("CurrentRound", "err", common.ErrBadHookResponseFromNode)
		return 0
	}

	response := rawResponse.(*common.MessageBlockchainCurrentRoundResponse)
	return response.Result
}

// CurrentTimeStamp forwards a message to the actual hook
func (blockchain *BlockchainHookGateway) CurrentTimeStamp() uint64 {
	request := common.NewMessageBlockchainCurrentTimeStampRequest()
	rawResponse, err := blockchain.messenger.SendHookCallRequest(request)
	if err != nil {
		return 0
	}

	if rawResponse.GetKind() != common.BlockchainCurrentTimeStampResponse {
		log.Error("CurrentTimeStamp", "err", common.ErrBadHookResponseFromNode)
		return 0
	}

	response := rawResponse.(*common.MessageBlockchainCurrentTimeStampResponse)
	return response.Result
}

// CurrentRandomSeed forwards a message to the actual hook
func (blockchain *BlockchainHookGateway) CurrentRandomSeed() []byte {
	request := common.NewMessageBlockchainCurrentRandomSeedRequest()
	rawResponse, err := blockchain.messenger.SendHookCallRequest(request)
	if err != nil {
		return nil
	}

	if rawResponse.GetKind() != common.BlockchainCurrentRandomSeedResponse {
		log.Error("CurrentRandomSeed", "err", common.ErrBadHookResponseFromNode)
		return nil
	}

	response := rawResponse.(*common.MessageBlockchainCurrentRandomSeedResponse)
	return response.Result
}

// CurrentEpoch forwards a message to the actual hook
func (blockchain *BlockchainHookGateway) CurrentEpoch() uint32 {
	request := common.NewMessageBlockchainCurrentEpochRequest()
	rawResponse, err := blockchain.messenger.SendHookCallRequest(request)
	if err != nil {
		return 0
	}

	if rawResponse.GetKind() != common.BlockchainCurrentEpochResponse {
		log.Error("CurrentEpoch", "err", common.ErrBadHookResponseFromNode)
		return 0
	}

	response := rawResponse.(*common.MessageBlockchainCurrentEpochResponse)
	return response.Result
}

// ProcessBuiltInFunction forwards a message to the actual hook
func (blockchain *BlockchainHookGateway) ProcessBuiltInFunction(input *vmcommon.ContractCallInput) (*vmcommon.VMOutput, error) {
	request := common.NewMessageBlockchainProcessBuiltinFunctionRequest(*input)
	rawResponse, err := blockchain.messenger.SendHookCallRequest(request)
	if err != nil {
		return nil, err
	}

	if rawResponse.GetKind() != common.BlockchainProcessBuiltinFunctionResponse {
		return nil, common.ErrBadHookResponseFromNode
	}

	response := rawResponse.(*common.MessageBlockchainProcessBuiltinFunctionResponse)
	return response.SerializableVMOutput.ConvertToVMOutput(), response.GetError()
}

// GetBuiltinFunctionNames forwards a message to the actual hook
func (blockchain *BlockchainHookGateway) GetBuiltinFunctionNames() vmcommon.FunctionNames {
	request := common.NewMessageBlockchainGetBuiltinFunctionNamesRequest()
	rawResponse, err := blockchain.messenger.SendHookCallRequest(request)
	if err != nil {
		return make(vmcommon.FunctionNames)
	}

	if rawResponse.GetKind() != common.BlockchainGetBuiltinFunctionNamesResponse {
		return make(vmcommon.FunctionNames)
	}

	response := rawResponse.(*common.MessageBlockchainGetBuiltinFunctionNamesResponse)
	return response.FunctionNames
}

// GetAllState forwards a message to the actual hook
func (blockchain *BlockchainHookGateway) GetAllState(address []byte) (map[string][]byte, error) {
	request := common.NewMessageBlockchainGetAllStateRequest(address)
	rawResponse, err := blockchain.messenger.SendHookCallRequest(request)
	if err != nil {
		return nil, err
	}

	if rawResponse.GetKind() != common.BlockchainGetAllStateResponse {
		return nil, common.ErrBadHookResponseFromNode
	}

	response := rawResponse.(*common.MessageBlockchainGetAllStateResponse)
	return response.SerializableAllState.ConvertToMap(), response.GetError()
}

// GetUserAccount forwards a message to the actual hook
// TODO: Perhaps cache GetUserAccount()? Since when it is called with address == contract address, the whole code is fetched.
func (blockchain *BlockchainHookGateway) GetUserAccount(address []byte) (vmcommon.UserAccountHandler, error) {
	request := common.NewMessageBlockchainGetUserAccountRequest(address)
	rawResponse, err := blockchain.messenger.SendHookCallRequest(request)
	if err != nil {
		return nil, err
	}

	if rawResponse.GetKind() != common.BlockchainGetUserAccountResponse {
		return nil, common.ErrBadHookResponseFromNode
	}

	response := rawResponse.(*common.MessageBlockchainGetUserAccountResponse)
	return response.Account, response.GetError()
}

// GetCode forwards a message to the actual hook
func (blockchain *BlockchainHookGateway) GetCode(account vmcommon.UserAccountHandler) []byte {
	requestAccount := &common.Account{
		Nonce:           account.GetNonce(),
		Balance:         account.GetBalance(),
		CodeHash:        account.GetCodeHash(),
		RootHash:        account.GetRootHash(),
		Address:         account.AddressBytes(),
		DeveloperReward: account.GetDeveloperReward(),
		OwnerAddress:    account.GetOwnerAddress(),
		UserName:        account.GetUserName(),
		CodeMetadata:    account.GetCodeMetadata(),
	}

	request := common.NewMessageBlockchainGetCodeRequest(requestAccount)
	rawResponse, err := blockchain.messenger.SendHookCallRequest(request)
	if err != nil {
		return nil
	}

	if rawResponse.GetKind() != common.BlockchainGetCodeResponse {
		return nil
	}

	response := rawResponse.(*common.MessageBlockchainGetCodeResponse)
	return response.Code
}

// GetShardOfAddress forwards a message to the actual hook
func (blockchain *BlockchainHookGateway) GetShardOfAddress(address []byte) uint32 {
	request := common.NewMessageBlockchainGetShardOfAddressRequest(address)
	rawResponse, err := blockchain.messenger.SendHookCallRequest(request)
	if err != nil {
		return 0
	}

	if rawResponse.GetKind() != common.BlockchainGetShardOfAddressResponse {
		log.Error("GetShardOfAddress", "err", common.ErrBadHookResponseFromNode)
		return 0
	}

	response := rawResponse.(*common.MessageBlockchainGetShardOfAddressResponse)
	return response.Shard
}

// IsSmartContract forwards a message to the actual hook
func (blockchain *BlockchainHookGateway) IsSmartContract(address []byte) bool {
	request := common.NewMessageBlockchainIsSmartContractRequest(address)
	rawResponse, err := blockchain.messenger.SendHookCallRequest(request)
	if err != nil {
		return false
	}

	if rawResponse.GetKind() != common.BlockchainIsSmartContractResponse {
		log.Error("IsSmartContract", "err", common.ErrBadHookResponseFromNode)
		return false
	}

	response := rawResponse.(*common.MessageBlockchainIsSmartContractResponse)
	return response.Result
}

// IsPayable forwards a message to the actual hook
func (blockchain *BlockchainHookGateway) IsPayable(address []byte) (bool, error) {
	request := common.NewMessageBlockchainIsPayableRequest(address)
	rawResponse, err := blockchain.messenger.SendHookCallRequest(request)
	if err != nil {
		return false, err
	}

	if rawResponse.GetKind() != common.BlockchainIsPayableResponse {
		return false, common.ErrBadHookResponseFromNode
	}

	response := rawResponse.(*common.MessageBlockchainIsPayableResponse)
	return response.Result, response.GetError()
}

// SaveCompiledCode forwards a message to the actual hook
func (blockchain *BlockchainHookGateway) SaveCompiledCode(codeHash []byte, code []byte) {
	request := common.NewMessageBlockchainSaveCompiledCodeRequest(codeHash, code)
	rawResponse, err := blockchain.messenger.SendHookCallRequest(request)
	if err != nil {
		return
	}

	if rawResponse.GetKind() != common.BlockchainSaveCompiledCodeResponse {
		log.Error("SaveCompiledCode", "err", common.ErrBadHookResponseFromNode)
	}
}

// GetCompiledCode forwards a message to the actual hook
func (blockchain *BlockchainHookGateway) GetCompiledCode(codeHash []byte) (bool, []byte) {
	request := common.NewMessageBlockchainGetCompiledCodeRequest(codeHash)
	rawResponse, err := blockchain.messenger.SendHookCallRequest(request)
	if err != nil {
		return false, nil
	}

	if rawResponse.GetKind() != common.BlockchainGetCompiledCodeResponse {
		return false, nil
	}

	response := rawResponse.(*common.MessageBlockchainGetCompiledCodeResponse)
	return response.Found, response.Code
}

// ClearCompiledCodes nothing to do - this needs to be called by nodepart only
func (blockchain *BlockchainHookGateway) ClearCompiledCodes() {
}

// IsInterfaceNil returns true if underlying implementation is nil
func (blockchain *BlockchainHookGateway) IsInterfaceNil() bool {
	return blockchain == nil
}
