package vm

import (
	"luago54/api"
)

// R[A+1] := R[B]; R[A] := R[B][RK(C):string]
func self(i Instruction, vm api.LuaVM) {
	a, k, b, c := i.ABC()
	a += 1
	b += 1
	vm.Copy(b, a+1)
	vm.GetRK(c, k)
	vm.GetTable(b)
	vm.Replace(a)
}

// R(A) := closure(KPROTO[Bx])
func closure(i Instruction, vm api.LuaVM) {
	a, bx := i.ABx()
	a += 1

	vm.LoadProto(bx)
	vm.Replace(a)
}

// R[A], R[A+1], ..., R[A+C-2] = vararg
func vararg(i Instruction, vm api.LuaVM) {
	a, _, _, c := i.ABC()
	a += 1

	if c != 1 { // c==0 or c>1
		vm.LoadVararg(c - 1)
		_popResults(a, c, vm)
	}
}

// adjust vararg parameters: put func and fixparams to the top of the stack
func varargPrep(i Instruction, vm api.LuaVM) {
	nfixparams, _, _, _ := i.ABC()
	// seem useless, for we already pushed fixparams.
	if nfixparams != 0 {
		// do prep/
	}

}

// return R(A)(R(A+1), ... ,R(A+B-1))
func tailCall(i Instruction, vm api.LuaVM) {
	a, k, b, _ := i.ABC()
	a += 1
	if k != 0 {
		vm.CloseUpvalues(0)
	}

	nArgs := _pushFuncAndArgs(a, b, vm)
	vm.TailCall(nArgs)
	// tailcall  will be followed by a return
}

// return R(A)(R(A+1), ... ,R(A+B-1))
// func tailCall(i Instruction, vm api.LuaVM) {
// 	a, k, b, _ := i.ABC()
// 	a += 1
// 	if k != 0 {
// 		vm.CloseUpvalues(0)
// 	}
// 	// todo: optimize tail call!
// 	c := 0
// 	nArgs := _pushFuncAndArgs(a, b, vm)
// 	vm.Call(nArgs, c-1)
//  //must be all out
// 	_popResults(a, c, vm)
// }

// R(A), ... ,R(A+C-2) := R(A)(R(A+1), ... ,R(A+B-1))
func call(i Instruction, vm api.LuaVM) {
	a, _, b, c := i.ABC()
	loca := a + 1

	// fmt.Println(":::"+ vm.StackToString())
	nArgs := _pushFuncAndArgs(loca, b, vm)
	vm.Call(nArgs, c-1)
	_popResults(loca, c, vm)
}

func _pushFuncAndArgs(a, b int, vm api.LuaVM) (nArgs int) {
	if b >= 1 {
		vm.CheckStack(b)
		for i := a; i < a+b; i++ {
			vm.PushValue(i)
		}
		return b - 1
	} else { // b== 0
		// e.g.
		// f(1,g())
		//g'return.c == 0(all out) and f'call.b == 0(all in)
		//x  [top]
		//g() -return
		//...
		//g (put g()'s return here) [x]
		//1
		//f [a]
		_fixStack(a, vm)
		return vm.GetTop() - vm.RegisterCount() - 1
	}
}

func _fixStack(a int, vm api.LuaVM) {
	x := int(vm.ToInteger(-1)) // reg that params should be sent
	vm.Pop(1)

	vm.CheckStack(x - a) //push : x-a-1 +1(f)
	for i := a; i < x; i++ {
		vm.PushValue(i)
	}
	vm.Rotate(vm.RegisterCount()+1, x-a)
}

func _popResults(a, c int, vm api.LuaVM) {
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
func _return(i Instruction, vm api.LuaVM) {
	a, k, b, c := i.ABC()
	a += 1
	if k != 0 {
		vm.CloseUpvalues(0)
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

func return0(i Instruction, vm api.LuaVM) {
	//
}

// return R[a]
func return1(i Instruction, vm api.LuaVM) {
	a, _, _, _ := i.ABC()
	a += 1
	vm.CheckStack(1)
	vm.PushValue(a)
	//back to pos
}
