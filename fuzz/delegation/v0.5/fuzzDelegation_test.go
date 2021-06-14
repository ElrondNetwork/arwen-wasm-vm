package delegation

import (
	"flag"
	"math/big"
	"math/rand"
	"os"
	"path/filepath"
	"testing"
	"time"

	fuzzutil "github.com/ElrondNetwork/arwen-wasm-vm/fuzz/util"
	mc "github.com/ElrondNetwork/arwen-wasm-vm/mandos-go/controller"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var fuzz = flag.Bool("fuzz", false, "Enable fuzz test")

var seedFlag = flag.Int64("seed", 0, "Random seed, use it to replay fuzz scenarios")

var iterationsFlag = flag.Int("iterations", 1000, "Number of iterations")

func getTestRoot() string {
	exePath, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	arwenTestRoot := filepath.Join(exePath, "../../../test")
	return arwenTestRoot
}

func newExecutorWithPaths() *fuzzDelegationExecutor {
	fileResolver := mc.NewDefaultFileResolver().
		ReplacePath(
			"delegation.wasm",
			filepath.Join(getTestRoot(), "delegation/v0_5_latest_full/output/delegation_latest_full.wasm")).
		ReplacePath(
			"auction-mock.wasm",
			filepath.Join(getTestRoot(), "delegation/auction-mock/output/auction-mock.wasm"))

	pfe, err := newFuzzDelegationExecutor(fileResolver)
	if err != nil {
		panic(err)
	}
	return pfe
}

func TestFuzzDelegation_v0_5(t *testing.T) {
	if !*fuzz {
		t.Skip("skipping test; only run with --fuzz argument")
	}

	pfe := newExecutorWithPaths()
	defer pfe.saveGeneratedScenario()

	var seed int64
	if *seedFlag == 0 {
		seed = time.Now().UnixNano()
	} else {
		seed = *seedFlag
	}
	pfe.log("Random seed: %d\n", seed)
	r := rand.New(rand.NewSource(seed))
	r.Seed(seed)

	stakePerNode := big.NewInt(1000000000)
	numDelegators := 10
	maxDelegationCap := big.NewInt(0).Mul(stakePerNode, big.NewInt(int64(4)))

	err := pfe.init(
		&fuzzDelegationExecutorInitArgs{
			serviceFee:                  r.Intn(10000),
			ownerMinStake:               0,
			minStake:                    r.Intn(1000000),
			numBlocksBeforeForceUnstake: r.Intn(1000),
			numBlocksBeforeUnbond:       r.Intn(1000),
			numDelegators:               numDelegators,
			stakePerNode:                stakePerNode,
			totalDelegationCap:          big.NewInt(0).Rand(r, maxDelegationCap),
		},
	)
	require.Nil(t, err)

	err = pfe.increaseBlockNonce(r.Intn(10000))
	require.Nil(t, err)

	re := fuzzutil.NewRandomEventProvider(r)
	for stepIndex := 0; stepIndex < *iterationsFlag; stepIndex++ {
		generateRandomEvent(t, pfe, r, re, maxDelegationCap)
	}

	err = pfe.increaseBlockNonce(r.Intn(pfe.numBlocksBeforeUnbond + 1))
	require.Nil(t, err)

	// all delegators (incl. owner) claim all rewards
	for delegatorIdx := 0; delegatorIdx <= pfe.numDelegators; delegatorIdx++ {
		err := pfe.claimRewards(delegatorIdx)
		require.Nil(t, err)
	}

	// check that delegators got all rewards out
	totalDelegatorBalance := pfe.getAllDelegatorsBalance()
	rewardsDifference := big.NewInt(0).Sub(pfe.totalRewards, totalDelegatorBalance)
	require.True(t, rewardsDifference.Cmp(big.NewInt(100)) == -1,
		"Rewards don't match. Total rewards: %d. Total delegator balance: %d.",
		pfe.totalRewards, totalDelegatorBalance,
	)

	pfe.printTotalStakeByType()

	// all delegators (incl. owner) unStake all from waiting
	for delegatorIdx := 0; delegatorIdx <= pfe.numDelegators; delegatorIdx++ {
		waitingFunds, err := pfe.getUserStakeOfType(delegatorIdx, UserWaiting)
		assert.Nil(t, err)

		err = pfe.unStake(delegatorIdx, waitingFunds)
		require.Nil(t, err)
	}

	pfe.printTotalStakeByType()

	err = pfe.increaseBlockNonce(pfe.numBlocksBeforeUnbond + 1)
	require.Nil(t, err)

	// all delegators (incl. owner) unBond
	for delegatorIdx := 0; delegatorIdx <= pfe.numDelegators; delegatorIdx++ {
		err = pfe.unBond(delegatorIdx)
		require.Nil(t, err)
	}

	pfe.printTotalStakeByType()
	for delegatorIdx := 0; delegatorIdx <= pfe.numDelegators; delegatorIdx++ {
		pfe.printUserStakeByType(delegatorIdx)
	}

	// auction SC should have no more funds
	auctionBalanceAfterUnbond := pfe.getAuctionBalance()
	require.True(t, auctionBalanceAfterUnbond.Sign() == 0,
		"Auction still has balance after full unbond: %d",
		auctionBalanceAfterUnbond)

	withdrawnAtTheEnd := pfe.getWithdrawTargetBalance()

	totalActiveStake := big.NewInt(0)
	for delegatorIdx := 0; delegatorIdx <= pfe.numDelegators; delegatorIdx++ {
		activeUserStake, err := pfe.getUserStakeOfType(delegatorIdx, UserActive)
		assert.Nil(t, err)

		totalActiveStake = totalActiveStake.Add(totalActiveStake, activeUserStake)
	}

	activeAndWithdrawn := big.NewInt(0).Add(withdrawnAtTheEnd, totalActiveStake)
	require.True(t, activeAndWithdrawn.Cmp(pfe.totalStakeAdded) == 0,
		"Stake added and withdrawn doesn't match. Staked: %d. Active+Withdrawn: %d. Off by: %d",
		pfe.totalStakeAdded, activeAndWithdrawn,
		big.NewInt(0).Sub(pfe.totalStakeAdded, activeAndWithdrawn))
}

func generateRandomEvent(
	t *testing.T,
	pfe *fuzzDelegationExecutor,
	r *rand.Rand,
	re *fuzzutil.RandomEventProvider,
	maxDelegationCap *big.Int,
) {
	maxStake := big.NewInt(0).Mul(pfe.stakePerNode, big.NewInt(2))
	maxSystemReward := big.NewInt(1000000000)
	maxServiceFee := 10000
	re.Reset()

	switch {
	case re.WithProbability(0.05):
		// increment block nonce
		err := pfe.increaseBlockNonce(r.Intn(1000))
		require.Nil(t, err)

		pfe.checkInvariants(t)
	case re.WithProbability(0.05):
		// add nodes
		err := pfe.addNodes(r.Intn(3))
		require.Nil(t, err)

		pfe.checkInvariants(t)
	case re.WithProbability(0.05):
		// add nodes
		err := pfe.removeNodes(r.Intn(2))
		require.Nil(t, err)
	case re.WithProbability(0.05):
		// stake
		delegatorIdx := r.Intn(pfe.numDelegators + 1)
		stake := big.NewInt(0).Rand(r, maxStake)

		err := pfe.stake(delegatorIdx, stake)
		require.Nil(t, err)

		pfe.checkInvariants(t)
	case re.WithProbability(0.05):
		ok, err := pfe.isBootstrapMode()
		require.Nil(t, err)

		if !ok {
			// add system rewards
			rewards := big.NewInt(0).Rand(r, maxSystemReward)

			err := pfe.addRewards(rewards)
			require.Nil(t, err)

			pfe.checkInvariants(t)
		}
	case re.WithProbability(0.2):
		// claim rewards
		delegatorIdx := r.Intn(pfe.numDelegators + 1)

		err := pfe.claimRewards(delegatorIdx)
		require.Nil(t, err)
	case re.WithProbability(0.05):
		// unStake
		delegatorIdx := r.Intn(pfe.numDelegators + 1)
		stake := big.NewInt(0).Rand(r, maxStake)

		err := pfe.unStake(delegatorIdx, stake)
		require.Nil(t, err)

		pfe.checkInvariants(t)
	case re.WithProbability(0.05):
		// unBond
		delegatorIdx := r.Intn(pfe.numDelegators + 1)
		err := pfe.unBond(delegatorIdx)
		require.Nil(t, err)
	case re.WithProbability(0.05):
		err := pfe.modifyDelegationCap(big.NewInt(0).Rand(r, maxDelegationCap))
		require.Nil(t, err)

		err = pfe.continueGlobalOperation()
		require.Nil(t, err)

		pfe.printServiceFeeAndDelegationCap(t)
		pfe.printTotalStakeByType()

		pfe.checkInvariants(t)
	case re.WithProbability(0.05):
		err := pfe.setServiceFee(r.Intn(maxServiceFee))
		require.Nil(t, err)

		err = pfe.continueGlobalOperation()
		require.Nil(t, err)

		pfe.printServiceFeeAndDelegationCap(t)
		pfe.printTotalStakeByType()

		pfe.checkInvariants(t)
	default:
	}
}

func (pfe *fuzzDelegationExecutor) checkInvariants(t *testing.T) {
	err := pfe.validateOwnerStakeShare()
	require.Nil(t, err)

	err = pfe.validateDelegationCapInvariant()
	require.Nil(t, err)
}
