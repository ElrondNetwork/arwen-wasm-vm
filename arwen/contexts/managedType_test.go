package contexts

import (
	"crypto/elliptic"
	"math/big"
	"testing"

	"github.com/ElrondNetwork/arwen-wasm-vm/v1_3/arwen"
	contextmock "github.com/ElrondNetwork/arwen-wasm-vm/v1_3/mock/context"
	"github.com/stretchr/testify/require"
)

func TestNewManagedTypes(t *testing.T) {
	t.Parallel()

	host := &contextmock.VMHostStub{}

	managedTypesContext, err := NewManagedTypesContext(host)

	require.Nil(t, err)
	require.False(t, managedTypesContext.IsInterfaceNil())
	require.NotNil(t, managedTypesContext.bigIntValues)
	require.NotNil(t, managedTypesContext.bigIntStateStack)
	require.NotNil(t, managedTypesContext.ecValues)
	require.NotNil(t, managedTypesContext.ecStateStack)
	require.Equal(t, 0, len(managedTypesContext.bigIntValues))
	require.Equal(t, 0, len(managedTypesContext.bigIntStateStack))
	require.Equal(t, 0, len(managedTypesContext.ecValues))
	require.Equal(t, 0, len(managedTypesContext.ecStateStack))
}

func TestManagedTypesContext_InitPushPopState(t *testing.T) {
	t.Parallel()
	host := &contextmock.VMHostStub{}
	value1, value2, value3 := int64(100), int64(200), int64(-42)
	p224ec, p256ec, p384ec, p521ec := elliptic.P224().Params(), elliptic.P256().Params(), elliptic.P384().Params(), elliptic.P521().Params()
	managedTypesContext, _ := NewManagedTypesContext(host)
	managedTypesContext.InitState()

	// Create 2 bigInt and 2 EC values on the active state
	index1 := managedTypesContext.PutBigInt(value1)
	require.Equal(t, int32(0), index1)
	index2 := managedTypesContext.PutBigInt(value2)
	require.Equal(t, int32(1), index2)

	bigValue1, err := managedTypesContext.GetBigInt(index1)
	require.Equal(t, big.NewInt(value1), bigValue1)
	require.Equal(t, nil, err)
	bigValue2, err := managedTypesContext.GetBigInt(index2)
	require.Equal(t, big.NewInt(value2), bigValue2)
	require.Equal(t, nil, err)

	ecIndex1 := managedTypesContext.PutEllipticCurve(p224ec)
	require.Equal(t, int32(0), ecIndex1)
	ecIndex2 := managedTypesContext.PutEllipticCurve(p256ec)
	require.Equal(t, int32(1), ecIndex2)

	ec1, err := managedTypesContext.GetEllipticCurve(ecIndex1)
	require.Nil(t, err)
	require.Equal(t, p224ec, ec1)
	ec2, err := managedTypesContext.GetEllipticCurve(ecIndex2)
	require.Nil(t, err)
	require.Equal(t, p256ec, ec2)

	p224NormalGasCostMultiplier := managedTypesContext.Get100xCurveGasCostMultiplier(ecIndex1)
	require.Equal(t, int32(100), p224NormalGasCostMultiplier)
	p256NormalGasCostMultiplier := managedTypesContext.Get100xCurveGasCostMultiplier(ecIndex2)
	require.Equal(t, int32(135), p256NormalGasCostMultiplier)
	p224ScalarMultGasCostMultiplier := managedTypesContext.GetScalarMult100xCurveGasCostMultiplier(ecIndex1)
	require.Equal(t, int32(100), p224ScalarMultGasCostMultiplier)
	p256ScalarMultGasCostMultiplier := managedTypesContext.GetScalarMult100xCurveGasCostMultiplier(ecIndex2)
	require.Equal(t, int32(110), p256ScalarMultGasCostMultiplier)
	p224UCompressedGasCostMultiplier := managedTypesContext.GetUCompressed100xCurveGasCostMultiplier(ecIndex1)
	require.Equal(t, int32(2000), p224UCompressedGasCostMultiplier)
	p256UCompressedGasCostMultiplier := managedTypesContext.GetUCompressed100xCurveGasCostMultiplier(ecIndex2)
	require.Equal(t, int32(100), p256UCompressedGasCostMultiplier)

	// Copy active state to stack, then clean it. The previous 2 values should not
	// be accessible.
	managedTypesContext.PushState()
	require.Equal(t, 1, len(managedTypesContext.bigIntStateStack))
	managedTypesContext.InitState()

	bigValue1, err = managedTypesContext.GetBigInt(index1)
	require.Nil(t, bigValue1)
	require.Equal(t, arwen.ErrNoBigIntUnderThisHandle, err)
	bigValue2, err = managedTypesContext.GetBigInt(index2)
	require.Nil(t, bigValue2)
	require.Equal(t, arwen.ErrNoBigIntUnderThisHandle, err)

	ec1, err = managedTypesContext.GetEllipticCurve(ecIndex1)
	require.Nil(t, ec1)
	require.Equal(t, arwen.ErrNoEllipticCurveUnderThisHandle, err)
	ec2, err = managedTypesContext.GetEllipticCurve(ecIndex2)
	require.Nil(t, ec2)
	require.Equal(t, arwen.ErrNoEllipticCurveUnderThisHandle, err)

	p224NormalGasCostMultiplier = managedTypesContext.Get100xCurveGasCostMultiplier(ecIndex1)
	require.Equal(t, int32(-1), p224NormalGasCostMultiplier)
	p256NormalGasCostMultiplier = managedTypesContext.Get100xCurveGasCostMultiplier(ecIndex2)
	require.Equal(t, int32(-1), p256NormalGasCostMultiplier)
	p224ScalarMultGasCostMultiplier = managedTypesContext.GetScalarMult100xCurveGasCostMultiplier(ecIndex1)
	require.Equal(t, int32(-1), p224ScalarMultGasCostMultiplier)
	p256ScalarMultGasCostMultiplier = managedTypesContext.GetScalarMult100xCurveGasCostMultiplier(ecIndex2)
	require.Equal(t, int32(-1), p256ScalarMultGasCostMultiplier)
	p224UCompressedGasCostMultiplier = managedTypesContext.GetUCompressed100xCurveGasCostMultiplier(ecIndex1)
	require.Equal(t, int32(-1), p224UCompressedGasCostMultiplier)
	p256UCompressedGasCostMultiplier = managedTypesContext.GetUCompressed100xCurveGasCostMultiplier(ecIndex2)
	require.Equal(t, int32(-1), p256UCompressedGasCostMultiplier)

	// Add a value on the current active state
	index3 := managedTypesContext.PutBigInt(value3)
	require.Equal(t, int32(0), index3)
	bigValue3, err := managedTypesContext.GetBigInt(index3)
	require.Equal(t, big.NewInt(value3), bigValue3)
	require.Equal(t, nil, err)

	ecIndex3 := managedTypesContext.PutEllipticCurve(p384ec)
	require.Equal(t, int32(0), ecIndex3)
	ec3, err := managedTypesContext.GetEllipticCurve(ecIndex3)
	require.Nil(t, err)
	require.Equal(t, p384ec, ec3)

	p384NormalGasCostMultiplier := managedTypesContext.Get100xCurveGasCostMultiplier(ecIndex3)
	require.Equal(t, int32(200), p384NormalGasCostMultiplier)
	p384ScalarMultGasCostMultiplier := managedTypesContext.GetScalarMult100xCurveGasCostMultiplier(ecIndex3)
	require.Equal(t, int32(150), p384ScalarMultGasCostMultiplier)
	p384UCompressedGasCostMultiplier := managedTypesContext.GetUCompressed100xCurveGasCostMultiplier(ecIndex3)
	require.Equal(t, int32(200), p384UCompressedGasCostMultiplier)

	// Copy active state to stack, then clean it. The previous 3 values should not
	// be accessible.
	managedTypesContext.PushState()
	require.Equal(t, 2, len(managedTypesContext.bigIntStateStack))
	managedTypesContext.InitState()

	bigValue1, err = managedTypesContext.GetBigInt(index1)
	require.Nil(t, bigValue1)
	require.Equal(t, arwen.ErrNoBigIntUnderThisHandle, err)
	bigValue2, err = managedTypesContext.GetBigInt(index2)
	require.Nil(t, bigValue2)
	require.Equal(t, arwen.ErrNoBigIntUnderThisHandle, err)
	bigValue3, err = managedTypesContext.GetBigInt(index3)
	require.Nil(t, bigValue3)
	require.Equal(t, arwen.ErrNoBigIntUnderThisHandle, err)

	ec1, err = managedTypesContext.GetEllipticCurve(ecIndex1)
	require.Nil(t, ec1)
	require.Equal(t, arwen.ErrNoEllipticCurveUnderThisHandle, err)
	ec2, err = managedTypesContext.GetEllipticCurve(ecIndex2)
	require.Nil(t, ec2)
	require.Equal(t, arwen.ErrNoEllipticCurveUnderThisHandle, err)
	ec3, err = managedTypesContext.GetEllipticCurve(ecIndex3)
	require.Nil(t, ec3)
	require.Equal(t, arwen.ErrNoEllipticCurveUnderThisHandle, err)

	value4 := int64(84)
	index4 := managedTypesContext.PutBigInt(value4)
	require.Equal(t, int32(0), index4)
	bigValue4, err := managedTypesContext.GetBigInt(index4)
	require.Equal(t, big.NewInt(value4), bigValue4)
	require.Equal(t, nil, err)

	ecIndex4 := managedTypesContext.PutEllipticCurve(p521ec)
	require.Equal(t, int32(0), ecIndex4)
	ec4, err := managedTypesContext.GetEllipticCurve(ecIndex4)
	require.Equal(t, p521ec, ec4)
	require.Nil(t, err)

	p521NormalGasCostMultiplier := managedTypesContext.Get100xCurveGasCostMultiplier(ecIndex4)
	require.Equal(t, int32(250), p521NormalGasCostMultiplier)
	p521ScalarMultGasCostMultiplier := managedTypesContext.GetScalarMult100xCurveGasCostMultiplier(ecIndex4)
	require.Equal(t, int32(190), p521ScalarMultGasCostMultiplier)
	p521UCompressedGasCostMultiplier := managedTypesContext.GetUCompressed100xCurveGasCostMultiplier(ecIndex4)
	require.Equal(t, int32(400), p521UCompressedGasCostMultiplier)

	// Discard the top of the stack, losing value3; value4 should still be
	// accessible, since its in the active state.
	managedTypesContext.PopDiscard()
	require.Equal(t, 1, len(managedTypesContext.bigIntStateStack))
	bigValue4, err = managedTypesContext.GetBigInt(index4)
	require.Equal(t, big.NewInt(value4), bigValue4)
	require.Equal(t, nil, err)

	ec4, err = managedTypesContext.GetEllipticCurve(ecIndex4)
	require.Equal(t, p521ec, ec4)
	require.Nil(t, err)
	// Restore the first active state by popping to the active state (which is
	// lost).
	managedTypesContext.PopSetActiveState()
	require.Equal(t, 0, len(managedTypesContext.bigIntStateStack))

	bigValue1, err = managedTypesContext.GetBigInt(index1)
	require.Equal(t, big.NewInt(value1), bigValue1)
	require.Equal(t, nil, err)
	bigValue2, err = managedTypesContext.GetBigInt(index2)
	require.Equal(t, big.NewInt(value2), bigValue2)
	require.Equal(t, nil, err)

	ec1, err = managedTypesContext.GetEllipticCurve(ecIndex1)
	require.Nil(t, err)
	require.Equal(t, p224ec, ec1)
	ec2, err = managedTypesContext.GetEllipticCurve(ecIndex2)
	require.Nil(t, err)
	require.Equal(t, p256ec, ec2)
}

