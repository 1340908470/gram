package base

const (
	TERM    = iota // 终结
	NONTERM        // 非终结
)

// Tag 标识符，包括终结符和非终结符
type Tag struct {
	Type  int    // 类型
	Value string // 值
}
