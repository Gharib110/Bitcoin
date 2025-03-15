package transaction

import (
	"bufio"
	"bytes"
	"fmt"
	ecc "github.com/Gharib110/Bitcoin/elliptic_curve"
	"math/big"
)

const (
	OP_0 = 0
)
const (
	OP_1NEGATE = iota + 79
)
const (
	OP_1 = iota + 81
	OP_2
	OP_3
	OP_4
	OP_5
	OP_6
	OP_7
	OP_8
	OP_9
	OP_10
	OP_11
	OP_12
	OP_13
	OP_14
	OP_15
	OP_16
	OP_NOP
)

const (
	OP_IF = iota + 99
	OP_NOTIf
)

const (
	OP_VERIFY = iota + 105
	OP_RETURN
	OP_TOTALSTACK
	OP_FROMALTSTACK
	OP_2DROP
	OP_2DUP
	OP_3DUP
	OP_2OVER
	OP_2ROT
	OP_2SWAP
	OP_IFDUP
	OP_DEPTH
	OP_DROP
	OP_DUP
	OP_NIP
	OP_OVER
	OP_PICK
	OP_ROLL
	OP_ROT
	OP_SWAP
	OP_TUCK
)

const (
	OP_SIZE = iota + 130
)

const (
	OP_EQUAL = iota + 135
	OP_EQUALVERIFY
)

const (
	OP_1ADD = iota + 139
	OP_1SUB
)

const (
	OP_NEGATE = iota + 143
	OP_ABS
	OP_NOT
	OP_0NOTEQUAL
	OP_ADD
	OP_SUB
	OP_MUL
)

const (
	OP_BOOLAND = iota + 154
	OP_BOOLOR
	OP_NUMEQUAL
	OP_NUMEQUALVERIFY
	OP_NUMNOTEQUAL
	OP_LESSTHAN
	OP_GREATERTHAN
	OP_LESSTHANOREQUAL
	OP_GREATERTHANOREQUAL
	OP_MIN
	OP_MAX
	OP_WITHIN
	OP_RIPEMD160
	OP_SHA1
	OP_SHA256
	OP_HASH160
	OP_HASH256
)

const (
	OP_CHECKSIG = iota + 172
	OP_HECKSIGVERIFY
	OP_CHECKMULTISIG
	OP_CHECKMULTISIGVERIFY
	OP_NOP1
	OP_CHECKLOGTIMEVERIFY
	OP_CHECKSEQUENCEVERIFY
	OP_NOP4
	OP_NOP5
	OP_NOP6
	OP_NOP7
	OP_NOP8
	OP_NOP9
	OP_NOP10
)

const (
	/*
	   this is not a bitcoin script command, it is defined by ourselves,
	   if we encounter the P2SH pattern on the script stack, that is the
	   first element is data chunk, the second element is OP_HASH160,
	   the third element is a chunk of data, the fourth element is
	   OP_EQUAL, then we will use this command to do P2SH parsing
	*/
	OP_P2SH = 254
)

type BitCoinOpCode struct {
	opCodeNames map[int]string
	stack       [][]byte
	altStack    [][]byte
	commands    [][]byte
}

