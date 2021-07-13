#![no_std]

pub use elrond_wasm::{Address, Vec};

pub const PARENT_ADDRESS: [u8; 32] = [
    0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x0F, 0x0F, b'p', b'a', b'r', b'e', b'n', b't',
    b'S', b'C', b'.', b'.', b'.', b'.', b'.', b'.', b'.', b'.', b'.', b'.', b'.', b'.', b'.', b'.',
];
pub const CHILD_ADDRESS: [u8; 32] = [
    0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x0F, 0x0F, b'c', b'h', b'i', b'l', b'd', b'S',
    b'C', b'.', b'.', b'.', b'.', b'.', b'.', b'.', b'.', b'.', b'.', b'.', b'.', b'.', b'.', b'.',
];

pub const FIRST_CONTRACT_ADDRESS: [u8; 32] = [
    0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x0F, 0x0F, b'f', b'i', b'r', b's', b't', b'C',
    b'o', b'n', b't', b'r', b'a', b'c', b't', b'.', b'.', b'.', b'.', b'.', b'.', b'.', b'.', b'.',
];
pub const SECOND_CONTRACT_ADDRESS: [u8; 32] = [
    0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x0F, 0x0F, b's', b'e', b'c', b'o', b'n', b'd',
    b'C',b'o', b'n', b't', b'r', b'a', b'c', b't', b'.', b'.', b'.', b'.', b'.', b'.', b'.', b'.',
];
pub const THIRD_CONTRACT_ADDRESS: [u8; 32] = [
    0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x0F, 0x0F, b't', b'h', b'i', b'r', b'd', b'C',
    b'o', b'n', b't', b'r', b'a', b'c', b't', b'.', b'.', b'.', b'.', b'.', b'.', b'.', b'.', b'.',
];

pub const GAS_100K: i64 = 100000;
pub const GAS_500K: i64 = 500000;
pub const GAS_5M: i64 = 5000000;
pub const GAS_10M: i64 = 10000000;
pub const GAS_50M: i64 = 50000000;
pub const GAS_100M: i64 = 100000000;

pub const ZERO: [u8; 32] = [0u8; 32];
pub const EMPTY_SLICE: &[u8] = &[];

pub fn construct_storage_key(key_parts: &[&[u8]]) -> Vec<u8> {
    let mut key = Vec::new();

    for part in key_parts {
        key.extend_from_slice(part);
    }

    key
}

#[inline(always)]
pub fn create_async_call(
    group_id: &[u8],
    destination: &Address,
    value: &[u8],
    data: &[u8],
    success_callback_name: &[u8],
    error_callback_name: &[u8],
    gas: i64) 
{
    unsafe {
        createAsyncCall(
            group_id.as_ref().as_ptr(),
            group_id.len() as i32,
            destination.as_ref().as_ptr(),
            value.as_ref().as_ptr(),
            data.as_ref().as_ptr(),
            data.len() as i32,
            success_callback_name.as_ref().as_ptr(),
            success_callback_name.len() as i32,
            error_callback_name.as_ref().as_ptr(),
            error_callback_name.len() as i32,
            gas
        )
    }
}

extern {
    fn createAsyncCall(
        groupIDOffset: *const u8, groupIDLength: i32,
        destOffset: *const u8,
        valueOffset: *const u8,
        dataOffset: *const u8,
        dataLength: i32,
        successCallbackNameOffset: *const u8, successCallbackNameLen: i32,
        errorCallbackNameOffset: *const u8, errorCallbackNameLen: i32,
        gas: i64
    );
}
