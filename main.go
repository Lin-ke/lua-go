package main

// import (
//
//	"math"
//
// )
func ShiftRight(a, n int64) int64 {

	return int64(uint64(a) >> uint64(n))

}
func main() {
	print(ShiftRight(-1, 63))
}
