package lr1

import (
	"errors"
	"fmt"
	"github.com/liushuochen/gotable"
	"gram/base"
	"reflect"
)

// StateStack 状态栈
var StateStack = []int{0}

// PopStateStack 从栈顶弹出一个state，并返回该state
func PopStateStack() int {
	if len(StateStack) == 0 {
		panic(errors.New("解析失败，State栈中已没有元素！"))
	}
	state := StateStack[len(StateStack)-1]
	StateStack = StateStack[:len(StateStack)-1]
	return state
}

// PushStateStack 向栈顶添加一个tag
func PushStateStack(state int) {
	StateStack = append(StateStack, state)
}

// SymbolStack 符号栈
var SymbolStack = []base.Tag{
	{
		Type:  base.TERM,
		Value: "—",
	},
}

// PopSymbolStack 从栈顶弹出一个tag，并返回该tag
func PopSymbolStack() base.Tag {
	if len(SymbolStack) == 0 {
		panic(errors.New("解析失败，Tag栈中已没有元素！"))
	}
	tag := SymbolStack[len(SymbolStack)-1]
	SymbolStack = SymbolStack[:len(SymbolStack)-1]
	return tag
}

// PushSymbolStack 向栈顶添加一个tag
func PushSymbolStack(tag base.Tag) {
	SymbolStack = append(SymbolStack, tag)
}

type Proc struct {
	Step   string
	Stack  string
	Input  string
	Output string
}

var Procedures []Proc

// LRAnalyze LR1分析程序
func LRAnalyze(input string, table LRTable) error {
	step := 0
	index := 0

	for index < len(input) {
		step++

		var tag base.Tag
		// 将input的第一个串匹配为一个tag
		for _, t := range base.GetTags() {
			if index+len(t.Value) < len(input) && t.Value == input[index:index+len(t.Value)] {
				tag = t
				break
			}
		}

		// 匹配 "$" 的情况
		if input[index:index+1] == "$" {
			tag = base.Tag{
				Type:  base.TERM,
				Value: "$",
			}
		}

		ag := table[StateStack[len(StateStack)-1]][tag]

		// 记录计算过程
		proc := Proc{
			Step: fmt.Sprintf("%v", step),
			Stack: fmt.Sprintf("State: %v ; Symbol: %v",
				ConvertIntSliceToStr(StateStack),
				base.ConvertTagsToStr(SymbolStack)),
			Input:  input[index:],
			Output: ag.ToString(),
		}
		Procedures = append(Procedures, proc)

		// S操作：tag入符号栈，ag.Value入状态栈
		if ag.Type == SHIFT {
			PushSymbolStack(tag)
			PushStateStack(ag.Value)
			index += len(tag.Value)
		}
		if ag.Type == REDUCE {
			// 首先根据对应产生式右部的数量，从状态栈中弹出对应个数的元素
			// 并倒着一个一个匹配，将符号栈中的元素，随着右部弹出
			p := base.GetProductions()[ag.Value]
			for i := len(p.Right) - 1; i >= 0; i-- {
				PopStateStack()
				if !reflect.DeepEqual(p.Right[i], PopSymbolStack()) {
					return errors.New("进行Reduce操作时，出现符号栈中符号不匹配的情况")
				}
			}
			// 然后将左部入栈
			PushSymbolStack(p.Left)
			// 然后根据当前状态栈、左部，得到GOTO内容，将其Value入状态栈
			PushStateStack(table[StateStack[len(StateStack)-1]][p.Left].Value)
		}
		if ag.Type == ACC {
			break
		}
	}

	return nil
}

func PrintProcedure() {
	column := []string{
		"STEP",
		"STACK",
		"INPUT",
		"OUTPUT",
	}
	table, err := gotable.Create(column...)
	if err != nil {
		panic(errors.New("创建LR计算过程表格失败！"))
	}

	for _, procedure := range Procedures {
		row := make(map[string]string)
		row["STEP"] = procedure.Step
		row["STACK"] = procedure.Stack
		row["INPUT"] = procedure.Input
		row["OUTPUT"] = procedure.Output

		err = table.AddRow(row)
		if err != nil {
			panic(errors.New("LR计算过程表格生成行时失败！"))
		}
	}

	fmt.Printf("-------------------------- 打印LR分析过程 -------------------------\n")
	table.PrintTable()
	fmt.Printf("-----------------------------------------------------------------\n")
}

func ConvertIntSliceToStr(nums []int) string {
	str := ""
	for _, num := range nums {
		str += fmt.Sprintf("%v ", num)
	}
	return str
}
