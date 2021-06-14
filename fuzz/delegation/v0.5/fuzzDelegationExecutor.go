package delegation

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"math/big"
	"math/rand"
	"strings"
	"testing"

	am "github.com/ElrondNetwork/arwen-wasm-vm/arwenmandos"
	fr "github.com/ElrondNetwork/arwen-wasm-vm/mandos-go/fileresolver"
	mj "github.com/ElrondNetwork/arwen-wasm-vm/mandos-go/json/model"
	mjparse "github.com/ElrondNetwork/arwen-wasm-vm/mandos-go/json/parse"
	mjwrite "github.com/ElrondNetwork/arwen-wasm-vm/mandos-go/json/write"
	worldhook "github.com/ElrondNetwork/arwen-wasm-vm/mock/world"
	vmi "github.com/ElrondNetwork/elrond-go/core/vmcommon"
	"github.com/stretchr/testify/require"
)

const (
	UserWithdrawOnly    = "getUserWithdrawOnlyStake"
	UserWaiting         = "getUserWaitingStake"
	UserActive          = "getUserActiveStake"
	UserDeferredPayment = "getUserDeferredPaymentStake"
	UserUnbondable      = "getUnBondable"
)

type fuzzDelegationExecutorInitArgs struct {
	serviceFee                  int
	ownerMinStake               int
	minStake                    int
	numBlocksBeforeForceUnstake int
	numBlocksBeforeUnbond       int
	numDelegators               int
	stakePerNode                *big.Int
	numGenesisNodes             int
	totalDelegationCap          *big.Int
}

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
	ownerAddress                string
	delegationContractAddress   string
	auctionMockAddress          string
	faucetAddress               string
	withdrawTargetAddress       string
	stakePurchaseForwardAddress string
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
	mandosGasSchedule := mj.GasScheduleV3
	arwenTestExecutor.SetMandosGasSchedule(mandosGasSchedule)

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
			Name:        "fuzz generated",
			GasSchedule: mandosGasSchedule,
		},
	}, nil
}

