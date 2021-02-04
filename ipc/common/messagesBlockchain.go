package common

import (
	"github.com/ElrondNetwork/elrond-go/core/vmcommon"
)

// MessageBlockchainNewAddressRequest represents a request message
type MessageBlockchainNewAddressRequest struct {
	Message
	CreatorAddress []byte
	CreatorNonce   uint64
	VMType         []byte
}

// NewMessageBlockchainNewAddressRequest creates a request message
func NewMessageBlockchainNewAddressRequest(creatorAddress []byte, creatorNonce uint64, vmType []byte) *MessageBlockchainNewAddressRequest {
	message := &MessageBlockchainNewAddressRequest{}
	message.Kind = BlockchainNewAddressRequest
	message.CreatorAddress = creatorAddress
	message.CreatorNonce = creatorNonce
	message.VMType = vmType
	return message
}

// MessageBlockchainNewAddressResponse represents a response message
type MessageBlockchainNewAddressResponse struct {
	Message
	Result []byte
}

// NewMessageBlockchainNewAddressResponse creates a response message
func NewMessageBlockchainNewAddressResponse(result []byte, err error) *MessageBlockchainNewAddressResponse {
	message := &MessageBlockchainNewAddressResponse{}
	message.Kind = BlockchainNewAddressResponse
	message.Result = result
	message.SetError(err)
	return message
}

// MessageBlockchainGetStorageDataRequest represents a request message
type MessageBlockchainGetStorageDataRequest struct {
	Message
	Address []byte
	Index   []byte
}

// NewMessageBlockchainGetStorageDataRequest creates a request message
func NewMessageBlockchainGetStorageDataRequest(address []byte, index []byte) *MessageBlockchainGetStorageDataRequest {
	message := &MessageBlockchainGetStorageDataRequest{}
	message.Kind = BlockchainGetStorageDataRequest
	message.Address = address
	message.Index = index
	return message
}

// MessageBlockchainGetStorageDataResponse represents a response message
type MessageBlockchainGetStorageDataResponse struct {
	Message
	Data []byte
}

// NewMessageBlockchainGetStorageDataResponse creates a response message
func NewMessageBlockchainGetStorageDataResponse(data []byte, err error) *MessageBlockchainGetStorageDataResponse {
	message := &MessageBlockchainGetStorageDataResponse{}
	message.Kind = BlockchainGetStorageDataResponse
	message.Data = data
	message.SetError(err)
	return message
}

// MessageBlockchainGetBlockhashRequest represents a request message
type MessageBlockchainGetBlockhashRequest struct {
	Message
	Nonce uint64
}

// NewMessageBlockchainGetBlockhashRequest creates a request message
func NewMessageBlockchainGetBlockhashRequest(nonce uint64) *MessageBlockchainGetBlockhashRequest {
	message := &MessageBlockchainGetBlockhashRequest{}
	message.Kind = BlockchainGetBlockhashRequest
	message.Nonce = nonce
	return message
}

// MessageBlockchainGetBlockhashResponse represents a response message
type MessageBlockchainGetBlockhashResponse struct {
	Message
	Result []byte
}

// NewMessageBlockchainGetBlockhashResponse creates a response message
func NewMessageBlockchainGetBlockhashResponse(result []byte, err error) *MessageBlockchainGetBlockhashResponse {
	message := &MessageBlockchainGetBlockhashResponse{}
	message.Kind = BlockchainGetBlockhashResponse
	message.Result = result
	message.SetError(err)
	return message
}

// MessageBlockchainLastNonceRequest represents a request message
type MessageBlockchainLastNonceRequest struct {
	Message
}

// NewMessageBlockchainLastNonceRequest creates a request message
func NewMessageBlockchainLastNonceRequest() *MessageBlockchainLastNonceRequest {
	message := &MessageBlockchainLastNonceRequest{}
	message.Kind = BlockchainLastNonceRequest

	return message
}

// MessageBlockchainLastNonceResponse represents a response message
type MessageBlockchainLastNonceResponse struct {
	Message
	Result uint64
}

