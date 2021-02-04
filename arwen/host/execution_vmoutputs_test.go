package host

import (
	"fmt"
	"math/big"

	"github.com/ElrondNetwork/arwen-wasm-vm/arwen"
	"github.com/ElrondNetwork/elrond-go/core/vmcommon"
)

var parentKeyA = []byte("parentKeyA......................")
var parentKeyB = []byte("parentKeyB......................")
var childKey = []byte("childKey........................")
var parentDataA = []byte("parentDataA")
var parentDataB = []byte("parentDataB")
var childData = []byte("childData")
var parentFinishA = []byte("parentFinishA")
var parentFinishB = []byte("parentFinishB")
var childFinish = []byte("childFinish")
var parentTransferReceiver = []byte("\x00\x00\x00\x00\x00\x00\x00\x00\x0f\x0fparentTransferReceiver")
var childTransferReceiver = []byte("\x00\x00\x00\x00\x00\x00\x00\x00\x0f\x0fchildTransferReceiver.")
var parentTransferValue = int64(42)
var parentTransferData = []byte("parentTransferData")

var recursiveIterationCounterKey = []byte("recursiveIterationCounter.......")
var recursiveIterationBigCounterKey = []byte("recursiveIterationBigCounter....")

var gasProvided = uint64(1000000)

var parentCompilationCostSameCtx uint64
var childCompilationCostSameCtx uint64

var parentCompilationCostDestCtx uint64
var childCompilationCostDestCtx uint64

var vaultAddress = []byte("\x00\x00\x00\x00\x00\x00\x00\x00\x0f\x0fvaultAddress..........")
var thirdPartyAddress = []byte("\x00\x00\x00\x00\x00\x00\x00\x00\x0f\x0fthirdPartyAddress.....")

func init() {
	parentCompilationCostSameCtx = uint64(len(arwen.GetTestSCCode("exec-same-ctx-parent", "../../")))
	childCompilationCostSameCtx = uint64(len(arwen.GetTestSCCode("exec-same-ctx-child", "../../")))

	parentCompilationCostDestCtx = uint64(len(arwen.GetTestSCCode("exec-dest-ctx-parent", "../../")))
	childCompilationCostDestCtx = uint64(len(arwen.GetTestSCCode("exec-dest-ctx-child", "../../")))
}

func expectedVMOutputSameCtxPrepare(_ []byte) *vmcommon.VMOutput {
	vmOutput := arwen.MakeVMOutput()
	vmOutput.ReturnCode = vmcommon.Ok

	parentAccount := arwen.AddNewOutputAccount(
		vmOutput,
		parentAddress,
		-parentTransferValue,
		nil,
	)
	parentAccount.Balance = big.NewInt(1000)
	parentAccount.GasUsed = 3405

	_ = arwen.AddNewOutputAccount(
		vmOutput,
		parentTransferReceiver,
		parentTransferValue,
		parentTransferData,
	)

	arwen.SetStorageUpdate(parentAccount, parentKeyA, parentDataA)
	arwen.SetStorageUpdate(parentAccount, parentKeyB, parentDataB)

	arwen.AddFinishData(vmOutput, parentFinishA)
	arwen.AddFinishData(vmOutput, parentFinishB)
	arwen.AddFinishData(vmOutput, []byte("succ"))

	expectedExecutionCost := uint64(138)
	gas := gasProvided
	gas -= parentCompilationCostSameCtx
	gas -= expectedExecutionCost
	vmOutput.GasRemaining = gas

	return vmOutput
}

func expectedVMOutputSameCtxWrongContractCalled(code []byte) *vmcommon.VMOutput {
	vmOutput := expectedVMOutputSameCtxPrepare(code)

	arwen.AddFinishData(vmOutput, []byte("fail"))

	executionCostBeforeExecuteAPI := uint64(156)
	executeAPICost := uint64(39)
	gasLostOnFailure := uint64(50000)
	finalCost := uint64(44)
	gas := gasProvided
	gas -= parentCompilationCostSameCtx
	gas -= executionCostBeforeExecuteAPI
	gas -= executeAPICost
	gas -= gasLostOnFailure
	gas -= finalCost
	vmOutput.GasRemaining = gas

	return vmOutput
}

func expectedVMOutputSameCtxOutOfGas(_ []byte, _ []byte) *vmcommon.VMOutput {
	vmOutput := arwen.MakeVMOutput()

	vmOutput.ReturnCode = vmcommon.Ok

	parentAccount := arwen.AddNewOutputAccount(
		vmOutput,
		parentAddress,
		0,
		nil,
	)

	arwen.SetStorageUpdate(parentAccount, parentKeyA, parentDataA)
	arwen.SetStorageUpdate(parentAccount, parentKeyB, parentDataB)
	arwen.AddFinishData(vmOutput, parentFinishA)

	arwen.AddFinishData(vmOutput, []byte("fail"))

	executionCostBeforeExecuteAPI := uint64(90)
	executeAPICost := uint64(1)
	gasLostOnFailure := uint64(3500)
	finalCost := uint64(54)
	gas := gasProvided
	gas -= parentCompilationCostSameCtx
	gas -= executionCostBeforeExecuteAPI
	gas -= executeAPICost
	gas -= gasLostOnFailure
	gas -= finalCost
	vmOutput.GasRemaining = gas

	return vmOutput
}

