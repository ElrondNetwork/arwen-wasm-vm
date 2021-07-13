package elrond_ethereum_bridge

import (
	"flag"
	"math/big"
	"math/rand"
	"os"
	"path/filepath"
	"strconv"
	"testing"
	"time"

	fuzzutil "github.com/ElrondNetwork/arwen-wasm-vm/v1_3/fuzz/util"
	mc "github.com/ElrondNetwork/arwen-wasm-vm/v1_3/mandos-go/controller"
)

var fuzz = flag.Bool("fuzz", true, "Enable fuzz test")

var seedFlag = flag.Int64("seed", 0, "Random seed, use it to replay fuzz scenarios")

var iterationsFlag = flag.Int("iterations", 1000, "Number of iterations")

const EXECUTE_ACTION_PROBABILITY = 0.25

func getTestRoot() string {
	exePath, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	arwenTestRoot := filepath.Join(exePath, "../../test")
	return arwenTestRoot
}

func newExecutorWithPaths() *fuzzExecutor {
	fileResolver := mc.NewDefaultFileResolver().
		ReplacePath(
			"price-aggregator.wasm",
			filepath.Join(getTestRoot(), "elrond-ethereum-bridge/price-aggregator/price-aggregator.wasm")).
		ReplacePath(
			"multisig.wasm",
			filepath.Join(getTestRoot(), "elrond-ethereum-bridge/multisig/multisig.wasm")).
		ReplacePath(
			"egld-esdt-swap.wasm",
			filepath.Join(getTestRoot(), "elrond-ethereum-bridge/egld-esdt-swap/egld-esdt-swap.wasm")).
		ReplacePath(
			"esdt-safe.wasm",
			filepath.Join(getTestRoot(), "elrond-ethereum-bridge/esdt-safe/esdt-safe.wasm")).
		ReplacePath(
			"multi-transfer-esdt.wasm",
			filepath.Join(getTestRoot(), "elrond-ethereum-bridge/multi-transfer-esdt/multi-transfer-esdt.wasm")).
		ReplacePath(
			"ethereum-fee-prepay.wasm",
			filepath.Join(getTestRoot(), "elrond-ethereum-bridge/ethereum-fee-prepay/ethereum-fee-prepay.wasm"))

	fe, err := newFuzzExecutor(fileResolver)
	if err != nil {
		panic(err)
	}
	return fe
}

