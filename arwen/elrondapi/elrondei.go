package elrondapi

// // Declare the function signatures (see [cgo](https://golang.org/cmd/cgo/)).
//
// #include <stdlib.h>
// typedef unsigned char uint8_t;
// typedef int int32_t;
// typedef unsigned int uint32_t;
// typedef long long int int64;
//
// extern void getSCAddress(void *context, int32_t resultOffset);
// extern void getOwnerAddress(void *context, int32_t resultOffset);
// extern int32_t getShardOfAddress(void *context, int32_t addressOffset);
// extern int32_t isSmartContract(void *context, int32_t addressOffset);
// extern void getExternalBalance(void *context, int32_t addressOffset, int32_t resultOffset);
// extern int32_t blockHash(void *context, int64 nonce, int32_t resultOffset);
// extern int32_t transferValue(void *context, int32_t dstOffset, int32_t valueOffset, int32_t dataOffset, int32_t length);
// extern int32_t getArgumentLength(void *context, int32_t id);
// extern int32_t getArgument(void *context, int32_t id, int32_t argOffset);
// extern int32_t getFunction(void *context, int32_t functionOffset);
// extern int32_t getNumArguments(void *context);
// extern int32_t storageStore(void *context, int32_t keyOffset, int32_t keyLength , int32_t dataOffset, int32_t dataLength);
// extern int32_t storageLoadLength(void *context, int32_t keyOffset, int32_t keyLength );
// extern int32_t storageLoad(void *context, int32_t keyOffset, int32_t keyLength , int32_t dataOffset);
// extern int32_t storageLoadFromAddress(void *context, int32_t addressOffset, int32_t keyOffset, int32_t keyLength , int32_t dataOffset);
// extern void getCaller(void *context, int32_t resultOffset);
// extern void checkNoPayment(void *context);
// extern int32_t callValue(void *context, int32_t resultOffset);
// extern int32_t getESDTValue(void *context, int32_t resultOffset);
// extern int32_t getESDTTokenName(void *context, int32_t resultOffset);
// extern int32_t getCallValueTokenName(void *context, int32_t callValueOffset, int32_t tokenNameOffset);
// extern void writeLog(void *context, int32_t pointer, int32_t length, int32_t topicPtr, int32_t numTopics);
// extern void returnData(void* context, int32_t dataOffset, int32_t length);
// extern void signalError(void* context, int32_t messageOffset, int32_t messageLength);
// extern int64 getGasLeft(void *context);
//
// extern int32_t executeOnDestContext(void *context, long long gas, int32_t addressOffset, int32_t valueOffset, int32_t functionOffset, int32_t functionLength, int32_t numArguments, int32_t argumentsLengthOffset, int32_t dataOffset);
// extern int32_t executeOnDestContextByCaller(void *context, long long gas, int32_t addressOffset, int32_t valueOffset, int32_t functionOffset, int32_t functionLength, int32_t numArguments, int32_t argumentsLengthOffset, int32_t dataOffset);
// extern int32_t executeOnSameContext(void *context, long long gas, int32_t addressOffset, int32_t valueOffset, int32_t functionOffset, int32_t functionLength, int32_t numArguments, int32_t argumentsLengthOffset, int32_t dataOffset);
// extern int32_t delegateExecution(void *context, long long gas, int32_t addressOffset, int32_t functionOffset, int32_t functionLength, int32_t numArguments, int32_t argumentsLengthOffset, int32_t dataOffset);
// extern int32_t executeReadOnly(void *context, long long gas, int32_t addressOffset, int32_t functionOffset, int32_t functionLength, int32_t numArguments, int32_t argumentsLengthOffset, int32_t dataOffset);
// extern int32_t createContract(void *context, long long gas, int32_t valueOffset, int32_t codeOffset, int32_t codeMetadataOffset, int32_t length, int32_t resultOffset, int32_t numArguments, int32_t argumentsLengthOffset, int32_t dataOffset);
// extern void upgradeContract(void *context, int32_t dstOffset, long long gas, int32_t valueOffset, int32_t codeOffset, int32_t codeMetadataOffset, int32_t length, int32_t numArguments, int32_t argumentsLengthOffset, int32_t dataOffset);
// extern void asyncCall(void *context, int32_t dstOffset, int32_t valueOffset, int32_t dataOffset, int32_t length);
// extern void createAsyncCall(void *context, int32_t groupIDOffset, int32_t groupIDLength, int32_t dstOffset, int32_t valueOffset, int32_t dataOffset, int32_t length, int32_t successCallback, int32_t successLength, int32_t errorCallback, int32_t errorLength, int64 gas);
// extern int32_t setAsyncGroupCallback(void *context, int32_t groupIDOffset, int32_t groupIDLength, int32_t callback, int32_t callbackLength, int32_t data, int32_t dataLength, int32_t gas);
//
// extern int32_t getNumReturnData(void *context);
// extern int32_t getReturnDataSize(void *context, int32_t resultID);
// extern int32_t getReturnData(void *context, int32_t resultID, int32_t dataOffset);
//
// extern int32_t setStorageLock(void *context, int32_t keyOffset, int32_t keyLength, int64 lockTimestamp);
// extern int64 getStorageLock(void *context, int32_t keyOffset, int32_t keyLength);
// extern int32_t isStorageLocked(void *context, int32_t keyOffset, int32_t keyLength);
// extern int32_t clearStorageLock(void *context, int32_t keyOffset, int32_t keyLength);
//
// extern int64 getBlockTimestamp(void *context);
// extern int64 getBlockNonce(void *context);
// extern int64 getBlockRound(void *context);
// extern int64 getBlockEpoch(void *context);
// extern void getBlockRandomSeed(void *context, int32_t resultOffset);
// extern void getStateRootHash(void *context, int32_t resultOffset);
//
// extern int64 getPrevBlockTimestamp(void *context);
// extern int64 getPrevBlockNonce(void *context);
// extern int64 getPrevBlockRound(void *context);
// extern int64 getPrevBlockEpoch(void *context);
// extern void getPrevBlockRandomSeed(void *context, int32_t resultOffset);
// extern void getOriginalTxHash(void *context, int32_t resultOffset);
// extern void getPrevTxHash(void *context, int32_t resultOffset);
import "C"

import (
	"encoding/binary"
	"encoding/hex"
	"errors"
	"fmt"
	"math/big"
	"unsafe"

	"github.com/ElrondNetwork/arwen-wasm-vm/arwen"
	"github.com/ElrondNetwork/arwen-wasm-vm/math"
	"github.com/ElrondNetwork/arwen-wasm-vm/wasmer"
	"github.com/ElrondNetwork/elrond-go/core/vmcommon"
)

