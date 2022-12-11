package vm

// R(A) := R(B)
func move(i Instruction, vm LuaVM) {
	a, b, _ := i.ABC()
	// reg_index + 1 = stack_index.
	a += 1
	b += 1

	vm.Copy(b, a)
}

// pc+=sBx
func jmp(i Instruction, vm LuaVM) {
	sj := i.SJ()

	vm.AddPC(sj)
}
