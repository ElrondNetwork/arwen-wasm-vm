package wasmer

// #cgo LDFLAGS: -Wl,-rpath,${SRCDIR} -L${SRCDIR}
// #cgo linux,amd64 LDFLAGS:-lwasmer_linux_amd64
// #cgo linux,arm64 LDFLAGS:-lwasmer_linux_arm64
// #cgo darwin,amd64 LDFLAGS:-lwasmer_darwin_amd64
// #include "./wasmer.h"
//
import "C"
import "unsafe"

type cBool C.bool
type cChar C.char
type cInt C.int
type cUchar C.uchar
type cUint C.uint
type cUint32T C.uint32_t
type cUint8T C.uint8_t
type cWasmerByteArray C.wasmer_byte_array
type cWasmerExportDescriptorT C.wasmer_export_descriptor_t
type cWasmerExportDescriptorsT C.wasmer_export_descriptors_t
type cWasmerExportFuncT C.wasmer_export_func_t
type cWasmerExportT C.wasmer_export_t
type cWasmerExportsT C.wasmer_exports_t
type cWasmerImportDescriptorT C.wasmer_import_descriptor_t
type cWasmerImportDescriptorsT C.wasmer_import_descriptors_t
type cWasmerImportExportKind C.wasmer_import_export_kind
type cWasmerImportExportValue C.wasmer_import_export_value
type cWasmerImportFuncT C.wasmer_import_func_t
type cWasmerImportT C.wasmer_import_t
type cWasmerImportObjectT C.wasmer_import_object_t
type cWasmerInstanceContextT C.wasmer_instance_context_t
type cWasmerInstanceT C.wasmer_instance_t
type cWasmerMemoryT C.wasmer_memory_t
type cWasmerModuleT C.wasmer_module_t
type cWasmerResultT C.wasmer_result_t
type cWasmerSerializedModuleT C.wasmer_serialized_module_t
type cWasmerValueT C.wasmer_value_t
type cWasmerValueTag C.wasmer_value_tag
type cWasmerCompilationOptions C.wasmer_compilation_options_t

const cWasmFunction = C.WASM_FUNCTION
const cWasmGlobal = C.WASM_GLOBAL
const cWasmI32 = C.WASM_I32
const cWasmI64 = C.WASM_I64
const cWasmMemory = C.WASM_MEMORY
const cWasmTable = C.WASM_TABLE
const cWasmerOk = C.WASMER_OK

func cNewWasmerImportT(moduleName string, importName string, function *cWasmerImportFuncT) cWasmerImportT {
	var importedFunction C.wasmer_import_t
	importedFunction.module_name = (C.wasmer_byte_array)(cGoStringToWasmerByteArray(moduleName))
	importedFunction.import_name = (C.wasmer_byte_array)(cGoStringToWasmerByteArray(importName))
	importedFunction.tag = cWasmFunction

	var pointer = (**C.wasmer_import_func_t)(unsafe.Pointer(&importedFunction.value))
	*pointer = (*C.wasmer_import_func_t)(function)

	return (cWasmerImportT)(importedFunction)
}

func cWasmerInstanceGetPointsUsed(instance *cWasmerInstanceT) uint64 {
	return uint64(C.wasmer_instance_get_points_used(
		(*C.wasmer_instance_t)(instance),
	))
}

func cWasmerInstanceSetPointsUsed(instance *cWasmerInstanceT, points uint64) {
	C.wasmer_instance_set_points_used(
		(*C.wasmer_instance_t)(instance),
		(C.uint64_t)(points),
	)
}

func cWasmerInstanceSetGasLimit(instance *cWasmerInstanceT, gasLimit uint64) {
	C.wasmer_instance_set_points_limit(
		(*C.wasmer_instance_t)(instance),
		(C.uint64_t)(gasLimit),
	)
}

func cWasmerInstanceSetBreakpointValue(instance *cWasmerInstanceT, value uint64) {
	C.wasmer_instance_set_runtime_breakpoint_value(
		(*C.wasmer_instance_t)(instance),
		(C.uint64_t)(value),
	)
}

func cWasmerInstanceGetBreakpointValue(instance *cWasmerInstanceT) uint64 {
	return uint64(C.wasmer_instance_get_runtime_breakpoint_value(
		(*C.wasmer_instance_t)(instance),
	))
}

func cWasmerInstanceCache(
	instance *cWasmerInstanceT,
	cacheBytes **cUchar,
	cacheLen *cUint32T,
) cWasmerResultT {
	return (cWasmerResultT)(C.wasmer_instance_cache(
		(*C.wasmer_instance_t)(instance),
		(**C.uchar)(unsafe.Pointer(cacheBytes)),
		(*C.uint32_t)(cacheLen),
	))
}

