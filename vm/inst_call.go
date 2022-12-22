package vm

import . "luago54/api"

// R[A+1] := R[B]; R[A] := R[B][RK(C):string]
func self(i Instruction, vm LuaVM) {
	a, k, b, c := i.ABC()
	a += 1
	b += 1

	vm.Copy(b, a+1)
	vm.GetRK(c, k)
	vm.GetTable(b)
	vm.Replace(a)
}

// R(A) := closure(KPROTO[Bx])
func closure(i Instruction, vm LuaVM) {
	a, bx := i.ABx()
	a += 1

	vm.LoadProto(bx)
	vm.Replace(a)
}

// R[A], R[A+1], ..., R[A+C-2] = vararg
func vararg(i Instruction, vm LuaVM) {
	a, _, _, c := i.ABC()
	a += 1

	if c != 1 { // b==0 or b>1
		vm.LoadVararg(c - 1)
		_popResults(a, c, vm)
	}
}

// return R(A)(R(A+1), ... ,R(A+B-1))
func tailCall(i Instruction, vm LuaVM) {
	a, _, b, _ := i.ABC()
	a += 1

	// todo: optimize tail call!
	c := 0
	nArgs := _pushFuncAndArgs(a, b, vm)
	vm.Call(nArgs, c-1)
	_popResults(a, c, vm)
}

// R(A), ... ,R(A+C-2) := R(A)(R(A+1), ... ,R(A+B-1))
func call(i Instruction, vm LuaVM) {
	a, _, b, c := i.ABC()
	loca := a + 1

	// println(":::"+ vm.StackToString())
	nArgs := _pushFuncAndArgs(loca, b, vm)
	vm.Call(nArgs, c-1)
	_popResults(loca, c, vm)
}

func _pushFuncAndArgs(a, b int, vm LuaVM) (nArgs int) {
	if b >= 1 {
		vm.CheckStack(b)
		for i := a; i < a+b; i++ {
			vm.PushValue(i)
		}
		return b - 1
	} else { // b== 0
		// e.g.
		// f(1,g())
		//g().c == 0 and f().b == 0
		//x  [top]
		//g() -ret
		//...
		//g (put g()'s return here) [x]
		//1
		//f [a]
		_fixStack(a, vm)
		return vm.GetTop() - vm.RegisterCount() - 1
	}
}

func _fixStack(a int, vm LuaVM) {
	x := int(vm.ToInteger(-1)) // reg that params should be sent
	vm.Pop(1)

	vm.CheckStack(x - a) //push : x-a-1 +1(f)
	for i := a; i < x; i++ {
		vm.PushValue(i)
	}
	vm.Rotate(vm.RegisterCount()+1, x-a) //keep old stack safe.
}

func _popResults(a, c int, vm LuaVM) {
	if c == 1 {
		// no results
	} else if c > 1 {
		for i := a + c - 2; i >= a; i-- {
			vm.Replace(i)
		}
	} else { //c == 0
		// leave results on stack
		vm.CheckStack(1)
		// need to push to A reg but we leave it here.
		vm.PushInteger(int64(a))
	}
}

// return R(A), ... ,R(A+B-2) b-1 returns
func _return(i Instruction, vm LuaVM) {
	a, k, b, c := i.ABC()
	a += 1
	if k != 0 {
		// todo (upvalue)
	}
	if c != 0 {
		// todo (vararg)
	}

	if b == 1 {
		// no return values
	} else if b > 1 {
		// b-1 return values
		vm.CheckStack(b - 1)
		for i := a; i <= a+b-2; i++ {
			vm.PushValue(i)
		}
	} else { // b== 0
		_fixStack(a, vm)
	}
}

func return0(i Instruction, vm LuaVM) {
	//
}

// return R[a]
func return1(i Instruction, vm LuaVM) {
	a, _, _, _ := i.ABC()
	a += 1
	vm.PushValue(a)
	//back to poscall
}
