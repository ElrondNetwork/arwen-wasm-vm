package elrond_ethereum_bridge

import (
	"errors"
	"fmt"
	"io/ioutil"
	"math/rand"

	am "github.com/ElrondNetwork/arwen-wasm-vm/v1_3/arwenmandos"
	fr "github.com/ElrondNetwork/arwen-wasm-vm/v1_3/mandos-go/fileresolver"
	mj "github.com/ElrondNetwork/arwen-wasm-vm/v1_3/mandos-go/json/model"
	mjparse "github.com/ElrondNetwork/arwen-wasm-vm/v1_3/mandos-go/json/parse"
	mjwrite "github.com/ElrondNetwork/arwen-wasm-vm/v1_3/mandos-go/json/write"
	worldhook "github.com/ElrondNetwork/arwen-wasm-vm/v1_3/mock/world"
	vmi "github.com/ElrondNetwork/elrond-vm-common"
)

type fuzzData struct {
	actorAddresses     *ActorAddresses
	wrappedEgldTokenId string
	tokenWhitelist     []string
	multisigState      *MultisigState
}

type fuzzExecutor struct {
	arwenTestExecutor *am.ArwenTestExecutor
	world             *worldhook.MockWorld
	vm                vmi.VMExecutionHandler
	mandosParser      mjparse.Parser
	txIndex           int
	ethereumBatchId   int
	generatedScenario *mj.Scenario
	randSource        rand.Rand
	data              *fuzzData
}

func newFuzzExecutor(fileResolver fr.FileResolver) (*fuzzExecutor, error) {
	arwenTestExecutor, err := am.NewArwenTestExecutorWithFileResolver(&fileResolver)
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

func (fe *fuzzExecutor) executeStep(stepSnippet string) error {
	step, err := fe.mandosParser.ParseScenarioStep(stepSnippet)
	if err != nil {
		return err
	}

	fe.addStep(step)
	return fe.arwenTestExecutor.ExecuteStep(step)
}

func (fe *fuzzExecutor) addStep(step mj.Step) {
	fe.generatedScenario.Steps = append(fe.generatedScenario.Steps, step)
}

func (fe *fuzzExecutor) saveGeneratedScenario() {
	serialized := mjwrite.ScenarioToJSONString(fe.generatedScenario)

	err := ioutil.WriteFile("fuzz_gen.scen.json", []byte(serialized), 0644)
	if err != nil {
		fmt.Println(err)
	}
}

func (fe *fuzzExecutor) executeTxStep(stepSnippet string) (*vmi.VMOutput, error) {
	step, err := fe.mandosParser.ParseScenarioStep(stepSnippet)
	if err != nil {
		return nil, err
	}

	txStep, isTx := step.(*mj.TxStep)
	if !isTx {
		return nil, errors.New("tx step expected")
	}

	fe.addStep(step)

	return fe.arwenTestExecutor.ExecuteTxStep(txStep)
}

func (fe *fuzzExecutor) log(info string, args ...interface{}) {
	fmt.Printf(info+"\n", args...)
}
