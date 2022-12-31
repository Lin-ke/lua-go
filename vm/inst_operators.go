package vm

import (
	"luago54/api"
)

// "a+1" is reasonable here, for it's vm's (private) rule to
// keep indexes begin at 0.
/* arith */

func add(i Instruction, vm api.LuaVM)  { _binaryArith(i, vm, api.LUA_OPADD) }  // +
func sub(i Instruction, vm api.LuaVM)  { _binaryArith(i, vm, api.LUA_OPSUB) }  // -
func mul(i Instruction, vm api.LuaVM)  { _binaryArith(i, vm, api.LUA_OPMUL) }  // *
func mod(i Instruction, vm api.LuaVM)  { _binaryArith(i, vm, api.LUA_OPMOD) }  // %
func pow(i Instruction, vm api.LuaVM)  { _binaryArith(i, vm, api.LUA_OPPOW) }  // ^
func div(i Instruction, vm api.LuaVM)  { _binaryArith(i, vm, api.LUA_OPDIV) }  // /
func idiv(i Instruction, vm api.LuaVM) { _binaryArith(i, vm, api.LUA_OPIDIV) } // //
func band(i Instruction, vm api.LuaVM) { _binaryArith(i, vm, api.LUA_OPBAND) } // &
func bor(i Instruction, vm api.LuaVM)  { _binaryArith(i, vm, api.LUA_OPBOR) }  // |
func bxor(i Instruction, vm api.LuaVM) { _binaryArith(i, vm, api.LUA_OPBXOR) } // ~
func shl(i Instruction, vm api.LuaVM)  { _binaryArith(i, vm, api.LUA_OPSHL) }  // <<
func shr(i Instruction, vm api.LuaVM)  { _binaryArith(i, vm, api.LUA_OPSHR) }  // >>
func unm(i Instruction, vm api.LuaVM)  { _unaryArith(i, vm, api.LUA_OPUNM) }   // -
func bnot(i Instruction, vm api.LuaVM) { _unaryArith(i, vm, api.LUA_OPBNOT) }  // ~

func addi(i Instruction, vm api.LuaVM) { _binaryscArith(i, vm, api.LUA_OPADD) }
func shri(i Instruction, vm api.LuaVM) { _binaryscArith(i, vm, api.LUA_OPSHR) }
func shli(i Instruction, vm api.LuaVM) { _binaryscArith(i, vm, api.LUA_OPSHL) }

func addk(i Instruction, vm api.LuaVM)  { _binaryKArith(i, vm, api.LUA_OPADD) }
func subk(i Instruction, vm api.LuaVM)  { _binaryKArith(i, vm, api.LUA_OPSUB) }
func mulk(i Instruction, vm api.LuaVM)  { _binaryKArith(i, vm, api.LUA_OPMUL) }
func modk(i Instruction, vm api.LuaVM)  { _binaryKArith(i, vm, api.LUA_OPMOD) }
func powk(i Instruction, vm api.LuaVM)  { _binaryKArith(i, vm, api.LUA_OPPOW) }
func divk(i Instruction, vm api.LuaVM)  { _binaryKArith(i, vm, api.LUA_OPDIV) }
func idivk(i Instruction, vm api.LuaVM) { _binaryKArith(i, vm, api.LUA_OPIDIV) }
func bandk(i Instruction, vm api.LuaVM) { _binaryKArith(i, vm, api.LUA_OPBAND) }
func bork(i Instruction, vm api.LuaVM)  { _binaryKArith(i, vm, api.LUA_OPBOR) }
func bxork(i Instruction, vm api.LuaVM) { _binaryKArith(i, vm, api.LUA_OPBXOR) }

func mmbin(i Instruction, vm api.LuaVM)  { _binaryKArith(i, vm, api.LUA_OPBXOR) }
func mmbini(i Instruction, vm api.LuaVM) { _binaryKArith(i, vm, api.LUA_OPBXOR) }
func mmbink(i Instruction, vm api.LuaVM) { _binaryKArith(i, vm, api.LUA_OPBXOR) }

// R(A) := R[B] op R[C]
func _binaryArith(i Instruction, vm api.LuaVM, op api.ArithOp) {
	a, _, b, c := i.ABC()
	a += 1
	b += 1
	c += 1
	vm.PushValue(b)
	vm.PushValue(c)
	vm.Arith(op)
	vm.Replace(a)
}

// R(A) = R(B) op sC
func _binaryscArith(i Instruction, vm api.LuaVM, op api.ArithOp) {
	a, _, b, c := i.ABC()
	a += 1
	b += 1
	vm.PushValue(b)
	vm.PushInteger((int64(c)))
	vm.Arith(op)
	vm.Replace(a)
}

