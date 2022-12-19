package main

const MAXARG_J = 1<<25 - 1      //33554431
const MAXARG_sJ = MAXARG_J >> 1 //16777215

func main() {
	// var a interface{}
	// a = math.MaxFloat64
	// b, _ := a.(float64)
	// print((int64)(b))
	var a int64 = -1
	var b int64 = -1
	print((uint64)(a + b))
	print((uint64)(a) - uint64(b))
}
func SJ(i uint32) int {
	return int(i>>7) - MAXARG_sJ
}