func expectedVMOutputSameCtxSimple(_ []byte, _ []byte) *vmcommon.VMOutput {
	vmOutput := arwen.MakeVMOutput()

	arwen.AddFinishData(vmOutput, []byte("child"))
	arwen.AddFinishData(vmOutput, []byte{})
	for i := 1; i < 100; i++ {
		arwen.AddFinishData(vmOutput, []byte{byte(i)})
	}
	arwen.AddFinishData(vmOutput, []byte{})
	arwen.AddFinishData(vmOutput, []byte("child"))
	arwen.AddFinishData(vmOutput, []byte{})
	for i := 1; i < 100; i++ {
		arwen.AddFinishData(vmOutput, []byte{byte(i)})
	}
	arwen.AddFinishData(vmOutput, []byte{})
	arwen.AddFinishData(vmOutput, []byte("parent"))

	parentAccount := arwen.AddNewOutputAccount(
		vmOutput,
		parentAddress,
		-198,
		nil,
	)
	parentAccount.Balance = big.NewInt(1000)
	parentAccount.GasUsed = 521

	childAccount := arwen.AddNewOutputAccount(
		vmOutput,
		childAddress,
		198,
		nil,
	)
	childAccount.GasUsed = 3435 // TODO: double this when fixed

	executionCost := parentAccount.GasUsed + 2*childAccount.GasUsed
	vmOutput.GasRemaining = gasProvided - executionCost

	return vmOutput
}

func expectedVMOutputSameCtxSuccessfulChildCall(parentCode []byte, _ []byte) *vmcommon.VMOutput {
	vmOutput := expectedVMOutputSameCtxPrepare(parentCode)

	parentAccount := vmOutput.OutputAccounts[string(parentAddress)]
	parentAccount.BalanceDelta = big.NewInt(-141)
	parentAccount.GasUsed = 3612

	childAccount := arwen.AddNewOutputAccount(
		vmOutput,
		childAddress,
		3,
		nil,
	)
	childAccount.Balance = big.NewInt(1000)

	executeAPICost := uint64(39)
	childExecutionCost := uint64(437)
	childAccount.GasUsed = childCompilationCostSameCtx + childExecutionCost

	_ = arwen.AddNewOutputAccount(
		vmOutput,
		childTransferReceiver,
		96,
		[]byte("qwerty"),
	)

	// The child SC has stored data on the parent's storage
	arwen.SetStorageUpdate(parentAccount, childKey, childData)

	// The called child SC will output some arbitrary data, but also data that it
	// has read from the storage of the parent.
	arwen.AddFinishData(vmOutput, childFinish)
	arwen.AddFinishData(vmOutput, parentDataA)
	for _, c := range parentDataA {
		arwen.AddFinishData(vmOutput, []byte{c})
	}
	arwen.AddFinishData(vmOutput, parentDataB)
	for _, c := range parentDataB {
		arwen.AddFinishData(vmOutput, []byte{c})
	}
	arwen.AddFinishData(vmOutput, []byte("child ok"))
	arwen.AddFinishData(vmOutput, []byte("succ"))
	arwen.AddFinishData(vmOutput, []byte("succ"))

	parentGasBeforeExecuteAPI := uint64(172)
	finalCost := uint64(134)
	gas := gasProvided
	gas -= parentCompilationCostSameCtx
	gas -= parentGasBeforeExecuteAPI
	gas -= executeAPICost
	gas -= childCompilationCostSameCtx
	gas -= childExecutionCost
	gas -= finalCost
	vmOutput.GasRemaining = gas

	return vmOutput
}

func expectedVMOutputSameCtxSuccessfulChildCallBigInts(_ []byte, _ []byte) *vmcommon.VMOutput {
	vmOutput := arwen.MakeVMOutput()
	vmOutput.ReturnCode = vmcommon.Ok

	parentAccount := arwen.AddNewOutputAccount(
		vmOutput,
		parentAddress,
		-99,
		nil,
	)
	parentAccount.Balance = big.NewInt(1000)
	parentAccount.GasUsed = 3461

	childAccount := arwen.AddNewOutputAccount(
		vmOutput,
		childAddress,
		99,
		nil,
	)
	childExecutionCost := uint64(108)
	childAccount.GasUsed = childCompilationCostSameCtx + childExecutionCost

	// The child SC will output "child ok" if it could read some expected Big
	// Ints directly from the parent's context.
	arwen.AddFinishData(vmOutput, []byte("child ok"))
	arwen.AddFinishData(vmOutput, []byte("succ"))
	arwen.AddFinishData(vmOutput, []byte("succ"))

	parentGasBeforeExecuteAPI := uint64(114)
	executeAPICost := uint64(13)
	finalCost := uint64(67)
	gas := gasProvided
	gas -= parentCompilationCostSameCtx
	gas -= parentGasBeforeExecuteAPI
	gas -= executeAPICost
	gas -= childCompilationCostSameCtx
	gas -= childExecutionCost
	gas -= finalCost
	vmOutput.GasRemaining = gas
	return vmOutput
}