// ElrondEIImports creates a new wasmer.Imports populated with the ElrondEI API methods
func ElrondEIImports() (*wasmer.Imports, error) {
	imports := wasmer.NewImports()
	imports = imports.Namespace("env")

	imports, err := imports.Append("getSCAddress", getSCAddress, C.getSCAddress)
	if err != nil {
		return nil, err
	}

	imports, err = imports.Append("getOwnerAddress", getOwnerAddress, C.getOwnerAddress)
	if err != nil {
		return nil, err
	}

	imports, err = imports.Append("getShardOfAddress", getShardOfAddress, C.getShardOfAddress)
	if err != nil {
		return nil, err
	}

	imports, err = imports.Append("isSmartContract", isSmartContract, C.isSmartContract)
	if err != nil {
		return nil, err
	}

	imports, err = imports.Append("getExternalBalance", getExternalBalance, C.getExternalBalance)
	if err != nil {
		return nil, err
	}

	imports, err = imports.Append("getBlockHash", blockHash, C.blockHash)
	if err != nil {
		return nil, err
	}

	imports, err = imports.Append("transferValue", transferValue, C.transferValue)
	if err != nil {
		return nil, err
	}

	imports, err = imports.Append("asyncCall", asyncCall, C.asyncCall)
	if err != nil {
		return nil, err
	}

	imports, err = imports.Append("createAsyncCall", createAsyncCall, C.createAsyncCall)
	if err != nil {
		return nil, err
	}

	imports, err = imports.Append("setAsyncGroupCallback", setAsyncGroupCallback, C.setAsyncGroupCallback)
	if err != nil {
		return nil, err
	}

	imports, err = imports.Append("getArgumentLength", getArgumentLength, C.getArgumentLength)
	if err != nil {
		return nil, err
	}

	imports, err = imports.Append("getArgument", getArgument, C.getArgument)
	if err != nil {
		return nil, err
	}

	imports, err = imports.Append("getFunction", getFunction, C.getFunction)
	if err != nil {
		return nil, err
	}

	imports, err = imports.Append("getNumArguments", getNumArguments, C.getNumArguments)
	if err != nil {
		return nil, err
	}

	imports, err = imports.Append("storageStore", storageStore, C.storageStore)
	if err != nil {
		return nil, err
	}

	imports, err = imports.Append("storageLoadLength", storageLoadLength, C.storageLoadLength)
	if err != nil {
		return nil, err
	}

	imports, err = imports.Append("storageLoad", storageLoad, C.storageLoad)
	if err != nil {
		return nil, err
	}

	imports, err = imports.Append("storageLoadFromAddress", storageLoadFromAddress, C.storageLoadFromAddress)
	if err != nil {
		return nil, err
	}

	imports, err = imports.Append("getStorageLock", getStorageLock, C.getStorageLock)
	if err != nil {
		return nil, err
	}

	imports, err = imports.Append("setStorageLock", setStorageLock, C.setStorageLock)
	if err != nil {
		return nil, err
	}

	imports, err = imports.Append("isStorageLocked", isStorageLocked, C.isStorageLocked)
	if err != nil {
		return nil, err
	}

	imports, err = imports.Append("clearStorageLock", clearStorageLock, C.clearStorageLock)
	if err != nil {
		return nil, err
	}

	imports, err = imports.Append("getCaller", getCaller, C.getCaller)
	if err != nil {
		return nil, err
	}

	imports, err = imports.Append("checkNoPayment", checkNoPayment, C.checkNoPayment)
	if err != nil {
		return nil, err
	}

	imports, err = imports.Append("getCallValue", callValue, C.callValue)
	if err != nil {
		return nil, err
	}

	imports, err = imports.Append("getESDTValue", getESDTValue, C.getESDTValue)
	if err != nil {
		return nil, err
	}

	imports, err = imports.Append("getESDTTokenName", getESDTTokenName, C.getESDTTokenName)
	if err != nil {
		return nil, err
	}

	imports, err = imports.Append("getCallValueTokenName", getCallValueTokenName, C.getCallValueTokenName)
	if err != nil {
		return nil, err
	}

	imports, err = imports.Append("writeLog", writeLog, C.writeLog)
	if err != nil {
		return nil, err
	}

	imports, err = imports.Append("finish", returnData, C.returnData)
	if err != nil {
		return nil, err
	}

	imports, err = imports.Append("signalError", signalError, C.signalError)
	if err != nil {
		return nil, err
	}

	imports, err = imports.Append("getBlockTimestamp", getBlockTimestamp, C.getBlockTimestamp)
	if err != nil {
		return nil, err
	}

	imports, err = imports.Append("getBlockNonce", getBlockNonce, C.getBlockNonce)
	if err != nil {
		return nil, err
	}

	imports, err = imports.Append("getBlockRound", getBlockRound, C.getBlockRound)
	if err != nil {
		return nil, err
	}

	imports, err = imports.Append("getBlockEpoch", getBlockEpoch, C.getBlockEpoch)
	if err != nil {
		return nil, err
	}

	imports, err = imports.Append("getBlockRandomSeed", getBlockRandomSeed, C.getBlockRandomSeed)
	if err != nil {
		return nil, err
	}

	imports, err = imports.Append("getStateRootHash", getStateRootHash, C.getStateRootHash)
	if err != nil {
		return nil, err
	}

	imports, err = imports.Append("getPrevBlockTimestamp", getPrevBlockTimestamp, C.getPrevBlockTimestamp)
	if err != nil {
		return nil, err
	}

	imports, err = imports.Append("getPrevBlockNonce", getPrevBlockNonce, C.getPrevBlockNonce)
	if err != nil {
		return nil, err
	}

	imports, err = imports.Append("getPrevBlockRound", getPrevBlockRound, C.getPrevBlockRound)
	if err != nil {
		return nil, err
	}

	imports, err = imports.Append("getPrevBlockEpoch", getPrevBlockEpoch, C.getPrevBlockEpoch)
	if err != nil {
		return nil, err
	}

	imports, err = imports.Append("getPrevBlockRandomSeed", getPrevBlockRandomSeed, C.getPrevBlockRandomSeed)
	if err != nil {
		return nil, err
	}

	imports, err = imports.Append("getOriginalTxHash", getOriginalTxHash, C.getOriginalTxHash)
	if err != nil {
		return nil, err
	}

	imports, err = imports.Append("getPrevTxHash", getPrevTxHash, C.getPrevTxHash)
	if err != nil {
		return nil, err
	}

	imports, err = imports.Append("getGasLeft", getGasLeft, C.getGasLeft)
	if err != nil {
		return nil, err
	}

	imports, err = imports.Append("executeOnDestContext", executeOnDestContext, C.executeOnDestContext)
	if err != nil {
		return nil, err
	}

	imports, err = imports.Append("executeOnDestContextByCaller", executeOnDestContextByCaller, C.executeOnDestContextByCaller)
	if err != nil {
		return nil, err
	}

	imports, err = imports.Append("executeOnSameContext", executeOnSameContext, C.executeOnSameContext)
	if err != nil {
		return nil, err
	}

	imports, err = imports.Append("delegateExecution", delegateExecution, C.delegateExecution)
	if err != nil {
		return nil, err
	}

	imports, err = imports.Append("createContract", createContract, C.createContract)
	if err != nil {
		return nil, err
	}

	imports, err = imports.Append("upgradeContract", upgradeContract, C.upgradeContract)
	if err != nil {
		return nil, err
	}

	imports, err = imports.Append("executeReadOnly", executeReadOnly, C.executeReadOnly)
	if err != nil {
		return nil, err
	}

	imports, err = imports.Append("getNumReturnData", getNumReturnData, C.getNumReturnData)
	if err != nil {
		return nil, err
	}

	imports, err = imports.Append("getReturnDataSize", getReturnDataSize, C.getReturnDataSize)
	if err != nil {
		return nil, err
	}

	imports, err = imports.Append("getReturnData", getReturnData, C.getReturnData)
	if err != nil {
		return nil, err
	}

	return imports, nil
}

//export getGasLeft
func getGasLeft(context unsafe.Pointer) int64 {
	metering := arwen.GetMeteringContext(context)

	gasToUse := metering.GasSchedule().ElrondAPICost.GetGasLeft
	metering.UseGas(gasToUse)

	return int64(metering.GasLeft())
}

//export getSCAddress
func getSCAddress(context unsafe.Pointer, resultOffset int32) {
	runtime := arwen.GetRuntimeContext(context)
	metering := arwen.GetMeteringContext(context)

	gasToUse := metering.GasSchedule().ElrondAPICost.GetSCAddress
	metering.UseGas(gasToUse)

	owner := runtime.GetSCAddress()
	err := runtime.MemStore(resultOffset, owner)
	if arwen.WithFault(err, context, runtime.ElrondAPIErrorShouldFailExecution()) {
		return
	}
}

//export getOwnerAddress
func getOwnerAddress(context unsafe.Pointer, resultOffset int32) {
	blockchain := arwen.GetBlockchainContext(context)
	runtime := arwen.GetRuntimeContext(context)
	metering := arwen.GetMeteringContext(context)

	gasToUse := metering.GasSchedule().ElrondAPICost.GetOwnerAddress
	metering.UseGas(gasToUse)

	owner, err := blockchain.GetOwnerAddress()
	if arwen.WithFault(err, context, runtime.ElrondAPIErrorShouldFailExecution()) {
		return
	}

	err = runtime.MemStore(resultOffset, owner)
	if arwen.WithFault(err, context, runtime.ElrondAPIErrorShouldFailExecution()) {
		return
	}
}

