package delegation

import (
	"errors"
	"fmt"
	"io/ioutil"
	"math/big"
	"strings"

	am "github.com/ElrondNetwork/arwen-wasm-vm/arwenmandos"
	fr "github.com/ElrondNetwork/arwen-wasm-vm/mandos-go/fileresolver"
	mj "github.com/ElrondNetwork/arwen-wasm-vm/mandos-go/json/model"
	mjparse "github.com/ElrondNetwork/arwen-wasm-vm/mandos-go/json/parse"
	mjwrite "github.com/ElrondNetwork/arwen-wasm-vm/mandos-go/json/write"
	worldhook "github.com/ElrondNetwork/arwen-wasm-vm/mock/world"
	vmi "github.com/ElrondNetwork/elrond-go/core/vmcommon"
)

type fuzzDelegationExecutor struct {
	arwenTestExecutor *am.ArwenTestExecutor
	world             *worldhook.MockWorld
	vm                vmi.VMExecutionHandler
	mandosParser      mjparse.Parser
	txIndex           int

	serviceFee                  int
	numBlocksBeforeForceUnstake int
	numBlocksBeforeUnbond       int
	numDelegators               int
	stakePerNode                *big.Int
	ownerAddress                []byte
	delegationContractAddress   []byte
	auctionMockAddress          []byte
	faucetAddress               []byte
	withdrawTargetAddress       []byte
	stakePurchaseForwardAddress []byte
	numNodes                    int
	totalStakeAdded             *big.Int
	totalStakeWithdrawn         *big.Int
	totalRewards                *big.Int
	generatedScenario           *mj.Scenario
}

func newFuzzDelegationExecutor(fileResolver fr.FileResolver) (*fuzzDelegationExecutor, error) {
	arwenTestExecutor, err := am.NewArwenTestExecutor()
	if err != nil {
		return nil, err
	}
	parser := mjparse.NewParser(fileResolver)
	return &fuzzDelegationExecutor{
		arwenTestExecutor:   arwenTestExecutor,
		world:               arwenTestExecutor.World,
		vm:                  arwenTestExecutor.GetVM(),
		mandosParser:        parser,
		txIndex:             0,
		numNodes:            0,
		totalStakeAdded:     big.NewInt(0),
		totalStakeWithdrawn: big.NewInt(0),
		totalRewards:        big.NewInt(0),
		generatedScenario: &mj.Scenario{
			Name: "fuzz generated",
		},
	}, nil
}

func (pfe *fuzzDelegationExecutor) log(info string, args ...interface{}) {
	fmt.Printf(info+"\n", args...)
}

type fuzzDelegationExecutorInitArgs struct {
	serviceFee                  int
	ownerMinStake               int
	numBlocksBeforeForceUnstake int
	numBlocksBeforeUnbond       int
	numDelegators               int
	stakePerNode                *big.Int
	numGenesisNodes             int
}

func (pfe *fuzzDelegationExecutor) addStep(step mj.Step) {
	pfe.generatedScenario.Steps = append(pfe.generatedScenario.Steps, step)
}

func (pfe *fuzzDelegationExecutor) saveGeneratedScenario() {
	serialized := mjwrite.ScenarioToJSONString(pfe.generatedScenario)

	err := ioutil.WriteFile("fuzz_gen.scen.json", []byte(serialized), 0644)
	if err != nil {
		fmt.Println(err)
	}
}

func (pfe *fuzzDelegationExecutor) nextTxIndex() int {
	pfe.txIndex++
	return pfe.txIndex
}

func (pfe *fuzzDelegationExecutor) getContractBalance() *big.Int {
	acct := pfe.world.AcctMap.GetAccount(pfe.delegationContractAddress)
	return acct.Balance
}

func (pfe *fuzzDelegationExecutor) getDelegatorBalance(delegIndex int) *big.Int {
	delegAddr := []byte(pfe.delegatorAddress(delegIndex))
	acct := pfe.world.AcctMap.GetAccount(delegAddr)
	return acct.Balance
}

func (pfe *fuzzDelegationExecutor) getAllDelegatorsBalance() *big.Int {
	totalDelegatorBalance := big.NewInt(0)
	for delegatorIdx := 0; delegatorIdx <= pfe.numDelegators; delegatorIdx++ {
		balance := pfe.getDelegatorBalance(delegatorIdx)
		totalDelegatorBalance.Add(totalDelegatorBalance, balance)
	}
	return totalDelegatorBalance
}

func (pfe *fuzzDelegationExecutor) getAuctionBalance() *big.Int {
	acct := pfe.world.AcctMap.GetAccount(pfe.auctionMockAddress)
	return acct.Balance
}

func (pfe *fuzzDelegationExecutor) getWithdrawTargetBalance() *big.Int {
	acct := pfe.world.AcctMap.GetAccount(pfe.withdrawTargetAddress)
	return acct.Balance
}

