package duration

import (
	"fmt"
	"math"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/ElrondNetwork/arwen-wasm-vm/arwen/crypto"
	"github.com/ElrondNetwork/arwen-wasm-vm/arwen/elrondapi"
	"github.com/ElrondNetwork/arwen-wasm-vm/arwen/ethapi"
	"github.com/ElrondNetwork/arwen-wasm-vm/arwen/host"
	"github.com/ElrondNetwork/arwen-wasm-vm/config"
	"github.com/ElrondNetwork/arwen-wasm-vm/wasmer"
	"github.com/stretchr/testify/require"
)

func TestCompilationTimeToSize(t *testing.T) {
	initializeWasmer()

	fmt.Println("Size(bytes),Name,Avg(ms),Min(ms),Max(ms)")

	for _, path := range getWasmFiles() {
		code := host.GetSCCode(path)
		min, max, avg := doBenchmark(t, path, code)
		fmt.Printf("%d,%s,%f,%f,%f\n", len(code), filepath.Base(path), avg, min, max)
	}
}

func doBenchmark(t *testing.T, name string, code []byte) (min float64, max float64, avg float64) {
	options := wasmer.CompilationOptions{
		GasLimit:           math.MaxUint32,
		OpcodeTrace:        false,
		Metering:           false,
		RuntimeBreakpoints: false,
	}

	durations := make([]float64, 0)
	for i := 0; i < 100; i++ {
		startIteration := time.Now()

		instance, err := wasmer.NewInstanceWithOptions(code, options)
		require.Nil(t, err)
		require.NotNil(t, instance)

		durations = append(durations, float64(time.Since(startIteration))/float64(time.Millisecond))
	}

	min, max, avg = analyseDurations(durations)
	return
}

func initializeWasmer() *wasmer.Imports {
	imports, _ := elrondapi.ElrondEIImports()
	imports, _ = elrondapi.BigIntImports(imports)
	imports, _ = ethapi.EthereumImports(imports)
	imports, _ = crypto.CryptoImports(imports)

	_ = wasmer.SetImports(imports)

	gasSchedule := config.MakeGasMapForTests()
	gasCostConfig, _ := config.CreateGasConfig(gasSchedule)
	opcodeCosts := gasCostConfig.WASMOpcodeCost.ToOpcodeCostsArray()
	wasmer.SetOpcodeCosts(&opcodeCosts)
	return imports
}

func getWasmFiles() []string {
	files, err := WalkMatch("../test", "*.wasm")
	if err != nil {
		panic("could not get wasm files")
	}

	toExclude := map[string]interface{}{
		"num-with-fp.wasm": nil,
		"promises.wasm":    nil,
		"train.wasm":       nil,
	}

	filtered := make([]string, 0)

	for _, file := range files {
		filename := filepath.Base(file)
		if _, ok := toExclude[filename]; !ok {
			filtered = append(filtered, file)
		}
	}

	return filtered
}

// https://stackoverflow.com/a/55300382/1475331
func WalkMatch(root, pattern string) ([]string, error) {
	var matches []string

	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		if matched, err := filepath.Match(pattern, filepath.Base(path)); err != nil {
			return err
		} else if matched {
			matches = append(matches, path)
		}
		return nil
	})

	if err != nil {
		return nil, err
	}
	return matches, nil
}

func analyseDurations(durations []float64) (min float64, max float64, avg float64) {
	min = float64(math.MaxUint32)

	sum := float64(0)
	for _, duration := range durations {
		if duration < min {
			min = duration
		}

		if duration > max {
			max = duration
		}

		sum += duration
	}

	avg = sum / float64(len(durations))
	return
}
