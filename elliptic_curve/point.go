package elliptic_curve

import (
	"fmt"
	"math/big"
)

type OP_TPYE int

const (
	ADD OP_TPYE = iota
	SUB
	MUL
	DIV
	EXP
)

type Point struct {
	//for coefficients for elliptic curve
	a *FieldElement
	b *FieldElement
	//the value of x, y may be ver huge
	x *FieldElement
	y *FieldElement
}

func OpOnBig(x *FieldElement, y *FieldElement, scalar *big.Int, opType OP_TPYE) *FieldElement {
	switch opType {
	case ADD:
		return x.Add(y)
	case SUB:
		return x.Sub(y)
	case MUL:
		if y != nil {
			return x.Mul(y)
		}
		if scalar != nil {
			return x.ScalarMul(scalar)
		}
		panic("error in multiply")
	case DIV:
		return x.Div(y)
	case EXP:
		if scalar == nil {
			panic("scalar should not be nil for EXP")
		}
		return x.Pow(scalar)
	}

	panic("should not come here")
}

func S256Point(x *big.Int, y *big.Int) *Point {
	a := S256Field(big.NewInt(int64(0)))
	b := S256Field(big.NewInt(int64(7)))

	if x == nil && y == nil {
		return &Point{
			x: nil,
			y: nil,
			a: a,
			b: b,
		}
	}

	return &Point{
		x: S256Field(x),
		y: S256Field(y),
		a: a,
		b: b,
	}
}

func (p *Point) Verify(z *FieldElement, sig *Signature) bool {
	/*
			7. any one who want to verify message z is created by owner of e:
			    1, compute u = z/s, v=r/s,
				2, compute u*G + v*P = (z/s)*G + (r/s)*P = (z/s)*G+(r/s)*eG
					=(z/s)*P + (r*e/s)*G = ((z+r*e)/s))*G = k*G = R'
				3, take the x coordinate of R' compare with r
					if the same => verify the message z is created by owner of e

				notice we have shown that n * G is identity, therefore the above computation related to
		        z, s, r, e need to do base on modular of n, and remember the operator
				"/" is not the normal arithmetic divide , it's the inverse of  multiplication.
	*/
	sInverse := sig.s.Inverse()
	u := z.Mul(sInverse)
	v := sig.r.Mul(sInverse)
	G := GetGenerator()
	total := (G.ScalarMul(u.num)).Add(p.ScalarMul(v.num))

	return total.x.num.Cmp(sig.r.num) == 0
}

func NewECPoint(x *FieldElement, y *FieldElement, a *FieldElement, b *FieldElement) *Point {
	if x == nil && y == nil {
		return &Point{
			a: a,
			b: b,
			x: x,
			y: y,
		}
	}
	//first check (x,y) on the curve defined by a, b
	left := OpOnBig(y, nil, big.NewInt(int64(2)), EXP)
	x3 := OpOnBig(x, nil, big.NewInt(int64(3)), EXP)
	ax := OpOnBig(a, x, nil, MUL)
	right := OpOnBig(OpOnBig(x3, ax, nil, ADD), b, nil, ADD)
	//if x and y are nil, then its identity point and
	//we don't need to check it on curve
	if left.EqualTo(right) != true {
		err := fmt.Sprintf("point:(%v, %v) is not on the curve with a: %v, b:%v\n", x, y, a, b)
		panic(err)
	}

	return &Point{
		a: a,
		b: b,
		x: x,
		y: y,
	}
}

func (p *Point) ScalarMul(scalar *big.Int) *Point {
	if scalar == nil {
		panic("scalar mul error ofr nil scalar")
	}
	/*
		turn scalar into binary string form, for example 13 will turn into "1101"
	*/
	binaryForm := fmt.Sprintf("%b", scalar)
	current := p
	result := NewECPoint(nil, nil, p.a, p.b)
	for i := len(binaryForm) - 1; i >= 0; i-- {
		if binaryForm[i] == '1' {
			result = result.Add(current)
		}
		//add itself is the same as left shift 1 bit
		current = current.Add(current)
	}

	return result
}

