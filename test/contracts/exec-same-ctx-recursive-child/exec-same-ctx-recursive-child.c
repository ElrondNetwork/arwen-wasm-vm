#include "../elrond/context.h"
#include "../elrond/bigInt.h"
#include "../elrond/types.h"
#include "../elrond/test_utils.h"

u64 maxGasForCalls = 100000;
 
byte parentAddress[32] = "\0\0\0\0\0\0\0\0\x0F\x0F" "parentSC..............";
byte bigIntCounterKey[] = "recursiveIterationBigCounter....";
bigInt bigIntCounterID = 88;
byte executeValue[] = {0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,5};
byte smallCounterKey[] = "recursiveIterationCounter.......";

byte arguments[1] = {0};
int argumentsLengths[1] = {1};

void childCallsParent() {
	int numArgs = getNumArguments();
	if (numArgs != 1) {
		byte message[] = "wrong number of arguments";
		signalError(message, 25);
	}

	byte iteration = (byte) int64getArgument(0);

    bigIntGetOrCreateInt64(bigIntCounterID);
	storeIterationNumber(iteration, 'C');
	finishIterationNumber(iteration, 'C');

    incrementIterCounter(smallCounterKey);
	incrementBigIntCounter(bigIntCounterID);

    // Run next iteration.
	byte functionName[] = "parentCallsChild";
	if (iteration > 0) {
		arguments[0] = iteration - 1;
        int result = executeOnSameContext(
            maxGasForCalls,
            parentAddress,
            executeValue,
            functionName,
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
