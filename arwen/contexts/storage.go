package contexts

import (
	"bytes"
	"errors"

	"github.com/ElrondNetwork/arwen-wasm-vm/v1_3/arwen"
	"github.com/ElrondNetwork/arwen-wasm-vm/v1_3/math"
	logger "github.com/ElrondNetwork/elrond-go-logger"
	"github.com/ElrondNetwork/elrond-go-logger/check"
	vmcommon "github.com/ElrondNetwork/elrond-vm-common"
)

var logStorage = logger.GetOrCreate("arwen/storage")

type storageContext struct {
	host                          arwen.VMHost
	blockChainHook                vmcommon.BlockchainHook
	address                       []byte
	stateStack                    [][]byte
	elrondProtectedKeyPrefix      []byte
	arwenStorageProtectionEnabled bool
}

// NewStorageContext creates a new storageContext
func NewStorageContext(
	host arwen.VMHost,
	blockChainHook vmcommon.BlockchainHook,
	elrondProtectedKeyPrefix []byte,
) (*storageContext, error) {
	if len(elrondProtectedKeyPrefix) == 0 {
		return nil, errors.New("elrondProtectedKeyPrefix cannot be empty")
	}
	context := &storageContext{
		host:                          host,
		blockChainHook:                blockChainHook,
		stateStack:                    make([][]byte, 0),
		elrondProtectedKeyPrefix:      elrondProtectedKeyPrefix,
		arwenStorageProtectionEnabled: true,
	}

	return context, nil
}

// InitState does nothing
func (context *storageContext) InitState() {
}

// PushState appends the current address to the state stack.
func (context *storageContext) PushState() {
	context.stateStack = append(context.stateStack, context.address)
}

// PopSetActiveState removes the latest entry from the state stack and sets it as the current address
func (context *storageContext) PopSetActiveState() {
	stateStackLen := len(context.stateStack)
	if stateStackLen == 0 {
		return
	}

	prevAddress := context.stateStack[stateStackLen-1]
	context.stateStack = context.stateStack[:stateStackLen-1]

	context.address = prevAddress
}

// PopDiscard removes the latest entry from the state stack
func (context *storageContext) PopDiscard() {
	stateStackLen := len(context.stateStack)
	if stateStackLen == 0 {
		return
	}

	context.stateStack = context.stateStack[:stateStackLen-1]
}

// ClearStateStack clears the state stack from the current context.
func (context *storageContext) ClearStateStack() {
	context.stateStack = make([][]byte, 0)
}

// SetAddress sets the given address as the address for the current context.
func (context *storageContext) SetAddress(address []byte) {
	context.address = address
	logStorage.Trace("storage under address set", "address", address)
}

// GetStorageUpdates returns the storage updates for the account mapped to the given address.
func (context *storageContext) GetStorageUpdates(address []byte) map[string]*vmcommon.StorageUpdate {
	account, _ := context.host.Output().GetOutputAccount(address)
	return account.StorageUpdates
}

// GetStorage returns the storage data mapped to the given key.
func (context *storageContext) GetStorage(key []byte) []byte {
	metering := context.host.Metering()

	extraBytes := len(key) - arwen.AddressLen
	if extraBytes > 0 {
		gasToUse := math.MulUint64(metering.GasSchedule().BaseOperationCost.DataCopyPerByte, uint64(extraBytes))
		metering.UseGas(gasToUse)
	}

	value := context.GetStorageUnmetered(key)

	gasToUse := math.MulUint64(metering.GasSchedule().BaseOperationCost.DataCopyPerByte, uint64(len(value)))
	metering.UseGas(gasToUse)

	logStorage.Trace("get", "key", key, "value", value)

	return value
}

// GetStorageFromAddress returns the data under the given key from the account mapped to the given address.
func (context *storageContext) GetStorageFromAddress(address []byte, key []byte) []byte {
	metering := context.host.Metering()

	extraBytes := len(key) - arwen.AddressLen
	if extraBytes > 0 {
		gasToUse := math.MulUint64(metering.GasSchedule().BaseOperationCost.DataCopyPerByte, uint64(extraBytes))
		metering.UseGas(gasToUse)
	}

	if !bytes.Equal(address, context.address) {
		userAcc, err := context.blockChainHook.GetUserAccount(address)
		if err != nil || check.IfNil(userAcc) {
			return nil
		}

		metadata := vmcommon.CodeMetadataFromBytes(userAcc.GetCodeMetadata())
		if !metadata.Readable {
			return nil
		}
	}

	// If the requested key is protected by the Elrond node, the stored value
	// could have been changed by a built-in function in the meantime, even if
	// contracts themselves cannot change protected values. Values stored under
	// protected keys must always be retrieved from the node, not from the cached
	// StorageUpdates.
	var value []byte
	if context.isElrondReservedKey(key) {
		value, _ = context.blockChainHook.GetStorageData(address, key)
	} else {
		value = context.getStorageFromAddressUnmetered(address, key)
	}

	costPerByte := metering.GasSchedule().BaseOperationCost.DataCopyPerByte
	gasToUse := math.MulUint64(costPerByte, uint64(len(value)))
	metering.UseGas(gasToUse)

	logStorage.Trace("get from address", "address", address, "key", key, "value", value)
	return value
}

