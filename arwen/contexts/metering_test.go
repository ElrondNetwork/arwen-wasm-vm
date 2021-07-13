package contexts

import (
	"math/big"
	"testing"

	"github.com/ElrondNetwork/arwen-wasm-vm/v1_3/arwen"
	"github.com/ElrondNetwork/arwen-wasm-vm/v1_3/config"
	"github.com/ElrondNetwork/arwen-wasm-vm/v1_3/math"
	contextmock "github.com/ElrondNetwork/arwen-wasm-vm/v1_3/mock/context"
	vmcommon "github.com/ElrondNetwork/elrond-vm-common"
	"github.com/stretchr/testify/require"
)

func TestNewMeteringContext(t *testing.T) {
	t.Parallel()
	const BlockGasLimit = uint64(15000)
	host := &contextmock.VMHostMock{}

	meteringContext, err := NewMeteringContext(host, config.MakeGasMapForTests(), BlockGasLimit)
	require.Nil(t, err)
	require.NotNil(t, meteringContext)
}

func TestNewMeteringContext_NilGasSchedule(t *testing.T) {
	t.Parallel()
	const BlockGasLimit = uint64(15000)
	host := &contextmock.VMHostMock{}

	meteringContext, err := NewMeteringContext(host, nil, BlockGasLimit)
	require.NotNil(t, err)
	require.Nil(t, meteringContext)
}

func TestMeteringContext_GasSchedule(t *testing.T) {
	t.Parallel()
	const BlockGasLimit = uint64(15000)

	host := &contextmock.VMHostStub{}
	meteringContext, _ := NewMeteringContext(host, config.MakeGasMapForTests(), BlockGasLimit)

	schedule := meteringContext.GasSchedule()
	require.NotNil(t, schedule)
}

func TestMeteringContext_UseGas(t *testing.T) {
	t.Parallel()
	const BlockGasLimit = uint64(15000)

	mockRuntime := &contextmock.RuntimeContextMock{}
	host := &contextmock.VMHostMock{
		RuntimeContext: mockRuntime,
	}
	meteringContext, _ := NewMeteringContext(host, config.MakeGasMapForTests(), BlockGasLimit)

	gasProvided := uint64(1001)
	meteringContext.gasForExecution = gasProvided
	gasUsed := uint64(1000)
	meteringContext.UseGas(gasUsed)
	require.Equal(t, mockRuntime.GetPointsUsed(), gasUsed)
	require.Equal(t, gasProvided-gasUsed, meteringContext.GasLeft())

	gasProvided = uint64(10000)
	mockRuntime.SetPointsUsed(0)
	meteringContext, _ = NewMeteringContext(host, config.MakeGasMapForTests(), BlockGasLimit)
	meteringContext.gasForExecution = gasProvided

	require.Equal(t, gasProvided, meteringContext.GasLeft())
	meteringContext.UseGas(gasUsed)
	require.Equal(t, gasProvided-gasUsed, meteringContext.GasLeft())
}

func TestMeteringContext_FreeGas(t *testing.T) {
	t.Parallel()
	const BlockGasLimit = uint64(15000)

	mockOutput := &contextmock.OutputContextMock{}
	host := &contextmock.VMHostMock{
		OutputContext: mockOutput,
	}

	meteringContext, _ := NewMeteringContext(host, config.MakeGasMapForTests(), BlockGasLimit)

	gasToFree := uint64(1000)
	mockOutput.GasRefund = big.NewInt(0)
	meteringContext.FreeGas(gasToFree)
	gasRefunded := mockOutput.GetRefund()
	require.Equal(t, gasToFree, gasRefunded)

	moreGasToFree := uint64(100)
	meteringContext.FreeGas(moreGasToFree)
	gasRefunded = mockOutput.GetRefund()
	require.Equal(t, gasToFree+moreGasToFree, gasRefunded)
}

func TestMeteringContext_BoundGasLimit(t *testing.T) {
	t.Parallel()
	const BlockGasLimit = uint64(15000)

	mockRuntime := &contextmock.RuntimeContextMock{}
	host := &contextmock.VMHostMock{
		RuntimeContext: mockRuntime,
	}
	meteringContext, _ := NewMeteringContext(host, config.MakeGasMapForTests(), BlockGasLimit)

	gasProvided := uint64(10000)
	meteringContext.gasForExecution = gasProvided
	mockRuntime.SetPointsUsed(0)

	gasLimit := 5000
	limit := meteringContext.BoundGasLimit(int64(gasLimit))
	require.Equal(t, uint64(gasLimit), limit)

	gasLimit = 25000
	limit = meteringContext.BoundGasLimit(int64(gasLimit))
	require.Equal(t, meteringContext.GasLeft(), limit)

	blockLimit := meteringContext.BlockGasLimit()
	require.Equal(t, BlockGasLimit, blockLimit)
}