func NewBitCoinOpCode() *BitCoinOpCode {
	opCodeNames := map[int]string{
		0:   "OP_0",
		76:  "OP_PUSHDATA1",
		77:  "OP_PUSHDATA2",
		78:  "OP_PUSHDATA4",
		79:  "OP_1NEGATE",
		81:  "OP_1",
		82:  "OP_2",
		83:  "OP_3",
		84:  "OP_4",
		85:  "OP_5",
		86:  "OP_6",
		87:  "OP_7",
		88:  "OP_8",
		89:  "OP_9",
		90:  "OP_10",
		91:  "OP_11",
		92:  "OP_12",
		93:  "OP_13",
		94:  "OP_14",
		95:  "OP_15",
		96:  "OP_16",
		97:  "OP_NOP",
		99:  "OP_IF",
		100: "OP_NOTIF",
		103: "OP_ELSE",
		104: "OP_ENDIF",
		105: "OP_VERIFY",
		106: "OP_RETURN",
		107: "OP_TOALTSTACK",
		108: "OP_FROMALTSTACK",
		109: "OP_2DROP",
		110: "OP_2DUP",
		111: "OP_3DUP",
		112: "OP_2OVER",
		113: "OP_2ROT",
		114: "OP_2SWAP",
		115: "OP_IFDUP",
		116: "OP_DEPTH",
		117: "OP_DROP",
		118: "OP_DUP",
		119: "OP_NIP",
		120: "OP_OVER",
		121: "OP_PICK",
		122: "OP_ROLL",
		123: "OP_ROT",
		124: "OP_SWAP",
		125: "OP_TUCK",
		130: "OP_SIZE",
		135: "OP_EQUAL",
		136: "OP_EQUALVERIFY",
		139: "OP_1ADD",
		140: "OP_1SUB",
		143: "OP_NEGATE",
		144: "OP_ABS",
		145: "OP_NOT",
		146: "OP_0NOTEQUAL",
		147: "OP_ADD",
		148: "OP_SUB",
		149: "OP_MUL",
		154: "OP_BOOLAND",
		155: "OP_BOOLOR",
		156: "OP_NUMEQUAL",
		157: "OP_NUMEQUALVERIFY",
		158: "OP_NUMNOTEQUAL",
		159: "OP_LESSTHAN",
		160: "OP_GREATERTHAN",
		161: "OP_LESSTHANOREQUAL",
		162: "OP_GREATERTHANOREQUAL",
		163: "OP_MIN",
		164: "OP_MAX",
		165: "OP_WITHIN",
		166: "OP_RIPEMD160",
		167: "OP_SHA1",
		168: "OP_SHA256",
		169: "OP_HASH160",
		170: "OP_HASH256",
		171: "OP_CODESEPARATOR",
		172: "OP_CHECKSIG",
		173: "OP_CHECKSIGVERIFY",
		174: "OP_CHECKMULTISIG",
		175: "OP_CHECKMULTISIGVERIFY",
		176: "OP_NOP1",
		177: "OP_CHECKLOCKTIMEVERIFY",
		178: "OP_CHECKSEQUENCEVERIFY",
		179: "OP_NOP4",
		180: "OP_NOP5",
		181: "OP_NOP6",
		182: "OP_NOP7",
		183: "OP_NOP8",
		184: "OP_NOP9",
		185: "OP_NOP10",
	}
	return &BitCoinOpCode{
		opCodeNames: opCodeNames,
		stack:       make([][]byte, 0),
		altStack:    make([][]byte, 0),
		commands:    make([][]byte, 0),
	}
}

// OP_0, OP_1, OP_2, push the given value on the top of parsing
func (b *BitCoinOpCode) opNum(op byte) bool {
	opNum := byte(0)
	if op >= OP_1 && op <= OP_16 {
		opNum = (op - OP_1) + 1
	}
	b.stack = append(b.stack, b.EncodeNum(int64(opNum)))
	return true
}

/*
isP2sh
check the patter for P2SH
*/
func (b *BitCoinOpCode) isP2sh() bool {
	/*
		the command stack has three elements; first is OP_HASH160, second is data chunk,
		the third is OP_EQUAL
	*/
	if len(b.commands[0]) != 1 && b.commands[0][0] != OP_HASH160 {
		return false
	}

	if len(b.commands[1]) == 1 {
		return false
	}

	if len(b.commands[2]) != 1 && b.commands[2][0] != OP_EQUAL {
		return false
	}

	return true
}

func (b *BitCoinOpCode) opDup() bool {
	if len(b.stack) < 1 {
		return false
	}

	b.stack = append(b.stack, b.stack[len(b.stack)-1])
	return true
}

func (b *BitCoinOpCode) opHash160() bool {
	if len(b.stack) < 1 {
		return false
	}

	element := b.stack[len(b.stack)-1]
	b.stack = b.stack[0 : len(b.stack)-1]
	hash160 := ecc.Hash160(element)
	b.stack = append(b.stack, hash160)
	return true
}

func (b *BitCoinOpCode) opEqual() bool {
	if len(b.stack) < 2 {
		return false
	}

	elem1 := b.stack[len(b.stack)-1]
	b.stack = b.stack[0 : len(b.stack)-1]
	elem2 := b.stack[len(b.stack)-1]
	b.stack = b.stack[0 : len(b.stack)-1]
	if bytes.Equal(elem1, elem2) {
		b.stack = append(b.stack, b.EncodeNum(1))
	} else {
		b.stack = append(b.stack, b.EncodeNum(0))
	}

	return true
}

func (b *BitCoinOpCode) opVerify() bool {
	if len(b.stack) < 1 {
		return false
	}

	elem := b.stack[len(b.stack)-1]
	b.stack = b.stack[0 : len(b.stack)-1]
	if b.DecodeNum(elem) == 0 {
		return false
	}

	return true
}

