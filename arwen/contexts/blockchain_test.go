package contexts

import (
	"errors"
	"math/big"
	"testing"

	"github.com/ElrondNetwork/arwen-wasm-vm/v1_3/arwen"
	contextmock "github.com/ElrondNetwork/arwen-wasm-vm/v1_3/mock/context"
	worldmock "github.com/ElrondNetwork/arwen-wasm-vm/v1_3/mock/world"
	"github.com/ElrondNetwork/elrond-vm-common"
	"github.com/stretchr/testify/require"
)

var errTestError = errors.New("some test error")

var testAccounts = []*worldmock.Account{
	{Address: []byte("account_old"), Nonce: 12, Balance: big.NewInt(0)},
	{Address: []byte("account_newer"), Nonce: 8, Balance: big.NewInt(0)},
	{Address: []byte("account_new"), Nonce: 0, Balance: big.NewInt(0)},
	{Address: []byte("account_new_with_money"), Nonce: 0, Balance: big.NewInt(1000)},
	{Address: []byte("account_with_code"), Nonce: 4, Balance: big.NewInt(512), Code: []byte("somecode")},
	{Address: []byte("account_old_with_money"), Nonce: 56, Balance: big.NewInt(1024)},
}

func TestNewBlockchainContext(t *testing.T) {
	t.Parallel()

	host := &contextmock.VMHostStub{}
	mockWorld := worldmock.NewMockWorld()

	blockchainContext, err := NewBlockchainContext(host, mockWorld)
	require.Nil(t, err)
	require.NotNil(t, blockchainContext)
}

func TestBlockchainContext_AccountExists(t *testing.T) {
	t.Parallel()

	host := &contextmock.VMHostStub{}
	mockWorld := worldmock.NewMockWorld()
	mockWorld.AcctMap.PutAccounts(testAccounts)

	blockchainContext, _ := NewBlockchainContext(host, mockWorld)

	require.False(t, blockchainContext.AccountExists([]byte("account_missing")))
	require.False(t, blockchainContext.AccountExists([]byte("account_faulty")))
	require.True(t, blockchainContext.AccountExists([]byte("account_old")))

	mockWorld.Err = errTestError
	require.False(t, blockchainContext.AccountExists([]byte("account_something")))
}

func TestBlockchainContext_GetBalance(t *testing.T) {
	t.Parallel()

	mockWorld := worldmock.NewMockWorld()
	mockWorld.AcctMap.PutAccounts(testAccounts)
	mockOutput := &contextmock.OutputContextMock{}
	host := &contextmock.VMHostMock{}
	host.OutputContext = mockOutput
	blockchainContext, _ := NewBlockchainContext(host, mockWorld)

	// Act as if the OutputContext has no OutputAccounts cached
	// (mockOutput.GetOutputAccount() always returns "is new")
	account := &vmcommon.OutputAccount{}
	account.BalanceDelta = big.NewInt(0)
	mockOutput.OutputAccountMock = account
	mockOutput.OutputAccountIsNew = true

	// Test if error is propagated from BlockchainHook
	mockWorld.Err = errTestError
	balanceBytes := blockchainContext.GetBalance([]byte("account_new_with_money"))
	value := big.NewInt(0).SetBytes(balanceBytes)
	require.Equal(t, arwen.Zero, value)
	mockWorld.Err = nil

	// Test that an account that doesn't exist will not be updated with any kind
	// of balance in the OutputAccounts, but GetBalance() must return 0
	account.Balance = big.NewInt(15)
	balanceBytes = blockchainContext.GetBalance([]byte("account_missing"))
	value = big.NewInt(0).SetBytes(balanceBytes)
	require.Equal(t, arwen.Zero, value)
	require.Equal(t, big.NewInt(15), account.Balance)

	// Act as if the OutputContext has the requested OutputAccount cached
	account.Balance = big.NewInt(42)
	mockOutput.OutputAccountIsNew = false
	balanceBytes = blockchainContext.GetBalance([]byte("any account"))
	value = big.NewInt(0).SetBytes(balanceBytes)
	require.Equal(t, big.NewInt(42), value)

	// GetBalance must add Balance and BalanceDelta together
	account.Balance = big.NewInt(10)
	account.BalanceDelta = big.NewInt(32)
	mockOutput.OutputAccountIsNew = false
	balanceBytes = blockchainContext.GetBalance([]byte("any account"))
	value = big.NewInt(0).SetBytes(balanceBytes)
	require.Equal(t, big.NewInt(42), value)
	require.Equal(t, big.NewInt(10), account.Balance)
	require.Equal(t, big.NewInt(32), account.BalanceDelta)
}

