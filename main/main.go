package main

import (
	"gram/base"
	"gram/lr1"
)

func main() {
	// 初始化定义
	err := base.InitDef()
	if err != nil {
		panic(err)
	}

	base.PrintProductions()

	base.PrintFirst()
	base.PrintFollow()

	///* 运算LL(1) */
	//// 消除左递归
	//base.RemoveLeftRecursion()
	//base.PrintProductions()
	//base.PrintFirst()
	//base.PrintFollow()
	//// 获取LL分析表
	//table := ll1.GenerateLLTable()
	//err = ll1.PrintLLTable(table)
	//// LL分析
	//err = ll1.LLAnalyze("num+num$", table)
	//ll1.PrintProcedure()

	/* 运算LR(1) */
	// 构造拓广文法
	base.GenerateExtension()

	// 构造项目集规范族
	lr1.GenerateFamily()
	lr1.PrintFamily()

	table := lr1.GenerateLRTable()
	err = lr1.PrintLRTable(table)

	if err != nil {
		panic(err)
	}
}
