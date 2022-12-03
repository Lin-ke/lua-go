package binchunk

import (
	"encoding/binary"
	"math"
)

type reader struct {
	data []byte
}

func (self *reader) readByte() byte {
	b := self.data[0]
	self.data = self.data[1:]
	return b
}

func (self *reader) readBytes(n uint64) []byte {
	bytes := self.data[:n]
	self.data = self.data[n:]
	return bytes
}

func (self *reader) readUint32() uint32 {
	i := binary.LittleEndian.Uint32(self.data)
	self.data = self.data[4:]
	return i
}

func (self *reader) readUint64() uint64 {
	i := binary.LittleEndian.Uint64(self.data)
	self.data = self.data[8:]
	return i
}

// size_t -> varint ( uint64 )
// TODO: 异常处理
func (self *reader) readVar(limit uint64) uint64 {
	var x uint64
	var a uint8 = 0x00
	for (a & 0x80) == 0 {

		a = self.data[0]

		x = (x << 7) | (uint64)(a&0x7f)

		if x > limit {
			panic("overflow")
		}

		self.data = self.data[1:]
	}
	return x
}
func (self *reader) readVarint() int32 {
	return int32(self.readVar(0x7fffffff))
}
func (self *reader) readVaruint() uint32 {
	return uint32(self.readVar(0xffffffff))
}
func (self *reader) readVarSizet() uint64 {
	return uint64(self.readVar(^uint64(0)))
}

func (self *reader) readLuaInteger() int64 {
	return int64(self.readUint64())
}

func (self *reader) readLuaNumber() float64 {
	return math.Float64frombits(self.readUint64())
}

func (self *reader) readString() string {
	size := self.readVarSizet()
	if size == 0 {
		return ""
	}
	if size == 0xFF {
		size = self.readVarSizet() // size_t
	}
	bytes := self.readBytes(size - 1)
	return string(bytes)
}

func (self *reader) checkHeader() {
	if string(self.readBytes(4)) != LUA_SIGNATURE {
		panic("not a precompiled chunk!")
	}
	if self.readByte() != LUAC_VERSION {
		panic("version mismatch!")
	}
	if self.readByte() != LUAC_FORMAT {
		panic("format mismatch!")
	}
	if string(self.readBytes(6)) != LUAC_DATA {
		panic("corrupted!")
	}
	if self.readByte() != INSTRUCTION_SIZE {
		panic("instruction size mismatch!")
	}
	if self.readByte() != LUA_INTEGER_SIZE {
		panic("lua_Integer size mismatch!")
	}
	if self.readByte() != LUA_NUMBER_SIZE {
		panic("lua_Number size mismatch!")
	}
	if self.readLuaInteger() != LUAC_INT {
		panic("endianness mismatch!")
	}
	if self.readLuaNumber() != LUAC_NUM {
		panic("float format mismatch!")
	}
}

func (self *reader) readProto(parentSource string) *Prototype {
	source := self.readString()
	if source == "" {
		source = parentSource
	}
	return &Prototype{
		Source:          source,
		LineDefined:     self.readVarint(),
		LastLineDefined: self.readVarint(),
		NumParams:       self.readByte(),
		IsVararg:        self.readByte(),
		MaxStackSize:    self.readByte(),
		Code:            self.readCode(),
		Constants:       self.readConstants(),
		Upvalues:        self.readUpvalues(),
		Protos:          self.readProtos(source),
		LineInfo:        self.readLineInfo(),
		AbsLineInfo:     self.readAbsLineInfo(),
		LocVars:         self.readLocVars(),
		UpvalueNames:    self.readUpvalueNames(),
	}
}

func (self *reader) readCode() []uint32 {
	code := make([]uint32, self.readVarint())
	for i := range code {
		code[i] = self.readUint32()
	}
	return code
}

func (self *reader) readConstants() []interface{} {
	constants := make([]interface{}, self.readVarint())
	for i := range constants {
		constants[i] = self.readConstant()
	}
	return constants
}

func (self *reader) readConstant() interface{} {
	switch self.readByte() {
	case TAG_NIL:
		return nil
	case TAG_FALSE:
		return true
	case TAG_TRUE:
		return false
	case TAG_INTEGER:
		return self.readLuaInteger()
	case TAG_NUMBER:
		return self.readLuaNumber()
	case TAG_SHORT_STR, TAG_LONG_STR: //04 ,14
		return self.readString()
	default:
		panic("corrupted!") // todo
	}
}

func (self *reader) readUpvalues() []Upvalue {
	upvalues := make([]Upvalue, self.readVarint())
	for i := range upvalues {
		upvalues[i] = Upvalue{
			Instack: self.readByte(),
			Idx:     self.readByte(),
		}
	}
	return upvalues
}

func (self *reader) readProtos(parentSource string) []*Prototype {
	protos := make([]*Prototype, self.readVarint())
	for i := range protos {
		protos[i] = self.readProto(parentSource)
	}
	return protos
}

func (self *reader) readLineInfo() []byte {
	return nil
}
func (self *reader) readAbsLineInfo() []AbsLine {
	absLineInfos := make([]AbsLine, self.readVarint())
	for i := range absLineInfos {
		absLineInfos[i] = AbsLine{
			Pc:   self.readVarint(),
			Line: self.readVarint(),
		}
	}
	return absLineInfos
}
func (self *reader) readLocVars() []LocVar {
	locVars := make([]LocVar, self.readVarint())
	for i := range locVars {
		locVars[i] = LocVar{
			VarName: self.readString(),
			StartPC: self.readVarint(),
			EndPC:   self.readVarint(),
		}
	}
	return locVars
}

func (self *reader) readUpvalueNames() []string {
	names := make([]string, self.readVarint())
	for i := range names {
		names[i] = self.readString()
	}
	return names
}
