package cryptoapi

// // Declare the function signatures (see [cgo](https://golang.org/cmd/cgo/)).
//
// #include <stdlib.h>
// typedef unsigned char uint8_t;
// typedef int int32_t;
//
// extern int32_t v1_3_sha256(void* context, int32_t dataOffset, int32_t length, int32_t resultOffset);
// extern int32_t v1_3_keccak256(void *context, int32_t dataOffset, int32_t length, int32_t resultOffset);
// extern int32_t v1_3_ripemd160(void *context, int32_t dataOffset, int32_t length, int32_t resultOffset);
// extern int32_t v1_3_verifyBLS(void *context, int32_t keyOffset, int32_t messageOffset, int32_t messageLength, int32_t sigOffset);
// extern int32_t v1_3_verifyEd25519(void *context, int32_t keyOffset, int32_t messageOffset, int32_t messageLength, int32_t sigOffset);
// extern int32_t v1_3_verifySecp256k1(void *context, int32_t keyOffset, int32_t keyLength, int32_t messageOffset, int32_t messageLength, int32_t sigOffset);
// extern void v1_3_addEC(void *context, int32_t xResultHandle, int32_t yResultHandle, int32_t ecHandle, int32_t fstPointXHandle, int32_t fstPointYHandle, int32_t sndPointXHandle, int32_t sndPointYHandle);
// extern void v1_3_doubleEC(void *context, int32_t xResultHandle, int32_t yResultHandle, int32_t ecHandle, int32_t pointXHandle, int32_t pointYHandle);
// extern int32_t v1_3_isOnCurveEC(void *context, int32_t ecHandle, int32_t pointXHandle, int32_t pointYHandle);
// extern int32_t v1_3_scalarBaseMultEC(void *context, int32_t xResultHandle, int32_t yResultHandle, int32_t ecHandle, int32_t dataOffset, int32_t length);
// extern int32_t v1_3_scalarMultEC(void *context, int32_t xResultHandle, int32_t yResultHandle, int32_t ecHandle, int32_t pointXHandle, int32_t pointYHandle, int32_t dataOffset, int32_t length);
// extern int32_t v1_3_marshalEC(void *context, int32_t xPairHandle, int32_t yPairHandle, int32_t ecHandle, int32_t resultOffset);
// extern int32_t v1_3_unmarshalEC(void *context, int32_t xResultHandle, int32_t yResultHandle, int32_t ecHandle, int32_t dataOffset, int32_t length);
// extern int32_t v1_3_marshalCompressedEC(void *context, int32_t xPairHandle, int32_t yPairHandle, int32_t ecHandle, int32_t resultOffset);
// extern int32_t v1_3_unmarshalCompressedEC(void *context, int32_t xResultHandle, int32_t yResultHandle, int32_t ecHandle, int32_t dataOffset, int32_t length);
// extern int32_t v1_3_generateKeyEC(void *context, int32_t xPubKeyHandle, int32_t yPubKeyHandle, int32_t ecHandle, int32_t resultOffset);
// extern int32_t v1_3_p224Ec(void *context);
// extern int32_t v1_3_p256Ec(void *context);
// extern int32_t v1_3_p384Ec(void *context);
// extern int32_t v1_3_p521Ec(void *context);
// extern int32_t v1_3_getCurveLengthEC(void *context, int32_t ecHandle);
// extern int32_t v1_3_getPrivKeyByteLengthEC(void *context, int32_t ecHandle);
// extern int32_t v1_3_ellipticCurveGetValues(void *context, int32_t ecHandle, int32_t fieldOrderHandle, int32_t basePointOrderHandle, int32_t eqConstantHandle, int32_t xBasePointHandle, int32_t yBasePointHandle);
import "C"

import (
	"crypto/elliptic"
	"crypto/rand"
	"unsafe"

	"github.com/ElrondNetwork/arwen-wasm-vm/v1_3/arwen"
	"github.com/ElrondNetwork/arwen-wasm-vm/v1_3/math"
	"github.com/ElrondNetwork/arwen-wasm-vm/v1_3/wasmer"
)

const blsPublicKeyLength = 96
const blsSignatureLength = 48
const ed25519PublicKeyLength = 32
const ed25519SignatureLength = 64
const secp256k1CompressedPublicKeyLength = 33
const secp256k1UncompressedPublicKeyLength = 65
const secp256k1SignatureLength = 64