// NewMessageBlockchainLastNonceResponse creates a response message
func NewMessageBlockchainLastNonceResponse(result uint64) *MessageBlockchainLastNonceResponse {
	message := &MessageBlockchainLastNonceResponse{}
	message.Kind = BlockchainLastNonceResponse
	message.Result = result
	return message
}

// MessageBlockchainLastRoundRequest represents a request message
type MessageBlockchainLastRoundRequest struct {
	Message
}

// NewMessageBlockchainLastRoundRequest creates a request message
func NewMessageBlockchainLastRoundRequest() *MessageBlockchainLastRoundRequest {
	message := &MessageBlockchainLastRoundRequest{}
	message.Kind = BlockchainLastRoundRequest

	return message
}

// MessageBlockchainLastRoundResponse represents a response message
type MessageBlockchainLastRoundResponse struct {
	Message
	Result uint64
}

// NewMessageBlockchainLastRoundResponse creates a response message
func NewMessageBlockchainLastRoundResponse(result uint64) *MessageBlockchainLastRoundResponse {
	message := &MessageBlockchainLastRoundResponse{}
	message.Kind = BlockchainLastRoundResponse
	message.Result = result
	return message
}

// MessageBlockchainLastTimeStampRequest represents a request message
type MessageBlockchainLastTimeStampRequest struct {
	Message
}

// NewMessageBlockchainLastTimeStampRequest creates a request message
func NewMessageBlockchainLastTimeStampRequest() *MessageBlockchainLastTimeStampRequest {
	message := &MessageBlockchainLastTimeStampRequest{}
	message.Kind = BlockchainLastTimeStampRequest

	return message
}

// MessageBlockchainLastTimeStampResponse represents a response message
type MessageBlockchainLastTimeStampResponse struct {
	Message
	Result uint64
}

// NewMessageBlockchainLastTimeStampResponse creates a response message
func NewMessageBlockchainLastTimeStampResponse(result uint64) *MessageBlockchainLastTimeStampResponse {
	message := &MessageBlockchainLastTimeStampResponse{}
	message.Kind = BlockchainLastTimeStampResponse
	message.Result = result
	return message
}

// MessageBlockchainLastRandomSeedRequest represents a request message
type MessageBlockchainLastRandomSeedRequest struct {
	Message
}

// NewMessageBlockchainLastRandomSeedRequest creates a request message
func NewMessageBlockchainLastRandomSeedRequest() *MessageBlockchainLastRandomSeedRequest {
	message := &MessageBlockchainLastRandomSeedRequest{}
	message.Kind = BlockchainLastRandomSeedRequest

	return message
}

// MessageBlockchainLastRandomSeedResponse represents a response message
type MessageBlockchainLastRandomSeedResponse struct {
	Message
	Result []byte
}

// NewMessageBlockchainLastRandomSeedResponse creates a response message
func NewMessageBlockchainLastRandomSeedResponse(result []byte) *MessageBlockchainLastRandomSeedResponse {
	message := &MessageBlockchainLastRandomSeedResponse{}
	message.Kind = BlockchainLastRandomSeedResponse
	message.Result = result
	return message
}

// MessageBlockchainLastEpochRequest represents a request message
type MessageBlockchainLastEpochRequest struct {
	Message
}

// NewMessageBlockchainLastEpochRequest creates a request message
func NewMessageBlockchainLastEpochRequest() *MessageBlockchainLastEpochRequest {
	message := &MessageBlockchainLastEpochRequest{}
	message.Kind = BlockchainLastEpochRequest

	return message
}

// MessageBlockchainLastEpochResponse represents a response message
type MessageBlockchainLastEpochResponse struct {
	Message
	Result uint32
}

// NewMessageBlockchainLastEpochResponse creates a response message
func NewMessageBlockchainLastEpochResponse(result uint32) *MessageBlockchainLastEpochResponse {
	message := &MessageBlockchainLastEpochResponse{}
	message.Kind = BlockchainLastEpochResponse
	message.Result = result
	return message
}

// MessageBlockchainGetStateRootHashRequest represents a request message
type MessageBlockchainGetStateRootHashRequest struct {
	Message
}

