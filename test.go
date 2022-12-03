package main

type reader struct {
	data []byte
}

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
