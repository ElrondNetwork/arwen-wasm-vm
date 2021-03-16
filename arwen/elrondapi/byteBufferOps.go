package elrondapi

// // Declare the function signatures (see [cgo](https://golang.org/cmd/cgo/)).
//
// #include <stdlib.h>
// typedef unsigned char uint8_t;
// typedef int int32_t;
//
// extern int32_t byteBufferNew(void* context, int32_t byteOffset, int32_t byteLength);
//
// extern int32_t byteBufferLength(void* context, int32_t reference);
// extern int32_t byteBufferGet(void* context, int32_t reference, int32_t byteOffset);
// extern void byteBufferSet(void* context, int32_t destination, int32_t byteOffset, int32_t byteLength);
import "C"

import (
	"unsafe"

	"github.com/ElrondNetwork/arwen-wasm-vm/arwen"
	"github.com/ElrondNetwork/arwen-wasm-vm/math"
	"github.com/ElrondNetwork/arwen-wasm-vm/wasmer"
)

// ByteBufferImports creates a new wasmer.Imports populated with the BytesHeap API methods
func ByteBufferImports(imports *wasmer.Imports) (*wasmer.Imports, error) {
	imports = imports.Namespace("env")

	imports, err := imports.Append("byteBufferNew", byteBufferNew, C.byteBufferNew)
	if err != nil {
		return nil, err
	}

	imports, err = imports.Append("byteBufferLength", byteBufferLength, C.byteBufferLength)
	if err != nil {
		return nil, err
	}

	return imports, nil
}

//export byteBufferNew
func byteBufferNew(context unsafe.Pointer, byteOffset int32, byteLength int32) int32 {
	bytesHeap := arwen.GetBytesHeapContext(context)
	runtime := arwen.GetRuntimeContext(context)
	metering := arwen.GetMeteringContext(context)

	gasToUse := metering.GasSchedule().BigIntAPICost.BigIntNew // TODO: change
	metering.UseGas(gasToUse)

	bytes, err := runtime.MemLoad(byteOffset, byteLength)
	if arwen.WithFault(err, context, runtime.ByteBufferAPIErrorShouldFailExecution()) {
		return -1
	}

	return bytesHeap.NewByteBuffer(bytes)
}

//export byteBufferLength
func byteBufferLength(context unsafe.Pointer, reference int32) int32 {
	bytesHeap := arwen.GetBytesHeapContext(context)
	metering := arwen.GetMeteringContext(context)

	gasToUse := metering.GasSchedule().BigIntAPICost.BigIntUnsignedByteLength // TODO: change
	metering.UseGas(gasToUse)

	bytes := bytesHeap.GetByteBuffer(reference)

	return int32(len(bytes))
}

//export byteBufferGet
func byteBufferGet(context unsafe.Pointer, reference int32, byteOffset int32) int32 {
	bytesHeap := arwen.GetBytesHeapContext(context)
	runtime := arwen.GetRuntimeContext(context)
	metering := arwen.GetMeteringContext(context)

	gasToUse := metering.GasSchedule().BigIntAPICost.BigIntGetUnsignedBytes
	metering.UseGas(gasToUse)

	bytes := bytesHeap.GetByteBuffer(reference)

	err := runtime.MemStore(byteOffset, bytes)
	if arwen.WithFault(err, context, runtime.ByteBufferAPIErrorShouldFailExecution()) {
		return 0
	}

	gasToUse = math.MulUint64(metering.GasSchedule().BaseOperationCost.DataCopyPerByte, uint64(len(bytes)))
	metering.UseGas(gasToUse)

	return int32(len(bytes))
}

//export byteBufferSet
func byteBufferSet(context unsafe.Pointer, destination int32, byteOffset int32, byteLength int32) {
	bytesHeap := arwen.GetBytesHeapContext(context)
	runtime := arwen.GetRuntimeContext(context)
	metering := arwen.GetMeteringContext(context)

	gasToUse := metering.GasSchedule().BigIntAPICost.BigIntSetUnsignedBytes
	metering.UseGas(gasToUse)

	bytes, err := runtime.MemLoad(byteOffset, byteLength)
	if arwen.WithFault(err, context, runtime.ByteBufferAPIErrorShouldFailExecution()) {
		return
	}

	gasToUse = math.MulUint64(metering.GasSchedule().BaseOperationCost.DataCopyPerByte, uint64(len(bytes)))
	metering.UseGas(gasToUse)

	bytesHeap.SetByteBuffer(destination, bytes)
}
