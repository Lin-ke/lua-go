package vm

import api "luago54/api"

// For apis are exposed to Users, VM needs to hide binchunk details.
type LuaVM interface {
	api.LuaState
	PC() int
	AddPC(n int)
	Fetch() uint32
	GetConst(idx int)
	GetRK(rk int)
	GetReg() // get register(the stack)
}
