package vmjsonintegrationtest

import (
	"testing"
)

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
	// TODO fix excluded tests and include them back
	runTestsInFolder(t, "features/async/mandos", []string{
		"features/async/mandos/recursive_caller_esdt_2.scen.json",
		"features/async/mandos/recursive_caller_esdt_x.scen.json",
	})
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
