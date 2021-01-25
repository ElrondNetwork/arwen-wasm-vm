package host

import (
	"math/big"
	"sync"
	"testing"

	"github.com/ElrondNetwork/arwen-wasm-vm/arwen"
	worldmock "github.com/ElrondNetwork/arwen-wasm-vm/mock/world"
	"github.com/ElrondNetwork/elrond-go/core/vmcommon"
	"github.com/stretchr/testify/require"
)

func TestExecution_MultipleArwenInstances(t *testing.T) {
	if testing.Short() {
		t.Skip("not a short test")
	}

	numArwenInstances := 40

	hosts := make([]*vmHost, 0)
	worlds := make([]*worldmock.MockWorld, 0)

	// The mutexes are required to ensure that each Arwen instance executes serially.
	mutexes := make([]*sync.Mutex, 0)

	// Create all Arwen instances *before* making any SC deplyoments or calls,
	// because each new instance will overwrite the Wasmer import object.
	for i := 0; i < numArwenInstances; i++ {
		host, mockWorld := createArwenInstanceForDeployment(t)
		hosts = append(hosts, host)
		worlds = append(worlds, mockWorld)
		mutexes = append(mutexes, &sync.Mutex{})
	}

	// Perform all deployments, one for each Arwen instance
	for i := 0; i < numArwenInstances; i++ {
		deployERC20OnArwenInstance(t, hosts[i], worlds[i])
	}

	numRepeats := 3000
	for i := 0; i < numRepeats; i++ {
		// Call a contract on each instance
		for i := 0; i < numArwenInstances; i++ {
			host := hosts[i]
			world := worlds[i]
			input := createERC20TransferInput(owner, receiver)
			mutex := mutexes[i]

			go func() {
				mutex.Lock()
				vmOutput, err := host.RunSmartContractCall(input)
				require.Nil(t, err)
				require.NotNil(t, vmOutput)

				_ = world.UpdateAccounts(vmOutput.OutputAccounts, nil)
				mutex.Unlock()
			}()
		}
	}
}

func createArwenInstanceForDeployment(tb testing.TB) (*vmHost, *worldmock.MockWorld) {
	mockWorld := worldmock.NewMockWorld()
	gasMap, err := loadGasScheduleConfig("../../test/gasSchedule.toml")
	require.Nil(tb, err)

	host, err := NewArwenVM(mockWorld, &arwen.VMHostParameters{
		VMType:                   defaultVMType,
		BlockGasLimit:            uint64(1000),
		GasSchedule:              gasMap,
		ProtocolBuiltinFunctions: make(vmcommon.FunctionNames),
		ElrondProtectedKeyPrefix: []byte("ELROND"),
	})
	require.Nil(tb, err)

	return host, mockWorld
}

func deployERC20OnArwenInstance(tb testing.TB, host *vmHost, world *worldmock.MockWorld) {
	totalTokenSupply := big.NewInt(1000)
	ownerAccount := &worldmock.Account{
		Address: owner,
		Nonce:   1024,
		Balance: big.NewInt(0),
	}
	world.AcctMap.PutAccount(ownerAccount)
	world.NewAddressMocks = append(world.NewAddressMocks, &worldmock.NewAddressMock{
		CreatorAddress: owner,
		CreatorNonce:   ownerAccount.Nonce,
		NewAddress:     scAddress,
	})

	deployInput := &vmcommon.ContractCreateInput{
		VMInput: vmcommon.VMInput{
			CallerAddr: owner,
			Arguments: [][]byte{
				totalTokenSupply.Bytes(),
			},
			CallValue:   big.NewInt(0),
			CallType:    vmcommon.DirectCall,
			GasPrice:    0,
			GasProvided: 0xFFFFFFFFFFFFFFFF,
		},
		ContractCode: GetTestSCCode("erc20", "../../"),
	}

	ownerAccount.Nonce++ // nonce increases before deploy
	vmOutput, err := host.RunSmartContractCreate(deployInput)
	require.Nil(tb, err)
	require.NotNil(tb, vmOutput)
	require.Equal(tb, "", vmOutput.ReturnMessage)
	require.Equal(tb, vmcommon.Ok, vmOutput.ReturnCode)

	// Ensure the deployment persists in the mock BlockchainHook
	_ = world.UpdateAccounts(vmOutput.OutputAccounts, nil)
}

func createERC20TransferInput(caller []byte, receiver []byte) *vmcommon.ContractCallInput {
	transferInput := &vmcommon.ContractCallInput{
		VMInput: vmcommon.VMInput{
			CallerAddr: caller,
			Arguments: [][]byte{
				receiver,
				big.NewInt(1).Bytes(),
			},
			CallValue:   big.NewInt(10),
			CallType:    vmcommon.DirectCall,
			GasPrice:    100000000000000,
			GasProvided: gasProvided,
		},
		RecipientAddr: scAddress,
		Function:      "transferToken",
	}

	return transferInput
}
