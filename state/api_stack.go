package state

// [-0, +0, –]
// http://www.lua.org/manual/5.4/manual.html#lua_gettop
func (L *luaState) GetTop() int {
	return L.stack.top
}

// [-0, +0, –]
// http://www.lua.org/manual/5.4/manual.html#lua_absindex
func (L *luaState) AbsIndex(idx int) int {
	return L.stack.absIndex(idx)
}

// [-0, +0, –]
// http://www.lua.org/manual/5.6/manual.html#lua_checkstack
func (L *luaState) CheckStack(n int) bool {
	L.stack.check(n)
	return true // never fails
}

// [-n, +0, –]
// http://www.lua.org/manual/5.4/manual.html#lua_pop
func (L *luaState) Pop(n int) {
	for i := 0; i < n; i++ {
		L.stack.pop()
	}
}

// [-0, +0, –]
// http://www.lua.org/manual/5.4/manual.html#lua_copy
func (L *luaState) Copy(fromIdx, toIdx int) {
	val := L.stack.get(fromIdx)
	L.stack.set(toIdx, val)
}

// [-0, +1, –]
// http://www.lua.org/manual/5.4/manual.html#lua_pushvalue
func (L *luaState) PushValue(idx int) {
	val := L.stack.get(idx)
	L.stack.push(val)
}

// [-1, +0, –]
// http://www.lua.org/manual/5.4/manual.html#lua_replace
func (L *luaState) Replace(idx int) {
	val := L.stack.pop()
	L.stack.set(idx, val)
}

// [-1, +1, –]
// http://www.lua.org/manual/5.4/manual.html#lua_insert
func (L *luaState) Insert(idx int) {
	L.Rotate(idx, 1)
}

// [-1, +0, –]
// http://www.lua.org/manual/5.4/manual.html#lua_remove
func (L *luaState) Remove(idx int) {
	L.Rotate(idx, -1)
	L.Pop(1)
}

// [-0, +0, –]
// http://www.lua.org/manual/5.4/manual.html#lua_rotate
func (L *luaState) Rotate(idx, n int) {
	t := L.stack.top - 1           /* end of stack segment being rotated */
	p := L.stack.absIndex(idx) - 1 /* start of segment */
	var m int                      /* end of prefix */
	if n >= 0 {
		m = t - n
	} else {
		m = p - n - 1
	}
	L.stack.reverse(p, m)   /* reverse the prefix with length 'n' */
	L.stack.reverse(m+1, t) /* reverse the suffix */
	L.stack.reverse(p, t)   /* reverse the entire segment */
}

// [-?, +?, –]
// http://www.lua.org/manual/5.4/manual.html#lua_settop
func (L *luaState) SetTop(idx int) {
	newTop := L.stack.absIndex(idx)
	if newTop < 0 {
		panic("stack underflow!")
	}

	n := L.stack.top - newTop
	if n > 0 {
		for i := 0; i < n; i++ {
			L.stack.pop()
		}
	} else if n < 0 {
		for i := 0; i > n; i-- {
			L.stack.push(nil)
		}
	}
}
