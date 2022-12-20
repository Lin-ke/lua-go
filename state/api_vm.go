package state

// secretly added apis for vm implementation
func (L *luaState) PC() int {
	return L.pc
}

func (L *luaState) AddPC(n int) {
	L.pc += n
}

// get pc from proto, and pc++
func (L *luaState) Fetch() uint32 {
	i := L.proto.Code[L.pc]
	L.pc++
	return i
}

// get const from proto according to index
func (L *luaState) GetConst(idx int) {
	c := L.proto.Constants[idx]
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
