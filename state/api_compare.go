package state

import "luago54/api"

// NAN 的比较：
// NAN return false to everything, nan == nan = false.
// golang has achieved it.
// [-0, +0, –]
// http://www.lua.org/manual/5.4/manual.html#lua_rawequal
func (L *luaState) RawEqual(idx1, idx2 int) bool {
	if !L.stack.isValid(idx1) || !L.stack.isValid(idx2) {
		return false
	}

	a := L.stack.get(idx1)
	b := L.stack.get(idx2)
	return _eq(a, b, nil)
}

// [-0, +0, e]
// http://www.lua.org/manual/5.4/manual.html#lua_compare
func (L *luaState) Compare(idx1, idx2 int, op api.CompareOp) bool {
	if !L.stack.isValid(idx1) || !L.stack.isValid(idx2) {
		return false
	}

	a := L.stack.get(idx1)
	b := L.stack.get(idx2)
	switch op {
	case api.LUA_OPEQ:
		return _eq(a, b, L)
	case api.LUA_OPLT:
		return _lt(a, b, L)
	case api.LUA_OPLE:
		return _le(a, b, L)
	default:
		panic("invalid compare op!")
	}
}

func _eq(a, b luaValue, ls *luaState) bool {
	switch x := a.(type) {
	case nil:
		return b == nil
	case bool:
		y, ok := b.(bool)
		return ok && x == y
	case string:
		y, ok := b.(string)
		return ok && x == y
	case int64:
		switch y := b.(type) {
		case int64:
			return x == y
		case float64:
			return float64(x) == y
		default:
			return false
		}
	case float64:
		switch y := b.(type) {
		case float64:
			return x == y
		case int64:
			return x == float64(y)
		default:
			return false
		}
	case *luaTable:
		// call __eq when a!=b and type a == type b == *table.
		if y, ok := b.(*luaTable); ok && x != y && ls != nil {
			if result, ok := callMetamethod(x, y, "__eq", ls); ok {
				return convertToBoolean(result)
			}
		}
		return a == b
	default:
		return a == b
	}
}

func _lt(a, b luaValue, ls *luaState) bool {
	switch x := a.(type) {
	case string:
		if y, ok := b.(string); ok {
			return x < y
		}
	case int64:
		switch y := b.(type) {
		case int64:
			return x < y
		case float64:
			return float64(x) < y
		}
	case float64:
		switch y := b.(type) {
		case float64:
			return x < y
		case int64:
			return x < float64(y)
		}
	}

	if result, ok := callMetamethod(a, b, "__lt", ls); ok {
		return convertToBoolean(result)
	} else {
		panic("comparison error!")
	}
}

func _le(a, b luaValue, ls *luaState) bool {
	switch x := a.(type) {
	case string:
		if y, ok := b.(string); ok {
			return x <= y
		}
	case int64:
		switch y := b.(type) {
		case int64:
			return x <= y
		case float64:
			return float64(x) <= y
		}
	case float64:
		switch y := b.(type) {
		case float64:
			return x <= y
		case int64:
			return x <= float64(y)
		}
	}

	if result, ok := callMetamethod(a, b, "__le", ls); ok {
		return convertToBoolean(result)
	} else if result, ok := callMetamethod(b, a, "__lt", ls); ok {
		return !convertToBoolean(result)
	} else {
		panic("comparison error!")
	}
}
