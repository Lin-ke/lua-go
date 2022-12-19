package vm

import api "luago54/api"

/* OpMode */
/* basic instruction format */
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
const (
	IABC  = iota // [  C:8  ][  B:8  ]k[ A:8  ][OP:7]
	IABx         // [      Bx:17     ][ A:8  ][OP:7]
	IAsBx        // [     sBx:17     ][ A:8  ][OP:7]
	IAx          // [           Ax:25        ][OP:7]
	IsJ          // [           sJ:25        ][OP:7]
)

/* OpCode */
const (
	OP_MOVE       = iota /*	A B	R[A] := R[B]					*/
	OP_LOADI             /*	A sBx	R[A] := sBx					*/
	OP_LOADF             /*	A sBx	R[A] := (lua_Number)sBx				*/
	OP_LOADK             /*	A Bx	R[A] := K[Bx]					*/
	OP_LOADKX            /*	A	R[A] := K[extra arg]				*/
	OP_LOADFALSE         /*	A	R[A] := false					*/
	OP_LFALSESKIP        /*A	R[A] := false; pc++	(*)			*/
	OP_LOADTRUE          /*	A	R[A] := true					*/
	OP_LOADNIL           /*	A B	R[A], R[A+1], ..., R[A+B] := nil		*/
	OP_GETUPVAL          /*	A B	R[A] := UpValue[B]				*/
	OP_SETUPVAL          /*	A B	UpValue[B] := R[A]				*/

	OP_GETTABUP /*	A B C	R[A] := UpValue[B][K[C]:string]			*/
	OP_GETTABLE /*	A B C	R[A] := R[B][R[C]]				*/
	OP_GETI     /*	A B C	R[A] := R[B][C]					*/
	OP_GETFIELD /*	A B C	R[A] := R[B][K[C]:string]			*/

	OP_SETTABUP /*	A B C	UpValue[A][K[B]:string] := RK(C)		*/
	OP_SETTABLE /*	A B C	R[A][R[B]] := RK(C)				*/
	OP_SETI     /*	A B C	R[A][B] := RK(C)				*/
	OP_SETFIELD /*	A B C	R[A][K[B]:string] := RK(C)			*/

	OP_NEWTABLE /*	A B C k	R[A] := {}					*/

	OP_SELF /*	A B C	R[A+1] := R[B]; R[A] := R[B][RK(C):string]	*/

	OP_ADDI /*	A B sC	R[A] := R[B] + sC				*/

	OP_ADDK  /*	A B C	R[A] := R[B] + K[C]:number			*/
	OP_SUBK  /*	A B C	R[A] := R[B] - K[C]:number			*/
	OP_MULK  /*	A B C	R[A] := R[B] * K[C]:number			*/
	OP_MODK  /*	A B C	R[A] := R[B] % K[C]:number			*/
	OP_POWK  /*	A B C	R[A] := R[B] ^ K[C]:number			*/
	OP_DIVK  /*	A B C	R[A] := R[B] / K[C]:number			*/
	OP_IDIVK /*	A B C	R[A] := R[B] // K[C]:number			*/

	OP_BANDK /*	A B C	R[A] := R[B] & K[C]:integer			*/
	OP_BORK  /*	A B C	R[A] := R[B] | K[C]:integer			*/
	OP_BXORK /*	A B C	R[A] := R[B] ~ K[C]:integer			*/

	OP_SHRI /*	A B sC	R[A] := R[B] >> sC				*/
	OP_SHLI /*	A B sC	R[A] := sC << R[B]				*/

	OP_ADD  /*	A B C	R[A] := R[B] + R[C]				*/
	OP_SUB  /*	A B C	R[A] := R[B] - R[C]				*/
	OP_MUL  /*	A B C	R[A] := R[B] * R[C]				*/
	OP_MOD  /*	A B C	R[A] := R[B] % R[C]				*/
	OP_POW  /*	A B C	R[A] := R[B] ^ R[C]				*/
	OP_DIV  /*	A B C	R[A] := R[B] / R[C]				*/
	OP_IDIV /*	A B C	R[A] := R[B] // R[C]				*/

	OP_BAND /*	A B C	R[A] := R[B] & R[C]				*/
	OP_BOR  /*	A B C	R[A] := R[B] | R[C]				*/
	OP_BXOR /*	A B C	R[A] := R[B] ~ R[C]				*/
	OP_SHL  /*	A B C	R[A] := R[B] << R[C]				*/
	OP_SHR  /*	A B C	R[A] := R[B] >> R[C]				*/

	OP_MMBIN  /*	A B C	call C metamethod over R[A] and R[B]	(*)	*/
	OP_MMBINI /*	A sB C k	call C metamethod over R[A] and sB	*/
	OP_MMBINK /*	A B C k		call C metamethod over R[A] and K[B]	*/

	OP_UNM  /*	A B	R[A] := -R[B]					*/
	OP_BNOT /*	A B	R[A] := ~R[B]					*/
	OP_NOT  /*	A B	R[A] := not R[B]				*/
	OP_LEN  /*	A B	R[A] := #R[B] (length operator)			*/

	OP_CONCAT /*	A B	R[A] := R[A].. ... ..R[A + B - 1]		*/

	OP_CLOSE /*	A	close all upvalues >= R[A]			*/
	OP_TBC   /*	A	mark variable A "to be closed"			*/
	OP_JMP   /*	sJ	pc += sJ					*/
	OP_EQ    /*	A B k	if ((R[A] == R[B]) ~= k) then pc++		*/
	OP_LT    /*	A B k	if ((R[A] <  R[B]) ~= k) then pc++		*/
	OP_LE    /*	A B k	if ((R[A] <= R[B]) ~= k) then pc++		*/

	OP_EQK /*	A B k	if ((R[A] == K[B]) ~= k) then pc++		*/
	OP_EQI /*	A sB k	if ((R[A] == sB) ~= k) then pc++		*/
	OP_LTI /*	A sB k	if ((R[A] < sB) ~= k) then pc++			*/
	OP_LEI /*	A sB k	if ((R[A] <= sB) ~= k) then pc++		*/
	OP_GTI /*	A sB k	if ((R[A] > sB) ~= k) then pc++			*/
	OP_GEI /*	A sB k	if ((R[A] >= sB) ~= k) then pc++		*/

	OP_TEST    /*	A k	if (not R[A] == k) then pc++			*/
	OP_TESTSET /*	A B k	if (not R[B] == k) then pc++ else R[A] := R[B] (*) */

	OP_CALL     /*	A B C	R[A], ... ,R[A+C-2] := R[A](R[A+1], ... ,R[A+B-1]) */
	OP_TAILCALL /*	A B C k	return R[A](R[A+1], ... ,R[A+B-1])		*/

	OP_RETURN  /*	A B C k	return R[A], ... ,R[A+B-2]	(see note)	*/
	OP_RETURN0 /*		return						*/
	OP_RETURN1 /*	A	return R[A]					*/

	OP_FORLOOP /*	A Bx	update counters; if loop continues then pc-=Bx; */
	OP_FORPREP /*	A Bx	<check values and prepare counters>;
		if not to run then pc+=Bx+1;			*/

	OP_TFORPREP /*	A Bx	create upvalue for R[A + 3]; pc+=Bx		*/
	OP_TFORCALL /*	A C	R[A+4], ... ,R[A+3+C] := R[A](R[A+1], R[A+2]);	*/
	OP_TFORLOOP /*	A Bx	if R[A+2] ~= nil then { R[A]=R[A+2]; pc -= Bx }	*/

	OP_SETLIST /*	A B C k	R[A][C+i] := R[A+i], 1 <= i <= B		*/

	OP_CLOSURE /*	A Bx	R[A] := closure(KPROTO[Bx])			*/

	OP_VARARG /*	A C	R[A], R[A+1], ..., R[A+C-2] = vararg		*/

	OP_VARARGPREP /*A	(adjust vararg parameters)			*/

	OP_EXTRAARG /*	Ax	extra (larger) argument for previous opcode	*/
)

