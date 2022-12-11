package vm

// lvm.c

// R(A), R(A+1), ..., R(A+B) := nil
func loadNil(i Instruction, vm LuaVM) {
	a, b, _ := i.ABC()
	a += 1

	vm.PushNil()
	for i := a; i <= a+b; i++ {
		vm.Copy(-1, i)
	}
	vm.Pop(1)
}

// R[A] := false
func loadFalse(i Instruction, vm LuaVM) {
	a, _, _ := i.ABC()
	vm.PushBoolean(false)
	vm.Replace(a)
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
