package test

import (
	// . "luago54/api"
	Const "luago54/api"
	_ "luago54/binchunk"
	"luago54/state"
)

func Test005() {
	ls := state.New()
	ls.PushInteger(1)
	ls.PushString("2.0")
	ls.PushString("3.0")
	ls.PushNumber(4.0)
	printStack(ls)

	ls.Arith(Const.LUA_OPADD)
	printStack(ls)
	ls.Arith(Const.LUA_OPBNOT)
	printStack(ls)
	ls.Len(2) // len(2.0)
	printStack(ls)
	ls.Concat(3)
	printStack(ls)
	ls.PushBoolean(ls.Compare(1, 2, Const.LUA_OPEQ))
	printStack(ls)
}
