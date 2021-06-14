package main

import (
	"fmt"
	"os"

	"github.com/ElrondNetwork/arwen-wasm-vm/ipc/arwenpart"
	"github.com/ElrondNetwork/arwen-wasm-vm/ipc/common"
	"github.com/ElrondNetwork/arwen-wasm-vm/ipc/marshaling"
	"github.com/ElrondNetwork/elrond-go-logger/pipes"
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

	readLogProfileFile := getPipeFile(fileDescriptorReadLogProfile)
	if readLogProfileFile == nil {
		return common.ErrCodeCannotCreateFile, "Cannot get pipe file: [readLogProfileFile]"
	}

	logToNodeFile := getPipeFile(fileDescriptorLogToNode)
	if logToNodeFile == nil {
		return common.ErrCodeCannotCreateFile, "Cannot get pipe file: [logToNodeFile]"
	}

	arwenArguments, err := common.GetArwenArguments(arwenInitFile)
	if err != nil {
		return common.ErrCodeInit, fmt.Sprintf("Cannot receive gasSchedule: %v", err)
	}

	messagesMarshalizer := marshaling.CreateMarshalizer(arwenArguments.MessagesMarshalizer)
	logsMarshalizer := marshaling.CreateMarshalizer(arwenArguments.LogsMarshalizer)

	logsPart, err := pipes.NewChildPart(readLogProfileFile, logToNodeFile, logsMarshalizer)
	if err != nil {
		return common.ErrCodeInit, fmt.Sprintf("Cannot create logs part: %v", err)
	}

	err = logsPart.StartLoop()
	if err != nil {
		return common.ErrCodeInit, fmt.Sprintf("Cannot start logs loop: %v", err)
	}

	defer logsPart.StopLoop()

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

func getPipeFile(fileDescriptor uintptr) *os.File {
	file := os.NewFile(fileDescriptor, fmt.Sprintf("/proc/self/fd/%d", fileDescriptor))
	return file
}
