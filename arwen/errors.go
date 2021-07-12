package arwen

import (
	"errors"
	"fmt"
)

// ErrReturnCodeNotOk signals that the returned code is different than vmcommon.Ok
var ErrReturnCodeNotOk = errors.New("return code is not ok")

// ErrInvalidCallOnReadOnlyMode signals that an operation is not permitted due to read only mode
var ErrInvalidCallOnReadOnlyMode = errors.New("operation not permitted in read only mode")

// ErrNotEnoughGas signals that there is not enough gas for the operation
var ErrNotEnoughGas = errors.New("not enough gas")

// ErrUnhandledRuntimeBreakpoint signals that the runtime breakpoint is unhandled
var ErrUnhandledRuntimeBreakpoint = errors.New("unhandled runtime breakpoint")

// ErrSignalError is given when the smart contract signals an error
var ErrSignalError = errors.New("error signalled by smartcontract")

// ErrExecutionFailed signals that the execution failed
var ErrExecutionFailed = errors.New("execution failed")

// ErrBadBounds signals that a certain variable is out of bounds
var ErrBadBounds = errors.New("bad bounds")

// ErrBadLowerBounds signals that a certain variable is lower than allowed
var ErrBadLowerBounds = fmt.Errorf("%w (lower)", ErrBadBounds)

// ErrBadUpperBounds signals that a certain variable is higher than allowed
var ErrBadUpperBounds = fmt.Errorf("%w (upper)", ErrBadBounds)

// ErrNegativeLength signals that the given length is less than 0
var ErrNegativeLength = errors.New("negative length")

// ErrFailedTransfer signals that the transfer operation has failed
var ErrFailedTransfer = errors.New("failed transfer")

// ErrTransferInsufficientFunds signals that the transfer has failed due to insufficient funds
var ErrTransferInsufficientFunds = fmt.Errorf("%w (insufficient funds)", ErrFailedTransfer)

// ErrTransferNegativeValue signals that the transfer has failed due to the fact that the value is less than 0
var ErrTransferNegativeValue = fmt.Errorf("%w (negative value)", ErrFailedTransfer)

// ErrUpgradeFailed signals that the upgrade encountered an error
var ErrUpgradeFailed = errors.New("upgrade failed")

// ErrInvalidUpgradeArguments signals that the upgrade process failed due to invalid arguments
var ErrInvalidUpgradeArguments = fmt.Errorf("%w (invalid arguments)", ErrUpgradeFailed)

// ErrInvalidFunction signals that the function is invalid
var ErrInvalidFunction = errors.New("invalid function")

// ErrInitFuncCalledInRun signals that the init func was called directly, which is forbidden
var ErrInitFuncCalledInRun = fmt.Errorf("%w (calling init() directly is forbidden)", ErrInvalidFunction)

// ErrCallBackFuncCalledInRun signals that a callback func was called directly, which is forbidden
var ErrCallBackFuncCalledInRun = fmt.Errorf("%w (calling callBack() directly is forbidden)", ErrInvalidFunction)

// ErrCallBackFuncNotExpected signals that an unexpected callback was received
var ErrCallBackFuncNotExpected = fmt.Errorf("%w (unexpected callback was received)", ErrInvalidFunction)

// ErrFuncNotFound signals that the the function does not exist
var ErrFuncNotFound = fmt.Errorf("%w (not found)", ErrInvalidFunction)

// ErrInvalidFunctionName signals that the function name is invalid
var ErrInvalidFunctionName = fmt.Errorf("%w (invalid name)", ErrInvalidFunction)

// ErrFunctionNonvoidSignature signals that the signature for the function is invalid
var ErrFunctionNonvoidSignature = fmt.Errorf("%w (nonvoid signature)", ErrInvalidFunction)

// ErrContractInvalid signals that the contract code is invalid
var ErrContractInvalid = fmt.Errorf("invalid contract code")

// ErrContractNotFound signals that the contract was not found
var ErrContractNotFound = fmt.Errorf("%w (not found)", ErrContractInvalid)

// ErrMemoryDeclarationMissing signals that a memory declaration is missing
var ErrMemoryDeclarationMissing = fmt.Errorf("%w (missing memory declaration)", ErrContractInvalid)

// ErrMaxInstancesReached signals that the max number of Wasmer instances has been reached.
var ErrMaxInstancesReached = fmt.Errorf("%w (max instances reached)", ErrExecutionFailed)

