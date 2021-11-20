package ll1

import (
	"errors"
	"fmt"
	"github.com/liushuochen/gotable"
	"gram/base"
	"reflect"
)

// TagStack Tag栈，并在最开始放上$和开始符
var TagStack = []base.Tag{
	{Type: base.TERM, Value: "$"},
}

// PopStack 从栈顶弹出一个tag，并返回该tag
func PopStack() base.Tag {
	if len(TagStack) == 0 {
		panic(errors.New("解析失败，Tag栈中已没有元素！"))
	}
	tag := TagStack[len(TagStack)-1]
	TagStack = TagStack[:len(TagStack)-1]
	return tag
}

// PushStack 向栈顶添加一个tag
func PushStack(tag base.Tag) {
	TagStack = append(TagStack, tag)
}

type Proc struct {
	Step   string
	Stack  string
	Input  string
	Output string
}

var Procedures []Proc

// PrintProcedure 打印LL的分析过程，打印的前提是进行了 LLAnalyze
func PrintProcedure() {
	column := []string{
		"STEP",
		"STACK",
		"INPUT",
		"OUTPUT",
	}
	table, err := gotable.Create(column...)
	if err != nil {
		panic(errors.New("创建LL计算过程表格失败！"))
	}

	for _, procedure := range Procedures {
		row := make(map[string]string)
		row["STEP"] = procedure.Step
		row["STACK"] = procedure.Stack
		row["INPUT"] = procedure.Input
		row["OUTPUT"] = procedure.Output

		err = table.AddRow(row)
		if err != nil {
			panic(errors.New("LL计算过程表格生成行时失败！"))
		}
	}

	fmt.Printf("-------------------------- 打印LL分析过程 -------------------------\n")
	table.PrintTable()
	fmt.Printf("-----------------------------------------------------------------\n")
}

// LLAnalyze LL分析程序
func LLAnalyze(input string, table LLTable) error {
	TagStack = append(TagStack, base.GetProductions()[0].Left)

	// 输入字符串的扫描指针
	index := 0

	step := 0

	// 当且仅当输入队列指针所指字符为 $ 时，才会跳出循环
	for input[index:index+1] != "$" {
		step++
		proc := Proc{
			Step:  fmt.Sprintf("(%v)", step),
			Stack: base.ConvertTagsToStr(TagStack),
			Input: input[index:],
		}

		var tag base.Tag
		// 将input的第一个串匹配为一个tag
		for _, t := range base.GetTags() {
			if t.Value == input[index:index+len(t.Value)] {
				tag = t
				break
			}
		}

		// 如果栈顶是终结符
		stackTop := PopStack()
		if stackTop.Type == base.TERM {
			if reflect.DeepEqual(tag, stackTop) {
				// 若其与tag相同，则同时消去
				index += len(tag.Value)
				continue
			} else {
				// 若其与tag不同，则报错
				return errors.New("计算过程中，出现栈顶元素与输入串元素不匹配的情况")
			}
		}

		// 如果栈顶是非终结符
		if stackTop.Type == base.NONTERM {
			// 根据分析表，得到 栈顶元素-tag 对应的产生式
			production := table[stackTop][tag]

			// 将产生式字符串放到proc中，用于打印计算过程
			proc.Output = production.ToString()

			// 将产生式的右部反序入栈（非空入栈，空则只弹不入）
			for i := len(production.Right) - 1; i >= 0; i-- {
				if !base.IsEmptyTag(production.Right[i]) {
					PushStack(production.Right[i])
				}
			}
		}

		Procedures = append(Procedures, proc)
	}

	// 检查Tag栈中剩余的元素，是否都有 M(E, $) = E -> ε
	for len(TagStack) != 1 {
		if !reflect.DeepEqual(table[PopStack()][base.Tag{
			Type: base.TERM, Value: "$"}].Right, []base.Tag{{
			Type:  base.TERM,
			Value: "ε",
		}}) {
			return errors.New("分析失败：符号栈内还有元素，输入串不符合LL(1)文法")
		}
	}

	fmt.Printf("\n\nLL分析成功！接受输入语句\n\n")

	return nil
}
