package transaction

import (
	"bufio"
	"io"
	"math/big"
)

type ScriptSig struct {
	commands      [][]byte
	bitcoinOpCode *BitCoinOpCode
}

const (
	// SCRIPT_DATA_LENGTH_BEGIN [0x1, 0x4b] -> [1, 75]
	SCRIPT_DATA_LENGTH_BEGIN = 1
	SCRIPT_DATA_LENGTH_END   = 75
	OP_PUSHDATA1             = 76
	OP_PUSHDATA2             = 77
)

func InitScriptSig(commands [][]byte) *ScriptSig {
	bitcoinOpCode := NewBitCoinOpCode()
	bitcoinOpCode.commands = commands
	return &ScriptSig{
		bitcoinOpCode: bitcoinOpCode,
	}
}

func NewScriptSig(reader *bufio.Reader) *ScriptSig {
	commands := [][]byte{}
	/*
		In the beginning is the total length for script field
	*/
	scriptLen := ReadVariant(reader).Int64()
	count := int64(0)
	current := make([]byte, 1)
	var currentByte byte
	for count < scriptLen {
		//buf fix
		//reader.Read(current)
		io.ReadFull(reader, current)
		//operation
		count += 1
		currentByte = current[0]
		if currentByte >= SCRIPT_DATA_LENGTH_BEGIN &&
			currentByte <= SCRIPT_DATA_LENGTH_END {
			//push the following bytes of data onto stack
			data := make([]byte, currentByte)
			//bug fix
			//reader.Read(data)
			io.ReadFull(reader, data)
			commands = append(commands, data)
			count += int64(currentByte)
		} else if currentByte == OP_PUSHDATA1 {
			/*
				read the following byte as the length of data
			*/
			length := make([]byte, 1)
			//bug fix
			//reader.Read(length)
			io.ReadFull(reader, length)

			data := make([]byte, length[0])
			//reader.Read(data)
			io.ReadFull(reader, data)
			commands = append(commands, data)
			count += int64(length[0] + 1)
		} else if currentByte == OP_PUSHDATA2 {
			/*
				read the following 2 bytes as length of data
			*/
			lenBuf := make([]byte, 2)
			//buf ifx
			//reader.Read(lenBuf)
			io.ReadFull(reader, lenBuf)
			length := LittleEndianToBigInt(lenBuf, LittleEndian2Bytes)
			data := make([]byte, length.Int64())
			//reader.Read(data)
			io.ReadFull(reader, data)
			commands = append(commands, data)
			count += int64(2 + length.Int64())
		} else {
			//it is the data processing instruction
			commands = append(commands, []byte{currentByte})
		}
	}

	if count != scriptLen {
		panic("parsing script field fail")
	}

	return InitScriptSig(commands)
}

func (s *ScriptSig) Evaluate(z []byte) bool {
	for s.bitcoinOpCode.HasCmd() {
		cmd := s.bitcoinOpCode.RemoveCmd()
		if len(cmd) == 1 {
			//this is op code, run it
			opRes := s.bitcoinOpCode.ExecuteOperation(int(cmd[0]), z)
			if opRes != true {
				return false
			}
		} else {
			s.bitcoinOpCode.AppendDataElement(cmd)
		}
	}

	/*
		After running all the operations in the scripts and the stack is empty,
		then evaluation fails, otherwise we check the top element of the stack,
		if its value is 0, then fail, if the value is not 0, then success
	*/
	if len(s.bitcoinOpCode.stack) == 0 {
		return false
	}
	if len(s.bitcoinOpCode.stack[len(s.bitcoinOpCode.stack)-1]) == 0 {
		return false
	}

	return true
}

func (s *ScriptSig) rawSerialize() []byte {
	result := []byte{}
	for _, cmd := range s.bitcoinOpCode.commands {
		if len(cmd) == 1 {
			//only one byte means it is an instruction
			result = append(result, cmd...)
		} else {
			length := len(cmd)
			if length <= SCRIPT_DATA_LENGTH_END {
				//length in [0x01, 0x4b]
				result = append(result, byte(length))
			} else if length > SCRIPT_DATA_LENGTH_END && length < 0x100 {
				//this is OP_PUSHDATA1 command,
				//push the command and then the next byte is the length of the data
				result = append(result, OP_PUSHDATA1)
				result = append(result, byte(length))
			} else if length >= 0x100 && length <= 520 {
				/*
					this is OP_PUSHDATA2 command, we push the command
					and then two bytes for the data length but in little endian format
				*/
				result = append(result, OP_PUSHDATA2)
				lenBuf := BigIntToLittleEndian(big.NewInt(int64(length)), LittleEndian2Bytes)
				result = append(result, lenBuf...)
			} else {
				panic("too long an cmd")
			}

			//append the chunk of data with given length
			result = append(result, cmd...)
		}
	}

	return result
}

func (s *ScriptSig) Serialize() []byte {
	rawResult := s.rawSerialize()
	total := len(rawResult)
	result := []byte{}
	//encode the total length of the script at the head
	result = append(result, EncodeVariant(big.NewInt(int64(total)))...)
	result = append(result, rawResult...)
	return result
}

func (s *ScriptSig) Add(script *ScriptSig) *ScriptSig {
	commands := make([][]byte, 0)
	commands = append(commands, s.bitcoinOpCode.commands...)
	commands = append(commands, script.bitcoinOpCode.commands...)
	return InitScriptSig(commands)
}