// ErrStoreElrondReservedKey signals that an attempt to write under an reserved key has been made
var ErrStoreElrondReservedKey = errors.New("cannot write to storage under Elrond reserved key")

// ErrCannotWriteProtectedKey signals an attempt to write to a protected key, while storage protection is enforced
var ErrCannotWriteProtectedKey = errors.New("cannot write to protected key")

// ErrNonPayableFunctionEgld signals that a non-payable function received non-zero call value
var ErrNonPayableFunctionEgld = errors.New("function does not accept EGLD payment")

// ErrNonPayableFunctionEsdt signals that a non-payable function received non-zero ESDT call value
var ErrNonPayableFunctionEsdt = errors.New("function does not accept ESDT payment")

// ErrArgIndexOutOfRange signals that the argument index is out of range
var ErrArgIndexOutOfRange = errors.New("argument index out of range")

// ErrArgOutOfRange signals that the argument is out of range
var ErrArgOutOfRange = errors.New("argument out of range")

// ErrStorageValueOutOfRange signals that the storage value is out of range
var ErrStorageValueOutOfRange = errors.New("storage value out of range")

// ErrDivZero signals that an attempt to divide by 0 has been made
var ErrDivZero = errors.New("division by 0")

// ErrBitwiseNegative signals that an attempt to apply a bitwise operation on negative numbers has been made
var ErrBitwiseNegative = errors.New("bitwise operations only allowed on positive integers")

// ErrShiftNegative signals that an attempt to apply a bitwise shift operation on negative numbers has been made
var ErrShiftNegative = errors.New("bitwise shift operations only allowed on positive integers and by a positive amount")

// ErrAsyncContextDoesNotExist signals that the async context does not exist
var ErrAsyncContextDoesNotExist = errors.New("async context does not exist")

// ErrInvalidAccount signals that a certain account does not exist
var ErrInvalidAccount = errors.New("account does not exist")

// ErrDeploymentOverExistingAccount signals that an attempt to deploy a new SC over an already existing account has been made
var ErrDeploymentOverExistingAccount = errors.New("cannot deploy over existing account")

// ErrAccountNotPayable signals that the value transfer to a non payable contract is not possible
var ErrAccountNotPayable = errors.New("sending value to non payable contract")

// ErrInvalidPublicKeySize signals that the public key size is invalid
var ErrInvalidPublicKeySize = errors.New("invalid public key size")

// ErrNilCallbackFunction signals that a nil callback function has been provided
var ErrNilCallbackFunction = errors.New("nil callback function")

// ErrUpgradeNotAllowed signals that an upgrade is not allowed
var ErrUpgradeNotAllowed = errors.New("upgrade not allowed")

// ErrNilContract signals that the contract is nil
var ErrNilContract = errors.New("nil contract")

// ErrBuiltinCallOnSameContextDisallowed signals that calling a built-in function on the same context is not allowed
var ErrBuiltinCallOnSameContextDisallowed = errors.New("calling built-in function on the same context is disallowed")

// ErrSyncExecutionNotInSameShard signals that the sync execution request is not in the same shard
var ErrSyncExecutionNotInSameShard = errors.New("sync execution request is not in the same shard")

// ErrInputAndOutputGasDoesNotMatch is raised when the output gas (gas used + gas locked + gas remaining)
// is not equal to the input gas
var ErrInputAndOutputGasDoesNotMatch = errors.New("input and output gas does not match")

// ErrTransferValueOnESDTCall signals that balance transfer was given in esdt call
var ErrTransferValueOnESDTCall = errors.New("transfer value on esdt call")

// ErrNoBigIntUnderThisHandle signals that there is no bigInt for the given handle
var ErrNoBigIntUnderThisHandle = errors.New("no bigInt under the given handle")

// ErrLengthOfBufferNotCorrect signals that length of the buffer is not correct
var ErrLengthOfBufferNotCorrect = errors.New("length of buffer is not correct")

// ErrNoEllipticCurveUnderThisHandle singals that there is no elliptic curve for the given handle
var ErrNoEllipticCurveUnderThisHandle = errors.New("no elliptic curve under the given handle")

// ErrPointNotOnCurve signals that the point to be used is not on curve
var ErrPointNotOnCurve = errors.New("point is not on curve")
