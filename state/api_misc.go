package state

// [-0, +1, e]
// http://www.lua.org/manual/5.4/manual.html#lua_len
// TODO: complete for other datastruct
func (L *luaState) Len(idx int) {
	val := L.stack.get(idx)

	if s, ok := val.(string); ok {
		L.stack.push(int64(len(s)))
	} else if result, ok := callMetamethod(val, val, "__len", L); ok {

		L.stack.push(result)
	} else if t, ok := val.(*luaTable); ok {
		L.stack.push((int64)(t.len()))

	} else {
		panic("length error!")
	}
}

// [-0,+0,-]
// http://www.lua.org/manual/5.4/manual.html#lua_rawlen
func (L *luaState) RawLen(idx int) int {
	val := L.stack.get(idx)
	if s, ok := val.(string); ok {
		return len(s)
	} else if t, ok := val.(*luaTable); ok {
		return t.len()

	} else {
		panic("lenth error")
	}
}

// [-1, +(2|0), e]
// http://www.lua.org/manual/5.4/manual.html#lua_next
func (L *luaState) Next(idx int) bool {
	val := L.stack.get(idx)
	if t, ok := val.(*luaTable); ok {
		key := L.stack.pop()
		if nextKey := t.nextKey(key); nextKey != nil {
			L.stack.push(nextKey)
			L.stack.push(t.get(nextKey))
			return true
		}
		return false
	}
	panic("table expected!")
}

// [-n, +1, e]
// http://www.lua.org/manual/5.4/manual.html#lua_concat
// TODO: complete for other datastruct
func (L *luaState) Concat(n int) {
	if n == 0 {
		L.stack.push("")
	} else if n >= 2 {

		for i := 1; i < n; i++ {
			sb, okb := L.ToStringX(-1)
			sa, oka := L.ToStringX(-2)
			b := L.stack.pop()
			a := L.stack.pop()
			if okb && oka {
				L.stack.push(sa + sb)
				continue
			} else if result, ok := callMetamethod(a, b, "__concat", L); ok {
				L.stack.push(result)
				continue
			}

			panic("concatenation error!")
		}
	}
	// n == 1, do nothing
}

// [-1, +0, v]
// http://www.lua.org/manual/5.4/manual.html#lua_error
func (L *luaState) Error() int {
	err := L.stack.pop()
	panic(err)
}