// NewMessageBlockchainGetStateRootHashRequest creates a request message
func NewMessageBlockchainGetStateRootHashRequest() *MessageBlockchainGetStateRootHashRequest {
	message := &MessageBlockchainGetStateRootHashRequest{}
	message.Kind = BlockchainGetStateRootHashRequest

	return message
}

// MessageBlockchainGetStateRootHashResponse represents a response message
type MessageBlockchainGetStateRootHashResponse struct {
	Message
	Result []byte
}

// NewMessageBlockchainGetStateRootHashResponse creates a response message
func NewMessageBlockchainGetStateRootHashResponse(result []byte) *MessageBlockchainGetStateRootHashResponse {
	message := &MessageBlockchainGetStateRootHashResponse{}
	message.Kind = BlockchainGetStateRootHashResponse
	message.Result = result
	return message
}

// MessageBlockchainCurrentNonceRequest represents a request message
type MessageBlockchainCurrentNonceRequest struct {
	Message
}

// NewMessageBlockchainCurrentNonceRequest creates a request message
func NewMessageBlockchainCurrentNonceRequest() *MessageBlockchainCurrentNonceRequest {
	message := &MessageBlockchainCurrentNonceRequest{}
	message.Kind = BlockchainCurrentNonceRequest

	return message
}

// MessageBlockchainCurrentNonceResponse represents a response message
type MessageBlockchainCurrentNonceResponse struct {
	Message
	Result uint64
}

// NewMessageBlockchainCurrentNonceResponse creates a response message
func NewMessageBlockchainCurrentNonceResponse(result uint64) *MessageBlockchainCurrentNonceResponse {
	message := &MessageBlockchainCurrentNonceResponse{}
	message.Kind = BlockchainCurrentNonceResponse
	message.Result = result
	return message
}

// MessageBlockchainCurrentRoundRequest represents a request message
type MessageBlockchainCurrentRoundRequest struct {
	Message
}

// NewMessageBlockchainCurrentRoundRequest creates a request message
func NewMessageBlockchainCurrentRoundRequest() *MessageBlockchainCurrentRoundRequest {
	message := &MessageBlockchainCurrentRoundRequest{}
	message.Kind = BlockchainCurrentRoundRequest

	return message
}

// MessageBlockchainCurrentRoundResponse represents a response message
type MessageBlockchainCurrentRoundResponse struct {
	Message
	Result uint64
}

// NewMessageBlockchainCurrentRoundResponse creates a response message
func NewMessageBlockchainCurrentRoundResponse(result uint64) *MessageBlockchainCurrentRoundResponse {
	message := &MessageBlockchainCurrentRoundResponse{}
	message.Kind = BlockchainCurrentRoundResponse
	message.Result = result
	return message
}

// MessageBlockchainCurrentTimeStampRequest represents a request message
type MessageBlockchainCurrentTimeStampRequest struct {
	Message
}

// NewMessageBlockchainCurrentTimeStampRequest creates a request message
func NewMessageBlockchainCurrentTimeStampRequest() *MessageBlockchainCurrentTimeStampRequest {
	message := &MessageBlockchainCurrentTimeStampRequest{}
	message.Kind = BlockchainCurrentTimeStampRequest

	return message
}

// MessageBlockchainCurrentTimeStampResponse represents a response message
type MessageBlockchainCurrentTimeStampResponse struct {
	Message
	Result uint64
}

// NewMessageBlockchainCurrentTimeStampResponse creates a response message
func NewMessageBlockchainCurrentTimeStampResponse(result uint64) *MessageBlockchainCurrentTimeStampResponse {
	message := &MessageBlockchainCurrentTimeStampResponse{}
	message.Kind = BlockchainCurrentTimeStampResponse
	message.Result = result
	return message
}

// MessageBlockchainCurrentRandomSeedRequest represents a request message
type MessageBlockchainCurrentRandomSeedRequest struct {
	Message
}

// NewMessageBlockchainCurrentRandomSeedRequest creates a request message
func NewMessageBlockchainCurrentRandomSeedRequest() *MessageBlockchainCurrentRandomSeedRequest {
	message := &MessageBlockchainCurrentRandomSeedRequest{}
	message.Kind = BlockchainCurrentRandomSeedRequest

	return message
}

