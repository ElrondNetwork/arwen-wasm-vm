package host

import (
	"math/big"

	mock "github.com/ElrondNetwork/arwen-wasm-vm/mock/context"
	"github.com/stretchr/testify/require"
)

func wasteGasChildMock(instanceMock *mock.InstanceMock, config interface{}) {
	testConfig := config.(directCallGasTestConfig)
	instanceMock.AddMockMethod("wasteGas", simpleWasteGasMockMethod(instanceMock, testConfig.gasUsedByChild))
}

func execOnSameCtxParentMock(instanceMock *mock.InstanceMock, config interface{}) {
	testConfig := config.(directCallGasTestConfig)
	instanceMock.AddMockMethod("execOnSameCtx", func() *mock.InstanceMock {
		host := instanceMock.Host
		instance := mock.GetMockInstance(host)
		t := instance.T
		host.Metering().UseGas(testConfig.gasUsedByParent)

		arguments := host.Runtime().Arguments()
		input := DefaultTestContractCallInput()
		input.GasProvided = testConfig.gasProvidedToChild
		input.CallerAddr = instance.Address
		input.RecipientAddr = arguments[0]
		input.Function = string(arguments[1])
		numCalls := big.NewInt(0).SetBytes(arguments[2]).Uint64()

		for i := uint64(0); i < numCalls; i++ {
			_, err := host.ExecuteOnSameContext(input)
			require.Nil(t, err)
		}

		return instance
	})
}

func execOnDestCtxParentMock(instanceMock *mock.InstanceMock, config interface{}) {
	testConfig := config.(directCallGasTestConfig)
	instanceMock.AddMockMethod("execOnDestCtx", func() *mock.InstanceMock {
		host := instanceMock.Host
		instance := mock.GetMockInstance(host)
		t := instance.T
		host.Metering().UseGas(testConfig.gasUsedByParent)

		argsPerCall := 3
		arguments := host.Runtime().Arguments()
		if len(arguments)%argsPerCall != 0 {
			host.Runtime().SignalUserError("need 3 arguments per individual call")
			return instance
		}

		input := DefaultTestContractCallInput()
		input.GasProvided = testConfig.gasProvidedToChild
		input.CallerAddr = instance.Address

		for callIndex := 0; callIndex < len(arguments); callIndex += argsPerCall {
			input.RecipientAddr = arguments[callIndex+0]
			input.Function = string(arguments[callIndex+1])
			numCalls := big.NewInt(0).SetBytes(arguments[callIndex+2]).Uint64()

			for i := uint64(0); i < numCalls; i++ {
				_, _, err := host.ExecuteOnDestContext(input)
				require.Nil(t, err)
			}
		}

		return instance
	})
}

func execOnDestCtxSingleCallParentMock(instanceMock *mock.InstanceMock, config interface{}) {
	testConfig := config.(directCallGasTestConfig)
	instanceMock.AddMockMethod("execOnDestCtxSingleCall", func() *mock.InstanceMock {
		host := instanceMock.Host
		instance := mock.GetMockInstance(host)
		host.Metering().UseGas(testConfig.gasUsedByParent)

		arguments := host.Runtime().Arguments()
		if len(arguments) != 2 {
			host.Runtime().SignalUserError("need 2 arguments")
			return instance
		}

		input := DefaultTestContractCallInput()
		input.GasProvided = testConfig.gasProvidedToChild
		input.CallerAddr = instance.Address

		input.RecipientAddr = arguments[0]
		input.Function = string(arguments[1])

		_, _, err := host.ExecuteOnDestContext(input)
		if err != nil {
			host.Runtime().FailExecution(err)
		}

		return instance
	})
}

func wasteGasParentMock(instanceMock *mock.InstanceMock, config interface{}) {
	testConfig := config.(directCallGasTestConfig)
	instanceMock.AddMockMethod("wasteGas", simpleWasteGasMockMethod(instanceMock, testConfig.gasUsedByParent))
}
