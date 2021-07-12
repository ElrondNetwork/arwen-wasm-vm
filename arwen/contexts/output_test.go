package contexts

import (
	"math/big"
	"testing"

	"github.com/ElrondNetwork/arwen-wasm-vm/v1_3/arwen"
	contextmock "github.com/ElrondNetwork/arwen-wasm-vm/v1_3/mock/context"
	worldmock "github.com/ElrondNetwork/arwen-wasm-vm/v1_3/mock/world"
	"github.com/ElrondNetwork/elrond-vm-common"
	"github.com/stretchr/testify/require"
)

func TestNewOutputContext(t *testing.T) {
	t.Parallel()

	host := &contextmock.VMHostStub{}

	outputContext, err := NewOutputContext(host)
	require.Nil(t, err)
	require.NotNil(t, outputContext)

	require.Equal(t, vmcommon.Ok, outputContext.ReturnCode())
	require.NotNil(t, outputContext.ReturnData())
	require.Equal(t, 0, len(outputContext.ReturnData()))
	require.Equal(t, "", outputContext.ReturnMessage())
	require.Equal(t, 0, len(outputContext.stateStack))

	require.Equal(t, uint64(0), outputContext.GetRefund())

	require.NotNil(t, outputContext.outputState.OutputAccounts)
	require.Equal(t, 0, len(outputContext.outputState.OutputAccounts))

	vmOutput := outputContext.outputState
	require.NotNil(t, vmOutput.DeletedAccounts)
	require.Equal(t, 0, len(vmOutput.DeletedAccounts))
	require.NotNil(t, vmOutput.TouchedAccounts)
	require.Equal(t, 0, len(vmOutput.TouchedAccounts))
	require.NotNil(t, vmOutput.Logs)
	require.Equal(t, 0, len(vmOutput.Logs))
}

func TestOutputContext_PushPopState(t *testing.T) {
	t.Parallel()

	host := &contextmock.VMHostStub{}
	host.RuntimeCalled = func() arwen.RuntimeContext {
		return &contextmock.RuntimeContextMock{VMInput: &vmcommon.VMInput{}}
	}
	outputContext, _ := NewOutputContext(host)

	address1 := []byte("address1")
	address2 := []byte("address2")

	// Create an account with nonce 99 on the active state.
	account, isNew := outputContext.GetOutputAccount(address1)
	account.Nonce = 99
	require.True(t, isNew)
	require.Equal(t, 1, len(outputContext.outputState.OutputAccounts))

	// Copy active state onto the stack.
	outputContext.PushState()
	require.Equal(t, 1, len(outputContext.stateStack))

	// Clear the active state and create a new account with the same address as
	// the previous; the new account must not have nonce 99.
	outputContext.InitState()
	account, isNew = outputContext.GetOutputAccount(address1)
	require.True(t, isNew)
	require.Equal(t, 1, len(outputContext.outputState.OutputAccounts))
	require.Equal(t, uint64(0), account.Nonce)

	account.Nonce = 84

	// Copy active state onto the stack, then create a new account with nonce 42.
	outputContext.PushState()
	require.Equal(t, 2, len(outputContext.stateStack))

	account, isNew = outputContext.GetOutputAccount(address2)
	account.Nonce = 42
	require.True(t, isNew)
	require.Equal(t, 2, len(outputContext.outputState.OutputAccounts))

	// Revert to the previous state: account with nonce 42 is lost, and the
	// account with "address1" has nonce 84.
	outputContext.PopSetActiveState()
	account, isNew = outputContext.GetOutputAccount(address1)
	require.False(t, isNew)
	require.Equal(t, uint64(84), account.Nonce)
	require.Equal(t, 1, len(outputContext.outputState.OutputAccounts))
	require.Equal(t, 1, len(outputContext.stateStack))

	outputContext.PushState()
	require.Equal(t, 2, len(outputContext.stateStack))

	outputContext.PopDiscard()
	require.Equal(t, 1, len(outputContext.stateStack))

	account, isNew = outputContext.GetOutputAccount(address1)
	require.False(t, isNew)
	require.Equal(t, uint64(84), account.Nonce)
	require.Equal(t, 1, len(outputContext.outputState.OutputAccounts))
	require.Equal(t, 1, len(outputContext.stateStack))

	outputContext.ClearStateStack()
	require.Equal(t, 0, len(outputContext.stateStack))
}