func (pfe *fuzzDelegationExecutor) executeStep(stepSnippet string) error {
	step, err := pfe.mandosParser.ParseScenarioStep(stepSnippet)
	if err != nil {
		return err
	}
	pfe.addStep(step)
	return pfe.arwenTestExecutor.ExecuteStep(step)
}

func (pfe *fuzzDelegationExecutor) executeTxStep(stepSnippet string) (*vmi.VMOutput, error) {
	step, err := pfe.mandosParser.ParseScenarioStep(stepSnippet)
	if err != nil {
		return nil, err
	}
	pfe.addStep(step)
	txStep, isTx := step.(*mj.TxStep)
	if !isTx {
		return nil, errors.New("tx step expected")
	}
	return pfe.arwenTestExecutor.ExecuteTxStep(txStep)
}

func (pfe *fuzzDelegationExecutor) querySingleResult(funcName string, args string) (*big.Int, error) {
	output, err := pfe.executeTxStep(fmt.Sprintf(`
	{
		"step": "scCall",
		"txId": "%d",
		"tx": {
			"from": "''%s",
			"to": "''%s",
			"value": "0",
			"function": "%s",
			"arguments": [
				%s
			],
			"gasLimit": "10,000,000",
			"gasPrice": "0"
		},
		"expect": {
			"out": [ "*" ],
			"status": "",
			"logs": [],
			"gas": "*",
			"refund": "*"
		}
	}`,
		pfe.nextTxIndex(),
		string(pfe.ownerAddress),
		string(pfe.delegationContractAddress),
		funcName,
		args,
	))
	if err != nil {
		return nil, err
	}

	result := big.NewInt(0).SetBytes(output.ReturnData[0])
	pfe.log("query: %s -> %d", funcName, result)
	return result, nil
}

func (pfe *fuzzDelegationExecutor) simpleQuery(funcName string) (*big.Int, error) {
	return pfe.querySingleResult(funcName, "")
}

func (pfe *fuzzDelegationExecutor) delegatorQuery(funcName string, delegIndex int) (*big.Int, error) {
	delegAddr := fmt.Sprintf(`"''%s"`, string(pfe.delegatorAddress(delegIndex)))
	return pfe.querySingleResult(funcName, delegAddr)
}

func (pfe *fuzzDelegationExecutor) increaseBlockNonce(nonceDelta int) error {
	curentBlockNonce := uint64(0)
	if pfe.world.CurrentBlockInfo != nil {
		curentBlockNonce = pfe.world.CurrentBlockInfo.BlockNonce
	}

	err := pfe.executeStep(fmt.Sprintf(`
	{
		"step": "setState",
		"comment": "%d - increase block nonce",
		"currentBlockInfo": {
			"blockNonce": "%d"
		}
	}`,
		pfe.nextTxIndex(),
		curentBlockNonce+uint64(nonceDelta),
	))
	if err != nil {
		return err
	}
	pfe.log("block nonce: %d ---> %d", curentBlockNonce, curentBlockNonce+uint64(nonceDelta))
	return nil
}

func (pfe *fuzzDelegationExecutor) dumpState() error {
	return pfe.executeStep(`
	{
		"step": "dumpState"
	}`)
}

func (pfe *fuzzDelegationExecutor) delegatorAddress(delegIndex int) []byte {
	if delegIndex == 0 {
		return pfe.ownerAddress
	}
	return []byte(fmt.Sprintf("delegator %5d               s1", delegIndex))
}

func blsKey(index int) string {
	return fmt.Sprintf(
		"bls key %5d ..................................................................................",
		index)
}

func blsSignature(index int) string {
	return fmt.Sprintf(
		"bls key signature %5d ........",
		index)
}

func blsKeySignatureArgsString(startIndex, numNodes int) string {
	var blsKeyArgs []string
	for i := startIndex; i < startIndex+numNodes; i++ {
		blsKeyArgs = append(blsKeyArgs, "\"''"+blsKey(i)+"\"")
		blsKeyArgs = append(blsKeyArgs, "\"''"+blsSignature(i)+"\"")
	}
	return strings.Join(blsKeyArgs, ",")
}

func blsKeyArgsString(numNodes int) string {
	var blsKeyArgs []string
	for i := 0; i < numNodes; i++ {
		blsKey := fmt.Sprintf(
			"bls key %5d ..................................................................................",
			i)
		blsKeyArg := "\"''" + blsKey + "\""
		blsKeyArgs = append(blsKeyArgs, blsKeyArg)
	}
	return strings.Join(blsKeyArgs, ",")
}

func blsSignatureArgsString(numNodes int) string {
	var blsSigArgs []string
	for i := 0; i < numNodes; i++ {
		blsSig := fmt.Sprintf(
			"bls key signature %5d ........",
			i)
		blsSigArg := "\"''" + blsSig + "\""
		blsSigArgs = append(blsSigArgs, blsSigArg)
	}
	return strings.Join(blsSigArgs, ",")
}