func expectedVMOutputSameCtxRecursiveDirect(_ []byte, recursiveCalls int) *vmcommon.VMOutput {
	vmOutput := arwen.MakeVMOutput()

	account := arwen.AddNewOutputAccount(
		vmOutput,
		parentAddress,
		0,
		nil,
	)
	account.Balance = big.NewInt(1000)
	account.BalanceDelta = big.NewInt(0).Sub(big.NewInt(1), big.NewInt(1))
	account.GasUsed = 25424

	for i := recursiveCalls; i >= 0; i-- {
		finishString := fmt.Sprintf("Rfinish%03d", i)
		arwen.AddFinishData(vmOutput, []byte(finishString))
	}

	for i := recursiveCalls - 1; i >= 0; i-- {
		arwen.AddFinishData(vmOutput, []byte("succ"))
	}

	for i := 0; i <= recursiveCalls; i++ {
		key := fmt.Sprintf("Rkey%03d.........................", i)
		value := fmt.Sprintf("Rvalue%03d", i)
		arwen.SetStorageUpdateStrings(account, key, value)
	}

	arwen.SetStorageUpdate(account, recursiveIterationCounterKey, []byte{byte(recursiveCalls + 1)})
	arwen.SetStorageUpdate(account, recursiveIterationBigCounterKey, big.NewInt(int64(recursiveCalls+1)).Bytes())

	return vmOutput
}

func expectedVMOutputSameCtxRecursiveDirectErrMaxInstances(_ []byte, recursiveCalls int) *vmcommon.VMOutput {
	vmOutput := arwen.MakeVMOutput()

	account := arwen.AddNewOutputAccount(
		vmOutput,
		parentAddress,
		0,
		nil,
	)

	finishString := fmt.Sprintf("Rfinish%03d", recursiveCalls)
	arwen.AddFinishData(vmOutput, []byte(finishString))

	key := fmt.Sprintf("Rkey%03d.........................", recursiveCalls)
	value := fmt.Sprintf("Rvalue%03d", recursiveCalls)
	arwen.SetStorageUpdateStrings(account, key, value)

	arwen.AddFinishData(vmOutput, []byte("fail"))
	arwen.SetStorageUpdate(account, recursiveIterationCounterKey, []byte{byte(1)})

	return vmOutput
}

func expectedVMOutputSameCtxRecursiveMutualMethods(_ []byte, recursiveCalls int) *vmcommon.VMOutput {
	vmOutput := arwen.MakeVMOutput()

	account := arwen.AddNewOutputAccount(
		vmOutput,
		parentAddress,
		0,
		nil,
	)
	account.Balance = big.NewInt(1000)
	account.BalanceDelta = big.NewInt(0).Sub(big.NewInt(1), big.NewInt(1))
	account.GasUsed = 29593

	arwen.SetStorageUpdate(account, recursiveIterationCounterKey, []byte{byte(recursiveCalls + 1)})
	arwen.SetStorageUpdate(account, recursiveIterationBigCounterKey, big.NewInt(int64(recursiveCalls+1)).Bytes())

	arwen.AddFinishData(vmOutput, []byte("start recursive mutual calls"))

	for i := 0; i <= recursiveCalls; i++ {
		var finishData string
		var key string
		var value string
		iteration := recursiveCalls - i
		if i%2 == 0 {
			finishData = fmt.Sprintf("Afinish%03d", iteration)
			key = fmt.Sprintf("Akey%03d.........................", iteration)
			value = fmt.Sprintf("Avalue%03d", iteration)
		} else {
			finishData = fmt.Sprintf("Bfinish%03d", iteration)
			key = fmt.Sprintf("Bkey%03d.........................", iteration)
			value = fmt.Sprintf("Bvalue%03d", iteration)
		}
		arwen.SetStorageUpdateStrings(account, key, value)
		arwen.AddFinishData(vmOutput, []byte(finishData))
	}

	for i := recursiveCalls; i >= 0; i-- {
		arwen.AddFinishData(vmOutput, []byte("succ"))
	}

	arwen.AddFinishData(vmOutput, []byte("end recursive mutual calls"))

	return vmOutput
}

