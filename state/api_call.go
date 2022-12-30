package state

import (
	"fmt"
	"luago54/api"
	"luago54/binchunk"
	"luago54/vm"
)

// [-0, +1, –]
// http://www.lua.org/manual/5.4/manual.html#lua_load
func (L *luaState) Load(chunk []byte, chunkName, mode string) int {
	proto := binchunk.Undump(chunk) // use undump's reader
	c := newLuaClosure(proto)
	L.stack.push(c)
	if len(proto.Upvalues) > 0 {
		env := L.registry.get(api.LUA_RIDX_GLOBALS) // all references.
		c.upvals[0] = &upvalue{&env}
	}
	return 0
}

// [-(nargs+1), +nresults, e]
// http://www.lua.org/manual/5.4/manual.html#lua_call
// rely on __call metamethod.
func (L *luaState) Call(nArgs, nResults int) {
	val := L.stack.get(-(nArgs + 1))
	if c, ok := val.(*closure); ok {
		if c.proto != nil {
			if DEBUG.printCall {
				fmt.Printf("call %s<%d,%d>\n", c.proto.Source,
					c.proto.LineDefined, c.proto.LastLineDefined)
			}
			L.callLuaClosure(nArgs, nResults, c)
		} else {
			L.callGoClosure(nArgs, nResults, c)
		}

	} else {
		panic("not function!")
	}
}

func (L *luaState) callLuaClosure(nArgs, nResults int, c *closure) {
	nRegs := int(c.proto.MaxStackSize)
	nParams := int(c.proto.NumParams)
	isVararg := c.proto.IsVararg == 1

	// create new lua stack
	newStack := newLuaStack(nRegs+api.LUA_MINSTACK, L)
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
		// run will keep the returns on the top.
		results := newStack.popN(newStack.top - nRegs)
		if DEBUG.printReturn {
			fmt.Printf("return : ")
			for _, k := range results {
				printLuaval(k)
			}
			println()
		}

		L.stack.check(len(results))
		L.stack.pushN(results, nResults)
	}
}
func (L *luaState) tailCallLuaClosure(nArgs int, c *closure) {
	nRegs := int(c.proto.MaxStackSize)
	nParams := int(c.proto.NumParams)
	isVararg := c.proto.IsVararg == 1

	// store args
	args := L.stack.popN(nArgs)
	// clean the stack
	L.SetTop(0)
	// check if stack space is enough
	L.stack.check(nRegs + api.LUA_MINSTACK)
	// substitue the closure to new one
	L.stack.closure = c
	L.stack.pc = 0
	// push fixed args
	L.stack.pushN(args, nParams)
	L.stack.top = nRegs

	// store varargs
	if nArgs > nParams && isVararg {
		L.stack.varargs = args[nParams:]
	}

	// run closure
	L.runLuaClosure()
}
func (L *luaState) runLuaClosure() {
	for {
		inst := vm.Instruction(L.Fetch())
		inst.Execute(L)
		if DEBUG.printInst {
			fmt.Printf("[%02d] %s ", L.stack.pc+1, inst.OpName())

		}
		if DEBUG.printStack {
			printStack(L.stack)
			println()
		}

		if inst.Opcode() == vm.OP_RETURN || inst.Opcode() == vm.OP_TAILCALL || inst.Opcode() == vm.OP_RETURN0 || inst.Opcode() == vm.OP_RETURN1 {
			break
		}
	}
}

func (L *luaState) callGoClosure(nArgs, nResults int, c *closure) {
	// create new lua stack
	newStack := newLuaStack(nArgs+api.LUA_MINSTACK, L)
	newStack.closure = c

	// pass args, pop func
	if nArgs > 0 {
		args := L.stack.popN(nArgs)
		newStack.pushN(args, nArgs)
	}
	L.stack.pop()

	// run closure
	L.pushLuaStack(newStack)
	r := c.goFunc(L)
	L.popLuaStack()

	// return results
	if nResults != 0 {
		results := newStack.popN(r)
		L.stack.check(len(results))
		L.stack.pushN(results, nResults)
	}
}