func TestOutputContext_GetOutputAccount(t *testing.T) {
	t.Parallel()

	host := &contextmock.VMHostStub{}
	outputContext, _ := NewOutputContext(host)
	require.Zero(t, len(outputContext.outputState.OutputAccounts))

	// Request an account that is missing from OutputAccounts
	account, isNew := outputContext.GetOutputAccount([]byte("account"))
	require.Equal(t, 1, len(outputContext.outputState.OutputAccounts))
	require.True(t, isNew)
	require.Equal(t, []byte("account"), account.Address)
	require.Zero(t, account.Nonce)
	require.Equal(t, arwen.Zero, account.BalanceDelta)
	require.Nil(t, account.Balance)
	require.Zero(t, len(account.StorageUpdates))

	// Change fields of the OutputAccount to ensure it will be returned on the
	// next call to GetOutputAccount(), from the OutputAccounts cache
	account.Address = []byte("changed address")
	account.Nonce = 88
	account.Balance = big.NewInt(94)
	cachedAccount, isNew := outputContext.GetOutputAccount([]byte("account"))
	require.False(t, isNew)
	require.Equal(t, []byte("changed address"), cachedAccount.Address)
	require.Equal(t, uint64(88), cachedAccount.Nonce)
	require.Equal(t, big.NewInt(94), cachedAccount.Balance)
	require.Zero(t, len(cachedAccount.StorageUpdates))
}

func TestOutputContext_GettersAndSetters(t *testing.T) {
	host := &contextmock.VMHostStub{}
	outputContext, _ := NewOutputContext(host)

	outputContext.SetRefund(24)
	require.Equal(t, uint64(24), outputContext.GetRefund())

	outputContext.SetReturnCode(vmcommon.ExecutionFailed)
	require.Equal(t, vmcommon.ExecutionFailed, outputContext.ReturnCode())

	outputContext.SetReturnMessage("rockets")
	require.Equal(t, "rockets", outputContext.ReturnMessage())
}

func TestOutputContext_FinishReturnData(t *testing.T) {
	host := &contextmock.VMHostStub{}
	outputContext, _ := NewOutputContext(host)

	require.Zero(t, len(outputContext.ReturnData()))

	outputContext.Finish([]byte("something"))
	expectedData := [][]byte{
		[]byte("something"),
	}
	require.Equal(t, expectedData, outputContext.ReturnData())

	outputContext.Finish([]byte("something else"))
	expectedData = append(expectedData, []byte("something else"))
	require.Equal(t, expectedData, outputContext.ReturnData())

	outputContext.Finish(big.NewInt(1234567).Bytes())
	expectedData = append(expectedData, big.NewInt(1234567).Bytes())
	require.Equal(t, expectedData, outputContext.ReturnData())

	outputContext.Finish([]byte{})
	expectedData = append(expectedData, []byte{})
	require.Equal(t, expectedData, outputContext.ReturnData())

	outputContext.ClearReturnData()
	require.Zero(t, len(outputContext.ReturnData()))
}