func (p *Point) Add(other *Point) *Point {
	//check points are on the same curve
	if p.a.EqualTo(other.a) != true || p.b.EqualTo(other.b) != true {
		panic("given two points are not on the same curve")
	}

	if p.x == nil {
		//current point is identity point
		return other
	}

	if other.x == nil {
		//the other point is identity
		return p
	}

	/*
		another simple case, two points on the same vertical line, that is
		they have the same x but inverse y, the addition of them should be
		identity
	*/
	zero := NewFieldElement(p.x.order, big.NewInt(int64(0)))
	if p.x.EqualTo(other.x) == true &&
		OpOnBig(p.y, other.y, nil, ADD).EqualTo(zero) == true {
		return &Point{
			x: nil,
			y: nil,
			a: p.a,
			b: p.b,
		}
	}

	//find slope of line AB
	//x1 = p.x, y1 = p.y, x2 = other.x, y2 = other.y
	var numerator *FieldElement
	var denominator *FieldElement
	if p.x.EqualTo(other.x) == true && p.y.EqualTo(other.y) == true {
		//two points are the same and compute the slope of tangent line
		//numerator is (3x^2+a)
		xSqrt := OpOnBig(p.x, nil, big.NewInt(int64(2)), EXP)
		threeXSqrt := OpOnBig(xSqrt, nil, big.NewInt(int64(3)), MUL)
		numerator = OpOnBig(threeXSqrt, p.a, nil, ADD)
		//demoninator is 2y
		denominator = OpOnBig(p.y, nil, big.NewInt(int64(2)), MUL)
	} else {
		//s = (y2-y2)/(x2-x1)
		numerator = OpOnBig(other.y, p.y, nil, SUB)
		denominator = OpOnBig(other.x, p.x, nil, SUB)
	}

	slope := OpOnBig(numerator, denominator, nil, DIV)

	//-s^2
	slopeSqrt := OpOnBig(slope, nil, big.NewInt(int64(2)), EXP)
	x3 := OpOnBig(OpOnBig(slopeSqrt, p.x, nil, SUB), other.x, nil, SUB)
	//x3-x1
	x3MinusX1 := OpOnBig(x3, p.x, nil, SUB)
	//y3=s(x3-x1)+y1
	y3 := OpOnBig(OpOnBig(slope, x3MinusX1, nil, MUL), p.y, nil, ADD)
	//-y3
	minusY3 := OpOnBig(y3, nil, big.NewInt(int64(-1)), MUL)

	return &Point{
		x: x3,
		y: minusY3,
		a: p.a,
		b: p.b,
	}
}

func (p *Point) String() string {
	xString := "nil"
	yString := "nil"
	if p.x != nil {
		xString = p.x.String()
	}
	if p.y != nil {
		yString = p.y.String()
	}
	return fmt.Sprintf("(x: %s, y: %s, a: %s, b: %s)\n", xString,
		yString, p.a.String(), p.b.String())
}

func (p *Point) Equal(other *Point) bool {
	if p.a.EqualTo(other.a) == true && p.b.EqualTo(other.b) == true &&
		p.x.EqualTo(other.x) == true && p.y.EqualTo(other.y) == true {
		return true
	}

	return false
}

func (p *Point) NoEqual(other *Point) bool {
	if p.a.EqualTo(other.a) != true || p.b.EqualTo(other.b) != true ||
		p.x.EqualTo(other.x) != true || p.y.EqualTo(other.y) != true {
		return true
	}

	return false
}

func (p *Point) Sec(compressed bool) (string, []byte) {
	secBytes := []byte{}
	if !compressed {
		/*
			uncompressed sec:
			1. first byte 04
			2. x in big endian hex string
			3. y in big endian hex string
			padding x,y with leading 0
		*/
		secBytes = append(secBytes, 0x04)
		secBytes = append(secBytes, p.x.num.Bytes()...)
		secBytes = append(secBytes, p.y.num.Bytes()...)
		return fmt.Sprintf("04%064x%064x", p.x.num, p.y.num), secBytes
	}

	var opMod big.Int
	if opMod.Mod(p.y.num, big.NewInt(int64(2))).Cmp(big.NewInt(int64(0))) == 0 {
		//y is even, set first byte t0 0x02
		secBytes = append(secBytes, 0x02)
		secBytes = append(secBytes, p.x.num.Bytes()...)
		return fmt.Sprintf("02%064x", p.x.num), secBytes
	} else {
		secBytes = append(secBytes, 0x03)
		secBytes = append(secBytes, p.x.num.Bytes()...)
		return fmt.Sprintf("03%064x", p.x.num), secBytes
	}
}

func (p *Point) hash160(compressed bool) []byte {
	_, secBytes := p.Sec(compressed)
	return Hash160(secBytes)
}

func (p *Point) Address(compressed bool, testnet bool) string {
	hash160 := p.hash160(compressed)
	prefix := []byte{}
	if testnet {
		prefix = append(prefix, 0x6f)
	} else {
		prefix = append(prefix, 0x00)
	}

	return Base58Checksum(append(prefix, hash160...))
}