//export getShardOfAddress
func getShardOfAddress(context unsafe.Pointer, addressOffset int32) int32 {
	blockchain := arwen.GetBlockchainContext(context)
	runtime := arwen.GetRuntimeContext(context)
	metering := arwen.GetMeteringContext(context)

	gasToUse := metering.GasSchedule().ElrondAPICost.GetShardOfAddress
	metering.UseGas(gasToUse)

	address, err := runtime.MemLoad(addressOffset, arwen.AddressLen)
	if arwen.WithFault(err, context, runtime.ElrondAPIErrorShouldFailExecution()) {
		return 0
	}

	return int32(blockchain.GetShardOfAddress(address))
}

//export isSmartContract
func isSmartContract(context unsafe.Pointer, addressOffset int32) int32 {
	blockchain := arwen.GetBlockchainContext(context)
	runtime := arwen.GetRuntimeContext(context)
	metering := arwen.GetMeteringContext(context)

	gasToUse := metering.GasSchedule().ElrondAPICost.IsSmartContract
	metering.UseGas(gasToUse)

	address, err := runtime.MemLoad(addressOffset, arwen.AddressLen)
	if arwen.WithFault(err, context, runtime.ElrondAPIErrorShouldFailExecution()) {
		return 0
	}

	isSmartContract := blockchain.IsSmartContract(address)
	return int32(arwen.BooleanToInt(isSmartContract))
}

//export signalError
func signalError(context unsafe.Pointer, messageOffset int32, messageLength int32) {
	runtime := arwen.GetRuntimeContext(context)
	metering := arwen.GetMeteringContext(context)

	gasToUse := metering.GasSchedule().ElrondAPICost.SignalError
	metering.UseGas(gasToUse)

	message, err := runtime.MemLoad(messageOffset, messageLength)
	if arwen.WithFault(err, context, runtime.ElrondAPIErrorShouldFailExecution()) {
		return
	}
	runtime.SignalUserError(string(message))
}

//export getExternalBalance
func getExternalBalance(context unsafe.Pointer, addressOffset int32, resultOffset int32) {
	blockchain := arwen.GetBlockchainContext(context)
	runtime := arwen.GetRuntimeContext(context)
	metering := arwen.GetMeteringContext(context)

	gasToUse := metering.GasSchedule().ElrondAPICost.GetExternalBalance
	metering.UseGas(gasToUse)

	address, err := runtime.MemLoad(addressOffset, arwen.AddressLen)
	if arwen.WithFault(err, context, runtime.ElrondAPIErrorShouldFailExecution()) {
		return
	}

	balance := blockchain.GetBalance(address)

	err = runtime.MemStore(resultOffset, balance)
	if arwen.WithFault(err, context, runtime.ElrondAPIErrorShouldFailExecution()) {
		return
	}
}

//export blockHash
func blockHash(context unsafe.Pointer, nonce int64, resultOffset int32) int32 {
	blockchain := arwen.GetBlockchainContext(context)
	runtime := arwen.GetRuntimeContext(context)
	metering := arwen.GetMeteringContext(context)

	gasToUse := metering.GasSchedule().ElrondAPICost.GetBlockHash
	metering.UseGas(gasToUse)

	hash := blockchain.BlockHash(uint64(nonce))
	err := runtime.MemStore(resultOffset, hash)
	if arwen.WithFault(err, context, runtime.ElrondAPIErrorShouldFailExecution()) {
		return 1
	}

	return 0
}

//export transferValue
func transferValue(context unsafe.Pointer, destOffset int32, valueOffset int32, dataOffset int32, length int32) int32 {
	host := arwen.GetVMContext(context)
	runtime := host.Runtime()
	metering := host.Metering()
	output := host.Output()

	gasToUse := metering.GasSchedule().ElrondAPICost.TransferValue
	metering.UseGas(gasToUse)

	send := runtime.GetSCAddress()
	dest, err := runtime.MemLoad(destOffset, arwen.AddressLen)
	if arwen.WithFault(err, context, runtime.ElrondAPIErrorShouldFailExecution()) {
		return 1
	}

	valueBytes, err := runtime.MemLoad(valueOffset, arwen.BalanceLen)
	if arwen.WithFault(err, context, runtime.ElrondAPIErrorShouldFailExecution()) {
		return 1
	}

	gasToUse = math.MulUint64(metering.GasSchedule().BaseOperationCost.PersistPerByte, uint64(length))
	metering.UseGas(gasToUse)

	data, err := runtime.MemLoad(dataOffset, length)
	if arwen.WithFault(err, context, runtime.ElrondAPIErrorShouldFailExecution()) {
		return 1
	}

	// TODO write test for this, after removing vmContextMap
	functionName, _, err := host.CallArgsParser().ParseData(string(data))
	if host.IsBuiltinFunctionName(functionName) {
		return 1
	}

	err = output.Transfer(dest, send, 0, 0, big.NewInt(0).SetBytes(valueBytes), data, vmcommon.DirectCall)
	if arwen.WithFault(err, context, runtime.ElrondAPIErrorShouldFailExecution()) {
		return 1
	}

	return 0
}

//export createAsyncCall
func createAsyncCall(context unsafe.Pointer,
	groupIDOffset int32,
	groupIDLength int32,
	destOffset int32,
	valueOffset int32,
	dataOffset int32,
	dataLength int32,
	successOffset int32,
	successLength int32,
	errorOffset int32,
	errorLength int32,
	gas int64,
) {
	// TODO Create new API cost for this method

	host := arwen.GetVMContext(context)
	runtime := host.Runtime()
	async := host.Async()

	groupIDBytes, err := runtime.MemLoad(groupIDOffset, groupIDLength)
	if arwen.WithFault(err, context, runtime.ElrondAPIErrorShouldFailExecution()) {
		return
	}

	groupID := string(groupIDBytes)
	if groupID == arwen.LegacyAsyncCallGroupID {
		err = arwen.ErrInvalidAsyncCallGroupID
		arwen.WithFault(err, context, runtime.ElrondAPIErrorShouldFailExecution())
		return
	}

	calledSCAddress, err := runtime.MemLoad(destOffset, arwen.AddressLen)
	if arwen.WithFault(err, context, runtime.ElrondAPIErrorShouldFailExecution()) {
		return
	}

	value, err := runtime.MemLoad(valueOffset, arwen.BalanceLen)
	if arwen.WithFault(err, context, runtime.ElrondAPIErrorShouldFailExecution()) {
		return
	}

	data, err := runtime.MemLoad(dataOffset, dataLength)
	if arwen.WithFault(err, context, runtime.ElrondAPIErrorShouldFailExecution()) {
		return
	}

	successFunc, err := runtime.MemLoad(successOffset, successLength)
	if arwen.WithFault(err, context, runtime.ElrondAPIErrorShouldFailExecution()) {
		return
	}

	errorFunc, err := runtime.MemLoad(errorOffset, errorLength)
	if arwen.WithFault(err, context, runtime.ElrondAPIErrorShouldFailExecution()) {
		return
	}

	err = async.RegisterAsyncCall(groupID, &arwen.AsyncCall{
		Status:          arwen.AsyncCallPending,
		Destination:     calledSCAddress,
		Data:            data,
		ValueBytes:      value,
		GasLimit:        uint64(gas),
		SuccessCallback: string(successFunc),
		ErrorCallback:   string(errorFunc),
	})
	if arwen.WithFault(err, context, runtime.ElrondAPIErrorShouldFailExecution()) {
		return
	}
}

