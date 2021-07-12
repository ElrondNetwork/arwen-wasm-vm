package elliptic_curves

import (
	"crypto/elliptic"
	"crypto/rand"
	"math/big"
)

type ellipticCurve struct {
	curve elliptic.CurveParams
}

// NewEllipticCurve returns an ellipticCurve struct with CurveParams set with the input parameters given.
func NewEllipticCurve(P *big.Int, N *big.Int, B *big.Int, Gx *big.Int, Gy *big.Int, BitSize int) *ellipticCurve {
	myEc := &ellipticCurve{}
	myEc.curve = elliptic.CurveParams{P: P, N: N, B: B, Gx: Gx, Gy: Gy, BitSize: BitSize, Name: "EC"}
	return myEc
}

// Add returns the sum of (x1,y1) and (x2,y2)
func Add(P *big.Int, N *big.Int, B *big.Int, Gx *big.Int, Gy *big.Int, BitSize int, x1 *big.Int, y1 *big.Int, x2 *big.Int, y2 *big.Int) (x *big.Int, y *big.Int) {
	myEc := NewEllipticCurve(P, N, B, Gx, Gy, BitSize)
	return myEc.curve.Add(x1, y1, x2, y2)
}

// Double returns 2*(x,y)
func Double(P *big.Int, N *big.Int, B *big.Int, Gx *big.Int, Gy *big.Int, BitSize int, x *big.Int, y *big.Int) (*big.Int, *big.Int) {
	myEc := NewEllipticCurve(P, N, B, Gx, Gy, BitSize)
	return myEc.curve.Double(x, y)
}

// IsOnCurve reports whether the given (x,y) lies on the curve.
func IsOnCurve(P *big.Int, N *big.Int, B *big.Int, Gx *big.Int, Gy *big.Int, BitSize int, x *big.Int, y *big.Int) bool {
	myEc := NewEllipticCurve(P, N, B, Gx, Gy, BitSize)
	return myEc.curve.IsOnCurve(x, y)
}

// ScalarBaseMult returns k*G, where G is the base point of the group and k is an integer in big-endian form.
func ScalarBaseMult(P *big.Int, N *big.Int, B *big.Int, Gx *big.Int, Gy *big.Int, BitSize int, k []byte) (*big.Int, *big.Int) {
	myEc := NewEllipticCurve(P, N, B, Gx, Gy, BitSize)
	return myEc.curve.ScalarBaseMult(k)
}

// ScalarMult returns k*(Bx,By) where k is a number in big-endian form.
func ScalarMult(P *big.Int, N *big.Int, B *big.Int, Gx *big.Int, Gy *big.Int, BitSize int, Bx *big.Int, By *big.Int, k []byte) (*big.Int, *big.Int) {
	myEc := NewEllipticCurve(P, N, B, Gx, Gy, BitSize)
	return myEc.curve.ScalarMult(Bx, By, k)
}

// Marshal converts a point on the curve into the uncompressed form specified in section 4.3.6 of ANSI X9.62.
func Marshal(P *big.Int, N *big.Int, B *big.Int, Gx *big.Int, Gy *big.Int, BitSize int, x *big.Int, y *big.Int) []byte {
	myEc := NewEllipticCurve(P, N, B, Gx, Gy, BitSize)
	return elliptic.Marshal(&myEc.curve, x, y)
}

// Unmarshal converts a point, serialized by Marshal, into an x, y pair.
// It is an error if the point is not in uncompressed form or is not on the curve.
// On error, x = nil.
func Unmarshal(P *big.Int, N *big.Int, B *big.Int, Gx *big.Int, Gy *big.Int, BitSize int, data []byte) (*big.Int, *big.Int) {
	myEc := NewEllipticCurve(P, N, B, Gx, Gy, BitSize)
	return elliptic.Unmarshal(&myEc.curve, data)
}

// GenerateKey returns a public/private key pair. The private key is
// generated using the given reader, which must return random data.
func GenerateKey(P *big.Int, N *big.Int, B *big.Int, Gx *big.Int, Gy *big.Int, BitSize int) (priv []byte, x, y *big.Int, err error) {
	myEc := NewEllipticCurve(P, N, B, Gx, Gy, BitSize)
	return elliptic.GenerateKey(&myEc.curve, rand.Reader)
}

// MarshalCompressed converts a point on the curve into the compressed form specified in section 4.3.6 of ANSI X9.62.
func MarshalCompressed(P *big.Int, N *big.Int, B *big.Int, Gx *big.Int, Gy *big.Int, BitSize int, x *big.Int, y *big.Int) []byte {
	myEc := NewEllipticCurve(P, N, B, Gx, Gy, BitSize)
	return elliptic.MarshalCompressed(&myEc.curve, x, y)
}

// UnmarshalCompressed converts a point, serialized by MarshalCompressed, into an x, y pair.
// It is an error if the point is not in compressed form or is not on the curve.
// On error, x = nil.
func UnmarshalCompressed(P *big.Int, N *big.Int, B *big.Int, Gx *big.Int, Gy *big.Int, BitSize int, data []byte) (*big.Int, *big.Int) {
	myEc := NewEllipticCurve(P, N, B, Gx, Gy, BitSize)
	return elliptic.UnmarshalCompressed(&myEc.curve, data)
}