func TestOutputContext_MergeCompleteAccounts(t *testing.T) {
	t.Parallel()

	transfer1 := vmcommon.OutputTransfer{
		Value:    big.NewInt(0),
		GasLimit: 9999,
		Data:     []byte("data1"),
	}
	left := &vmcommon.OutputAccount{
		Address:         []byte("addr1"),
		Nonce:           1,
		Balance:         big.NewInt(1000),
		BalanceDelta:    big.NewInt(10000),
		StorageUpdates:  nil,
		Code:            []byte("code1"),
		OutputTransfers: []vmcommon.OutputTransfer{transfer1},
	}
	right := &vmcommon.OutputAccount{
		Address:         []byte("addr2"),
		Nonce:           2,
		Balance:         big.NewInt(2000),
		BalanceDelta:    big.NewInt(20000),
		StorageUpdates:  map[string]*vmcommon.StorageUpdate{"key": {Data: []byte("data"), Offset: []byte("offset")}},
		Code:            []byte("code2"),
		OutputTransfers: []vmcommon.OutputTransfer{transfer1, transfer1},
	}

	expected := &vmcommon.OutputAccount{
		Address:         []byte("addr2"),
		Nonce:           2,
		Balance:         big.NewInt(2000),
		BalanceDelta:    big.NewInt(20000),
		StorageUpdates:  map[string]*vmcommon.StorageUpdate{"key": {Data: []byte("data"), Offset: []byte("offset")}},
		Code:            []byte("code2"),
		OutputTransfers: []vmcommon.OutputTransfer{transfer1, transfer1},
	}

	mergeOutputAccounts(left, right)
	require.Equal(t, expected, left)
}

func TestOutputContext_MergeIncompleteAccounts(t *testing.T) {
	t.Parallel()

	left := &vmcommon.OutputAccount{}
	right := &vmcommon.OutputAccount{}
	expected := &vmcommon.OutputAccount{
		StorageUpdates: make(map[string]*vmcommon.StorageUpdate),
		BalanceDelta:   big.NewInt(0),
	}
	mergeOutputAccounts(left, right)
	require.Equal(t, expected, left)

	left = &vmcommon.OutputAccount{
		OutputTransfers: []vmcommon.OutputTransfer{{GasLimit: 92}},
	}
	right = &vmcommon.OutputAccount{
		BalanceDelta: big.NewInt(42),
	}
	expected = &vmcommon.OutputAccount{
		StorageUpdates:  make(map[string]*vmcommon.StorageUpdate),
		BalanceDelta:    big.NewInt(42),
		OutputTransfers: []vmcommon.OutputTransfer{{GasLimit: 92}},
	}
	mergeOutputAccounts(left, right)
	require.Equal(t, expected, left)

	left = &vmcommon.OutputAccount{
		BalanceDelta: big.NewInt(48),
	}
	right = &vmcommon.OutputAccount{}
	expected = &vmcommon.OutputAccount{
		StorageUpdates: make(map[string]*vmcommon.StorageUpdate),
		BalanceDelta:   big.NewInt(48),
	}
	mergeOutputAccounts(left, right)
	require.Equal(t, expected, left)

	left = &vmcommon.OutputAccount{
		Address: []byte("left address"),
		Code:    []byte("left code"),
	}
	right = &vmcommon.OutputAccount{
		BalanceDelta: big.NewInt(42),
	}
	expected = &vmcommon.OutputAccount{
		Code:           []byte("left code"),
		StorageUpdates: make(map[string]*vmcommon.StorageUpdate),
		BalanceDelta:   big.NewInt(42),
		Address:        left.Address,
	}
	mergeOutputAccounts(left, right)
	require.Equal(t, expected, left)

	left = &vmcommon.OutputAccount{
		OutputTransfers: []vmcommon.OutputTransfer{{Data: []byte("left data")}},
	}
	right = &vmcommon.OutputAccount{}
	expected = &vmcommon.OutputAccount{
		OutputTransfers: []vmcommon.OutputTransfer{{Data: []byte("left data")}},
		StorageUpdates:  make(map[string]*vmcommon.StorageUpdate),
		BalanceDelta:    big.NewInt(0),
	}
	mergeOutputAccounts(left, right)
	require.Equal(t, expected, left)

	left = &vmcommon.OutputAccount{
		Nonce: 44,
	}
	right = &vmcommon.OutputAccount{
		Nonce: 42,
	}
	expected = &vmcommon.OutputAccount{
		Nonce:          44,
		StorageUpdates: make(map[string]*vmcommon.StorageUpdate),
		BalanceDelta:   big.NewInt(0),
	}
	mergeOutputAccounts(left, right)
	require.Equal(t, expected, left)
}