func TestMeteringContext_DeductInitialGasForExecution(t *testing.T) {
	t.Parallel()

	mockRuntime := &contextmock.RuntimeContextMock{}
	gasProvided := uint64(10000)
	vmInput := &vmcommon.VMInput{
		GasProvided: gasProvided,
	}

	mockRuntime.SetVMInput(vmInput)

	host := &contextmock.VMHostMock{
		RuntimeContext: mockRuntime,
	}

	meteringContext, _ := NewMeteringContext(host, config.MakeGasMapForTests(), uint64(15000))

	contract := []byte("contract")
	err := meteringContext.DeductInitialGasForExecution(contract)
	require.Nil(t, err)

	vmInput.GasProvided = 1
	err = meteringContext.DeductInitialGasForExecution(contract)
	require.Equal(t, arwen.ErrNotEnoughGas, err)
}

func TestDeductInitialGasForDirectDeployment(t *testing.T) {
	t.Parallel()
	mockRuntime := &contextmock.RuntimeContextMock{}
	gasProvided := uint64(10000)
	contractCode := []byte("contractCode")
	input := &vmcommon.ContractCreateInput{
		VMInput: vmcommon.VMInput{
			GasProvided: gasProvided,
		},
		ContractCode: contractCode,
	}

	mockRuntime.SetVMInput(&input.VMInput)

	host := &contextmock.VMHostMock{
		RuntimeContext: mockRuntime,
	}

	meteringContext, _ := NewMeteringContext(host, config.MakeGasMapForTests(), uint64(15000))

	mockRuntime.SetPointsUsed(0)
	err := meteringContext.DeductInitialGasForDirectDeployment(arwen.CodeDeployInput{ContractCode: contractCode})
	require.Nil(t, err)
	remainingGas := meteringContext.GasLeft()
	require.Equal(t, gasProvided-uint64(len(contractCode))-1, remainingGas)

	input.GasProvided = 2
	mockRuntime.SetPointsUsed(0)
	err = meteringContext.DeductInitialGasForDirectDeployment(arwen.CodeDeployInput{ContractCode: contractCode})
	require.Equal(t, arwen.ErrNotEnoughGas, err)
}

func TestDeductInitialGasForIndirectDeployment(t *testing.T) {
	t.Parallel()

	mockRuntime := &contextmock.RuntimeContextMock{}
	gasProvided := uint64(10000)
	contractCode := []byte("contractCode")
	input := &vmcommon.ContractCreateInput{
		VMInput: vmcommon.VMInput{
			GasProvided: gasProvided,
		},
		ContractCode: contractCode,
	}

	mockRuntime.SetVMInput(&input.VMInput)

	host := &contextmock.VMHostMock{
		RuntimeContext: mockRuntime,
	}

	meteringContext, _ := NewMeteringContext(host, config.MakeGasMapForTests(), uint64(15000))

	mockRuntime.SetPointsUsed(0)
	err := meteringContext.DeductInitialGasForIndirectDeployment(arwen.CodeDeployInput{ContractCode: contractCode})
	require.Nil(t, err)
	remainingGas := meteringContext.GasLeft()
	require.Equal(t, gasProvided-uint64(len(contractCode)), remainingGas)

	input.GasProvided = 2
	mockRuntime.SetPointsUsed(0)
	err = meteringContext.DeductInitialGasForDirectDeployment(arwen.CodeDeployInput{ContractCode: contractCode})
	require.Equal(t, arwen.ErrNotEnoughGas, err)
}