// CryptoImports adds some crypto imports to the Wasmer Imports map
func CryptoImports(imports *wasmer.Imports) (*wasmer.Imports, error) {
	imports = imports.Namespace("env")
	imports, err := imports.Append("sha256", v1_3_sha256, C.v1_3_sha256)
	if err != nil {
		return nil, err
	}

	imports, err = imports.Append("keccak256", v1_3_keccak256, C.v1_3_keccak256)
	if err != nil {
		return nil, err
	}

	imports, err = imports.Append("ripemd160", v1_3_ripemd160, C.v1_3_ripemd160)
	if err != nil {
		return nil, err
	}

	imports, err = imports.Append("verifyBLS", v1_3_verifyBLS, C.v1_3_verifyBLS)
	if err != nil {
		return nil, err
	}

	imports, err = imports.Append("verifyEd25519", v1_3_verifyEd25519, C.v1_3_verifyEd25519)
	if err != nil {
		return nil, err
	}

	imports, err = imports.Append("verifySecp256k1", v1_3_verifySecp256k1, C.v1_3_verifySecp256k1)
	if err != nil {
		return nil, err
	}

	imports, err = imports.Append("addEC", v1_3_addEC, C.v1_3_addEC)
	if err != nil {
		return nil, err
	}

	imports, err = imports.Append("doubleEC", v1_3_doubleEC, C.v1_3_doubleEC)
	if err != nil {
		return nil, err
	}

	imports, err = imports.Append("isOnCurveEC", v1_3_isOnCurveEC, C.v1_3_isOnCurveEC)
	if err != nil {
		return nil, err
	}

	imports, err = imports.Append("scalarBaseMultEC", v1_3_scalarBaseMultEC, C.v1_3_scalarBaseMultEC)
	if err != nil {
		return nil, err
	}

	imports, err = imports.Append("scalarMultEC", v1_3_scalarMultEC, C.v1_3_scalarMultEC)
	if err != nil {
		return nil, err
	}

	imports, err = imports.Append("marshalEC", v1_3_marshalEC, C.v1_3_marshalEC)
	if err != nil {
		return nil, err
	}

	imports, err = imports.Append("unmarshalEC", v1_3_unmarshalEC, C.v1_3_unmarshalEC)
	if err != nil {
		return nil, err
	}

	imports, err = imports.Append("marshalCompressedEC", v1_3_marshalCompressedEC, C.v1_3_marshalCompressedEC)
	if err != nil {
		return nil, err
	}

	imports, err = imports.Append("unmarshalCompressedEC", v1_3_unmarshalCompressedEC, C.v1_3_unmarshalCompressedEC)
	if err != nil {
		return nil, err
	}

	imports, err = imports.Append("generateKeyEC", v1_3_generateKeyEC, C.v1_3_generateKeyEC)
	if err != nil {
		return nil, err
	}

	imports, err = imports.Append("p224Ec", v1_3_p224Ec, C.v1_3_p224Ec)
	if err != nil {
		return nil, err
	}

	imports, err = imports.Append("p256Ec", v1_3_p256Ec, C.v1_3_p256Ec)
	if err != nil {
		return nil, err
	}

	imports, err = imports.Append("p384Ec", v1_3_p384Ec, C.v1_3_p384Ec)
	if err != nil {
		return nil, err
	}

	imports, err = imports.Append("p521Ec", v1_3_p521Ec, C.v1_3_p521Ec)
	if err != nil {
		return nil, err
	}

	imports, err = imports.Append("getCurveLengthEC", v1_3_getCurveLengthEC, C.v1_3_getCurveLengthEC)
	if err != nil {
		return nil, err
	}

	imports, err = imports.Append("getPrivKeyByteLengthEC", v1_3_getPrivKeyByteLengthEC, C.v1_3_getPrivKeyByteLengthEC)
	if err != nil {
		return nil, err
	}

	imports, err = imports.Append("ellipticCurveGetValues", v1_3_ellipticCurveGetValues, C.v1_3_ellipticCurveGetValues)
	if err != nil {
		return nil, err
	}
	return imports, nil
}

//export v1_3_sha256
func v1_3_sha256(context unsafe.Pointer, dataOffset int32, length int32, resultOffset int32) int32 {
	runtime := arwen.GetRuntimeContext(context)
	crypto := arwen.GetCryptoContext(context)
	metering := arwen.GetMeteringContext(context)

	memLoadGas := math.MulUint64(metering.GasSchedule().BaseOperationCost.DataCopyPerByte, uint64(length))
	gasToUse := math.AddUint64(metering.GasSchedule().CryptoAPICost.SHA256, memLoadGas)
	metering.UseGas(gasToUse)

	data, err := runtime.MemLoad(dataOffset, length)
	if arwen.WithFault(err, context, runtime.CryptoAPIErrorShouldFailExecution()) {
		return 1
	}

	result, err := crypto.Sha256(data)
	if err != nil {
		return 1
	}

	err = runtime.MemStore(resultOffset, result)
	if arwen.WithFault(err, context, runtime.CryptoAPIErrorShouldFailExecution()) {
		return 1
	}

	return 0
}

//export v1_3_keccak256
func v1_3_keccak256(context unsafe.Pointer, dataOffset int32, length int32, resultOffset int32) int32 {
	runtime := arwen.GetRuntimeContext(context)
	crypto := arwen.GetCryptoContext(context)
	metering := arwen.GetMeteringContext(context)

	memLoadGas := math.MulUint64(metering.GasSchedule().BaseOperationCost.DataCopyPerByte, uint64(length))
	gasToUse := math.AddUint64(metering.GasSchedule().CryptoAPICost.Keccak256, memLoadGas)
	metering.UseGas(gasToUse)

	data, err := runtime.MemLoad(dataOffset, length)
	if arwen.WithFault(err, context, runtime.CryptoAPIErrorShouldFailExecution()) {
		return 1
	}

	result, err := crypto.Keccak256(data)
	if err != nil {
		return 1
	}

	err = runtime.MemStore(resultOffset, result)
	if arwen.WithFault(err, context, runtime.CryptoAPIErrorShouldFailExecution()) {
		return 1
	}

	return 0
}

//export v1_3_ripemd160
func v1_3_ripemd160(context unsafe.Pointer, dataOffset int32, length int32, resultOffset int32) int32 {
	runtime := arwen.GetRuntimeContext(context)
	crypto := arwen.GetCryptoContext(context)
	metering := arwen.GetMeteringContext(context)

	memLoadGas := math.MulUint64(metering.GasSchedule().BaseOperationCost.DataCopyPerByte, uint64(length))
	gasToUse := math.AddUint64(metering.GasSchedule().CryptoAPICost.Ripemd160, memLoadGas)
	metering.UseGas(gasToUse)

	data, err := runtime.MemLoad(dataOffset, length)
	if arwen.WithFault(err, context, runtime.CryptoAPIErrorShouldFailExecution()) {
		return 1
	}

	result, err := crypto.Ripemd160(data)
	if err != nil {
		return 1
	}

	err = runtime.MemStore(resultOffset, result)
	if arwen.WithFault(err, context, runtime.CryptoAPIErrorShouldFailExecution()) {
		return 1
	}

	return 0
}

