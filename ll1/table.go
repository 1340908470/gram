package ll1

import (
	"errors"
	"fmt"
	"github.com/liushuochen/gotable"
	"gram/base"
)

// LLTable LL(1)预测分析表 第一个Tag是非终结符，第二个Tag是终结符，值是对应产生式
type LLTable map[base.Tag]map[base.Tag]base.Production

// GenerateLLTable 生成LL(1)预测分析表
func GenerateLLTable() LLTable {
	llTable := make(LLTable)

	for left, productions := range base.GetProdMap() {
		llTable[left] = map[base.Tag]base.Production{}
		for _, production := range productions {
			// 终结符则将该式子添加到LLTable中
			if production.Right[0].Type == base.TERM {
				if !base.IsEmptyTag(production.Right[0]) {
					// 如果First非空
					llTable[left][production.Right[0]] = production
				} else {
					// 如果First为空，则在Follow(production.Left)中加入 production.Left -> ε
					for _, tag := range base.GetFollow(production.Left) {
						llTable[left][tag] = production
					}
				}
			}
			// 否则求右部首元素的FIRST集，并将对应位置填上该生成式
			if production.Right[0].Type == base.NONTERM {
				for _, tag := range base.GetFirst(production.Right[0]) {
					llTable[left][tag] = production
				}
			}
		}
	}

	return llTable
}

// PrintLLTable 打印LL分析表
func PrintLLTable(table LLTable) error {
	colNames := []string{" "}
	for _, tag := range base.GetTags() {
		if tag.Type == base.TERM && !base.IsEmptyTag(tag) {
			colNames = append(colNames, tag.Value)
		} else if base.IsEmptyTag(tag) {
			colNames = append(colNames, "$")
		}
	}
	gt, err := gotable.Create(colNames...)
	if err != nil {
		return errors.New("创建表格失败")
	}

	for tag, m := range table {
		row := make(map[string]string)
		row[" "] = tag.Value

		for t, production := range m {
			row[t.Value] = production.ToString()
		}

		err = gt.AddRow(row)
		if err != nil {
			return err
		}
	}

	fmt.Printf("-------------------------------------- 打印LL分析表 -------------------------------------\n")
	gt.PrintTable()
	fmt.Printf("-----------------------------------------------------------------------------------------\n")

	return err
}
