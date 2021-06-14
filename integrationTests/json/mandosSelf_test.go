package vmjsonintegrationtest

import (
	"testing"

	"github.com/stretchr/testify/require"
)

// Tests Mandos consistency, no smart contracts.
func TestMandosSelfTest(t *testing.T) {
	runTestsInFolder(t, "mandos-self-test", []string{
		"mandos-self-test/builtin-func-esdt-transfer.scen.json",
	})
}

func TestMandoSetAccountAddressLengthErr1(t *testing.T) {
	err := runSingleTest(t, "mandos-self-test/set-check", "set-account-addr-len.err1.json")
	require.EqualError(t, err,
		"error processing steps: cannot parse set state step: account address is not 32 bytes in length")
}

func TestMandoSetAccountAddressLengthErr2(t *testing.T) {
	err := runSingleTest(t, "mandos-self-test/set-check", "set-account-addr-len.err2.json")
	require.EqualError(t, err,
		"error processing steps: error parsing new addresses: account address is not 32 bytes in length")
}

func TestMandoSetAccountSCAddressErr1(t *testing.T) {
	err := runSingleTest(t, "mandos-self-test/set-check", "set-account-sc-addr.err1.json")
	require.EqualError(t, err,
		"\"setState\" step validation failed for account \"address:not-a-sc-address\": account has a smart contract address, but has no code: 0x6e6f742d612d73632d616464726573735f5f5f5f5f5f5f5f5f5f5f5f5f5f5f5f")
}

func TestMandoSetAccountSCAddressErr2(t *testing.T) {
	err := runSingleTest(t, "mandos-self-test/set-check", "set-account-sc-addr.err2.json")
	require.EqualError(t, err,
		"\"setState\" step validation failed for account \"sc:should-be-sc\": account has code but not a smart contract address: 000000000000000073686f756c642d62652d73635f5f5f5f5f5f5f5f5f5f5f5f")
}

func TestMandoSetAccountSCAddressErr3(t *testing.T) {
	err := runSingleTest(t, "mandos-self-test/set-check", "set-account-sc-addr.err3.json")
	require.EqualError(t, err,
		"address in \"setState\" \"newAddresses\" field should have SC format: address:not-a-sc-address")
}

func TestMandosCheckNonceErr(t *testing.T) {
	err := runSingleTest(t, "mandos-self-test/set-check", "set-check-nonce.err.json")
	require.EqualError(t, err,
		"bad account nonce. Account: address:the-address. Want: \"1002\". Have: \"1001\"")
}

func TestMandosCheckBalanceErr(t *testing.T) {
	err := runSingleTest(t, "mandos-self-test/set-check", "set-check-balance.err.json")
	require.EqualError(t, err,
		"bad account balance. Account: address:the-address. Want: \"1,000,002\". Have: \"1000001\"")
}

func TestMandosCheckUsernameErr(t *testing.T) {
	err := runSingleTest(t, "mandos-self-test/set-check", "set-check-username.err.json")
	require.EqualError(t, err,
		"bad account username. Account: address:the-address. Want: \"str:wrong.elrond\". Have: \"str:theusername.elrond\"")
}

func TestMandosCheckCodeErr(t *testing.T) {
	err := runSingleTest(t, "mandos-self-test/set-check", "set-check-code.err.json")
	require.EqualError(t, err,
		"bad account code. Account: sc:contract-address. Want: \"file:set-check-code.scen.json\". Have: \"0x7b0a2020202022636f6d...\"")
}

func TestMandosCheckStorageErr1(t *testing.T) {
	err := runSingleTest(t, "mandos-self-test/set-check", "set-check-storage.err1.json")
	require.EqualError(t, err,
		"wrong account storage for account \"address:the-address\":\n"+
			"  for key 0x6b65792d63 (str:key-c): Want: \"str:another-value\". Have: \"0x76616c75652d63 (str:value-c)\"")
}

func TestMandosCheckStorageErr2(t *testing.T) {
	err := runSingleTest(t, "mandos-self-test/set-check", "set-check-storage.err2.json")
	require.EqualError(t, err,
		"wrong account storage for account \"address:the-address\":\n"+
			"  for key 0x6b65792d63 (str:key-c): Want: \"\". Have: \"0x76616c75652d63 (str:value-c)\"")
}

func TestMandosCheckStorageErr3(t *testing.T) {
	err := runSingleTest(t, "mandos-self-test/set-check", "set-check-storage.err3.json")
	require.EqualError(t, err,
		"wrong account storage for account \"address:the-address\":\n"+
			"  for key 0x6b65792d64 (str:key-d): Want: \"str:value-d\". Have: \"\"")
}

func TestMandosCheckStorageErr4(t *testing.T) {
	err := runSingleTest(t, "mandos-self-test/set-check", "set-check-storage.err4.json")
	require.EqualError(t, err,
		"wrong account storage for account \"address:the-address\":\n"+
			"  for key 0x6b65792d63 (str:key-c): Want: \"\". Have: \"0x76616c75652d63 (str:value-c)\"")
}

func TestMandosCheckStorageErr5(t *testing.T) {
	err := runSingleTest(t, "mandos-self-test/set-check", "set-check-storage.err5.json")
	require.EqualError(t, err,
		"wrong account storage for account \"address:the-address\":\n"+
			"  for key 0x6b65792d62 (str:key-b): Want: \"str:another-b\". Have: \"0x76616c75652d62 (str:value-b)\"")
}

func TestMandosCheckESDTErr1(t *testing.T) {
	err := runSingleTest(t, "mandos-self-test/set-check", "set-check-esdt.err1.json")
	require.EqualError(t, err,
		`mismatch for account "address:the-address":
  for token: NFT-123456, nonce: 1: Bad balance. Want: "4". Have: "1"
  for token: NFT-123456, nonce: 1: Bad creator. Want: "address:another-address". Have: "address:the-address"
  for token: NFT-123456, nonce: 1: Bad royalties. Want: "2001". Have: "2000"
  for token: NFT-123456, nonce: 1: Bad hash. Want: "keccak256:str:another_hash". Have: 0x54e3ea4bdef3b22154767a2cae081fca2bec2eae1ec62ee71308cb2a300d675d (str:"T\xe3\xeaK\xde\xf3\xb2!Tvz,\xae\b\x1f\xca+\xec.\xae\x1e\xc6.\xe7\x13\b\xcb*0\rg]")
  for token: NFT-123456, nonce: 1: Bad URI. Want: [
    "str:www.cool_nft.com/another_nft.jpg"
]. Have: "str:www.cool_nft.com/my_nft.jpg"
  for token: NFT-123456, nonce: 1: Bad attributes. Want: "str:other_attributes". Have: "str:serialized_attributes"`)
}
