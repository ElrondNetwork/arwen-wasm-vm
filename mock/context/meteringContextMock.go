package mock

import (
	"github.com/ElrondNetwork/arwen-wasm-vm/v1_3/arwen"
	"github.com/ElrondNetwork/arwen-wasm-vm/v1_3/config"
	"github.com/ElrondNetwork/elrond-vm-common"
)

var _ arwen.MeteringContext = (*MeteringContextMock)(nil)

// MeteringContextMock is used in tests to check the MeteringContext interface method calls
type MeteringContextMock struct {
	GasCost           *config.GasCost
	GasLeftMock       uint64
	GasLockedMock     uint64
	GasProvidedMock   uint64
	GasComputedToLock uint64
	BlockGasLimitMock uint64
	Err               error
}

// InitState mocked method
func (m *MeteringContextMock) InitState() {
}

// PushState mocked method
func (m *MeteringContextMock) PushState() {
}

// PopSetActiveState mocked method
func (m *MeteringContextMock) PopSetActiveState() {
}

// PopMergeActiveState mocked method
func (m *MeteringContextMock) PopMergeActiveState() {
}

// PopDiscard mocked method
func (m *MeteringContextMock) PopDiscard() {
}

// ClearStateStack mocked method
func (m *MeteringContextMock) ClearStateStack() {
}

// SetGasSchedule mocked method
func (m *MeteringContextMock) SetGasSchedule(gasSchedule config.GasScheduleMap) {
	gasCostConfig, _ := config.CreateGasConfig(gasSchedule)
	m.GasCost = gasCostConfig
}

// GasSchedule mocked method
func (m *MeteringContextMock) GasSchedule() *config.GasCost {
	return m.GasCost
}

// UseGas mocked method
func (m *MeteringContextMock) UseGas(_ uint64) {
}

// FreeGas mocked method
func (m *MeteringContextMock) FreeGas(_ uint64) {
}

// RestoreGas mocked method
func (m *MeteringContextMock) RestoreGas(_ uint64) {
}

// GasLeft mocked method
func (m *MeteringContextMock) GasLeft() uint64 {
	return m.GasLeftMock
}

// UpdateGasStateOnSuccess mocked method
func (m *MeteringContextMock) UpdateGasStateOnSuccess(_ *vmcommon.VMOutput) error {
	return nil
}

// UpdateGasStateOnFailure mocked method
func (m *MeteringContextMock) UpdateGasStateOnFailure(_ *vmcommon.VMOutput) {
}

// InitStateFromContractCallInput mocked method
func (m *MeteringContextMock) InitStateFromContractCallInput(_ *vmcommon.VMInput) {
}

// TrackGasUsedByBuiltinFunction mocked method
func (m *MeteringContextMock) TrackGasUsedByBuiltinFunction(_ *vmcommon.ContractCallInput, _ *vmcommon.VMOutput, _ *vmcommon.ContractCallInput) {
}

// GasUsedByContract mocked method
func (m *MeteringContextMock) GasUsedByContract() (uint64, uint64) {
	return 0, 0
}

// GasUsedForExecution mocked method
func (m *MeteringContextMock) GasUsedForExecution() uint64 {
	return 0
}

// GasSpentByContract mocked method
func (m *MeteringContextMock) GasSpentByContract() uint64 {
	return 0
}

// GetGasForExecution mocked method
func (m *MeteringContextMock) GetGasForExecution() uint64 {
	return 0
}

// GetGasProvided mocked method
func (m *MeteringContextMock) GetGasProvided() uint64 {
	return 0
}

// GetSCPrepareInitialCost mocked method
func (m *MeteringContextMock) GetSCPrepareInitialCost() uint64 {
	return 0
}

// BoundGasLimit mocked method
func (m *MeteringContextMock) BoundGasLimit(value int64) uint64 {
	gasLeft := m.GasLeft()
	limit := uint64(value)

	if gasLeft < limit {
		return gasLeft
	}
	return limit
}

// ComputeGasLockedForAsync mocked method
func (m *MeteringContextMock) ComputeGasLockedForAsync() uint64 {
	return m.GasComputedToLock
}

// DeductGasIfAsyncStep mocked method
func (m *MeteringContextMock) DeductGasIfAsyncStep() error {
	return m.Err
}

// UseGasBounded mocked method
func (m *MeteringContextMock) UseGasBounded(_ uint64) error {
	return m.Err
}

// UnlockGasIfAsyncCallback mocked method
func (m *MeteringContextMock) UnlockGasIfAsyncCallback() {
}

// UseGasForAsyncStep mocked method
func (m *MeteringContextMock) UseGasForAsyncStep() error {
	return m.Err
}

// UnlockGasIfAsyncStep mocked method
func (m *MeteringContextMock) UnlockGasIfAsyncStep() {
}

// GetGasLocked mocked method
func (m *MeteringContextMock) GetGasLocked() uint64 {
	return m.GasLockedMock
}

// BlockGasLimit mocked method
func (m *MeteringContextMock) BlockGasLimit() uint64 {
	return m.BlockGasLimitMock
}

// DeductInitialGasForExecution mocked method
func (m *MeteringContextMock) DeductInitialGasForExecution(_ []byte) error {
	return m.Err
}

// DeductInitialGasForDirectDeployment mocked method
func (m *MeteringContextMock) DeductInitialGasForDirectDeployment(_ arwen.CodeDeployInput) error {
	return m.Err
}

// DeductInitialGasForIndirectDeployment mocked method
func (m *MeteringContextMock) DeductInitialGasForIndirectDeployment(_ arwen.CodeDeployInput) error {
	return m.Err
}