func expectedVMOutputSameCtxRecursiveMutualSCs(_ []byte, _ []byte, recursiveCalls int) *vmcommon.VMOutput {
	vmOutput := arwen.MakeVMOutput()

	parentAccount := arwen.AddNewOutputAccount(
		vmOutput,
		parentAddress,
		0,
		nil,
	)
	parentAccount.Balance = big.NewInt(1000)
	parentAccount.GasUsed = 5426

	childAccount := arwen.AddNewOutputAccount(
		vmOutput,
		childAddress,
		0,
		nil,
	)
	childAccount.Balance = big.NewInt(1000)
	childAccount.GasUsed = 3652

	if recursiveCalls%2 == 1 {
		parentAccount.BalanceDelta = big.NewInt(-5)
		childAccount.BalanceDelta = big.NewInt(5)
	} else {
		parentAccount.BalanceDelta = big.NewInt(0).Sub(big.NewInt(1), big.NewInt(1))
		childAccount.BalanceDelta = big.NewInt(0).Sub(big.NewInt(1), big.NewInt(1))
	}

	for i := 0; i <= recursiveCalls; i++ {
		var finishData string
		var key string
		var value string
		iteration := recursiveCalls - i
		if i%2 == 0 {
			finishData = fmt.Sprintf("Pfinish%03d", iteration)
			key = fmt.Sprintf("Pkey%03d.........................", iteration)
			value = fmt.Sprintf("Pvalue%03d", iteration)
		} else {
			finishData = fmt.Sprintf("Cfinish%03d", iteration)
			key = fmt.Sprintf("Ckey%03d.........................", iteration)
			value = fmt.Sprintf("Cvalue%03d", iteration)
		}
		arwen.SetStorageUpdateStrings(parentAccount, key, value)
		arwen.AddFinishData(vmOutput, []byte(finishData))
	}

	for i := recursiveCalls - 1; i >= 0; i-- {
		arwen.AddFinishData(vmOutput, []byte("succ"))
	}

	arwen.SetStorageUpdate(parentAccount, recursiveIterationCounterKey, []byte{byte(recursiveCalls + 1)})
	arwen.SetStorageUpdate(parentAccount, recursiveIterationBigCounterKey, big.NewInt(int64(recursiveCalls+1)).Bytes())

	return vmOutput
}

func expectedVMOutputDestCtxBuiltinFunctions1(_ []byte) *vmcommon.VMOutput {
	vmOutput := arwen.MakeVMOutput()

	gasProvided := uint64(100000)
	account := arwen.AddNewOutputAccount(
		vmOutput,
		parentAddress,
		42,
		nil,
	)
	account.Balance = big.NewInt(1000)
	account.GasUsed = 1541

	vmOutput.GasRemaining = gasProvided - account.GasUsed

	arwen.AddFinishData(vmOutput, []byte("succ"))

	return vmOutput
}

func expectedVMOutputDestCtxBuiltinFunctions2(_ []byte) *vmcommon.VMOutput {
	vmOutput := arwen.MakeVMOutput()

	gasProvided := uint64(100000)
	account := arwen.AddNewOutputAccount(
		vmOutput,
		parentAddress,
		0,
		nil,
	)
	account.Balance = big.NewInt(1000)
	account.BalanceDelta = big.NewInt(0)
	account.GasUsed = 1541

	arwen.AddFinishData(vmOutput, []byte("succ"))

	vmOutput.GasRemaining = 98459
	vmOutput.GasRemaining = gasProvided - account.GasUsed

	return vmOutput
}

func expectedVMOutputDestCtxBuiltinFunctions3(_ []byte) *vmcommon.VMOutput {
	vmOutput := arwen.MakeVMOutput()

	vmOutput.ReturnCode = vmcommon.ExecutionFailed
	vmOutput.ReturnMessage = "not enough gas"
	vmOutput.GasRemaining = 0
	vmOutput.ReturnData = nil
	vmOutput.OutputAccounts = nil
	vmOutput.TouchedAccounts = nil
	vmOutput.DeletedAccounts = nil
	vmOutput.Logs = nil

	return vmOutput
}

func expectedVMOutputDestCtxPrepare(_ []byte) *vmcommon.VMOutput {
	vmOutput := arwen.MakeVMOutput()
	vmOutput.ReturnCode = vmcommon.Ok

	parentAccount := arwen.AddNewOutputAccount(
		vmOutput,
		parentAddress,
		-parentTransferValue,
		nil,
	)
	parentAccount.Balance = big.NewInt(1000)
	parentAccount.GasUsed = 3093

	_ = arwen.AddNewOutputAccount(
		vmOutput,
		parentTransferReceiver,
		parentTransferValue,
		parentTransferData,
	)

	arwen.SetStorageUpdate(parentAccount, parentKeyA, parentDataA)
	arwen.SetStorageUpdate(parentAccount, parentKeyB, parentDataB)

	arwen.AddFinishData(vmOutput, parentFinishA)
	arwen.AddFinishData(vmOutput, parentFinishB)
	arwen.AddFinishData(vmOutput, []byte("succ"))

	expectedExecutionCost := uint64(138)
	gas := gasProvided
	gas -= parentCompilationCostDestCtx
	gas -= expectedExecutionCost
	vmOutput.GasRemaining = gas

	return vmOutput
}

