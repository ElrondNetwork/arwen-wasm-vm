package contexts

import (
	"bytes"
	"errors"
	"math/big"
	"testing"

	"github.com/ElrondNetwork/arwen-wasm-vm/v1_3/arwen"
	"github.com/ElrondNetwork/arwen-wasm-vm/v1_3/config"
	contextmock "github.com/ElrondNetwork/arwen-wasm-vm/v1_3/mock/context"
	worldmock "github.com/ElrondNetwork/arwen-wasm-vm/v1_3/mock/world"
	"github.com/ElrondNetwork/elrond-vm-common"
	"github.com/stretchr/testify/require"
)

var elrondReservedTestPrefix = []byte("RESERVED")

func TestNewStorageContext(t *testing.T) {
	t.Parallel()

	host := &contextmock.VMHostMock{}
	mockBlockchain := worldmock.NewMockWorld()

	storageContext, err := NewStorageContext(host, mockBlockchain, elrondReservedTestPrefix)
	require.Nil(t, err)
	require.NotNil(t, storageContext)
}

func TestStorageContext_SetAddress(t *testing.T) {
	t.Parallel()

	addressA := []byte("accountA")
	addressB := []byte("accountB")
	stubOutput := &contextmock.OutputContextStub{}
	accountA := &vmcommon.OutputAccount{
		Address:        addressA,
		Nonce:          0,
		BalanceDelta:   big.NewInt(0),
		Balance:        big.NewInt(0),
		StorageUpdates: make(map[string]*vmcommon.StorageUpdate),
	}
	accountB := &vmcommon.OutputAccount{
		Address:        addressB,
		Nonce:          0,
		BalanceDelta:   big.NewInt(0),
		Balance:        big.NewInt(0),
		StorageUpdates: make(map[string]*vmcommon.StorageUpdate),
	}
	stubOutput.GetOutputAccountCalled = func(address []byte) (*vmcommon.OutputAccount, bool) {
		if bytes.Equal(address, addressA) {
			return accountA, false
		}
		if bytes.Equal(address, addressB) {
			return accountB, false
		}
		return nil, false
	}

	mockRuntime := &contextmock.RuntimeContextMock{}
	mockMetering := &contextmock.MeteringContextMock{}
	mockMetering.SetGasSchedule(config.MakeGasMapForTests())
	mockMetering.BlockGasLimitMock = uint64(15000)

	host := &contextmock.VMHostMock{
		OutputContext:   stubOutput,
		MeteringContext: mockMetering,
		RuntimeContext:  mockRuntime,
	}
	bcHook := &contextmock.BlockchainHookStub{}

	storageContext, _ := NewStorageContext(host, bcHook, elrondReservedTestPrefix)

	keyA := []byte("keyA")
	valueA := []byte("valueA")

	storageContext.SetAddress(addressA)
	storageStatus, err := storageContext.SetStorage(keyA, valueA)
	require.Nil(t, err)
	require.Equal(t, arwen.StorageAdded, storageStatus)
	require.Equal(t, valueA, storageContext.GetStorage(keyA))
	require.Len(t, storageContext.GetStorageUpdates(addressA), 1)
	require.Len(t, storageContext.GetStorageUpdates(addressB), 0)

	keyB := []byte("keyB")
	valueB := []byte("valueB")
	storageContext.SetAddress(addressB)
	storageStatus, err = storageContext.SetStorage(keyB, valueB)
	require.Nil(t, err)
	require.Equal(t, arwen.StorageAdded, storageStatus)
	require.Equal(t, valueB, storageContext.GetStorage(keyB))
	require.Len(t, storageContext.GetStorageUpdates(addressA), 1)
	require.Len(t, storageContext.GetStorageUpdates(addressB), 1)
	require.Equal(t, []byte(nil), storageContext.GetStorage(keyA))
}

func TestStorageContext_GetStorageUpdates(t *testing.T) {
	t.Parallel()

	mockOutput := &contextmock.OutputContextMock{}
	account := mockOutput.NewVMOutputAccount([]byte("account"))
	mockOutput.OutputAccountMock = account
	mockOutput.OutputAccountIsNew = false

	account.StorageUpdates["update"] = &vmcommon.StorageUpdate{
		Offset: []byte("update"),
		Data:   []byte("some data"),
	}

	host := &contextmock.VMHostMock{
		OutputContext: mockOutput,
	}

	mockBlockchainHook := worldmock.NewMockWorld()
	storageContext, _ := NewStorageContext(host, mockBlockchainHook, elrondReservedTestPrefix)

	storageUpdates := storageContext.GetStorageUpdates([]byte("account"))
	require.Equal(t, 1, len(storageUpdates))
	require.Equal(t, []byte("update"), storageUpdates["update"].Offset)
	require.Equal(t, []byte("some data"), storageUpdates["update"].Data)
}

