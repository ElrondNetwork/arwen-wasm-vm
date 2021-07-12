package config

import "github.com/ElrondNetwork/arwen-wasm-vm/v1_3/wasmer"

type GasCost struct {
	BaseOperationCost BaseOperationCost
	BigIntAPICost     BigIntAPICost
	EthAPICost        EthAPICost
	ElrondAPICost     ElrondAPICost
	CryptoAPICost     CryptoAPICost
	WASMOpcodeCost    WASMOpcodeCost
}

type BaseOperationCost struct {
	StorePerByte      uint64
	ReleasePerByte    uint64
	DataCopyPerByte   uint64
	PersistPerByte    uint64
	CompilePerByte    uint64
	AoTPreparePerByte uint64
	GetCode           uint64
}

type ElrondAPICost struct {
	GetSCAddress         uint64
	GetOwnerAddress      uint64
	IsSmartContract      uint64
	GetShardOfAddress    uint64
	GetExternalBalance   uint64
	GetBlockHash         uint64
	TransferValue        uint64
	GetArgument          uint64
	GetFunction          uint64
	GetNumArguments      uint64
	StorageStore         uint64
	StorageLoad          uint64
	GetCaller            uint64
	GetCallValue         uint64
	Log                  uint64
	Finish               uint64
	SignalError          uint64
	GetBlockTimeStamp    uint64
	GetGasLeft           uint64
	Int64GetArgument     uint64
	Int64StorageStore    uint64
	Int64StorageLoad     uint64
	Int64Finish          uint64
	GetStateRootHash     uint64
	GetBlockNonce        uint64
	GetBlockEpoch        uint64
	GetBlockRound        uint64
	GetBlockRandomSeed   uint64
	ExecuteOnSameContext uint64
	ExecuteOnDestContext uint64
	DelegateExecution    uint64
	ExecuteReadOnly      uint64
	AsyncCallStep        uint64
	AsyncCallbackGasLock uint64
	CreateContract       uint64
	GetReturnData        uint64
	GetNumReturnData     uint64
	GetReturnDataSize    uint64
}

type EthAPICost struct {
	UseGas              uint64
	GetAddress          uint64
	GetExternalBalance  uint64
	GetBlockHash        uint64
	Call                uint64
	CallDataCopy        uint64
	GetCallDataSize     uint64
	CallCode            uint64
	CallDelegate        uint64
	CallStatic          uint64
	StorageStore        uint64
	StorageLoad         uint64
	GetCaller           uint64
	GetCallValue        uint64
	CodeCopy            uint64
	GetCodeSize         uint64
	GetBlockCoinbase    uint64
	Create              uint64
	GetBlockDifficulty  uint64
	ExternalCodeCopy    uint64
	GetExternalCodeSize uint64
	GetGasLeft          uint64
	GetBlockGasLimit    uint64
	GetTxGasPrice       uint64
	Log                 uint64
	GetBlockNumber      uint64
	GetTxOrigin         uint64
	Finish              uint64
	Revert              uint64
	GetReturnDataSize   uint64
	ReturnDataCopy      uint64
	SelfDestruct        uint64
	GetBlockTimeStamp   uint64
}

type BigIntAPICost struct {
	BigIntNew                  uint64
	BigIntUnsignedByteLength   uint64
	BigIntSignedByteLength     uint64
	BigIntGetUnsignedBytes     uint64
	BigIntGetSignedBytes       uint64
	BigIntSetUnsignedBytes     uint64
	BigIntSetSignedBytes       uint64
	BigIntIsInt64              uint64
	BigIntGetInt64             uint64
	BigIntSetInt64             uint64
	BigIntAdd                  uint64
	BigIntSub                  uint64
	BigIntMul                  uint64
	BigIntSqrt                 uint64
	BigIntPow                  uint64
	BigIntLog                  uint64
	BigIntTDiv                 uint64
	BigIntTMod                 uint64
	BigIntEDiv                 uint64
	BigIntEMod                 uint64
	BigIntAbs                  uint64
	BigIntNeg                  uint64
	BigIntSign                 uint64
	BigIntCmp                  uint64
	BigIntNot                  uint64
	BigIntAnd                  uint64
	BigIntOr                   uint64
	BigIntXor                  uint64
	BigIntShr                  uint64
	BigIntShl                  uint64
	BigIntFinishUnsigned       uint64
	BigIntFinishSigned         uint64
	BigIntStorageLoadUnsigned  uint64
	BigIntStorageStoreUnsigned uint64
	BigIntGetUnsignedArgument  uint64
	BigIntGetSignedArgument    uint64
	BigIntGetCallValue         uint64
	BigIntGetExternalBalance   uint64
}

type CryptoAPICost struct {
	SHA256                 uint64
	Keccak256              uint64
	Ripemd160              uint64
	VerifyBLS              uint64
	VerifyEd25519          uint64
	VerifySecp256k1        uint64
	EllipticCurveNew       uint64
	AddECC                 uint64
	DoubleECC              uint64
	IsOnCurveECC           uint64
	ScalarMultECC          uint64
	MarshalECC             uint64
	MarshalCompressedECC   uint64
	UnmarshalECC           uint64
	UnmarshalCompressedECC uint64
	GenerateKeyECC         uint64
}