func cWasmerInstanceFromCache(
	instance **cWasmerInstanceT,
	cacheBytes *cUchar,
	cacheLen cUint32T,
	options *cWasmerCompilationOptions,
) cWasmerResultT {
	return (cWasmerResultT)(C.wasmer_instance_from_cache(
		(**C.wasmer_instance_t)(unsafe.Pointer(instance)),
		(*C.uchar)(cacheBytes),
		(C.uint32_t)(cacheLen),
		(*C.wasmer_compilation_options_t)(options),
	))
}

func cWasmerCacheImportObjectSet(
	imports *cWasmerImportT,
	importsLength cInt,
	index uint64,
) cWasmerResultT {
	return (cWasmerResultT)(C.wasmer_import_object_set(
		(*C.wasmer_import_t)(imports),
		(C.uint32_t)(importsLength),
		(C.uint64_t)(index),
	))
}

func cWasmerImportObjectNewIndex() uint64 {
	return uint64(C.wasmer_import_object_new_index())
}

func cWasmerSetOpcodeCosts(opcode_costs *[OPCODE_COUNT]uint32) {
	C.wasmer_set_opcode_costs(
		(*C.uint32_t)(unsafe.Pointer(opcode_costs)),
	)
}

func cWasmerExportFuncParams(function *cWasmerExportFuncT, parameters *cWasmerValueTag, parametersLength cUint32T) cWasmerResultT {
	return (cWasmerResultT)(C.wasmer_export_func_params(
		(*C.wasmer_export_func_t)(function),
		(*C.wasmer_value_tag)(parameters),
		(C.uint32_t)(parametersLength),
	))
}

func cWasmerExportFuncParamsArity(function *cWasmerExportFuncT, result *cUint32T) cWasmerResultT {
	return (cWasmerResultT)(C.wasmer_export_func_params_arity(
		(*C.wasmer_export_func_t)(function),
		(*C.uint32_t)(result),
	))
}

func cWasmerExportFuncResultsArity(function *cWasmerExportFuncT, result *cUint32T) cWasmerResultT {
	return (cWasmerResultT)(C.wasmer_export_func_returns_arity(
		(*C.wasmer_export_func_t)(function),
		(*C.uint32_t)(result),
	))
}

func cWasmerExportKind(export *cWasmerExportT) cWasmerImportExportKind {
	return (cWasmerImportExportKind)(C.wasmer_export_kind(
		(*C.wasmer_export_t)(export),
	))
}

func cWasmerExportName(export *cWasmerExportT) cWasmerByteArray {
	return (cWasmerByteArray)(C.wasmer_export_name(
		(*C.wasmer_export_t)(export),
	))
}

func cWasmerExportToFunc(export *cWasmerExportT) *cWasmerExportFuncT {
	return (*cWasmerExportFuncT)(C.wasmer_export_to_func(
		(*C.wasmer_export_t)(export),
	))
}

func cWasmerExportToMemory(export *cWasmerExportT, memory **cWasmerMemoryT) cWasmerResultT {
	return (cWasmerResultT)(C.wasmer_export_to_memory(
		(*C.wasmer_export_t)(export),
		(**C.wasmer_memory_t)(unsafe.Pointer(memory)),
	))
}

func cWasmerExportsDestroy(exports *cWasmerExportsT) {
	C.wasmer_exports_destroy(
		(*C.wasmer_exports_t)(exports),
	)
}

func cWasmerExportsGet(exports *cWasmerExportsT, index cInt) *cWasmerExportT {
	return (*cWasmerExportT)(C.wasmer_exports_get(
		(*C.wasmer_exports_t)(exports),
		(C.int)(index),
	))
}

func cWasmerExportsLen(exports *cWasmerExportsT) cInt {
	return (cInt)(C.wasmer_exports_len(
		(*C.wasmer_exports_t)(exports),
	))
}

func cWasmerImportFuncDestroy(function *cWasmerImportFuncT) {
	C.wasmer_import_func_destroy(
		(*C.wasmer_import_func_t)(function),
	)
}

func cWasmerImportFuncNew(
	function unsafe.Pointer,
	parametersSignature *cWasmerValueTag,
	parametersLength cUint,
	resultsSignature *cWasmerValueTag,
	resultsLength cUint,
) *cWasmerImportFuncT {
	return (*cWasmerImportFuncT)(C.wasmer_import_func_new(
		(*[0]byte)(function),
		(*C.wasmer_value_tag)(parametersSignature),
		(C.uint)(parametersLength),
		(*C.wasmer_value_tag)(resultsSignature),
		(C.uint)(resultsLength),
	))
}

