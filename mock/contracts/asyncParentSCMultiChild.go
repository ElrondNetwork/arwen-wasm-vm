package contracts

import (
	"math/big"

	"github.com/ElrondNetwork/arwen-wasm-vm/v1_3/arwen"
	mock "github.com/ElrondNetwork/arwen-wasm-vm/v1_3/mock/context"
	test "github.com/ElrondNetwork/arwen-wasm-vm/v1_3/testcommon"
	"github.com/ElrondNetwork/elrond-vm-common/txDataBuilder"
	"github.com/stretchr/testify/require"
)

// ForwardAsyncCallMultiChildMock is an exposed mock contract method
func ForwardAsyncCallMultiChildMock(instanceMock *mock.InstanceMock, testConfig *test.TestConfig) {
	instanceMock.AddMockMethod("forwardAsyncCall", func() *mock.InstanceMock {
		host := instanceMock.Host
		instance := mock.GetMockInstance(host)
		t := instance.T
		arguments := host.Runtime().Arguments()
		destination := arguments[0]
		function := string(arguments[1])
		value := big.NewInt(testConfig.TransferFromParentToChild).Bytes()

		host.Metering().UseGas(testConfig.GasUsedByParent)

		for childCall := 0; childCall < testConfig.ChildCalls; childCall++ {
			callData := txDataBuilder.NewBuilder()
			callData.Func(function)
			// recursiveChildCalls
			callData.BigInt(big.NewInt(1))
			// child will return this
			callData.BigInt(big.NewInt(int64(childCall)))

			async := host.Async()
			err := async.RegisterAsyncCall("myAsyncGroup", &arwen.AsyncCall{
				Status:          arwen.AsyncCallPending,
				Destination:     destination,
				Data:            callData.ToBytes(),
				ValueBytes:      value,
				GasLimit:        uint64(300),
				SuccessCallback: "callBack",
				ErrorCallback:   "callBack",
			})
			require.Nil(t, err)
		}

		return instance

	})
}

// CallBackMultiChildMock is an exposed mock contract method
func CallBackMultiChildMock(instanceMock *mock.InstanceMock, testConfig *test.TestConfig) {
	instanceMock.AddMockMethod("callBack", test.SimpleWasteGasMockMethod(instanceMock, testConfig.GasUsedByCallback))
}