func TestOutputContext_MergeVMOutputs(t *testing.T) {
	t.Parallel()

	left := newVMOutput()
	right := newVMOutput()
	expected := newVMOutput()
	mergeVMOutputs(left, right)
	require.Equal(t, expected, left)

	left = newVMOutput()
	right = newVMOutput()
	right.OutputAccounts["address"] = NewVMOutputAccount([]byte("address"))
	right.OutputAccounts["address"].Nonce = 84
	expected = newVMOutput()
	expected.OutputAccounts["address"] = NewVMOutputAccount([]byte("address"))
	expected.OutputAccounts["address"].Nonce = 84
	mergeVMOutputs(left, right)
	require.Equal(t, expected, left)

	left = newVMOutput()
	left.OutputAccounts["address"] = NewVMOutputAccount([]byte("address"))
	left.OutputAccounts["address"].Nonce = 84
	right = newVMOutput()
	right.OutputAccounts["address"] = NewVMOutputAccount([]byte("address"))
	right.OutputAccounts["address"].Nonce = 92
	expected = newVMOutput()
	expected.OutputAccounts["address"] = NewVMOutputAccount([]byte("address"))
	expected.OutputAccounts["address"].Nonce = 92
	mergeVMOutputs(left, right)
	require.Equal(t, expected, left)

	left = newVMOutput()
	left.OutputAccounts["left address"] = NewVMOutputAccount([]byte("left address"))
	right = newVMOutput()
	right.OutputAccounts["right address"] = NewVMOutputAccount([]byte("right address"))
	expected = newVMOutput()
	expected.OutputAccounts["left address"] = NewVMOutputAccount([]byte("left address"))
	expected.OutputAccounts["right address"] = NewVMOutputAccount([]byte("right address"))
	mergeVMOutputs(left, right)
	require.Equal(t, expected, left)

	left = newVMOutput()
	left.GasRemaining = 99
	left.GasRefund = big.NewInt(42)
	left.ReturnCode = vmcommon.ContractNotFound
	left.ReturnMessage = "fireflies"
	right = newVMOutput()
	right.ReturnData = [][]byte{[]byte("rockets"), []byte("albatross")}
	right.GasRemaining = 100
	right.GasRefund = big.NewInt(84)
	right.ReturnCode = vmcommon.ExecutionFailed
	right.ReturnMessage = "turtles all the way down"
	expected = newVMOutput()
	expected.ReturnData = [][]byte{[]byte("rockets"), []byte("albatross")}
	expected.GasRemaining = 100
	expected.GasRefund = big.NewInt(84)
	expected.ReturnCode = vmcommon.ExecutionFailed
	expected.ReturnMessage = "turtles all the way down"
	mergeVMOutputs(left, right)
	require.Equal(t, expected, left)
}

func TestOutputContext_VMOutputError(t *testing.T) {
	t.Parallel()

	host := &contextmock.VMHostMock{
		MeteringContext: &contextmock.MeteringContextMock{},
		RuntimeContext: &contextmock.RuntimeContextMock{
			VMInput: &vmcommon.VMInput{},
		},
	}

	outputContext, _ := NewOutputContext(host)

	returnCode := vmcommon.ContractNotFound
	returnMessage := arwen.ErrContractNotFound.Error()

	expected := &vmcommon.VMOutput{
		GasRemaining:  0,
		GasRefund:     big.NewInt(0),
		ReturnCode:    returnCode,
		ReturnMessage: returnMessage,
	}
	vmOutput := outputContext.CreateVMOutputInCaseOfError(arwen.ErrContractNotFound)
	require.Equal(t, expected, vmOutput)
}

