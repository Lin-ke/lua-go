package binchunk

const (
	LUA_SIGNATURE    = "\x1bLua"
	LUAC_VERSION     = 0x53
	LUAC_FORMAT      = 0
	LUAC_DATA        = "\x19\x93\r\n\x1a\n"
	INSTRUCTION_SIZE = 4
	LUA_INTEGER_SIZE = 8
	LUA_NUMBER_SIZE  = 8
	LUAC_INT         = 0x5678
	LUAC_NUM         = 370.5
)

const (
	TAG_NIL       = 0x00
	TAG_FALSE     = 0x01
	TAG_TRUE      = 0x11
	TAG_NUMBER    = 0x03
	TAG_INTEGER   = 0x13
	TAG_SHORT_STR = 0x04
	TAG_LONG_STR  = 0x14
)

type binaryChunk struct {
	header
	sizeUpvalues byte //
	mainFunc     *Prototype
}

// head of luac
type header struct {
	signature [4]byte //magic number
	version   byte    // 5.4 - 0x54
	format    byte    // 00 vm format
	luacData  [6]byte // 0x19 93 0D 0A 1A 0A
	// cintSize        byte    //平台相关，现已被放弃
	// sizetSize       byte    //平台相关，现已被放弃
	instructionSize byte    //04
	luaIntegerSize  byte    //08
	luaNumberSize   byte    //08
	luacInt         int64   //0x5678 (0x78 56 00 00 00 00 00 00)
	luacNum         float64 //370.5
}

// function prototype
// 读入还是int，写的时候转成varint

type Prototype struct {
	Source          string // debug, 源文件名。 short : 0x（name_length+1）+1 @ name;
	LineDefined     uint32 //  start
	LastLineDefined uint32 // end。如果是main则都是0x80
	NumParams       byte
	IsVararg        byte          // 变长参数
	MaxStackSize    byte          // 寄存器
	Code            []uint32      //下面的，前面都有varint 的size
	Constants       []interface{} //空接口模拟联合体。前面有type
	Upvalues        []Upvalue
	Protos          []*Prototype
	LineInfo        []byte    // debug 和指令表中的指令一一对应
	AbsLineInfo     []AbsLine // debug
	LocVars         []LocVar  // debug 局部变量
	UpvalueNames    []string  // debug Upvalue名
}

type Upvalue struct {
	Instack byte
	Idx     byte
	Kind    byte
}

type LocVar struct {
	VarName string
	StartPC uint32
	EndPC   uint32
}

type AbsLine struct {
	Pc   int32
	Line int32
}

func Undump(data []byte) *Prototype {
	reader := &reader{data}
	reader.checkHeader()
	reader.readByte() // size_upvalues
	return reader.readProto("")
}