func TestStorageContext_SetStorage(t *testing.T) {
	t.Parallel()

	address := []byte("account")
	mockOutput := &contextmock.OutputContextMock{}
	account := mockOutput.NewVMOutputAccount(address)
	mockOutput.OutputAccountMock = account
	mockOutput.OutputAccountIsNew = false

	mockRuntime := &contextmock.RuntimeContextMock{}
	mockMetering := &contextmock.MeteringContextMock{}
	mockMetering.SetGasSchedule(config.MakeGasMapForTests())
	mockMetering.BlockGasLimitMock = uint64(15000)

	host := &contextmock.VMHostMock{
		OutputContext:   mockOutput,
		MeteringContext: mockMetering,
		RuntimeContext:  mockRuntime,
	}
	bcHook := &contextmock.BlockchainHookStub{}

	storageContext, _ := NewStorageContext(host, bcHook, elrondReservedTestPrefix)
	storageContext.SetAddress(address)

	key := []byte("key")
	value := []byte("value")

	storageStatus, err := storageContext.SetStorage(key, value)
	require.Nil(t, err)
	require.Equal(t, arwen.StorageAdded, storageStatus)
	require.Equal(t, value, storageContext.GetStorage(key))
	require.Len(t, storageContext.GetStorageUpdates(address), 1)

	value = []byte("newValue")
	storageStatus, err = storageContext.SetStorage(key, value)
	require.Nil(t, err)
	require.Equal(t, arwen.StorageModified, storageStatus)
	require.Equal(t, value, storageContext.GetStorage(key))
	require.Len(t, storageContext.GetStorageUpdates(address), 1)

	value = []byte("newValue")
	storageStatus, err = storageContext.SetStorage(key, value)
	require.Nil(t, err)
	require.Equal(t, arwen.StorageUnchanged, storageStatus)
	require.Equal(t, value, storageContext.GetStorage(key))
	require.Len(t, storageContext.GetStorageUpdates(address), 1)

	value = nil
	storageStatus, err = storageContext.SetStorage(key, value)
	require.Nil(t, err)
	require.Equal(t, arwen.StorageDeleted, storageStatus)
	require.Equal(t, []byte{}, storageContext.GetStorage(key))
	require.Len(t, storageContext.GetStorageUpdates(address), 1)

	mockRuntime.SetReadOnly(true)
	value = []byte("newValue")
	storageStatus, err = storageContext.SetStorage(key, value)
	require.Nil(t, err)
	require.Equal(t, arwen.StorageUnchanged, storageStatus)
	require.Equal(t, []byte{}, storageContext.GetStorage(key))
	require.Len(t, storageContext.GetStorageUpdates(address), 1)

	mockRuntime.SetReadOnly(false)
	key = []byte("other_key")
	value = []byte("other_value")
	storageStatus, err = storageContext.SetStorage(key, value)
	require.Nil(t, err)
	require.Equal(t, arwen.StorageAdded, storageStatus)
	require.Equal(t, value, storageContext.GetStorage(key))
	require.Len(t, storageContext.GetStorageUpdates(address), 2)

	key = []byte("RESERVEDkey")
	value = []byte("doesn't matter")
	storageStatus, err = storageContext.SetStorage(key, value)
	require.Equal(t, arwen.StorageUnchanged, storageStatus)
	require.Equal(t, arwen.ErrStoreElrondReservedKey, err)

	key = []byte("RESERVED")
	value = []byte("doesn't matter")
	storageStatus, err = storageContext.SetStorage(key, value)
	require.Equal(t, arwen.StorageUnchanged, storageStatus)
	require.Equal(t, arwen.ErrStoreElrondReservedKey, err)
}

