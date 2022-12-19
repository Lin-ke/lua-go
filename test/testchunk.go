package test

import (
	"io/ioutil"
	"luago54/binchunk"
	"os"
)

func Test002() {
	if len(os.Args) > 1 {
		data, err := ioutil.ReadFile(os.Args[1])
		if err != nil {
			panic(err)
		}

		proto := binchunk.Undump(data)
		list(proto)
	}
}
