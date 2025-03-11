package elliptic_curve

import (
	"fmt"
	"math/big"
	"os"
)

type FieldElement struct {
	order *big.Int
	num   *big.Int
}

// S256Field Construct a FieldElement of s^256 - 2^32 - 977
func S256Field(num *big.Int) *FieldElement {
	var op big.Int
	twoExp256 := op.Exp(big.NewInt(int64(2)), big.NewInt(int64(256)), nil)
	var op1 big.Int
	twoExp32 := op1.Exp(big.NewInt(int64(3)), big.NewInt(int64(32)), nil)
	var op2 big.Int
	p := op2.Sub(twoExp256, twoExp32)
	var op3 big.Int
	pp := op3.Sub(p, big.NewInt(int64(977)))

	return NewFieldElement(pp, num)
}

func NewFieldElement(order *big.Int, num *big.Int) *FieldElement {
	if order.Cmp(num) == -1 {
		err := fmt.Sprintf("num %s is greater than order %s",
			num.String(), order.String())
		fmt.Println(err)
		os.Exit(1)
	}

	return &FieldElement{order, num}
}

func (el *FieldElement) String() string {
	return fmt.Sprintf("FieldElement{order: %s, num: %s}",
		el.order.String(), el.num.String())
}

func (el *FieldElement) EqualTo(other *FieldElement) bool {
	return el.order.Cmp(other.order) == 0 && el.num.Cmp(other.num) == 0
}

func (el *FieldElement) Add(other *FieldElement) *FieldElement {
	el.CheckOrder(other)

	var op big.Int
	return &FieldElement{el.order,
		op.Mod(op.Add(el.num, other.num), el.order)}
}

func (el *FieldElement) Negate() *FieldElement {
	var op big.Int
	return &FieldElement{el.order, op.Sub(el.order, el.num)}
}

func (el *FieldElement) Sub(other *FieldElement) *FieldElement {
	el.CheckOrder(other)
	return el.Add(other.Negate())
}

func (el *FieldElement) CheckOrder(other *FieldElement) {
	if el.order.Cmp(other.order) != 0 {
		err := fmt.Sprintf("Adding different order of %d and %d", el.order, other.order)
		fmt.Println(err)
		os.Exit(1)
	}
}

func (el *FieldElement) Mul(other *FieldElement) *FieldElement {
	el.CheckOrder(other)
	var op big.Int
	mul := op.Mul(el.num, other.num)
	return NewFieldElement(el.order, op.Mod(mul, el.order))
}

// Pow Arithmetic power over modular of the order
// Pow k ^ (p - 1) % p = 1, power > p - 1 => power % (p - 1)
func (el *FieldElement) Pow(pow *big.Int) *FieldElement {
	var op big.Int
	t := op.Mod(pow, op.Sub(el.order, big.NewInt(int64(1))))
	powRes := op.Exp(el.num, t, el.order)
	//modRes := op.Mod(powRes, el.order)
	return NewFieldElement(el.order, powRes)
}

func (el *FieldElement) ScalarMul(r *big.Int) *FieldElement {
	var op big.Int
	res := op.Mul(el.num, r)
	res = op.Mod(res, el.order)
	return NewFieldElement(el.order, res)
}

func (el *FieldElement) Div(other *FieldElement) *FieldElement {
	el.CheckOrder(other)
	var op big.Int
	otherReverse := other.Pow(op.Sub(el.order, big.NewInt(int64(2))))
	return el.Mul(otherReverse)
}
