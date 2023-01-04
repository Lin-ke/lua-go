package state

import (
	"fmt"
	"luago54/api"
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

func (L *luaState) LastInst() uint32 {
	if L.stack.pc == 0 {
		panic("no lastInst")
	}

	return L.stack.closure.proto.Code[L.stack.pc-1]
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

func (L *luaState) CallMetaMethod(mmName string) {
	var mm luaValue
	a := L.Get(-2)
	b := L.Get(-1)
	if mm = getMetafield(a, mmName, L); mm == nil {
		if mm = getMetafield(b, mmName, L); mm == nil {
			return // call failed.
		}
	}
	L.stack.check(2)
	L.Push(mm)
	L.Insert(-3)
	L.Call(2, 1)

}

func (L *luaState) RawArith(op api.ArithOp) bool {
	var a, b luaValue // operands
	b = L.stack.pop()
	if op != api.LUA_OPUNM && op != api.LUA_OPBNOT {
		a = L.stack.pop()
	} else {
		a = b
	}
	// pop, then push result.
	operator := operators[op]
	if result := _rawarith(a, b, operator); result != nil {
		L.stack.push(result)
	} else {
		// result == nil
		return false
	}
	return true
}
func _rawarith(a, b luaValue, op operator) luaValue {
	if op.floatFunc == nil { // bitwise
		if x, ok := convertToInteger(a); ok {
			if y, ok := convertToInteger(b); ok {
				return op.integerFunc(x, y)
			}
		}
	} else { // arith
		if op.integerFunc != nil { // add,sub,mul,mod,idiv,unm
			if x, ok := a.(int64); ok {
				if y, ok := b.(int64); ok {
					return op.integerFunc(x, y)
				}
			}
		}
		if x, ok := convertToFloat(a); ok {
			if y, ok := convertToFloat(b); ok {
				return op.floatFunc(x, y)
			}
		}
	}
	return nil
}

// Insert a variable in the list of to-be-closed variables.
func (L *luaState) NewTbcUpval(idx int) {

	val := L.stack.get(idx)
	flag := false // have __close method
	if val == nil || val == false {
		flag = true

	} else if mm := getMetafield(val, "__close", L); mm != nil {
		flag = true
	}
	if flag {
		if L.stack.tbcuvs == nil {
			L.stack.tbcuvs = make([]int, 1)
			L.stack.tbcuvs[0] = idx
			return
		} else {
			L.stack.tbcuvs = append(L.stack.tbcuvs, idx)
			return
		}
	}

	panic("non-closable value")

}
func (L *luaState) CloseTbc(idx int) {
	for p, a := range L.stack.tbcuvs {
		if a >= idx {
			val := L.stack.get(a)
			if val == nil || val == false {
				continue
			}
			callOneArgMM(val, "__close", L) // nil or false returns false.
			if p == len(L.stack.tbcuvs)-1 {
				L.stack.tbcuvs = L.stack.tbcuvs[:p]
			} else {
				L.stack.tbcuvs = append(L.stack.tbcuvs[:p], L.stack.tbcuvs[p+1:]...)
			}
		}

	}
}
