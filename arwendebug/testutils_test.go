package arwendebug

import (
	"fmt"
	"math/big"
	"math/rand"
	"testing"
	"time"

	vmcommon "github.com/ElrondNetwork/elrond-vm-common"
	"github.com/stretchr/testify/require"
)

const gasLimit = 50000000

type testContext struct {
	t       *testing.T
	worldID string
	facade  *DebugFacade
}

func newTestContext(t *testing.T) *testContext {
	worldID := fmt.Sprintf("%s_%d", time.Now().Format("20060102150405"), rand.Intn(100))

	return &testContext{
		t:       t,
		worldID: worldID,
		facade:  NewDebugFacade(),
	}
}

func (context *testContext) createAccount(address string, balance string) {
	request := CreateAccountRequest{
		RequestBase: context.createRequestBase(),
		AddressHex:  address,
		Balance:     balance,
		Nonce:       0,
	}

	response, err := context.facade.CreateAccount(request)

	t := context.t
	require.Nil(t, err)
	require.NotNil(t, response)
}

func (context *testContext) accountExists(address []byte) bool {
	world := context.loadWorld()
	account, err := world.blockchainHook.GetUserAccount(address)
	return err == nil && account != nil
}

func (context *testContext) deployContract(codePath string, impersonated string, arguments ...string) *DeployResponse {
	request := DeployRequest{
		ContractRequestBase: ContractRequestBase{
			RequestBase:     context.createRequestBase(),
			ImpersonatedHex: impersonated,
			GasLimit:        gasLimit,
		},
		CodePath:     codePath,
		ArgumentsHex: arguments,
	}

	response, err := context.facade.DeploySmartContract(request)

	t := context.t
	require.Nil(t, err)
	require.NotNil(t, response)
	require.NotNil(t, response.Output)
	require.Nil(t, response.Error)
	require.Equal(t, vmcommon.Ok.String(), response.Output.ReturnCode.String(), response.Output.ReturnMessage)

	return response
}

func (context *testContext) upgradeContract(contract string, codePath string, impersonated string, arguments ...string) *UpgradeResponse {
	request := UpgradeRequest{
		DeployRequest: DeployRequest{
			ContractRequestBase: ContractRequestBase{
				RequestBase:     context.createRequestBase(),
				ImpersonatedHex: impersonated,
				GasLimit:        gasLimit,
			},
			CodePath:     codePath,
			ArgumentsHex: arguments,
		},
		ContractAddressHex: contract,
	}

	response, err := context.facade.UpgradeSmartContract(request)

	t := context.t
	require.Nil(t, err)
	require.NotNil(t, response)
	require.NotNil(t, response.Output)
	require.Nil(t, response.Error)
	require.Equal(t, vmcommon.Ok.String(), response.Output.ReturnCode.String(), response.Output.ReturnMessage)

	return response
}

func (context *testContext) runContract(contract string, impersonated string, function string, arguments ...string) *RunResponse {
	request := RunRequest{
		ContractRequestBase: ContractRequestBase{
			RequestBase:     context.createRequestBase(),
			ImpersonatedHex: impersonated,
			GasLimit:        gasLimit,
		},
		ContractAddressHex: contract,
		Function:           function,
		ArgumentsHex:       arguments,
	}

	response, err := context.facade.CallSmartContract(request)

	t := context.t
	require.Nil(t, err)
	require.NotNil(t, response)
	require.NotNil(t, response.Output)
	require.Nil(t, response.Error)
	require.Equal(t, vmcommon.Ok.String(), response.Output.ReturnCode.String(), response.Output.ReturnMessage)

	return response
}

func (context *testContext) queryContract(contract string, impersonated string, function string, arguments ...string) *QueryResponse {
	request := QueryRequest{
		RunRequest: RunRequest{
			ContractRequestBase: ContractRequestBase{
				RequestBase:     context.createRequestBase(),
				ImpersonatedHex: impersonated,
				GasLimit:        gasLimit,
			},
			ContractAddressHex: contract,
			Function:           function,
			ArgumentsHex:       arguments,
		},
	}

	response, err := context.facade.QuerySmartContract(request)

	t := context.t
	require.Nil(t, err)
	require.NotNil(t, response)
	require.NotNil(t, response.Output)
	require.Nil(t, response.Error)
	require.Equal(t, vmcommon.Ok.String(), response.ReturnCodeString, response.Output.ReturnMessage)

	return response
}

func (response *ContractResponseBase) getFirstResultAsInt64() int64 {
	result, err := response.Output.ConvertToVMOutput().GetFirstReturnData(vmcommon.AsBigInt)
	if err != nil {
		return 0
	}

	asBigInt := result.(*big.Int)
	return asBigInt.Int64()
}

func (context *testContext) createRequestBase() RequestBase {
	randomOutcome := fmt.Sprintf("%s_%d", time.Now().Format("20060102150405"), rand.Intn(100))

	return RequestBase{
		DatabasePath: databasePath,
		World:        context.worldID,
		Outcome:      randomOutcome,
	}
}

func (context *testContext) loadWorld() *world {
	database := newDatabase(databasePath)
	world, err := database.loadWorld(context.worldID)
	require.Nil(context.t, err)

	return world
}

type dummyAddress struct {
	hex string
	raw []byte
}

// Left-pads the input address with character "0" (not \0)
func newDummyAddress(address string) *dummyAddress {
	rawString := fmt.Sprintf("%032s", address)
	return &dummyAddress{
		hex: toHex([]byte(rawString)),
		raw: []byte(rawString),
	}
}
