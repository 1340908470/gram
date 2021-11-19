package ll1

import "gram/base"

// TagStack Tag栈，并在最开始放上$和开始符
var TagStack = []base.Tag{
	{Type: base.TERM, Value: "$"},
	base.GetProductions()[0].Left,
}

// LLAnalyze LL分析程序
func LLAnalyze(input string) {
	// 当且仅当输入队列指针
	for {

	}
}
