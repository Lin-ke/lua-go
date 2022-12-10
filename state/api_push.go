package state

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