func TestMeteringContext_AsyncCallGasLocking(t *testing.T) {
	t.Parallel()

	mockRuntime := &contextmock.RuntimeContextMock{}
	contractSize := uint64(1000)
	input := &vmcommon.ContractCallInput{
		VMInput: vmcommon.VMInput{
			CallType: vmcommon.AsynchronousCall,
		},
	}

	mockRuntime.SCCodeSize = contractSize
	mockRuntime.SetVMInput(&input.VMInput)
	mockRuntime.SetPointsUsed(0)

	host := &contextmock.VMHostMock{
		RuntimeContext: mockRuntime,
	}

	meteringContext, _ := NewMeteringContext(host, config.MakeGasMapForTests(), uint64(15000))

	input.GasProvided = 1
	err := meteringContext.UseGasForAsyncStep()
	require.Equal(t, arwen.ErrNotEnoughGas, err)

	mockRuntime.SetPointsUsed(0)
	gasProvided := uint64(1_000_000)
	input.GasProvided = gasProvided
	meteringContext.gasForExecution = gasProvided
	gasToLock := meteringContext.ComputeGasLockedForAsync()
	err = meteringContext.UseGasBounded(gasToLock)
	require.Nil(t, err)
	expectedGasLeft := gasProvided - gasToLock
	require.Equal(t, expectedGasLeft, meteringContext.GasLeft())

	mockRuntime.VMInput.CallType = vmcommon.AsynchronousCallBack
	mockRuntime.VMInput.GasLocked = gasToLock
	require.Equal(t, gasToLock, meteringContext.GetGasLocked())

	meteringContext.unlockGasIfAsyncCallback(&input.VMInput)
	err = meteringContext.UseGasForAsyncStep()
	require.Nil(t, err)
	require.Equal(t, gasProvided-1, meteringContext.GasLeft())
}

func TestMeteringContext_GasUsed_NoStacking(t *testing.T) {
	t.Parallel()
	const BlockGasLimit = uint64(15000)

	mockRuntime := &contextmock.RuntimeContextMock{}
	host := &contextmock.VMHostMock{
		RuntimeContext: mockRuntime,
	}

	contractSize := uint64(1000)
	contract := make([]byte, contractSize)
	input := &vmcommon.ContractCallInput{VMInput: vmcommon.VMInput{}}

	mockRuntime.SCCodeSize = contractSize
	mockRuntime.SetVMInput(&input.VMInput)
	mockRuntime.SetPointsUsed(0)

	metering, _ := NewMeteringContext(host, config.MakeGasMapForTests(), BlockGasLimit)

	input.GasProvided = 2000
	metering.InitStateFromContractCallInput(&input.VMInput)
	require.Equal(t, input.GasProvided, metering.initialGasProvided)

	costPerByte := uint64(1)
	gasAfterDeductingInitial := metering.initialGasProvided - costPerByte - contractSize
	_ = metering.DeductInitialGasForExecution(contract)
	require.Equal(t, gasAfterDeductingInitial, metering.GasLeft())

	gasUsed := uint64(400)
	metering.UseGas(gasUsed)
	require.Equal(t, gasAfterDeductingInitial-gasUsed, metering.GasLeft())

	totalGasUsed := metering.initialGasProvided - metering.GasLeft()
	gasUsedByContract := metering.GasSpentByContract()
	require.Equal(t, totalGasUsed, gasUsedByContract)
}

func setUpStackOneLevel(t *testing.T, parentInput *vmcommon.ContractCallInput, childInput *vmcommon.ContractCallInput) (*contextmock.VMHostMock, *contextmock.RuntimeContextMock, uint64) {
	t.Parallel()

	host := &contextmock.VMHostMock{}

	contractSize := uint64(1000)
	contract := make([]byte, contractSize)

	mockRuntime := &contextmock.RuntimeContextMock{}
	host.RuntimeContext = mockRuntime

	mockRuntime.SCCodeSize = contractSize
	mockRuntime.SCAddress = []byte("parent")

	mockRuntime.SetPointsUsed(0)
	mockRuntime.SetVMInput(&parentInput.VMInput)

	metering, _ := NewMeteringContext(host, config.MakeGasMapForTests(), uint64(15000))
	host.MeteringContext = metering
	zeroCodeCosts(metering)

	output, _ := NewOutputContext(host)
	host.OutputContext = output

	parentInput.GasProvided = 4000
	host.MeteringContext.InitStateFromContractCallInput(&parentInput.VMInput)
	_ = host.MeteringContext.DeductInitialGasForExecution(contract)

	parentPointsBeforeStacking := initStateFromChildGetParentPointsBeforeStacking(t, host, childInput)
	return host, mockRuntime, parentPointsBeforeStacking
}