func TestManagedTypesContext_PutGet(t *testing.T) {
	t.Parallel()
	host := &contextmock.VMHostStub{}

	value1, value2, value3, value4 := int64(100), int64(200), int64(-42), int64(-80)
	p224ec, p256ec, p384ec, p521ec := elliptic.P224().Params(), elliptic.P256().Params(), elliptic.P384().Params(), elliptic.P521().Params()
	managedTypesContext, _ := NewManagedTypesContext(host)

	index1 := managedTypesContext.PutBigInt(value1)
	require.Equal(t, int32(0), index1)
	index2 := managedTypesContext.PutBigInt(value2)
	require.Equal(t, int32(1), index2)
	index3 := managedTypesContext.PutBigInt(value3)
	require.Equal(t, int32(2), index3)

	ecIndex1 := managedTypesContext.PutEllipticCurve(p224ec)
	require.Equal(t, int32(0), ecIndex1)
	ecIndex2 := managedTypesContext.PutEllipticCurve(p256ec)
	require.Equal(t, int32(1), ecIndex2)
	ecIndex3 := managedTypesContext.PutEllipticCurve(p384ec)
	require.Equal(t, int32(2), ecIndex3)

	bigValue1, err := managedTypesContext.GetBigInt(index1)
	require.Equal(t, big.NewInt(value1), bigValue1)
	require.Equal(t, nil, err)
	bigValue2, err := managedTypesContext.GetBigInt(index2)
	require.Equal(t, big.NewInt(value2), bigValue2)
	require.Equal(t, nil, err)
	bigValue4, err := managedTypesContext.GetBigInt(int32(3))
	require.Nil(t, bigValue4)
	require.Equal(t, arwen.ErrNoBigIntUnderThisHandle, err)
	bigValue4 = managedTypesContext.GetBigIntOrCreate(3)
	require.Equal(t, big.NewInt(0), bigValue4)

	ec1, err := managedTypesContext.GetEllipticCurve(ecIndex1)
	require.Nil(t, err)
	require.Equal(t, p224ec, ec1)
	ec2, err := managedTypesContext.GetEllipticCurve(ecIndex2)
	require.Nil(t, err)
	require.Equal(t, p256ec, ec2)
	ec4, err := managedTypesContext.GetEllipticCurve(int32(3))
	require.Nil(t, ec4)
	require.Equal(t, arwen.ErrNoEllipticCurveUnderThisHandle, err)

	index4 := managedTypesContext.PutBigInt(value4)
	require.Equal(t, int32(4), index4)
	bigValue4 = managedTypesContext.GetBigIntOrCreate(4)
	require.Equal(t, big.NewInt(value4), bigValue4)

	ecIndex4 := managedTypesContext.PutEllipticCurve(p521ec)
	require.Equal(t, int32(3), ecIndex4)
	ec4, err = managedTypesContext.GetEllipticCurve(ecIndex4)
	require.Nil(t, err)
	require.Equal(t, p521ec, ec4)

	bigValue, err := managedTypesContext.GetBigInt(123)
	require.Nil(t, bigValue)
	require.Equal(t, arwen.ErrNoBigIntUnderThisHandle, err)
	bigValue1, err = managedTypesContext.GetBigInt(index1)
	require.Equal(t, big.NewInt(value1), bigValue1)
	require.Equal(t, nil, err)
	bigValue2, err = managedTypesContext.GetBigInt(index2)
	require.Equal(t, big.NewInt(value2), bigValue2)
	require.Equal(t, nil, err)

	bigValue1, err = managedTypesContext.GetBigInt(index1)
	require.Equal(t, big.NewInt(value1), bigValue1)
	require.Equal(t, nil, err)
	bigValue2, err = managedTypesContext.GetBigInt(index2)
	require.Equal(t, big.NewInt(value2), bigValue2)
	require.Equal(t, nil, err)
	bigValue3, err := managedTypesContext.GetBigInt(index3)
	require.Equal(t, big.NewInt(value3), bigValue3)
	require.Equal(t, nil, err)

}

func TestManagedTypesContext_PopSetActiveStateIfStackIsEmptyShouldNotPanic(t *testing.T) {
	t.Parallel()
	host := &contextmock.VMHostStub{}

	managedTypesContext, _ := NewManagedTypesContext(host)
	managedTypesContext.PopSetActiveState()

	require.Equal(t, 0, len(managedTypesContext.bigIntStateStack))
	require.Equal(t, 0, len(managedTypesContext.ecStateStack))

}

func TestManagedTypesContext_PopDiscardIfStackIsEmptyShouldNotPanic(t *testing.T) {
	t.Parallel()
	host := &contextmock.VMHostStub{}

	managedTypesContext, _ := NewManagedTypesContext(host)
	managedTypesContext.PopDiscard()

	require.Equal(t, 0, len(managedTypesContext.bigIntStateStack))
	require.Equal(t, 0, len(managedTypesContext.ecStateStack))

}
