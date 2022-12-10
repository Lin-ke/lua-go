package state

import (
	"fmt"
	Const "luago54/api"
)

// [-0, +0, –]
// http://www.lua.org/manual/5.4/manual.html#lua_typename
func (L *luaState) TypeName(tp Const.LuaType) string {
	switch tp {
	case Const.LUA_TNONE:
		return "no value"
	case Const.LUA_TNIL:
		return "nil"
	case Const.LUA_TBOOLEAN:
		return "boolean"
	case Const.LUA_TNUMBER:
		return "number"
	case Const.LUA_TSTRING:
		return "string"
	case Const.LUA_TTABLE:
		return "table"
	case Const.LUA_TFUNCTION:
		return "function"
	case Const.LUA_TTHREAD:
		return "thread"
	default:
		return "userdata"
	}
}

// [-0, +0, –]
// http://www.lua.org/manual/5.4/manual.html#lua_type
func (L *luaState) Type(idx int) Const.LuaType {
	// 计算了两次索引
	if L.stack.isValid(idx) {
		val := L.stack.tget(idx)
		return typeOf(val)
	}
	return Const.LUA_TNONE
}

// [-0, +0, –]
// http://www.lua.org/manual/5.4/manual.html#lua_isnone
func (L *luaState) IsNone(idx int) bool {
	return L.Type(idx) == Const.LUA_TNONE
}

// [-0, +0, –]
// http://www.lua.org/manual/5.4/manual.html#lua_isnil
func (L *luaState) IsNil(idx int) bool {
	return L.Type(idx) == Const.LUA_TNIL
}

// [-0, +0, –]
// http://www.lua.org/manual/5.4/manual.html#lua_isnoneornil
func (L *luaState) IsNoneOrNil(idx int) bool {
	return L.Type(idx) <= Const.LUA_TNIL
}

// [-0, +0, –]
// http://www.lua.org/manual/5.4/manual.html#lua_isboolean
func (L *luaState) IsBoolean(idx int) bool {
	return L.Type(idx) == Const.LUA_TBOOLEAN
}

// [-0, +0, –]
// http://www.lua.org/manual/5.4/manual.html#lua_istable
func (L *luaState) IsTable(idx int) bool {
	return L.Type(idx) == Const.LUA_TTABLE
}

// [-0, +0, –]
// http://www.lua.org/manual/5.4/manual.html#lua_isfunction
func (L *luaState) IsFunction(idx int) bool {
	return L.Type(idx) == Const.LUA_TFUNCTION
}

// [-0, +0, –]
// http://www.lua.org/manual/5.4/manual.html#lua_isthread
func (L *luaState) IsThread(idx int) bool {
	return L.Type(idx) == Const.LUA_TTHREAD
}

// [-0, +0, –]
// http://www.lua.org/manual/5.4/manual.html#lua_isstring
func (L *luaState) IsString(idx int) bool {
	t := L.Type(idx)
	return t == Const.LUA_TSTRING || t == Const.LUA_TNUMBER
}

// [-0, +0, –]
// http://www.lua.org/manual/5.4/manual.html#lua_isnumber
func (L *luaState) IsNumber(idx int) bool {
	_, ok := L.ToNumberX(idx)
	return ok
}

// [-0, +0, –]
// http://www.lua.org/manual/5.4/manual.html#lua_isinteger
func (L *luaState) IsInteger(idx int) bool {
	val := L.stack.get(idx)
	_, ok := val.(int64)
	return ok
}

// [-0, +0, –]
// http://www.lua.org/manual/5.4/manual.html#lua_toboolean
func (L *luaState) ToBoolean(idx int) bool {
	val := L.stack.get(idx)
	return convertToBoolean(val)
}

// [-0, +0, –]
// http://www.lua.org/manual/5.4/manual.html#lua_tointeger
func (L *luaState) ToInteger(idx int) int64 {
	i, _ := L.ToIntegerX(idx)
	return i
}

// [-0, +0, –]
// http://www.lua.org/manual/5.4/manual.html#lua_tointegerx
func (L *luaState) ToIntegerX(idx int) (int64, bool) {
	val := L.stack.get(idx)
	i, ok := val.(int64)
	return i, ok
}

// [-0, +0, –]
// http://www.lua.org/manual/5.4/manual.html#lua_tonumber
func (L *luaState) ToNumber(idx int) float64 {
	n, _ := L.ToNumberX(idx)
	return n
}

// [-0, +0, –]
// http://www.lua.org/manual/5.4/manual.html#lua_tonumberx
func (L *luaState) ToNumberX(idx int) (float64, bool) {
	val := L.stack.get(idx)
	return convertToFloat(val)
}

// [-0, +0, m]
// http://www.lua.org/manual/5.4/manual.html#lua_tostring
func (L *luaState) ToString(idx int) string {
	s, _ := L.ToStringX(idx)
	return s
}

func (L *luaState) ToStringX(idx int) (string, bool) {
	val := L.stack.get(idx)

	switch x := val.(type) {
	case string:
		return x, true
	case int64, float64:
		s := fmt.Sprintf("%v", x) // todo
		L.stack.set(idx, s)
		return s, true
	default:
		return "", false
	}
}