func initStateFromChildGetParentPointsBeforeStacking(t *testing.T, host *contextmock.VMHostMock, childInput *vmcommon.ContractCallInput) uint64 {
	parentGasProvided := uint64(4000)
	parentExecutionGas := uint64(1000)

	require.Equal(t, parentGasProvided-parentExecutionGas, host.MeteringContext.GasLeft())

	parentUsedGas := uint64(400)
	host.MeteringContext.UseGas(parentUsedGas)
	require.Equal(t, parentGasProvided-parentExecutionGas-parentUsedGas, host.MeteringContext.GasLeft())

	gasSpentByContract := host.MeteringContext.GasSpentByContract()
	require.Equal(t, parentExecutionGas+parentUsedGas, gasSpentByContract)

	childProvidedGas := childInput.GasProvided
	host.MeteringContext.UseGas(childProvidedGas)
	parentPointsBeforeStacking := host.RuntimeContext.GetPointsUsed()
	require.Equal(t, childProvidedGas+parentUsedGas, parentPointsBeforeStacking)
	require.Equal(t, parentGasProvided-parentExecutionGas-parentPointsBeforeStacking, host.MeteringContext.GasLeft())

	host.RuntimeContext.SetSCAddress([]byte("child"))
	host.RuntimeContext.SetPointsUsed(0)
	host.RuntimeContext.SetVMInput(&childInput.VMInput)
	host.MeteringContext.PushState()
	host.MeteringContext.InitStateFromContractCallInput(&childInput.VMInput)
	require.Equal(t, childProvidedGas, host.MeteringContext.GetGasProvided())

	return parentPointsBeforeStacking
}

func TestMeteringContext_GasUsed_StackOneLevel(t *testing.T) {
	parentExecutionGas := uint64(1000)
	parentUsedGas := uint64(400)
	parentInput := &vmcommon.ContractCallInput{VMInput: vmcommon.VMInput{}}
	parentInput.CallerAddr = []byte("user")
	parentInput.RecipientAddr = []byte("parent")

	childInput := &vmcommon.ContractCallInput{VMInput: vmcommon.VMInput{}}
	childInput.GasProvided = 500
	childInput.CallerAddr = parentInput.RecipientAddr
	childInput.RecipientAddr = []byte("child")

	host, mockRuntime, parentPointsBeforeStacking := setUpStackOneLevel(t, parentInput, childInput)
	metering := host.MeteringContext
	output := host.OutputContext

	childExecutionGas := uint64(100)
	_ = metering.DeductInitialGasForExecution(make([]byte, childExecutionGas))
	require.Equal(t, childInput.GasProvided-childExecutionGas, metering.GasLeft())

	childUsedGas := uint64(50)
	metering.UseGas(childUsedGas)
	gasRemaining := metering.GasLeft()
	require.Equal(t, childInput.GasProvided-childExecutionGas-childUsedGas, metering.GasLeft())

	gasSpentByContract := metering.GasSpentByContract()
	require.Equal(t, childExecutionGas+childUsedGas, gasSpentByContract)

	_ = output.GetVMOutput()

	// return to the parent
	metering.PopMergeActiveState()
	mockRuntime.SCAddress = []byte("parent")
	mockRuntime.SetPointsUsed(parentPointsBeforeStacking)
	mockRuntime.SetVMInput(&parentInput.VMInput)

	metering.RestoreGas(gasRemaining)
	mockRuntime.IsContractOnStack = false
	require.Equal(t, parentInput.GasProvided-parentExecutionGas-parentPointsBeforeStacking+gasRemaining, metering.GasLeft())

	gasSpentByContract = metering.GasSpentByContract()
	require.Equal(t, parentExecutionGas+parentUsedGas+childUsedGas+childExecutionGas, gasSpentByContract)

	metering.UseGas(50)
	parentUsedGas += 50
	require.Equal(t, parentInput.GasProvided-parentExecutionGas-parentUsedGas-childExecutionGas-childUsedGas, metering.GasLeft())

	gasSpentByContract = metering.GasSpentByContract()
	require.Equal(t, parentExecutionGas+parentUsedGas+childExecutionGas+childUsedGas, gasSpentByContract)

	vmOutput := output.GetVMOutput()

	gasUsed := vmOutput.OutputAccounts["parent"].GasUsed
	require.Equal(t, parentExecutionGas+parentUsedGas, gasUsed)

	gasUsed = vmOutput.OutputAccounts["child"].GasUsed
	require.Equal(t, childExecutionGas+childUsedGas, gasUsed)

	gasRemaining = math.SubUint64(parentInput.GasProvided, gasSpentByContract)
	// calculate gas remaining

	require.Equal(t, gasRemaining, metering.GasLeft())
}

