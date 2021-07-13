package host

import (
	"fmt"
	"sync"

	"github.com/ElrondNetwork/arwen-wasm-vm/v1_3/arwen"
	"github.com/ElrondNetwork/arwen-wasm-vm/v1_3/arwen/contexts"
	"github.com/ElrondNetwork/arwen-wasm-vm/v1_3/arwen/cryptoapi"
	"github.com/ElrondNetwork/arwen-wasm-vm/v1_3/arwen/elrondapi"
	"github.com/ElrondNetwork/arwen-wasm-vm/v1_3/config"
	"github.com/ElrondNetwork/arwen-wasm-vm/v1_3/crypto"
	"github.com/ElrondNetwork/arwen-wasm-vm/v1_3/crypto/factory"
	"github.com/ElrondNetwork/arwen-wasm-vm/v1_3/wasmer"
	logger "github.com/ElrondNetwork/elrond-go-logger"
	vmcommon "github.com/ElrondNetwork/elrond-vm-common"
	"github.com/ElrondNetwork/elrond-vm-common/atomic"
	"github.com/ElrondNetwork/elrond-vm-common/parsers"
)

var log = logger.GetOrCreate("arwen/host")

// MaximumWasmerInstanceCount specifies the maximum number of allowed Wasmer
// instances on the InstanceStack of the RuntimeContext
var MaximumWasmerInstanceCount = uint64(10)

// TryFunction corresponds to the try() part of a try / catch block
type TryFunction func()

// CatchFunction corresponds to the catch() part of a try / catch block
type CatchFunction func(error)

// vmHost implements HostContext interface.
type vmHost struct {
	cryptoHook   crypto.VMCrypto
	mutExecution sync.RWMutex

	ethInput []byte

	blockchainContext arwen.BlockchainContext
	runtimeContext    arwen.RuntimeContext
	asyncContext      arwen.AsyncContext
	outputContext     arwen.OutputContext
	meteringContext   arwen.MeteringContext
	storageContext    arwen.StorageContext
	bigIntContext     arwen.BigIntContext

	gasSchedule              config.GasScheduleMap
	scAPIMethods             *wasmer.Imports
	protocolBuiltinFunctions vmcommon.FunctionNames

	arwenV2EnableEpoch uint32
	flagArwenV2        atomic.Flag

	aotEnableEpoch  uint32
	flagAheadOfTime atomic.Flag

	dynGasLockEnableEpoch uint32
	flagDynGasLock        atomic.Flag

	arwenV3EnableEpoch uint32
	flagArwenV3        atomic.Flag

	eSDTFunctionsEnableEpoch uint32
	flagESDTFunctions        atomic.Flag

	callArgsParser arwen.CallArgsParser
}

