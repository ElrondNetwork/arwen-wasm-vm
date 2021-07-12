#include "../elrond/context.h"
#include "../elrond/bigInt.h"
#include "../elrond/types.h"
#include "../elrond/test_utils.h"

u64 maxGasForCalls = 100000;

byte selfAddress[] = "\0\0\0\0\0\0\0\0\x0F\x0F" "parentSC..............";
byte executeValue[] = {0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,5};
byte arguments[1] = {0};
int argumentsLengths[1] = {1};

byte smallCounterKey[] = "recursiveIterationCounter.......";
byte bigIntCounterKey[] = "recursiveIterationBigCounter....";
bigInt bigIntCounterID = 16;

void callRecursive() {
	int numArgs = getNumArguments();
	if (numArgs != 1) {
		byte message[] = "wrong number of arguments";
		signalError(message, 25);
	}

	byte iteration = (byte) int64getArgument(0);

    bigIntGetOrCreateInt64(bigIntCounterID);
    finishIterationNumber(iteration, 'R');
    storeIterationNumber(iteration, 'R');
    incrementIterCounter(smallCounterKey);
    incrementBigIntCounter(bigIntCounterID);

    // Run next iteration.
    byte functionName[] = "callRecursive";
    if (iteration > 0) {
        arguments[0] = iteration - 1;
        int result = executeOnDestContext(
            maxGasForCalls,
            selfAddress,
            executeValue,
            functionName,
            13,
            1,
            (byte*)argumentsLengths,
            arguments
        );
        finishResult(result);
    } else {
        bigIntStorageStoreUnsigned(bigIntCounterKey, 32, bigIntCounterID);
    }
}

void callRecursiveMutualMethods() {
	int numArgs = getNumArguments();
	if (numArgs != 1) {
		byte message[] = "wrong number of arguments";
		signalError(message, 25);
	}

	byte iteration = (byte) int64getArgument(0);
    if (iteration < 2) {
        byte message[] = "need number of recursive calls >= 2";
        signalError(message, 25);
    }

    byte startMsg[] = "start recursive mutual calls";
    finish(startMsg, 28);

    byte functionNameB[] = "recursiveMethodA";
    if (iteration > 0) {
        arguments[0] = iteration;
        int result = executeOnDestContext(
            maxGasForCalls,
            selfAddress,
            executeValue,
            functionNameB,
            16,
            1,
            (byte*)argumentsLengths,
            arguments
        );
        finishResult(result);
    }

    byte endMsg[] = "end recursive mutual calls";
    finish(endMsg, 26);
}

void recursiveMethodA() {
	byte iteration = (byte) int64getArgument(0);

    bigIntGetOrCreateInt64(bigIntCounterID);
    finishIterationNumber(iteration, 'A');
    storeIterationNumber(iteration, 'A');
    incrementIterCounter(smallCounterKey);
    incrementBigIntCounter(bigIntCounterID);

    byte functionNameB[] = "recursiveMethodB";
    if (iteration > 0) {
        arguments[0] = iteration - 1;
        int result = executeOnDestContext(
            maxGasForCalls,
            selfAddress,
            executeValue,
            functionNameB,
            16,
            1,
            (byte*)argumentsLengths,
            arguments
        );
        finishResult(result);
    } else {
        bigIntStorageStoreUnsigned(bigIntCounterKey, 32, bigIntCounterID);
    }
}

void recursiveMethodB() {
	byte iteration = (byte) int64getArgument(0);

    bigIntGetOrCreateInt64(bigIntCounterID);
    finishIterationNumber(iteration, 'B');
    storeIterationNumber(iteration, 'B');
    incrementIterCounter(smallCounterKey);
    incrementBigIntCounter(bigIntCounterID);

    byte functionNameB[] = "recursiveMethodA";
    if (iteration > 0) {
        arguments[0] = iteration - 1;
        int result = executeOnDestContext(
            maxGasForCalls,
            selfAddress,
            executeValue,
            functionNameB,
            16,
            1,
            (byte*)argumentsLengths,
            arguments
        );
        finishResult(result);
    } else {
        bigIntStorageStoreUnsigned(bigIntCounterKey, 32, bigIntCounterID);
    }
}