func (pfe *fuzzDelegationExecutor) executeStep(stepSnippet string) error {
	step, err := pfe.mandosParser.ParseScenarioStep(stepSnippet)
	if err != nil {
		return err
	}

	pfe.addStep(step)
	return pfe.arwenTestExecutor.ExecuteStep(step)
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

func (pfe *fuzzDelegationExecutor) executeTxStep(stepSnippet string) (*vmi.VMOutput, error) {
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

func (pfe *fuzzDelegationExecutor) log(info string, args ...interface{}) {
	fmt.Printf(info+"\n", args...)
}

func (pfe *fuzzDelegationExecutor) addNodes(numNodesToAdd int) error {
	pfe.log("addNodes %d -> %d", numNodesToAdd, pfe.numNodes+numNodesToAdd)

	_, err := pfe.executeTxStep(fmt.Sprintf(`
	{
		"step": "scCall",
		"txId": "%d",
		"tx": {
			"from": "%s",
			"to": "%s",
			"value": "0",
			"function": "addNodes",
			"arguments": [
				%s
			],
			"gasLimit": "1,000,000,000",
			"gasPrice": "0"
		},
		"expect": {
			"out": [],
			"status": "",
			"logs": [],
			"gas": "*",
			"refund": "*"
		}
	}`,
		pfe.nextTxIndex(),
		pfe.ownerAddress,
		pfe.delegationContractAddress,
		blsKeySignatureArgsString(pfe.numNodes, numNodesToAdd),
	))
	pfe.numNodes += numNodesToAdd
	return err
}

func (pfe *fuzzDelegationExecutor) removeNodes(numNodesToRemove int) error {
	pfe.log("removeNodes %d -> %d", numNodesToRemove, pfe.numNodes-numNodesToRemove)

	output, err := pfe.executeTxStep(fmt.Sprintf(`
	{
		"step": "scCall",
		"txId": "%d",
		"tx": {
			"from": "%s",
			"to": "%s",
			"value": "0",
			"function": "removeNodes",
			"arguments": [
				%s
			],
			"gasLimit": "1,000,000,000",
			"gasPrice": "0"
		}
	}`,
		pfe.nextTxIndex(),
		pfe.ownerAddress,
		pfe.delegationContractAddress,
		blsKeysToBeRemoved(pfe.numNodes, numNodesToRemove),
	))
	if err != nil {
		return err
	}

	if output.ReturnCode != vmi.Ok {
		pfe.log("could not remove node because %s", output.ReturnMessage)
		return nil
	}

	return nil
}

func (pfe *fuzzDelegationExecutor) nextTxIndex() int {
	pfe.txIndex++
	return pfe.txIndex
}

func blsKeysToBeRemoved(totalNumNodes, numKeysToBeRemoved int) string {
	var blsKeys []string
	for i := 0; i < numKeysToBeRemoved; i++ {
		keyIndex := rand.Intn(totalNumNodes + 1)
		blsKeys = append(blsKeys, "\"str:"+blsKey(keyIndex)+"\"")
	}
	return strings.Join(blsKeys, ",")
}

func blsKeySignatureArgsString(startIndex, numNodes int) string {
	var blsKeyArgs []string
	for i := startIndex; i < startIndex+numNodes; i++ {
		blsKeyArgs = append(blsKeyArgs, "\"str:"+blsKey(i)+"\"")
		blsKeyArgs = append(blsKeyArgs, "\"str:"+blsSignature(i)+"\"")
	}
	return strings.Join(blsKeyArgs, ",")
}

func blsKey(index int) string {
	return fmt.Sprintf(
		"bls key %5d ..................................................................................",
		index)
}

func blsSignature(index int) string {
	return fmt.Sprintf(
		"bls key signature %5d ........................",
		index)
}

func (pfe *fuzzDelegationExecutor) getCurrentBlockNonce() uint64 {
	curentBlockNonce := uint64(0)
	if pfe.world.CurrentBlockInfo != nil {
		curentBlockNonce = pfe.world.CurrentBlockInfo.BlockNonce
	}
	return curentBlockNonce
}

func (pfe *fuzzDelegationExecutor) setBlockNonce(oldNonce, newNonce uint64) error {
	err := pfe.executeStep(fmt.Sprintf(`
	{
		"step": "setState",
		"comment": "%d - increase block nonce",
		"currentBlockInfo": {
			"blockNonce": "%d"
		}
	}`,
		pfe.nextTxIndex(),
		newNonce,
	))
	if err != nil {
		return err
	}

	pfe.log("block nonce: %d ---> %d", oldNonce, newNonce)
	return nil
}

func (pfe *fuzzDelegationExecutor) increaseBlockNonce(nonceDelta int) error {
	curentBlockNonce := pfe.getCurrentBlockNonce()
	return pfe.setBlockNonce(curentBlockNonce, curentBlockNonce+uint64(nonceDelta))
}

func (pfe *fuzzDelegationExecutor) querySingleResult(funcName string, args string) (*big.Int, error) {
	output, err := pfe.executeTxStep(fmt.Sprintf(`
	{
		"step": "scQuery",
		"txId": "%d",
		"tx": {
			"to": "%s",
			"function": "%s",
			"arguments": [
				%s
			]
		},
		"expect": {
			"out": [ "*" ],
			"status": ""
		}
	}`,
		pfe.nextTxIndex(),
		pfe.delegationContractAddress,
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

func (pfe *fuzzDelegationExecutor) delegatorQuery(funcName string, delegatorIndex int) (*big.Int, error) {
	delegatorAddr := fmt.Sprintf(`"str:%s"`, pfe.delegatorAddress(delegatorIndex))
	return pfe.querySingleResult(funcName, delegatorAddr)
}

func (pfe *fuzzDelegationExecutor) getAllDelegatorsBalance() *big.Int {
	totalDelegatorBalance := big.NewInt(0)
	for delegatorIdx := 0; delegatorIdx <= pfe.numDelegators; delegatorIdx++ {
		balance := pfe.getDelegatorBalance(delegatorIdx)
		totalDelegatorBalance.Add(totalDelegatorBalance, balance)
	}

	return totalDelegatorBalance
}

func (pfe *fuzzDelegationExecutor) getDelegatorBalance(delegatorIndex int) *big.Int {
	delegatorAddr := pfe.delegatorAddress(delegatorIndex)
	acct := pfe.world.AcctMap.GetAccount(pfe.interpretExpr(delegatorAddr))

	return acct.Balance
}

func (pfe *fuzzDelegationExecutor) modifyDelegationCap(newCap *big.Int) error {
	output, err := pfe.executeTxStep(fmt.Sprintf(`
	{
		"step": "scCall",
		"txId": "-modify-delegation-cap-",
		"tx": {
			"from": "%s",
			"to": "%s",
			"value": "0",
			"function": "modifyTotalDelegationCap",
			"arguments": ["%d"],
			"gasLimit": "100,000,000",
			"gasPrice": "0"
		}
	}`,
		pfe.ownerAddress,
		pfe.delegationContractAddress,
		newCap,
	))
	if err != nil {
		return err
	}

	pfe.log("modify delegation cap: returned code %s, returned message %s, newDelegationCap %d", output.ReturnCode, output.ReturnMessage, newCap)

	return nil
}

func (pfe *fuzzDelegationExecutor) setServiceFee(newServiceFee int) error {
	output, err := pfe.executeTxStep(fmt.Sprintf(`
	{
		"step": "scCall",
		"txId": "-set-service-fee-",
		"tx": {
			"from": "%s",
			"to": "%s",
			"value": "0",
			"function": "setServiceFee",
			"arguments": ["%d"],
			"gasLimit": "100,000,000",
			"gasPrice": "0"
		}
	}`,
		pfe.ownerAddress,
		pfe.delegationContractAddress,
		newServiceFee,
	))
	if err != nil {
		return err
	}

	pfe.log("modify service fee: returned code %s, newServiceFee %d", output.ReturnCode, newServiceFee)

	return nil
}

func (pfe *fuzzDelegationExecutor) continueGlobalOperation() error {
	completed := false
	for !completed {
		output, err := pfe.executeTxStep(fmt.Sprintf(`
		{
			"step": "scCall",
			"txId": "-continue-global-operation-",
			"tx": {
				"from": "%s",
				"to": "%s",
				"value": "0",
				"function": "continueGlobalOperation",
				"arguments": [],
				"gasLimit": "200,000,000",
				"gasPrice": "0"
			},
			"expect": {
				"out": [ "*" ],
				"refund": "*"
			}
		}`,
			pfe.ownerAddress,
			pfe.delegationContractAddress,
		))
		if err != nil {
			return err
		}
		pfe.log("continue global operation %s", string(output.ReturnData[0]))

		if bytes.Equal(output.ReturnData[0], []byte("completed")) {
			completed = true
		} else if bytes.Equal(output.ReturnData[0], []byte("interrupted")) {
			completed = false
		} else {
			return fmt.Errorf("unexpected global operation status: %x", output.ReturnData[0])
		}
	}

	return nil
}

func (pfe *fuzzDelegationExecutor) isBootstrapMode() (bool, error) {
	output, err := pfe.executeTxStep(fmt.Sprintf(`
	{
		"step": "scCall",
		"txId": "-is-bootstrap-mode-",
		"tx": {
			"from": "%s",
			"to": "%s",
			"value": "0",
			"function": "isBootstrapMode",
			"arguments": [],
			"gasLimit": "50,000,000",
			"gasPrice": "0"
		}
	}`,
		pfe.ownerAddress,
		pfe.delegationContractAddress,
	))
	if err != nil {
		return false, err
	}

	if bytes.Equal(output.ReturnData[0], []byte{1}) {
		return true, nil
	} else {
		return false, nil
	}
}

func (pfe *fuzzDelegationExecutor) printServiceFeeAndDelegationCap(t *testing.T) {
	_, err := pfe.querySingleResult("getTotalDelegationCap", "")
	require.Nil(t, err)

	_, err = pfe.querySingleResult("getServiceFee", "")
	require.Nil(t, err)
}
