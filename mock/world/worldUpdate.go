package worldmock

import (
	"errors"
	"math/big"

	"github.com/ElrondNetwork/arwen-wasm-vm/crypto/hashing"
	"github.com/ElrondNetwork/elrond-go/core/vmcommon"
)

// UpdateBalance sets a new balance to an account
func (b *MockWorld) UpdateBalance(address []byte, newBalance *big.Int) error {
	acct := b.AcctMap.GetAccount(address)
	if acct == nil {
		return errors.New("method UpdateBalance expects an existing address")
	}
	acct.Balance = newBalance
	return nil
}

// UpdateBalanceWithDelta changes balance of an account by a given amount
func (b *MockWorld) UpdateBalanceWithDelta(address []byte, balanceDelta *big.Int) error {
	acct := b.AcctMap.GetAccount(address)
	if acct == nil {
		return errors.New("method UpdateBalanceWithDelta expects an existing address")
	}
	acct.Balance = big.NewInt(0).Add(acct.Balance, balanceDelta)
	return nil
}

// UpdateWorldStateBefore performs gas payment, before transaction
func (b *MockWorld) UpdateWorldStateBefore(
	fromAddr []byte,
	gasLimit uint64,
	gasPrice uint64) error {

	acct := b.AcctMap.GetAccount(fromAddr)
	if acct == nil {
		return errors.New("method UpdateWorldStateBefore expects an existing address")
	}
	acct.Nonce++
	gasPayment := big.NewInt(0).Mul(
		big.NewInt(0).SetUint64(gasLimit),
		big.NewInt(0).SetUint64(gasPrice))
	if acct.Balance.Cmp(gasPayment) < 0 {
		return errors.New("not enough balance to pay gas upfront")
	}
	acct.Balance.Sub(acct.Balance, gasPayment)
	return nil
}

// UpdateAccounts should be called after the VM test has run, to update world state
func (b *MockWorld) UpdateAccounts(
	outputAccounts map[string]*vmcommon.OutputAccount,
	accountsToDelete [][]byte) error {

	for _, modAcct := range outputAccounts {
		acct := b.AcctMap.GetAccount(modAcct.Address)
		if acct == nil {
			acct = &Account{
				Exists:       false,
				Address:      modAcct.Address,
				Nonce:        0,
				Balance:      zero,
				Storage:      make(map[string][]byte),
				Code:         nil,
				OwnerAddress: modAcct.CodeDeployerAddress,
				ESDTData:     make(map[string]*ESDTData),
			}
			b.AcctMap.PutAccount(acct)
		}
		acct.Exists = true
		if modAcct.BalanceDelta != nil {
			acct.Balance = big.NewInt(0).Add(acct.Balance, modAcct.BalanceDelta)
		} else {
			acct.Balance = modAcct.Balance
		}
		if modAcct.Nonce > acct.Nonce {
			acct.Nonce = modAcct.Nonce
		}
		if len(modAcct.Code) > 0 {
			acct.Code = modAcct.Code
			hasher := hashing.NewHasher()
			acct.CodeHash, _ = hasher.Sha256(acct.Code)
		}
		if len(modAcct.OutputTransfers) > 0 && len(modAcct.OutputTransfers[0].Data) > 0 {
			acct.AsyncCallData = string(modAcct.OutputTransfers[0].Data)
		}

		for _, stu := range modAcct.StorageUpdates {
			acct.Storage[string(stu.Offset)] = stu.Data
		}
	}

	// commit ESDT data
	for _, acct := range b.AcctMap {
		for _, esdtData := range acct.ESDTData {
			esdtData.Balance = esdtData.Balance.Add(esdtData.Balance, esdtData.BalanceDelta)
			esdtData.BalanceDelta = big.NewInt(0)
		}
	}

	for _, delAddr := range accountsToDelete {
		b.AcctMap.DeleteAccount(delAddr)
	}

	return nil

}

// RollbackChanges should be called after the VM test has run, if the tx has failed
func (b *MockWorld) RollbackChanges() {
	// discard ESDT deltas
	// so they don't interfere with future txs
	for _, acct := range b.AcctMap {
		for _, esdtData := range acct.ESDTData {
			esdtData.BalanceDelta = big.NewInt(0)
		}
	}
}