func TestStorageContext_StorageProtection(t *testing.T) {
	address := []byte("account")
	mockOutput := &contextmock.OutputContextMock{}
	account := mockOutput.NewVMOutputAccount(address)
	mockOutput.OutputAccountMock = account
	mockOutput.OutputAccountIsNew = false

	mockRuntime := &contextmock.RuntimeContextMock{}
	mockMetering := &contextmock.MeteringContextMock{}
	mockMetering.SetGasSchedule(config.MakeGasMapForTests())
	mockMetering.BlockGasLimitMock = uint64(15000)

	host := &contextmock.VMHostMock{
		OutputContext:   mockOutput,
		MeteringContext: mockMetering,
		RuntimeContext:  mockRuntime,
	}
	bcHook := &contextmock.BlockchainHookStub{}

	storageContext, _ := NewStorageContext(host, bcHook, elrondReservedTestPrefix)
	storageContext.SetAddress(address)

	key := []byte(arwen.ProtectedStoragePrefix + "something")
	value := []byte("data")

	storageStatus, err := storageContext.SetStorage(key, value)
	require.Equal(t, arwen.StorageUnchanged, storageStatus)
	require.True(t, errors.Is(err, arwen.ErrCannotWriteProtectedKey))
	require.Len(t, storageContext.GetStorageUpdates(address), 0)

	storageContext.disableStorageProtection()
	storageStatus, err = storageContext.SetStorage(key, value)
	require.Nil(t, err)
	require.Equal(t, arwen.StorageAdded, storageStatus)
	require.Len(t, storageContext.GetStorageUpdates(address), 1)

	storageContext.enableStorageProtection()
	storageStatus, err = storageContext.SetStorage(key, value)
	require.Equal(t, arwen.StorageUnchanged, storageStatus)
	require.True(t, errors.Is(err, arwen.ErrCannotWriteProtectedKey))
	require.Len(t, storageContext.GetStorageUpdates(address), 1)
}

func TestStorageContext_GetStorageFromAddress(t *testing.T) {
	t.Parallel()

	scAddress := []byte("account")
	mockOutput := &contextmock.OutputContextMock{}
	account := mockOutput.NewVMOutputAccount(scAddress)
	mockOutput.OutputAccountMock = account
	mockOutput.OutputAccountIsNew = false

	mockRuntime := &contextmock.RuntimeContextMock{}
	mockMetering := &contextmock.MeteringContextMock{}
	mockMetering.SetGasSchedule(config.MakeGasMapForTests())
	mockMetering.BlockGasLimitMock = uint64(15000)

	host := &contextmock.VMHostMock{
		OutputContext:   mockOutput,
		MeteringContext: mockMetering,
		RuntimeContext:  mockRuntime,
	}

	readable := []byte("readable")
	nonreadable := []byte("nonreadable")
	internalData := []byte("internalData")

	bcHook := &contextmock.BlockchainHookStub{
		GetUserAccountCalled: func(address []byte) (vmcommon.UserAccountHandler, error) {
			if bytes.Equal(readable, address) {
				return &worldmock.Account{CodeMetadata: []byte{4, 0}}, nil
			}
			if bytes.Equal(nonreadable, address) || bytes.Equal(scAddress, address) {
				return &worldmock.Account{CodeMetadata: []byte{0, 0}}, nil
			}
			return nil, nil
		},
		GetStorageDataCalled: func(accountsAddress []byte, index []byte) ([]byte, error) {
			return internalData, nil
		},
	}

	storageContext, _ := NewStorageContext(host, bcHook, elrondReservedTestPrefix)
	storageContext.SetAddress(scAddress)

	key := []byte("key")
	data := storageContext.GetStorageFromAddress(scAddress, key)
	require.Equal(t, data, internalData)

	data = storageContext.GetStorageFromAddress(readable, key)
	require.Equal(t, data, internalData)

	data = storageContext.GetStorageFromAddress(nonreadable, key)
	require.Nil(t, data)
}

func TestStorageContext_LoadGasStoreGasPerKey(t *testing.T) {
	// TODO
}

func TestStorageContext_StoreGasPerKey(t *testing.T) {
	// TODO
}

func TestStorageContext_PopSetActiveStateIfStackIsEmptyShouldNotPanic(t *testing.T) {
	t.Parallel()

	storageContext, _ := NewStorageContext(&contextmock.VMHostMock{}, &contextmock.BlockchainHookStub{}, elrondReservedTestPrefix)
	storageContext.PopSetActiveState()

	require.Equal(t, 0, len(storageContext.stateStack))
}

func TestStorageContext_PopDiscardIfStackIsEmptyShouldNotPanic(t *testing.T) {
	t.Parallel()

	storageContext, _ := NewStorageContext(&contextmock.VMHostMock{}, &contextmock.BlockchainHookStub{}, elrondReservedTestPrefix)
	storageContext.PopDiscard()

	require.Equal(t, 0, len(storageContext.stateStack))
}
