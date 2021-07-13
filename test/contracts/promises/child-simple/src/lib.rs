#![no_std]
#![no_main]
#![allow(unused_attributes)]
#![feature(lang_items)]

use elrond_wasm::api::{EndpointArgumentApi, EndpointFinishApi, ErrorApi};
use elrond_wasm_node::ArwenApiImpl;

pub static EEI: ArwenApiImpl = ArwenApiImpl{};

#[no_mangle]
pub extern "C" fn answer() {
    EEI.finish_u64(42);
}

#[no_mangle]
pub extern "C" fn answer_wrong() {
    EEI.finish_u64(24);
}

// receives u64 as argument and returns it back
#[no_mangle]
pub extern "C" fn echo() {
    EEI.check_num_arguments(1);

    let arg = EEI.get_argument_u64(0);

    EEI.finish_u64(arg);
}

#[no_mangle]
pub extern "C" fn fail() {
    EEI.signal_error(&b"fail"[..]);
}
