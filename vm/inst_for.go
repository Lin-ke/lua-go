package vm

import (
	"luago54/api"
	"math"
)

// cut
func _float2int(a float64) (int64, bool) {
	if a > math.MaxInt64 {
		return math.MaxInt64, false
	} else if a < math.MinInt64 {
		return math.MinInt64, false
	}
	return (int64)(a), true
}

// int可以转，float不能转
// make sure that floatvalue can be used.
func _toNumberType(a interface{}) (int64, float64, bool) {
	fa, fi := a.(float64)
	ia, ii := a.(int64)
	if fi {

		return 0, fa, true
	} else if ii {
		return ia, float64(ia), false
	}
	panic("err")
}

//lvm.c #197
/*
** Prepare a numerical for loop (opcode OP_FORPREP).
** Return true to skip the loop. Otherwise,
** after preparation, stack will be as follows:
**   ra : internal index (safe copy of the control variable)
**   ra + 1 : loop counter (integer loops) or limit (float loops)
**   ra + 2 : step
**   ra + 3 : control variable
 */
// panic will replaced by print errors
// prior to Number ,if step/init is string.
func forPrep(i Instruction, vm api.LuaVM) {
	a, Bx := i.ABx()
	a += 1
	if _forPrep(a, vm) {
		vm.AddPC(Bx + 1)
	}
}
func _forPrep(a int, vm api.LuaVM) bool {

	// type of init and step

	if vm.Type(a) == api.LUA_TSTRING {
		vm.PushNumber(vm.ToNumber(a))
		vm.Replace(a)

	}
	if vm.Type(a+1) == api.LUA_TSTRING {
		vm.PushNumber(vm.ToNumber(a + 1))
		vm.Replace(a + 1)

	}
	if vm.Type(a+2) == api.LUA_TSTRING {
		vm.PushNumber(vm.ToNumber(a + 2))
		vm.Replace(a + 2)

	}
	init := vm.Get(a)
	limit := vm.Get(a + 1)
	step := vm.Get(a + 2)

	istep, fstep, fs := _toNumberType(step)
	iinit, finit, fi := _toNumberType(init)

	if !fs && !fi {
		// int loop
		if istep == 0 {
			panic("'for' step is zero ")
		}
		vm.Copy(a, a+3)
		// modify limit
		fjump, ilimit := forilimit(istep, iinit, limit)

		if fjump {
			return true
		} /* prepare loop */
		var count uint64
		if istep > 0 { /* ascending loop? */
			count = (uint64)(ilimit) - (uint64)(iinit)
			if istep != 1 { /* avoid division in the too common case */
				count = count / (uint64)(istep)
			}
		} else {
			count = (uint64)(iinit) - (uint64)(ilimit)
			count /= (uint64)(-istep-1) + (uint64)(1)
		}

		vm.Set(a+1, (int64)(count))

	} else {
		/*
			In official implementation, nan will not throw a error instantly.
			Loop will carry out once.
		*/
		_, flimit, _ := _toNumberType(limit)
		if fstep == 0 {
			panic("'for' step is zero ")
		}

		// float loop
		if (fstep < 0 && finit < flimit) || (fstep > 0 && finit > flimit) {
			return true
		}
		vm.Set(a+1, flimit)
		vm.Set(a+2, fstep)
		vm.Set(a, finit) /* internal index */
		vm.Copy(a, a+3)  /* control variable */
	}
	return false
}
func forilimit(istep, iinit int64, limit interface{}) (bool, int64) {

	ilimit, flimit, fl := _toNumberType(limit)
	if fl {
		ilimit, cutflag := _float2int(flimit)
		if !cutflag && (istep < 0 && ilimit < 0 || istep > 0 && ilimit > 0) {
			return true, 0

		}
	}
	if (istep < 0 && iinit < ilimit) || (istep > 0 && iinit > ilimit) {
		return true, 0
	} else {
		return false, ilimit
	}
}

func forLoop(i Instruction, vm api.LuaVM) {
	a, Bx := i.ABx()
	a += 1
	if vm.IsInteger(a + 2) { // int loop
		if vm.ToNumber(a+1) > 0 {
			// execute loop
			//count = count - 1
			vm.PushValue(a + 1)
			vm.PushInteger(1)
			vm.Arith(api.LUA_OPSUB)
			vm.Replace(a + 1)
			// i = i + step
			vm.PushValue(a)
			vm.PushValue(a + 2)
			vm.Arith(api.LUA_OPADD)
			vm.Replace(a)
			vm.Copy(a, a+3) //i
			vm.AddPC(-Bx)
		}
	} else {
		if _floatforloop(a, vm) {
			vm.AddPC(-Bx)
		}
		//finish

	}
}

// inf will carry out forever and nan will carry out once.
func _floatforloop(a int, vm api.LuaVM) bool {
	step := vm.ToNumber(a + 2)
	limit := vm.ToNumber(a + 1)
	idx := vm.ToNumber(a)
	idx += step
	// carry out
	if (step < 0 && idx > limit) || (step > 0 && idx < limit) {
		vm.Set(a, idx)
		vm.Copy(a+3, a)
		return true
	}
	return false
}
