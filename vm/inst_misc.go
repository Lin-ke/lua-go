package vm

import (
	"luago54/api"
)

// R(A) := R(B)
func move(i Instruction, vm api.LuaVM) {
	a, _, b, _ := i.ABC()
	// reg_index + 1 = stack_index.
	a += 1
	b += 1

	vm.Copy(b, a)
}

// pc+=sBx
func jmp(i Instruction, vm api.LuaVM) {
	sj := i.SJ()

	vm.AddPC(sj)
}
