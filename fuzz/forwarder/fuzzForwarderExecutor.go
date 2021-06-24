package fuzzForwarder

import (
	"errors"
	"fmt"
	"io/ioutil"

	am "github.com/ElrondNetwork/arwen-wasm-vm/v1_3/arwenmandos"
	fr "github.com/ElrondNetwork/arwen-wasm-vm/v1_3/mandos-go/fileresolver"
	mj "github.com/ElrondNetwork/arwen-wasm-vm/v1_3/mandos-go/json/model"
	mjparse "github.com/ElrondNetwork/arwen-wasm-vm/v1_3/mandos-go/json/parse"
	mjwrite "github.com/ElrondNetwork/arwen-wasm-vm/v1_3/mandos-go/json/write"
	worldhook "github.com/ElrondNetwork/arwen-wasm-vm/v1_3/mock/world"
	vmi "github.com/ElrondNetwork/elrond-go/core/vmcommon"
)

type fuzzExecutor struct {
	arwenTestExecutor *am.ArwenTestExecutor
	world             *worldhook.MockWorld
	vm                vmi.VMExecutionHandler
	mandosParser      mjparse.Parser
	txIndex           int
	generatedScenario *mj.Scenario
	data              *fuzzData
}

func newFuzzExecutor(fileResolver fr.FileResolver) (*fuzzExecutor, error) {
	arwenTestExecutor, err := am.NewArwenTestExecutor()
	if err != nil {
		return nil, err
	}
	mandosGasSchedule := mj.GasScheduleV3
	arwenTestExecutor.SetMandosGasSchedule(mandosGasSchedule)

	parser := mjparse.NewParser(fileResolver)

	return &fuzzExecutor{
		arwenTestExecutor: arwenTestExecutor,
		world:             arwenTestExecutor.World,
		vm:                arwenTestExecutor.GetVM(),
		mandosParser:      parser,
		txIndex:           0,
		generatedScenario: &mj.Scenario{
			Name:        "fuzz generated",
			GasSchedule: mandosGasSchedule,
		},
		data: nil,
	}, nil
}

func (pfe *fuzzExecutor) executeStep(stepSnippet string) error {
	step, err := pfe.mandosParser.ParseScenarioStep(stepSnippet)
	if err != nil {
		return err
	}

	pfe.addStep(step)
	return pfe.arwenTestExecutor.ExecuteStep(step)
}

func (pfe *fuzzExecutor) addStep(step mj.Step) {
	pfe.generatedScenario.Steps = append(pfe.generatedScenario.Steps, step)
}

func (pfe *fuzzExecutor) saveGeneratedScenario() {
	serialized := mjwrite.ScenarioToJSONString(pfe.generatedScenario)

	err := ioutil.WriteFile("fuzz_gen.scen.json", []byte(serialized), 0644)
	if err != nil {
		fmt.Println(err)
	}
}

func (pfe *fuzzExecutor) executeTxStep(stepSnippet string) (*vmi.VMOutput, error) {
	step, err := pfe.mandosParser.ParseScenarioStep(stepSnippet)
	if err != nil {
		return nil, err
	}

	txStep, isTx := step.(*mj.TxStep)
	if !isTx {
		return nil, errors.New("tx step expected")
	}

	pfe.addStep(step)

	return pfe.arwenTestExecutor.ExecuteTxStep(txStep)
}

func (pfe *fuzzExecutor) log(info string, args ...interface{}) {
	fmt.Printf(info+"\n", args...)
}
