package main

import (
	"gram/base"
)

func main() {
	err := base.InitDef()
	if err != nil {
		panic(err)
	}

	tags := base.GetTags()

	base.GetFirst(tags[8])
	d := base.GetFollow(tags[8], 0)

	print(len(d))
}
