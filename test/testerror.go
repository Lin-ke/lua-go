package test

import (
	"io/ioutil"
	"luago54/api"
	"luago54/state"
	"os"
)

func error(ls api.LuaState) int {
	return ls.Error()
}

func pCall(ls api.LuaState) int {
	nArgs := ls.GetTop() - 1
	status := ls.PCall(nArgs, -1, 0)
	ls.PushBoolean(status == api.LUA_OK)
	ls.Insert(1)
	return ls.GetTop()
}

func Test010() {
	if len(os.Args) > 1 {
		data, err := ioutil.ReadFile(os.Args[1])
		if err != nil {
			panic(err)
		}

		ls := state.New()
		ls.Register("print", print)
		ls.Register("getmetatable", getMetatable)
		ls.Register("setmetatable", setMetatable)
		ls.Register("next", next)
		ls.Register("pairs", pairs)
		ls.Register("ipairs", iPairs)
		ls.Register("error", error)
		ls.Register("pcall", pCall)
		ls.Load(data, os.Args[1], "b")
		ls.Call(0, 0)
	}
}
