package arwendebug

import (
	"math/big"

	"github.com/ElrondNetwork/arwen-wasm-vm/ipc/common"
	vmcommon "github.com/ElrondNetwork/elrond-vm-common"
)

// RequestBase is a CLI / REST request message
type RequestBase struct {
	DatabasePath string
	World        string
	Outcome      string
}

func (request *RequestBase) digest() error {
	if request.DatabasePath == "" {
		request.DatabasePath = "./db"
	}

	if request.World == "" {
		request.World = "default"
	}

	return nil
}

// ResponseBase is a CLI / REST response message
type ResponseBase struct {
	Error error
}

// ContractRequestBase is a CLI / REST request message
type ContractRequestBase struct {
	RequestBase
	ImpersonatedHex string
	Impersonated    []byte
	Value           string
	ValueAsBigInt   *big.Int
	GasPrice        uint64
	GasLimit        uint64
}

func (request *ContractRequestBase) digest() error {
	err := request.RequestBase.digest()
	if err != nil {
		return err
	}

	if request.ImpersonatedHex == "" {
		return NewRequestError("empty impersonated address")
	}

	request.Impersonated, err = fromHex(request.ImpersonatedHex)
	if err != nil {
		return NewRequestErrorMessageInner("invalid impersonated address", err)
	}

	if request.GasPrice == 0 {
		request.GasPrice = DefaultGasPrice
	}

	if request.GasLimit == 0 {
		return NewRequestError("invalid gas limit")
	}

	request.ValueAsBigInt, err = parseValue(request.Value)
	if err != nil {
		return err
	}

	return nil
}

// ContractResponseBase is a CLI / REST response message
type ContractResponseBase struct {
	ResponseBase
	Input            *vmcommon.VMInput
	Output           *common.SerializableVMOutput
	ReturnCodeString string
}

func createContractResponseBase(input *vmcommon.VMInput, output *vmcommon.VMOutput) ContractResponseBase {
	response := ContractResponseBase{
		Input:  input,
		Output: common.NewSerializableVMOutput(output),
	}

	if output != nil {
		response.ReturnCodeString = output.ReturnCode.String()
	}

	return response
}
