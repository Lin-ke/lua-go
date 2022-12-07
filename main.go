package main

import (
	"luago54/number"
)

// import (
//
//	"math"
//
// )
const MaxFloat64 = 0x1p1023 * (1 + (1 - 0x1p-52)) // 1.79769313486231570814527423731704356798070e+308
func IFloorDiv(a, b float64) float64 {
	return a / b
}
func lower(c byte) byte {
	return c | ('x' - 'X')
}
func main() {
	print(number.ParseInteger("123"))
}
