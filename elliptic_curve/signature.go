package elliptic_curve

import (
	"fmt"
)

type Signature struct {
	r *FieldElement
	s *FieldElement
}

func NewSignature(r, s *FieldElement) *Signature {
	return &Signature{r: r, s: s}
}

func (sig *Signature) String() string {
	return fmt.Sprintf("Signature(r: %v, s: %v)",
		sig.r, sig.s)
}

/*
DER
1. Set the first byte to 0x30
2. second byte is the total length of s and r
3. the first byte is 0x02 that indicate beginning of the byte array for r
4. if the first byte of r is => 0x80, then we need to append 0x00 as the beginning byte
of the byte array of r, compute the len of the byte array of r and append the length
behind the 0x02 of step 2
5. insert 0x02 behind the last byte of the r byte array, as indicator for the beginning
6. do the same for s as step 4

total length of 0x44 or 0x45
*/
func (sig *Signature) DER() []byte {
	rBin := sig.r.num.Bytes()
	//if the first byte >= 0x80, append 0x00 at the beginning
	if rBin[0] >= 0x80 {
		rBin = append([]byte{0x00}, rBin...)
	}
	rBin = append([]byte{0x02, byte(len(rBin))}, rBin...)
	sBin := sig.s.num.Bytes()
	if sBin[0] >= 0x80 {
		sBin = append([]byte{0x80}, sBin...)
	}
	sBin = append([]byte{0x02, byte(len(sBin))}, sBin...)
	derBin := append([]byte{0x30, byte(len(rBin) + len(sBin))}, rBin...)
	derBin = append(derBin, sBin...)

	return derBin
}
