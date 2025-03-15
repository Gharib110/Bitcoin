package networking

import (
	"bufio"
	"bytes"
	"encoding/hex"
	"fmt"
	tx "github.com/Gharib110/Bitcoin/transaction"
	"math/big"
)

type GetHeaderMessage struct {
	command    string
	version    *big.Int
	numbHashes *big.Int
	startBlock []byte
	endBlock   []byte
}

func GetGenesisBlockHash() []byte {
	genesisBlockRawData, err := hex.DecodeString("0100000000000000000000000000000000000000000000000000000000000000000000003ba3edfd7a7b12b27ac72c3e67768f617fc81bc3888a51323a9fb8aa4b1e5e4a29ab5f49ffff001d1dac2b7c")
	if err != nil {
		panic(err)
	}
	genesisBlock := tx.ParseBlock(genesisBlockRawData)
	return genesisBlock.Hash()
}

func NewGetHeaderMessage(startBlock []byte) *GetHeaderMessage {
	return &GetHeaderMessage{
		command:    "getheaders",
		version:    big.NewInt(70015),
		numbHashes: big.NewInt(1),
		startBlock: startBlock,
		endBlock:   make([]byte, 32),
	}
}

func (g *GetHeaderMessage) Command() string {
	return g.command
}

func (g *GetHeaderMessage) Serialize() []byte {
	result := make([]byte, 0)
	result = append(result, tx.BigIntToLittleEndian(g.version, tx.LittleEndian4Bytes)...)
	result = append(result, tx.EncodeVariant(g.numbHashes)...)
	result = append(result, tx.ReverseByteSlice(g.startBlock)...)
	result = append(result, tx.ReverseByteSlice(g.endBlock)...)
	return result
}

/*
02

00000020df3b053dc46f162a9b00c7f0d5124e2676d47bbe7c5d0793a500000000000000ef445fef2ed495c275892206ca533e7411907971013ab83e3b47bd0d692d14d4dc7c835b67d8001ac157e67

00

00000002030eb2540c41025690160a1014c577061596e32e426b712c7ca00000000000000768b89f07044e6130ead292a3f51951adbd2202df447d98789339937fd006bd44880835b67d8001ade092046

00

1. variant value, only one byte here, number of block headers returned: 0x02

2. the following 80 bytes is the first block header:
00000020df3b053dc46f162a9b00c7f0d5124e2676d47bbe7c5d0793a500000000000000ef445fef2ed495c275892206ca533e7411907971013ab83e3b47bd0d692d14d4dc7c835b67d8001ac157e67

3. the following bytes in variant is the number of transactions for the given block,
it is always fixed to 0x00

4， the following 80 bytes is the second block header:
00000002030eb2540c41025690160a1014c577061596e32e426b712c7ca00000000000000768b89f07044e6130ead292a3f51951adbd2202df447d98789339937fd006bd44880835b67d8001ade092046


5. finally is variant value, but is always with one byte and value set to 0x00

*/

func LenOfVariant(val *big.Int) int {
	//returning how many bytes needed by the variant value
	shiftBytes := len(val.Bytes())
	if val.Cmp(big.NewInt(0xfd)) >= 0 {
		shiftBytes += 1
	}

	return shiftBytes
}

func ParseGetHeader(rawData []byte) []*tx.Block {
	reader := bytes.NewReader(rawData)
	bufReader := bufio.NewReader(reader)
	numHeads := tx.ReadVariant(bufReader)
	fmt.Printf("header count%d\n", numHeads)
	shiftBytes := LenOfVariant(numHeads)
	rawData = rawData[shiftBytes:]

	blocks := make([]*tx.Block, 0)
	for i := 0; i < int(numHeads.Int64()); i++ {
		block := tx.ParseBlock(rawData)
		blocks = append(blocks, block)

		rawData = rawData[len(block.Serialize()):]
		reader := bytes.NewReader(rawData)
		bufReader := bufio.NewReader(reader)
		numTxs := tx.ReadVariant(bufReader)
		if numTxs.Cmp(big.NewInt(0)) != 0 {
			panic("number of transaction is not 0")
		}

		shift := LenOfVariant(numTxs)
		if shift == 0 {
			/*
				a := big.NewInt(0), len(a.Bytes()) == 0
				big.Int will remove any prefix 0 in the buffer, if the value is 0,
				then the buffer will only contain 0, it will remove the 0 in the buffer
				which leads to 0 length of buffer
				buf := [0x00] => []
				b := make(big.Int)
				b.SetBytes(buf),
				len(b.Bytes()) => 0
			*/
			shift = 1
		}
		rawData = rawData[shift:]
	}

	return blocks
}
