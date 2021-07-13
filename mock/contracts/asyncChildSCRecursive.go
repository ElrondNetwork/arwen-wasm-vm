package contracts

import (
	"math/big"

	mock "github.com/ElrondNetwork/arwen-wasm-vm/v1_3/mock/context"
	test "github.com/ElrondNetwork/arwen-wasm-vm/v1_3/testcommon"
	"github.com/ElrondNetwork/elrond-vm-common/txDataBuilder"
	"github.com/stretchr/testify/require"
)

// RecursiveAsyncCallRecursiveChildMock is an exposed mock contract method
func RecursiveAsyncCallRecursiveChildMock(instanceMock *mock.InstanceMock, testConfig *test.TestConfig) {
	instanceMock.AddMockMethod("recursiveAsyncCall", func() *mock.InstanceMock {
		host := instanceMock.Host
		instance := mock.GetMockInstance(host)
		t := instance.T
		arguments := host.Runtime().Arguments()

		host.Metering().UseGas(testConfig.GasUsedByChild)

		var recursiveChildCalls uint64
		if len(arguments) > 0 {
			recursiveChildCalls = big.NewInt(0).SetBytes(arguments[0]).Uint64()
		} else {
			recursiveChildCalls = 1
		}
		recursiveChildCalls = recursiveChildCalls - 1
		returnValue := big.NewInt(int64(recursiveChildCalls)).Bytes()
		if len(arguments) == 2 {
			returnValue = arguments[1]
		}
		host.Output().Finish(returnValue)
		if recursiveChildCalls == 0 {
			return instance
		}

		destination := host.Runtime().GetSCAddress()
		function := string("recursiveAsyncCall")
		value := big.NewInt(testConfig.TransferFromParentToChild).Bytes()

		callData := txDataBuilder.NewBuilder()
		callData.Func(function)
		callData.BigInt(big.NewInt(int64(recursiveChildCalls)))

		async := host.Async()
		err := async.RegisterLegacyAsyncCall(destination, callData.ToBytes(), value)
		require.Nil(t, err)

		return instance
	})
}

// CallBackRecursiveChildMock is an exposed mock contract method
func CallBackRecursiveChildMock(instanceMock *mock.InstanceMock, testConfig *test.TestConfig) {
	instanceMock.AddMockMethod("callBack", test.SimpleWasteGasMockMethod(instanceMock, testConfig.GasUsedByCallback))
}
