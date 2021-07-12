#include "../elrond/context.h"
#include "../elrond/bigInt.h"
#include "../elrond/test_utils.h"

byte dataA[20] = {0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0};
byte dataB[20] = {0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0};
byte parentKeyA[] =  "parentKeyA......................";
byte parentDataA[] = "parentDataA";
byte parentKeyB[] =  "parentKeyB......................";
byte parentDataB[] = "parentDataB";
byte childKey[] =  "childKey........................";
byte childData[] = "childData";
byte childFinish[] = "childFinish";

byte recipient[32] = "\0\0\0\0\0\0\0\0\x0F\x0F" "childTransferReceiver.";
byte value[] = {0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,12};

void childFunction() {
	int numArgs = getNumArguments();
	if (numArgs != 3) {
		byte message[] = "wrong number of arguments";
		signalError(message, 25);
	}

	didCallerPay(99);

	// This transfer will appear alongside the transfers made by the parent.
	byte transferData[100];
	getArgument(1, transferData);
	int dataLength = getArgumentLength(1);
	int result = transferValue(recipient, value, transferData, dataLength);
	if (result != 0) {
		not_ok();
	}
	
	// This storage update will appear separate from the storage updates made by the parent.
	storageStore(childKey, 32, childData, 9);

	// This finish value will appear alongside the finish values set by the parent.
	finish(childFinish, 11);
}

void childFunction_BigInts() {
	int numArgs = getNumArguments();
	if (numArgs != 3) {
		byte message[] = "wrong number of arguments";
		signalError(message, 25);
	}

	didCallerPay(99);

	int status = 0;

	bigInt intA = int64getArgument(0);
	bigInt intB = int64getArgument(1);
	bigInt intC = int64getArgument(2);

	long long a = bigIntGetOrCreateInt64(intA);
	long long b = bigIntGetOrCreateInt64(intB);
	long long c = bigIntGetOrCreateInt64(intC);

	// The parent sent bigInt ID 0 as argument, but since the parent bigInt context is
	// separate from the child, the ID 0 was already taken inside didCallerPay(),
	// and now it equals to 99, the call value.
	if (a != 99) {
		not_ok();
		byte msg[] = "nr a";
		finish(msg, 4);
		int64finish(a);
	}
	if (b != 0) {
		not_ok();
		byte msg[] = "nr b";
		finish(msg, 4);
		int64finish(b);
	}
	if (c != 0) {
		not_ok();
		byte msg[] = "nr c";
		finish(msg, 4);
		int64finish(c);
	}

	bigInt intX = bigIntNew(256);
	if (intX != 3) {
		not_ok();
		byte msg[] = "nr x";
		finish(msg, 4);
		int64finish(intX);
		status = 1;
	}

	if (status == 0) {
		byte msg[] = "child ok";
		finish(msg, 8);
	} else {
		byte msg[] = "child not ok";
		finish(msg, 12);
	}
}

void childFunction_OutOfGas() {
	int numArgs = getNumArguments();
	if (numArgs != 0) {
		byte message[] = "wrong number of arguments";
		signalError(message, 25);
	}

	didCallerPay(99);

	storageStore(childKey, 32, childData, 9);
	finish(childFinish, 11);
	bigIntSetInt64(12, 88);

	// Start infinite loop.
	byte msg[] = "rockets";
	while (1) {
		finish(msg, 7);
	}
}
