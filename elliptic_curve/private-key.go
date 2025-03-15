package elliptic_curve

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"math/big"
)

type PrivateKey struct {
	secret *big.Int
	point  *Point
}

func NewPrivateKey(secret *big.Int) *PrivateKey {
	G := GetGenerator()
	return &PrivateKey{
		secret: secret,
		//public key
		point: G.ScalarMul(secret),
	}
}

func (p *PrivateKey) String() string {
	return fmt.Sprintf("private key hex:{%s}", p.secret)
}

func (p *PrivateKey) GetPublicKey() *Point {
	return p.point
}

func (p *PrivateKey) Sign(z *big.Int) *Signature {
	//(s, r)
	//s = (z + r * e) / k
	// k is a strong random number
	n := GetBitcoinValueN()
	k, err := rand.Int(rand.Reader, n)
	if err != nil {
		panic(fmt.Sprintf("Sign err with rand int: %s", err))
	}
	kField := NewFieldElement(n, k)
	G := GetGenerator()
	// s = (z + r * e) / k
	// r = G * k
	r := G.ScalarMul(k).x.num
	rField := NewFieldElement(n, r)
	eField := NewFieldElement(n, p.secret)
	zField := NewFieldElement(n, z)
	// r*e
	rMulSecret := rField.Mul(eField)
	// z+r*e
	zAddRMulSecret := zField.Add(rMulSecret)
	// /k
	kInverse := kField.Inverse()
	sField := zAddRMulSecret.Mul(kInverse)
	/*
	   if s > n / 2 we need to change it to n - s, when doing signature
	   verify, s and n - s are equivalence doing this change is for malleability reasons, detail:
	   https://bitcoin.stackexchange.com/questions/85946/low-s-value-in-bitcoin-signature
	*/
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

/*
Wif
1, set first byte to 0x80 MainNet, 0xef testnet
2, append the bytes array of private key behind first byte,
if the length of bytes array < 32 bytes, append leading 0 to 32 bytes
3, if public SEC compressed, add suffix byte 0x01
4, do hash256 on result of step 3 and get its first 4 bytes
5, append the 4 bytes from step 4 with result from step 3 and do base58 encoding

4,5 base58checksum
*/
func (p *PrivateKey) Wif(compressed bool, testnet bool) string {
	bytes := []byte{}
	if testnet {
		bytes = append(bytes, 0xef)
	} else {
		bytes = append(bytes, 0x80)
	}

	secretBytes := p.secret.Bytes()
	if len(secretBytes) < 32 {
		//two chars turn into one byte
		s := fmt.Sprintf("%064x", secretBytes)
		paddingBytes, err := hex.DecodeString(s)
		if err != nil {
			panic(fmt.Sprintf("padding secret bytes err: %v\n", err))
		}
		secretBytes = paddingBytes
	}

	bytes = append(bytes, secretBytes...)
	if compressed {
		bytes = append(bytes, 0x01)
	}

	return Base58Checksum(bytes)

}