type WASMOpcodeCost struct {
	Unreachable            uint32
	Nop                    uint32
	Block                  uint32
	Loop                   uint32
	If                     uint32
	Else                   uint32
	End                    uint32
	Br                     uint32
	BrIf                   uint32
	BrTable                uint32
	Return                 uint32
	Call                   uint32
	CallIndirect           uint32
	Drop                   uint32
	Select                 uint32
	TypedSelect            uint32
	LocalGet               uint32
	LocalSet               uint32
	LocalTee               uint32
	GlobalGet              uint32
	GlobalSet              uint32
	I32Load                uint32
	I64Load                uint32
	F32Load                uint32
	F64Load                uint32
	I32Load8S              uint32
	I32Load8U              uint32
	I32Load16S             uint32
	I32Load16U             uint32
	I64Load8S              uint32
	I64Load8U              uint32
	I64Load16S             uint32
	I64Load16U             uint32
	I64Load32S             uint32
	I64Load32U             uint32
	I32Store               uint32
	I64Store               uint32
	F32Store               uint32
	F64Store               uint32
	I32Store8              uint32
	I32Store16             uint32
	I64Store8              uint32
	I64Store16             uint32
	I64Store32             uint32
	MemorySize             uint32
	MemoryGrow             uint32
	I32Const               uint32
	I64Const               uint32
	F32Const               uint32
	F64Const               uint32
	RefNull                uint32
	RefIsNull              uint32
	RefFunc                uint32
	I32Eqz                 uint32
	I32Eq                  uint32
	I32Ne                  uint32
	I32LtS                 uint32
	I32LtU                 uint32
	I32GtS                 uint32
	I32GtU                 uint32
	I32LeS                 uint32
	I32LeU                 uint32
	I32GeS                 uint32
	I32GeU                 uint32
	I64Eqz                 uint32
	I64Eq                  uint32
	I64Ne                  uint32
	I64LtS                 uint32
	I64LtU                 uint32
	I64GtS                 uint32
	I64GtU                 uint32
	I64LeS                 uint32
	I64LeU                 uint32
	I64GeS                 uint32
	I64GeU                 uint32
	F32Eq                  uint32
	F32Ne                  uint32
	F32Lt                  uint32
	F32Gt                  uint32
	F32Le                  uint32
	F32Ge                  uint32
	F64Eq                  uint32
	F64Ne                  uint32
	F64Lt                  uint32
	F64Gt                  uint32
	F64Le                  uint32
	F64Ge                  uint32
	I32Clz                 uint32
	I32Ctz                 uint32
	I32Popcnt              uint32
	I32Add                 uint32
	I32Sub                 uint32
	I32Mul                 uint32
	I32DivS                uint32
	I32DivU                uint32
	I32RemS                uint32
	I32RemU                uint32
	I32And                 uint32
	I32Or                  uint32
	I32Xor                 uint32
	I32Shl                 uint32
	I32ShrS                uint32
	I32ShrU                uint32
	I32Rotl                uint32
	I32Rotr                uint32
	I64Clz                 uint32
	I64Ctz                 uint32
	I64Popcnt              uint32
	I64Add                 uint32
	I64Sub                 uint32
	I64Mul                 uint32
	I64DivS                uint32
	I64DivU                uint32
	I64RemS                uint32
	I64RemU                uint32
	I64And                 uint32
	I64Or                  uint32
	I64Xor                 uint32
	I64Shl                 uint32
	I64ShrS                uint32
	I64ShrU                uint32
	I64Rotl                uint32
	I64Rotr                uint32
	F32Abs                 uint32
	F32Neg                 uint32
	F32Ceil                uint32
	F32Floor               uint32
	F32Trunc               uint32
	F32Nearest             uint32
	F32Sqrt                uint32
	F32Add                 uint32
	F32Sub                 uint32
	F32Mul                 uint32
	F32Div                 uint32
	F32Min                 uint32
	F32Max                 uint32
	F32Copysign            uint32
	F64Abs                 uint32
	F64Neg                 uint32
	F64Ceil                uint32
	F64Floor               uint32
	F64Trunc               uint32
	F64Nearest             uint32
	F64Sqrt                uint32
	F64Add                 uint32
	F64Sub                 uint32
	F64Mul                 uint32
	F64Div                 uint32
	F64Min                 uint32
	F64Max                 uint32
	F64Copysign            uint32
	I32WrapI64             uint32
	I32TruncF32S           uint32
	I32TruncF32U           uint32
	I32TruncF64S           uint32
	I32TruncF64U           uint32
	I64ExtendI32S          uint32
	I64ExtendI32U          uint32
	I64TruncF32S           uint32
	I64TruncF32U           uint32
	I64TruncF64S           uint32
	I64TruncF64U           uint32
	F32ConvertI32S         uint32
	F32ConvertI32U         uint32
	F32ConvertI64S         uint32
	F32ConvertI64U         uint32
	F32DemoteF64           uint32
	F64ConvertI32S         uint32
	F64ConvertI32U         uint32
	F64ConvertI64S         uint32
	F64ConvertI64U         uint32
	F64PromoteF32          uint32
	I32ReinterpretF32      uint32
	I64ReinterpretF64      uint32
	F32ReinterpretI32      uint32
	F64ReinterpretI64      uint32
	I32Extend8S            uint32
	I32Extend16S           uint32
	I64Extend8S            uint32
	I64Extend16S           uint32
	I64Extend32S           uint32
	I32TruncSatF32S        uint32
	I32TruncSatF32U        uint32
	I32TruncSatF64S        uint32
	I32TruncSatF64U        uint32
	I64TruncSatF32S        uint32
	I64TruncSatF32U        uint32
	I64TruncSatF64S        uint32
	I64TruncSatF64U        uint32
	MemoryInit             uint32
	DataDrop               uint32
	MemoryCopy             uint32
	MemoryFill             uint32
	TableInit              uint32
	ElemDrop               uint32
	TableCopy              uint32
	TableFill              uint32
	TableGet               uint32
	TableSet               uint32
	TableGrow              uint32
	TableSize              uint32
	AtomicNotify           uint32
	I32AtomicWait          uint32
	I64AtomicWait          uint32
	AtomicFence            uint32
	I32AtomicLoad          uint32
	I64AtomicLoad          uint32
	I32AtomicLoad8U        uint32
	I32AtomicLoad16U       uint32
	I64AtomicLoad8U        uint32
	I64AtomicLoad16U       uint32
	I64AtomicLoad32U       uint32
	I32AtomicStore         uint32
	I64AtomicStore         uint32
	I32AtomicStore8        uint32
	I32AtomicStore16       uint32
	I64AtomicStore8        uint32
	I64AtomicStore16       uint32
	I64AtomicStore32       uint32
	I32AtomicRmwAdd        uint32
	I64AtomicRmwAdd        uint32
	I32AtomicRmw8AddU      uint32
	I32AtomicRmw16AddU     uint32
	I64AtomicRmw8AddU      uint32
	I64AtomicRmw16AddU     uint32
	I64AtomicRmw32AddU     uint32
	I32AtomicRmwSub        uint32
	I64AtomicRmwSub        uint32
	I32AtomicRmw8SubU      uint32
	I32AtomicRmw16SubU     uint32
	I64AtomicRmw8SubU      uint32
	I64AtomicRmw16SubU     uint32
	I64AtomicRmw32SubU     uint32
	I32AtomicRmwAnd        uint32
	I64AtomicRmwAnd        uint32
	I32AtomicRmw8AndU      uint32
	I32AtomicRmw16AndU     uint32
	I64AtomicRmw8AndU      uint32
	I64AtomicRmw16AndU     uint32
	I64AtomicRmw32AndU     uint32
	I32AtomicRmwOr         uint32
	I64AtomicRmwOr         uint32
	I32AtomicRmw8OrU       uint32
	I32AtomicRmw16OrU      uint32
	I64AtomicRmw8OrU       uint32
	I64AtomicRmw16OrU      uint32
	I64AtomicRmw32OrU      uint32
	I32AtomicRmwXor        uint32
	I64AtomicRmwXor        uint32
	I32AtomicRmw8XorU      uint32
	I32AtomicRmw16XorU     uint32
	I64AtomicRmw8XorU      uint32
	I64AtomicRmw16XorU     uint32
	I64AtomicRmw32XorU     uint32
	I32AtomicRmwXchg       uint32
	I64AtomicRmwXchg       uint32
	I32AtomicRmw8XchgU     uint32
	I32AtomicRmw16XchgU    uint32
	I64AtomicRmw8XchgU     uint32
	I64AtomicRmw16XchgU    uint32
	I64AtomicRmw32XchgU    uint32
	I32AtomicRmwCmpxchg    uint32
	I64AtomicRmwCmpxchg    uint32
	I32AtomicRmw8CmpxchgU  uint32
	I32AtomicRmw16CmpxchgU uint32
	I64AtomicRmw8CmpxchgU  uint32
	I64AtomicRmw16CmpxchgU uint32
	I64AtomicRmw32CmpxchgU uint32
	V128Load               uint32
	V128Store              uint32
	V128Const              uint32
	I8x16Splat             uint32
	I8x16ExtractLaneS      uint32
	I8x16ExtractLaneU      uint32
	I8x16ReplaceLane       uint32
	I16x8Splat             uint32
	I16x8ExtractLaneS      uint32
	I16x8ExtractLaneU      uint32
	I16x8ReplaceLane       uint32
	I32x4Splat             uint32
	I32x4ExtractLane       uint32
	I32x4ReplaceLane       uint32
	I64x2Splat             uint32
	I64x2ExtractLane       uint32
	I64x2ReplaceLane       uint32
	F32x4Splat             uint32
	F32x4ExtractLane       uint32
	F32x4ReplaceLane       uint32
	F64x2Splat             uint32
	F64x2ExtractLane       uint32
	F64x2ReplaceLane       uint32
	I8x16Eq                uint32
	I8x16Ne                uint32
	I8x16LtS               uint32
	I8x16LtU               uint32
	I8x16GtS               uint32
	I8x16GtU               uint32
	I8x16LeS               uint32
	I8x16LeU               uint32
	I8x16GeS               uint32
	I8x16GeU               uint32
	I16x8Eq                uint32
	I16x8Ne                uint32
	I16x8LtS               uint32
	I16x8LtU               uint32
	I16x8GtS               uint32
	I16x8GtU               uint32
	I16x8LeS               uint32
	I16x8LeU               uint32
	I16x8GeS               uint32
	I16x8GeU               uint32
	I32x4Eq                uint32
	I32x4Ne                uint32
	I32x4LtS               uint32
	I32x4LtU               uint32
	I32x4GtS               uint32
	I32x4GtU               uint32
	I32x4LeS               uint32
	I32x4LeU               uint32
	I32x4GeS               uint32
	I32x4GeU               uint32
	F32x4Eq                uint32
	F32x4Ne                uint32
	F32x4Lt                uint32
	F32x4Gt                uint32
	F32x4Le                uint32
	F32x4Ge                uint32
	F64x2Eq                uint32
	F64x2Ne                uint32
	F64x2Lt                uint32
	F64x2Gt                uint32
	F64x2Le                uint32
	F64x2Ge                uint32
	V128Not                uint32
	V128And                uint32
	V128AndNot             uint32
	V128Or                 uint32
	V128Xor                uint32
	V128Bitselect          uint32
	I8x16Neg               uint32
	I8x16AnyTrue           uint32
	I8x16AllTrue           uint32
	I8x16Shl               uint32
	I8x16ShrS              uint32
	I8x16ShrU              uint32
	I8x16Add               uint32
	I8x16AddSaturateS      uint32
	I8x16AddSaturateU      uint32
	I8x16Sub               uint32
	I8x16SubSaturateS      uint32
	I8x16SubSaturateU      uint32
	I8x16MinS              uint32
	I8x16MinU              uint32
	I8x16MaxS              uint32
	I8x16MaxU              uint32
	I8x16Mul               uint32
	I16x8Neg               uint32
	I16x8AnyTrue           uint32
	I16x8AllTrue           uint32
	I16x8Shl               uint32
	I16x8ShrS              uint32
	I16x8ShrU              uint32
	I16x8Add               uint32
	I16x8AddSaturateS      uint32
	I16x8AddSaturateU      uint32
	I16x8Sub               uint32
	I16x8SubSaturateS      uint32
	I16x8SubSaturateU      uint32
	I16x8Mul               uint32
	I16x8MinS              uint32
	I16x8MinU              uint32
	I16x8MaxS              uint32
	I16x8MaxU              uint32
	I32x4Neg               uint32
	I32x4AnyTrue           uint32
	I32x4AllTrue           uint32
	I32x4Shl               uint32
	I32x4ShrS              uint32
	I32x4ShrU              uint32
	I32x4Add               uint32
	I32x4Sub               uint32
	I32x4Mul               uint32
	I32x4MinS              uint32
	I32x4MinU              uint32
	I32x4MaxS              uint32
	I32x4MaxU              uint32
	I64x2Neg               uint32
	I64x2AnyTrue           uint32
	I64x2AllTrue           uint32
	I64x2Shl               uint32
	I64x2ShrS              uint32
	I64x2ShrU              uint32
	I64x2Add               uint32
	I64x2Sub               uint32
	I64x2Mul               uint32
	F32x4Abs               uint32
	F32x4Neg               uint32
	F32x4Sqrt              uint32
	F32x4Add               uint32
	F32x4Sub               uint32
	F32x4Mul               uint32
	F32x4Div               uint32
	F32x4Min               uint32
	F32x4Max               uint32
	F64x2Abs               uint32
	F64x2Neg               uint32
	F64x2Sqrt              uint32
	F64x2Add               uint32
	F64x2Sub               uint32
	F64x2Mul               uint32
	F64x2Div               uint32
	F64x2Min               uint32
	F64x2Max               uint32
	I32x4TruncSatF32x4S    uint32
	I32x4TruncSatF32x4U    uint32
	I64x2TruncSatF64x2S    uint32
	I64x2TruncSatF64x2U    uint32
	F32x4ConvertI32x4S     uint32
	F32x4ConvertI32x4U     uint32
	F64x2ConvertI64x2S     uint32
	F64x2ConvertI64x2U     uint32
	V8x16Swizzle           uint32
	V8x16Shuffle           uint32
	V8x16LoadSplat         uint32
	V16x8LoadSplat         uint32
	V32x4LoadSplat         uint32
	V64x2LoadSplat         uint32
	I8x16NarrowI16x8S      uint32
	I8x16NarrowI16x8U      uint32
	I16x8NarrowI32x4S      uint32
	I16x8NarrowI32x4U      uint32
	I16x8WidenLowI8x16S    uint32
	I16x8WidenHighI8x16S   uint32
	I16x8WidenLowI8x16U    uint32
	I16x8WidenHighI8x16U   uint32
	I32x4WidenLowI16x8S    uint32
	I32x4WidenHighI16x8S   uint32
	I32x4WidenLowI16x8U    uint32
	I32x4WidenHighI16x8U   uint32
	I16x8Load8x8S          uint32
	I16x8Load8x8U          uint32
	I32x4Load16x4S         uint32
	I32x4Load16x4U         uint32
	I64x2Load32x2S         uint32
	I64x2Load32x2U         uint32
	I8x16RoundingAverageU  uint32
	I16x8RoundingAverageU  uint32
	LocalAllocate          uint32
	LocalsUnmetered        uint32
}