//export v1_3_verifyBLS
func v1_3_verifyBLS(
	context unsafe.Pointer,
	keyOffset int32,
	messageOffset int32,
	messageLength int32,
	sigOffset int32,
) int32 {
	runtime := arwen.GetRuntimeContext(context)
	crypto := arwen.GetCryptoContext(context)
	metering := arwen.GetMeteringContext(context)

	gasToUse := metering.GasSchedule().CryptoAPICost.VerifyBLS
	metering.UseGas(gasToUse)

	key, err := runtime.MemLoad(keyOffset, blsPublicKeyLength)
	if arwen.WithFault(err, context, runtime.CryptoAPIErrorShouldFailExecution()) {
		return 1
	}

	gasToUse = math.MulUint64(metering.GasSchedule().BaseOperationCost.DataCopyPerByte, uint64(messageLength))
	metering.UseGas(gasToUse)

	message, err := runtime.MemLoad(messageOffset, messageLength)
	if arwen.WithFault(err, context, runtime.CryptoAPIErrorShouldFailExecution()) {
		return 1
	}

	sig, err := runtime.MemLoad(sigOffset, blsSignatureLength)
	if arwen.WithFault(err, context, runtime.CryptoAPIErrorShouldFailExecution()) {
		return 1
	}

	invalidSigErr := crypto.VerifyBLS(key, message, sig)
	if invalidSigErr != nil {
		return -1
	}

	return 0
}

//export v1_3_verifyEd25519
func v1_3_verifyEd25519(
	context unsafe.Pointer,
	keyOffset int32,
	messageOffset int32,
	messageLength int32,
	sigOffset int32,
) int32 {
	runtime := arwen.GetRuntimeContext(context)
	crypto := arwen.GetCryptoContext(context)
	metering := arwen.GetMeteringContext(context)

	gasToUse := metering.GasSchedule().CryptoAPICost.VerifyEd25519
	metering.UseGas(gasToUse)

	key, err := runtime.MemLoad(keyOffset, ed25519PublicKeyLength)
	if arwen.WithFault(err, context, runtime.CryptoAPIErrorShouldFailExecution()) {
		return 1
	}

	gasToUse = math.MulUint64(metering.GasSchedule().BaseOperationCost.DataCopyPerByte, uint64(messageLength))
	metering.UseGas(gasToUse)

	message, err := runtime.MemLoad(messageOffset, messageLength)
	if arwen.WithFault(err, context, runtime.CryptoAPIErrorShouldFailExecution()) {
		return 1
	}

	sig, err := runtime.MemLoad(sigOffset, ed25519SignatureLength)
	if arwen.WithFault(err, context, runtime.CryptoAPIErrorShouldFailExecution()) {
		return 1
	}

	invalidSigErr := crypto.VerifyEd25519(key, message, sig)
	if invalidSigErr != nil {
		return -1
	}

	return 0
}

//export v1_3_verifySecp256k1
func v1_3_verifySecp256k1(
	context unsafe.Pointer,
	keyOffset int32,
	keyLength int32,
	messageOffset int32,
	messageLength int32,
	sigOffset int32,
) int32 {
	runtime := arwen.GetRuntimeContext(context)
	crypto := arwen.GetCryptoContext(context)
	metering := arwen.GetMeteringContext(context)

	gasToUse := metering.GasSchedule().CryptoAPICost.VerifySecp256k1
	metering.UseGas(gasToUse)

	if keyLength != secp256k1CompressedPublicKeyLength && keyLength != secp256k1UncompressedPublicKeyLength {
		arwen.WithFault(arwen.ErrInvalidPublicKeySize, context, runtime.ElrondAPIErrorShouldFailExecution())
		return 1
	}

	key, err := runtime.MemLoad(keyOffset, keyLength)
	if arwen.WithFault(err, context, runtime.CryptoAPIErrorShouldFailExecution()) {
		return 1
	}

	gasToUse = math.MulUint64(metering.GasSchedule().BaseOperationCost.DataCopyPerByte, uint64(messageLength))
	metering.UseGas(gasToUse)

	message, err := runtime.MemLoad(messageOffset, messageLength)
	if arwen.WithFault(err, context, runtime.CryptoAPIErrorShouldFailExecution()) {
		return 1
	}

	// read the 2 leading bytes first
	// byte1: 0x30, header
	// byte2: the remaining buffer length
	const sigHeaderLength = 2
	sigHeader, err := runtime.MemLoad(sigOffset, sigHeaderLength)
	if arwen.WithFault(err, context, runtime.CryptoAPIErrorShouldFailExecution()) {
		return 1
	}
	sigLength := int32(sigHeader[1]) + sigHeaderLength
	sig, err := runtime.MemLoad(sigOffset, sigLength)
	if arwen.WithFault(err, context, runtime.CryptoAPIErrorShouldFailExecution()) {
		return 1
	}

	invalidSigErr := crypto.VerifySecp256k1(key, message, sig)
	if invalidSigErr != nil {
		return -1
	}

	return 0
}

