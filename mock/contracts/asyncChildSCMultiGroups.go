package contracts

import (
	mock "github.com/ElrondNetwork/arwen-wasm-vm/v1_3/mock/context"
	test "github.com/ElrondNetwork/arwen-wasm-vm/v1_3/testcommon"
)

// ChildAsyncMultiGroupsMock is an exposed mock contract method
func ChildAsyncMultiGroupsMock(instanceMock *mock.InstanceMock, testConfig *test.TestConfig) {
	for _, groupConfig := range AsyncGroupsConfig {
		for g := 1; g < len(groupConfig); g++ {
			functionName := groupConfig[g]
			instanceMock.AddMockMethod(functionName,
				test.WasteGasWithReturnDataMockMethod(
					instanceMock,
					testConfig.GasUsedByCallback,
					[]byte(functionName+AsyncReturnDataSuffix)))
		}
	}
}
