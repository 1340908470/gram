package base

import "reflect"

const (
	TERM    = iota // 终结
	NONTERM        // 非终结
)

// Tag 标识符，包括终结符和非终结符
type Tag struct {
	Type  int    // 类型
	Value string // 值
}

func IsEmptyTag(tag Tag) bool {
	return reflect.DeepEqual(tag, Tag{
		Type:  TERM,
		Value: "ε",
	})
}

func ConvertTagsToStr(tags []Tag) string {
	str := ""

	for _, tag := range tags {
		str += tag.Value
	}

	return str
}