func TestOutputContext_Transfer(t *testing.T) {
	t.Parallel()

	sender := []byte("sender")
	receiver := []byte("receiver")
	balance := big.NewInt(10000)
	valueToTransfer := big.NewInt(1000)

	host := &contextmock.VMHostMock{}
	host.RuntimeContext = &contextmock.RuntimeContextMock{VMInput: &vmcommon.VMInput{}}
	mockWorld := worldmock.NewMockWorld()
	mockWorld.AcctMap.PutAccount(&worldmock.Account{
		Address: sender,
		Nonce:   42,
		Balance: balance,
	})

	blockchainContext, _ := NewBlockchainContext(host, mockWorld)
	outputContext, _ := NewOutputContext(host)

	host.OutputContext = outputContext
	host.BlockchainContext = blockchainContext

	err := outputContext.Transfer(receiver, sender, 54, 0, valueToTransfer, []byte("txdata"), 0)
	require.Nil(t, err)

	senderAccount, isNew := outputContext.GetOutputAccount(sender)
	require.False(t, isNew)
	require.Equal(t, big.NewInt(-1000), senderAccount.BalanceDelta)

	destAccount, isNew := outputContext.GetOutputAccount(receiver)
	require.False(t, isNew)
	require.Equal(t, valueToTransfer, destAccount.BalanceDelta)
	require.Equal(t, uint64(54), destAccount.OutputTransfers[0].GasLimit)
	require.Equal(t, []byte("txdata"), destAccount.OutputTransfers[0].Data)
}

func TestOutputContext_Transfer_Errors_And_Checks(t *testing.T) {
	t.Parallel()

	sender := []byte("sender")
	receiver := []byte("receiver")

	mockWorld := worldmock.NewMockWorld()
	mockWorld.AcctMap.PutAccount(&worldmock.Account{
		Address: sender,
		Nonce:   88,
		Balance: big.NewInt(2000),
	})

	host := &contextmock.VMHostMock{}
	outputContext, _ := NewOutputContext(host)
	blockchainContext, _ := NewBlockchainContext(host, mockWorld)

	host.RuntimeContext = &contextmock.RuntimeContextMock{VMInput: &vmcommon.VMInput{}}
	host.OutputContext = outputContext
	host.BlockchainContext = blockchainContext

	senderOutputAccount, _ := outputContext.GetOutputAccount(sender)
	require.Nil(t, senderOutputAccount.Balance)
	require.Equal(t, arwen.Zero, senderOutputAccount.BalanceDelta)

	// negative transfers are disallowed
	valueToTransfer := big.NewInt(-1000)
	err := outputContext.Transfer(receiver, sender, 54, 0, valueToTransfer, []byte("txdata"), 0)
	require.Equal(t, arwen.ErrTransferNegativeValue, err)
	require.Nil(t, senderOutputAccount.Balance)
	require.Equal(t, arwen.Zero, senderOutputAccount.BalanceDelta)

	// account must have enough money to transfer
	valueToTransfer = big.NewInt(5000)
	err = outputContext.Transfer(receiver, sender, 54, 0, valueToTransfer, []byte("txdata"), 0)
	require.Equal(t, arwen.ErrTransferInsufficientFunds, err)
	require.Equal(t, big.NewInt(2000), senderOutputAccount.Balance)
	require.Equal(t, arwen.Zero, senderOutputAccount.BalanceDelta)

	senderOutputAccount.BalanceDelta = big.NewInt(4000)
	valueToTransfer = big.NewInt(5000)
	err = outputContext.Transfer(receiver, sender, 54, 0, valueToTransfer, []byte("txdata"), 0)
	require.Nil(t, err)
	require.Equal(t, big.NewInt(-1000), senderOutputAccount.BalanceDelta)

	require.Equal(t, big.NewInt(1000), blockchainContext.GetBalanceBigInt(sender))
}

