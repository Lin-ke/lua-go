package api

// For apis are exposed to Users, VM needs to hide binchunk details.
type LuaVM interface {
	LuaState
	PC() int
	AddPC(n int)
	Fetch() uint32
	GetConst(idx int)
	GetRK(rk int)
	Set(idx int, val interface{})
	Push(val interface{})
}