// R(A) = R(B) op k[C]
func _binaryKArith(i Instruction, vm api.LuaVM, op api.ArithOp) {
	a, _, b, c := i.ABC()
	a += 1
	b += 1
	vm.PushValue(b)
	vm.GetConst(c)
	vm.Arith(op)
	vm.Replace(a)
}

// R(A) := op R(B)
func _unaryArith(i Instruction, vm api.LuaVM, op api.ArithOp) {
	a, _, b, _ := i.ABC()
	a += 1
	b += 1

	vm.PushValue(b)
	vm.Arith(op)
	vm.Replace(a)
}

/* compare */

func eq(i Instruction, vm api.LuaVM) { _compare(i, vm, api.LUA_OPEQ) } // ==
func lt(i Instruction, vm api.LuaVM) { _compare(i, vm, api.LUA_OPLT) } // <
func le(i Instruction, vm api.LuaVM) { _compare(i, vm, api.LUA_OPLE) } // <=

func eqi(i Instruction, vm api.LuaVM) { _comparesb(i, vm, api.LUA_OPEQ, false) } // ==
func lti(i Instruction, vm api.LuaVM) { _comparesb(i, vm, api.LUA_OPLT, false) } // <
func lei(i Instruction, vm api.LuaVM) { _comparesb(i, vm, api.LUA_OPLE, false) } // <=
func gti(i Instruction, vm api.LuaVM) { _comparesb(i, vm, api.LUA_OPLT, true) }  // >
func gei(i Instruction, vm api.LuaVM) { _comparesb(i, vm, api.LUA_OPLE, true) }  // >=

func eqk(i Instruction, vm api.LuaVM) { _compareK(i, vm, api.LUA_OPEQ) } // ==

// if ((RK(B) op RK(C)) ~= A) then pc++
func _compare(i Instruction, vm api.LuaVM, op api.CompareOp) {
	a, k, b, _ := i.ABC()
	a += 1
	b += 1
	vm.PushValue(a)
	vm.PushValue(b)
	if vm.Compare(-2, -1, op) != (k != 0) {
		vm.AddPC(1)
	}
	vm.Pop(2)
}
func _compareK(i Instruction, vm api.LuaVM, op api.CompareOp) {
	a, k, b, _ := i.ABC()
	a += 1
	vm.GetConst(b)
	vm.PushValue(a)
	if vm.Compare(-2, -1, op) != (k != 0) {
		vm.AddPC(1)
	}
	vm.Pop(2)
}
func _comparesb(i Instruction, vm api.LuaVM, op api.CompareOp, inv bool) {
	a, k, sb, _ := i.ABC()
	a += 1
	vm.PushValue(a)
	vm.PushInteger(int64(sb))
	if inv {
		if vm.Compare(-1, -2, op) != (k != 0) {
			vm.AddPC(1)
		}
	} else {
		if vm.Compare(-2, -1, op) != (k != 0) {
			vm.AddPC(1)
		}
	}

	vm.Pop(2)
}

/* logical */

// R(A) := not R(B)
func not(i Instruction, vm api.LuaVM) {
	a, _, b, _ := i.ABC()
	a += 1
	b += 1

	vm.PushBoolean(!vm.ToBoolean(b))
	vm.Replace(a)
}

// if (not R(A) == k) then pc++
func test(i Instruction, vm api.LuaVM) {
	a, k, _, _ := i.ABC()
	a += 1

	if vm.ToBoolean(a) != (k != 0) {
		vm.AddPC(1)
	}
}

// if (not R[B] == k) then pc++ else R[A] := R[B]
func testSet(i Instruction, vm api.LuaVM) {
	a, k, b, _ := i.ABC()
	a += 1
	b += 1

	if vm.ToBoolean(b) == (k != 0) {
		vm.Copy(b, a)
	} else {
		vm.AddPC(1)
	}
}

/* len & concat */

// R(A) := length of R(B)
func len(i Instruction, vm api.LuaVM) {
	a, _, b, _ := i.ABC()
	a += 1
	b += 1

	vm.Len(b)
	vm.Replace(a)
}

// R[A] := R[A].. ... ..R[A + B - 1]
func concat(i Instruction, vm api.LuaVM) {
	a, _, b, _ := i.ABC()

	a += 1
	vm.CheckStack(b)
	for i := a; i <= a+b-1; i++ {
		vm.PushValue(i)
	}
	vm.Concat(b)
	vm.Replace(a)
}