//export v1_3_addEC
func v1_3_addEC(
	context unsafe.Pointer,
	xResultHandle int32,
	yResultHandle int32,
	ecHandle int32,
	fstPointXHandle int32,
	fstPointYHandle int32,
	sndPointXHandle int32,
	sndPointYHandle int32,
) {
	managedType := arwen.GetManagedTypesContext(context)
	metering := arwen.GetMeteringContext(context)
	runtime := arwen.GetRuntimeContext(context)

	curveMultiplier := managedType.Get100xCurveGasCostMultiplier(ecHandle)
	if curveMultiplier < 0 {
		arwen.WithFault(arwen.ErrNoEllipticCurveUnderThisHandle, context, runtime.CryptoAPIErrorShouldFailExecution())
		return
	}
	gasToUse := metering.GasSchedule().CryptoAPICost.AddECC * uint64(curveMultiplier) / 100
	metering.UseGas(gasToUse)

	ec, err1 := managedType.GetEllipticCurve(ecHandle)
	if err1 != nil {
		arwen.WithFault(arwen.ErrNoEllipticCurveUnderThisHandle, context, runtime.CryptoAPIErrorShouldFailExecution())
		return
	}

	xResult, yResult, err1 := managedType.GetTwoBigInt(xResultHandle, yResultHandle)
	x1, y1, err2 := managedType.GetTwoBigInt(fstPointXHandle, fstPointYHandle)
	x2, y2, err3 := managedType.GetTwoBigInt(sndPointXHandle, sndPointYHandle)
	if err1 != nil || err2 != nil || err3 != nil {
		arwen.WithFault(arwen.ErrNoBigIntUnderThisHandle, context, runtime.BigIntAPIErrorShouldFailExecution())
		return
	}
	if !ec.IsOnCurve(x1, y1) || !ec.IsOnCurve(x2, y2) {
		arwen.WithFault(arwen.ErrPointNotOnCurve, context, runtime.CryptoAPIErrorShouldFailExecution())
	}

	managedType.ConsumeGasForBigIntCopy(xResult, yResult, ec.P, ec.N, ec.B, ec.Gx, ec.Gy, x1, y1, x2, y2)
	xResultAdd, yResultAdd := ec.Add(x1, y1, x2, y2)
	xResult.Set(xResultAdd)
	yResult.Set(yResultAdd)
}

//export v1_3_doubleEC
func v1_3_doubleEC(
	context unsafe.Pointer,
	xResultHandle int32,
	yResultHandle int32,
	ecHandle int32,
	pointXHandle int32,
	pointYHandle int32,
) {
	managedType := arwen.GetManagedTypesContext(context)
	metering := arwen.GetMeteringContext(context)
	runtime := arwen.GetRuntimeContext(context)

	curveMultiplier := managedType.Get100xCurveGasCostMultiplier(ecHandle)
	if curveMultiplier < 0 {
		arwen.WithFault(arwen.ErrNoEllipticCurveUnderThisHandle, context, runtime.CryptoAPIErrorShouldFailExecution())
		return
	}
	gasToUse := metering.GasSchedule().CryptoAPICost.DoubleECC * uint64(curveMultiplier) / 100
	metering.UseGas(gasToUse)

	ec, err1 := managedType.GetEllipticCurve(ecHandle)
	if err1 != nil {
		arwen.WithFault(arwen.ErrNoEllipticCurveUnderThisHandle, context, runtime.CryptoAPIErrorShouldFailExecution())
		return
	}

	xResult, yResult, err1 := managedType.GetTwoBigInt(xResultHandle, yResultHandle)
	x, y, err2 := managedType.GetTwoBigInt(pointXHandle, pointYHandle)
	if err1 != nil || err2 != nil {
		arwen.WithFault(arwen.ErrNoBigIntUnderThisHandle, context, runtime.CryptoAPIErrorShouldFailExecution())
		return
	}
	if !ec.IsOnCurve(x, y) {
		arwen.WithFault(arwen.ErrPointNotOnCurve, context, runtime.CryptoAPIErrorShouldFailExecution())
	}

	managedType.ConsumeGasForBigIntCopy(xResult, yResult, ec.P, ec.N, ec.B, ec.Gx, ec.Gy, x, y)
	xResultDouble, yResultDouble := ec.Double(x, y)
	xResult.Set(xResultDouble)
	yResult.Set(yResultDouble)
}

//export v1_3_isOnCurveEC
func v1_3_isOnCurveEC(
	context unsafe.Pointer,
	ecHandle int32,
	pointXHandle int32,
	pointYHandle int32,
) int32 {
	managedType := arwen.GetManagedTypesContext(context)
	metering := arwen.GetMeteringContext(context)
	runtime := arwen.GetRuntimeContext(context)

	curveMultiplier := managedType.Get100xCurveGasCostMultiplier(ecHandle)
	if curveMultiplier < 0 {
		arwen.WithFault(arwen.ErrNoEllipticCurveUnderThisHandle, context, runtime.CryptoAPIErrorShouldFailExecution())
		return 1
	}
	gasToUse := metering.GasSchedule().CryptoAPICost.IsOnCurveECC * uint64(curveMultiplier) / 100
	metering.UseGas(gasToUse)

	ec, err := managedType.GetEllipticCurve(ecHandle)
	if err != nil {
		arwen.WithFault(arwen.ErrNoEllipticCurveUnderThisHandle, context, runtime.CryptoAPIErrorShouldFailExecution())
		return -1
	}

	x, y, err := managedType.GetTwoBigInt(pointXHandle, pointYHandle)
	if err != nil || x == nil || y == nil {
		arwen.WithFault(arwen.ErrNoBigIntUnderThisHandle, context, runtime.CryptoAPIErrorShouldFailExecution())
		return -1
	}

	managedType.ConsumeGasForBigIntCopy(ec.P, ec.N, ec.B, ec.Gx, ec.Gy, x, y)
	if ec.IsOnCurve(x, y) {
		return 1
	}

	return 0
}

