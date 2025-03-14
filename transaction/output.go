package transaction

import (
	"bufio"
	"math/big"
)

type TransactionOutput struct {
	//satoshi
	amount       *big.Int
	scriptPubKey *ScriptSig
}

func NewTractionOutput(reader *bufio.Reader) *TransactionOutput {
	/*
		amount is in satoshi 1/100,000,0000 of one bitcoin
	*/
	amountBuf := make([]byte, 8)
	reader.Read(amountBuf)
	amount := LittleEndianToBigInt(amountBuf, LittleEndian8Bytes)
	script := NewScriptSig(reader)
	return &TransactionOutput{
		amount:       amount,
		scriptPubKey: script,
	}
}

func (t *TransactionOutput) Serialize() []byte {
	result := make([]byte, 0)
	result = append(result, BigIntToLittleEndian(t.amount, LittleEndian8Bytes)...)
	result = append(result, t.scriptPubKey.Serialize()...)
	return result
}
