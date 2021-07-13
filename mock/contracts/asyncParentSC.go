package contracts

import (
	"bytes"
	"encoding/hex"
	"errors"
	"math/big"

	"github.com/ElrondNetwork/arwen-wasm-vm/v1_3/arwen"
	mock "github.com/ElrondNetwork/arwen-wasm-vm/v1_3/mock/context"
	test "github.com/ElrondNetwork/arwen-wasm-vm/v1_3/testcommon"
	"github.com/ElrondNetwork/elrond-vm-common/txDataBuilder"
	"github.com/stretchr/testify/require"
)

var AsyncChildFunction = "transferToThirdParty"
var AsyncChildData = " there"

// PerformAsyncCallParentMock is an exposed mock contract method
func PerformAsyncCallParentMock(instanceMock *mock.InstanceMock, testConfig *test.TestConfig) {
	instanceMock.AddMockMethod("performAsyncCall", func() *mock.InstanceMock {
		host := instanceMock.Host
		instance := mock.GetMockInstance(host)
		t := instance.T
		host.Metering().UseGas(testConfig.GasUsedByParent)

		host.Storage().SetStorage(test.ParentKeyA, test.ParentDataA)
		host.Storage().SetStorage(test.ParentKeyB, test.ParentDataB)
		host.Output().Finish(test.ParentFinishA)
		host.Output().Finish(test.ParentFinishB)

		scAddress := host.Runtime().GetSCAddress()
		transferValue := big.NewInt(testConfig.TransferToThirdParty)
		err := host.Output().Transfer(test.ThirdPartyAddress, scAddress, 0, 0, transferValue, []byte("hello"), 0)
		require.Nil(t, err)

		// function to be called on child
		// value to send to third party
		// data for child -> third party tx
		// behavior param for child
		// amount to transfer from parent to child
		err = RegisterAsyncCallToChild(host, testConfig, host.Runtime().Arguments())
		require.Nil(t, err)

		return instance

	})
}

// RegisterAsyncCallToChild is resued also in some tests before async context serialization
func RegisterAsyncCallToChild(host arwen.VMHost, testConfig *test.TestConfig, arguments [][]byte) error {
	callData := txDataBuilder.NewBuilder()
	callData.Func(AsyncChildFunction)
	callData.Int64(testConfig.TransferToThirdParty)
	callData.Str(AsyncChildData)
	callData.Bytes(append(arguments[0]))

	value := big.NewInt(testConfig.TransferFromParentToChild).Bytes()
	async := host.Async()
	return async.RegisterLegacyAsyncCall(test.ChildAddress, callData.ToBytes(), value)
}

// SimpleCallbackMock is an exposed mock contract method
func SimpleCallbackMock(instanceMock *mock.InstanceMock, testConfig *test.TestConfig) {
	instanceMock.AddMockMethod("callBack", func() *mock.InstanceMock {
		host := instanceMock.Host
		instance := mock.GetMockInstance(host)
		arguments := host.Runtime().Arguments()

		host.Metering().UseGas(testConfig.GasUsedByCallback)

		if string(arguments[1]) == "fail" {
			host.Runtime().SignalUserError("wrong num of arguments")
			return instance
		}

		return instance
	})
}

// CallBackParentMock is an exposed mock contract method
func CallBackParentMock(instanceMock *mock.InstanceMock, testConfig *test.TestConfig) {
	instanceMock.AddMockMethod("callBack", func() *mock.InstanceMock {
		host := instanceMock.Host
		instance := mock.GetMockInstance(host)
		t := instance.T
		arguments := host.Runtime().Arguments()

		host.Metering().UseGas(testConfig.GasUsedByCallback)

		if len(arguments) < 2 {
			host.Runtime().SignalUserError("wrong num of arguments")
			return instance
		}

		loadedData := host.Storage().GetStorage(test.ParentKeyB)

		status := bytes.Compare(loadedData, test.ParentDataB)
		if status != 0 {
			status = 1
		}

		if len(arguments) >= 4 {
			err := handleParentBehaviorArgument(host, big.NewInt(0).SetBytes(arguments[1]))
			if err != nil {
				return instance
			}
		}
		err := handleTransferToVault(host, arguments)
		require.Nil(t, err)

		finishResult(host, status)

		return instance
	})
}

func handleParentBehaviorArgument(host arwen.VMHost, behavior *big.Int) error {
	if behavior.Cmp(big.NewInt(3)) == 0 {
		host.Runtime().SignalUserError("callBack error")
		return errors.New("behavior / parent error")
	}
	if behavior.Cmp(big.NewInt(4)) == 0 {
		for {
			host.Output().Finish([]byte("loop"))
		}
	}

	behaviorBytes := behavior.Bytes()
	if len(behaviorBytes) == 0 {
		behaviorBytes = []byte{0}
	}
	host.Output().Finish(behaviorBytes)

	return nil
}

func mustTransferToVault(arguments [][]byte) bool {
	vault := "vault"
	numArgs := len(arguments)
	if numArgs == 3 {
		if string(arguments[2]) == vault {
			return false
		}
	}

	if numArgs == 4 {
		if string(arguments[3]) == vault {
			return false
		}
	}

	return true
}

func handleTransferToVault(host arwen.VMHost, arguments [][]byte) error {
	err := error(nil)
	if mustTransferToVault(arguments) {
		valueToTransfer := big.NewInt(4)
		err = host.Output().Transfer(test.VaultAddress, host.Runtime().GetSCAddress(), 0, 0, valueToTransfer, arguments[1], 0)
	}

	return err
}

func finishResult(host arwen.VMHost, result int) {
	outputContext := host.Output()
	if result == 0 {
		outputContext.Finish([]byte("succ"))
	}
	if result == 1 {
		outputContext.Finish([]byte("fail"))
	}
	if result != 0 && result != 1 {
		outputContext.Finish([]byte("unkn"))
	}
}

func argumentsToHexString(functionName string, args ...[]byte) []byte {
	separator := byte('@')
	output := append([]byte(functionName))
	for _, arg := range args {
		output = append(output, separator)
		output = append(output, hex.EncodeToString(arg)...)
	}
	return output
}