func expectedVMOutputDestCtxWrongContractCalled(parentCode []byte) *vmcommon.VMOutput {
	vmOutput := expectedVMOutputSameCtxPrepare(parentCode)

	parentAccount := vmOutput.OutputAccounts[string(parentAddress)]
	parentAccount.BalanceDelta = big.NewInt(-42)

	arwen.AddFinishData(vmOutput, []byte("fail"))

	executionCostBeforeExecuteAPI := uint64(156)
	executeAPICost := uint64(42)
	gasLostOnFailure := uint64(10000)
	finalCost := uint64(44)
	gas := gasProvided
	gas -= parentCompilationCostDestCtx
	gas -= executionCostBeforeExecuteAPI
	gas -= executeAPICost
	gas -= gasLostOnFailure
	gas -= finalCost
	vmOutput.GasRemaining = gas

	return vmOutput
}

func expectedVMOutputDestCtxOutOfGas(_ []byte) *vmcommon.VMOutput {
	vmOutput := arwen.MakeVMOutput()

	vmOutput.ReturnCode = vmcommon.Ok

	parentAccount := arwen.AddNewOutputAccount(
		vmOutput,
		parentAddress,
		0,
		nil,
	)

	arwen.SetStorageUpdate(parentAccount, parentKeyA, parentDataA)
	arwen.SetStorageUpdate(parentAccount, parentKeyB, parentDataB)
	arwen.AddFinishData(vmOutput, parentFinishA)

	arwen.AddFinishData(vmOutput, []byte("fail"))

	executionCostBeforeExecuteAPI := uint64(90)
	executeAPICost := uint64(1)
	gasLostOnFailure := uint64(3500)
	finalCost := uint64(54)
	gas := gasProvided
	gas -= parentCompilationCostDestCtx
	gas -= executionCostBeforeExecuteAPI
	gas -= executeAPICost
	gas -= gasLostOnFailure
	gas -= finalCost
	vmOutput.GasRemaining = gas

	return vmOutput
}

func expectedVMOutputDestCtxSuccessfulChildCall(parentCode []byte, _ []byte) *vmcommon.VMOutput {
	vmOutput := expectedVMOutputSameCtxPrepare(parentCode)

	parentAccount := vmOutput.OutputAccounts[string(parentAddress)]
	parentAccount.BalanceDelta = big.NewInt(-141)
	parentAccount.GasUsed = 3228

	childAccount := arwen.AddNewOutputAccount(
		vmOutput,
		childAddress,
		99-12,
		nil,
	)
	childAccount.Balance = big.NewInt(1000)
	childAccount.GasUsed = 2256

	_ = arwen.AddNewOutputAccount(
		vmOutput,
		childTransferReceiver,
		12,
		[]byte("Second sentence."),
	)

	arwen.SetStorageUpdate(childAccount, childKey, childData)

	arwen.AddFinishData(vmOutput, childFinish)
	arwen.AddFinishData(vmOutput, []byte("succ"))
	arwen.AddFinishData(vmOutput, []byte("succ"))

	parentGasBeforeExecuteAPI := uint64(168)
	executeAPICost := uint64(42)
	childExecutionCost := uint64(91)
	finalCost := uint64(65)
	gas := gasProvided
	gas -= parentCompilationCostDestCtx
	gas -= parentGasBeforeExecuteAPI
	gas -= executeAPICost
	gas -= childCompilationCostDestCtx
	gas -= childExecutionCost
	gas -= finalCost
	vmOutput.GasRemaining = gas
	return vmOutput
}

func expectedVMOutputDestCtxSuccessfulChildCallBigInts(_ []byte, _ []byte) *vmcommon.VMOutput {
	vmOutput := arwen.MakeVMOutput()
	vmOutput.ReturnCode = vmcommon.Ok

	parentAccount := arwen.AddNewOutputAccount(
		vmOutput,
		parentAddress,
		-99,
		nil,
	)
	parentAccount.Balance = big.NewInt(1000)
	parentAccount.GasUsed = 3150

	childAccount := arwen.AddNewOutputAccount(
		vmOutput,
		childAddress,
		99,
		nil,
	)
	childAccount.GasUsed = 2265

	// The child SC will output "child ok" if it could NOT read the Big Ints from
	// the parent's context.
	arwen.AddFinishData(vmOutput, []byte("child ok"))
	arwen.AddFinishData(vmOutput, []byte("succ"))
	arwen.AddFinishData(vmOutput, []byte("succ"))

	parentGasBeforeExecuteAPI := uint64(115)
	executeAPICost := uint64(13)
	childExecutionCost := uint64(101)
	finalCost := uint64(68)
	gas := gasProvided
	gas -= parentCompilationCostDestCtx
	gas -= parentGasBeforeExecuteAPI
	gas -= executeAPICost
	gas -= childCompilationCostDestCtx
	gas -= childExecutionCost
	gas -= finalCost
	vmOutput.GasRemaining = gas
	return vmOutput
}

