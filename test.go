package main

import (
	"fmt"
)

type reader struct {
	data []byte
}

func (self *reader) readVar(limit uint64) uint64 {
	var x uint64
	var a uint8 = 0x00
	for (a & 0x80) == 0 {

		a = self.data[0]

		x = (x << 7) | (uint64)(a&0x7f)

		if x > limit {
			panic("overflow")
		}

		self.data = self.data[1:]
	}
	return x
}

//go:noinline
func opmode(callMetamethod byte, argCMode byte, argBMode byte, testFlag byte, setAFlag byte, opMode byte) byte {
	return (((callMetamethod) << 7) | ((argCMode) << 6) | ((argBMode) << 5) | ((testFlag) << 4) | ((setAFlag) << 3) | (opMode))
}

const (
	IABC  = iota // [  C:8  ][  B:8  ]k[ A:8  ][OP:7]
	IABx         // [      Bx:17     ][ A:8  ][OP:7]
	IAsBx        // [     sBx:17     ][ A:8  ][OP:7]
	IAx          // [           Ax:25        ][OP:7]
	IsJ          // [           sJ:25        ][OP:7]
)

var opcodes = []byte{

	opmode(0, 0, 0, 0, 1, IABC),  /* OP_MOVE */
	opmode(0, 0, 0, 0, 1, IAsBx), /* OP_LOADI */
	opmode(0, 0, 0, 0, 1, IAsBx), /* OP_LOADF */
	opmode(0, 0, 0, 0, 1, IABx),  /* OP_LOADK */
	opmode(0, 0, 0, 0, 1, IABx),  /* OP_LOADKX */
	opmode(0, 0, 0, 0, 1, IABC),  /* OP_LOADFALSE */
	opmode(0, 0, 0, 0, 1, IABC),  /* OP_LFALSESKIP */
	opmode(0, 0, 0, 0, 1, IABC),  /* OP_LOADTRUE */
	opmode(0, 0, 0, 0, 1, IABC),  /* OP_LOADNIL */
	opmode(0, 0, 0, 0, 1, IABC),  /* OP_GETUPVAL */
	opmode(0, 0, 0, 0, 0, IABC),  /* OP_SETUPVAL */
	opmode(0, 0, 0, 0, 1, IABC),  /* OP_GETTABUP */
	opmode(0, 0, 0, 0, 1, IABC),  /* OP_GETTABLE */
	opmode(0, 0, 0, 0, 1, IABC),  /* OP_GETI */
	opmode(0, 0, 0, 0, 1, IABC),  /* OP_GETFIELD */
	opmode(0, 0, 0, 0, 0, IABC),  /* OP_SETTABUP */
	opmode(0, 0, 0, 0, 0, IABC),  /* OP_SETTABLE */
	opmode(0, 0, 0, 0, 0, IABC),  /* OP_SETI */
	opmode(0, 0, 0, 0, 0, IABC),  /* OP_SETFIELD */
	opmode(0, 0, 0, 0, 1, IABC),  /* OP_NEWTABLE */
	opmode(0, 0, 0, 0, 1, IABC),  /* OP_SELF */
	opmode(0, 0, 0, 0, 1, IABC),  /* OP_ADDI */
	opmode(0, 0, 0, 0, 1, IABC),  /* OP_ADDK */
	opmode(0, 0, 0, 0, 1, IABC),  /* OP_SUBK */
	opmode(0, 0, 0, 0, 1, IABC),  /* OP_MULK */
	opmode(0, 0, 0, 0, 1, IABC),  /* OP_MODK */
	opmode(0, 0, 0, 0, 1, IABC),  /* OP_POWK */
	opmode(0, 0, 0, 0, 1, IABC),  /* OP_DIVK */
	opmode(0, 0, 0, 0, 1, IABC),  /* OP_IDIVK */
	opmode(0, 0, 0, 0, 1, IABC),  /* OP_BANDK */
	opmode(0, 0, 0, 0, 1, IABC),  /* OP_BORK */
	opmode(0, 0, 0, 0, 1, IABC),  /* OP_BXORK */
	opmode(0, 0, 0, 0, 1, IABC),  /* OP_SHRI */
	opmode(0, 0, 0, 0, 1, IABC),  /* OP_SHLI */
	opmode(0, 0, 0, 0, 1, IABC),  /* OP_ADD */
	opmode(0, 0, 0, 0, 1, IABC),  /* OP_SUB */
	opmode(0, 0, 0, 0, 1, IABC),  /* OP_MUL */
	opmode(0, 0, 0, 0, 1, IABC),  /* OP_MOD */
	opmode(0, 0, 0, 0, 1, IABC),  /* OP_POW */
	opmode(0, 0, 0, 0, 1, IABC),  /* OP_DIV */
	opmode(0, 0, 0, 0, 1, IABC),  /* OP_IDIV */
	opmode(0, 0, 0, 0, 1, IABC),  /* OP_BAND */
	opmode(0, 0, 0, 0, 1, IABC),  /* OP_BOR */
	opmode(0, 0, 0, 0, 1, IABC),  /* OP_BXOR */
	opmode(0, 0, 0, 0, 1, IABC),  /* OP_SHL */
	opmode(0, 0, 0, 0, 1, IABC),  /* OP_SHR */
	opmode(1, 0, 0, 0, 0, IABC),  /* OP_MMBIN */
	opmode(1, 0, 0, 0, 0, IABC),  /* OP_MMBINI*/
	opmode(1, 0, 0, 0, 0, IABC),  /* OP_MMBINK*/
	opmode(0, 0, 0, 0, 1, IABC),  /* OP_UNM */
	opmode(0, 0, 0, 0, 1, IABC),  /* OP_BNOT */
	opmode(0, 0, 0, 0, 1, IABC),  /* OP_NOT */
	opmode(0, 0, 0, 0, 1, IABC),  /* OP_LEN */
	opmode(0, 0, 0, 0, 1, IABC),  /* OP_CONCAT */
	opmode(0, 0, 0, 0, 0, IABC),  /* OP_CLOSE */
	opmode(0, 0, 0, 0, 0, IABC),  /* OP_TBC */
	opmode(0, 0, 0, 0, 0, IsJ),   /* OP_JMP */
	opmode(0, 0, 0, 1, 0, IABC),  /* OP_EQ */
	opmode(0, 0, 0, 1, 0, IABC),  /* OP_LT */
	opmode(0, 0, 0, 1, 0, IABC),  /* OP_LE */
	opmode(0, 0, 0, 1, 0, IABC),  /* OP_EQK */
	opmode(0, 0, 0, 1, 0, IABC),  /* OP_EQI */
	opmode(0, 0, 0, 1, 0, IABC),  /* OP_LTI */
	opmode(0, 0, 0, 1, 0, IABC),  /* OP_LEI */
	opmode(0, 0, 0, 1, 0, IABC),  /* OP_GTI */
	opmode(0, 0, 0, 1, 0, IABC),  /* OP_GEI */
	opmode(0, 0, 0, 1, 0, IABC),  /* OP_TEST */
	opmode(0, 0, 0, 1, 1, IABC),  /* OP_TESTSET */
	opmode(0, 1, 1, 0, 1, IABC),  /* OP_CALL */
	opmode(0, 1, 1, 0, 1, IABC),  /* OP_TAILCALL */
	opmode(0, 0, 1, 0, 0, IABC),  /* OP_RETURN */
	opmode(0, 0, 0, 0, 0, IABC),  /* OP_RETURN0 */
	opmode(0, 0, 0, 0, 0, IABC),  /* OP_RETURN1 */
	opmode(0, 0, 0, 0, 1, IABx),  /* OP_FORLOOP */
	opmode(0, 0, 0, 0, 1, IABx),  /* OP_FORPREP */
	opmode(0, 0, 0, 0, 0, IABx),  /* OP_TFORPREP */
	opmode(0, 0, 0, 0, 0, IABC),  /* OP_TFORCALL */
	opmode(0, 0, 0, 0, 1, IABx),  /* OP_TFORLOOP */
	opmode(0, 0, 1, 0, 0, IABC),  /* OP_SETLIST */
	opmode(0, 0, 0, 0, 1, IABx),  /* OP_CLOSURE */
	opmode(0, 1, 0, 0, 1, IABC),  /* OP_VARARG */
	opmode(0, 0, 1, 0, 1, IABC),  /* OP_VARARGPREP */
	opmode(0, 0, 0, 0, 0, IAx),   /* OP_EXTRAARG */
}

func main() {
	fmt.Printf("%x", opcodes[1])
}