// NewArwenVM creates a new Arwen vmHost
func NewArwenVM(
	blockChainHook vmcommon.BlockchainHook,
	hostParameters *arwen.VMHostParameters,
) (arwen.VMHost, error) {

	cryptoHook := factory.NewVMCrypto()
	host := &vmHost{
		cryptoHook:               cryptoHook,
		meteringContext:          nil,
		runtimeContext:           nil,
		asyncContext:             nil,
		blockchainContext:        nil,
		storageContext:           nil,
		bigIntContext:            nil,
		gasSchedule:              hostParameters.GasSchedule,
		scAPIMethods:             nil,
		protocolBuiltinFunctions: hostParameters.ProtocolBuiltinFunctions,
		arwenV2EnableEpoch:       hostParameters.ArwenV2EnableEpoch,
		aotEnableEpoch:           hostParameters.AheadOfTimeEnableEpoch,
		arwenV3EnableEpoch:       hostParameters.ArwenV3EnableEpoch,
		dynGasLockEnableEpoch:    hostParameters.DynGasLockEnableEpoch,
		eSDTFunctionsEnableEpoch: hostParameters.ArwenESDTFunctionsEnableEpoch,
		callArgsParser:           parsers.NewCallArgsParser(),
	}

	var err error

	imports, err := elrondapi.ElrondEIImports()
	if err != nil {
		return nil, err
	}

	imports, err = elrondapi.BigIntImports(imports)
	if err != nil {
		return nil, err
	}

	imports, err = elrondapi.SmallIntImports(imports)
	if err != nil {
		return nil, err
	}

	imports, err = cryptoapi.CryptoImports(imports)
	if err != nil {
		return nil, err
	}

	err = wasmer.SetImports(imports)
	if err != nil {
		return nil, err
	}

	host.scAPIMethods = imports

	host.blockchainContext, err = contexts.NewBlockchainContext(host, blockChainHook)
	if err != nil {
		return nil, err
	}

	host.runtimeContext, err = contexts.NewRuntimeContext(
		host,
		hostParameters.VMType,
		hostParameters.UseWarmInstance,
	)
	if err != nil {
		return nil, err
	}

	host.asyncContext = contexts.NewAsyncContext(host)

	host.meteringContext, err = contexts.NewMeteringContext(host, hostParameters.GasSchedule, hostParameters.BlockGasLimit)
	if err != nil {
		return nil, err
	}

	host.outputContext, err = contexts.NewOutputContext(host)
	if err != nil {
		return nil, err
	}

	host.storageContext, err = contexts.NewStorageContext(host, blockChainHook, hostParameters.ElrondProtectedKeyPrefix)
	if err != nil {
		return nil, err
	}

	host.bigIntContext, err = contexts.NewBigIntContext()
	if err != nil {
		return nil, err
	}

	gasCostConfig, err := config.CreateGasConfig(host.gasSchedule)
	if err != nil {
		return nil, err
	}

	host.runtimeContext.SetMaxInstanceCount(MaximumWasmerInstanceCount)

	opcodeCosts := gasCostConfig.WASMOpcodeCost.ToOpcodeCostsArray()
	wasmer.SetOpcodeCosts(&opcodeCosts)

	host.initContexts()

	return host, nil
}

// GetVersion returns the Arwen version string
func (host *vmHost) GetVersion() string {
	return arwen.ArwenVersion
}

// Crypto returns the VMCrypto instance of the host
func (host *vmHost) Crypto() crypto.VMCrypto {
	return host.cryptoHook
}

// Blockchain returns the BlockchainContext instance of the host
func (host *vmHost) Blockchain() arwen.BlockchainContext {
	return host.blockchainContext
}

// Runtime returns the RuntimeContext instance of the host
func (host *vmHost) Runtime() arwen.RuntimeContext {
	return host.runtimeContext
}

// Output returns the OutputContext instance of the host
func (host *vmHost) Output() arwen.OutputContext {
	return host.outputContext
}

// Metering returns the MeteringContext instance of the host
func (host *vmHost) Metering() arwen.MeteringContext {
	return host.meteringContext
}

// Storage returns the StorageContext instance of the host
func (host *vmHost) Storage() arwen.StorageContext {
	return host.storageContext
}

// BigInt returns the BigIntContext instance of the host
func (host *vmHost) BigInt() arwen.BigIntContext {
	return host.bigIntContext
}

// IsArwenV2Enabled returns whether the Arwen V2 mode is enabled
func (host *vmHost) IsArwenV2Enabled() bool {
	return host.flagArwenV2.IsSet()
}

// IsArwenV3Enabled returns whether the V3 features are enabled
func (host *vmHost) IsArwenV3Enabled() bool {
	return host.flagArwenV3.IsSet()
}

// IsAheadOfTimeCompileEnabled returns whether ahead-of-time compilation is enabled
func (host *vmHost) IsAheadOfTimeCompileEnabled() bool {
	return host.flagAheadOfTime.IsSet()
}

// IsDynamicGasLockingEnabled returns whether dynamic gas locking mode is enabled
func (host *vmHost) IsDynamicGasLockingEnabled() bool {
	return host.flagDynGasLock.IsSet()
}

// IsESDTFunctionsEnabled returns whether ESDT functions are enabled
func (host *vmHost) IsESDTFunctionsEnabled() bool {
	return host.flagESDTFunctions.IsSet()
}

