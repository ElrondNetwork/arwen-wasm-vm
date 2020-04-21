package main

import (
	"os"
	"path/filepath"
	"testing"

	controller "github.com/ElrondNetwork/elrond-vm-util/test-util/testcontroller"
)

func TestErc20FromC(t *testing.T) {
	testExec := newArwenTestExecutor().replaceCode(
		"erc20.wasm",
		filepath.Join(getTestRoot(), "contracts/erc20-c.wasm"))

	err := controller.RunAllJSONTestsInDirectory(
		getTestRoot(),
		"erc20",
		".json",
		[]string{},
		testExec)

	if err != nil {
		t.Error(err)
	}
}

func TestErc20FromRust(t *testing.T) {
	testExec := newArwenTestExecutor().replaceCode(
		"erc20.wasm",
		filepath.Join(getTestRoot(), "contracts/simple-coin.wasm"))

	err := controller.RunAllJSONTestsInDirectory(
		getTestRoot(),
		"erc20",
		".json",
		[]string{},
		testExec)

	if err != nil {
		t.Error(err)
	}
}

func TestAdderFromRust(t *testing.T) {
	testExec := newArwenTestExecutor()

	err := controller.RunAllJSONTestsInDirectory(
		getTestRoot(),
		"adder",
		".json",
		[]string{},
		testExec)

	if err != nil {
		t.Error(err)
	}
}

func TestCryptoBubbles(t *testing.T) {
	testExec := newArwenTestExecutor()
	excludedTests := []string{}

	err := controller.RunAllJSONTestsInDirectory(
		getTestRoot(),
		"crypto_bubbles_min_v1",
		".json",
		excludedTests,
		testExec)

	if err != nil {
		t.Error(err)
	}
}

func TestFeaturesFromRust(t *testing.T) {
	testExec := newArwenTestExecutor()

	err := controller.RunAllJSONTestsInDirectory(
		getTestRoot(),
		"features",
		".json",
		[]string{},
		testExec)

	if err != nil {
		t.Error(err)
	}
}

func TestAsyncCalls(t *testing.T) {
	testExec := newArwenTestExecutor().replaceCode(
		"features.wasm",
		filepath.Join(getTestRoot(), "contracts/features.wasm"))

	err := controller.RunAllJSONTestsInDirectory(
		getTestRoot(),
		"async",
		".json",
		[]string{},
		testExec)

	if err != nil {
		t.Error(err)
	}
}

func TestDelegationContract(t *testing.T) {
	testExec := newArwenTestExecutor()
	err := controller.RunAllJSONTestsInDirectory(
		getTestRoot(),
		"delegation",
		".json",
		[]string{},
		testExec)

	if err != nil {
		t.Error(err)
	}
}

func TestDnsContract(t *testing.T) {
	testExec := newArwenTestExecutor()
	err := controller.RunAllJSONTestsInDirectory(
		getTestRoot(),
		"dns",
		".json",
		[]string{},
		testExec)

	if err != nil {
		t.Error(err)
	}
}

func getTestRoot() string {
	exePath, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	arwenTestRoot := filepath.Join(exePath, "../../test")
	return arwenTestRoot
}