func TestMeteringContext_UpdateGasStateOnFailure_StackOneLevel(t *testing.T) {

	parentExecutionGas := uint64(1000) // this is the contract size, but I chose to keep the convention used on child
	parentUsedGas := uint64(400)
	parentInput := &vmcommon.ContractCallInput{
		VMInput: vmcommon.VMInput{},
	}

	childInput := &vmcommon.ContractCallInput{VMInput: vmcommon.VMInput{}}
	childInput.GasProvided = 500
	childInput.CallerAddr = parentInput.RecipientAddr

	host, mockRuntime, parentPointsBeforeStacking := setUpStackOneLevel(t, parentInput, childInput)

	metering := host.MeteringContext
	output, _ := NewOutputContext(host)
	host.OutputContext = output

	childExecutionGas := uint64(600)
	_ = metering.DeductInitialGasForExecution(make([]byte, childExecutionGas))
	require.Equal(t, childInput.GasProvided, metering.GasLeft()) // not enough gas provided. It will remain the same for now

	gasRemaining := metering.GasLeft()
	gasSpentByContract := metering.GasSpentByContract()
	require.Equal(t, uint64(0), gasSpentByContract)

	metering.UpdateGasStateOnFailure(output.outputState)

	// return to the parent
	metering.PopSetActiveState()
	mockRuntime.SCAddress = []byte("parent")
	mockRuntime.SetPointsUsed(parentPointsBeforeStacking)
	mockRuntime.SetVMInput(&parentInput.VMInput)

	metering.RestoreGas(gasRemaining)
	mockRuntime.IsContractOnStack = false
	require.Equal(t, parentInput.GasProvided-parentExecutionGas-parentPointsBeforeStacking+gasRemaining, metering.GasLeft())

	gasSpentByContract = metering.GasSpentByContract()
	require.Equal(t, parentExecutionGas+parentUsedGas, gasSpentByContract) // child execution will fail due to insufficient gas

	metering.UpdateGasStateOnFailure(output.outputState)

	// after update all gas will be used

	gasUsed := output.outputState.OutputAccounts["parent"].GasUsed
	require.Equal(t, parentInput.GasProvided, gasUsed)

	gasUsed = output.outputState.OutputAccounts["child"].GasUsed
	require.Equal(t, childInput.GasProvided, gasUsed)

	gasRemaining = math.SubUint64(parentInput.GasProvided, gasSpentByContract)
	// calculate gas remaining

	require.Equal(t, int(gasRemaining), int(metering.GasLeft()))
}

func zeroCodeCosts(context *meteringContext) {
	context.GasSchedule().BaseOperationCost.GetCode = 0
}

func TestMeteringContext_TrackGasUsedByBuiltinFunction_GasRemaining(t *testing.T) {
	t.Parallel()

	mockRuntime := &contextmock.RuntimeContextMock{}

	host := &contextmock.VMHostMock{
		RuntimeContext: mockRuntime,
	}

	names := make(vmcommon.FunctionNames)
	contractSize := uint64(1000)
	mockRuntime.SCCodeSize = contractSize
	mockRuntime.SCAddress = []byte("parent")

	mockRuntime.SetPointsUsed(0)

	input := &vmcommon.ContractCallInput{
		VMInput:  vmcommon.VMInput{},
		Function: "callBuiltinClaim",
	}
	mockRuntime.SetVMInput(&input.VMInput)

	metering, _ := NewMeteringContext(host, config.MakeGasMapForTests(), uint64(15000))
	host.MeteringContext = metering
	zeroCodeCosts(metering)

	input.GasProvided = 10000
	metering.InitStateFromContractCallInput(&input.VMInput)
	require.Equal(t, input.GasProvided, metering.GasLeft())

	var empty struct{}
	names["builtinClaim"] = empty
	host.SetProtocolBuiltinFunctions(names)

	vmOutput := &vmcommon.VMOutput{
		GasRemaining: 5000,
	}

	postBuiltinInput := &vmcommon.ContractCallInput{
		VMInput: vmcommon.VMInput{
			GasProvided: 300,
		},
	}

	metering.TrackGasUsedByBuiltinFunction(input, vmOutput, postBuiltinInput)
	require.Equal(t, vmOutput.GasRemaining+postBuiltinInput.GasProvided, metering.GasLeft())
}
