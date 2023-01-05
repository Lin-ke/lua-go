package api

// For apis are exposed to Users, VM needs to hide binchunk details.
// index here keep same with state. (from 1 to top)
type LuaVM interface {
	LuaState
	PC() int
	AddPC(n int)
	Fetch() uint32
	GetConst(idx int)
	LastInst() uint32

	Get(idx int) interface{}
	GetRK(rk int, k int)
	Set(idx int, val interface{})
	Push(val interface{})

	RegisterCount() int
	LoadVararg(n int)
	LoadProto(idx int)

	CloseUpvalues(a int)

	RawArith(op ArithOp) bool // 执行成功
	CallMetaMethod(mmIdx int)

	NewTbcUpval(idx int)
	CloseTbc(idx int)
}
