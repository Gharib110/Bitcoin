package elliptic_curve

import (
	"fmt"
)

type Signature struct {
	r *FieldElement
	s *FieldElement
}

func NewSignature(r *FieldElement, s *FieldElement) *Signature {
	return &Signature{
		r: r,
		s: s,
	}
}

func (s *Signature) String() string {
	return fmt.Sprintf("Signature(r: {%s}, s: {%s})", s.r, s.s)
}

/*
DER:
1. set the first byte to 0x30

2, second byte is the total length of s and r

3, the first byte is 0x02 it is an indicator for the beginning of the byte array for r

4. if the first byte of r is >= 0x80, then we need to append 0x00 as beginning byte
of the bytes array of r, compute the length of the bytes array of r and append the length
behind the 0x02 of step 2

5. insert 0x02 behind the last byte of the r byte array, as indicator for the beginning of s

6. do the same for s as step 4

total length of 0x44 or 0x45

30 45 02 21 00 ed81ff192e75a3fd2304004dcadb746fa5e24c5031ccfcf21320b0277457c98f 02 20 7a986d955c6e0cb35d446a89d3f56100f4d7f67801c31967743a9c8e10615bed

first byte is 0x30
second byte 0x45 is the total length of r and s,
third byte 02 is indicator for beginning of r
fourth byte 0x21 is the total length of r
the fifth byte is 00, because the first byte of r is 0xed >= 0x80
r: ed81ff192e75a3fd2304004dcadb746fa5e24c5031ccfcf21320b0277457c98f
following byte 0x02 is indicator for the beginning of s
following byte 0x20 is length of s
the first byte of s is 7a < 0x80, we don't insert 0x00 at the beginning

s: 7a986d955c6e0cb35d446a89d3f56100f4d7f67801c31967743a9c8e10615bed
*/

func (s *Signature) Der() []byte {
	rBin := s.r.num.Bytes()
	//if the first byte >= 0x80, append 0x00 at the beginning
	if rBin[0] >= 0x80 {
		rBin = append([]byte{0x00}, rBin...)
	}
	//insert indicator 0x02 and the length of rBin
	rBin = append([]byte{0x02, byte(len(rBin))}, rBin...)
	//do the same to s
	sBin := s.s.num.Bytes()
	//if the first byte of s >= 0x80, append 0x00 at the beginning of s
	if sBin[0] >= 0x80 {
		sBin = append([]byte{0x00}, sBin...)
	}
	//insert indicator 0x02 and the length of sBin
	sBin = append([]byte{0x02, byte(len(sBin))}, sBin...)
	//combine rBin, sBin and insert 0x30 and the total length of sBin rBin at the beginning
	derBin := append([]byte{0x30, byte(len(rBin) + len(sBin))}, rBin...)
	derBin = append(derBin, sBin...)

	return derBin
}