// GetContexts returns the main contexts of the host
func (host *vmHost) GetContexts() (
	arwen.BigIntContext,
	arwen.BlockchainContext,
	arwen.MeteringContext,
	arwen.OutputContext,
	arwen.RuntimeContext,
	arwen.AsyncContext,
	arwen.StorageContext,
) {
	return host.bigIntContext,
		host.blockchainContext,
		host.meteringContext,
		host.outputContext,
		host.runtimeContext,
		host.asyncContext,
		host.storageContext
}

// InitState resets the contexts of the host and reconfigures its flags
func (host *vmHost) InitState() {
	host.initContexts()
	currentEpoch := host.Blockchain().CurrentEpoch()
	host.flagArwenV2.Toggle(currentEpoch >= host.arwenV2EnableEpoch)
	log.Trace("arwenV2", "enabled", host.flagArwenV2.IsSet())

	host.flagAheadOfTime.Toggle(currentEpoch >= host.aotEnableEpoch)
	log.Trace("aheadOfTime compile", "enabled", host.flagAheadOfTime.IsSet())

	host.flagDynGasLock.Toggle(currentEpoch >= host.dynGasLockEnableEpoch)
	log.Trace("dynamic gas locking", "enabled", host.flagDynGasLock.IsSet())

	host.flagArwenV3.Toggle(currentEpoch >= host.arwenV3EnableEpoch)
	log.Trace("arwen v3 improvement", "enabled", host.flagArwenV3.IsSet())

	host.flagESDTFunctions.Toggle(currentEpoch >= host.eSDTFunctionsEnableEpoch)
	log.Trace("esdt functions", "enabled", host.flagESDTFunctions.IsSet())
}

func (host *vmHost) initContexts() {
	host.ClearContextStateStack()
	host.bigIntContext.InitState()
	host.outputContext.InitState()
	host.meteringContext.InitState()
	host.runtimeContext.InitState()
	host.asyncContext.InitState()
	host.storageContext.InitState()
	host.blockchainContext.InitState()
	host.ethInput = nil
}

// ClearContextStateStack cleans the state stacks of all the contexts of the host
func (host *vmHost) ClearContextStateStack() {
	host.bigIntContext.ClearStateStack()
	host.outputContext.ClearStateStack()
	host.meteringContext.ClearStateStack()
	host.runtimeContext.ClearStateStack()
	host.asyncContext.ClearStateStack()
	host.storageContext.ClearStateStack()
	host.blockchainContext.ClearStateStack()
}

// Clean closes the currently running Wasmer instance
func (host *vmHost) Clean() {
	if host.runtimeContext.IsWarmInstance() {
		return
	}
	host.runtimeContext.CleanWasmerInstance()
}

// GetAPIMethods returns the EEI as a set of imports for Wasmer
func (host *vmHost) GetAPIMethods() *wasmer.Imports {
	return host.scAPIMethods
}

// GetProtocolBuiltinFunctions returns the names of the built-in functions, reserved by the protocol
func (host *vmHost) GetProtocolBuiltinFunctions() vmcommon.FunctionNames {
	return host.protocolBuiltinFunctions
}

// SetProtocolBuiltinFunctions sets the names of build-in functions, reserved by the protocol
func (host *vmHost) SetProtocolBuiltinFunctions(functionNames vmcommon.FunctionNames) {
	host.protocolBuiltinFunctions = functionNames
}

// GasScheduleChange applies a new gas schedule to the host
func (host *vmHost) GasScheduleChange(newGasSchedule config.GasScheduleMap) {
	host.mutExecution.Lock()
	defer host.mutExecution.Unlock()

	host.gasSchedule = newGasSchedule
	gasCostConfig, err := config.CreateGasConfig(newGasSchedule)
	if err != nil {
		log.Error("cannot apply new gas config remained with old one")
		return
	}

	opcodeCosts := gasCostConfig.WASMOpcodeCost.ToOpcodeCostsArray()
	wasmer.SetOpcodeCosts(&opcodeCosts)

	host.meteringContext.SetGasSchedule(newGasSchedule)
}

