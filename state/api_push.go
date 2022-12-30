package state

import "luago54/api"

// [-0, +1, –]
// http://www.lua.org/manual/5.4/manual.html#lua_pushnil
func (L *luaState) PushNil() {
	L.stack.push(nil)
}

// [-0, +1, –]
// http://www.lua.org/manual/5.4/manual.html#lua_pushboolean
func (L *luaState) PushBoolean(b bool) {
	L.stack.push(b)
}

// [-0, +1, –]
// http://www.lua.org/manual/5.4/manual.html#lua_pushinteger
func (L *luaState) PushInteger(n int64) {
	L.stack.push(n)
}

// [-0, +1, –]
// http://www.lua.org/manual/5.4/manual.html#lua_pushnumber
func (L *luaState) PushNumber(n float64) {
	L.stack.push(n)
}

// [-0, +1, m]
// http://www.lua.org/manual/5.4/manual.html#lua_pushstring
func (L *luaState) PushString(s string) {
	L.stack.push(s)
}

// [-0, +1, –]
// http://www.lua.org/manual/5.4/manual.html#lua_pushcfunction
func (L *luaState) PushGoFunction(f api.GoFunction) {
	L.stack.push(newGoClosure(f, 0))
}

// [-n, +1, m]
// http://www.lua.org/manual/5.4/manual.html#lua_pushcclosure
func (L *luaState) PushGoClosure(f api.GoFunction, n int) {
	closure := newGoClosure(f, n)
	for i := n; i > 0; i-- {
		val := L.stack.pop()
		closure.upvals[i-1] = &upvalue{&val}
	}
	L.stack.push(closure)
}

// [-0, +1, –]
// http://www.lua.org/manual/5.4/manual.html#lua_pushglobaltable
func (L *luaState) PushGlobalTable() {
	global := L.registry.get(api.LUA_RIDX_GLOBALS)
	L.stack.push(global)
}
