package fuzzForwarder

import (
	"flag"
	"math/rand"
	"os"
	"path/filepath"
	"testing"
	"time"

	mc "github.com/ElrondNetwork/arwen-wasm-vm/v1_3/mandos-go/controller"
	"github.com/stretchr/testify/require"
)

var fuzz = flag.Bool("fuzz", false, "Enable fuzz test")

var seedFlag = flag.Int64("seed", 0, "Random seed, use it to replay fuzz scenarios")

var iterationsFlag = flag.Int("iterations", 250, "Number of iterations")

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
			"../forwarder/output/forwarder.wasm",
			filepath.Join(getTestRoot(), "features/composability/forwarder/output/forwarder.wasm"))

	pfe, err := newFuzzExecutor(fileResolver)
	if err != nil {
		panic(err)
	}
	return pfe
}

func TestFuzzForwarder(t *testing.T) {

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

	err := pfe.initData()
	require.Nil(t, err)

	err = pfe.setUp()
	require.Nil(t, err)

	for stepIndex := 0; stepIndex < *iterationsFlag; stepIndex++ {
		callType := pfe.randomCallType(r)
		tokenName, nonce := pfe.randomTokenNameAndNonce(r)
		fromIndex := r.Intn(pfe.data.numForwarders) + 1
		toIndex := r.Intn(pfe.data.numForwarders) + 1
		if nonce > 0 {
			for toIndex == fromIndex {
				toIndex = r.Intn(pfe.data.numForwarders) + 1
			}
		}
		err = pfe.programCall(callType, fromIndex, toIndex, tokenName, nonce, "10")
		require.Nil(t, err)
	}

	err = pfe.executeCallCheckLogs(1)
	require.Nil(t, err)
}