func expectedVMOutputDestCtxRecursiveDirect(_ []byte, recursiveCalls int) *vmcommon.VMOutput {
	vmOutput := arwen.MakeVMOutput()

	account := arwen.AddNewOutputAccount(
		vmOutput,
		parentAddress,
		0,
		nil,
	)
	account.Balance = big.NewInt(1000)
	account.BalanceDelta = big.NewInt(0).Sub(big.NewInt(1), big.NewInt(1))
	account.GasUsed = 29670

	for i := recursiveCalls; i >= 0; i-- {
		finishString := fmt.Sprintf("Rfinish%03d", i)
		arwen.AddFinishData(vmOutput, []byte(finishString))
	}

	for i := recursiveCalls - 1; i >= 0; i-- {
		arwen.AddFinishData(vmOutput, []byte("succ"))
	}

	for i := 0; i <= recursiveCalls; i++ {
		key := fmt.Sprintf("Rkey%03d.........................", i)
		value := fmt.Sprintf("Rvalue%03d", i)
		arwen.SetStorageUpdateStrings(account, key, value)
	}

	arwen.SetStorageUpdate(account, recursiveIterationCounterKey, []byte{byte(recursiveCalls + 1)})
	arwen.SetStorageUpdate(account, recursiveIterationBigCounterKey, big.NewInt(int64(1)).Bytes())

	return vmOutput
}

func expectedVMOutputDestCtxRecursiveMutualMethods(_ []byte, recursiveCalls int) *vmcommon.VMOutput {
	vmOutput := arwen.MakeVMOutput()

	account := arwen.AddNewOutputAccount(
		vmOutput,
		parentAddress,
		0,
		nil,
	)
	account.Balance = big.NewInt(1000)
	account.BalanceDelta = big.NewInt(0).Sub(big.NewInt(1), big.NewInt(1))
	account.GasUsed = 38083

	arwen.SetStorageUpdate(account, recursiveIterationCounterKey, []byte{byte(recursiveCalls + 1)})
	arwen.SetStorageUpdate(account, recursiveIterationBigCounterKey, big.NewInt(int64(1)).Bytes())

	arwen.AddFinishData(vmOutput, []byte("start recursive mutual calls"))

	for i := 0; i <= recursiveCalls; i++ {
		var finishData string
		var key string
		var value string
		iteration := recursiveCalls - i
		if i%2 == 0 {
			finishData = fmt.Sprintf("Afinish%03d", iteration)
			key = fmt.Sprintf("Akey%03d.........................", iteration)
			value = fmt.Sprintf("Avalue%03d", iteration)
		} else {
			finishData = fmt.Sprintf("Bfinish%03d", iteration)
			key = fmt.Sprintf("Bkey%03d.........................", iteration)
			value = fmt.Sprintf("Bvalue%03d", iteration)
		}
		arwen.SetStorageUpdateStrings(account, key, value)
		arwen.AddFinishData(vmOutput, []byte(finishData))
	}

	for i := recursiveCalls; i >= 0; i-- {
		arwen.AddFinishData(vmOutput, []byte("succ"))
	}

	arwen.AddFinishData(vmOutput, []byte("end recursive mutual calls"))

	return vmOutput
}

func expectedVMOutputDestCtxRecursiveMutualSCs(_ []byte, _ []byte, recursiveCalls int) *vmcommon.VMOutput {
	vmOutput := arwen.MakeVMOutput()

	parentIterations := (recursiveCalls / 2) + (recursiveCalls % 2)
	childIterations := recursiveCalls - parentIterations

	balanceDelta := int64(5*parentIterations - 3*childIterations)

	parentAccount := arwen.AddNewOutputAccount(
		vmOutput,
		parentAddress,
		0,
		nil,
	)
	parentAccount.Balance = big.NewInt(1000)
	parentAccount.BalanceDelta = big.NewInt(-balanceDelta)
	parentAccount.GasUsed = 7252

	childAccount := arwen.AddNewOutputAccount(
		vmOutput,
		childAddress,
		0,
		nil,
	)
	childAccount.Balance = big.NewInt(1000)
	childAccount.BalanceDelta = big.NewInt(balanceDelta)
	childAccount.GasUsed = 5464

	for i := 0; i <= recursiveCalls; i++ {
		var finishData string
		var key string
		var value string
		iteration := recursiveCalls - i
		if i%2 == 0 {
			finishData = fmt.Sprintf("Pfinish%03d", iteration)
			key = fmt.Sprintf("Pkey%03d.........................", iteration)
			value = fmt.Sprintf("Pvalue%03d", iteration)
			arwen.SetStorageUpdateStrings(parentAccount, key, value)
		} else {
			finishData = fmt.Sprintf("Cfinish%03d", iteration)
			key = fmt.Sprintf("Ckey%03d.........................", iteration)
			value = fmt.Sprintf("Cvalue%03d", iteration)
			arwen.SetStorageUpdateStrings(childAccount, key, value)
		}
		arwen.AddFinishData(vmOutput, []byte(finishData))
	}

	for i := recursiveCalls - 1; i >= 0; i-- {
		arwen.AddFinishData(vmOutput, []byte("succ"))
	}

	counterValue := (recursiveCalls + recursiveCalls%2) / 2
	arwen.SetStorageUpdate(parentAccount, recursiveIterationCounterKey, []byte{byte(counterValue + 1)})
	arwen.SetStorageUpdate(childAccount, recursiveIterationCounterKey, []byte{byte(counterValue)})
	if recursiveCalls%2 == 0 {
		arwen.SetStorageUpdate(parentAccount, recursiveIterationBigCounterKey, big.NewInt(int64(1)).Bytes())
	} else {
		arwen.SetStorageUpdate(childAccount, recursiveIterationBigCounterKey, big.NewInt(int64(1)).Bytes())
	}

	return vmOutput
}

