package transaction

import (
	"bufio"
	"fmt"
	"math/big"
)

type TransactionOutput struct {
	//satoshi
	amount       *big.Int
	scriptPubKey *ScriptSig
}

func InitTransactionOutput(amount *big.Int, script *ScriptSig) *TransactionOutput {
	return &TransactionOutput{
		amount:       amount,
		scriptPubKey: script,
	}
}

func (t *TransactionOutput) String() string {
	return fmt.Sprintf("amount:%v\n scriptPubKey: %x\n", t.amount, t.scriptPubKey.Serialize())
}

func NewTractionOutput(reader *bufio.Reader) *TransactionOutput {
	/*
		the amount is in satoshi 1/100,000,0000 of one bitcoin
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
