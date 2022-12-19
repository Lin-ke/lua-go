package vm

import . "luago54/api"

// "a+1" is reasonable here, for it's vm's (private) rule to
// keep indexes begin at 0.
/* arith */

func add(i Instruction, vm LuaVM)  { _binaryArith(i, vm, LUA_OPADD) }  // +
func sub(i Instruction, vm LuaVM)  { _binaryArith(i, vm, LUA_OPSUB) }  // -
func mul(i Instruction, vm LuaVM)  { _binaryArith(i, vm, LUA_OPMUL) }  // *
func mod(i Instruction, vm LuaVM)  { _binaryArith(i, vm, LUA_OPMOD) }  // %
func pow(i Instruction, vm LuaVM)  { _binaryArith(i, vm, LUA_OPPOW) }  // ^
func div(i Instruction, vm LuaVM)  { _binaryArith(i, vm, LUA_OPDIV) }  // /
func idiv(i Instruction, vm LuaVM) { _binaryArith(i, vm, LUA_OPIDIV) } // //
func band(i Instruction, vm LuaVM) { _binaryArith(i, vm, LUA_OPBAND) } // &
func bor(i Instruction, vm LuaVM)  { _binaryArith(i, vm, LUA_OPBOR) }  // |
func bxor(i Instruction, vm LuaVM) { _binaryArith(i, vm, LUA_OPBXOR) } // ~
func shl(i Instruction, vm LuaVM)  { _binaryArith(i, vm, LUA_OPSHL) }  // <<
func shr(i Instruction, vm LuaVM)  { _binaryArith(i, vm, LUA_OPSHR) }  // >>
func unm(i Instruction, vm LuaVM)  { _unaryArith(i, vm, LUA_OPUNM) }   // -
func bnot(i Instruction, vm LuaVM) { _unaryArith(i, vm, LUA_OPBNOT) }  // ~

func addi(i Instruction, vm LuaVM) { _binaryscArith(i, vm, LUA_OPADD) }
func shri(i Instruction, vm LuaVM) { _binaryscArith(i, vm, LUA_OPSHR) }
func shli(i Instruction, vm LuaVM) { _binaryscArith(i, vm, LUA_OPSHL) }

func addk(i Instruction, vm LuaVM)  { _binaryKArith(i, vm, LUA_OPADD) }
func subk(i Instruction, vm LuaVM)  { _binaryKArith(i, vm, LUA_OPSUB) }
func mulk(i Instruction, vm LuaVM)  { _binaryKArith(i, vm, LUA_OPMUL) }
func modk(i Instruction, vm LuaVM)  { _binaryKArith(i, vm, LUA_OPMOD) }
func powk(i Instruction, vm LuaVM)  { _binaryKArith(i, vm, LUA_OPPOW) }
func divk(i Instruction, vm LuaVM)  { _binaryKArith(i, vm, LUA_OPDIV) }
func idivk(i Instruction, vm LuaVM) { _binaryKArith(i, vm, LUA_OPIDIV) }
func bandk(i Instruction, vm LuaVM) { _binaryKArith(i, vm, LUA_OPBAND) }
func bork(i Instruction, vm LuaVM)  { _binaryKArith(i, vm, LUA_OPBOR) }
func bxork(i Instruction, vm LuaVM) { _binaryKArith(i, vm, LUA_OPBXOR) }

// R(A) := R[B] op R[C]
func _binaryArith(i Instruction, vm LuaVM, op ArithOp) {
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
func _binaryscArith(i Instruction, vm LuaVM, op ArithOp) {
	a, _, b, c := i.ABC()
	a += 1
	b += 1
	vm.PushValue(b)
	vm.Push(c)
	vm.Arith(op)
	vm.Replace(a)
}

// R(A) = R(B) op k[C]
func _binaryKArith(i Instruction, vm LuaVM, op ArithOp) {
	a, _, b, c := i.ABC()
	a += 1
	b += 1
	vm.PushValue(b)
	vm.GetConst(c)
	vm.Arith(op)
	vm.Replace(a)
}

// R(A) := op R(B)
func _unaryArith(i Instruction, vm LuaVM, op ArithOp) {
	a, _, b, _ := i.ABC()
	a += 1
	b += 1

	vm.PushValue(b)
	vm.Arith(op)
	vm.Replace(a)
}

/* compare */

func eq(i Instruction, vm LuaVM) { _compare(i, vm, LUA_OPEQ) } // ==
func lt(i Instruction, vm LuaVM) { _compare(i, vm, LUA_OPLT) } // <
func le(i Instruction, vm LuaVM) { _compare(i, vm, LUA_OPLE) } // <=

func eqi(i Instruction, vm LuaVM) { _comparesb(i, vm, LUA_OPEQ, false) } // ==
func lti(i Instruction, vm LuaVM) { _comparesb(i, vm, LUA_OPLT, false) } // <
func lei(i Instruction, vm LuaVM) { _comparesb(i, vm, LUA_OPLE, false) } // <=
func gti(i Instruction, vm LuaVM) { _comparesb(i, vm, LUA_OPLT, true) }  // >
func gei(i Instruction, vm LuaVM) { _comparesb(i, vm, LUA_OPLE, true) }  // >=

func eqk(i Instruction, vm LuaVM) { _compareK(i, vm, LUA_OPEQ) } // ==

// if ((RK(B) op RK(C)) ~= A) then pc++
func _compare(i Instruction, vm LuaVM, op CompareOp) {
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
func _compareK(i Instruction, vm LuaVM, op CompareOp) {
	a, k, b, _ := i.ABC()
	a += 1
	vm.GetConst(b)
	vm.PushValue(a)
	if vm.Compare(-2, -1, op) != (k != 0) {
		vm.AddPC(1)
	}
	vm.Pop(2)
}
func _comparesb(i Instruction, vm LuaVM, op CompareOp, inv bool) {
	a, k, sb, _ := i.ABC()
	a += 1
	vm.PushValue(a)
	vm.Push(sb)
	if inv {
		if vm.Compare(-1, -2, op) != (k != 0) {
			vm.AddPC(1)
		}
	} else {
		if vm.Compare(-1, -2, op) != (k != 0) {
			vm.AddPC(1)
		}
	}

	vm.Pop(2)
}

/* logical */

// R(A) := not R(B)
func not(i Instruction, vm LuaVM) {
	a, _, b, _ := i.ABC()
	a += 1
	b += 1

	vm.PushBoolean(!vm.ToBoolean(b))
	vm.Replace(a)
}

// if (not R(A) == k) then pc++
func test(i Instruction, vm LuaVM) {
	a, k, _, _ := i.ABC()
	a += 1

	if vm.ToBoolean(a) != (k != 0) {
		vm.AddPC(1)
	}
}

// if (not R[B] == k) then pc++ else R[A] := R[B]
func testSet(i Instruction, vm LuaVM) {
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
func len(i Instruction, vm LuaVM) {
	a, _, b, _ := i.ABC()
	a += 1
	b += 1

	vm.Len(b)
	vm.Replace(a)
}

// R(A) := R(B).. ... ..R(C)
func concat(i Instruction, vm LuaVM) {
	a, _, b, c := i.ABC()
	a += 1
	b += 1
	c += 1

	n := c - b + 1
	vm.CheckStack(n)
	for i := b; i <= c; i++ {
		vm.PushValue(i)
	}
	vm.Concat(n)
	vm.Replace(a)
}