func expectedVMOutputDestCtxByCallerSimpleTransfer(value int64) *vmcommon.VMOutput {
	vmOutput := arwen.MakeVMOutput()

	parentAccount := arwen.AddNewOutputAccount(
		vmOutput,
		parentAddress,
		0,
		nil,
	)
	parentAccount.Balance = nil
	parentAccount.GasUsed = 762

	childAccount := arwen.AddNewOutputAccount(
		vmOutput,
		childAddress,
		0,
		nil,
	)
	childAccount.Balance = big.NewInt(1000)
	childAccount.BalanceDelta = big.NewInt(-value)
	childAccount.GasUsed = 667

	userAccount := arwen.AddNewOutputAccount(
		vmOutput,
		userAddress,
		0,
		nil,
	)
	userAccount.BalanceDelta = big.NewInt(value)
	userAccount.OutputTransfers = append(userAccount.OutputTransfers, vmcommon.OutputTransfer{
		Value:     big.NewInt(value),
		GasLimit:  0,
		GasLocked: 0,
		Data:      []byte{},
		CallType:  vmcommon.DirectCall,
	})

	arwen.AddFinishData(vmOutput, []byte("sent"))
	arwen.AddFinishData(vmOutput, []byte("child called"))
	return vmOutput
}

func expectedVMOutputAsyncCall(_ []byte, _ []byte) *vmcommon.VMOutput {
	vmOutput := arwen.MakeVMOutput()

	parentAccount := arwen.AddNewOutputAccount(
		vmOutput,
		parentAddress,
		-10,
		nil,
	)
	parentAccount.Balance = big.NewInt(1000)
	parentAccount.GasUsed = 104753
	arwen.SetStorageUpdate(parentAccount, parentKeyA, parentDataA)
	arwen.SetStorageUpdate(parentAccount, parentKeyB, parentDataB)
	arwen.AddFinishData(vmOutput, parentFinishA)
	arwen.AddFinishData(vmOutput, parentFinishB)

	thirdPartyAccount := arwen.AddNewOutputAccount(
		vmOutput,
		thirdPartyAddress,
		3,
		[]byte("hello"),
	)
	outTransfer := vmcommon.OutputTransfer{Data: []byte(" there"), Value: big.NewInt(3)}
	thirdPartyAccount.OutputTransfers = append(thirdPartyAccount.OutputTransfers, outTransfer)
	thirdPartyAccount.BalanceDelta = big.NewInt(6)

	childAccount := arwen.AddNewOutputAccount(
		vmOutput,
		childAddress,
		0,
		nil,
	)
	childAccount.Balance = big.NewInt(1000)
	childAccount.GasUsed = 1296
	arwen.SetStorageUpdate(childAccount, childKey, childData)

	_ = arwen.AddNewOutputAccount(
		vmOutput,
		vaultAddress,
		4,
		[]byte{},
	)

	arwen.AddFinishData(vmOutput, []byte{0})
	arwen.AddFinishData(vmOutput, []byte("thirdparty"))
	arwen.AddFinishData(vmOutput, []byte("vault"))
	arwen.AddFinishData(vmOutput, []byte{0})
	arwen.AddFinishData(vmOutput, []byte("succ"))

	return vmOutput
}

func expectedVMOutputAsyncCallChildFails(_ []byte, _ []byte) *vmcommon.VMOutput {
	vmOutput := arwen.MakeVMOutput()

	parentAccount := arwen.AddNewOutputAccount(
		vmOutput,
		parentAddress,
		-7,
		nil,
	)
	parentAccount.Balance = big.NewInt(1000)
	parentAccount.GasUsed = 3928
	arwen.SetStorageUpdate(parentAccount, parentKeyA, parentDataA)
	arwen.SetStorageUpdate(parentAccount, parentKeyB, parentDataB)
	arwen.AddFinishData(vmOutput, parentFinishA)
	arwen.AddFinishData(vmOutput, parentFinishB)

	_ = arwen.AddNewOutputAccount(
		vmOutput,
		thirdPartyAddress,
		3,
		[]byte("hello"),
	)

	childAccount := arwen.AddNewOutputAccount(
		vmOutput,
		childAddress,
		0,
		nil,
	)
	childAccount.Balance = big.NewInt(1000)

	_ = arwen.AddNewOutputAccount(
		vmOutput,
		vaultAddress,
		4,
		[]byte{},
	)

	arwen.AddFinishData(vmOutput, []byte("succ"))

	return vmOutput
}

