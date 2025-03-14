package elliptic_curve

import (
	"bufio"
	"bytes"
	"crypto/sha256"
	"github.com/tsuna/endian"
	"golang.org/x/crypto/ripemd160"
	"math/big"
)

type LittleEndianLength int

const (
	LittleEndian2Bytes = iota
	LittleEndian4Bytes
	LittleEndian8Bytes
)

func BigIntToLittleEndian(v *big.Int, length LittleEndianLength) []byte {
	switch length {
	case LittleEndian2Bytes:
		val := v.Int64()
		littleEndianVal := endian.HostToNetUint16(uint16(val))
		p := big.NewInt(int64(littleEndianVal))
		return p.Bytes()

	case LittleEndian4Bytes:
		val := v.Int64()
		littleEndianVal := endian.HostToNetUint32(uint32(val))
		p := big.NewInt(int64(littleEndianVal))
		return p.Bytes()

	case LittleEndian8Bytes:
		val := v.Int64()
		littleEndianVal := endian.HostToNetUint64(uint64(val))
		p := big.NewInt(int64(littleEndianVal))
		return p.Bytes()
	}

	return nil
}

func LittleEndianToBigInt(bytes []byte, length LittleEndianLength) *big.Int {
	switch length {
	case LittleEndian2Bytes:
		p := new(big.Int)
		p.SetBytes(bytes)
		val := endian.NetToHostUint16(uint16(p.Uint64()))
		return big.NewInt(int64(val))

	case LittleEndian4Bytes:
		p := new(big.Int)
		p.SetBytes(bytes)
		val := endian.NetToHostUint32(uint32(p.Uint64()))
		return big.NewInt(int64(val))

	case LittleEndian8Bytes:
		p := new(big.Int)
		p.SetBytes(bytes)
		val := endian.NetToHostUint64(uint64(p.Uint64()))
		return big.NewInt(int64(val))
	}

	return nil
}

/*
z, sha256(sha256(create a text)) -> 256bits-> 32bytes integer
*/

func Hash256(text string) []byte {
	hashOnce := sha256.Sum256([]byte(text))
	hashTwice := sha256.Sum256(hashOnce[:])
	return hashTwice[:]
}

func GetGenerator() *Point {
	Gx := new(big.Int)
	Gx.SetString("79be667ef9dcbbac55a06295ce870b07029bfcdb2dce28d959f2815b16f81798", 16)
	Gy := new(big.Int)
	Gy.SetString("483ada7726a3c4655da4fbfc0e1108a8fd17b448a68554199c47d08ffb10d4b8", 16)
	G := S256Point(Gx, Gy)
	return G
}

func GetBitCoinValueN() *big.Int {
	n := new(big.Int)
	n.SetString("fffffffffffffffffffffffffffffffebaaedce6af48a03bbfd25e8cd0364141", 16)
	return n
}

func ParseSEC(secBin []byte) *Point {
	//check the first byte to decide it is compressed or uncompressed
	if secBin[0] == 4 {
		//uncompressed
		x := new(big.Int)
		x.SetBytes(secBin[1:33])
		y := new(big.Int)
		y.SetBytes(secBin[33:65])
		return S256Point(x, y)
	}

	//check first byte for y is odd or even
	isEven := secBin[0] == 2
	x := new(big.Int)
	x.SetBytes(secBin[1:])
	y2 := S256Field(x).Pow(big.NewInt(int64(3))).Add(S256Field(big.NewInt(int64(7))))
	y := y2.Sqrt()
	var modOp big.Int
	var yEven *FieldElement
	var yOdd *FieldElement
	if modOp.Mod(y.num, big.NewInt(int64(2))).Cmp(big.NewInt(int64(0))) == 0 {
		yEven = y
		yOdd = y.Negate() //p - y
	} else {
		yOdd = y
		yEven = y.Negate()
	}

	if isEven {
		return S256Point(x, yEven.num)
	} else {
		return S256Point(x, yOdd.num)
	}
}

/*
EncodeBase58
base58 it removes 0 O, l I
*/
func EncodeBase58(s []byte) string {
	Base58Alphabet := "123456789ABCDEFGHJKLMNPQRSTUVWXYZabcdefghijkmnopqrstuvwxyz"
	count := 0
	for idx := range s {
		if s[idx] == 0 {
			count += 1
		} else {
			break
		}
	}

	prefix := ""
	for i := 0; i < count; i++ {
		prefix += "1"
	}
	result := ""
	num := new(big.Int)
	num.SetBytes(s)
	for num.Cmp(big.NewInt(int64(0))) > 0 {
		var divOp big.Int
		var modOp big.Int
		mod := modOp.Mod(num, big.NewInt(int64(58)))
		num = divOp.Div(num, big.NewInt(int64(58)))
		result = string(Base58Alphabet[mod.Int64()]) + result
	}

	return prefix + result
}

func Base58Checksum(s []byte) string {
	hash256 := Hash256(string(s))
	return EncodeBase58(append(s, hash256[:4]...))
}

func Hash160(s []byte) []byte {
	sum256 := sha256.Sum256(s)
	hasher := ripemd160.New()
	hasher.Write(sum256[:])
	hashBytes := hasher.Sum(nil)

	return hashBytes
}

func ParseSigBin(sigBin []byte) *Signature {
	reader := bytes.NewReader(sigBin)
	bufReader := bufio.NewReader(reader)
	//first byte should be 0x30
	firstByte := make([]byte, 1)
	bufReader.Read(firstByte)
	if firstByte[0] != 0x30 {
		panic("Bad Signature, first byte is not 0x30")
	}
	//second byte is the length of r and s
	lenBuf := make([]byte, 1)
	bufReader.Read(lenBuf)
	//first two byte with the length of r and s should be the total length of sigBin
	if lenBuf[0]+2 != byte(len(sigBin)) {
		panic("Bad Signature length")
	}

	//marker 0x02 as the beginning of r
	marker := make([]byte, 1)
	bufReader.Read(marker)
	if marker[0] != 0x02 {
		panic("signature marker for r is not 0x02")
	}
	//following is the length of r bin
	lenBuf = make([]byte, 1)
	bufReader.Read(lenBuf)
	rLength := lenBuf[0]
	rBin := make([]byte, rLength)
	//it may have 0x00 append at the head, but it does not affect the value of r
	bufReader.Read(rBin)
	r := new(big.Int)
	r.SetBytes(rBin)
	//marketer 0x02 for the beginning of s
	marker = make([]byte, 1)
	bufReader.Read(marker)
	if marker[0] != 0x02 {
		panic("signature marker for s is not 0x02")
	}
	//following is length of s bin
	lenBuf = make([]byte, 1)
	bufReader.Read(lenBuf)
	sLength := lenBuf[0]
	sBin := make([]byte, sLength)
	bufReader.Read(sBin)
	s := new(big.Int)
	s.SetBytes(sBin)
	if len(sigBin) != int(6+rLength+sLength) {
		panic("signature wrong length")
	}

	n := GetBitCoinValueN()
	return NewSignature(NewFieldElement(n, r), NewFieldElement(n, s))
}
