package merkletree

import (
	"bufio"
	"bytes"
	"fmt"
	"github.com/Gharib110/Bitcoin/transaction"
	"math/big"
	"strings"

	"golang.org/x/example/hello/reverse"
)

/*
command : merkle-block
binary for command of merkle-block:
00000020

df3b053dc46f162a9b00c7f0d5124e2676d47bbe7c5d0793a500000000000000

ef445fef2ed495c275892206ca533e7411907971013ab83e3b47bd0d692d14d4

dc7c835b

67d8001a

c157e670

bf0d0000

0a

ba412a0d1480e370173072c9562becffe87aa661c1e4a6dbc305d38ec5dc088a7cf92e6458aca7b32edae818f9c2c
98c37e06bf72ae0ce80649a38655ee1e27d34d9421d940b16732f24b94023e9d572a7f9ab8023434a4feb532d
2adfc8c2c2158785d1bd04eb99df2e86c54bc13e139862897217400def5d72c280222c4cbaee7261831e1550d
bb8fa82853e9fe506fc5fda3f7b919d8fe74b6282f92763cef8e625f977af7c8619c32a369b832bc2d051ecd9c
73c51e76370ceabd4f25097c256597fa898d404ed53425de608ac6bfe426f6e2bb457f1c554866eb69dcb8d6bf
6f880e9a59b3cd053e6c7060eeacaacf4dac6697dac20e4bd3f38a2ea2543d1ab7953e3430790a9f81e1c67f5b
58c825acf46bd02848384eebe9af917274cdfbb1a28a5d58a23a17977def0de10d644258d9c54f886d47d293a41
1cb622610

3

b55635

1, first 4 bytes in LE(little endian): 00000020 version

2, 32 bytes in LE is previous block hash,or id:
df3b053dc46f162a9b00c7f0d5124e2676d47bbe7c5d0793a500000000000000

3, 32 bytes in LE merkle root:
ef445fef2ed495c275892206ca533e7411907971013ab83e3b47bd0d692d14d4

4, 4 bytes is timestamp: dc7c835b

5, 4 bytes is named bits: 67d8001a

6, 4 bytes nonce: c157e670

7, 4 bytes in LE is number of total transaction of the block that contains the
given transaction: bf0d0000

8, variant int: number of hashes, blue boxes:0a

9, following chunk of data are hash values of the blue boxes:
ba412a0d1480e370173072c9562becffe87aa661c1e4a6dbc305d38ec5dc088a7cf92e6458aca7b32edae818f9c2c
98c37e06bf72ae0ce80649a38655ee1e27d34d9421d940b16732f24b94023e9d572a7f9ab8023434a4feb532d
2adfc8c2c2158785d1bd04eb99df2e86c54bc13e139862897217400def5d72c280222c4cbaee7261831e1550d
bb8fa82853e9fe506fc5fda3f7b919d8fe74b6282f92763cef8e625f977af7c8619c32a369b832bc2d051ecd9c
73c51e76370ceabd4f25097c256597fa898d404ed53425de608ac6bfe426f6e2bb457f1c554866eb69dcb8d6bf
6f880e9a59b3cd053e6c7060eeacaacf4dac6697dac20e4bd3f38a2ea2543d1ab7953e3430790a9f81e1c67f5b
58c825acf46bd02848384eebe9af917274cdfbb1a28a5d58a23a17977def0de10d644258d9c54f886d47d293a41
1cb622610

10:
variant int: length of flags: 3


11. flags: b55635, it will tell the position for the blue boxes
*/

type MerkleBlock struct {
	version           *big.Int
	previousBlock     []byte
	merkleRoot        []byte
	timeStamp         *big.Int
	bits              []byte
	nonce             []byte
	totalTransactions *big.Int
	numHashes         *big.Int
	hashes            [][]byte
	flagBits          []byte
}

func ErrorPanic(err error, msg string) {
	if err != nil {
		panic(msg)
	}
}

// BytesToBitsField convert bytes to bits
func BytesToBitsField(bytes []byte) []string {
	flagBits := make([]string, 0)
	for _, byteVal := range bytes {
		flagBits = append(flagBits, reverse.String(fmt.Sprintf("%08b", byteVal)))
	}

	return flagBits
}

