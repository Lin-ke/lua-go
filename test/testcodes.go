package test

import (
	"io/ioutil"
	"luago54/binchunk"
	"luago54/state"
	vm "luago54/vm"
	"os"
)

func luaMain(proto *binchunk.Prototype) {
	nRegs := int(proto.MaxStackSize)
	ls := state.New(nRegs+8, proto)
	ls.SetTop(nRegs)
	for {
		//pc := ls.PC()
		inst := vm.Instruction(ls.Fetch())
		if inst.Opcode() != vm.OP_RETURN {
			inst.Execute(ls)

			// fmt.Printf("[%02d] %s ", pc+1, inst.OpName())
			// printStack(ls)
		} else {
			break
		}
	}
}

func Test006() {
	if len(os.Args) > 1 {
		data, err := ioutil.ReadFile(os.Args[1])
		if err != nil {
			panic(err)
		}

		proto := binchunk.Undump(data)
		luaMain(proto)
	}
}
