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

func NewECPoint(x, y, a, b *big.Int) *Point  {
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

func (p *Point) Equal(other *Point) bool  {
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