func (b *BitCoinOpCode) opEqualVerify() bool {
	resEqual := b.opEqual()
	resVerify := b.opVerify()
	return resEqual && resVerify
}

/*
popStack
OP_CHECKMULTISIG, structure of the evaluated stack:
1, number of public keys (2 pubKeys)
2, pubKey 2
3, pubKey 1
4. number of signatures (1 signature)
5. signature
6. []
*/
func (b *BitCoinOpCode) popStack() []byte {
	elem := b.stack[len(b.stack)-1]
	b.stack = b.stack[0 : len(b.stack)-1]
	return elem
}

func (b *BitCoinOpCode) opCheckMultiSig(zBin []byte) bool {
	if len(b.stack) < 1 {
		return false
	}
	//read the top element to get the number of public keys
	pubKeyCounts := int(b.DecodeNum(b.popStack()))
	if len(b.stack) < pubKeyCounts+1 {
		return false
	}
	secPubKeys := make([][]byte, 0)
	for i := 0; i < pubKeyCounts; i++ {
		secPubKeys = append(secPubKeys, b.popStack())
	}

	//get the number of signatures
	sigCounts := int(b.DecodeNum(b.popStack()))
	if len(b.stack) < sigCounts+1 {
		return false
	}

	derSignatures := make([][]byte, 0)
	for i := 0; i < sigCounts; i++ {
		signature := b.popStack()
		//remove the last byte, it is a hash type
		signature = signature[0 : len(signature)-1]
		derSignatures = append(derSignatures, signature)
	}

	points := make([]*ecc.Point, 0)
	sigs := make([]*ecc.Signature, 0)
	for i := 0; i < pubKeyCounts; i++ {
		points = append(points, ecc.ParseSEC(secPubKeys[i]))
	}
	for i := 0; i < sigCounts; i++ {
		sigs = append(sigs, ecc.ParseSigBin(derSignatures[i]))
	}

	/*
		m public keys, n signatures, m >= n, given the signature with index i,
		we need to find the paring key with index after i,
	*/
	z := new(big.Int)
	z.SetBytes(zBin)
	n := ecc.GetBitcoinValueN()
	zField := ecc.NewFieldElement(n, z)
	for _, sig := range sigs {
		if len(points) == 0 {
			return false
		}

		for len(points) > 0 {
			point := points[0]
			points = points[1:]
			if point.Verify(zField, sig) {
				break
			}
		}
	}

	b.stack = append(b.stack, b.EncodeNum(1))
	return true
}

func (b *BitCoinOpCode) opCheckSig(zBin []byte) bool {
	/*
			OP_CHECKSIG verify the validity of the message z,
			DER binary data of the signature and the uncompressed sec public key
			are top two elements of the stack

			notice!!
		    we need to remove the last byte of the der binary data because
			this byte is used for a hash type

			if the signature verification success, push 1 on the stack, otherwise
			push 0 on the stack

			if the script is using uncompressed sec format for a public key,
			then the script is called P2PK (pay for a public key)
	*/
	if len(b.stack) < 2 {
		return false
	}
	pubKey := b.stack[len(b.stack)-1]
	b.stack = b.stack[0 : len(b.stack)-1]
	derSig := b.stack[len(b.stack)-1]
	derSig = derSig[0 : len(derSig)-1]
	b.stack = b.stack[0 : len(b.stack)-1]

	point := ecc.ParseSEC(pubKey)
	sig := ecc.ParseSigBin(derSig)

	z := new(big.Int)
	z.SetBytes(zBin)
	n := ecc.GetBitcoinValueN()
	zField := ecc.NewFieldElement(n, z)
	if point.Verify(zField, sig) == true {
		b.stack = append(b.stack, b.EncodeNum(1))
	} else {
		b.stack = append(b.stack, b.EncodeNum(0))
	}

	return true
}

func (b *BitCoinOpCode) RemoveCmd() []byte {
	cmd := b.commands[0]
	b.commands = b.commands[1:]
	return cmd
}

func (b *BitCoinOpCode) HasCmd() bool {
	return len(b.commands) > 0
}

