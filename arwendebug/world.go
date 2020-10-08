package arwendebug

import (
	"github.com/ElrondNetwork/arwen-wasm-vm/arwen"
	"github.com/ElrondNetwork/arwen-wasm-vm/arwen/host"
	"github.com/ElrondNetwork/arwen-wasm-vm/config"
	vmcommon "github.com/ElrondNetwork/elrond-vm-common"
)

type worldDataModel struct {
	ID       string
	Accounts AccountsMap
}

type world struct {
	id             string
	blockchainHook *BlockchainHookMock
	vm             vmcommon.VMExecutionHandler
}

func newWorldDataModel(worldID string) *worldDataModel {
	return &worldDataModel{
		ID:       worldID,
		Accounts: make(AccountsMap),
	}
}

// newWorld creates a new debugging world
func newWorld(dataModel *worldDataModel) (*world, error) {
	blockchainHook := NewBlockchainHookMock()
	blockchainHook.Accounts = dataModel.Accounts

	vm, err := host.NewArwenVM(
		blockchainHook,
		getHostParameters(),
	)
	if err != nil {
		return nil, err
	}

	return &world{
		id:             dataModel.ID,
		blockchainHook: blockchainHook,
		vm:             vm,
	}, nil
}

func getHostParameters() *arwen.VMHostParameters {
	return &arwen.VMHostParameters{
		VMType:                   []byte{5, 0},
		BlockGasLimit:            uint64(10000000),
		GasSchedule:              config.MakeGasMap(1, 1),
		ElrondProtectedKeyPrefix: []byte("ELROND"),
	}
}

func (w *world) deploySmartContract(request DeployRequest) *DeployResponse {
	input := w.prepareDeployInput(request)
	log.Trace("w.deploySmartContract()", "input", prettyJson(input))

	vmOutput, err := w.vm.RunSmartContractCreate(input)
	if err == nil {
		w.blockchainHook.UpdateAccounts(vmOutput.OutputAccounts)
	}

	response := &DeployResponse{}
	response.ContractResponseBase = createContractResponseBase(&input.VMInput, vmOutput)
	response.Error = err
	response.ContractAddress = w.blockchainHook.LastCreatedContractAddress
	response.ContractAddressHex = toHex(response.ContractAddress)
	return response
}

func (w *world) upgradeSmartContract(request UpgradeRequest) *UpgradeResponse {
	input := w.prepareUpgradeInput(request)
	log.Trace("w.upgradeSmartContract()", "input", prettyJson(input))

	vmOutput, err := w.vm.RunSmartContractCall(input)
	if err == nil {
		w.blockchainHook.UpdateAccounts(vmOutput.OutputAccounts)
	}

	response := &UpgradeResponse{}
	response.ContractResponseBase = createContractResponseBase(&input.VMInput, vmOutput)
	response.Error = err

	return response
}

func (w *world) callSmartContract(request RunRequest) *RunResponse {
	input := w.prepareCallInput(request)
	log.Trace("w.callSmartContract()", "input", prettyJson(input))

	vmOutput, err := w.vm.RunSmartContractCall(input)
	if err == nil {
		w.blockchainHook.UpdateAccounts(vmOutput.OutputAccounts)
	}

	response := &RunResponse{}
	response.ContractResponseBase = createContractResponseBase(&input.VMInput, vmOutput)
	response.Error = err

	return response
}

func (w *world) querySmartContract(request QueryRequest) *QueryResponse {
	input := w.prepareCallInput(request.RunRequest)
	log.Trace("w.querySmartContract()", "input", prettyJson(input))

	vmOutput, err := w.vm.RunSmartContractCall(input)

	response := &QueryResponse{}
	response.ContractResponseBase = createContractResponseBase(&input.VMInput, vmOutput)
	response.Error = err

	return response
}

func (w *world) createAccount(request CreateAccountRequest) *CreateAccountResponse {
	log.Trace("w.createAccount()", "request", prettyJson(request))

	account := NewAccount(request.Address, request.Nonce, request.BalanceAsBigInt)
	w.blockchainHook.AddAccount(account)
	return &CreateAccountResponse{Account: account}
}

func (w *world) toDataModel() *worldDataModel {
	return &worldDataModel{
		ID:       w.id,
		Accounts: w.blockchainHook.Accounts,
	}
}
