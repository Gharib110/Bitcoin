package transaction

import (
	"bufio"
	"bytes"
	"fmt"
	ecc "github.com/Gharib110/Bitcoin/elliptic_curve"
	"io"
	"math/big"
)

type Block struct {
	version         []byte
	previousBlockID []byte
	merkleRoot      []byte
	timeStamp       []byte
	bits            []byte
	nonce           []byte
}

const (
	TwoWeeks = 60 * 60 * 24 * 14
)

func ComputeNewTarget(firstBlockBytes []byte, lastBlockBytes []byte) *big.Int {
	firstBlock := ParseBlock(firstBlockBytes)
	lastBlock := ParseBlock(lastBlockBytes)

	firstBlockTime := new(big.Int)
	firstBlockTime.SetBytes(firstBlock.timeStamp)

	lastBlockTime := new(big.Int)
	lastBlockTime.SetBytes(lastBlock.timeStamp)

	var opSub big.Int
	timeDifferential := opSub.Sub(lastBlockTime, firstBlockTime)
	if timeDifferential.Cmp(big.NewInt(TwoWeeks*4)) > 0 {
		timeDifferential = big.NewInt(TwoWeeks * 4)
	}
	if timeDifferential.Cmp(big.NewInt(TwoWeeks/4)) < 0 {
		timeDifferential = big.NewInt(TwoWeeks / 4)
	}

	var opMul big.Int
	var opDiv big.Int
	newTarget := opDiv.Div(opMul.Mul(lastBlock.Target(), timeDifferential), big.NewInt(TwoWeeks))
	return newTarget
}

func TargetToBits(target *big.Int) []byte {
	targetBytes := target.Bytes()
	exponent := len(targetBytes)
	coefficient := targetBytes[0:3]
	bits := make([]byte, 0)
	bits = append(bits, reverseByteSlice(coefficient)...)
	bits = append(bits, byte(exponent))
	return bits
}

func ParseBlock(rawBlock []byte) *Block {
	block := &Block{}

	reader := bytes.NewReader(rawBlock)
	bufReader := bufio.NewReader(reader)

	buffer := make([]byte, 4)
	io.ReadFull(bufReader, buffer)
	block.version = reverseByteSlice(buffer)

	buffer = make([]byte, 32)
	io.ReadFull(bufReader, buffer)
	block.previousBlockID = reverseByteSlice(buffer)

	buffer = make([]byte, 32)
	io.ReadFull(bufReader, buffer)
	block.merkleRoot = reverseByteSlice(buffer)

	buffer = make([]byte, 4)
	io.ReadFull(bufReader, buffer)
	block.timeStamp = reverseByteSlice(buffer)

	buffer = make([]byte, 4)
	io.ReadFull(bufReader, buffer)
	block.bits = buffer

	buffer = make([]byte, 4)
	io.ReadFull(bufReader, buffer)
	block.nonce = buffer

	return block
}

func (b *Block) Serialize() []byte {
	result := make([]byte, 0)
	//four bytes version in little endian format
	version := new(big.Int)
	version.SetBytes(b.version)
	result = append(result, BigIntToLittleEndian(version, LittleEndian4Bytes)...)
	//previous block header hash and merkle root in little endian
	result = append(result, reverseByteSlice(b.previousBlockID)...)
	result = append(result, reverseByteSlice(b.merkleRoot)...)

	timeStamp := new(big.Int)
	timeStamp.SetBytes(b.timeStamp)
	result = append(result, BigIntToLittleEndian(timeStamp, LittleEndian4Bytes)...)

	result = append(result, b.bits...)
	result = append(result, b.nonce...)

	return result
}

func (b *Block) Hash() []byte {
	s := b.Serialize()
	sha := ecc.Hash256(string(s))
	return reverseByteSlice(sha)
}

func (b *Block) String() string {
	s := fmt.Sprintf("version:%x\nprevious block id:%x\nmerkle root:%x\ntimestamp:%x\nbits:%x\nnonce:%x\nhash:%x\n",
		b.version, b.previousBlockID, b.merkleRoot, b.timeStamp, b.bits, b.nonce, b.Hash())

	return s
}

/*

bits in the field of version used to indicate which protocol the current miner can support
BIP0009, BIP0091, BIP0141
*/

func (b *Block) Bip9() bool {
	version := new(big.Int)
	version.SetBytes(b.version)
	ver := version.Int64()

	return (ver >> 29) == 0b001
}

func (b *Block) Bip91() bool {
	version := new(big.Int)
	version.SetBytes(b.version)
	ver := version.Int64()

	return (ver >> 4 & 1) == 1
}

func (b *Block) Bip141() bool {
	version := new(big.Int)
	version.SetBytes(b.version)
	ver := version.Int64()
	return (ver >> 1 & 1) == 1
}

func (b *Block) Target() *big.Int {
	//exponent, last byte - 3
	var opSub big.Int
	exponentPart := opSub.Sub(big.NewInt(int64(b.bits[len(b.bits)-1])), big.NewInt(3))
	//left most significant three bytes
	coefficientBuf := b.bits[0 : len(b.bits)-1]
	coefficientBytes := reverseByteSlice(coefficientBuf)
	coefficient := new(big.Int)
	coefficient.SetBytes(coefficientBytes)

	var opPow big.Int
	var opMul big.Int
	exponent := opPow.Exp(big.NewInt(256), exponentPart, nil)
	result := opMul.Mul(coefficient, exponent)
	return result
}

func (b *Block) Difficulty() *big.Int {
	//difficulty = 0xffff * 256^(0x1d-3) / target
	target := b.Target()
	var opMul big.Int
	var opExp big.Int
	var opDiv big.Int
	numerator := opMul.Mul(big.NewInt(0xffff), opExp.Exp(big.NewInt(256), big.NewInt(0x1d-3), nil))
	denominator := target
	difficulty := opDiv.Div(numerator, denominator)
	return difficulty
}
