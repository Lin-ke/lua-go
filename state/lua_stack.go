package state

type luaStack struct {
	slots []luaValue
	top   int
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
func (L *luaStack) check(n int) {
	free := len(L.slots) - L.top
	for i := free; i < n; i++ {
		L.slots = append(L.slots, nil)
	}
}

func (L *luaStack) push(val luaValue) {
	if L.top == len(L.slots) {
		panic("stack overflow!")
	}
	L.slots[L.top] = val
	L.top++
} // 需要考虑线程么？

func (L *luaStack) pop() luaValue {
	if L.top < 1 {
		panic("stack underflow!")
	}
	L.top--
	val := L.slots[L.top]
	L.slots[L.top] = nil
	return val
}

// idx -L.top => -1 => 1 => L.top
// 到后面会完善

func (L *luaStack) absIndex(idx int) int {
	if idx >= 0 {
		return idx
	}
	return idx + L.top + 1
}

func (L *luaStack) isValid(idx int) bool {
	absIdx := L.absIndex(idx)
	return absIdx > 0 && absIdx <= L.top
}

func (L *luaStack) get(idx int) luaValue {
	absIdx := L.absIndex(idx)
	if absIdx > 0 && absIdx <= L.top {
		return L.slots[absIdx-1]
	}
	return nil
}

func (L *luaStack) tget(idx int) luaValue {

	return L.slots[L.absIndex(idx)-1]

}

func (L *luaStack) set(idx int, val luaValue) {
	absIdx := L.absIndex(idx)
	if absIdx > 0 && absIdx <= L.top {
		L.slots[absIdx-1] = val // 设定开始是1的问题
		return
	}
	panic("invalid index!")
}

// 反转
// 注意这些index是go的（0-> top-1)
func (L *luaStack) reverse(from, to int) {
	slots := L.slots
	for from < to {
		slots[from], slots[to] = slots[to], slots[from]
		from++
		to--
	}
}
