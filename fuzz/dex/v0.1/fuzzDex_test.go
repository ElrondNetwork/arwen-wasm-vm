package dex

import (
	"flag"
	"math/rand"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	roulette "github.com/ElrondNetwork/arwen-wasm-vm/v1_3/fuzz/weightedroulette"
	mc "github.com/ElrondNetwork/arwen-wasm-vm/v1_3/mandos-go/controller"
	"github.com/stretchr/testify/require"
)

var fuzz = flag.Bool("fuzz", false, "Enable fuzz test")

var seedFlag = flag.Int64("seed", 0, "Random seed, use it to replay fuzz scenarios")

func newExecutorWithPaths() *fuzzDexExecutor {
	pwd, err := os.Getwd()
	if err != nil {
		return nil
	}

	fileResolver := mc.NewDefaultFileResolver().
		ReplacePath(
			"elrond_dex_router.wasm",
			filepath.Join(pwd, "wasms/elrond_dex_router.wasm")).
		ReplacePath(
			"elrond_dex_pair.wasm",
			filepath.Join(pwd, "wasms/elrond_dex_pair.wasm")).
		ReplacePath(
			"elrond_dex_farm.wasm",
			filepath.Join(pwd, "wasms/elrond_dex_farm.wasm"))

	pfe, err := newFuzzDexExecutor(fileResolver)
	if err != nil {
		panic(err)
	}
	return pfe
}

func TestFuzzDex_v0_1(t *testing.T) {
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

	err := pfe.init(
		&fuzzDexExecutorInitArgs{
			wegldTokenId:            "WEGLD-abcdef",
			mexTokenId:              "MEX-abcdef",
			busdTokenId:             "BUSD-abcdef",
			wemeLpTokenId:           "WEMELP-abcdef",
			webuLpTokenId:           "WEBULP-abcdef",
			wemeFarmTokenId:         "WEMEFARM-abcdef",
			webuFarmTokenId:         "WEBUFARM-abcdef",
			mexFarmTokenId:          "MEXFARM-abcdef",
			numUsers:                10,
			numEvents:               500,
			removeLiquidityProb:     5,
			addLiquidityProb:        20,
			swapProb:                25,
			queryPairsProb:          5,
			enterFarmProb:           18,
			exitFarmProb:            6,
			claimRewardsProb:        20,
			compoundRewardsProb:     10,
			removeLiquidityMaxValue: 1000000000,
			addLiquidityMaxValue:    1000000000,
			swapMaxValue:            10000000,
			enterFarmMaxValue:       100000000,
			exitFarmMaxValue:        100000000,
			claimRewardsMaxValue:    50000000,
			compoundRewardsMaxValue: 50000000,
			tokenDepositMaxValue:    50000000,
			blockNonceIncrease:      1,
		},
	)
	require.Nil(t, err)

	stats := eventsStatistics{
		swapFixedInputHits:         0,
		swapFixedInputMisses:       0,
		swapFixedOutputHits:        0,
		swapFixedOutputMisses:      0,
		addLiquidityHits:           0,
		addLiquidityMisses:         0,
		addLiquidityPriceChecks:    0,
		removeLiquidityHits:        0,
		removeLiquidityMisses:      0,
		removeLiquidityPriceChecks: 0,
		queryPairsHits:             0,
		queryPairsMisses:           0,
		enterFarmHits:              0,
		enterFarmMisses:            0,
		exitFarmHits:               0,
		exitFarmMisses:             0,
		exitFarmWithRewards:        0,
		claimRewardsHits:           0,
		claimRewardsMisses:         0,
		claimRewardsWithRewards:    0,
	}

	for stepIndex := 0; stepIndex < pfe.numEvents; stepIndex++ {
		generateRandomEvent(t, pfe, r, &stats)
		err := pfe.increaseBlockNonce(1)
		require.Nil(t, err)
	}

	printStatistics(&stats, pfe)
}