//export v1_3_scalarBaseMultEC
func v1_3_scalarBaseMultEC(
	context unsafe.Pointer,
	xResultHandle int32,
	yResultHandle int32,
	ecHandle int32,
	dataOffset int32,
	length int32,
) int32 {
	runtime := arwen.GetRuntimeContext(context)
	metering := arwen.GetMeteringContext(context)
	managedType := arwen.GetManagedTypesContext(context)

	if length < 0 {
		arwen.WithFault(arwen.ErrNegativeLength, context, runtime.CryptoAPIErrorShouldFailExecution())
		return 1
	}

	curveMultiplier := managedType.GetScalarMult100xCurveGasCostMultiplier(ecHandle)
	if curveMultiplier < 0 {
		arwen.WithFault(arwen.ErrNoEllipticCurveUnderThisHandle, context, runtime.CryptoAPIErrorShouldFailExecution())
		return 1
	}
	oneByteScalarGasCost := metering.GasSchedule().CryptoAPICost.ScalarMultECC * uint64(curveMultiplier) / 100
	gasToUse := oneByteScalarGasCost + uint64(length)*oneByteScalarGasCost
	metering.UseGas(gasToUse)

	data, err := runtime.MemLoad(dataOffset, length)
	if arwen.WithFault(err, context, runtime.CryptoAPIErrorShouldFailExecution()) {
		return 1
	}

	ec, err := managedType.GetEllipticCurve(ecHandle)
	if err != nil {
		arwen.WithFault(arwen.ErrNoEllipticCurveUnderThisHandle, context, runtime.CryptoAPIErrorShouldFailExecution())
		return 1
	}

	xResult, yResult, err := managedType.GetTwoBigInt(xResultHandle, yResultHandle)
	if err != nil {
		arwen.WithFault(arwen.ErrNoBigIntUnderThisHandle, context, runtime.CryptoAPIErrorShouldFailExecution())
		return 1
	}

	managedType.ConsumeGasForBigIntCopy(ec.P, ec.N, ec.B, ec.Gx, ec.Gy, xResult, yResult)
	xResultSBM, yResultSBM := ec.ScalarBaseMult(data)
	xResult.Set(xResultSBM)
	yResult.Set(yResultSBM)

	return 0
}

//export v1_3_scalarMultEC
func v1_3_scalarMultEC(
	context unsafe.Pointer,
	xResultHandle int32,
	yResultHandle int32,
	ecHandle int32,
	pointXHandle int32,
	pointYHandle int32,
	dataOffset int32,
	length int32,
) int32 {
	runtime := arwen.GetRuntimeContext(context)
	metering := arwen.GetMeteringContext(context)
	managedType := arwen.GetManagedTypesContext(context)

	if length < 0 {
		arwen.WithFault(arwen.ErrNegativeLength, context, runtime.CryptoAPIErrorShouldFailExecution())
		return 1
	}

	curveMultiplier := managedType.GetScalarMult100xCurveGasCostMultiplier(ecHandle)
	if curveMultiplier < 0 {
		arwen.WithFault(arwen.ErrNoEllipticCurveUnderThisHandle, context, runtime.CryptoAPIErrorShouldFailExecution())
		return 1
	}
	oneByteScalarGasCost := metering.GasSchedule().CryptoAPICost.ScalarMultECC * uint64(curveMultiplier) / 100
	gasToUse := oneByteScalarGasCost + uint64(length)*oneByteScalarGasCost
	metering.UseGas(gasToUse)

	data, err := runtime.MemLoad(dataOffset, length)
	if arwen.WithFault(err, context, runtime.CryptoAPIErrorShouldFailExecution()) {
		return 1
	}

	ec, err1 := managedType.GetEllipticCurve(ecHandle)
	if err1 != nil {
		arwen.WithFault(arwen.ErrNoEllipticCurveUnderThisHandle, context, runtime.CryptoAPIErrorShouldFailExecution())
		return 1
	}

	xResult, yResult, err1 := managedType.GetTwoBigInt(xResultHandle, yResultHandle)
	x, y, err2 := managedType.GetTwoBigInt(pointXHandle, pointYHandle)
	if err1 != nil || err2 != nil {
		arwen.WithFault(arwen.ErrNoBigIntUnderThisHandle, context, runtime.CryptoAPIErrorShouldFailExecution())
		return 1
	}
	if !ec.IsOnCurve(x, y) {
		arwen.WithFault(arwen.ErrPointNotOnCurve, context, runtime.CryptoAPIErrorShouldFailExecution())
	}

	managedType.ConsumeGasForBigIntCopy(xResult, yResult, ec.P, ec.N, ec.B, ec.Gx, ec.Gy, x, y)
	xResultSM, yResultSM := ec.ScalarMult(x, y, data)
	xResult.Set(xResultSM)
	yResult.Set(yResultSM)

	return 0
}

