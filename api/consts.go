package api

const LUA_MINSTACK = 20
const LUAI_MAXSTACK = 1000000
const LUA_REGISTRYINDEX = -LUAI_MAXSTACK - 1000
const LUA_RIDX_GLOBALS int64 = 2

/* basic types */
const (
	LUA_TNONE = iota - 1 // -1
	LUA_TNIL
	LUA_TBOOLEAN
	LUA_TLIGHTUSERDATA
	LUA_TNUMBER
	LUA_TSTRING
	LUA_TTABLE
	LUA_TFUNCTION
	LUA_TUSERDATA
	LUA_TTHREAD
	LUA_TNUMBERTYPES
)

/* arithmetic functions */
const (
	LUA_OPADD  = iota // +
	LUA_OPSUB         // -
	LUA_OPMUL         // *
	LUA_OPMOD         // %
	LUA_OPPOW         // ^
	LUA_OPDIV         // /
	LUA_OPIDIV        // //
	LUA_OPBAND        // &
	LUA_OPBOR         // |
	LUA_OPBXOR        // ~
	LUA_OPSHL         // <<
	LUA_OPSHR         // >>
	LUA_OPUNM         // -
	LUA_OPBNOT        // ~
)
const (
	LUA_OPEQ = iota
	LUA_OPLT
	LUA_OPLE
)

const (
	LUA_TMINDEX = iota
	LUA_TMNEWINDEX
	LUA_TMGC
	LUA_TMMODE
	LUA_TMLEN
	LUA_TMEQ /* last tag method with fast access */
	LUA_TMADD
	LUA_TMSUB
	LUA_TMMUL
	LUA_TMMOD
	LUA_TMPOW
	LUA_TMDIV
	LUA_TMIDIV
	LUA_TMBAND
	LUA_TMBOR
	LUA_TMBXOR
	LUA_TMSHL
	LUA_TMSHR
	LUA_TMUNM
	LUA_TMBNOT
	LUA_TMLT
	LUA_TMLE
	LUA_TMCONCAT
	LUA_TMCALL
	LUA_TMCLOSE
	LUA_TMN /* number of elements in the enum */
)
