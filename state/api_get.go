package state

// [-0, +1, m]
// http://www.lua.org/manual/5.4/manual.html#lua_newtable
func (L *luaState) NewTable() {
	L.CreateTable(0, 0)
}

// [-0, +1, m]
// http://www.lua.org/manual/5.4/manual.html#lua_createtable
func (L *luaState) CreateTable(nArr, nRec int) {
	t := newLuaTable(nArr, nRec)
	L.stack.push(t)
}

// [-1, +1, e]
// http://www.lua.org/manual/5.4/manual.html#lua_gettable
func (L *luaState) GetTable(idx int) LuaType {
	t := L.stack.get(idx)
	k := L.stack.pop()
	return L.getTable(t, k)
}

// [-0, +1, e] field is hashkey, I is index
// http://www.lua.org/manual/5.4/manual.html#lua_getfield
func (L *luaState) GetField(idx int, k string) LuaType {
	t := L.stack.get(idx)
	return L.getTable(t, k)
}

// [-0, +1, e]
// http://www.lua.org/manual/5.4/manual.html#lua_geti
func (L *luaState) GetI(idx int, i int64) LuaType {
	t := L.stack.get(idx)
	return L.getTable(t, i)
}

// push(t[k])
func (L *luaState) getTable(t, k luaValue) LuaType {
	if tbl, ok := t.(*luaTable); ok {
		v := tbl.get(k)
		L.stack.push(v)
		return typeOf(v)
	}

	panic("not a table!") // todo
}
