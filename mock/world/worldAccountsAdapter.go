package worldmock

import (
	"bytes"
	"errors"
	"fmt"

	"github.com/ElrondNetwork/arwen-wasm-vm/v1_3/arwen"
	"github.com/ElrondNetwork/elrond-go/core"
	"github.com/ElrondNetwork/elrond-go/data"
	"github.com/ElrondNetwork/elrond-go/data/state"
)

var _ state.AccountsAdapter = (*MockAccountsAdapter)(nil)

// ErrTrieHandlingNotImplemented indicates that no trie-related operations are
// currently implemented.
var ErrTrieHandlingNotImplemented = errors.New("trie handling not implemented")

// MockAccountsAdapter is an implementation of AccountsAdapter based on
// MockWorld and the accounts within it.
type MockAccountsAdapter struct {
	World     *MockWorld
	Snapshots []AccountMap
}

// NewMockAccountsAdapter instantiates a new MockAccountsAdapter.
func NewMockAccountsAdapter(world *MockWorld) *MockAccountsAdapter {
	return &MockAccountsAdapter{
		World:     world,
		Snapshots: make([]AccountMap, 0),
	}
}

// GetExistingAccount -
func (m *MockAccountsAdapter) GetExistingAccount(address []byte) (state.AccountHandler, error) {
	account, exists := m.World.AcctMap[string(address)]
	if !exists {
		return nil, arwen.ErrInvalidAccount
	}

	return account, nil
}

// LoadAccount -
func (m *MockAccountsAdapter) LoadAccount(address []byte) (state.AccountHandler, error) {
	return m.GetExistingAccount(address)
}

// SaveAccount -
func (m *MockAccountsAdapter) SaveAccount(account state.AccountHandler) error {
	mockAccount, ok := account.(*Account)
	if !ok {
		return errors.New("invalid account to save")
	}

	m.World.AcctMap.PutAccount(mockAccount)
	return nil
}

// RemoveAccount -
func (m *MockAccountsAdapter) RemoveAccount(address []byte) error {
	_, exists := m.World.AcctMap[string(address)]
	if !exists {
		return arwen.ErrInvalidAccount
	}

	m.World.AcctMap.DeleteAccount(address)
	return nil
}

// Commit -
func (m *MockAccountsAdapter) Commit() ([]byte, error) {
	m.Snapshots = make([]AccountMap, 0)
	return nil, nil
}

// Close -
func (m *MockAccountsAdapter) Close() error {
	return nil
}

// JournalLen -
func (m *MockAccountsAdapter) JournalLen() int {
	return len(m.Snapshots) - 1
}

// RevertToSnapshot -
func (m *MockAccountsAdapter) RevertToSnapshot(snapshotIndex int) error {
	if len(m.Snapshots) == 0 {
		return errors.New("no snapshots")
	}

	if snapshotIndex >= len(m.Snapshots) || snapshotIndex < 0 {
		return fmt.Errorf(
			"snapshot %d out of bounds (min 0, max %d)",
			snapshotIndex,
			len(m.Snapshots)-1)
	}

	snapshot := m.Snapshots[snapshotIndex]
	m.Snapshots = m.Snapshots[:snapshotIndex]

	// TODO should probably set BalanceDelta of all accounts to 0 as well?
	return m.World.AcctMap.LoadAccountStorageFrom(snapshot)
}

// GetNumCheckpoints -
func (m *MockAccountsAdapter) GetNumCheckpoints() uint32 {
	return uint32(len(m.Snapshots))
}

// GetCode -
func (m *MockAccountsAdapter) GetCode(codeHash []byte) []byte {
	for _, account := range m.World.AcctMap {
		if bytes.Equal(account.GetCodeHash(), codeHash) {
			return account.GetCode()
		}
	}

	return nil
}

// RootHash -
func (m *MockAccountsAdapter) RootHash() ([]byte, error) {
	return nil, ErrTrieHandlingNotImplemented
}

// RecreateTrie -
func (m *MockAccountsAdapter) RecreateTrie(_ []byte) error {
	return ErrTrieHandlingNotImplemented
}

// PruneTrie -
func (m *MockAccountsAdapter) PruneTrie(_ []byte, _ data.TriePruningIdentifier) {
}

// CancelPrune -
func (m *MockAccountsAdapter) CancelPrune(_ []byte, _ data.TriePruningIdentifier) {
}

// SnapshotState -
func (m *MockAccountsAdapter) SnapshotState(_ []byte) {
	snapshot := m.World.AcctMap.Clone()
	m.Snapshots = append(m.Snapshots, snapshot)
}

// SetStateCheckpoint -
func (m *MockAccountsAdapter) SetStateCheckpoint(_ []byte) {
}

// IsPruningEnabled -
func (m *MockAccountsAdapter) IsPruningEnabled() bool {
	return false
}

// GetAllLeaves -
func (m *MockAccountsAdapter) GetAllLeaves(_ []byte) (chan core.KeyValueHolder, error) {
	return nil, ErrTrieHandlingNotImplemented
}

// RecreateAllTries -
func (m *MockAccountsAdapter) RecreateAllTries(_ []byte) (map[string]data.Trie, error) {
	return nil, ErrTrieHandlingNotImplemented
}

// GetTrie -
func (m *MockAccountsAdapter) GetTrie(_ []byte) (data.Trie, error) {
	return nil, ErrTrieHandlingNotImplemented
}

// IsInterfaceNil -
func (m *MockAccountsAdapter) IsInterfaceNil() bool {
	return m == nil
}
