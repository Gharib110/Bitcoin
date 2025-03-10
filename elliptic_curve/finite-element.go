package elliptic_curve

import (
	"fmt"
	"math"
	"os"
)

type FieldElement struct {
	order uint64
	num   uint64
}

func NewFieldElement(order uint64, num uint64) *FieldElement {
	if num >= order {
		err := fmt.Sprintf("num %d is greater than order %d", num, order)
		fmt.Println(err)
		os.Exit(1)
	}

	return &FieldElement{order, num}
}

func (el *FieldElement) String() string {
	return fmt.Sprintf("FieldElement{order: %d, num: %d}", el.order, el.num)
}

func (el *FieldElement) EqualTo(other *FieldElement) bool {
	return el.order == other.order && el.num == other.num
}

func (el *FieldElement) Add(other *FieldElement) *FieldElement {
	el.CheckOrder(other)

	return &FieldElement{el.order, (el.num + other.num) % el.order}
}

func (el *FieldElement) Negate() *FieldElement {
	return &FieldElement{el.order, el.order - el.num}
}

func (el *FieldElement) Sub(other *FieldElement) *FieldElement {
	el.CheckOrder(other)
	return el.Add(other.Negate())
}

func (el *FieldElement) CheckOrder(other *FieldElement) {
	if el.order != other.order {
		err := fmt.Sprintf("Adding different order of %d and %d", el.order, other.order)
		fmt.Println(err)
		os.Exit(1)
	}
}

func (el *FieldElement) Mul(other *FieldElement) *FieldElement {
	el.CheckOrder(other)
	return NewFieldElement(el.order, (el.num*other.num)%el.order)
}

func (el *FieldElement) Pow(pow int64) *FieldElement {
	return NewFieldElement(el.order,
		uint64(math.Pow(float64(el.num), float64(pow)))%el.order)
}

func (el *FieldElement) ScalarMul(r uint64) *FieldElement {
	return NewFieldElement(el.order, (r*el.num)%el.order)
}
