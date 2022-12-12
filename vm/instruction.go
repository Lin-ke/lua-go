package vm

const MAXARG_Bx = 1<<17 - 1       // 262143
const MAXARG_sBx = MAXARG_Bx >> 1 // 131071

const MAXARG_J = 1<<25 - 1      //33554431
const MAXARG_sJ = MAXARG_J >> 1 //16777215

/*
===========================================================================

	We assume that instructions are unsigned 32-bit integers.
	All instructions have an opcode in the first 7 bits.
	Instructions can have the following formats:

	      3 3 2 2 2 2 2 2 2 2 2 2 1 1 1 1 1 1 1 1 1 1 0 0 0 0 0 0 0 0 0 0
	      1 0 9 8 7 6 5 4 3 2 1 0 9 8 7 6 5 4 3 2 1 0 9 8 7 6 5 4 3 2 1 0

iABC          C(8)     |      B(8)     |k|     A(8)      |   Op(7)     |
iABx                Bx(17)               |     A(8)      |   Op(7)     |
iAsBx              sBx (signed)(17)      |     A(8)      |   Op(7)     |
iAx                           Ax(25)                     |   Op(7)     |
isJ                           sJ(25)                     |   Op(7)     |

	A signed argument is represented in excess K: the represented value is
	the written unsigned value minus K, where K is half the maximum for the
	corresponding unsigned argument.

===========================================================================
*/
type Instruction uint32

func (i Instruction) Opcode() int {
	return int(i & 0x7F) //0x111 1111
}

// BC 是寄存器索引
func (i Instruction) ABC() (a, b, k, c int) {
	a = int(i >> 7 & 0xFF)
	k = int(i >> 15 & 0x1)
	b = int(i >> 16 & 0xFF)

	c = int(i >> 24 & 0xFF)
	return
}

func (i Instruction) ABx() (a, bx int) {
	a = int(i >> 7 & 0xFF)
	bx = int(i >> 15)
	return
}

func (i Instruction) AsBx() (a, sbx int) {
	a, bx := i.ABx()
	return a, bx - MAXARG_sBx
}

func (i Instruction) Ax() int {
	return int(i >> 7)
}
func (i Instruction) SJ() int {
	return int(i>>7) - MAXARG_sJ
}

func (i Instruction) OpName() string {
	return opcodes[i.Opcode()].name
}

func (i Instruction) OpMode() byte {
	return opcodes[i.Opcode()].opmode & 0x7
}

//	func opmode(callMetamethod byte, argCMode byte, argBMode byte, testFlag byte, setAFlag byte, opMode byte) byte {
//		return (((callMetamethod) << 7) | ((argCMode) << 6) | ((argBMode) << 5) | ((testFlag) << 4) | ((setAFlag) << 3) | (opMode))
//	}
func (i Instruction) BMode() byte {
	// boolean
	return (opcodes[i.Opcode()].opmode) >> 5 & 0x1
}
func (i Instruction) CMode() byte {
	// boolean
	return (opcodes[i.Opcode()].opmode) >> 6 & 0x1
}
func (i Instruction) SetA() byte {
	// boolean
	return (opcodes[i.Opcode()].opmode) >> 3 & 0x1
}

func (i Instruction) MM() byte {
	// boolean
	return (opcodes[i.Opcode()].opmode) >> 7 & 0x1
}
func (i Instruction) TMode() byte {
	// boolean
	return (opcodes[i.Opcode()].opmode) >> 4 & 0x1
}
