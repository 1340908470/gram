package main

import (
	"gram/base"
)

func main() {
	err := base.InitDef()
	if err != nil {
		panic(err)
	}

	base.GetFirst(base.GetTags()[8])
}
