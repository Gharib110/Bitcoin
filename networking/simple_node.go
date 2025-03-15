package networking

import (
	"bytes"
	"fmt"
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

	s.GetHeaders(conn)
}

func (s *SimpleNode) GetHeaders(conn net.Conn) {
	getHeaderMsg := NewGetHeaderMessage(GetGenesisBlockHash())
	s.Send(conn, getHeaderMsg)

	receivedGetHeader := false
	for !receivedGetHeader {
		//let the peer have a rest
		time.Sleep(2 * time.Second)
		msgs := s.Read(conn)
		for i := 0; i < len(msgs); i++ {
			msg := msgs[i]
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
	var msgs []*NetworkEnvelope
	parsedLen := 0
	for {
		if parsedLen >= totalLen {
			break
		}
		msg := ParseNetwork(receivedBuf, s.testnet)
		msgs = append(msgs, msg)
		if parsedLen < totalLen {
			parsedLen += len(msg.Serialize())
			receivedBuf = receivedBuf[len(msg.Serialize()):]
		}
	}

	return msgs
}

func (s *SimpleNode) WaitFor(conn net.Conn) {
	s.Send(conn, NewVersionMessage())

	verackReceived := false
	versionReceived := false
	for !verackReceived || !versionReceived {
		msgs := s.Read(conn)
		for i := 0; i < len(msgs); i++ {
			msg := msgs[i]
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
