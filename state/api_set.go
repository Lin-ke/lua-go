package state

// [-2, +0, e] top == v , k  ... 0
// http://www.lua.org/manual/5.4/manual.html#lua_settable
func (L *luaState) SetTable(idx int) {
	t := L.stack.get(idx)
	v := L.stack.pop()
	k := L.stack.pop()
	L.setTable(t, k, v)
}

// [-1, +0, e]
// http://www.lua.org/manual/5.4/manual.html#lua_setfield
func (L *luaState) SetField(idx int, k string) {
	t := L.stack.get(idx)
	v := L.stack.pop()
	L.setTable(t, k, v)
}

// [-1, +0, e]
// http://www.lua.org/manual/5.4/manual.html#lua_seti
func (L *luaState) SetI(idx int, i int64) {
	t := L.stack.get(idx)
	v := L.stack.pop()
	L.setTable(t, i, v)
}

// t[k]=v
func (L *luaState) setTable(t, k, v luaValue) {
	if tbl, ok := t.(*luaTable); ok {
		tbl.put(k, v)
		return
	}

	panic("not a table!")
}
