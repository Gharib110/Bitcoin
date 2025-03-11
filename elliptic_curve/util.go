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
