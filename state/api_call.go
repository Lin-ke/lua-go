package state

import (
	"fmt"
	"luago54/binchunk"
	"luago54/vm"
)

// [-0, +1, –]
// http://www.lua.org/manual/5.4/manual.html#lua_load
func (L *luaState) Load(chunk []byte, chunkName, mode string) int {
	proto := binchunk.Undump(chunk) // use undump's reader
	c := newLuaClosure(proto)
	L.stack.push(c)
	return 0
}

// [-(nargs+1), +nresults, e]
// http://www.lua.org/manual/5.4/manual.html#lua_call
// rely on __call metamethod.
func (L *luaState) Call(nArgs, nResults int) {
	val := L.stack.get(-(nArgs + 1))
	if c, ok := val.(*closure); ok {
		fmt.Printf("call %s<%d,%d>\n", c.proto.Source,
			c.proto.LineDefined, c.proto.LastLineDefined)
		L.callLuaClosure(nArgs, nResults, c)
	} else {
		panic("not function!")
	}
}

func (L *luaState) callLuaClosure(nArgs, nResults int, c *closure) {
	nRegs := int(c.proto.MaxStackSize)
	nParams := int(c.proto.NumParams)
	isVararg := c.proto.IsVararg == 1

	// create new lua stack
	newStack := newLuaStack(nRegs + 20)
	newStack.closure = c

	// pass args, pop func
	funcAndArgs := L.stack.popN(nArgs + 1)
	newStack.pushN(funcAndArgs[1:], nParams)
	newStack.top = nRegs
	//var args.
	if nArgs > nParams && isVararg {
		newStack.varargs = funcAndArgs[nParams+1:]
	}

	// run closure
	L.pushLuaStack(newStack)
	L.runLuaClosure()
	L.popLuaStack()

	// return results
	if nResults != 0 {
		results := newStack.popN(newStack.top - nRegs)
		L.stack.check(len(results))
		L.stack.pushN(results, nResults)
	}
}

func (L *luaState) runLuaClosure() {
	for {
		inst := vm.Instruction(L.Fetch())
		inst.Execute(L)
		if inst.Opcode() == vm.OP_RETURN {
			break
		}
	}
}