func TestBlockchainContext_GetBalance_Updates(t *testing.T) {
	t.Parallel()

	mockWorld := worldmock.NewMockWorld()
	mockWorld.AcctMap.PutAccounts(testAccounts)
	mockOutput := &contextmock.OutputContextMock{}
	host := &contextmock.VMHostMock{}
	host.OutputContext = mockOutput
	blockchainContext, _ := NewBlockchainContext(host, mockWorld)

	// Act as if the OutputContext has no OutputAccounts cached
	// (mockOutput.GetOutputAccount() always returns "is new")
	account := &vmcommon.OutputAccount{
		Address:        []byte("account_new_with_money"),
		Nonce:          2,
		BalanceDelta:   big.NewInt(0),
		Balance:        nil,
		StorageUpdates: make(map[string]*vmcommon.StorageUpdate),
	}

	mockOutput.OutputAccountMock = account
	mockOutput.OutputAccountIsNew = false

	balanceBytes := blockchainContext.GetBalance([]byte("account_new_with_money"))
	value := big.NewInt(0).SetBytes(balanceBytes)
	require.Equal(t, big.NewInt(1000), value)
	require.Equal(t, big.NewInt(1000), account.Balance)
}

func TestBlockchainContext_GetNonceAndIncrease(t *testing.T) {
	t.Parallel()

	account := &vmcommon.OutputAccount{
		Nonce: 3,
	}

	host := &contextmock.VMHostMock{}

	mockOutput := &contextmock.OutputContextMock{}
	host.OutputContext = mockOutput

	mockWorld := worldmock.NewMockWorld()
	mockWorld.AcctMap.PutAccounts(testAccounts)
	blockchainContext, _ := NewBlockchainContext(host, mockWorld)

	// GetNonce: Test if error is propagated from BlockchainHook, and that the
	// cached OutputAccount doesn't lose its Nonce due to the error.
	mockOutput.OutputAccountMock = account
	mockOutput.OutputAccountIsNew = true
	mockWorld.Err = errTestError
	nonce, err := blockchainContext.GetNonce([]byte("any account"))
	require.Equal(t, errTestError, err)
	require.Equal(t, uint64(0), nonce)
	require.Equal(t, uint64(3), account.Nonce)
	mockWorld.Err = nil

	// GetNonce: Test requesting the nonce of an Account not yet cached by
	// OutputAccounts
	mockOutput.OutputAccountIsNew = true
	nonce, err = blockchainContext.GetNonce([]byte("account_old"))
	require.Equal(t, nil, err)
	require.Equal(t, uint64(12), nonce)
	require.Equal(t, uint64(12), account.Nonce)

	// GetNonce: Test requesting the nonce of an Account already cached by
	// OutputAccounts
	account.Nonce = 88
	mockOutput.OutputAccountIsNew = false
	nonce, err = blockchainContext.GetNonce([]byte("any account"))
	require.Nil(t, err)
	require.Equal(t, uint64(88), nonce)

	// IncreaseNonce: Test increasing the nonce of an Account already cached by
	// OutputAccounts
	account.Nonce = 88
	mockOutput.OutputAccountIsNew = false
	blockchainContext.IncreaseNonce([]byte("any account"))
	blockchainContext.IncreaseNonce([]byte("any account"))
	blockchainContext.IncreaseNonce([]byte("any account"))
	nonce, err = blockchainContext.GetNonce([]byte("any account"))
	require.Nil(t, err)
	require.Equal(t, uint64(91), nonce)
}

func TestBlockchainContext_GetCodeHashAndSize(t *testing.T) {
	t.Parallel()

	mockCrypto := &contextmock.CryptoHookMock{}

	mockWorld := worldmock.NewMockWorld()
	mockWorld.AcctMap.PutAccounts(testAccounts)

	outputContext := &contextmock.OutputContextMock{}

	host := &contextmock.VMHostMock{}
	host.CryptoHook = mockCrypto
	host.OutputContext = outputContext

	blockchainContext, _ := NewBlockchainContext(host, mockWorld)

	address := []byte("account_with_code")
	expectedCode := []byte("somecode")

	// GetCode: Test if error is propagated from blockchain hook
	outputContext.OutputAccountIsNew = true
	outputContext.OutputAccountMock = &vmcommon.OutputAccount{}
	mockWorld.Err = errTestError
	code, err := blockchainContext.GetCode(address)
	require.NotNil(t, err)
	require.Nil(t, code)
	mockWorld.Err = nil

	// GetCode: Test success
	outputContext.OutputAccountIsNew = true
	outputContext.OutputAccountMock = &vmcommon.OutputAccount{}
	code, err = blockchainContext.GetCode(address)
	require.Nil(t, err)
	require.Equal(t, expectedCode, code)

	// GetCodeHash: Test if error is propagated from blockchain hook
	outputContext.OutputAccountIsNew = true
	outputContext.OutputAccountMock = &vmcommon.OutputAccount{}
	mockWorld.Err = errTestError
	codeHash := blockchainContext.GetCodeHash(address)
	require.Nil(t, codeHash)
	mockWorld.Err = nil

	// GetCodeHash: Test success
	outputContext.OutputAccountIsNew = true
	outputContext.OutputAccountMock = &vmcommon.OutputAccount{}
	codeHash = blockchainContext.GetCodeHash(address)
	require.Equal(t, 0, len(codeHash))
	require.Nil(t, err)

	// GetCodeSize: Test if error is propagated from blockchain hook
	outputContext.OutputAccountIsNew = true
	outputContext.OutputAccountMock = &vmcommon.OutputAccount{}
	mockWorld.Err = errTestError
	size, err := blockchainContext.GetCodeSize(address)
	require.Equal(t, errTestError, err)
	require.Equal(t, int32(0), size)
	mockWorld.Err = nil

	// GetCodeSize: Test success
	outputContext.OutputAccountIsNew = true
	outputContext.OutputAccountMock = &vmcommon.OutputAccount{}
	size, err = blockchainContext.GetCodeSize(address)
	require.Nil(t, err)
	require.Equal(t, int32(len(expectedCode)), size)
}

