package elliptic_curve

import (
	"crypto/rand"
	"fmt"
	"math/big"
)

type PrivateKey struct {
	secret *big.Int
	point  *Point // Public Key
}

func NewPrivateKey(sec *big.Int) *PrivateKey {
	G := GetGenerator()
	return &PrivateKey{
		secret: sec,
		// Public Key
		point: G.ScalarMul(sec),
	}
}

func (p *PrivateKey) String() string {
	return fmt.Sprintf("private key hex : %x", p.secret)
}

func (p *PrivateKey) GetPublicKey() *Point {
	return p.point
}

func (p *PrivateKey) Sign(z *big.Int) *Signature {
	n := GetBitCoinValueN()
	k, err := rand.Int(rand.Reader, n)
	if err != nil {
		panic(err)
	}
	kField := NewFieldElement(n, k)
	G := GetGenerator()

	r := G.ScalarMul(k).x.num
	rField := NewFieldElement(n, r)
	eField := NewFieldElement(n, p.secret)
	zField := NewFieldElement(n, z)

	rMulSecret := rField.Mul(eField)
	zAddMulSecret := zField.Add(rMulSecret)
	KInverse := kField.Inverse()
	sField := zAddMulSecret.Mul(KInverse)

	var opDiv big.Int
	if sField.num.Cmp(opDiv.Div(n, big.NewInt(int64(2)))) > 0 {
		var opSub big.Int
		sField = NewFieldElement(n, opSub.Sub(n, sField.num))
	}

	return &Signature{
		r: NewFieldElement(n, r),
		s: sField,
	}
}
