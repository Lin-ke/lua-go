package main

// import (
// 	. "luago54/api"
// 	_ "luago54/binchunk"
// 	"luago54/state"
// )

// func main() {
// 	ls := state.New()
// 	ls.PushInteger(1)
// 	ls.PushString("2.0")
// 	ls.PushString("3.0")
// 	ls.PushNumber(4.0)
// 	printStack(ls)

// 	ls.Arith(LUA_OPADD)
// 	printStack(ls)
// 	ls.Arith(LUA_OPBNOT)
// 	printStack(ls)
// 	ls.Len(2)
// 	printStack(ls)
// 	ls.Concat(3)
// 	printStack(ls)
// 	ls.PushBoolean(ls.Compare(1, 2, LUA_OPEQ))
// 	printStack(ls)
// }