//export v1_3_marshalEC
func v1_3_marshalEC(
	context unsafe.Pointer,
	xPairHandle int32,
	yPairHandle int32,
	ecHandle int32,
	resultOffset int32,
) int32 {
	runtime := arwen.GetRuntimeContext(context)
	metering := arwen.GetMeteringContext(context)
	managedType := arwen.GetManagedTypesContext(context)

	curveMultiplier := managedType.Get100xCurveGasCostMultiplier(ecHandle)
	if curveMultiplier < 0 {
		arwen.WithFault(arwen.ErrNoEllipticCurveUnderThisHandle, context, runtime.CryptoAPIErrorShouldFailExecution())
		return 1
	}
	gasToUse := metering.GasSchedule().CryptoAPICost.MarshalECC * uint64(curveMultiplier) / 100
	metering.UseGas(gasToUse)

	ec, err := managedType.GetEllipticCurve(ecHandle)
	if err != nil {
		arwen.WithFault(arwen.ErrNoEllipticCurveUnderThisHandle, context, runtime.CryptoAPIErrorShouldFailExecution())
		return 1
	}

	x, y, err := managedType.GetTwoBigInt(xPairHandle, yPairHandle)
	if err != nil {
		arwen.WithFault(arwen.ErrNoBigIntUnderThisHandle, context, runtime.CryptoAPIErrorShouldFailExecution())
		return 1
	}
	if !ec.IsOnCurve(x, y) {
		arwen.WithFault(arwen.ErrPointNotOnCurve, context, runtime.CryptoAPIErrorShouldFailExecution())
	}
	if x.BitLen() > int(ec.BitSize) || y.BitLen() > int(ec.BitSize) {
		arwen.WithFault(arwen.ErrLengthOfBufferNotCorrect, context, runtime.CryptoAPIErrorShouldFailExecution())
		return 1
	}

	managedType.ConsumeGasForBigIntCopy(ec.P, ec.N, ec.B, ec.Gx, ec.Gy, x, y)
	result := elliptic.Marshal(ec, x, y)
	err = runtime.MemStore(resultOffset, result)
	if arwen.WithFault(err, context, runtime.CryptoAPIErrorShouldFailExecution()) {
		return int32(len(result))
	}
	return 0
}

//export v1_3_marshalCompressedEC
func v1_3_marshalCompressedEC(
	context unsafe.Pointer,
	xPairHandle int32,
	yPairHandle int32,
	ecHandle int32,
	resultOffset int32,
) int32 {
	runtime := arwen.GetRuntimeContext(context)
	metering := arwen.GetMeteringContext(context)
	managedType := arwen.GetManagedTypesContext(context)

	curveMultiplier := managedType.Get100xCurveGasCostMultiplier(ecHandle)
	if curveMultiplier < 0 {
		arwen.WithFault(arwen.ErrNoEllipticCurveUnderThisHandle, context, runtime.CryptoAPIErrorShouldFailExecution())
		return 1
	}
	gasToUse := metering.GasSchedule().CryptoAPICost.MarshalCompressedECC * uint64(curveMultiplier) / 100
	metering.UseGas(gasToUse)

	ec, err := managedType.GetEllipticCurve(ecHandle)
	if err != nil {
		arwen.WithFault(arwen.ErrNoEllipticCurveUnderThisHandle, context, runtime.CryptoAPIErrorShouldFailExecution())
		return 1
	}

	x, y, err := managedType.GetTwoBigInt(xPairHandle, yPairHandle)
	if err != nil || x == nil || y == nil {
		arwen.WithFault(arwen.ErrNoBigIntUnderThisHandle, context, runtime.CryptoAPIErrorShouldFailExecution())
		return 1
	}
	if !ec.IsOnCurve(x, y) {
		arwen.WithFault(arwen.ErrPointNotOnCurve, context, runtime.CryptoAPIErrorShouldFailExecution())
	}
	if x.BitLen() > int(ec.BitSize) || y.BitLen() > int(ec.BitSize) {
		arwen.WithFault(arwen.ErrLengthOfBufferNotCorrect, context, runtime.CryptoAPIErrorShouldFailExecution())
		return 1
	}

	managedType.ConsumeGasForBigIntCopy(ec.P, ec.N, ec.B, ec.Gx, ec.Gy, x, y)
	result := elliptic.MarshalCompressed(ec, x, y)
	err = runtime.MemStore(resultOffset, result)
	if arwen.WithFault(err, context, runtime.CryptoAPIErrorShouldFailExecution()) {
		return int32(len(result))
	}
	return 0
}

//export v1_3_unmarshalEC
func v1_3_unmarshalEC(
	context unsafe.Pointer,
	xResultHandle int32,
	yResultHandle int32,
	ecHandle int32,
	dataOffset int32,
	length int32,
) int32 {
	runtime := arwen.GetRuntimeContext(context)
	metering := arwen.GetMeteringContext(context)
	managedType := arwen.GetManagedTypesContext(context)

	curveMultiplier := managedType.Get100xCurveGasCostMultiplier(ecHandle)
	if curveMultiplier < 0 {
		arwen.WithFault(arwen.ErrNoEllipticCurveUnderThisHandle, context, runtime.CryptoAPIErrorShouldFailExecution())
		return 1
	}
	gasToUse := metering.GasSchedule().CryptoAPICost.UnmarshalECC * uint64(curveMultiplier) / 100
	metering.UseGas(gasToUse)

	data, err := runtime.MemLoad(dataOffset, length)
	if arwen.WithFault(err, context, runtime.CryptoAPIErrorShouldFailExecution()) {
		return 1
	}

	ec, err := managedType.GetEllipticCurve(ecHandle)
	if err != nil {
		arwen.WithFault(arwen.ErrNoEllipticCurveUnderThisHandle, context, runtime.CryptoAPIErrorShouldFailExecution())
		return 1
	}
	byteLen := (ec.BitSize + 7) / 8
	if int(length) != 1+2*byteLen {
		arwen.WithFault(arwen.ErrLengthOfBufferNotCorrect, context, runtime.CryptoAPIErrorShouldFailExecution())
		return 1
	}

	xResult, yResult, err := managedType.GetTwoBigInt(xResultHandle, yResultHandle)
	if err != nil {
		arwen.WithFault(arwen.ErrNoBigIntUnderThisHandle, context, runtime.CryptoAPIErrorShouldFailExecution())
		return 1
	}

	managedType.ConsumeGasForBigIntCopy(ec.P, ec.N, ec.B, ec.Gx, ec.Gy, xResult, yResult)
	xResultU, yResultU := elliptic.Unmarshal(ec, data)
	if xResultU == nil || yResultU == nil || !ec.IsOnCurve(xResultU, yResultU) {
		arwen.WithFault(arwen.ErrPointNotOnCurve, context, runtime.CryptoAPIErrorShouldFailExecution())
		return 1
	}
	xResult.Set(xResultU)
	yResult.Set(yResultU)

	return 0
}

