package vm

import (
	. "luago54/api"
)

// lvm.c

// R(A), R(A+1), ..., R(A+B) := nil
func loadNil(i Instruction, vm LuaVM) {
	a, b, _ := i.ABC()
	for i := a + 1; i <= a+b+1; i++ {
		vm.Set(i, nil)
	}
}

// R[A] := false
func loadFalse(i Instruction, vm LuaVM) {
	a, _, _ := i.ABC()
	vm.Set(a+1, false)
}
func LOADTRUE(i Instruction, vm LuaVM) {
	a, _, _ := i.ABC()
	vm.Set(a+1, true)
}

// R(A) := Kst(Bx)
func loadK(i Instruction, vm LuaVM) {
	a, bx := i.ABx()
	a += 1

	vm.GetConst(bx)
	vm.Replace(a)
}

// R(A) := Kst(extra arg)
func loadKx(i Instruction, vm LuaVM) {
	a, _ := i.ABx()
	a += 1
	ax := Instruction(vm.Fetch()).Ax()

	//vm.CheckStack(1)
	vm.GetConst(ax)
	vm.Replace(a)
}

// R[A] := false; pc++	(*)
func lFalseSkip(i Instruction, vm LuaVM) {
	a, _, _ := i.ABC()
	a += 1
	vm.Set(a, false)
	vm.AddPC(1)
}

func loadI(i Instruction, vm LuaVM) {
	a, sBx := i.AsBx()
	a += 1
	vm.Set(a, sBx)
}
func loadF(i Instruction, vm LuaVM) {
	a, sBx := i.AsBx()
	a += 1
	vm.Set(a, (float64)(sBx))
}
