package elliptic_curve

import (
	"fmt"
	"math/big"
)

type OpType int

const (
	ADD OpType = iota
	SUB
	MUL
	DIV
	EXP
)

type Point struct {
	x, y *FieldElement
	a, b *FieldElement
}

func OpOnBig(x, y *FieldElement,
	scalar *big.Int, op_type OpType) *FieldElement {

	switch op_type {
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
		panic("error in Multiply")
	case DIV:
		return x.Div(y)
	case EXP:
		if scalar == nil {
			panic("scalar should not be nil")
		}

		return x.Pow(scalar)
	}

	panic("unknown operation type")
}

func NewECPoint(x, y, a, b *FieldElement) *Point {
	if x == nil && y == nil {
		return &Point{
			x: x,
			y: y,
			a: a,
			b: b,
		}
	}
	left := OpOnBig(y, nil, big.NewInt(int64(2)), EXP)
	x3 := OpOnBig(x, nil, big.NewInt(int64(3)), EXP)
	ax := OpOnBig(a, x, nil, MUL)
	right := OpOnBig(OpOnBig(x3, ax, nil, ADD), b, nil, ADD)
	if left.EqualTo(right) != true {
		err := fmt.Sprintf("Point(%v, %v) is not on the curve !")
		panic(err)
	}

	return &Point{
		x: x,
		y: y,
		a: a,
		b: b,
	}
}

func (p *Point) Add(other *Point) *Point {
	if p.a.EqualTo(other.a) != true &&
		p.b.EqualTo(other.b) != true {
		panic("points are not on the same curve")
	}

	if p.x == nil {
		return other
	}

	if other.x == nil {
		return p
	}

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

	var numerator *FieldElement
	var denominator *FieldElement
	if p.x.EqualTo(other.x) == true &&
		p.y.EqualTo(other.y) == true {
		xSqrt := OpOnBig(p.x, nil, big.NewInt(int64(2)), EXP)
		threeXSqrt := OpOnBig(xSqrt, nil, big.NewInt(int64(3)), MUL)
		numerator = OpOnBig(threeXSqrt, p.a, nil, ADD)
		denominator = OpOnBig(p.y, nil, big.NewInt(int64(2)), MUL)
	} else {
		numerator = OpOnBig(other.y, p.y, nil, SUB)
		denominator = OpOnBig(other.x, p.x, nil, SUB)
	}

	slope := OpOnBig(numerator, denominator, nil, DIV)

	slopeSqrt := OpOnBig(slope, nil, big.NewInt(int64(2)), EXP)
	x3 := OpOnBig(OpOnBig(slopeSqrt, p.x, nil, SUB), other.x, nil, SUB)
	x3MinusX1 := OpOnBig(x3, p.x, nil, SUB)

	y3 := OpOnBig(OpOnBig(slope, x3MinusX1, nil, MUL),
		p.y, nil, ADD)
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
	return fmt.Sprintf("(X: %v, y: %v, a: %v, b: %v)",
		xString, yString, p.a.String(), p.b.String())
}

func (p *Point) Equal(other *Point) bool {

	if p.a.EqualTo(other.a) == true && p.b.EqualTo(other.b) == true &&
		p.y.EqualTo(other.y) == true && p.x.EqualTo(other.x) == true {
		return true
	}

	return false
}

func (p *Point) NotEqual(other *Point) bool {
	if p.a.EqualTo(other.a) != true || p.b.EqualTo(other.b) != true ||
		p.y.EqualTo(other.y) != true || p.x.EqualTo(other.x) != true {
		return true
	}

	return false
}
