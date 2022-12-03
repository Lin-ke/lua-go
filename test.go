package main

type reader struct {
	data []byte
}

func (self *reader) readVarint() uint32 {
	var x uint32
	var a uint8 = 0x00
	for (a & 0x80) == 0 {
		a = self.data[0]
		x = (x << 7) | (uint32)(a&0x7f)
		self.data = self.data[1:]
	}

	return x
}
func main() {
	reader := &reader{[]byte{0x11, 0x11, 0x11, 0x11, 0x11, 0x11, 0x80}}
	print(reader.readVarint())
}
