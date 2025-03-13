package elliptic_curve

import (
	"crypto/sha256"
	"math/big"
)

// Hash256  z, sha256(sha256(create a text0) -> 256bits integer -> 32bytes integer
func Hash256(text string) []byte {
	hashOnce := sha256.Sum256([]byte(text))
	hashTwice := sha256.Sum256(hashOnce[:])
	return hashTwice[:]
}

func GetGenerator() *Point {
	Gx := new(big.Int)
	Gy := new(big.Int)
	Gx.SetString(
		"79be667ef9dcbbac55a06295ce870b07029bfcdb2dce28d959f2815b16f81798",
		16)
	Gy.SetString(
		"483ada7726a3c4655da4fbfc0e1108a8fd17b448a68554199c47d08ffb10d4b8",
		16)
	G := S256Point(Gx, Gy)
	return G
}

func GetBitCoinValueN() *big.Int {
	n := new(big.Int)
	n.SetString(
		"fffffffffffffffffffffffffffffffebaaedce6af48a03bbfd25e8cd0364141",
		16)
	return n
}

func ParseSECEncoding(secBin []byte) *Point {
	//Check the first byte to decide it is compressed or uncompressed
	if secBin[0] == 4 {
		//uncompress
		x := new(big.Int)
		x.SetBytes(secBin[1:33])
		y := new(big.Int)
		y.SetBytes(secBin[33:65])

		return S256Point(x, y)
	}

	//check first byte for y is odd or even
	isEven := secBin[0] == 2
	x := new(big.Int)
	x.SetBytes(secBin[1:33])
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

// EncodeBase58 which removes 0 O, l I
func EncodeBase58(s []byte) string {
	base58Alphabets := "123456789ABCDEFGHJKLMNPQRSTUVWXYZabcdefghijkmnopqrstuvwxyz"
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
	for num.Cmp(big.NewInt(0)) > 0 {
		var divOp big.Int
		var modOp big.Int
		mod := modOp.Mod(num, big.NewInt(int64(58)))
		num = divOp.Div(num, big.NewInt(int64(58)))
		result = string(base58Alphabets[mod.Int64()]) + result
	}

	return prefix + result
}