func expectedVMOutputAsyncCallCallBackFails(_ []byte, _ []byte) *vmcommon.VMOutput {
	vmOutput := arwen.MakeVMOutput()

	parentAccount := arwen.AddNewOutputAccount(
		vmOutput,
		parentAddress,
		-10,
		nil,
	)
	parentAccount.Balance = big.NewInt(1000)
	parentAccount.GasUsed = 197173
	arwen.SetStorageUpdate(parentAccount, parentKeyA, parentDataA)
	arwen.SetStorageUpdate(parentAccount, parentKeyB, parentDataB)
	arwen.AddFinishData(vmOutput, parentFinishA)
	arwen.AddFinishData(vmOutput, parentFinishB)

	_ = arwen.AddNewOutputAccount(
		vmOutput,
		thirdPartyAddress,
		3,
		[]byte("hello"),
	)
	outTransfer2 := vmcommon.OutputTransfer{Value: big.NewInt(3), Data: []byte(" there")}
	outAcc := vmOutput.OutputAccounts[string(thirdPartyAddress)]
	outAcc.OutputTransfers = append(outAcc.OutputTransfers, outTransfer2)
	outAcc.BalanceDelta = big.NewInt(6)

	childAccount := arwen.AddNewOutputAccount(
		vmOutput,
		childAddress,
		0,
		nil,
	)
	childAccount.Balance = big.NewInt(1000)
	childAccount.BalanceDelta = big.NewInt(0).Sub(big.NewInt(1), big.NewInt(1))
	childAccount.GasUsed = 1296
	arwen.SetStorageUpdate(childAccount, childKey, childData)

	_ = arwen.AddNewOutputAccount(
		vmOutput,
		vaultAddress,
		4,
		[]byte{},
	)

	arwen.AddFinishData(vmOutput, []byte{3})
	arwen.AddFinishData(vmOutput, []byte("thirdparty"))
	arwen.AddFinishData(vmOutput, []byte("vault"))
	arwen.AddFinishData(vmOutput, []byte("user error"))
	arwen.AddFinishData(vmOutput, []byte("txhash"))

	vmOutput.ReturnMessage = "callBack error"

	return vmOutput
}

func expectedVMOutputCreateNewContractSuccess(_ []byte, childCode []byte) *vmcommon.VMOutput {
	vmOutput := arwen.MakeVMOutput()
	parentAccount := arwen.AddNewOutputAccount(
		vmOutput,
		parentAddress,
		-42,
		nil,
	)
	parentAccount.Balance = big.NewInt(1000)
	parentAccount.GasUsed = 885
	parentAccount.Nonce = 1
	arwen.SetStorageUpdate(parentAccount, []byte{'A'}, childCode)

	childAccount := arwen.AddNewOutputAccount(
		vmOutput,
		[]byte("newAddress"),
		42,
		nil,
	)
	childAccount.Code = childCode
	childAccount.GasUsed = 472
	childAccount.CodeMetadata = []byte{1, 0}
	childAccount.CodeDeployerAddress = parentAddress

	l := len(childCode)
	arwen.AddFinishData(vmOutput, []byte{byte(l / 256), byte(l % 256)})
	arwen.AddFinishData(vmOutput, []byte("init successful"))
	arwen.AddFinishData(vmOutput, []byte("succ"))

	return vmOutput
}

func expectedVMOutputCreateNewContractFail(_ []byte, childCode []byte) *vmcommon.VMOutput {
	vmOutput := arwen.MakeVMOutput()
	parentAccount := arwen.AddNewOutputAccount(
		vmOutput,
		parentAddress,
		0,
		nil,
	)
	parentAccount.Nonce = 0
	parentAccount.GasUsed = 2885
	arwen.SetStorageUpdate(parentAccount, []byte{'A'}, childCode)

	l := len(childCode)
	arwen.AddFinishData(vmOutput, []byte{byte(l / 256), byte(l % 256)})
	arwen.AddFinishData(vmOutput, []byte("fail"))

	return vmOutput
}

func expectedVMOutputMockedWasmerInstances() *vmcommon.VMOutput {
	vmOutput := arwen.MakeVMOutput()

	parentAccount := arwen.AddNewOutputAccount(
		vmOutput,
		parentAddress,
		0,
		nil,
	)
	parentAccount.Balance = big.NewInt(1000)
	parentAccount.BalanceDelta = big.NewInt(-4)
	parentAccount.GasUsed = 547
	arwen.SetStorageUpdate(parentAccount, []byte("parent"), []byte("parent storage"))

	childAccount := arwen.AddNewOutputAccount(
		vmOutput,
		childAddress,
		0,
		nil,
	)
	childAccount.BalanceDelta = big.NewInt(4)
	childAccount.GasUsed = 146
	arwen.SetStorageUpdate(childAccount, []byte("child"), []byte("child storage"))

	arwen.AddFinishData(vmOutput, []byte("parent returns this"))
	arwen.AddFinishData(vmOutput, []byte("child returns this"))

	return vmOutput
}