func TestBlockchainContext_NewAddress(t *testing.T) {
	t.Parallel()

	mockOutput := &contextmock.OutputContextMock{}

	mockWorld := worldmock.NewMockWorld()
	mockWorld.AcctMap.PutAccounts(testAccounts)

	mockRuntime := &contextmock.RuntimeContextMock{}
	mockRuntime.VMType = []byte{0xF, 0xF}

	host := &contextmock.VMHostMock{
		OutputContext:  mockOutput,
		RuntimeContext: mockRuntime,
	}

	// Test error propagation from GetNonce()
	blockchainContext, _ := NewBlockchainContext(host, mockWorld)
	creatorAddress := []byte("account_new")
	creatorAccount := mockWorld.AcctMap.GetAccount(creatorAddress)
	creatorOutputAccount := mockOutput.NewVMOutputAccountFromMockAccount(creatorAccount)
	mockOutput.OutputAccountMock = creatorOutputAccount
	mockOutput.OutputAccountIsNew = true
	mockWorld.Err = errTestError

	address, err := blockchainContext.NewAddress(creatorAddress)
	require.Equal(t, errTestError, err)
	require.Nil(t, address)

	mockWorld.Err = nil

	// Test if nonce is not deducted if 0, before calling BlockchainHook.NewAddres()
	creatorAddress = []byte("account_new")
	creatorAccount = mockWorld.AcctMap.GetAccount(creatorAddress)
	creatorOutputAccount = mockOutput.NewVMOutputAccountFromMockAccount(creatorAccount)
	mockOutput.OutputAccountMock = creatorOutputAccount
	mockOutput.OutputAccountIsNew = true

	expectedCreatorAddres := creatorAddress
	stubBlockchain := &contextmock.BlockchainHookStub{
		GetUserAccountCalled: mockWorld.GetUserAccount,
		NewAddressCalled: func(creatorAddress []byte, creatorNonce uint64, vmType []byte) ([]byte, error) {
			require.Equal(t, expectedCreatorAddres, creatorAddress)
			require.Equal(t, uint64(0), creatorNonce)
			require.Equal(t, mockRuntime.VMType, vmType)
			return []byte("new_address"), nil
		},
	}
	blockchainContext, _ = NewBlockchainContext(host, stubBlockchain)

	address, err = blockchainContext.NewAddress(creatorAddress)
	require.Nil(t, err)
	require.Equal(t, []byte("new_address"), address)

	// Test if nonce is correctly deducted if greater than 0, before calling BlockchainHook.NewAddres()
	creatorAddress = []byte("account_old_with_money")
	creatorAccount = mockWorld.AcctMap.GetAccount(creatorAddress)
	creatorOutputAccount = mockOutput.NewVMOutputAccountFromMockAccount(creatorAccount)
	mockOutput.OutputAccountMock = creatorOutputAccount
	mockOutput.OutputAccountIsNew = false

	expectedCreatorAddres = creatorAddress
	stubBlockchain = &contextmock.BlockchainHookStub{
		GetUserAccountCalled: mockWorld.GetUserAccount,
		NewAddressCalled: func(creatorAddress []byte, creatorNonce uint64, vmType []byte) ([]byte, error) {
			require.Equal(t, expectedCreatorAddres, creatorAddress)
			require.Equal(t, uint64(55), creatorNonce)
			require.Equal(t, mockRuntime.VMType, vmType)
			return []byte("new_address"), nil
		},
	}
	blockchainContext, _ = NewBlockchainContext(host, stubBlockchain)

	address, err = blockchainContext.NewAddress(creatorAddress)
	require.Nil(t, err)
	require.Equal(t, []byte("new_address"), address)

	// Test if error is propagated from Blockchain.NewAddress
	creatorAddress = []byte("account_with_code")
	creatorAccount = mockWorld.AcctMap.GetAccount(creatorAddress)
	creatorOutputAccount = mockOutput.NewVMOutputAccountFromMockAccount(creatorAccount)
	mockOutput.OutputAccountMock = creatorOutputAccount
	mockOutput.OutputAccountIsNew = false

	expectedCreatorAddres = creatorAddress
	stubBlockchain = &contextmock.BlockchainHookStub{
		GetUserAccountCalled: mockWorld.GetUserAccount,
		NewAddressCalled: func(creatorAddress []byte, creatorNonce uint64, vmType []byte) ([]byte, error) {
			require.Equal(t, expectedCreatorAddres, creatorAddress)
			require.Equal(t, uint64(3), creatorNonce)
			require.Equal(t, mockRuntime.VMType, vmType)
			return nil, errTestError
		},
	}
	blockchainContext, _ = NewBlockchainContext(host, stubBlockchain)

	address, err = blockchainContext.NewAddress(creatorAddress)
	require.Equal(t, errTestError, err)
	require.Nil(t, address)
}

