package api

type LuaType = int
type ArithOp = int
type CompareOp int

// idx : 1->top ; -top -> -1
type LuaState interface {
	/* basic stack manipulation */
	GetTop() int
	AbsIndex(idx int) int
	CheckStack(n int) bool
	Pop(n int)
	Copy(fromIdx, toIdx int)
	PushValue(idx int)
	Replace(idx int)
	Insert(idx int)
	Remove(idx int)
	Rotate(idx, n int)
	SetTop(idx int)
	/* access functions (stack -> Go) */
	TypeName(tp LuaType) string
	Type(idx int) LuaType
	IsNone(idx int) bool
	IsNil(idx int) bool
	IsNoneOrNil(idx int) bool
	IsBoolean(idx int) bool
	IsInteger(idx int) bool
	IsNumber(idx int) bool
	IsString(idx int) bool
	IsTable(idx int) bool
	IsThread(idx int) bool
	IsFunction(idx int) bool
	ToBoolean(idx int) bool
	ToInteger(idx int) int64
	ToIntegerX(idx int) (int64, bool)
	ToNumber(idx int) float64
	ToNumberX(idx int) (float64, bool)
	ToString(idx int) string
	ToStringX(idx int) (string, bool)
	/* push functions (Go -> stack) */
	PushNil()
	PushBoolean(b bool)
	PushInteger(n int64)
	PushNumber(n float64)
	PushString(s string)
	/* get functions */
	NewTable()
	CreateTable(nArr, nRec int)
	GetTable(idx int) LuaType
	GetField(idx int, k string) LuaType
	GetI(idx int, i int64) LuaType
	RawGet(idx int) LuaType
	RawGetI(idx int, i int64) LuaType
	GetMetatable(idx int) bool
	GetGlobal(name string) LuaType
	/* set functions */
	SetTable(idx int)
	SetField(idx int, k string)
	SetI(idx int, i int64)
	RawSet(idx int)
	RawSetI(idx int, i int64)
	SetMetatable(idx int)
	SetGlobal(name string)
	/* Comparison and arithmetic functions */
	Arith(op ArithOp)
	Compare(idx1, idx2 int, op CompareOp) bool // 两个索引
	/* miscellaneous functions */
	Len(idx int)
	Concat(n int)
}
