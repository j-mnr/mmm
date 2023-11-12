package code

type Instructions []byte

type Instruction struct {
	op Opcode
	operands [3]byte
}

// Opcode is the first part of an Instruction that's used to tell the VM what to
// do with the bytes after the Opcode.
type Opcode byte

const (
	_ Opcode = iota
	// OpConstant takes in 1 uint16 operand.
	OpConstant
)

func Make(i Instruction) Instructions {
	
}
