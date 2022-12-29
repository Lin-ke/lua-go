package state

type luaState struct {
	stack *luaStack // current call stack.
}
type LuaType = int
type ArithOp = int
type CompareOp int

func New() *luaState {
	return &luaState{
		stack: newLuaStack(20),
	}
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