//export setAsyncGroupCallback
func setAsyncGroupCallback(context unsafe.Pointer,
	groupIDOffset int32,
	groupIDLength int32,
	callback int32,
	callbackLength int32,
	data int32,
	dataLength int32,
	gas int32,
) int32 {
	// TODO Create new API cost for this method

	host := arwen.GetVMContext(context)
	runtime := host.Runtime()
	async := host.Async()

	groupIDBytes, err := runtime.MemLoad(groupIDOffset, groupIDLength)
	if arwen.WithFault(err, context, runtime.ElrondAPIErrorShouldFailExecution()) {
		return -1
	}

	groupID := string(groupIDBytes)
	if groupID == arwen.LegacyAsyncCallGroupID {
		err = arwen.ErrInvalidAsyncCallGroupID
		arwen.WithFault(err, context, runtime.ElrondAPIErrorShouldFailExecution())
		return -1
	}

	callbackNameBytes, err := runtime.MemLoad(callback, callbackLength)
	if arwen.WithFault(err, context, runtime.ElrondAPIErrorShouldFailExecution()) {
		return -1
	}

	dataBytes, err := runtime.MemLoad(data, dataLength)
	if arwen.WithFault(err, context, runtime.ElrondAPIErrorShouldFailExecution()) {
		return -1
	}

	_, exists := async.GetCallGroup(groupID)
	if !exists {
		err = arwen.ErrAsyncCallGroupDoesNotExist
		arwen.WithFault(err, context, runtime.ElrondAPIErrorShouldFailExecution())
		return -1
	}

	callbackName := string(callbackNameBytes)
	err = async.SetGroupCallback(groupID, callbackName, dataBytes, uint64(gas))
	if arwen.WithFault(err, context, runtime.ElrondAPIErrorShouldFailExecution()) {
		return -1
	}

	return 0
}

//export upgradeContract
func upgradeContract(
	context unsafe.Pointer,
	destOffset int32,
	gasLimit int64,
	valueOffset int32,
	codeOffset int32,
	codeMetadataOffset int32,
	length int32,
	numArguments int32,
	argumentsLengthOffset int32,
	dataOffset int32,
) {
	host := arwen.GetVMContext(context)
	runtime := host.Runtime()
	metering := host.Metering()

	gasToUse := metering.GasSchedule().ElrondAPICost.CreateContract
	metering.UseGas(gasToUse)

	// value, err := runtime.MemLoad(valueOffset, arwen.BalanceLen)
	// if arwen.WithFault(err, context, runtime.ElrondAPIErrorShouldFailExecution()) {
	// 	return
	// }

	code, err := runtime.MemLoad(codeOffset, length)
	if arwen.WithFault(err, context, runtime.ElrondAPIErrorShouldFailExecution()) {
		return
	}

	codeMetadata, err := runtime.MemLoad(codeMetadataOffset, arwen.CodeMetadataLen)
	if arwen.WithFault(err, context, runtime.ElrondAPIErrorShouldFailExecution()) {
		return
	}

	_, data, actualLen, err := getArgumentsFromMemory(
		host,
		0,
		0,
		numArguments,
		argumentsLengthOffset,
		dataOffset,
	)

	gasToUse = math.MulUint64(metering.GasSchedule().BaseOperationCost.DataCopyPerByte, uint64(actualLen))
	metering.UseGas(gasToUse)

	if arwen.WithFault(err, context, runtime.ElrondAPIErrorShouldFailExecution()) {
		return
	}

	gasSchedule := metering.GasSchedule()
	gasToUse = gasSchedule.ElrondAPICost.AsyncCallStep
	metering.UseGas(gasToUse)

	// calledSCAddress, err := runtime.MemLoad(destOffset, arwen.AddressLen)
	// if arwen.WithFault(err, context, runtime.ElrondAPIErrorShouldFailExecution()) {
	// 	return
	// }

	gasToUse = math.MulUint64(gasSchedule.BaseOperationCost.DataCopyPerByte, uint64(length))
	metering.UseGas(gasToUse)

	minAsyncCallCost := math.AddUint64(
		math.MulUint64(2, gasSchedule.ElrondAPICost.AsyncCallStep),
		gasSchedule.ElrondAPICost.AsyncCallbackGasLock)
	if uint64(gasLimit) < minAsyncCallCost {
		runtime.SetRuntimeBreakpointValue(arwen.BreakpointOutOfGas)
		return
	}

	// TODO replace the following code with runtime.AddAsyncCall()

	// Set up the async call as if it is not known whether the called SC
	// is in the same shard with the caller or not. This will be later resolved
	// in the handler for BreakpointAsyncCall.
	codeEncoded := hex.EncodeToString(code)
	codeMetadataEncoded := hex.EncodeToString(codeMetadata)
	finalData := arwen.UpgradeFunctionName + "@" + codeEncoded + "@" + codeMetadataEncoded
	for _, arg := range data {
		finalData += "@" + string(arg)
	}

	// runtime.SetAsyncCallInfo(&arwen.AsyncCallInfo{
	// 	Destination: calledSCAddress,
	// 	Data:        []byte(finalData),
	// 	GasLimit:    uint64(gasLimit),
	// 	ValueBytes:  value,
	// })

	// Instruct Wasmer to interrupt the execution of the caller SC.
	runtime.SetRuntimeBreakpointValue(arwen.BreakpointAsyncCall)
}

//export asyncCall
func asyncCall(context unsafe.Pointer, destOffset int32, valueOffset int32, dataOffset int32, length int32) {
	host := arwen.GetVMContext(context)
	runtime := host.Runtime()
	async := host.Async()
	metering := host.Metering()

	gasSchedule := metering.GasSchedule()
	gasToUse := gasSchedule.ElrondAPICost.AsyncCallStep
	metering.UseGas(gasToUse)

	calledSCAddress, err := runtime.MemLoad(destOffset, arwen.AddressLen)
	if arwen.WithFault(err, context, runtime.ElrondAPIErrorShouldFailExecution()) {
		return
	}

	value, err := runtime.MemLoad(valueOffset, arwen.BalanceLen)
	if arwen.WithFault(err, context, runtime.ElrondAPIErrorShouldFailExecution()) {
		return
	}

	gasToUse = math.MulUint64(gasSchedule.BaseOperationCost.DataCopyPerByte, uint64(length))
	metering.UseGas(gasToUse)

	data, err := runtime.MemLoad(dataOffset, length)
	if arwen.WithFault(err, context, runtime.ElrondAPIErrorShouldFailExecution()) {
		return
	}

	err = async.RegisterLegacyAsyncCall(calledSCAddress, data, value)
	if errors.Is(err, arwen.ErrNotEnoughGas) {
		runtime.SetRuntimeBreakpointValue(arwen.BreakpointOutOfGas)
		return
	}
	if arwen.WithFault(err, context, runtime.ElrondAPIErrorShouldFailExecution()) {
		return
	}
}

//export getArgumentLength
func getArgumentLength(context unsafe.Pointer, id int32) int32 {
	runtime := arwen.GetRuntimeContext(context)
	metering := arwen.GetMeteringContext(context)

	gasToUse := metering.GasSchedule().ElrondAPICost.GetArgument
	metering.UseGas(gasToUse)

	args := runtime.Arguments()
	if id < 0 || int32(len(args)) <= id {
		return -1
	}

	return int32(len(args[id]))
}

//export getArgument
func getArgument(context unsafe.Pointer, id int32, argOffset int32) int32 {
	runtime := arwen.GetRuntimeContext(context)
	metering := arwen.GetMeteringContext(context)

	gasToUse := metering.GasSchedule().ElrondAPICost.GetArgument
	metering.UseGas(gasToUse)

	args := runtime.Arguments()
	if id < 0 || int32(len(args)) <= id {
		return -1
	}

	err := runtime.MemStore(argOffset, args[id])
	if arwen.WithFault(err, context, runtime.ElrondAPIErrorShouldFailExecution()) {
		return -1
	}

	return int32(len(args[id]))
}

//export getFunction
func getFunction(context unsafe.Pointer, functionOffset int32) int32 {
	runtime := arwen.GetRuntimeContext(context)
	metering := arwen.GetMeteringContext(context)

	gasToUse := metering.GasSchedule().ElrondAPICost.GetFunction
	metering.UseGas(gasToUse)

	function := runtime.Function()
	err := runtime.MemStore(functionOffset, []byte(function))
	if arwen.WithFault(err, context, runtime.ElrondAPIErrorShouldFailExecution()) {
		return -1
	}

	return int32(len(function))
}

//export getNumArguments
func getNumArguments(context unsafe.Pointer) int32 {
	runtime := arwen.GetRuntimeContext(context)
	metering := arwen.GetMeteringContext(context)

	gasToUse := metering.GasSchedule().ElrondAPICost.GetNumArguments
	metering.UseGas(gasToUse)

	args := runtime.Arguments()
	return int32(len(args))
}

