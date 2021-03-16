package vmjsonintegrationtest

import (
	"os"
	"path"
	"path/filepath"
	"testing"

	am "github.com/ElrondNetwork/arwen-wasm-vm/arwenmandos"
	mc "github.com/ElrondNetwork/arwen-wasm-vm/mandos-go/controller"
	logger "github.com/ElrondNetwork/elrond-go-logger"
	"github.com/stretchr/testify/require"
)

func init() {
	_ = logger.SetLogLevel("*:NONE")
}

func getTestRoot() string {
	exePath, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	arwenTestRoot := filepath.Join(exePath, "../../test")
	return arwenTestRoot
}

func TestRustErc20(t *testing.T) {
	runAllTestsInFolder(t, "erc20-rust/mandos")
}

func TestCErc20(t *testing.T) {
	runAllTestsInFolder(t, "erc20-c")
}

func TestRustAdder(t *testing.T) {
	runAllTestsInFolder(t, "adder/mandos")
}

func TestMultisig(t *testing.T) {
	runAllTestsInFolder(t, "multisig/mandos")
}

func TestRustBasicFeaturesLatest(t *testing.T) {
	if testing.Short() {
		t.Skip("not a short test")
	}

	runAllTestsInFolder(t, "features/basic-features/mandos")
}

func TestRustBasicFeaturesNoSmallIntApi(t *testing.T) {
	if testing.Short() {
		t.Skip("not a short test")
	}

	runAllTestsInFolder(t, "features/basic-features-no-small-int-api/mandos")
}

// Backwards compatibility.
func TestRustBasicFeaturesLegacy(t *testing.T) {
	if testing.Short() {
		t.Skip("not a short test")
	}

	runAllTestsInFolder(t, "features/basic-features-legacy/mandos")
}

func TestRustPayableFeaturesLatest(t *testing.T) {
	if testing.Short() {
		t.Skip("not a short test")
	}

	runAllTestsInFolder(t, "features/payable-features/mandos")
}

func TestRustAsyncCalls(t *testing.T) {
	runTestsInFolder(t, "features/async/mandos", []string{
		"features/async/mandos/forwarder_send_twice_esdt.scen.json",
		"features/async/mandos/recursive_caller_esdt_2.scen.json",
		"features/async/mandos/recursive_caller_esdt_x.scen.json",
	})
}

func TestDelegation_v0_2(t *testing.T) {
	if testing.Short() {
		t.Skip("not a short test")
	}

	runAllTestsInFolder(t, "delegation/v0_2")
}

func TestDelegation_v0_3(t *testing.T) {
	if testing.Short() {
		t.Skip("not a short test")
	}

	runTestsInFolder(t, "delegation/v0_3", []string{
		"delegation/v0_3/test/integration/genesis/genesis.scen.json",
	})
}

func TestDelegation_v0_4_genesis(t *testing.T) {
	if testing.Short() {
		t.Skip("not a short test")
	}

	runAllTestsInFolder(t, "delegation/v0_4_genesis")
}

func TestDelegation_v0_5_2_full(t *testing.T) {
	if testing.Short() {
		t.Skip("not a short test")
	}

	runAllTestsInFolder(t, "delegation/v0_5_2_full")
}

func TestDelegation_v0_5_2_update(t *testing.T) {
	if testing.Short() {
		t.Skip("not a short test")
	}

	runAllTestsInFolder(t, "delegation/v0_5_2_update")
}

func TestDnsContract(t *testing.T) {
	if testing.Short() {
		t.Skip("not a short test")
	}

	runAllTestsInFolder(t, "dns")
}

func TestTimelocks(t *testing.T) {
	runAllTestsInFolder(t, "timelocks")
}

// func TestPromises(t *testing.T) {
// 	executor, err := am.NewArwenTestExecutor()
// 	require.Nil(t, err)
// 	runner := mc.NewScenarioRunner(
// 		executor,
// 		mc.NewDefaultFileResolver(),
// 	)
// 	err = runner.RunAllJSONScenariosInDirectory(
// 		getTestRoot(),
// 		"promises",
// 		".scen.json",
// 		[]string{})

// 	if err != nil {
// 		t.Error(err)
// 	}
// }

func TestCrowdfundingEsdt(t *testing.T) {
	runAllTestsInFolder(t, "crowdfunding-esdt")
}

func TestEgldEsdtSwap(t *testing.T) {
	runAllTestsInFolder(t, "egld-esdt-swap")
}

func TestPingPongEgld(t *testing.T) {
	runAllTestsInFolder(t, "ping-pong-egld")
}

func runAllTestsInFolder(t *testing.T, folder string) {
	runTestsInFolder(t, folder, []string{})
}

func runTestsInFolder(t *testing.T, folder string, exclusions []string) {
	executor, err := am.NewArwenTestExecutor()
	require.Nil(t, err)
	runner := mc.NewScenarioRunner(
		executor,
		mc.NewDefaultFileResolver(),
	)

	err = runner.RunAllJSONScenariosInDirectory(
		getTestRoot(),
		folder,
		".scen.json",
		exclusions)

	if err != nil {
		t.Error(err)
	}
}

func runSingleTest(t *testing.T, folder string, filename string) {
	executor, err := am.NewArwenTestExecutor()
	require.Nil(t, err)
	runner := mc.NewScenarioRunner(
		executor,
		mc.NewDefaultFileResolver(),
	)

	fullPath := path.Join(getTestRoot(), folder)
	fullPath = path.Join(fullPath, filename)

	err = runner.RunSingleJSONScenario(fullPath)
	if err != nil {
		t.Error(err)
	}
}
