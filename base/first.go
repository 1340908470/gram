package base

import (
	"fmt"
	"reflect"
)

// GetFirst 根据推导式的左部，得到其对应的FIRST集
func GetFirst(left Tag) []Tag {
	var ansTags []Tag

	// 如果传入了终结符，则直接返回本身（终结符的First为本身）
	if left.Type == TERM {
		return []Tag{left}
	}

	getFirstRE(left, []Tag{}, &ansTags)
	return ansTags
}

// getFirstRE 递归查找First集，并将该次调用得到的tag加到ansTags中
func getFirstRE(symbol Tag, tmpTags []Tag, ansTags *[]Tag) {
	// 如果出现了与之前tag相同的tag，则说明进入了循环，直接跳出
	for _, production := range GetProdMap()[symbol] {
		// 如果是终结符，则加入ansTag
		if production.Right[0].Type == TERM {
			if !HasReTags(production.Right[0], *ansTags) {
				*ansTags = append(*ansTags, production.Right[0])
			}
		}
		// 如果是非终结符，则 若tmpTags里已经有了则结束，如果没有，则递归查询
		if production.Right[0].Type == NONTERM {
			if !HasReTags(production.Right[0], tmpTags) {
				tmpTags = append(tmpTags, production.Right[0])
				getFirstRE(production.Right[0], tmpTags, ansTags)
			}
		}
	}
}

// HasReTags 是否有重复的tag
func HasReTags(tag Tag, tags []Tag) bool {
	tmp := 0
	for _, t := range tags {
		if reflect.DeepEqual(t, tag) {
			tmp++
			break
		}
	}
	if tmp == 0 {
		return false
	} else {
		return true
	}
}

func PrintFirst() {
	fmt.Printf("--------- 打印FIRST集 ---------\n")
	for _, tag := range GetTags() {
		if tag.Type == NONTERM {
			fmt.Printf("%v: ", tag.Value)
			for _, t := range GetFirst(tag) {
				fmt.Printf("%v ", t.Value)
			}
			fmt.Printf("\n")
		}
	}
	fmt.Printf("-----------------------------\n\n")
}