//export storageStore
func storageStore(context unsafe.Pointer, keyOffset int32, keyLength int32, dataOffset int32, dataLength int32) int32 {
	runtime := arwen.GetRuntimeContext(context)
	storage := arwen.GetStorageContext(context)
	metering := arwen.GetMeteringContext(context)

	gasToUse := metering.GasSchedule().ElrondAPICost.StorageStore
	metering.UseGas(gasToUse)

	key, err := runtime.MemLoad(keyOffset, keyLength)
	if arwen.WithFault(err, context, runtime.ElrondAPIErrorShouldFailExecution()) {
		return -1
	}

	data, err := runtime.MemLoad(dataOffset, dataLength)
	if arwen.WithFault(err, context, runtime.ElrondAPIErrorShouldFailExecution()) {
		return -1
	}

	storageStatus, err := storage.SetStorage(key, data)
	if arwen.WithFault(err, context, runtime.ElrondAPIErrorShouldFailExecution()) {
		return -1
	}

	return int32(storageStatus)
}

//export storageLoadLength
func storageLoadLength(context unsafe.Pointer, keyOffset int32, keyLength int32) int32 {
	runtime := arwen.GetRuntimeContext(context)
	storage := arwen.GetStorageContext(context)
	metering := arwen.GetMeteringContext(context)

	gasToUse := metering.GasSchedule().ElrondAPICost.StorageLoad
	metering.UseGas(gasToUse)

	key, err := runtime.MemLoad(keyOffset, keyLength)
	if arwen.WithFault(err, context, runtime.ElrondAPIErrorShouldFailExecution()) {
		return -1
	}

	data := storage.GetStorageUnmetered(key)

	return int32(len(data))
}

//export storageLoadFromAddress
func storageLoadFromAddress(context unsafe.Pointer, addressOffset int32, keyOffset int32, keyLength int32, dataOffset int32) int32 {
	runtime := arwen.GetRuntimeContext(context)
	storage := arwen.GetStorageContext(context)
	metering := arwen.GetMeteringContext(context)

	gasToUse := metering.GasSchedule().ElrondAPICost.StorageLoad
	metering.UseGas(gasToUse)

	key, err := runtime.MemLoad(keyOffset, keyLength)
	if arwen.WithFault(err, context, runtime.ElrondAPIErrorShouldFailExecution()) {
		return -1
	}

	address, err := runtime.MemLoad(addressOffset, arwen.AddressLen)
	if arwen.WithFault(err, context, runtime.ElrondAPIErrorShouldFailExecution()) {
		return -1
	}

	data := storage.GetStorageFromAddress(address, key)

	err = runtime.MemStore(dataOffset, data)
	if arwen.WithFault(err, context, runtime.ElrondAPIErrorShouldFailExecution()) {
		return -1
	}

	return int32(len(data))
}

//export storageLoad
func storageLoad(context unsafe.Pointer, keyOffset int32, keyLength int32, dataOffset int32) int32 {
	runtime := arwen.GetRuntimeContext(context)
	storage := arwen.GetStorageContext(context)
	metering := arwen.GetMeteringContext(context)

	gasToUse := metering.GasSchedule().ElrondAPICost.StorageLoad
	metering.UseGas(gasToUse)

	key, err := runtime.MemLoad(keyOffset, keyLength)
	if arwen.WithFault(err, context, runtime.ElrondAPIErrorShouldFailExecution()) {
		return -1
	}

	data := storage.GetStorage(key)

	err = runtime.MemStore(dataOffset, data)
	if arwen.WithFault(err, context, runtime.ElrondAPIErrorShouldFailExecution()) {
		return -1
	}

	return int32(len(data))
}

//export setStorageLock
func setStorageLock(context unsafe.Pointer, keyOffset int32, keyLength int32, lockTimestamp int64) int32 {
	runtime := arwen.GetRuntimeContext(context)
	storage := arwen.GetStorageContext(context)
	metering := arwen.GetMeteringContext(context)

	gasToUse := metering.GasSchedule().ElrondAPICost.Int64StorageStore
	metering.UseGas(gasToUse)

	key, err := runtime.MemLoad(keyOffset, keyLength)
	if arwen.WithFault(err, context, runtime.ElrondAPIErrorShouldFailExecution()) {
		return -1
	}

	timeLockKey := arwen.CustomStorageKey(arwen.TimeLockKeyPrefix, key)
	bigTimestamp := big.NewInt(0).SetInt64(lockTimestamp)
	storageStatus, err := storage.SetProtectedStorage(timeLockKey, bigTimestamp.Bytes())
	if arwen.WithFault(err, context, runtime.ElrondAPIErrorShouldFailExecution()) {
		return -1
	}
	return int32(storageStatus)
}

//export getStorageLock
func getStorageLock(context unsafe.Pointer, keyOffset int32, keyLength int32) int64 {
	runtime := arwen.GetRuntimeContext(context)
	metering := arwen.GetMeteringContext(context)
	storage := arwen.GetStorageContext(context)

	gasToUse := metering.GasSchedule().ElrondAPICost.StorageLoad
	metering.UseGas(gasToUse)

	key, err := runtime.MemLoad(keyOffset, keyLength)
	if arwen.WithFault(err, context, runtime.ElrondAPIErrorShouldFailExecution()) {
		return -1
	}

	timeLockKey := arwen.CustomStorageKey(arwen.TimeLockKeyPrefix, key)
	data := storage.GetStorage(timeLockKey)
	timeLock := big.NewInt(0).SetBytes(data).Int64()

	return timeLock
}

//export isStorageLocked
func isStorageLocked(context unsafe.Pointer, keyOffset int32, keyLength int32) int32 {
	timeLock := getStorageLock(context, keyOffset, keyLength)
	if timeLock < 0 {
		return -1
	}

	currentTimestamp := int64(getBlockTimestamp(context))
	if timeLock <= currentTimestamp {
		return 0
	}

	return 1
}

//export clearStorageLock
func clearStorageLock(context unsafe.Pointer, keyOffset int32, keyLength int32) int32 {
	return setStorageLock(context, keyOffset, keyLength, 0)
}

//export getCaller
func getCaller(context unsafe.Pointer, resultOffset int32) {
	runtime := arwen.GetRuntimeContext(context)
	metering := arwen.GetMeteringContext(context)

	gasToUse := metering.GasSchedule().ElrondAPICost.GetCaller
	metering.UseGas(gasToUse)

	caller := runtime.GetVMInput().CallerAddr

	err := runtime.MemStore(resultOffset, caller)
	if arwen.WithFault(err, context, runtime.ElrondAPIErrorShouldFailExecution()) {
		return
	}
}

//export checkNoPayment
func checkNoPayment(context unsafe.Pointer) {
	runtime := arwen.GetRuntimeContext(context)
	metering := arwen.GetMeteringContext(context)

	gasToUse := metering.GasSchedule().ElrondAPICost.GetCallValue
	metering.UseGas(gasToUse)

	vmInput := runtime.GetVMInput()
	if vmInput.CallValue.Sign() > 0 {
		runtime := arwen.GetRuntimeContext(context)
		arwen.WithFault(arwen.ErrNonPayableFunctionEgld, context, runtime.ElrondAPIErrorShouldFailExecution())
		return
	}
	if vmInput.ESDTValue != nil && vmInput.ESDTValue.Sign() > 0 {
		runtime := arwen.GetRuntimeContext(context)
		arwen.WithFault(arwen.ErrNonPayableFunctionEsdt, context, runtime.ElrondAPIErrorShouldFailExecution())
		return
	}
}

//export callValue
func callValue(context unsafe.Pointer, resultOffset int32) int32 {
	runtime := arwen.GetRuntimeContext(context)
	metering := arwen.GetMeteringContext(context)

	gasToUse := metering.GasSchedule().ElrondAPICost.GetCallValue
	metering.UseGas(gasToUse)

	value := runtime.GetVMInput().CallValue.Bytes()
	value = arwen.PadBytesLeft(value, arwen.BalanceLen)

	err := runtime.MemStore(resultOffset, value)
	if arwen.WithFault(err, context, runtime.ElrondAPIErrorShouldFailExecution()) {
		return -1
	}

	return int32(len(value))
}

