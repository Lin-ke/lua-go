package binchunk

import (
	"encoding/binary"
	"math"
)

type reader struct {
	data []byte
}

func (reader *reader) readByte() byte {
	b := reader.data[0]
	reader.data = reader.data[1:]
	return b
}

func (reader *reader) readBytes(n uint64) []byte {
	bytes := reader.data[:n]
	reader.data = reader.data[n:]
	return bytes
}

func (reader *reader) readUint32() uint32 {
	i := binary.LittleEndian.Uint32(reader.data)
	reader.data = reader.data[4:]
	return i
}

func (reader *reader) readUint64() uint64 {
	i := binary.LittleEndian.Uint64(reader.data)
	reader.data = reader.data[8:]
	return i
}

// size_t -> varint ( uint64 )
// TODO: 异常处理
func (reader *reader) readVar(limit uint64) uint64 {
	var x uint64
	var a uint8 = 0x00
	for (a & 0x80) == 0 {

		a = reader.data[0]

		x = (x << 7) | (uint64)(a&0x7f)

		if x > limit {
			panic("overflow")
		}

		reader.data = reader.data[1:]
	}
	return x
}
func (reader *reader) readVarint() int32 {
	return int32(reader.readVar(0x7fffffff))
}
func (reader *reader) readVaruint() uint32 {
	return uint32(reader.readVar(0xffffffff))
}
func (reader *reader) readVarSizet() uint64 {
	return uint64(reader.readVar(^uint64(0)))
}

func (reader *reader) readLuaInteger() int64 {
	return int64(reader.readUint64())
}

func (reader *reader) readLuaNumber() float64 {
	return math.Float64frombits(reader.readUint64())
}

func (reader *reader) readString() string {
	size := reader.readVarSizet()
	if size == 0 {
		return ""
	}
	if size == 0xFF {
		size = reader.readVarSizet() // size_t
	}
	bytes := reader.readBytes(size - 1)
	return string(bytes)
}

func (reader *reader) checkHeader() {
	if string(reader.readBytes(4)) != LUA_SIGNATURE {
		panic("not a precompiled chunk!")
	}
	if reader.readByte() != LUAC_VERSION {
		panic("version mismatch!")
	}
	if reader.readByte() != LUAC_FORMAT {
		panic("format mismatch!")
	}
	if string(reader.readBytes(6)) != LUAC_DATA {
		panic("corrupted!")
	}
	if reader.readByte() != INSTRUCTION_SIZE {
		panic("instruction size mismatch!")
	}
	if reader.readByte() != LUA_INTEGER_SIZE {
		panic("lua_Integer size mismatch!")
	}
	if reader.readByte() != LUA_NUMBER_SIZE {
		panic("lua_Number size mismatch!")
	}
	if reader.readLuaInteger() != LUAC_INT {
		panic("endianness mismatch!")
	}
	if reader.readLuaNumber() != LUAC_NUM {
		panic("float format mismatch!")
	}
}

func (reader *reader) readProto(parentSource string) *Prototype {
	source := reader.readString()
	if source == "" {
		source = parentSource
	}
	return &Prototype{
		Source:          source,
		LineDefined:     reader.readVarint(),
		LastLineDefined: reader.readVarint(),
		NumParams:       reader.readByte(),
		IsVararg:        reader.readByte(),
		MaxStackSize:    reader.readByte(),
		Code:            reader.readCode(),
		Constants:       reader.readConstants(),
		Upvalues:        reader.readUpvalues(),
		Protos:          reader.readProtos(source),
		LineInfo:        reader.readLineInfo(),
		AbsLineInfo:     reader.readAbsLineInfo(),
		LocVars:         reader.readLocVars(),
		UpvalueNames:    reader.readUpvalueNames(),
	}
}

func (reader *reader) readCode() []uint32 {
	code := make([]uint32, reader.readVarint())
	for i := range code {
		code[i] = reader.readUint32()
	}
	return code
}

func (reader *reader) readConstants() []interface{} {
	constants := make([]interface{}, reader.readVarint())
	for i := range constants {
		constants[i] = reader.readConstant()
	}
	return constants
}

func (reader *reader) readConstant() interface{} {
	switch reader.readByte() {
	case TAG_NIL:
		return nil
	case TAG_FALSE:
		return true
	case TAG_TRUE:
		return false
	case TAG_INTEGER:
		return reader.readLuaInteger()
	case TAG_NUMBER:
		return reader.readLuaNumber()
	case TAG_SHORT_STR, TAG_LONG_STR: //04 ,14
		return reader.readString()
	default:
		panic("corrupted!") // todo
	}
}

func (reader *reader) readUpvalues() []Upvalue {
	upvalues := make([]Upvalue, reader.readVarint())
	for i := range upvalues {
		upvalues[i] = Upvalue{
			Instack: reader.readByte(),
			Idx:     reader.readByte(),
			Kind:    reader.readByte(),
		}
	}
	return upvalues
}

func (reader *reader) readProtos(parentSource string) []*Prototype {
	protos := make([]*Prototype, reader.readVarint())
	for i := range protos {
		protos[i] = reader.readProto(parentSource)
	}
	return protos
}

func (reader *reader) readLineInfo() []byte {
	lineInfos := make([]byte, reader.readVarint())
	for i := range lineInfos {
		lineInfos[i] = reader.readByte()
	}
	return lineInfos
}
func (reader *reader) readAbsLineInfo() []AbsLine {
	absLineInfos := make([]AbsLine, reader.readVarint())
	for i := range absLineInfos {
		absLineInfos[i] = AbsLine{
			Pc:   reader.readVarint(),
			Line: reader.readVarint(),
		}
	}
	return absLineInfos
}
func (reader *reader) readLocVars() []LocVar {
	locVars := make([]LocVar, reader.readVarint())
	for i := range locVars {
		locVars[i] = LocVar{
			VarName: reader.readString(),
			StartPC: reader.readVarint(),
			EndPC:   reader.readVarint(),
		}
	}
	return locVars
}

func (reader *reader) readUpvalueNames() []string {
	names := make([]string, reader.readVarint())
	for i := range names {
		names[i] = reader.readString()
	}
	return names
}