//export v1_3_unmarshalCompressedEC
func v1_3_unmarshalCompressedEC(
	context unsafe.Pointer,
	xResultHandle int32,
	yResultHandle int32,
	ecHandle int32,
	dataOffset int32,
	length int32,
) int32 {
	runtime := arwen.GetRuntimeContext(context)
	metering := arwen.GetMeteringContext(context)
	managedType := arwen.GetManagedTypesContext(context)

	curveMultiplier := managedType.GetUCompressed100xCurveGasCostMultiplier(ecHandle)
	if curveMultiplier < 0 {
		arwen.WithFault(arwen.ErrNoEllipticCurveUnderThisHandle, context, runtime.CryptoAPIErrorShouldFailExecution())
		return 1
	}
	gasToUse := metering.GasSchedule().CryptoAPICost.UnmarshalCompressedECC * uint64(curveMultiplier) / 100
	metering.UseGas(gasToUse)

	data, err := runtime.MemLoad(dataOffset, length)
	if arwen.WithFault(err, context, runtime.CryptoAPIErrorShouldFailExecution()) {
		return int32(len(data))
	}

	ec, err := managedType.GetEllipticCurve(ecHandle)
	if err != nil {
		arwen.WithFault(arwen.ErrNoEllipticCurveUnderThisHandle, context, runtime.CryptoAPIErrorShouldFailExecution())
		return 1
	}
	byteLen := (ec.BitSize + 7) / 8
	if int(length) != 1+byteLen {
		arwen.WithFault(arwen.ErrLengthOfBufferNotCorrect, context, runtime.CryptoAPIErrorShouldFailExecution())
		return 1
	}

	xResult, yResult, err := managedType.GetTwoBigInt(xResultHandle, yResultHandle)
	if err != nil {
		arwen.WithFault(arwen.ErrNoBigIntUnderThisHandle, context, runtime.CryptoAPIErrorShouldFailExecution())
		return 1
	}

	managedType.ConsumeGasForBigIntCopy(ec.P, ec.N, ec.B, ec.Gx, ec.Gy, xResult, yResult)
	xResultUC, yResultUC := elliptic.UnmarshalCompressed(ec, data)
	if xResultUC == nil || yResultUC == nil || !ec.IsOnCurve(xResultUC, yResultUC) {
		arwen.WithFault(arwen.ErrPointNotOnCurve, context, runtime.CryptoAPIErrorShouldFailExecution())
		return 1
	}
	xResult.Set(xResultUC)
	yResult.Set(yResultUC)
	return 0
}

//export v1_3_generateKeyEC
func v1_3_generateKeyEC(
	context unsafe.Pointer,
	xPubKeyHandle int32,
	yPubKeyHandle int32,
	ecHandle int32,
	resultOffset int32,
) int32 {
	runtime := arwen.GetRuntimeContext(context)
	metering := arwen.GetMeteringContext(context)
	managedType := arwen.GetManagedTypesContext(context)

	curveMultiplier := managedType.Get100xCurveGasCostMultiplier(ecHandle)
	if curveMultiplier < 0 {
		arwen.WithFault(arwen.ErrNoEllipticCurveUnderThisHandle, context, runtime.CryptoAPIErrorShouldFailExecution())
		return 1
	}
	if curveMultiplier == 250 {
		curveMultiplier = 500
	}
	gasToUse := metering.GasSchedule().CryptoAPICost.MarshalCompressedECC * uint64(curveMultiplier) / 100
	metering.UseGas(gasToUse)

	ec, err := managedType.GetEllipticCurve(ecHandle)
	if err != nil {
		arwen.WithFault(arwen.ErrNoEllipticCurveUnderThisHandle, context, runtime.CryptoAPIErrorShouldFailExecution())
		return 1
	}

	xPubKey, yPubKey, err := managedType.GetTwoBigInt(xPubKeyHandle, yPubKeyHandle)
	if err != nil {
		arwen.WithFault(arwen.ErrNoBigIntUnderThisHandle, context, runtime.CryptoAPIErrorShouldFailExecution())
		return 1
	}
	managedType.ConsumeGasForBigIntCopy(ec.P, ec.N, ec.B, ec.Gx, ec.Gy, xPubKey, yPubKey)

	result, xPubKeyGK, yPubKeyGK, err := elliptic.GenerateKey(ec, rand.Reader)
	if arwen.WithFault(err, context, runtime.CryptoAPIErrorShouldFailExecution()) {
		return int32(len(result))
	}

	err = runtime.MemStore(resultOffset, result)
	if arwen.WithFault(err, context, runtime.CryptoAPIErrorShouldFailExecution()) {
		return int32(len(result))
	}

	xPubKey.Set(xPubKeyGK)
	yPubKey.Set(yPubKeyGK)
	return 0
}