func (context *storageContext) getStorageFromAddressUnmetered(address []byte, key []byte) []byte {
	var value []byte

	storageUpdates := context.GetStorageUpdates(address)
	if storageUpdate, ok := storageUpdates[string(key)]; ok {
		value = storageUpdate.Data
	} else {
		value, _ = context.blockChainHook.GetStorageData(address, key)
		storageUpdates[string(key)] = &vmcommon.StorageUpdate{
			Offset: key,
			Data:   value,
		}
	}

	return value
}

// GetStorageUnmetered returns the data under the given key.
func (context *storageContext) GetStorageUnmetered(key []byte) []byte {
	return context.getStorageFromAddressUnmetered(context.address, key)
}

// enableStorageProtection will prevent writing to protected keys
func (context *storageContext) enableStorageProtection() {
	context.arwenStorageProtectionEnabled = true
}

// disableStorageProtection will prevent writing to protected keys
func (context *storageContext) disableStorageProtection() {
	context.arwenStorageProtectionEnabled = false
}

func (context *storageContext) isArwenProtectedKey(key []byte) bool {
	return bytes.HasPrefix(key, []byte(arwen.ProtectedStoragePrefix))
}

func (context *storageContext) isElrondReservedKey(key []byte) bool {
	prefixElrond := bytes.HasPrefix(key, context.elrondProtectedKeyPrefix)
	return prefixElrond
}

func (context *storageContext) SetProtectedStorage(key []byte, value []byte) (arwen.StorageStatus, error) {
	context.disableStorageProtection()
	defer context.enableStorageProtection()

	return context.SetStorage(key, value)
}

// SetStorage sets the given value at the given key.
func (context *storageContext) SetStorage(key []byte, value []byte) (arwen.StorageStatus, error) {
	if context.host.Runtime().ReadOnly() {
		logStorage.Trace("storage set", "error", "cannot set storage in readonly mode")
		return arwen.StorageUnchanged, nil
	}
	if context.isElrondReservedKey(key) {
		logStorage.Trace("storage set", "error", arwen.ErrStoreElrondReservedKey, "key", key)
		return arwen.StorageUnchanged, arwen.ErrStoreElrondReservedKey
	}
	if context.isArwenProtectedKey(key) && context.arwenStorageProtectionEnabled {
		logStorage.Trace("storage set", "error", arwen.ErrCannotWriteProtectedKey, "key", key)
		return arwen.StorageUnchanged, arwen.ErrCannotWriteProtectedKey
	}

	metering := context.host.Metering()

	extraBytes := len(key) - arwen.AddressLen
	if extraBytes > 0 {
		gasToUse := math.MulUint64(metering.GasSchedule().BaseOperationCost.DataCopyPerByte, uint64(extraBytes))
		metering.UseGas(gasToUse)
	}

	var zero []byte
	strKey := string(key)
	length := len(value)

	var oldValue []byte
	storageUpdates := context.GetStorageUpdates(context.address)
	if update, ok := storageUpdates[strKey]; !ok {
		oldValue = context.GetStorageUnmetered(key)
		storageUpdates[strKey] = &vmcommon.StorageUpdate{
			Offset: key,
			Data:   oldValue,
		}
	} else {
		oldValue = update.Data
	}

	lengthOldValue := len(oldValue)
	if bytes.Equal(oldValue, value) {
		useGas := math.MulUint64(metering.GasSchedule().BaseOperationCost.DataCopyPerByte, uint64(length))
		metering.UseGas(useGas)
		logStorage.Trace("storage set to identical value")
		return arwen.StorageUnchanged, nil
	}

	newUpdate := &vmcommon.StorageUpdate{
		Offset: key,
		Data:   make([]byte, length),
	}
	copy(newUpdate.Data[:length], value[:length])
	storageUpdates[strKey] = newUpdate

	if bytes.Equal(oldValue, zero) {
		useGas := math.MulUint64(metering.GasSchedule().BaseOperationCost.StorePerByte, uint64(length))
		metering.UseGas(useGas)
		logStorage.Trace("storage added", "key", key, "value", value)
		return arwen.StorageAdded, nil
	}
	if bytes.Equal(value, zero) {
		freeGas := math.MulUint64(metering.GasSchedule().BaseOperationCost.ReleasePerByte, uint64(lengthOldValue))
		metering.FreeGas(freeGas)
		logStorage.Trace("storage deleted", "key", key)
		return arwen.StorageDeleted, nil
	}

	newValueExtraLength := math.SubInt(length, lengthOldValue)

	if newValueExtraLength > 0 {
		useGas := math.MulUint64(metering.GasSchedule().BaseOperationCost.PersistPerByte, uint64(lengthOldValue))
		newValStoreUseGas := math.MulUint64(metering.GasSchedule().BaseOperationCost.StorePerByte, uint64(newValueExtraLength))
		gasUsed := math.AddUint64(useGas, newValStoreUseGas)

		metering.UseGas(gasUsed)
	}

	if newValueExtraLength < 0 {
		newValueExtraLength = -newValueExtraLength

		useGas := math.MulUint64(metering.GasSchedule().BaseOperationCost.PersistPerByte, uint64(length))
		metering.UseGas(useGas)

		freeGas := math.MulUint64(metering.GasSchedule().BaseOperationCost.ReleasePerByte, uint64(newValueExtraLength))
		metering.FreeGas(freeGas)
	}

	logStorage.Trace("storage modified", "key", key, "value", value, "lengthDelta", newValueExtraLength)
	return arwen.StorageModified, nil
}
