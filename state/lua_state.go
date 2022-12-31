package state

import "luago54/api"

type luaState struct {
	registry *luaTable
	stack    *luaStack // current call stack.
}

func New() *luaState {
	registry := newLuaTable(0, 0)
	registry.put(api.LUA_RIDX_GLOBALS, newLuaTable(0, 0)) // global table

	ls := &luaState{registry: registry}
	ls.pushLuaStack(newLuaStack(api.LUA_MINSTACK, ls))
	return ls
}

// state works as  a callstack.
/*state ->
stack f() ->
stack g()-> nil
*/

func (L *luaState) pushLuaStack(stack *luaStack) {
	stack.prev = L.stack
	L.stack = stack
}

func (L *luaState) popLuaStack() {
	stack := L.stack
	L.stack = stack.prev
	stack.prev = nil
}