//export v1_3_p224Ec
func v1_3_p224Ec(context unsafe.Pointer) int32 {
	managedType := arwen.GetManagedTypesContext(context)
	metering := arwen.GetMeteringContext(context)

	gasToUse := metering.GasSchedule().CryptoAPICost.EllipticCurveNew
	metering.UseGas(gasToUse)

	curveParams := elliptic.P224().Params()
	return managedType.PutEllipticCurve(curveParams)
}

//export v1_3_p256Ec
func v1_3_p256Ec(context unsafe.Pointer) int32 {
	managedType := arwen.GetManagedTypesContext(context)
	metering := arwen.GetMeteringContext(context)

	gasToUse := metering.GasSchedule().CryptoAPICost.EllipticCurveNew
	metering.UseGas(gasToUse)

	curveParams := elliptic.P256().Params()
	return managedType.PutEllipticCurve(curveParams)
}

//export v1_3_p384Ec
func v1_3_p384Ec(context unsafe.Pointer) int32 {
	managedType := arwen.GetManagedTypesContext(context)
	metering := arwen.GetMeteringContext(context)

	gasToUse := metering.GasSchedule().CryptoAPICost.EllipticCurveNew
	metering.UseGas(gasToUse)

	curveParams := elliptic.P384().Params()
	return managedType.PutEllipticCurve(curveParams)
}

//export v1_3_p521Ec
func v1_3_p521Ec(context unsafe.Pointer) int32 {
	managedType := arwen.GetManagedTypesContext(context)
	metering := arwen.GetMeteringContext(context)

	gasToUse := metering.GasSchedule().CryptoAPICost.EllipticCurveNew
	metering.UseGas(gasToUse)

	curveParams := elliptic.P521().Params()
	return managedType.PutEllipticCurve(curveParams)
}

//export v1_3_getCurveLengthEC
func v1_3_getCurveLengthEC(context unsafe.Pointer, ecHandle int32) int32 {
	managedType := arwen.GetManagedTypesContext(context)
	metering := arwen.GetMeteringContext(context)
	runtime := arwen.GetRuntimeContext(context)

	gasToUse := metering.GasSchedule().BigIntAPICost.BigIntGetInt64
	metering.UseGas(gasToUse)

	ecLength := managedType.GetEllipticCurveSizeOfField(ecHandle)
	if ecLength == -1 {
		arwen.WithFault(arwen.ErrNoEllipticCurveUnderThisHandle, context, runtime.BigIntAPIErrorShouldFailExecution())
	}

	return ecLength
}

//export v1_3_getPrivKeyByteLengthEC
func v1_3_getPrivKeyByteLengthEC(context unsafe.Pointer, ecHandle int32) int32 {
	managedType := arwen.GetManagedTypesContext(context)
	metering := arwen.GetMeteringContext(context)
	runtime := arwen.GetRuntimeContext(context)

	gasToUse := metering.GasSchedule().BigIntAPICost.BigIntGetInt64
	metering.UseGas(gasToUse)

	byteLength := managedType.GetPrivateKeyByteLengthEC(ecHandle)
	if byteLength == -1 {
		arwen.WithFault(arwen.ErrNoEllipticCurveUnderThisHandle, context, runtime.BigIntAPIErrorShouldFailExecution())
	}

	return byteLength
}

//export v1_3_ellipticCurveGetValues
func v1_3_ellipticCurveGetValues(context unsafe.Pointer, ecHandle int32, fieldOrderHandle int32, basePointOrderHandle int32, eqConstantHandle int32, xBasePointHandle int32, yBasePointHandle int32) int32 {
	managedType := arwen.GetManagedTypesContext(context)
	metering := arwen.GetMeteringContext(context)
	runtime := arwen.GetRuntimeContext(context)

	gasToUse := metering.GasSchedule().BigIntAPICost.BigIntGetInt64 * 5
	metering.UseGas(gasToUse)

	ec, err := managedType.GetEllipticCurve(ecHandle)
	if err != nil {
		arwen.WithFault(arwen.ErrNoEllipticCurveUnderThisHandle, context, runtime.CryptoAPIErrorShouldFailExecution())
		return -1
	}
	fieldOrder, basePointOrder, err := managedType.GetTwoBigInt(fieldOrderHandle, basePointOrderHandle)
	if err != nil {
		arwen.WithFault(arwen.ErrNoBigIntUnderThisHandle, context, runtime.CryptoAPIErrorShouldFailExecution())
		return -1
	}
	eqConstant, err := managedType.GetBigInt(eqConstantHandle)
	if err != nil {
		arwen.WithFault(arwen.ErrNoBigIntUnderThisHandle, context, runtime.CryptoAPIErrorShouldFailExecution())
		return -1
	}
	xBasePoint, yBasePoint, err := managedType.GetTwoBigInt(xBasePointHandle, yBasePointHandle)
	if err != nil {
		arwen.WithFault(arwen.ErrNoBigIntUnderThisHandle, context, runtime.CryptoAPIErrorShouldFailExecution())
		return -1
	}
	fieldOrder.Set(ec.P)
	basePointOrder.Set(ec.N)
	eqConstant.Set(ec.B)
	xBasePoint.Set(ec.Gx)
	yBasePoint.Set(ec.Gy)
	return ecHandle
}
