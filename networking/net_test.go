package networking

import (
	"encoding/hex"
	"fmt"
	"testing"
)

func TestOne(t *testing.T) {
	networkRawData, err := hex.DecodeString("f9beb4d976657273696f6e0000000000650000005f1a69d2721101000100000000000000bc8f5e5400000000010000000000000000000000000000000000ffffc61b6409208d010000000000000000000000000000000000ffffcb0071c0208d128035cbc97953f80f2f5361746f7368693a302e392e332fcf05050001")
	if err != nil {
		panic(err)
	}

	network := ParseNetwork(networkRawData, false)
	fmt.Printf("%s\n", network)
	fmt.Printf("%x\n", network.Serialize())
	version := NewVersionMessage()
	fmt.Printf("version: %x\n", version.Serialize())
}

func TestTwo(t *testing.T) {
	getHeaderMsg := NewGetHeaderMessage(GetGenesisBlockHash())
	fmt.Printf("raw data for get headers msg:%x\n", getHeaderMsg.Serialize())
}
