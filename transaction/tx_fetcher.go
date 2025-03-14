package transaction

import (
	"encoding/hex"
	"fmt"
	"io"
	"net/http"
)

type TransactionFetcher struct{}

func NewTransactionFetch() *TransactionFetcher {
	return &TransactionFetcher{}
}

func (t *TransactionFetcher) getURL(testnet bool) string {
	if testnet {
		return "https://blockstream.info/testnet/api/tx"
	}

	return "https://blockstream.info/api/tx"
}

func (t *TransactionFetcher) Fetch(txID string, testnet bool) []byte {
	url := fmt.Sprintf("%s/%s/hex", t.getURL(testnet), txID)
	fmt.Printf("fetching url: %s\n", url)
	resp, err := http.Get(url)
	if err != nil {
		panic(fmt.Sprintf("fetch transaction err:%v\n", err))
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		panic(fmt.Sprintf("read body err:%v\n", err))
	}
	fmt.Printf("response: %s\n", string(body))
	buf, err := hex.DecodeString(string(body))

	if err != nil {
		panic(err)
	}
	return buf
}
