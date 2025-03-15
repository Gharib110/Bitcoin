package networking

import (
	tx "github.com/Gharib110/Bitcoin/transaction"
	"math/big"
	"math/rand"
	"time"
)

type VersionMessage struct {
	command          string
	version          *big.Int
	services         *big.Int
	timestamp        *big.Int
	receiverServices *big.Int
	receiverIP       []byte
	receiverPort     uint16
	senderServices   *big.Int
	senderIP         []byte
	senderPort       uint16
	nonce            []byte
	userAgent        string
	latestBlock      *big.Int
	relay            bool
}

func NewVersionMessage() *VersionMessage {
	nonceBuf := make([]byte, 8)
	rand.Read(nonceBuf)
	return &VersionMessage{
		command:          "version",
		version:          big.NewInt(70015),
		services:         big.NewInt(0),
		timestamp:        big.NewInt(time.Now().Unix()),
		receiverServices: big.NewInt(0),
		receiverIP:       []byte{0x00, 0x00, 0x00, 0x00},
		receiverPort:     8333, //default for main-net nodes
		senderServices:   big.NewInt(0),
		senderIP:         []byte{0x00, 0x00, 0x00, 0x00},
		senderPort:       8333,
		nonce:            nonceBuf,
		userAgent:        "golang_bitcoin_lib",
		latestBlock:      big.NewInt(0),
		relay:            false,
	}
}

func (v *VersionMessage) Command() string {
	return v.command
}

func (v *VersionMessage) Serialize() []byte {
	result := make([]byte, 0)
	result = append(result, tx.BigIntToLittleEndian(v.version, tx.LittleEndian4Bytes)...)
	result = append(result, tx.BigIntToLittleEndian(v.services, tx.LittleEndian8Bytes)...)
	result = append(result, tx.BigIntToLittleEndian(v.timestamp, tx.LittleEndian8Bytes)...)
	result = append(result, tx.BigIntToLittleEndian(v.receiverServices, tx.LittleEndian8Bytes)...)
	//for ip with ip4 format, we need to append 12 bytes 0x00…ffff as prefix
	//then follow four bytes for ip4 address
	ipBuf := make([]byte, 16)
	for i := 0; i < 12; i++ {
		if i < 10 {
			ipBuf = append(ipBuf, 0x00)
		} else {
			ipBuf = append(ipBuf, 0xff)
		}
	}
	ipBuf = append(ipBuf, v.receiverIP...)
	result = append(result, ipBuf...)

	result = append(result, big.NewInt(int64(v.receiverPort)).Bytes()...)
	result = append(result, tx.BigIntToLittleEndian(v.senderServices, tx.LittleEndian8Bytes)...)
	ipBuf = make([]byte, 16)
	for i := 0; i < 12; i++ {
		if i < 10 {
			ipBuf = append(ipBuf, 0x00)
		} else {
			ipBuf = append(ipBuf, 0xff)
		}
	}
	ipBuf = append(ipBuf, v.senderIP...)
	result = append(result, ipBuf...)
	result = append(result, big.NewInt(int64(v.senderPort)).Bytes()...)

	result = append(result, v.nonce...)
	agentLen := tx.EncodeVariant(big.NewInt(int64(len(v.userAgent))))
	result = append(result, agentLen...)
	result = append(result, []byte(v.userAgent)...)
	result = append(result, tx.BigIntToLittleEndian(v.latestBlock, tx.LittleEndian4Bytes)...)

	if v.relay {
		result = append(result, 0x01)
	} else {
		result = append(result, 0x00)
	}

	return result
}

type VerAckMessage struct {
	command string
}

func NewVerAckMessage() *VerAckMessage {
	return &VerAckMessage{
		command: "verack",
	}
}

func (v *VerAckMessage) Command() string {
	return v.command
}

func (v *VerAckMessage) Serialize() []byte {
	return []byte{}
}
