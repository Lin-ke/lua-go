package state

// [-0, +1, e]
// http://www.lua.org/manual/5.4/manual.html#lua_len
// TODO: complete for other datastruct
func (L *luaState) Len(idx int) {
	val := L.stack.get(idx)

	if s, ok := val.(string); ok {
		L.stack.push(int64(len(s)))
	} else {
		panic("length error!")
	}
}

// [-n, +1, e]
// http://www.lua.org/manual/5.4/manual.html#lua_concat
// TODO: complete for other datastruct
func (L *luaState) Concat(n int) {
	if n == 0 {
		L.stack.push("")
	} else if n >= 2 {
		s0 := L.ToString(-1)
		L.stack.pop()

		for i := 1; i < n; i++ {
			if L.IsString(-1) {
				next := L.ToString(-1)
				L.stack.pop()
				s0 = next + s0
				continue
			}
			L.stack.push(s0)
			panic("concatenation error!")
		}
		L.stack.push(s0)
	}
	// n == 1, do nothing
}
