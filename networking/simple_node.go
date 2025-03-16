package networking

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"github.com/Gharib110/Bitcoin/bloom-filter"
	"github.com/Gharib110/Bitcoin/merkle-tree"
	"net"
	"time"
)

type Message interface {
	Command() string
	Serialize() []byte
}

type SimpleNode struct {
	host    string
	port    uint16
	testnet bool
}

func NewSimpleNode(host string, port uint16, testnet bool) *SimpleNode {
	return &SimpleNode{
		host:    host,
		port:    port,
		testnet: testnet,
	}
}

func (s *SimpleNode) Run() {
	/*
		using socket connect to given host with given port, then
		construct package with payload is a version message and send to
		the peer, waiting peer to send back its version message and verack,
		and we send verack back to peer and close the connection
	*/
	conStr := fmt.Sprintf("%s:%d", s.host, s.port)
	conn, err := net.Dial("tcp", conStr)
	if err != nil {
		panic(err)
	}
	s.WaitFor(conn)

	s.GetData(conn)
}

func (s *SimpleNode) GetData(conn net.Conn) {
	//prepare the bloom filter,
	txHash, err := hex.DecodeString("1df77b894e1910628714bb73df59e20fb9114f9dcc051d8c03ca197dd112cc8a")
	if err != nil {
		panic(err)
	}
	bf := bloomfilter.NewBloomFilter(30, 5, 90210)
	//set up the bloom filter map the transaction hash into buckets
	bf.Add(txHash)
	//send filterload command to fullnode
	s.Send(conn, bf.FilterLoadMsg())
	getData := bloomfilter.NewGetDataMessage()
	receiveMerkleBlock := false
	/*
		    ask full node to search all transactions in the given block with the given id,
			put all those transactions in the block through the filter we sent by using
			the filterload command, then collect all transactions that can pass through the
			filter and put them in to merkleblock command
	*/
	blockHash, err := hex.DecodeString("0000000000000138f016a6fc1666fd667b7d282d65ad14b7f0b16a75a2e90e50")
	getData.AddData(bloomfilter.FilteredDataType(), blockHash)
	s.Send(conn, getData)

	for !receiveMerkleBlock {
		time.Sleep(2 * time.Second)
		messages := s.Read(conn)
		for i := 0; i < len(messages); i++ {
			msg := messages[i]
			fmt.Printf("receiving command %s\n", msg.command)
			command := string(bytes.Trim(msg.command, "\x00"))

			if command == "merkleblock" {
				merkleBlock := merkletree.ParseMerkleBlock(msg.payload)
				fmt.Printf("merkleblock received: %s\n", merkleBlock)
				fmt.Printf("merkleblock valid: %v\n", merkleBlock.IsValid())
				receiveMerkleBlock = true
			}
		}
	}

}

func (s *SimpleNode) GetHeaders(conn net.Conn) {
	getHeaderMsg := NewGetHeaderMessage(GetGenesisBlockHash())
	s.Send(conn, getHeaderMsg)

	receivedGetHeader := false
	for !receivedGetHeader {
		//let the peer have a rest
		time.Sleep(2 * time.Second)
		messages := s.Read(conn)
		for i := 0; i < len(messages); i++ {
			msg := messages[i]
			fmt.Printf("receiving command:%s\n", msg.command)
			command := string(bytes.Trim(msg.command, "\x00"))
			if command == "headers" {
				receivedGetHeader = true
				blocks := ParseGetHeader(msg.payload)
				for i := 0; i < len(blocks); i++ {
					fmt.Printf("block header:\n%s\n", blocks[i])
				}
			}
		}
	}
}

func (s *SimpleNode) Send(conn net.Conn, msg Message) {
	envelop := NewNetworkEnvelope([]byte(msg.Command()), msg.Serialize(), s.testnet)
	n, err := conn.Write(envelop.Serialize())
	if err != nil {
		panic(err)
	}
	fmt.Printf("write to %d\n bytes", n)
}

func (s *SimpleNode) Read(conn net.Conn) []*NetworkEnvelope {
	receivedBuf := make([]byte, 0)
	totalLen := 0
	for {
		buf := make([]byte, 4096)
		n, err := conn.Read(buf)
		if err != nil {
			panic(err)
		}
		totalLen += n
		receivedBuf = append(receivedBuf, buf...)
		if n < 4096 {
			break
		}
	}
	/*
		the peer node may return a version and verack
		at once
	*/
	var messages []*NetworkEnvelope
	parsedLen := 0
	for {
		if parsedLen >= totalLen {
			break
		}
		msg := ParseNetwork(receivedBuf, s.testnet)
		messages = append(messages, msg)
		if parsedLen < totalLen {
			parsedLen += len(msg.Serialize())
			receivedBuf = receivedBuf[len(msg.Serialize()):]
		}
	}

	return messages
}

func (s *SimpleNode) WaitFor(conn net.Conn) {
	s.Send(conn, NewVersionMessage())

	verackReceived := false
	versionReceived := false
	for !verackReceived || !versionReceived {
		messages := s.Read(conn)
		for i := 0; i < len(messages); i++ {
			msg := messages[i]
			command := string(bytes.Trim(msg.command, "\x00"))
			fmt.Printf("command:%s\n", command)
			if command == "verack" {
				fmt.Printf("receiving verack from peer\n")
				verackReceived = true
			}
			if command == "version" {
				versionReceived = true
				fmt.Printf("receiving version message from peer\n:%s", msg)
				s.Send(conn, NewVerAckMessage())
			}
		}
	}
}

/*
get header command, following is the payload for command getheaders:
command: getheaders
payload:
7f110100

01

a35bd0ca2f4a88c4eda6d213e2378a5758dfcd6af43712000000000000000000

0000000000000000000000000000000000000000000000000000000000000000

1. the first 4 bytes is the version as we have seen before:7f110100 little endian

2. variant int to indicate the amount of block hash, in the data above, it only
has one byte which 01

3. then following 32 bytes is the block hash we want to get its header:
starting block header
a35bd0ca2f4a88c4eda6d213e2378a5758dfcd6af43712000000000000000000

4. the ending block hash:
0000000000000000000000000000000000000000000000000000000000000000
which means the full node peer will return as many block headers as possible

not more than 2000 block headers
*/