// MessageBlockchainCurrentRandomSeedResponse represents a response message
type MessageBlockchainCurrentRandomSeedResponse struct {
	Message
	Result []byte
}

// NewMessageBlockchainCurrentRandomSeedResponse creates a response message
func NewMessageBlockchainCurrentRandomSeedResponse(result []byte) *MessageBlockchainCurrentRandomSeedResponse {
	message := &MessageBlockchainCurrentRandomSeedResponse{}
	message.Kind = BlockchainCurrentRandomSeedResponse
	message.Result = result
	return message
}

// MessageBlockchainCurrentEpochRequest represents a request message
type MessageBlockchainCurrentEpochRequest struct {
	Message
}

// NewMessageBlockchainCurrentEpochRequest creates a request message
func NewMessageBlockchainCurrentEpochRequest() *MessageBlockchainCurrentEpochRequest {
	message := &MessageBlockchainCurrentEpochRequest{}
	message.Kind = BlockchainCurrentEpochRequest

	return message
}

// MessageBlockchainCurrentEpochResponse represents a response message
type MessageBlockchainCurrentEpochResponse struct {
	Message
	Result uint32
}

// NewMessageBlockchainCurrentEpochResponse creates a response message
func NewMessageBlockchainCurrentEpochResponse(result uint32) *MessageBlockchainCurrentEpochResponse {
	message := &MessageBlockchainCurrentEpochResponse{}
	message.Kind = BlockchainCurrentEpochResponse
	message.Result = result
	return message
}

// MessageBlockchainProcessBuiltinFunctionRequest represents a request message
type MessageBlockchainProcessBuiltinFunctionRequest struct {
	Message
	CallInput vmcommon.ContractCallInput
}

// NewMessageBlockchainProcessBuiltinFunctionRequest creates a request message
func NewMessageBlockchainProcessBuiltinFunctionRequest(callInput vmcommon.ContractCallInput) *MessageBlockchainProcessBuiltinFunctionRequest {
	message := &MessageBlockchainProcessBuiltinFunctionRequest{}
	message.Kind = BlockchainProcessBuiltinFunctionRequest
	message.CallInput = callInput

	return message
}

// MessageBlockchainProcessBuiltinFunctionResponse represents a response message
type MessageBlockchainProcessBuiltinFunctionResponse struct {
	Message
	SerializableVMOutput *SerializableVMOutput
}

// NewMessageBlockchainProcessBuiltinFunctionResponse creates a response message
func NewMessageBlockchainProcessBuiltinFunctionResponse(vmOutput *vmcommon.VMOutput, err error) *MessageBlockchainProcessBuiltinFunctionResponse {
	message := &MessageBlockchainProcessBuiltinFunctionResponse{}
	message.Kind = BlockchainProcessBuiltinFunctionResponse
	message.SerializableVMOutput = NewSerializableVMOutput(vmOutput)
	message.SetError(err)
	return message
}

// MessageBlockchainGetBuiltinFunctionNamesRequest represents a request message
type MessageBlockchainGetBuiltinFunctionNamesRequest struct {
	Message
}

// NewMessageBlockchainGetBuiltinFunctionNamesRequest creates a request message
func NewMessageBlockchainGetBuiltinFunctionNamesRequest() *MessageBlockchainGetBuiltinFunctionNamesRequest {
	message := &MessageBlockchainGetBuiltinFunctionNamesRequest{}
	message.Kind = BlockchainGetBuiltinFunctionNamesRequest

	return message
}

// MessageBlockchainGetBuiltinFunctionNamesResponse represents a response message
type MessageBlockchainGetBuiltinFunctionNamesResponse struct {
	Message
	FunctionNames vmcommon.FunctionNames
}

// NewMessageBlockchainGetBuiltinFunctionNamesResponse creates a response message
func NewMessageBlockchainGetBuiltinFunctionNamesResponse(functionNames vmcommon.FunctionNames) *MessageBlockchainGetBuiltinFunctionNamesResponse {
	message := &MessageBlockchainGetBuiltinFunctionNamesResponse{}
	message.Kind = BlockchainGetBuiltinFunctionNamesResponse
	message.FunctionNames = functionNames

	return message
}