func generateRandomEvent(
	t *testing.T,
	pfe *fuzzDexExecutor,
	r *rand.Rand,
	statistics *eventsStatistics,
) {
	roulette.RandomChoice(
		r,
		roulette.Outcome{
			Weight: pfe.removeLiquidityProb,
			Event: func() {
				err := pfe.removeLiquidity(r, statistics)
				assert.Nil(t, err)
			},
		},
		roulette.Outcome{
			Weight: pfe.addLiquidityProb,
			Event: func() {
				err := pfe.addLiquidity(r, statistics)
				assert.Nil(t, err)
			},
		},
		roulette.Outcome{
			Weight: pfe.swapProb,
			Event: func() {
				err := pfe.swap(r, statistics)
				assert.Nil(t, err)
			},
		},
		roulette.Outcome{
			Weight: pfe.queryPairsProb,
			Event: func() {
				err := pfe.checkPairViews(r, statistics)
				assert.Nil(t, err)
			},
		},
		roulette.Outcome{
			Weight: pfe.enterFarmProb,
			Event: func() {
				err := pfe.enterFarm(r, statistics)
				assert.Nil(t, err)
			},
		},
		roulette.Outcome{
			Weight: pfe.exitFarmProb,
			Event: func() {
				err := pfe.exitFarm(r, statistics)
				assert.Nil(t, err)
			},
		},
		roulette.Outcome{
			Weight: pfe.claimRewardsProb,
			Event: func() {
				err := pfe.claimRewards(r, statistics)
				assert.Nil(t, err)
			},
		},
		roulette.Outcome{
			Weight: pfe.compoundRewardsProb,
			Event: func() {
				err := pfe.compoundRewards(r, statistics)
				assert.Nil(t, err)
			},
		},
	)
}

func printStatistics(statistics *eventsStatistics, pfe *fuzzDexExecutor) {
	pfe.log("\nStatistics:")
	pfe.log("\tswapFixedInputHits %d", statistics.swapFixedInputHits)
	pfe.log("\tswapFixedInputMisses %d", statistics.swapFixedInputMisses)
	pfe.log("")
	pfe.log("\tswapFixedOutputHits %d", statistics.swapFixedOutputHits)
	pfe.log("\tswapFixedOutputMissed %d", statistics.swapFixedOutputMisses)
	pfe.log("")
	pfe.log("\taddLiquidityHits %d", statistics.addLiquidityHits)
	pfe.log("\taddLiquidityMisses %d", statistics.addLiquidityMisses)
	pfe.log("\taddLiquidityPriceChecks %d", statistics.addLiquidityPriceChecks)
	pfe.log("")
	pfe.log("\tremoveLiquidityHits %d", statistics.removeLiquidityHits)
	pfe.log("\tremoveLiquidityMisses %d", statistics.removeLiquidityMisses)
	pfe.log("\tremoveLiquidityPriceChecks %d", statistics.removeLiquidityPriceChecks)
	pfe.log("")
	pfe.log("\tqueryPairHits %d", statistics.queryPairsHits)
	pfe.log("\tqueryPairMisses %d", statistics.queryPairsMisses)
	pfe.log("")
	pfe.log("\tenterFarmHits %d", statistics.enterFarmHits)
	pfe.log("\tenterFarmMisses %d", statistics.enterFarmMisses)
	pfe.log("")
	pfe.log("\texitFarmHits %d", statistics.exitFarmHits)
	pfe.log("\texitFarmMisses %d", statistics.exitFarmMisses)
	pfe.log("\texitFarmWithRewards %d", statistics.exitFarmWithRewards)
	pfe.log("")
	pfe.log("\tclaimRewardsHits %d", statistics.claimRewardsHits)
	pfe.log("\tclaimRewardsMisses %d", statistics.claimRewardsMisses)
	pfe.log("\tclaimRewardsWithRewards %d", statistics.claimRewardsWithRewards)
	pfe.log("")
	pfe.log("\tcompoundRewardsHits %d", statistics.compoundRewardsHits)
	pfe.log("\tcompoundRewardsMisses %d", statistics.compoundRewardsMisses)
	pfe.log("")
}
