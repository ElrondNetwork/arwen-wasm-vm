package common

import (
	"fmt"
	"math"
)

// MessageKind is the kind of a message (that is passed between the Node and Arwen)
type MessageKind uint32

const (
	FirstKind MessageKind = iota
	Initialize
	Stop
	ContractDeployRequest
	ContractCallRequest
	ContractResponse
	GasScheduleChangeRequest
	GasScheduleChangeResponse
	BlockchainNewAddressRequest
	BlockchainNewAddressResponse
	BlockchainGetStorageDataRequest
	BlockchainGetStorageDataResponse
	BlockchainGetBlockhashRequest
	BlockchainGetBlockhashResponse
	BlockchainLastNonceRequest
	BlockchainLastNonceResponse
	BlockchainLastRoundRequest
	BlockchainLastRoundResponse
	BlockchainLastTimeStampRequest
	BlockchainLastTimeStampResponse
	BlockchainLastRandomSeedRequest
	BlockchainLastRandomSeedResponse
	BlockchainLastEpochRequest
	BlockchainLastEpochResponse
	BlockchainGetStateRootHashRequest
	BlockchainGetStateRootHashResponse
	BlockchainCurrentNonceRequest
	BlockchainCurrentNonceResponse
	BlockchainCurrentRoundRequest
	BlockchainCurrentRoundResponse
	BlockchainCurrentTimeStampRequest
	BlockchainCurrentTimeStampResponse
	BlockchainCurrentRandomSeedRequest
	BlockchainCurrentRandomSeedResponse
	BlockchainCurrentEpochRequest
	BlockchainCurrentEpochResponse
	BlockchainProcessBuiltinFunctionRequest
	BlockchainProcessBuiltinFunctionResponse
	BlockchainGetBuiltinFunctionNamesRequest
	BlockchainGetBuiltinFunctionNamesResponse
	BlockchainGetAllStateRequest
	BlockchainGetAllStateResponse
	BlockchainGetUserAccountRequest
	BlockchainGetUserAccountResponse
	BlockchainGetCodeRequest
	BlockchainGetCodeResponse
	BlockchainGetShardOfAddressRequest
	BlockchainGetShardOfAddressResponse
	BlockchainIsPayableRequest
	BlockchainIsPayableResponse
	BlockchainIsSmartContractRequest
	BlockchainIsSmartContractResponse
	BlockchainSaveCompiledCodeRequest
	BlockchainSaveCompiledCodeResponse
	BlockchainGetCompiledCodeRequest
	BlockchainGetCompiledCodeResponse
	DiagnoseWaitRequest
	DiagnoseWaitResponse
	VersionRequest
	VersionResponse
	UndefinedRequestOrResponse
	LastKind
)

var messageKindNameByID = map[MessageKind]string{}