func TestElrondEthereumBridge(t *testing.T) {
	if !*fuzz {
		t.Skip("skipping test; only run with --fuzz argument")
	}

	fe := newExecutorWithPaths()
	defer fe.saveGeneratedScenario()

	err := fe.initData()
	if err != nil {
		t.Error(err)
	}

	// TODO: Uncomment once aggregator is integrated
	// The current version doesn't have relayer incentives and user fees

	/*
		err = fe.setupAggregator()
		if err != nil {
			t.Error(err)
		}
	*/

	nrRelayers := 2
	nrUsers := 2
	initialBalance := big.NewInt(INIT_BALANCE)
	err = fe.initAccounts(nrRelayers, nrUsers, initialBalance)
	if err != nil {
		t.Error(err)
	}

	multisigInitArgs := MultisigInitArgs{
		requiredStake: big.NewInt(1000),
		slashAmount:   big.NewInt(500),
		quorum:        len(fe.data.actorAddresses.relayers) / 2,
		boardMembers:  fe.data.actorAddresses.relayers,
	}
	err = fe.deployMultisig(&multisigInitArgs)
	if err != nil {
		t.Error(err)
	}

	deployChildContractsArgs := DeployChildContractsArgs{
		egldEsdtSwapCodePath:      "file:egld-esdt-swap.wasm",
		multiTransferEsdtCodePath: "file:multi-transfer-esdt.wasm",
		ethereumFeePrepayCodePath: "file:ethereum-fee-prepay.wasm",
		esdtSafeCodePath:          "file:esdt-safe.wasm",
		priceAggregatorAddress:    "sc:price-aggregator",
		wrappedEgldTokenId:        "str:WEGLD-123456",
		wrappedEthTokenId:         "str:WETH-abcdef",
		tokenWhitelist:            []string{},
	}
	err = fe.setupChildContracts(&deployChildContractsArgs)
	if err != nil {
		t.Error(err)
	}

	var seed int64
	if *seedFlag == 0 {
		seed = time.Now().UnixNano()
	} else {
		seed = *seedFlag
	}
	fe.log("Random seed: %d\n", seed)
	r := rand.New(rand.NewSource(seed))
	r.Seed(seed)

	fe.randSource = *r

	re := fuzzutil.NewRandomEventProvider(r)
	for stepIndex := 1; stepIndex <= *iterationsFlag; stepIndex++ {
		re.Reset()

		switch {
		case re.WithProbability(0.2):
			userAcc := fe.getRandomUser()
			wrapAmount := big.NewInt(int64(fe.randSource.Intn(100) + 1))

			err = fe.wrapEgld(userAcc, wrapAmount)
			if err != nil {
				t.Error(err)
			}
		case re.WithProbability(0.16):
			userAcc := fe.getRandomUser()
			userWrappedEgldBalance := fe.getEsdtBalance(userAcc, fe.data.wrappedEgldTokenId)

			// user has no wrapped eGLD
			if userWrappedEgldBalance.Cmp(big.NewInt(0)) == 0 {
				stepIndex--
				continue
			}

			unwrapAmount := fe.getRandomBigInt(userWrappedEgldBalance)
			scEgldBalance := fe.getBalance(fe.data.actorAddresses.egldEsdtSwap)

			// EgldEsdtSwap does not have enough funds
			expectedError := ""
			if unwrapAmount.Cmp(scEgldBalance) > 0 {
				expectedError = "Contract does not have enough funds"
			}

			err = fe.unwrapEgld(userAcc, unwrapAmount)
			if err != nil && err.Error() != expectedError {
				t.Error(err)
			}
		case re.WithProbability(0.16):
			userAcc := fe.getRandomUser()
			tokenId, amount, err := fe.generateValidRandomEsdtPayment(userAcc)

			// user has no ESDT to send
			if err != nil {
				stepIndex--
				continue
			}

			destAddress := fe.getEthAddress()
			err = fe.createEsdtSafeTransaction(userAcc, tokenId, amount, destAddress)
			if err != nil {
				t.Error(err)
			}
		case re.WithProbability(0.16):
			// must execute current transaction batch first, so this scCall would fail
			expectedError := ""
			if len(fe.data.multisigState.currentEsdtSafeTransactionBatch) > 0 {
				expectedError = "Must execute and set status for current tx batch first"
			}

			err := fe.getNextTransactionBatch()
			if err != nil && err.Error() != expectedError {
				t.Error(err)
			}
		case re.WithProbability(0.16):
			// must get batch first
			expectedError := ""
			if len(fe.data.multisigState.currentEsdtSafeTransactionBatch) == 0 {
				expectedError = "There is no transaction to set status for"
			}

			// generate random statuses for action
			statuses := []TransactionStatus{}
			for i := 0; i < len(fe.data.multisigState.currentEsdtSafeTransactionBatch); i++ {
				randNr := fe.randSource.Int31n(2)
				if randNr == 0 {
					statuses = append(statuses, Executed)
				} else {
					statuses = append(statuses, Rejected)
				}
			}

			expectedBalances, err := fe.GetExpectedBalancesAfterBridgeTransferToEthereum(
				fe.data.multisigState.currentEsdtSafeTransactionBatch,
				statuses,
			)
			if err != nil {
				t.Error(err)
			}

			wasActionProposed := false
			actionAlreadyProposedErrMsg := "Action already proposed"
			actionId, err := fe.proposeEsdtSafeSetCurrentTransactionBatchStatus(
				fe.getRandomRelayer(),
				fe.data.multisigState.currentEsdtSafeBatchId,
				statuses...,
			)

			// in this case, we need to get the action ID by query
			if err != nil && err.Error() == actionAlreadyProposedErrMsg {
				actionId, err = fe.getActionIdForSetCurrentTransactionBatchStatus(
					fe.data.multisigState.currentEsdtSafeBatchId,
					statuses...,
				)
				if err != nil {
					t.Error(err)

					continue
				}
				if actionId == 0 {
					t.Errorf("Action ID returned was 0 for EsdtSafe batchID %s",
						strconv.Itoa(fe.data.multisigState.currentEsdtSafeBatchId))

					continue
				}

				wasActionProposed = true

			} else if err != nil {
				// only signal error in test if the message is not the expected one
				if err.Error() != expectedError {
					t.Error(err)
				}

				continue
			}

			// To test multiple proposals, we only execute the proposed action with a probability
			// of 25%, or if the exact action was proposed already
			randNr := fe.randSource.Float32()
			if !(wasActionProposed || randNr < EXECUTE_ACTION_PROBABILITY) {
				continue
			}

			_, err = fe.performAction(fe.getRandomRelayer(), actionId)
			if err != nil {
				t.Error(err)
			}

			fe.data.multisigState.currentEsdtSafeBatchId = 0
			fe.data.multisigState.currentEsdtSafeTransactionBatch = []*Transaction{}

			for address := range expectedBalances {
				for tokenId := range expectedBalances[address] {
					expectedBalance := expectedBalances[address][tokenId]
					actualBalance := fe.getEsdtBalance(address, tokenId)

					if expectedBalance.Cmp(actualBalance) != 0 {
						t.Errorf("Expected and actual balances do not match. Address: %s. Expected %s. Actual: %s",
							address,
							expectedBalance.String(),
							actualBalance.String(),
						)
					}
				}
			}
		case re.WithProbability(0.16):
			nrTransfers := fe.randSource.Intn(10) + 1

			var transfers []*SimpleTransfer
			for i := 0; i < nrTransfers; i++ {
				transfers = append(transfers, fe.generateValidBridgedEsdtPayment())
			}

			var batchId int
			if fe.data.multisigState.currentEthereumBatchId != 0 {
				batchId = fe.data.multisigState.currentEthereumBatchId
			} else {
				batchId = fe.nextEthereumBatchId()
				fe.data.multisigState.currentEthereumBatchId = batchId
			}

			expectedBalances := fe.GetExpectedBalancesAfterBridgeTransferToElrond(transfers)

			wasActionProposed := false
			actionAlreadyProposedErrMsg := "This batch was already proposed"
			actionId, err := fe.proposeMultiTransferEsdtBatch(
				fe.getRandomRelayer(),
				batchId,
				transfers,
			)

			// in this case, we need to get the action ID by query
			if err != nil && err.Error() == actionAlreadyProposedErrMsg {
				actionId, err = fe.getActionIdForTransferBatch(
					fe.data.multisigState.currentEthereumBatchId,
					transfers,
				)
				if err != nil {
					t.Error(err)

					continue
				}
				if actionId == 0 {
					t.Errorf("Action ID returned was 0 for MultiTransferEsdt batchID %s",
						strconv.Itoa(fe.data.multisigState.currentEthereumBatchId))

					continue
				}

				wasActionProposed = true

			} else if err != nil {
				t.Error(err)

				continue
			}

			// To test multiple proposals, we only execute the proposed action with a probability
			// of 25%, or if the exact action was proposed already
			randNr := fe.randSource.Float32()
			if !(wasActionProposed || randNr < EXECUTE_ACTION_PROBABILITY) {
				continue
			}

			output, err := fe.performAction(fe.getRandomRelayer(), actionId)
			if err != nil {
				t.Error(err)
			}

			// output contains the status for each transfer
			for i := 0; i < nrTransfers; i++ {
				status := TransactionStatus(fe.bytesToInt(output[i]))

				switch status {
				case Executed:
					// no change needed
				case Rejected:
					// deduct balance from user for the specific transfer
					transfer := transfers[i]

					newUserBalance := big.NewInt(0)
					newUserBalance.Sub(expectedBalances[transfer.to][transfer.tokenId], transfer.amount)

					expectedBalances[transfer.to][transfer.tokenId] = newUserBalance
				default:
					t.Errorf("Invalid status parsed from output: %s", strconv.Itoa(int(status)))
				}
			}

			// check to see if expected and actual balances match
			for address := range expectedBalances {
				for tokenId := range expectedBalances[address] {
					expectedBalance := expectedBalances[address][tokenId]
					actualBalance := fe.getEsdtBalance(address, tokenId)

					if expectedBalance.Cmp(actualBalance) != 0 {
						t.Errorf("Expected and actual balances do not match. Address: %s. Expected %s. Actual: %s",
							address,
							expectedBalance.String(),
							actualBalance.String(),
						)
					}
				}
			}

			fe.data.multisigState.currentEthereumBatchId = 0

		default:
		}
	}
}
