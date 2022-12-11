package main

const MAXARG_J = 1<<25 - 1      //33554431
const MAXARG_sJ = MAXARG_J >> 1 //16777215

func main() {
	a := 0x7FFFFF38
	println(a >> 7)
	print(MAXARG_J)
	print(SJ(uint32(a)))
}
func SJ(i uint32) int {
	return int(i>>7) - MAXARG_sJ
}
