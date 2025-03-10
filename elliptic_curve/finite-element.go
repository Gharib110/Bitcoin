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

func (el *FieldElement) Pow(pow *big.Int) *FieldElement {
	var op big.Int
	t := op.Mod(pow, op.Sub(el.order, big.NewInt(int64(1))))
	powRes := op.Exp(el.num, t, nil)
	modRes := op.Mod(powRes, el.order)
	return NewFieldElement(el.order, modRes)
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