func TestOutputContext_Transfer_IsAccountPayable(t *testing.T) {
	t.Parallel()

	sender := []byte("sender")
	receiverNonPayable := make([]byte, 32)
	receiverPayable := make([]byte, 32)
	receiverPayable[31] = 1

	mockWorld := worldmock.NewMockWorld()
	mockWorld.AcctMap.PutAccounts([]*worldmock.Account{
		{
			Address: sender,
			Nonce:   0,
			Balance: big.NewInt(2000),
		},
		{
			Address:         receiverNonPayable,
			Nonce:           0,
			Balance:         big.NewInt(0),
			Code:            []byte("contract_code"),
			IsSmartContract: true,
		},
		{
			Address:         receiverPayable,
			Nonce:           0,
			Balance:         big.NewInt(0),
			Code:            []byte("contract_code"),
			CodeMetadata:    []byte{0, vmcommon.MetadataPayable},
			IsSmartContract: true,
		},
	})

	host := &contextmock.VMHostMock{}
	oc, _ := NewOutputContext(host)
	bc, _ := NewBlockchainContext(host, mockWorld)

	host.OutputContext = oc
	host.BlockchainContext = bc
	host.RuntimeContext = &contextmock.RuntimeContextMock{VMInput: &vmcommon.VMInput{}}

	valueToTransfer := big.NewInt(10)
	err := oc.Transfer(receiverNonPayable, sender, 54, 0, valueToTransfer, []byte("txdata"), 0)

	require.Equal(t, arwen.ErrAccountNotPayable, err)

	valueToTransfer = big.NewInt(0)
	err = oc.Transfer(receiverNonPayable, sender, 54, 0, valueToTransfer, []byte("txdata"), 0)

	require.Nil(t, err)

	valueToTransfer = big.NewInt(10)
	err = oc.Transfer(receiverPayable, sender, 54, 0, valueToTransfer, []byte("txdata"), 0)

	require.Nil(t, err)
}

func TestOutputContext_WriteLog(t *testing.T) {
	t.Parallel()

	host := &contextmock.VMHostMock{
		RuntimeContext: &contextmock.RuntimeContextMock{},
	}
	outputContext, _ := NewOutputContext(host)

	address := []byte("address")
	data := []byte("data")
	topics := make([][]byte, 0)

	outputContext.WriteLog(address, topics, data)
	require.Equal(t, len(outputContext.outputState.Logs), 1)
	require.Equal(t, outputContext.outputState.Logs[0].Address, address)
	require.Equal(t, outputContext.outputState.Logs[0].Data, data)
	require.Empty(t, outputContext.outputState.Logs[0].Identifier)
	require.Empty(t, outputContext.outputState.Logs[0].Topics)

	identifier := []byte("identifier")
	topic := []byte("topic")
	topics = append(topics, identifier)
	outputContext.WriteLog(address, topics, data)

	require.Equal(t, outputContext.outputState.Logs[1].Identifier, identifier)
	require.Empty(t, outputContext.outputState.Logs[1].Topics)

	topics = append(topics, topic)
	outputContext.WriteLog(address, topics, data)

	require.Equal(t, outputContext.outputState.Logs[2].Topics, [][]byte{topic})
}

func TestOutputContext_PopSetActiveStateIfStackIsEmptyShouldNotPanic(t *testing.T) {
	t.Parallel()

	outputContext, _ := NewOutputContext(&contextmock.VMHostMock{})
	outputContext.PopSetActiveState()

	require.Equal(t, 0, len(outputContext.stateStack))
}

func TestOutputContext_PopMergeActiveStateIfStackIsEmptyShouldNotPanic(t *testing.T) {
	t.Parallel()

	outputContext, _ := NewOutputContext(&contextmock.VMHostMock{})
	outputContext.PopMergeActiveState()

	require.Equal(t, 0, len(outputContext.stateStack))
}

func TestOutputContext_PopDiscardIfStackIsEmptyShouldNotPanic(t *testing.T) {
	t.Parallel()

	outputContext, _ := NewOutputContext(&contextmock.VMHostMock{})
	outputContext.PopDiscard()

	require.Equal(t, 0, len(outputContext.stateStack))
}