func init() {
	messageKindNameByID[FirstKind] = "FirstKind"
	messageKindNameByID[Initialize] = "Initialize"
	messageKindNameByID[Stop] = "Stop"
	messageKindNameByID[ContractDeployRequest] = "ContractDeployRequest"
	messageKindNameByID[ContractCallRequest] = "ContractCallRequest"
	messageKindNameByID[ContractResponse] = "ContractResponse"
	messageKindNameByID[GasScheduleChangeRequest] = "GasScheduleChangeRequest"
	messageKindNameByID[GasScheduleChangeResponse] = "GasScheduleChangeResponse"
	messageKindNameByID[BlockchainNewAddressRequest] = "BlockchainNewAddressRequest"
	messageKindNameByID[BlockchainNewAddressResponse] = "BlockchainNewAddressResponse"
	messageKindNameByID[BlockchainGetStorageDataRequest] = "BlockchainGetStorageDataRequest"
	messageKindNameByID[BlockchainGetStorageDataResponse] = "BlockchainGetStorageDataResponse"
	messageKindNameByID[BlockchainGetBlockhashRequest] = "BlockchainGetBlockhashRequest"
	messageKindNameByID[BlockchainGetBlockhashResponse] = "BlockchainGetBlockhashResponse"
	messageKindNameByID[BlockchainLastNonceRequest] = "BlockchainLastNonceRequest"
	messageKindNameByID[BlockchainLastNonceResponse] = "BlockchainLastNonceResponse"
	messageKindNameByID[BlockchainLastRoundRequest] = "BlockchainLastRoundRequest"
	messageKindNameByID[BlockchainLastRoundResponse] = "BlockchainLastRoundResponse"
	messageKindNameByID[BlockchainLastTimeStampRequest] = "BlockchainLastTimeStampRequest"
	messageKindNameByID[BlockchainLastTimeStampResponse] = "BlockchainLastTimeStampResponse"
	messageKindNameByID[BlockchainLastRandomSeedRequest] = "BlockchainLastRandomSeedRequest"
	messageKindNameByID[BlockchainLastRandomSeedResponse] = "BlockchainLastRandomSeedResponse"
	messageKindNameByID[BlockchainLastEpochRequest] = "BlockchainLastEpochRequest"
	messageKindNameByID[BlockchainLastEpochResponse] = "BlockchainLastEpochResponse"
	messageKindNameByID[BlockchainGetStateRootHashRequest] = "BlockchainGetStateRootHashRequest"
	messageKindNameByID[BlockchainGetStateRootHashResponse] = "BlockchainGetStateRootHashResponse"
	messageKindNameByID[BlockchainCurrentNonceRequest] = "BlockchainCurrentNonceRequest"
	messageKindNameByID[BlockchainCurrentNonceResponse] = "BlockchainCurrentNonceResponse"
	messageKindNameByID[BlockchainCurrentRoundRequest] = "BlockchainCurrentRoundRequest"
	messageKindNameByID[BlockchainCurrentRoundResponse] = "BlockchainCurrentRoundResponse"
	messageKindNameByID[BlockchainCurrentTimeStampRequest] = "BlockchainCurrentTimeStampRequest"
	messageKindNameByID[BlockchainCurrentTimeStampResponse] = "BlockchainCurrentTimeStampResponse"
	messageKindNameByID[BlockchainCurrentRandomSeedRequest] = "BlockchainCurrentRandomSeedRequest"
	messageKindNameByID[BlockchainCurrentRandomSeedResponse] = "BlockchainCurrentRandomSeedResponse"
	messageKindNameByID[BlockchainCurrentEpochRequest] = "BlockchainCurrentEpochRequest"
	messageKindNameByID[BlockchainCurrentEpochResponse] = "BlockchainCurrentEpochResponse"
	messageKindNameByID[BlockchainProcessBuiltinFunctionRequest] = "BlockchainProcessBuiltinFunctionRequest"
	messageKindNameByID[BlockchainProcessBuiltinFunctionResponse] = "BlockchainProcessBuiltinFunctionResponse"
	messageKindNameByID[BlockchainGetBuiltinFunctionNamesRequest] = "BlockchainGetBuiltinFunctionNamesRequest"
	messageKindNameByID[BlockchainGetBuiltinFunctionNamesResponse] = "BlockchainGetBuiltinFunctionNamesResponse"
	messageKindNameByID[BlockchainGetAllStateRequest] = "BlockchainGetAllStateRequest"
	messageKindNameByID[BlockchainGetAllStateResponse] = "BlockchainGetAllStateResponse"
	messageKindNameByID[BlockchainGetUserAccountRequest] = "BlockchainGetUserAccountRequest"
	messageKindNameByID[BlockchainGetUserAccountResponse] = "BlockchainGetUserAccountResponse"
	messageKindNameByID[BlockchainGetCodeRequest] = "BlockchainGetCodeRequest"
	messageKindNameByID[BlockchainGetCodeResponse] = "BlockchainGetCodeResponse"
	messageKindNameByID[BlockchainGetShardOfAddressRequest] = "BlockchainGetShardOfAddressRequest"
	messageKindNameByID[BlockchainGetShardOfAddressResponse] = "BlockchainGetShardOfAddressResponse"
	messageKindNameByID[BlockchainIsSmartContractRequest] = "BlockchainIsSmartContractRequest"
	messageKindNameByID[BlockchainIsSmartContractResponse] = "BlockchainIsSmartContractResponse"
	messageKindNameByID[BlockchainIsPayableRequest] = "BlockchainIsPayableRequest"
	messageKindNameByID[BlockchainIsPayableResponse] = "BlockchainIsPayableResponse"
	messageKindNameByID[BlockchainGetCompiledCodeResponse] = "BlockchainGetCompiledCodeResponse"
	messageKindNameByID[BlockchainGetCompiledCodeRequest] = "BlockchainGetCompiledCodeRequest"
	messageKindNameByID[BlockchainSaveCompiledCodeRequest] = "BlockchainSaveCompiledCodeRequest"
	messageKindNameByID[BlockchainSaveCompiledCodeResponse] = "BlockchainSaveCompiledCodeResponse"
	messageKindNameByID[DiagnoseWaitRequest] = "DiagnoseWaitRequest"
	messageKindNameByID[DiagnoseWaitResponse] = "DiagnoseWaitResponse"
	messageKindNameByID[VersionRequest] = "VersionRequest"
	messageKindNameByID[VersionResponse] = "VersionResponse"
	messageKindNameByID[UndefinedRequestOrResponse] = "UndefinedRequestOrResponse"
	messageKindNameByID[LastKind] = "LastKind"
}

