package vm

import (
	"fmt"
	. "luago54/api"
)

// R(A) := {}
func newTable(i Instruction, vm LuaVM) {
	a, k, b, c := i.ABC()
	// b:log2(hash size) + 1
	// c:arraysize
	a += 1
	if b > 0 {
		b = 1 << (b - 1)
	}
	//assert : k==0 == pc.Ax()==0
	if k != 0 { //extra arguments
		// EXTRAARG, which is the next instruction.
		c += Instruction(vm.Fetch()).Ax() << 8
		// 256*ax
	} else {
		vm.AddPC(1)
	}
	//vm.SetTop(a) /* correct top in case of emergency GC */
	vm.CreateTable(c, b)
	vm.Replace(a)
}

// R[A] := R[B][C]
func getI(i Instruction, vm LuaVM) {
	a, _, b, c := i.ABC()
	a += 1
	b += 1
	vm.GetI(b, int64(c))

	vm.Replace(a)
}

// R[A] := R[B][K[C]:string]
func getField(i Instruction, vm LuaVM) {
	a, _, b, c := i.ABC()
	a += 1
	b += 1
	vm.GetConst(c)
	if !vm.IsString(-1) {
		panic("getField's key must be a string")
	}
	vm.GetTable(b)
	vm.Replace(a)
}

// R(A) := R[B][R[C]]
func getTable(i Instruction, vm LuaVM) {
	a, _, b, c := i.ABC()
	a += 1
	b += 1
	c += 1
	vm.PushValue(c)
	vm.GetTable(b)
	vm.Replace(a)
}

// R(A)[R(B)] := RK(C)
func setTable(i Instruction, vm LuaVM) {
	a, k, b, c := i.ABC()
	a += 1
	b += 1

	vm.PushValue(b)
	vm.GetRK(c, k)
	vm.SetTable(a)
}

// R[A][B] := RK(C)
func setI(i Instruction, vm LuaVM) {
	a, k, b, c := i.ABC()
	a += 1

	vm.GetRK(c, k)
	vm.SetI(a, int64(b))
}
func printLuaval(val interface{}) {

	switch v := val.(type) {
	case bool:
		fmt.Printf("[%t]", v)
	case int64:
		fmt.Printf("[%d]", v)
	case float64:
		fmt.Printf("[%f]", v)
	case string:
		fmt.Printf("[%q]", v)
	case nil:
		fmt.Printf("[nil]")
	default: // other values
		fmt.Printf("[%s]", v)
	}
}

// R[A][K[B]:string] := RK(C)
func setField(i Instruction, vm LuaVM) {

	a, k, b, c := i.ABC()
	a += 1
	vm.GetConst(b)
	if !vm.IsString(-1) {
		panic("getField's key must be a string")
	}
	vm.GetRK(c, k) //push v
	vm.SetTable(a)
}

// R[A][C+i] := R[A+i], 1 <= i <= B
func setList(i Instruction, vm LuaVM) {
	a, k, b, c := i.ABC()
	loca := a + 1
	if b == 0 { /* get up to the top */
		b = vm.GetTop() - loca
	} else {
		// todo (set top)
	}
	if k != 0 {
		c += Instruction(vm.Fetch()).Ax() * 1 << 8
	}
	// c:last
	vm.CheckStack(b)
	for i := 1; i <= b; i++ {
		vm.PushValue(loca + i)
		vm.SetI(loca, (int64)(c+i))

	}
	// after :
}
