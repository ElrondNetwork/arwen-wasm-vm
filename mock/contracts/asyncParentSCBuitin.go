package contracts

import (
	"math/big"

	mock "github.com/ElrondNetwork/arwen-wasm-vm/v1_3/mock/context"
	test "github.com/ElrondNetwork/arwen-wasm-vm/v1_3/testcommon"
	"github.com/stretchr/testify/require"
)

// ForwardAsyncCallParentBuiltinMock is an exposed mock contract method
func ForwardAsyncCallParentBuiltinMock(instanceMock *mock.InstanceMock, testConfig *test.TestConfig) {
	instanceMock.AddMockMethod("forwardAsyncCall", func() *mock.InstanceMock {
		host := instanceMock.Host
		instance := mock.GetMockInstance(host)
		host.Metering().UseGas(testConfig.GasUsedByParent)

		arguments := host.Runtime().Arguments()
		destination := arguments[0]
		function := arguments[1]
		value := big.NewInt(testConfig.TransferFromParentToChild).Bytes()

		err := host.Async().RegisterLegacyAsyncCall(destination, function, value)
		require.Nil(instance.T, err)

		return instance
	})
}

// CallBackParentBuiltinMock is an exposed mock contract method
func CallBackParentBuiltinMock(instanceMock *mock.InstanceMock, testConfig *test.TestConfig) {
	instanceMock.AddMockMethod("callBack", test.SimpleWasteGasMockMethod(instanceMock, testConfig.GasUsedByCallback))
}