func (opcode_costs_struct *WASMOpcodeCost) ToOpcodeCostsArray() [wasmer.OPCODE_COUNT]uint32 {
	opcode_costs := [wasmer.OPCODE_COUNT]uint32{}

	opcode_costs[wasmer.OpcodeUnreachable] = opcode_costs_struct.Unreachable
	opcode_costs[wasmer.OpcodeNop] = opcode_costs_struct.Nop
	opcode_costs[wasmer.OpcodeBlock] = opcode_costs_struct.Block
	opcode_costs[wasmer.OpcodeLoop] = opcode_costs_struct.Loop
	opcode_costs[wasmer.OpcodeIf] = opcode_costs_struct.If
	opcode_costs[wasmer.OpcodeElse] = opcode_costs_struct.Else
	opcode_costs[wasmer.OpcodeEnd] = opcode_costs_struct.End
	opcode_costs[wasmer.OpcodeBr] = opcode_costs_struct.Br
	opcode_costs[wasmer.OpcodeBrIf] = opcode_costs_struct.BrIf
	opcode_costs[wasmer.OpcodeBrTable] = opcode_costs_struct.BrTable
	opcode_costs[wasmer.OpcodeReturn] = opcode_costs_struct.Return
	opcode_costs[wasmer.OpcodeCall] = opcode_costs_struct.Call
	opcode_costs[wasmer.OpcodeCallIndirect] = opcode_costs_struct.CallIndirect
	opcode_costs[wasmer.OpcodeDrop] = opcode_costs_struct.Drop
	opcode_costs[wasmer.OpcodeSelect] = opcode_costs_struct.Select
	opcode_costs[wasmer.OpcodeTypedSelect] = opcode_costs_struct.TypedSelect
	opcode_costs[wasmer.OpcodeLocalGet] = opcode_costs_struct.LocalGet
	opcode_costs[wasmer.OpcodeLocalSet] = opcode_costs_struct.LocalSet
	opcode_costs[wasmer.OpcodeLocalTee] = opcode_costs_struct.LocalTee
	opcode_costs[wasmer.OpcodeGlobalGet] = opcode_costs_struct.GlobalGet
	opcode_costs[wasmer.OpcodeGlobalSet] = opcode_costs_struct.GlobalSet
	opcode_costs[wasmer.OpcodeI32Load] = opcode_costs_struct.I32Load
	opcode_costs[wasmer.OpcodeI64Load] = opcode_costs_struct.I64Load
	opcode_costs[wasmer.OpcodeF32Load] = opcode_costs_struct.F32Load
	opcode_costs[wasmer.OpcodeF64Load] = opcode_costs_struct.F64Load
	opcode_costs[wasmer.OpcodeI32Load8S] = opcode_costs_struct.I32Load8S
	opcode_costs[wasmer.OpcodeI32Load8U] = opcode_costs_struct.I32Load8U
	opcode_costs[wasmer.OpcodeI32Load16S] = opcode_costs_struct.I32Load16S
	opcode_costs[wasmer.OpcodeI32Load16U] = opcode_costs_struct.I32Load16U
	opcode_costs[wasmer.OpcodeI64Load8S] = opcode_costs_struct.I64Load8S
	opcode_costs[wasmer.OpcodeI64Load8U] = opcode_costs_struct.I64Load8U
	opcode_costs[wasmer.OpcodeI64Load16S] = opcode_costs_struct.I64Load16S
	opcode_costs[wasmer.OpcodeI64Load16U] = opcode_costs_struct.I64Load16U
	opcode_costs[wasmer.OpcodeI64Load32S] = opcode_costs_struct.I64Load32S
	opcode_costs[wasmer.OpcodeI64Load32U] = opcode_costs_struct.I64Load32U
	opcode_costs[wasmer.OpcodeI32Store] = opcode_costs_struct.I32Store
	opcode_costs[wasmer.OpcodeI64Store] = opcode_costs_struct.I64Store
	opcode_costs[wasmer.OpcodeF32Store] = opcode_costs_struct.F32Store
	opcode_costs[wasmer.OpcodeF64Store] = opcode_costs_struct.F64Store
	opcode_costs[wasmer.OpcodeI32Store8] = opcode_costs_struct.I32Store8
	opcode_costs[wasmer.OpcodeI32Store16] = opcode_costs_struct.I32Store16
	opcode_costs[wasmer.OpcodeI64Store8] = opcode_costs_struct.I64Store8
	opcode_costs[wasmer.OpcodeI64Store16] = opcode_costs_struct.I64Store16
	opcode_costs[wasmer.OpcodeI64Store32] = opcode_costs_struct.I64Store32
	opcode_costs[wasmer.OpcodeMemorySize] = opcode_costs_struct.MemorySize
	opcode_costs[wasmer.OpcodeMemoryGrow] = opcode_costs_struct.MemoryGrow
	opcode_costs[wasmer.OpcodeI32Const] = opcode_costs_struct.I32Const
	opcode_costs[wasmer.OpcodeI64Const] = opcode_costs_struct.I64Const
	opcode_costs[wasmer.OpcodeF32Const] = opcode_costs_struct.F32Const
	opcode_costs[wasmer.OpcodeF64Const] = opcode_costs_struct.F64Const
	opcode_costs[wasmer.OpcodeRefNull] = opcode_costs_struct.RefNull
	opcode_costs[wasmer.OpcodeRefIsNull] = opcode_costs_struct.RefIsNull
	opcode_costs[wasmer.OpcodeRefFunc] = opcode_costs_struct.RefFunc
	opcode_costs[wasmer.OpcodeI32Eqz] = opcode_costs_struct.I32Eqz
	opcode_costs[wasmer.OpcodeI32Eq] = opcode_costs_struct.I32Eq
	opcode_costs[wasmer.OpcodeI32Ne] = opcode_costs_struct.I32Ne
	opcode_costs[wasmer.OpcodeI32LtS] = opcode_costs_struct.I32LtS
	opcode_costs[wasmer.OpcodeI32LtU] = opcode_costs_struct.I32LtU
	opcode_costs[wasmer.OpcodeI32GtS] = opcode_costs_struct.I32GtS
	opcode_costs[wasmer.OpcodeI32GtU] = opcode_costs_struct.I32GtU
	opcode_costs[wasmer.OpcodeI32LeS] = opcode_costs_struct.I32LeS
	opcode_costs[wasmer.OpcodeI32LeU] = opcode_costs_struct.I32LeU
	opcode_costs[wasmer.OpcodeI32GeS] = opcode_costs_struct.I32GeS
	opcode_costs[wasmer.OpcodeI32GeU] = opcode_costs_struct.I32GeU
	opcode_costs[wasmer.OpcodeI64Eqz] = opcode_costs_struct.I64Eqz
	opcode_costs[wasmer.OpcodeI64Eq] = opcode_costs_struct.I64Eq
	opcode_costs[wasmer.OpcodeI64Ne] = opcode_costs_struct.I64Ne
	opcode_costs[wasmer.OpcodeI64LtS] = opcode_costs_struct.I64LtS
	opcode_costs[wasmer.OpcodeI64LtU] = opcode_costs_struct.I64LtU
	opcode_costs[wasmer.OpcodeI64GtS] = opcode_costs_struct.I64GtS
	opcode_costs[wasmer.OpcodeI64GtU] = opcode_costs_struct.I64GtU
	opcode_costs[wasmer.OpcodeI64LeS] = opcode_costs_struct.I64LeS
	opcode_costs[wasmer.OpcodeI64LeU] = opcode_costs_struct.I64LeU
	opcode_costs[wasmer.OpcodeI64GeS] = opcode_costs_struct.I64GeS
	opcode_costs[wasmer.OpcodeI64GeU] = opcode_costs_struct.I64GeU
	opcode_costs[wasmer.OpcodeF32Eq] = opcode_costs_struct.F32Eq
	opcode_costs[wasmer.OpcodeF32Ne] = opcode_costs_struct.F32Ne
	opcode_costs[wasmer.OpcodeF32Lt] = opcode_costs_struct.F32Lt
	opcode_costs[wasmer.OpcodeF32Gt] = opcode_costs_struct.F32Gt
	opcode_costs[wasmer.OpcodeF32Le] = opcode_costs_struct.F32Le
	opcode_costs[wasmer.OpcodeF32Ge] = opcode_costs_struct.F32Ge
	opcode_costs[wasmer.OpcodeF64Eq] = opcode_costs_struct.F64Eq
	opcode_costs[wasmer.OpcodeF64Ne] = opcode_costs_struct.F64Ne
	opcode_costs[wasmer.OpcodeF64Lt] = opcode_costs_struct.F64Lt
	opcode_costs[wasmer.OpcodeF64Gt] = opcode_costs_struct.F64Gt
	opcode_costs[wasmer.OpcodeF64Le] = opcode_costs_struct.F64Le
	opcode_costs[wasmer.OpcodeF64Ge] = opcode_costs_struct.F64Ge
	opcode_costs[wasmer.OpcodeI32Clz] = opcode_costs_struct.I32Clz
	opcode_costs[wasmer.OpcodeI32Ctz] = opcode_costs_struct.I32Ctz
	opcode_costs[wasmer.OpcodeI32Popcnt] = opcode_costs_struct.I32Popcnt
	opcode_costs[wasmer.OpcodeI32Add] = opcode_costs_struct.I32Add
	opcode_costs[wasmer.OpcodeI32Sub] = opcode_costs_struct.I32Sub
	opcode_costs[wasmer.OpcodeI32Mul] = opcode_costs_struct.I32Mul
	opcode_costs[wasmer.OpcodeI32DivS] = opcode_costs_struct.I32DivS
	opcode_costs[wasmer.OpcodeI32DivU] = opcode_costs_struct.I32DivU
	opcode_costs[wasmer.OpcodeI32RemS] = opcode_costs_struct.I32RemS
	opcode_costs[wasmer.OpcodeI32RemU] = opcode_costs_struct.I32RemU
	opcode_costs[wasmer.OpcodeI32And] = opcode_costs_struct.I32And
	opcode_costs[wasmer.OpcodeI32Or] = opcode_costs_struct.I32Or
	opcode_costs[wasmer.OpcodeI32Xor] = opcode_costs_struct.I32Xor
	opcode_costs[wasmer.OpcodeI32Shl] = opcode_costs_struct.I32Shl
	opcode_costs[wasmer.OpcodeI32ShrS] = opcode_costs_struct.I32ShrS
	opcode_costs[wasmer.OpcodeI32ShrU] = opcode_costs_struct.I32ShrU
	opcode_costs[wasmer.OpcodeI32Rotl] = opcode_costs_struct.I32Rotl
	opcode_costs[wasmer.OpcodeI32Rotr] = opcode_costs_struct.I32Rotr
	opcode_costs[wasmer.OpcodeI64Clz] = opcode_costs_struct.I64Clz
	opcode_costs[wasmer.OpcodeI64Ctz] = opcode_costs_struct.I64Ctz
	opcode_costs[wasmer.OpcodeI64Popcnt] = opcode_costs_struct.I64Popcnt
	opcode_costs[wasmer.OpcodeI64Add] = opcode_costs_struct.I64Add
	opcode_costs[wasmer.OpcodeI64Sub] = opcode_costs_struct.I64Sub
	opcode_costs[wasmer.OpcodeI64Mul] = opcode_costs_struct.I64Mul
	opcode_costs[wasmer.OpcodeI64DivS] = opcode_costs_struct.I64DivS
	opcode_costs[wasmer.OpcodeI64DivU] = opcode_costs_struct.I64DivU
	opcode_costs[wasmer.OpcodeI64RemS] = opcode_costs_struct.I64RemS
	opcode_costs[wasmer.OpcodeI64RemU] = opcode_costs_struct.I64RemU
	opcode_costs[wasmer.OpcodeI64And] = opcode_costs_struct.I64And
	opcode_costs[wasmer.OpcodeI64Or] = opcode_costs_struct.I64Or
	opcode_costs[wasmer.OpcodeI64Xor] = opcode_costs_struct.I64Xor
	opcode_costs[wasmer.OpcodeI64Shl] = opcode_costs_struct.I64Shl
	opcode_costs[wasmer.OpcodeI64ShrS] = opcode_costs_struct.I64ShrS
	opcode_costs[wasmer.OpcodeI64ShrU] = opcode_costs_struct.I64ShrU
	opcode_costs[wasmer.OpcodeI64Rotl] = opcode_costs_struct.I64Rotl
	opcode_costs[wasmer.OpcodeI64Rotr] = opcode_costs_struct.I64Rotr
	opcode_costs[wasmer.OpcodeF32Abs] = opcode_costs_struct.F32Abs
	opcode_costs[wasmer.OpcodeF32Neg] = opcode_costs_struct.F32Neg
	opcode_costs[wasmer.OpcodeF32Ceil] = opcode_costs_struct.F32Ceil
	opcode_costs[wasmer.OpcodeF32Floor] = opcode_costs_struct.F32Floor
	opcode_costs[wasmer.OpcodeF32Trunc] = opcode_costs_struct.F32Trunc
	opcode_costs[wasmer.OpcodeF32Nearest] = opcode_costs_struct.F32Nearest
	opcode_costs[wasmer.OpcodeF32Sqrt] = opcode_costs_struct.F32Sqrt
	opcode_costs[wasmer.OpcodeF32Add] = opcode_costs_struct.F32Add
	opcode_costs[wasmer.OpcodeF32Sub] = opcode_costs_struct.F32Sub
	opcode_costs[wasmer.OpcodeF32Mul] = opcode_costs_struct.F32Mul
	opcode_costs[wasmer.OpcodeF32Div] = opcode_costs_struct.F32Div
	opcode_costs[wasmer.OpcodeF32Min] = opcode_costs_struct.F32Min
	opcode_costs[wasmer.OpcodeF32Max] = opcode_costs_struct.F32Max
	opcode_costs[wasmer.OpcodeF32Copysign] = opcode_costs_struct.F32Copysign
	opcode_costs[wasmer.OpcodeF64Abs] = opcode_costs_struct.F64Abs
	opcode_costs[wasmer.OpcodeF64Neg] = opcode_costs_struct.F64Neg
	opcode_costs[wasmer.OpcodeF64Ceil] = opcode_costs_struct.F64Ceil
	opcode_costs[wasmer.OpcodeF64Floor] = opcode_costs_struct.F64Floor
	opcode_costs[wasmer.OpcodeF64Trunc] = opcode_costs_struct.F64Trunc
	opcode_costs[wasmer.OpcodeF64Nearest] = opcode_costs_struct.F64Nearest
	opcode_costs[wasmer.OpcodeF64Sqrt] = opcode_costs_struct.F64Sqrt
	opcode_costs[wasmer.OpcodeF64Add] = opcode_costs_struct.F64Add
	opcode_costs[wasmer.OpcodeF64Sub] = opcode_costs_struct.F64Sub
	opcode_costs[wasmer.OpcodeF64Mul] = opcode_costs_struct.F64Mul
	opcode_costs[wasmer.OpcodeF64Div] = opcode_costs_struct.F64Div
	opcode_costs[wasmer.OpcodeF64Min] = opcode_costs_struct.F64Min
	opcode_costs[wasmer.OpcodeF64Max] = opcode_costs_struct.F64Max
	opcode_costs[wasmer.OpcodeF64Copysign] = opcode_costs_struct.F64Copysign
	opcode_costs[wasmer.OpcodeI32WrapI64] = opcode_costs_struct.I32WrapI64
	opcode_costs[wasmer.OpcodeI32TruncF32S] = opcode_costs_struct.I32TruncF32S
	opcode_costs[wasmer.OpcodeI32TruncF32U] = opcode_costs_struct.I32TruncF32U
	opcode_costs[wasmer.OpcodeI32TruncF64S] = opcode_costs_struct.I32TruncF64S
	opcode_costs[wasmer.OpcodeI32TruncF64U] = opcode_costs_struct.I32TruncF64U
	opcode_costs[wasmer.OpcodeI64ExtendI32S] = opcode_costs_struct.I64ExtendI32S
	opcode_costs[wasmer.OpcodeI64ExtendI32U] = opcode_costs_struct.I64ExtendI32U
	opcode_costs[wasmer.OpcodeI64TruncF32S] = opcode_costs_struct.I64TruncF32S
	opcode_costs[wasmer.OpcodeI64TruncF32U] = opcode_costs_struct.I64TruncF32U
	opcode_costs[wasmer.OpcodeI64TruncF64S] = opcode_costs_struct.I64TruncF64S
	opcode_costs[wasmer.OpcodeI64TruncF64U] = opcode_costs_struct.I64TruncF64U
	opcode_costs[wasmer.OpcodeF32ConvertI32S] = opcode_costs_struct.F32ConvertI32S
	opcode_costs[wasmer.OpcodeF32ConvertI32U] = opcode_costs_struct.F32ConvertI32U
	opcode_costs[wasmer.OpcodeF32ConvertI64S] = opcode_costs_struct.F32ConvertI64S
	opcode_costs[wasmer.OpcodeF32ConvertI64U] = opcode_costs_struct.F32ConvertI64U
	opcode_costs[wasmer.OpcodeF32DemoteF64] = opcode_costs_struct.F32DemoteF64
	opcode_costs[wasmer.OpcodeF64ConvertI32S] = opcode_costs_struct.F64ConvertI32S
	opcode_costs[wasmer.OpcodeF64ConvertI32U] = opcode_costs_struct.F64ConvertI32U
	opcode_costs[wasmer.OpcodeF64ConvertI64S] = opcode_costs_struct.F64ConvertI64S
	opcode_costs[wasmer.OpcodeF64ConvertI64U] = opcode_costs_struct.F64ConvertI64U
	opcode_costs[wasmer.OpcodeF64PromoteF32] = opcode_costs_struct.F64PromoteF32
	opcode_costs[wasmer.OpcodeI32ReinterpretF32] = opcode_costs_struct.I32ReinterpretF32
	opcode_costs[wasmer.OpcodeI64ReinterpretF64] = opcode_costs_struct.I64ReinterpretF64
	opcode_costs[wasmer.OpcodeF32ReinterpretI32] = opcode_costs_struct.F32ReinterpretI32
	opcode_costs[wasmer.OpcodeF64ReinterpretI64] = opcode_costs_struct.F64ReinterpretI64
	opcode_costs[wasmer.OpcodeI32Extend8S] = opcode_costs_struct.I32Extend8S
	opcode_costs[wasmer.OpcodeI32Extend16S] = opcode_costs_struct.I32Extend16S
	opcode_costs[wasmer.OpcodeI64Extend8S] = opcode_costs_struct.I64Extend8S
	opcode_costs[wasmer.OpcodeI64Extend16S] = opcode_costs_struct.I64Extend16S
	opcode_costs[wasmer.OpcodeI64Extend32S] = opcode_costs_struct.I64Extend32S
	opcode_costs[wasmer.OpcodeI32TruncSatF32S] = opcode_costs_struct.I32TruncSatF32S
	opcode_costs[wasmer.OpcodeI32TruncSatF32U] = opcode_costs_struct.I32TruncSatF32U
	opcode_costs[wasmer.OpcodeI32TruncSatF64S] = opcode_costs_struct.I32TruncSatF64S
	opcode_costs[wasmer.OpcodeI32TruncSatF64U] = opcode_costs_struct.I32TruncSatF64U
	opcode_costs[wasmer.OpcodeI64TruncSatF32S] = opcode_costs_struct.I64TruncSatF32S
	opcode_costs[wasmer.OpcodeI64TruncSatF32U] = opcode_costs_struct.I64TruncSatF32U
	opcode_costs[wasmer.OpcodeI64TruncSatF64S] = opcode_costs_struct.I64TruncSatF64S
	opcode_costs[wasmer.OpcodeI64TruncSatF64U] = opcode_costs_struct.I64TruncSatF64U
	opcode_costs[wasmer.OpcodeMemoryInit] = opcode_costs_struct.MemoryInit
	opcode_costs[wasmer.OpcodeDataDrop] = opcode_costs_struct.DataDrop
	opcode_costs[wasmer.OpcodeMemoryCopy] = opcode_costs_struct.MemoryCopy
	opcode_costs[wasmer.OpcodeMemoryFill] = opcode_costs_struct.MemoryFill
	opcode_costs[wasmer.OpcodeTableInit] = opcode_costs_struct.TableInit
	opcode_costs[wasmer.OpcodeElemDrop] = opcode_costs_struct.ElemDrop
	opcode_costs[wasmer.OpcodeTableCopy] = opcode_costs_struct.TableCopy
	opcode_costs[wasmer.OpcodeTableFill] = opcode_costs_struct.TableFill
	opcode_costs[wasmer.OpcodeTableGet] = opcode_costs_struct.TableGet
	opcode_costs[wasmer.OpcodeTableSet] = opcode_costs_struct.TableSet
	opcode_costs[wasmer.OpcodeTableGrow] = opcode_costs_struct.TableGrow
	opcode_costs[wasmer.OpcodeTableSize] = opcode_costs_struct.TableSize
	opcode_costs[wasmer.OpcodeAtomicNotify] = opcode_costs_struct.AtomicNotify
	opcode_costs[wasmer.OpcodeI32AtomicWait] = opcode_costs_struct.I32AtomicWait
	opcode_costs[wasmer.OpcodeI64AtomicWait] = opcode_costs_struct.I64AtomicWait
	opcode_costs[wasmer.OpcodeAtomicFence] = opcode_costs_struct.AtomicFence
	opcode_costs[wasmer.OpcodeI32AtomicLoad] = opcode_costs_struct.I32AtomicLoad
	opcode_costs[wasmer.OpcodeI64AtomicLoad] = opcode_costs_struct.I64AtomicLoad
	opcode_costs[wasmer.OpcodeI32AtomicLoad8U] = opcode_costs_struct.I32AtomicLoad8U
	opcode_costs[wasmer.OpcodeI32AtomicLoad16U] = opcode_costs_struct.I32AtomicLoad16U
	opcode_costs[wasmer.OpcodeI64AtomicLoad8U] = opcode_costs_struct.I64AtomicLoad8U
	opcode_costs[wasmer.OpcodeI64AtomicLoad16U] = opcode_costs_struct.I64AtomicLoad16U
	opcode_costs[wasmer.OpcodeI64AtomicLoad32U] = opcode_costs_struct.I64AtomicLoad32U
	opcode_costs[wasmer.OpcodeI32AtomicStore] = opcode_costs_struct.I32AtomicStore
	opcode_costs[wasmer.OpcodeI64AtomicStore] = opcode_costs_struct.I64AtomicStore
	opcode_costs[wasmer.OpcodeI32AtomicStore8] = opcode_costs_struct.I32AtomicStore8
	opcode_costs[wasmer.OpcodeI32AtomicStore16] = opcode_costs_struct.I32AtomicStore16
	opcode_costs[wasmer.OpcodeI64AtomicStore8] = opcode_costs_struct.I64AtomicStore8
	opcode_costs[wasmer.OpcodeI64AtomicStore16] = opcode_costs_struct.I64AtomicStore16
	opcode_costs[wasmer.OpcodeI64AtomicStore32] = opcode_costs_struct.I64AtomicStore32
	opcode_costs[wasmer.OpcodeI32AtomicRmwAdd] = opcode_costs_struct.I32AtomicRmwAdd
	opcode_costs[wasmer.OpcodeI64AtomicRmwAdd] = opcode_costs_struct.I64AtomicRmwAdd
	opcode_costs[wasmer.OpcodeI32AtomicRmw8AddU] = opcode_costs_struct.I32AtomicRmw8AddU
	opcode_costs[wasmer.OpcodeI32AtomicRmw16AddU] = opcode_costs_struct.I32AtomicRmw16AddU
	opcode_costs[wasmer.OpcodeI64AtomicRmw8AddU] = opcode_costs_struct.I64AtomicRmw8AddU
	opcode_costs[wasmer.OpcodeI64AtomicRmw16AddU] = opcode_costs_struct.I64AtomicRmw16AddU
	opcode_costs[wasmer.OpcodeI64AtomicRmw32AddU] = opcode_costs_struct.I64AtomicRmw32AddU
	opcode_costs[wasmer.OpcodeI32AtomicRmwSub] = opcode_costs_struct.I32AtomicRmwSub
	opcode_costs[wasmer.OpcodeI64AtomicRmwSub] = opcode_costs_struct.I64AtomicRmwSub
	opcode_costs[wasmer.OpcodeI32AtomicRmw8SubU] = opcode_costs_struct.I32AtomicRmw8SubU
	opcode_costs[wasmer.OpcodeI32AtomicRmw16SubU] = opcode_costs_struct.I32AtomicRmw16SubU
	opcode_costs[wasmer.OpcodeI64AtomicRmw8SubU] = opcode_costs_struct.I64AtomicRmw8SubU
	opcode_costs[wasmer.OpcodeI64AtomicRmw16SubU] = opcode_costs_struct.I64AtomicRmw16SubU
	opcode_costs[wasmer.OpcodeI64AtomicRmw32SubU] = opcode_costs_struct.I64AtomicRmw32SubU
	opcode_costs[wasmer.OpcodeI32AtomicRmwAnd] = opcode_costs_struct.I32AtomicRmwAnd
	opcode_costs[wasmer.OpcodeI64AtomicRmwAnd] = opcode_costs_struct.I64AtomicRmwAnd
	opcode_costs[wasmer.OpcodeI32AtomicRmw8AndU] = opcode_costs_struct.I32AtomicRmw8AndU
	opcode_costs[wasmer.OpcodeI32AtomicRmw16AndU] = opcode_costs_struct.I32AtomicRmw16AndU
	opcode_costs[wasmer.OpcodeI64AtomicRmw8AndU] = opcode_costs_struct.I64AtomicRmw8AndU
	opcode_costs[wasmer.OpcodeI64AtomicRmw16AndU] = opcode_costs_struct.I64AtomicRmw16AndU
	opcode_costs[wasmer.OpcodeI64AtomicRmw32AndU] = opcode_costs_struct.I64AtomicRmw32AndU
	opcode_costs[wasmer.OpcodeI32AtomicRmwOr] = opcode_costs_struct.I32AtomicRmwOr
	opcode_costs[wasmer.OpcodeI64AtomicRmwOr] = opcode_costs_struct.I64AtomicRmwOr
	opcode_costs[wasmer.OpcodeI32AtomicRmw8OrU] = opcode_costs_struct.I32AtomicRmw8OrU
	opcode_costs[wasmer.OpcodeI32AtomicRmw16OrU] = opcode_costs_struct.I32AtomicRmw16OrU
	opcode_costs[wasmer.OpcodeI64AtomicRmw8OrU] = opcode_costs_struct.I64AtomicRmw8OrU
	opcode_costs[wasmer.OpcodeI64AtomicRmw16OrU] = opcode_costs_struct.I64AtomicRmw16OrU
	opcode_costs[wasmer.OpcodeI64AtomicRmw32OrU] = opcode_costs_struct.I64AtomicRmw32OrU
	opcode_costs[wasmer.OpcodeI32AtomicRmwXor] = opcode_costs_struct.I32AtomicRmwXor
	opcode_costs[wasmer.OpcodeI64AtomicRmwXor] = opcode_costs_struct.I64AtomicRmwXor
	opcode_costs[wasmer.OpcodeI32AtomicRmw8XorU] = opcode_costs_struct.I32AtomicRmw8XorU
	opcode_costs[wasmer.OpcodeI32AtomicRmw16XorU] = opcode_costs_struct.I32AtomicRmw16XorU
	opcode_costs[wasmer.OpcodeI64AtomicRmw8XorU] = opcode_costs_struct.I64AtomicRmw8XorU
	opcode_costs[wasmer.OpcodeI64AtomicRmw16XorU] = opcode_costs_struct.I64AtomicRmw16XorU
	opcode_costs[wasmer.OpcodeI64AtomicRmw32XorU] = opcode_costs_struct.I64AtomicRmw32XorU
	opcode_costs[wasmer.OpcodeI32AtomicRmwXchg] = opcode_costs_struct.I32AtomicRmwXchg
	opcode_costs[wasmer.OpcodeI64AtomicRmwXchg] = opcode_costs_struct.I64AtomicRmwXchg
	opcode_costs[wasmer.OpcodeI32AtomicRmw8XchgU] = opcode_costs_struct.I32AtomicRmw8XchgU
	opcode_costs[wasmer.OpcodeI32AtomicRmw16XchgU] = opcode_costs_struct.I32AtomicRmw16XchgU
	opcode_costs[wasmer.OpcodeI64AtomicRmw8XchgU] = opcode_costs_struct.I64AtomicRmw8XchgU
	opcode_costs[wasmer.OpcodeI64AtomicRmw16XchgU] = opcode_costs_struct.I64AtomicRmw16XchgU
	opcode_costs[wasmer.OpcodeI64AtomicRmw32XchgU] = opcode_costs_struct.I64AtomicRmw32XchgU
	opcode_costs[wasmer.OpcodeI32AtomicRmwCmpxchg] = opcode_costs_struct.I32AtomicRmwCmpxchg
	opcode_costs[wasmer.OpcodeI64AtomicRmwCmpxchg] = opcode_costs_struct.I64AtomicRmwCmpxchg
	opcode_costs[wasmer.OpcodeI32AtomicRmw8CmpxchgU] = opcode_costs_struct.I32AtomicRmw8CmpxchgU
	opcode_costs[wasmer.OpcodeI32AtomicRmw16CmpxchgU] = opcode_costs_struct.I32AtomicRmw16CmpxchgU
	opcode_costs[wasmer.OpcodeI64AtomicRmw8CmpxchgU] = opcode_costs_struct.I64AtomicRmw8CmpxchgU
	opcode_costs[wasmer.OpcodeI64AtomicRmw16CmpxchgU] = opcode_costs_struct.I64AtomicRmw16CmpxchgU
	opcode_costs[wasmer.OpcodeI64AtomicRmw32CmpxchgU] = opcode_costs_struct.I64AtomicRmw32CmpxchgU
	opcode_costs[wasmer.OpcodeV128Load] = opcode_costs_struct.V128Load
	opcode_costs[wasmer.OpcodeV128Store] = opcode_costs_struct.V128Store
	opcode_costs[wasmer.OpcodeV128Const] = opcode_costs_struct.V128Const
	opcode_costs[wasmer.OpcodeI8x16Splat] = opcode_costs_struct.I8x16Splat
	opcode_costs[wasmer.OpcodeI8x16ExtractLaneS] = opcode_costs_struct.I8x16ExtractLaneS
	opcode_costs[wasmer.OpcodeI8x16ExtractLaneU] = opcode_costs_struct.I8x16ExtractLaneU
	opcode_costs[wasmer.OpcodeI8x16ReplaceLane] = opcode_costs_struct.I8x16ReplaceLane
	opcode_costs[wasmer.OpcodeI16x8Splat] = opcode_costs_struct.I16x8Splat
	opcode_costs[wasmer.OpcodeI16x8ExtractLaneS] = opcode_costs_struct.I16x8ExtractLaneS
	opcode_costs[wasmer.OpcodeI16x8ExtractLaneU] = opcode_costs_struct.I16x8ExtractLaneU
	opcode_costs[wasmer.OpcodeI16x8ReplaceLane] = opcode_costs_struct.I16x8ReplaceLane
	opcode_costs[wasmer.OpcodeI32x4Splat] = opcode_costs_struct.I32x4Splat
	opcode_costs[wasmer.OpcodeI32x4ExtractLane] = opcode_costs_struct.I32x4ExtractLane
	opcode_costs[wasmer.OpcodeI32x4ReplaceLane] = opcode_costs_struct.I32x4ReplaceLane
	opcode_costs[wasmer.OpcodeI64x2Splat] = opcode_costs_struct.I64x2Splat
	opcode_costs[wasmer.OpcodeI64x2ExtractLane] = opcode_costs_struct.I64x2ExtractLane
	opcode_costs[wasmer.OpcodeI64x2ReplaceLane] = opcode_costs_struct.I64x2ReplaceLane
	opcode_costs[wasmer.OpcodeF32x4Splat] = opcode_costs_struct.F32x4Splat
	opcode_costs[wasmer.OpcodeF32x4ExtractLane] = opcode_costs_struct.F32x4ExtractLane
	opcode_costs[wasmer.OpcodeF32x4ReplaceLane] = opcode_costs_struct.F32x4ReplaceLane
	opcode_costs[wasmer.OpcodeF64x2Splat] = opcode_costs_struct.F64x2Splat
	opcode_costs[wasmer.OpcodeF64x2ExtractLane] = opcode_costs_struct.F64x2ExtractLane
	opcode_costs[wasmer.OpcodeF64x2ReplaceLane] = opcode_costs_struct.F64x2ReplaceLane
	opcode_costs[wasmer.OpcodeI8x16Eq] = opcode_costs_struct.I8x16Eq
	opcode_costs[wasmer.OpcodeI8x16Ne] = opcode_costs_struct.I8x16Ne
	opcode_costs[wasmer.OpcodeI8x16LtS] = opcode_costs_struct.I8x16LtS
	opcode_costs[wasmer.OpcodeI8x16LtU] = opcode_costs_struct.I8x16LtU
	opcode_costs[wasmer.OpcodeI8x16GtS] = opcode_costs_struct.I8x16GtS
	opcode_costs[wasmer.OpcodeI8x16GtU] = opcode_costs_struct.I8x16GtU
	opcode_costs[wasmer.OpcodeI8x16LeS] = opcode_costs_struct.I8x16LeS
	opcode_costs[wasmer.OpcodeI8x16LeU] = opcode_costs_struct.I8x16LeU
	opcode_costs[wasmer.OpcodeI8x16GeS] = opcode_costs_struct.I8x16GeS
	opcode_costs[wasmer.OpcodeI8x16GeU] = opcode_costs_struct.I8x16GeU
	opcode_costs[wasmer.OpcodeI16x8Eq] = opcode_costs_struct.I16x8Eq
	opcode_costs[wasmer.OpcodeI16x8Ne] = opcode_costs_struct.I16x8Ne
	opcode_costs[wasmer.OpcodeI16x8LtS] = opcode_costs_struct.I16x8LtS
	opcode_costs[wasmer.OpcodeI16x8LtU] = opcode_costs_struct.I16x8LtU
	opcode_costs[wasmer.OpcodeI16x8GtS] = opcode_costs_struct.I16x8GtS
	opcode_costs[wasmer.OpcodeI16x8GtU] = opcode_costs_struct.I16x8GtU
	opcode_costs[wasmer.OpcodeI16x8LeS] = opcode_costs_struct.I16x8LeS
	opcode_costs[wasmer.OpcodeI16x8LeU] = opcode_costs_struct.I16x8LeU
	opcode_costs[wasmer.OpcodeI16x8GeS] = opcode_costs_struct.I16x8GeS
	opcode_costs[wasmer.OpcodeI16x8GeU] = opcode_costs_struct.I16x8GeU
	opcode_costs[wasmer.OpcodeI32x4Eq] = opcode_costs_struct.I32x4Eq
	opcode_costs[wasmer.OpcodeI32x4Ne] = opcode_costs_struct.I32x4Ne
	opcode_costs[wasmer.OpcodeI32x4LtS] = opcode_costs_struct.I32x4LtS
	opcode_costs[wasmer.OpcodeI32x4LtU] = opcode_costs_struct.I32x4LtU
	opcode_costs[wasmer.OpcodeI32x4GtS] = opcode_costs_struct.I32x4GtS
	opcode_costs[wasmer.OpcodeI32x4GtU] = opcode_costs_struct.I32x4GtU
	opcode_costs[wasmer.OpcodeI32x4LeS] = opcode_costs_struct.I32x4LeS
	opcode_costs[wasmer.OpcodeI32x4LeU] = opcode_costs_struct.I32x4LeU
	opcode_costs[wasmer.OpcodeI32x4GeS] = opcode_costs_struct.I32x4GeS
	opcode_costs[wasmer.OpcodeI32x4GeU] = opcode_costs_struct.I32x4GeU
	opcode_costs[wasmer.OpcodeF32x4Eq] = opcode_costs_struct.F32x4Eq
	opcode_costs[wasmer.OpcodeF32x4Ne] = opcode_costs_struct.F32x4Ne
	opcode_costs[wasmer.OpcodeF32x4Lt] = opcode_costs_struct.F32x4Lt
	opcode_costs[wasmer.OpcodeF32x4Gt] = opcode_costs_struct.F32x4Gt
	opcode_costs[wasmer.OpcodeF32x4Le] = opcode_costs_struct.F32x4Le
	opcode_costs[wasmer.OpcodeF32x4Ge] = opcode_costs_struct.F32x4Ge
	opcode_costs[wasmer.OpcodeF64x2Eq] = opcode_costs_struct.F64x2Eq
	opcode_costs[wasmer.OpcodeF64x2Ne] = opcode_costs_struct.F64x2Ne
	opcode_costs[wasmer.OpcodeF64x2Lt] = opcode_costs_struct.F64x2Lt
	opcode_costs[wasmer.OpcodeF64x2Gt] = opcode_costs_struct.F64x2Gt
	opcode_costs[wasmer.OpcodeF64x2Le] = opcode_costs_struct.F64x2Le
	opcode_costs[wasmer.OpcodeF64x2Ge] = opcode_costs_struct.F64x2Ge
	opcode_costs[wasmer.OpcodeV128Not] = opcode_costs_struct.V128Not
	opcode_costs[wasmer.OpcodeV128And] = opcode_costs_struct.V128And
	opcode_costs[wasmer.OpcodeV128AndNot] = opcode_costs_struct.V128AndNot
	opcode_costs[wasmer.OpcodeV128Or] = opcode_costs_struct.V128Or
	opcode_costs[wasmer.OpcodeV128Xor] = opcode_costs_struct.V128Xor
	opcode_costs[wasmer.OpcodeV128Bitselect] = opcode_costs_struct.V128Bitselect
	opcode_costs[wasmer.OpcodeI8x16Neg] = opcode_costs_struct.I8x16Neg
	opcode_costs[wasmer.OpcodeI8x16AnyTrue] = opcode_costs_struct.I8x16AnyTrue
	opcode_costs[wasmer.OpcodeI8x16AllTrue] = opcode_costs_struct.I8x16AllTrue
	opcode_costs[wasmer.OpcodeI8x16Shl] = opcode_costs_struct.I8x16Shl
	opcode_costs[wasmer.OpcodeI8x16ShrS] = opcode_costs_struct.I8x16ShrS
	opcode_costs[wasmer.OpcodeI8x16ShrU] = opcode_costs_struct.I8x16ShrU
	opcode_costs[wasmer.OpcodeI8x16Add] = opcode_costs_struct.I8x16Add
	opcode_costs[wasmer.OpcodeI8x16AddSaturateS] = opcode_costs_struct.I8x16AddSaturateS
	opcode_costs[wasmer.OpcodeI8x16AddSaturateU] = opcode_costs_struct.I8x16AddSaturateU
	opcode_costs[wasmer.OpcodeI8x16Sub] = opcode_costs_struct.I8x16Sub
	opcode_costs[wasmer.OpcodeI8x16SubSaturateS] = opcode_costs_struct.I8x16SubSaturateS
	opcode_costs[wasmer.OpcodeI8x16SubSaturateU] = opcode_costs_struct.I8x16SubSaturateU
	opcode_costs[wasmer.OpcodeI8x16MinS] = opcode_costs_struct.I8x16MinS
	opcode_costs[wasmer.OpcodeI8x16MinU] = opcode_costs_struct.I8x16MinU
	opcode_costs[wasmer.OpcodeI8x16MaxS] = opcode_costs_struct.I8x16MaxS
	opcode_costs[wasmer.OpcodeI8x16MaxU] = opcode_costs_struct.I8x16MaxU
	opcode_costs[wasmer.OpcodeI8x16Mul] = opcode_costs_struct.I8x16Mul
	opcode_costs[wasmer.OpcodeI16x8Neg] = opcode_costs_struct.I16x8Neg
	opcode_costs[wasmer.OpcodeI16x8AnyTrue] = opcode_costs_struct.I16x8AnyTrue
	opcode_costs[wasmer.OpcodeI16x8AllTrue] = opcode_costs_struct.I16x8AllTrue
	opcode_costs[wasmer.OpcodeI16x8Shl] = opcode_costs_struct.I16x8Shl
	opcode_costs[wasmer.OpcodeI16x8ShrS] = opcode_costs_struct.I16x8ShrS
	opcode_costs[wasmer.OpcodeI16x8ShrU] = opcode_costs_struct.I16x8ShrU
	opcode_costs[wasmer.OpcodeI16x8Add] = opcode_costs_struct.I16x8Add
	opcode_costs[wasmer.OpcodeI16x8AddSaturateS] = opcode_costs_struct.I16x8AddSaturateS
	opcode_costs[wasmer.OpcodeI16x8AddSaturateU] = opcode_costs_struct.I16x8AddSaturateU
	opcode_costs[wasmer.OpcodeI16x8Sub] = opcode_costs_struct.I16x8Sub
	opcode_costs[wasmer.OpcodeI16x8SubSaturateS] = opcode_costs_struct.I16x8SubSaturateS
	opcode_costs[wasmer.OpcodeI16x8SubSaturateU] = opcode_costs_struct.I16x8SubSaturateU
	opcode_costs[wasmer.OpcodeI16x8Mul] = opcode_costs_struct.I16x8Mul
	opcode_costs[wasmer.OpcodeI16x8MinS] = opcode_costs_struct.I16x8MinS
	opcode_costs[wasmer.OpcodeI16x8MinU] = opcode_costs_struct.I16x8MinU
	opcode_costs[wasmer.OpcodeI16x8MaxS] = opcode_costs_struct.I16x8MaxS
	opcode_costs[wasmer.OpcodeI16x8MaxU] = opcode_costs_struct.I16x8MaxU
	opcode_costs[wasmer.OpcodeI32x4Neg] = opcode_costs_struct.I32x4Neg
	opcode_costs[wasmer.OpcodeI32x4AnyTrue] = opcode_costs_struct.I32x4AnyTrue
	opcode_costs[wasmer.OpcodeI32x4AllTrue] = opcode_costs_struct.I32x4AllTrue
	opcode_costs[wasmer.OpcodeI32x4Shl] = opcode_costs_struct.I32x4Shl
	opcode_costs[wasmer.OpcodeI32x4ShrS] = opcode_costs_struct.I32x4ShrS
	opcode_costs[wasmer.OpcodeI32x4ShrU] = opcode_costs_struct.I32x4ShrU
	opcode_costs[wasmer.OpcodeI32x4Add] = opcode_costs_struct.I32x4Add
	opcode_costs[wasmer.OpcodeI32x4Sub] = opcode_costs_struct.I32x4Sub
	opcode_costs[wasmer.OpcodeI32x4Mul] = opcode_costs_struct.I32x4Mul
	opcode_costs[wasmer.OpcodeI32x4MinS] = opcode_costs_struct.I32x4MinS
	opcode_costs[wasmer.OpcodeI32x4MinU] = opcode_costs_struct.I32x4MinU
	opcode_costs[wasmer.OpcodeI32x4MaxS] = opcode_costs_struct.I32x4MaxS
	opcode_costs[wasmer.OpcodeI32x4MaxU] = opcode_costs_struct.I32x4MaxU
	opcode_costs[wasmer.OpcodeI64x2Neg] = opcode_costs_struct.I64x2Neg
	opcode_costs[wasmer.OpcodeI64x2AnyTrue] = opcode_costs_struct.I64x2AnyTrue
	opcode_costs[wasmer.OpcodeI64x2AllTrue] = opcode_costs_struct.I64x2AllTrue
	opcode_costs[wasmer.OpcodeI64x2Shl] = opcode_costs_struct.I64x2Shl
	opcode_costs[wasmer.OpcodeI64x2ShrS] = opcode_costs_struct.I64x2ShrS
	opcode_costs[wasmer.OpcodeI64x2ShrU] = opcode_costs_struct.I64x2ShrU
	opcode_costs[wasmer.OpcodeI64x2Add] = opcode_costs_struct.I64x2Add
	opcode_costs[wasmer.OpcodeI64x2Sub] = opcode_costs_struct.I64x2Sub
	opcode_costs[wasmer.OpcodeI64x2Mul] = opcode_costs_struct.I64x2Mul
	opcode_costs[wasmer.OpcodeF32x4Abs] = opcode_costs_struct.F32x4Abs
	opcode_costs[wasmer.OpcodeF32x4Neg] = opcode_costs_struct.F32x4Neg
	opcode_costs[wasmer.OpcodeF32x4Sqrt] = opcode_costs_struct.F32x4Sqrt
	opcode_costs[wasmer.OpcodeF32x4Add] = opcode_costs_struct.F32x4Add
	opcode_costs[wasmer.OpcodeF32x4Sub] = opcode_costs_struct.F32x4Sub
	opcode_costs[wasmer.OpcodeF32x4Mul] = opcode_costs_struct.F32x4Mul
	opcode_costs[wasmer.OpcodeF32x4Div] = opcode_costs_struct.F32x4Div
	opcode_costs[wasmer.OpcodeF32x4Min] = opcode_costs_struct.F32x4Min
	opcode_costs[wasmer.OpcodeF32x4Max] = opcode_costs_struct.F32x4Max
	opcode_costs[wasmer.OpcodeF64x2Abs] = opcode_costs_struct.F64x2Abs
	opcode_costs[wasmer.OpcodeF64x2Neg] = opcode_costs_struct.F64x2Neg
	opcode_costs[wasmer.OpcodeF64x2Sqrt] = opcode_costs_struct.F64x2Sqrt
	opcode_costs[wasmer.OpcodeF64x2Add] = opcode_costs_struct.F64x2Add
	opcode_costs[wasmer.OpcodeF64x2Sub] = opcode_costs_struct.F64x2Sub
	opcode_costs[wasmer.OpcodeF64x2Mul] = opcode_costs_struct.F64x2Mul
	opcode_costs[wasmer.OpcodeF64x2Div] = opcode_costs_struct.F64x2Div
	opcode_costs[wasmer.OpcodeF64x2Min] = opcode_costs_struct.F64x2Min
	opcode_costs[wasmer.OpcodeF64x2Max] = opcode_costs_struct.F64x2Max
	opcode_costs[wasmer.OpcodeI32x4TruncSatF32x4S] = opcode_costs_struct.I32x4TruncSatF32x4S
	opcode_costs[wasmer.OpcodeI32x4TruncSatF32x4U] = opcode_costs_struct.I32x4TruncSatF32x4U
	opcode_costs[wasmer.OpcodeI64x2TruncSatF64x2S] = opcode_costs_struct.I64x2TruncSatF64x2S
	opcode_costs[wasmer.OpcodeI64x2TruncSatF64x2U] = opcode_costs_struct.I64x2TruncSatF64x2U
	opcode_costs[wasmer.OpcodeF32x4ConvertI32x4S] = opcode_costs_struct.F32x4ConvertI32x4S
	opcode_costs[wasmer.OpcodeF32x4ConvertI32x4U] = opcode_costs_struct.F32x4ConvertI32x4U
	opcode_costs[wasmer.OpcodeF64x2ConvertI64x2S] = opcode_costs_struct.F64x2ConvertI64x2S
	opcode_costs[wasmer.OpcodeF64x2ConvertI64x2U] = opcode_costs_struct.F64x2ConvertI64x2U
	opcode_costs[wasmer.OpcodeV8x16Swizzle] = opcode_costs_struct.V8x16Swizzle
	opcode_costs[wasmer.OpcodeV8x16Shuffle] = opcode_costs_struct.V8x16Shuffle
	opcode_costs[wasmer.OpcodeV8x16LoadSplat] = opcode_costs_struct.V8x16LoadSplat
	opcode_costs[wasmer.OpcodeV16x8LoadSplat] = opcode_costs_struct.V16x8LoadSplat
	opcode_costs[wasmer.OpcodeV32x4LoadSplat] = opcode_costs_struct.V32x4LoadSplat
	opcode_costs[wasmer.OpcodeV64x2LoadSplat] = opcode_costs_struct.V64x2LoadSplat
	opcode_costs[wasmer.OpcodeI8x16NarrowI16x8S] = opcode_costs_struct.I8x16NarrowI16x8S
	opcode_costs[wasmer.OpcodeI8x16NarrowI16x8U] = opcode_costs_struct.I8x16NarrowI16x8U
	opcode_costs[wasmer.OpcodeI16x8NarrowI32x4S] = opcode_costs_struct.I16x8NarrowI32x4S
	opcode_costs[wasmer.OpcodeI16x8NarrowI32x4U] = opcode_costs_struct.I16x8NarrowI32x4U
	opcode_costs[wasmer.OpcodeI16x8WidenLowI8x16S] = opcode_costs_struct.I16x8WidenLowI8x16S
	opcode_costs[wasmer.OpcodeI16x8WidenHighI8x16S] = opcode_costs_struct.I16x8WidenHighI8x16S
	opcode_costs[wasmer.OpcodeI16x8WidenLowI8x16U] = opcode_costs_struct.I16x8WidenLowI8x16U
	opcode_costs[wasmer.OpcodeI16x8WidenHighI8x16U] = opcode_costs_struct.I16x8WidenHighI8x16U
	opcode_costs[wasmer.OpcodeI32x4WidenLowI16x8S] = opcode_costs_struct.I32x4WidenLowI16x8S
	opcode_costs[wasmer.OpcodeI32x4WidenHighI16x8S] = opcode_costs_struct.I32x4WidenHighI16x8S
	opcode_costs[wasmer.OpcodeI32x4WidenLowI16x8U] = opcode_costs_struct.I32x4WidenLowI16x8U
	opcode_costs[wasmer.OpcodeI32x4WidenHighI16x8U] = opcode_costs_struct.I32x4WidenHighI16x8U
	opcode_costs[wasmer.OpcodeI16x8Load8x8S] = opcode_costs_struct.I16x8Load8x8S
	opcode_costs[wasmer.OpcodeI16x8Load8x8U] = opcode_costs_struct.I16x8Load8x8U
	opcode_costs[wasmer.OpcodeI32x4Load16x4S] = opcode_costs_struct.I32x4Load16x4S
	opcode_costs[wasmer.OpcodeI32x4Load16x4U] = opcode_costs_struct.I32x4Load16x4U
	opcode_costs[wasmer.OpcodeI64x2Load32x2S] = opcode_costs_struct.I64x2Load32x2S
	opcode_costs[wasmer.OpcodeI64x2Load32x2U] = opcode_costs_struct.I64x2Load32x2U
	opcode_costs[wasmer.OpcodeI8x16RoundingAverageU] = opcode_costs_struct.I8x16RoundingAverageU
	opcode_costs[wasmer.OpcodeI16x8RoundingAverageU] = opcode_costs_struct.I16x8RoundingAverageU
	opcode_costs[wasmer.OpcodeLocalAllocate] = opcode_costs_struct.LocalAllocate
	// opcode_costs_struct.LocalsUnmetered is not added to the opcode_costs
	// array; the value will be sent to Wasmer as a compilation option instead

	return opcode_costs
}
