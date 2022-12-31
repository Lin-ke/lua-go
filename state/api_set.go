package state

import "luago54/api"

// [-2, +0, e] top == v , k  ... 0
// http://www.lua.org/manual/5.4/manual.html#lua_settable
func (L *luaState) SetTable(idx int) {
	t := L.stack.get(idx)
	v := L.stack.pop()
	k := L.stack.pop()
	L.setTable(t, k, v, false)
}

// [-1, +0, e]
// http://www.lua.org/manual/5.4/manual.html#lua_setfield
func (L *luaState) SetField(idx int, k string) {
	t := L.stack.get(idx)
	v := L.stack.pop()
	L.setTable(t, k, v, false)
}

// [-1, +0, e]
// http://www.lua.org/manual/5.4/manual.html#lua_seti
func (L *luaState) SetI(idx int, i int64) {
	t := L.stack.get(idx)
	v := L.stack.pop()
	L.setTable(t, i, v, false)
}

// [-1, +0, e]
// http://www.lua.org/manual/5.4/manual.html#lua_setglobal
func (L *luaState) SetGlobal(name string) {
	t := L.registry.get(api.LUA_RIDX_GLOBALS)
	v := L.stack.pop()
	L.setTable(t, name, v, true)
}

// [-0, +0, e]
// http://www.lua.org/manual/5.4/manual.html#lua_register
func (L *luaState) Register(name string, f api.GoFunction) {
	L.PushGoFunction(f)
	L.SetGlobal(name)
}

// [-1, +0, –]
// http://www.lua.org/manual/5.4/manual.html#lua_setmetatable
func (L *luaState) SetMetatable(idx int) {
	val := L.stack.get(idx)
	mtVal := L.stack.pop()

	if mtVal == nil {
		setMetatable(val, nil, L)
	} else if mt, ok := mtVal.(*luaTable); ok {
		setMetatable(val, mt, L)
	} else {
		panic("table expected!") // todo
	}
}

// [-2, +0, m]
// http://www.lua.org/manual/5.4/manual.html#lua_rawset
func (L *luaState) RawSet(idx int) {
	t := L.stack.get(idx)
	v := L.stack.pop()
	k := L.stack.pop()
	L.setTable(t, k, v, true)
}

// [-1, +0, m]
// http://www.lua.org/manual/5.4/manual.html#lua_rawseti
func (L *luaState) RawSetI(idx int, i int64) {
	t := L.stack.get(idx)
	v := L.stack.pop()
	L.setTable(t, i, v, true)
}

// t[k]=v
func (L *luaState) setTable(t, k, v luaValue, raw bool) {
	if tbl, ok := t.(*luaTable); ok {
		if raw || tbl.get(k) != nil || !tbl.hasMetafield("__newindex") {
			tbl.put(k, v)
			return
		}
		if !raw {
			if mf := getMetafield(t, "__newindex", L); mf != nil {
				switch x := mf.(type) {
				case *luaTable:
					L.setTable(x, k, v, false)
					return
				case *closure:
					L.stack.push(mf)
					L.stack.push(t)
					L.stack.push(k)
					L.stack.push(v)
					L.Call(3, 0)
					return
				}
			}
		}

		panic("not a table!")
	}
}
