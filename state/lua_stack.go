package state

import "fmt"

var printstack bool = false

// stack contains pc and prototype.
type luaStack struct {
	/* virtual stack */
	slots []luaValue
	top   int
	/* call info */
	closure *closure
	varargs []luaValue
	pc      int
	/* linked list */
	prev *luaStack
}

func newLuaStack(size int) *luaStack {
	return &luaStack{
		slots: make([]luaValue, size),
		top:   0,
	}
}

/*
检查枝 的 空 闲 空间是否还可以容纳（推入）至少 n 个值，如果不满足这个条 件 ， 则 调用
Go 语 言 内 置的 append （） 函数进行扩容
*/
func (stack *luaStack) check(n int) {
	free := len(stack.slots) - stack.top
	for i := free; i < n; i++ {
		stack.slots = append(stack.slots, nil)
	}
}

func (stack *luaStack) push(val luaValue) {
	if stack.top == len(stack.slots) {
		panic("stack overflow!")
	}
	stack.slots[stack.top] = val
	stack.top++
} // 需要考虑线程么？

func (stack *luaStack) pop() luaValue {
	if stack.top < 1 {
		panic("stack underflow!")
	}
	stack.top--
	val := stack.slots[stack.top]
	stack.slots[stack.top] = nil
	return val
}

// idx -stack.top => -1 => 1 => stack.top
// 到后面会完善

func (stack *luaStack) absIndex(idx int) int {
	if idx >= 0 {
		return idx
	}
	return idx + stack.top + 1
}

func (stack *luaStack) isValid(idx int) bool {
	absIdx := stack.absIndex(idx)
	return absIdx > 0 && absIdx <= stack.top
}

func (stack *luaStack) get(idx int) luaValue {
	absIdx := stack.absIndex(idx)
	if absIdx > 0 && absIdx <= stack.top {
		return stack.slots[absIdx-1]
	}
	// or panic?
	return nil
}

func (stack *luaStack) tget(idx int) luaValue {

	return stack.slots[stack.absIndex(idx)-1]

}

func (stack *luaStack) set(idx int, val luaValue) {
	absIdx := stack.absIndex(idx)
	if absIdx > 0 && absIdx <= stack.top {
		stack.slots[absIdx-1] = val // 设定开始是1的问题
		//debug
		if printstack {
			fmt.Printf("Set slot %d\n", absIdx)
			printStack(stack)
		}

		return
	}
	panic("invalid index!")
}

// 反转
// 注意这些index是go的（0-> top-1)
func (stack *luaStack) reverse(from, to int) {
	slots := stack.slots
	for from < to {
		slots[from], slots[to] = slots[to], slots[from]
		from++
		to--
	}
}

func (stack *luaStack) pushN(vals []luaValue, n int) {
	nVals := len(vals)
	if n < 0 {
		n = nVals
	}

	for i := 0; i < n; i++ {
		if i < nVals {
			stack.push(vals[i])
		} else {
			stack.push(nil)
		}
	}
}

func (stack *luaStack) popN(n int) []luaValue {
	vals := make([]luaValue, n)
	for i := n - 1; i >= 0; i-- {
		vals[i] = stack.pop()
	}
	return vals
}