// MessageBlockchainGetAllStateRequest represents a request message
type MessageBlockchainGetAllStateRequest struct {
	Message
	Address []byte
}

// NewMessageBlockchainGetAllStateRequest creates a request message
func NewMessageBlockchainGetAllStateRequest(address []byte) *MessageBlockchainGetAllStateRequest {
	message := &MessageBlockchainGetAllStateRequest{}
	message.Kind = BlockchainGetAllStateRequest
	message.Address = address
	return message
}

// MessageBlockchainGetAllStateResponse represents a response message
type MessageBlockchainGetAllStateResponse struct {
	Message
	SerializableAllState *SerializableMapStringBytes
}

// NewMessageBlockchainGetAllStateResponse creates a response message
func NewMessageBlockchainGetAllStateResponse(state map[string][]byte, err error) *MessageBlockchainGetAllStateResponse {
	message := &MessageBlockchainGetAllStateResponse{}
	message.Kind = BlockchainGetAllStateResponse
	message.SerializableAllState = NewSerializableMapStringBytes(state)
	message.SetError(err)
	return message
}

// MessageBlockchainGetUserAccountRequest represents a request message
type MessageBlockchainGetUserAccountRequest struct {
	Message
	Address []byte
}

// NewMessageBlockchainGetUserAccountRequest creates a request message
func NewMessageBlockchainGetUserAccountRequest(address []byte) *MessageBlockchainGetUserAccountRequest {
	message := &MessageBlockchainGetUserAccountRequest{}
	message.Kind = BlockchainGetUserAccountRequest
	message.Address = address
	return message
}

// MessageBlockchainGetUserAccountResponse represents a response message
type MessageBlockchainGetUserAccountResponse struct {
	Message
	Account *Account
}

// NewMessageBlockchainGetUserAccountResponse creates a response message
func NewMessageBlockchainGetUserAccountResponse(account *Account, err error) *MessageBlockchainGetUserAccountResponse {
	message := &MessageBlockchainGetUserAccountResponse{}
	message.Kind = BlockchainGetUserAccountResponse
	message.Account = account
	message.SetError(err)
	return message
}

// NewMessageBlockchainGetCodeRequest represents a request message
type MessageBlockchainGetCodeRequest struct {
	Message
	Account *Account
}

// NewMessageBlockchainGetCodeRequest creates a request message
func NewMessageBlockchainGetCodeRequest(account *Account) *MessageBlockchainGetCodeRequest {
	message := &MessageBlockchainGetCodeRequest{}
	message.Kind = BlockchainGetCodeRequest
	message.Account = account
	return message
}

// MessageBlockchainGetCodeResponse represents a response message
type MessageBlockchainGetCodeResponse struct {
	Message
	Code []byte
}

// NewMessageBlockchainGetCodeResponse creates a response message
func NewMessageBlockchainGetCodeResponse(code []byte) *MessageBlockchainGetCodeResponse {
	message := &MessageBlockchainGetCodeResponse{}
	message.Kind = BlockchainGetCodeResponse
	message.Code = code
	return message
}

// MessageBlockchainGetShardOfAddressRequest represents a request message
type MessageBlockchainGetShardOfAddressRequest struct {
	Message
	Address []byte
}

// NewMessageBlockchainGetShardOfAddressRequest creates a request message
func NewMessageBlockchainGetShardOfAddressRequest(address []byte) *MessageBlockchainGetShardOfAddressRequest {
	message := &MessageBlockchainGetShardOfAddressRequest{}
	message.Kind = BlockchainGetShardOfAddressRequest
	message.Address = address
	return message
}

// MessageBlockchainGetShardOfAddressResponse represents a response message
type MessageBlockchainGetShardOfAddressResponse struct {
	Message
	Shard uint32
}

// NewMessageBlockchainGetShardOfAddressResponse creates a response message
func NewMessageBlockchainGetShardOfAddressResponse(shard uint32) *MessageBlockchainGetShardOfAddressResponse {
	message := &MessageBlockchainGetShardOfAddressResponse{}
	message.Kind = BlockchainGetShardOfAddressResponse
	message.Shard = shard
	return message
}

