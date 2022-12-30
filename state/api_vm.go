package state

import (
	"fmt"
)

// secretly added apis for vm implementation
func (L *luaState) PC() int {
	return L.stack.pc
}

func (L *luaState) AddPC(n int) {
	L.stack.pc += n
}

// get pc from proto, and pc++
func (L *luaState) Fetch() uint32 {
	i := L.stack.closure.proto.Code[L.stack.pc]
	L.stack.pc++
	return i
}

// get const from proto according to index
func (L *luaState) GetConst(idx int) {
	c := L.stack.closure.proto.Constants[idx]
	L.stack.push(c)
}

// 使用常量或者使用寄存器
func (L *luaState) GetRK(rk, k int) {
	if k != 0 {
		L.GetConst(rk)
		return

	}
	L.PushValue(rk + 1)
}

func (L *luaState) Set(idx int, val interface{}) {
	L.stack.set(idx, val)
}

func (L *luaState) Push(val interface{}) {
	L.stack.push(val)
}

func (L *luaState) Get(idx int) interface{} {
	return L.stack.get(idx)
}

func (L *luaState) LoadProto(idx int) {
	stack := L.stack
	subProto := stack.closure.proto.Protos[idx]
	closure := newLuaClosure(subProto)
	stack.push(closure)

	for i, uvInfo := range subProto.Upvalues {
		uvIdx := int(uvInfo.Idx)
		if uvInfo.Instack == 1 { // open upvalue
			if stack.openuvs == nil {
				stack.openuvs = map[int]*upvalue{}
			}

			if openuv, found := stack.openuvs[uvIdx]; found {
				closure.upvals[i] = openuv
			} else {
				closure.upvals[i] = &upvalue{&stack.slots[uvIdx]}
				stack.openuvs[uvIdx] = closure.upvals[i]
			}
		} else { // closed upvalue
			closure.upvals[i] = stack.closure.upvals[uvIdx]
		}
	}
}

func (L *luaState) RegisterCount() int {
	return int(L.stack.closure.proto.MaxStackSize)
}

func (L *luaState) LoadVararg(n int) {
	if n < 0 { // all out
		n = len(L.stack.varargs)
	}

	L.stack.check(n)
	L.stack.pushN(L.stack.varargs, n)
}
func (L *luaState) TailCall(nArgs int) {
	val := L.stack.get(-(nArgs + 1))
	if c, ok := val.(*closure); ok {
		fmt.Printf("tailcall %s<%d,%d>\n", c.proto.Source,
			c.proto.LineDefined, c.proto.LastLineDefined)
		L.tailCallLuaClosure(nArgs, c)
	} else {
		panic("not function!")
	}
}

func (L *luaState) CloseUpvalues(a int) {
	for i, openuv := range L.stack.openuvs {
		if i >= a-1 {
			val := *openuv.val
			openuv.val = &val
			delete(L.stack.openuvs, i)
			if DEBUG.printUpval {
				fmt.Printf("closeupval:[%d]", i)
				printLuaval(val)
				fmt.Println()
			}
		}

	}
}