//export getESDTValue
func getESDTValue(context unsafe.Pointer, resultOffset int32) int32 {
	runtime := arwen.GetRuntimeContext(context)
	metering := arwen.GetMeteringContext(context)

	gasToUse := metering.GasSchedule().ElrondAPICost.GetCallValue
	metering.UseGas(gasToUse)

	var value []byte

	esdtValue := runtime.GetVMInput().ESDTValue
	if esdtValue != nil {
		value = esdtValue.Bytes()
		value = arwen.PadBytesLeft(value, arwen.BalanceLen)
	}

	err := runtime.MemStore(resultOffset, value)
	if arwen.WithFault(err, context, runtime.ElrondAPIErrorShouldFailExecution()) {
		return -1
	}

	return int32(len(value))
}

//export getESDTTokenName
func getESDTTokenName(context unsafe.Pointer, resultOffset int32) int32 {
	runtime := arwen.GetRuntimeContext(context)
	metering := arwen.GetMeteringContext(context)

	gasToUse := metering.GasSchedule().ElrondAPICost.GetCallValue
	metering.UseGas(gasToUse)

	tokenName := runtime.GetVMInput().ESDTTokenName

	err := runtime.MemStore(resultOffset, tokenName)
	if arwen.WithFault(err, context, runtime.ElrondAPIErrorShouldFailExecution()) {
		return -1
	}

	return int32(len(tokenName))
}

//export getCallValueTokenName
func getCallValueTokenName(context unsafe.Pointer, callValueOffset int32, tokenNameOffset int32) int32 {
	runtime := arwen.GetRuntimeContext(context)
	metering := arwen.GetMeteringContext(context)

	gasToUse := metering.GasSchedule().ElrondAPICost.GetCallValue
	metering.UseGas(gasToUse)

	callValue := runtime.GetVMInput().CallValue.Bytes()
	tokenName := []byte("eGLD")
	if len(runtime.GetVMInput().ESDTTokenName) > 0 {
		tokenName = runtime.GetVMInput().ESDTTokenName
		callValue = runtime.GetVMInput().ESDTValue.Bytes()
	}
	callValue = arwen.PadBytesLeft(callValue, arwen.BalanceLen)

	err := runtime.MemStore(tokenNameOffset, tokenName)
	if arwen.WithFault(err, context, runtime.ElrondAPIErrorShouldFailExecution()) {
		return -1
	}

	err = runtime.MemStore(callValueOffset, callValue)
	if arwen.WithFault(err, context, runtime.ElrondAPIErrorShouldFailExecution()) {
		return -1
	}

	return int32(len(tokenName))
}

//export writeLog
func writeLog(context unsafe.Pointer, pointer int32, length int32, topicPtr int32, numTopics int32) {
	runtime := arwen.GetRuntimeContext(context)
	output := arwen.GetOutputContext(context)
	metering := arwen.GetMeteringContext(context)

	gasToUse := metering.GasSchedule().ElrondAPICost.Log
	gas := math.MulUint64(metering.GasSchedule().BaseOperationCost.PersistPerByte, uint64(numTopics*arwen.HashLen+length))
	gasToUse = math.AddUint64(gasToUse, gas)
	metering.UseGas(gasToUse)

	log, err := runtime.MemLoad(pointer, length)
	if arwen.WithFault(err, context, runtime.ElrondAPIErrorShouldFailExecution()) {
		return
	}

	topics, err := arwen.GuardedMakeByteSlice2D(numTopics)
	if arwen.WithFault(err, context, runtime.ElrondAPIErrorShouldFailExecution()) {
		return
	}

	for i := int32(0); i < numTopics; i++ {
		topics[i], err = runtime.MemLoad(topicPtr+i*arwen.HashLen, arwen.HashLen)
		if arwen.WithFault(err, context, runtime.ElrondAPIErrorShouldFailExecution()) {
			return
		}
	}

	output.WriteLog(runtime.GetSCAddress(), topics, log)
}

//export getBlockTimestamp
func getBlockTimestamp(context unsafe.Pointer) int64 {
	blockchain := arwen.GetBlockchainContext(context)
	metering := arwen.GetMeteringContext(context)

	gasToUse := metering.GasSchedule().ElrondAPICost.GetBlockTimeStamp
	metering.UseGas(gasToUse)

	return int64(blockchain.CurrentTimeStamp())
}

//export getBlockNonce
func getBlockNonce(context unsafe.Pointer) int64 {
	blockchain := arwen.GetBlockchainContext(context)
	metering := arwen.GetMeteringContext(context)

	gasToUse := metering.GasSchedule().ElrondAPICost.GetBlockNonce
	metering.UseGas(gasToUse)

	return int64(blockchain.CurrentNonce())
}

//export getBlockRound
func getBlockRound(context unsafe.Pointer) int64 {
	blockchain := arwen.GetBlockchainContext(context)
	metering := arwen.GetMeteringContext(context)

	gasToUse := metering.GasSchedule().ElrondAPICost.GetBlockRound
	metering.UseGas(gasToUse)

	return int64(blockchain.CurrentRound())
}

//export getBlockEpoch
func getBlockEpoch(context unsafe.Pointer) int64 {
	blockchain := arwen.GetBlockchainContext(context)
	metering := arwen.GetMeteringContext(context)

	gasToUse := metering.GasSchedule().ElrondAPICost.GetBlockEpoch
	metering.UseGas(gasToUse)

	return int64(blockchain.CurrentEpoch())
}

//export getBlockRandomSeed
func getBlockRandomSeed(context unsafe.Pointer, pointer int32) {
	runtime := arwen.GetRuntimeContext(context)
	blockchain := arwen.GetBlockchainContext(context)
	metering := arwen.GetMeteringContext(context)

	gasToUse := metering.GasSchedule().ElrondAPICost.GetBlockRandomSeed
	metering.UseGas(gasToUse)

	randomSeed := blockchain.CurrentRandomSeed()
	err := runtime.MemStore(pointer, randomSeed)
	arwen.WithFault(err, context, runtime.ElrondAPIErrorShouldFailExecution())
}

//export getStateRootHash
func getStateRootHash(context unsafe.Pointer, pointer int32) {
	runtime := arwen.GetRuntimeContext(context)
	blockchain := arwen.GetBlockchainContext(context)
	metering := arwen.GetMeteringContext(context)

	gasToUse := metering.GasSchedule().ElrondAPICost.GetStateRootHash
	metering.UseGas(gasToUse)

	stateRootHash := blockchain.GetStateRootHash()
	err := runtime.MemStore(pointer, stateRootHash)
	arwen.WithFault(err, context, runtime.ElrondAPIErrorShouldFailExecution())
}

//export getPrevBlockTimestamp
func getPrevBlockTimestamp(context unsafe.Pointer) int64 {
	blockchain := arwen.GetBlockchainContext(context)
	metering := arwen.GetMeteringContext(context)

	gasToUse := metering.GasSchedule().ElrondAPICost.GetBlockTimeStamp
	metering.UseGas(gasToUse)

	return int64(blockchain.LastTimeStamp())
}

//export getPrevBlockNonce
func getPrevBlockNonce(context unsafe.Pointer) int64 {
	blockchain := arwen.GetBlockchainContext(context)
	metering := arwen.GetMeteringContext(context)

	gasToUse := metering.GasSchedule().ElrondAPICost.GetBlockNonce
	metering.UseGas(gasToUse)

	return int64(blockchain.LastNonce())
}

//export getPrevBlockRound
func getPrevBlockRound(context unsafe.Pointer) int64 {
	blockchain := arwen.GetBlockchainContext(context)
	metering := arwen.GetMeteringContext(context)

	gasToUse := metering.GasSchedule().ElrondAPICost.GetBlockRound
	metering.UseGas(gasToUse)

	return int64(blockchain.LastRound())
}

//export getPrevBlockEpoch
func getPrevBlockEpoch(context unsafe.Pointer) int64 {
	blockchain := arwen.GetBlockchainContext(context)
	metering := arwen.GetMeteringContext(context)

	gasToUse := metering.GasSchedule().ElrondAPICost.GetBlockEpoch
	metering.UseGas(gasToUse)

	return int64(blockchain.LastEpoch())
}