// MessageHandler is a message abstraction
type MessageHandler interface {
	GetNonce() uint32
	SetNonce(nonce uint32)
	GetKind() MessageKind
	SetKind(kind MessageKind)
	GetError() error
	SetError(err error)
	GetKindName() string
	DebugString() string
}

// Message is the implementation of the abstraction
type Message struct {
	DialogueNonce uint32
	Kind          MessageKind
	ErrorMessage  string
}

// GetNonce gets the dialogue nonce
func (message *Message) GetNonce() uint32 {
	return message.DialogueNonce
}

// SetNonce sets the dialogue nonce
func (message *Message) SetNonce(nonce uint32) {
	message.DialogueNonce = nonce
}

// GetKind gets the message kind
func (message *Message) GetKind() MessageKind {
	return message.Kind
}

// SetKind sets the message kind
func (message *Message) SetKind(kind MessageKind) {
	message.Kind = kind
}

// GetError gets the error within the message
func (message *Message) GetError() error {
	if len(message.ErrorMessage) == 0 {
		return nil
	}

	return fmt.Errorf(message.ErrorMessage)
}

// SetError sets the error within the message
func (message *Message) SetError(err error) {
	if err != nil {
		message.ErrorMessage = err.Error()
	}
}

// GetKindName gets the kind name
func (message *Message) GetKindName() string {
	kindName := messageKindNameByID[message.Kind]
	return kindName
}

// DebugString is a debug representation of the message
func (message *Message) DebugString() string {
	kindName := messageKindNameByID[message.Kind]
	return fmt.Sprintf("[kind=%s nonce=%d err=%s]", kindName, message.DialogueNonce, message.ErrorMessage)
}

// MessageInitialize is a message sent by Node to initialize Arwen
type MessageInitialize struct {
	Message
	Arguments ArwenArguments
}

// NewMessageInitialize creates a new message
func NewMessageInitialize(arguments ArwenArguments) *MessageInitialize {
	message := &MessageInitialize{}
	message.Kind = Initialize
	message.Arguments = arguments
	return message
}

// MessageStop is a message sent by Node to stop Arwen
type MessageStop struct {
	Message
}

// NewMessageStop creates a new message
func NewMessageStop() *MessageStop {
	message := &MessageStop{}
	message.Kind = Stop
	return message
}

// UndefinedMessage is an undefined message
type UndefinedMessage struct {
	Message
}

// NewUndefinedMessage creates an undefined message
func NewUndefinedMessage() *UndefinedMessage {
	message := &UndefinedMessage{}
	message.Kind = UndefinedRequestOrResponse
	message.SetNonce(math.MaxUint32)
	return message
}

// MessageReplier is a callback signature
type MessageReplier func(MessageHandler) MessageHandler

// CreateReplySlots creates a slice of no-operation repliers, to be substituted with actual repliers (by message listeners)
func CreateReplySlots(noopReplier MessageReplier) []MessageReplier {
	slots := make([]MessageReplier, LastKind)
	for i := 0; i < len(slots); i++ {
		slots[i] = noopReplier
	}

	return slots
}

// IsHookCall returns whether a message is a hook call
func IsHookCall(message MessageHandler) bool {
	kind := message.GetKind()
	return kind >= BlockchainNewAddressRequest && kind <= BlockchainGetCompiledCodeResponse
}

// IsStopRequest returns whether a message is a stop request
func IsStopRequest(message MessageHandler) bool {
	return message.GetKind() == Stop
}

// IsVersionResponse returns version response
func IsVersionResponse(message MessageHandler) bool {
	return message.GetKind() == VersionResponse
}

// IsContractResponse returns whether a message is a contract response
func IsContractResponse(message MessageHandler) bool {
	return message.GetKind() == ContractResponse
}

// IsGasScheduleChangeResponse returns a message with gas schedule response
func IsGasScheduleChangeResponse(message MessageHandler) bool {
	return message.GetKind() == GasScheduleChangeResponse
}

// IsDiagnose returns whether a message is a diagnose request
func IsDiagnose(message MessageHandler) bool {
	kind := message.GetKind()
	return kind >= DiagnoseWaitRequest && kind <= DiagnoseWaitResponse
}
