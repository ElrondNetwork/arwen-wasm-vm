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
    int int64storageStore(byte *key, int keyLength, long long value);
    long long int64storageLoad(byte *key, int keyLength);
    long long getBlockTimestamp();
    void int64finish(long long value);

    int setStorageLock(byte *key, int keyLen, long long timeLock);
    long long getStorageLock(byte *key, int keyLen);
    int isStorageLocked(byte *key, int keyLen);
    int clearStorageLock(byte *key, int keyLen);
    void bookTrainSuccess();
    void bookTrainError();
    void lockTrain();
}

byte isTrainBooked[] = "storage";
byte databaseAddress[] = "\0\0\0\0\0\0\0\0\x0f\x0f" "dataSC................";


void init() {
    int64storageStore(isTrainBooked, sizeof(isTrainBooked), 0);
}

extern "C" void bookTrain() {
    if (isStorageLocked(isTrainBooked, sizeof(isTrainBooked)) == 1) {
        // This means somebody else is booking this but the database has not recorded this yet
        return; // Should finish with error message
    }

    lockTrain();
    // Now we call external database to save the train booking
    byte asyncContext[] = "somebody_booking_train";
    byte bookTrainCallData[] = "bookTrain";
    createAsyncCall(
        asyncContext,
        sizeof(asyncContext),
        databaseAddress,
        0,
        (byte*)"bookTrain",
        9,
        (byte*)"bookTrainSuccess",
        16,
        (byte*)"bookTrainError",
        14,
        2000000
    );
}

extern "C" void lockTrain() {
    long long lockTimestamp = getBlockTimestamp();
    lockTimestamp += 3600*24;
    setStorageLock(isTrainBooked, sizeof(isTrainBooked), lockTimestamp);
}

extern "C" void bookTrainSuccess() {
    int64storageStore(isTrainBooked, sizeof(isTrainBooked), 1);
}

extern "C" void bookTrainError() {
    clearStorageLock(isTrainBooked, sizeof(isTrainBooked));
    int64storageStore(isTrainBooked, sizeof(isTrainBooked), 0);
}

extern "C" void cancelTrainBooking() {
    clearStorageLock(isTrainBooked, sizeof(isTrainBooked));
    int64storageStore(isTrainBooked, sizeof(isTrainBooked), 0);
}

extern "C" i64 isMyTrainBooked() {
    i64 isBooked = int64storageLoad(isTrainBooked, sizeof(isTrainBooked));
    int64finish(isBooked);
    return isBooked;
}
