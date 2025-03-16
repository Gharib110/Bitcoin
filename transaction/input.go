package transaction

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"math/big"
)

type TransactionInput struct {
	previousTransactionID    []byte
	previousTransactionIndex *big.Int
	scriptSig                *ScriptSig
	sequence                 *big.Int
	fetcher                  *TransactionFetcher
	//add new here
	witness [][]byte
}

func InitTransactionInput(previousTx []byte, previousIndex *big.Int) *TransactionInput {
	return &TransactionInput{
		previousTransactionID:    previousTx,
		previousTransactionIndex: previousIndex,
		scriptSig:                nil,
		sequence:                 big.NewInt(int64(0xffffffff)),
	}
}

func (t *TransactionInput) String() string {
	return fmt.Sprintf("previous transaction: %x\n previous tx index:%x\n",
		t.previousTransactionID, t.previousTransactionIndex)
}

func (t *TransactionInput) SetScriptSig(script *ScriptSig) {
	t.scriptSig = script
}

func (t *TransactionInput) SetString(sig *ScriptSig) {
	t.scriptSig = sig
}

func reverseByteSlice(bytes []byte) []byte {
	reverseBytes := []byte{}
	for i := len(bytes) - 1; i >= 0; i-- {
		reverseBytes = append(reverseBytes, bytes[i])
	}

	return reverseBytes
}

func NewTractionInput(reader *bufio.Reader) *TransactionInput {
	//the first 32 bytes are hash256 of previous transaction
	transactionInput := &TransactionInput{}
	transactionInput.fetcher = NewTransactionFetch()

	previousTransaction := make([]byte, 32)
	//bug fix
	//reader.Read(previousTransaction)
	io.ReadFull(reader, previousTransaction)
	//convert it from little endian to big endian
	//reverse the byte array [0x01, 0x02, 0x03, 0x04] -> [0x04, 0x03, 0x02, 0x01]
	transactionInput.previousTransactionID = reverseByteSlice(previousTransaction)

	//4 bytes for previous transaction index
	idx := make([]byte, 4)
	//bug fix, io.ReadFull
	//reader.Read(idx)
	io.ReadFull(reader, idx)
	transactionInput.previousTransactionIndex = LittleEndianToBigInt(idx, LittleEndian4Bytes)

	transactionInput.scriptSig = NewScriptSig(reader)

	//last four bytes for the sequence
	seqBytes := make([]byte, 4)
	//bug fix, io.ReadFull
	//reader.Read(seqBytes)
	io.ReadFull(reader, seqBytes)
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

func (t *TransactionInput) isP2sh(script *ScriptSig) bool {
	isP2sh := true
	if len(script.bitcoinOpCode.commands[0]) != 1 || script.bitcoinOpCode.commands[0][0] != OP_HASH160 {
		isP2sh = false
	}

	if len(script.bitcoinOpCode.commands[1]) == 1 {
		isP2sh = false
	}

	if len(script.bitcoinOpCode.commands[2]) != 1 || script.bitcoinOpCode.commands[2][0] != OP_EQUAL {
		isP2sh = false
	}

	return isP2sh
}

func (t *TransactionInput) ReplaceWithScriptPubKey(testnet bool) {
	/*
		if it is a P2SH transaction, we use the redeem script to replace the
		scriptSig of the current input
	*/
	script := t.scriptPubKey(testnet)
	isP2sh := t.isP2sh(script)
	if isP2sh != true {
		t.scriptSig = script
	} else {
		/*
			for P2SH, we need to use the redeem script to replace the input script;
			the redeem script is at the bottom of scriptSig command stack
		*/
		redeemScriptBinary := t.scriptSig.bitcoinOpCode.commands[len(t.scriptSig.bitcoinOpCode.commands)-1]
		//bug here, append total length to the head
		redeemScriptBinary = append([]byte{byte(len(redeemScriptBinary))}, redeemScriptBinary...)
		redeemScriptReader := bytes.NewReader(redeemScriptBinary)
		redeemScript := NewScriptSig(bufio.NewReader(redeemScriptReader))
		t.scriptSig = redeemScript
	}

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
