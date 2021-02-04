#include "../elrond/context.h"
#include "../elrond/test_utils.h"

byte value[] = {0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0};

void testBuiltins1() {
	byte scAddress[] = "\0\0\0\0\0\0\0\0\x0f\x0fparentSC..............";
	byte functionName[] = "builtinClaim";
	byte functionLength = 12;

	value[31] = 96;
	u64 result = executeOnDestContext(
			500,
			scAddress,
			value,
			functionName,
			functionLength,
			0,
			0,
			0
	);

	finishResult(result);
}

void testBuiltins2() {
	byte scAddress[] = "\0\0\0\0\0\0\0\0\x0f\x0fparentSC..............";
	byte functionName[] = "builtinDoSomething";
	byte functionLength = 18;

	value[31] = 100;
	u64 result = executeOnDestContext(
			500,
			scAddress,
			value,
			functionName,
			functionLength,
			0,
			0,
			0
	);

	finishResult(result);
}

void testBuiltins3() {
	byte scAddress[] = "\0\0\0\0\0\0\0\0\x0f\x0fparentSC..............";
	byte functionName[] = "builtinDoesntExist";
	byte functionLength = 18;

	value[31] = 11;
	u64 result = executeOnDestContext(
			500,
			scAddress,
			value,
			functionName,
			functionLength,
			0,
			0,
			0
	);

	finishResult(result);
}
