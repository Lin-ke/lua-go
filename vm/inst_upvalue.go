package vm

import (
	"luago54/api"
)

// R(A) := UpValue[B]
func getUpval(i Instruction, vm api.LuaVM) {
	a, _, b, _ := i.ABC()
	a += 1
	b += 1

	vm.Copy(api.LuaUpvalueIndex(b), a)
}

// UpValue[B] := R(A)
func setUpval(i Instruction, vm api.LuaVM) {
	a, _, b, _ := i.ABC()
	a += 1
	b += 1

	vm.Copy(a, api.LuaUpvalueIndex(b))
}

// R(A) := UpValue[B][RK(C)]
func getTabUp(i Instruction, vm api.LuaVM) {
	a, _, b, c := i.ABC()
	a += 1
	b += 1

	vm.GetConst(c)
	vm.GetTable(api.LuaUpvalueIndex(b))
	vm.Replace(a)
}

// UpValue[A][K[B]:string] := RK(C)
func setTabUp(i Instruction, vm api.LuaVM) {
	a, k, b, c := i.ABC()
	a += 1

	vm.GetConst(b)
	vm.GetRK(c, k)
	vm.SetTable(api.LuaUpvalueIndex(a))
}

func close(i Instruction, vm api.LuaVM) {
	a, _, _, _ := i.ABC()
	vm.CloseUpvalues(a)

}
