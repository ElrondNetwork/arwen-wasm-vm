package arwenmandos

import (
	mc "github.com/ElrondNetwork/arwen-wasm-vm/mandos-go/controller"
	fr "github.com/ElrondNetwork/arwen-wasm-vm/mandos-go/fileresolver"
	mj "github.com/ElrondNetwork/arwen-wasm-vm/mandos-go/json/model"
	vmi "github.com/ElrondNetwork/elrond-go/core/vmcommon"
)

// Reset clears state/world.
// Is called in RunAllJSONScenariosInDirectory, but not in RunSingleJSONScenario.
func (ae *ArwenTestExecutor) Reset() {
	ae.World.Clear()
}

// ExecuteScenario executes an individual test.
func (ae *ArwenTestExecutor) ExecuteScenario(scenario *mj.Scenario, fileResolver fr.FileResolver) error {
	ae.fileResolver = fileResolver
	ae.checkGas = scenario.CheckGas
	err := ae.SetMandosGasSchedule(scenario.GasSchedule)
	if err != nil {
		return err
	}

	txIndex := 0
	for _, generalStep := range scenario.Steps {
		err := ae.ExecuteStep(generalStep)
		if err != nil {
			return err
		}

		txIndex++
	}

	return nil
}

// ExecuteStep executes an individual step from a scenario.
func (ae *ArwenTestExecutor) ExecuteStep(generalStep mj.Step) error {
	err := error(nil)

	switch step := generalStep.(type) {
	case *mj.ExternalStepsStep:
		err = ae.ExecuteExternalStep(step)
	case *mj.SetStateStep:
		err = ae.ExecuteSetStateStep(step)
	case *mj.CheckStateStep:
		err = ae.ExecuteCheckStateStep(step)
	case *mj.TxStep:
		_, err = ae.ExecuteTxStep(step)
	case *mj.DumpStateStep:
		err = ae.DumpWorld()
	}

	return err
}

// ExecuteExternalStep executes an external step referenced by the scenario.
func (ae *ArwenTestExecutor) ExecuteExternalStep(step *mj.ExternalStepsStep) error {
	log.Trace("ExternalStepsStep", "path", step.Path)
	if len(step.Comment) > 0 {
		log.Trace("ExternalStepsStep", "comment", step.Comment)
	}

	fileResolverBackup := ae.fileResolver
	clonedFileResolver := ae.fileResolver.Clone()
	externalStepsRunner := mc.NewScenarioRunner(ae, clonedFileResolver)

	extAbsPth := ae.fileResolver.ResolveAbsolutePath(step.Path)
	err := externalStepsRunner.RunSingleJSONScenario(extAbsPth)
	if err != nil {
		return err
	}

	ae.fileResolver = fileResolverBackup

	return nil
}

// ExecuteSetStateStep executes a SetStateStep.
func (ae *ArwenTestExecutor) ExecuteSetStateStep(step *mj.SetStateStep) error {
	if len(step.Comment) > 0 {
		log.Trace("SetStateStep", "comment", step.Comment)
	}

	// append accounts
	for _, mandosAccount := range step.Accounts {
		worldAccount, err := convertAccount(mandosAccount)
		if err != nil {
			return err
		}
		err = validateSetStateAccount(mandosAccount, worldAccount)
		if err != nil {
			return err
		}

		ae.World.AcctMap.PutAccount(worldAccount)
	}

	// replace block info
	ae.World.PreviousBlockInfo = convertBlockInfo(step.PreviousBlockInfo)
	ae.World.CurrentBlockInfo = convertBlockInfo(step.CurrentBlockInfo)
	ae.World.Blockhashes = mj.JSONBytesFromStringValues(step.BlockHashes)

	// append NewAddressMocks
	err := validateNewAddressMocks(step.NewAddressMocks)
	if err != nil {
		return err
	}
	addressMocksToAdd := convertNewAddressMocks(step.NewAddressMocks)
	ae.World.NewAddressMocks = append(ae.World.NewAddressMocks, addressMocksToAdd...)

	return nil
}

// ExecuteTxStep executes a TxStep.
func (ae *ArwenTestExecutor) ExecuteTxStep(step *mj.TxStep) (*vmi.VMOutput, error) {
	log.Trace("ExecuteTxStep", "id", step.TxIdent)
	if len(step.Comment) > 0 {
		log.Trace("ExecuteTxStep", "comment", step.Comment)
	}

	output, err := ae.executeTx(step.TxIdent, step.Tx)
	if err != nil {
		return nil, err
	}

	// check results
	if step.ExpectedResult != nil {
		err = checkTxResults(step.TxIdent, step.ExpectedResult, ae.checkGas, output)
		if err != nil {
			return nil, err
		}
	}

	return output, nil
}
