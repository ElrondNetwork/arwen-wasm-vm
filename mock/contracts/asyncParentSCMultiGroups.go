package contracts

import (
	"math/big"

	"github.com/ElrondNetwork/arwen-wasm-vm/v1_3/arwen"
	mock "github.com/ElrondNetwork/arwen-wasm-vm/v1_3/mock/context"
	test "github.com/ElrondNetwork/arwen-wasm-vm/v1_3/testcommon"
	"github.com/ElrondNetwork/elrond-vm-common/txDataBuilder"
	"github.com/stretchr/testify/require"
)

var AsyncGroupsConfig = [][]string{
	{"reserveHousingGroup", "reserveMotel", "reserveHotel"},
	{"reserveTravelGroup", "reserveTrain", "reserveCar", "reserveAirplane"},
}

var AsyncReturnDataSuffix = "_returnData"
var AsyncCallbackPrefix = "callback_"

var AsyncContextCallbackFunction = "contextCallback"

// ForwardAsyncCallMultiGroupsMock is an exposed mock contract method
func ForwardAsyncCallMultiGroupsMock(instanceMock *mock.InstanceMock, testConfig *test.TestConfig) {
	instanceMock.AddMockMethod("forwardMultiGroupAsyncCall", func() *mock.InstanceMock {
		host := instanceMock.Host
		instance := mock.GetMockInstance(host)
		t := instance.T
		arguments := host.Runtime().Arguments()
		destination := arguments[0]
		value := big.NewInt(testConfig.TransferFromParentToChild).Bytes()

		host.Metering().UseGas(testConfig.GasUsedByParent)

		async := host.Async()
		for _, groupConfig := range AsyncGroupsConfig {
			groupName := groupConfig[0]
			for g := 1; g < len(groupConfig); g++ {
				callData := txDataBuilder.NewBuilder()
				functionName := groupConfig[g]
				callData.Func(functionName)
				// child will return this
				callData.Str(functionName + AsyncReturnDataSuffix)

				err := async.RegisterAsyncCall(groupName, &arwen.AsyncCall{
					Status:          arwen.AsyncCallPending,
					Destination:     destination,
					Data:            callData.ToBytes(),
					ValueBytes:      value,
					GasLimit:        uint64(300),
					SuccessCallback: AsyncCallbackPrefix + functionName,
					ErrorCallback:   AsyncCallbackPrefix + functionName,
				})
				require.Nil(t, err)
			}

			async.SetGroupCallback(groupName, AsyncCallbackPrefix+groupName, nil, uint64(100))
		}

		async.SetContextCallback(AsyncContextCallbackFunction, nil, uint64(100))

		return instance

	})
}

// CallBackMultiGroupsMock is an exposed mock contract method
func CallBackMultiGroupsMock(instanceMock *mock.InstanceMock, testConfig *test.TestConfig) {
	for _, groupConfig := range AsyncGroupsConfig {
		groupName := groupConfig[0]
		for g := 1; g < len(groupConfig); g++ {
			functionName := groupConfig[g]
			instanceMock.AddMockMethod(AsyncCallbackPrefix+functionName,
				test.WasteGasWithReturnDataMockMethod(
					instanceMock,
					testConfig.GasUsedByCallback,
					[]byte(AsyncCallbackPrefix+functionName+AsyncReturnDataSuffix)))
		}

		instanceMock.AddMockMethod(AsyncCallbackPrefix+groupName,
			test.WasteGasWithReturnDataMockMethod(
				instanceMock,
				testConfig.GasUsedByCallback,
				[]byte(AsyncCallbackPrefix+groupName+AsyncReturnDataSuffix)))

		instanceMock.AddMockMethod(AsyncContextCallbackFunction,
			test.WasteGasWithReturnDataMockMethod(
				instanceMock,
				testConfig.GasUsedByCallback,
				[]byte(AsyncContextCallbackFunction+AsyncReturnDataSuffix)))
	}
}
