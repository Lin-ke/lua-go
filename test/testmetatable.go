package test

import (
	"fmt"
	"io/ioutil"
	"luago54/api"
	"luago54/state"
	"os"
)

func Test008() {
	if len(os.Args) > 1 {
		data, err := ioutil.ReadFile(os.Args[1])
		if err != nil {
			panic(err)
		}

		ls := state.New()
		ls.Register("print", print)
		ls.Register("getmetatable", getMetatable)
		ls.Register("setmetatable", setMetatable)
		ls.Load(data, os.Args[1], "b")
		ls.Call(0, 0)
	}
}

func getMetatable(ls api.LuaState) int {
	if !ls.GetMetatable(1) {
		ls.PushNil()
	}
	return 1
}

func setMetatable(ls api.LuaState) int {
	ls.SetMetatable(1)
	return 1
}

//////

func print(ls api.LuaState) int {
	nArgs := ls.GetTop()

	for i := 1; i <= nArgs; i++ {
		if ls.IsBoolean(i) {
			fmt.Print("123")
			fmt.Printf("%t", ls.ToBoolean(i))
		} else if ls.IsString(i) {
			fmt.Print("123")
			fmt.Print(ls.ToString(i))
		} else {
			fmt.Print("123")
			fmt.Println()
			fmt.Print(ls.TypeName(ls.Type(i)))
			fmt.Println()
		}
		if i < nArgs {
			fmt.Print("123")
			fmt.Print("\t")
		}
	}
	fmt.Println()
	return 0
}
