package state

import (
	"fmt"
)

func printStack(L *luaStack) {
	fmt.Printf("top:%d", L.top)
	for i := 0; i < L.top; i++ {
		printLuaval(L.slots[i])
	}
	fmt.Println()
}

func printLuaval(val luaValue) {

	switch v := val.(type) {
	case bool:
		fmt.Printf("[%t]", v)
	case int64:
		fmt.Printf("[%d]", v)
	case float64:
		fmt.Printf("[%f]", v)
	case string:
		fmt.Printf("[%q]", v)
	case nil:
		fmt.Printf("[nil]")
	default: // other values
		fmt.Printf("[%s]", v)
	}
}

func printTable(Table *luaTable) {
	fmt.Printf("\narray:%d\n", Table.len())
	for _, itr := range Table.arr {
		printLuaval(itr)
	}
	fmt.Println()
	fmt.Printf("map:%d\n", len(Table._map))
	for k, v := range Table._map {
		printLuaval(k)
		print(":")
		printLuaval(v)
		print("|")
	}
	fmt.Println()
}