func (b *BitCoinOpCode) opP2sh() bool {
	//the first command is OP_HASH160
	b.RemoveCmd()
	//second element is data chunk
	h160 := b.RemoveCmd()
	//buf fix, remove OP_EQUAL
	b.RemoveCmd()
	/*
		the top element of stack is the content of the redeem script, cache it then
		do hash160 on it
	*/
	redeemScriptBinary := b.stack[len(b.stack)-1]
	if b.opHash160() != true {
		return false
	}
	//append the h160 on the top of stack
	b.stack = append(b.stack, h160)
	//compare the two top hash elements
	if b.opEqual() != true {
		return false
	}

	//make sure the top element is 1
	if b.opVerify() != true {
		return false
	}

	//need to parse the redeem script and execute its command
	//bug here, append the total length at the beginning
	redeemScriptBinary = append([]byte{byte(len(redeemScriptBinary))}, redeemScriptBinary...)
	scriptReader := bytes.NewReader(redeemScriptBinary)
	redeemScriptSig := NewScriptSig(bufio.NewReader(scriptReader))
	//bug here
	b.commands = append(b.commands, redeemScriptSig.bitcoinOpCode.commands...)
	return true

}

func (b *BitCoinOpCode) AppendDataElement(element []byte) {
	b.stack = append(b.stack, element)
	/*
		everytime we push a data element, we need to check the command stack
		meet the pattern of P2SH
	*/
	if b.isP2sh() {
		//insert the OP_P2SH to the head of the command stack
		b.commands = append([][]byte{[]byte{OP_P2SH}}, b.commands...)
	}
}

func (b *BitCoinOpCode) ExecuteOperation(cmd int, z []byte) bool {
	/*
		if the operation executed successfully then return true,
		otherwise return false
	*/
	switch cmd {
	case OP_CHECKSIG:
		return b.opCheckSig(z)
	case OP_DUP:
		return b.opDup()
	case OP_HASH160:
		return b.opHash160()
	case OP_EQUALVERIFY:
		return b.opEqualVerify()
	case OP_CHECKMULTISIG:
		//bug fix here
		return b.opCheckMultiSig(z)
	case OP_P2SH:
		return b.opP2sh()
	case OP_0:
		fallthrough
	case OP_1:
		fallthrough
	case OP_2:
		fallthrough
	case OP_3:
		fallthrough
	case OP_4:
		fallthrough
	case OP_5:
		fallthrough
	case OP_6:
		fallthrough
	case OP_7:
		fallthrough
	case OP_8:
		fallthrough
	case OP_9:
		fallthrough
	case OP_10:
		fallthrough
	case OP_11:
		fallthrough
	case OP_12:
		fallthrough
	case OP_13:
		fallthrough
	case OP_14:
		fallthrough
	case OP_15:
		fallthrough
	case OP_16:
		return b.opNum(byte(cmd))
	case OP_EQUAL:
		return b.opEqual()
	default:
		errStr := fmt.Sprintf("opeation %s not implemented\n", b.opCodeNames[cmd])
		panic(errStr)
	}

	return false
}

func (b *BitCoinOpCode) EncodeNum(num int64) []byte {
	if num == 0 {
		//not push 0x00 but empty byte string
		return []byte("")
	}

	result := []byte{}
	absNum := num
	negative := false
	if num < 0 {
		absNum = -num
		negative = true
	}

	for absNum > 0 {
		/*
					append the last byte of asbNum into a result,
			        a notice result will be little endian byte array of absNum
		*/
		result = append(result, byte(absNum&0xff))
		absNum >>= 8
	}

	/*
		check the most significant bit, notice the most significant byte is
		at the end of result
		0x8080 -> 32896-32896
	*/
	if (result[len(result)-1] & 0x80) != 0 {
		if negative {
			//need to insert 0x80 at the head, the most significant byte is at the end
			//of result, we should insert 0x80 at the end
			result = append(result, 0x80)
		} else {
			result = append(result, 0x00)
		}
	} else if negative {
		//set the most significant bit to 1
		result[len(result)-1] |= 0x80
	}

	return result
}

func (b *BitCoinOpCode) DecodeNum(element []byte) int64 {
	//check empty byte string
	if len(element) == 0 {
		return 0
	}

	bigEndian := reverseByteSlice(element)
	negative := false
	result := int64(0)

	//if the most significant bit is 1, it is negative value
	if (bigEndian[0] & 0x80) != 0 {
		negative = true
		//reset the most significant bit to 0,
		//0x7f is 0111 111
		result = int64(bigEndian[0] & 0x7f)
	} else {
		negative = false
		result = int64(bigEndian[0])
	}

	for i := 1; i < len(bigEndian); i++ {
		result <<= 8
		result += int64(bigEndian[i])
	}

	if negative {
		return -result
	}

	return result
}
