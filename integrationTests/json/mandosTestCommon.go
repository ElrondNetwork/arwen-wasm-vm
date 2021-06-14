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
	_ = logger.SetLogLevel("*:DEBUG")
}

func getTestRoot() string {
	exePath, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	arwenTestRoot := filepath.Join(exePath, "../../test")
	return arwenTestRoot
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

func runSingleTest(t *testing.T, folder string, filename string) error {
	executor, err := am.NewArwenTestExecutor()
	require.Nil(t, err)
	runner := mc.NewScenarioRunner(
		executor,
		mc.NewDefaultFileResolver(),
	)

	fullPath := path.Join(getTestRoot(), folder)
	fullPath = path.Join(fullPath, filename)

	return runner.RunSingleJSONScenario(fullPath)
}
