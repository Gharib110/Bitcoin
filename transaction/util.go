package transaction

import (
	"bufio"
	"fmt"
	"math/big"

	"github.com/tsuna/endian"
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

func ReadVariant(reader *bufio.Reader) *big.Int {
	/*
		1. check the byte after the version, < 0xfd,
		then the value of the byte is the count of input

		2, if the byte value >=0xfd < fe, read the following 2 bytes as the count of input

		3, if the byte following the version is >=0xfe < 0xff
		read the following 4 bytes as the count of input

		4, if the byte following version is == 0xff, we read the following 8 bytes as count
		of input
	*/
	i := make([]byte, 1)
	reader.Read(i)
	v := new(big.Int)
	v.SetBytes(i)
	if v.Cmp(big.NewInt(int64(0xfd))) < 0 {
		return v
	}

	if v.Cmp(big.NewInt(int64(0xfd))) == 0 {
		i1 := make([]byte, 2)
		reader.Read(i1)
		return LittleEndianToBigInt(i1, LittleEndian2Bytes)
	}

	if v.Cmp(big.NewInt(int64(0xfe))) == 0 {
		i1 := make([]byte, 4)
		reader.Read(i1)
		return LittleEndianToBigInt(i1, LittleEndian4Bytes)
	}

	i1 := make([]byte, 8)
	reader.Read(i1)
	return LittleEndianToBigInt(i1, LittleEndian8Bytes)
}

func EncodeVariant(v *big.Int) []byte {
	//if the value < 0xfd, one byte is enough
	if v.Cmp(big.NewInt(int64(0xfd))) < 0 {
		vBytes := v.Bytes()
		return []byte{vBytes[0]}
	} else if v.Cmp(big.NewInt(int64(0x10000))) < 0 {
		//if value >= 0xfd and < 0x10000, then need 2 bytes
		buf := []byte{0xfd}
		vBuf := BigIntToLittleEndian(v, LittleEndian2Bytes)
		buf = append(buf, vBuf...)
		return buf
	} else if v.Cmp(big.NewInt(int64(0x100000000))) < 0 {
		//value >= 0xFFFF and <= 0xFFFFFFFF, then need 4 bytes
		buf := []byte{0xfe}
		vBuf := BigIntToLittleEndian(v, LittleEndian4Bytes)
		buf = append(buf, vBuf...)
		return buf
	}

	p := new(big.Int)
	p.SetString("10000000000000000", 16)
	if v.Cmp(p) < 0 {
		//need 8 bytes
		buf := []byte{0xff}
		vBuf := BigIntToLittleEndian(v, LittleEndian8Bytes)
		buf = append(buf, vBuf...)
		return buf
	}

	panic(fmt.Sprintf("integer too large: %x\n", v))
}