//export getPrevBlockRandomSeed
func getPrevBlockRandomSeed(context unsafe.Pointer, pointer int32) {
	runtime := arwen.GetRuntimeContext(context)
	blockchain := arwen.GetBlockchainContext(context)
	metering := arwen.GetMeteringContext(context)

	gasToUse := metering.GasSchedule().ElrondAPICost.GetBlockRandomSeed
	metering.UseGas(gasToUse)

	randomSeed := blockchain.LastRandomSeed()
	err := runtime.MemStore(pointer, randomSeed)
	arwen.WithFault(err, context, runtime.ElrondAPIErrorShouldFailExecution())
}

//export returnData
func returnData(context unsafe.Pointer, pointer int32, length int32) {
	runtime := arwen.GetRuntimeContext(context)
	output := arwen.GetOutputContext(context)
	metering := arwen.GetMeteringContext(context)

	gasToUse := metering.GasSchedule().ElrondAPICost.Finish
	gas := math.MulUint64(metering.GasSchedule().BaseOperationCost.PersistPerByte, uint64(length))
	gasToUse = math.AddUint64(gasToUse, gas)
	metering.UseGas(gasToUse)

	data, err := runtime.MemLoad(pointer, length)
	if arwen.WithFault(err, context, runtime.ElrondAPIErrorShouldFailExecution()) {
		return
	}

	output.Finish(data)
}

//export executeOnSameContext
func executeOnSameContext(
	context unsafe.Pointer,
	gasLimit int64,
	addressOffset int32,
	valueOffset int32,
	functionOffset int32,
	functionLength int32,
	numArguments int32,
	argumentsLengthOffset int32,
	dataOffset int32,
) int32 {
	host := arwen.GetVMContext(context)
	runtime := host.Runtime()
	metering := host.Metering()

	gasToUse := metering.GasSchedule().ElrondAPICost.ExecuteOnSameContext
	metering.UseGas(gasToUse)

	value, err := runtime.MemLoad(valueOffset, arwen.BalanceLen)
	if arwen.WithFault(err, context, runtime.ElrondSyncExecAPIErrorShouldFailExecution()) {
		return 1
	}

	sender := runtime.GetSCAddress()
	bigIntVal := big.NewInt(0).SetBytes(value)
	contractCallInput, err := prepareIndirectContractCallInput(
		host,
		sender,
		bigIntVal,
		gasLimit,
		addressOffset,
		functionOffset,
		functionLength,
		numArguments,
		argumentsLengthOffset,
		dataOffset,
	)
	if arwen.WithFault(err, context, runtime.ElrondSyncExecAPIErrorShouldFailExecution()) {
		return 1
	}

	err = host.ExecuteOnSameContext(contractCallInput)
	if arwen.WithFault(err, context, runtime.ElrondSyncExecAPIErrorShouldFailExecution()) {
		return 1
	}

	return 0
}

//export executeOnDestContext
func executeOnDestContext(
	context unsafe.Pointer,
	gasLimit int64,
	addressOffset int32,
	valueOffset int32,
	functionOffset int32,
	functionLength int32,
	numArguments int32,
	argumentsLengthOffset int32,
	dataOffset int32,
) int32 {
	host := arwen.GetVMContext(context)
	runtime := host.Runtime()
	metering := host.Metering()

	gasToUse := metering.GasSchedule().ElrondAPICost.ExecuteOnDestContext
	metering.UseGas(gasToUse)

	value, err := runtime.MemLoad(valueOffset, arwen.BalanceLen)
	if arwen.WithFault(err, context, runtime.ElrondSyncExecAPIErrorShouldFailExecution()) {
		return 1
	}

	sender := runtime.GetSCAddress()
	bigIntVal := big.NewInt(0).SetBytes(value)
	contractCallInput, err := prepareIndirectContractCallInput(
		host,
		sender,
		bigIntVal,
		gasLimit,
		addressOffset,
		functionOffset,
		functionLength,
		numArguments,
		argumentsLengthOffset,
		dataOffset,
	)
	if arwen.WithFault(err, context, runtime.ElrondSyncExecAPIErrorShouldFailExecution()) {
		return 1
	}

	_, _, err = host.ExecuteOnDestContext(contractCallInput)
	if arwen.WithFault(err, context, runtime.ElrondSyncExecAPIErrorShouldFailExecution()) {
		return 1
	}

	return 0
}

//export executeOnDestContextByCaller
func executeOnDestContextByCaller(
	context unsafe.Pointer,
	gasLimit int64,
	addressOffset int32,
	valueOffset int32,
	functionOffset int32,
	functionLength int32,
	numArguments int32,
	argumentsLengthOffset int32,
	dataOffset int32,
) int32 {
	host := arwen.GetVMContext(context)
	runtime := host.Runtime()
	metering := host.Metering()

	gasToUse := metering.GasSchedule().ElrondAPICost.ExecuteOnDestContext
	metering.UseGas(gasToUse)

	value, err := runtime.MemLoad(valueOffset, arwen.BalanceLen)
	if arwen.WithFault(err, context, runtime.ElrondSyncExecAPIErrorShouldFailExecution()) {
		return 1
	}

	sender := runtime.GetVMInput().CallerAddr
	bigIntVal := big.NewInt(0).SetBytes(value)
	contractCallInput, err := prepareIndirectContractCallInput(
		host,
		sender,
		bigIntVal,
		gasLimit,
		addressOffset,
		functionOffset,
		functionLength,
		numArguments,
		argumentsLengthOffset,
		dataOffset,
	)
	if arwen.WithFault(err, context, runtime.ElrondSyncExecAPIErrorShouldFailExecution()) {
		return 1
	}

	_, _, err = host.ExecuteOnDestContext(contractCallInput)
	if arwen.WithFault(err, context, runtime.ElrondSyncExecAPIErrorShouldFailExecution()) {
		return 1
	}

	return 0
}

//export delegateExecution
func delegateExecution(
	context unsafe.Pointer,
	gasLimit int64,
	addressOffset int32,
	functionOffset int32,
	functionLength int32,
	numArguments int32,
	argumentsLengthOffset int32,
	dataOffset int32,
) int32 {
	host := arwen.GetVMContext(context)
	runtime := host.Runtime()
	metering := host.Metering()

	gasToUse := metering.GasSchedule().ElrondAPICost.DelegateExecution
	metering.UseGas(gasToUse)

	sender := runtime.GetSCAddress()
	value := runtime.GetVMInput().CallValue
	bigIntVal := big.NewInt(0).Set(value)
	contractCallInput, err := prepareIndirectContractCallInput(
		host,
		sender,
		bigIntVal,
		gasLimit,
		addressOffset,
		functionOffset,
		functionLength,
		numArguments,
		argumentsLengthOffset,
		dataOffset,
	)
	if arwen.WithFault(err, context, runtime.ElrondSyncExecAPIErrorShouldFailExecution()) {
		return 1
	}

	err = host.ExecuteOnSameContext(contractCallInput)
	if arwen.WithFault(err, context, runtime.ElrondSyncExecAPIErrorShouldFailExecution()) {
		return 1
	}

	return 0
}

//export executeReadOnly
func executeReadOnly(
	context unsafe.Pointer,
	gasLimit int64,
	addressOffset int32,
	functionOffset int32,
	functionLength int32,
	numArguments int32,
	argumentsLengthOffset int32,
	dataOffset int32,
) int32 {
	host := arwen.GetVMContext(context)
	runtime := host.Runtime()
	metering := host.Metering()

	gasToUse := metering.GasSchedule().ElrondAPICost.ExecuteReadOnly
	metering.UseGas(gasToUse)

	sender := runtime.GetSCAddress()
	value := runtime.GetVMInput().CallValue
	bigIntVal := big.NewInt(0).Set(value)

	contractCallInput, err := prepareIndirectContractCallInput(
		host,
		sender,
		bigIntVal,
		gasLimit,
		addressOffset,
		functionOffset,
		functionLength,
		numArguments,
		argumentsLengthOffset,
		dataOffset,
	)
	if arwen.WithFault(err, context, runtime.ElrondSyncExecAPIErrorShouldFailExecution()) {
		return 1
	}

	runtime.SetReadOnly(true)
	err = host.ExecuteOnSameContext(contractCallInput)
	runtime.SetReadOnly(false)
	if arwen.WithFault(err, context, runtime.ElrondSyncExecAPIErrorShouldFailExecution()) {
		return 1
	}

	return 0
}

