package contracts

import (
	"github.com/ElrondNetwork/arwen-wasm-vm/v1_3/arwen"
	"github.com/ElrondNetwork/arwen-wasm-vm/v1_3/arwen/elrondapi"
	vmcommon "github.com/ElrondNetwork/elrond-vm-common"
)

// ExecuteOnSameContextInMockContracts - calls the corresponding method in elrond api
func ExecuteOnSameContextInMockContracts(host arwen.VMHost, input *vmcommon.ContractCallInput) int32 {
	return elrondapi.ExecuteOnSameContextWithTypedArgs(host, int64(input.GasProvided), input.CallValue, []byte(input.Function), input.RecipientAddr, input.Arguments)
}

// ExecuteOnDestContextInMockContracts - calls the corresponding method in elrond api
func ExecuteOnDestContextInMockContracts(host arwen.VMHost, input *vmcommon.ContractCallInput) int32 {
	return elrondapi.ExecuteOnDestContextWithTypedArgs(host, int64(input.GasProvided), input.CallValue, []byte(input.Function), input.RecipientAddr, input.Arguments)
}
