#![no_std]
#![no_main]
#![allow(unused_attributes)]
#![feature(lang_items)]

use elrond_wasm::api::{ContractHookApi, EndpointArgumentApi, EndpointFinishApi, StorageReadApi, StorageWriteApi};
use elrond_wasm_node::ArwenApiImpl;

use promises_common::*;

pub static EEI: ArwenApiImpl = ArwenApiImpl{};

const SUCCESS_CALLBACK_ARGUMENT_KEY: &[u8] = b"SuccessCallbackArg";
const FAIL_CALLBACK_ARGUMENT_KEY: &[u8] = b"FailCallbackArg";
const CURRENT_STORAGE_INDEX_KEY: &[u8] = b"CurrentStorageIndex";

const COMMON_GROUP_ID: &[u8] = b"testgroup";
const SUCCESS_CALLBACK_NAME: &[u8] = b"success_callback";
const FAIL_CALLBACK_NAME: &[u8] = b"fail_callback";

#[no_mangle]
pub extern "C" fn answer() {
    EEI.finish_u64(42);
}

#[no_mangle]
pub extern "C" fn call_caller() {
    let caller = EEI.get_caller();

    create_async_call(COMMON_GROUP_ID,
        &caller,
        &ZERO,
        b"answer",
        SUCCESS_CALLBACK_NAME,
        FAIL_CALLBACK_NAME,
        GAS_100K);
}

// receives call data as argument
#[no_mangle]
pub extern "C" fn call_second_contract() {
    EEI.check_num_arguments(1);

    let call_data = EEI.get_argument_vec_u8(0);

    create_async_call(COMMON_GROUP_ID,
        &Address::from(SECOND_CONTRACT_ADDRESS),
        &ZERO,
        call_data.as_slice(),
        SUCCESS_CALLBACK_NAME,
        FAIL_CALLBACK_NAME,
        GAS_100M);
}

// receives call data as argument
#[no_mangle]
pub extern "C" fn call_third_contract() {
    create_async_call(COMMON_GROUP_ID,
        &Address::from(SECOND_CONTRACT_ADDRESS),
        &ZERO,
        b"answer",
        SUCCESS_CALLBACK_NAME,
        FAIL_CALLBACK_NAME,
        GAS_10M);
}

// receives call data as arguments
#[no_mangle]
pub extern "C" fn call_second_and_third_contract() {
    EEI.check_num_arguments(2);

    let call_data_for_second_contract = EEI.get_argument_vec_u8(0);
    let call_data_for_third_contract = EEI.get_argument_vec_u8(1);

    create_async_call(COMMON_GROUP_ID,
        &Address::from(SECOND_CONTRACT_ADDRESS),
        &ZERO,
        call_data_for_second_contract.as_slice(),
        SUCCESS_CALLBACK_NAME,
        FAIL_CALLBACK_NAME,
        GAS_100K);

    create_async_call(COMMON_GROUP_ID,
        &Address::from(THIRD_CONTRACT_ADDRESS),
        &ZERO,
        call_data_for_third_contract.as_slice(),
        SUCCESS_CALLBACK_NAME,
        FAIL_CALLBACK_NAME,
        GAS_100K);
}

// callbacks

// first argument is "0" for success, followed by data passed by finish() in callee contract
#[no_mangle]
pub extern "C" fn success_callback() {
    let num_args = EEI.get_num_arguments();
    let mut storage_index = EEI.storage_load_u64(&CURRENT_STORAGE_INDEX_KEY);

    for arg_index in 0..num_args {
        let arg = EEI.get_argument_u64(arg_index);
        let storage_key = construct_storage_key(&[SUCCESS_CALLBACK_ARGUMENT_KEY, &[storage_index as u8]]);

        storage_index += 1;
        EEI.storage_store_u64(&storage_key, arg);
    }

    EEI.storage_store_u64(&CURRENT_STORAGE_INDEX_KEY, storage_index);
}

// first argument is error code, followed by error message
#[no_mangle]
pub extern "C" fn fail_callback() {
    let expected_num_args = 2;
    EEI.check_num_arguments(expected_num_args);

    let mut storage_index = EEI.storage_load_u64(&CURRENT_STORAGE_INDEX_KEY);

    for arg_index in 0..expected_num_args {
        let arg = EEI.get_argument_vec_u8(arg_index);
        let storage_key = construct_storage_key(&[FAIL_CALLBACK_ARGUMENT_KEY, &[storage_index as u8]]);
    
        storage_index += 1;
        EEI.storage_store_slice_u8(&storage_key, &arg);
    }

    EEI.storage_store_u64(&CURRENT_STORAGE_INDEX_KEY, storage_index);
}
