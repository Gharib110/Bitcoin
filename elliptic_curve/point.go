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
	x, y *big.Int
	a, b *big.Int
}

func OpOnBig(x, y *big.Int, op_type OpType) *big.Int {
	var op big.Int
	switch op_type {
	case ADD:
		return op.Add(x, y)
	case SUB:
		return op.Sub(x, y)
	case MUL:
		return op.Mul(x, y)
	case DIV:
		return op.Div(x, y)
	case EXP:
		return op.Exp(x, y, nil)
	}

	panic("unknown operation type")
}

func NewECPoint(x, y, a, b *big.Int) *Point {
	if x == nil && y == nil {
		return &Point{
			x: x,
			y: y,
			a: a,
			b: b,
		}
	}
	left := OpOnBig(y, big.NewInt(int64(2)), EXP)
	x3 := OpOnBig(x, big.NewInt(int64(3)), EXP)
	ax := OpOnBig(a, x, MUL)
	right := OpOnBig(OpOnBig(x3, ax, ADD), b, ADD)
	if left.Cmp(right) != 0 {
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
	if p.a.Cmp(other.a) != 0 && p.b.Cmp(other.b) != 0 {
		panic("points are not on the same curve")
	}

	if p.x == nil {
		return other
	}

	if other.x == nil {
		return p
	}

	if p.x.Cmp(other.x) == 0 &&
		OpOnBig(p.y, other.y, ADD).Cmp(big.NewInt(int64(0))) == 0 {
		return &Point{
			x: nil,
			y: nil,
			a: p.a,
			b: p.b,
		}
	}

	var numerator *big.Int
	var denominator *big.Int
	if p.x.Cmp(other.x) == 0 && p.y.Cmp(other.y) == 0 {
		xSqrt := OpOnBig(p.x, big.NewInt(int64(2)), EXP)
		threeXSqrt := OpOnBig(xSqrt, big.NewInt(int64(3)), MUL)
		numerator = OpOnBig(threeXSqrt, p.a, ADD)
		denominator = OpOnBig(p.y, big.NewInt(int64(2)), MUL)
	} else {
		numerator = OpOnBig(other.y, p.y, SUB)
		denominator = OpOnBig(other.x, p.x, SUB)
	}

	slope := OpOnBig(numerator, denominator, DIV)

	slopeSqrt := OpOnBig(slope, big.NewInt(int64(2)), EXP)
	x3 := OpOnBig(OpOnBig(slopeSqrt, p.x, SUB), other.x, SUB)
	x3MinusX1 := OpOnBig(x3, p.x, SUB)

	y3 := OpOnBig(OpOnBig(slope, x3MinusX1, MUL), p.y, ADD)
	minusY3 := OpOnBig(y3, big.NewInt(int64(-1)), MUL)

	return &Point{
		x: x3,
		y: minusY3,
		a: p.a,
		b: p.b,
	}
}

func (p *Point) String() string {
	return fmt.Sprintf("(X: %v, y: %v, a: %v, b: %v)",
		p.x.String(), p.y.String(), p.a.String(), p.b.String())
}

func (p *Point) Equal(other *Point) bool {
	if p.a.Cmp(other.a) == 0 && p.b.Cmp(other.b) == 0 &&
		p.y.Cmp(other.y) == 0 && p.x.Cmp(other.x) == 0 {
		return true
	}

	return false
}

func (p *Point) NotEqual(other *Point) bool {
	if p.a.Cmp(other.a) != 0 || p.b.Cmp(other.b) != 0 ||
		p.y.Cmp(other.y) != 0 || p.x.Cmp(other.x) != 0 {
		return true
	}

	return false
}
