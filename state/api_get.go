package state

import "luago54/api"

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
func (L *luaState) GetTable(idx int) api.LuaType {
	t := L.stack.get(idx)
	k := L.stack.pop()
	return L.getTable(t, k, false)
}

// [-0, +1, e] field is hashkey, I is index
// http://www.lua.org/manual/5.4/manual.html#lua_getfield
func (L *luaState) GetField(idx int, k string) api.LuaType {
	t := L.stack.get(idx)
	return L.getTable(t, k, false)
}

// [-0, +1, e]
// http://www.lua.org/manual/5.4/manual.html#lua_geti
func (L *luaState) GetI(idx int, i int64) api.LuaType {
	t := L.stack.get(idx)
	return L.getTable(t, i, false)
}

// [-0, +1, e]
// http://www.lua.org/manual/5.4/manual.html#lua_getglobal
func (L *luaState) GetGlobal(name string) api.LuaType {
	t := L.registry.get(api.LUA_RIDX_GLOBALS)
	return L.getTable(t, name, false)
}

// [-1, +1, –]
// http://www.lua.org/manual/5.4/manual.html#lua_rawget
func (L *luaState) RawGet(idx int) api.LuaType {
	t := L.stack.get(idx)
	k := L.stack.pop()
	return L.getTable(t, k, true)
}

// [-0, +1, –]
// http://www.lua.org/manual/5.4/manual.html#lua_rawgeti
func (L *luaState) RawGetI(idx int, i int64) api.LuaType {
	t := L.stack.get(idx)
	return L.getTable(t, i, true)
}

// [-0, +(0|1), –]
// http://www.lua.org/manual/5.4/manual.html#lua_getmetatable
func (L *luaState) GetMetatable(idx int) bool {
	val := L.stack.get(idx)

	if mt := getMetatable(val, L); mt != nil {
		L.stack.push(mt)
		return true
	} else {
		return false
	}
}

// push(t[k])
func (L *luaState) getTable(t, k luaValue, raw bool) api.LuaType {

	printTable(t.(*luaTable))
	if tbl, ok := t.(*luaTable); ok {
		v := tbl.get(k)
		// neglect metamethod
		if raw || v != nil || !tbl.hasMetafield("__index") {
			L.stack.push(v)
			return typeOf(v)
		}
	}
	// (t is not a table || k not exists) and raw
	if !raw {
		if mf := getMetafield(t, "__index", L); mf != nil {
			switch x := mf.(type) {
			case *luaTable:
				return L.getTable(x, k, false)
			case *closure:
				L.stack.push(mf)
				L.stack.push(t)
				L.stack.push(k)
				L.Call(2, 1)
				v := L.stack.get(-1)
				return typeOf(v)
			}
		}
	}

	panic("index error!")
}
