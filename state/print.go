package state

import (
	"fmt"
	"luago54/vm"
	"runtime"
)

type DebugInfo struct {
	printUpval, printStack, printTable, printInst, printInstDetail bool
	printSetSlot, printReturn, printCall, printFunc                bool
}

var DEBUG DebugInfo = DebugInfo{
	true, true, false, true, false,
	true, true, true, true,
}

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
	case *closure:
		fmt.Printf("[function")
		if DEBUG.printFunc {
			if v.proto != nil {
				fmt.Printf("<%d,%d>", v.proto.LineDefined, v.proto.LastLineDefined)
			}
		}
		fmt.Printf("]")
	case *luaTable:
		fmt.Printf("[table]")
		if DEBUG.printTable {
			printTable(v)
		}

	default: // other values
		fmt.Printf("[%s]", v)
	}
}

func printTable(Table *luaTable) {
	if Table.metatable != nil {
		fmt.Print("[metatable]")
	}
	fmt.Printf("\narray:%d\n", Table.len())
	for _, itr := range Table.arr {
		printLuaval(itr)
	}
	fmt.Println()
	fmt.Printf("map:%d\n", len(Table._map))
	for k, v := range Table._map {
		printLuaval(k)
		fmt.Printf(":")
		printLuaval(v)
		fmt.Printf(" | ")
	}
	fmt.Println()
	fmt.Println()
}
func printInst(i vm.Instruction) {
	_, k, b, c := i.ABC()
	_, sbx := i.AsBx()
	a, bx := i.ABx()
	sj := i.SJ()
	sc := (int)(c)
	sb := (int)(b)
	ax := i.Ax()
	isk := ""
	if k != 0 {
		isk = "k"
	}

	switch i.OpMode() {
	case vm.OP_MOVE:
		fmt.Printf("%d %d", a, b)

	case vm.OP_LOADI:
		fmt.Printf("%d %d", a, sbx)

	case vm.OP_LOADF:
		fmt.Printf("%d %d", a, sbx)

	case vm.OP_LOADK:
		fmt.Printf("%d %d", a, bx)

	case vm.OP_LOADKX:
		fmt.Printf("%d", a)

	case vm.OP_LOADFALSE:
		fmt.Printf("%d", a)

	case vm.OP_LFALSESKIP:
		fmt.Printf("%d", a)

	case vm.OP_LOADTRUE:
		fmt.Printf("%d", a)

	case vm.OP_LOADNIL:
		fmt.Printf("%d %d", a, b)

	case vm.OP_GETUPVAL:
		fmt.Printf("%d %d", a, b)

	case vm.OP_SETUPVAL:
		fmt.Printf("%d %d", a, b)

	case vm.OP_GETTABUP:
		fmt.Printf("%d %d %d", a, b, c)

	case vm.OP_GETTABLE:
		fmt.Printf("%d %d %d", a, b, c)

	case vm.OP_GETI:
		fmt.Printf("%d %d %d", a, b, c)

	case vm.OP_GETFIELD:
		fmt.Printf("%d %d %d", a, b, c)

	case vm.OP_SETTABUP:
		fmt.Printf("%d %d %d%s", a, b, c, isk)

	case vm.OP_SETTABLE:
		fmt.Printf("%d %d %d%s", a, b, c, isk)

	case vm.OP_SETI:
		fmt.Printf("%d %d %d%s", a, b, c, isk)

	case vm.OP_SETFIELD:
		fmt.Printf("%d %d %d%s", a, b, c, isk)

	case vm.OP_NEWTABLE:
		fmt.Printf("%d %d %d", a, b, c)

	case vm.OP_SELF:
		fmt.Printf("%d %d %d%s", a, b, c, isk)

	case vm.OP_ADDI:
		fmt.Printf("%d %d %d", a, b, sc)

	case vm.OP_ADDK:
		fmt.Printf("%d %d %d", a, b, c)
		fmt.Println()

	case vm.OP_SUBK:
		fmt.Printf("%d %d %d", a, b, c)
		fmt.Println()

	case vm.OP_MULK:
		fmt.Printf("%d %d %d", a, b, c)
		fmt.Println()

	case vm.OP_MODK:
		fmt.Printf("%d %d %d", a, b, c)
		fmt.Println()

	case vm.OP_POWK:
		fmt.Printf("%d %d %d", a, b, c)
		fmt.Println()

	case vm.OP_DIVK:
		fmt.Printf("%d %d %d", a, b, c)
		fmt.Println()

	case vm.OP_IDIVK:
		fmt.Printf("%d %d %d", a, b, c)
		fmt.Println()

	case vm.OP_BANDK:
		fmt.Printf("%d %d %d", a, b, c)
		fmt.Println()

	case vm.OP_BORK:
		fmt.Printf("%d %d %d", a, b, c)
		fmt.Println()

	case vm.OP_BXORK:
		fmt.Printf("%d %d %d", a, b, c)
		fmt.Println()

	case vm.OP_SHRI:
		fmt.Printf("%d %d %d", a, b, sc)

	case vm.OP_SHLI:
		fmt.Printf("%d %d %d", a, b, sc)

	case vm.OP_ADD:
		fmt.Printf("%d %d %d", a, b, c)

	case vm.OP_SUB:
		fmt.Printf("%d %d %d", a, b, c)

	case vm.OP_MUL:
		fmt.Printf("%d %d %d", a, b, c)

	case vm.OP_MOD:
		fmt.Printf("%d %d %d", a, b, c)

	case vm.OP_POW:
		fmt.Printf("%d %d %d", a, b, c)

	case vm.OP_DIV:
		fmt.Printf("%d %d %d", a, b, c)

	case vm.OP_IDIV:
		fmt.Printf("%d %d %d", a, b, c)

	case vm.OP_BAND:
		fmt.Printf("%d %d %d", a, b, c)

	case vm.OP_BOR:
		fmt.Printf("%d %d %d", a, b, c)

	case vm.OP_BXOR:
		fmt.Printf("%d %d %d", a, b, c)

	case vm.OP_SHL:
		fmt.Printf("%d %d %d", a, b, c)

	case vm.OP_SHR:
		fmt.Printf("%d %d %d", a, b, c)

	case vm.OP_MMBIN:
		fmt.Printf("%d %d %d", a, b, c)
		fmt.Println()

	case vm.OP_MMBINI:
		fmt.Printf("%d %d %d %d", a, sb, c, k)
		fmt.Println()

	case vm.OP_MMBINK:
		fmt.Printf("%d %d %d %d", a, b, c, k)

	case vm.OP_UNM:
		fmt.Printf("%d %d", a, b)

	case vm.OP_BNOT:
		fmt.Printf("%d %d", a, b)

	case vm.OP_NOT:
		fmt.Printf("%d %d", a, b)

	case vm.OP_LEN:
		fmt.Printf("%d %d", a, b)

	case vm.OP_CONCAT:
		fmt.Printf("%d %d", a, b)

	case vm.OP_CLOSE:
		fmt.Printf("%d", a)

	case vm.OP_TBC:
		fmt.Printf("%d", a)

	case vm.OP_JMP:
		fmt.Printf("%d", sj)
		fmt.Println()

	case vm.OP_EQ:
		fmt.Printf("%d %d %d", a, b, k)

	case vm.OP_LT:
		fmt.Printf("%d %d %d", a, b, k)

	case vm.OP_LE:
		fmt.Printf("%d %d %d", a, b, k)

	case vm.OP_EQK:
		fmt.Printf("%d %d %d", a, b, k)
		fmt.Println()

	case vm.OP_EQI:
		fmt.Printf("%d %d %d", a, sb, k)

	case vm.OP_LTI:
		fmt.Printf("%d %d %d", a, sb, k)

	case vm.OP_LEI:
		fmt.Printf("%d %d %d", a, sb, k)

	case vm.OP_GTI:
		fmt.Printf("%d %d %d", a, sb, k)

	case vm.OP_GEI:
		fmt.Printf("%d %d %d", a, sb, k)

	case vm.OP_TEST:
		fmt.Printf("%d %d", a, k)

	case vm.OP_TESTSET:
		fmt.Printf("%d %d %d", a, b, k)

	case vm.OP_CALL:
		fmt.Printf("%d %d %d", a, b, c)
		fmt.Println()
		if b == 0 {
			fmt.Printf("all in ")
		} else {
			fmt.Printf("%d in ", b-1)
		}
		if c == 0 {
			fmt.Printf("all out")
		} else {
			fmt.Printf("%d out", c-1)
		}

	case vm.OP_TAILCALL:
		fmt.Printf("%d %d %d%s", a, b, c, isk)
		fmt.Println()

	case vm.OP_RETURN:
		fmt.Printf("%d %d %d%s", a, b, c, isk)
		fmt.Println()
		if b == 0 {
			fmt.Printf("all out")
		} else {
			fmt.Printf("%d out", b-1)
		}

	case vm.OP_RETURN0:

	case vm.OP_RETURN1:
		fmt.Printf("%d", a)

	case vm.OP_FORLOOP:
		fmt.Printf("%d %d", a, bx)
		fmt.Println()

	case vm.OP_FORPREP:
		fmt.Printf("%d %d", a, bx)
		fmt.Println()

	case vm.OP_TFORPREP:
		fmt.Printf("%d %d", a, bx)
		fmt.Println()

	case vm.OP_TFORCALL:
		fmt.Printf("%d %d", a, c)

	case vm.OP_TFORLOOP:
		fmt.Printf("%d %d", a, bx)
		fmt.Println()

	case vm.OP_SETLIST:
		fmt.Printf("%d %d %d", a, b, c)

	case vm.OP_CLOSURE:
		fmt.Printf("%d %d", a, bx)
		fmt.Println()

	case vm.OP_VARARG:
		fmt.Printf("%d %d", a, c)
		fmt.Println()
		if c == 0 {
			fmt.Printf("all out")
		} else {
			fmt.Printf("%d out", c-1)
		}

	case vm.OP_VARARGPREP:
		fmt.Printf("%d", a)

	case vm.OP_EXTRAARG:
		fmt.Printf("%d", ax)

	}
}
func runFuncName() string {
	pc := make([]uintptr, 1)
	runtime.Callers(2, pc)
	f := runtime.FuncForPC(pc[0])
	return f.Name()
}
