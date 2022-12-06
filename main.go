package main

// import (
// 	"math"
// )

const MaxFloat64 = 0x1p1023 * (1 + (1 - 0x1p-52)) // 1.79769313486231570814527423731704356798070e+308
func IFloorDiv(a, b float64) float64 {
	return a / b
}
func main() {
	print(IFloorDiv(0, 0))
}
