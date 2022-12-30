package state

import (
	"luago54/api"
	"luago54/binchunk"
)

type closure struct {
	proto  *binchunk.Prototype
	goFunc api.GoFunction
	upvals []*upvalue
}
type upvalue struct {
	val *luaValue
}

func newLuaClosure(proto *binchunk.Prototype) *closure {
	c := &closure{proto: proto}
	if nUpvals := len(proto.Upvalues); nUpvals > 0 {
		c.upvals = make([]*upvalue, nUpvals)
	}
	return c
}

func newGoClosure(f api.GoFunction, nUpvals int) *closure {
	c := &closure{goFunc: f}
	if nUpvals > 0 {
		c.upvals = make([]*upvalue, nUpvals)
	}
	return c
}
