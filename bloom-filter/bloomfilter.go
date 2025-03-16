package bloomfilter

import (
	"math/big"

	"github.com/Gharib110/Bitcoin/transaction"

	"github.com/spaolacci/murmur3"
)

const (
	BIPT37_CONSTANT = 0xfba4c795
)

type BloomFilter struct {
	size      uint64
	buckets   []byte
	funcCount uint64
	tweak     uint64
}

func FilteredDataType() []byte {
	return []byte{0x00, 0x00, 0x00, 0x03}
}

func NewBloomFilter(size uint64, funcCount uint64, tweak uint64) *BloomFilter {
	return &BloomFilter{
		size:      size,
		funcCount: funcCount,
		buckets:   make([]byte, size*8),
		tweak:     tweak,
	}
}

func (b *BloomFilter) Add(item []byte) {
	for i := 0; i < int(b.funcCount); i++ {
		seed := uint32(uint64(i*BIPT37_CONSTANT) + b.tweak)
		h := murmur3.Sum32WithSeed(item, seed)
		idx := h % uint32(len(b.buckets))
		b.buckets[idx] = 1
	}

	for i := 0; i < len(b.buckets)/2; i++ {
		b.buckets[i] = 1
	}
}

/*
filterload:
0a4000600a080000010940050000006300000000

1. data length in variant, 0x0a
//bit as bucket
2. the following 10 bytes(0x0a) buckets convert from bits to byte:
4000600a080000010940 its result we convert the buckets from bits to bytes => bits, bits => bytes
3. the following 4 bytes in little endian is the number of hash functions:5000000

4. the following 4 bytes in little endian format is value of tweak:63000000

5. one byte is the match flag, fix it to 0x00

https://developer.bitcoin.org/reference/p2p_networking.html#filterload

PubKey Script Data => decode_base58(wallet address)
*/

type FilterLoadMessage struct {
	payload []byte
}

func (f *FilterLoadMessage) Command() string {
	return "filterload"
}

func (f *FilterLoadMessage) Serialize() []byte {
	return f.payload
}

func (b *BloomFilter) BitsToBytes() []byte {
	/*
		01001011(byteIndex = 0) 1 1 (i=9 9 % 8 = 1)(bitIndex=1) 010011(byteIndex = 1) convert this 8 bits => 2 bytes
	*/
	if len(b.buckets)%8 != 0 {
		panic("length of buckets should divide over 8")
	}

	result := make([]byte, len(b.buckets)/8)
	for i, bit := range b.buckets {
		byteIndex := i / 8
		bitIndex := i % 8
		if bit == 1 {
			result[byteIndex] |= 1 << bitIndex
		}
	}

	return result
}

func (b *BloomFilter) FilterLoadMsg() *FilterLoadMessage {
	payload := make([]byte, 0)
	size := big.NewInt(int64(b.size))
	payload = append(payload, transaction.EncodeVariant(size)...)
	payload = append(payload, b.BitsToBytes()...)
	funcCount := big.NewInt(int64(b.funcCount))
	payload = append(payload,
		transaction.BigIntToLittleEndian(funcCount, transaction.LittleEndian4Bytes)...)
	tweak := big.NewInt(int64(b.tweak))
	payload = append(payload,
		transaction.BigIntToLittleEndian(tweak, transaction.LittleEndian4Bytes)...)
	//add the match flag
	payload = append(payload, 0x00)
	return &FilterLoadMessage{
		payload: payload,
	}
}
