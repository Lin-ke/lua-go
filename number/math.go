package number

import "math"

// lua official : lvm.c
// Integer : x // 0 and x mod 0 panic
// Float : x(not 0) / 0 and x mod 0 infinity ,and 0 / 0 and 0 mod 0 Nan (-nan)
// Same as Golang.

// todo: correct?
func FloatToInteger(f float64) (int64, bool) {
	i := int64(f)
	return i, float64(i) == f
}

// a % b == a - ((a // b) * b)
func IMod(a, b int64) int64 {
	if b == 0 {
		panic("attempt to perform 'n%0'")
	}
	return a - IFloorDiv(a, b)*b
}

// lvm.c l756
func FMod(a, b float64) float64 {

	return math.Mod(a, b)
}

// 1//0 panic
func IFloorDiv(a, b int64) int64 {
	if b == 0 {
		panic("attempt to perform 'n%0'")
	}
	if a > 0 && b > 0 || a < 0 && b < 0 || a%b == 0 {
		return a / b
	} else {
		return a/b - 1
	}
}

func FFloorDiv(a, b float64) float64 {
	return math.Floor(a / b)
}

func ShiftLeft(a, n int64) int64 {
	if n >= 0 {
		return a << uint64(n)
	} else {
		return ShiftRight(a, -n)
	}
}

func ShiftRight(a, n int64) int64 {
	if n >= 0 {
		return int64(uint64(a) >> uint64(n))
	} else {
		return ShiftLeft(a, -n)
	}
}

// logic ：如果值可以转换为
//false
//，那么 整 个表达式的结果就是该值