// MessageBlockchainIsSmartContractRequest represents a request message
type MessageBlockchainIsSmartContractRequest struct {
	Message
	Address []byte
}

// NewMessageBlockchainIsSmartContractRequest creates a request message
func NewMessageBlockchainIsSmartContractRequest(address []byte) *MessageBlockchainIsSmartContractRequest {
	message := &MessageBlockchainIsSmartContractRequest{}
	message.Kind = BlockchainIsSmartContractRequest
	message.Address = address
	return message
}

// MessageBlockchainIsSmartContractResponse represents a response message
type MessageBlockchainIsSmartContractResponse struct {
	Message
	Result bool
}

// NewMessageBlockchainIsSmartContractResponse creates a response message
func NewMessageBlockchainIsSmartContractResponse(result bool) *MessageBlockchainIsSmartContractResponse {
	message := &MessageBlockchainIsSmartContractResponse{}
	message.Kind = BlockchainIsSmartContractResponse
	message.Result = result
	return message
}

// MessageBlockchainIsPayableRequest represents a request message
type MessageBlockchainIsPayableRequest struct {
	Message
	Address []byte
}

// NewMessageBlockchainIsPayableRequest creates a request message
func NewMessageBlockchainIsPayableRequest(address []byte) *MessageBlockchainIsPayableRequest {
	message := &MessageBlockchainIsPayableRequest{}
	message.Kind = BlockchainIsPayableRequest
	message.Address = address
	return message
}

// MessageBlockchainIsPayableResponse represents a response message
type MessageBlockchainIsPayableResponse struct {
	Message
	Result bool
}

// NewMessageBlockchainIsPayableResponse creates a response message
func NewMessageBlockchainIsPayableResponse(result bool, err error) *MessageBlockchainIsPayableResponse {
	message := &MessageBlockchainIsPayableResponse{}
	message.Kind = BlockchainIsPayableResponse
	message.Result = result
	message.SetError(err)
	return message
}

// MessageBlockchainSaveCompiledCodeRequest represents a request message
type MessageBlockchainSaveCompiledCodeRequest struct {
	Message
	CodeHash []byte
	Code     []byte
}

// NewMessageBlockchainSaveCompiledCodeRequest creates a request message
func NewMessageBlockchainSaveCompiledCodeRequest(codeHash []byte, code []byte) *MessageBlockchainSaveCompiledCodeRequest {
	message := &MessageBlockchainSaveCompiledCodeRequest{}
	message.Kind = BlockchainSaveCompiledCodeRequest
	message.CodeHash = codeHash
	message.Code = code
	return message
}

// MessageBlockchainSaveCompiledCodeResponse represents a response message
type MessageBlockchainSaveCompiledCodeResponse struct {
	Message
}

// NewMessageBlockchainSaveCompiledCodeResponse creates a response message
func NewMessageBlockchainSaveCompiledCodeResponse() *MessageBlockchainSaveCompiledCodeResponse {
	message := &MessageBlockchainSaveCompiledCodeResponse{}
	message.Kind = BlockchainSaveCompiledCodeResponse
	return message
}

// MessageBlockchainGetCompiledCodeRequest represents a request message
type MessageBlockchainGetCompiledCodeRequest struct {
	Message
	CodeHash []byte
}

// NewMessageBlockchainGetCompiledCodeRequest creates a request message
func NewMessageBlockchainGetCompiledCodeRequest(codeHash []byte) *MessageBlockchainGetCompiledCodeRequest {
	message := &MessageBlockchainGetCompiledCodeRequest{}
	message.Kind = BlockchainGetCompiledCodeRequest
	message.CodeHash = codeHash
	return message
}

// MessageBlockchainGetCompiledCodeResponse represents a response message
type MessageBlockchainGetCompiledCodeResponse struct {
	Message
	Found bool
	Code  []byte
}

// NewMessageBlockchainGetCompiledCodeResponse creates a response message
func NewMessageBlockchainGetCompiledCodeResponse(result bool, code []byte) *MessageBlockchainGetCompiledCodeResponse {
	message := &MessageBlockchainGetCompiledCodeResponse{}
	message.Kind = BlockchainGetCompiledCodeResponse
	message.Found = result
	message.Code = code
	return message
}
