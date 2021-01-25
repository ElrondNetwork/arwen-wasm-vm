package arwen

import (
	"sync"
	"unsafe"

	"github.com/ElrondNetwork/arwen-wasm-vm/crypto"
	"github.com/ElrondNetwork/arwen-wasm-vm/wasmer"
)

const (
	// AddressLen specifies the length of the address
	AddressLen = 32

	// HashLen specifies the lenghth of a hash
	HashLen = 32

	// BalanceLen specifies the number of bytes on which the balance is stored
	BalanceLen = 32

	// CodeMetadataLen specifies the length of the code metadata
	CodeMetadataLen = 2

	// InitFunctionName specifies the name for the init function
	InitFunctionName = "init"

	// InitFunctionNameEth specifies the name for the init function on Ethereum
	InitFunctionNameEth = "solidity.ctor"

	// UpgradeFunctionName specifies if the call is an upgradeContract call
	UpgradeFunctionName = "upgradeContract"
)

var (
	vmContextCounter uint8
	vmContextMap     map[uint8]VMHost
)

var vmContextMapMutex = sync.Mutex{}

// AddHostContext adds the given context to the context map, and returns the context id
func AddHostContext(ctx VMHost) int {
	vmContextMapMutex.Lock()
	defer vmContextMapMutex.Unlock()

	id := vmContextCounter
	vmContextCounter++
	if vmContextMap == nil {
		vmContextMap = make(map[uint8]VMHost)
	}
	vmContextMap[id] = ctx
	return int(id)
}

// RemoveAllHostContexts reinitializes the vm context map
func RemoveAllHostContexts() {
	vmContextMapMutex.Lock()
	defer vmContextMapMutex.Unlock()
	vmContextMap = make(map[uint8]VMHost)
}

// RemoveHostContext deletes the context at the given id from the map
func RemoveHostContext(idx int) {
	vmContextMapMutex.Lock()
	defer vmContextMapMutex.Unlock()
	delete(vmContextMap, uint8(idx))
}

// GetVMContext returns the vm Context from the vm context map
func GetVMContext(context unsafe.Pointer) VMHost {
	vmContextMapMutex.Lock()
	defer vmContextMapMutex.Unlock()

	instCtx := wasmer.IntoInstanceContext(context)
	var idx = *(*int)(instCtx.Data())

	return vmContextMap[uint8(idx)]
}

// GetBlockchainContext returns the blockchain context
func GetBlockchainContext(context unsafe.Pointer) BlockchainContext {
	return GetVMContext(context).Blockchain()
}

// GetRuntimeContext returns the runtime context
func GetRuntimeContext(context unsafe.Pointer) RuntimeContext {
	return GetVMContext(context).Runtime()
}

// GetCryptoContext returns the crypto context
func GetCryptoContext(context unsafe.Pointer) crypto.VMCrypto {
	return GetVMContext(context).Crypto()
}

// GetBigIntContext returns the big int context
func GetBigIntContext(context unsafe.Pointer) BigIntContext {
	return GetVMContext(context).BigInt()
}

// GetOutputContext returns the output context
func GetOutputContext(context unsafe.Pointer) OutputContext {
	return GetVMContext(context).Output()
}

// GetMeteringContext returns the metering context
func GetMeteringContext(context unsafe.Pointer) MeteringContext {
	return GetVMContext(context).Metering()
}

// GetStorageContext returns the storage context
func GetStorageContext(context unsafe.Pointer) StorageContext {
	return GetVMContext(context).Storage()
}