//export createContract
func createContract(
	context unsafe.Pointer,
	gasLimit int64,
	valueOffset int32,
	codeOffset int32,
	codeMetadataOffset int32,
	length int32,
	resultOffset int32,
	numArguments int32,
	argumentsLengthOffset int32,
	dataOffset int32,
) int32 {
	host := arwen.GetVMContext(context)
	runtime := host.Runtime()
	metering := host.Metering()

	gasToUse := metering.GasSchedule().ElrondAPICost.CreateContract
	metering.UseGas(gasToUse)

	sender := runtime.GetSCAddress()
	value, err := runtime.MemLoad(valueOffset, arwen.BalanceLen)
	if arwen.WithFault(err, context, runtime.ElrondAPIErrorShouldFailExecution()) {
		return 1
	}

	code, err := runtime.MemLoad(codeOffset, length)
	if arwen.WithFault(err, context, runtime.ElrondAPIErrorShouldFailExecution()) {
		return 1
	}

	codeMetadata, err := runtime.MemLoad(codeMetadataOffset, arwen.CodeMetadataLen)
	if arwen.WithFault(err, context, runtime.ElrondAPIErrorShouldFailExecution()) {
		return 1
	}

	_, data, actualLen, err := getArgumentsFromMemory(
		host,
		0,
		0,
		numArguments,
		argumentsLengthOffset,
		dataOffset,
	)

	gasToUse = math.MulUint64(metering.GasSchedule().BaseOperationCost.DataCopyPerByte, uint64(actualLen))
	metering.UseGas(gasToUse)

	if arwen.WithFault(err, context, runtime.ElrondAPIErrorShouldFailExecution()) {
		return 1
	}

	contractCreate := &vmcommon.ContractCreateInput{
		VMInput: vmcommon.VMInput{
			CallerAddr:  sender,
			Arguments:   data,
			CallValue:   big.NewInt(0).SetBytes(value),
			GasPrice:    0,
			GasProvided: metering.BoundGasLimit(uint64(gasLimit)),
		},
		ContractCode:         code,
		ContractCodeMetadata: codeMetadata,
	}

	newAddress, err := host.CreateNewContract(contractCreate)
	if err != nil {
		return 1
	}

	err = runtime.MemStore(resultOffset, newAddress)
	if arwen.WithFault(err, context, runtime.ElrondAPIErrorShouldFailExecution()) {
		return 1
	}

	return 0
}

//export getNumReturnData
func getNumReturnData(context unsafe.Pointer) int32 {
	output := arwen.GetOutputContext(context)
	metering := arwen.GetMeteringContext(context)

	gasToUse := metering.GasSchedule().ElrondAPICost.GetNumReturnData
	metering.UseGas(gasToUse)

	returnData := output.ReturnData()
	return int32(len(returnData))
}

//export getReturnDataSize
func getReturnDataSize(context unsafe.Pointer, resultID int32) int32 {
	output := arwen.GetOutputContext(context)
	metering := arwen.GetMeteringContext(context)

	gasToUse := metering.GasSchedule().ElrondAPICost.GetReturnDataSize
	metering.UseGas(gasToUse)

	returnData := output.ReturnData()
	if resultID >= int32(len(returnData)) {
		return 0
	}

	return int32(len(returnData[resultID]))
}

//export getReturnData
func getReturnData(context unsafe.Pointer, resultID int32, dataOffset int32) int32 {
	runtime := arwen.GetRuntimeContext(context)
	output := arwen.GetOutputContext(context)
	metering := arwen.GetMeteringContext(context)

	gasToUse := metering.GasSchedule().ElrondAPICost.GetReturnData
	metering.UseGas(gasToUse)

	returnData := output.ReturnData()
	if resultID >= int32(len(returnData)) {
		return 0
	}

	err := runtime.MemStore(dataOffset, returnData[resultID])
	if arwen.WithFault(err, context, runtime.ElrondAPIErrorShouldFailExecution()) {
		return 0
	}

	return int32(len(returnData[resultID]))
}

//export getOriginalTxHash
func getOriginalTxHash(context unsafe.Pointer, dataOffset int32) {
	runtime := arwen.GetRuntimeContext(context)
	metering := arwen.GetMeteringContext(context)

	gasToUse := metering.GasSchedule().ElrondAPICost.GetBlockHash
	metering.UseGas(gasToUse)

	err := runtime.MemStore(dataOffset, runtime.GetOriginalTxHash())
	_ = arwen.WithFault(err, context, runtime.ElrondAPIErrorShouldFailExecution())
}

//export getPrevTxHash
func getPrevTxHash(context unsafe.Pointer, dataOffset int32) {
	runtime := arwen.GetRuntimeContext(context)
	metering := arwen.GetMeteringContext(context)

	gasToUse := metering.GasSchedule().ElrondAPICost.GetBlockHash
	metering.UseGas(gasToUse)

	err := runtime.MemStore(dataOffset, runtime.GetPrevTxHash())
	_ = arwen.WithFault(err, context, runtime.ElrondAPIErrorShouldFailExecution())
}

func prepareIndirectContractCallInput(
	host arwen.VMHost,
	sender []byte,
	value *big.Int,
	gasLimit int64,
	addressOffset int32,
	functionOffset int32,
	functionLength int32,
	numArguments int32,
	argumentsLengthOffset int32,
	dataOffset int32,
) (*vmcommon.ContractCallInput, error) {
	runtime := host.Runtime()
	metering := host.Metering()

	destination, err := runtime.MemLoad(addressOffset, arwen.AddressLen)
	if err != nil {
		return nil, err
	}

	if !host.AreInSameShard(runtime.GetSCAddress(), destination) {
		return nil, arwen.ErrSyncExecutionNotInSameShard
	}

	function, data, actualLen, err := getArgumentsFromMemory(
		host,
		functionOffset,
		functionLength,
		numArguments,
		argumentsLengthOffset,
		dataOffset,
	)

	gasToUse := math.MulUint64(metering.GasSchedule().BaseOperationCost.DataCopyPerByte, uint64(actualLen))
	metering.UseGas(gasToUse)
	if err != nil {
		return nil, err
	}

	contractCallInput := &vmcommon.ContractCallInput{
		VMInput: vmcommon.VMInput{
			CallerAddr:  sender,
			Arguments:   data,
			CallValue:   value,
			GasPrice:    0,
			GasProvided: metering.BoundGasLimit(uint64(gasLimit)),
		},
		RecipientAddr: destination,
		Function:      function,
	}

	return contractCallInput, nil
}

func getArgumentsFromMemory(
	host arwen.VMHost,
	functionOffset int32,
	functionLength int32,
	numArguments int32,
	argumentsLengthOffset int32,
	dataOffset int32,
) (string, [][]byte, int32, error) {
	runtime := host.Runtime()

	if numArguments < 0 {
		return "", nil, 0, fmt.Errorf("negative numArguments (%d)", numArguments)
	}

	function, err := runtime.MemLoad(functionOffset, functionLength)
	if err != nil {
		return "", nil, 0, err
	}

	argumentsLengthData, err := runtime.MemLoad(argumentsLengthOffset, numArguments*4)
	if err != nil {
		return "", nil, 0, err
	}

	argumentLengths := createInt32Array(argumentsLengthData, numArguments)
	data, err := runtime.MemLoadMultiple(dataOffset, argumentLengths)
	if err != nil {
		return "", nil, 0, err
	}

	totalArgumentBytes := int32(0)
	for _, length := range argumentLengths {
		totalArgumentBytes += length
	}

	return string(function), data, totalArgumentBytes, nil
}

func createInt32Array(rawData []byte, numIntegers int32) []int32 {
	integers := make([]int32, numIntegers)
	index := 0
	for cursor := 0; cursor < len(rawData); cursor += 4 {
		rawInt := rawData[cursor : cursor+4]
		actualInt := binary.LittleEndian.Uint32(rawInt)
		integers[index] = int32(actualInt)
		index++
	}
	return integers
}
