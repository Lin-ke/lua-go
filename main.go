package main

import (
	"luago54/test"
)

func main() {
	test.Test007()

}

//old tests :
// var a interface{}
// a = math.MaxFloat64
// b, _ := a.(float64)
// print((int64)(b))
// var a int64 = -1
// var b int64 = -1
// print((uint64)(a + b))
// print((uint64)(a) - uint64(b))

// a, f := number.FloatToInteger(1e100)
// print(a, f)

// const MAXARG_J = 1<<25 - 1      //33554431
// const MAXARG_sJ = MAXARG_J >> 1 //16777215
