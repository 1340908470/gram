package main

import (
	"gram/base"
	"gram/ll1"
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

	if base.GetMethod() == "LL" {
		/* 运算LL(1) */
		// 消除左递归
		base.RemoveLeftRecursion()

		// 打印中间用到的表
		base.PrintProductions()
		base.PrintFirst()
		base.PrintFollow()

		// 获取LL分析表
		table := ll1.GenerateLLTable()
		err = ll1.PrintLLTable(table)
		if err != nil {
			panic(err)
		}
		// LL分析
		err = ll1.LLAnalyze(base.GetInput(), table)
		if err != nil {
			panic(err)
		}
		ll1.PrintProcedure()
	} else if base.GetMethod() == "LR" {
		/* 运算LR(1) */
		// 构造拓广文法
		base.GenerateExtension()

		// 打印中间用到的表
		base.PrintProductions()
		base.PrintFirst()
		base.PrintFollow()

		// 构造项目集规范族
		lr1.GenerateFamily()
		lr1.PrintFamily()

		table := lr1.GenerateLRTable()
		err = lr1.PrintLRTable(table)
		if err != nil {
			panic(err)
		}
		err = lr1.LRAnalyze(base.GetInput(), table)
		if err != nil {
			panic(err)
		}
		lr1.PrintProcedure()
	}

	if err != nil {
		panic(err)
	}
}
