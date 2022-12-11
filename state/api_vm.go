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

func (L *luaState) GetRK(rk int) {
	if rk > 0xFF { // constant
		L.GetConst(rk & 0xFF)
	} else { // register
		L.PushValue(rk + 1)
	}
}

func (L *luaState) Set(idx int, val luaValue) {
	L.stack.set(idx, val)
}