//	type opcode struct {
//		testFlag byte // operator is a test (next instruction must be a jump) T
//		setAFlag byte // instruction set register A
//		argBMode byte // B arg mode
//		argCMode byte // C arg mode
//		opMode   byte // op mode
//		name     string
//	}
//
//go:noinline
func opmode(callMetamethod byte, argCMode byte, argBMode byte, testFlag byte, setAFlag byte, opMode byte) byte {
	return (((callMetamethod) << 7) | ((argCMode) << 6) | ((argBMode) << 5) | ((testFlag) << 4) | ((setAFlag) << 3) | (opMode))
}

type opcode struct {
	opmode byte
	name   string
	action func(i Instruction, vm api.LuaVM)
}

// 指令表
var opcodes = []opcode{
	{opmode(0, 0, 0, 0, 1, IABC), "MOVE", move},
	{opmode(0, 0, 0, 0, 1, IAsBx), "LOADI", loadI},
	{opmode(0, 0, 0, 0, 1, IAsBx), "LOADF", loadF},
	{opmode(0, 0, 0, 0, 1, IABx), "LOADK", loadK},
	{opmode(0, 0, 0, 0, 1, IABx), "LOADKX", loadKx},
	{opmode(0, 0, 0, 0, 1, IABC), "LOADFALSE", loadFalse},
	{opmode(0, 0, 0, 0, 1, IABC), "LFALSESKIP", lFalseSkip},
	{opmode(0, 0, 0, 0, 1, IABC), "LOADTRUE", loadTrue},
	{opmode(0, 0, 0, 0, 1, IABC), "LOADNIL", loadNil},
	{opmode(0, 0, 0, 0, 1, IABC), "GETUPVAL", nil},
	{opmode(0, 0, 0, 0, 0, IABC), "SETUPVAL", nil},
	{opmode(0, 0, 0, 0, 1, IABC), "GETTABUP", nil},
	{opmode(0, 0, 0, 0, 1, IABC), "GETTABLE", nil},
	{opmode(0, 0, 0, 0, 1, IABC), "GETI", nil},
	{opmode(0, 0, 0, 0, 1, IABC), "GETFIELD", nil},
	{opmode(0, 0, 0, 0, 0, IABC), "SETTABUP", nil},
	{opmode(0, 0, 0, 0, 0, IABC), "SETTABLE", nil},
	{opmode(0, 0, 0, 0, 0, IABC), "SETI", nil},
	{opmode(0, 0, 0, 0, 0, IABC), "SETFIELD", nil},
	{opmode(0, 0, 0, 0, 1, IABC), "NEWTABLE", nil},
	{opmode(0, 0, 0, 0, 1, IABC), "SELF", nil},
	{opmode(0, 0, 0, 0, 1, IABC), "ADDI", addi},
	{opmode(0, 0, 0, 0, 1, IABC), "ADDK", addk},
	{opmode(0, 0, 0, 0, 1, IABC), "SUBK", subk},
	{opmode(0, 0, 0, 0, 1, IABC), "MULK", mulk},
	{opmode(0, 0, 0, 0, 1, IABC), "MODK", modk},
	{opmode(0, 0, 0, 0, 1, IABC), "POWK", powk},
	{opmode(0, 0, 0, 0, 1, IABC), "DIVK", divk},
	{opmode(0, 0, 0, 0, 1, IABC), "IDIVK", idivk},
	{opmode(0, 0, 0, 0, 1, IABC), "BANDK", bandk},
	{opmode(0, 0, 0, 0, 1, IABC), "BORK", bork},
	{opmode(0, 0, 0, 0, 1, IABC), "BXORK", bxork},
	{opmode(0, 0, 0, 0, 1, IABC), "SHRI", shri},
	{opmode(0, 0, 0, 0, 1, IABC), "SHLI", shli},
	{opmode(0, 0, 0, 0, 1, IABC), "ADD", add},
	{opmode(0, 0, 0, 0, 1, IABC), "SUB", sub},
	{opmode(0, 0, 0, 0, 1, IABC), "MUL", mul},
	{opmode(0, 0, 0, 0, 1, IABC), "MOD", mod},
	{opmode(0, 0, 0, 0, 1, IABC), "POW", pow},
	{opmode(0, 0, 0, 0, 1, IABC), "DIV", div},
	{opmode(0, 0, 0, 0, 1, IABC), "IDIV", idiv},
	{opmode(0, 0, 0, 0, 1, IABC), "BAND", band},
	{opmode(0, 0, 0, 0, 1, IABC), "BOR", bor},
	{opmode(0, 0, 0, 0, 1, IABC), "BXOR", bxor},
	{opmode(0, 0, 0, 0, 1, IABC), "SHL", shl},
	{opmode(0, 0, 0, 0, 1, IABC), "SHR", shr},
	{opmode(1, 0, 0, 0, 0, IABC), "MMBIN", nil},
	{opmode(1, 0, 0, 0, 0, IABC), "MMBINI", nil},
	{opmode(1, 0, 0, 0, 0, IABC), "MMBINK", nil},
	{opmode(0, 0, 0, 0, 1, IABC), "UNM", nil},
	{opmode(0, 0, 0, 0, 1, IABC), "BNOT", bnot},
	{opmode(0, 0, 0, 0, 1, IABC), "NOT", not},
	{opmode(0, 0, 0, 0, 1, IABC), "LEN", len},
	{opmode(0, 0, 0, 0, 1, IABC), "CONCAT", concat},
	{opmode(0, 0, 0, 0, 0, IABC), "CLOSE", nil},
	{opmode(0, 0, 0, 0, 0, IABC), "TBC", nil},
	{opmode(0, 0, 0, 0, 0, IsJ), "JMP", jmp},
	{opmode(0, 0, 0, 1, 0, IABC), "EQ", eq},
	{opmode(0, 0, 0, 1, 0, IABC), "LT", lt},
	{opmode(0, 0, 0, 1, 0, IABC), "LE", le},
	{opmode(0, 0, 0, 1, 0, IABC), "EQK", eqk},
	{opmode(0, 0, 0, 1, 0, IABC), "EQI", eqi},
	{opmode(0, 0, 0, 1, 0, IABC), "LTI", lti},
	{opmode(0, 0, 0, 1, 0, IABC), "LEI", lei},
	{opmode(0, 0, 0, 1, 0, IABC), "GTI", gti},
	{opmode(0, 0, 0, 1, 0, IABC), "GEI", gei},
	{opmode(0, 0, 0, 1, 0, IABC), "TEST", test},
	{opmode(0, 0, 0, 1, 1, IABC), "TESTSET", testSet},
	{opmode(0, 1, 1, 0, 1, IABC), "CALL", nil},
	{opmode(0, 1, 1, 0, 1, IABC), "TAILCALL", nil},
	{opmode(0, 0, 1, 0, 0, IABC), "RETURN", nil},
	{opmode(0, 0, 0, 0, 0, IABC), "RETURN0", nil},
	{opmode(0, 0, 0, 0, 0, IABC), "RETURN1", nil},
	{opmode(0, 0, 0, 0, 1, IABx), "FORLOOP", forLoop},
	{opmode(0, 0, 0, 0, 1, IABx), "FORPREP", forPrep},
	{opmode(0, 0, 0, 0, 0, IABx), "TFORPREP", nil},
	{opmode(0, 0, 0, 0, 0, IABC), "TFORCALL", nil},
	{opmode(0, 0, 0, 0, 1, IABx), "TFORLOOP", nil},
	{opmode(0, 0, 1, 0, 0, IABC), "SETLIST", nil},
	{opmode(0, 0, 0, 0, 1, IABx), "CLOSURE", nil},
	{opmode(0, 1, 0, 0, 1, IABC), "VARARG", nil},
	{opmode(0, 0, 1, 0, 1, IABC), "VARARGPREP", nil},
	{opmode(0, 0, 0, 0, 0, IAx), "EXTRAARG", nil},
}
