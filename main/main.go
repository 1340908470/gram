package main

import (
	"gram/base"
	"gram/ll1"
)

func main() {
	// 初始化定义
	err := base.InitDef()
	if err != nil {
		panic(err)
	}

	/* 运算LL(1) */
	// 消除左递归
	base.RemoveLeftRecursion()

	// 获取LL分析表
	err = ll1.PrintLLTable(ll1.GenerateLLTable())
	if err != nil {
		panic(err)
	}
}
