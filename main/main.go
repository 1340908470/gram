package main

import (
	"gram/def"
)

func main() {
	err := def.InitDef()
	if err != nil {
		panic(err)
	}

	def.GetTags()
	def.GetProductions()
}
