package main

import (
	"fmt"
	// "io/ioutil"
	// "luago54/binchunk"
	// "os"
	. "luago54/vm"
)

// func main() {
// 	if len(os.Args) > 1 {
// 		data, err := ioutil.ReadFile(os.Args[1])
// 		if err != nil {
// 			panic(err)
// 		}

// 		proto := binchunk.Undump(data)
// 		list(proto)
// 	}
// }

func printOperands(i Instruction) {
	switch i.OpMode() {
	case IABC:
		a, b, c := i.ABC()

		fmt.Printf("%d", a)
		if i.BMode() != 0 {
			if b > 0xFF {
				fmt.Printf(" %d", -1-b&0xFF)
			} else {
				fmt.Printf(" %d", b)
			}
		}
		if i.CMode() != 0 {
			if c > 0xFF {
				fmt.Printf(" %d", -1-c&0xFF)
			} else {
				fmt.Printf(" %d", c)
			}
		}
	case IABx:
		a, bx := i.ABx()
		// TODO: complete BX
		fmt.Printf("%d", a)
		if i.BMode() == 0 {
			fmt.Printf(" %d", -1-bx)
		}
	case IAsBx:
		a, sbx := i.AsBx()
		fmt.Printf("%d %d", a, sbx)
	case IAx:
		ax := i.Ax()
		fmt.Printf("%d", -1-ax)

	case IsJ:
		sJ := i.SJ()
		fmt.Printf("%d", sJ)
	}
}