func TestBlockchainContext_BlockHash(t *testing.T) {
	t.Parallel()

	// TODO rewrite this test to use absolute block nonces
	host := &contextmock.VMHostMock{}
	mockWorld := worldmock.NewMockWorld()
	blockchainContext, _ := NewBlockchainContext(host, mockWorld)

	mockWorld.Err = errTestError
	hash := blockchainContext.BlockHash(42)
	require.Nil(t, hash)
	mockWorld.Err = nil

	mockWorld.SetCurrentBlockHash([]byte("1234fa"))
	hash = blockchainContext.BlockHash(3)
	require.Nil(t, hash)

	mockWorld.SetCurrentBlockHash([]byte("1234fb"))
	hash = blockchainContext.BlockHash(0)
	require.Equal(t, []byte("1234fb"), hash)

	mockWorld.SetCurrentBlockHash([]byte("1234fc"))
	hash = blockchainContext.BlockHash(42)
	require.Nil(t, hash)
}

func TestBlockchainContext_IsPayable(t *testing.T) {
	t.Parallel()

	host := &contextmock.VMHostMock{}
	mockWorld := worldmock.NewMockWorld()
	accounts := []*worldmock.Account{
		{Address: []byte("test"), CodeMetadata: []byte{0, vmcommon.MetadataPayable}},
	}
	mockWorld.AcctMap.PutAccounts(accounts)

	bc, _ := NewBlockchainContext(host, mockWorld)

	isPayable, err := bc.IsPayable([]byte("test"))
	require.Nil(t, err)
	require.True(t, isPayable)
}

func TestBlockchainContext_Getters(t *testing.T) {
	t.Parallel()

	host := &contextmock.VMHostMock{}
	var randomSeed1 [48]byte
	copy(randomSeed1[:], "last random seed                                ")
	var randomSeed2 [48]byte
	copy(randomSeed2[:], "current random seed                             ")

	mockWorld := &worldmock.MockWorld{
		PreviousBlockInfo: &worldmock.BlockInfo{
			BlockTimestamp: 6749,
			BlockNonce:     90,
			BlockRound:     96,
			BlockEpoch:     3,
			RandomSeed:     &randomSeed1,
		},
		CurrentBlockInfo: &worldmock.BlockInfo{
			BlockTimestamp: 6800,
			BlockNonce:     98,
			BlockRound:     99,
			BlockEpoch:     4,
			RandomSeed:     &randomSeed2,
		},
		StateRootHash: []byte("root hash"),
	}

	blockchainContext, _ := NewBlockchainContext(host, mockWorld)

	require.Equal(t, uint32(3), blockchainContext.LastEpoch())
	require.Equal(t, uint32(4), blockchainContext.CurrentEpoch())

	require.Equal(t, uint64(90), blockchainContext.LastNonce())
	require.Equal(t, uint64(98), blockchainContext.CurrentNonce())

	require.Equal(t, uint64(96), blockchainContext.LastRound())
	require.Equal(t, uint64(99), blockchainContext.CurrentRound())

	require.Equal(t, uint64(6749), blockchainContext.LastTimeStamp())
	require.Equal(t, uint64(6800), blockchainContext.CurrentTimeStamp())

	require.Equal(t, []byte("root hash"), blockchainContext.GetStateRootHash())
	require.Equal(t, randomSeed1[:], blockchainContext.LastRandomSeed())
	require.Equal(t, randomSeed2[:], blockchainContext.CurrentRandomSeed())
}