func cWasmerInstanceCall(
	instance *cWasmerInstanceT,
	name *cChar,
	parameters *cWasmerValueT,
	parametersLength cUint32T,
	results *cWasmerValueT,
	resultsLength cUint32T,
) cWasmerResultT {
	return (cWasmerResultT)(C.wasmer_instance_call(
		(*C.wasmer_instance_t)(instance),
		(*C.char)(name),
		(*C.wasmer_value_t)(parameters),
		(C.uint32_t)(parametersLength),
		(*C.wasmer_value_t)(results),
		(C.uint32_t)(resultsLength),
	))
}

func cWasmerInstanceContextGet(instance *cWasmerInstanceT) *cWasmerInstanceContextT {
	return (*cWasmerInstanceContextT)(C.wasmer_instance_context_get(
		(*C.wasmer_instance_t)(instance),
	))
}

func cWasmerInstanceContextDataGet(instanceContext *cWasmerInstanceContextT) unsafe.Pointer {
	return unsafe.Pointer(C.wasmer_instance_context_data_get(
		(*C.wasmer_instance_context_t)(instanceContext),
	))
}

func cWasmerInstanceContextDataSet(instance *cWasmerInstanceT, dataPointer unsafe.Pointer) {
	C.wasmer_instance_context_data_set(
		(*C.wasmer_instance_t)(instance),
		dataPointer,
	)
}

func cWasmerInstanceContextMemory(instanceContext *cWasmerInstanceContextT) *cWasmerMemoryT {
	return (*cWasmerMemoryT)(C.wasmer_instance_context_memory(
		(*C.wasmer_instance_context_t)(instanceContext),
		0,
	))
}

func cWasmerInstanceDestroy(instance *cWasmerInstanceT) {
	C.wasmer_instance_destroy(
		(*C.wasmer_instance_t)(instance),
	)
}

func cWasmerInstanceExports(instance *cWasmerInstanceT, exports **cWasmerExportsT) {
	C.wasmer_instance_exports(
		(*C.wasmer_instance_t)(instance),
		(**C.wasmer_exports_t)(unsafe.Pointer(exports)),
	)
}

func cWasmerInstantiateWithOptions(
	instance **cWasmerInstanceT,
	wasmBytes *cUchar,
	wasmBytesLength cUint,
	options *cWasmerCompilationOptions,
) cWasmerResultT {
	return (cWasmerResultT)(C.wasmer_instantiate_with_options(
		(**C.wasmer_instance_t)(unsafe.Pointer(instance)),
		(*C.uchar)(wasmBytes),
		(C.uint)(wasmBytesLength),
		(*C.wasmer_compilation_options_t)(options),
	))
}

func cWasmerLastErrorLength() cInt {
	return (cInt)(C.wasmer_last_error_length())
}

func cWasmerLastErrorMessage(buffer *cChar, length cInt) cInt {
	return (cInt)(C.wasmer_last_error_message(
		(*C.char)(buffer),
		(C.int)(length),
	))
}

func cWasmerMemoryData(memory *cWasmerMemoryT) *cUint8T {
	return (*cUint8T)(C.wasmer_memory_data(
		(*C.wasmer_memory_t)(memory),
	))
}

func cWasmerMemoryDataLength(memory *cWasmerMemoryT) cUint32T {
	return (cUint32T)(C.wasmer_memory_data_length(
		(*C.wasmer_memory_t)(memory),
	))
}

func cWasmerMemoryGrow(memory *cWasmerMemoryT, numberOfPages cUint32T) cWasmerResultT {
	return (cWasmerResultT)(C.wasmer_memory_grow(
		(*C.wasmer_memory_t)(memory),
		(C.uint32_t)(numberOfPages),
	))
}

func cWasmerMemoryDestroy(memory *cWasmerMemoryT) {
	C.wasmer_memory_destroy(
		(*C.wasmer_memory_t)(memory),
	)
}

func cCString(string string) *cChar {
	return (*cChar)(C.CString(string))
}

func cFree(pointer unsafe.Pointer) {
	C.free(pointer)
}

func cGoString(string *cChar) string {
	return C.GoString((*C.char)(string))
}

func cGoStringN(string *cChar, length cInt) string {
	return C.GoStringN((*C.char)(string), (C.int)(length))
}

func cGoStringToWasmerByteArray(string string) cWasmerByteArray {
	var cString = cCString(string)

	var byteArray cWasmerByteArray
	byteArray.bytes = (*C.uchar)(unsafe.Pointer(cString))
	byteArray.bytes_len = (C.uint)(len(string))

	return byteArray
}