func ParseMerkleBlock(payload []byte) *MerkleBlock {
	merkleBlock := &MerkleBlock{}
	reader := bytes.NewReader(payload)
	bufReader := bufio.NewReader(reader)
	version := make([]byte, 4)
	_, err := bufReader.Read(version)
	ErrorPanic(err, "MerkleBlock read version")
	merkleBlock.version = transaction.LittleEndianToBigInt(version, transaction.LittleEndian4Bytes)

	prevBlock := make([]byte, 32)
	_, err = bufReader.Read(prevBlock)
	ErrorPanic(err, "MerkleBlock read previous block")
	merkleBlock.previousBlock = transaction.ReverseByteSlice(prevBlock)

	merkleRoot := make([]byte, 32)
	_, err = bufReader.Read(merkleRoot)
	ErrorPanic(err, "MerkleBlock read merkle root")
	merkleBlock.merkleRoot = transaction.ReverseByteSlice(merkleRoot)

	timeStamp := make([]byte, 4)
	_, err = bufReader.Read(timeStamp)
	ErrorPanic(err, "MerkleBlock read timestamp")
	merkleBlock.timeStamp = transaction.LittleEndianToBigInt(timeStamp, transaction.LittleEndian4Bytes)

	bits := make([]byte, 4)
	_, err = bufReader.Read(bits)
	ErrorPanic(err, "MerkleBlock read bits")
	merkleBlock.bits = bits

	nonce := make([]byte, 4)
	_, err = bufReader.Read(nonce)
	ErrorPanic(err, "MerkleBlock read nonce")
	merkleBlock.nonce = nonce

	total := make([]byte, 4)
	_, err = bufReader.Read(total)
	ErrorPanic(err, "MerkleBlock read total")
	merkleBlock.totalTransactions = transaction.LittleEndianToBigInt(total, transaction.LittleEndian4Bytes)

	numHashes := transaction.ReadVariant(bufReader)
	merkleBlock.numHashes = numHashes

	hashes := make([][]byte, 0)
	for i := 0; i < int(numHashes.Int64()); i++ {
		hash := make([]byte, 32)
		_, err = bufReader.Read(hash)
		ErrorPanic(err, "MerkleBlock read hash")
		hashes = append(hashes, transaction.ReverseByteSlice(hash))
	}
	merkleBlock.hashes = hashes

	flagLen := transaction.ReadVariant(bufReader)
	flags := make([]byte, flagLen.Int64())
	_, err = bufReader.Read(flags)
	ErrorPanic(err, "MerkleBlock read flags")
	merkleBlock.flagBits = flags

	return merkleBlock

}

func (m *MerkleBlock) String() string {
	result := make([]string, 0)
	result = append(result, fmt.Sprintf("version:%x", m.version))
	result = append(result, fmt.Sprintf("previous block:%x", m.previousBlock))
	result = append(result, fmt.Sprintf("merkle root:%x", m.merkleRoot))
	bitsString := strings.Join(BytesToBitsField(m.bits), ",")
	result = append(result, fmt.Sprintf("bits: %s", bitsString))
	result = append(result, fmt.Sprintf("nonce: %x", m.nonce))
	result = append(result, fmt.Sprintf("total tx:%x", m.totalTransactions))
	result = append(result, fmt.Sprintf("number of hashes:%d", m.numHashes.Int64()))
	for i := 0; i < int(m.numHashes.Int64()); i++ {
		result = append(result, fmt.Sprintf("%x", m.hashes[i]))
	}

	flagToBits := strings.Join(BytesToBitsField(m.flagBits), "")
	result = append(result, fmt.Sprintf("flags: %s", flagToBits))

	return strings.Join(result, "\n")
}

func (m *MerkleBlock) IsValid() bool {
	flagBits := strings.Join(BytesToBitsField(m.flagBits), "")
	merkleTree := InitEmptyMerkleTree(int(m.totalTransactions.Int64()))
	//when compute the merkle root, we need all hash value in little endian format
	hashes := make([][]byte, 0)
	for _, hash := range m.hashes {
		hashes = append(hashes, transaction.ReverseByteSlice(hash))
	}
	merkleTree.PopulateTree(flagBits, hashes)
	fmt.Printf("%x\n", merkleTree.Root())
	//need to reverse the byte order of the merkle root
	return bytes.Equal(m.merkleRoot, transaction.ReverseByteSlice(merkleTree.Root()))
}
