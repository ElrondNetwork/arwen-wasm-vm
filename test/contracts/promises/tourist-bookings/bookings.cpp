#include "../elrond/types.h"

extern "C"
{
    extern void createAsyncCall(
        // Async call info
        byte *contextIdentifier,
        int identifierLength,
        byte *destination,
        byte *value,
        byte *data,
        int dataLength,
        // Callbacks
        byte *successCallback,
        int successCallbackLen,
        byte *errorCallback,
        int errorCallbackLen,
        long long gas
    );
    long long getBlockTimestamp();
    int setStorageLock(byte *key, int keyLen, long long timeLock);
    int clearStorageLock(byte *key, int keyLen);
    int isStorageLocked(byte *key, int keyLen);
    int int64storageStore(byte *key, int keyLength, long long value);
    long long int64storageLoad(byte *key, int keyLength);
    void int64finish(long long value);
    void myTrainSuccess();
    void myTrainError();
}

void lockMyStorage();

byte successInteractionsKey[] = "storage";
byte trainAddress[] = "\0\0\0\0\0\0\0\0\x0f\x0ftrainSC...............";
// Maybe setup an initial state
extern "C" void init() {
    int64storageStore(successInteractionsKey, sizeof(successInteractionsKey), 0);
}

extern "C" void bookMyStuff() {
    lockMyStorage();
    byte asyncContext[] = "my_first_vacation";
    createAsyncCall(
        asyncContext,
        sizeof(asyncContext),
        trainAddress,
        0,
        (byte*)"bookTrain",
        9, // bookTrain
        (byte*)"myTrainSuccess",
        14,
        (byte*)"myTrainError",
        12,
        4000000
    );
}

extern "C" void myTrainSuccess() {
    clearStorageLock(successInteractionsKey, sizeof(successInteractionsKey));
    i64 counter = int64storageLoad(successInteractionsKey, sizeof(successInteractionsKey));
    counter++;
    int64storageStore(successInteractionsKey, sizeof(successInteractionsKey), counter);
}

extern "C" void myTrainError() {
    clearStorageLock(successInteractionsKey, sizeof(successInteractionsKey));
}

extern "C" void isMyStorageLocked() {
    i64 isLocked = isStorageLocked(successInteractionsKey, sizeof(successInteractionsKey));
    int64finish(isLocked);
}

void myFirstVacationBooked() {
}

void lockMyStorage() {
    long long lockTimestamp = getBlockTimestamp();
    lockTimestamp += 3600*24;
    setStorageLock(successInteractionsKey, sizeof(successInteractionsKey), lockTimestamp);
}

void releaseMyStorage() {
    clearStorageLock(successInteractionsKey, sizeof(successInteractionsKey));
}
