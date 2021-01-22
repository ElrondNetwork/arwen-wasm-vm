package main

import (
	"fmt"
	"os"

	"github.com/ElrondNetwork/arwen-wasm-vm/ipc/arwenpart"
	"github.com/ElrondNetwork/arwen-wasm-vm/ipc/common"
	"github.com/ElrondNetwork/arwen-wasm-vm/ipc/marshaling"
	logger "github.com/ElrondNetwork/elrond-go-logger"
	nodeConfig "github.com/ElrondNetwork/elrond-go/config"
	"github.com/ElrondNetwork/elrond-go/core"
	"github.com/ElrondNetwork/elrond-go/health"
)

const (
	fileDescriptorArwenInit      = 3
	fileDescriptorNodeToArwen    = 4
	fileDescriptorArwenToNode    = 5
	fileDescriptorReadLogProfile = 6
	fileDescriptorLogToNode      = 7
)

var appVersion = "undefined"

func main() {
	errCode, errMessage := doMain()
	if errCode != common.ErrCodeSuccess {
		fmt.Fprintln(os.Stderr, errMessage)
		os.Exit(errCode)
	}
}

// doMain returns (error code, error message)
func doMain() (int, string) {
	startHealthService()

	arwenInitFile := getPipeFile(fileDescriptorArwenInit)
	if arwenInitFile == nil {
		return common.ErrCodeCannotCreateFile, "Cannot get pipe file: [arwenInitFile]"
	}

	nodeToArwenFile := getPipeFile(fileDescriptorNodeToArwen)
	if nodeToArwenFile == nil {
		return common.ErrCodeCannotCreateFile, "Cannot get pipe file: [nodeToArwenFile]"
	}

	arwenToNodeFile := getPipeFile(fileDescriptorArwenToNode)
	if arwenToNodeFile == nil {
		return common.ErrCodeCannotCreateFile, "Cannot get pipe file: [arwenToNodeFile]"
	}

	arwenArguments, err := common.GetArwenArguments(arwenInitFile)
	if err != nil {
		return common.ErrCodeInit, fmt.Sprintf("Cannot receive gasSchedule: %v", err)
	}

	err = startLogging()
	if err != nil {
		return common.ErrCodeInit, fmt.Sprintf("Cannot initialize logging: %v", err)
	}

	messagesMarshalizer := marshaling.CreateMarshalizer(arwenArguments.MessagesMarshalizer)

	part, err := arwenpart.NewArwenPart(
		appVersion,
		nodeToArwenFile,
		arwenToNodeFile,
		&arwenArguments.VMHostParameters,
		messagesMarshalizer,
	)
	if err != nil {
		return common.ErrCodeInit, fmt.Sprintf("Cannot create ArwenPart: %v", err)
	}

	err = part.StartLoop()
	if err != nil {
		return common.ErrCodeTerminated, fmt.Sprintf("Ended Arwen loop: %v", err)
	}

	// This is never reached, actually. Arwen is supposed to run an infinite message loop.
	return common.ErrCodeSuccess, ""
}

func startHealthService() {
	healthService := health.NewHealthService(nodeConfig.HealthServiceConfig{
		IntervalVerifyMemoryInSeconds:             5,
		IntervalDiagnoseComponentsInSeconds:       60,
		IntervalDiagnoseComponentsDeeplyInSeconds: 60,
		MemoryUsageToCreateProfiles:               1 * 1024 * 1024,
		NumMemoryUsageRecordsToKeep:               100,
		FolderPath:                                "health-records",
	}, getWorkingDirectory())

	healthService.Start()
}

func getWorkingDirectory() string {
	workingDirectory := fmt.Sprintf("arwen_%d", os.Getpid())
	os.MkdirAll(workingDirectory, os.ModePerm)
	return workingDirectory
}

func getPipeFile(fileDescriptor uintptr) *os.File {
	file := os.NewFile(fileDescriptor, fmt.Sprintf("/proc/self/fd/%d", fileDescriptor))
	return file
}

func startLogging() error {
	logsFile, err := core.CreateFile(
		core.ArgCreateFileArgument{
			Prefix:        "logviewer",
			Directory:     getWorkingDirectory(),
			FileExtension: "log",
		},
	)
	if err != nil {
		return err
	}

	err = logger.AddLogObserver(logsFile, &logger.PlainFormatter{})
	if err != nil {
		return err
	}

	logger.SetLogLevel("*:TRACE")

	return nil
}
