package arwendebug

import (
	"github.com/ElrondNetwork/arwen-wasm-vm/arwen"
	"github.com/ElrondNetwork/arwen-wasm-vm/arwen/host"
	"github.com/ElrondNetwork/arwen-wasm-vm/config"
	"github.com/ElrondNetwork/arwen-wasm-vm/ipc/arwenpart"
	vmcommon "github.com/ElrondNetwork/elrond-vm-common"
)

type worldDataModel struct {
	ID        string
	CreatedOn string
	Accounts  AccountsMap
}

type world struct {
	id             string
	blockchainHook BlockchainHookMock
	vm             vmcommon.VMExecutionHandler
}

func newWorldDataModel(worldID string) *worldDataModel {
	return &worldDataModel{
		ID:        worldID,
		CreatedOn: "now",
		Accounts:  make(AccountsMap),
	}
}

// NewWorld -
func NewWorld(dataModel *worldDataModel) (*world, error) {
	blockchainHook := NewBlockchainHookMock()
	blockchainHook.Accounts = dataModel.Accounts

	vm, err := host.NewArwenVM(
		blockchainHook,
		arwenpart.NewCryptoHookGateway(),
		getHostParameters(),
	)
	if err != nil {
		return nil, err
	}

	return &world{
		id:             dataModel.ID,
		blockchainHook: *blockchainHook,
		vm:             vm,
	}, nil
}

func getHostParameters() *arwen.VMHostParameters {
	return &arwen.VMHostParameters{
		VMType:        []byte{5, 0},
		BlockGasLimit: uint64(10000000),
		GasSchedule:   config.MakeGasMap(1),
	}
}

// DeploySmartContract -
func (world *world) DeploySmartContract(request DeployRequest) (*DeployResponse, error) {
	log.Debug("world.DeploySmartContract()")

	createInput, err := world.prepareDeployInput(request)
	if err != nil {
		return nil, err
	}

	vmOutput, err := world.vm.RunSmartContractCreate(createInput)
	if err == nil {
		world.blockchainHook.UpdateAccounts(vmOutput.OutputAccounts)
	}

	response := &DeployResponse{}
	response.Output = vmOutput
	response.Error = err

	return response, nil
}

func (world *world) prepareDeployInput(request DeployRequest) (*vmcommon.ContractCreateInput, error) {
	var err error

	createInput := &vmcommon.ContractCreateInput{}
	createInput.CallerAddr = request.getImpersonated()
	createInput.CallValue = request.getValue()
	createInput.ContractCode, err = request.getCode()
	if err != nil {
		return nil, err
	}

	createInput.ContractCodeMetadata, err = request.getCodeMetadata()
	if err != nil {
		return nil, err
	}

	createInput.GasProvided = request.getGasLimit()
	createInput.GasPrice = request.getGasPrice()

	return createInput, nil
}

// UpgradeSmartContract -
func (world *world) UpgradeSmartContract() (*UpgradeResponse, error) {
	return &UpgradeResponse{}, nil
}

// RunSmartContract -
func (world *world) RunSmartContract(request RunRequest) (*RunResponse, error) {
	callInput, err := world.prepareCallInput(request)
	if err != nil {
		return nil, err
	}

	vmOutput, err := world.vm.RunSmartContractCall(callInput)
	if err == nil {
		world.blockchainHook.UpdateAccounts(vmOutput.OutputAccounts)
	}

	response := &RunResponse{}
	response.Output = vmOutput
	response.Error = err

	return response, nil
}

func (world *world) prepareCallInput(request RunRequest) (*vmcommon.ContractCallInput, error) {
	callInput := &vmcommon.ContractCallInput{}
	callInput.RecipientAddr = []byte(request.ContractAddress)
	callInput.CallerAddr = request.getImpersonated()
	callInput.CallValue = request.getValue()
	callInput.Function = request.Function
	callInput.Arguments = nil
	callInput.GasProvided = request.getGasLimit()
	callInput.GasPrice = request.getGasPrice()

	return callInput, nil
}

// QuerySmartContract -
func (world *world) QuerySmartContract() error {
	return nil
}

func (world *world) CreateAccount(request CreateAccountRequest) (*CreateAccountResponse, error) {
	balance, err := request.getBalance()
	if err != nil {
		return nil, err
	}

	account := &Account{
		Address: request.getAddress(),
		Nonce:   request.Nonce,
		Balance: balance,
	}

	world.blockchainHook.AddAccount(account)
	return &CreateAccountResponse{}, nil
}

func (world *world) toDataModel() *worldDataModel {
	return &worldDataModel{
		ID:        world.id,
		CreatedOn: "test",
		Accounts:  world.blockchainHook.Accounts,
	}
}