// GetGasScheduleMap returns the currently stored gas schedule
func (host *vmHost) GetGasScheduleMap() config.GasScheduleMap {
	return host.gasSchedule
}

// RunSmartContractCreate executes the deployment of a new contract
func (host *vmHost) RunSmartContractCreate(input *vmcommon.ContractCreateInput) (vmOutput *vmcommon.VMOutput, err error) {
	host.mutExecution.RLock()
	defer host.mutExecution.RUnlock()

	log.Trace("RunSmartContractCreate begin", "len(code)", len(input.ContractCode), "metadata", input.ContractCodeMetadata)

	try := func() {
		vmOutput = host.doRunSmartContractCreate(input)
	}

	catch := func(caught error) {
		err = caught
		log.Error("RunSmartContractCreate", "error", err)
	}

	TryCatch(try, catch, "arwen.RunSmartContractCreate")
	if vmOutput != nil {
		log.Trace("RunSmartContractCreate end", "returnCode", vmOutput.ReturnCode, "returnMessage", vmOutput.ReturnMessage)
	}

	return
}

// RunSmartContractCall executes the call of an existing contract
func (host *vmHost) RunSmartContractCall(input *vmcommon.ContractCallInput) (vmOutput *vmcommon.VMOutput, err error) {
	host.mutExecution.RLock()
	defer host.mutExecution.RUnlock()

	log.Trace("RunSmartContractCall begin", "function", input.Function)

	tryUpgrade := func() {
		vmOutput = host.doRunSmartContractUpgrade(input)
	}

	tryCall := func() {
		vmOutput = host.doRunSmartContractCall(input)

		if host.hasRetriableExecutionError(vmOutput) {
			log.Error("Retriable execution error detected. Will reset warm Wasmer instance.")
			host.runtimeContext.ResetWarmInstance()
		}
	}

	catch := func(caught error) {
		err = caught
		log.Error("RunSmartContractCall", "error", err)
	}

	isUpgrade := input.Function == arwen.UpgradeFunctionName
	if isUpgrade {
		TryCatch(tryUpgrade, catch, "arwen.RunSmartContractUpgrade")
	} else {
		TryCatch(tryCall, catch, "arwen.RunSmartContractCall")
	}

	return
}

// TryCatch simulates a try/catch block using golang's recover() functionality
func TryCatch(try TryFunction, catch CatchFunction, catchFallbackMessage string) {
	defer func() {
		if r := recover(); r != nil {
			err, ok := r.(error)
			if !ok {
				err = fmt.Errorf("%s, panic: %v", catchFallbackMessage, r)
			}

			catch(err)
		}
	}()

	try()
}

func (host *vmHost) hasRetriableExecutionError(vmOutput *vmcommon.VMOutput) bool {
	if !host.runtimeContext.IsWarmInstance() {
		return false
	}

	return vmOutput.ReturnMessage == "allocation error"
}

// AreInSameShard returns true if the provided addresses are part of the same shard
func (host *vmHost) AreInSameShard(leftAddress []byte, rightAddress []byte) bool {
	blockchain := host.Blockchain()
	leftShard := blockchain.GetShardOfAddress(leftAddress)
	rightShard := blockchain.GetShardOfAddress(rightAddress)

	return leftShard == rightShard
}

// IsInterfaceNil returns true if there is no value under the interface
func (host *vmHost) IsInterfaceNil() bool {
	return host == nil
}

// SetRuntimeContext sets the runtimeContext for this host, used in tests
func (host *vmHost) SetRuntimeContext(runtime arwen.RuntimeContext) {
	host.runtimeContext = runtime
}

// GetRuntimeErrors obtains the cumultated error object after running the SC
func (host *vmHost) GetRuntimeErrors() error {
	if host.runtimeContext != nil {
		return host.runtimeContext.GetAllErrors()
	}
	return nil
}

// CallArgsParser returns the CallArgsParser instance of the host
func (host *vmHost) CallArgsParser() arwen.CallArgsParser {
	return host.callArgsParser
}

// Async returns the AsyncContext instance of the host
func (host *vmHost) Async() arwen.AsyncContext {
	return host.asyncContext
}
