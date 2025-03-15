package transaction

import (
	"bufio"
	"fmt"
	"math/big"
)

type TransactionInput struct {
	previousTransactionID    []byte
	previousTransactionIndex *big.Int
	scriptSig                *ScriptSig
	sequence                 *big.Int
	fetcher                  *TransactionFetcher
}

func reverseByteSlice(bytes []byte) []byte {
	reverseBytes := []byte{}
	for i := len(bytes) - 1; i >= 0; i-- {
		reverseBytes = append(reverseBytes, bytes[i])
	}

	return reverseBytes
}

func NewTractionInput(reader *bufio.Reader) *TransactionInput {
	//first 32 bytes are hash256 of previous transaction
	transactionInput := &TransactionInput{}
	transactionInput.fetcher = NewTransactionFetch()

	previousTransaction := make([]byte, 32)
	reader.Read(previousTransaction)
	//convert it from little endian to big endian
	//reverse the byte array [0x01, 0x02, 0x03, 0x04] -> [0x04, 0x03, 0x02, 0x01]
	transactionInput.previousTransactionID = reverseByteSlice(previousTransaction)
	fmt.Printf("previous transaction id:%x\n", transactionInput.previousTransactionID)
	//4 bytes for previous transaction index
	idx := make([]byte, 4)
	reader.Read(idx)
	transactionInput.previousTransactionIndex = LittleEndianToBigInt(idx, LittleEndian4Bytes)
	fmt.Printf("previous transaction index:%x\n", transactionInput.previousTransactionIndex)

	transactionInput.scriptSig = NewScriptSig(reader)
	scriptBuf := transactionInput.scriptSig.Serialize()
	fmt.Printf("script byte:%x\n", scriptBuf)

	//last four bytes for sequence
	seqBytes := make([]byte, 4)
	reader.Read(seqBytes)
	transactionInput.sequence = LittleEndianToBigInt(seqBytes, LittleEndian4Bytes)

	return transactionInput
}

func (t *TransactionInput) getPreviousTx(testnet bool) *Transaction {
	previousTxID := fmt.Sprintf("%x", t.previousTransactionID)
	previousTX := t.fetcher.Fetch(previousTxID, testnet)
	tx := ParseTransaction(previousTX)
	return tx
}

func (t *TransactionInput) Value(testnet bool) *big.Int {
	tx := t.getPreviousTx(testnet)

	return tx.txOutputs[t.previousTransactionIndex.Int64()].amount
}

func (t *TransactionInput) Script(testnet bool) *ScriptSig {
	previousTxID := fmt.Sprintf("%x", t.previousTransactionID)
	previousTX := t.fetcher.Fetch(previousTxID, testnet)
	tx := ParseTransaction(previousTX)

	scriptPubKey := tx.txOutputs[t.previousTransactionIndex.Int64()].scriptPubKey
	return t.scriptSig.Add(scriptPubKey)
}

func (t *TransactionInput) scriptPubKey(testnet bool) *ScriptSig {
	tx := t.getPreviousTx(testnet)
	return tx.txOutputs[t.previousTransactionIndex.Int64()].scriptPubKey
}

func (t *TransactionInput) ReplaceWithScriptPubKey(testnet bool) {
	t.scriptSig = t.scriptPubKey(testnet)
	fmt.Printf("scriptpubkey: %x\n", t.scriptSig)
}

func (t *TransactionInput) Serialize() []byte {
	result := make([]byte, 0)
	result = append(result, reverseByteSlice(t.previousTransactionID)...)
	result = append(result,
		BigIntToLittleEndian(t.previousTransactionIndex, LittleEndian4Bytes)...)
	result = append(result, t.scriptSig.Serialize()...)
	result = append(result, BigIntToLittleEndian(t.sequence, LittleEndian4Bytes)...)
	return result
}
