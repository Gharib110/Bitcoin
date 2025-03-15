package transaction

import (
	"bufio"
	"bytes"
	"fmt"
	ecc "github.com/Gharib110/Bitcoin/elliptic_curve"
	"io"
	"math/big"
)

const (
	SIGHASH_ALL = 1
)

type Transaction struct {
	version   *big.Int
	txInputs  []*TransactionInput
	txOutputs []*TransactionOutput
	lockTime  *big.Int
	testnet   bool
}

func InitTransaction(version *big.Int, txInputs []*TransactionInput, txOutputs []*TransactionOutput,
	lockTime *big.Int, testnet bool) *Transaction {
	return &Transaction{
		version:   version,
		txInputs:  txInputs,
		txOutputs: txOutputs,
		lockTime:  lockTime,
		testnet:   testnet,
	}
}

func (t *Transaction) String() string {
	txIns := ""
	for i := 0; i < len(t.txInputs); i++ {
		txIns += t.txInputs[i].String()
		txIns += "\n"
	}

	txOuts := ""
	for i := 0; i < len(t.txOutputs); i++ {
		txOuts += t.txOutputs[i].String()
		txOuts += "\n"
	}

	return fmt.Sprintf("tx version: %x\n transaction inputs: %s\n transaction outputs:%s\n locktime:%x\n",
		t.version, txIns, txOuts, t.lockTime)
}

func getInputCount(bufReader *bufio.Reader) *big.Int {
	/*
		if the first byte of input is 0, then witness transaction,
		we need to skip the first two bytes(0x00, 0x01)
	*/
	firstByte, err := bufReader.Peek(1)
	if err != nil {
		panic(err)
	}
	if firstByte[0] == 0x00 {
		//skip the first two bytes
		skipBuf := make([]byte, 2)
		//_, err = bufReader.Read(skipBuf)
		_, err = io.ReadFull(bufReader, skipBuf)
		if err != nil {
			panic(err)
		}
	}

	count := ReadVariant(bufReader)
	fmt.Printf("input count is: %x\n", count)
	return count
}

func (t *Transaction) SerializeWithSign(inputIdx int) []byte {
	/*
		construct a signature message for the given input indicate by inputIdx,
		we need to change the given scriptSig with the scriptPubKey from the
		output of the previous transaction, and then do hash256 on the binary transaction
		data
	*/
	signBinary := make([]byte, 0)
	signBinary = append(signBinary, BigIntToLittleEndian(t.version, LittleEndian4Bytes)...)

	inputCount := big.NewInt(int64(len(t.txInputs)))
	signBinary = append(signBinary, EncodeVariant(inputCount)...)
	/*
		serialize inputs, need to replace the scriptSig of the given input
		to scriptPubKey of previous transaction
	*/
	for i := 0; i < len(t.txInputs); i++ {
		if i == inputIdx {
			t.txInputs[i].ReplaceWithScriptPubKey(t.testnet)
			signBinary = append(signBinary, t.txInputs[i].Serialize()...)
		} else {
			signBinary = append(signBinary, t.txInputs[i].Serialize()...)
		}
	}

	outputCount := big.NewInt(int64(len(t.txOutputs)))
	signBinary = append(signBinary, EncodeVariant(outputCount)...)
	for i := 0; i < len(t.txOutputs); i++ {
		signBinary = append(signBinary, t.txOutputs[i].Serialize()...)
	}

	signBinary = append(signBinary, BigIntToLittleEndian(t.lockTime, LittleEndian4Bytes)...)
	signBinary = append(signBinary,
		BigIntToLittleEndian(big.NewInt(int64(SIGHASH_ALL)), LittleEndian4Bytes)...)

	return signBinary
}

func (t *Transaction) SignHash(inputIdx int) []byte {
	signBinary := t.SerializeWithSign(inputIdx)
	h256 := ecc.Hash256(string(signBinary))
	return h256
}

func (t *Transaction) VerifyInput(inputIndex int) bool {
	verifyScript := t.GetScript(inputIndex, t.testnet)
	z := t.SignHash(inputIndex)
	return verifyScript.Evaluate(z)
}

func (t *Transaction) Verify() bool {
	/*
		1. verify fee
		2. verify each transaction input
	*/
	if t.Fee().Cmp(big.NewInt(int64(0))) < 0 {
		return false
	}

	for i := 0; i < len(t.txInputs); i++ {
		if t.VerifyInput(i) != true {
			return false
		}
	}

	return true
}

func ParseTransaction(binary []byte) *Transaction {
	transaction := &Transaction{}
	reader := bytes.NewReader(binary)
	bufReader := bufio.NewReader(reader)

	verBuf := make([]byte, 4)
	//bufReader.Read(verBuf)
	io.ReadFull(bufReader, verBuf)

	version := LittleEndianToBigInt(verBuf, LittleEndian4Bytes)
	fmt.Printf("transaction version:%x\n", version)
	transaction.version = version

	inputs := getInputCount(bufReader)
	transactionInputs := []*TransactionInput{}
	for i := 0; i < int(inputs.Int64()); i++ {
		input := NewTractionInput(bufReader)
		transactionInputs = append(transactionInputs, input)
	}
	transaction.txInputs = transactionInputs

	//read output counts
	outputs := ReadVariant(bufReader)
	transactionOutputs := []*TransactionOutput{}
	for i := 0; i < int(outputs.Int64()); i++ {
		output := NewTractionOutput(bufReader)
		transactionOutputs = append(transactionOutputs, output)
	}
	transaction.txOutputs = transactionOutputs

	//get last four bytes for lock time
	lockTimeBytes := make([]byte, 4)
	//bufReader.Read(lockTimeBytes)
	io.ReadFull(bufReader, lockTimeBytes)
	transaction.lockTime = LittleEndianToBigInt(lockTimeBytes, LittleEndian4Bytes)

	return transaction

}

func (t *Transaction) GetScript(idx int, testnet bool) *ScriptSig {
	if idx < 0 || idx >= len(t.txInputs) {
		panic("invalid idx for transaction input")
	}

	txInput := t.txInputs[idx]
	return txInput.Script(testnet)
}

func (t *Transaction) Fee() *big.Int {
	//amount of input - amount of output > 0
	inputSum := big.NewInt(int64(0))
	outputSum := big.NewInt(int64(0))

	for i := 0; i < len(t.txInputs); i++ {
		addOp := new(big.Int)
		value := t.txInputs[i].Value(t.testnet)
		inputSum = addOp.Add(inputSum, value)
	}

	for i := 0; i < len(t.txOutputs); i++ {
		addOp := new(big.Int)
		outputSum = addOp.Add(outputSum, t.txOutputs[i].amount)
	}

	opSub := new(big.Int)
	return opSub.Sub(inputSum, outputSum)
}

func (t *Transaction) IsCoinBase() bool {
	/*
			1, input count always 1
		    2, previous transaction id or hash is 32 bytes with all 0s:
		    3, previous output index is always 0xffffffff
	*/
	if len(t.txInputs) != 1 {
		return false
	}

	for i := 0; i < 32; i++ {
		if t.txInputs[0].previousTransactionID[i] != 0x00 {
			return false
		}
	}

	coinBaseIdx := big.NewInt(int64(0xffffffff))
	if t.txInputs[0].previousTransactionIndex.Cmp(coinBaseIdx) != 0 {
		return false
	}

	return true
}
